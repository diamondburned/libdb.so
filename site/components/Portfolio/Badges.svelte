<script lang="ts">
  export let badges: {
    image: string;
    link?: string;
    alt?: string;
  }[];
</script>

<div class="badges">
  {#each badges as badge}
    <a class="link" href={badge.link}>
      <img class="badge" src={badge.image} alt={badge.alt} />
    </a>
  {/each}
</div>

<style lang="scss">
  .badges {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(88px, 1fr));
    gap: 0.25em;
    justify-content: center;
  }

  .badge {
    aspect-ratio: 88 / 31;
    image-rendering: pixelated;
  }
</style>
