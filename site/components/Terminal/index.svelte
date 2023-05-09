<script lang="ts">
  import "xterm/css/xterm.css";

  import { WebglAddon as WebGLAddon } from "xterm-addon-webgl";
  import { ImageAddon, IImageAddonOptions } from "xterm-addon-image";
  import { FitAddon } from "xterm-addon-fit";
  import * as svelte from "svelte";
  import colorScheme from "./color-schemes.json";
  import type * as xterm from "xterm";

  let terminalElement: HTMLElement;

  export let id: string;
  export let done: (_: xterm.Terminal) => void;
  export let colors: Record<string, string> = {};

  let title = "libdb.so";
  $: combinedColors = { ...colorScheme, ...colors };

  const imageAddonSettings: IImageAddonOptions = {
    enableSizeReports: true,
    sixelSupport: true,
    sixelScrolling: true,
    sixelPaletteLimit: 4096,
    showPlaceholder: true,
  };

  svelte.onMount(async () => {
    const xterm = await import("xterm");

    const terminal = new xterm.Terminal({
      fontFamily: "monospace",
      fontWeight: "500",
      fontWeightBold: "700",
      lineHeight: 1.1,
      allowTransparency: true,
      convertEol: true,
      theme: combinedColors,
      drawBoldTextInBrightColors: false,
      linkHandler: {
        activate: (event: MouseEvent, uri: string) => {
          window.open(uri, "_blank");
        },
      },
    });

    terminal.attachCustomKeyEventHandler(function (e) {
      // Bind Ctrl + C to copy if there is a selection.
      if (e.ctrlKey && e.key == "c") {
        if (terminal.hasSelection()) {
          console.log("copying");
          navigator.clipboard.writeText(terminal.getSelection());
          return false;
        }
      }

      // Bind Ctrl + V to paste for consistency with Ctrl + C. We don't actually
      // need to do anything here because the browser will automatically paste
      // the clipboard contents into the terminal.
      if (e.ctrlKey && e.key == "v") {
        return false;
      }

      return true;
    });

    const fitAddon = new FitAddon();
    const onResize = () => fitAddon.fit();
    terminal.loadAddon(fitAddon);

    const imageAddon = new ImageAddon(imageAddonSettings);
    terminal.loadAddon(imageAddon);

    try {
      const webglAddon = new WebGLAddon();
      webglAddon.onContextLoss(() => webglAddon.dispose());
      terminal.loadAddon(webglAddon);
    } catch (err) {
      console.log("WebGL not supported, using canvas fallback");
    }

    terminal.open(terminalElement);
    terminal.write("Starting VM...\r\n");

    onResize();
    window.addEventListener("resize", onResize);

    const onTitleChange = terminal.onTitleChange((t) => {
      title = t;
    });

    done(terminal);

    return () => {
      terminal.dispose();
      webglAddon.dispose();
      fitAddon.dispose();
      onTitleChange.dispose();
      window.removeEventListener("resize", onResize);
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
  <div class="terminal" bind:this={terminalElement} />
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

  @media (max-width: 500px) {
    header h3 {
      margin: 0.45rem;
    }
  }
</style>
