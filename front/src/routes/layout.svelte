<script lang="ts">
  import Sidebar from "$src/lib/navigation/sidebar.svelte";
  import Topbar from "$src/lib/navigation/topbar.svelte";
  import { onMount, type Snippet } from "svelte";
  import { Toaster } from 'svelte-sonner'
  import { themeChange } from "theme-change";

  let { children }: { children: Snippet } = $props();

  onMount(() => {
    themeChange(false);
  });
</script>

<Toaster 
  	toastOptions={{
		class: '!bg-base-300 !text-base-content !border-base-200',
	}}
/>
<div class="drawer min-h-svh">
  <input id="menu-toggle" type="checkbox" class="drawer-toggle" />
  <div class="drawer-content">
    <Topbar />
    <div class="mt-20"></div>
    {@render children()}
  </div>
  <Sidebar />
</div>
