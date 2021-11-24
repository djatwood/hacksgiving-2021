package main

import (
	"net/http"
	"syscall/js"

	"github.com/djatwood/hacksgiving-2021/packages/generator"
	"github.com/ojrac/opensimplex-go"
)

func streamGeneration() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := args[0].Int()
		height := args[1].Int()
		scale := float32(args[2].Int())
		seed := args[3].Int()
		chunkSize := args[4].Int()
		// simple := args[5].Bool()
		simple := false

		gen := opensimplex.NewNormalized32(int64(seed))

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			if chunkSize < 1 {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("chunkSize must be greater than or equal to 1")
				reject.Invoke(errorObject)
				return nil
			}

			go func() {
				streamSource := map[string]interface{}{
					"start": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						controller := args[0]

						go func() {
							buffer := make([]byte, chunkSize*width*4)
							Uint8ClampedArray := js.Global().Get("Uint8ClampedArray")

							for cursor := 0; cursor < height; cursor++ {
								row := []byte(generator.GenerateRow(gen, simple, width, height, scale, float32(cursor)))
								copy(buffer[(cursor%chunkSize)*4*width:], row)

								if (cursor+1)%chunkSize == 0 {
									data := Uint8ClampedArray.New(len(buffer))
									js.CopyBytesToJS(data, buffer)
									controller.Call("enqueue", data)
								}
							}

							data := Uint8ClampedArray.New(len(buffer))
							js.CopyBytesToJS(data, buffer)
							controller.Call("enqueue", data)

							controller.Call("close")
						}()

						return nil
					}),
				}

				readableStreamConstructor := js.Global().Get("ReadableStream")
				readableStream := readableStreamConstructor.New(streamSource)

				responseInitObj := map[string]interface{}{
					"status":     http.StatusOK,
					"statusText": http.StatusText(http.StatusOK),
				}

				responseConstructor := js.Global().Get("Response")
				response := responseConstructor.New(readableStream, responseInitObj)

				resolve.Invoke(response)
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func main() {
	js.Global().Set("generate", streamGeneration())
	<-make(chan bool)
}
