import "./wasm_exec.js";
import consoleBlob from "#/libdb.so/build/vm.wasm?url";
import * as xtermpty from "xterm-pty";
import type * as xterm from "xterm";

declare global {
  interface Termios {
    iflag: number;
    oflag: number;
    cflag: number;
    lflag: number;
    cc: number[];
  }
  function vm_write_stdin(data: string): void;
  function vm_update_terminal(_: {
    row: number;
    col: number;
    xpixel: number;
    ypixel: number;
    sixel: boolean;
  }): void;
  function vm_start(): void;
  function vm_init(opts: {
    public_fs: { json: string; base: string };
    tcsets: (termios: Termios) => void;
    tcgets: () => Termios;
  }): void;
  function vm_set_public_fs(json: string, basePath: string): void;
  var console_write: null | ((fd: number, bytes: Uint8Array) => void);
}

let running: Promise<void> | null = null;

class TerminalProxy {
  private onDataDisposer: xterm.IDisposable;
  private onResizeDisposer: xterm.IDisposable;

  constructor(
    public readonly terminal: xterm.Terminal,
    public readonly pty: xtermpty.Slave
  ) {
    const lineBuffer: number[] = [];

    globalThis.console_write = (fd: number, bytes: Uint8Array) => {
      switch (fd) {
        case 1: // stdout
          this.pty.write(Array.from(bytes));
          break;
        case 2: // stderr
          this.pty.write(Array.from(bytes));
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

    this.onDataDisposer = pty.onReadable(() => this.onData(pty.read()));
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

  private onData(data: number[]) {
    const write_stdin = globalThis.vm_write_stdin;
    if (write_stdin) {
      write_stdin(String.fromCharCode(...data));
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

export async function start(
  terminal: xterm.Terminal,
  slave: xtermpty.Slave,
  publicFS: FilesystemJSON
) {
  if (running) return;

  // @ts-ignore
  const go = new globalThis.Go();
  const proxy = new TerminalProxy(terminal, slave);

  const resp = await fetch(consoleBlob);
  const module = await WebAssembly.compileStreaming(resp);
  const instance = await WebAssembly.instantiate(module, go.importObject);

  console.log("loaded wasm blob from", consoleBlob);

  console.log("starting wasm...");
  running = go.run(instance).catch((err: any) => {
    console.error("error running wasm blob", err);
  });

  console.log("initialize public httpfs");
  globalThis.vm_init({
    public_fs: {
      json: JSON.stringify(publicFS.tree),
      base: publicFS.base,
    },
    tcsets: (termios: Termios) =>
      slave.ioctl(
        "TCSETS",
        new xtermpty.Termios(
          termios.iflag,
          termios.oflag,
          termios.cflag,
          termios.lflag,
          termios.cc
        )
      ),
    tcgets: () => slave.ioctl("TCGETS"),
  });

  console.log("starting console...");
  proxy.updateQuery();
  globalThis.vm_start();

  console.log("done");
}
