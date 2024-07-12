<script lang="ts">
  import Popover, { type FlyProps } from "./Popover.svelte";

  // direction is the direction that the popover will be shown.
  export let direction: "top" | "bottom" = "top";
  export let title = "";
  export let open = false;
  export let fly: FlyProps | undefined = undefined;

  export let onclick = () => {
    open = !open;
  };

  let buttonClass = "";
  export { buttonClass as class };
</script>

<div class="popover-container">
  <slot name="button">
    <button
      {title}
      class="popover-button {buttonClass}"
      class:active={open}
      on:click={onclick}
      on:mousedown|stopPropagation
    >
      <slot>
        <span>{title}</span>
      </slot>
    </button>
  </slot>

  <Popover bind:open {direction} {fly}>
    <slot name="popover" />
  </Popover>
</div>

<style lang="scss">
  .popover-container {
    position: relative;
    display: flex;

    :global(.popover-button) {
      box-sizing: border-box;
    }

    :global(.popover) {
      z-index: 100;
      font-size: initial;
      font-weight: initial;
    }
  }
</style>
