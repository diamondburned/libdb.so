<script lang="ts">
  import * as svelte from "svelte";
  import * as vm from "#/libdb.so/site/lib/vm.js";
  import favicon from "#/libdb.so/public/favicon.ico?url";
  import type * as xterm from "xterm";

  import Terminal from "#/libdb.so/site/components/Terminal/index.svelte";
</script>

<svelte:head>
  <meta name="darkreader-lock" />
  <link rel="icon" href={favicon} />
  <link rel="stylesheet" href="normalize.css" />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Inconsolata:wght@400;500;600;700;900&display=swap"
  />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Lato:wght@400;700;900&display=swap"
  />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Source+Code+Pro:wght@400;500;700;900&display=swap"
  />
</svelte:head>

<main>
  <div class="backdrop">
    <Terminal
      id="terminal"
      done={(terminal) => {
        vm.start(terminal, "/_fs.json").catch((err) => console.error(err));
      }}
    />
  </div>
</main>

<style global>
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
    background-image: url("/_fs/Pictures/background.jpg");
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
    width: min(calc(100% - clamp(6px, 5vw, 3rem)), max(80vw, 1000px));
    height: min(calc(100% - clamp(6px, 7vw, 5rem)), max(70vh, 800px));
  }

  @media (max-width: 500px) {
    #terminal {
      width: 100%;
      height: 100%;
      border-radius: 0;
    }
  }
</style>
