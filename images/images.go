// Package images provides image processing utilities, including compression.
//
// This package currently supports JPEG compression for PNG and JPEG images,
// converting all output to JPEG format for consistent compression.
//
// Example usage:
//
//	// Read an image file
//	imageData, _ := os.ReadFile("photo.png")
//
//	// Compress with 80% quality
//	compressed, newFilename, err := images.CompressImage(imageData, 80, "photo.png")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Save compressed image (filename will be photo.jpeg)
//	os.WriteFile(newFilename, compressed, 0644)
//
//	// Use with HTTP upload
//	func UploadHandler(w http.ResponseWriter, r *http.Request) {
//		file, header, _ := r.FormFile("image")
//		data, _ := io.ReadAll(file)
//
//		compressed, filename, err := images.CompressImage(data, 75, header.Filename)
//		if err != nil {
//			http.Error(w, err.Error(), 400)
//			return
//		}
//		// Save or process compressed image
//	}
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

// Re is a wrapper for regexp.Regexp (unused, consider removing)
type Re struct {
	R *regexp.Regexp
}

// supported defines the file extensions accepted for compression
var supported = regexp.MustCompile(`\.(png|jpg|jpeg)$`)

// CompressImage compresses an image to JPEG format with the specified quality.
//
// Parameters:
//   - data: raw image bytes (PNG or JPEG)
//   - quality: JPEG quality from 1 (worst) to 100 (best), recommended 60-85
//   - filename: original filename, used to validate format and generate output name
//
// Returns:
//   - compressed image bytes
//   - new filename with .jpeg extension
//   - error if compression fails or format is unsupported
//
// The function converts PNG images to JPEG (with white background for transparency)
// and re-encodes JPEG images at the specified quality level.
//
// Example:
//
//	// Compress a PNG with 80% quality
//	compressed, filename, err := images.CompressImage(pngData, 80, "avatar.png")
//	// filename = "avatar.jpeg"
//
//	// Compress uploaded file
//	data := readUploadedFile(r)
//	compressed, _, err := images.CompressImage(data, 75, uploadedFilename)
//	if err != nil {
//		// Handle unsupported format or compression error
//	}
//
//	// Note: Returns error if compressed size is larger than original
//	// (can happen with already-optimized images at high quality settings)
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
