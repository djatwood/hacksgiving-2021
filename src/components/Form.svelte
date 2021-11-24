<script lang="ts">
  export let canvasWidth, loading, noiseScale, noiseSeed, onSubmit, renderMode;

  let rotate = 0;

  const handleOnSubmit = (event: Event) => {
    event.preventDefault();
    onSubmit();
  };

  const newRandomSeed = (event: Event) => {
    event.preventDefault();
    if (loading) return;
    noiseSeed = Math.floor(Math.random() * 100);
    rotate++;
    onSubmit();
  };
</script>

<form on:submit={handleOnSubmit}>
  <label for="width">Canvas Width</label>
  <input id="width" type="number" required bind:value={canvasWidth} />
  <label for="scale">Noise Scale</label>
  <input id="scale" type="number" required bind:value={noiseScale} />
  <label for="seed">Seed </label>
  <div class="seed">
    <button on:click={newRandomSeed} style="--rotate: {rotate * 360}deg">
      <img src="rotate-cw.svg" alt="" />
    </button>
    <input id="seed" type="number" required bind:value={noiseSeed} />
  </div>

  <div class="render">
    <p>Render mode</p>
    <div class="row">
      <div class="radio">
        <label for="live">Live</label>
        <input
          type="radio"
          name="render"
          id="live"
          bind:group={renderMode}
          value="live"
        />
      </div>
      <div class="radio">
        <label for="smooth">Smooth</label>
        <input
          type="radio"
          name="render"
          id="smooth"
          bind:group={renderMode}
          value="smooth"
        />
      </div>
      <div class="radio">
        <label for="full">Full</label>
        <input
          type="radio"
          name="render"
          id="full"
          bind:group={renderMode}
          value="full"
        />
      </div>
    </div>
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
    width: 100%;
    padding: 0.5rem 0.6rem;
    border: 1px solid #ccc;
    border-radius: 6px;
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
    margin-top: 4px;
  }

  input[type="radio"] {
    width: auto;
  }

  .seed {
    display: flex;
    flex-flow: row nowrap;
    justify-content: space-between;
    align-items: stretch;
  }
  button {
    display: flex;
    flex-flow: column nowrap;
    justify-content: center;
    align-items: center;

    border: none;
    margin-right: 1rem;
    background: transparent;

    cursor: pointer;
  }

  button img {
    transform: rotate(var(--rotate));
    transition: transform 200ms ease-in-out;
  }

  .render {
    grid-column-start: span 2;
  }

  .render p {
    margin-top: 0;
  }

  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    grid-column-start: span 2;
  }

  .radio {
    display: flex;
    flex-flow: row nowrap;
    justify-content: space-between;
  }
  .radio label {
    order: 1;
    margin-left: 4px;
  }
</style>
