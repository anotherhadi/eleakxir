<script lang="ts">
  import { ExternalLink } from "@lucide/svelte";

  const {
    row,
  }: { row: Record<string, string> | Array<Record<string, string>>| string[][] } = $props();
</script>

<div class="overflow-x-auto">
  <table class="table">
    {#if Array.isArray(row) && row.length > 0 && row[0]}
      {@const head = Object.entries(row[0])}
      <!-- head -->
      <thead>
        <tr>
          {#each head as [key, _]}
            <th
              class="text-xs whitespace-nowrap font-semibold opacity-60 capitalize"
            >
              {key}
            </th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each row as item}
          <tr>
            {#each Object.entries(item) as [key, value]}
              <th class="text-xs whitespace-nowrap font-semibold opacity-60">
                {#if (key.toLowerCase() == "url" || key
                    .toLowerCase()
                    .endsWith("_url")) && value !== null && value !== ""}
                  <a
                    href={value}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="link link-primary gap-2 items-center flex"
                  >
                    {value}
                    <ExternalLink size={12} />
                  </a>
                {:else}
                  {value}
                {/if}
              </th>
            {/each}
          </tr>
        {/each}
      </tbody>
    {:else if row && Object.keys(row).length > 0}
      <tbody>
        {#each Object.entries(row) as [key, value]}
          {#if key !== "source" && value !== null && value !== ""}
            <tr class="">
              <th
                class="text-xs whitespace-nowrap font-semibold opacity-60 capitalize"
                >{key.replace(/_/g, " ")}</th
              >

              <td class="w-fit overflow-x-auto whitespace-nowrap">
                {#if key.toLowerCase() == "url" || key
                    .toLowerCase()
                    .endsWith("_url")}
                  <a
                    href={value}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="link link-primary gap-2 items-center flex"
                  >
                    {value}
                    <ExternalLink size={12} />
                  </a>
                {:else}
                  {value}
                {/if}
              </td>
            </tr>
          {/if}
        {/each}
      </tbody>
    {/if}
  </table>
</div>
