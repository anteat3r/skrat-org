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
  
  function detailAlertCallback() {
    document.getElementById('fade').style.display='none';
    document.getElementById('light').style.display='none';

  }
  function forwardButtonPress(e: KeyboardEvent) {
    e.preventDefault();
    e.target.dispatchEvent(new MouseEvent("click")); 
  }
</script>

<div id="fade" > </div>
<div id="light" 
    onclick={detailAlertCallback} 
    role="button" tabindex="-1" onkeypress={forwardButtonPress}
><div id="light2">kdsdfhjf</div> </div>

{#key reload}
{#if pb.authStore.isValid && ( pb.authStore.isSuperuser || pb.authStore.record.collectionName == "users" ) }
  <h1>
    Hello mr. {pb.authStore.record.name}
    <img src={avatarUrl} alt="avatar">
    <button onclick={logout}>Logout</button>
  </h1>
  <br>
  {#await pb.collection("users").authRefresh()}
    <!-- promise is pending -->
  {:then value}
    {#if value.record.bakavalid}
        <Home />
    {:else}
        <BakaLogin />
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

  #fade{
    display: none;
    position: fixed;
    width: 100vw;
    height: 100vh;
    background-color: black;
    opacity: 20%;
  }
  #light{
    display: none;
    position: fixed;
    /*background-color: white;*/
    color: black;
    height: 100vw;
    width: 100vw;
  }
  
  #light2{
    /*display: none;*/
    position: relative;
    background-color: white;
    color: black;
    height: auto;
    min-width: none;
    max-width: 1000px;
    /*width: 50vw;*/
    margin-left: auto;
    margin-right: auto
  }

</style>
