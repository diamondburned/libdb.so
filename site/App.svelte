<script lang="ts">
  import * as svelte from "svelte";
  import * as vm from "#/libdb.so/site/lib/vm.js";
  import type * as xterm from "xterm";

  import Terminal from "#/libdb.so/site/components/Terminal/index.svelte";
  import publicFS from "#/libdb.so/build/publicfs.json";

  let screen: HTMLElement;

  function init(terminal: xterm.Terminal) {
    vm.start(terminal, publicFS).catch((err) => console.error(err));
  }
</script>

<main>
  <div class="backdrop">
    <Terminal id="terminal" done={init} />
  </div>
</main>

<style global>
  @import "normalize.css";
  @import url("https://fonts.googleapis.com/css2?family=Inconsolata:wght@400;500;600;700&display=swap");
  @import url("https://fonts.googleapis.com/css2?family=Lato:wght@400;700;900&display=swap");

  body {
    font-family: "Lato", "Source Sans Pro", "Noto Sans", "Helvetica", "Segoe UI",
      sans-serif;
  }

  .monospace {
    font-family: "Inconsolata", "Noto Mono", "Source Code Pro", monospace;
  }

  main {
    width: 100vw;
    height: 100vh;

    /* https://www.joshwcomeau.com/gradient-generator?colors=f690dc|4e98fa&angle=55&colorMode=hcl&precision=20&easingCurve=0.25|0.75|0.75|0.25 */
    background-color: dimgray;
    background-image: url("/Pictures/background.jpg");
    background-size: cover;
  }

  main > div {
    width: 100%;
    height: 100%;

    display: flex;
    justify-content: center;
    align-items: center;

    backdrop-filter: blur(20px);
  }

  #terminal {
    width: min(800px, calc(100% - clamp(1rem, 10vw, 3rem)));
    height: min(600px, calc(100% - clamp(1rem, 10vw, 3rem)));
  }

  #screen {
    position: absolute;
    border: 2px solid aliceblue;
    border-radius: 8px;
    padding: 4px;
  }
</style>
