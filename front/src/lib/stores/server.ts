import { writable } from "svelte/store";

function persistent(key: string, initial: any) {
  const stored = localStorage.getItem(key);
  const data = writable(stored ? stored : initial);

  data.subscribe((value) => {
    localStorage.setItem(key, value);
  });

  return data;
}

export const serverUrl = persistent("serverUrl", "");
export const serverPassword = persistent("serverPassword", "");
