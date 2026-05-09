<script lang="ts">
  import "../app.css";
  import AppShell from "$lib/components/AppShell.svelte";
  import { page } from "$app/stores";
</script>

<svelte:head>
  <meta charset="utf-8" />
</svelte:head>

<AppShell>
  {#if $page.url.pathname === "/"}
    <slot />
  {:else}
    <div class="routePane">
      <div class="routeContent">
        <slot />
      </div>
    </div>
  {/if}
</AppShell>

<style>
  /* Fix: non-terminal pages must not float 100% transparent on the Matrix */
  .routePane {
    width: 100%;
    height: 100%;
    min-height: 0;

    display: grid;
    place-items: center;

    background: var(--surface, rgba(15, 20, 25, 0.35));
    border: 1px solid var(--border, rgba(255, 255, 255, 0.10));
    border-radius: var(--radius, 14px);

    backdrop-filter: blur(var(--blur, 18px));
    -webkit-backdrop-filter: blur(var(--blur, 18px));

    overflow: auto;
    padding: clamp(10px, 1.8vw, 20px);
    box-sizing: border-box;
  }

  .routeContent {
    width: min(1220px, 100%);
    min-height: min-content;
    margin: auto;
    overflow: visible;
  }

  @media (max-width: 700px) {
    .routePane {
      padding: 8px;
      place-items: start center;
    }

    .routeContent {
      width: 100%;
      margin: 0;
    }
  }
</style>
