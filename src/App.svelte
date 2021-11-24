<script lang="ts">
  import { onMount } from "svelte";
  import type { writable } from "svelte/store";

  import Form from "./components/Form.svelte";

  export let generate: ReturnType<typeof writable>;

  let loading = false;
  let canvas: HTMLCanvasElement;
  let canvasCtx: CanvasRenderingContext2D;

  let canvasWidth = Math.round(1024 * window.devicePixelRatio);
  let canvasHeight: number;
  $: canvasHeight = Math.round(canvasWidth * (2 / 3));

  let noiseSeed = 73;
  let noiseScale = 128 * window.devicePixelRatio;
  let renderMode: "live" | "smooth" | "full" = "live";
  let simpleMode = false;

  const generateNoise = () => {
    if (!$generate || loading) return;

    loading = true;
    setTimeout(async () => {
      console.time("Rendered terrain");

      const response = await $generate(
        canvasWidth,
        canvasHeight,
        noiseScale,
        noiseSeed,
        renderMode === "smooth" ? Math.floor(canvasHeight * 0.1) : 1
      );

      if (canvas.width != canvasWidth) canvas.width = canvasWidth;
      if (canvas.height != canvasHeight) canvas.height = canvasHeight;

      const body: ReadableStream<Uint8ClampedArray> = await response.body;
      const reader = body.getReader();

      const drawRow = (row: number, data?: Uint8ClampedArray): number => {
        if (!data) {
          return;
        }
        const imageData = new ImageData(data, canvasWidth);
        canvasCtx.putImageData(imageData, 0, row);
        return row + data.length / 4 / canvasWidth;
      };

      const smoothDraw = () => {
        let row = 0;
        const draw = async (t: number) => {
          const { done, value } = await reader.read();
          row = drawRow(row, value);
          if (done) {
            return;
          }

          window.requestAnimationFrame(draw);
        };
        draw(0);
      };

      if (renderMode === "smooth") {
        smoothDraw();
      } else {
        let row = 0;
        while (row < canvasHeight) {
          const { done, value } = await reader.read();

          if (renderMode === "live") {
            setTimeout(() => (row = drawRow(row, value)), 0);
          } else {
            row = drawRow(row, value);
          }

          if (done) {
            break;
          }
        }
      }

      console.log(
        "chunkSize:",
        renderMode === "smooth" ? Math.floor(canvasHeight * 0.1) : 1
      );
      console.timeEnd("Rendered terrain");

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
    bind:renderMode
    bind:simpleMode
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
