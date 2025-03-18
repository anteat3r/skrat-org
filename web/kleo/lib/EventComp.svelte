<script lang="ts">
  let { eventRaw } = $props();
  let event = eventRaw as any;

  let expanded = $state(false);

  function btnClicked() {
    expanded = !expanded;
  }
</script>

<div class="event-border">
  <button onclick={btnClicked}>{event.Title}</button>
  {#if expanded}
    <p>
      <span>{event.Description}</span>
      {#if event.Times.length > 0 }
        <br><br>
      {/if}
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
      {event.Classes.map((e: { Abbrev: any; }) => e.Abbrev).join(", ")}
      <br>
      {#if event.Teachers.length > 0} <span style="font-weight: bold;">Učitelé:</span> {/if}
      {event.Teachers.map((e: { Name: any; }) => e.Name).join(", ")}
      <br>
      {#if event.Rooms.length > 0} <span style="font-weight: bold;">Místnosti:</span> {/if}
      {event.Rooms.map((e: { Abbrev: any; }) => e.Abbrev).join(", ")}
      <br>
      {#if event.Students.length > 0} <span style="font-weight: bold;">Studenti:</span> {/if}
      {event.Students.map((e: { Name: any; }) => e.Name).join(", ")}
    </p>
  {/if}
</div>

<style>
  .event-border {
    border: 1px solid white;
    padding: 10px;
  }
</style>
