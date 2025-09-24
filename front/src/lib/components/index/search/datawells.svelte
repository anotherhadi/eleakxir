<script lang="ts">
  import type { Dataleak } from "$src/lib/types";
  import { Replace, Search } from "@lucide/svelte";

  let {
    dataleaks,
    perPage = 5,
    showColumns = false,
  }: {
    dataleaks: Dataleak[];
    perPage?: number;
    showColumns?: boolean;
  } = $props();

  let page = $state(1);
  let filter = $state("");
  let filteredDataleaks = $state<Dataleak[]>(dataleaks);
  let paginatedDataleaks = $state<Dataleak[]>([]);
  let totalPages = $state(0);

  $effect(() => {
    if (filter.trim() === "") {
      filteredDataleaks = dataleaks;
    } else {
      const lowerFilter = filter.toLowerCase();
      filteredDataleaks = dataleaks.filter((item) =>
        item.Name.toLowerCase().includes(lowerFilter),
      );
    }
    page = 1;
  });

  $effect(() => {
    if (filteredDataleaks) {
      totalPages = Math.ceil(filteredDataleaks.length / perPage);
      const start = (page - 1) * perPage;
      const end = start + perPage;
      paginatedDataleaks = filteredDataleaks.slice(start, end);
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
          <th>Name</th>
          <th>Number of rows</th>
          {#if showColumns}
            <th>Columns</th>
          {/if}
        </tr>
      </thead>
      <tbody>
        {#if paginatedDataleaks.length > 0}
          {#each paginatedDataleaks as item}
            <tr class="hover:bg-base-300">
              <th>
                {item.Name}
              </th>
              <td>{item.Length.toLocaleString("fr")}</td>
              {#if showColumns}
                <td class="capitalize">
                  {item.Columns.map((col) => col.replace(/_/g, " ")).join(", ")}
                </td>
              {/if}
            </tr>
          {/each}
        {:else}
          <tr class="hover:bg-base-300">
            <td colspan={3} class="text-center leading-9"
              ><span class="text-3xl">(·.·)</span><br />No data wells found</td
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
