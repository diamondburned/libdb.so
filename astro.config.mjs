// @ts-check
import { defineConfig, fontProviders } from "astro/config";
import { satteri } from "@astrojs/markdown-satteri";
import node from "@astrojs/node";
import mdx from "@astrojs/mdx";

// https://astro.build/config
export default defineConfig({
  site: "https://libdb.so",
  srcDir: "site",
  adapter: node({ mode: "standalone" }),

  vite: {
    clearScreen: false,
    css: {
      transformer: "lightningcss",
    },
  },

  fonts: [
    {
      provider: fontProviders.google(),
      name: "Inconsolata",
      cssVariable: "--font-inconsolata",
      fallbacks: ["monospace"],
      optimizedFallbacks: false,
      // https://docs.astro.build/en/guides/fonts/#using-variable-fonts
      styles: ["normal"],
      weights: ["200 900"],
    },
  ],

  integrations: [mdx()],

  markdown: {
    processor: satteri({
      features: { superscript: true, subscript: true, smartPunctuation: true },
    }),
  },
});
