<script lang="ts">
  import { fly } from "svelte/transition";
  import { onMount } from "svelte";

  // direction is the direction that the popover will be shown.
  export let direction: "top" | "bottom" = "top";
  export let open = false;

  let buttonHeight = 0;
  let maxX = 0;
  let maxY = 0;

  function getMargins(el: HTMLElement) {
    const style = window.getComputedStyle(el);
    return {
      marginTop: parseInt(style.marginTop, 10),
      marginBottom: parseInt(style.marginBottom, 10),
      marginLeft: parseInt(style.marginLeft, 10),
      marginRight: parseInt(style.marginRight, 10),
    };
  }

  function clamp(min: number, val: number, max: number) {
    return Math.min(Math.max(val, min), max);
  }

  let x = 0;
  let y = 0;

  let popover: HTMLDivElement;
  let popoverWidth = 0;
  let popoverHeight = 0;

  let container: HTMLDivElement;

  $: {
    popoverWidth;
    popoverHeight;

    if (container && popover) {
      const buttonRect = container.getBoundingClientRect();
      const { width: popoverWidth, height: popoverHeight } =
        popover.getBoundingClientRect();
      const { marginLeft, marginRight, marginTop, marginBottom } =
        getMargins(popover);

      switch (direction) {
        case "top": {
          const anchorX = buttonRect.left + buttonRect.width / 2;
          const anchorY = buttonRect.top;

          x = anchorX - popoverWidth / 2 - marginLeft;
          y = anchorY - popoverHeight - marginBottom * 2;

          break;
        }
        case "bottom": {
          const anchorX = buttonRect.left + buttonRect.width / 2;
          const anchorY = buttonRect.bottom;

          x = anchorX - popoverWidth / 2 - marginLeft;
          y = anchorY + marginBottom * 2;

          break;
        }
      }

      x = clamp(0, x, maxX - popoverWidth - marginRight * 2 - 2);
      y = clamp(0, y, maxY - popoverHeight - marginBottom * 2 - 2);
    }
  }
</script>

<svelte:window
  bind:innerWidth={maxX}
  bind:innerHeight={maxY}
  on:mousedown={() => {
    if (open) {
      open = false;
    }
  }}
/>

<div class="popover-container" bind:this={container}>
  <button
    class="popover-button"
    class:active={open}
    on:click={(ev) => (open = !open)}
    on:mousedown={(ev) => ev.stopPropagation()}
    bind:offsetHeight={buttonHeight}
  >
    <slot />
  </button>

  {#if open}
    <div
      role="tooltip"
      class="popover"
      class:popover-top={direction === "top"}
      class:popover-bottom={direction === "bottom"}
      transition:fly={{ y: buttonHeight / 2, duration: 100 }}
      bind:this={popover}
      bind:offsetWidth={popoverWidth}
      bind:offsetHeight={popoverHeight}
      on:mousedown={(ev) => ev.stopPropagation()}
      style="--top: {y}px; --left: {x}px;"
    >
      <slot name="popover" />
    </div>
  {/if}
</div>

<style lang="scss">
  .popover-container {
    position: relative;
    display: flex;
  }

  .popover-button {
    box-sizing: border-box;
  }

  .popover {
    z-index: 100;

    position: fixed;
    top: var(--top);
    left: var(--left);

    font-size: initial;
    font-weight: initial;

    margin: 0.5em;
    padding: 0 1em;
    outline: 1px solid rgba(255, 255, 255, 0.1);
    outline-offset: -1px;
    background-color: #0f0f0f;
    border-radius: 12px;
  }
</style>
