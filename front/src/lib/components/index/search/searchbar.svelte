<script lang="ts">
  import { serverPassword, serverUrl } from "$src/lib/stores/server";
  import { cn } from "$src/lib/utils";
  import { Equal, EqualNot, Search } from "@lucide/svelte";
  import axios from "axios";
  import { navigate } from "sv-router/generated";
  import { toast } from "svelte-sonner";

  const {
    initialQuery = "",
    initialFilter = "all",
    initialExactMatch = false,
  }: {
    initialQuery?: string;
    initialFilter?: string;
    initialExactMatch?: boolean;
  } = $props();

  let filters = [
    "all",
    "username",
    "email",
    "name",
    "phone",
    "url",
    "password",
    "password hash",
    "full_text",
  ];
  let activeFilter = $state<string>(initialFilter);
  let query = $state<string>(initialQuery);
  let exactMatch = $state<boolean>(initialExactMatch);

  function NewSearch() {
    axios
      .post(
        `${$serverUrl}/search`,
        { Text: query, Column: activeFilter, ExactMatch: exactMatch },
        {
          headers: {
            "Content-Type": "application/json",
            "X-Password": $serverPassword,
          },
        },
      )
      .then((r) => {
        const id = r.data.Id;
        window.location.href = `/search/${id}`;
      })
      .catch((e) => {
        if (e.response.data.Error !== undefined) {
          toast.error(e.response.data.Error);
        } else {
          toast.error("An error occurred");
        }
      });
  }
</script>

<div class="flex gap-5 flex-col">
  <div
    class="flex gap-3 justify-start items-center w-full overflow-y-hidden overflow-x-auto"
  >
    {#each filters as filter}
      <button
        class={cn(
          "btn btn-md capitalize",
          activeFilter === filter
            ? "btn-primary"
            : "btn-ghost btn-neutral text-base-content/80 hover:text-neutral-content",
        )}
        onclick={() => (activeFilter = filter)}>{filter.replace("_", " ")}</button
      >
    {/each}
  </div>

  <form
    class="join w-full"
    onsubmit={(e) => {
      e.preventDefault();
      NewSearch();
    }}
  >
    <label class="grow input input-xl input-primary join-item w-full">
      <Search size={16} />
      <input
        class="grow input-xl"
        type="text"
        bind:value={query}
        placeholder="Search..."
        required
      />

      <div class="tooltip" data-tip="Exact Match">
        <label class="toggle text-base-content toggle-xs">
          <input type="checkbox" bind:checked={exactMatch} />
          <EqualNot aria-label="disable" size={12} />
          <Equal aria-label="enabled" size={12} />
        </label>
      </div>
    </label>
    <button class="btn btn-primary btn-xl join-item">Search</button>
  </form>
</div>
