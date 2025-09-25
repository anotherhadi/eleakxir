<script lang="ts">
  import type { Result } from "$src/lib/types";
  import axios from "axios";
  import { navigate, route } from "sv-router/generated";
  import { serverPassword, serverUrl } from "$src/lib/stores/server";
  import Searchbar from "$src/lib/components/index/search/searchbar.svelte";
  import { toast } from "svelte-sonner";
  import { onMount } from "svelte";
  import Stats from "$src/lib/components/index/search/id/stats.svelte";
  import Rows from "$src/lib/components/index/search/id/rows.svelte";
  import {
    ChevronDown,
    CircleAlert,
    CircleCheck,
    CircleMinus,
    CircleX,
    Database,
    Github,
  } from "@lucide/svelte";
  import { convertNanoSeconds } from "$src/lib/utils";
  import GithubResult from "$src/lib/components/index/search/id/githubResult.svelte";

  route.getParams("/search/:id");

  let { id } = route.params;

  let result = $state<Result | null>(null);

  function loadData() {
    if (id === undefined) {
      return;
    }
    if (id === "") {
      return;
    }
    axios
      .get(`${$serverUrl}/search/${id}`, {
        headers: {
          "X-Password": $serverPassword,
        },
      })
      .then((r) => {
        result = r.data;
        console.log(r.data);
        if (result && result.Status !== "pending") {
          clearInterval(intervalId);
        }
      })
      .catch((e) => {
        toast.error("Failed to fetch search result!");
        clearInterval(intervalId);
        navigate("/search");
      });
  }

  let intervalId: ReturnType<typeof setInterval>;
  let elapsedTime = 0;
  let pollingInterval = 10000; // Start with a 10-second interval

  onMount(() => {
    if ($serverUrl === "") {
      toast.error("Please, configure your server first!");
      navigate("/");
      return;
    }

    loadData();

    intervalId = setInterval(() => {
      elapsedTime += pollingInterval;

      // Check for status change inside the interval
      if (result && result.Status !== "pending") {
        clearInterval(intervalId);
        return;
      }

      // Change polling frequency based on elapsed time
      if (elapsedTime >= 120000 && pollingInterval !== 10000) {
        clearInterval(intervalId);
        pollingInterval = 15000;
        intervalId = setInterval(loadData, pollingInterval);
        return;
      } else if (elapsedTime >= 600000 && pollingInterval !== 30000) {
        clearInterval(intervalId);
        pollingInterval = 30000;
        intervalId = setInterval(loadData, pollingInterval);
        return;
      }

      loadData();
    }, pollingInterval);

    return () => {
      clearInterval(intervalId);
    };
  });
</script>

<main>
  {#if result}
    <header class="flex gap-5 flex-col">
      <a href="/search">
        <h1 class="h1"><span class="text-2xl align-middle">🔍</span> Search</h1>
      </a>

      <Searchbar
        initialQuery={result.Query.Text}
        initialFilter={result.Query.Column}
        initialExactMatch={result.Query.ExactMatch}
        initialDatawells={result.Query.Datawells}
        initialGithubRecon={result.Query.GithubRecon}
      />
    </header>

    <div class="my-10"></div>

    <div class="grid grid-cols-1 gap-5 [&>div]:border-neutral">
      <div class="card card-border shadow col-span-full">
        <Stats {result} />
      </div>

      {#if result.LeakResult.Error !== "not enabled" }
        <div class="collapse collapse-arrow bg-base-100 border">
          <input type="radio" name="my-accordion-2" checked={true} />
          <div
            class="collapse-title font-semibold text-xl flex justify-between items-center"
          >
            <div class="flex items-center gap-2">
              <Database size={18} class="text-base-content/60" />
              Data wells lookup
            </div>
            {#if result.LeakResult.Error !== ""}
              <CircleX size={16} class="text-error" />
            {:else if result.LeakResult.Duration === 0}
              <span class="loading loading-dots loading-xs"></span>
            {:else if result.LeakResult.Rows.length > 0}
              <CircleCheck size={16} class="text-success" />
            {:else}
              <CircleMinus size={16} class="text-base-content/60" />
            {/if}
          </div>
          <div class="collapse-content">
            {#if result.LeakResult.Error !== ""}
              <div role="alert" class="alert alert-soft alert-error">
                <CircleAlert size={20} />
                <span>Error! {result.LeakResult.Error}</span>
              </div>
            {:else if result.LeakResult.Duration === 0}
              <ul class="list rounded-box">
                {#each Array(5) as _}
                  <div class="list-row text-left">
                    <div>
                      <div
                        class="skeleton size-10 rounded-box items-center justify-center flex"
                      ></div>
                    </div>
                    <div>
                      <div class="skeleton h-5 mb-1 w-52"></div>
                      <div
                        class="text-xs skeleton h-4 w-34 uppercase font-semibold opacity-60"
                      ></div>
                    </div>
                    <div class="btn btn-square btn-ghost">
                      <ChevronDown size={12} />
                    </div>
                  </div>
                {/each}
              </ul>
            {:else}
              <p class="text-base-content/60">
                {result.LeakResult.Rows.length} results in {convertNanoSeconds(
                  result.LeakResult.Duration,
                )}
              </p>
              {#if result.LeakResult.LimitHit}
                <div role="alert" class="alert alert-soft my-4">
                  <CircleAlert size={20} />
                  <div>
                    <span class="font-semibold">Limit hit!</span> Consider refining
                    your search query for more specific results.
                  </div>
                </div>
              {/if}
              <Rows {result} />
            {/if}
          </div>
        </div>
      {/if}
      {#if result.GithubResult.Error !== "not enabled" }
        <div class="collapse collapse-arrow bg-base-100 border">
          <input type="radio" name="my-accordion-2" />
          <div
            class="collapse-title font-semibold text-xl flex justify-between items-center"
          >
            <div class="flex items-center gap-2">
              <Github size={18} class="text-base-content/60" />
              Github Recon
            </div>
            {#if result.GithubResult.Error !== ""}
              <CircleX size={16} class="text-error" />
            {:else if result.GithubResult.Duration === 0}
              <span class="loading loading-dots loading-xs"></span>
            {:else if !result.GithubResult.EmailResult?.Commits && !result.GithubResult.EmailResult?.Spoofing && !result.GithubResult.UsernameResult?.User}
              <CircleMinus size={16} class="text-base-content/60" />
            {:else if result.GithubResult.UsernameResult || result.GithubResult.EmailResult}
              <CircleCheck size={16} class="text-success" />
            {/if}
          </div>
          <div class="collapse-content">
            {#if result.GithubResult.Error !== ""}
              <div role="alert" class="alert alert-soft alert-error">
                <CircleAlert size={20} />
                <span>Error! {result.GithubResult.Error}</span>
              </div>
            {:else if result.GithubResult.Duration === 0}
              <div role="alert" class="alert alert-soft">
                <span class="loading loading-dots loading-sm"></span>
                <span>Loading...</span>
              </div>
            {:else if !result.GithubResult.EmailResult?.Commits && !result.GithubResult.EmailResult?.Spoofing && !result.GithubResult.UsernameResult?.User}
              <div role="alert" class="alert alert-soft">
                <CircleMinus size={20} />
                <span>No result</span>
              </div>
            {:else}
              <p class="text-base-content/60 mb-4">
                Found a result in {convertNanoSeconds(
                  result.GithubResult.Duration,
                )}
              </p>
              <GithubResult githubResult={result.GithubResult} />
            {/if}
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <div class="mb-10"></div>
</main>
