package types

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime"
	"os"
	"strings"
)

// 根据文件后缀，获取文件类型
func getFileType(filename string) string {
	s := strings.Split(filename, ".")
	if len(s) < 2 {
		return ""
	}

	ext := "." + s[len(s)-1]
	return mime.TypeByExtension(ext)
}

// 根据 image 文件名，判断 image 类型
func imageType(filename string) string {
	switch ft := getFileType(filename); ft {
	case "image/jpeg", "image/gif", "image/png":
		return ft
	}
	return "unknown"
}

// 获取 Image 内容
func DecodeImg(img []byte, filename ...string) (image.Image, error) {
	f := bytes.NewBuffer(img)
	if len(filename) == 0 {
		m, _, err := image.Decode(f)
		return m, err
	}

	switch filename[0] {
	case "image/jpeg":
		return jpeg.Decode(f)
	case "image/gif":
		return gif.Decode(f)
	case "image/png":
		return png.Decode(f)
	}

	return nil, fmt.Errorf("invalid file type")
}

// 裁剪图片
//func subImage(img []byte) image.Image {
//	f := bytes.NewBuffer(img)
//	m, _, _ := image.Decode(f) // 图片文件解码
//	rgbImg := m.(*image.YCbCr)
//	subImg := rgbImg.SubImage(image.Rect(0, 0, 200, 200)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
//	return subImg
//}

// 获取 image 文件对应的 base64 值
func Base64ImageFile(file string) ([]byte, error) {
	ff, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dist := make([]byte, base64.StdEncoding.EncodedLen(len(ff)))
	// 文件转base64
	base64.StdEncoding.Encode(dist, ff)
	return dist, nil
}

// 获取 image 对象对应的 base64 值
func Base64Image(img image.Image) ([]byte, error) {
	//开辟一个新的空buff
	emptyBuff := bytes.NewBuffer(nil)
	//img写入到buff
	err := jpeg.Encode(emptyBuff, img, nil)
	if err != nil {
		return nil, err
	}
	//开辟存储空间
	dist := make([]byte, base64.StdEncoding.EncodedLen(emptyBuff.Len()))
	//buff转成base64
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	return dist, nil
}

// 将 image 对象写入文件中
func WriteImage(img image.Image, filename string) error {
	ft := imageType(filename)
	if ft == "unknown" {
		return fmt.Errorf("filename type is not support")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	//写入文件
	switch ft {
	case "image/jpeg":
		return jpeg.Encode(f, img, nil)
	case "image/gif":
		return gif.Encode(f, img, nil)
	case "image/png":
		return png.Encode(f, img)
	}
	return nil
}
