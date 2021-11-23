package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/ojrac/opensimplex-go"
)

func colorAt(gen opensimplex.Noise32, x, y, offset float32) uint8 {
	return uint8(255 * (gen.Eval2(x+offset, y+offset)))
}

func someNoise(width, height, scale, seed int) []interface{} {
	gen := opensimplex.NewNormalized32(int64(seed))

	start := time.Now()
	canvas := make([]interface{}, width*height*4)
	for i := 0; i < len(canvas); i += 4 {
		x := float32((i/4)%(width)) / float32(scale)
		y := float32((i/4)/(width)) / float32(scale)
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

		return someNoise(width, height, scale, seed)
	})
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("generate", genWrapper())
	<-make(chan bool)
}
