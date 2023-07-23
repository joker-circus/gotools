package gotools

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
	"net/http"
	"os"
	"strings"
)

const (
	ImageTypeJpeg    = "image/jpeg"
	ImageTypeGif     = "image/gif"
	ImageTypePng     = "image/png"
	ImageTypeUnknown = "unknown"
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
func ImageType(filename string) string {
	switch ft := getFileType(filename); ft {
	case ImageTypeJpeg, ImageTypeGif, ImageTypePng:
		return ft
	}
	return ImageTypeUnknown
}

// 通过 img 文件数据，探查图片的真正格式，判断 image 类型。
func ImageTypeByDetectContentType(img []byte) string {
	contentType := http.DetectContentType(img)
	switch contentType {
	case ImageTypeJpeg, ImageTypeGif, ImageTypePng:
		return contentType
	}
	return ImageTypeUnknown
}

// 获取 image 对象，img 是文件读取的字节数据。
func DecodeImg(img []byte, imageType ...string) (image.Image, error) {
	if len(imageType) == 0 {
		m, _, err := image.Decode(bytes.NewBuffer(img))
		return m, err
	}

	return DecodeImgByImageType(img, imageType[0])
}

// 获取 Image 内容，通过指定 ImageType，img 是文件读取的字节数据。
func DecodeImgByImageType(img []byte, imageType string) (image.Image, error) {
	f := bytes.NewBuffer(img)
	switch imageType {
	case ImageTypeJpeg:
		return jpeg.Decode(f)
	case ImageTypeGif:
		return gif.Decode(f)
	case ImageTypePng:
		return png.Decode(f)
	}
	return nil, fmt.Errorf("invalid image type")
}

// 获取 Image 内容。
// img 是文件读取的字节数据。
// 因为有人为修改后缀的可能性，所以会先探查一下图片的真正格式。
func DecodeImgByDetectContentType(img []byte) (image.Image, error) {
	return DecodeImg(img, ImageTypeByDetectContentType(img))
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

// 获取 image 文件对应的 Image 对象
func ImageFile(file string) (image.Image, error) {
	ff, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return DecodeImgByDetectContentType(ff)
}

// base64 值转 Image 对象。
func Base64ToImage(base64Data []byte) (image.Image, error) {
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(base64Data)))
	n, err := base64.StdEncoding.Decode(dbuf, base64Data)
	if err != nil {
		return nil, err
	}
	dbuf = dbuf[:n]
	return DecodeImg(dbuf)
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

// 去除 base64 字符前缀，例如：`data:image/png;base64,iVBO...uQmCC` => `iVBO...uQmCC`
func Base64TrimData(urlBase64 string) string {
	i := strings.Index(urlBase64, ",")
	return urlBase64[i+1:]
}

// 给 base64 加上 `data:image/png;base64,` 这种前缀，方便在 html 上直接使用。
func UrlBase64(base64Data, imageType string) string {
	return fmt.Sprintf("data:%s;base64,%s", imageType, base64Data)
}

// 将图片 base64 数据，写入到文件中。
// 文件后缀名最好与图片的类型保持一致。
func WriteBase64(base64Data string, fileName string) error {
	dec, err := base64.StdEncoding.DecodeString(Base64TrimData(base64Data))
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, dec, 0666)
}

// 将 image 对象写入文件中
func WriteImage(img image.Image, filename string) error {
	ft := ImageType(filename)
	if ft == ImageTypeUnknown {
		return fmt.Errorf("filename type is not support")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	//写入文件
	switch ft {
	case ImageTypeJpeg:
		return jpeg.Encode(f, img, nil)
	case ImageTypeGif:
		return gif.Encode(f, img, nil)
	case ImageTypePng:
		return png.Encode(f, img)
	}
	return nil
}
