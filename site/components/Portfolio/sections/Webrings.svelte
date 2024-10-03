<script lang="ts">
  import { FetchedJSONLD, usingContext } from "#/libdb.so/site/lib/jsonld";

  import Webring from "#/libdb.so/site/components/Webring.svelte";

  export let doc: FetchedJSONLD;
</script>

<section class="webring">
  <h2>Webrings</h2>
  <div class="content card">
    {#each doc.self["libdb:webring"] as webring}
      {#await usingContext(doc, webring, "https://0xd14.id#Webring/") then webring}
        <Webring class="content-item" data={webring} />
      {/await}
    {/each}
  </div>
</section>
