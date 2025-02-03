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
      margin-top: -54px;
      margin-left: -4px;
      position: fixed;
      background-color: #090306;
      border: 4px solid #090306;
    }
    .container {
        background-color: #090306;
        /*background-color: #ffffff;*/
        width: 700px;
        height: auto;
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
        flex-grow: 1;
    }
</style>
