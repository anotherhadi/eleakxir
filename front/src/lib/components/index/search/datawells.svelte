<script lang="ts">
  import type { Dataleak } from "$src/lib/types";
  import { Search } from "@lucide/svelte";
  import FaviconOrIcon from "../../favicon-or-icon.svelte";

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
  let filteredDataleaks = $state<Dataleak[]>([]);
  let paginatedDataleaks = $state<Dataleak[]>([]);
  let totalPages = $state(0);

  function sortByModTime(arr: Dataleak[]): Dataleak[] {
    return arr.slice().sort((a, b) => {
      return new Date(b.ModTime).getTime() - new Date(a.ModTime).getTime();
    });
  }

  $effect(() => {
    const sortedData = sortByModTime(dataleaks);

    if (filter.trim() === "") {
      filteredDataleaks = sortedData;
    } else {
      const lowerFilter = filter.toLowerCase();
      filteredDataleaks = sortedData.filter((item) =>
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

  function getDomain(dataleakName: string) {
    const firstPart = dataleakName.split(" ")[0].toLowerCase();
    const domainRegex =
      /^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$/;
    if (domainRegex.test(firstPart)) {
      return firstPart.replace("_", "-");
    }
    return null;
  }

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
          <th></th>
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
                <FaviconOrIcon
                  url={getDomain(item.Name)}
                  icon={item.Columns.includes("password")
                    ? "password"
                    : item.Columns.includes("email")
                      ? "email"
                      : ""}
                />
              </th>
              <th class="text-nowrap">
                {item.Name}
              </th>
              <td>{item.Length.toLocaleString("fr")}</td>
              {#if showColumns}
                <td class="capitalize text-nowrap">
                  {item.Columns.map((col) => col.replace(/_/g, " ")).join(", ")}
                </td>
              {/if}
            </tr>
          {/each}
        {:else}
          <tr class="hover:bg-base-300">
            <td colspan={100} class="text-center leading-9"
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
