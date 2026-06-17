// @ts-check
import { defineConfig, fontProviders } from "astro/config";
import { satteri } from "@astrojs/markdown-satteri";
import node from "@astrojs/node";
import mdx from "@astrojs/mdx";
import svelte from "@astrojs/svelte";

const { PUBLIC_INCONSOLATA_PATH } = process.env;

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

  fonts: PUBLIC_INCONSOLATA_PATH
    ? [
        {
          provider: fontProviders.local(),
          name: "Inconsolata",
          cssVariable: "--font-inconsolata",
          fallbacks: ["monospace"],
          optimizedFallbacks: false,
          options: {
            variants: [
              {
                weight: "200 900",
                style: "normal",
                src: [`${PUBLIC_INCONSOLATA_PATH}/Inconsolata[wdth,wght].ttf`],
                featureSettings: "'dlig' on",
              },
            ],
          },
        },
      ]
    : undefined,

  integrations: [mdx(), svelte()],

  markdown: {
    processor: satteri({
      features: { superscript: true, subscript: true, smartPunctuation: true },
    }),
  },
});
