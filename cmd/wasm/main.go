package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/aquilax/go-perlin"
)

func colorAt(gen *perlin.Perlin, x, y, offset float64) uint8 {
	return uint8(255 * ((1 + gen.Noise2D(x+offset, y+offset)) / 2))
}

func someNoise(width, height, scale, seed int) []uint8 {
	gen := perlin.NewPerlin(2, 2, 3, int64(seed))

	start := time.Now()
	canvas := make([]uint8, width*height*4)
	for i := 0; i < len(canvas); i += 4 {
		x := float64((i/4)%(width)) / float64(scale)
		y := float64((i/4)/(width)) / float64(scale)
		canvas[i] = colorAt(gen, x, y, 0.25)
		canvas[i+1] = colorAt(gen, x, y, 0.5)
		canvas[i+2] = colorAt(gen, x, y, 0.75)
		canvas[i+3] = 255
	}

	fmt.Println("Generated noise:", time.Since(start))

	return canvas
}

func genWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := args[0].Int()
		height := args[1].Int()
		scale := args[2].Int()
		seed := args[3].Int()

		data := someNoise(width, height, scale, seed)
		returned := make([]interface{}, len(data))
		for i := 0; i < len(data); i++ {
			returned[i] = data[i]
		}

		return returned
	})
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("generate", genWrapper())
	<-make(chan bool)
}
