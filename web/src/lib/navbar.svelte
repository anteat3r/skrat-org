<script lang="ts">
  import { link } from "svelte-routing";
</script>

<div class="begin">
  <div class="title">
    <h1 class="header">SKRAT</h1>
  </div>
  <div class="navbar">
    <a class="menu" href="/" use:link>Domů</a>
    <a class="menu" href="/about_us" use:link>O nás</a>
    <a class="menu" href="/photos" use:link>Fotky</a>
    <a class="menu" href="/achieve" use:link>Úspěchy</a>
  </div>
</div>

<style>
    a {
      color: white;
    }
    .navbar {
      display: flex;
      flex-direction: column;
      grid-area: navbar;
      justify-content: center;
      align-items: center;
    }
    .begin {
      display: grid;
      height: 100vh;
      width: 100vw;
      grid-template-areas:
        'title navbar'
        'title navbar';
      grid-template-columns: 80% 20%;
      @media (max-width: 1000px) {
        grid-template-areas:
          'title'
          'navbar';
        grid-template-columns: 100%;
      }
      background-attachment: local;
    }
    .title {
      grid-area: title;
      display: flex;
      justify-content: center;
      align-items: center;
      margin-top: 100px;
      @media (min-width: 1000px) {
        margin-bottom: 100px;
        margin-left: 100px;
      }
    }
    .menu {
      font-size: 30px;
      width: fit-content;
      /* transition: font-weight cubic-bezier(0.075, 0.82, 0.165, 1) 100ms; */
    }
    .menu:hover {
      font-weight: bolder;
    }
    .header {
      font-size: 100px;
      text-shadow: 0 0 5px white;
      width: fit-content;
    }
</style>
