<script lang="ts">
  import AnimatedBeamMultiple from "$src/lib/components/index/AnimatedBeamMultiple.svelte";
  import Logo from "$src/lib/components/logo.svelte";
  import ServerDialog from "$src/lib/components/server-dialog.svelte";
  import { serverUrl } from "$src/lib/stores/server";
  import { ArrowRight, Github, Search } from "@lucide/svelte";
  import { onMount } from "svelte";

  let open = $state(false);

  onMount(() => {
    open = true;
  });
</script>

<header
  class="min-h-[80vh] relative flex justify-center items-center px-6 py-10"
>
  <div
    class="absolute top-0 left-0 z-[-10] w-full h-full bg-top transition-opacity duration-[2000ms]"
    class:opacity-0={!open}
    class:opacity-100={open}
    style="
    background-image:
      linear-gradient(to bottom, rgba(255,255,255,0) 50%, var(--color-base-100) 100%),
      url('https://lovable.dev/img/background/gradient-optimized.svg');
  "
  ></div>
  <div class="mx-auto max-w-3xl flex gap-8 flex-col">
    <a href="https://github.com/anotherhadi/eleakxir" target="_blank">
      <span class="badge badge-lg hover:opacity-90"
        >✨ Check the Github repo <ArrowRight size={16} /></span
      >
    </a>
    <div class="flex gap-6 items-center">
      <Logo size={46} class="fill-primary" />
      <h1 class="font-bold text-7xl">Eleakxir</h1>
    </div>
    <p>
      Eleakxir is a self-hosted search engine that lets you connect to your own
      private and secure server, explore data wells (parquet files) from
      multiple sources, and visualize results in a clean, modern web interface.
    </p>
    <div class="flex gap-6 items-center">
      <a href="/search">
        <button class="btn btn-primary">
          <Search size={16} />
          Let's search</button
        >
      </a>
      <ServerDialog text="Connect to my server" />
    </div>
  </div>
</header>

<main class="flex flex-col gap-24 max-w-7xl m-auto mt-10">
  <div class="card card-dash bg-base-300">
    <div class="card-body flex flex-col gap-10 lg:flex-row">
      <div class="flex gap-5 flex-col">
        <h2 class="card-title text-3xl">⚙️ How Eleakxir works?</h2>
        <p>
          You run an Elixir server that manages parquet files from various
          leaked data sources and multiple OSINT tools. The web client connects
          to your server via HTTPS and authenticated headers then you can search
          across indexed leaks and OSINT tools, browse results interactively and
          review history and stats
          <br />
          <br />
          And it's open source!
        </p>

        <div class="flex items-center gap-2">
          {#if $serverUrl === "https://" || $serverUrl === ""}
            <ServerDialog
              text="Connect your server"
              class="grow btn-outline btn btn-accent btn-sm"
            />
          {/if}
          <a href="https://github.com/anotherhadi/eleakxir">
            <button class="btn btn-outline btn-sm hover:bg-base-200 grow"
              ><Github size={16} /> Check the Github repo</button
            >
          </a>
        </div>
      </div>
      <AnimatedBeamMultiple />
    </div>
  </div>

  <div>
    <h2 class="text-3xl font-bold text-center mb-10">🚀 Features</h2>
    <div
      class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 justify-center m-auto gap-5"
    >
      {#each [{ title: "🔐 Private by design", content: "connect to your own Eleakxir server with a custom URL + password." }, { title: "🛠 Open source & extensible", content: "hack it, self-host it, extend it." }, { title: "📁 Efficient File Format", content: "Uses the columnar Parquet format for high compression and rapid query performance." }, { title: "🔍 OSINT Tools", content: "Includes Github-recon, GHunt, sherlock and more." }, { title: "📜 Standardized Schema", content: "Includes a detailed guide on how to normalize your data leaks for consistent and effective searching across different breaches." }] as value}
        <div class="card bg-base-200 shadow-sm">
          <div class="card-body">
            <h2 class="card-title">{value.title}</h2>
            <p>
              {value.content}
            </p>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <div>
    <h2 class="text-3xl font-bold text-center mb-10">🐢 Speed</h2>
    <p class="max-w-2xl m-auto">
      While Eleakxir is designed to be storage-efficient rather than
      lightning-fast, searches will naturally take longer compared to an indexed
      engine like Elasticsearch. Indexing systems can provide near-instant
      results, but at the cost of massive disk usage — often requiring multiple
      terabytes even for relatively modest datasets. In contrast, Eleakxir
      trades some speed for compactness: for example, I’m able to store 25
      billion rows in just over 600 GB on entry-level hardware. A query might
      take around an hour to complete, but the key point is that it’s actually
      possible to run such searches at home — something that would be completely
      out of reach if I had to maintain Elasticsearch’s much larger index
      footprint.
    </p>
  </div>

  <div>
    <h2 class="text-3xl font-bold text-center mb-10">🚨 Disclaimer</h2>
    <p class="max-w-lg m-auto">
      Eleakxir is provided for educational and research purposes only. You are
      solely responsible for how you use this software. Accessing, storing, or
      distributing leaked data may be illegal in your jurisdiction. The authors
      and contributors do not condone or promote illegal activity. Use
      responsibly and only with data you are legally permitted to process.
    </p>
  </div>
</main>

<div class="pb-24"></div>
