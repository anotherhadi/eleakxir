<script lang="ts">
  import {
    ChevronDown,
    ChevronUp,
    type Icon as IconType,
  } from "@lucide/svelte";
  import type { Snippet } from "svelte";
    import { cn } from "../utils";

  let isOpen = $state<boolean>(false);

  const {
    imageUrl,
    icon,
    title,
    subtitle,
    children,
  }: {
    imageUrl?: string | null;
    icon: typeof IconType;
    title: string;
    subtitle?: string;
    children?: Snippet;
  } = $props();
</script>

<button
  class={cn("list-row text-left bg-base-200/40",
    children != null ? "cursor-pointer hover:bg-base-300/75" : ""
  )}
  class:bg-base-300={isOpen}
  class:rounded-b-none={isOpen}
  onclick={() => {
    if (children != null) {
      isOpen = !isOpen;
    }
  }}
>
  <div>
    {#if imageUrl && imageUrl.length > 0}
      <img
        src="https://icons.duckduckgo.com/ip3/{imageUrl}.ico"
        class="size-10 rounded-box bg-neutral"
        alt="Favicon of {imageUrl}"
      />
    {:else}
      {@const Icon = icon}
      <div
        class="size-10 rounded-box bg-neutral items-center justify-center flex"
      >
        <Icon />
      </div>
    {/if}
  </div>
  <div class="flex flex-col justify-center">
    <div class="font-semibold">{title}</div>
    {#if subtitle != null && subtitle.length !== 0}
      <div class="text-xs uppercase font-semibold opacity-60">
        {subtitle}
      </div>
    {/if}
  </div>
  {#if children != null}
    <div class="btn btn-square btn-ghost">
      {#if isOpen}
        <ChevronUp size={12} />
      {:else}
        <ChevronDown size={12} />
      {/if}
    </div>
  {/if}
</button>
{#if children != null}
  {#if isOpen}
    <li class="list-row bg-base-200 rounded-t-none mb-2">
      {@render children()}
    </li>
  {/if}
{/if}
