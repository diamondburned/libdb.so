<script lang="ts">
  import * as svelte from "svelte";
  import * as vm from "#/libdb.so/site/lib/vm.js";
  import favicon from "#/libdb.so/public/favicon.ico?url";
  import type * as xterm from "xterm";

  import Terminal from "#/libdb.so/site/components/Terminal/index.svelte";
  import Portfolio from "#/libdb.so/site/components/Portfolio/index.svelte";

  type View = null | "terminal" | "portfolio";

  let currentView: View = "terminal";
  let lastView: View = null;
  function switchView(view: View) {
    currentView = currentView == view ? null : view;
    lastView = null;
  }
  function viewDesktop() {
    if (!lastView) lastView = currentView;
    currentView = currentView == null ? lastView : null;
  }

  let currentTime = "00:00";
  function updateTime() {
    currentTime = new Date().toLocaleTimeString(undefined, {
      hour: "2-digit",
      minute: "2-digit",
    });
  }
  updateTime();

  const updateTimer = setInterval(() => updateTime(), 5000);
  svelte.onDestroy(() => clearInterval(updateTimer));
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
  <div class="backdrop" />

  <div class="content">
    <div class="centered" class:active={currentView == "terminal"}>
      <Terminal
        id="terminal"
        done={(terminal) => {
          vm.start(terminal, "/_fs.json").catch((err) => console.error(err));
        }}
      />
    </div>
    <div class:active={currentView == "portfolio"}>
      <Portfolio />
    </div>
  </div>

  <nav id="navbar">
    <div class="left">
      <button class="start" on:click={() => alert("hii!!")}>
        <img src="/favicon.ico" alt="diamondburned's eye" />
      </button>
      <div class="window-list">
        <button
          class:active={currentView == "portfolio"}
          on:click={() => switchView("portfolio")}
          disabled
        >
          <img src="/_assets/papirus/system-users.svg" alt="Portfolio icon" />
          Portfolio
        </button>
        <button
          class:active={currentView == "terminal"}
          on:click={() => switchView("terminal")}
        >
          <img src="/_assets/papirus/terminal.svg" alt="Terminal icon" />
          Terminal
        </button>
      </div>
    </div>
    <div class="right">
      <span class="clock">{currentTime}</span>
      <button class="view-desktop" on:click={() => viewDesktop()} />
    </div>
  </nav>
</main>

<style global lang="scss">
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

    display: flex;
    flex-direction: column;

    .content {
      width: 100%;
      height: 100%;
      overflow: auto;

      & > div {
        transition: opacity 0.075s ease-in-out;

        &:not(.active) {
          opacity: 0;
          pointer-events: none;
        }
      }
    }

    .centered {
      width: 100%;
      height: 100%;

      display: flex;
      justify-content: center;
      align-items: center;
    }

    & > * {
      z-index: 1;
    }

    & > .backdrop {
      width: 100%;
      height: 100%;
      z-index: 0;
      position: absolute;
      backdrop-filter: blur(20px);
    }
  }

  #navbar {
    background-color: rgba(0, 0, 0, 0.5);
    background-image: linear-gradient(
      to bottom,
      rgba(0, 0, 0, 0.25),
      rgba(0, 0, 0, 0.65)
    );

    display: flex;
    justify-content: space-between;
    gap: 0.5rem;

    font-weight: 700;
    user-select: none;

    .left {
      width: 100%;
    }

    .left,
    .right {
      display: flex;
      flex-direction: row;
    }

    .window-list {
      width: 100%;

      display: grid;
      grid-template-columns: repeat(auto-fit, clamp(100px, 20%, 200px));
      grid-template-rows: 1fr;

      overflow-x: auto;
    }

    button {
      --bg-hover: rgba(255, 255, 255, 0.05);
      --bg-active: rgba(255, 255, 255, 0.1);
      --bg-active-hover: rgba(255, 255, 255, 0.15);

      border: none;
      background-color: transparent;
      font-weight: inherit;
      color: white;

      display: flex;
      align-items: center;

      padding: 0.35em 0.5em;
      padding-left: 0;

      border-top: 2px solid transparent;
      border-bottom: 2px solid transparent;

      &:hover:not(:disabled) {
        background-color: var(--bg-hover);
      }

      &:disabled {
        opacity: 0.5;
      }

      &.active {
        background-color: var(--bg-active);
        border-bottom: 2px solid white;

        &:hover {
          background-color: var(--bg-active-hover);
        }
      }

      img {
        width: 1.5em;
        height: 1.5em;
        margin: 0 0.5em;
      }
    }

    .start {
      padding-right: 0;
    }

    .clock {
      height: 100%;
      display: flex;
      align-items: center;
      padding: 0 0.5rem;
    }

    .view-desktop {
      width: 0.35em;
      border-left: 1px solid rgba(255, 255, 255, 0.25);

      &:hover {
        background-color: var(--bg-hover);
      }
    }
  }

  #terminal {
    width: min(calc(100% - clamp(6px, 5vw, 3rem)), max(80vw, 1000px));
    height: min(calc(100% - clamp(6px, 7vw, 5rem)), max(70vh, 800px));

    @media (max-width: 500px) {
      width: 100%;
      height: 100%;
      border-radius: 0;
    }
  }
</style>
