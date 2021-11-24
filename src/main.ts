import App from "./App.svelte";

import { writable } from 'svelte/store';

const generate = writable(null);

const go = new Go();
WebAssembly.instantiateStreaming(fetch("wasm/index.wasm"), go.importObject).then((result) => {
  go.run(result.instance);
  generate.set(window.generate)
});

const app = new App({
  target: document.body,
  props: {
    generate: generate
  },
});

export default app;
