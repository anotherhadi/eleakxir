<script lang="ts">
    import { isActiveLink } from 'sv-router';
  let { item } = $props();
  import Self from './sidebar-menu-item.svelte';
</script>

<li>
  {#if item.items}
    <details open>
      <summary class="flex gap-2 items-center">
      {#if item.icon}
        <item.icon size={16} />
      {/if}
      {item.title}
      </summary>
      <ul>
        {#each item.items as subitem}
          <Self item={subitem} />
        {/each}
      </ul>
    </details>
  {:else}
    <a href={item.url} class="flex gap-2 items-center" use:isActiveLink={{ className: 'menu-active' }}>
      {#if item.icon}
        <item.icon size={16} />
      {/if}
      {item.title}
    </a>
  {/if}
</li>
