const go = new Go();
WebAssembly.instantiateStreaming(
  fetch("/wasm/index.wasm"),
  go.importObject
).then((result) => {
  go.run(result.instance);
});
