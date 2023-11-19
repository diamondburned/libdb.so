<script lang="ts">
  export let checked = false;
  export let width = "35px";

  let maxX = 0;
  let maxY = 0;
</script>

<label style="--width: {width}">
  <span class="label"><slot /></span>
  <div class="switch">
    <input type="checkbox" bind:checked />
    <span class="slider" />
  </div>
</label>

<style lang="scss">
  label {
    white-space: nowrap;

    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 0.5em;
    cursor: pointer;
  }

  .label {
    flex: 1;
    margin-right: 0.25em;
  }

  .switch {
    --width: 45px;
    --height: calc(var(--width) / 1.65);
    --transition: 0.15s ease-in-out;

    display: flex;
    flex-direction: row;
    align-items: center;

    position: relative;
    width: var(--width, 0);
    height: var(--height);
    margin-left: auto;
  }

  .switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  /* The slider */
  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #535353;
    transition: 0.4s;
    border-radius: var(--height);
  }

  .slider:before {
    --slider-size: calc(var(--height) * 0.85);
    --slider-offset: calc((var(--height) - var(--slider-size)) / 2);

    position: absolute;
    content: "";

    width: var(--slider-size);
    height: var(--slider-size);
    left: var(--slider-offset);
    bottom: var(--slider-offset);

    transition: var(--transition);
    background-color: white;
    border-radius: var(--slider-size);
    box-shadow: 0 0 2px rgba(0, 0, 0, 0.25);
  }

  input:checked + .slider {
    background-color: #2196f3;
  }

  input:focus + .slider {
    box-shadow: 0 0 1px #2196f3;
  }

  input:checked + .slider:before {
    transform: translateX(
      calc(
        var(--width) - var(--slider-size) - var(--slider-offset) -
          var(--slider-offset)
      )
    );
  }
</style>
