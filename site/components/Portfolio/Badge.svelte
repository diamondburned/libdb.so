<script lang="ts" context="module">
  export type Badge = {
    image: string;
    link?: string;
    alt?: string;
  };
</script>

<script lang="ts">
  import { fly } from "svelte/transition";

  import Popover from "../Popover.svelte";

  export let badge: Badge;
  let hover = false;

  function hostname(url: string): string | null {
    try {
      const u = new URL(url);
      return u.host || null;
    } catch (_) {
      return null;
    }
  }
</script>

<div class="link">
  <a href={badge.link} on:mouseenter={() => (hover = true)} on:mouseleave={() => (hover = false)}>
    <img class="badge" src={badge.image} alt={badge.alt} />
  </a>
  <div class="tooltip">
    <Popover open={hover} arrowX={0.15}>
      <b>88x31</b>
      {#if badge.alt}
        <br />
        <span class="alt">{badge.alt}</span>
      {/if}
      {#if badge.link && hostname(badge.link) && hostname(badge.link) != badge.alt}
        <br />
        <span class="host">{hostname(badge.link)}</span>
      {/if}
    </Popover>
  </div>
</div>

<style lang="scss">
  .badge {
    aspect-ratio: 88 / 31;
    image-rendering: pixelated;
  }

  .link {
    position: relative;

    .tooltip {
      position: absolute;
      z-index: 10;
      bottom: 42px;

      font-size: 0.85em;
      line-height: 1.35;

      box-sizing: border-box;
      pointer-events: none;

      span {
        opacity: 0.75;
      }
    }
  }
</style>
