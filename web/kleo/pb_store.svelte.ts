import PocketBase from "pocketbase";

export const pb: PocketBase = $state(new PocketBase("https://skrat.org"));
