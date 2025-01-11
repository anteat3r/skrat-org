<script lang="ts">
  import skratLogo_color from './assets/svelte.svg'
  import skratLogo_cnb from '/scrat_logo_cnb.png'
  import Counter from './lib/Counter.svelte'
  import Navbar from './lib/navbar.svelte';
  let page = $state(0);
</script>

<main>
  <Navbar bind:page={page} />
  {#if page === 0}
  <div>
    <a href="about_us" target="_blank" rel="noreferrer">
      <img src={skratLogo_cnb} class="logo" alt="scrat_logo_cnb" />
    </a>
    <a href="https://svelte.dev" target="_blank" rel="noreferrer">
      <img src={skratLogo_color} class="logo svelte" alt="scrat_logo_color" />
    </a>
  </div>
  <h1>Vite + Svelte</h1>

  <div class="card">
    <Counter />
  </div>


  <p>
    Check out <a href="https://github.com/sveltejs/kit#readme" target="_blank" rel="noreferrer">SvelteKit</a>, the official Svelte app framework powered by Vite!
  </p>

  {/if}

  <p class="read-the-docs">
    Skrat© 2025
  </p>
</main>

<style>
  .logo {
    height: 6em;
    padding: 1.5em;
    will-change: filter;
    transition: filter 300ms;
  }
  .logo:hover {
    filter: drop-shadow(0 0 2em #646cffaa);
  }
  .logo.svelte:hover {
    filter: drop-shadow(0 0 2em #ff3e00aa);
  }
  .read-the-docs {
    color: #888;
  }
</style>
