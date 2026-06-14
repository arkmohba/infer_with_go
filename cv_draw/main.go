package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func main() {
	inputFile := "image.jpg"
	img := gocv.IMRead(inputFile, gocv.IMReadColor)
	if img.Empty() {
		log.Fatalf("Failt to load image: %s", inputFile)
	}
	fmt.Printf("Image Loaded: %s\n", inputFile)

	bbox := image.Rect(100, 50, 300, 250)

	borderColor := color.RGBA{R: 0, G: 255, B: 0, A: 0}
	thickness := 3

	gocv.Rectangle(&img, bbox, borderColor, thickness)

	outputFile := "output.jpg"
	if ok := gocv.IMWrite(outputFile, img); !ok {
		log.Fatalf("Failed to save image: %s", outputFile)
	}
	fmt.Printf("Saved output: %s\n", outputFile)
}
