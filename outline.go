package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"golang.org/x/image/webp"
)

var (
	outlineColor = color.RGBA{255, 255, 255, 255} // White color
	outlineWidth = 5                              // Adjust the width of the outline here
)

func ImageOutline(inputFileName string, outputFileName string) {
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	img, err := webp.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	// // Create a new image to draw the outline and original image on
	// outlinedImg := image.NewRGBA(bounds)

	// Calculate new image dimensions to include the outline
	newWidth := bounds.Dx() + 2*outlineWidth
	newHeight := bounds.Dy() + 2*outlineWidth

	// Create a new image with the increased dimensions
	extendedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Fill the new image with a fully transparent background
	draw.Draw(extendedImg, extendedImg.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Determine the offset to center the original image on the new canvas
	offsetX := (newWidth - bounds.Dx()) / 2
	offsetY := (newHeight - bounds.Dy()) / 2

	// Draw the original image centered on the new canvas
	draw.Draw(extendedImg, img.Bounds().Add(image.Point{X: offsetX, Y: offsetY}), img, bounds.Min, draw.Over)

	// Iterate over the image to find non-transparent pixels and draw an outline around them
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			_, _, _, a := originalPixel.RGBA()
			if a > 0 { // If the pixel is not fully transparent
				// Draw an outline around the pixel
				for dy := -outlineWidth; dy <= outlineWidth; dy++ {
					for dx := -outlineWidth; dx <= outlineWidth; dx++ {
						if x+dx >= bounds.Min.X && x+dx < bounds.Max.X && y+dy >= bounds.Min.Y && y+dy < bounds.Max.Y {
							extendedImg.Set(x+dx, y+dy, outlineColor)
						}
					}
				}
			}
		}
	}

	// Redraw the original image over the outline to ensure it's on top
	draw.Draw(extendedImg, img.Bounds().Add(image.Point{X: offsetX, Y: offsetY}), img, bounds.Min, draw.Over)

	// Save the new image with the outline
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, extendedImg); err != nil {
		panic(err)
	}
}
