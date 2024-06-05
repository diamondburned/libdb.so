<script lang="ts" context="module">
  export type Badge = {
    image: string;
    link?: string;
    alt?: string;
  };
</script>

<script lang="ts">
  import { fly } from "svelte/transition";

  export let badge: Badge;
  let hover = false;
</script>

<a
  class="link"
  href={badge.link}
  on:mouseenter={() => (hover = true)}
  on:mouseleave={() => (hover = false)}
>
  <img class="badge" src={badge.image} alt={badge.alt} />
  {#if badge.alt && hover}
    <div class="popover tooltip" transition:fly={{ y: 20, duration: 250 }}>{badge.alt}</div>
  {/if}
</a>

<style lang="scss">
  .badge {
    aspect-ratio: 88 / 31;
    image-rendering: pixelated;
  }

  .link {
    position: relative;

    & .tooltip {
      position: absolute;
      z-index: 10;
      top: -2.5em;
      width: 100%;
      padding: 0.35em;

      font-size: 0.85em;
      background: var(--adw-popover-bg-color);
      border: 1px solid var(--adw-popover-shade-color);
      color: #fff;

      box-sizing: border-box;
      text-align: center;
      pointer-events: none;
    }
  }
</style>
