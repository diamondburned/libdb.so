<script lang="ts">
  import "./style.scss";

  import * as svelte from "svelte";
  import * as console_ from "#/libdb.so/site/lib/console.js";
  import type * as xterm from "xterm";

  import Terminal from "#/libdb.so/site/components/Terminal.svelte";

  let screen: HTMLElement;

  function init(terminal: xterm.Terminal) {
    console_.start(terminal).catch((err) => console.error(err));
  }
</script>

<main>
  <Terminal id="terminal" done={init} />
  <div id="screen" style="display: none" bind:this={screen}>
    <div style="white-space: pre; font: 14px monospace; line-height: 14px" />
    <canvas style="display: none" />
  </div>
</main>

<style global>
  main {
    width: 100vw;
    height: 100vh;

    display: flex;
    justify-content: center;
    align-items: center;
  }

  #terminal {
    width: calc(100% - 1rem);
    height: calc(100% - 2rem);
  }

  #screen {
    position: absolute;
    border: 2px solid aliceblue;
    border-radius: 8px;
    padding: 4px;
  }
</style>
