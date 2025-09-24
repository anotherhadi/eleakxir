<script lang="ts">
  import Datawells from "$src/lib/components/index/search/datawells.svelte";
  import History from "$src/lib/components/index/search/history.svelte";
  import HowToSearch from "$src/lib/components/index/search/howToSearch.svelte";
  import Searchbar from "$src/lib/components/index/search/searchbar.svelte";
  import Services from "$src/lib/components/index/search/services.svelte";
  import Stats from "$src/lib/components/index/search/stats.svelte";
  import { serverPassword, serverUrl } from "$src/lib/stores/server";
  import type { Server, History as HistoryT } from "$src/lib/types";
  import axios from "axios";
  import { navigate } from "sv-router/generated";
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";

  let serverInfo = $state<Server | null>(null);
  let history = $state<HistoryT>([]);

  onMount(() => {
    if ($serverUrl === "") {
      toast.error("Please, configure your server first!");
      navigate("/");
      return;
    }
    axios
      .get(`${$serverUrl}/`, {
        headers: {
          "X-Password": $serverPassword,
        },
      })
      .then((r) => {
        serverInfo = r.data;
        console.log(serverInfo);
      })
      .catch((e) => {
        toast.error(
          "Failed to fetch server info. Please, change your server configuration!",
        );
        navigate("/");
      });
    axios
      .get(`${$serverUrl}/history`, {
        headers: {
          "X-Password": $serverPassword,
        },
      })
      .then((r) => {
        history = r.data.History;
      })
      .catch((e) => {
        toast.error("Failed to fetch history");
      });
  });
</script>

<main>
  <header class="flex gap-5 flex-col">
    <h1 class="h1"><span class="text-2xl align-middle">🔍</span> Search</h1>
    <Searchbar />
  </header>

  <div class="my-10"></div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
    <div class="card card-border border-neutral shadow col-span-full">
      <Stats {serverInfo} />
    </div>
    <div class="card card-border border-neutral shadow card-body">
      <h2 class="h2">History</h2>
      <History {history} />
    </div>
    <div class="card card-border border-neutral shadow card-body">
      <h2 class="h2">Active services</h2>
      <div class="overflow-x-auto">
        {#if !serverInfo}
          <p>Loading...</p>
        {:else}
          <Services {serverInfo} />
        {/if}
      </div>
    </div>
    <div class="card card-border border-neutral shadow card-body">
      <h2 class="h2">Last data wells added</h2>
      <div class="overflow-x-auto">
        <Datawells dataleaks={serverInfo?.Dataleaks || []} />
      </div>
    </div>
    <div class="card card-border border-neutral shadow card-body">
      <h2 class="h2">How to search</h2>
      <HowToSearch />
    </div>
  </div>

  <div class="mb-10"></div>
</main>
