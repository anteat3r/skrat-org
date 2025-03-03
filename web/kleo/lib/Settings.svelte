<script lang="ts">
  import { pb } from '../pb_store.svelte';
  import githubLogo from '/github.svg'

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
    await navigator.serviceWorker.register("/kleo/service-worker.js");
    let perm = await Notification.requestPermission();
    console.log(perm);
    if (perm !== "granted") { return }
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
  }

  async function vapidTest() {
    let resp = await pb.send("/api/kleo/vapidtest", {
      method: "POST",
    });
    console.log(resp);
  }

  function urlBase64ToUint8Array(base64String: string): Uint8Array {
      const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
      const base64 = (base64String + padding)
        .replace(/\-/g, '+')
        .replace(/_/g, '/');
      const rawData = window.atob(base64);
      return Uint8Array.from([...rawData].map(char => char.charCodeAt(0)));
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
<button onclick={setupNotifs}>Setup Notifs</button>
<button onclick={vapidTest}>Send Test Notif</button>



<p>{JSON.stringify(pb.authStore.record)}</p>
