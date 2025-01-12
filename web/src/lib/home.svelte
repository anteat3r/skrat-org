<script>
  import skratLogo_color from '/scrat_logo_color.svg'
  import skratLogo_cnb from '/scrat_logo_cnb.png'
  import Counter from './Counter.svelte'

</script>


<div>
  <a href="about_us" target="_blank" rel="noreferrer">
    <img src={skratLogo_cnb} class="logo" alt="scrat_logo_cnb" />
  </a>
  <a href="https://svelte.dev" target="_blank" rel="noreferrer">
    <img src={skratLogo_color} class="logo svelte" alt="scrat_logo_color" />
  </a>
</div>

<div class="card">
  <Counter />
</div>


<h2>Vite + Svelte</h2>

<style>
  .logo {
    height: 10em;
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
</style>
