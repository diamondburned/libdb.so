<script lang="ts">
  import { fly } from "svelte/transition";

  export let open = false;
  export let y = 20;
  export let arrowX = 0.5;
</script>

{#if open}
  <div
    class="popover"
    transition:fly={{ y, duration: 200 }}
    style="--arrow-x: calc({arrowX} * 100%)"
  >
    <slot />
  </div>
{/if}

<style lang="scss">
  .popover {
    position: relative;
    margin-bottom: 8px;

    background: var(--adw-popover-bg-color);

    border: 1px solid rgba(0 0 0 / 14%);
    box-shadow: 0 1px 5px 1px rgba(0 0 0 / 9%), 0 2px 14px 3px rgba(0 0 0 / 5%);
    padding: 8px;
    border-radius: var(--adw-popover-radius);

    &::after {
      content: " ";
      border-top: 8px solid var(--adw-popover-bg-color);
      border-left: 8px solid transparent;
      border-right: 8px solid transparent;
      position: absolute;
      bottom: -8px;
      left: var(--arrow-x);
      transform: translateX(var(--arrow-x));
    }
  }
</style>
