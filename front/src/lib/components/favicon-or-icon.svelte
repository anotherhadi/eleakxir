<script lang="ts">
  import { Database, Key, Mail } from "@lucide/svelte";
  import { cn } from "../utils";
    import { onMount } from "svelte";
    import axios from "axios";

  let { url, icon="", class: className = "" } = $props();

  let hasError = $state(false);
</script>

<div
  class={cn(
    "w-10 h-10 rounded-2xl bg-neutral text-neutral-content items-center justify-center flex",
    className,
  )}
>
  {#if url && !hasError}
    <img
      src={"https://icons.duckduckgo.com/ip3/" + url + ".ico"}
      class="bg-neutral rounded-2xl w-2/3 h-2/3"
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
