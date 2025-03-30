<script lang="ts">
  import { pb } from '../pb_store.svelte';
  import githubLogo from '/github.svg'
  import CustomEndp from './CustomEndp.svelte';

  function logout() {
    pb.authStore.clear();
  }

  async function githubLogin() {
    await pb.collection("users").authWithOAuth2({ provider: "github" });
  }

  async function discordLogin() {
    await pb.collection("users").authWithOAuth2({ provider: "discord" });
  }

  async function setupNotifs() {
    try {
      await navigator.serviceWorker.register("/kleo/service-worker.js");
      let perm = await Notification.requestPermission();
      console.log(perm);
      if (perm !== "granted") {
        alert(perm);
        return;
      }
      let reg = await navigator.serviceWorker.ready;
      let old_sub = await reg.pushManager.getSubscription()
      if (old_sub !== null) {
        let res = await old_sub.unsubscribe()
        console.log(res);
      }
      const subscription = await reg.pushManager.subscribe({
        userVisibleOnly: false,
        applicationServerKey: "BGl8lG0dFZxVzpEwgnPQlHaqDuaBojbFJHJzh2CMYi8mZshivG7RRkGDLKAC6E23E6ELtp3ikBXuepRJBMRlbwc",
      });
      let resp = await pb.send("/api/kleo/setupnotifs", {
        method: "POST",
        body: {
          vapid: JSON.stringify(subscription.toJSON())
        }
      });
      console.log(resp);
      alert("alles OK, přihlásili jsme tě na notifkace");
    } catch (e) {
      alert(e);
    }
  }

  async function vapidTest() {
    let resp = await pb.send("/api/kleo/vapidtest", {
      method: "POST",
    });
    console.log(resp);
  }

  let personalReloadInterval = $state(pb.authStore.record.refresh_interval);

  async function setPersonalReloadInterval() {
    await pb.collection("users").update(pb.authStore.record.id, {
      refresh_interval: personalReloadInterval,
    });
    alert("ok");
  }

  async function unsubscribeNotifs() {
    await pb.collection("users").update(pb.authStore.record.id, {
      wants_refresh: false,
    });
    alert("ok když myslíš 😒");
  }
</script>

<button onclick={logout}>Logout</button>
<br><br>
<button onclick={githubLogin}>
  <img src={githubLogo} alt="github logo">
  Login with Github
</button>
<button onclick={discordLogin}>
  Login with Discord
</button>
<br><br>
{#if pb.authStore.record.wants_refresh}
  <p>přihlášen na notifikace yipeeee 🤩</p>
{:else}
  <p>nepřihlášen na notifikace 😨</p>
{/if}
<button onclick={setupNotifs}>Setup Notifs</button>
<button onclick={vapidTest}>Send Test Notif</button>
<button onclick={unsubscribeNotifs}>Odhlásit</button>
<br><br>
<label for="refresh_interval">refresh_interval</label>
<input type="number" id="refresh_interval" bind:value={personalReloadInterval}>
<button onclick={setPersonalReloadInterval}>Nastavit</button>
<br><br>
<CustomEndp />

<p>{__COMMIT_HASH__}</p>
