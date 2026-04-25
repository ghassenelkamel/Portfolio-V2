<script lang="ts">
  import "./Terminal.css";
  import { goto } from "$app/navigation";
  import KAGHA_ASCII_RAW from "$lib/assets/ascii.txt?raw";
  import { wmClose, wmCloseN, wmDuplicate, wmFocusN, wmList, wmSplitOpen, type AppKind } from "$lib/wm/store";

  const props = $props<{ visible?: boolean; wmMode?: boolean }>();
  const visible = $derived(props.visible ?? true);
  const wmMode = $derived(props.wmMode ?? false);

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

  const programs = ["AboutMe", "Experience", "Work", "Skills", "ContactMe"] as const;
  type Program = typeof programs[number];

  const programToKind: Record<Program, AppKind> = {
    AboutMe: "about",
    Experience: "experience",
    Work: "work",
    Skills: "skills",
    ContactMe: "contact"
  };

  const baseCmds = ["help", "clear", "ls", "exec", "split", "list", "focus", "close", "dup", "ws"] as const;

  function esc(s: string) {
    return s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");
  }

  const starter: Row[] = [{ kind: "banner" }];

  let out = $state<Row[]>([...starter]);
  let bodyEl = $state<HTMLDivElement | null>(null);
  let inputEl = $state<HTMLInputElement | null>(null);

  let cmdline = $state("");

  // completion state (token-based)
  let matches = $state<string[]>([]);
  let matchIndex = $state(0);
  let ghost = $state("");

  // history
  let history = $state<string[]>([]);
  let historyIndex = $state(-1);

  function scrollDown() {
    queueMicrotask(() => {
      if (bodyEl) bodyEl.scrollTop = bodyEl.scrollHeight;
    });
  }

  $effect(() => {
    if (visible) queueMicrotask(() => inputEl?.focus());
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
    if (wmMode) wmSplitOpen(kind, mode);
    else void goto(routeForKind(kind));
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
      return ["terminal", "about", "experience", "work", "skills", "contact", "1", "2", "3", "4", "5", "6"];
    }

    // allow program as second token in "./" not needed
    return [];
  }

  function setSuggestionFrom(value: string) {
    const raw = value;
    const { token } = currentTokenInfo(raw);
    const cands = candidatesFor(raw);

    const pref = token;
    const m = cands
      .filter((w) => w.toLowerCase().startsWith(pref.toLowerCase()))
      .slice(0, 20);

    matches = m;
    matchIndex = 0;

    if (m.length) {
      const first = m[0];
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

    if (c === "help") {
      println(
        `<div class="help modernHelp">
          <div class="helpHead">
            <span class="helpTitle">Command Palette</span>
            <span class="helpMeta">window manager · terminal</span>
          </div>

          <div class="helpGrid">
            <div class="helpRow"><span class="cmd">exec X</span><span class="desc">open X (auto dwindle split)</span></div>
            <div class="helpRow"><span class="cmd">./X</span><span class="desc">alias of <span class="cmd">exec X</span></span></div>
            <div class="helpRow"><span class="cmd">split v X</span><span class="desc">force vertical split (left/right)</span></div>
            <div class="helpRow"><span class="cmd">split h X</span><span class="desc">force horizontal split (top/bottom)</span></div>
            <div class="helpRow"><span class="cmd">dup</span><span class="desc">duplicate focused window (auto)</span></div>
            <div class="helpRow"><span class="cmd">close</span><span class="desc">close focused window</span></div>
            <div class="helpRow"><span class="cmd">close N</span><span class="desc">close window #N from list</span></div>
            <div class="helpRow"><span class="cmd">list</span><span class="desc">list windows</span></div>
            <div class="helpRow"><span class="cmd">focus N</span><span class="desc">focus window #N</span></div>
            <div class="helpRow"><span class="cmd">ws name|N</span><span class="desc">fullscreen route page</span></div>
            <div class="helpRow"><span class="cmd">clear</span><span class="desc">clear terminal output</span></div>
          </div>
        </div>`
      );
      return;
    }

    if (c === "ls") {
      println(
        `<span class="muted">AboutMe</span>  <span class="muted">Experience</span>  <span class="muted">Work</span>  <span class="muted">Skills</span>  <span class="muted">ContactMe</span>`
      );
      return;
    }

    if (c === "clear") {
      out = [...starter];
      scrollDown();
      return;
    }

    if (c === "list") {
      const xs = wmList();
      const rows = xs
        .map((w) => `<div><span class="cmd">${w.n}</span> ${esc(w.title)}${w.focused ? ` <span class="muted">(focused)</span>` : ""}</div>`)
        .join("");
      println(`<div class="help">${rows || "<div class='muted'>no windows</div>"}</div>`);
      return;
    }

    if (c === "dup") {
      wmDuplicate(undefined, "auto");
      println(`<span class="muted">duplicated</span>`);
      return;
    }

    if (c === "close") {
      wmClose();
      println(`<span class="muted">closed focused window</span>`);
      return;
    }

    if (c.startsWith("close ")) {
      const n = Number(c.slice(6).trim());
      if (Number.isFinite(n) && n > 0) {
        wmCloseN(n);
        println(`<span class="muted">closed #${n}</span>`);
      } else {
        println(`<span class="err">invalid close target</span>`);
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
      const arg = c.slice(3).trim().toLowerCase();
      const map: Record<string, string> = {
        "1": "/", "terminal": "/",
        "2": "/about", "about": "/about",
        "3": "/experience", "experience": "/experience",
        "4": "/work", "work": "/work",
        "5": "/skills", "skills": "/skills",
        "6": "/contact", "contact": "/contact"
      };
      const href = map[arg];
      if (href) void goto(href);
      else println(`<span class="err">unknown workspace</span>`);
      return;
    }

    // ./Program alias
    if (c.startsWith("./")) {
      const prog = c.slice(2).trim() as Program;
      const kind = (programToKind as any)[prog] as AppKind | undefined;
      if (!kind) {
        println(`<span class="err">unknown program</span>`);
        return;
      }
      open(kind, "auto");
      return;
    }

    if (c.startsWith("split ")) {
      const rest = c.slice(6).trim();
      const [d, progRaw] = rest.split(/\s+/, 2);
      const prog = progRaw as Program;
      const kind = (programToKind as any)[prog] as AppKind | undefined;
      if (!kind) {
        println(`<span class="err">unknown program</span>`);
        return;
      }
      const mode = d === "h" ? "col" : "row";
      open(kind, mode);
      return;
    }

    if (c.startsWith("exec ")) {
      const prog = c.slice(5).trim() as Program;
      const kind = (programToKind as any)[prog] as AppKind | undefined;
      if (!kind) {
        println(`<span class="err">unknown program</span>`);
        return;
      }
      open(kind, "auto");
      return;
    }

    // allow typing Program directly
    const asProg = c as Program;
    if ((programToKind as any)[asProg]) {
      open(programToKind[asProg], "auto");
      return;
    }

    println(`command not found: <span class="err">${esc(c)}</span>`);
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();
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

<div class="term" onclick={() => inputEl?.focus()}>
  <div class="termBody" bind:this={bodyEl}>
    {#each out as row, i (i)}
      {#if row.kind === "banner"}
        <div class="line">
          <div class="banner">
            <div class="kaghaHero">
              <div
                class="kaghaHit"
                style={`--hx:${hx}; --hy:${hy}; --tx:${hx * 6}px; --ty:${hy * 4}px;`}
                onpointermove={onHeroMove}
                onpointerleave={onHeroLeave}
              >
                <pre class="kaghaAscii base">{KAGHA_ASCII}</pre>
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

    {#if matches.length > 0}
      <div class="suggest">
        {#each matches.slice(0, 10) as s, idx (s)}
          <span class={"sug " + (idx === matchIndex ? "sel" : "")}>{s}</span>
        {/each}
      </div>
    {/if}
  </div>
</div>
