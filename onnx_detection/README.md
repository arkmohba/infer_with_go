# cv_draw

## ライブラリ

```
go get "github.com/nfnt/resize"
go get github.com/yalue/onnxruntime_go
```

## モデルファイル作成

```
pip install ultralytics "onnx>=1.12.0,<2.0.0" "onnxruntime" "onnxslim>=0.1.82" --no-cache-dir  --break-system-packages
pip install
yolo export model=yolo26m.pt format=onnx
```

## 参考

https://github.com/yalue/onnxruntime_go_examples/blob/master/image_object_detect/image_object_detect.go
