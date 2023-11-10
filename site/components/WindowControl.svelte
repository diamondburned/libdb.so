<script lang="ts">
  import * as svelte from "svelte";

  let className = "";
  export { className as class };

  export let clicked: () => void = () => {};
</script>

<button class="window-control {className}" on:click={clicked}>
  <slot />
</button>

<style lang="scss">
  button.window-control {
    color: white;
    border: none;
    border-radius: 99px;
    width: 1.5rem;
    height: 1.5rem;
    font-size: 0.75rem;
    font-weight: 900;
    line-height: 0;
    background-color: rgba(255, 255, 255, 0.1);
    transition: all 0.1s ease-in-out;
    padding: 0;

    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      background-color: rgba(255, 255, 255, 0.2);
    }

    &.minimize {
      :global(svg *) {
        fill: var(--blue);
      }
    }

    &.maximize {
      :global(svg *) {
        fill: var(--pink);
      }
      @media (max-width: 500px) {
        /* Always maximize on mobile */
        display: none;
      }
    }
  }
</style>
