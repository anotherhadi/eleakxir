<script lang="ts">
    import { Moon, Sun, SunMoon } from "@lucide/svelte";
</script>

<div class="dropdown dropdown-end">
  <div tabindex="0" role="button" class="btn btn-ghost btn-sm btn-square m-1">
    <SunMoon size={16 }/>
  </div>
  <div
    class="dropdown-content bg-base-300 rounded-box z-1 w-52 p-2 shadow-2xl grid gap-2"
  >
    <button
      data-set-theme="light"
      class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
      aria-label="Light"
    >
      <Sun />
      Light
    </button>
    <button
      data-set-theme="dark"
      class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
      aria-label="Dark"
    >
      <Moon />
      Dark
    </button>
  </div>
</div>
