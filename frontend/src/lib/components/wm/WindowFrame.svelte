<script lang="ts">
  import { goto } from "$app/navigation";
  import { wm, wmClose, wmDuplicate, wmFocus, wmSplitOpen, type AppKind, type SplitMode } from "$lib/wm/store";

  import Terminal from "$lib/components/Terminal.svelte";
  import AboutPane from "$lib/wm/panes/AboutPane.svelte";
  import WorkPane from "$lib/wm/panes/WorkPane.svelte";
  import SkillsPane from "$lib/wm/panes/SkillsPane.svelte";
  import ContactPane from "$lib/wm/panes/ContactPane.svelte";
  import ExperiencePane from "$lib/wm/panes/ExperiencePane.svelte";

  const props = $props<{ winId: string }>();
  const wmState = $derived($wm);
  const win = $derived(wmState.windows[props.winId]);
  const focused = $derived(wmState.focusWinId === props.winId);

  function routeFor(kind: AppKind) {
    switch (kind) {
      case "terminal": return "/";
      case "about": return "/about";
      case "experience": return "/experience";
      case "work": return "/work";
      case "skills": return "/skills";
      case "contact": return "/contact";
    }
  }

  const launcher: { kind: AppKind; label: string }[] = [
    { kind: "terminal", label: "Terminal" },
    { kind: "about", label: "About" },
    { kind: "experience", label: "Experience" },
    { kind: "work", label: "Work" },
    { kind: "skills", label: "Skills" },
    { kind: "contact", label: "Contact" }
  ];

  let open = $state(false);
  let mode = $state<SplitMode>("auto");

  function openLauncher(m: SplitMode) { mode = m; open = true; }
  function launch(kind: AppKind) { wmSplitOpen(kind, mode, props.winId); open = false; }
  function fullscreen() { if (win) void goto(routeFor(win.kind)); }
</script>

<div class={"win " + (focused ? "focus" : "")} onclick={() => wmFocus(props.winId)}>
  <div class="titlebar">
    <div class="left">
      <span class="tag">{win?.title ?? "?"}</span>
      {#if focused}<span class="dot">●</span>{/if}
    </div>

    <div class="right">
      <button type="button" class="btn" title="open (auto dwindle split)" onclick={(e) => { e.stopPropagation(); openLauncher("auto"); }}>+</button>
      <button type="button" class="btn" title="force vertical split (left/right)" onclick={(e) => { e.stopPropagation(); openLauncher("row"); }}>v</button>
      <button type="button" class="btn" title="force horizontal split (top/bottom)" onclick={(e) => { e.stopPropagation(); openLauncher("col"); }}>h</button>

      <button type="button" class="btn" title="duplicate (auto)" onclick={(e) => { e.stopPropagation(); wmDuplicate(props.winId, "auto"); }}>dup</button>
      <button type="button" class="btn" title="fullscreen route" onclick={(e) => { e.stopPropagation(); fullscreen(); }}>fs</button>
      <button type="button" class="btn danger" title="close" onclick={(e) => { e.stopPropagation(); wmClose(props.winId); }}>x</button>
    </div>
  </div>

  {#if open}
    <div class="launcher" onclick={(e) => e.stopPropagation()}>
      <div class="launchTitle">
        {#if mode === "auto"}auto split{/if}
        {#if mode === "row"}force vertical{/if}
        {#if mode === "col"}force horizontal{/if}
      </div>

      {#each launcher as item (item.kind)}
        <button type="button" class="launchItem" onclick={() => launch(item.kind)}>{item.label}</button>
      {/each}
    </div>
  {/if}

  <div class="body">
    {#if win?.kind === "terminal"}
      <Terminal visible={true} wmMode={true} />
    {:else if win?.kind === "about"}
      <AboutPane />
    {:else if win?.kind === "experience"}
      <ExperiencePane />
    {:else if win?.kind === "work"}
      <WorkPane />
    {:else if win?.kind === "skills"}
      <SkillsPane />
    {:else if win?.kind === "contact"}
      <ContactPane />
    {:else}
      <div class="pad muted">unknown</div>
    {/if}
  </div>
</div>

<style>
  .win {
    position: relative; /* REQUIRED so launcher stacks correctly */
    width: 100%;
    height: 100%;
    display: grid;
    grid-template-rows: 36px 1fr;
    min-width: 0;
    min-height: 0;

    background: var(--surface, rgba(15, 20, 25, 0.35));
    border: 1px solid var(--border, rgba(255, 255, 255, 0.10));
    border-radius: 12px;
    box-shadow:
      0 10px 34px rgba(0, 0, 0, 0.35),
      0 0 0 1px rgba(255, 255, 255, 0.03) inset;

    backdrop-filter: blur(var(--blur, 18px));
    -webkit-backdrop-filter: blur(var(--blur, 18px));
    overflow: hidden;
  }

  .win.focus {
    border-color: rgba(228, 90, 90, 0.30);
    box-shadow:
      0 10px 34px rgba(0, 0, 0, 0.35),
      0 0 0 1px rgba(228, 90, 90, 0.12) inset;
  }

  .titlebar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    padding: 0 12px;

    background: rgba(10, 14, 18, 0.55);
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);

    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12.5px;
    color: rgba(255, 255, 255, 0.78);
    user-select: none;
  }

  .left {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
  }

  .tag {
    color: rgba(255, 255, 255, 0.92);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .dot {
    color: var(--accent, #e45a5a);
  }

  .right {
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }

  .btn {
    border: 1px solid rgba(255, 255, 255, 0.10);
    background: rgba(255, 255, 255, 0.04);
    color: rgba(255, 255, 255, 0.86);
    font-family: inherit;
    font-size: 11px;
    padding: 5px 8px;
    cursor: pointer;
    border-radius: 8px;
    transition: border-color 0.12s ease, background 0.12s ease, transform 0.12s ease;
  }

  .btn:hover {
    border-color: rgba(228, 90, 90, 0.28);
    background: rgba(228, 90, 90, 0.10);
    transform: translateY(-1px);
  }

  .danger:hover {
    border-color: rgba(255, 140, 140, 0.55);
    background: rgba(255, 140, 140, 0.10);
  }

  /* IMPORTANT: allow scrolling when content is taller */
  .body {
    overflow: auto;
    min-width: 0;
    min-height: 0;
    padding: 12px;
  }

  .launcher {
    position: absolute;
    top: 40px;
    right: 10px;
    z-index: 999;
    display: grid;
    gap: 8px;
    padding: 10px;

    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 12px;
    background: rgba(10, 14, 18, 0.78);
    backdrop-filter: blur(var(--blur, 18px));
    -webkit-backdrop-filter: blur(var(--blur, 18px));
    box-shadow: 0 10px 24px rgba(0, 0, 0, 0.25);
  }

  .launchTitle {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 11px;
    color: rgba(228, 90, 90, 0.85);
  }

  .launchItem {
    border: 1px solid rgba(255, 255, 255, 0.10);
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.04);
    color: rgba(255, 255, 255, 0.90);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    padding: 10px 12px;
    cursor: pointer;
    text-align: left;
    transition: background 0.12s ease, border-color 0.12s ease;
  }

  .launchItem:hover {
    border-color: rgba(228, 90, 90, 0.25);
    background: rgba(228, 90, 90, 0.08);
  }

  .pad {
    padding: 12px;
  }

  .muted {
    color: rgba(255, 255, 255, 0.60);
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
  }

  @media (max-width: 700px) {
    .win {
      grid-template-rows: 40px 1fr;
      border-radius: 14px;
    }

    .body {
      padding: 10px;
    }

    .right {
      gap: 4px;
    }

    .btn {
      padding: 6px 8px;
    }
  }
</style>
