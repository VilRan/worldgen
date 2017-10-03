package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()

	w := flag.Int("w", 2048, "width")
	h := flag.Int("h", 1024, "height")
	r := flag.Int("r", 250, "region count")
	n := flag.Int("n", 1, "number of worlds to generate")
	s := flag.Int64("s", start.Unix(), "seed")
	p := flag.String("p", fmt.Sprintf("img/%v", start.Unix()), "output file path")
	f := flag.String("f", "png", "output format: png, jpg/jpeg, gif")
	flag.Parse()

	rand.Seed(*s)

	worlds := make(chan *world)
	done := make(chan bool)

	go func() {
		for i := 0; i < *n; i++ {
			world := <-worlds
			t := time.Now()
			world.image().save(fmt.Sprintf("%v-%v.%v", *p, i, *f))
			fmt.Printf("Saved world, took 		%v\n", time.Now().Sub(t))
		}
		done <- true
	}()
	for i := 0; i < *n; i++ {
		//go func(i int) {
		t := time.Now()
		world := newWorld(*w, *h, *r)
		world.expandRegions(-1, -1)
		worlds <- world
		fmt.Printf("Created world %v, took 		%v\n", i, time.Now().Sub(t))
		//}(i)
	}

	<-done
	fmt.Printf("Total time elapsed:		%v\n", time.Now().Sub(start))
}
