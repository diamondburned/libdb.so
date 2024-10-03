<script lang="ts">
  import OpenInNew from "#/libdb.so/site/components/MaterialIcons/open_in_new.svelte";

  import ResumeWork from "./ResumeWork.svelte";
  import ResumeProjects from "./ResumeProjects.svelte";

  export let resume: any;
</script>

<section class="resume">
  <h2>Resume</h2>
  <div class="links card content" role="list">
    <a
      role="button"
      href="https://github.com/diamondburned/resume/blob/main/resume.pdf"
      class="content-item"
      target="_blank"
    >
      <OpenInNew /><span class="filename">PDF</span>
      <span class="source">(github.com)</span>
    </a>
    <a
      role="button"
      href="https://github.com/diamondburned/resume/blob/main/resume.json"
      class="content-item"
      style="--color: var(--pink-rgb);"
      target="_blank"
    >
      <OpenInNew /><span class="filename">JSON</span>
      <span class="source">(github.com)</span>
    </a>
  </div>
</section>

<ResumeWork {resume} />
<ResumeProjects {resume} />

<style lang="scss">
  section.resume {
    .links {
      display: flex;
      flex-direction: row;

      a[role="button"] {
        --color: var(--blue-rgb);

        width: 100%;
        align-self: center;
      }
    }

    .source {
      display: inline;
      font-size: 0.8em;
      vertical-align: baseline;
      opacity: 0.65;
    }

    .links {
      display: flex;
      flex-direction: column;

      @media (max-width: 400px) {
        flex-wrap: wrap;
      }
    }

    :global(svg) {
      vertical-align: top;
      margin-right: 0.5em;
      min-width: 1.75em;
    }
  }
</style>
