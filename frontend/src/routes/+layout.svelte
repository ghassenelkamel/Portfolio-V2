<script lang="ts">
  import "../app.css";
  import AppShell from "$lib/components/AppShell.svelte";
  import { page } from "$app/stores";

  const siteOrigin = "https://ghassenelkamel.fr";

  const routeTitleMap: Record<string, string> = {
    "/": "Ghassen El Kamel",
    "/about": "About",
    "/experience": "Experience",
    "/work": "Work",
    "/skills": "Skills",
    "/contact": "Contact"
  };

  const routeDescriptionMap: Record<string, string> = {
    "/": "IoT and software systems engineer focused on secure backend infrastructure, embedded systems, and production reliability.",
    "/about": "Profile summary, specialties, and current goals in IoT, backend engineering, and infrastructure delivery.",
    "/experience": "Professional journey across IoT engineering, security, backend systems, and product delivery roles.",
    "/work": "Project portfolio with repositories, technology stacks, and recent engineering activity.",
    "/skills": "Live capability monitor covering core backend, infrastructure, and IoT technology skills.",
    "/contact": "Professional contact channels for collaboration, project scope discussions, and engineering opportunities."
  };

  const activeTitle = $derived(routeTitleMap[$page.url.pathname] ?? routeTitleMap["/"]);
  const activeDescription = $derived(routeDescriptionMap[$page.url.pathname] ?? routeDescriptionMap["/"]);
  const seoTitle = $derived(`${activeTitle} | Ghassen El Kamel`);
  const canonicalPath = $derived.by(() => {
    const lang = ($page.url.searchParams.get("lang") || "").toLowerCase();
    return lang.startsWith("fr") ? `${$page.url.pathname}?lang=fr` : $page.url.pathname;
  });
  const canonicalUrl = $derived(`${siteOrigin}${canonicalPath}`);
</script>

<svelte:head>
  <meta charset="utf-8" />
  <title>{seoTitle}</title>
  <meta name="description" content={activeDescription} />
  <link rel="canonical" href={canonicalUrl} />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content="Ghassen El Kamel" />
  <meta property="og:url" content={canonicalUrl} />
  <meta property="og:title" content={seoTitle} />
  <meta property="og:description" content={activeDescription} />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content={seoTitle} />
  <meta name="twitter:description" content={activeDescription} />
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
