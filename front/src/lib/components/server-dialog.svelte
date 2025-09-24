<script lang="ts">
  import { Key, Link, RefreshCw, Server } from "@lucide/svelte";
  import { cn } from "../utils";
  import { serverUrl, serverPassword } from "$lib/stores/server";
  import { toast } from "svelte-sonner";
  import axios from "axios";

  let { text = "", class: className = "" } = $props();

  let isModalOpen = $state(false);
  let needToTest = $state(true);

  let url = $state($serverUrl || "https://");
  let password = $state($serverPassword);

  let working = $state<boolean | null>(null);

  function save() {
    isModalOpen = false;
    $serverUrl = url;
    $serverPassword = password;
    toast.success("Server settings saved!");
  }

  function test() {
    axios
      .get(`${url}/`)
      .then(() => {
        toast.success("Server is working!");
        needToTest = false;
        working = true;
      })
      .catch(() => {
        toast.error("Server is not working!");
        needToTest = true;
        working = false;
      });
  }

  function reset() {
    $serverUrl = "";
    $serverPassword = "";
    url = "https://";
    password = "";
    needToTest = true;
    working = null;
  }

  $effect(() => {
    if (isModalOpen) {
      url = $serverUrl || "https://";
      needToTest = true;
      working = null;
    }
  });
</script>

<div class="indicator">
  <span class="indicator-item">
    <div class="inline-grid *:[grid-area:1/1]">
      {#if $serverUrl !== ""}
        <div class="status status-success"></div>
      {:else}
        <div class="status status-error animate-ping"></div>
        <div class="status status-error"></div>
      {/if}
    </div>
  </span>
  <button
    onclick={() => {
      isModalOpen = !isModalOpen;
    }}
    class={cn(className, "btn btn-ghost btn-primary")}
  >
    <Server size={16} />
    {text}
  </button>
</div>

<dialog
  class="modal modal-bottom sm:modal-middle"
  class:modal-open={isModalOpen}
>
  <div class="modal-box">
    <form method="dialog">
      <button
        onclick={() => (isModalOpen = false)}
        class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button
      >
    </form>

    <div class="flex flex-col gap-5">
      <div>
        <h2 class="card-title">Connect to your server</h2>
        <p>
          You can connect to your own Eleakxir server by providing the server
          URL and an optional password.
        </p>
      </div>

      <label
        class="input w-full"
        class:input-error={working === false}
        class:input-success={working === true}
      >
        <Link size={16} />
        <input
          class="grow"
          type="url"
          required
          placeholder="https://"
          bind:value={url}
        />

        <button class="btn btn-xs btn-square btn-ghost" onclick={reset}
          ><RefreshCw size={8} /></button
        >
      </label>

      <label class="input w-full">
        <Key />
        <input
          type="password"
          class="grow"
          placeholder="Password"
          bind:value={password}
        />
        <span class="badge badge-neutral badge-xs">Optional</span>
      </label>

      <div class="card-actions flex gap-2">
        <button onclick={test} class="btn btn-primary btn-outline">Test</button>
        <button
          onclick={save}
          disabled={needToTest}
          class="btn btn-primary grow">Save</button
        >
      </div>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button onclick={() => (isModalOpen = false)}>close</button>
  </form>
</dialog>
