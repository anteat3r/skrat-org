<script lang="ts">
  import { pb } from '../pb_store.svelte';

  let endpoint = $state("");

  function login() {
    pb.send("/api/kleo/endp/" + endpoint, {
      method: "GET",
    })
  }
</script>

<input type="text" bind:value={endpoint} placeholder="endpoint">
<br>
<button onclick={login}>Query</button>
