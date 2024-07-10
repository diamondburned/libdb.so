<script lang="ts">
  import "libwebring/dist/webring.css";
  import { WebringElement } from "libwebring/dist/webring-element.js";
  import { WebringData } from "libwebring/lib/webring";

  export let data: any;
  $: typedData = data as WebringData & {
    link: string;
    self: string;
  };

  let element: WebringElement;
  $: {
    if (element) {
      element.removeAttribute("src");
      element.removeAttribute("statusSrc");
      element.setAttribute("data", data ? JSON.stringify(data) : undefined);
    }
  }
</script>

<webring-element name={typedData.self ?? "diamond"} bind:this={element}>
  <section class="webring">
    <a href={typedData.link ?? ""} class="ring">...</a>
    <div>
      <a class="left" target="_blank" href={"#"}>_</a>
      <span class="middle" />
      <a class="right" target="_blank" href={"#"}>_</a>
    </div>
  </section>
</webring-element>

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
        color: inherit;
        text-decoration: underline dashed;
        text-decoration-color: #fff5;
      }
    }
  }
</style>
