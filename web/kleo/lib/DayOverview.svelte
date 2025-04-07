<script lang="ts">
  import { pb } from "../pb_store.svelte";
  import STimeTableComp from "./STimeTableComp.svelte";
  
  let dayIdx = $state(0)
  let tt = $state(null)

  async function dayIdxChange() {
    if (dayIdx == 0) return;
    tt = await pb.send(`/api/kleo/daytt/${ttype}`, { query: { day: dayIdx } });
    tt.data = Object.values(tt.data).toSorted((a: { owner: string; }, b: { owner: any; }) => a.owner.localeCompare(b.owner));
  }

  let ttype = $state("Class");

</script>


<select bind:value={ttype} onchange={dayIdxChange}>
  <option value="Class">Class</option>
  <option value="Room">Room</option>
  <option value="Teacher">Teacher</option>
</select>
<select bind:value={dayIdx} onchange={dayIdxChange}>
  <option value={0}></option>
  {#each [...Array(5).keys()].map((x) => x+1) as idx}
    <option value={idx}>{new Date(Date.UTC(2025, 11, idx)).toLocaleDateString('cs-CZ', { weekday: "short" })}</option>
  {/each}
</select>

{#if tt !== null}
  <STimeTableComp ttable={tt} useOwner={true} />
{/if}
