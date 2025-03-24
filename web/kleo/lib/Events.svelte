<script lang="ts">
  import EventComp from "./EventComp.svelte";
  import { pb } from "../pb_store.svelte";

  let evts = $state(null);
  let date = $state(null);
  let useDate = $state(false);
  let teacher = $state("");
  let useTeacher = $state(false);
  let classs = $state("");
  let useClass = $state(false);
  let room = $state("");
  let useRoom = $state(false);
  let student = $state("");
  let useStudent = $state(false);
  let string = $state("");
  let useString = $state(false);
  let etype = $state("all");
  let useCached = $state(false);
  let refreshKey = $state({});

  async function getEvts() {
    evts = await pb.send("/api/kleo/events", {
      query: {
        type: etype,
        ...(useDate) && { "date" : date },
        ...(useString) && { string },
        ...(useClass) && { "class": classs },
        ...(useTeacher) && { teacher },
        ...(useRoom) && { room },
        ...(useStudent) && { student },
        ...(useCached) && { cached: "true" }
      },
    });
    refreshKey = {};
  }
</script>

<input type="checkbox" bind:checked={useDate}>
<input type="date" bind:value={date}>
<br>
{#await pb.send("/api/kleo/websrcs", {})}
  <h1>Loading... 🙄</h1>
{:then srcs}
  <input type="checkbox" bind:checked={useClass}>
  <select bind:value={classs}>
    <option value=""></option>
    {#each srcs.classes as item}
      <option value={item.id}>{item.name}</option>
    {/each}
  </select>
  <br>
  <input type="checkbox" bind:checked={useTeacher}>
  <select bind:value={teacher}>
    <option value=""></option>
    {#each srcs.teachers as item}
      <option value={item.id}>{item.name}</option>
    {/each}
  </select>
  <br>
  <input type="checkbox" bind:checked={useRoom}>
  <select bind:value={room}>
    <option value=""></option>
    {#each srcs.rooms as item}
      <option value={item.id}>{item.name}</option>
    {/each}
  </select>
{:catch e}
  <h1>error {e}</h1>
{/await}
<br>
<input type="checkbox" bind:checked={useStudent}>
<input type="text" bind:value={student} placeholder="student">
<br>
<input type="checkbox" bind:checked={useString}>
<input type="text" bind:value={string}>
<br>
<input type="checkbox" bind:checked={useCached} id="cached-check">
<label for="cached-check">cached</label>
<br>
<br>
<select bind:value={etype}>
  <option value="all">All</option>
  <option value="my">My</option>
  <option value="public">Public</option>
</select>
<button onclick={getEvts}>Send</button>
{#key refreshKey}
{#if evts !== null }
  {#each evts.Events as event}
    <EventComp eventRaw={event} />
  {/each}
{/if}
{/key}
