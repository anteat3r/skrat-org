<script lang="ts">
  import DayComp from "./DayComp.svelte";

  let { ttable, useOwner = false } = $props();

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
    
    //document.innerHTML = detail
    }
  }

  let selDayIdx: number | null = $state(null);

  function dayOnclick(idx: number) {
    return function() {
      if (selDayIdx === idx) {
        selDayIdx = null;
      } else {
        selDayIdx = idx;
      }
    }
  }

  function parseDateFromDayTitle(title: string): Date {
    let nums = title.split("\n")[1].split(".");
    console.log(nums);
    let nw = new Date();
    return new Date(nw.getFullYear(), parseInt(nums[1])-1, parseInt(nums[0]));
  }

</script>

<main>
<table>
  <tbody>
    <tr>
      <th>
        {#if useOwner}
          <h1>{Object.values(ttable.data).map((e: any) => e.title)[Math.floor(Math.random() * Object.keys(ttable.data).length)]}</h1>
        {/if}
      </th>
      {#each ttable.hours as hour}
        <th>
          <h1>{hour.idx}</h1>
          <p>{hour.dur.split(" - ")[0]}</p>
          <p>{hour.dur.split(" - ")[1]}</p>
          <p class="spacer"></p>
        </th>
      {/each}
    </tr>
    {#each (useOwner ? Object.values(ttable.data) : ttable.days) as day, dayIdx}
      <tr>
        <th>
          <div class="cell"
              onclick={dayOnclick(dayIdx)}
              role="button" tabindex="-1"
              onkeypress={forwardButtonPress}
          >
            {#if useOwner}
              <h1>{day.owner} {day.special}</h1>
            {:else}
              {#if day.special === ""}
                <h1>{day.title.split(" ")[0]}</h1>
                <h1>{day.title.split(" ")[1]}</h1>
              {:else}
                  <h1>{day.title.replace("\n", " ")}</h1>
                  <h1>{day.special}</h1>
              {/if}
            {/if}
            {#if day.events.length > 0}
              <p>{day.events.length} evts</p>
            {/if}
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
      {#if selDayIdx === dayIdx}
        <DayComp events={day.events} type="personal" date={parseDateFromDayTitle(day.title)} />
      {/if}
    {/each}
  </tbody>
  </table>
</main>

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
