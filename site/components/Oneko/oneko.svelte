<script lang="ts">
  import * as svelte from "svelte";
  import {
    Window,
    focusedView,
    viewWindows,
  } from "#/libdb.so/site/lib/views.js";
  import { speed, spriteSets } from "./oneko.ts";
  import { fade } from "svelte/transition";

  let x = 16;
  let y = 16;

  let spriteX = 0;
  let spriteY = 0;

  let mouseX = 0;
  let mouseY = 0;

  let screenWidth = 0;
  let screenHeight = 0;

  let frameCount = 0;
  let idleTime = 0;
  let idleAnimation: string | null = null;
  let idleAnimationFrame = 0;
  let idleLastX = 0;
  let idleLastY = 0;

  export let windows: Window[] = [];

  function windowToRect(window: Window) {
    return {
      top: window.y,
      left: window.x,
      bottom: window.y + window.height,
      right: window.x + window.width,
    };
  }

  const mouseSize = 16;

  $: activeWindow = $viewWindows[$focusedView];
  $: allWindows = [
    // Main screen
    {
      top: 0,
      left: 0,
      bottom: screenHeight,
      right: screenWidth,
    },
    // Mouse cursor
    {
      top: mouseY - mouseSize / 2,
      left: mouseX - mouseSize / 2,
      bottom: mouseY + mouseSize / 2,
      right: mouseX + mouseSize / 2,
    },
    // Active window
    windowToRect(activeWindow),
    // Other windows, usually the navigation bar
    ...windows.map(windowToRect),
  ] as Rectangle[];

  function setSprite(name: string, frame: number) {
    const sprites = spriteSets[name];
    const sprite = sprites[frame % sprites.length];
    spriteX = sprite[0] * 32;
    spriteY = sprite[1] * 32;
  }

  function randomIntn(n: number) {
    return Math.floor(Math.random() * n);
  }

  function clamp(min: number, val: number, max: number) {
    return Math.min(Math.max(val, min), max);
  }

  function within(min: number, val: number, max: number) {
    return val >= min && val <= max;
  }

  function idle() {
    idleTime++;

    // every ~ 20 seconds
    if (idleTime > 10 && Math.random() < 0.05 && !idleAnimation) {
      const edgeThreshold = 16;

      let idleAnimations = ["sleeping", "scratchSelf"];
      for (const window of allWindows) {
        if (within(-edgeThreshold, x - window.right, 0)) {
          // inside window left edge or outside window right edge
          idleAnimations.push("scratchWallW");
        }
        if (within(0, x - window.right, edgeThreshold)) {
          // inside window right edge or outside window left edge
          idleAnimations.push("scratchWallE");
        }
        if (within(-edgeThreshold, y - window.top, 0)) {
          // outside window top edge
          idleAnimations.push("scratchWallN");
        }
        if (within(0, y - window.bottom, edgeThreshold)) {
          // outside window bottom edge
          idleAnimations.push("scratchWallS");
        }
      }
      idleAnimation = idleAnimations[randomIntn(idleAnimations.length)];
    }

    function resetIdleAnimation() {
      idleAnimation = null;
      idleAnimationFrame = 0;
      x = idleLastX;
      y = idleLastY;
    }

    if (idleAnimationFrame == 0) {
      idleLastX = x;
      idleLastY = y;
    }

    switch (idleAnimation) {
      case "sleeping": {
        if (idleAnimationFrame < 8) {
          setSprite("tired", 0);
          break;
        }
        setSprite("sleeping", Math.floor(idleAnimationFrame / 4));
        if (idleAnimationFrame > 192) {
          resetIdleAnimation();
        }
        break;
      }
      case "scratchWallN":
      case "scratchWallS":
      case "scratchWallE":
      case "scratchWallW":
      case "scratchSelf": {
        setSprite(idleAnimation, idleAnimationFrame);
        if (idleAnimationFrame > 9) {
          resetIdleAnimation();
        }
        break;
      }
      default: {
        setSprite("idle", 0);
        return;
      }
    }
    idleAnimationFrame++;
  }

  function update() {
    frameCount += 1;
    const diffX = x - mouseX;
    const diffY = y - mouseY;
    const distance = Math.sqrt(diffX ** 2 + diffY ** 2);

    if (distance < speed || distance < 16) {
      idle();
      return;
    }

    idleAnimation = null;
    idleAnimationFrame = 0;

    if (idleTime > 1) {
      setSprite("alert", 0);
      // count down after being alerted before moving
      idleTime = Math.min(idleTime, 7);
      idleTime--;
      return;
    }

    let direction: string;
    direction = diffY / distance > 0.5 ? "N" : "";
    direction += diffY / distance < -0.5 ? "S" : "";
    direction += diffX / distance > 0.5 ? "W" : "";
    direction += diffX / distance < -0.5 ? "E" : "";
    setSprite(direction, frameCount);

    x -= (diffX / distance) * speed;
    y -= (diffY / distance) * speed;

    x = clamp(16, x, screenWidth - 16);
    y = clamp(16, y, screenHeight - 16);
  }

  let drawing = false;
  let lastTimestamp = 0;
  function onFrame(timestamp: number) {
    if (!lastTimestamp) {
      lastTimestamp = timestamp;
    }
    if (timestamp - lastTimestamp > 100) {
      lastTimestamp = timestamp;
      update();
    }
    if (drawing) {
      window.requestAnimationFrame(onFrame);
    }
  }

  svelte.onMount(() => {
    drawing = true;
    window.requestAnimationFrame(onFrame);
    return () => (drawing = false);
  });
</script>

<svelte:window
  bind:innerWidth={screenWidth}
  bind:innerHeight={screenHeight}
  on:mousemove={({ clientX, clientY }) => {
    mouseX = clientX;
    mouseY = clientY;
  }}
/>

<div
  aria-hidden="true"
  class="neko"
  style="--x: {x}px; --y: {y}px; background-position: {spriteX}px {spriteY}px;"
  transition:fade={{ duration: 100 }}
/>

<style>
  .neko {
    width: 32px;
    height: 32px;
    position: absolute;
    pointer-events: none;
    background-image: url(./oneko.gif);
    image-rendering: pixelated;
    z-index: 9999;
    left: calc(var(--x) - 16px);
    top: calc(var(--y) - 16px);
  }
</style>
