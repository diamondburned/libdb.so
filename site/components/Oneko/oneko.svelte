<script lang="ts">
  import * as svelte from "svelte";
  import type { Window } from "#/libdb.so/site/lib/views.js";
  import { speed, sakura } from "./oneko.js";
  import { fade } from "svelte/transition";

  const character = sakura;

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

  export let windows: Window[] = [];

  function windowToRect(window: Window) {
    return {
      top: window.y,
      left: window.x,
      bottom: window.y + window.height,
      right: window.x + window.width,
    };
  }

  $: rectangles = [
    // Main screen
    {
      top: 0,
      left: 0,
      bottom: screenHeight,
      right: screenWidth,
    },
    ...windows.map(windowToRect),
  ];

  $: midpoints = rectangles.map((rect) => ({
    x: rect.left + (rect.right - rect.left) / 2,
    y: rect.top + (rect.bottom - rect.top) / 2,
  }));

  function setSprite(name: string, frame: number) {
    const sprites = character.spriteSets[name];
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

  function resetIdleAnimation() {
    idleAnimation = null;
    idleAnimationFrame = 0;
  }

  function idle() {
    idleTime++;

    // every ~ 20 seconds
    if (idleTime > 10 && Math.random() < 0.05 && !idleAnimation) {
      const idleAnimations = ["sleeping", "scratchSelf", "scratchWallS"];
      idleAnimation = idleAnimations[randomIntn(idleAnimations.length)];
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
      case "scratchWallS":
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
    frameCount++;

    let nearest = rectangles[0];
    let nearestMid = midpoints[0];

    for (let i = 1; i < rectangles.length; i++) {
      const rect = rectangles[i];
      const mid = midpoints[i];
      if (
        i == 1 ||
        Math.abs(mid.y - mouseY) < Math.abs(nearestMid.y - mouseY)
      ) {
        nearest = rect;
      }
    }

    const nextX = clamp(nearest.left + 16, mouseX, nearest.right - 16);
    const nextY = nearest.top - 14;

    const diffX = x - nextX;
    const diffY = y - nextY;

    const distance = Math.sqrt(diffX ** 2 + diffY ** 2);
    if (distance < speed) {
      idle();
      return;
    }

    resetIdleAnimation();

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
  style="
    --x: {x}px;
    --y: {y}px;
    --url: url({character.url});
    background-position: {spriteX}px {spriteY}px;
  "
  transition:fade={{ duration: 100 }}
/>

<style>
  .neko {
    width: 32px;
    height: 32px;
    position: absolute;
    pointer-events: none;
    background-image: var(--url);
    image-rendering: pixelated;
    z-index: 9999;
    left: calc(var(--x) - 16px);
    top: calc(var(--y) - 16px);
  }
</style>
