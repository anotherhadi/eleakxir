import { svelte } from "@sveltejs/vite-plugin-svelte";
import { router } from "sv-router/vite-plugin";
import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
  plugins: [tailwindcss(), svelte({}), router()],
  resolve: {
    alias: {
      $lib: path.resolve("./src/lib"),
      $src: path.resolve("./src"),
    },
  },
});
