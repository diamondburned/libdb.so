import * as xterm from "xterm";
import { ImageAddon, IImageAddonOptions } from "xterm-addon-image";
import { FitAddon } from "xterm-addon-fit";

if (document?.fonts) {
  await document.fonts.ready;
}

export type InitOptions = Omit<
  xterm.ITerminalOptions & xterm.ITerminalInitOnlyOptions,
  "linkHandler"
>;

export type LinkHandlers = {
  [scheme: string]: (uri: string) => void;
};

export class Terminal extends xterm.Terminal {
  linkHandlers: LinkHandlers = {
    https: (uri: string) => window.open(uri, "_blank"),
    http: (uri: string) => window.open(uri, "_blank"),
    mailto: (uri: string) => window.open(uri, "_blank"),
    terminal: this.handleTerminalLink.bind(this),
  };

  private fitAddon = new FitAddon();
  private imageAddon = new ImageAddon({
    enableSizeReports: true,
    sixelSupport: true,
    sixelScrolling: true,
    sixelPaletteLimit: 4096,
    showPlaceholder: true,
  });

  private onResize_ = () => this.fitAddon.fit();

  constructor(options: InitOptions) {
    super({
      ...options,
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
      convertEol: true,
    });

    this.attachCustomKeyEventHandler((e) => {
      // Bind Ctrl + C to copy if there is a selection.
      if (e.ctrlKey && e.key == "c") {
        if (this.hasSelection()) {
          console.log("copying");
          navigator.clipboard.writeText(this.getSelection());
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

    this.loadAddon(this.fitAddon);
    this.loadAddon(this.imageAddon);
  }

  open(e: HTMLElement) {
    super.open(e);
    this.fitAddon.fit();

    if (window) {
      window.addEventListener("resize", this.onResize_);
    }
  }

  fit() {
    this.fitAddon.fit();
    console.log("fitting");
  }

  dispose() {
    super.dispose();
    this.fitAddon.dispose();
    this.imageAddon.dispose();

    if (window) {
      window.removeEventListener("resize", this.onResize_);
    }
  }

  private handleTerminalLink(uri: string) {
    const parsed = new URL(uri);
    // JS is incapable of competently parsing for the opaque part, so it gets
    // put into the pathname instead.
    switch (parsed.pathname) {
      case "write": {
        this.paste(parsed.searchParams.get("data") ?? "");
        break;
      }
      default: {
        console.log("unknown terminal command", parsed.pathname);
        return;
      }
    }
  }
}
