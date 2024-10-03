<script lang="ts">
  import { nsfw } from "#/libdb.so/site/lib/prefs.js";
  import { pronouns, selectPronouns } from "#/libdb.so/site/lib/pronouns.js";

  import Avatar from "#/libdb.so/public/_fs/pics/portfolio/avatar.webp?url";
  import Banner from "#/libdb.so/public/_fs/pics/portfolio/banner.webp?url";

  $: pronoun = (content: string) => selectPronouns(content, $pronouns);
</script>

<section class="banner" class:nsfw={$nsfw}>
  <img src={$nsfw ? "/_fs/.nsfw/banner.webp" : Banner} alt="Banner" />
</section>

<section class="about card content">
  <div class="intro">
    <img src={$nsfw ? "/_fs/.nsfw/avatar.jpg" : Avatar} alt="Diamond" />
    <div>
      <span>Hi, this is</span>
      <h1>
        Diamond{#if $nsfw}
          <span class="text-pink-glow">:3</span>{:else}!{/if}
      </h1>
    </div>
  </div>
  <p class="i-am">
    {pronoun("She")} is a <b>software engineer at Google</b> with over
    <b>5 years of experience</b>.
  </p>
  <p>
    {pronoun("She")} considers {pronoun("herself")} the world's biggest "open source
    {#if $nsfw}<span class="text-pink-glow">slut</span>{:else}cheerleader{/if}"!
    {pronoun("She's")} very passionate about making the world a better place through technology and open
    source!
  </p>
</section>

<style lang="scss">
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

    p.i-am {
      margin-bottom: var(--adw-menu-margin);
    }
  }
</style>
