import * as store from "svelte/store";

export type View = "terminal" | "portfolio";

export const activeViews = store.writable<Record<View, boolean>>();
export const focusedView = store.writable<View | null>("portfolio");

// Load from local storage.
const localStorageActiveViews = localStorage.getItem("active_views");
activeViews.set(
  localStorageActiveViews
    ? (JSON.parse(localStorageActiveViews) as Record<View, boolean>)
    : { portfolio: true, terminal: false }
);
const localStorageTopView = localStorage.getItem("top_view");
focusedView.set(
  localStorageTopView ? (JSON.parse(localStorageTopView) as View | null) : null
);

// Bind to local storage.
activeViews.subscribe((activeViews) => {
  localStorage.setItem("active_views", JSON.stringify(activeViews));
});
focusedView.subscribe((topView) => {
  localStorage.setItem("top_view", JSON.stringify(topView));
});

export function toggleView(view: View) {
  focusedView.update((focusedView) => {
    activeViews.update((activeViews) => {
      if (!activeViews[view]) {
        // Currently not active, so bring to focus.
        activeViews[view] = true;
        focusedView = view;
      } else if (view == focusedView) {
        // Already on top, so hide it.
        activeViews[view] = false;
        // Bring the next visible view to focus. If there is none, then set to
        // null.
        focusedView = Object.keys(activeViews).find(
          (view) => activeViews[view as View]
        ) as View | null;
      } else {
        // Otherwise, bring to focus.
        focusedView = view;
      }
      return activeViews;
    });
    return focusedView;
  });
}

export function bringToFocus(view: View): void {
  focusedView.set(view);
}

export function viewIsActive(view: View): store.Readable<boolean> {
  return store.derived(activeViews, (activeViews) => activeViews[view]);
}

export function viewIsFocused(view: View): store.Readable<boolean> {
  return store.derived(focusedView, (topView) => topView === view);
}

// Save some states for Show Desktop.
let savedActiveViews: Record<View, boolean> | null = null;
let savedTopView: View | null = null;
let showDesktop = false;

export function toggleShowDesktop() {
  if (showDesktop) {
    // Save the current states.
    savedActiveViews = store.get(activeViews);
    savedTopView = store.get(focusedView);
    // Hide all views.
    activeViews.set({ portfolio: false, terminal: false });
    focusedView.set(null);
  } else if (savedActiveViews && savedTopView) {
    // Restore the saved states.
    activeViews.set(savedActiveViews);
    focusedView.set(savedTopView);
  }
  showDesktop = !showDesktop;
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
    public readonly setPosition: (x: number, y: number) => void
  ) {
    this.initialOffsetX = posX - cursorX;
    this.initialOffsetY = posY - cursorY;
  }

  // update updates the position of the window based on the current cursor
  // position.
  update(cursorX: number, cursorY: number) {
    this.setPosition(
      this.initialOffsetX + cursorX,
      this.initialOffsetY + cursorY
    );
  }
}
