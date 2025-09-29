<script lang="ts">
  import { Database, Key, Mail } from "@lucide/svelte";
  import { cn } from "../utils";

  let { url, icon="", class: className = "" } = $props();

  let hasError = $state(false);
</script>

<div
  class={cn(
    "w-10 h-10 rounded-box bg-neutral text-neutral-content items-center justify-center flex",
    className,
  )}
>
  {#if url && !hasError}
    <img
      src={"https://icons.duckduckgo.com/ip3/" + url + ".ico"}
      class="rounded-box w-full h-full"
      alt="Favicon de {url}"
      onerror={() => { hasError = true }}
    />
  {:else if icon == "password"}
    <Key />
  {:else if icon == "email"}
    <Mail />
  {:else}
    <Database />
  {/if}
</div>
