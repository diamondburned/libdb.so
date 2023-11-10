import * as store from "svelte/store";

export type Toast = {
  text: string;
  class?: string;
  dismissed?: () => void;

  _timeout?: number;
};

export class ToastStore {
  readonly store = store.writable<Toast[]>([]);

  constructor() {
    console.log(this);
  }

  add(toast: Toast, timeout = 0) {
    this.store.update((toasts) => {
      toasts.push(toast);
      if (timeout > 0) {
        toast._timeout = window.setTimeout(() => this.remove(toast), timeout);
      }
      return toasts;
    });
  }

  remove(toast: Toast) {
    let dismissed = false;

    this.store.update((toasts) => {
      const toastIx = toasts.indexOf(toast);
      if (toastIx == -1) {
        return toasts;
      }

      toasts.splice(toastIx, 1);
      dismissed = true;
      return toasts;
    });

    if (dismissed) {
      if (toast._timeout) {
        window.clearTimeout(toast._timeout);
      }
      if (toast.dismissed) {
        toast.dismissed();
      }
    }
  }

  subscribe(
    run: store.Subscriber<Toast[]>,
    invalidate?: (value?: Toast[]) => void
  ): store.Unsubscriber {
    return this.store.subscribe(run, invalidate);
  }
}
