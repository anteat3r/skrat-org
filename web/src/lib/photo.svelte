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
    {#each value as item}
      {#if item.url != ""} 
      <tr>
        <td>{item.year}</td>
        <td>
            <a href="{item.url}" target="_blank" rel="noreferrer">{item.name}</a>
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
    width: calc(100%);
  }


</style>
