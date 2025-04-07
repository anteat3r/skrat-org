<script lang="ts">
  import { pb } from "../pb_store.svelte";
  import STimeTableComp from "./STimeTableComp.svelte";

  function srcChange(ttype: string) {
    return async function() {
      if (src === "") {
        ttable = null;
        return;
      }
      let resp = await pb.send(`/api/kleo/web/${ttime}/${ttype}/${src}`, {});
      ttable = resp;
      last_ttype = ttype;
    }
  }

  async function ttimeChange() {
    if (src === "") {
      ttable = null;
      return;
    }
    if (last_ttype === "") {
      return;
    }
    let resp = await pb.send(`/api/kleo/web/${ttime}/${last_ttype}/${src}`, {});
    ttable = resp;
  }

  let ttable = $state(null);
  let src = $state("");
  let ttime = $state("Actual");
  let last_ttype = $state("");

  const srcts = [
    { name: "classes",
      url: "Class", },
    { name: "teachers",
      url: "Teacher", },
    { name: "rooms",
      url: "Room", },
  ];

</script>

{#await pb.send("/api/kleo/websrcs", {})}
  <h1>Loading... 🙄</h1>
{:then srcs}
  
  <select bind:value={ttime} onchange={ttimeChange}>
    <option value="Actual">Actual</option>
    <option value="Next">Next</option>
    <option value="Permanent">Permanent</option>
  </select>
  {#each srcts as srct}
    <select bind:value={src} onchange={srcChange(srct.url)}>
      <option value=""></option>
      {#each srcs[srct.name] as item}
        <option value={item.id}>{item.name}</option>
      {/each}
    </select>
  {/each}

  {#if ttable !== null}
    <STimeTableComp ttable={ttable} />
  {/if}

{:catch error}
  <h1>error: {error}</h1>
{/await}
