<script lang="ts">
  import "libwebring/dist/webring.css";
  import "libwebring/dist/webring-element.js";
  import { WebringData } from "libwebring/lib/webring";

  export let data: any;
  $: typedData = data as WebringData & {
    self?: string;
  };

  let className: string;
  export { className as class };

  let element: HTMLElement;
  $: {
    if (element) {
      element.removeAttribute("src");
      element.removeAttribute("statusSrc");
      if (data) {
        element.setAttribute("data", JSON.stringify(data));
        console.log(JSON.stringify(data));
      } else {
        element.removeAttribute("data");
      }
    }
  }
</script>

<webring-element class={className} name={typedData.self ?? "diamond"} bind:this={element}>
  <a href={"#"} target="_blank" class="ring">...</a>
  <div class="links">
    <a class="left" target="_blank" href={"#"}>_</a>
    <span class="middle" />
    <a class="right" target="_blank" href={"#"}>_</a>
  </div>
</webring-element>

<style global lang="scss">
  webring-element {
    display: flex;
    padding: 1em 0;
    flex-direction: column;

    & > * {
      margin: 0;
    }

    .ring {
      align-self: center;
    }

    .links {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
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
  }
</style>
