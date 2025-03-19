<script lang="ts">
  import { pb } from "../pb_store.svelte";
  
  let dayIdx = $state(0)
  let tt = $state(null)

  async function dayIdxChange() {
    if (dayIdx == 0) return;
    tt = await pb.send(`/api/kleo/daytt/${ttype}`, { query: { day: dayIdx } });
  }

  function detailAlertCallback(detail: string) {
    return function() {
    //  alert(detail);
      var podrobnosti = JSON.parse(detail)
      
      document.getElementById('popup').style.display='block';
      //document.getElementById('fade').style.display='block';
      //document.getElementById("infotst").innerHTML = detail;
      if( podrobnosti.type == "atom"){
        document.getElementById("infor").innerHTML = "místnost: " + podrobnosti.room ;
        document.getElementById("infos").innerHTML =  podrobnosti.subjecttext ;
        document.getElementById("infoth").innerHTML = "téma: " + podrobnosti.theme ;
        if (podrobnosti.teacher != null){
          document.getElementById("infou").innerHTML = "učitel: " + podrobnosti.teacher ;
        }
        else{
          document.getElementById("infou").innerHTML = "" ;
        }
        document.getElementById("infog").innerHTML = "skupina: " + podrobnosti.group ;
        if (podrobnosti.changeinfo != ""){
          document.getElementById("infoch").innerHTML = "změna: " + podrobnosti.changeinfo ;
        }
        else{
          document.getElementById("infoch").innerHTML = "" ;
        }
      }
      if(podrobnosti.type == "removed"){
        document.getElementById("infor").innerHTML = "" ;
        document.getElementById("infos").innerHTML =  podrobnosti.subjecttext ;
        document.getElementById("infoth").innerHTML = "popis: " + podrobnosti.absentinfo ;
        document.getElementById("infou").innerHTML = "" ;
        document.getElementById("infog").innerHTML = "" ;
        document.getElementById("infoch").innerHTML = "změna: " + podrobnosti.removedinfo ;
      }
      if(podrobnosti.type == "absent"){
        document.getElementById("infor").innerHTML = "" ;
        document.getElementById("infos").innerHTML =  podrobnosti.subjecttext ;
        document.getElementById("infoth").innerHTML = "popis: " + podrobnosti.InfoAbsentName ;
        document.getElementById("infou").innerHTML = "zkratka: " + podrobnosti.absentinfo ;
        document.getElementById("infog").innerHTML = "" ;
        document.getElementById("infoch").innerHTML = "" ;
      }
      
    }
    //document.innerHTML = detail
  }

  function forwardButtonPress(e: KeyboardEvent) {
    e.preventDefault();
    e.target.dispatchEvent(new MouseEvent("click")); 
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
  <main>
  <table>
    <tbody>
      <tr>
        <th></th>
        {#each tt.hours as hour}
          <th>
            <h1>{hour.idx}</h1>
            <p>{hour.dur.split(" - ")[0]}</p>
            <p>{hour.dur.split(" - ")[1]}</p>
            <p class="spacer"></p>
          </th>
        {/each}
      </tr>
      {#each Object.entries(tt.data) as [title, day]}
        <tr>
          <th>
            <div class="cell">
              <h1>{title}</h1>
              <p>{(day as any).special}</p>
            </div>
          </th>
          {#each (day as any).hours as hour}
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
