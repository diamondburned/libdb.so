<script lang="ts">
  import { fly } from "svelte/transition";
  import { onMount } from "svelte";
  import { ToastStore } from "#/libdb.so/site/lib/toasts.js";

  import WindowControl from "#/libdb.so/site/components/WindowControl.svelte";
  import WindowClose from "#/libdb.so/site/components/Papirus/window-close.svelte";

  export let toasts: ToastStore;
  export let toastClass = "";
</script>

{#each $toasts as toast}
  <div
    class="toast {toastClass} {toast.class ?? ''}"
    transition:fly={{ y: -100, duration: 150 }}
  >
    <span>{toast.text}</span>
    <WindowControl class="dismiss" clicked={() => toasts.remove(toast)}>
      <WindowClose />
    </WindowControl>
  </div>
{/each}

<style lang="scss">
  .toast {
    background-color: #030303;
    border-radius: 10px;
    box-shadow: 0 2px 14px -4px rgba(0, 0, 0, 1);

    max-width: 300px;
    width: fit-content;

    padding: 0.5em;
    margin: 0 auto;
    margin-top: 0.5em;

    display: flex;
    gap: 0.5em;

    span {
      margin-left: 0.5em;
    }
  }
</style>
