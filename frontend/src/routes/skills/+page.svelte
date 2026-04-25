<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { skillsData } from "$lib/data/skills";

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

  let skills = $state<LiveSkill[]>(
    [...skillsData]
      .sort((a, b) => b.pct - a.pct)
      .map((s) => {
        const v = Math.round(s.pct * 100);
        return {
          name: s.name,
          base: v,
          live: v,
          history: Array.from({ length: points }, () => v)
        };
      })
  );

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
  <title>Skills</title>
</svelte:head>

<div class="wrap">
  <header class="title">
    <h1>Skills</h1>
    <div class="meta">
      <span>avg {avg}%</span>
      <span>peak {peak?.name ?? "-"}</span>
      <span>uptime {fmtUptime(uptime)}</span>
    </div>
  </header>

  <section class="board" aria-label="Live skills dashboard">
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
    gap: 6px;
  }

  h1 {
    margin: 0;
    color: rgba(255, 255, 255, 0.95);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: clamp(22px, 2.4vw, 30px);
    letter-spacing: -0.4px;
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
    .name {
      font-size: 12px;
    }
  }
</style>
