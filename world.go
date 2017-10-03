package main

import (
	"image"
	"image/color"
	"math/rand"
)

type world struct {
	width     int
	height    int
	tiles     []tile    // All tiles in a world, arranged row by row
	regions   []region  // All regions in a world
	expanders []*region // Regions that have room left to grow
}

func newWorld(width, height, regionCount int) *world {
	var w world
	w.width = width
	w.height = height

	w.tiles = make([]tile, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			w.tileAt(x, y).pos = point{x, y}
		}
	}

	b := make([]biome, 5)
	b[0] = biome{color.RGBA{0x00, 0x00, 0xFF, 0xFF}}
	b[1] = biome{color.RGBA{0x00, 0x00, 0xFF, 0xFF}}
	b[2] = biome{color.RGBA{0x00, 0x00, 0xFF, 0xFF}}
	b[3] = biome{color.RGBA{0x00, 0xFF, 0x00, 0xFF}}
	b[4] = biome{color.RGBA{0xFF, 0xFF, 0x00, 0xFF}}

	w.regions = make([]region, regionCount)
	w.expanders = make([]*region, regionCount)
	for i := 0; i < regionCount; i++ {
		w.regions[i].initialize(&w, b)
		w.expanders[i] = &w.regions[i]
	}

	w.expandRegions(-1, -1)

	return &w
}

// iterations < 0 to expand random regions until there are no free tiles in the world.
// maxPerRegion <= 0 to use defaults.
func (w *world) expandRegions(iterations, maxPerRegion int) {
	if maxPerRegion <= 0 {
		maxPerRegion = len(w.tiles) / len(w.regions)
	}

	for iterations != 0 && len(w.expanders) > 0 {
		ri := rand.Intn(len(w.expanders))
		r := w.expanders[ri]
		n := rand.Intn(maxPerRegion)
		if iterations >= 0 && n > iterations {
			n = iterations
		}

		for i := 0; i < n; i++ {
			if !r.expand(w) {
				w.expanders[ri] = w.expanders[len(w.expanders)-1]
				w.expanders = w.expanders[:len(w.expanders)-1]
				break
			}
			iterations--
		}
	}
}

func (w *world) image() worldImage {
	return worldImage{w}
}

func (w *world) tileAt(x, y int) *tile {
	return &w.tiles[x+y*w.width]
}

func (w *world) wrappedTileAt(x, y int) *tile {
	if y < 0 || y >= w.height {
		return nil
	}
	if x < 0 {
		x = w.width + x%w.width
	} else if x >= w.width {
		x %= w.width
	}

	return w.tileAt(x, y)
}

type point = image.Point

type region struct {
	border []*tile
	biome  *biome
	color  color.RGBA
	origin point
}

func (r *region) initialize(w *world, b []biome) {
	var t *tile
	for {
		r.origin.X = rand.Intn(w.width)
		r.origin.Y = rand.Intn(w.height)
		t = w.tileAt(r.origin.X, r.origin.Y)
		if t.region == nil {
			break
		}
	}

	r.color = color.RGBA{
		uint8(rand.Intn(0xFF)),
		uint8(rand.Intn(0xFF)),
		uint8(rand.Intn(0xFF)),
		0xFF,
	}

	r.biome = &b[rand.Intn(len(b))]

	r.expandTo(t)
	n := rand.Intn(len(w.tiles) / len(w.regions))
	for i := 0; i < n; i++ {
		if !r.expand(w) {
			break
		}
	}
}

// Returns false if the region can't expand anymore, returns true otherwise.
func (r *region) expand(w *world) bool {
	if len(r.border) == 0 {
		return false
	}

	i := rand.Intn(len(r.border))
	t := r.border[i]

	// This loop is intentionally unrolled for performance.
	// Testing showed significant time savings.
	r.expandTo(w.wrappedTileAt(t.pos.X-1, t.pos.Y))
	r.expandTo(w.wrappedTileAt(t.pos.X+1, t.pos.Y))
	r.expandTo(w.wrappedTileAt(t.pos.X, t.pos.Y-1))
	r.expandTo(w.wrappedTileAt(t.pos.X, t.pos.Y+1))
	r.expandTo(w.wrappedTileAt(t.pos.X-1, t.pos.Y-1))
	r.expandTo(w.wrappedTileAt(t.pos.X+1, t.pos.Y-1))
	r.expandTo(w.wrappedTileAt(t.pos.X-1, t.pos.Y+1))
	r.expandTo(w.wrappedTileAt(t.pos.X+1, t.pos.Y+1))

	r.border[i] = r.border[len(r.border)-1]
	r.border = r.border[:len(r.border)-1]

	return true
}

func (r *region) expandTo(t *tile) {
	if t != nil && t.region == nil {
		r.border = append(r.border, t)
		t.region = r
	}
}

type tile struct {
	region *region
	pos    point
}

type biome struct {
	color color.RGBA
}
