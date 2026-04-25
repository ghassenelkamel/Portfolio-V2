<script lang="ts">
  import "$lib/appshell.css";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  import Matrix from "./Matrix.svelte";
  import WindowManager from "./wm/WindowManager.svelte";

  const props = $props<{ children?: import("svelte").Snippet }>();

  type WSKey = "terminal" | "about" | "experience" | "work" | "skills" | "contact";

  const workspaces: { key: WSKey; n: string; label: string; href: string }[] = [
    { key: "terminal", n: "1", label: "wm", href: "/" },
    { key: "about", n: "2", label: "about", href: "/about" },
    { key: "experience", n: "3", label: "exp", href: "/experience" },
    { key: "work", n: "4", label: "work", href: "/work" },
    { key: "skills", n: "5", label: "skills", href: "/skills" },
    { key: "contact", n: "6", label: "contact", href: "/contact" }
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

  function openWS(ws: WSKey) {
    const w = workspaces.find((x) => x.key === ws);
    if (!w) return;
    void goto(w.href);
  }

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
        >
          <span class="wsNum">{ws.n}</span>
          <span class="wsLabel">{ws.label}</span>
        </button>
      {/each}
    </nav>

    <div class="barRight">
      <a class="link" href="https://linkedin.com/in/ghassenelkamel" target="_blank" rel="noreferrer">linkedin</a>
      <a class="link" href="https://github.com/ghassenelkamel" target="_blank" rel="noreferrer">github</a>
    </div>
  </header>

  <main class="content">
    {#if isWM}
      <WindowManager />
    {:else}
      <div class="single">
        {@render props.children?.()}
      </div>
    {/if}
  </main>
</div>
