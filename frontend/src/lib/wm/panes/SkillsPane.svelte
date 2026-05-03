<script lang="ts">
  import { page } from "$app/stores";
  import { onDestroy, onMount } from "svelte";

  type SkillContent = {
    id: string;
    name: string;
    level: number;
    order: number;
  };

  type SkillsContent = {
    eyebrow?: string;
    title?: string;
    description?: string;
    items?: SkillContent[];
  };

  type LiveSkill = {
    name: string;
    base: number;
    live: number;
    history: number[];
  };

  const points = 22;

  const fallbackItems: SkillContent[] = [
    { id: "go", name: "Go (Backend)", level: 88, order: 1 },
    { id: "mqtt-vpn-tls", name: "MQTT / VPN / TLS", level: 86, order: 2 },
    { id: "linux", name: "Linux / Unix", level: 86, order: 3 },
    { id: "postgresql", name: "PostgreSQL", level: 82, order: 4 },
    { id: "docker", name: "Docker", level: 78, order: 5 },
    { id: "javascript-svelte", name: "JavaScript / Svelte", level: 78, order: 6 },
    { id: "c-cpp", name: "C / C++", level: 72, order: 7 },
    { id: "python", name: "Python", level: 70, order: 8 }
  ];

  let content = $state<SkillsContent>({
    eyebrow: "Skills",
    title: "Capability Monitor",
    description: "Track live proficiency signals across core tools and technologies.",
    items: fallbackItems
  });

  let skills = $state<LiveSkill[]>(
    fallbackItems.map((s) => {
      const v = Math.round(s.level);
      return { name: s.name, base: v, live: v, history: Array.from({ length: points }, () => v) };
    })
  );

  let tick = 0;
  let timer: ReturnType<typeof setInterval> | null = null;
  const langQuery = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr") ? "?lang=fr" : "");
  const isFr = $derived(langQuery === "?lang=fr");

  const ui = $derived.by(() => isFr
    ? {
        eyebrow: "Competences",
        title: "Moniteur de capacites",
        desc: "Suivi des niveaux de maitrise sur les outils et technologies cles.",
        avg: "moy",
        top: "top",
        boardAria: "Moniteur de compétences en direct"
      }
    : {
        eyebrow: "Skills",
        title: "Capability Monitor",
        desc: "Track live proficiency signals across core tools and technologies.",
        avg: "avg",
        top: "top",
        boardAria: "Live skill monitor"
      });

  function clamp(v: number, a: number, b: number) {
    return Math.max(a, Math.min(b, v));
  }

  function step() {
    tick += 1;
    skills = skills.map((s, i) => {
      const wave = Math.sin((tick + i * 2.4) / 4.2) * 4.4;
      const jitter = (Math.random() - 0.5) * 2.4;
      const target = clamp(s.base + wave + jitter, 35, 99);
      const live = clamp(s.live + (target - s.live) * 0.4, 35, 99);
      return {
        ...s,
        live,
        history: [...s.history.slice(1), Math.round(live)]
      };
    });
  }

  $effect(() => {
    void (async () => {
      try {
        const res = await fetch(`/api/content/skills${langQuery}`);
        const j = await res.json();
        const d = (j?.data ?? {}) as SkillsContent;
        const items = Array.isArray(d.items) && d.items.length
          ? [...d.items].sort((a, b) => a.order - b.order)
          : fallbackItems;

        content = {
          eyebrow: d.eyebrow || ui.eyebrow,
          title: d.title || ui.title,
          description: d.description || ui.desc,
          items
        };

        skills = items.map((s) => {
          const v = Math.round(s.level);
          return { name: s.name, base: v, live: v, history: Array.from({ length: points }, () => v) };
        });
      } catch {
        // Keep fallback content.
      }
    })();
  });

  onMount(() => {
    timer = setInterval(step, 280);
  });

  onDestroy(() => {
    if (timer) clearInterval(timer);
  });

  const avg = $derived(Math.round(skills.reduce((sum, s) => sum + s.live, 0) / Math.max(skills.length, 1)));
  const top = $derived(skills.length ? [...skills].sort((a, b) => b.live - a.live)[0] : null);

  function barClass(v: number) {
    if (v >= 85) return "hot";
    if (v >= 75) return "warm";
    return "cool";
  }
</script>

<div class="pane">
  <header class="head">
    <div class="eyebrow">{content.eyebrow || ui.eyebrow}</div>
    <div class="t">{content.title || ui.title}</div>
    <div class="s">{content.description || ui.desc}</div>
    <div class="meta">
      <span>{ui.avg} {avg}%</span>
      <span>{ui.top} {top?.name ?? "-"}</span>
    </div>
  </header>

  <section class="board" aria-label={ui.boardAria}>
    {#each skills as s (s.name)}
      <article class="row">
        <div class="rhead">
          <span class="name">{s.name}</span>
          <span class="val">{Math.round(s.live)}%</span>
        </div>

        <div class="meter">
          <span class={`fill ${barClass(s.live)}`} style={`width:${s.live}%;`}></span>
          <span class="ticks" aria-hidden="true"></span>
        </div>

      </article>
    {/each}
  </section>
</div>

<style>
  .pane {
    display: grid;
    gap: 12px;
    color: rgba(255, 255, 255, 0.9);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
  }

  .head {
    display: grid;
    gap: 7px;
  }

  .eyebrow {
    font-size: 12px;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: rgba(245, 245, 245, 0.58);
  }

  .t {
    font-size: 20px;
    line-height: 1.2;
    letter-spacing: 0.2px;
    color: var(--accent, #e45a5a);
  }

  .s {
    margin-top: 1px;
    max-width: 52ch;
    font-size: 12px;
    line-height: 1.45;
    color: rgba(245, 245, 245, 0.7);
  }

  .meta {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.65);
  }

  .board {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    background: rgba(8, 14, 20, 0.3);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    padding: 12px;
    display: grid;
    gap: 11px;
  }

  .row {
    display: grid;
    gap: 6px;
  }

  .rhead {
    display: flex;
    justify-content: space-between;
    gap: 8px;
    align-items: baseline;
  }

  .name {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.88);
  }

  .val {
    font-size: 12px;
    font-weight: 700;
    color: rgba(255, 170, 125, 0.92);
  }

  .meter {
    position: relative;
    height: 14px;
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.11);
    background: rgba(255, 255, 255, 0.05);
    overflow: hidden;
  }

  .ticks {
    position: absolute;
    inset: 0;
    background: repeating-linear-gradient(
      90deg,
      transparent 0,
      transparent 18px,
      rgba(255, 255, 255, 0.08) 18px,
      rgba(255, 255, 255, 0.08) 19px
    );
    opacity: 0.28;
    pointer-events: none;
  }

  .fill {
    position: absolute;
    inset: 0 auto 0 0;
    border-radius: 6px;
    transition: width 260ms linear;
    box-shadow: 0 0 8px rgba(228, 90, 90, 0.1);
  }

  .fill.hot {
    background: linear-gradient(90deg, rgba(228, 90, 90, 0.82), rgba(228, 124, 124, 0.76));
  }

  .fill.warm {
    background: linear-gradient(90deg, rgba(228, 104, 104, 0.72), rgba(228, 132, 132, 0.66));
  }

  .fill.cool {
    background: linear-gradient(90deg, rgba(228, 116, 116, 0.62), rgba(228, 148, 148, 0.56));
  }
</style>
