<script lang="ts">
  import { pb } from '../pb_store.svelte';

  let endpoint = $state("");

  async function login() {
    let resp = await pb.send("/api/kleo/endp", {
      query: {
        endp: endpoint,
      },
      method: "GET",
    })
    alert(resp);
  }
</script>

<input type="text" bind:value={endpoint} placeholder="endpoint">
<br>
<button onclick={login}>Query</button>
