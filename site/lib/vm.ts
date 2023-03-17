const RAMSize = 128 * 1024 * 1024; // 128 MB
const VGASize = 8 * 1024 * 1024; // 8 MB

export async function spawn() {
  // @ts-ignore
  const v86 = await import("v86");
  // @ts-ignore
  const v86wasm = await import("v86/build/v86.wasm");
  const v86bios = await import("v86/bios/seabios.bin?url");

  const vm = v86.V86Starter({
    // TODO: swap this out for a wasm loader
    wasm_fn: v86wasm,
    memory_size: RAMSize,
    vga_memory_size: VGASize,
    autostart: true,
  });
}
