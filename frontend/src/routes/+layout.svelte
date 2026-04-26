<script lang="ts">
  import { browser } from "$app/environment";
  import "../app.css";
  import AppShell from "$lib/components/AppShell.svelte";
  import Terminal from "$lib/components/Terminal.svelte";
  import { page } from "$app/stores";

  function hasDockedTerminal(pathname: string) {
    return pathname === "/about"
      || pathname === "/experience"
      || pathname === "/work"
      || pathname === "/skills"
      || pathname === "/contact";
  }

  function workspaceMeta(pathname: string) {
    switch (pathname) {
      case "/": return { n: 1, name: "Terminal" };
      case "/about": return { n: 2, name: "About" };
      case "/experience": return { n: 3, name: "Experience" };
      case "/work": return { n: 4, name: "Work" };
      case "/skills": return { n: 5, name: "Skills" };
      case "/contact": return { n: 6, name: "Contact" };
      default: return null;
    }
  }

  const navSourceStorageKey = "kagha:last-ws-source";
  let dockSuggest = $state<string | undefined>(undefined);
  let showRouteContent = $state(true);

  $effect(() => {
    $page.url.pathname;
    showRouteContent = true;
  });

  $effect(() => {
    const current = workspaceMeta($page.url.pathname);
    if (!browser || !current) {
      dockSuggest = undefined;
      return;
    }

    try {
      const source = (window.sessionStorage.getItem(navSourceStorageKey) ?? "").trim();
      if (!source || source.toLowerCase() === current.name.toLowerCase()) {
        dockSuggest = undefined;
        return;
      }
      dockSuggest = `ws ${source}`;
    } catch {
      dockSuggest = undefined;
    }
  });

  function onRouteTerminalAction(action: "hide-content" | "reset-layout" | "show-content") {
    if (action === "hide-content") {
      showRouteContent = false;
      return;
    }

    if (action === "show-content" || action === "reset-layout") {
      showRouteContent = true;
    }
  }
</script>

<AppShell>
  {#if $page.url.pathname === "/"}
    <slot />
  {:else}
    <div class={"routePane" + (showRouteContent ? "" : " onlyTerm")}>
      {#if showRouteContent}
        <div class="routeContent">
          <slot />
        </div>
      {/if}

      {#if hasDockedTerminal($page.url.pathname)}
        <div class="routeTerm">
          <Terminal
            visible={true}
            showBanner={false}
            embedded={true}
            suggestCommand={dockSuggest}
            onLayoutAction={onRouteTerminalAction}
            sessionKey={`route:${$page.url.pathname}:dock`}
          />
        </div>
      {/if}
    </div>
  {/if}
</AppShell>

<style>
  /* Fix: non-terminal pages must not float 100% transparent on the Matrix */
  .routePane {
    width: 100%;
    height: auto;
    min-height: 100%;

    display: flex;
    flex-direction: column;
    gap: 0;

    background: var(--surface, rgba(15, 20, 25, 0.35));
    border: 1px solid var(--border, rgba(255, 255, 255, 0.10));
    border-radius: var(--radius, 14px);

    backdrop-filter: blur(var(--blur, 18px));
    -webkit-backdrop-filter: blur(var(--blur, 18px));

    overflow: visible;
    padding: 8px;
    box-sizing: border-box;
  }

  .routeContent {
    min-height: 0;
    flex: 1 1 auto;
    overflow: visible;
  }

  .routePane.onlyTerm {
    justify-content: flex-end;
  }

  .routeTerm {
    min-height: 56px;
    height: auto;
    overflow: visible;
    border-radius: 0;
    border-top: 0;
  }

  @media (max-width: 700px) {
    .routePane {
      gap: 0;
      padding: 6px;
    }

    .routeTerm {
      min-height: 52px;
      height: auto;
    }
  }
</style>
