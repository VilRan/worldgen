package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/VilRan/worldgen/world"
)

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handle(writer http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	q.Add("w", "2048")
	q.Add("h", "1024")
	q.Add("r", "250")
	q.Add("s", fmt.Sprint(time.Now().Unix()))
	q.Add("f", "png")

	w, err := strconv.ParseInt(q.Get("w"), 10, 32)
	h, err := strconv.ParseInt(q.Get("h"), 10, 32)
	r, err := strconv.ParseInt(q.Get("r"), 10, 32)
	s, err := strconv.ParseInt(q.Get("s"), 10, 64)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	rand.Seed(int64(s))
	world := world.NewWorld(int(w), int(h), int(r))
	err = world.Image().Encode(writer, q.Get("f"))
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
}
