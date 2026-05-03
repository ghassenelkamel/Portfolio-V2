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

  const { data } = $props<{ data: { content?: SkillsContent } }>();
  const content = $derived.by(() => data?.content ?? {});
  const isFr = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr"));

  const ui = $derived.by(() => isFr
    ? {
        title: "Competences",
        fallbackEyebrow: "Competences",
        fallbackHeader: "Moniteur de capacites",
        fallbackDesc: "Suivi des niveaux de maitrise sur les outils et technologies cles.",
        avg: "moy",
        peak: "pic",
        uptime: "actif",
        dashboardAria: "Tableau de bord des competences en direct"
      }
    : {
        title: "Skills",
        fallbackEyebrow: "Skills",
        fallbackHeader: "Capability Monitor",
        fallbackDesc: "Track live proficiency signals across core tools and technologies.",
        avg: "avg",
        peak: "peak",
        uptime: "uptime",
        dashboardAria: "Live skills dashboard"
      });
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

  const baseItems = $derived.by<SkillContent[]>(() => {
    const items = content.items;
    if (Array.isArray(items) && items.length) {
      return [...items].sort((a, b) => a.order - b.order);
    }
    return fallbackItems;
  });

  type LiveSkill = {
    name: string;
    base: number;
    live: number;
    history: number[];
  };

  const points = 34;
  let tick = 0;
  let uptime = $state(0);
  let timer: ReturnType<typeof setInterval> | null = null;

  function toLiveSkills(items: SkillContent[]): LiveSkill[] {
    return [...items]
      .sort((a, b) => b.level - a.level)
      .map((s) => {
        const v = Math.round(s.level);
        return {
          name: s.name,
          base: v,
          live: v,
          history: Array.from({ length: points }, () => v)
        };
      });
  }

  let skills = $state<LiveSkill[]>([]);

  $effect(() => {
    skills = toLiveSkills(baseItems);
    tick = 0;
    uptime = 0;
  });

  function clamp(v: number, a: number, b: number) {
    return Math.max(a, Math.min(b, v));
  }

  function updateLive() {
    tick += 1;
    uptime += 1;

    skills = skills.map((s, i) => {
      const wave = Math.sin((tick + i * 2.2) / 5.1) * 4.8;
      const burst = Math.sin((tick + i * 11) / 16) * 2.2;
      const noise = (Math.random() - 0.5) * 2.6;
      const target = clamp(s.base + wave + burst + noise, 35, 99);
      const live = clamp(s.live + (target - s.live) * 0.35, 35, 99);
      return {
        ...s,
        live,
        history: [...s.history.slice(1), Math.round(live)]
      };
    });
  }

  onMount(() => {
    timer = setInterval(updateLive, 260);
  });

  onDestroy(() => {
    if (timer) clearInterval(timer);
  });

  const avg = $derived(Math.round(skills.reduce((a, s) => a + s.live, 0) / skills.length));
  const peak = $derived([...skills].sort((a, b) => b.live - a.live)[0] ?? null);

  function level(v: number) {
    if (v >= 85) return "hot";
    if (v >= 75) return "warm";
    return "cool";
  }

  function fmtUptime(sec: number) {
    const m = Math.floor(sec / 60);
    const s = sec % 60;
    return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
  }
</script>

<svelte:head>
  <title>{ui.title}</title>
</svelte:head>

<div class="wrap">
  <header class="title">
    <div class="eyebrow">{content.eyebrow || ui.fallbackEyebrow}</div>
    <div class="t">{content.title || ui.fallbackHeader}</div>
    <div class="s">{content.description || ui.fallbackDesc}</div>
    <div class="meta">
      <span>{ui.avg} {avg}%</span>
      <span>{ui.peak} {peak?.name ?? "-"}</span>
      <span>{ui.uptime} {fmtUptime(uptime)}</span>
    </div>
  </header>

  <section class="board" aria-label={ui.dashboardAria}>
    {#each skills as s, i (s.name)}
      <article class="row">
        <div class="head">
          <div class="left">
            <span class="idx">{String(i + 1).padStart(2, "0")}</span>
            <span class="name">{s.name}</span>
          </div>
          <span class="val">{Math.round(s.live)}%</span>
        </div>

        <div class="meter">
          <span class={`fill ${level(s.live)}`} style={`width:${s.live}%;`}></span>
          <span class="ticks" aria-hidden="true"></span>
        </div>

      </article>
    {/each}
  </section>
</div>

<style>
  .wrap {
    display: grid;
    gap: 14px;
  }

  .title {
    display: grid;
    gap: 7px;
  }

  .eyebrow {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: rgba(245, 245, 245, 0.58);
  }

  .t {
    color: rgba(228, 90, 90, 0.95);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: clamp(22px, 2.4vw, 30px);
    line-height: 1.2;
    letter-spacing: 0.2px;
  }

  .s {
    margin-top: 1px;
    max-width: 66ch;
    color: rgba(245, 245, 245, 0.72);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 13px;
    line-height: 1.5;
  }

  .meta {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    color: rgba(255, 255, 255, 0.66);
  }

  .board {
    border: 1px solid rgba(255, 255, 255, 0.10);
    border-radius: 14px;
    background: linear-gradient(180deg, rgba(8, 14, 20, 0.42), rgba(8, 14, 20, 0.28));
    backdrop-filter: blur(14px);
    -webkit-backdrop-filter: blur(14px);
    padding: 12px;
    display: grid;
    gap: 11px;
  }

  .row {
    display: grid;
    gap: 7px;
  }

  .head {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 10px;
  }

  .left {
    display: flex;
    gap: 8px;
    align-items: baseline;
  }

  .idx {
    color: rgba(255, 255, 255, 0.44);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 11px;
  }

  .name {
    color: rgba(255, 255, 255, 0.9);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 13px;
  }

  .val {
    color: rgba(255, 170, 125, 0.92);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    font-weight: 700;
  }

  .meter {
    position: relative;
    height: 15px;
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.11);
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
    transition: width 240ms linear;
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

  @media (max-width: 760px) {
    .wrap {
      gap: 12px;
    }

    .t {
      font-size: 22px;
    }

    .s {
      max-width: none;
      font-size: 12px;
    }

    .meta {
      gap: 8px;
      font-size: 11px;
    }

    .board {
      padding: 10px;
      gap: 10px;
    }

    .head {
      gap: 8px;
    }

    .name {
      font-size: 12px;
    }

    .val {
      font-size: 11px;
    }

    .meter {
      height: 14px;
    }
  }

  @media (max-width: 520px) {
    .left {
      flex-wrap: wrap;
      row-gap: 4px;
    }

    .idx {
      font-size: 10px;
    }
  }
</style>
