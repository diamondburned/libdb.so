import { writable, derived } from "svelte/store";

const reducedMotionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");

export const animate = writable(true);

export const shouldAnimate = derived(
  animate,
  ($animate, set) => {
    const updateMotionPreference = () => set($animate && !reducedMotionQuery.matches);
    updateMotionPreference();

    reducedMotionQuery.addEventListener("change", updateMotionPreference);
    return () => reducedMotionQuery.removeEventListener("change", updateMotionPreference);
  },
  reducedMotionQuery.matches,
);

// Export these two stores globally:
window.animate = animate;
window.shouldAnimate = shouldAnimate;
