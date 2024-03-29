<script lang="ts">
  import * as svelte from "svelte";
  import * as store from "svelte/store";

  import {
    View,
    DragState,
    toggleView,
    viewWindows,
    bringToFocus,
    viewIsActive,
    viewIsFocused,
  } from "#/libdb.so/site/lib/views.js";
  import { dragWindows } from "#/libdb.so/site/lib/prefs.js";
  import WindowMinimize from "#/libdb.so/site/components/Papirus/window-minimize.svelte";
  import WindowMaximize from "#/libdb.so/site/components/Papirus/window-maximize.svelte";
  import WindowRestore from "#/libdb.so/site/components/Papirus/window-restore.svelte";
  import WindowControl from "#/libdb.so/site/components/WindowControl.svelte";

  export let view: View;
  export let maximized = false;
  export let scrollable = false; // allows content scrolling up and down

  export let windowClass = "";
  export let headerClass = "";
  export let contentClass = "";

  export let maxWidth = "max(80vw, 1000px)";
  export let maxHeight = "max(80vh, 600px)";

  // minimize sets the minimize callback. If null, then the minimize button is
  // hidden. If undefined, then the minimize button minimizes the window.
  export let minimize: undefined | null | (() => void) = undefined;
  // maximize behaves like minimize, except undefined means the button is
  // visible but disabled.
  export let maximize: undefined | null | (() => void) = undefined;

  function onMinimize() {
    if (minimize) {
      minimize();
      return;
    }
    toggleView(view);
  }

  function onMaximize() {
    if (maximize) {
      maximize();
      return;
    }
    maximized = !maximized;
  }

  export let x = 0; // position of the window
  export let y = 0; // position of the window
  let moved = false; // whether the window has been moved
  let active = viewIsActive(view); // whether the window is visible
  let focused = viewIsFocused(view); // whether the window is focused

  let windowWidth = 0; // width of the window
  let windowHeight = 0; // height of the window
  let containerWidth = 0; // width of the container
  let containerHeight = 0; // height of the container

  $: {
    viewWindows.update((windows) => {
      windows[view] = {
        x,
        y,
        width: windowWidth,
        height: windowHeight,
      };
      return windows;
    });
  }

  function clamp(min: number, val: number, max: number) {
    return Math.min(Math.max(val, min), max);
  }

  function clampPositions() {
    const minVisible = 100; // minimum visible area per dimension
    x = clamp(minVisible - windowWidth, x, containerWidth - minVisible);
    // Specifically prevent the window from being dragged above the top of the
    // screen. This is to prevent the headerbar from being inaccessible.
    y = clamp(0, y, containerHeight - minVisible);
  }

  $: {
    if (!moved) {
      x = (containerWidth - windowWidth) / 2;
      y = (containerHeight - windowHeight) / 2;
      clampPositions();
    }
  }

  let dragState: DragState | null = null; // null when not dragging

  function dragBegin(ev: MouseEvent) {
    if (!$dragWindows) {
      return;
    }
    dragState = new DragState(x, y, ev.clientX, ev.clientY, (newX, newY) => {
      x = newX;
      y = newY;
      clampPositions();
    });
  }

  function dragEnd() {
    dragState = null;
  }

  function drag(ev: MouseEvent) {
    if (dragState) {
      if (!moved) moved = true;
      dragState.update(ev.clientX, ev.clientY);
    }
  }

  function checkMouseEnter(ev: MouseEvent) {
    if (dragState && !ev.buttons) {
      // Stop dragging if the mouse enters the window while not dragging.
      // This prevents the window from being stuck in a dragging state.
      dragEnd();
    }
  }
</script>

<svelte:body
  on:mouseenter={checkMouseEnter}
  on:mouseup={dragEnd}
  on:mousemove={drag}
/>

<div
  bind:clientWidth={containerWidth}
  bind:clientHeight={containerHeight}
  class="window-container {windowClass}"
  class:maximized
  class:focused={$focused}
  class:active={$active}
>
  <main
    id={view}
    bind:clientWidth={windowWidth}
    bind:clientHeight={windowHeight}
    on:mousedown={(ev) => bringToFocus(view)}
    class="window"
    class:maximized
    class:dragging={dragState !== null}
    style="
      --max-width: {maxWidth};
      --max-height: {maxHeight};
      top: {y}px;
      left: {x}px;
    "
  >
    <header
      class="titlebar {headerClass}"
      on:dblclick={onMaximize}
      on:mousedown={(ev) => dragBegin(ev)}
    >
      <div class="title">
        <slot name="title" />
      </div>
      <div class="controls">
        {#if minimize !== null}
          <WindowControl class="minimize" clicked={onMinimize}>
            <WindowMinimize />
          </WindowControl>
        {/if}
        {#if maximize !== null}
          <WindowControl class="maximize" clicked={onMaximize}>
            {#if maximized}
              <WindowRestore />
            {:else}
              <WindowMaximize />
            {/if}
          </WindowControl>
        {/if}
      </div>
    </header>
    <div class="content-wrapper">
      <div class="overlays">
        <slot name="overlay" />
      </div>
      <div class="content {contentClass}" class:scrollable>
        <slot />
      </div>
    </div>
  </main>
</div>

<style lang="scss">
  div.window-container {
    width: 100%;
    height: 100%;
    position: absolute;

    pointer-events: none;

    --ease-function: cubic-bezier(0, 1.005, 0.165, 1);
    --ease-duration: 0.2s;

    @keyframes minimize-animation {
      /*
       * Do a silly little hack here. When the window is minimized, its
       * scale is 0, which causes xterm.js to get the completely wrong
       * size.
       *
       * To get around this, we'll scale to 0.0001 at the 99% mark, and
       * then instantly set opacity to 0 so that we can scale it back to 1
       * at the end.
       */
      0% {
        opacity: 1;
        transform: translateY(0) scale(1);
      }
      98% {
        opacity: 1;
      }
      99% {
        opacity: 0;
        transform: translateY(50vh) scale(0.0001);
      }
      100% {
        opacity: 0;
        transform: translateY(50vh) scale(1);
      }
    }

    &:not(.active) {
      opacity: 0;
      animation: minimize-animation var(--ease-duration) var(--ease-function);
    }

    @keyframes unminimize-animation {
      /*
       * Repeat the same silly hack here.
       */
      0% {
        opacity: 0;
        transform: translateY(50vh) scale(1);
      }
      1% {
        opacity: 0;
        transform: translateY(50vh) scale(0.0001);
      }
      100% {
        opacity: 1;
        transform: translateY(0) scale(1);
      }
    }

    &.active {
      animation: unminimize-animation var(--ease-duration) var(--ease-function);

      main.window {
        pointer-events: auto;
      }
    }

    @media (prefers-reduced-motion: reduce) {
      --ease-duration: 0s;
    }

    &.focused {
      z-index: 10;
    }
  }

  main.window {
    --window-border-color: rgba(255, 255, 255, 0.1);

    overflow: hidden;
    display: flex;
    flex-direction: column;

    box-shadow: 0 2px 16px -6px rgba(0, 0, 0, 0.77);
    border-radius: 12px;

    outline: 1px solid var(--window-border-color);
    outline-offset: -1px;

    width: min(calc(100% - clamp(6px, 5vw, 3rem)), var(--max-width));
    height: min(calc(100% - clamp(6px, 7vw, 5rem)), var(--max-height));

    color: #fcfcfc;
    background-color: #0f0f0f;

    position: absolute;
    top: 0;
    left: 0;

    @mixin maximize {
      position: initial;
      border-radius: 0;
      width: 100%;
      height: 100%;
      outline: none;
    }

    &.maximized {
      @include maximize;
    }

    @media (max-width: 500px) {
      @include maximize;
    }

    &.dragging {
      cursor: move;
    }

    .content-wrapper {
      flex: 1;
      overflow: hidden;

      /* Hack to make position: fixed work.
       * See https://stackoverflow.com/a/38796408/5041327. */
      transform: translateZ(0);
    }

    .overlays {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      z-index: 10;
      pointer-events: none;

      & > :global(*) {
        pointer-events: auto;
      }
    }

    .content {
      height: 100%;
      overflow: hidden;

      &.scrollable {
        overflow-y: auto;
        height: 100%;
      }
    }

    header.titlebar {
      color: white;
      background-color: #030303;
      box-shadow: 0 -4px 6px 6px rgba(0, 0, 0, 0.47);

      text-align: center;

      display: flex;
      flex-direction: row;
      justify-content: space-between;

      user-select: none;

      :global(h1),
      :global(h2),
      :global(h3) {
        font-weight: bold;
        font-size: 1rem;

        margin: 0.65rem;
        @media (max-width: 500px) {
          margin: 0.45rem;
        }
      }

      .title {
        flex: 1;
      }

      .controls {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 0.5rem;

        margin: 0 0.65rem;
        @media (max-width: 500px) {
          margin: 0 0.45rem;
        }
      }
    }
  }
</style>
