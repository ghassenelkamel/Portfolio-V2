import { writable, get } from "svelte/store";

export type AppKind =
  | "terminal"
  | "about"
  | "experience"
  | "work"
  | "skills"
  | "contact";

export type SplitMode = "auto" | "row" | "col";

export type WMWindow = {
  id: string;
  kind: AppKind;
  title: string;
  createdAt: number;
};

export type WMNode =
  | { type: "leaf"; winId: string }
  | { type: "split"; id: string; mode: SplitMode; ratio: number; a: WMNode; b: WMNode };

export type WMState = {
  windows: Record<string, WMWindow>;
  order: string[];
  root: WMNode;
  focusWinId: string;
};

let wSeq = 0;
let sSeq = 0;

function newWinId() {
  wSeq += 1;
  return `w_${Date.now().toString(36)}_${wSeq.toString(36)}`;
}
function newSplitId() {
  sSeq += 1;
  return `s_${Date.now().toString(36)}_${sSeq.toString(36)}`;
}

function titleFor(kind: AppKind) {
  switch (kind) {
    case "terminal": return "Terminal";
    case "about": return "About";
    case "experience": return "Experience";
    case "work": return "Work";
    case "skills": return "Skills";
    case "contact": return "Contact";
  }
}

function makeWindow(kind: AppKind): WMWindow {
  const id = newWinId();
  return { id, kind, title: titleFor(kind), createdAt: Date.now() };
}

export const wm = writable<WMState>({
  windows: {},
  order: [],
  root: { type: "leaf", winId: "" },
  focusWinId: ""
});

function findPath(node: WMNode, winId: string, path: ("a" | "b")[] = []): ("a" | "b")[] | null {
  if (node.type === "leaf") return node.winId === winId ? path : null;
  return findPath(node.a, winId, [...path, "a"]) ?? findPath(node.b, winId, [...path, "b"]) ?? null;
}

function replaceAtPath(node: WMNode, path: ("a" | "b")[], repl: WMNode): WMNode {
  if (path.length === 0) return repl;
  if (node.type !== "split") return node;
  const [head, ...rest] = path;
  if (head === "a") return { ...node, a: replaceAtPath(node.a, rest, repl) };
  return { ...node, b: replaceAtPath(node.b, rest, repl) };
}

function prune(node: WMNode, removeWinId: string): WMNode | null {
  if (node.type === "leaf") return node.winId === removeWinId ? null : node;

  const a = prune(node.a, removeWinId);
  const b = prune(node.b, removeWinId);

  if (!a && !b) return null;
  if (!a) return b!;
  if (!b) return a!;
  return { ...node, a, b };
}

function updateSplitRatio(node: WMNode, splitId: string, ratio: number): WMNode {
  if (node.type === "leaf") return node;
  if (node.id === splitId) return { ...node, ratio };
  return { ...node, a: updateSplitRatio(node.a, splitId, ratio), b: updateSplitRatio(node.b, splitId, ratio) };
}

export function wmInit() {
  const s = get(wm);
  if (Object.keys(s.windows).length > 0 && s.focusWinId) return;

  const term = makeWindow("terminal");
  wm.set({
    windows: { [term.id]: term },
    order: [term.id],
    root: { type: "leaf", winId: term.id },
    focusWinId: term.id
  });
}

export function wmFocus(winId: string) {
  wm.update((s) => (s.windows[winId] ? { ...s, focusWinId: winId } : s));
}

export function wmList() {
  const s = get(wm);
  return s.order
    .filter((id) => s.windows[id])
    .map((id, idx) => ({
      n: idx + 1,
      id,
      title: s.windows[id].title,
      kind: s.windows[id].kind,
      focused: s.focusWinId === id
    }));
}

export function wmFocusN(n: number) {
  const s = get(wm);
  const id = s.order[n - 1];
  if (id) wmFocus(id);
}

export function wmClose(winId?: string) {
  wm.update((s) => {
    const target = winId ?? s.focusWinId;
    if (!target || !s.windows[target]) return s;

    const nextRoot = prune(s.root, target);
    const nextWindows = { ...s.windows };
    delete nextWindows[target];

    const nextOrder = s.order.filter((id) => id !== target);
    const nextFocus =
      (nextOrder.length ? nextOrder[nextOrder.length - 1] : "") ||
      Object.keys(nextWindows)[0] ||
      "";

    if (!nextRoot || !nextFocus) {
      const term = makeWindow("terminal");
      return {
        windows: { [term.id]: term },
        order: [term.id],
        root: { type: "leaf", winId: term.id },
        focusWinId: term.id
      };
    }

    return { ...s, windows: nextWindows, order: nextOrder, root: nextRoot, focusWinId: nextFocus };
  });
}

export function wmCloseN(n: number) {
  const s = get(wm);
  const id = s.order[n - 1];
  if (id) wmClose(id);
}

export function wmSetSplitRatio(splitId: string, ratio: number) {
  const r = Math.max(0.12, Math.min(0.88, ratio));
  wm.update((s) => ({ ...s, root: updateSplitRatio(s.root, splitId, r) }));
}

export function wmSplitOpen(kind: AppKind, mode: SplitMode = "auto", baseWinId?: string, focus: "new" | "base" = "new") {
  wm.update((s) => {
    const base = baseWinId ?? s.focusWinId;
    if (!base || !s.windows[base]) return s;

    const basePath = findPath(s.root, base);
    if (!basePath) return s;

    const nw = makeWindow(kind);
    const split: WMNode = {
      type: "split",
      id: newSplitId(),
      mode,
      ratio: 0.5,
      a: { type: "leaf", winId: base },
      b: { type: "leaf", winId: nw.id }
    };

    const nextRoot = replaceAtPath(s.root, basePath, split);

    return {
      ...s,
      windows: { ...s.windows, [nw.id]: nw },
      order: [...s.order, nw.id],
      root: nextRoot,
      focusWinId: focus === "base" ? base : nw.id
    };
  });
}

export function wmDuplicate(winId?: string, mode: SplitMode = "auto") {
  const s = get(wm);
  const target = winId ?? s.focusWinId;
  const w = s.windows[target];
  if (!w) return;
  wmSplitOpen(w.kind, mode, target);
}
