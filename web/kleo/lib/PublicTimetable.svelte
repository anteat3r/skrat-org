<script lang="ts">
  import { pb } from "../pb_store.svelte";

  function srcChange(ttype: string) {
    return async function() {
      if (src === "") {
        ttable = null;
        return;
      }
      let resp = await pb.send(`/api/kleo/web/Actual/${ttype}/${src}`, {});
      ttable = resp;
    }
  }

  let ttable = $state(null);
  let src = $state("");

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
  
  {#each srcts as srct}
    <select bind:value={src} onchange={srcChange(srct.url)}>
      <option value=""></option>
      {#each srcs[srct.name] as item}
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
