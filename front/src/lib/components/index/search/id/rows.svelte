<script lang="ts">
  import type { Result } from "$src/lib/types";
  import { Search } from "@lucide/svelte";
  import Row from "./row.svelte";

  const { result }: { result: Result } = $props();

  let page = $state(1);
  let totalPages = $state(0);
  const perPage = 20;

  let paginated = $state<Record<string, string>[]>([]);

  $effect(() => {
    if (!result || !result.LeakResult.Rows) {
      return;
    }
    if (filter.trim() !== "") {
      let rowsFiltered = result.LeakResult.Rows.filter((row) => {
        const rowText = flattenObjectValues(row);
        const keywords = filter
          .toLowerCase()
          .split(/\s+/)
          .filter((t) => t.length > 0);
        if (keywords.length === 0) return true;
        for (const term of keywords) {
          if (!rowText.includes(term)) {
            return false;
          }
        }
        return true;
      });

      totalPages = Math.ceil(rowsFiltered.length / perPage);
      const start = (page - 1) * perPage;
      const end = start + perPage;
      paginated = rowsFiltered.slice(start, end);
      if (page > totalPages) {
        page = totalPages > 0 ? totalPages : 1;
      }
    } else {
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

  function flattenObjectValues(obj: any): string {
    let resultText = "";

    function recurse(current: any) {
      for (const key in current) {
        if (Object.prototype.hasOwnProperty.call(current, key)) {
          const value = current[key];
          if (typeof value === "object" && value !== null) {
            recurse(value);
          } else if (value !== undefined && value !== null) {
            resultText += String(value).toLowerCase() + " ";
          }
        }
      }
    }

    recurse(obj);
    return resultText;
  }

  let top: any = $state();
  let filter = $state("");
</script>

<div bind:this={top} class="absolute -mt-[100px]"></div>
<label class="input input-xs w-full mb-2">
  <Search size={12} />
  <input class="grow" placeholder="Filter" bind:value={filter} />
</label>
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
