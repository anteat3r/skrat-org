<script lang="ts">
  import BakaLogin from './lib/BakaLogin.svelte';
  import Home from './lib/Home.svelte'
  import githubLogo from '/github.svg'
  import { pb } from './pb_store.svelte';
  let reload = $state({});

  function logout() {
    pb.authStore.clear();
    reload = {};
  }

  async function githubLogin() {
    await pb.collection("users").authWithOAuth2({ provider: "github" });
    reload = {};
  }

  async function discordLogin() {
    await pb.collection("users").authWithOAuth2({ provider: "discord" });
    reload = {};
  }

  let avatarUrl = $derived(pb.baseURL + "/api/files/users/" + pb.authStore.record.id + "/" + pb.authStore.record.avatar + "?thumb=50x50");

  function detailback() {
    //document.getElementById('fade').style.display='none';
    document.getElementById('popup').style.display='none';

  }
  function forwardButtonPress(e: KeyboardEvent) {
    e.preventDefault();
    e.target.dispatchEvent(new MouseEvent("click")); 
  }
</script>

<div id="popup">
  <div id="fade" 
    onclick={detailback} 
    role="button" tabindex="-1" onkeypress={forwardButtonPress}
  > </div>
  <div class="info" id="info">
  <div class="info" id="infos"></div> 
  <div class="info" id="infor"></div> 
  <div class="info" id="infoth"></div> 
  <div class="info" id="infou"></div> 
  <div class="info" id="infog"></div> 
  <div class="info" id="infoch"></div>
  </div> 
</div>

{#key reload}
{#if pb.authStore.isValid && ( pb.authStore.isSuperuser || pb.authStore.record.collectionName == "users" ) }
  <h1>
    Hello mr. {pb.authStore.record.name}
    <img src={avatarUrl} alt="avatar">
  </h1>
  <br>
  {#await pb.collection("users").authRefresh()}
    <!-- promise is pending -->
  {:then value}
    {#if value.record.bakavalid}
        <Home />
    {:else}
        <BakaLogin reload={reload} />
    {/if}
  {/await}
{:else}
  <button onclick={githubLogin}>
    <img src={githubLogo} alt="github logo">
    Login with Github
  </button>
  <br>
  <button onclick={discordLogin}>
    <!-- <img src={githubLogo} alt="github logo"> -->
    Login with Discord
  </button>
{/if}
{/key}

<style>

  button{
    font-size: 30px;
    margin: 20px 20px 20px 20px;
  }

  img{
    width: 30px;
  }

  #fade{
    /*display: none;*/
    position: fixed;
    width: 100vw;
    height: 100vh;
    background-color: black;
    opacity: 20%;
  }
  #popup{
    display: none;
    position: fixed;
    /*background-color: white;*/
    /*color: black;*/
    height: 100vw;
    width: 100vw;
  }
  
  #info{
    top: 20%;
    -ms-transform: translateY(-50%);
    transform: translateY(-50%);
    border: white solid 2px;
    padding: 20px;
    padding-top: 20px;

  }
  .info{
    /*display: none;*/
    position: relative;
    background-color: black;
    font-size: 22px;
    color: white;
    height: auto;
    min-width: none;
    max-width: 20cm;
    padding: 2px;
    padding-left: 30px;
    /*width: 50vw;*/
    /*margin-top: 200px;*/
    /*margin: auto;*/
    margin-left: auto;
    margin-right: auto
  }
  
  @media(max-aspect-ratio: 1){
    #info{
      /*top: 50vw;*/
      margin-top: 50%;
    }
    #fade{
      width: 15300vw;
      height: 15300vw;
    }
  }

</style>
