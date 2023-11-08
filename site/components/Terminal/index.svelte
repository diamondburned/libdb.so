<script lang="ts">
  import "xterm/css/xterm.css";

  import * as svelte from "svelte";
  import colorScheme from "./color-schemes.json";
  import type * as xterm from "xterm";

  import Window from "#/libdb.so/site/components/Window.svelte";

  let terminalElement: HTMLElement;

  export let id: string;
  export let done: (_: xterm.Terminal) => void;
  export let colors: Record<string, string> = {};

  let title = "";
  let terminal: xterm.Terminal;
  $: combinedColors = { ...colorScheme, ...colors };

  svelte.onMount(async () => {
    const libterminal = await import("#/libdb.so/site/lib/terminal.js");

    terminal = new libterminal.Terminal({
      fontFamily: `"Inconsolata", "Noto Mono", "Source Code Pro", monospace`,
      fontWeight: "500",
      fontWeightBold: "700",
      lineHeight: 1.1,
      theme: combinedColors,
      drawBoldTextInBrightColors: false,
      linkHandler: {
        activate: (event: MouseEvent, uri: string) => {
          window.open(uri, "_blank");
        },
      },
    });

    terminal.open(terminalElement);
    terminal.write("Starting VM...\r\n");

    const onTitleChange = terminal.onTitleChange((t) => {
      title = t;
    });

    const resizer = new ResizeObserver(() => terminal.fit());
    resizer.observe(terminalElement);

    done(terminal);

    return () => {
      resizer.disconnect();
      terminal.dispose();
      onTitleChange.dispose();
    };
  });
</script>

<Window {id}>
  <h3 slot="title">{title ? `${title} – xterm.js` : "xterm.js"}</h3>
  <div
    class="terminal-box monospace"
    style="
      --background: {combinedColors.background};
      --foreground: {combinedColors.foreground};
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

  div.terminal-box
    :global(.xterm-underline-5[style="text-decoration: underline;"]) {
    text-decoration: underline !important;
    text-decoration-thickness: 0.05em !important;
  }
</style>
