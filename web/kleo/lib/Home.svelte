<script lang="ts">
  // import { pb } from './pb_store.svelte';
  import Marks from './Marks.svelte';
  import PublicTimetable from './PublicTimetable.svelte';

  let page = $state(0);
</script>

<button onclick={() => { page = 0; }}>Home</button>
<button onclick={() => { page = 1; }}>Marks</button>
<button onclick={() => { page = 2; }}>PublicTimetable</button>
<br>
{#if page == 0}
  cs lil bro
{:else if page == 1}
  <Marks />
{:else if page == 2}
  <PublicTimetable />
{/if}
