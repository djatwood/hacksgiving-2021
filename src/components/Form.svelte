<script lang="ts">
  export let canvasWidth, noiseScale, noiseSeed, loading, onSubmit;

  const handleOnSubmit = (event: Event) => {
    event.preventDefault();
    onSubmit();
  };
</script>

<form on:submit={handleOnSubmit}>
  <label for="width"> Canvas Width </label>
  <input id="width" type="number" required bind:value={canvasWidth} />
  <label for="scale"> Noise Scale </label>
  <input id="scale" type="number" required bind:value={noiseScale} />
  <label for="seed"> Seed </label>
  <input id="seed" type="number" required bind:value={noiseSeed} />

  <input
    type="submit"
    value={loading ? "Generating" : "Generate"}
    disabled={loading}
  />
</form>

<style>
  form {
    display: grid;
    grid-template-columns: auto 1fr;
    grid-gap: 1rem;

    color: white;
  }

  label {
    display: flex;
    flex-flow: column nowrap;
    justify-content: center;
    align-items: flex-start;

    cursor: pointer;
  }

  input {
    background: transparent;
    border-radius: 4px;

    color: currentColor;
    transition: background 150ms ease-in-out;
  }

  input:disabled {
    background: #ccc;
  }

  input[type="submit"] {
    grid-column-start: span 2;
    cursor: pointer;
  }
</style>
