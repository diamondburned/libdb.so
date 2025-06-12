import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import * as path from "path";
import viteCompression from "vite-plugin-compression";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

const root = path.resolve(__dirname);

export default defineConfig({
  // Why the FUCK is clearScreen true by default? That is fucking stupid.
  clearScreen: false,
  plugins: [
    viteCompression({
      filter: /\.(js|mjs|json|css|html|wasm)$/i,
    }),
    svelte({
      preprocess: vitePreprocess(),
      onwarn: (warning, handler) => {
        const ignoredCodes = ["css-unused-selector", "a11y/no-noninteractive-element-interactions"];
        if (ignoredCodes.find((code) => code == warning.code)) return;
        if (handler) handler(warning);
      },
    }),
  ],
  root: path.join(root, "site"),
  publicDir: path.join(root, "build", "public"),
  server: {
    host: "0.0.0.0",
    port: 5001,
    allowedHosts: true,
  },
  build: {
    outDir: path.join(root, "build", "dist"),
    assetsDir: "_assets",
    emptyOutDir: true,
    rollupOptions: {
      output: {
        format: "esm",
      },
    },
    target: "esnext",
    sourcemap: true,
    reportCompressedSize: false, // viteCompression does this for us
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
