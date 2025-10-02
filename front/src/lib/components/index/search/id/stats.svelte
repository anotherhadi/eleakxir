<script lang="ts">
  import type { Result } from "$src/lib/types";
  import { formatDate } from "$src/lib/utils";
  import { BadgeInfo, Clock, File } from "@lucide/svelte";

  const { result }: { result: Result } = $props();
</script>

<div class="stats stats-vertical md:stats-horizontal">
  <div class="stat">
    <div class="stat-figure text-secondary">
      <File />
    </div>
    <div class="stat-title">Results</div>
    <div class="stat-value" class:animate-pulse={result.Status === "pending"}>
      {result.ResultsCount.toLocaleString("fr")}
      {#if result.Status === "pending"}
        <span class="loading loading-dots loading-xs ml-2"></span>
      {/if}
    </div>
  </div>
  <div class="stat">
    <div class="stat-figure text-secondary">
      <Clock />
    </div>
    <div class="stat-title">Date</div>
    <div class="stat-value">
      {formatDate(result.Date)}
    </div>
  </div>
  <div class="stat">
    <div class="stat-figure text-secondary">
      <BadgeInfo />
    </div>
    <div class="stat-title">Status</div>
    <div class="stat-value" class:animate-pulse={result.Status === "pending"}>
      {#if result.Status === "pending"}
        Pending
        <span class="loading loading-dots loading-xs ml-2"></span>
      {:else if result.Status === "completed" && result.ResultsCount === 0}
        No results
      {:else if result.Status === "completed"}
        Completed
      {:else}
        {result.Status}
      {/if}
    </div>
  </div>
</div>
