<script lang="ts">
  import type { Dataleak } from "$src/lib/types";
  import { Clipboard, Search } from "@lucide/svelte";
  import FaviconOrIcon from "../../favicon-or-icon.svelte";
  import DatawellPopup from "./datawell-popup.svelte";
  import { serverUrl, serverPassword } from "$src/lib/stores/server";
  import axios from "axios";

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

  let copyText = $state("Copy to clipboard");

  async function copyDataleaksInformation(withSample: boolean) {
    if (!filteredDataleaks || filteredDataleaks.length === 0) {
      copyText = "No dataleaks to copy";
      return;
    }

    let fullText = "";

    for (const dataleak of filteredDataleaks) {
      fullText += `Name: ${dataleak.Name}\n`;
      fullText += `Length: ${dataleak.Length.toLocaleString()}\n`;

      if (withSample) {
        await axios
          .get(`${$serverUrl}/dataleak/sample`, {
            params: { path: dataleak.Path },
            headers: { "X-Password": $serverPassword },
          })
          .then((r) => {
            const samples: string[][] = r.data.Sample;

            if (samples && samples.length > 0) {
              fullText += "Sample:\n";
              const headers = samples[0];
              const rows = samples.slice(1);

              fullText += headers.join(", ") + "\n";
              for (const row of rows) {
                fullText += row.join(", ") + "\n";
              }
            }
          })
          .catch((err) => {
            console.error("Failed to fetch sample for", dataleak.Name, err);
            fullText += "Sample: [Failed to fetch]\n";
          });
      }
      fullText += "\n";
    }

    try {
      await navigator.clipboard.writeText(fullText);
      copyText = "Copied!";
      setTimeout(() => (copyText = "Copy to clipboard"), 2000);
    } catch (err) {
      console.error("Failed to copy:", err);
      copyText = "Copy failed";
      setTimeout(() => (copyText = "Copy to clipboard"), 2000);
    }
  }

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
  <div class="flex items-center gap-1">
    <label class="input input-xs w-full grow">
      <Search size={12} />
      <input class="grow" placeholder="Filter" bind:value={filter} />
    </label>

    <div class="dropdown dropdown-end">
      <div
        tabindex="0"
        role="button"
        class="btn btn-ghost btn-sm size-6 btn-square m-1"
      >
        <Clipboard size={12} />
      </div>
      <div
        class="dropdown-content bg-base-300 rounded-box z-1 w-52 p-2 shadow-2xl grid gap-2"
      >
        <button
          class="btn btn-xs"
          onclick={() => copyDataleaksInformation(false)}>{copyText}</button
        >
        <button
          class="btn btn-xs"
          onclick={() => copyDataleaksInformation(true)}
          >{copyText} (w/ sample)</button
        >
      </div>
    </div>
  </div>

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
              <th>
                <DatawellPopup dataleak={item} />
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
