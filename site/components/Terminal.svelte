<script lang="ts">
  import { WebglAddon as WebGLAddon } from "xterm-addon-webgl";
  import { FitAddon } from "xterm-addon-fit";
  import * as xterm from "xterm";
  import * as svelte from "svelte";

  let terminalElement: HTMLElement;

  export let id: string;
  export let terminal = new xterm.Terminal({
    fontFamily: "monospace",
    allowTransparency: true,
  });

  svelte.onMount(() => {
    const webglAddon = new WebGLAddon();
    webglAddon.onContextLoss(() => webglAddon.dispose());

    const fitAddon = new FitAddon();
    const onResize = () => fitAddon.fit();

    // terminal.loadAddon(webglAddon);
    terminal.loadAddon(fitAddon);
    terminal.open(terminalElement);
    terminal.write("Initializing...");

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
