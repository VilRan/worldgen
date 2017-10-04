package main

import "testing"

func TestNewWorld(t *testing.T) {
	w := newWorld(10, 10, 1)
	if w == nil {
		t.Fatal("World should not be nil.")
	}
}

func TestInitializeTiles(t *testing.T) {
	w := worldWithTiles()
	if len(w.tiles) != w.width*w.height {
		t.Fatal("")
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
func worldWithTiles() *world {
	var w world
	w.width = 10
	w.height = 10
	w.initializeTiles()
	return &w
}
