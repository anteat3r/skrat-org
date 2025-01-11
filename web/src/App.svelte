<script lang="ts">
  import Home from './lib/home.svelte'
  import Navbar from './lib/navbar.svelte';
  let page = $state(0);
</script>

<main>
  <Navbar bind:page={page} />
  {#if page === 0}
    <Home />
  {/if}

  <p class="read-the-docs">
    Skrat© 2025
  </p>
</main>

<style>
  .read-the-docs {
    color: #888;
  }
</style>
