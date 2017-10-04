package world

import "testing"

func TestNewWorld(t *testing.T) {
	w := NewWorld(10, 10, 1)
	if w == nil {
		t.Fatal("World should not be nil.")
	}
}

func BenchmarkNewWorld(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewWorld(2024, 1024, 250)
	}
}

func TestInitializeTiles(t *testing.T) {
	w := worldWithTiles()
	size := w.width * w.height
	length := len(w.tiles)
	if length != size {
		t.Fatalf("len(w.tiles) was %v, expected %v", length, size)
	}
}

func TestTileAt(t *testing.T) {
	w := worldWithTiles()

	x := w.width - 1
	y := 0
	tile := w.tileAt(x, y)
	pos := point{x, y}
	if tile == nil {
		t.Fatal("Tile should not be nil.")
	}
	if tile.pos != pos {
		t.Fatalf("tile.pos was %v, expected %v", tile.pos, pos)
	}
}

func TestWrappedTileAt(t *testing.T) {
	w := worldWithTiles()

	x := 0
	y := -1
	tile := w.wrappedTileAt(x, y)
	if tile != nil {
		t.Fatal("Tile should be nil.")
	}

	x = -1
	y = 0
	tile = w.wrappedTileAt(x, y)
	pos := point{x + w.width, y}
	if tile.pos != pos {
		t.Fatalf("tile.pos was %v, expected %v", tile.pos, pos)
	}
}

/*
func TestExpand(t *testing.T) {
	w := worldWithTiles()


}
*/

func worldWithTiles() *World {
	var w World
	w.width = 10
	w.height = 10
	w.initializeTiles()
	return &w
}
