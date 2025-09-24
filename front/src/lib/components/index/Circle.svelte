<script lang="ts">
  import { cn } from "$lib/utils";

  let { children, label ="", ref=$bindable()} = $props();
  let className: any = $state("");
  export { className as class };
</script>

<div class="tooltip z-10" data-tip={label}>
  <div
    bind:this={ref}
    class={cn(
      "bg-base-100 hover:bg-base-200 border-base-200 z-10 flex h-12 w-12 items-center justify-center rounded-full border-2 p-3 transition-all duration-200 cursor-pointer",
      className,
    )}
  >
    {@render children()}
  </div>
</div>
