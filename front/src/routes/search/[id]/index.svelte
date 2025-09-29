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
  import GravatarResult from "$src/lib/components/index/search/id/gravatarResult.svelte";

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
        console.log(result);
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
      <a href="/search" class="w-fit">
        <h1 class="h1 "><span class="text-2xl align-middle">🔍</span> Search</h1>
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

      {#if !result.LeakResult.Inactive}
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
            {:else if !result.LeakResult.Rows || result.LeakResult.Rows.length == 0}
              <div role="alert" class="alert alert-soft">
                <CircleMinus size={20} />
                <span>No result</span>
              </div>
            {:else}
              <p class="text-base-content/60 mb-2">
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
      {#if !result.GithubResult.Inactive}
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
      {#if !result.GravatarResult.Inactive}
        <div class="collapse collapse-arrow bg-base-100 border">
          <input type="radio" name="my-accordion-2" />
          <div
            class="collapse-title font-semibold text-xl flex justify-between items-center"
          >
            <div class="flex items-center gap-2">
              <svg
                width="18"
                height="18"
                viewBox="0 0 18 18"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M7.20008 1.79933V8.09932C7.20008 8.57653 7.38965 9.0342 7.72709 9.37164C8.06453 9.70908 8.5222 9.89865 8.99941 9.89865C9.47662 9.89865 9.93429 9.70908 10.2717 9.37164C10.6092 9.0342 10.7987 8.57653 10.7987 8.09932V3.90799C11.9031 4.29735 12.851 5.03509 13.4996 6.01006C14.1482 6.98502 14.4623 8.14438 14.3947 9.31342C14.327 10.4825 13.8812 11.5978 13.1245 12.4915C12.3678 13.3851 11.3411 14.0086 10.1992 14.2679C9.05725 14.5273 7.86198 14.4084 6.79347 13.9294C5.72497 13.4503 4.84112 12.6369 4.27513 11.6117C3.70914 10.5866 3.49168 9.40529 3.6555 8.24581C3.81933 7.08634 4.35557 6.01152 5.18342 5.18333C5.51545 4.84434 5.70032 4.38803 5.69786 3.91353C5.69541 3.43902 5.50582 2.98465 5.17029 2.64912C4.83476 2.31359 4.38039 2.12401 3.90589 2.12155C3.43138 2.11909 2.97508 2.30396 2.63609 2.636C1.16373 4.10834 0.247437 6.04566 0.043349 8.11786C-0.160739 10.1901 0.360003 12.2689 1.51684 14.0002C2.67368 15.7315 4.39505 17.0081 6.38762 17.6125C8.38019 18.2169 10.5207 18.1117 12.4444 17.3148C14.3681 16.5179 15.956 15.0786 16.9374 13.2422C17.9189 11.4059 18.2333 9.28595 17.827 7.24376C17.4207 5.20156 16.3188 3.36344 14.7091 2.04258C13.0995 0.721724 11.0816 -0.000136192 8.99941 1.92733e-08C8.5222 1.92733e-08 8.06453 0.189572 7.72709 0.527012C7.38965 0.864452 7.20008 1.32212 7.20008 1.79933Z"
                  class="fill-base-content/60"
                />
              </svg>

              Gravatar Recon
            </div>
            {#if result.GravatarResult.Error !== ""}
              <CircleX size={16} class="text-error" />
            {:else if result.GravatarResult.Duration === 0}
              <span class="loading loading-dots loading-xs"></span>
            {:else if !result.GravatarResult.Results || result.GravatarResult.Results.length == 0}
              <CircleMinus size={16} class="text-base-content/60" />
            {:else if result.GravatarResult.Results}
              <CircleCheck size={16} class="text-success" />
            {/if}
          </div>
          <div class="collapse-content">
            {#if result.GravatarResult.Error !== ""}
              <div role="alert" class="alert alert-soft alert-error">
                <CircleAlert size={20} />
                <span>Error! {result.GravatarResult.Error}</span>
              </div>
            {:else if result.GravatarResult.Duration === 0}
              <div role="alert" class="alert alert-soft">
                <span class="loading loading-dots loading-sm"></span>
                <span>Loading...</span>
              </div>
            {:else if !result.GravatarResult.Results || result.GravatarResult.Results.length == 0}
              <div role="alert" class="alert alert-soft">
                <CircleMinus size={20} />
                <span>No result</span>
              </div>
            {:else}
              <p class="text-base-content/60 mb-4">
                Found a result in {convertNanoSeconds(
                  result.GravatarResult.Duration,
                )}
              </p>
              <GravatarResult result={result.GravatarResult} />
            {/if}
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <div class="mb-10"></div>
</main>
