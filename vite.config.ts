import { defineConfig, loadEnv } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import type * as vite from "vite";
import * as path from "path";
import viteCompression from "vite-plugin-compression";
import sveltePreprocess from "svelte-preprocess";

const root = path.resolve(__dirname);

export default defineConfig({
  // Why the FUCK is clearScreen true by default? That is fucking stupid.
  clearScreen: false,
  plugins: [
    viteCompression(),
    svelte({
      preprocess: sveltePreprocess(),
    }),
  ],
  root: path.join(root, "site"),
  publicDir: path.join(root, "dist"),
  server: {
    port: 5000,
  },
  build: {
    outDir: path.join(root, "dist"),
    emptyOutDir: true,
    rollupOptions: {
      output: {
        format: "esm",
      },
    },
    target: "esnext",
    sourcemap: true,
  },
  esbuild: {
    sourcemap: true,
  },
  // https://github.com/vitejs/vite/issues/7385#issuecomment-1286606298
  resolve: {
    alias: {
      "#/libdb.so": root,
    },
  },
});

if (import.meta.hot) {
  // always reload the page on change because v86 is fragile
  import.meta.hot.accept(() => import.meta.hot!.invalidate());
}
