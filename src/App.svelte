<script lang="ts">
  import { onMount } from "svelte";

  export let name: string;

  const canvasWidth = 1024;

  let canvas;

  onMount(() => {
    canvas.width = canvasWidth;
    canvas.height = canvasWidth * (9 / 16);

    console.log(canvas.width, canvas.height);

    const ctx = canvas.getContext("2d");
    // ctx.fillRect(0, 0, canvas.width, canvas.height);

    const data = new Uint8ClampedArray(4 * canvas.width * canvas.height);

    for (let i = 0; i < data.length; i += 4) {
      const value = (255 * ((i / 4) % canvas.width)) / canvasWidth;
      data[i] = value;
      data[i + 1] = value;
      data[i + 2] = value;
      data[i + 3] = 255;
    }

    console.log(data);

    const imageData = new ImageData(data, canvas.width, canvas.height);

    console.log(imageData);
    ctx.putImageData(imageData, 0, 0);
  });
</script>

<main>
  <canvas
    bind:this={canvas}
    width={canvasWidth}
    style="--max-width: {canvasWidth}px"
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
    /* aspect-ratio: 16/9; */
    width: 100%;
  }
</style>
