<script lang="ts">
  import "../app.css";
  import AppShell from "$lib/components/AppShell.svelte";
  import { page } from "$app/stores";
</script>

<AppShell>
  {#if $page.url.pathname === "/"}
    <slot />
  {:else}
    <div class="routePane">
      <slot />
    </div>
  {/if}
</AppShell>

<style>
  /* Fix: non-terminal pages must not float 100% transparent on the Matrix */
  .routePane {
    width: 100%;
    height: 100%;

    background: var(--surface, rgba(15, 20, 25, 0.35));
    border: 1px solid var(--border, rgba(255, 255, 255, 0.10));
    border-radius: var(--radius, 14px);

    backdrop-filter: blur(var(--blur, 18px));
    -webkit-backdrop-filter: blur(var(--blur, 18px));
  }
</style>
