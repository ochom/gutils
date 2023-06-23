package images

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"

	_ "image/png" // register png decoder
	"regexp"
)

// Re is a regexp
type Re struct {
	R *regexp.Regexp
}

var supported = regexp.MustCompile(`\.(png|jpg|jpeg)$`)

// CompressImage compresses an image to a given size
func CompressImage(data []byte, quality int, filename string) ([]byte, string, error) {
	if !supported.MatchString(filename) {
		return nil, filename, fmt.Errorf("unsupported file type")
	}

	filename = supported.ReplaceAllString(filename, ".jpeg")

	imgSrc, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, filename, fmt.Errorf("failed to decode image: %v", err)
	}

	newImg := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)

	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, filename, fmt.Errorf("error encoding image: %s", err.Error())
	}

	if buf.Len() > len(data) {
		return nil, filename, fmt.Errorf("image is too big")
	}

	return buf.Bytes(), filename, nil
}
