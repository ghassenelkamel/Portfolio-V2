<script lang="ts">
  type Summary = {
    headline?: string;
    subtitle?: string;
    about?: string;
    working?: string;
    goals?: string;
    email?: string;
    portfolio?: string;
  };

  // SvelteKit provides `data` to the page component
  const { data } = $props<{ data: any }>();

  // Normalize shape: either { ok, summary } or just summary-like
  const summary: Summary =
    data && typeof data === "object" && "summary" in data ? (data.summary ?? {}) : (data ?? {});

  function stripMarkdown(s: string) {
    return (s ?? "")
      .replaceAll(/\*\*(.*?)\*\*/g, "$1")
      .replaceAll(/`([^`]+)`/g, "$1")
      .replaceAll(/\[(.*?)\]\((.*?)\)/g, "$1 ($2)")
      .replaceAll(/\r\n/g, "\n")
      .trim();
  }

  function linesToBullets(s: string) {
    const t = stripMarkdown(s);
    const out: string[] = [];
    for (const raw of t.split("\n")) {
      const line = raw.trim();
      if (!line) continue;
      if (!line.startsWith("-")) continue;

      let b = line.replace(/^-+\s*/, "");
      // remove leading emoji/icons if present
      b = b.replace(/^[\p{Extended_Pictographic}\uFE0F]+\s*/u, "");
      if (b) out.push(b.trim());
    }
    return out;
  }

  function pickParagraphs(s: string) {
    const t = stripMarkdown(s);
    const paras: string[] = [];
    let buf: string[] = [];

    const flush = () => {
      const p = buf.join(" ").replace(/\s+/g, " ").trim();
      if (p) paras.push(p);
      buf = [];
    };

    for (const raw of t.split("\n")) {
      const line = raw.trim();
      if (!line) {
        flush();
        continue;
      }
      if (line.startsWith("-")) continue;
      buf.push(line);
    }
    flush();
    return paras.slice(0, 4);
  }

  const headline = summary.headline || "Ghassen El Kamel";
  const subtitle = summary.subtitle || "IoT & Systems Engineer · Backend · Infrastructure";
  const email = summary.email || "Ghassenelkamel@live.fr";
  const portfolio = summary.portfolio || "https://ghassenelkamel.fr";

  const aboutParas = pickParagraphs(summary.about || "");
  const specialties = linesToBullets(summary.about || "");
  const now = linesToBullets(summary.working || "");
  const goals = linesToBullets(summary.goals || "");

  let copied = $state<string | null>(null);

  async function copy(text: string, label: string) {
    try {
      await navigator.clipboard.writeText(text);
      copied = `${label} copied`;
      setTimeout(() => (copied = null), 1200);
    } catch {
      copied = "Clipboard not available";
      setTimeout(() => (copied = null), 1400);
    }
  }
</script>

<svelte:head>
  <title>About</title>
</svelte:head>

<div class="about">
  <section class="hero">
    <img class="avatar" src="/assets/avatar.png" alt="Ghassen El Kamel" />

    <div class="heroText">
      <div class="kicker">about</div>
      <h1>{headline}</h1>
      <div class="sub">{subtitle}</div>

      <div class="intro">
        <p>
          I build the systems that make connected devices behave like products: reliable data flow, secure connectivity,
          and backends that stay fast under real-world constraints.
        </p>

        {#each aboutParas as p}
          <p class="muted">{p}</p>
        {/each}
      </div>

      <div class="links">
        <div class="linkRow">
          <span class="label">Email</span>
          <a class="link" href={"mailto:" + email}>{email}</a>
          <button class="mini" type="button" onclick={() => copy(email, "Email")}>copy</button>
        </div>

        <div class="linkRow">
          <span class="label">Portfolio</span>
          <a class="link" href={portfolio} target="_blank" rel="noreferrer">{portfolio}</a>
          <button class="mini" type="button" onclick={() => copy(portfolio, "Link")}>copy</button>
        </div>

        <div class="linkRow">
          <span class="label">GitHub</span>
          <a class="link" href="https://github.com/ghassenelkamel" target="_blank" rel="noreferrer">
            github.com/ghassenelkamel
          </a>
        </div>

        {#if copied}
          <div class="toast">{copied}</div>
        {/if}
      </div>
    </div>
  </section>

  <section class="grid">
    <div class="card">
      <div class="cardTitle">Focus</div>
      <div class="chips">
        <span class="chip">Go services</span>
        <span class="chip">Embedded Linux</span>
        <span class="chip">MQTT / TLS</span>
        <span class="chip">VPN networking</span>
        <span class="chip">APIs & automation</span>
        <span class="chip">Reliability</span>
      </div>
      <div class="hr"></div>
      <div class="listTitle">Core strengths</div>
      {#if specialties.length}
        <ul class="list">
          {#each specialties as it}
            <li>{it}</li>
          {/each}
        </ul>
      {:else}
        <div class="empty">Add bullet points to your profile summary to populate this section.</div>
      {/if}
    </div>

    <div class="card">
      <div class="cardTitle">Now</div>
      {#if now.length}
        <ul class="list">
          {#each now as it}
            <li>{it}</li>
          {/each}
        </ul>
      {:else}
        <div class="empty">No “Now” items yet.</div>
      {/if}
    </div>

    <div class="card">
      <div class="cardTitle">Goals</div>
      {#if goals.length}
        <ul class="list">
          {#each goals as it}
            <li>{it}</li>
          {/each}
        </ul>
      {:else}
        <div class="empty">No “Goals” items yet.</div>
      {/if}
    </div>
  </section>
</div>

<style>
  .about { padding: 18px; display: grid; gap: 16px; }

  .hero {
    display: grid;
    grid-template-columns: 84px 1fr;
    gap: 14px;
    align-items: start;

    padding: 14px;
    border-radius: 14px;
    border: 1px solid rgba(245, 245, 245, 0.10);
    background: rgba(12, 14, 18, 0.34);
    backdrop-filter: blur(22px) saturate(150%);
    -webkit-backdrop-filter: blur(22px) saturate(150%);
  }

  .avatar {
    width: 84px;
    height: 84px;
    border-radius: 14px;
    object-fit: cover;
    border: 1px solid rgba(204, 102, 102, 0.30);
    box-shadow: 0 10px 30px rgba(0,0,0,0.35);
  }

  .kicker {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    letter-spacing: 1.2px;
    color: rgba(228, 90, 90, 0.75);
    text-transform: uppercase;
    margin-bottom: 6px;
  }

  h1 {
    margin: 0;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 22px;
    letter-spacing: 0.6px;
    color: rgba(245, 245, 245, 0.94);
  }

  .sub {
    margin-top: 6px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 13px;
    color: rgba(245, 245, 245, 0.70);
  }

  .intro { margin-top: 12px; display: grid; gap: 10px; }
  .intro p { margin: 0; line-height: 1.55; color: rgba(245, 245, 245, 0.88); font-size: 14px; }
  .intro .muted { color: rgba(245, 245, 245, 0.72); font-size: 13px; }

  .links { margin-top: 12px; display: grid; gap: 8px; }
  .linkRow { display: grid; grid-template-columns: 70px 1fr auto; gap: 10px; align-items: center; }

  .label {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.55);
  }

  .link {
    text-decoration: none;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.86);
    border-bottom: 1px dashed rgba(228, 90, 90, 0.35);
    padding-bottom: 2px;
    overflow-wrap: anywhere;
  }
  .link:hover { border-bottom-color: rgba(228, 90, 90, 0.75); }

  .mini {
    border: 1px solid rgba(204, 102, 102, 0.25);
    background: rgba(12, 14, 18, 0.20);
    color: rgba(245, 245, 245, 0.82);
    border-radius: 999px;
    padding: 6px 10px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    cursor: pointer;
  }
  .mini:hover { background: rgba(204, 102, 102, 0.14); }

  .toast {
    margin-top: 6px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.70);
  }

  .grid { display: grid; grid-template-columns: 1.2fr 1fr 1fr; gap: 14px; }

  .card {
    padding: 14px;
    border-radius: 14px;
    border: 1px solid rgba(245, 245, 245, 0.10);
    background: rgba(12, 14, 18, 0.28);
    backdrop-filter: blur(22px) saturate(150%);
    -webkit-backdrop-filter: blur(22px) saturate(150%);
    overflow: hidden;
  }

  .cardTitle {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    letter-spacing: 1.1px;
    color: rgba(228, 90, 90, 0.75);
    text-transform: uppercase;
    margin-bottom: 10px;
  }

  .chips { display: flex; flex-wrap: wrap; gap: 8px; }
  .chip {
    padding: 7px 10px;
    border-radius: 999px;
    border: 1px solid rgba(204, 102, 102, 0.22);
    background: rgba(204, 102, 102, 0.10);
    color: rgba(245, 245, 245, 0.86);
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    white-space: nowrap;
  }

  .hr { height: 1px; background: rgba(204, 102, 102, 0.18); margin: 12px 0; }

  .listTitle {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.72);
    margin-bottom: 8px;
  }

  .list {
    margin: 0;
    padding-left: 18px;
    display: grid;
    gap: 8px;
    color: rgba(245, 245, 245, 0.80);
    font-size: 13px;
    line-height: 1.45;
  }

  .empty { color: rgba(245, 245, 245, 0.62); font-size: 13px; line-height: 1.45; }

  @media (max-width: 950px) {
    .grid { grid-template-columns: 1fr; }
    .hero { grid-template-columns: 72px 1fr; }
    .avatar { width: 72px; height: 72px; }
  }
</style>
