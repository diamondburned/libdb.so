<script lang="ts" context="module">
  // Hard-coded size of the popover arrow.
  const arrowSize = 8;
  const borderWidth = 1;

  // Enforce that there has to be at least 8px of margin between the popover
  // and the edge of the screen.
  const popoverEdgeMargin = 8;

  type Pt = { x: number; y: number };
  const zeroPt: Pt = { x: 0, y: 0 };

  type Size = { width: number; height: number };
  const zeroSize: Size = { width: 0, height: 0 };

  export type FlyProps = {
    x?: number;
    y?: number;
    duration?: number;
  };
</script>

<script lang="ts">
  import { onMount, tick } from "svelte";
  import { fly } from "svelte/transition";
  import { boundingWindow } from "./Window.svelte";

  // direction is the direction that the popover will be shown.
  export let direction: "top" | "bottom" = "top";

  // open is whether the popover is open.
  export let open = true;

  // The settings for the fly transition.
  // duration is optional and defaults to 200.
  let flyProps2: FlyProps | undefined = undefined;
  export { flyProps2 as fly };

  $: flyProps = {
    x: flyProps2?.x,
    y: flyProps2?.y ?? (direction === "top" ? 10 : -10),
    duration: flyProps2?.duration ?? 200,
  };

  // This element is never rendered, but it's always in the DOM tree.
  // It is used to get the parent elements of the popover.
  let dummyElement: HTMLElement;
  let parentSize: DOMRect = new DOMRect();
  onMount(() => {
    const parent = boundingWindow(dummyElement);
    const update = () => {
      parentSize = parent.getBoundingClientRect();
    };
    update();

    const observer = new ResizeObserver(update);
    observer.observe(parent);
    return () => observer.disconnect();
  });

  let popover: HTMLDivElement | undefined = undefined;
  let popoverSize = zeroSize;
  let overflowOffset = zeroPt;
  let arrowX = 0;

  $: {
    open;
    popover;
    parentSize;

    (async () => {
      await tick();
      if (popover && open) {
        const rect = popover.getBoundingClientRect();
        const overflowX = rect.right - parentSize.left - parentSize.width + popoverEdgeMargin;
        const overflowY = rect.bottom - parentSize.top - parentSize.height + popoverEdgeMargin;
        overflowOffset = {
          x: Math.max(0, overflowX),
          y: Math.max(0, overflowY),
        };
      } else {
        overflowOffset = zeroPt;
        return;
      }

      await tick();
      if (popover && open) {
        const popoverOffset = {
          x: popover.offsetLeft,
          y: popover.offsetTop,
        };
        const parent = dummyElement.parentElement;
        const parentWidth = parent?.clientWidth ?? 0;
        arrowX = -popoverOffset.x + parentWidth / 2 - arrowSize / 4;
      } else {
        arrowX = 0;
      }
    })();
  }
</script>

<!--
  Close the popover if the user clicks outside of it.
  This only works if the popover also has an on:mousedown|stopPropagation
  so that it swallows mouse clicks inside the popover.
-->
<svelte:window
  on:mousedown|stopPropagation={() => {
    if (open) {
      open = false;
    }
  }}
/>

<div
  role="presentation"
  class="popover-dummy"
  bind:this={dummyElement}
  on:mousedown|stopPropagation
>
  {#if open}
    <div
      role="tooltip"
      class="popover"
      bind:this={popover}
      bind:offsetWidth={popoverSize.width}
      bind:offsetHeight={popoverSize.height}
      class:popover-top={direction === "top"}
      class:popover-bottom={direction === "bottom"}
      class:arrow-relative={arrowX == undefined}
      class:arrow-absolute={arrowX != undefined}
      style="
      --arrow-x: {arrowX}px;
      --arrow-size: {arrowSize}px;
      --width: {popoverSize.width}px;
      --height: {popoverSize.height}px;
      --overflow-x: {overflowOffset.x}px;
      --overflow-y: {overflowOffset.y}px;
    "
      transition:fly={flyProps}
    >
      <slot />
    </div>
  {/if}
</div>

<style lang="scss">
  .popover {
    position: absolute;
    left: calc(50% - var(--width) / 2 - var(--overflow-x, 0px));

    &:not(.no-arrow) {
      --arrow-y: calc(var(--arrow-size));
    }

    &.popover-top {
      top: calc(-1 * var(--height) - var(--offset-y, 0px) - var(--arrow-y, 0));
    }

    &.popover-bottom {
      bottom: calc(-1 * var(--height) - var(--offset-y, 0px) - var(--arrow-y, 0));
    }

    z-index: 10;

    background: var(--adw-popover-bg-color);
    border: 1px solid rgba(0 0 0 / 14%);
    box-shadow:
      0 1px 5px 1px rgba(0 0 0 / 9%),
      0 2px 14px 3px rgba(0 0 0 / 5%);
    box-sizing: border-box;
    padding: 8px;
    margin: 0 var(--margin-x);
    border-radius: var(--adw-popover-radius);

    max-width: min(250px, calc(100vw - 32px));

    &:not(.no-arrow)::after {
      content: " ";
      border-top: var(--arrow-size) solid var(--adw-popover-bg-color);
      border-left: var(--arrow-size) solid transparent;
      border-right: var(--arrow-size) solid transparent;
      position: absolute;
      bottom: calc(-1 * var(--arrow-size));
      left: calc(var(--arrow-x) - var(--adw-menu-margin) - var(--margin-x, 0px) - 2px);
    }
  }
</style>
