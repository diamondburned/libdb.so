<script lang="ts">
  import "libwebring/dist/webring.css";
  import "libwebring/dist/webring-element.js";

  export let src: string | null = null;
  export let data: Record<string, unknown> | null = null;
  export let name = "diamond";

  $: src_ = src ?? undefined;
  $: data_ = data ? JSON.stringify(data) : undefined;
  $: console.log({ src_, data_ });
</script>

{#if src_ || data_}
  <webring-element {name} src={src_} data={data_}>
    <section class="webring">
      <span class="ring" />
      <div>
        <a class="left" target="_blank" href={"#"}>_</a>
        <span class="middle" />
        <a class="right" target="_blank" href={"#"}>_</a>
      </div>
    </section>
  </webring-element>
{/if}

<style global lang="scss">
  webring-element {
    section.webring {
      display: flex;
      padding: 1em 0;
      flex-direction: column;

      & > * {
        margin: 0;
      }

      & > :nth-child(1) {
        align-self: center;
      }

      & > :nth-child(2) {
        display: grid;
        grid-template-columns: repeat(3, 1fr);

        a {
          color: var(--blue);
          text-decoration: none;

          &:hover {
            text-decoration: underline;
          }
        }
      }

      .left,
      .middle,
      .right {
        font-size: 0.9em;
      }

      .left {
        text-align: left;
        &::before {
          content: "‹ ";
        }
      }

      .middle {
        text-align: center;
        opacity: 0.5;
      }

      .right {
        text-align: right;
        &::after {
          content: " ›";
        }
      }

      .ring {
        opacity: 0.75;
      }
    }
  }
</style>
