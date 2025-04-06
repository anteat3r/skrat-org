<script lang="ts">
  import Home from './lib/home.svelte'
  import Navbar from './lib/navbar.svelte';
  import AboutUs from './lib/about_us.svelte';
  import Achieve from './lib/achieve.svelte'
  import Footer from './lib/footer.svelte'
  import Photo from './lib/photo.svelte'
  import UrlShortener from './lib/url_shortener.svelte';

  import { Router, Route } from 'svelte-routing';
    import { preventDefault } from 'svelte/legacy';

  let url = $state("");

  function scrollIntoView({ target }) {
		const el: Element = document.querySelector(target.getAttribute('href'));
		if (!el) return;
    el.scrollIntoView({
      behavior: 'smooth'
    });
  }

  function scrollCarousel(e: MouseEvent) {
    let f: Element = e.target as Element;
    let imgs = f.parentNode?.querySelectorAll("img");
    if (imgs == null) return;
    let aidx = imgs.values().toArray().findIndex((g) => g.hasAttribute("data-active"));
    console.log(imgs.values().toArray());
    let len = imgs.length;
    imgs[aidx].removeAttribute("data-active");
    imgs[(aidx + 1) % len].setAttribute("data-active", "data-active");
  }

  function forwardButtonPress(e: KeyboardEvent) {
    e.preventDefault();
    e.target.dispatchEvent(new MouseEvent("click")); 
  }

</script>

{#snippet carImg(url: string, alt: string)}
  <img onclick={scrollCarousel} data-active src={url} alt={alt} role="button" tabindex="0" onkeypress={forwardButtonPress}>
{/snippet}

<main>
  <div class="begin">
    <div id="title">
      <h1 class="header">SKRAT</h1>
    </div>
    <div class="navbar">
      <a class="menu" href="#title" onclick={scrollIntoView}>Domů</a>
      <a class="menu" href="#about" onclick={scrollIntoView}>O nás</a>
      <a class="menu" href="#photos" onclick={scrollIntoView}>Fotky</a>
      <a class="menu" href="#achievments" onclick={scrollIntoView}>Úspěchy</a>
    </div>
  </div>
  <div class="spacer"></div>
  <div id="about">
    <h1>Kdo jsme?</h1>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Alias explicabo facilis, rerum repudiandae voluptatibus suscipit nisi dolorem neque? Doloribus quo odit aut. Commodi dignissimos eius in itaque sequi aspernatur mollitia? Lorem ipsum dolor, sit amet consectetur adipisicing elit. Voluptatibus, optio harum doloremque nihil quia a neque minima consequuntur delectus quisquam ipsa cumque. Molestiae magni voluptatibus ut iste incidunt. Asperiores, alias.</p>
    {@render carImg("https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fpleated-jeans.com%2Fwp-content%2Fuploads%2F2023%2F05%2Ffunny-random-pictures-may13.jpg&f=1&nofb=1&ipt=cbcaa87d2fecd633c978bda634cc21e6323a8ccfd0976d91d2718cac4cfaf5da&ipo=images", "fyziklani")}
  </div>
  <div class="spacer"></div>
  <div class="article">
    <h1>Fyziklani</h1>
    <p>Lorem ipsum dolor sit, amet consectetur adipisicing elit. Quasi sed voluptatum, deserunt ea nobis distinctio natus quo placeat magnam ut expedita dicta corporis obcaecati nostrum nihil soluta eum doloremque quis.</p>
    {@render carImg("https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fpleated-jeans.com%2Fwp-content%2Fuploads%2F2023%2F05%2Ffunny-random-pictures-may13.jpg&f=1&nofb=1&ipt=cbcaa87d2fecd633c978bda634cc21e6323a8ccfd0976d91d2718cac4cfaf5da&ipo=images", "fyziklani")}
    {@render carImg("https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fi.pinimg.com%2Foriginals%2F6d%2Fc4%2Fc8%2F6dc4c867fc9454a688aa72578db03568.png&f=1&nofb=1&ipt=25bbeff5cfe561af2bc632b8c2e8ea08604b47a7e1f1f007c8381047ea0cda76&ipo=images", "idk bro")}
  </div>
  <div class="article left">
    <h1>FLL 2025</h1>
    <p>Lorem ipsum dolor sit, amet consectetur adipisicing elit. Quasi sed voluptatum, deserunt ea nobis distinctio natus quo placeat magnam ut expedita dicta corporis obcaecati nostrum nihil soluta eum doloremque quis.</p>
    <img src="https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fpleated-jeans.com%2Fwp-content%2Fuploads%2F2023%2F05%2Ffunny-random-pictures-may13.jpg&f=1&nofb=1&ipt=cbcaa87d2fecd633c978bda634cc21e6323a8ccfd0976d91d2718cac4cfaf5da&ipo=images" alt="fyziklani">
  </div>
  <div class="spacer"></div>
  <div id="members">
    <h1 id="robotitle">Robotický tým ze srdce Prahy</h1>
    <br>
    <div class="medallion" id="simon">
      <img src="/simon.png" alt="simon">
      <p class="medallion-text">proste simon</p>
    </div>
    <div class="medallion" id="ondra">
      <img src="/simon.png" alt="ondra">
      <p class="medallion-text">proste simon</p>
    </div>
    <div class="medallion" id="rosta">
      <img src="/simon.png" alt="rosta">
      <p class="medallion-text">proste simon</p>
    </div>
    <div class="medallion" id="tomas">
      <img src="/simon.png" alt="tomas">
      <p class="medallion-text">proste simon</p>
    </div>
    <div class="medallion" id="filip">
      <img src="/simon.png" alt="filip">
      <p class="medallion-text">proste simon</p>
    </div>
    <div class="medallion" id="vojta">
      <img src="/simon.png" alt="vojta">
      <p class="medallion-text">proste simon</p>
    </div>
  </div>
  <div class="spacer"></div>
  <hr>
  <div id="footer">
    <a href="https://instagram.com/skrat_org">
      <img src="https://strela-vlna.gchd.cz/images/instagram_icon.png" alt="instagram">
    </a>
  </div>
</main>

<style>
  img {
    border-radius: 10px;
  }
  main {
    width: 100svw;
    background-repeat: no-repeat;
    background-image: url("/lego1.svg"), url("/lego2.svg"), url("/lego3.svg"), url("/lego4.svg");
    background-position: right 0px top 100px, 200px 200px, 20% calc(100vh - 200px), calc(100vw - 500px) calc(100vh - 200px);
    background-size: 400px, 100px, 400px, 300px;
    @media (max-width: 1000px) {
      background-image: url("/lego1.svg"), url("/lego3.svg"), url("/lego2.svg");
      background-position: right 0px top 80px, -50px 450px, right 10% top 550px;
      background-size: 400px, 400px, 100px;
    }
    display: flex;
    flex-direction: column;
    justify-content: start;
    align-items: center;
  }
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
    width: 100svw;
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
  #title {
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
  }
  .menu:hover {
    font-weight: bolder;
  }
  .header {
    font-size: 100px;
    text-shadow: 0 0 5px white;
    width: fit-content;
  }
  .spacer {
    height: 200px;
  }
  #simon { grid-area: simon; }
  #ondra { grid-area: ondra; }
  #rosta { grid-area: rosta; }
  #tomas { grid-area: tomas; }
  #filip { grid-area: filip; }
  #vojta { grid-area: vojta; }
  #robotitle {
    grid-area: robotitle;
    @media (max-width: 700px) {
      font-size: 40px;
    }
  }
  .medallion {
    padding-inline: 50px;
    padding-top: 40px;
    display: flex;
    flex-direction: row;
  }
  .medallion-text {
    padding-left: 10px;
    width: 100px;
  }
  #members {
    display: grid;
    grid-template-areas: 
      'robotitle'
      'simon'
      'ondra'
      'rosta'
      'tomas'
      'filip'
      'vojta';
    @media (min-width: 800px) {
      grid-template-areas: 
        'robotitle robotitle'
        'simon ondra'
        'rosta tomas'
        'filip vojta';
    }
    @media (min-width: 1300px) {
      grid-template-areas: 
        'robotitle robotitle robotitle'
        'simon ondra rosta'
        'tomas filip vojta';
    }
  }
  hr {
    width: 100vw;
  }
  .article {
    display: grid;
    grid-template-areas:
      'title image'
      'text image';
    @media (max-width: 1000px) {
      grid-template-areas:
        'title'
        'text'
        'image';
    }
    margin: 50px;
  }
  .article.left {
    grid-template-areas:
      'image title'
      'image text';
  }
  .article > h1 {
    grid-area: title;
  }
  .article > p {
    grid-area: text;
    margin-bottom: 30px;
    @media (min-width: 1000px) {
      width: 30vw;
    }
  }
  .article > img {
    grid-area: image;
    width: calc(100vw - 100px);
    @media (min-width: 1000px) {
      width: 40vw;
      margin-left: 50px;
    }
  }
  .article.left > img {
    @media (min-width: 1000px) {
      margin-left: 0px;
      margin-right: 50px;
    }
  }
  img[data-active] {
    display: block;
  }
  .article > img:not([data-active]) {
    display: none;
  }
  #about {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  #about > p {
    margin-inline: 100px;
  }
  #about > img {
    margin-inline: 100px;
    width: 70vw;
  }
  #footer > a > img {
    width: 30px;
    height: 30px;
  }
</style>
