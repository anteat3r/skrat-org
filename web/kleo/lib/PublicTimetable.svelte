<script lang="ts">
  import { pb } from "../pb_store.svelte";
</script>

{#await pb.send("/api/kleo/websrcs", {})}
  <h1>Loading... 🙄</h1>
{:then srcs}
  <h1>{srcs}</h1>
{:catch error}
  <h1>error: {error}</h1>
{/await}
