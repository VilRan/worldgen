package main

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

type worldImage struct {
	*world
}

func (w *world) image() worldImage {
	return worldImage{w}
}

func (w worldImage) save(path string) error {
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
		err = jpeg.Encode(writer, w, nil)
	case ".gif":
		err = gif.Encode(writer, w, nil)
	default:
		err = png.Encode(writer, w)
	}
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

func (w worldImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (w worldImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.width, w.height)
}

func (w worldImage) At(x, y int) color.Color {
	t := w.tileAt(x, y)
	if t.region == nil {
		return color.RGBA{0x00, 0x00, 0x00, 0xFF}
	}
	return w.tileAt(x, y).region.biome.color
}
