import * as xterm from "@xterm/xterm";
import FontFaceObserver from "fontfaceobserver";
import { ImageAddon } from "@xterm/addon-image";
import { FitAddon } from "@xterm/addon-fit";
import { writable } from "svelte/store";

if (document?.fonts) {
  await document.fonts.ready;
}

export type InitOptions = Omit<
  xterm.ITerminalOptions & xterm.ITerminalInitOnlyOptions,
  "linkHandler" | "fontFamily" | "fontWeight" | "fontWeightBold"
>;

export type LinkHandlers = {
  [scheme: string]: (uri: string) => void;
};

// Load modules in parallel.
const webglModule = import("@xterm/addon-webgl");
const canvasModule = import("@xterm/addon-canvas");

export class Terminal {
  static supportsWebGL(): boolean {
    const canvas = document.createElement("canvas");
    const gl = canvas.getContext("webgl");
    return gl instanceof WebGLRenderingContext;
  }

  linkHandlers: LinkHandlers = {
    https: (uri: string) => window.open(uri, "_blank"),
    http: (uri: string) => window.open(uri, "_blank"),
    mailto: (uri: string) => window.open(uri, "_blank"),
    terminal: this.handleTerminalLink.bind(this),
  };

  title = writable("");
  xterm?: xterm.Terminal;

  private fitAddon = new FitAddon();
  private imageAddon = new ImageAddon({
    enableSizeReports: true,
    sixelSupport: true,
    sixelScrolling: true,
    sixelPaletteLimit: 4096,
    showPlaceholder: true,
  });

  private onResize_ = () => this.fit();
  private resizeObserver = new ResizeObserver(() => this.fit());

  constructor(private options: InitOptions) {}

  async open(e: HTMLElement) {
    // Wait for the fonts to load before creating the terminal.
    await Promise.all([
      new FontFaceObserver("Inconsolata").load(),
      new FontFaceObserver("Inconsolata", { weight: "500" }).load(),
      new FontFaceObserver("Inconsolata", { weight: "700" }).load(),
    ]);

    this.xterm = new xterm.Terminal({
      ...this.options,
      fontFamily: `"Inconsolata", monospace`,
      fontWeight: "500",
      fontWeightBold: "700",
      lineHeight: 1.1,
      linkHandler: {
        activate: (_: MouseEvent, uri: string) => {
          const parsed = new URL(uri);
          const scheme = parsed.protocol.replace(/:$/, "");
          if (scheme in this.linkHandlers) {
            this.linkHandlers[scheme](uri);
          }
        },
        allowNonHttpProtocols: true,
      },
      allowTransparency: true,
      customGlyphs: true,
      convertEol: true,
    });

    this.xterm.onTitleChange((t) => this.title.set(t));

    this.xterm.attachCustomKeyEventHandler((e) => {
      if (!this.xterm) {
        return true;
      }

      // Bind Ctrl + C to copy if there is a selection.
      if (e.ctrlKey && e.key == "c") {
        if (this.xterm.hasSelection()) {
          console.log("copying");
          navigator.clipboard.writeText(this.xterm.getSelection());
          return false;
        }
      }

      // Bind Ctrl + V to paste for consistency with Ctrl + C. We don't actually
      // need to do anything here because the browser will automatically paste
      // the clipboard contents into the terminal.
      if (e.ctrlKey && e.key == "v") {
        return false;
      }

      return true;
    });

    this.xterm.loadAddon(this.fitAddon);
    this.xterm.loadAddon(this.imageAddon);

    try {
      const addonModule = await webglModule;
      const addon = new addonModule.WebglAddon();
      addon.onContextLoss(() => {
        console.error("WebGL context lost, falling back to non-GL rendered terminal.");
        addon.dispose();
      });
      this.xterm.loadAddon(addon);
    } catch (err) {
      console.info("WebGL is not supported on this browser.", err);
      const addonModule = await canvasModule;
      const addon = new addonModule.CanvasAddon();
      this.xterm.loadAddon(addon);
    }

    this.xterm.open(e);
    this.resizeObserver.observe(e);
    this.fit();

    if (window) {
      window.addEventListener("resize", this.onResize_);
    }
  }

  fit() {
    this.fitAddon.fit();
    console.log("fitting");
  }

  setTheme(theme: xterm.ITheme) {
    this.options.theme = theme;
    if (this.xterm) this.xterm.options.theme = { ...theme };
  }

  dispose() {
    this.resizeObserver.disconnect();

    if (window) {
      window.removeEventListener("resize", this.onResize_);
    }

    if (this.xterm) {
      this.xterm.dispose();
      this.fitAddon.dispose();
      this.imageAddon.dispose();
    }
  }

  private handleTerminalLink(uri: string) {
    if (!this.xterm) {
      return;
    }

    const parsed = new URL(uri);
    // JS is incapable of competently parsing for the opaque part, so it gets
    // put into the pathname instead.
    switch (parsed.pathname) {
      case "write": {
        this.xterm.paste(parsed.searchParams.get("data") ?? "");
        break;
      }
      default: {
        console.log("unknown terminal command", parsed.pathname);
        return;
      }
    }
  }
}
