<script lang="ts">
  import { nsfw } from "#/libdb.so/site/lib/prefs.js";
  import { ToastStore } from "#/libdb.so/site/lib/toasts.js";

  import Toasts from "#/libdb.so/site/components/Toasts.svelte";
  import Window from "#/libdb.so/site/components/Window.svelte";
  import OpenInNew from "#/libdb.so/site/components/MaterialIcons/open_in_new.svelte";

  import GitHubIcon from "super-tiny-icons/images/svg/github.svg";
  import GitLabIcon from "super-tiny-icons/images/svg/gitlab.svg";
  import MastodonIcon from "super-tiny-icons/images/svg/mastodon.svg";
  import DiscordIcon from "super-tiny-icons/images/svg/discord.svg";
  import MatrixIcon from "super-tiny-icons/images/svg/matrix.svg";
  import LinkedInIcon from "super-tiny-icons/images/svg/linkedin.svg";
  import EmailIcon from "super-tiny-icons/images/svg/email.svg";

  const toasts = new ToastStore();

  const resumeURL =
    "https://raw.githubusercontent.com/diamondburned/resume/main/resume.json";
  const resume = fetch(resumeURL)
    .then((r) => r.json())
    .catch((err) => {
      console.error("Failed to fetch resume:", err);
      throw err;
    });

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

<Window
  view="portfolio"
  maxWidth="max(50vw, 600px)"
  maxHeight="max(90vh, 1000px)"
  scrollable
>
  <h3 slot="title">About</h3>

  <div slot="overlay">
    <div class="toasts">
      <Toasts {toasts} toastClass="portfolio-toast" />
    </div>
  </div>

  <div class="portfolio-content">
    <section class="banner" class:nsfw={$nsfw}>
      <img
        src={$nsfw ? "/_fs/.nsfw/banner.webp" : "/_assets/banner.webp"}
        alt="Banner"
      />
    </section>

    <section class="about">
      <div class="intro">
        <img
          src={$nsfw ? "/_fs/.nsfw/avatar.jpg" : "/_assets/avatar.webp"}
          alt="Diamond"
        />
        <div>
          <span>Hi, I'm</span>
          <h1>Diamond!</h1>
        </div>
      </div>
      <p>
        I'm a
        {#if $nsfw}<b class="text-pink-glow">catgirl pet</b>,{/if}
        <b>4th-year Computer Science major 👩🎓</b>
        and past <b>Software Engineer Intern 👩‍💻 🖥️</b>
      </p>
      <p>
        I am what ChatGPT calls the world's biggest "open source
        {#if $nsfw}<span class="text-pink-glow">slut</span
          >{:else}cheerleader{/if}"!
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
      <div class="links-list" role="list">
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

    <webring-element
      src="https://raw.githubusercontent.com/diamondburned/acmfriends-webring/%3C3-spring-2023/webring.json"
      name="diamond"
    >
      <section class="webring">
        <span class="ring" />
        <div>
          <a class="left" target="_blank" />
          <span class="middle" />
          <a class="right" target="_blank" />
        </div>
      </section>
    </webring-element>
  </div>
</Window>

<style lang="scss">
  .toasts {
    margin: 1.5em 0;
  }

  .portfolio-content {
    padding: 1em 0;
    margin: 0 auto;

    width: 100%;
    max-width: clamp(400px, 80vw, 550px);
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

      :global(h1),
      :global(h2),
      :global(h3),
      :global(h4),
      :global(h5),
      :global(h6) {
        margin: 1rem 0;
        font-family: "Nunito";
      }

      :global(h1),
      :global(h2),
      :global(h3) {
        line-height: 1.25;
      }

      & > * {
        margin: 1rem 0;
      }

      a {
        text-decoration: none;
        color: var(--blue);

        &:hover {
          text-decoration: underline;
        }
      }

      a[role="button"] {
        --color: 255 255 255;

        text-decoration: none;
        padding: 0.5em;
        border-radius: 5px;
        transition: all 0.1s ease-in-out;

        outline: 1px solid rgba(var(--color), 0.4);
        outline-offset: -1px;

        box-shadow: 0 0 0.35em -0.1em rgba(var(--color), 0.5);
        background-color: rgba(var(--color), 0.1);
        &:hover {
          box-shadow: 0 0 0.45em -0.1em rgba(var(--color), 0.8);
          background-color: rgba(var(--color), 0.2);
        }
      }

      @media (max-width: 400px) {
        font-size: 1em;
      }
    }

    section.banner {
      padding: 0;
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

    section.links {
      .links-list {
        list-style: none;
        padding: 0;

        width: 100%;
        display: grid;
        grid-template-columns: auto auto 1fr;
        grid-gap: 0.5em;
      }

      a[role="button"] {
        width: 100%;
        cursor: pointer;
        box-sizing: border-box;

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
      }

      .name {
        font-weight: bold;
        margin-right: 0.5em;
        color: rgb(var(--color));
      }

      .value {
        color: rgba(255, 255, 255, 0.75);
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
      & > div {
        display: flex;
        flex-direction: row;
      }

      a[role="button"] {
        --color: var(--blue-rgb);
        width: 100%;
        align-self: center;
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
            opacity: 0.65;
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

    webring-element {
      section.webring {
        display: flex;
        padding: 1em;
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

    .loading {
      opacity: 0.5;
      font-size: 0.9em;
      text-align: center;
    }

    footer {
      margin: 0 calc(0.5em + 1px);
      border-radius: 10px;
      border: 1px solid rgba(var(--blue-rgb), 0.4);
      background-color: rgba(var(--blue-rgb), 0.1);
      padding: 1em;
      font-size: 0.9em;
    }
  }
</style>
