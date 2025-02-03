<script lang="ts">
  import githubLogo from '/github.svg'
  import { pb } from './pb_store';
</script>

{#if $pb.authStore.isValid && ( $pb.authStore.isSuperuser || $pb.authStore.record.collectionName == "users" ) }
  Hello mr. {$pb.authStore.record.username}
{:else}
  <button onclick={() => $pb.collection("users").authWithOAuth2({provider: "github"})}>
    <img src={githubLogo} alt="github logo">
    Login with Discord
  </button>
{/if}
