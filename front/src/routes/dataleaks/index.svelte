<script lang="ts">
  import Datawells from "$src/lib/components/index/search/datawells.svelte";
  import Stats from "$src/lib/components/index/search/stats.svelte";

  import { serverPassword, serverUrl } from "$src/lib/stores/server";
  import type { Server } from "$src/lib/types";
  import axios from "axios";
  import { navigate } from "sv-router/generated";
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";

  let serverInfo = $state<Server | null>(null);

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
        console.log(e);
        navigate("/");
      });
  });
</script>

<main>
  <header class="flex flex-col gap-2 mb-8">
    <h1 class="h1"><span class="text-2xl align-middle">🗃️</span> Data wells</h1>
    <p>List of data wells (databases) available on the connected server.</p>
  </header>
  {#if serverInfo}
    <div class="card card-border border-neutral shadow col-span-full mb-5">
      <Stats {serverInfo} />
    </div>
    <Datawells
      dataleaks={serverInfo.Dataleaks}
      showColumns={true}
      perPage={20}
    />
  {:else}
    <p>Loading...</p>
  {/if}
</main>
