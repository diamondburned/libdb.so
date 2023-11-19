import { persisted } from "svelte-persisted-store";

const isReducedMotion = !!window.matchMedia(`(prefers-reduced-motion: reduce)`)
  ?.matches;

export const onekoCursor = persisted("prefs-oneko-cursor", !isReducedMotion);
export const dragWindows = persisted("prefs-move-window", true);
