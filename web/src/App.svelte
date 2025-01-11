<script lang="ts">
  import Home from './lib/home.svelte'
  import Navbar from './lib/navbar.svelte';
  import About_us from './lib/about_us.svelte';
  import Achieve from './lib/achieve.svelte'
  import Footer from './lib/footer.svelte'

  let page = $state(0);
</script>

<main>
  <div class="spacer"></div>
  <Navbar bind:page={page} />
  {#if page === 0}
    <Home />
  {:else if page === 1 }
    <About_us />
  {:else if page === 2 }
    <Achieve />
  {/if}

  <Footer />
</main>

<style>
  .read-the-docs {
    color: #888;
  }
  .spacer{
    width: 700px;
  }
</style>
