<script lang="ts">
  import EventComp from "./EventComp.svelte";
  import { pb } from "../pb_store.svelte";
</script>

{#await pb.send("/api/kleo/events", {})}
  <h1>Loading... 🙄</h1>
{:then events}
  {#each events.Events as event}
    <EventComp eventRaw={event} />
  {/each}
{:catch error}
  <h1>error: {error}</h1>
{/await}
