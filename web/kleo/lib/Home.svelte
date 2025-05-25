<script lang="ts">
  // import { pb } from './pb_store.svelte';
  import Marks from './Marks.svelte';
  import PublicTimetable from './PublicTimetable.svelte';
  import DayOverview from './DayOverview.svelte';
  import Events from './Events.svelte';
  import Settings from './Settings.svelte';
  import MyTimeTable from './MyTimeTable.svelte';

  let page = $state(0);

  async function loadXKCD(): Promise<string> {
    let resp = await fetch("https://skrat.org/api/kleo/xkcd");
    return await resp.text();
  }
</script>

<button onclick={() => { page = 0; }}>Home</button>
<button onclick={() => { page = 1; }}>Marks</button>
<button onclick={() => { page = 2; }}>PublicTimetable</button>
<button onclick={() => { page = 3; }}>DayOverview</button>
<button onclick={() => { page = 4; }}>Events</button>
<button onclick={() => { page = 5; }}>MyTimeTable</button>
<button onclick={() => { page = 6; }}>Settings</button>
<br> <br> <br>
{#if page == 0}
  {#await loadXKCD()}
    <h1>loading</h1>
  {:then value}
    <img src={value} alt="">
  {:catch error}
    <h1>error: {error}</h1>
  {/await}
{:else if page == 1}
  <Marks />
{:else if page == 2}
  <PublicTimetable />
{:else if page == 3}
  <DayOverview />
{:else if page == 4}
  <Events />
{:else if page == 5}
  <MyTimeTable />
{:else if page == 6}
  <Settings />
{/if}
