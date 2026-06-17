<script lang="ts">
  import { loadSpiral } from "@lib/wasm";
  import { shouldAnimate } from "@lib/motion.svelte";
  import { onMount, onDestroy } from "svelte";

  const asciiTable = [" ", ".", ":", "*", "#", "@"];
  const fps = 7;

  let {
    width,
    height,
    spinSpeed = 3.6,
    zoom = 1.2,
    blur = 0.4,
  }: {
    width: number;
    height: number;
    spinSpeed: number;
    zoom: number;
    blur: number;
  } = $props();

  let pixBuffer = $derived(new Uint8Array(width * height * 2));
  let txtBuffer = $derived(new Array((width + 1) * height).fill("\n"));

  function calcIdx(x: number, y: number, stride: number) {
    return y * stride + x;
  }

  let element: HTMLPreElement;
  let drawing = $state(false);
  let startTime = 0;

  function draw(ts: number) {
    if (!startTime) {
      startTime = ts;
    }

    const ms = ts - startTime;
    window.hypnospiral_draw(pixBuffer, width, height * 2, ms, { spinSpeed, zoom, blur });

    for (let y = 0; y < height; y++) {
      for (let x = 0; x < width; x++) {
        const v1 = pixBuffer[calcIdx(x, y * 2 + 0, width)];
        const v2 = pixBuffer[calcIdx(x, y * 2 + 1, width)];
        const v = (v1 / 255 + v2 / 255) / 2;

        const cx = Math.floor(v * (asciiTable.length - 1));
        const ch = asciiTable[cx];

        txtBuffer[calcIdx(x, y, width + 1)] = ch;
      }
    }

    element.textContent = txtBuffer.join("");

    setTimeout(() => {
      if (drawing && $shouldAnimate) requestAnimationFrame(draw);
    }, 1000 / fps);
  }

  onMount(async () => {
    await loadSpiral();
    drawing = true;
  });

  onDestroy(() => {
    drawing = false;
  });

  $effect(() => {
    if (drawing && $shouldAnimate) {
      startTime = 0;
      requestAnimationFrame(draw);
    }
  });
</script>

<pre class="spiral ascii-art" bind:this={element}></pre>

<style>
  .spiral {
    margin: 0;
    white-space: pre;
  }
</style>
