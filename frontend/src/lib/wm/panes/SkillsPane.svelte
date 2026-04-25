<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { avgSkill, skillsData, topSkill } from "$lib/data/skills";

  type LiveSkill = {
    name: string;
    base: number;
    live: number;
    history: number[];
  };

  const points = 22;
  let skills = $state<LiveSkill[]>(
    skillsData.map((s) => {
      const v = Math.round(s.pct * 100);
      return {
        name: s.name,
        base: v,
        live: v,
        history: Array.from({ length: points }, () => v)
      };
    })
  );

  let tick = 0;
  let timer: ReturnType<typeof setInterval> | null = null;

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

  onMount(() => {
    timer = setInterval(step, 280);
  });

  onDestroy(() => {
    if (timer) clearInterval(timer);
  });

  const avg = $derived(Math.round(avgSkill(skills.map((s) => ({ name: s.name, pct: s.live / 100 }))) * 100));
  const top = $derived(topSkill(skills.map((s) => ({ name: s.name, pct: s.live / 100 }))));

  function barClass(v: number) {
    if (v >= 85) return "hot";
    if (v >= 75) return "warm";
    return "cool";
  }
</script>

<div class="pane">
  <header class="head">
    <h2>Skills</h2>
    <div class="meta">
      <span>avg {avg}%</span>
      <span>top {top?.name ?? "-"}</span>
    </div>
  </header>

  <section class="board" aria-label="Live skill monitor">
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
    gap: 6px;
  }

  h2 {
    margin: 0;
    font-size: 20px;
    color: var(--accent, #e45a5a);
    letter-spacing: -0.2px;
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
