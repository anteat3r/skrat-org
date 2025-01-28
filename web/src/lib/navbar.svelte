<script lang="ts">
  import { link } from "svelte-routing";
</script>

<div class="container">
  <a class="menu" href="/" use:link>Domů</a>
  <a class="menu" href="/about_us" use:link>O nás</a>
  <a class="menu" href="/photos" use:link>Fotky</a>
  <a class="menu" href="/achieve" use:link>Úspechy</a>
</div>

<style>
    .container {
        background-color: #090306;
        /*background-color: #ffffff;*/
        width: 100%;
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
