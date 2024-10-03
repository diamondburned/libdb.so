<script lang="ts">
  import { pronouns, selectPronouns } from "#/libdb.so/site/lib/pronouns.js";
  import OpenInNew from "#/libdb.so/site/components/MaterialIcons/open_in_new.svelte";

  $: pronoun = (content: string) => selectPronouns(content, $pronouns);
</script>

<section class="docs">
  <h2>Documentation</h2>
  <div class="card content">
    <div class="content-item yap">
      <p>
        A <b>User manual</b> exists to help you understand the basics of operating, interacting with,
        and using Diamond. This
      </p>
      <p>
        The guide is maintained by Diamond
        {pronoun("herself")}.You may access it via the link below:
      </p>
    </div>
    <a class="content-item" role="button" href="https://docs.0xd14.id" target="_blank">
      <OpenInNew /><span>docs.0xd14.id</span>
    </a>
  </div>
</section>

<style lang="scss">
  section.docs {
    .content {
      display: flex;
      flex-direction: column;
      overflow: hidden;

      .yap {
        border-bottom: none;

        p {
          margin-top: 0;

          &:last-child {
            margin-bottom: 0;
          }
        }
      }

      a[role="button"] {
        width: 100%;

        :global(svg) {
          vertical-align: top;
          margin-right: 0.5em;
          min-width: 1.75em;
        }
      }
    }
  }
</style>
