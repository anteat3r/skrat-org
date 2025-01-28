<script lang="ts">
  import Home from './lib/home.svelte'
  import Navbar from './lib/navbar.svelte';
  import About_us from './lib/about_us.svelte';
  import Achieve from './lib/achieve.svelte'
  import Footer from './lib/footer.svelte'
  import Photo from './lib/photo.svelte'

  let page = $state(0);
</script>

<main>
  <div class="spacer"></div>
  <Navbar bind:page={page} />
  <h1>SKRAT</h1>
  {#if page === 0}
    <Home />
  {:else if page === 1 }
    <About_us />
  {:else if page === 2 }
    <Photo />
  {:else if page === 3 }
    <Achieve />
  {/if}

  <Footer />
</main>

<style>
  .spacer{
    width: 700px;
  }
</style>
