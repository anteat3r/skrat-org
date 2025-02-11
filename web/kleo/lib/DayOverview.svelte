<script lang="ts">
  import { pb } from "../pb_store.svelte";
  
  let dayIdx = $state(0)
  let tt = $state(null)

  async function dayIdxChange() {
    if (dayIdx == 0) return;
    tt = await pb.send(`/api/kleo/daytt`, { query: { day: dayIdx } });
  }
</script>

<select bind:value={dayIdx} onchange={dayIdxChange}>
  <option value={0}></option>
  {#each [...Array(5).keys()].map((x) => x+1) as idx}
    <option value={idx}>{idx}</option>
  {/each}
</select>

<pre>{JSON.stringify(tt)}</pre>
