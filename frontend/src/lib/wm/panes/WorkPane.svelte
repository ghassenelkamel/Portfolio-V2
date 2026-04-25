<script lang="ts">
  import { onMount } from "svelte";

  type Repo = {
    id: number;
    name: string;
    html_url: string;
    description: string | null;
    language: string | null;
    stargazers_count?: number;
    forks_count?: number;
  };

  let loading = $state(true);
  let err = $state<string | null>(null);
  let repos = $state<Repo[]>([]);

  onMount(async () => {
    try {
      const r = await fetch("/api/github-repos?user=ghassenelkamel&limit=100");
      const j = await r.json();
      if (j?.error) err = j.error;
      repos = j?.repos ?? [];
    } catch (e) {
      err = e instanceof Error ? e.message : "Unknown error.";
    } finally {
      loading = false;
    }
  });
</script>

<div class="pane">
  <h2>Work</h2>

  {#if loading}
    <div class="muted">loading…</div>
  {:else if err}
    <div class="err">{err}</div>
  {/if}

  <div class="list">
    {#each repos as r (r.id)}
      <a class="row" href={r.html_url} target="_blank" rel="noreferrer">
        <div class="name">{r.name}</div>
        <div class="meta">
          <span class="lang">{r.language ?? "—"}</span>
          <span class="desc">{r.description ?? ""}</span>
        </div>
      </a>
    {/each}
  </div>
</div>

<style>
  .pane {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    color: rgba(255, 255, 255, 0.88);
  }

  h2 {
    margin: 0 0 10px;
    color: var(--accent, #e45a5a);
    font-size: 18px;
    letter-spacing: -0.25px;
  }

  .muted {
    color: rgba(255, 255, 255, 0.55);
  }

  .err {
    color: rgba(255, 140, 140, 0.92);
  }

  .list {
    display: grid;
    gap: 10px;
    margin-top: 12px;
  }

  .row {
    display: grid;
    gap: 8px;
    padding: 12px;
    border: 1px solid rgba(255, 255, 255, 0.10);
    border-radius: 12px;
    background: rgba(15, 20, 25, 0.28);
    text-decoration: none;
    color: inherit;
    transition: transform 0.15s ease, border-color 0.15s ease, background 0.15s ease;
  }

  .row:hover {
    transform: translateY(-1px);
    border-color: rgba(228, 90, 90, 0.28);
    background: rgba(255, 255, 255, 0.04);
  }

  .name {
    color: rgba(255, 255, 255, 0.92);
    font-weight: 650;
  }

  .meta {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.65);
  }

  .lang {
    color: rgba(0, 255, 180, 0.80);
  }

  .desc {
    opacity: 0.95;
  }
</style>
