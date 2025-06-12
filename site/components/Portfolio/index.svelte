<script lang="ts">
  import { jsonldDoc } from "#/libdb.so/site/lib/jsonld.js";
  import { ToastStore } from "#/libdb.so/site/lib/toasts.js";

  import Toasts from "#/libdb.so/site/components/Toasts.svelte";
  import Window from "#/libdb.so/site/components/Window.svelte";

  import About from "#/libdb.so/site/components/Portfolio/sections/About.svelte";
  import Docs from "#/libdb.so/site/components/Portfolio/sections/Docs.svelte";
  import Links from "#/libdb.so/site/components/Portfolio/sections/Links.svelte";
  import Resume from "#/libdb.so/site/components/Portfolio/sections/Resume.svelte";
  import Webrings from "#/libdb.so/site/components/Portfolio/sections/Webrings.svelte";
  import Badges88x31 from "#/libdb.so/site/components/Portfolio/sections/88x31Badges.svelte";

  const toasts = new ToastStore();

  const resumeURL = "https://raw.githubusercontent.com/diamondburned/resume/main/resume.json";
  const resume = fetch(resumeURL)
    .then((r) => r.json())
    .catch((err) => {
      console.error("Failed to fetch resume:", err);
      throw err;
    });
</script>

<Window maxWidth="max(50vw, 700px)" maxHeight="max(90vh, 1000px)" scrollable>
  <h3 slot="title">About</h3>

  <div slot="overlay">
    <div class="toasts">
      <Toasts {toasts} toastClass="portfolio-toast" />
    </div>
  </div>

  <div class="portfolio-content">
    <section class="annoyance card">
      <b>Hey!!</b> You should totally check out the <b><u>xterm.js</u></b> window underneath!
    </section>

    <About />
    <Docs />
    <Links />

    {#await resume}
      <span class="loading">Hang on, I'm loading the rest of the resume!</span>
    {:then resume}
      <Resume {resume} />
    {:catch}
      <span class="loading">
        I couldn't load my resume {":("}
        <br />
        Maybe the console can help?
      </span>
    {/await}

    {#await jsonldDoc}
      <span class="loading">Give me a bit, I'm loading the rest!</span>
    {:then doc}
      <Webrings {doc} />
      <Badges88x31 {doc} />
    {:catch}
      <span class="loading">
        I couldn't load my JSON card-LD information either {":("}
        <br />
        Maybe the console can help?
      </span>
    {/await}
  </div>
</Window>

<style lang="scss">
  .toasts {
    margin: 1.5em 0;
  }

  .portfolio-content {
    padding-top: 2px;
    padding-bottom: 1em;
    margin: 0 auto;

    width: 100%;
    max-width: clamp(450px, 90vw, 650px);
    line-height: 1.5;

    position: relative;

    display: flex;
    gap: 1em;
    flex-direction: column;

    @media (max-width: 400px) {
      gap: 0.5em;
      padding: 0.5em 0;
    }

    & > * {
      margin-bottom: 1em;
    }
  }

  @mixin content {
    box-sizing: border-box;
    min-height: 32px; /* should be tall even when only containing a label */
    overflow: hidden;
  }

  @mixin content-item {
    margin: 0;
    padding: var(--adw-menu-padding);
    min-height: 32px;
    box-sizing: border-box;

    &:not(:last-child) {
      border-bottom: 1px solid var(--adw-card-shade-color);
    }
  }

  :global(section) {
    margin: 0 0.5em;
    box-sizing: border-box;

    :global(h1),
    :global(h2),
    :global(h3),
    :global(h4),
    :global(h5),
    :global(h6) {
      font-size: 1em;
    }

    & > :global(h1),
    & > :global(h2),
    & > :global(h3) {
      // line-height: 1.25;

      padding: 2px;
      padding-top: 18px;
      padding-bottom: 6px;

      margin: 0;
      margin-bottom: 6px;
    }

    @media (max-width: 400px) {
      font-size: 1em;
    }

    :global(button),
    :global(a[role="button"]) {
      background-color: var(--adw-card-bg-color);

      &:hover {
        background-color: var(--adw-button-hover-color);
      }
    }

    :global(.content) {
      @include content;
    }

    :global(.content-item) {
      @include content-item;
    }
  }

  section.annoyance {
    @include content;
    margin-bottom: 0;

    border: 1.5px solid rgba(var(--pink-rgb), 0.4);
    padding: var(--adw-card-padding);
    background-color: rgba(var(--pink-rgb), 0.1);
  }

  .loading {
    opacity: 0.5;
    font-size: 0.9em;
    text-align: center;
  }
</style>
