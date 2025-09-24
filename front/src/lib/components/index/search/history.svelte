<script lang="ts">
  import type { History } from "$src/lib/types";
    import { formatDate } from "$src/lib/utils";
  import { Search } from "@lucide/svelte";
  import { navigate, p } from "sv-router/generated";

  let { history, perPage = 5 }: { history: History; perPage?: number } =
    $props();

  let page = $state(1);
  let filter = $state("");
  let filteredHistory = $state<History>(history);
  let paginatedHistory = $state<History>([]);
  let totalPages = $state(0);

  $effect(() => {
    if (filter.trim() === "") {
      filteredHistory = history;
    } else {
      const lowerFilter = filter.toLowerCase();
      filteredHistory = history.filter((item) =>
        item.Query.Text.toLowerCase().includes(lowerFilter),
      );
    }
    page = 1;
  });

  $effect(() => {
    if (filteredHistory) {
      totalPages = Math.ceil(filteredHistory.length / perPage);
      const start = (page - 1) * perPage;
      const end = start + perPage;
      paginatedHistory = filteredHistory.slice(start, end);
      if (page > totalPages) {
        page = totalPages > 0 ? totalPages : 1;
      }
    }
  });

  function previousPage() {
    if (page > 1) {
      page--;
    }
  }

  function nextPage() {
    if (page < totalPages) {
      page++;
    }
  }
</script>

<div class="my-4 flex flex-col gap-2">
  <label class="input input-xs w-full">
    <Search size={12} />
    <input class="grow" placeholder="Filter" bind:value={filter} />
  </label>

  <div class="overflow-x-auto">
    <table class="table">
      <!-- head -->
      <thead>
        <tr>
          <th>Query</th>
          <th>Results</th>
          <th>Status</th>
          <th>Date</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#if paginatedHistory.length > 0}
          {#each paginatedHistory as item}
            <tr class="hover:bg-base-300">
              <th>
                <button
                  onclick={() => {
                    navigate(`/search/:id`, { params: { id: item.Id } });
                  }}
                  class="btn btn-link p-0 no-underline text-base-content"
                >
                  {item.Query.Text}
                </button>
              </th>
              <td>{item.Results}</td>
              <td
                ><div
                  class="badge badge-xs"
                  class:badge-success={item.Status === "completed"}
                  class:badge-warning={item.Status === "pending"}
                >
                  {item.Status}
                </div></td
              >
              <td>{formatDate(item.Date)}</td>
              <td
                onclick={() => {
                  navigate(`/search/:id`, { params: { id: item.Id } });
                }}
                ><button class="btn btn-xs btn-square"
                  ><Search size={11} /></button
                ></td
              >
            </tr>
          {/each}
        {:else}
          <tr class="hover:bg-base-300">
            <td colspan="5" class="text-center leading-9"
              ><span class="text-3xl">(·.·)</span><br />No history found</td
            >
          </tr>
        {/if}
      </tbody>
    </table>
  </div>

  {#if totalPages > 1}
    <div class="join m-auto mt-5">
      <button class="join-item btn" onclick={previousPage} disabled={page === 1}
        >«</button
      >
      <button class="join-item btn">Page {page} / {totalPages}</button>
      <button
        class="join-item btn"
        onclick={nextPage}
        disabled={page === totalPages}>»</button
      >
    </div>
  {/if}
</div>
