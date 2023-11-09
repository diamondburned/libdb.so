<script lang="ts">
  import Window from "#/libdb.so/site/components/Window.svelte";
</script>

<Window id="portfolio">
  <h3 slot="title">Portfolio</h3>

  <div class="portfolio-content">
    <section class="intro">
      <span>Hi, I'm</span>
      <h2>Diamond!</h2>
    </section>
  </div>
</Window>

<style lang="scss">
  .portfolio-content {
    overflow-x: hidden;
    overflow-y: auto;

    padding: clamp(0.5em, 1.5vw, 1.5em);

    max-width: clamp(500px, 80vw, 800px);
    margin: auto;

    .intro {
      span {
        font-weight: lighter;
        font-size: 1.5em;
      }
      h2 {
        font-weight: bold;
        font-size: 2.5em;
        margin: 0;
      }
    }
  }
</style>
