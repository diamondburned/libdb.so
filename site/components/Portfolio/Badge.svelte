<script lang="ts" context="module">
  export type Badge = {
    image: string;
    link?: string;
    alt?: string;
  };
</script>

<script lang="ts">
  import PopoverButton from "../PopoverButton.svelte";
  import OpenInNew from "#/libdb.so/site/components/MaterialIcons/open_in_new.svelte";

  export let badge: Badge;

  let hover = false;
  let popoverAnchor: HTMLElement;

  function hostname(url: string): string | null {
    try {
      const u = new URL(url);
      return u.host || null;
    } catch (_) {
      return null;
    }
  }
</script>

<div
  class="link"
  on:mouseenter={() => {
    hover = true;
  }}
  on:mouseleave={() => {
    hover = false;
  }}
>
  <PopoverButton open={hover} fly={{ y: 2 }}>
    <a href={badge.link} slot="button" role="button" class="badge popover-button">
      <img class="badge-image" bind:this={popoverAnchor} src={badge.image} alt={badge.alt} />
    </a>
    <div slot="popover">
      <span class="alt">{badge.alt ?? "88x31 badge"}</span>
      {#if badge.link && hostname(badge.link)}
        <br />
        <span class="host">{hostname(badge.link)} <OpenInNew /></span>
      {/if}
    </div>
  </PopoverButton>
</div>

<style lang="scss">
  .link {
    position: relative;

    :global(.badge-image) {
      aspect-ratio: 88 / 31;
      image-rendering: pixelated;
    }

    :global(.badge) {
      background: none;

      &:hover {
        background: none;

        filter: brightness(1.2) contrast(0.8);
        transform: translateY(-2px);
      }
    }

    :global(.popover) {
      // position: absolute;
      // bottom: 44px;

      width: max-content;
      min-width: 100px;

      padding: 0.5em 0.75em;
      translate: 0 -2px;

      font-size: 0.85em;
      line-height: 1.35;

      box-sizing: border-box;
      pointer-events: none;
    }

    .host {
      opacity: 0.55;

      :global(svg) {
        width: 0.85em;
        height: 0.85em;
        vertical-align: -0.15em;
      }
    }
  }
</style>
