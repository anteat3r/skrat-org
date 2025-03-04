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
      const subscription = await reg.pushManager.subscribe({
        userVisibleOnly: true,
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
<br>
<br>
<CustomEndp />
<br>
<br>
<button onclick={setupNotifs}>Setup Notifs</button>
<button onclick={vapidTest}>Send Test Notif</button>

verze JDU.SE.ZABÍT 3

<p>{JSON.stringify(pb.authStore.record)}</p>
