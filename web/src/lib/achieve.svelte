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
  <!--<table>-->
  <!--<tbody>-->
  <!--  {#each value as item}-->
  <!--    <tr>-->
  <!--      <td>{item.year}</td>-->
  <!--      <td>{item.name}</td>-->
  <!--      <td>{item.desc}</td>-->
  <!--    </tr>-->
  <!-- <!--{#if item.desc1 = ""}-->
      <!--  <tr>-->
      <!--      <td>{item.desc1}</td>-->
      <!--      <td>{item.desc2}</td>-->
      <!--      <td>{item.desc3}</td>-->
      <!--      <td>{item.desc4}</td>-->
      <!---->
      <!--  </tr>-->
      <!--{/if}--> 
  <!--  {/each}-->
  <!--</tbody>-->
  <!--</table>-->
  <!--konec-->

  <div class="container2">
  {#each value as item}
    <div class="border">
    <div class= "row bigger">
      <div class="column"> <div class= "menu bigger">{item.tdesc}</div></div>
      <div class="column1"> <div class= "menu bigger">{item.name}</div></div>
      <div class="column2"> <div class= "menu bigger">{item.year}</div></div>
    </div>
    {#if item.desc1 != ""}
    <div class= "row smaller">
     <div class="column3"><div class= "menu smaller">{item.desc1}</div></div>
     <div class="column3"><div class= "menu smaller">{item.desc2}</div></div>
     <div class="column3"><div class= "menu smaller">{item.desc3}</div></div>
     <div class="column3"><div class= "menu smaller">{item.desc4}</div></div>
    </div>
    {/if}
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

  .container2{
    width: 100%;
    height: auto;
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
    border:3px double  #aa00aa;


  }

  .row{
    /*flex-direction: column;*/
    width: calc(100%);
    height: auto;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    border:0.5px solid  #bbbbbb;


  }
  
  .border{
    border:1px solid  #ffffff;
    
  }

  .bigger{
    height: 50px;
    font-size: 22px;
    color: #4444ee;
  }

  .smaller{
    height: 30px;
    font-size: 16px;
    color: #aaaaaa;
  }

  .menu {
    justify-content: center;
    align-items: center;
    width: calc(100%);
    /*height: auto;*/
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
    /*background-color:#eee;*/
    border:1px solid  #666666;
    border-spacing:5px;/*cellspacing:poor IE support for  this*/
   /* border-collapse:separate;*/
    }
  
  .column{
    width: calc(27%);
    
  }
    
  .column1{
    width: calc(46%);
    
  }

  .column2{
    width: calc(27%);
    
  }

  .column3{
    width: calc(25%);
    
  }

  table, td {
    border: 1px solid;
  }

  table {
    width: calc(100%);
  }


</style>

