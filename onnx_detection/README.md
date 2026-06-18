# cv_draw

## ライブラリ

```
go get github.com/nfnt/resize
go get github.com/yalue/onnxruntime_go
go get github.com/fogleman/gg
```

## onnx runtimeのインストール

```
cd /tmp/build
wget -O - https://github.com/microsoft/onnxruntime/releases/download/v1.26.0/onnxruntime-linux-x64-gpu-1.26.0.tgz | tar -xz
cp -r onnxrcd untime-linux-x64-gpu-1.26.0/include/* /usr/local/include/
cp -r onnxruntime-linux-x64-gpu-1.26.0/lib/* /usr/local/lib/
cd /tmp
rm -r onnxruntime-linux-x64-gpu-1.26.0
ldconfig
```

## モデルファイル作成

```
pip install ultralytics "onnx>=1.12.0,<2.0.0" "onnxruntime" "onnxslim>=0.1.82" --no-cache-dir  --break-system-packages
yolo export model=yolo26m.pt format=onnx
```

## 参考

https://github.com/yalue/onnxruntime_go_examples/blob/master/image_object_detect/image_object_detect.go
