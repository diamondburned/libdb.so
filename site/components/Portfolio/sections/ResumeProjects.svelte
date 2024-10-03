<script lang="ts">
  export let resume: any;
</script>

<section class="projects">
  <h2>Projects</h2>
  <ul class="projects-list card content">
    {#each resume.projects as project}
      <li class="project-item content-item">
        <div class="header">
          <b class="name">{project.name ?? ""}</b>
          <span class="keywords">
            {(project.keywords ?? []).join(", ")}
          </span>
          {#if project.url}
            <a
              class="url"
              href={project.url.includes("://") ? project.url : `https://${project.url}`}
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

<style lang="scss">
  section.projects {
    .projects-list {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .project-item {
      p {
        margin: 0;
        margin-top: var(--adw-menu-margin);
      }
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
      }

      .keywords {
        grid-area: keywords;

        &:not(:empty) {
          opacity: 0.75;
          border-left: 1px solid rgba(255, 255, 255, 0.35);
          padding-left: 0.5em;
        }
      }

      .url {
        grid-area: url;
      }

      .url,
      .keywords {
        font-size: 0.9em;
      }

      @media (max-width: 500px) {
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
</style>
