<script lang="ts">
  export let canvasWidth, noiseScale, noiseSeed, loading, onSubmit;

  let rotate = 0;

  const handleOnSubmit = (event: Event) => {
    event.preventDefault();
    onSubmit();
  };

  const newRandomSeed = (event: Event) => {
    event.preventDefault();
    noiseSeed = Math.floor(Math.random() * 100);
    rotate++;
    onSubmit();
  };
</script>

<form on:submit={handleOnSubmit}>
  <label for="width"> Canvas Width </label>
  <input id="width" type="number" required bind:value={canvasWidth} />
  <label for="scale"> Noise Scale </label>
  <input id="scale" type="number" required bind:value={noiseScale} />
  <label for="seed"> Seed </label>
  <div class="seed">
    <input id="seed" type="number" required bind:value={noiseSeed} />
    <button on:click={newRandomSeed} style="--rotate: {rotate * 360}deg">
      <img src="/rotate-cw.svg" alt="" />
    </button>
  </div>

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
    padding: 0.4em;
    border: 1px solid #ccc;
    border-radius: 4px;
    background: transparent;

    font-family: inherit;
    font-size: inherit;

    color: currentColor;
    transition: background 150ms ease-in-out;
    box-sizing: border-box;
  }

  input:disabled {
    background: #ccc;
  }

  input[type="submit"] {
    grid-column-start: span 2;
    cursor: pointer;
  }

  .seed {
    display: flex;
    flex-flow: row nowrap;
    justify-content: space-between;
    align-items: stretch;
  }
  button {
    border: none;
    margin-left: 1rem;
    background: transparent;

    cursor: pointer;
  }

  button img {
    transform: rotate(var(--rotate));
    transition: transform 200ms ease-in-out;
  }
</style>
