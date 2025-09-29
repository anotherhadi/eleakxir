<script lang="ts">
  import type { Server } from "$src/lib/types";

  let { serverInfo }: { serverInfo: Server } = $props();
</script>

<div class="my-4">
  <div class="overflow-x-auto">
    <table class="table">
      <!-- head -->
      <thead>
        <tr>
          <th>Service</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr class="hover:bg-base-300">
          <th class="flex flex-wrap gap-2 items-center">
            Data wells lookup
            <div class="badge badge-xs badge-neutral">Limit: {serverInfo.Settings.Limit}</div>
          </th>
          <td class="text-nowrap">
            {#if serverInfo.Dataleaks.length !== 0}
              <div class="inline-grid *:[grid-area:1/1] mr-2">
                <div class="status status-success"></div>
                <div class="status status-success"></div>
              </div>
              Active
            {:else}
              <div class="inline-grid *:[grid-area:1/1] mr-2">
                <div class="status status-error animate-ping"></div>
                <div class="status status-error"></div>
              </div>
              Inactive
            {/if}
          </td>
        </tr>
        <tr class="hover:bg-base-300">
          <th class="flex flex-wrap gap-2 items-center">
            Github recon
            {#if serverInfo.Settings.GithubTokenLoaded === true}
              <div class="badge badge-xs badge-neutral">Token</div>
            {/if}
            {#if serverInfo.Settings.GithubDeepMode === true}
              <div class="badge badge-xs badge-neutral">Deep Mode</div>
            {/if}
          </th>
          <td class="text-nowrap">
            {#if serverInfo.Settings.GithubRecon === true}
              <div class="inline-grid *:[grid-area:1/1] mr-2">
                <div class="status status-success"></div>
                <div class="status status-success"></div>
              </div>
              Active
            {:else}
              <div class="inline-grid *:[grid-area:1/1] mr-2">
                <div class="status status-error animate-ping"></div>
                <div class="status status-error"></div>
              </div>
              Inactive
            {/if}
          </td>
        </tr>
        <tr class="hover:bg-base-300">
          <th>Gravatar recon</th>
          <td class="text-nowrap">
            {#if serverInfo.Settings.GravatarRecon === true}
              <div class="inline-grid *:[grid-area:1/1] mr-2">
                <div class="status status-success"></div>
                <div class="status status-success"></div>
              </div>
              Active
            {:else}
              <div class="inline-grid *:[grid-area:1/1] mr-2">
                <div class="status status-error animate-ping"></div>
                <div class="status status-error"></div>
              </div>
              Inactive
            {/if}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
