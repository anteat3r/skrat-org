<script lang="ts">
  import Home from './lib/home.svelte'
  import Navbar from './lib/navbar.svelte';
  import AboutUs from './lib/about_us.svelte';
  import Achieve from './lib/achieve.svelte'
  import Footer from './lib/footer.svelte'
  import Photo from './lib/photo.svelte'
  import UrlShortener from './lib/url_shortener.svelte';

  import { Router, Route } from 'svelte-routing';

  let url = $state("");
</script>

<main>
<Router {url}>
  <div class="spacer"></div>
  <Navbar />

  <h1>SKRAT</h1>

  <Route path="/" component={Home} />
  <Route path="/about_us" component={AboutUs} />
  <Route path="/photos" component={Photo} />
  <Route path="/achieve" component={Achieve} />
  <Route path="/urls" component={UrlShortener} />

  <Footer />
</Router>
</main>

<style>
  main{
    width: 100vw;
  }
  .spacer{
    /*position: absolute;*/
    /*left: 0px;*/
    height: 50px;
    /*width: 100vw;*/
  }
</style>
