<script lang="ts">
  import { pb } from "../pb_store.svelte";

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

  function forwardButtonPress(e: KeyboardEvent) {
    e.preventDefault();
    e.target.dispatchEvent(new MouseEvent("click")); 
  }

  function detailAlertCallback(detail: string) {
    return function() {
    //  alert(detail);
    var podrobnosti = JSON.parse(detail)
    
    document.getElementById('popup').style.display='block';
    //document.getElementById('fade').style.display='block';
    document.getElementById("infor").innerHTML = "místnost: " + podrobnosti.room ;
    document.getElementById("infos").innerHTML =  podrobnosti.subjecttext ;
    document.getElementById("infoth").innerHTML = "téma: " + podrobnosti.theme ;
    document.getElementById("infou").innerHTML = "učitel: " + podrobnosti.teacher ;
    document.getElementById("infog").innerHTML = "skupina: " + podrobnosti.group ;
    document.getElementById("infoch").innerHTML = "změna: " + podrobnosti.changeinfo ;
    //document.innerHTML = detail
    }
  }
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
    <main>
    <table>
      <tbody>
        <tr>
          <th></th>
          {#each ttable.hours as hour}
            <th>
              <h1>{hour.idx}</h1>
              <p>{hour.dur.split(" - ")[0]}</p>
              <p>{hour.dur.split(" - ")[1]}</p>
              <p class="spacer"></p>
            </th>
          {/each}
        </tr>
        {#each ttable.days as day}
          <tr>
            <th>
              <div class="cell">
                <h1>{day.title.split(" ")[0]}</h1>
                <h1>{day.title.split(" ")[1]}</h1>
              </div>
            </th>
            {#each day.hours as hour}
              <td>
                {#each hour.cells as cell}
                  <div
                      class="bk-{cell.color} cell" 
                      onclick={detailAlertCallback(cell.detail)} 
                      role="button" tabindex="-1" onkeypress={forwardButtonPress}
                  >
                    <div class="cell-top">
                      <div class="cell-topleft">{cell.group}</div>
                      <div class="cell-topright">{cell.room}</div>
                    </div>
                    <div class="cell-middle">
                      <div> {cell.subject} </div>
                    </div>
                    <div class="cell-bottom">
                      <div> {cell.teacher} </div>
                    </div>
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
    /* border: solid; */
  }
  main {
    overflow-x: scroll;
    width: fit-content;
  }
  td, th {
    border: solid 1px;
    display: flex;
    flex-direction: column;
    justify-content: stretch;
    flex-basis: 100%;
    width: 100px;
    padding: 2px;
  }
  .cell {
    width: 100%;
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    /* padding: 2px; */
    margin: 2px;
  }
  .cell-top {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    height: 20%;
    font-size: 12px;
    margin-inline: 3px;
  }
  .cell-middle {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    height: 60%;
    font-size: 20px;
  }
  .cell-bottom {
    display: flex;
    flex-direction: row;
    justify-content: start;
    align-items: end;
    height: 20%;
    font-size: 12px;
    margin-inline: 3px;
  }
  tr {
    display: flex;
    flex-direction: row;
    justify-content: stretch;
  }
  h1 {
    text-align: center;
    margin: 0px;
    white-space: pre-line;
  }
  p {
    margin: 0px;
  }
  .spacer {
    margin: 3px;
  }
  .bk-white {
    background-color: #31373a;
  }
  .bk-pink {
    background-color: maroon;
  }
  .bk-green {
    background-color: darkgreen;
  }
</style>
