<script lang="ts">
  import PocketBase from 'pocketbase';
  const pb = new PocketBase("https://skrat.org");
</script>


{#await pb.collection("accomplishments").getFullList()}
  4ek8m na data
{:then value}
  <table>
  <tbody>
    {#each value as item}
      <tr>
        <td>{item.name}</td>
        <td>{item.year}</td>
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
