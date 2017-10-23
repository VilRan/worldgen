package world

import (
	"image"
	"image/color"
	"math/rand"
)

// World ...
type World struct {
	width     int       // Width (x-axis) of the world in tiles (=pixels, for now)
	height    int       // Height (y-axis) of the world in tiles (=pixels, for now)
	tiles     []tile    // All tiles in the world, arranged row by row
	regions   []*region // All regions in the world
	expanders []*region // Regions that have room left to grow
}

// NewWorld ...
func NewWorld(width, height, regionCount int) *World {
	var w World
	w.width = width
	w.height = height

	w.initializeTiles()
	w.initializeRegions(regionCount, makeDefaultBiomes())
	w.expandRegions(-1, -1)

	return &w
}

func (w *World) initializeTiles() {
	w.tiles = make([]tile, w.width*w.height)
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			w.tileAt(x, y).pos = point{x, y}
		}
	}
}

func (w *World) initializeRegions(count int, b []*biome) {
	w.regions = make([]*region, 0, count)
	w.expanders = make([]*region, 0, count)
	for i := 0; i < count; i++ {
		r := newRegion(w, b)
		if r == nil {
			continue
		}
		w.regions = append(w.regions, r)
		w.expanders = append(w.expanders, r)
	}
}

// iterations < 0 to expand random regions until there are no free tiles in the world.
// maxPerRegion <= 0 to use defaults.
func (w *World) expandRegions(iterations, maxPerRegion int) {
	if maxPerRegion <= 0 {
		maxPerRegion = w.defaultMaxExpansionsPerRegion()
	}

	for iterations != 0 && len(w.expanders) > 0 {
		ri := rand.Intn(len(w.expanders))
		r := w.expanders[ri]
		n := rand.Intn(maxPerRegion)
		if iterations >= 0 && n > iterations {
			n = iterations
		}

		for i := 0; i < n; i++ {
			if !r.expandRandom(w) {
				w.removeExpanderAt(ri)
				break
			}
			iterations--
		}
	}
}

func (w *World) defaultMaxExpansionsPerRegion() int {
	return len(w.tiles) / (1 + cap(w.regions))
}

func (w *World) removeExpanderAt(index int) {
	w.expanders[index] = w.expanders[len(w.expanders)-1]
	w.expanders = w.expanders[:len(w.expanders)-1]
}

func (w *World) tileAt(x, y int) *tile {
	return &w.tiles[x+y*w.width]
}

func (w *World) wrappedTileAt(x, y int) *tile {
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

// Searches randomly for a tile not occupied by a region.
// If maxTries < 0, the function will not give up
// at least until it overflows and wraps back to 0.
func (w *World) findFreeTile(maxTries int) *tile {
	for maxTries != 0 {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		t := w.tileAt(x, y)
		if t.region == nil {
			return t
		}
		maxTries--
	}
	return nil
}

type point = image.Point

type region struct {
	border []*tile
	biome  *biome
	color  color.RGBA
	origin point
}

func newRegion(w *World, b []*biome) *region {
	t := w.findFreeTile(-1)
	if t == nil {
		return nil
	}

	var r region
	r.origin = t.pos
	r.color = randomColor()
	r.biome = b[rand.Intn(len(b))]

	r.expandTo(t)
	n := w.defaultMaxExpansionsPerRegion()
	for i := 0; i < n; i++ {
		if !r.expandRandom(w) {
			break
		}
	}

	return &r
}

func randomColor() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(0xFF)),
		uint8(rand.Intn(0xFF)),
		uint8(rand.Intn(0xFF)),
		0xFF,
	}
}

// Returns false if the region can't expand anymore, returns true otherwise.
func (r *region) expandRandom(w *World) bool {
	if len(r.border) == 0 {
		return false
	}
	i := rand.Intn(len(r.border))
	r.expand(i, w)
	return true
}

func (r *region) expand(borderTileIndex int, w *World) {
	t := r.border[borderTileIndex]
	r.removeBorderAt(borderTileIndex)

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
}

func (r *region) removeBorderAt(index int) {
	r.border[index] = r.border[len(r.border)-1]
	r.border = r.border[:len(r.border)-1]
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

func makeDefaultBiomes() []*biome {
	b := make([]*biome, 5)
	b[0] = &biome{color.RGBA{0x00, 0x00, 0xFF, 0xFF}}
	b[1] = &biome{color.RGBA{0x00, 0x00, 0xFF, 0xFF}}
	b[2] = &biome{color.RGBA{0x00, 0x00, 0xFF, 0xFF}}
	b[3] = &biome{color.RGBA{0x00, 0xFF, 0x00, 0xFF}}
	b[4] = &biome{color.RGBA{0xFF, 0xFF, 0x00, 0xFF}}
	return b
}
