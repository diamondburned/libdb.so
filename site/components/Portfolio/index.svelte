<script lang="ts">
  import Window from "#/libdb.so/site/components/Window.svelte";
  import OpenInNew from "#/libdb.so/site/components/MaterialIcons/open_in_new.svelte";

  const resumeURL =
    "https://raw.githubusercontent.com/diamondburned/resume/main/resume.json";
  const resume = fetch(resumeURL)
    .then((r) => r.json())
    .catch((err) => {
      console.error("Failed to fetch resume:", err);
      throw err;
    });
</script>

<Window
  view="portfolio"
  maxWidth="max(50vw, 600px)"
  maxHeight="max(90vh, 1000px)"
  scrollable
>
  <h3 slot="title">About</h3>

  <div class="portfolio-content">
    <section class="banner">
      <img src="/_assets/banner.png" alt="Banner" />
    </section>

    <section class="about">
      <div class="intro">
        <img src="/_assets/avatar.webp" alt="Diamond" />
        <div>
          <span>Hi, I'm</span>
          <h1>Diamond!</h1>
        </div>
      </div>
      <p>
        I'm a <b>4th-year Computer Science major 👩🎓</b>
        and past <b>Software Engineer Intern 👩‍💻 🖥️</b>
      </p>
      <p>
        I am what ChatGPT calls the world's biggest "open source cheerleader"!
        <br />
        I'm passionate about making the world a better place through technology and
        open source.
      </p>
    </section>

    <section class="annoyance">
      <p>
        <b>Hey!!</b> You should totally check out the <b><u>xterm.js</u></b> window
        underneath!
      </p>
    </section>

    <section class="links">
      <h2>Links</h2>
      <p>For quick access, here are my <b>links</b>!</p>
    </section>

    <section class="resume">
      <h2>Resume</h2>
      <div>
        <a
          role="button"
          href="https://github.com/diamondburned/resume/blob/main/resume.pdf"
          target="_blank"
        >
          <OpenInNew />
          <span class="filename">resume.pdf</span>
          <span class="source">(github.com)</span>
        </a>
      </div>
    </section>

    {#await resume}
      <span class="loading">Give me a bit, I'm loading the rest!</span>
    {:then resume}
      <section class="work">
        <h2>Experience</h2>
        <ol class="work-list">
          {#each resume.work as work}
            <li class="work-item">
              <h4>
                <b class="company-name">{work.company ?? ""}</b>
                <span class="location">{""}</span>
                <span class="position">{work.position ?? ""}</span>
                <span class="duration">
                  {#if work.startDate && work.endDate}
                    {work.startDate} - {work.endDate}
                  {:else if work.startDate}
                    {work.startDate} - now
                  {/if}
                </span>
              </h4>
              <ul class="highlights-list">
                {#each work.highlights as highlight}
                  <li class="highlight-item">{highlight}</li>
                {/each}
              </ul>
            </li>
          {/each}
        </ol>
      </section>

      <section class="projects">
        <h2>Projects</h2>
        <ul class="projects-list">
          {#each resume.projects as project}
            <li class="project-item">
              <div class="header">
                <b class="name">{project.name ?? ""}</b>
                <span class="keywords">
                  {(project.keywords ?? []).join(", ")}
                </span>
                {#if project.url}
                  <a
                    class="url"
                    href={project.url.includes("://")
                      ? project.url
                      : `https://${project.url}`}
                    target="_blank"
                  >
                    {project.url ?? ""}
                  </a>
                {/if}
              </div>
              <p class="description">{project.description ?? ""}</p>
            </li>
          {/each}
        </ul>
      </section>
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
    padding: 1em 0;
    margin: 0 auto;

    width: 100%;
    max-width: clamp(400px, 80vw, 550px);
    line-height: 1.5;

    display: flex;
    flex-direction: column;
    gap: 1em;

    @media (max-width: 400px) {
      gap: 0.5em;
      padding: 0.5em 0;
    }

    & > * {
      margin-bottom: 1em;
    }

    section {
      margin: 0 0.5em;
      padding: 0 1em;

      font-family: "Lato";
      font-size: 1.05em;

      background-color: rgba(255, 255, 255, 0.05);
      border: 1px solid var(--window-border-color);
      border-radius: 10px;
      box-shadow: 0 2px 16px -6px rgba(0, 0, 0, 0.52);
      box-sizing: border-box;

      h1,
      h2,
      h3,
      h4,
      h5,
      h6 {
        margin: 1rem 0;
        font-family: "Nunito";
      }

      h1,
      h2,
      h3 {
        line-height: 1.25;
      }

      & > * {
        margin: 1em 0;
      }

      a {
        text-decoration: none;
        color: var(--blue);

        &:hover {
          text-decoration: underline;
        }
      }

      @media (max-width: 400px) {
        font-size: 1em;
      }
    }

    section.banner {
      padding: 0;
      height: clamp(100px, 15vh, 150px);

      img {
        image-rendering: pixelated;
        width: 100%;
        height: 100%;
        margin: 0;
        object-fit: cover;
        border-radius: 10px;
      }
    }

    section.about {
      div.intro {
        display: flex;
        flex-direction: row;
        align-items: flex-end;
        gap: 1em;
        line-height: 1.15;
        margin-top: 1em;

        @media (max-width: 400px) {
          flex-direction: column;
          align-items: flex-start;
        }

        & > img {
          grid-area: img;
          width: 120px;
          aspect-ratio: 1/1;
          border-radius: 10px 0 10px 0;
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
            font-weight: 700;
            font-size: 2.5em;
            margin: 0;
          }
        }
      }

      p:last-child {
        max-width: 525px;
      }
    }

    section.annoyance {
      border: 1px solid rgba(var(--pink-rgb), 0.4);
      background-color: rgba(var(--pink-rgb), 0.1);
    }

    section.resume {
      & > div {
        display: flex;
        flex-direction: row;
      }

      a[role="button"] {
        text-decoration: none;

        width: 100%;
        padding: 0.5em;
        border-radius: 5px;
        align-self: center;

        background-color: rgba(var(--blue-rgb), 0.1);
        transition: background-color 0.1s ease-in-out;

        &:hover {
          background-color: rgba(var(--blue-rgb), 0.2);
        }
      }

      .source {
        display: inline;
        font-size: 0.8em;
        vertical-align: baseline;
        opacity: 0.65;
      }

      :global(svg) {
        vertical-align: top;
      }
    }

    section.work {
      .work-list {
        list-style: none;
        padding: 0;
      }

      .work-item {
        margin: 1em 0;

        h4 {
          display: grid;
          grid-template-columns: 1fr auto;
          grid-template-rows: auto auto;

          margin-bottom: 0.5em;

          & > *:nth-child(-n + 2) {
            font-size: 1.1em;
            font-weight: bold;
          }

          & > *:nth-last-child(-n + 2) {
            font-weight: normal;
          }

          @media (max-width: 400px) {
            grid-template-columns: 1fr;
          }
        }
      }

      .highlights-list {
        padding-left: 1.5em;
        padding-right: 0.5em;
        list-style: disc;
      }

      .highlight-item {
        margin: 0.25em 0;
        padding-left: 0.25em;
        font-size: 0.95em;
      }
    }

    section.projects {
      .projects-list {
        list-style: none;
        padding: 0;
      }

      .header {
        display: grid;
        grid-template-areas: "name keywords url";
        grid-template-rows: 1fr;
        grid-template-columns: auto 1fr auto;
        grid-gap: 0.5em;
        align-items: baseline;

        .name {
          grid-area: name;
          font-size: 1.1em;
        }

        .keywords {
          grid-area: keywords;

          &:not(:empty) {
            border-left: 1px solid rgba(255, 255, 255, 0.35);
            padding-left: 0.5em;
          }
        }

        .url {
          grid-area: url;
        }

        .url,
        .keywords {
          font-size: 0.95em;
        }

        @media (max-width: 400px) {
          grid-gap: 0;
          grid-template-areas:
            "name url"
            "keywords keywords";
          grid-template-columns: 1fr auto;
          grid-template-rows: auto auto;

          .keywords,
          .keywords:not(:empty) {
            border-left: none;
            padding: 0;
          }
        }
      }

      .description {
        margin-top: 0.5em;
      }
    }

    .loading {
      opacity: 0.5;
      font-size: 0.9em;
      text-align: center;
    }
  }
</style>
