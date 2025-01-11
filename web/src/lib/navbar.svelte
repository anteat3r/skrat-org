<script lang="ts">
  let { page = $bindable() } = $props();
</script>

<div class="container">
  <div class="menu"> <a href="/" on:click|preventDefault={() => (page = 0)}> domů </a></div>
  <div class="menu"> <a href="about_us" on:click|preventDefault={() => (page = 1)}> o nás </a></div>
  <div class="menu"> <a href="achieve" on:click|preventDefault={() => (page = 2)}> úspěchy </a></div>
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
        width: 32%;
        height: auto;
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        height: 50px;
        border: white solid 2px;
    }
</style>
