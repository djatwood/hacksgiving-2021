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

  let noiseSeed = 73;
  let noiseScale = 128 * window.devicePixelRatio;
  let simpleMode = false;
  let slowMode = true;

  const generateNoise = () => {
    if (!$generate || loading) return;

    loading = true;
    setTimeout(async () => {
      console.time("Generated terrain");

      const response = await $generate(
        canvasWidth,
        canvasHeight,
        noiseScale,
        noiseSeed,
        simpleMode
      );

      if (canvas.width != canvasWidth) canvas.width = canvasWidth;
      if (canvas.height != canvasHeight) canvas.height = canvasHeight;

      const body: ReadableStream<Uint8ClampedArray> = await response.body;
      const reader = body.getReader();

      const drawRow = (row: number, data?: Uint8ClampedArray) => {
        if (!data) {
          return;
        }
        const imageData = new ImageData(data, canvasWidth, 1);
        canvasCtx.putImageData(imageData, 0, row);
      };

      let row = 0;
      const slowDraw = async () => {
        const buffer = await Promise.all(
          Array(Math.floor(Math.min(canvasHeight - row, canvasHeight * 0.01)))
            .fill(0)
            .map(async () => await reader.read())
        );

        const data = new Uint8ClampedArray(buffer.length * 4 * canvasWidth);
        buffer.forEach(({ value }, index) => {
          data.set(value, index * 4 * canvasWidth);
        });

        const imageData = new ImageData(data, canvasWidth, buffer.length);
        canvasCtx.putImageData(imageData, 0, row);

        row += buffer.length;

        if (row >= canvasHeight) {
          return;
        }

        window.requestAnimationFrame(slowDraw);
      };

      if (slowMode) {
        slowDraw();
      } else {
        for (let row = 0; row < canvasHeight; row++) {
          const { done, value } = await reader.read();
          drawRow(row, value);
          if (done) {
            break;
          }
        }
      }

      console.timeEnd("Generated terrain");

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
    bind:simpleMode
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
