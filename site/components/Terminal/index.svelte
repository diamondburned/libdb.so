<script lang="ts">
  import "@xterm/xterm/css/xterm.css";

  import * as svelte from "svelte";
  import colorScheme from "./color-schemes.json";
  import { fade } from "svelte/transition";
  import { isDark } from "#/libdb.so/site/lib/prefs.js";
  import type * as libterminal from "#/libdb.so/site/lib/terminal.js";

  import Window from "#/libdb.so/site/components/Window.svelte";
  import { readable } from "svelte/store";

  let terminalElement: HTMLElement;

  let terminal: libterminal.Terminal;
  $: title = terminal?.title || readable("");
  $: theme = colorScheme[$isDark ? "dark" : "light"];

  let destroy: () => void | undefined;
  let destroyed = false;
  svelte.onDestroy(() => {
    destroy && destroy();
    destroyed = true;
  });

  async function initTerminal() {
    try {
      const [libterminal, vm] = await Promise.all([
        import("#/libdb.so/site/lib/terminal.js"),
        import("#/libdb.so/site/lib/vm.js"),
      ]);

      terminal = new libterminal.Terminal({
        theme,
        drawBoldTextInBrightColors: false,
      });

      await terminal.open(terminalElement);
      destroy = () => terminal.dispose();

      const url = new URL(location.href);
      const localhost = url.hostname == "localhost" || !url.hostname;

      await vm.start(terminal.xterm!, {
        publicFSURLs: localhost ? ["/_fs.json"] : [],
      });
    } catch (err) {
      console.error("Failed to initialize terminal:", err);
      throw err;
    } finally {
      // In case the component is destroyed before initialization:
      if (destroyed) destroy();
    }
  }

  // Watch for dark/light theme toggling.
  $: if (terminal && theme) terminal.setTheme(theme);
</script>

<Window view="terminal">
  <h3 slot="title">{$title ? `${$title} – xterm.js` : "xterm.js"}</h3>
  <div
    class="terminal-box"
    style="
      --background: {theme.background};
      --foreground: {theme.foreground};
    "
  >
    {#await initTerminal()}
      <p class="status loading" out:fade={{ duration: 150 }}>Initializing terminal...</p>
    {:catch}
      <p class="status error" transition:fade={{ duration: 150 }}>
        Failed to initialize terminal. Please check DevTools.
      </p>
    {/await}
    <div class="monospace terminal-box-content" bind:this={terminalElement} />
  </div>
</Window>

<style lang="scss">
  div.terminal-box {
    height: 100%;
    box-sizing: border-box;
    background-color: var(--background);

    position: relative;

    p.status {
      position: absolute;
      width: 100%;
      height: 100%;

      background-color: var(--background);

      display: flex;
      align-items: center;
      justify-content: center;

      &.error {
        color: var(--adw-destructive-color);
      }
    }

    .terminal-box-content {
      padding: clamp(4px, 1.5vh, 8px) clamp(0px, 0.5vw, 4px);
    }

    .terminal-box-content,
    :global(div.terminal),
    :global(div.xterm-viewport) {
      height: 100%;
    }

    :global(.xterm-screen) {
      margin: auto;
    }

    :global(.xterm-underline-5) {
      text-decoration: dotted underline !important;
      text-decoration-thickness: 0.05em !important;
    }

    :global(.xterm-underline-5[style="text-decoration: underline;"]) {
      text-decoration: underline !important;
      text-decoration-thickness: 0.05em !important;
    }
  }
</style>
