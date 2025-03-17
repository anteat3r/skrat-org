<script lang="ts">
  import { onMount } from "svelte";
  import { pb } from "../pb_store.svelte";
  import EventComp from "./EventComp.svelte";

  interface Props {
    date: Date,
    type: "personal" | "class" | "teacher" | "room",
    events: any[] | null,
  }
  let { date, type = "personal", events = null }: Props = $props();

  onMount(async function() {
    if (events !== null) return;
    events = await pb.send("/api/kleo/events", {});
  })
</script>

<div class="event-border">
  <h1>{date.toLocaleDateString("cs-CZ")}</h1>
  {#if events !== null}
      {#each events as event}
        <EventComp eventRaw={event}/>
      {/each}
  {/if}
</div>

<style>
  .event-border {
    border: solid 1px white;
    padding: 10px;
  }
</style>
