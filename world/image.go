package world

import (
	"bufio"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

// Image is an Image of a World created by worldgen.
// It provides methods for saving it to a file
// or encoding it into an io.Writer in any of the following formats:
// .png (default), .jpg, .gif
type Image struct {
	*World
}

// Image is a method for World that creates an Image based on the World.
func (w *World) Image() Image {
	return Image{w}
}

// Save saves the Image to the specified path.
// The file extension can be included as a part of the path,
// but the method defaults to .png.
func (i Image) Save(path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModeDir)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	err = i.Encode(writer, filepath.Ext(path))
	if err != nil {
		return err
	}

	return writer.Flush()
}

// Encode uses the specified format to encode the Image to an io.Writer.
// The following formats are supported:
// .png (default), .jpg, .gif
func (i Image) Encode(writer io.Writer, format string) error {
	var err error

	switch format {
	case ".jpg":
		fallthrough
	case ".jpeg":
		err = jpeg.Encode(writer, i, nil)
	case ".gif":
		err = gif.Encode(writer, i, nil)
	default:
		err = png.Encode(writer, i)
	}

	return err
}

// ColorModel implements the ColorModel method
// required by image.Image interface.
func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements the Bounds method
// required by image.Image interface.
func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.width, i.height)
}

// At implements the At method
// required by image.Image interface.
func (i Image) At(x, y int) color.Color {
	t := i.tileAt(x, y)
	if t.region == nil {
		return color.RGBA{0x00, 0x00, 0x00, 0xFF}
	}
	return i.tileAt(x, y).region.biome.color
}
