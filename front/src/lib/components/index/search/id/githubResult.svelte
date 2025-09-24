<script lang="ts">
  import Accordion from "$src/lib/components/accordion.svelte";
  import Table from "$src/lib/components/table.svelte";
  import type { GithubResult } from "$src/lib/types";
    import { FlattenObject } from "$src/lib/utils";
  import {
    Building,
    ExternalLink,
    GitCommitVertical,
    Handshake,
    Key,
    Mail,
    UserRoundPen,
  } from "@lucide/svelte";
  const { githubResult }: { githubResult: GithubResult } = $props();
</script>

{#if githubResult.UsernameResult}
  <div class="w-full">
    <div class="flex flex-wrap gap-5">
      <div class="avatar">
        <div class="w-24 h-24 rounded-xl">
          <img
            src={githubResult.UsernameResult.User.AvatarURL}
            alt="Avatar of {githubResult.UsernameResult.User.Username}"
          />
        </div>
      </div>
      <div class="flex flex-col gap-2">
        <div class="flex flex-col">
          <h3 class="h3">{githubResult.UsernameResult.User.Name}</h3>
          <p class="text-base-content/60">
            @{githubResult.UsernameResult.User.Username}
          </p>
        </div>
        <p class="max-w-sm">{githubResult.UsernameResult.User.Bio}</p>
      </div>
    </div>
    <div class="card card-border border-neutral shadow my-8">
      <div class="grid">
        <Table
          row={{
            publicRepos: githubResult.UsernameResult.User.PublicRepos,
            followers: githubResult.UsernameResult.User.Followers,
            following: githubResult.UsernameResult.User.Following,
            createdAt: new Date(
              githubResult.UsernameResult.User.CreatedAt,
            ).toLocaleDateString(),
            email: githubResult.UsernameResult.User.Email,
            location: githubResult.UsernameResult.User.Location,
            company: githubResult.UsernameResult.User.Company,
            url:
              "https://github.com/" + githubResult.UsernameResult.User.Username,
          }}
        />
      </div>
    </div>
    {#if githubResult.UsernameResult.Socials && githubResult.UsernameResult.Socials.length > 0}
      <div class="mt-4">
        <h4 class="h4 mb-2">Social Links</h4>
        <ul class="flex gap-4 flex-col mt-4 mb-6">
          {#each githubResult.UsernameResult.Socials as social}
            <a href={social.URL} target="_blank" rel="noopener noreferrer">
              <div class="badge bg-base-300">
                <ExternalLink size={12} />
                {social.URL}
              </div>
            </a>
          {/each}
        </ul>
      </div>
    {/if}
    {#if githubResult.UsernameResult.CloseFriends && githubResult.UsernameResult.CloseFriends.length > 0}
      <div class="mt-4">
        <ul class="list bg-base-100 rounded-box shadow-md">
            <Accordion
              icon={Handshake}
              title={"Close Friends"}
              subtitle={ githubResult.UsernameResult.CloseFriends.length + " close friends found"}
            >
              <Table
                row={githubResult.UsernameResult.CloseFriends}
              />
            </Accordion>
        </ul>
      </div>
    {/if}
    {#if githubResult.UsernameResult.Orgs && githubResult.UsernameResult.Orgs.length > 0}
      <div class="mt-4">
        <h4 class="h4 mb-2">Organizations</h4>
        <ul class="list bg-base-100 rounded-box shadow-md">
            <Accordion
              icon={Building}
              title="Organizations"
              subtitle={"Found " + githubResult.UsernameResult.Orgs.length + " organizations"}
            >
              <Table
                row={githubResult.UsernameResult.Orgs}
              />
            </Accordion>
        </ul>
      </div>
    {/if}
    {#if githubResult.UsernameResult.Commits && githubResult.UsernameResult.Commits.length > 0}
      <div class="mt-4">
        <h4 class="h4 mb-2">Commits</h4>
        <ul class="list bg-base-100 rounded-box shadow-md">
          {#each githubResult.UsernameResult.Commits as commit}
            <Accordion
              icon={GitCommitVertical}
              title={commit.Name + " <" + commit.Email + ">"}
              subtitle={"Occurrences: " + commit.Occurrences}
            >
              <Table
                row={{
                  name: commit.Name,
                  email: commit.Email,
                  url: "https://github.com/" + commit.FirstFoundIn,
                  occurrences: commit.Occurrences,
                }}
              />
            </Accordion>
          {/each}
        </ul>
      </div>
    {/if}
    {#if githubResult.UsernameResult.SshKeys && githubResult.UsernameResult.SshKeys.length > 0}
      <div class="mt-4">
        <h4 class="h4 mb-2">SSH Keys</h4>
        <ul class="list bg-base-100 rounded-box shadow-md">
          {#each githubResult.UsernameResult.SshKeys as key}
            <Accordion
              icon={Key}
              title={"Created At: " +
                new Date(key.CreatedAt).toLocaleDateString()}
              subtitle={"Last Used: " +
                (key.LastUsed !== "0001-01-01 00:00:00 +0000 UTC"
                  ? new Date(key.LastUsed).toLocaleDateString()
                  : "Never")}
            >
              <pre class="overflow-x-auto p-2 bg-base-200 rounded"><code
                  class="break-all">{key.Key}</code
                ></pre>
            </Accordion>
          {/each}
        </ul>
      </div>
    {/if}
    {#if githubResult.UsernameResult.SshSigningKeys && githubResult.UsernameResult.SshSigningKeys.length > 0}
      <div class="mt-4">
        <h4 class="h4 mb-2">SSH Signing Keys</h4>
        <ul class="list bg-base-100 rounded-box shadow-md">
          {#each githubResult.UsernameResult.SshSigningKeys as key}
            <Accordion
              icon={Key}
              title={key.Title}
              subtitle={"Created At: " + key.CreatedAt}
            >
              <pre class="overflow-x-auto p-2 bg-base-200 rounded"><code
                  class="break-all">{key.Key}</code
                ></pre>
            </Accordion>
          {/each}
        </ul>
      </div>
    {/if}
    {#if githubResult.UsernameResult.GpgKeys && githubResult.UsernameResult.GpgKeys.length > 0}
      <div class="mt-4">
        <h4 class="h4 mb-2">GPG Keys</h4>
        <ul class="list bg-base-100 rounded-box shadow-md">
          {#each githubResult.UsernameResult.GpgKeys as key}
            <Accordion
              icon={Key}
              title={key.Emails && key.Emails.length > 0 ? key.Emails[0].Email : key.KeyID}
              subtitle={"Created At: " + key.CreatedAt}
            >
              <Table
                row={FlattenObject(key)}
              />
            </Accordion>
          {/each}
        </ul>
      </div>
    {/if}
    {#if githubResult.UsernameResult.DeepScan}
      {#if githubResult.UsernameResult.DeepScan.Authors && githubResult.UsernameResult.DeepScan.Authors.length > 0}
        <div class="mt-4">
          <h4 class="h4 mb-2">Deep scan authors</h4>
          <ul class="list bg-base-100 rounded-box shadow-md">
              <Accordion
                icon={UserRoundPen}
                title="Authors"
                subtitle={"Found " + githubResult.UsernameResult.DeepScan.Authors.length + " authors"
                }
              >
                <Table
                  row={githubResult.UsernameResult.DeepScan.Authors}
                />
              </Accordion>
          </ul>
        </div>
      {/if}
      {#if githubResult.UsernameResult.DeepScan.Emails && githubResult.UsernameResult.DeepScan.Emails.length > 0}
        <div class="mt-4">
          <h4 class="h4 mb-2">Deep scan emails</h4>
          <ul class="list bg-base-100 rounded-box shadow-md">
              <Accordion
                icon={Mail}
                title="Emails"
                subtitle={"Found " + githubResult.UsernameResult.DeepScan.Emails.length + " emails"
                }
              >
                <Table
                  row={githubResult.UsernameResult.DeepScan.Emails}
                />
              </Accordion>
          </ul>
        </div>
      {/if}
      {#if githubResult.UsernameResult.DeepScan.Secrets && githubResult.UsernameResult.DeepScan.Secrets.length > 0}
        {@const flattenedSecrets = githubResult.UsernameResult.DeepScan.Secrets.map(FlattenObject)}
        <div class="mt-4">
          <h4 class="h4 mb-2">Deep scan secrets</h4>
          <ul class="list bg-base-100 rounded-box shadow-md">
              <Accordion
                icon={Mail}
                title="Secrets"
                subtitle={"Found " + githubResult.UsernameResult.DeepScan.Secrets.length + " secrets"
                }
              >
                <Table
                  row={flattenedSecrets}
                />
              </Accordion>
          </ul>
        </div>
      {/if}
    {/if}
  </div>
{:else if githubResult.EmailResult}
  <div class="w-full">
    {#if githubResult.EmailResult.Spoofing}
      <h4 class="h4 mb-4">From spoofing</h4>
      <div class="flex flex-wrap gap-5">
        <div class="avatar">
          <div class="w-24 h-24 rounded-xl">
            <img
              src={githubResult.EmailResult.Spoofing.AvatarURL}
              alt="Avatar of {githubResult.EmailResult.Spoofing.Username}"
            />
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <div class="flex flex-col gap-2">
            <h4 class="h4">@{githubResult.EmailResult.Spoofing.Username}</h4>
            {#if githubResult.EmailResult.Spoofing.Name}
              <p>
                <strong>Name:</strong>
                {githubResult.EmailResult.Spoofing.Name}
              </p>
            {/if}
            {#if githubResult.EmailResult.Spoofing.Email}
              <p>
                <strong>Public email:</strong>
                {githubResult.EmailResult.Spoofing.Email}
              </p>
            {/if}
            {#if githubResult.EmailResult.Target}
              <p class="break-all">
                <strong>Primary email:</strong>
                {githubResult.EmailResult.Target}
              </p>
            {/if}
            <a
              href={githubResult.EmailResult.Spoofing.Url}
              class="link link-primary flex gap-2 items-center"
              target="_blank"
            >
              {githubResult.EmailResult.Spoofing.Url}
              <ExternalLink size={12} />
            </a>
          </div>
        </div>
      </div>
    {/if}
    {#if githubResult.EmailResult.Commits}
      <div class="mt-4">
        <h4 class="h4 mb-2">Commits</h4>
        <ul class="list bg-base-100 rounded-box shadow-md">
          {#each githubResult.EmailResult.Commits as commit}
            <Accordion
              icon={GitCommitVertical}
              title={commit.Username && commit.Username !== ""
                ? commit.Name + " (@" + commit.Username + ")"
                : commit.Name}
              subtitle={"Occurrences: " + commit.Occurrences}
            >
              <Table
                row={{
                  name: commit.Name,
                  username: commit.Username,
                  email: commit.Email,
                  first_found_in: commit.FirstFoundIn,
                  occurrences: commit.Occurrences,
                }}
              />
            </Accordion>
          {/each}
        </ul>
      </div>
    {/if}
  </div>
{/if}
