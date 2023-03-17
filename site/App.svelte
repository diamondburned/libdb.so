<script lang="ts">
  import "./style.scss";

  import * as xterm from "xterm";
  import * as svelte from "svelte";

  import Terminal from "#/libdb.so/site/components/Terminal.svelte";

  let terminal: xterm.Terminal;

  svelte.onMount(async () => {
    const vm = await import("#/libdb.so/site/lib/vm.js");
    vm.spawn();
  });
</script>

<main>
  <Terminal id="terminal" bind:terminal />
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
    width: calc(100% - 2rem);
    height: calc(100% - 2rem);
  }
</style>
