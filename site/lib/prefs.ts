import { persisted } from "svelte-persisted-store";
import { derived, get, writable } from "svelte/store";

const isReducedMotion = !!window.matchMedia(`(prefers-reduced-motion: reduce)`)?.matches;

export const onekoCursor = persisted("oneko-cursor", !isReducedMotion);
export const dragWindows = persisted("drag-windows", true);
export const nsfw = persisted("nsfw-v1", false);
export const theme = persisted("theme", "Adwaita");
export const prefersDark = persisted<boolean | null>("prefers-dark", null);

for (const store of [onekoCursor, dragWindows]) {
  store.set(get(store));
}

const systemIsDark = writable(false);
if (window.matchMedia) {
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)");
  prefersDark.addEventListener("change", (event) => systemIsDark.set(event.matches));
  systemIsDark.set(prefersDark.matches);
}

export const isDark = derived([prefersDark, systemIsDark], ([prefersDark, systemIsDark]) => {
  return prefersDark !== null ? prefersDark : systemIsDark;
});
