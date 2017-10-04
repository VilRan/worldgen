package world

import (
	"bufio"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// Image ...
type Image struct {
	*World
}

// Image ...
func (w *World) Image() Image {
	return Image{w}
}

// Save ...
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

	switch filepath.Ext(path) {
	case ".jpg":
		fallthrough
	case ".jpeg":
		err = jpeg.Encode(writer, i, nil)
	case ".gif":
		err = gif.Encode(writer, i, nil)
	default:
		err = png.Encode(writer, i)
	}
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

// ColorModel ...
func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds ...
func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.width, i.height)
}

// At ...
func (i Image) At(x, y int) color.Color {
	t := i.tileAt(x, y)
	if t.region == nil {
		return color.RGBA{0x00, 0x00, 0x00, 0xFF}
	}
	return i.tileAt(x, y).region.biome.color
}