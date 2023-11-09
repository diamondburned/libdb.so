import * as store from "svelte/store";

type View = null | "terminal" | "portfolio";

export const defaultView: View = "terminal";

const currentView = store.writable<View>(
  (JSON.parse(localStorage.getItem("current_view") ?? "null") as View) ||
    defaultView
);
currentView.subscribe((currentView) => {
  localStorage.setItem("current_view", JSON.stringify(currentView));
});

const lastView = store.writable<View>(null);

// switchView switches to the given view, or minimizes the current view if it
// is already active.
export function switchView(view: View) {
  currentView.update((currentView) => (currentView == view ? null : view));
  lastView.set(null);
}

// viewDesktop switches to desktop, minimizing the current view.
export function viewDesktop() {
  lastView.update((lastView) => {
    return lastView == null ? store.get(currentView) : null;
  });
  currentView.update((currentView) => {
    return currentView == null ? store.get(lastView) : null;
  });
}

// view is the current view.
export const view: store.Readable<View> = currentView;
