package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"os"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	ort "github.com/yalue/onnxruntime_go"
)

type ModelSession struct {
	Session *ort.AdvancedSession
	Input   *ort.Tensor[float32]
	Output  *ort.Tensor[float32]
}

func initSession(modelPath string, inputShape ort.Shape, outputShape ort.Shape) (*ModelSession, error) {
	err := ort.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("Error initializing ORT environment: %w", err)
	}

	inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	if err != nil {
		return nil, fmt.Errorf("Error creating input tensor: %w", err)
	}
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		inputTensor.Destroy()
		return nil, fmt.Errorf("Error creating output tensor: %w", err)
	}
	options, err := ort.NewSessionOptions()
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		return nil, fmt.Errorf("Error creating ORT session options: %w", err)
	}
	defer options.Destroy()

	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"images"}, []string{"output0"},
		[]ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensor},
		options)
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		return nil, fmt.Errorf("Error creating ORT session: %w", err)
	}

	return &ModelSession{
		Session: session,
		Input:   inputTensor,
		Output:  outputTensor,
	}, nil
}

func (m *ModelSession) Destroy() {
	m.Session.Destroy()
	m.Input.Destroy()
	m.Output.Destroy()
}

func loadImageFile(filePath string) (image.Image, error) {
	f, e := os.Open(filePath)
	if e != nil {
		return nil, fmt.Errorf("Error opening %s: %w", filePath, e)
	}
	defer f.Close()
	pic, _, e := image.Decode(f)
	if e != nil {
		return nil, fmt.Errorf("Error decoding %s: %w", filePath, e)
	}
	return pic, nil
}

func prepareInput(pic image.Image, dst *ort.Tensor[float32], modelWidth uint, modelHeight uint) error {
	data := dst.GetData()
	channelSize := modelWidth * modelHeight
	if len(data) < (int(channelSize) * 3) {
		return fmt.Errorf("Destination tensor only holds %d floats, needs "+
			"%d (make sure it's the right shape!)", len(data), channelSize*3)
	}
	redChannel := data[0:channelSize]
	greenChannel := data[channelSize : channelSize*2]
	blueChannel := data[channelSize*2 : channelSize*3]

	pic = resize.Resize(modelWidth, modelHeight, pic, resize.Bilinear)

	i := 0
	for y := range modelHeight {
		for x := range modelWidth {
			r, g, b, _ := pic.At(int(x), int(y)).RGBA()
			redChannel[i] = float32(r>>8) / 255.0
			greenChannel[i] = float32(g>>8) / 255.0
			blueChannel[i] = float32(b>>8) / 255.0
			i++
		}
	}

	return nil
}

type BoundingBox struct {
	classID        int
	confidence     float32
	x1, y1, x2, y2 float32
}

func processOutput(output []float32, scaleX,
	scaleY float32) []BoundingBox {
	boundingBoxes := make([]BoundingBox, 0, 300)

	var classID int
	var probability float32

	for idx := range 300 {
		x1 := output[idx*6+0] * float32(scaleX)
		y1 := output[idx*6+1] * float32(scaleY)
		x2 := output[idx*6+2] * float32(scaleX)
		y2 := output[idx*6+3] * float32(scaleY)
		probability = output[idx*6+4]
		classID = int(output[idx*6+5])

		// If the probability is less than 0.5, continue to the next index
		if probability < 0.5 {
			continue
		}

		// Append the bounding box to the result
		boundingBoxes = append(boundingBoxes, BoundingBox{
			classID:    classID,
			confidence: probability,
			x1:         x1,
			y1:         y1,
			x2:         x2,
			y2:         y2,
		})
	}
	return boundingBoxes
}

func DrawBoundingBox(g *gg.Context, bbox BoundingBox) {
	color := CocoColors[bbox.classID%len(CocoColors)]
	g.SetRGB(color[0], color[1], color[2])
	g.SetLineWidth(2)
	g.DrawRectangle(float64(bbox.x1), float64(bbox.y1), float64(bbox.x2-bbox.x1), float64(bbox.y2-bbox.y1))
	g.Stroke()
}

func main() {
	ort.SetSharedLibraryPath("/usr/local/lib/libonnxruntime.so")
	// Create ONNX Session
	var modelWidth int64 = 640
	var modelHeight int64 = 640
	inputShape := ort.NewShape(1, 3, modelHeight, modelWidth)
	outputShape := ort.NewShape(1, 300, 6)
	session, err := initSession("yolo26m.onnx", inputShape, outputShape)
	if err != nil {
		panic(err)
	}
	defer session.Destroy()

	img, err := loadImageFile("image.jpg")
	if err != nil {
		panic(err)
	}
	originalWidth := img.Bounds().Canon().Dx()
	originalHeight := img.Bounds().Canon().Dy()

	err = prepareInput(img, session.Input, uint(modelWidth), uint(modelHeight))
	if err != nil {
		panic(err)
	}

	err = session.Session.Run()
	if err != nil {
		panic(err)
	}

	boxes := processOutput(session.Output.GetData(), float32(originalWidth)/float32(modelWidth),
		float32(originalHeight)/float32(modelHeight))

	g := gg.NewContextForImage(img)
	for idx := range boxes {
		fmt.Printf("x1:%f y1:%f x2:%f y2:%f\n", boxes[idx].x1, boxes[idx].y1, boxes[idx].x2, boxes[idx].y2)
		DrawBoundingBox(g, boxes[idx])
	}

	outFile, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, g.Image(), &jpeg.Options{Quality: 95})
	if err != nil {
		panic(err)
	}

}
