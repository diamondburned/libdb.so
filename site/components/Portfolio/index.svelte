<script lang="ts">
  import Window from "#/libdb.so/site/components/Window.svelte";
  import FolderDocuments from "#/libdb.so/site/components/Papirus/folder-documents.svelte";

  const resumeURL =
    "https://raw.githubusercontent.com/diamondburned/resume/main/resume.json";
  const resume = fetch(resumeURL)
    .then((r) => r.json())
    .catch((err) => {
      console.error("Failed to fetch resume:", err);
      throw err;
    });
</script>

<Window id="portfolio">
  <h3 slot="title">Portfolio</h3>

  <div class="portfolio-content">
    <section class="about">
      <div class="intro">
        <img
          src="https://avatars.githubusercontent.com/u/8463786?v=4"
          alt="Diamond"
        />
        <div>
          <span>Hi, I'm</span>
          <h1>Diamond</h1>
        </div>
      </div>
      <p>
        I am what ChatGPT calls the world's biggest "open source cheerleader!"
      </p>
    </section>

    <section class="links">
      <h2>Links</h2>
      <p>For quick access, here are my <b>links</b>!</p>
    </section>

    <section class="resume">
      <h2>Resume</h2>
      <p>
        <a href="https://github.com/diamondburned/resume/blob/main/resume.pdf">
          <FolderDocuments />
          resume.pdf
        </a>
        for quick access to my resume!
      </p>
    </section>

    {#await resume}
      <span class="loading">
        Give me a bit, I'm loading the rest from my resume!
      </span>
    {:then resume}
      <section class="projects" />
    {:catch _}
      <span class="loading">
        I couldn't load my resume {":("}
        <br />
        Maybe the console can help?
      </span>
    {/await}
  </div>
</Window>

<style lang="scss">
  .portfolio-content {
    overflow-x: hidden;
    overflow-y: auto;

    padding: 0;
    line-height: 1.5;

    max-width: clamp(500px, 80vw, 800px);
    margin: auto;

    & > * {
      margin-bottom: 1em;
    }

    section {
      padding: 0 clamp(1em, 1.5vw, 1.5em);
      padding-top: 1.35em;

      margin-top: 1.5em;
      margin-bottom: 1.5em;

      &:first-child {
        margin-top: 2em;
      }

      &:not(:first-child) {
        border-top: 1px solid var(--window-border-color);
      }

      h2 {
        margin: 0;
      }

      font-family: "Nunito";

      p {
        margin: 1em 0;
        font-family: "Source Sans Pro";
      }
    }

    section.about {
      div.intro {
        display: flex;
        flex-direction: row;
        align-items: flex-end;
        gap: 1em;
        line-height: 1.15;

        @media (max-width: 350px) {
          flex-direction: column;
          align-items: flex-start;
        }

        & > img {
          grid-area: img;
          width: 80px;
          aspect-ratio: 1/1;
        }

        & > div {
          flex: 1;

          & > span {
            grid-area: span;
            font-size: 1.5em;
            font-weight: lighter;
          }

          & > h1 {
            grid-area: h2;
            font-weight: bold;
            font-size: 2.5em;
            margin: 0;
          }
        }
      }
    }

    section.resume {
      a {
        text-decoration: none;
        margin-right: 0.15em;

        :global(svg) {
          width: 1.25em;
          height: 1.25em;
          vertical-align: sub;
        }
      }
    }

    .loading {
      opacity: 0.5;
      font-size: 0.9em;
    }
  }
</style>
