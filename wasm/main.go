package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "canvas")
	width := doc.Get("body").Get("clientWidth").Float()
	height := doc.Get("body").Get("clientHeight").Float()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	ctx := canvasEl.Call("getContext", "2d")
	done := make(chan struct{}, 0)

	ctx.Set("fillStyle", "#000")
	ctx.Call("beginPath")
	ctx.Call("fillRect", 20, 20, 150, 100)
	ctx.Call("stroke")

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("render!")

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)

	fmt.Println("Hello, WebAssembly!")

	<-done
}
