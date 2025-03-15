<script lang="ts">
  import { pb } from '../pb_store.svelte';

  let { reload } = $props();

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
    reload = {};
  }
</script>

<p>přihlášení do bakalář</p>
<p>vaše hesla neukládáme</p>
<input type="text" bind:value={username} placeholder="username">
<br>
<input type="password" bind:value={password} placeholder="password">
<br>
<button onclick={login}>Login</button>

<style>
  p{
    margin-left: 3px;
  }
  input{
    margin: 3px;
  }
  button{
    margin: 3px;
    margin-top: 7px;
    width: 180px;
    font-size: 16px;
  }
</style>
