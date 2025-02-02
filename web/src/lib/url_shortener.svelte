<script lang="ts">
  let key = $state("");
  let value = $state("");

  async function sets() {
    let resp = await fetch(`https://skrat.org/ss?k=${key}&v=${value}`);
    let body = await resp.text();
    alert(`${resp.status} ${body}`);
  }
</script>


<input type="text" placeholder="key" bind:value={key}>
<br>
<input type="text" placeholder="value" bind:value={value}>
<br>
<button onclick={sets}>Nastavit</button>
