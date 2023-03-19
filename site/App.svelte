<script lang="ts">
  import "./style.scss";

  import * as svelte from "svelte";
  import * as console_ from "#/libdb.so/site/lib/console.js";
  import type * as xterm from "xterm";

  import Terminal from "#/libdb.so/site/components/Terminal/index.svelte";

  let screen: HTMLElement;

  function init(terminal: xterm.Terminal) {
    console_.start(terminal).catch((err) => console.error(err));
  }
</script>

<main>
  <div class="backdrop">
    <Terminal id="terminal" done={init} />
  </div>
</main>

<style global>
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
