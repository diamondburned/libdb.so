import * as store from "svelte/store";

const key = "nsfw-v1";

let emitting = false;

const value = store.writable<boolean>(
  localStorage.getItem(key) == "true",
  (set) => {
    const listener = (ev: StorageEvent) => {
      if (!emitting && ev.key == key) {
        set(ev.newValue == "true");
      }
    };
    window.addEventListener("storage", listener);
    return () => window.removeEventListener("storage", listener);
  }
);

value.subscribe((value) => {
  const oldValue = localStorage.getItem(key);
  const newValue = JSON.stringify(value);

  localStorage.setItem(key, newValue);

  window.dispatchEvent(
    new StorageEvent("storage", {
      key,
      oldValue,
      newValue,
    })
  );
});

export default value;
