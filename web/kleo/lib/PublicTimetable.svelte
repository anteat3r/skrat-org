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
    <main>
    <table>
      <tbody>
        <tr>
          <th></th>
          {#each ttable.hours as hour}
            <th>
              <h1>{hour.idx}</h1>
              <h5>{hour.dur.replaceAll(" ", "\xa0")}</h5>
            </th>
          {/each}
        </tr>
        {#each ttable.days as day}
          <tr>
            <th>
              <h1>{day.title.replaceAll(" ", "\n")}</h1>
            </th>
            {#each day.hours as hour}
              <td>
                {#each hour.cells as cell}
                  <div class="bk-{cell.color}">
                    <h1>{cell.subject}</h1>
                  </div>
                {/each}
              </td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
    </main>
  {/if}

{:catch error}
  <h1>error: {error}</h1>
{/await}

<style>
  table {
    table-layout: fixed;
    width: max-content;
    border: solid;
  }
  main {
    overflow-x: scroll;
    width: fit-content;
  }
  td, th {
    border: solid;
    display: flex;
    flex-direction: column;
    justify-content: stretch;
    flex-basis: 100%;
  }
  div {
    width: 100%;
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
  }
  tr {
    display: flex;
    flex-direction: row;
    justify-content: stretch;
  }
  h1, h5 {
    text-align: center;
    margin: 5px;
    white-space: nowrap;
  }
  .bk-white {
    background-color: transparent;
  }
  .bk-pink {
    background-color: maroon;
  }
  .bk-green {
    background-color: darkgreen;
  }
</style>
