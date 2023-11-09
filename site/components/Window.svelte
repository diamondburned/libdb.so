<script lang="ts">
  import { viewDesktop } from "#/libdb.so/site/lib/views.js";
  import WindowMinimize from "#/libdb.so/site/components/Papirus/window-minimize.svelte";
  import WindowMaximize from "#/libdb.so/site/components/Papirus/window-maximize.svelte";
  import WindowRestore from "#/libdb.so/site/components/Papirus/window-restore.svelte";

  export let id = "";
  export let maximized = false;
  export let scrollable = false; // allows content scrolling up and down

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
    viewDesktop();
  }

  function onMaximize() {
    if (maximize) {
      maximize();
      return;
    }
    maximized = !maximized;
  }
</script>

<div class="window-container" class:maximized>
  <main
    class="window"
    class:maximized
    {id}
    style="
	  --max-width: {maxWidth};
	  --max-height: {maxHeight};
	"
  >
    <header class="titlebar" on:dblclick={() => onMaximize()}>
      <div class="title">
        <slot name="title" />
      </div>
      <div class="controls">
        {#if minimize !== null}
          <button class="minimize" on:click={onMinimize}>
            <WindowMinimize />
          </button>
        {/if}
        {#if maximize !== null}
          <button class="maximize" on:click={onMaximize}>
            {#if maximized}
              <WindowRestore />
            {:else}
              <WindowMaximize />
            {/if}
          </button>
        {/if}
      </div>
    </header>
    <div class="content" class:scrollable>
      <slot />
    </div>
  </main>
</div>

<style lang="scss">
  div.window-container {
    width: 100%;
    height: 100%;

    &:not(.maximized) {
      display: flex;
      justify-content: center;
      align-items: center;
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

    @mixin maximize {
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

    .content {
      flex: 1;
      overflow: hidden;

      &.scrollable {
        overflow-y: auto;
        height: 100%;
      }
    }

    header.titlebar {
      color: white;
      background-color: #030303;
      border-bottom: 1px solid var(--window-border-color);

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

        button {
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
      }
    }
  }
</style>
