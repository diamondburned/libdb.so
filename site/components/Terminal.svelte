<script lang="ts">
  import "xterm/css/xterm.css";

  import { WebglAddon as WebGLAddon } from "xterm-addon-webgl";
  import { FitAddon } from "xterm-addon-fit";
  import * as svelte from "svelte";
  import type * as xterm from "xterm";

  let terminalElement: HTMLElement;

  export let id: string;
  export let terminal: xterm.Terminal;

  svelte.onMount(async () => {
    const xterm = await import("xterm");

    terminal = new xterm.Terminal({
      fontFamily: "monospace",
      fontWeight: "500",
      fontWeightBold: "700",
      allowTransparency: true,
    });

    const webglAddon = new WebGLAddon();
    webglAddon.onContextLoss(() => webglAddon.dispose());

    const fitAddon = new FitAddon();
    const onResize = () => fitAddon.fit();

    // terminal.loadAddon(webglAddon);
    terminal.loadAddon(fitAddon);
    terminal.open(terminalElement);
    terminal.write("Initializing VM...\r\n");

    onResize();
    window.addEventListener("resize", onResize);

    return () => {
      terminal.dispose();
      webglAddon.dispose();
      fitAddon.dispose();
      window.removeEventListener("resize", onResize);
    };
  });
</script>

<div {id} bind:this={terminalElement} />

<style>
  div {
    overflow: hidden;

    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
