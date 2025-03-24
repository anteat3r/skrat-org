<script lang="ts">
    import { onMount } from "svelte";
  import { pb } from "../pb_store.svelte";
  import DayComp from "./DayComp.svelte";

  let ttable = $state(null);

  function forwardButtonPress(e: KeyboardEvent) {
    e.preventDefault();
    e.target.dispatchEvent(new MouseEvent("click")); 
  }

  function dayOnclick(idx: number) {
    return function() {}
  }

  function detailAlertCallback(idk: any) {
    return function() {}
  }

  let selDayIdx = $state(null);

  function getWeekDayByIdx(idx: number, next: boolean = false): Date {
    let nw = new Date();
    nw.setDate(nw.getDate() - nw.getDay() + 1 + idx + (next ? 7 : 0));
    return nw;
  }

  onMount(async function() {
    ttable = await pb.send("/api/kleo/mytt", {})
  })
</script>

{#if ttable !== null}
  <main>
  <table>
    <tbody>
      <tr>
        <th></th>
        {#each ttable.Hours as hour}
          <th>
            <h1>{hour.Caption}</h1>
            <p>{hour.BeginTime}</p>
            <p>{hour.EndTime}</p>
            <p class="spacer"></p>
          </th>
        {/each}
      </tr>
      {#each ttable.Days as day, dayIdx}
        <tr>
          <th>
            <div class="cell"
                onclick={dayOnclick(dayIdx)}
                role="button" tabindex="-1"
                onkeypress={forwardButtonPress}
            >
              <h1>{day.DayOfWeek}</h1>
              <h1>{day.Date}</h1>
              {#if day.events.length > 0}
                <p>{day.events.length} evts</p>
              {/if}
            </div>
          </th>
          {#each ttable.Hours.map((hour: { Id: any; }) => day.Atoms.filter((atom: { HourId: any; }) => atom.HourId == hour.Id)) as atoms}
            <td>
              {#each atoms as atom}
                <div
                    class="bk-{atom.color} cell" 
                    onclick={detailAlertCallback(atom)} 
                    role="button" tabindex="-1" onkeypress={forwardButtonPress}
                >
                  <div class="cell-top">
                    <div class="cell-topleft">{atom.GroupIds}</div>
                    <div class="cell-topright">{atom.RoomId}</div>
                  </div>
                  <div class="cell-middle">
                    <div> {atom.SubjectId} </div>
                  </div>
                  <div class="cell-bottom">
                    <div> {atom.TeacherId} </div>
                  </div>
                </div>
              {/each}
            </td>
          {/each}
        </tr>
        {#if selDayIdx === dayIdx}
          <DayComp events={day.events} type="personal" date={getWeekDayByIdx(dayIdx)} />
        {/if}
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
