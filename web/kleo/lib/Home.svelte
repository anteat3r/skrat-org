<script lang="ts">
  // import { pb } from './pb_store.svelte';
  import Marks from './Marks.svelte';
  import PublicTimetable from './PublicTimetable.svelte';
  import DayOverview from './DayOverview.svelte';
  import Events from './Events.svelte';
  import Settings from './Settings.svelte';

  let page = $state(0);
</script>

<button onclick={() => { page = 0; }}>Home</button>
<button onclick={() => { page = 1; }}>Marks</button>
<button onclick={() => { page = 2; }}>PublicTimetable</button>
<button onclick={() => { page = 3; }}>DayOverview</button>
<button onclick={() => { page = 4; }}>Events</button>
<button onclick={() => { page = 5; }}>Settings</button>
<br> <br> <br>
{#if page == 0}
  cs lil bro
{:else if page == 1}
  <Marks />
{:else if page == 2}
  <PublicTimetable />
{:else if page == 3}
  <DayOverview />
{:else if page == 4}
  <Events />
{:else if page == 5}
  <Settings />
{/if}
