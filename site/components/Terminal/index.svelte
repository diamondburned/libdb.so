<script lang="ts">
  import "xterm/css/xterm.css";

  import * as svelte from "svelte";
  import colorScheme from "./color-schemes.json";
  import type * as xterm from "xterm";

  let terminalElement: HTMLElement;

  export let id: string;
  export let done: (_: xterm.Terminal) => void;
  export let colors: Record<string, string> = {};

  let title = "libdb.so";
  $: combinedColors = { ...colorScheme, ...colors };

  svelte.onMount(async () => {
    const libterminal = await import("#/libdb.so/site/lib/terminal.js");

    const terminal = new libterminal.Terminal({
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

    done(terminal);

    return () => {
      terminal.dispose();
      onTitleChange.dispose();
    };
  });
</script>

<div
  {id}
  class="terminal-box"
  style="
    --background: {combinedColors.background};
    --foreground: {combinedColors.foreground};
  "
>
  <header>
    <h3>{title} – xterm.js</h3>
  </header>
  <div class="terminal monospace" bind:this={terminalElement} />
</div>

<style>
  div.terminal-box {
    overflow: hidden;
    display: flex;
    flex-direction: column;

    box-shadow: 0 2px 16px -6px rgba(0, 0, 0, 0.77);
    border-radius: 15px 15px 0 0;
  }

  header {
    color: black;
    background: linear-gradient(
      to right,
      rgba(85, 205, 252, 1) 0%,
      rgba(147, 194, 255, 1) 25%,
      rgba(200, 181, 245, 1) 50%,
      rgba(234, 171, 217, 1) 75%,
      rgba(247, 168, 184, 1) 100%
    );
    text-align: center;
  }

  header h3 {
    font-weight: bold;
    font-size: 1rem;
    margin: 0.75rem;
    user-select: none;
  }

  div.terminal {
    flex: 1;
    padding: clamp(6px, 1.5vh, 12px) clamp(0px, 0.5vw, 4px);
    background-color: var(--background);
  }

  div.terminal > :global(*) {
    height: 100%;
  }

  div.terminal :global(.xterm-screen) {
    margin: auto;
  }

  div.terminal :global(.xterm-underline-5) {
    text-decoration: dotted underline !important;
    text-decoration-thickness: 0.05em !important;
  }

  div.terminal
    :global(.xterm-underline-5[style="text-decoration: underline;"]) {
    text-decoration: underline !important;
    text-decoration-thickness: 0.05em !important;
  }

  @media (max-width: 500px) {
    header h3 {
      margin: 0.45rem;
    }
  }
</style>
