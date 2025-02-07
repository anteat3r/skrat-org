<script lang="ts">
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

  let avatarUrl = $derived(pb.baseURL + "/api/files/users/" + pb.authStore.record.id + "/" + pb.authStore.record.avatar + "?thumb=20x20");
</script>

{#key reload}
{#if pb.authStore.isValid && ( pb.authStore.isSuperuser || pb.authStore.record.collectionName == "users" ) }
  Hello mr. {pb.authStore.record.name}
  <img src={avatarUrl} alt="avatar">
  <br>
  <button onclick={logout}>Logout</button>
{:else}
  <button onclick={githubLogin}>
    <img src={githubLogo} alt="github logo">
    Login with Discord
  </button>
{/if}
{/key}
