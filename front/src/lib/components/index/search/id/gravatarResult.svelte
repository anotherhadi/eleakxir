<script lang="ts">
  import Accordion from "$src/lib/components/accordion.svelte";
  import Table from "$src/lib/components/table.svelte";
  import type { GravatarResult } from "$src/lib/types";
  import { Contact, ExternalLink, Mail, Phone } from "@lucide/svelte";

  const { result }: { result: GravatarResult } = $props();
</script>

<div class="w-full flex flex-col gap-10">
  {#each result.Results as r}
    <div class="flex flex-wrap gap-5">
      <div class="avatar">
        <div class="w-24 h-24 rounded-xl">
          <img src={r.thumbnailUrl} alt="Avatar of {r.preferredUsername}" />
        </div>
      </div>
      <div class="flex flex-col gap-2">
        <div class="flex flex-col">
          <h3 class="h3">{r.displayName}</h3>
          <p class="text-base-content/60">
            @{r.preferredUsername}
          </p>
        </div>
        <p class="max-w-sm">{r.aboutMe}</p>
      </div>
    </div>
    <div class="card card-border border-neutral shadow">
      <div class="grid">
        <Table
          row={{
            profile_url: r.profileUrl,
            current_location: r.currentLocation,
            job_title: r.job_title,
            company: r.company,
            pronouns: r.pronouns,
            pronunciation: r.pronunciation,
            photos: r.photos.length > 0 ? r.photos.length : "N/A",
          }}
        />
      </div>
    </div>

    <div class="flex flex-col gap-2">
      {#if r.accounts && r.accounts.length > 0}
        <div>
          <h4 class="h4 mb-2">Social Links</h4>
          <ul class="flex gap-4 flex-col mt-4 mb-6">
            {#each r.accounts as account}
              <a href={account.url} target="_blank" rel="noopener noreferrer">
                <div class="badge bg-base-300">
                  <ExternalLink size={12} />
                  {account.username} ({account.url})
                </div>
              </a>
            {/each}
          </ul>
        </div>
      {/if}

      {#if r.emails && r.emails.length > 0}
        <div>
          <ul class="list bg-base-100 rounded-box shadow-md">
            <Accordion
              icon={Mail}
              title={"Emails"}
              subtitle={r.emails.length + " email found"}
            >
              <Table row={r.emails} />
            </Accordion>
          </ul>
        </div>
      {/if}

      {#if r.phoneNumbers && r.phoneNumbers.length > 0}
        <div>
          <ul class="list bg-base-100 rounded-box shadow-md">
            <Accordion
              icon={Phone}
              title={"Phone Numbers"}
              subtitle={r.phoneNumbers.length + " phone numbers found"}
            >
              <Table row={r.phoneNumbers} />
            </Accordion>
          </ul>
        </div>
      {/if}

      {#if r.contactInfo && r.contactInfo.length > 0}
        <div>
          <ul class="list bg-base-100 rounded-box shadow-md">
            <Accordion
              icon={Contact}
              title={"Contact Info"}
              subtitle={r.contactInfo.length + " contact info found"}
            >
              <Table row={r.contactInfo} />
            </Accordion>
          </ul>
        </div>
      {/if}
    </div>
  {/each}
</div>
