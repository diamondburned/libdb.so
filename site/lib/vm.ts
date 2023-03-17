import v86wasm from "v86/build/v86.wasm?url";
import v86bios from "v86/bios/seabios.bin?url";
import v86vgabios from "v86/bios/vgabios.bin?url";
import nixosImage from "#/libdb.so/nixos/result/nixos.img?url";

import type * as xterm from "xterm";

const RAMSize = 128 * 1024 * 1024; // 128 MB
const VGASize = 8 * 1024 * 1024; // 8 MB

export class VM {
  constructor() {}

  async start(
    terminal: xterm.Terminal,
    screen?: HTMLElement
  ): Promise<() => void> {
    // @ts-ignore
    const v86 = await import("v86");

    // Since v86 cannot competently get the size of the disk image, we'll do
    // it ourselves.
    const imagesz = await contentLength(nixosImage);

    const vm = new v86.V86Starter({
      wasm_fn: async (env: any) => {
        // https://web.dev/loading-wasm/
        const resp = await fetch(v86wasm);
        const module = await WebAssembly.compileStreaming(resp);
        const instance = await WebAssembly.instantiate(module, env);
        console.log("loaded wasm blob from", v86wasm);
        return instance.exports;
      },
      memory_size: RAMSize,
      screen_container: screen,
      vga_memory_size: VGASize,
      disable_keyboard: true,
      disable_mouse: true,
      autostart: true,
      acpi: false,
      bios: {
        url: v86bios,
      },
      vga_bios: {
        url: v86vgabios,
      },
      // TODO: hda
      hda: {
        url: nixosImage,
        size: imagesz,
        async: false, // async is buggy as fuck
      },
      serial_adapter: {
        show: () => {
          terminal.write("Ready VM. Booting...\r\n");
        },
      },
    });

    const onDownloadProgress = (() => {
      let last = 0;
      const sizeInterval = 50 * 1024 * 1024; // 50 MB

      return (progress: any) => {
        const basename = progress.file_name.split("/").pop();
        const { loaded, total } = progress;

        if (loaded - last > sizeInterval) {
          terminal.write(`Downloading ${basename}: ${loaded}/${total}\r\n`);
          last = loaded;
        } else if (loaded === total) {
          terminal.write(`Downloaded ${basename}\r\n`);
        }

        console.debug("download progress", progress);
      };
    })();

    const onDownloadError = (progress: any) => {
      const basename = progress.file_name.split("/").pop();
      terminal.write(`Download ${basename} failed\r\n`);

      console.error("download error", progress);
    };

    vm.add_listener("download-progress", onDownloadProgress);
    vm.add_listener("download-error", onDownloadError);
    vm.add_listener("serial0-output-char", terminal.write);

    const onData = terminal.onData((data) => vm.serial0_send);

    return async () => {
      onData.dispose();

      vm.remove_listener("serial0-output-char", terminal.write);
      vm.remove_listener("download-progress", onDownloadProgress);
      vm.remove_listener("download-error", onDownloadError);

      await vm.stop();
    };
  }
}

async function contentLength(url: string): Promise<number> {
  const resp = await fetch(url, { method: "HEAD" });
  const contentLength = resp.headers.get("Content-Length");
  if (!contentLength) {
    throw new Error(`no content length for ${url}`);
  }
  return parseInt(contentLength);
}
