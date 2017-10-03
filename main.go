package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	t := start

	w := flag.Int("w", 2048, "width")
	h := flag.Int("h", 1024, "height")
	r := flag.Int("r", 250, "region count")
	s := flag.Int64("s", start.Unix(), "seed")
	f := flag.String(
		"f", fmt.Sprintf("img/%v.png", start.Unix()),
		"output file path, supported formats .png, .jpg/.jpeg, .gif",
	)
	flag.Parse()

	rand.Seed(*s)

	t = time.Now()
	fmt.Printf("Initializing...         ")
	world := newWorld(*w, *h, *r)
	fmt.Printf("%v\n", time.Now().Sub(t))

	t = time.Now()
	fmt.Printf("Filling gaps...         ")
	world.expandRegions(-1, -1)
	fmt.Printf("%v\n", time.Now().Sub(t))

	t = time.Now()
	fmt.Printf("Saving...               ")
	err := world.image().save(*f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", time.Now().Sub(t))

	fmt.Printf("Total time elapsed:     %v\n", time.Now().Sub(start))
}
