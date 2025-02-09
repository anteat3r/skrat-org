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
    <div class="bar">
    {#each value as item}
    <div class="card">
      <div class="title">{item.name}</div>
      <div class="year">{item.year}</div>
      <div class="rank">{item.tdesc}</div>
      <div class="desc">{item.text}</div>

    </div>
    {/each}
    </div>

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

  .bar{
    width: 99vw;
    /*max-width: 3000px;*/

    margin-left: 1vw;
    /*margin-right: 30px;*/
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 20px;
    /*margin: auto;*/
    justify-content: center;
    z-index: 0;
    /*justify-content: center;*/
  }

  .card{
    /*position: relative;*/
    padding: 20px;
    background-color: #440044;
    border-radius: 30px;
    transition: transform 0.3s;
    width: 300px;
    height: 350px;
    background-color: #440044;
    border: 5px solid #330033;
  }

  .title{
    margin: -13px 0px 25px 0px;
    font-size: 30px;
  }
  
  .year{
    font-size: 25px;
    color: #aabbaa;
  }

  .rank{
    font-size: 22px;
  }

  .desc{
    text-align: justify;
    font-size: 15px;
  }

  .card:nth-child(2n) {
    margin-top: 70px;
  }

  .card:nth-child(3n) {
    margin-top: 140px;
  }
  /**/
  /*.card:nth-child(4n) {*/
  /*  margin-top: 210px;*/
  /*}*/
  /**/
  /*.card:nth-child(5n) {*/
  /*  margin-top: 140px;*/
  /*}*/

</style>
