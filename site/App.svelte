<script lang="ts">
  import * as svelte from "svelte";
  import * as vm from "#/libdb.so/site/lib/vm.js";
  import favicon from "#/libdb.so/public/favicon.ico?url";
  import type * as xterm from "xterm";
  import { view, viewDesktop, switchView } from "#/libdb.so/site/lib/views.js";

  import Terminal from "#/libdb.so/site/components/Terminal/index.svelte";
  import Portfolio from "#/libdb.so/site/components/Portfolio/index.svelte";

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

  svelte.onMount(() => {
    // Default view is Terminal.
    switchView("terminal");
  });
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

<div class="screen">
  <div class="backdrop" />

  <div class="content">
    <div class:active={$view == "terminal"}>
      <Terminal
        id="terminal"
        done={(terminal) => {
          vm.start(terminal, "/_fs.json").catch((err) => console.error(err));
        }}
      />
    </div>
    <div class:active={$view == "portfolio"}>
      <Portfolio />
    </div>
  </div>

  <nav id="navbar">
    <div class="left">
      <button class="start" on:click={() => alert("hii!!")}>
        <img src={favicon} alt="diamondburned's eye" />
      </button>
      <div class="window-list">
        <button
          class:active={$view == "portfolio"}
          on:click={() => switchView("portfolio")}
          disabled
        >
          <img src="/_assets/papirus/system-users.svg" alt="Portfolio icon" />
          Portfolio
        </button>
        <button
          class:active={$view == "terminal"}
          on:click={() => switchView("terminal")}
        >
          <img src="/_assets/papirus/terminal.svg" alt="Terminal icon" />
          xterm.js
        </button>
      </div>
    </div>
    <div class="right">
      <span class="clock">{currentTime}</span>
      <button class="view-desktop" on:click={() => viewDesktop()} />
    </div>
  </nav>
</div>

<style global lang="scss">
  html,
  body {
    height: 100%;
  }

  body {
    font-family: "Lato", "Source Sans Pro", "Noto Sans", "Helvetica", "Segoe UI",
      sans-serif;
  }

  .monospace {
    font-family: "Inconsolata", "Noto Mono", "Source Code Pro", monospace;
  }

  div.screen {
    width: 100vw;
    height: 100%;

    /* https://www.joshwcomeau.com/gradient-generator?colors=f690dc|4e98fa&angle=55&colorMode=hcl&precision=20&easingCurve=0.25|0.75|0.75|0.25 */
    background-color: dimgray;
    background-image: url("/_fs/Pictures/background.jpg");
    background-size: cover;

    display: flex;
    flex-direction: column;

    & > .content {
      width: 100%;
      height: 100%;
      overflow: hidden;

      & > div {
        width: 100%;
        height: 100%;

        transition: all 0.2s cubic-bezier(0, 1.005, 0.165, 1);
        transform: translateY(0) scale(1);

        &:not(.active) {
          opacity: 0;
          transform: translateY(50vh) scale(0);
          pointer-events: none;
        }
      }
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
    color: white;

    .left {
      flex: 1;
    }

    .left,
    .right {
      display: flex;
      flex-direction: row;
    }

    .window-list {
      --button-width: clamp(125px, 20%, 200px);

      width: 100%;

      display: grid;
      grid-auto-flow: column;
      grid-template-columns: repeat(auto-fit, var(--button-width));
      grid-template-rows: 1fr;

      overflow: auto;

      button {
        min-width: var(--button-width);
      }
    }

    button {
      --bg-hover: rgba(255, 255, 255, 0.05);
      --bg-active: rgba(255, 255, 255, 0.1);
      --bg-active-hover: rgba(255, 255, 255, 0.15);

      border: none;
      background-color: transparent;
      font-weight: inherit;
      font-family: inherit;
      font-size: 0.9em;
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
</style>
