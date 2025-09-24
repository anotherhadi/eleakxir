<script lang="ts">
  import type { Result } from "$src/lib/types";
  import Row from "./row.svelte";

  const { result }: { result: Result } = $props();

  let page = $state(1);
  let totalPages = $state(0);
  const perPage = 20;

  let paginated = $state<Record<string, string>[]>([]);

  $effect(() => {
    if (result && result.LeakResult.Rows) {
      totalPages = Math.ceil(result.LeakResult.Rows.length / perPage);
      const start = (page - 1) * perPage;
      const end = start + perPage;
      paginated = result.LeakResult.Rows.slice(start, end);
      if (page > totalPages) {
        page = totalPages > 0 ? totalPages : 1;
      }
    }
  });

  function goToFirstPage() {
    page = 1;
    top.scrollIntoView();
  }

  function previousPage() {
    if (page > 1) {
      page--;
      top.scrollIntoView();
    }
  }

  function nextPage() {
    if (page < totalPages) {
      page++;
      top.scrollIntoView();
    }
  }

  let top: any = $state();
</script>

<div bind:this={top} class="absolute -mt-[100px]"></div>
{#if result}
  <ul class="list bg-base-100 rounded-box shadow-md">
    {#each paginated as row (row)}
      <Row {row} />
    {/each}
  </ul>

  {#if totalPages > 1}
    <div class="join m-auto mt-5">
      <button class="join-item btn" onclick={previousPage} disabled={page === 1}
        >«</button
      >
      <button class="join-item btn" onclick={goToFirstPage}
        >Page {page} / {totalPages}</button
      >
      <button
        class="join-item btn"
        onclick={nextPage}
        disabled={page === totalPages}>»</button
      >
    </div>
  {/if}
{:else}
  No result
{/if}
