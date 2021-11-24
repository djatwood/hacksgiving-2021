package main

import (
	"fmt"
	"net/http"
	"syscall/js"
	"time"

	"github.com/ojrac/opensimplex-go"
)

func interpolateNoiseTo255(gen opensimplex.Noise32, x, y, offset float32) uint8 {
	return uint8(255 * (gen.Eval2(x+offset, y+offset)))
}

func getColor(value float32) (channels [4]uint8) {
	channels[3] = 255

	// water
	if value < 0.3 {
		channels[0] = 127
		channels[1] = 205
		channels[2] = 255
		return
	}

	// beach
	if value < 0.35 {
		channels[0] = 210
		channels[1] = 255
		channels[2] = 228
		return
	}

	// grass
	if value < 0.7 {
		channels[0] = 2
		channels[1] = 161
		channels[2] = 78
		return
	}

	// forest
	if value < 0.9 {
		channels[0] = 1
		channels[1] = 85
		channels[2] = 41
		return
	}

	// deep forest
	return [4]uint8{
		60,
		42,
		42,
		255,
	}
}

func noiseRow(gen opensimplex.Noise32, width, height, scale, cursor int) []uint8 {
	row := make([]uint8, width*4)
	for col := 0; col < width; col++ {
		x := float32(col%width) / float32(scale)
		y := float32(cursor) / float32(scale)
		value := gen.Eval2(x+0.5, y+0.5)
		color := getColor(value)
		i := col * 4
		row[i] = color[0]
		row[i+1] = color[1]
		row[i+2] = color[2]
		row[i+3] = color[3]
	}
	return row
}

func streamGeneration() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := args[0].Int()
		height := args[1].Int()
		scale := args[2].Int()
		seed := args[3].Int()
		slow := args[4].Bool()

		gen := opensimplex.NewNormalized32(int64(seed))

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			// reject := args[1]

			go func() {
				// See: https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream/ReadableStream
				underlyingSource := map[string]interface{}{
					"start": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						controller := args[0]

						// Process the stream in yet another background goroutine,
						// because we can't block on a goroutine invoked by JS in Wasm
						// that is dealing with HTTP requests
						go func() {
							start := time.Now()
							for cursor := 0; cursor < height; cursor++ {
								buf := []byte(noiseRow(gen, width, height, scale, cursor))
								arrayConstructor := js.Global().Get("Uint8ClampedArray")
								dataJS := arrayConstructor.New(len(buf))
								js.CopyBytesToJS(dataJS, buf[:])
								controller.Call("enqueue", dataJS)

								if slow {
									time.Sleep(1 * time.Nanosecond)
								}
							}
							fmt.Println("Finished streaming", time.Since(start))

							controller.Call("close")
						}()

						return nil
					}),
				}

				// Create a ReadableStream object from the underlyingSource object
				readableStreamConstructor := js.Global().Get("ReadableStream")
				readableStream := readableStreamConstructor.New(underlyingSource)

				// Create the init argument for the Response constructor
				// This allows us to pass a custom status code (and optionally headers and more)
				// See: https://developer.mozilla.org/en-US/docs/Web/API/Response/Response
				responseInitObj := map[string]interface{}{
					"status":     http.StatusOK,
					"statusText": http.StatusText(http.StatusOK),
				}

				// Create a Response object with the stream inside
				responseConstructor := js.Global().Get("Response")
				response := responseConstructor.New(readableStream, responseInitObj)

				// Resolve the Promise
				resolve.Invoke(response)
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("generate", streamGeneration())
	<-make(chan bool)
}
