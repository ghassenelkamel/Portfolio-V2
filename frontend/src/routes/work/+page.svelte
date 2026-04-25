<script lang="ts">
  type Repo = {
    id: number;
    name: string;
    html_url: string;
    description: string | null;
    language: string | null;
    fork: boolean;
    archived: boolean;
    stargazers_count: number;
    forks_count: number;
    pushed_at?: string;
    updated_at?: string;
  };

  type Rate = { limit: string; remaining: string; reset: string };

  let { data } = $props<{ data: any }>();

  // IMPORTANT: use $derived.by for functions
  let repos = $derived.by((): Repo[] => {
    const d = data ?? {};
    const list = d?.repos;
    return Array.isArray(list) ? (list as Repo[]) : [];
  });

  let error = $derived.by((): string | null => {
    const d = data ?? {};
    const e = d?.error;
    return typeof e === "string" && e.trim() ? e : null;
  });

  let rate = $derived.by((): Rate | null => {
    const d = data ?? {};
    const r = d?.rate;
    if (!r || typeof r !== "object") return null;
    if (typeof r.limit !== "string" || typeof r.remaining !== "string" || typeof r.reset !== "string") return null;
    return r as Rate;
  });

  let q = $state("");
  let sort = $state<"recent" | "stars" | "name">("recent");
  let includeForks = $state(true);
  let includeArchived = $state(false);

  function safeLower(s: unknown) {
    return typeof s === "string" ? s.toLowerCase() : "";
  }

  function dateKey(r: Repo) {
    const s = r.pushed_at ?? r.updated_at ?? "";
    const t = Date.parse(s);
    return Number.isFinite(t) ? t : 0;
  }

  function ago(iso?: string) {
    if (!iso) return "—";
    const t = Date.parse(iso);
    if (!Number.isFinite(t)) return iso;

    const diff = Date.now() - t;
    const sec = Math.floor(diff / 1000);
    if (sec < 60) return `${sec}s ago`;
    const min = Math.floor(sec / 60);
    if (min < 60) return `${min}m ago`;
    const h = Math.floor(min / 60);
    if (h < 48) return `${h}h ago`;
    const d = Math.floor(h / 24);
    if (d < 30) return `${d}d ago`;
    const mo = Math.floor(d / 30);
    if (mo < 18) return `${mo}mo ago`;
    const y = Math.floor(mo / 12);
    return `${y}y ago`;
  }

  let filteredSorted = $derived.by(() => {
    const query = q.trim().toLowerCase();

    let list = repos;

    if (!includeForks) list = list.filter((r) => !r.fork);
    if (!includeArchived) list = list.filter((r) => !r.archived);

    if (query) {
      list = list.filter((r) => {
        const hay =
          `${r.name} ${r.description ?? ""} ${r.language ?? ""}`.toLowerCase();
        return hay.includes(query);
      });
    }

    const out = [...list];
    out.sort((a, b) => {
      if (sort === "name") return a.name.localeCompare(b.name);
      if (sort === "stars") return (b.stargazers_count ?? 0) - (a.stargazers_count ?? 0);
      // recent
      return dateKey(b) - dateKey(a);
    });

    return out;
  });

  let languages = $derived.by(() => {
    // top 6 languages (non-null)
    const m = new Map<string, number>();
    for (const r of repos) {
      if (!r.language) continue;
      m.set(r.language, (m.get(r.language) ?? 0) + 1);
    }
    return [...m.entries()].sort((a, b) => b[1] - a[1]).slice(0, 6);
  });

  let totalStars = $derived.by(() => repos.reduce((acc, r) => acc + (r.stargazers_count ?? 0), 0));

  function resetFilters() {
    q = "";
    sort = "recent";
    includeForks = true;
    includeArchived = false;
  }
</script>

<svelte:head>
  <title>Work</title>
</svelte:head>

<div class="wrap">
  <header class="hero">
    <div class="titleBlock">
      <div class="k">work --repos</div>
      <h1>Work</h1>
      <p class="sub">
        Live GitHub repository grid. Filter + sort locally. Open any card in a new tab.
      </p>
    </div>

    <div class="stats">
      <div class="stat">
        <div class="statLabel">Repos</div>
        <div class="statVal">{repos.length}</div>
      </div>
      <div class="stat">
        <div class="statLabel">Stars</div>
        <div class="statVal">{totalStars}</div>
      </div>
      <div class="stat">
        <div class="statLabel">Showing</div>
        <div class="statVal">{filteredSorted.length}</div>
      </div>
      {#if rate}
        <div class="stat muted">
          <div class="statLabel">API</div>
          <div class="statVal">{rate.remaining}/{rate.limit}</div>
        </div>
      {/if}
    </div>
  </header>

  <section class="controlBar">
    <div class="searchWrap">
      <span class="searchIcon" aria-hidden="true">⌕</span>
      <input
        class="search"
        type="text"
        placeholder="Search (name / description / language)…"
        bind:value={q}
        spellcheck="false"
      />
      {#if q.trim().length}
        <button type="button" class="miniBtn" onclick={() => (q = "")}>clear</button>
      {/if}
    </div>

    <div class="controls">
      <div class="selectWrap">
        <span class="lbl">sort</span>
        <select class="select" bind:value={sort} aria-label="Sort repositories">
          <option value="recent">recent</option>
          <option value="stars">stars</option>
          <option value="name">name</option>
        </select>
      </div>

      <button
        type="button"
        class={"chip " + (includeForks ? "on" : "")}
        onclick={() => (includeForks = !includeForks)}
        aria-pressed={includeForks}
        title="Toggle fork repositories"
      >
        forks
      </button>

      <button
        type="button"
        class={"chip " + (includeArchived ? "on" : "")}
        onclick={() => (includeArchived = !includeArchived)}
        aria-pressed={includeArchived}
        title="Toggle archived repositories"
      >
        archived
      </button>

      <button type="button" class="chip ghost" onclick={resetFilters}>reset</button>
    </div>
  </section>

  {#if languages.length}
    <section class="langRow" aria-label="Top languages">
      {#each languages as [lang, n] (lang)}
        <span class="langChip" title={`${n} repos`}>
          <span class="dot" aria-hidden="true"></span>
          {lang}
          <span class="count">{n}</span>
        </span>
      {/each}
    </section>
  {/if}

  {#if error}
    <div class="msg err">
      <div class="msgTitle">API error</div>
      <div class="msgBody">{error}</div>
    </div>
  {:else if repos.length === 0}
    <div class="msg">
      <div class="msgTitle">Loading</div>
      <div class="msgBody">Fetching repositories…</div>
      <div class="skeletonGrid">
        {#each Array(8) as _, i (i)}
          <div class="skCard"></div>
        {/each}
      </div>
    </div>
  {:else if filteredSorted.length === 0}
    <div class="msg">
      <div class="msgTitle">No match</div>
      <div class="msgBody">No repositories match your current filters.</div>
    </div>
  {:else}
    <div class="grid" role="list">
      {#each filteredSorted as r (r.id)}
        <a class="card" role="listitem" href={r.html_url} target="_blank" rel="noreferrer">
          <div class="cardTop">
            <div class="nameRow">
              <div class="name">{r.name}</div>
              <div class="badges">
                {#if r.language}<span class="badge">{r.language}</span>{/if}
                {#if r.fork}<span class="badge dim">fork</span>{/if}
                {#if r.archived}<span class="badge warn">archived</span>{/if}
              </div>
            </div>

            <div class={"desc " + (r.description ? "" : "dim")}>
              {r.description ?? "No description."}
            </div>
          </div>

          <div class="cardBot">
            <div class="meta">
              <span class="metaItem">★ {r.stargazers_count}</span>
              <span class="metaSep">·</span>
              <span class="metaItem">forks {r.forks_count}</span>
            </div>
            <div class="meta dim">
              updated {ago(r.pushed_at ?? r.updated_at)}
            </div>
          </div>

          <div class="shine" aria-hidden="true"></div>
        </a>
      {/each}
    </div>
  {/if}
</div>

<style>
  .wrap { display: grid; gap: 12px; }

  /* Header */
  .hero {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 14px;
    align-items: end;
  }

  .k {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.55);
  }

  h1 {
    margin: 6px 0 6px 0;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    color: rgba(228,90,90,0.95);
    letter-spacing: 0.5px;
  }

  .sub {
    margin: 0;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.70);
    font-size: 13px;
    max-width: 70ch;
  }

  .stats {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .stat {
    border: 1px solid rgba(204,102,102,0.20);
    background: rgba(12,14,18,0.14);
    border-radius: 14px;
    padding: 10px 12px;
    min-width: 94px;
  }
  .stat.muted { opacity: 0.85; }

  .statLabel {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 11px;
    color: rgba(245,245,245,0.58);
  }
  .statVal {
    margin-top: 2px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 14px;
    color: rgba(245,245,245,0.92);
  }

  /* Controls */
  .controlBar {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 10px;
    align-items: center;
  }

  .searchWrap {
    position: relative;
    display: flex;
    align-items: center;
    gap: 8px;

    border: 1px solid rgba(204,102,102,0.20);
    background: rgba(12,14,18,0.14);
    border-radius: 14px;
    padding: 10px 10px;
    overflow: hidden;
  }

  .searchIcon {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    color: rgba(228,90,90,0.75);
    user-select: none;
  }

  .search {
    flex: 1;
    border: 0;
    outline: none;
    background: transparent;
    color: rgba(245,245,245,0.92);
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 13px;
    min-width: 0;
  }

  .miniBtn {
    border: 1px solid rgba(228,90,90,0.30);
    background: rgba(228,90,90,0.10);
    color: rgba(245,245,245,0.86);
    border-radius: 12px;
    padding: 8px 10px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    cursor: pointer;
  }
  .miniBtn:hover { background: rgba(204,102,102,0.12); }

  .controls {
    display: flex;
    gap: 8px;
    align-items: center;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .selectWrap {
    display: flex;
    gap: 8px;
    align-items: center;
    border: 1px solid rgba(204,102,102,0.20);
    background: rgba(12,14,18,0.14);
    border-radius: 14px;
    padding: 8px 10px;
  }

  .lbl {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.60);
  }

  .select {
    border: 0;
    outline: none;
    background: transparent;
    color: rgba(245,245,245,0.90);
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    cursor: pointer;
  }

  .chip {
    border: 1px solid rgba(204,102,102,0.20);
    background: rgba(12,14,18,0.14);
    color: rgba(245,245,245,0.76);
    border-radius: 999px;
    padding: 9px 12px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    cursor: pointer;
    transition: transform 0.15s ease, background 0.15s ease, border-color 0.15s ease;
  }
  .chip:hover { transform: translateY(-1px); background: rgba(204,102,102,0.10); }
  .chip.on {
    border-color: rgba(228,90,90,0.35);
    color: rgba(228,90,90,0.92);
    background: rgba(228,90,90,0.08);
  }
  .chip.ghost {
    border-color: rgba(245,245,245,0.18);
    color: rgba(245,245,245,0.70);
  }

  /* Languages row */
  .langRow {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .langChip {
    display: inline-flex;
    gap: 8px;
    align-items: center;

    border: 1px solid rgba(204,102,102,0.18);
    background: rgba(12,14,18,0.12);
    border-radius: 999px;
    padding: 7px 10px;

    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.78);
  }

  .dot {
    width: 7px;
    height: 7px;
    border-radius: 999px;
    background: rgba(228,90,90,0.85);
    box-shadow: 0 0 0 3px rgba(228,90,90,0.10);
  }

  .count {
    color: rgba(245,245,245,0.60);
    font-size: 11px;
    border-left: 1px solid rgba(245,245,245,0.14);
    padding-left: 8px;
  }

  /* Messages + skeleton */
  .msg {
    border: 1px solid rgba(204,102,102,0.18);
    background: rgba(12,14,18,0.14);
    border-radius: 14px;
    padding: 14px;
  }
  .msg.err { border-color: rgba(255, 140, 140, 0.25); }

  .msgTitle {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.92);
    margin-bottom: 6px;
  }
  .msgBody {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.70);
    font-size: 12px;
    white-space: pre-wrap;
  }

  .skeletonGrid {
    margin-top: 12px;
    display: grid;
    gap: 12px;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  }

  .skCard {
    height: 110px;
    border-radius: 14px;
    border: 1px solid rgba(204,102,102,0.12);
    background: linear-gradient(
      90deg,
      rgba(12,14,18,0.10),
      rgba(204,102,102,0.10),
      rgba(12,14,18,0.10)
    );
    background-size: 200% 100%;
    animation: shimmer 1.2s linear infinite;
  }

  @keyframes shimmer {
    0% { background-position: 0% 0%; }
    100% { background-position: 200% 0%; }
  }

  /* Grid */
  .grid {
    display: grid;
    gap: 12px;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  }

  /* Cards */
  .card {
    position: relative;
    display: grid;
    grid-template-rows: 1fr auto;
    gap: 12px;

    text-decoration: none;

    border: 1px solid rgba(204,102,102,0.18);
    background: rgba(12,14,18,0.14);
    border-radius: 16px;
    padding: 14px;

    overflow: hidden;
    transition: transform 0.15s ease, border-color 0.15s ease, background 0.15s ease;
  }

  .card:hover {
    transform: translateY(-2px);
    border-color: rgba(228,90,90,0.35);
    background: rgba(204,102,102,0.10);
  }

  .shine {
    position: absolute;
    inset: -2px;
    background: radial-gradient(
      600px 220px at var(--mx, 30%) var(--my, 0%),
      rgba(228,90,90,0.10),
      transparent 55%
    );
    opacity: 0;
    transition: opacity 0.2s ease;
    pointer-events: none;
  }

  .card:hover .shine { opacity: 1; }

  .cardTop { display: grid; gap: 10px; }

  .nameRow {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 10px;
  }

  .name {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.94);
    font-size: 13px;
    word-break: break-word;
  }

  .badges { display: flex; gap: 6px; flex-wrap: wrap; justify-content: flex-end; }

  .badge {
    border: 1px solid rgba(204,102,102,0.22);
    background: rgba(12,14,18,0.10);
    border-radius: 999px;
    padding: 3px 8px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 11px;
    color: rgba(245,245,245,0.74);
    white-space: nowrap;
  }
  .badge.dim { opacity: 0.75; }
  .badge.warn {
    border-color: rgba(228,90,90,0.28);
    color: rgba(228,90,90,0.92);
  }

  .desc {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.78);
    font-size: 12px;
    line-height: 1.35;
  }

  .cardBot {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
  }

  .meta {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 11px;
    color: rgba(245,245,245,0.64);
    display: inline-flex;
    align-items: center;
    gap: 8px;
    white-space: nowrap;
  }
  .metaSep { opacity: 0.45; }

  .dim { color: rgba(245,245,245,0.52); }

  /* Responsive */
  @media (max-width: 760px) {
    .hero { grid-template-columns: 1fr; }
    .stats { justify-content: flex-start; }
    .controlBar { grid-template-columns: 1fr; }
    .controls { justify-content: flex-start; }
  }
</style>

<svelte:window
  on:mousemove={(e) => {
    // subtle shine tracking
    const t = e.target as HTMLElement | null;
    const card = t?.closest?.(".card") as HTMLElement | null;
    if (!card) return;
    const r = card.getBoundingClientRect();
    const x = ((e.clientX - r.left) / r.width) * 100;
    const y = ((e.clientY - r.top) / r.height) * 100;
    card.style.setProperty("--mx", `${x}%`);
    card.style.setProperty("--my", `${y}%`);
  }}
/>
