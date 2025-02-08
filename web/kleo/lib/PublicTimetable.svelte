<script lang="ts">
  import { pb } from "../pb_store.svelte";

  function srcChange(ttype: string) {
    return async function() {
      if (src === "") {
        ttable = null;
        return;
      }
      let resp = await pb.send(`/api/kleo/web/Actual/${ttype[0].toUpperCase() + ttype.substring(1).toLowerCase()}/${src}`, {});
      console.log(resp);
      console.log(typeof resp);
    }
  }

  let ttable = $state(null);
  let src = $state("");
</script>

{#await pb.send("/api/kleo/websrcs", {})}
  <h1>Loading... 🙄</h1>
{:then srcs}
  
  {#each ["classes", "teachers", "rooms"] as srct}
    <select bind:value={src} onchange={srcChange(srct)}>
      <option value=""></option>
      {#each srcs[srct] as item}
        <option value={item.id}>{item.name}</option>
      {/each}
    </select>
  {/each}

  {#if ttable !== null}
    <p>{JSON.stringify(ttable)}</p>
  {/if}

{:catch error}
  <h1>error: {error}</h1>
{/await}
