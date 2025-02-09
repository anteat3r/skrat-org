<script lang="ts">
  import PocketBase from 'pocketbase';
  const pb = new PocketBase("https://skrat.org");
</script>

<div class="container">
  <h1> Fotky </h1>
  <p>
    Fotky z akcí
    <br> Zde jsou:
  </p>
</div>

{#await pb.collection("accomplishments").getFullList({ sort: "-year, -ranked" })}
  čekám na data
{:then value}
  <table>
  <tbody>
    <tr>
        <td>rok</td>
        <td>Fotky Šimon</td>
        <td>Fotky Tomáš</td>
      </tr>
    {#each value as item}
      {#if item.url != "" || item.urlt != ""} 
      <tr>
        <td>{item.year}</td>
        <td> {#if item.url != ""}
            <a href="{item.url}" target="_blank" rel="noreferrer">{item.name}</a>
              {/if}
          </td>
        <td> {#if item.urlt != ""}
            <a href="{item.urlt}" target="_blank" rel="noreferrer">{item.name}</a>
              {/if}
          </td>
      </tr>
      {/if}
    {/each}
  </tbody>
  </table>
{:catch error}
  no data :( ({error})
{/await}


<style>
  .container{
    width: 100%;
    justify-content: center;
    height: auto;
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
  }

  table, td {
    border: 1px solid;
  }

  table {
    /*width: calc(40%);*/
    min-width: 600px;
    max-width: 1000px;
    margin-left: auto;
    margin-right: auto;
  }


</style>
