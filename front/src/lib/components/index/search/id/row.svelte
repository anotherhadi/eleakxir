<script lang="ts">
  import Table from "$src/lib/components/table.svelte";
  import { ChevronDown, ChevronUp, Database, Key, Mail } from "@lucide/svelte";

  const { row }: { row: Record<string, string> } = $props();

  let isOpen = $state<boolean>(false);

  function getDomain(dataleakName: string) {
    if (!dataleakName) return null;
    const firstPart = dataleakName.split(" ")[0].toLowerCase();
    const domainRegex =
      /^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$/;
    if (domainRegex.test(firstPart)) {
      return firstPart.replace("_", "-");
    }
    return null;
  }

  function getHighlightedContent(row: Record<string, string>): string {
    const prioritizedKeys = [
      "email",
      "username",
      "full_name",
      "first_name",
      "last_name",
      "phone",
      "password",
      "address",
    ];

    for (const key of prioritizedKeys) {
      if (row[key]) {
        return row[key];
      }
    }

    for (const key in row) {
      if (row[key]) {
        return row[key];
      }
    }

    return "No content";
  }
</script>

<button
  class="list-row hover:bg-base-300/75 text-left"
  class:bg-base-300={isOpen}
  class:rounded-b-none={isOpen}
  onclick={() => {
    isOpen = !isOpen;
  }}
>
  <div>
    {#if getDomain(row["source"])}
      <img
        src="https://icons.duckduckgo.com/ip3/{getDomain(row['source'])}.ico"
        class="size-10 rounded-box bg-neutral"
        alt="Favicon de {getDomain(row['source'])}"
      />
    {:else if row["password"] !== null}
      <div
        class="size-10 rounded-box bg-neutral items-center justify-center flex"
      >
        <Key class="text-neutral-content" />
      </div>
    {:else if row["email"] !== null}
      <div
        class="size-10 rounded-box bg-neutral items-center justify-center flex"
      >
        <Mail class="text-neutral-content" />
      </div>
    {:else}
      <div
        class="size-10 rounded-box bg-neutral items-center justify-center flex"
      >
        <Database class="text-neutral-content" />
      </div>
    {/if}
  </div>
  <div
    class="flex-1 flex flex-col min-w-0 items-start justify-center"
  >
    <div class="w-full overflow-hidden whitespace-nowrap text-ellipsis">
      {getHighlightedContent(row)}
    </div>
    <div class="text-xs uppercase font-semibold opacity-60">
      {row["source"]}
    </div>
  </div>
  <div class="btn btn-square btn-ghost">
    {#if isOpen}
      <ChevronUp size={12} />
    {:else}
      <ChevronDown size={12} />
    {/if}
  </div>
</button>
{#if isOpen}
  <li class="grid list-row bg-base-200 rounded-t-none mb-2">
    <Table {row} />
  </li>
{/if}
