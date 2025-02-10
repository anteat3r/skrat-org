<script lang="ts">
  import { pb } from '../pb_store.svelte.ts'
</script>

{#await pb.send("/api/kleo/endp", { query: { endp: "marks" } })}
  <h1>Loading... 🙄</h1>
{:then marks}
  {#each marks.Subjects as subject}
    <h1>{subject.Subject.Name}: {subject.AverageText}</h1>
    {#each subject.Marks as mark}
      <h3>{mark.MarkText} - {mark.Caption} ({mark.Theme})</h3>
    {/each}
  {/each}
{:catch error}
  <h1>error: {error}</h1>
{/await}
