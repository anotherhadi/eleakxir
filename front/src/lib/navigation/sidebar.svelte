<script lang="ts">
  import { BrushCleaning, Database, FileSearch, Home, Scale, Search } from "@lucide/svelte";
  import SidebarMenuItem from "./sidebar-menu-item.svelte";

  interface NavItem {
    title?: string;
    url?: string;
    items?: NavItem[];
    icon?: any;
    type?: "link" | "parent" | "divider";
  }

  const Nav: NavItem[] = [
    {
      title: "Home",
      url: "/",
      icon: Home,
    },
    {
      title: "Search",
      url: "/search",
      icon: Search,
    },
    {
      title: "Data wells",
      url: "/dataleaks",
      icon: Database,
    },
    { type: "divider" },
    {
      title: "Parquet files & Rules",
      url: "/parquet",
      icon: Scale,
    },
    {
      title: "Leak utils",
      url: "/leak-utils",
      icon: BrushCleaning,
    },
  ];
</script>

<div class="drawer-side z-[101]">
  <label for="menu-toggle" aria-label="close sidebar" class="drawer-overlay"
  ></label>
  <ul class="menu bg-base-200 text-base-content min-h-full w-56 p-4 gap-2">
    {#each Nav as item}
      {#if item.type === "divider"}
        <li class="menu-title pt-2">
          <div class="divider"></div>
        </li>
      {:else}
        <SidebarMenuItem {item} />
      {/if}
    {/each}
  </ul>
</div>
