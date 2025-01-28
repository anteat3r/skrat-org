<script lang="ts">
  import PocketBase from 'pocketbase';
  const pb = new PocketBase("https://skrat.org");
</script>

<div class="container">
  <h1> Úspěchy </h1>
  <p>
    Za posledních pár let se nám podařilo mnoho úspěchů v robotických soutěžích 
    <br> Zde jsou:
  </p>
</div>

{#await pb.collection("accomplishments").getFullList({ sort: "-year, -ranked" })}
  4ek8m na data
{:then value}
  <table>
  <tbody>
    {#each value as item}
      <tr>
        <td>{item.year}</td>
        <td>{item.name}</td>
        <td>{item.desc}</td>
      </tr>
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
