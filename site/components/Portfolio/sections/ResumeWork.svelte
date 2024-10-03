<script lang="ts">
  export let resume: any;
</script>

<section class="work hidden">
  <h2>Experience</h2>
  <ol class="work-list card content">
    {#each resume.work as work}
      <li class="work-item content-item">
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

<style lang="scss">
  section.work {
    .work-list {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .work-item {
      h4 {
        margin: 0;

        display: grid;
        grid-template-columns: 1fr auto;
        grid-template-rows: auto auto;

        & > *:nth-child(-n + 2) {
          font-weight: bold;
        }

        & > *:nth-last-child(-n + 2) {
          font-weight: normal;
          opacity: 0.75;
          font-size: 0.9em;
        }

        @media (max-width: 400px) {
          grid-template-columns: 1fr;
        }
      }
    }

    .highlights-list {
      padding-left: 1em;
      padding-right: 0.5em;
      list-style: disc;
    }

    .highlight-item {
      margin: var(--adw-menu-margin) 0;
      padding-left: 0.25em;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
</style>
