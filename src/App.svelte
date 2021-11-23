<script lang="ts">
  import { onMount } from "svelte";

  const canvasWidth = 1024 * window.devicePixelRatio;
  const canvasHeight = canvasWidth * (9 / 16);

  let canvas: HTMLCanvasElement;

  onMount(() => {
    canvas.width = canvasWidth;
    canvas.height = canvasHeight;

    const ctx = canvas.getContext("2d");

    const data = new Uint8ClampedArray(4 * canvasWidth * canvasHeight);

    for (let i = 0; i < data.length; i += 4) {
      const value = (255 * ((i / 4) % canvasWidth)) / canvasWidth;
      data[i] = value;
      data[i + 1] = value;
      data[i + 2] = value;
      data[i + 3] = 255;
    }

    const imageData = new ImageData(data, canvasWidth, canvasHeight);
    ctx.putImageData(imageData, 0, 0);
  });
</script>

<main>
  <canvas
    bind:this={canvas}
    width={canvasWidth}
    style="--max-width: {canvasWidth / window.devicePixelRatio}px"
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

    background: #567;
  }

  canvas {
    max-width: var(--max-width);
    width: 100%;
    background: black;
  }
</style>
