<script module lang="ts">
  type SessionSnapshotShared = {
    out: any[];
    history: string[];
    historyIndex: number;
    contextualSuggestion?: string | null;
  };

  // Shared across component remounts in the same SPA session.
  const terminalSessions = new Map<string, SessionSnapshotShared>();

  const sessionStoragePrefix = "kagha:terminal:session:";
</script>

<script lang="ts">
  import { browser } from "$app/environment";
  import { onDestroy, onMount } from "svelte";
  import "./Terminal.css";
  import { goto } from "$app/navigation";
  import KAGHA_ASCII_RAW from "$lib/assets/ascii.txt?raw";
  import { wmClose, wmDuplicate, wmFocusN, wmList, wmSplitOpen, type AppKind } from "$lib/wm/store";

  const props = $props<{ visible?: boolean; wmMode?: boolean; sessionKey?: string; showBanner?: boolean; compact?: boolean; active?: boolean; suggestCommand?: string; embedded?: boolean; onLayoutAction?: (action: "hide-content" | "reset-layout" | "show-content") => void }>();
  const visible = $derived(props.visible ?? true);
  const wmMode = $derived(props.wmMode ?? false);
  const compact = $derived(props.compact ?? false);
  const active = $derived(props.active ?? true);

  type SessionSnapshot = {
    out: Row[];
    history: string[];
    historyIndex: number;
    contextualSuggestion: string | null;
  };

  function storageKey(sessionId: string) {
    return `${sessionStoragePrefix}${sessionId}`;
  }

  function isRow(value: unknown): value is Row {
    if (!value || typeof value !== "object") return false;
    const kind = (value as { kind?: unknown }).kind;
    if (kind === "banner") return true;
    if (kind === "html") return typeof (value as { html?: unknown }).html === "string";
    return false;
  }

  function loadStoredSession(sessionId: string): SessionSnapshot | undefined {
    if (!browser) return undefined;

    try {
      const raw = window.sessionStorage.getItem(storageKey(sessionId));
      if (!raw) return undefined;

      const parsed = JSON.parse(raw) as {
        out?: unknown;
        history?: unknown;
        historyIndex?: unknown;
        contextualSuggestion?: unknown;
      };

      const out = Array.isArray(parsed.out) ? parsed.out.filter(isRow) : undefined;
      const history = Array.isArray(parsed.history) ? parsed.history.filter((h): h is string => typeof h === "string") : undefined;
      const historyIndex = typeof parsed.historyIndex === "number" ? parsed.historyIndex : undefined;
      const contextualSuggestion = typeof parsed.contextualSuggestion === "string" ? parsed.contextualSuggestion : null;

      if (!out || !history || historyIndex === undefined) return undefined;

      return {
        out,
        history,
        historyIndex,
        contextualSuggestion
      };
    } catch {
      return undefined;
    }
  }

  const PROMPT_USER = "kagha@nixos";
  const PROMPT_DIR = "~";
  const PROMPT_SYM = "❯";

  type Row =
    | { kind: "banner" }
    | { kind: "html"; html: string };

  const KAGHA_ASCII = KAGHA_ASCII_RAW
    .replaceAll("\r\n", "\n")
    // The file has a blank first line for spacing; drop it to save vertical space.
    .replace(/^[ \t]*\n/, "")
    // Avoid an extra trailing blank line in the <pre>.
    .replace(/\n+$/, "")
    // Trim trailing whitespace per-line so the banner centers correctly.
    .split("\n")
    .map((l) => l.replace(/[ \t]+$/, ""))
    .join("\n");

  let hx = $state(0);
  let hy = $state(0);
  let asciiHost = $state<HTMLDivElement | null>(null);
  let asciiBase = $state<HTMLElement | null>(null);
  let asciiFontPx = $state(20);

  const ASCII_MIN_FONT = 9;
  const ASCII_MAX_FONT = 20;

  function fitAscii() {
    if (!asciiHost || !asciiBase) return;

    // Measure at max size first, then scale down to available width.
    asciiFontPx = ASCII_MAX_FONT;
    queueMicrotask(() => {
      if (!asciiHost || !asciiBase) return;
      const available = Math.max(0, asciiHost.clientWidth - 4);
      const required = asciiBase.scrollWidth;
      if (available <= 0 || required <= 0) return;

      const ratio = available / required;
      const fitted = Math.max(ASCII_MIN_FONT, Math.min(ASCII_MAX_FONT, Math.floor(ASCII_MAX_FONT * ratio * 10) / 10));
      asciiFontPx = fitted;
    });
  }

  function onHeroMove(e: PointerEvent) {
    // If the user is dragging (likely selecting/copying), don't animate/tilt.
    if (e.buttons !== 0) return;

    const el = e.currentTarget as HTMLElement;
    const r = el.getBoundingClientRect();
    if (r.width <= 0 || r.height <= 0) return;

    const x = (e.clientX - r.left) / r.width - 0.5;
    const y = (e.clientY - r.top) / r.height - 0.5;

    // Clamp to keep motion subtle and controlled.
    hx = Math.max(-0.5, Math.min(0.5, x));
    hy = Math.max(-0.5, Math.min(0.5, y));
  }

  function onHeroLeave() {
    hx = 0;
    hy = 0;
  }

  onMount(() => {
    const rerenderFit = () => fitAscii();
    const ro = new ResizeObserver(rerenderFit);
    if (asciiHost) ro.observe(asciiHost);

    // First render + after font settles.
    fitAscii();
    const t = setTimeout(fitAscii, 180);
    window.addEventListener("resize", rerenderFit);

    return () => {
      ro.disconnect();
      window.removeEventListener("resize", rerenderFit);
      clearTimeout(t);
    };
  });

  function promptHtml() {
    const nixSvg = `<svg class="nixSvg" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true" focusable="false">
      <path d="M32 6l6.3 10.9h-12.6L32 6z" fill="currentColor" opacity="0.92" />
      <path d="M32 58l-6.3-10.9h12.6L32 58z" fill="currentColor" opacity="0.92" />
      <path d="M6 32l10.9-6.3v12.6L6 32z" fill="currentColor" opacity="0.92" />
      <path d="M58 32l-10.9 6.3V25.7L58 32z" fill="currentColor" opacity="0.92" />
      <path d="M14.9 14.9l12.1 3.2-8.9 8.9-3.2-12.1z" fill="currentColor" opacity="0.78" />
      <path d="M49.1 49.1l-12.1-3.2 8.9-8.9 3.2 12.1z" fill="currentColor" opacity="0.78" />
      <path d="M49.1 14.9l-3.2 12.1-8.9-8.9 12.1-3.2z" fill="currentColor" opacity="0.78" />
      <path d="M14.9 49.1l3.2-12.1 8.9 8.9-12.1 3.2z" fill="currentColor" opacity="0.78" />
    </svg>`;

    return `<span class="p10k" aria-hidden="true">
      <span class="seg user"><span class="nixIcon">${nixSvg}</span><span class="nixSep">|</span><span class="segTxt userTxt">${PROMPT_USER}</span></span>
      <span class="seg tail"><span class="segTxt dirTxt">${PROMPT_DIR}</span><span class="segTxt symTxt">${PROMPT_SYM}</span></span>
    </span>`;
  }

  const programs = ["Terminal", "AboutMe", "Experience", "Work", "Skills", "ContactMe"] as const;
  type Program = typeof programs[number];

  const programToKind: Record<Program, AppKind> = {
    Terminal: "terminal",
    AboutMe: "about",
    Experience: "experience",
    Work: "work",
    Skills: "skills",
    ContactMe: "contact"
  };

  const baseCmds = ["help", "clear", "ls", "exec", "split", "ps", "focus", "close", "dup", "ws"] as const;

  function esc(s: string) {
    return s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");
  }

  const programAliases: Record<string, Program> = {
    terminal: "Terminal",
    term: "Terminal",
    about: "AboutMe",
    aboutme: "AboutMe",
    experience: "Experience",
    work: "Work",
    skills: "Skills",
    contact: "ContactMe",
    contactme: "ContactMe"
  };

  type Workspace = {
    n: number;
    name: "Terminal" | "About" | "Experience" | "Work" | "Skills" | "Contact";
    path: string;
    kind: AppKind;
  };

  const workspaces: Workspace[] = [
    { n: 1, name: "Terminal", path: "/", kind: "terminal" },
    { n: 2, name: "About", path: "/about", kind: "about" },
    { n: 3, name: "Experience", path: "/experience", kind: "experience" },
    { n: 4, name: "Work", path: "/work", kind: "work" },
    { n: 5, name: "Skills", path: "/skills", kind: "skills" },
    { n: 6, name: "Contact", path: "/contact", kind: "contact" }
  ];

  const navSourceStorageKey = "kagha:last-ws-source";

  function workspaceForPath(pathname: string) {
    return workspaces.find((w) => w.path === pathname);
  }

  function rememberWorkspaceSource() {
    if (!browser) return;
    const current = workspaceForPath(window.location.pathname);
    if (!current) return;
    try {
      window.sessionStorage.setItem(navSourceStorageKey, current.name);
    } catch {
      // ignore storage errors
    }
  }

  function workspaceForKind(kind: AppKind) {
    return workspaces.find((w) => w.kind === kind);
  }

  function workspaceForArg(raw: string) {
    const token = raw.trim();
    if (!token) return undefined;

    const asNum = Number(token);
    if (Number.isFinite(asNum) && asNum > 0) return workspaces.find((w) => w.n === asNum);

    const lowered = token.toLowerCase();
    return workspaces.find((w) => w.name.toLowerCase() === lowered);
  }

  function workspaceCommandForKind(kind: AppKind) {
    const ws = workspaceForKind(kind);
    if (!ws) return undefined;
    return `ws ${ws.name}`;
  }

  function setContextWorkspaceHint(kind: AppKind) {
    const hint = workspaceCommandForKind(kind) ?? null;
    contextualSuggestion = hint;
    saveSession();

    if (!hint) return;
    if (cmdline.trim().length > 0) return;

    if (!matches.some((m) => m.toLowerCase() === hint.toLowerCase())) {
      matches = [hint, ...matches];
    }
    matchIndex = 0;
    ghost = hint;
  }

  function workspaceSuggestionForInput(raw: string) {
    const t = raw.trim();
    if (!t) return undefined;

    if (t.startsWith("./")) {
      const prog = resolveProgram(t.slice(2));
      if (!prog) return undefined;
      return workspaceCommandForKind(programToKind[prog]);
    }

    if (!/\s/.test(t)) {
      const prog = resolveProgram(t);
      if (!prog) return undefined;
      return workspaceCommandForKind(programToKind[prog]);
    }

    return undefined;
  }

  type HelpSpec = {
    usage: string;
    summary: string;
    options?: Array<{ flag: string; desc: string }>;
    examples?: string[];
  };

  const helpSpecs: Record<string, HelpSpec> = {
    help: {
      usage: "help [command]",
      summary: "Show help for commands.",
      examples: ["help", "help close", "ls --help"]
    },
    clear: {
      usage: "clear",
      summary: "Clear terminal output (history preserved).",
      examples: ["clear"]
    },
    ls: {
      usage: "ls [-a] [-l]",
      summary: "List available apps in short or long format.",
      options: [
        { flag: "-a", desc: "Include . and .. entries." },
        { flag: "-l", desc: "Use long listing format." },
        { flag: "-h, --help", desc: "Show help for ls." }
      ],
      examples: ["ls", "ls -l", "ls -al"]
    },
    exec: {
      usage: "exec <Program>",
      summary: "Open program (split in WM, route workspace outside WM).",
      examples: ["exec AboutMe", "exec Skills", "exec Terminal"]
    },
    split: {
      usage: "split <v|h> <Program>",
      summary: "Open program in forced split mode.",
      options: [
        { flag: "v", desc: "Vertical split (left/right)." },
        { flag: "h", desc: "Horizontal split (top/bottom)." }
      ],
      examples: ["split v Work", "split h Terminal"]
    },
    ps: {
      usage: "ps",
      summary: "List current windows in open order.",
      examples: ["ps"]
    },
    focus: {
      usage: "focus <N>",
      summary: "Focus window number N from list.",
      examples: ["focus 2"]
    },
    close: {
      usage: "close [N ...]",
      summary: "Close focused window or multiple windows by number.",
      examples: ["close", "close 2", "close 2 4 6", "close 2,4,6"]
    },
    dup: {
      usage: "dup",
      summary: "Duplicate focused window.",
      examples: ["dup"]
    },
    ws: {
      usage: "ws <Name|N>",
      summary: "Go directly to a workspace route.",
      examples: ["ws About", "ws 3", "ws Skills"]
    },
    "./X": {
      usage: "./<Program>",
      summary: "Program alias for exec.",
      examples: ["./about", "./work", "./Terminal"]
    },
    Program: {
      usage: "<Program>",
      summary: "Open program directly by name.",
      examples: ["AboutMe", "Experience", "Work", "Skills", "ContactMe"]
    }
  };

  function resolveProgram(raw: string): Program | undefined {
    const token = raw.trim();
    if (!token) return undefined;
    if ((programToKind as any)[token]) return token as Program;
    return programAliases[token.toLowerCase()];
  }

  function iconForKind(kind: AppKind) {
    switch (kind) {
      case "terminal": return "";
      case "about": return "󰀄";
      case "experience": return "󰄉";
      case "work": return "󰈔";
      case "skills": return "󰓹";
      case "contact": return "󰇮";
    }
  }

  type LsEntry = {
    name: string;
    kind: AppKind;
    perms: string;
    size: string;
    owner: string;
    date: string;
    hidden?: boolean;
  };

  const lsAppEntries: LsEntry[] = [
    { name: "Terminal", kind: "terminal", perms: "-rwxr-xr-x", size: "12k", owner: "kagha", date: "26 Apr 08:37" },
    { name: "AboutMe", kind: "about", perms: "drwxr-xr-x", size: "-", owner: "kagha", date: "26 Apr 08:34" },
    { name: "Experience", kind: "experience", perms: "drwxr-xr-x", size: "-", owner: "kagha", date: "26 Apr 08:31" },
    { name: "Work", kind: "work", perms: "drwxr-xr-x", size: "-", owner: "kagha", date: "26 Apr 08:29" },
    { name: "Skills", kind: "skills", perms: "drwxr-xr-x", size: "-", owner: "kagha", date: "26 Apr 08:28" },
    { name: "ContactMe", kind: "contact", perms: "drwxr-xr-x", size: "-", owner: "kagha", date: "26 Apr 08:26" }
  ];

  const lsDotEntries: LsEntry[] = [
    { name: ".", kind: "terminal", perms: "drwxr-xr-x", size: "-", owner: "kagha", date: "26 Apr 08:37", hidden: true },
    { name: "..", kind: "terminal", perms: "drwxr-xr-x", size: "-", owner: "kagha", date: "26 Apr 08:10", hidden: true }
  ];

  function lsEntries(showAll: boolean) {
    return showAll ? [...lsDotEntries, ...lsAppEntries] : lsAppEntries;
  }

  function lsIconForEntry(entry: LsEntry) {
    if (entry.perms.startsWith("d")) return "";
    return iconForKind(entry.kind);
  }

  function appTokenHtml(kind: AppKind, label: string) {
    return `<span class="appToken"><span class="appIcon kind-${kind}" aria-hidden="true">${iconForKind(kind)}</span><span class="appText">${esc(label)}</span></span>`;
  }

  function lsShortHtml(showAll: boolean) {
    return `<div class="lsShortRows">${lsEntries(showAll).map((entry) => {
      const icon = lsIconForEntry(entry);
      return `<div class="lsShortRow"><span class="appIcon kind-${entry.kind}" aria-hidden="true">${icon}</span><span class="${entry.hidden ? "muted" : "appText"}">${esc(entry.name)}</span></div>`;
    }).join("")}</div>`;
  }

  function lsLongHtml(showAll: boolean) {
    return `<div class="lsTable">${lsEntries(showAll).map((entry) => {
      const icon = lsIconForEntry(entry);
      return `<div class="lsRow"><span class="lsPerm">${entry.perms}</span><span class="lsSize">${entry.size}</span><span class="lsOwner">${entry.owner}</span><span class="lsDate">${entry.date}</span><span class="lsName"><span class="appIcon kind-${entry.kind}" aria-hidden="true">${icon}</span><span class="${entry.hidden ? "muted" : "appText"}">${esc(entry.name)}</span></span></div>`;
    }).join("")}</div>`;
  }

  function parseLsArgs(raw: string) {
    const parts = raw.split(/\s+/).filter(Boolean);
    let showAll = false;
    let longView = false;

    for (const token of parts.slice(1)) {
      if (token.startsWith("-")) {
        for (const f of token.slice(1)) {
          if (f === "a") showAll = true;
          else if (f === "l" || f === "1") longView = true;
          else if (f === "h") return { ok: false as const, msg: "__HELP__" };
          else return { ok: false as const, msg: `ls: unsupported option -- '${f}'` };
        }
        continue;
      }

      if (token === "." || token === "~" || token === "/") continue;
      return { ok: false as const, msg: `ls: cannot access '${token}': No such file or directory` };
    }

    return { ok: true as const, showAll, longView };
  }

  function printUnknownProgram(target?: string) {
    const t = target?.trim();
    println(`<span class="err">unknown program${t ? `: ${esc(t)}` : ""}</span>`);
    println(`<span class="muted">tip:</span> <span class="cmd">ls -al</span> <span class="muted">to list available apps</span>`);
  }

  function escUsage(v: string) {
    return esc(v).replaceAll(" ", "&nbsp;");
  }

  function printHelp(name: keyof typeof helpSpecs) {
    const spec = helpSpecs[name];
    println(`<span class="muted">usage:</span> <span class="cmd">${escUsage(spec.usage)}</span>`);
    println(`${esc(spec.summary)}`);

    if (spec.options?.length) {
      println(`<span class="muted">options:</span>`);
      for (const opt of spec.options) {
        println(`<span class="cmd">${escUsage(opt.flag)}</span> <span class="muted">${esc(opt.desc)}</span>`);
      }
    }

    if (spec.examples?.length) {
      println(`<span class="muted">examples:</span>`);
      for (const ex of spec.examples) {
        println(`<span class="cmd">${escUsage(ex)}</span>`);
      }
    }
  }

  function printGeneralHelp() {
    println(`<span class="muted">usage:</span> <span class="cmd">command [options]</span>`);
    println(`<span class="muted">commands:</span>`);
    println(`<span class="cmd">help</span> <span class="muted">show help pages</span>`);
    println(`<span class="cmd">clear</span> <span class="muted">clear output</span>`);
    println(`<span class="cmd">ls</span> <span class="muted">list apps</span>`);
    println(`<span class="cmd">exec</span> <span class="muted">open app</span>`);
    println(`<span class="cmd">split</span> <span class="muted">open app in split</span>`);
    println(`<span class="cmd">ps</span> <span class="muted">list windows</span>`);
    println(`<span class="cmd">focus</span> <span class="muted">focus window</span>`);
    println(`<span class="cmd">close</span> <span class="muted">close one or many windows</span>`);
    println(`<span class="cmd">dup</span> <span class="muted">duplicate focused window</span>`);
    println(`<span class="cmd">ws</span> <span class="muted">go to workspace</span>`);
    println(`<span class="muted">tip:</span> use <span class="cmd">&lt;command&gt; --help</span> for details.`);
  }

  function isHelpFlag(v: string) {
    return v === "-h" || v === "--help";
  }

  function normalizedHelpCommand(parts: string[]): keyof typeof helpSpecs | undefined {
    if (!parts.length) return undefined;
    const cmd0 = parts[0];
    if (cmd0 === "list") return "ps";
    if (cmd0 in helpSpecs) return cmd0 as keyof typeof helpSpecs;
    if (cmd0.startsWith("./")) return "./X";
    if (resolveProgram(cmd0)) return "Program";
    return undefined;
  }

  function maybeHandleHelp(c: string) {
    if (compact) return false;

    const parts = c.split(/\s+/).filter(Boolean);
    if (!parts.length) return false;

    if (parts.length === 1 && isHelpFlag(parts[0])) {
      printGeneralHelp();
      return true;
    }

    if (parts[0] === "help") {
      if (parts.length === 1) {
        printGeneralHelp();
        return true;
      }

      if (isHelpFlag(parts[1])) {
        printHelp("help");
        return true;
      }

      const k = normalizedHelpCommand([parts[1]]);
      if (k) {
        printHelp(k);
      } else {
        println(`<span class="err">no help topic for: ${esc(parts[1])}</span>`);
      }
      return true;
    }

    if (parts.slice(1).some(isHelpFlag)) {
      const k = normalizedHelpCommand(parts);
      if (k) {
        printHelp(k);
        return true;
      }
    }

    return false;
  }

  const starter: Row[] = (props.showBanner ?? true) ? [{ kind: "banner" }] : [];
  const sessionId = props.sessionKey ?? ((props.wmMode ?? false) ? "wm:default" : "route:main");
  const savedSession = terminalSessions.get(sessionId) ?? loadStoredSession(sessionId);

  let out = $state<Row[]>(savedSession ? [...savedSession.out] : [...starter]);
  let bodyEl = $state<HTMLDivElement | null>(null);
  let inputEl = $state<HTMLInputElement | null>(null);

  let cmdline = $state("");
  let contextualSuggestion = $state<string | null>(savedSession ? (savedSession.contextualSuggestion ?? null) : null);

  // completion state (token-based)
  let matches = $state<string[]>([]);
  let matchIndex = $state(0);
  let ghost = $state("");

  // history
  let history = $state<string[]>(savedSession ? [...savedSession.history] : []);
  let historyIndex = $state(savedSession ? savedSession.historyIndex : -1);
  let lastEmptyEnterAt = 0;

  function saveSession() {
    const snapshot = {
      out: [...out],
      history: [...history],
      historyIndex,
      contextualSuggestion
    };

    terminalSessions.set(sessionId, snapshot);

    if (!browser) return;
    try {
      window.sessionStorage.setItem(storageKey(sessionId), JSON.stringify(snapshot));
    } catch {
      // Ignore storage quota/privacy errors; in-memory map still persists in SPA session.
    }
  }

  $effect(() => {
    saveSession();
  });

  onDestroy(() => {
    saveSession();
  });

  function scrollDown() {
    queueMicrotask(() => {
      const el = bodyEl;
      if (!el) return;

      if (props.embedded) {
        let p: HTMLElement | null = el.parentElement;
        while (p) {
          const st = getComputedStyle(p);
          const canScroll = (st.overflowY === "auto" || st.overflowY === "scroll") && p.scrollHeight > p.clientHeight;
          if (canScroll) {
            p.scrollTop = p.scrollHeight;
            return;
          }
          p = p.parentElement;
        }

        const root = document.scrollingElement as HTMLElement | null;
        if (root) root.scrollTop = root.scrollHeight;
        return;
      }

      el.scrollTop = el.scrollHeight;
    });
  }

  $effect(() => {
    if (visible && active) queueMicrotask(() => inputEl?.focus());
  });

  function println(html: string) {
    out = [...out, { kind: "html", html }];
    scrollDown();
  }

  function routeForKind(kind: AppKind) {
    switch (kind) {
      case "terminal": return "/";
      case "about": return "/about";
      case "experience": return "/experience";
      case "work": return "/work";
      case "skills": return "/skills";
      case "contact": return "/contact";
    }
  }

  function open(kind: AppKind, mode: "auto" | "row" | "col") {
    // Persist immediately in case this action remounts the terminal.
    saveSession();
    if (wmMode) wmSplitOpen(kind, mode);
    else {
      rememberWorkspaceSource();
      void goto(routeForKind(kind));
    }
  }

  function currentTokenInfo(v: string) {
    const endsWithSpace = /\s$/.test(v);
    const parts = v.split(/\s+/).filter((x) => x.length > 0);
    const tokenIndex = endsWithSpace ? parts.length : Math.max(0, parts.length - 1);
    const token = endsWithSpace ? "" : (parts[parts.length - 1] ?? "");
    const head = parts[0] ?? "";
    return { endsWithSpace, parts, tokenIndex, token, head };
  }

  function candidatesFor(v: string) {
    const { parts, endsWithSpace, tokenIndex, head } = currentTokenInfo(v);

    // ./Program completion
    if (v.trimStart().startsWith("./")) {
      return programs.map((p) => "./" + p);
    }

    // token 0 -> base commands + direct programs + ./programs
    if (tokenIndex === 0) {
      return [
        ...baseCmds,
        ...programs,
        ...programs.map((p) => "./" + p)
      ];
    }

    // exec <Program>
    if (head === "exec") {
      return programs;
    }

    // split <v|h> <Program>
    if (head === "split") {
      if (tokenIndex === 1) return ["v", "h"];
      if (tokenIndex === 2) return programs;
      return [];
    }

    // focus <N>
    if (head === "focus") {
      const xs = wmList();
      return xs.map((x) => String(x.n));
    }

    // close [N]
    if (head === "close") {
      const xs = wmList();
      return xs.map((x) => String(x.n));
    }

    // ws <name|N>
    if (head === "ws") {
      return ["Terminal", "About", "Experience", "Work", "Skills", "Contact", "1", "2", "3", "4", "5", "6"];
    }

    if (head === "ls") {
      return ["-a", "-l", "-al", "-la", "-1"];
    }

    // allow program as second token in "./" not needed
    return [];
  }

  function setSuggestionFrom(value: string) {
    if (compact) {
      const raw = value;
      const { token } = currentTokenInfo(raw);
      const cands = candidatesFor(raw);
      const wsHint = workspaceSuggestionForInput(raw);
      const hint = (props.suggestCommand ?? "").trim();
      const pref = token;
      const compactMatches = cands
        .filter((w) => w.toLowerCase().startsWith(pref.toLowerCase()))
        .slice(0, 20);

      if (wsHint && !compactMatches.some((x) => x.toLowerCase() === wsHint.toLowerCase())) {
        compactMatches.unshift(wsHint);
      }

      if (!raw.trim()) {
        const ctxHint = contextualSuggestion;
        if (ctxHint && !compactMatches.some((x) => x.toLowerCase() === ctxHint.toLowerCase())) {
          compactMatches.unshift(ctxHint);
        }

        if (hint && !compactMatches.some((x) => x.toLowerCase() === hint.toLowerCase())) {
          compactMatches.unshift(hint);
        }
      }

      matches = compactMatches;
      matchIndex = 0;

      if (compactMatches.length) {
        ghost = compactMatches[0].slice(token.length);
      } else {
        ghost = "";
      }
      return;
    }

    const raw = value;
    const { token } = currentTokenInfo(raw);
    const cands = candidatesFor(raw);
    const hint = (props.suggestCommand ?? "").trim();
    const wsHint = compact ? undefined : workspaceSuggestionForInput(raw);

    const pref = token;
    const m = cands
      .filter((w) => w.toLowerCase().startsWith(pref.toLowerCase()))
      .slice(0, 20);

    const finalMatches = [...m];

    if (wsHint && !finalMatches.some((x) => x.toLowerCase() === wsHint.toLowerCase())) {
      finalMatches.unshift(wsHint);
    }

    const ctxHint = contextualSuggestion;
    if (!raw.trim() && ctxHint && !finalMatches.some((x) => x.toLowerCase() === ctxHint.toLowerCase())) {
      finalMatches.unshift(ctxHint);
    }

    if (!raw.trim() && hint && !finalMatches.some((x) => x.toLowerCase() === hint.toLowerCase())) {
      finalMatches.unshift(hint);
    }

    matches = finalMatches;
    matchIndex = 0;

    if (finalMatches.length) {
      const first = finalMatches[0];
      ghost = first.slice(pref.length);
    } else {
      ghost = "";
    }
  }

  function replaceCurrentToken(v: string, replacement: string) {
    const endsWithSpace = /\s$/.test(v);
    if (endsWithSpace) return v + replacement;

    const idx = v.lastIndexOf(" ");
    if (idx === -1) return replacement;
    return v.slice(0, idx + 1) + replacement;
  }

  function applyCurrentMatch() {
    if (!matches.length) return;
    const chosen = matches[matchIndex] ?? matches[0];
    const info = currentTokenInfo(cmdline);

    // if completing a base command exactly at token0, append a space
    if (info.tokenIndex === 0 && baseCmds.includes(chosen as any) && cmdline.trim() === chosen) {
      cmdline = chosen + " ";
      setSuggestionFrom(cmdline);
      return;
    }

    // replace only current token
    cmdline = replaceCurrentToken(cmdline, chosen);

    // if first token is exec/split/focus/close/ws and we just completed that token, add a space
    const after = currentTokenInfo(cmdline);
    if (after.parts.length === 1 && (after.head === "exec" || after.head === "split" || after.head === "focus" || after.head === "close" || after.head === "ws")) {
      if (!/\s$/.test(cmdline)) cmdline += " ";
    }

    setSuggestionFrom(cmdline);
  }

  function run(raw: string) {
    const c = raw.trim();
    println(`${promptHtml()}<span class="cmdText">${esc(c)}</span>`);
    if (!c) return;

    history = [...history, c];
    historyIndex = history.length;

    if (maybeHandleHelp(c)) {
      return;
    }

    if (c === "help") {
      if (compact) {
        println(`<span class="muted">help is available in the full terminal</span>`);
        return;
      }
      printGeneralHelp();
      return;
    }

    if (c === "ls" || c.startsWith("ls ")) {
      const parsed = parseLsArgs(c);
      if (!parsed.ok) {
        if (parsed.msg === "__HELP__") {
          if (compact) println(`<span class="muted">use full terminal for command help</span>`);
          else printHelp("ls");
          return;
        }
        println(`<span class="err">${esc(parsed.msg)}</span>`);
        return;
      }

      println(parsed.longView ? lsLongHtml(parsed.showAll) : lsShortHtml(parsed.showAll));
      return;
    }

    if (c === "clear") {
      out = [...starter];
      scrollDown();
      return;
    }

    if (c === "ps" || c === "list") {
      const xs = wmList();
      const rows = xs
        .map((w) => `<div class="appRow"><span class="appIdx">${w.n}.</span>${appTokenHtml(w.kind, w.title)}${w.focused ? `<span class="focusMark">*</span>` : ""}<span class="muted">${w.kind}</span></div>`)
        .join("");
      println(rows || "<span class='muted'>no windows</span>");
      return;
    }

    if (c === "dup") {
      wmDuplicate(undefined, "auto");
      println(`<span class="muted">duplicated</span>`);
      return;
    }

    if (!wmMode && props.onLayoutAction && (c === "close" || c === "close page" || c === "close data")) {
      props.onLayoutAction("hide-content");
      println(`<span class="muted">page content hidden (double Enter to reset)</span>`);
      return;
    }

    if (!wmMode && props.onLayoutAction && (c === "open" || c === "open page" || c === "show" || c === "show page")) {
      props.onLayoutAction("show-content");
      println(`<span class="muted">page content shown</span>`);
      return;
    }

    if (!wmMode && props.onLayoutAction && c === "reset") {
      props.onLayoutAction("reset-layout");
      out = [...starter];
      println(`<span class="muted">layout reset</span>`);
      return;
    }

    if (c === "close") {
      wmClose();
      println(`<span class="muted">closed focused window</span>`);
      return;
    }

    if (c.startsWith("close ")) {
      const rawTargets = c.slice(6).split(/[\s,]+/).map((x) => x.trim()).filter(Boolean);
      const listed = wmList();
      const nToId = new Map<number, string>(listed.map((w) => [w.n, w.id]));

      const ids: string[] = [];
      const invalid: string[] = [];

      for (const token of rawTargets) {
        const n = Number(token);
        if (!Number.isFinite(n) || n <= 0) {
          invalid.push(token);
          continue;
        }
        const id = nToId.get(n);
        if (!id) {
          invalid.push(token);
          continue;
        }
        if (!ids.includes(id)) ids.push(id);
      }

      for (const id of ids) wmClose(id);

      if (ids.length) {
        println(`<span class="muted">closed ${ids.length} window${ids.length > 1 ? "s" : ""}</span>`);
      }
      if (invalid.length) {
        println(`<span class="err">invalid close target: ${esc(invalid.join(", "))}</span>`);
      }
      return;
    }

    if (c.startsWith("focus ")) {
      const n = Number(c.slice(6).trim());
      if (Number.isFinite(n) && n > 0) {
        wmFocusN(n);
        println(`<span class="muted">focused #${n}</span>`);
      } else {
        println(`<span class="err">invalid focus</span>`);
      }
      return;
    }

    if (c.startsWith("ws ")) {
      const arg = c.slice(3).trim();
      const ws = workspaceForArg(arg);
      if (ws) {
        if (arg.toLowerCase() !== ws.name.toLowerCase() || arg !== ws.name) {
          println(`<span class="muted">workspace:</span> <span class="cmd">ws ${ws.name}</span>`);
        }
        rememberWorkspaceSource();
        saveSession();
        void goto(ws.path);
      }
      else println(`<span class="err">unknown workspace</span>`);
      return;
    }

    // ./Program alias
    if (c.startsWith("./")) {
      const prog = resolveProgram(c.slice(2).trim());
      const kind = prog ? programToKind[prog] : undefined;
      if (!kind) {
        printUnknownProgram(c.slice(2).trim());
        return;
      }
      setContextWorkspaceHint(kind);
      if (!wmMode) {
        const wsCmd = workspaceCommandForKind(kind);
        if (wsCmd) println(`<span class="muted">hint:</span> <span class="cmd">${esc(wsCmd)}</span>`);
      }
      open(kind, "auto");
      return;
    }

    if (c.startsWith("split ")) {
      const rest = c.slice(6).trim();
      const [d, progRaw] = rest.split(/\s+/, 2);
      const prog = resolveProgram(progRaw ?? "");
      const kind = prog ? programToKind[prog] : undefined;
      if (!kind) {
        printUnknownProgram(progRaw ?? "");
        return;
      }
      setContextWorkspaceHint(kind);
      const mode = d === "h" ? "col" : "row";
      open(kind, mode);
      return;
    }

    if (c.startsWith("exec ")) {
      const prog = resolveProgram(c.slice(5).trim());
      const kind = prog ? programToKind[prog] : undefined;
      if (!kind) {
        printUnknownProgram(c.slice(5).trim());
        return;
      }
      setContextWorkspaceHint(kind);
      if (!wmMode) {
        const wsCmd = workspaceCommandForKind(kind);
        if (wsCmd) println(`<span class="muted">hint:</span> <span class="cmd">${esc(wsCmd)}</span>`);
      }
      open(kind, "auto");
      return;
    }

    // allow typing Program directly
    const asProg = resolveProgram(c);
    if (asProg) {
      setContextWorkspaceHint(programToKind[asProg]);
      if (!wmMode) {
        const wsCmd = workspaceCommandForKind(programToKind[asProg]);
        if (wsCmd) println(`<span class="muted">hint:</span> <span class="cmd">${esc(wsCmd)}</span>`);
      }
      open(programToKind[asProg], "auto");
      return;
    }

    println(`command not found: <span class="err">${esc(c)}</span>`);
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();

      const now = Date.now();
      if (!cmdline.trim()) {
        if (!wmMode && props.onLayoutAction && now - lastEmptyEnterAt <= 450) {
          props.onLayoutAction("reset-layout");
          out = [...starter];
          matches = [];
          matchIndex = 0;
          ghost = "";
          lastEmptyEnterAt = 0;
          queueMicrotask(() => inputEl?.focus());
          return;
        }
        lastEmptyEnterAt = now;
      } else {
        lastEmptyEnterAt = 0;
      }

      const v = cmdline;
      cmdline = "";
      matches = [];
      matchIndex = 0;
      ghost = "";
      run(v);
      queueMicrotask(() => inputEl?.focus());
      return;
    }

    if (e.key === "Tab") {
      e.preventDefault();
      applyCurrentMatch();
      queueMicrotask(() => inputEl?.focus());
      return;
    }

    if (e.key === "ArrowUp") {
      e.preventDefault();
      if (!history.length) return;
      historyIndex = Math.max(0, historyIndex - 1);
      cmdline = history[historyIndex] ?? cmdline;
      setSuggestionFrom(cmdline);
      queueMicrotask(() => inputEl?.setSelectionRange(cmdline.length, cmdline.length));
      return;
    }

    if (e.key === "ArrowDown") {
      e.preventDefault();
      if (!history.length) return;
      historyIndex = Math.min(history.length, historyIndex + 1);
      cmdline = historyIndex >= history.length ? "" : (history[historyIndex] ?? "");
      setSuggestionFrom(cmdline);
      queueMicrotask(() => inputEl?.setSelectionRange(cmdline.length, cmdline.length));
      return;
    }

    // Keep cursor at end to match the rendered terminal cursor.
    queueMicrotask(() => inputEl?.setSelectionRange(cmdline.length, cmdline.length));
  }

  $effect(() => { setSuggestionFrom(cmdline); });
</script>

<div class={"term" + (compact ? " compact" : "") + ((props.embedded ?? false) ? " embedded" : "")} onclick={() => inputEl?.focus()}>
  <div class="termBody" bind:this={bodyEl}>
    {#each out as row, i (i)}
      {#if row.kind === "banner"}
        <div class="line">
          <div class="banner">
            <div class="kaghaHero">
              <div
                class="kaghaHit"
                bind:this={asciiHost}
                style={`--hx:${hx}; --hy:${hy}; --tx:${hx * 6}px; --ty:${hy * 4}px; --ascii-font:${asciiFontPx}px;`}
                onpointermove={onHeroMove}
                onpointerleave={onHeroLeave}
              >
                <pre class="kaghaAscii base" bind:this={asciiBase}>{KAGHA_ASCII}</pre>
                <pre class="kaghaAscii layer red" aria-hidden="true">{KAGHA_ASCII}</pre>
                <pre class="kaghaAscii layer amber" aria-hidden="true">{KAGHA_ASCII}</pre>
              </div>
            </div>
          </div>
        </div>
      {:else}
        <div class="line html">
          {@html row.html}
        </div>
      {/if}
    {/each}

    <div class="promptRow">
      <span class="p10k" aria-hidden="true">
        <span class="seg user">
          <span class="nixIcon" aria-hidden="true">
            <svg class="nixSvg" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M32 6l6.3 10.9h-12.6L32 6z" fill="currentColor" opacity="0.92" />
              <path d="M32 58l-6.3-10.9h12.6L32 58z" fill="currentColor" opacity="0.92" />
              <path d="M6 32l10.9-6.3v12.6L6 32z" fill="currentColor" opacity="0.92" />
              <path d="M58 32l-10.9 6.3V25.7L58 32z" fill="currentColor" opacity="0.92" />
              <path d="M14.9 14.9l12.1 3.2-8.9 8.9-3.2-12.1z" fill="currentColor" opacity="0.78" />
              <path d="M49.1 49.1l-12.1-3.2 8.9-8.9 3.2 12.1z" fill="currentColor" opacity="0.78" />
              <path d="M49.1 14.9l-3.2 12.1-8.9-8.9 12.1-3.2z" fill="currentColor" opacity="0.78" />
              <path d="M14.9 49.1l3.2-12.1 8.9 8.9-12.1 3.2z" fill="currentColor" opacity="0.78" />
            </svg>
          </span>
          <span class="nixSep" aria-hidden="true">|</span>
          <span class="segTxt userTxt">{PROMPT_USER}</span>
        </span>
        <span class="seg tail"><span class="segTxt dirTxt">{PROMPT_DIR}</span><span class="segTxt symTxt">{PROMPT_SYM}</span></span>
      </span>

      <span class="inputLayer">
        <span class="ghostLine" aria-hidden="true">
          <span class="typed">{cmdline}</span><span class="ghost">{ghost}</span>
        </span>

        <input
          class="in"
          bind:this={inputEl}
          value={cmdline}
          oninput={(e) => {
            cmdline = (e.currentTarget as HTMLInputElement).value;
            queueMicrotask(() => inputEl?.setSelectionRange(cmdline.length, cmdline.length));
          }}
          onkeydown={onKeydown}
          autocomplete="off"
          autocapitalize="off"
          spellcheck="false"
        />
      </span>
    </div>

    {#if !compact && matches.length > 0}
      <div class="suggest">
        {#each matches.slice(0, 10) as s, idx (s)}
          <span class={"sug " + (idx === matchIndex ? "sel" : "")}>{s}</span>
        {/each}
      </div>
    {/if}
  </div>
</div>
