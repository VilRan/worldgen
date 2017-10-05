package main

import (
	"image/png"
	"log"
	"net/http"

	"github.com/VilRan/worldgen/world"
)

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	world := world.NewWorld(2048, 1024, 250)
	png.Encode(w, world.Image())
}
