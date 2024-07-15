<script lang="ts">
  import "xterm/css/xterm.css";

  import * as svelte from "svelte";
  import colorScheme from "./color-schemes.json";
  import { isDark } from "#/libdb.so/site/lib/prefs.js";
  import type * as libterminal from "#/libdb.so/site/lib/terminal.js";

  import Window from "#/libdb.so/site/components/Window.svelte";

  let terminalElement: HTMLElement;

  export let onload: (_: libterminal.Terminal) => void;

  let title = "";
  let terminal: libterminal.Terminal;
  $: theme = colorScheme[$isDark ? "dark" : "light"];

  svelte.onMount(async () => {
    const libterminal = await import("#/libdb.so/site/lib/terminal.js");

    terminal = new libterminal.Terminal({
      fontFamily: `"Inconsolata", "Noto Mono", "Source Code Pro", monospace`,
      fontWeight: "500",
      fontWeightBold: "700",
      lineHeight: 1.1,
      theme,
      drawBoldTextInBrightColors: false,
    });

    terminal.open(terminalElement);
    terminal.write("Starting VM...\r\n");

    const onTitleChange = terminal.onTitleChange((t) => {
      title = t;
    });

    const resizer = new ResizeObserver(() => terminal.fit());
    resizer.observe(terminalElement);

    onload(terminal);

    return () => {
      resizer.disconnect();
      terminal.dispose();
      onTitleChange.dispose();
    };
  });

  // Watch for dark/light theme toggling.
  $: if (terminal && theme) terminal.options.theme = { ...theme };
</script>

<Window view="terminal">
  <h3 slot="title">{title ? `${title} – xterm.js` : "xterm.js"}</h3>
  <div
    class="terminal-box monospace"
    style="
      --background: {theme.background};
      --foreground: {theme.foreground};
    "
  >
    <div class="terminal-box-content" bind:this={terminalElement} />
  </div>
</Window>

<style>
  div.terminal-box {
    height: 100%;
    padding: clamp(4px, 1.5vh, 8px) clamp(0px, 0.5vw, 4px);
    box-sizing: border-box;
    background-color: var(--background);
  }

  div.terminal-box div.terminal-box-content,
  div.terminal-box :global(div.terminal),
  div.terminal-box :global(div.xterm-viewport) {
    height: 100%;
  }

  div.terminal-box :global(.xterm-screen) {
    margin: auto;
  }

  div.terminal-box :global(.xterm-underline-5) {
    text-decoration: dotted underline !important;
    text-decoration-thickness: 0.05em !important;
  }

  div.terminal-box :global(.xterm-underline-5[style="text-decoration: underline;"]) {
    text-decoration: underline !important;
    text-decoration-thickness: 0.05em !important;
  }
</style>
