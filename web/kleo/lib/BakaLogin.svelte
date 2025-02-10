<script lang="ts">
  import { pb } from '../pb_store.svelte';

  let username = $state("");
  let password = $state("");

  function login() {
    pb.send("/api/kleo/login", {
      method: "POST",
      body: {
        username,
        password,
      },
    })
  }
</script>

<input type="text" bind:value={username} placeholder="username">
<br>
<input type="password" bind:value={password} placeholder="password">
<br>
<button onclick={login}>Login</button>
