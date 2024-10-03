<script lang="ts" context="module">
  import { nsfw } from "#/libdb.so/site/lib/prefs.js";
  import { derived } from "svelte/store";

  import GitHubIcon from "super-tiny-icons/images/svg/github.svg";
  import GitLabIcon from "super-tiny-icons/images/svg/gitlab.svg";
  import MastodonIcon from "super-tiny-icons/images/svg/mastodon.svg";
  import DiscordIcon from "super-tiny-icons/images/svg/discord.svg";
  import MatrixIcon from "super-tiny-icons/images/svg/matrix.svg";
  import LinkedInIcon from "super-tiny-icons/images/svg/linkedin.svg";
  import EmailIcon from "super-tiny-icons/images/svg/email.svg";

  type Link = {
    url?: string; // copy name to clipboard if not present
    name: string;
    value: string;
    color?: string; // RGB triplet, make sure value=99 (HSV)!
    iconURL?: string;
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

  function asLinks(links: Link[]): Link[] {
    return links;
  }

  export const links = derived([nsfw], ([nsfw]) =>
    asLinks([
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
        hidden: !nsfw,
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
        hidden: nsfw,
      },
      {
        url: "mailto:diamond@libdb.so",
        name: "Email",
        value: "diamond@libdb.so",
        color: "var(--blue-rgb)",
        iconURL: EmailIcon,
      },
    ]),
  );
</script>

<script lang="ts">
  import OpenInNew from "#/libdb.so/site/components/MaterialIcons/open_in_new.svelte";
</script>

<section class="links">
  <h2>Links</h2>
  <div class="links-list card content" role="list">
    {#each $links.filter((link) => !link.hidden) as link}
      <a
        style={`--color: ${link.color};`}
        class="content-item {link.class ?? ''}"
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

<style lang="scss">
  section.links {
    .links-list {
      list-style: none;

      width: 100%;
      display: grid;
      grid-template-columns: auto auto 1fr;
    }

    a[role="button"] {
      width: 100%;

      display: grid;
      grid-gap: 0.5em;
      grid-column: span 3;
      grid-template-columns: subgrid;

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
</style>
