package main

import (
	"net/http"
	"syscall/js"
	"time"

	"github.com/djatwood/hacksgiving-2021/packages/generator"
	"github.com/ojrac/opensimplex-go"
)

func streamGeneration() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := args[0].Int()
		height := args[1].Int()
		scale := float32(args[2].Int())
		seed := args[3].Int()
		simple := args[4].Bool()
		slow := args[5].Bool()

		gen := opensimplex.NewNormalized32(int64(seed))

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]

			go func() {
				streamSource := map[string]interface{}{
					"start": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						controller := args[0]

						go func() {
							for cursor := 0; cursor < height; cursor++ {
								row := []byte(generator.GenerateRow(gen, simple, width, height, scale, float32(cursor)))
								array := js.Global().Get("Uint8ClampedArray")
								data := array.New(len(row))
								js.CopyBytesToJS(data, row)
								controller.Call("enqueue", data)

								if slow {
									time.Sleep(1 * time.Nanosecond)
								}
							}

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
