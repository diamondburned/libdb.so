import "./wasm_exec.js";
import consoleBlob from "#/libdb.so/build/vm.wasm?url";
import type * as xterm from "xterm";

declare global {
  function vm_write_stdin(data: string): void;
  function vm_update_terminal(_: {
    row: number;
    col: number;
    xpixel: number;
    ypixel: number;
    sixel: boolean;
  }): void;
  function vm_start(): void;
  function vm_set_public_fs(json: string, basePath: string): void;
  var console_write: null | ((fd: number, bytes: Uint8Array) => void);
}

let running: Promise<void> | null = null;

class TerminalProxy {
  private onDataDisposer: xterm.IDisposable;
  private onResizeDisposer: xterm.IDisposable;

  constructor(public readonly terminal: xterm.Terminal) {
    const lineBuffer: number[] = [];

    globalThis.console_write = (fd: number, bytes: Uint8Array) => {
      switch (fd) {
        case 1: // stdout
          this.terminal.write(bytes);
          break;
        case 2: // stderr
          this.terminal.write(bytes);
          while (true) {
            const index = lineBuffer.indexOf("\n".charCodeAt(0));
            if (index === -1) {
              break;
            }
            console.log("vm:", fd, bytes);
            lineBuffer.splice(0, lineBuffer.length);
          }
          break;
        default:
          console.log("unknown fd", fd, bytes);
      }
    };

    this.onDataDisposer = this.terminal.onData(this.onData.bind(this));
    this.onResizeDisposer = this.terminal.onResize(this.onResize.bind(this));
  }

  reset() {
    globalThis.console_write = null;
    this.onDataDisposer.dispose();
    this.onResizeDisposer.dispose();
  }

  updateQuery() {
    this.onResize(this.terminal);
  }

  private onData(data: string) {
    const write_stdin = globalThis.vm_write_stdin;
    if (write_stdin) {
      write_stdin(data);
    } else {
      console.log("write_stdin is not ready yet");
    }
  }

  private onResize(termsz: { rows: number; cols: number }) {
    if (globalThis.vm_update_terminal) {
      globalThis.vm_update_terminal({
        row: termsz.rows,
        col: termsz.cols,
        xpixel: 0,
        ypixel: 0,
        sixel: true,
      });
    } else {
      console.log("update_terminal is not ready yet");
    }
  }
}

type FileTreeKey =
  | `${string}/` // directory
  | `${string}`; // file

interface FileTree {
  [key: FileTreeKey]: FileTree | { size: number };
}

type FilesystemJSON = {
  base: string;
  tree: FileTree;
};

export async function start(terminal: xterm.Terminal, publicFSURL: string) {
  if (running) return;

  // @ts-ignore
  const go = new globalThis.Go();
  const proxy = new TerminalProxy(terminal);

  try {
    const resp = await fetch(consoleBlob);
    const module = await WebAssembly.compileStreaming(resp);
    const instance = await WebAssembly.instantiate(module, go.importObject);

    console.log("loaded wasm blob from", consoleBlob);

    console.log("starting wasm...");
    running = go.run(instance).catch((err: any) => {
      console.error("error running wasm blob", err);
    });

    console.log("initialize public httpfs");
    const publicFS = await fetch(publicFSURL).then((r) => r.json());
    globalThis.vm_set_public_fs(JSON.stringify(publicFS.tree), publicFS.base);

    console.log("starting console...");
    proxy.updateQuery();
    globalThis.vm_start();

    console.log("done");
  } catch (err) {
    console.error("error starting vm", err);
    terminal.write(`Error starting VM: ${err}\r\n`);
  }
}
