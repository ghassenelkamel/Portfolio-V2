<script lang="ts">
  import type { WMNode as Node, SplitMode } from "$lib/wm/store";
  import { wmSetSplitRatio } from "$lib/wm/store";
  import WMNode from "./WMNode.svelte";
  import WindowFrame from "./WindowFrame.svelte";

  const props = $props<{ node: Node }>();

  let box = $state<HTMLDivElement | null>(null);
  let w = $state(0);
  let h = $state(0);

  let dragging = $state(false);

  function modeToDir(mode: SplitMode, W: number, H: number): "row" | "col" {
    if (mode === "row") return "row";
    if (mode === "col") return "col";
    return W > H ? "row" : "col";
  }

  function startDrag(e: PointerEvent, splitId: string, dir: "row" | "col") {
    e.preventDefault();
    e.stopPropagation();
    dragging = true;
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);

    const onMove = (ev: PointerEvent) => {
      if (!box) return;
      const r = box.getBoundingClientRect();
      if (dir === "row") {
        wmSetSplitRatio(splitId, (ev.clientX - r.left) / r.width);
      } else {
        wmSetSplitRatio(splitId, (ev.clientY - r.top) / r.height);
      }
    };

    const onUp = () => {
      dragging = false;
      window.removeEventListener("pointermove", onMove);
      window.removeEventListener("pointerup", onUp);
    };

    window.addEventListener("pointermove", onMove);
    window.addEventListener("pointerup", onUp);
  }

  $effect(() => {
    if (!box) return;
    const ro = new ResizeObserver((entries) => {
      const cr = entries[0]?.contentRect;
      if (!cr) return;
      w = cr.width;
      h = cr.height;
    });
    ro.observe(box);
    return () => ro.disconnect();
  });
</script>

<div class="nodeBox" bind:this={box}>
  {#if props.node.type === "leaf"}
    <WindowFrame winId={props.node.winId} />
  {:else}
    {@const dir = modeToDir(props.node.mode, w, h)}
    {@const aFlex = Math.max(1, Math.round(props.node.ratio * 1000))}
    {@const bFlex = Math.max(1, Math.round((1 - props.node.ratio) * 1000))}

    <div class={"split " + dir} style={"--a:" + aFlex + "; --b:" + bFlex}>
      <div class="pane a">
        <WMNode node={props.node.a} />
      </div>

      <div
        class={"divider " + (dir === "row" ? "v" : "h") + (dragging ? " dragging" : "")}
        onpointerdown={(e) => startDrag(e, props.node.id, dir)}
        title="drag to resize"
      ></div>

      <div class="pane b">
        <WMNode node={props.node.b} />
      </div>
    </div>
  {/if}
</div>

<style>
  .nodeBox { width: 100%; height: 100%; min-width: 0; min-height: 0; }

  .split {
    width: 100%;
    height: 100%;
    display: flex;
    gap: 10px;
    min-width: 0;
    min-height: 0;
  }
  .split.row { flex-direction: row; }
  .split.col { flex-direction: column; }

  .pane { min-width: 0; min-height: 0; }
  .pane.a { flex: var(--a) 1 0; }
  .pane.b { flex: var(--b) 1 0; }

  .divider {
    flex: 0 0 auto;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.10);
    opacity: 0.85;
    transition: background 0.12s ease, opacity 0.12s ease, border-color 0.12s ease;
  }

  .divider.v { width: 8px; cursor: col-resize; }
  .divider.h { height: 8px; cursor: row-resize; }

  .divider:hover {
    background: rgba(228, 90, 90, 0.10);
    border-color: rgba(228, 90, 90, 0.18);
    opacity: 1;
  }

  .divider.dragging {
    background: rgba(0, 255, 180, 0.08);
    border-color: rgba(0, 255, 180, 0.20);
  }

  @media (max-width: 700px) {
    .split {
      gap: 12px;
    }

    .divider.v { width: 10px; }
    .divider.h { height: 10px; }
  }
</style>
