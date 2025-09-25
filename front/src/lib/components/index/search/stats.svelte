<script lang="ts">
  import type { Server } from "$src/lib/types";
  import { Database, File, Save } from "@lucide/svelte";
    import { unmount } from "svelte";

  const { serverInfo }: { serverInfo: Server | null } = $props();

  function mbToGb(mb: number): number {
    return Math.round((mb / 1024) * 100) / 100;
  }
</script>

<div class="stats  stats-vertical md:stats-horizontal">
  <div class="stat">
    <div class="stat-figure text-secondary">
      <File />
    </div>
    <div class="stat-title">Rows available</div>
    <div class="stat-value">
      {serverInfo?.TotalRows !== undefined
        ? serverInfo.TotalRows.toLocaleString("fr")
        : "-- --- --- ---"}
    </div>
  </div>
  <div class="stat">
    <div class="stat-figure text-secondary">
      <Database />
    </div>
    <div class="stat-title">Data wells available</div>
    <div class="stat-value">
      {serverInfo?.TotalDataleaks !== undefined
        ? serverInfo.TotalDataleaks.toLocaleString("fr")
        : "---"}
    </div>
  </div>
  <div class="stat">
    <div class="stat-figure text-secondary">
      <Save />
    </div>
    <div class="stat-title">Storage used</div>
    <div class="stat-value">
      {serverInfo?.TotalSize !== undefined
        ? mbToGb(serverInfo.TotalSize).toLocaleString("fr") + " Gb"
        : "--- Gb"}
    </div>
  </div>
</div>
