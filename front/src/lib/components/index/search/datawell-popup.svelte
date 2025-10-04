<script lang="ts">
  import { serverPassword, serverUrl } from "$src/lib/stores/server";
  import type { Dataleak } from "$src/lib/types";
  import { Info } from "@lucide/svelte";
  import axios from "axios";

  const { dataleak }: { dataleak: Dataleak } = $props();

  let popupOpen = $state(false);
  let samples = $state<string[][]>([]);
  let copyText = $state("Copy to clipboard")

  async function getSample() {
    if (!dataleak) return;

    await axios
      .get(`${$serverUrl}/dataleak/sample`, {
        params: { path: dataleak.Path },
        headers: {
          "X-Password": $serverPassword,
        },
      })
      .then((r) => {
        samples = r.data.Sample;
      })
      .catch((e) => {
        console.error("Erreur lors du fetch sample:", e);
      });
  }

  async function copyToClipboard() {
    if (!dataleak || !samples || samples.length === 0) {
      copyText = "No data to copy";
      return;
    }

    const leakName = dataleak.Name;
    const columns = samples[0].join(", ");
    const sampleRows = samples.slice(1).map(r => r.join(", ")).join("\n");

    const textToCopy = `Leak Name: ${leakName}\nColumns: ${columns}\nSample:\n${sampleRows}`;

    try {
      await navigator.clipboard.writeText(textToCopy);
      copyText = "Copied!";
      setTimeout(() => copyText = "Copy Sample", 2000);
    } catch (err) {
      console.error("Failed to copy: ", err);
      copyText = "Copy failed";
      setTimeout(() => copyText = "Copy Sample", 2000);
    }
  }
</script>

<button
  class="text-nowrap flex gap-2 items-center hover:text-base-content/70"
  onclick={() => {
    popupOpen = true;
    getSample();
  }}
>
  {dataleak.Name}
  <Info size={12} />
</button>

<dialog class="modal modal-bottom sm:modal-middle" class:modal-open={popupOpen}>
  <div class="modal-box">
    <form method="dialog">
      <button
        onclick={() => (popupOpen = false)}
        class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button
      >
    </form>

    <div class="flex flex-col gap-5">
      <div>
        <h2 class="card-title mb-6">{dataleak.Name}</h2>

        {#if samples && samples.length > 1}
          <div
            class="overflow-x-auto border border-base-300 rounded-xl shadow-sm"
          >
            <table class="table table-zebra w-full text-sm">
              <thead class="bg-base-200">
                <tr>
                  {#each samples[0] as header, i (i)}
                    <th
                      class="font-semibold text-xs whitespace-nowrap capitalize"
                      >{header}</th
                    >
                  {/each}
                </tr>
              </thead>

              <tbody>
                {#each samples.slice(1) as row, rowIndex (rowIndex)}
                  <tr>
                    {#each row as cell, cellIndex (cellIndex)}
                      <td class="whitespace-nowrap">{cell}</td>
                    {/each}
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else if samples && samples.length === 1}
          <p class="text-sm opacity-60">No data available in this file.</p>
        {:else}
          <p class="text-sm opacity-60 italic">Loading...</p>
        {/if}
      </div>
    </div>
    <button class="btn my-6 btn-primary btn-xs w-full" onclick={copyToClipboard}>{copyText}</button>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button onclick={() => (popupOpen = false)}>close</button>
  </form>
</dialog>
