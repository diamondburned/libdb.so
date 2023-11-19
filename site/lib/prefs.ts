import { persisted } from "svelte-persisted-store";
import { get } from "svelte/store";

const isReducedMotion = !!window.matchMedia(`(prefers-reduced-motion: reduce)`)
  ?.matches;

export const onekoCursor = persisted("oneko-cursor", !isReducedMotion);
export const dragWindows = persisted("drag-windows", true);
export const nsfw = persisted("nsfw-v1", false);

for (const store of [onekoCursor, dragWindows]) {
  store.set(get(store));
}
