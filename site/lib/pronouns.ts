import { nsfw } from "#/libdb.so/site/lib/prefs.js";
import { derived } from "svelte/store";

export const pronouns = derived(nsfw, (nsfw) => (nsfw ? "she/her/hers" : "it/its/its"));

export const selectPronouns = (content: string, option: string) => {
  const pronoun = content.toLowerCase();
  const replace = {
    "it/its/its": {
      it2: "it",
      its3: "its",
      she: "it",
      her: "its",
      hers: "its",
      herself: "itself",
      "she's": "it's",
      "she/her": "it/its",
    }[pronoun],
    "she/her/hers": {
      it: "she",
      it2: "her",
      its: "her",
      its3: "hers",
      itself: "herself",
      "it's": "she's",
      "it/its": "she/her",
    }[pronoun],
  }[option];
  return replace ? preserveCase(content, replace) : content;
};

export function preserveCase(original: string, replace: string): string {
  // First-letter capitalization.
  if (original[0].toUpperCase() == original[0]) {
    return replace[0].toUpperCase() + replace.slice(1);
  }
  // All caps.
  if (original == original.toUpperCase()) {
    return replace.toUpperCase();
  }
  // Lowercase.
  return replace;
}
