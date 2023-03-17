import { defineConfig, loadEnv } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import wasm from "vite-plugin-wasm";
import type * as vite from "vite";
import * as path from "path";
import sveltePreprocess from "svelte-preprocess";

const root = path.resolve(__dirname);

export default defineConfig({
  plugins: [
    svelte({
      preprocess: sveltePreprocess(),
    }),
    wasm(),
  ],
  root: path.join(root, "site"),
  envPrefix: ["BUILD_"],
  publicDir: path.join(root, "site", "public"),
  server: {
    port: 5000,
  },
  build: {
    emptyOutDir: true,
    rollupOptions: {
      output: {
        format: "esm",
        manualChunks: {
          vm: ["v86"],
          vmmisc: [],
          terminal: ["xterm", /xterm-addon-.*/],
        },
      },
      external: ["node_modules/v86/build/v86.wasm"],
    },
    target: "esnext",
  },
  // https://github.com/vitejs/vite/issues/7385#issuecomment-1286606298
  resolve: {
    alias: {
      "#/libdb.so": root,
    },
  },
});
