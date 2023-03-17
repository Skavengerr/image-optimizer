package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

func SaveImage(id string, quality int, imgBytes []byte) error {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("./images/%s-%d.jpg", id, quality))
	if err != nil {
		return err
	}
	defer file.Close()

	opts := jpeg.Options{Quality: quality}
	if err := jpeg.Encode(file, img, &opts); err != nil {
		return err
	}

	return nil
}

func LoadImage(id string) (image.Image, error) {
	file, err := os.Open(fmt.Sprintf("./images/%s.png", id))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ResizeImage(imagePath string, scale float64) ([]byte, error) {
	file, err := os.Open(imagePath)
	fmt.Println("err1", err)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	fmt.Println("err2", err)
	if err != nil {
		return nil, err
	}

	width := int(float64(img.Bounds().Dx()) * scale)
	height := int(float64(img.Bounds().Dy()) * scale)

	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, nil)
	fmt.Println("err3", err)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
