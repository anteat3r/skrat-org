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

  let avatarUrl = $derived(pb.baseURL + "/api/files/users/" + pb.authStore.record.id + "/" + pb.authStore.record.avatar + "?thumb=50x50");
</script>

<main>
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
      Login with Discord
    </button>
  {/if}
  {/key}
</main>
