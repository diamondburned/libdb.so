import type { Readable, Writable } from "svelte/store";

declare global {
  interface Window {
    animate: Writable<boolean>;
    shouldAnimate: Readable<boolean>;
  }
}
