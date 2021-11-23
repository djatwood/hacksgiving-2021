<script lang="ts">
  import { onMount } from "svelte";
  import { writable } from "svelte/store";

  import Form from "./components/Form.svelte";

  export let generate: ReturnType<typeof writable>;

  let loading = false;
  let canvas: HTMLCanvasElement;
  let canvasCtx: CanvasRenderingContext2D;

  let canvasWidth = 375;
  let canvasHeight: number;
  $: canvasHeight = canvasWidth * (2 / 3);

  let noiseSeed = 0;
  let noiseScale = 30;

  const canvasData = writable(
    new Uint8ClampedArray(canvasWidth * canvasHeight * 4)
  );

  canvasData.subscribe((canvasData) => {
    if (!canvasCtx) return;
    const imageData = new ImageData(canvasData, canvasWidth, canvasHeight);
    canvas.width = canvasWidth;
    canvas.height = canvasHeight;
    canvasCtx.putImageData(imageData, 0, 0);
  });

  const generateNoise = (event?: SubmitEvent) => {
    if (!$generate || loading) return;

    loading = true;
    setTimeout(() => {
      console.time("Saved generated noise");
      $canvasData = new Uint8ClampedArray(
        $generate(canvasWidth, canvasHeight, noiseScale, noiseSeed)
      );
      console.timeEnd("Saved generated noise");
      loading = false;
    }, 0);
  };

  generate.subscribe(() => generateNoise());

  onMount(() => {
    canvas.width = canvasWidth;
    canvas.height = canvasHeight;

    canvasCtx = canvas.getContext("2d");
  });
</script>

<main>
  <canvas
    bind:this={canvas}
    style="--max-width: {canvasWidth / window.devicePixelRatio}px"
  />

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

  canvas {
    width: 1000px;
    background: black;
    border-radius: 4px;
    margin-bottom: 2rem;

    box-shadow: 0 4px 10px #0004;

    image-rendering: pixelated; /* Chrome */
  }
</style>
