<script lang="ts">
  import { onMount } from "svelte";

  import Form from "./components/Form.svelte";

  export let generate: ReturnType<typeof writable>;

  let loading = false;
  let canvas: HTMLCanvasElement;
  let canvasCtx: CanvasRenderingContext2D;

  let canvasWidth = 1024;
  let canvasHeight: number;
  $: canvasHeight = canvasWidth * (2 / 3);

  let noiseSeed = 0;
  let noiseScale = 128;

  const generateNoise = (event?: SubmitEvent) => {
    if (!$generate || loading) return;

    loading = true;
    setTimeout(async () => {
      console.time("Saved generated noise");
      const data = await $generate(
        canvasWidth,
        canvasHeight,
        noiseScale,
        noiseSeed,
        false
      );

      const body: ReadableStream<Uint8ClampedArray> = await data.body;
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
  });
</script>

<main>
  <div class="container">
    <canvas
      class:loading
      bind:this={canvas}
      style="--max-width: {canvasWidth / window.devicePixelRatio}px"
    />
  </div>
  <Form
    bind:canvasWidth
    bind:noiseSeed
    bind:noiseScale
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

    background: #212121;
  }

  .container {
    width: 1000px;

    margin-bottom: 2rem;

    background: #fff;
    box-shadow: 0 4px 10px #0004;
  }

  canvas {
    display: block;
    width: 100%;

    image-rendering: pixelated;

    transition: opacity 200ms ease-in-out;
  }

  /* canvas.loading {
    opacity: 0.5;
  } */
</style>
