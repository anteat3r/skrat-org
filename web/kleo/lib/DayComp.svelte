<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "../pb_store.svelte";

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

<h1>{date.toLocaleDateString("cs-CZ")}</h1>
{#if events !== null}
  {#each events as event}
    <p>
      <span style="font-size: 20; font-weight: bold;">{event.Title}</span>
      &nbsp;
      <span style="font-size: 5;">({event.EventType.Name})</span>
      <br>
      <span>{event.Description}</span>
      <br><br>
      {#each event.Times as time}
        <span style="font-weight: bold;">{new Date(time.StartTime).toLocaleDateString("cs-CZ")}</span>
        {#if !time.WholeDay}
          <span>
            {new Date(time.StartTime).toLocaleDateString("cs-CZ")}
            {new Date(time.StartTime).toLocaleTimeString("cs-CZ").split(":").slice(0, -1).join(":")} - 
            {new Date(time.EndTime).toLocaleTimeString("cs-CZ").split(":").slice(0, -1).join(":")}
          </span>
        {/if}
      {/each}
      <br><br>
      {#if event.Classes.length > 0} <span style="font-weight: bold;">Třídy:</span> {/if}
      {event.Classes.map(e => e.Abbrev).join(", ")}
      <br>
      {#if event.Teachers.length > 0} <span style="font-weight: bold;">Učitelé:</span> {/if}
      {event.Teachers.map(e => e.Name).join(", ")}
      <br>
      {#if event.Rooms.length > 0} <span style="font-weight: bold;">Místnosti:</span> {/if}
      {event.Rooms.map(e => e.Abbrev).join(", ")}
      <br>
      {#if event.Students.length > 0} <span style="font-weight: bold;">Studenti:</span> {/if}
      {event.Students.map(e => e.Abbrev).join(", ")}
    </p>
  {/each}
{/if}
