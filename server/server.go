package main

import (
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/VilRan/worldgen/world"
)

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(wr http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	q.Add("w", "2048")
	q.Add("h", "1024")
	q.Add("r", "250")
	q.Add("s", fmt.Sprint(time.Now().Unix()))

	w, err := strconv.ParseInt(q.Get("w"), 10, 32)
	if err != nil {
		wr.Write([]byte(err.Error()))
	}
	h, err := strconv.ParseInt(q.Get("h"), 10, 32)
	if err != nil {
		wr.Write([]byte(err.Error()))
	}
	r, err := strconv.ParseInt(q.Get("r"), 10, 32)
	if err != nil {
		wr.Write([]byte(err.Error()))
	}
	s, err := strconv.ParseInt(q.Get("s"), 10, 64)
	if err != nil {
		wr.Write([]byte(err.Error()))
	}

	rand.Seed(int64(s))

	world := world.NewWorld(int(w), int(h), int(r))
	png.Encode(wr, world.Image())
}
