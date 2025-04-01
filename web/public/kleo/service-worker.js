/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
const sw = /** @type {ServiceWorkerGlobalScope} */ (/** @type {unknown} */ (self));

/* eslint-env browser, serviceworker */

sw.addEventListener('install', () => {
  sw.skipWaiting();
})

sw.addEventListener('push', (evt) => {
  const data = evt.data.json();

  if (data.type === "notif") {
    evt.waitUntil(
      sw.registration.showNotification(data.title, data.options),
    )
  }
})
