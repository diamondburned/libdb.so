<script lang="ts">
  import jsonld from "jsonld";
  import type { NodeObject } from "jsonld";

  import type { Badge } from "./Badges.svelte";
  import { nsfw } from "#/libdb.so/site/lib/prefs.js";
  import { ToastStore } from "#/libdb.so/site/lib/toasts.js";

  import Avatar from "#/libdb.so/site/assets/avatar.webp?url";
  import Banner from "#/libdb.so/site/assets/banner.webp?url";

  import Badges from "./Badges.svelte";
  import Toasts from "#/libdb.so/site/components/Toasts.svelte";
  import Window from "#/libdb.so/site/components/Window.svelte";
  import Webring from "#/libdb.so/site/components/Portfolio/Webring.svelte";
  import OpenInNew from "#/libdb.so/site/components/MaterialIcons/open_in_new.svelte";

  import GitHubIcon from "super-tiny-icons/images/svg/github.svg";
  import GitLabIcon from "super-tiny-icons/images/svg/gitlab.svg";
  import MastodonIcon from "super-tiny-icons/images/svg/mastodon.svg";
  import DiscordIcon from "super-tiny-icons/images/svg/discord.svg";
  import MatrixIcon from "super-tiny-icons/images/svg/matrix.svg";
  import LinkedInIcon from "super-tiny-icons/images/svg/linkedin.svg";
  import EmailIcon from "super-tiny-icons/images/svg/email.svg";

  const toasts = new ToastStore();

  const resumeURL = "https://raw.githubusercontent.com/diamondburned/resume/main/resume.json";
  const resume = fetch(resumeURL)
    .then((r) => r.json())
    .catch((err) => {
      console.error("Failed to fetch resume:", err);
      throw err;
    });

  type JSONLDDocument = {
    "@context": Record<string, string>;
    "@graph": NodeObject[];
  };

  const jsonldURL = "/_fs/0xd14.jsonld";
  const jsonldDoc = fetch(jsonldURL)
    .then((r) => r.json() as Promise<JSONLDDocument>)
    .then((r) => ({
      root: r,
      self: r["@graph"][0] as NodeObject & {
        // TypeScript. Shut up.
        [key: string]: any;
      },
    }))
    .catch((err) => {
      console.error("Failed to fetch self JSON-LD:", err);
      throw err;
    });

  async function usingContext<T = object>(
    doc: Awaited<typeof jsonldDoc>,
    object: any,
    vocab: string,
  ): Promise<T> {
    return (await jsonld.compact(
      {
        "@context": doc.root["@context"],
        ...object,
      },
      {
        "@vocab": vocab,
      },
    )) as T;
  }

  async function combined88x31s(doc: Awaited<typeof jsonldDoc>) {
    const x = await Promise.all(
      [
        ...doc.self["libdb:88x31"], //
        ...doc.self["libdb:other88x31"],
      ].map((o) => usingContext<Badge>(doc, o, "https://0xd14.id#88x31Badge/")),
    );
    console.log(x);
    return x;
  }

  type Link = {
    url?: string; // copy name to clipboard if not present
    name: string;
    value: string;
    color?: string; // RGB triplet, make sure value=99 (HSV)!
    iconURL: string;
    class?: string;
    hidden?: boolean;
  };

  function hexToTriplet(hex: string): string {
    hex = hex.slice(1);
    const r = parseInt(hex.slice(0, 2), 16);
    const g = parseInt(hex.slice(2, 4), 16);
    const b = parseInt(hex.slice(4, 6), 16);
    return `${r}, ${g}, ${b}`;
  }

  $: links = [
    {
      url: "https://blog.libdb.so",
      name: "Blog",
      value: "blog.libdb.so",
      color: "var(--pink-rgb)",
    },
    {
      url: "https://github.com/diamondburned",
      name: "GitHub",
      value: "@diamondburned",
      color: hexToTriplet("#984ffc"),
      iconURL: GitHubIcon,
      class: "icon-invert",
    },
    {
      url: "https://gitlab.com/diamondburned",
      name: "GitLab",
      value: "@diamondburned",
      color: hexToTriplet("#fca326"),
      iconURL: GitLabIcon,
    },
    {
      url: "https://tech.lgbt/@diamond",
      name: "Mastodon",
      value: "@diamond@tech.lgbt",
      color: hexToTriplet("#6263fc"),
      iconURL: MastodonIcon,
    },
    {
      url: "https://girlcock.club/@diamond",
      name: "Mastodon",
      value: "@diamond@girlcock.club",
      iconURL: MastodonIcon,
      color: hexToTriplet("#6263fc"),
      class: "nsfw",
      hidden: !$nsfw,
    },
    {
      name: "Discord",
      value: "@diamondburned",
      color: hexToTriplet("#849ffc"),
      iconURL: DiscordIcon,
    },
    {
      url: "https://matrix.to/#/@diamondburned:matrix.org",
      name: "Matrix",
      value: "@diamondburned:matrix.org",
      color: hexToTriplet("#fcfcfc"),
      iconURL: MatrixIcon,
    },
    {
      url: "https://www.linkedin.com/in/diamondnotburned",
      name: "LinkedIn",
      value: "Diamond Dinh",
      color: hexToTriplet("#00a6fc"),
      iconURL: LinkedInIcon,
      hidden: $nsfw,
    },
    {
      url: "mailto:diamond@libdb.so",
      name: "Email",
      value: "diamond@libdb.so",
      color: "var(--blue-rgb)",
      iconURL: EmailIcon,
    },
  ] as Link[];
</script>

<Window view="portfolio" maxWidth="max(50vw, 700px)" maxHeight="max(90vh, 1000px)" scrollable>
  <h3 slot="title">About</h3>

  <div slot="overlay">
    <div class="toasts">
      <Toasts {toasts} toastClass="portfolio-toast" />
    </div>
  </div>

  <div class="portfolio-content">
    <section class="banner" class:nsfw={$nsfw}>
      <img src={$nsfw ? "/_fs/.nsfw/banner.webp" : Banner} alt="Banner" />
    </section>

    <section class="about card">
      <div class="intro">
        <img src={$nsfw ? "/_fs/.nsfw/avatar.jpg" : Avatar} alt="Diamond" />
        <div>
          <span>Hi, I'm</span>
          <h1>Diamond!</h1>
        </div>
      </div>
      <p class="i-am">I am a:</p>
      <ul>
        {#if $nsfw}
          <li>
            <b class="text-pink-glow">Cat girlthing toy/object</b>
          </li>
        {/if}
        <li>
          <b>4th-year Computer Science major 👩🎓</b> and a
        </li>
        <li>
          <b>{">"}5 years Software Engineer 👩‍💻 🖥️</b>
        </li>
      </ul>
      <p>
        I consider myself the world's biggest "open source
        {#if $nsfw}<span class="text-pink-glow">slut</span>{:else}cheerleader{/if}"! I'm passionate
        about making the world a better place through technology and open source.
      </p>
    </section>

    <section class="annoyance card">
      <b>Hey!!</b> You should totally check out the <b><u>xterm.js</u></b> window underneath!
    </section>

    <section class="links">
      <h2>Links</h2>
      <div class="links-list card" role="list">
        {#each links.filter((link) => !link.hidden) as link}
          <a
            style={`--color: ${link.color};`}
            class={link.class ?? ""}
            role="button"
            href={link.url}
            target="_blank"
            on:click={(ev) => {
              if (!link.url) {
                ev.preventDefault();
                navigator.clipboard.writeText(link.value);
                toasts.add({ text: `Copied to clipboard!` }, 5000);
              }
            }}
          >
            <span class="icon">
              {#if link.iconURL}
                <img src={link.iconURL} alt={link.name} />
              {:else}
                <OpenInNew />
              {/if}
            </span>
            <span class="name">{link.name}</span>
            <span class="value">{link.value}</span>
          </a>
        {/each}
      </div>
    </section>

    <section class="resume">
      <h2>Resume</h2>
      <div class="links card" role="list">
        <a
          role="button"
          href="https://github.com/diamondburned/resume/blob/main/resume.pdf"
          target="_blank"
        >
          <OpenInNew /><span class="filename">PDF</span>
          <span class="source">(github.com)</span>
        </a>
        <a
          role="button"
          href="https://github.com/diamondburned/resume/blob/main/resume.json"
          style="--color: var(--pink-rgb);"
          target="_blank"
        >
          <OpenInNew /><span class="filename">JSON</span>
          <span class="source">(github.com)</span>
        </a>
      </div>
    </section>

    {#await resume}
      <span class="loading">Give me a bit, I'm loading the rest!</span>
    {:then resume}
      <section class="work hidden">
        <h2>Experience</h2>
        <ol class="work-list card">
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
        <ul class="projects-list card">
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
      <section class="webring">
        <h2>Webrings</h2>
        <div class="content card">
          {#each doc.self["libdb:webring"] as webring}
            {#await usingContext(doc, webring, "https://0xd14.id#Webring/") then webring}
              <Webring data={webring} />
            {/await}
          {/each}
        </div>
      </section>

      {#await combined88x31s(doc) then badges}
        <div class="badges">
          <Badges {badges} />
        </div>
      {/await}
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

  section {
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

    button,
    a[role="button"] {
      background-color: var(--adw-card-bg-color);

      &:hover {
        background-color: var(--adw-hover-color);
      }
    }
  }

  section.banner {
    height: clamp(150px, 20vw, 250px);

    &.nsfw {
      border-color: rgba(var(--pink-rgb), 0.4);
    }

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
    @include content;
    padding: 0 var(--adw-menu-padding);

    div.intro {
      display: flex;
      flex-direction: row;
      align-items: flex-end;
      gap: 1em;
      line-height: 1.15;
      margin-top: var(--adw-menu-padding);

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

    ul {
      list-style: "🌟  ";
      padding-left: 2em;
      margin: 0;
    }

    p.i-am {
      margin-bottom: var(--adw-menu-margin);
    }
  }

  section.annoyance {
    @include content;

    border: 1.5px solid rgba(var(--pink-rgb), 0.4);
    padding: var(--adw-card-padding);
    background-color: rgba(var(--pink-rgb), 0.1);
  }

  section.links {
    .links-list {
      @include content;

      list-style: none;

      width: 100%;
      display: grid;
      grid-template-columns: auto auto 1fr;
    }

    a[role="button"] {
      @include content-item;

      width: 100%;

      display: grid;
      grid-gap: 0.5em;
      grid-column: span 3;
      grid-template-columns: subgrid;
      /* grid-template-columns: auto auto 1fr auto; */

      &.nsfw .value::after {
        content: "(NSFW)";
        color: rgba(var(--pink-rgb), 0.75);
        text-shadow: var(--pink-glow);
        font-size: 0.75em;
        margin-left: 0.5em;
      }

      &.icon-invert .icon {
        filter: invert(1);
      }
    }

    .icon {
      user-select: none;
      min-width: 1.75em;
    }

    .name {
      font-weight: bold;
      margin-right: 0.5em;
    }

    :global(img),
    :global(svg) {
      width: 1.5em;
      height: 1.5em;
      vertical-align: bottom;
    }

    :global(svg) {
      color: rgb(var(--color));
    }

    :global(img) {
      border-radius: 40px;
    }

    @media (max-width: 400px) {
      .links-list {
        grid-template-columns: auto 1fr;
      }

      .name {
        display: none;
      }
    }
  }

  section.resume {
    .links {
      @include content;

      display: flex;
      flex-direction: row;

      a[role="button"] {
        @include content-item;

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

  section.work {
    .work-list {
      @include content;

      list-style: none;
      padding: 0;
      margin: 0;
    }

    .work-item {
      @include content-item;

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

  section.projects {
    .projects-list {
      @include content;

      list-style: none;
      padding: 0;
      margin: 0;
    }

    .project-item {
      @include content-item;

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

  section.webring {
    .content {
      @include content;
    }

    :global(webring-element) {
      @include content-item;
    }
  }

  .loading {
    opacity: 0.5;
    font-size: 0.9em;
    text-align: center;
  }

  .badges {
    margin: 0 0.5em;
    margin-top: 1em;
    padding: 1em;
    border-top: 1px solid #fff3;
  }
</style>
