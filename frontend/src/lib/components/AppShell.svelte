<script lang="ts">
  import { browser } from "$app/environment";
  import { onMount } from "svelte";
  import "$lib/appshell.css";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  import Matrix from "./Matrix.svelte";
  import WindowManager from "./wm/WindowManager.svelte";

  const props = $props<{ children?: import("svelte").Snippet }>();

  type WSKey = "terminal" | "about" | "experience" | "work" | "skills" | "contact";

  const workspaces: { key: WSKey; n: number; label: string; href: string }[] = [
    { key: "terminal", n: 1, label: "Terminal", href: "/" },
    { key: "about", n: 2, label: "About", href: "/about" },
    { key: "experience", n: 3, label: "Experience", href: "/experience" },
    { key: "work", n: 4, label: "Work", href: "/work" },
    { key: "skills", n: 5, label: "Skills", href: "/skills" },
    { key: "contact", n: 6, label: "Contact", href: "/contact" }
  ];

  const routeToWS: Record<string, WSKey> = {
    "/": "terminal",
    "/about": "about",
    "/experience": "experience",
    "/work": "work",
    "/skills": "skills",
    "/contact": "contact"
  };

  const active = $derived(routeToWS[$page.url.pathname] ?? "terminal");
  const activeLang = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr") ? "fr" : "en");
  const isFr = $derived(activeLang === "fr");

  const workspaceHints = $derived.by(() => isFr
    ? {
        terminal: "Gestionnaire de fenêtres",
        about: "Résumé du profil",
        experience: "Parcours et rôles",
        work: "Projets et dépôts",
        skills: "Moniteur de compétences",
        contact: "Me contacter"
      }
    : {
        terminal: "Window manager",
        about: "Profile summary",
        experience: "Timeline and roles",
        work: "Projects and repos",
        skills: "Live skill monitor",
        contact: "Reach out"
      });

  const langStateLabel = $derived(isFr ? `${activeLang.toUpperCase()} sélectionnée` : `${activeLang.toUpperCase()} selected`);
  const kbHintLabel = $derived.by(() => isFr
    ? {
      terminal: "Terminal",
      runCommands: "exécutez des commandes comme",
      switchWs: "Changer d'espace",
      switchSuffix: "(ou Ctrl+Shift+1..6)",
      joinOr: "ou"
    }
  : {
      terminal: "Terminal",
      runCommands: "run commands like",
      switchWs: "Switch WS",
      switchSuffix: "(or Ctrl+Shift+1..6)",
      joinOr: "or"
    });

  const navSourceStorageKey = "kagha:last-ws-source";
  let showKbHint = $state(false);
  let hintClosing = $state(false);

  function dismissHint() {
    if (!showKbHint || hintClosing) return;
    hintClosing = true;
    setTimeout(() => {
      showKbHint = false;
      hintClosing = false;
    }, 280);
  }

  function openWS(ws: WSKey) {
    const w = workspaces.find((x) => x.key === ws);
    if (!w) return;
    if (w.key === active) return;

    if (browser) {
      const cur = workspaces.find((x) => x.key === active);
      if (cur) {
        try {
          window.sessionStorage.setItem(navSourceStorageKey, cur.label);
        } catch {
          // ignore storage errors
        }
      }
    }

    const href = activeLang === "fr" ? `${w.href}?lang=fr` : w.href;
    void goto(href, { keepFocus: true, noScroll: true, invalidateAll: false });
  }

  function switchLang(lang: "en" | "fr") {
    if (lang === activeLang) return;
    const target = lang === "fr" ? `${$page.url.pathname}?lang=fr` : $page.url.pathname;
    void goto(target, { keepFocus: true, noScroll: true, invalidateAll: true });
  }

  onMount(() => {
    showKbHint = true;
    const auto = setTimeout(dismissHint, 5200);

    const onKey = (e: KeyboardEvent) => {
      if (e.defaultPrevented) return;

      if (e.key === "Escape") {
        dismissHint();
        return;
      }

      const k = e.key;
      if (!/^[1-6]$/.test(k)) return;

      const byAlt = e.altKey && !e.ctrlKey && !e.shiftKey && !e.metaKey;
      const byCtrlShift = e.ctrlKey && e.shiftKey && !e.altKey && !e.metaKey;
      if (!byAlt && !byCtrlShift) return;

      e.preventDefault();
      const n = Number(k);
      const ws = workspaces.find((w) => w.n === n);
      if (!ws) return;

      dismissHint();
      openWS(ws.key);
    };

    window.addEventListener("keydown", onKey, { capture: true });

    return () => {
      clearTimeout(auto);
      window.removeEventListener("keydown", onKey, { capture: true });
    };
  });

  const isWM = $derived($page.url.pathname === "/");
</script>

<div class="wm">
  <Matrix />

  <header class="bar">
    <div class="barLeft">
      <div class="brand" title="NixOS">
        <span class="brandIcon" aria-hidden="true">
          <svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <linearGradient id="nixRed" x1="8" y1="8" x2="56" y2="56" gradientUnits="userSpaceOnUse">
                <stop offset="0" stop-color="#ff5a5a" />
                <stop offset="0.55" stop-color="#e45a5a" />
                <stop offset="1" stop-color="#ff5a5a" />
              </linearGradient>
            </defs>
            <path
              d="M32 6l6.3 10.9h-12.6L32 6z"
              fill="url(#nixRed)"
              opacity="0.9"
            />
            <path
              d="M32 58l-6.3-10.9h12.6L32 58z"
              fill="url(#nixRed)"
              opacity="0.9"
            />
            <path
              d="M6 32l10.9-6.3v12.6L6 32z"
              fill="url(#nixRed)"
              opacity="0.9"
            />
            <path
              d="M58 32l-10.9 6.3V25.7L58 32z"
              fill="url(#nixRed)"
              opacity="0.9"
            />
            <path
              d="M14.9 14.9l12.1 3.2-8.9 8.9-3.2-12.1z"
              fill="url(#nixRed)"
              opacity="0.75"
            />
            <path
              d="M49.1 49.1l-12.1-3.2 8.9-8.9 3.2 12.1z"
              fill="url(#nixRed)"
              opacity="0.75"
            />
            <path
              d="M49.1 14.9l-3.2 12.1-8.9-8.9 12.1-3.2z"
              fill="url(#nixRed)"
              opacity="0.75"
            />
            <path
              d="M14.9 49.1l3.2-12.1 8.9 8.9-12.1 3.2z"
              fill="url(#nixRed)"
              opacity="0.75"
            />
          </svg>
        </span>
      </div>
    </div>

    <nav class="barCenter">
      {#each workspaces as ws (ws.key)}
        <button
          type="button"
          class={"ws " + (ws.key === active ? "wsActive" : "")}
          onclick={() => openWS(ws.key)}
          aria-current={ws.key === active ? "page" : undefined}
          title={workspaceHints[ws.key]}
        >
          <span class="wsIcon" aria-hidden="true">
            {#if ws.key === "terminal"}
              <svg viewBox="0 0 24 24" fill="none">
                <rect x="3.5" y="4.5" width="17" height="15" rx="2.5" stroke="currentColor" stroke-width="1.5" />
                <path d="M7 9l3 3-3 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                <path d="M12.5 15h4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
              </svg>
            {:else if ws.key === "about"}
              <svg viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="8" r="3" stroke="currentColor" stroke-width="1.5" />
                <path d="M6 18c1.2-2.4 3.2-3.6 6-3.6s4.8 1.2 6 3.6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
              </svg>
            {:else if ws.key === "experience"}
              <svg viewBox="0 0 24 24" fill="none">
                <rect x="4" y="7" width="16" height="11" rx="2" stroke="currentColor" stroke-width="1.5" />
                <path d="M9 7V5.5A1.5 1.5 0 0110.5 4h3A1.5 1.5 0 0115 5.5V7" stroke="currentColor" stroke-width="1.5" />
              </svg>
            {:else if ws.key === "work"}
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M3.5 8.5A2.5 2.5 0 016 6h4l1.2 1.5H18A2.5 2.5 0 0120.5 10v7A2.5 2.5 0 0118 19.5H6A2.5 2.5 0 013.5 17V8.5z" stroke="currentColor" stroke-width="1.5" />
              </svg>
            {:else if ws.key === "skills"}
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M6 17V9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                <path d="M12 17V6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                <path d="M18 17v-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                <circle cx="6" cy="8" r="1.5" fill="currentColor" />
                <circle cx="12" cy="11" r="1.5" fill="currentColor" />
                <circle cx="18" cy="12" r="1.5" fill="currentColor" />
              </svg>
            {:else}
              <svg viewBox="0 0 24 24" fill="none">
                <rect x="4" y="6" width="16" height="12" rx="2" stroke="currentColor" stroke-width="1.5" />
                <path d="M5.5 8l6.5 5L18.5 8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            {/if}
          </span>
          <span class="wsNum" aria-hidden="true">{ws.n}</span>
          <span class="wsLabel">{ws.label}</span>
        </button>
      {/each}
    </nav>

    <div class="barRight">
      <div class="langSwitch" role="group" aria-label="Language">
        <button type="button" class={"langBtn " + (activeLang === "en" ? "active" : "")} onclick={() => switchLang("en")}>EN</button>
        <button type="button" class={"langBtn " + (activeLang === "fr" ? "active" : "")} onclick={() => switchLang("fr")}>FR</button>
      </div>
      <span class="langState">{langStateLabel}</span>
      <a class="link" href="https://linkedin.com/in/ghassenelkamel" target="_blank" rel="noreferrer">linkedin</a>
      <a class="link" href="https://github.com/ghassenelkamel" target="_blank" rel="noreferrer">github</a>
    </div>
  </header>

  <main class="content">
    {#if showKbHint}
      <div class={"kbHint " + (hintClosing ? "out" : "")}> 
        <div class="hintLine"><span class="hintKey">{kbHintLabel.terminal}</span> {kbHintLabel.runCommands} <span class="hintCmd">./AboutMe</span> {kbHintLabel.joinOr} <span class="hintCmd">exec Work</span></div>
        <div class="hintLine"><span class="hintKey">{kbHintLabel.switchWs}</span> <span class="hintCmd">Alt+1..6</span> <span class="hintMuted">{kbHintLabel.switchSuffix}</span></div>
      </div>
    {/if}

    {#if isWM}
      <WindowManager />
    {:else}
      <div class="single">
        {@render props.children?.()}
      </div>
    {/if}
  </main>
</div>
