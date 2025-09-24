<script lang="ts">
  import { Github, Menu, Search } from "@lucide/svelte";
  import Logo from "../components/logo.svelte";
  import DarkModeToggle from "../components/dark-mode-toggle.svelte";
  import { cn } from "../utils";
  import ServerDialog from "../components/server-dialog.svelte";

  let y = $state(0);
</script>

<svelte:window bind:scrollY={y} />

<nav
  class={cn(
    "w-full h-20 flex gap-5 items-center justify-between px-6 fixed transition-colors duration-1000 z-[100]",
    y === 0 || "bg-base-200",
  )}
>
  <div class="flex gap-2 items-center">
    <label
      for="menu-toggle"
      class="btn btn-ghost btn-sm btn-square drawer-button"
    >
      <Menu size={16} />
    </label>
    <a href="/" class="flex gap-2 items-center">
    <Logo size={16} />
    <p>Eleakxir</p>
    </a>
  </div>

  <div class="flex gap-2 items-center">
    <a href="/search">
      <button class="btn btn-sm btn-ghost btn-square">
        <Search size={16} />
      </button>
    </a>
    <a href="https://github.com/anotherhadi/eleakxir">
      <button class="btn btn-sm btn-ghost btn-square">
        <Github size={16} />
      </button>
    </a>
    <DarkModeToggle />

    <ServerDialog class="btn-sm btn-square" />
  </div>
</nav>
