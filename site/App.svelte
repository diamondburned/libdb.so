<script lang="ts">
  import * as svelte from "svelte";
  import * as vm from "#/libdb.so/site/lib/vm.js";
  import nsfw from "#/libdb.so/internal/nsfw/nsfw.js";
  import favicon from "#/libdb.so/public/favicon.ico?url";
  import normalizeCSS from "normalize.css/normalize.css?url";
  import type * as xterm from "xterm";
  import {
    activeViews,
    toggleView,
    focusedView,
    toggleShowDesktop,
  } from "#/libdb.so/site/lib/views.js";

  import Terminal from "#/libdb.so/site/components/Terminal/index.svelte";
  import Portfolio from "#/libdb.so/site/components/Portfolio/index.svelte";
  import VisibilityIcon from "@material-design-icons/svg/filled/visibility.svg?raw";
  import VisibilityOffIcon from "@material-design-icons/svg/filled/visibility_off.svg?raw";

  import "libwebring/dist/webring.css";
  import "libwebring/dist/webring-element.js";

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

  const fonts = ["Inconsolata", "Lato", "Nunito", "Source Code Pro"];

  let visibilityIcon = VisibilityIcon;
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <link rel="stylesheet" href={normalizeCSS} />
  {#each fonts as font}
    <link
      rel="stylesheet"
      href={"https://fonts.googleapis.com/css2?family=" +
        font.replaceAll(" ", "+") +
        ":wght@400;500;600;700;900&display=swap"}
    />
  {/each}
  <meta name="description" content="Hi, I'm Diamond!" />
  <meta name="author" content="diamondburned" />
</svelte:head>

<div class="screen">
  <div class="backdrop" />

  <div class="content">
    <Terminal
      done={(terminal) => {
        vm.start(terminal, "/_fs.json").catch((err) => console.error(err));
      }}
    />
    <Portfolio />
  </div>

  <nav id="navbar">
    <div class="left">
      <button class="start" on:click={() => alert("hii!!")}>
        <img src={favicon} alt="diamondburned's eye" />
      </button>
      <div class="window-list">
        <button
          class:active={$focusedView == "portfolio"}
          on:click={() => toggleView("portfolio")}
        >
          <img src="/_assets/papirus/system-users.svg" alt="Portfolio icon" />
          About
        </button>
        <button
          class:active={$focusedView == "terminal"}
          on:click={() => toggleView("terminal")}
        >
          <img src="/_assets/papirus/terminal.svg" alt="Terminal icon" />
          xterm.js
        </button>
      </div>
    </div>
    <div class="right">
      {#if $nsfw}
        <button
          class="icon toggle-nsfw"
          title="Disable NSFW"
          on:click={() => ($nsfw = false)}
          on:mouseenter={() => (visibilityIcon = VisibilityOffIcon)}
          on:mouseleave={() => (visibilityIcon = VisibilityIcon)}
        >
          {@html visibilityIcon}
        </button>
      {/if}
      <span class="clock">{currentTime}</span>
      <button class="view-desktop" on:click={() => toggleShowDesktop()} />
    </div>
  </nav>
</div>

<style global lang="scss">
  :root {
    --blue: rgba(85, 205, 252, 1);
    --pink: rgba(247, 168, 184, 1);

    --blue-rgb: 85, 205, 252;
    --pink-rgb: 247, 168, 184;

    --glow-radius: 0.35em;
    --glow-alpha: 0.65;
    --pink-glow: 0 0 var(--glow-radius) rgba(var(--pink-rgb), var(--glow-alpha));
    --blue-glow: 0 0 var(--glow-radius) rgba(var(--blue-rgb), var(--glow-alpha));
  }

  .text-pink {
    color: var(--pink);
  }

  .text-blue {
    color: var(--blue);
  }

  .text-pink-glow {
    @extend .text-pink;
    text-shadow: var(--pink-glow);
  }

  .text-blue-glow {
    @extend .text-blue;
    text-shadow: var(--blue-glow);
  }

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
      position: relative;
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

      border-top: 2px solid transparent;
      border-bottom: 2px solid transparent;

      transition: all 0.075s ease-in-out;

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

      img,
      :global(svg) {
        width: 1.5em;
        height: 1.5em;
      }
    }

    .window-list button {
      min-width: var(--button-width);
      padding-left: 0;

      img {
        margin: 0 0.5em;
      }
    }

    .toggle-nsfw {
      :global(svg) {
        --color: var(--pink-rgb);

        fill: rgb(var(--color));
        width: 1.5em;
        height: 1.5em;
        filter: drop-shadow(
          0 0 var(--glow-radius) rgba(var(--color), var(--glow-alpha))
        );
        opacity: 0.75;
      }

      &:hover :global(svg) {
        --color: var(--blue-rgb);
        opacity: 1;
      }
    }

    .clock {
      height: 100%;
      display: flex;
      align-items: center;
      padding: 0 0.5rem;
    }

    .view-desktop {
      padding: 0;
      padding-left: 0.5em;
      border-left: 1px solid rgba(255, 255, 255, 0.25);

      &:hover {
        background-color: var(--bg-hover);
      }
    }
  }
</style>
