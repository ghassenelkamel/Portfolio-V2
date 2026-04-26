<script lang="ts">
  type Summary = {
    headline: string;
    subtitle: string;
    summary: string[];
    specialties: string[];
    currently: string[];
    goals: string[];
    email?: string;
    site?: string;
    github?: string;
  };

  let loading = $state(true);
  let err = $state<string | null>(null);
  let s = $state<Summary | null>(null);

  const fallback: Summary = {
    headline: "Ghassen El Kamel",
    subtitle: "IoT & Systems Engineer · Backend · Infrastructure",
    summary: [],
    specialties: [],
    currently: [],
    goals: [],
    email: "Ghassenelkamel@live.fr",
    site: "https://ghassenelkamel.fr",
    github: "https://github.com/ghassenelkamel"
  };

  $effect(() => {
    (async () => {
      try {
        loading = true;
        err = null;
        const res = await fetch("/api/content/about");
        const j = await res.json();
        const d = (j?.data ?? {}) as Partial<Summary>;
        s = {
          ...fallback,
          ...d,
          summary: Array.isArray(d.summary) ? d.summary : fallback.summary,
          specialties: Array.isArray(d.specialties) ? d.specialties : fallback.specialties,
          currently: Array.isArray(d.currently) ? d.currently : fallback.currently,
          goals: Array.isArray(d.goals) ? d.goals : fallback.goals
        };
      } catch (e: any) {
        err = e?.message ?? "failed";
      } finally {
        loading = false;
      }
    })();
  });
</script>

<div class="surface">
  <div class="hdr">
    <div class="k">/about</div>
    <h1>{s?.headline ?? "About"}</h1>
    <p>{s?.subtitle ?? "—"}</p>
  </div>

  {#if loading}
    <div class="muted">loading…</div>
  {:else if err}
    <div class="muted">error: {err}</div>
  {:else if s}
    <div class="grid">
      <div class="card">
        <div class="label">About</div>
        {#each s.summary as line (line)}
          <div class="p">{line}</div>
        {/each}

        <div class="label label2">Specialties</div>
        {#each s.specialties as line (line)}
          <div class="li">{line}</div>
        {/each}
      </div>

      <div class="card">
        <div class="label">Currently</div>
        {#each s.currently as line (line)}
          <div class="li">{line}</div>
        {/each}
      </div>

      <div class="card">
        <div class="label">Goals</div>
        {#each s.goals as line (line)}
          <div class="li">{line}</div>
        {/each}

        <div class="sep"></div>

        <div class="kv">
          <div><span class="k2">email</span> {s.email ?? "-"}</div>
          <div><span class="k2">site</span> {s.site ?? "-"}</div>
          <div><span class="k2">github</span> {s.github ?? "-"}</div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .surface {
    position: relative;
    overflow: hidden;
    border: 1px solid var(--border, rgba(255, 255, 255, 0.10));
    background: linear-gradient(180deg, rgba(12, 18, 24, 0.36), rgba(12, 18, 24, 0.24));
    backdrop-filter: blur(var(--blur, 18px));
    -webkit-backdrop-filter: blur(var(--blur, 18px));
    border-radius: var(--radius, 14px);
    padding: 16px;
  }

  .surface::before {
    content: "";
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: repeating-linear-gradient(
      to bottom,
      rgba(255, 255, 255, 0.035) 0,
      rgba(255, 255, 255, 0.035) 1px,
      transparent 1px,
      transparent 7px
    );
    opacity: 0.06;
    mix-blend-mode: overlay;
  }

  .hdr {
    display: grid;
    gap: 8px;
    margin-bottom: 12px;
  }

  .k {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    color: rgba(255, 255, 255, 0.55);
  }

  h1 {
    margin: 0;
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 22px;
    font-weight: 650;
    color: rgba(255, 255, 255, 0.92);
    letter-spacing: -0.25px;
  }

  p {
    margin: 0;
    font-family: var(--font-sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Ubuntu, Cantarell, "Noto Sans", Arial);
    color: rgba(255, 255, 255, 0.72);
    font-size: 14px;
    line-height: 1.55;
  }

  .grid {
    display: grid;
    gap: 14px;
  }

  .card {
    position: relative;
    overflow: hidden;
    border: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(8, 14, 20, 0.30);
    border-radius: 12px;
    padding: 15px;
  }

  .card::before {
    content: "";
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 2px;
    background: linear-gradient(to bottom, rgba(228, 90, 90, 0.9), rgba(228, 90, 90, 0));
  }

  .label {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    letter-spacing: 0.6px;
    text-transform: uppercase;
    color: rgba(228, 90, 90, 0.85);
    margin-bottom: 10px;
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }

  .label2 {
    margin-top: 10px;
  }

  .label::before {
    content: "#";
    color: rgba(255, 255, 255, 0.45);
  }

  .p,
  .li {
    font-family: var(--font-sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Ubuntu, Cantarell, "Noto Sans", Arial);
    font-size: 13px;
    line-height: 1.7;
    color: rgba(255, 255, 255, 0.86);
    margin-bottom: 8px;
  }

  .li {
    padding-left: 12px;
    border-left: 2px solid rgba(228, 90, 90, 0.22);
    position: relative;
  }

  .li::before {
    content: "";
    position: absolute;
    left: -2px;
    top: 0.8em;
    width: 6px;
    height: 1px;
    background: rgba(228, 90, 90, 0.75);
  }

  .muted {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    color: rgba(255, 255, 255, 0.62);
  }

  .sep {
    height: 1px;
    background: rgba(255, 255, 255, 0.08);
    margin: 12px 0;
  }

  .kv {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    color: rgba(255, 255, 255, 0.78);
    display: grid;
    gap: 6px;
  }

  .k2 {
    color: rgba(255, 255, 255, 0.55);
    margin-right: 8px;
  }
</style>
