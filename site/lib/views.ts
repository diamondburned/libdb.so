import { persisted } from "svelte-persisted-store";
import { writable, derived, get } from "svelte/store";
import type * as store from "svelte/store";
import type { ComponentType, SvelteComponent } from "svelte";

import Terminal from "#/libdb.so/site/components/Terminal/index.svelte";
import TerminalIcon from "#/libdb.so/public/_fs/.icons/papirus/terminal.svg?url";

import Portfolio from "#/libdb.so/site/components/Portfolio/index.svelte";
import SystemUsersIcon from "#/libdb.so/public/_fs/.icons/papirus/system-users.svg?url";

export type Application = {
  component: ComponentType<SvelteComponent<{ state: store.Readable<WindowState> }>>;
  iconURL: string;
  title: string;
};
export type ApplicationID = keyof typeof applications;

// applications contains the list of available applications.
// At its core, it map application names to their Svelte components.
// It is constant always.
export const applications = {
  terminal: {
    component: Terminal,
    iconURL: TerminalIcon,
    title: "Terminal",
  },
  portfolio: {
    component: Portfolio,
    iconURL: SystemUsersIcon,
    title: "About Me",
  },
} satisfies Record<string, Application>;

export type WindowState = {
  applicationID: ApplicationID;
  windowID: WindowID;
  dimensions: Dimensions;
  visible: boolean;
  focused: boolean;
  maximized: boolean;
};

export type Dimensions = {
  x: number;
  y: number;
  width: number;
  height: number;
};

// A local storage-persisted list of opened windows.
// Only windows registered as applications are allowed.
export const windowStates = persisted<WindowState[]>("window_states_v2", []);

export function windowState();

// Like windows, but joins with the applications to get the native Svelte
// component.
export const windows = derived(windowStates, (windowStates) => {
  return windowStates.map((window) => ({
    ...window,
    ...applications[window.applicationID],
  }));
});

export function toggleVisibility(window: WindowState) {
  windowStates.update((windows) => {
    if (!window.visible) {
      // Currently not visible, so bring to focus.
      window.visible = true;
      window.focused = true;
    } else if (window.focused) {
      // Find the next visible window to focus.
      const nextVisible = windows.find((w) => w.visible && !w.focused);

      // Already on top, so hide it.
      window.visible = false;
      window.focused = false;

      // Focus the next visible window if there is one.
      if (nextVisible) {
        nextVisible.focused = true;
      }
    } else {
      // Otherwise, bring to focus.
      window.focused = true;
    }
    return windows;
  });
}

export function bringToFocus(window: WindowState) {
  windowStates.update((windows) => {
    windows.forEach((w) => {
      w.focused = false;
    });

    // Likely but not necessarily in `windows`.
    window.focused = true;

    return windows;
  });
}

// DragState is the state of a drag operation. It helps implement window
// dragging using the cursor.
export class DragState {
  initialOffsetX: number;
  initialOffsetY: number;

  constructor(
    // posX is the X offset of the window at the start of the drag.
    public posX: number,
    // posY is the Y offset of the window at the start of the drag.
    public posY: number,
    // cursorX is the X coordinate of the cursor at the start of the drag.
    public cursorX: number,
    // cursorY is the Y coordinate of the cursor at the start of the drag.
    public cursorY: number,
    // setPosition is the callback to set the position of the window.
    // The calculated position is passed as arguments.
    public readonly setPosition: (x: number, y: number) => void,
  ) {
    this.initialOffsetX = posX - cursorX;
    this.initialOffsetY = posY - cursorY;
  }

  // update updates the position of the window based on the current cursor
  // position.
  update(cursorX: number, cursorY: number) {
    this.setPosition(this.initialOffsetX + cursorX, this.initialOffsetY + cursorY);
  }
}

export type WindowID = string;

export function generateWindowID(): WindowID {
  return crypto.randomUUID();
}
