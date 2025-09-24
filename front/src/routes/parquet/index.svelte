<script lang="ts">
  import axios from "axios";
  import { marked } from "marked";
  import { onMount } from "svelte";
  const url = "https://raw.githubusercontent.com/anotherhadi/eleakxir/refs/heads/main/leak-utils/DATALEAKS-NORMALIZATION.md"

  let text = $state<string>("");

  onMount(() => {
    axios.get(url).then((r) => {
      text = r.data;
    });
  });
</script>

<main class="prose max-w-7xl pt-16 pb-28">
  {@html marked(text)}
</main>
