<script lang="ts">
  import type { Result } from "$src/lib/types";
  import { formatDate } from "$src/lib/utils";
  import { BadgeInfo, Clock, File } from "@lucide/svelte";

  const { result }: { result: Result } = $props();

  let nresult = $state(0);
  $effect(() => {
    const r = [
      result.LeakResult.Rows?.length | 0,
      result.GithubResult.EmailResult?.Commits?.length | 0,
      result.GithubResult.EmailResult?.Spoofing ? 1 : 0,
      result.GithubResult.UsernameResult?.Commits?.length | 0,
    ];
    nresult = r.reduce((a, b) => a + b, 0);
  });
</script>

<div class="stats stats-vertical md:stats-horizontal">
  <div class="stat">
    <div class="stat-figure text-secondary">
      <File />
    </div>
    <div class="stat-title">Results</div>
    <div class="stat-value" class:animate-pulse={result.Status === "pending"}>
      {nresult.toLocaleString("fr")}
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
      {result.Status}
      {#if result.Status === "pending"}
        <span class="loading loading-dots loading-xs ml-2"></span>
      {/if}
    </div>
  </div>
</div>
