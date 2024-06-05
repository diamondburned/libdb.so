<script lang="ts" context="module">
  import type { Badge as BadgeType } from "./Badge.svelte";
  export type { BadgeType as Badge };
</script>

<script lang="ts">
  import Badge from "./Badge.svelte";
  export let badges: BadgeType[];
</script>

<div class="badges">
  {#each badges as badge}
    <Badge {badge} />
  {/each}
</div>

<style lang="scss">
  .badges {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(88px, 1fr));
    grid-gap: 0.25em;
    justify-content: center;
    justify-items: center;
  }
</style>
