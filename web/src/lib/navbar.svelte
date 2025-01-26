<script lang="ts">
  let { page = $bindable() } = $props();
</script>

<div class="container">
  <div class="menu"> <a href="/" on:click|preventDefault={() => (page = 0)}> Domů </a></div>
  <div class="menu"> <a href="about_us" on:click|preventDefault={() => (page = 1)}> O nás </a></div>
  <div class="menu"> <a href="photo" on:click|preventDefault={() => (page = 2)}> Fotky </a></div>
  <div class="menu"> <a href="achieve" on:click|preventDefault={() => (page = 3)}> Úspěchy </a></div>
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
