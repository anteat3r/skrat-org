<script lang="ts">
  import { link } from "svelte-routing";
</script>

<div class="floating">
<div class="container">
  <a class="menu" href="/" use:link>Domů</a>
  <a class="menu" href="/about_us" use:link>O nás</a>
  <a class="menu" href="/photos" use:link>Fotky</a>
  <a class="menu" href="/achieve" use:link>Úspěchy</a>
</div>
</div>

<style>
    .floating{
      position: fixed;
      width: 100%;
      margin-top: -52px;
      /*left: 34vw;*/
      /*margin-left: 30vw;*/
      /*background-color: #090306;*/
      /*border: 2px solid #090306;*/
      /*margin: auto,0;*/
      text-align: center;
      flex-grow: 2;
      z-index: 999;
    }
    .container {
        border: 2px solid #090306;
        background-color: #090306;
        /*background-color: #ffffff;*/
        min-width: 400px;
        max-width: 700px;
        margin-left: auto;
        margin-right: auto;
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        height: 50px;
    }
    .menu {
        /*background-color: #090306;*/
        justify-content: center;
        align-items: center;
        /*background-color: #ffffff;*/
        font-size: 22px;
        /*width: 100%;*/
        height: auto;
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        height: 50px;
        border: white solid 2px;
        flex-grow: 2;
        z-index: 999;
    }
</style>
