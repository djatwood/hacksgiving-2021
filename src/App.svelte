<script lang="ts">
  import { onMount } from "svelte";
  import type { writable } from "svelte/store";

  import Form from "./components/Form.svelte";

  export let generate: ReturnType<typeof writable>;

  let loading = false;
  let canvas: HTMLCanvasElement;
  let canvasCtx: CanvasRenderingContext2D;

  let canvasWidth = 1024 * window.devicePixelRatio;
  let canvasHeight: number;
  $: canvasHeight = Math.round(canvasWidth * (2 / 3));

  let noiseSeed = 0;
  let noiseScale = 128 * window.devicePixelRatio;
  let slowMode = false;

  const generateNoise = () => {
    if (!$generate || loading) return;

    loading = true;
    setTimeout(async () => {
      console.time("Saved generated noise");
      const response = await $generate(
        canvasWidth,
        canvasHeight,
        noiseScale,
        noiseSeed,
        slowMode
      );

      if (canvas.width != canvasWidth) canvas.width = canvasWidth;
      if (canvas.height != canvasHeight) canvas.height = canvasHeight;

      const body: ReadableStream<Uint8ClampedArray> = await response.body;
      const reader = body.getReader();
      let row = 0;
      while (true) {
        const { done, value } = await reader.read();
        if (value) {
          const rowData = new Uint8ClampedArray(value);
          const imageData = new ImageData(rowData, canvasWidth, 1);
          canvasCtx.putImageData(imageData, 0, row);
          row++;
        }

        if (done) {
          break;
        }
      }

      console.timeEnd("Saved generated noise");

      loading = false;
    }, 0);
  };

  generate.subscribe(() => generateNoise());

  onMount(() => {
    canvas.width = canvasWidth;
    canvas.height = canvasHeight;

    canvasCtx = canvas.getContext("2d", { alpha: false });
    canvasCtx.imageSmoothingEnabled = true;
  });
</script>

<main>
  <canvas bind:this={canvas} />

  <Form
    bind:canvasWidth
    bind:noiseSeed
    bind:noiseScale
    bind:slowMode
    {loading}
    onSubmit={generateNoise}
  />
</main>

<style>
  main {
    display: flex;
    flex-flow: column nowrap;
    justify-content: center;
    align-items: center;
    width: 100vw;
    height: 100vh;

    padding: 1rem;
    box-sizing: border-box;

    background: #212121;
  }

  canvas {
    max-width: 1000px;
    width: 100%;
    margin-bottom: 2rem;

    background: #fff;

    box-shadow: 0 4px 10px #0004;
  }
</style>
