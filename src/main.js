import App from "./App.svelte";

const go = new Go();
WebAssembly.instantiateStreaming(fetch("/wasm/index.wasm"), go.importObject).then(
  (result) => {
    go.run(result.instance);
  }
);

const app = new App({
  target: document.body,
  props: {
    name: "world",
  },
});

export default app;
