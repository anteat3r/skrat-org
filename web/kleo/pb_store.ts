import { readable } from "svelte/store";
import PocketBase from "pocketbase";

export const pb = readable(new PocketBase("https://skrat.org"));
