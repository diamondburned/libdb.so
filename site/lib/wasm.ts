import "@lib/extern/wasm_exec.js";
import spiralWasm from "@dist/spiral.wasm?url";

declare global {
  class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<void>;
  }

  function hypnospiral_draw(
    buf: Uint8Array,
    width: number,
    height: number,
    ms: number,
    opts: {
      spinSpeed?: number;
      zoom?: number;
      blur?: number;
    },
  ): void;

  // function vm_write_stdin(data: string): void;
  // function vm_update_terminal(_: {
  //   row: number;
  //   col: number;
  //   xpixel: number;
  //   ypixel: number;
  //   sixel: boolean;
  // }): void;
  // function vm_start(): void;
  // function vm_stop(): void;
  // function vm_set_public_fs(json: string, basePath: string): void;
  // function vm_add_public_fs(json: string, basePath: string): void;
  // function vm_add_public_fs_url(url: string): void;
  //
  // var console_write: null | ((fd: number, bytes: Uint8Array) => void);
}

export async function loadSpiral() {
  await load(spiralWasm);

  const eachWait = 100; // ms
  const maxWait = 4000; // ms
  let waited = 0;
  while (waited < maxWait && !("hypnospiral_draw" in globalThis)) {
    waited += eachWait;
    await sleep(eachWait);
  }
}

const loaded: Record<string, loadResult> = {};

type loadResult = {
  go: Go;
  load: Promise<void>;
  exit: Promise<void> | null;
};

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function load(wasmURL: string) {
  if (wasmURL in loaded) {
    await loaded[wasmURL].load;
    return;
  }

  const go = new Go();
  const load = (async () => {
    const resp = await fetch(wasmURL);
    const module = await WebAssembly.compileStreaming(resp);
    const wasmInstance = await WebAssembly.instantiate(module, go.importObject);

    console.debug("loaded wasm blob from", wasmURL);
    console.debug("starting wasm...");
    const goInstance = go.run(wasmInstance).catch((err: any) => {
      console.error("error running wasm blob", err);
      throw err;
    });

    loaded[wasmURL].exit = goInstance;
  })();
  loaded[wasmURL] = { go, load, exit: null };
}
