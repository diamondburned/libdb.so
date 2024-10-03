<script lang="ts">
  import { type FetchedJSONLD, usingContext } from "#/libdb.so/site/lib/jsonld";

  import Badge from "#/libdb.so/site/components/Portfolio/sections/88x31Badge.svelte";
  import type { _88x31Badge } from "#/libdb.so/site/components/Portfolio/sections/88x31Badge.svelte";

  async function combined88x31s(doc: Awaited<FetchedJSONLD>) {
    return await Promise.all(
      [
        ...doc.self["libdb:88x31"], //
        ...doc.self["libdb:other88x31"],
      ].map((o) => usingContext<_88x31Badge>(doc, o, "https://0xd14.id#88x31Badge/")),
    );
  }

  export let doc: FetchedJSONLD;
</script>

{#await combined88x31s(doc) then badges}
  <div class="badges">
    {#each badges as badge}
      <Badge {badge} />
    {/each}
  </div>
{/await}

<style lang="scss">
  .badges {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(88px, 1fr));
    grid-gap: 0.25em;
    justify-content: center;
    justify-items: center;

    margin: 0 0.5em;
    margin-top: 1em;
    padding: 1em;
    border-top: 1px solid #fff3;
  }
</style>
