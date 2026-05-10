<script module lang="ts">
  type SessionSnapshotShared = {
    out: any[];
    history: string[];
    historyIndex: number;
    contextualSuggestion?: string | null;
  };

  type BannerSessionShared = any;

  // Shared across component remounts in the same SPA session.
  const terminalSessions = new Map<string, SessionSnapshotShared>();
  const terminalBannerSessions = new Map<string, BannerSessionShared>();

  const sessionStoragePrefix = "kagha:terminal:session:";
</script>

<script lang="ts">
  import { browser } from "$app/environment";
  import { onDestroy, onMount } from "svelte";
  import "./Terminal.css";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import KAGHA_ASCII_RAW from "$lib/assets/ascii.txt?raw";
  import { wmClose, wmDuplicate, wmFocus, wmFocusN, wmList, wmSplitOpen, type AppKind } from "$lib/wm/store";

  const props = $props<{ visible?: boolean; wmMode?: boolean; sessionKey?: string; showBanner?: boolean; compact?: boolean; active?: boolean; suggestCommand?: string; embedded?: boolean; onLayoutAction?: (action: "hide-content" | "reset-layout" | "show-content") => void }>();
  const visible = $derived(props.visible ?? true);
  const wmMode = $derived(props.wmMode ?? false);
  const compact = $derived(props.compact ?? false);
  const active = $derived(props.active ?? true);
  const isFr = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr"));

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

  const KAGHA_LINES = KAGHA_ASCII.split("\n");

  type BannerChunk = { x: number; y: number; life: number; line: number };
  type BannerParticle = { x: number; y: number; vx: number; vy: number; life: number; size: number };
  type BannerRuntimeSnapshot = {
    gameActive: boolean;
    won: boolean;
    total: number;
    remaining: number;
    chunks: BannerChunk[];
    lineAlive: number[];
    lineTotals: number[];
    activeLines: number;
    step: number;
    lineHeight: number;
    sourceW: number;
    sourceH: number;
    snake: {
      x: number;
      y: number;
      radius: number;
      speed: number;
      dirX: number;
      dirY: number;
      targetX: number;
      targetY: number;
      trail: Array<{ x: number; y: number }>;
      maxTrailLen: number;
    };
  };

  let hx = $state(0);
  let hy = $state(0);
  let asciiHost = $state<HTMLDivElement | null>(null);
  let asciiBase = $state<HTMLElement | null>(null);
  let asciiFontPx = $state(20);
  let bannerCanvas = $state<HTMLCanvasElement | null>(null);

  let bannerCtx: CanvasRenderingContext2D | null = null;
  let bannerChunks: BannerChunk[] = [];
  let bannerParticles: BannerParticle[] = [];
  let bannerTotal = $state(0);
  let bannerRemaining = $state(0);
  let bannerWon = $state(false);
  let bannerGameActive = $state(false);
  const bannerHasScars = $derived(bannerTotal > 0 && bannerRemaining < bannerTotal);
   const bannerShowCanvas = $derived(bannerGameActive || bannerWon || bannerTotal > 0);
  const bannerProgressPct = $derived(bannerTotal ? Math.round(((bannerTotal - bannerRemaining) / bannerTotal) * 100) : 0);

  let bannerViewW = $state(1);
  let bannerViewH = $state(1);
  let bannerSourceW = $state(1);
  let bannerSourceH = $state(1);
  let bannerScaleX = $state(1);
  let bannerScaleY = $state(1);
  let bannerStep = $state(3);
  let bannerLineHeight = $state(1);
  let bannerLineAlive: number[] = [];
  let bannerLineTotals: number[] = [];
  let bannerActiveLines = $state(0);
  let bannerLastTs = 0;
  let bannerRaf = 0;
  let bannerGameHeight = $derived(activeBannerHeight());
  let bannerAudioCtx: AudioContext | null = null;
  let bannerLastEatToneAt = 0;
  let bannerLastWrapToneAt = 0;
  let bannerWinSoundPlayed = false;

  const redScaleY = $derived.by(() => {
    if (bannerLineAlive.length === 0) return 0;
    const alive = bannerLineAlive.filter(Boolean).length;
    return alive / KAGHA_LINES.length;
  });

  const amberScaleY = $derived.by(() => {
    if (bannerLineAlive.length === 0) return 0;
    const alive = bannerLineAlive.filter(Boolean).length;
    return alive / KAGHA_LINES.length;
  });

  const WIN_ASCII_FRAMES = [
    ["*** KAGHA CLEARED ***"],
    ["<\\o/> KAGHA WIN <\\o/>"],
    ["*** GREAT RUN! ***"]
  ] as const;

  const bannerSnake = {
    x: 0,
    y: 0,
    radius: 9,
    speed: 200,
    dirX: 1,
    dirY: 0,
    targetX: 1,
    targetY: 0,
    trail: [] as Array<{ x: number; y: number }>,
    maxTrailLen: 220
  };

  const ASCII_MIN_FONT = 9;
  const ASCII_MAX_FONT = 20;

  function clamp(n: number, lo: number, hi: number) {
    return Math.max(lo, Math.min(hi, n));
  }

  function normalize2d(x: number, y: number) {
    const len = Math.hypot(x, y) || 1;
    return { x: x / len, y: y / len };
  }

  function fitAscii() {
    if (!asciiHost || !asciiBase) return;

    // Measure at max size first, then scale down to available width.
    asciiFontPx = ASCII_MAX_FONT;
    queueMicrotask(() => {
      if (!asciiHost || !asciiBase) return;
      requestAnimationFrame(() => {
        if (!asciiHost || !asciiBase) return;
        const available = Math.max(0, asciiHost.clientWidth - 4);
        const required = asciiBase.scrollWidth;
        if (available <= 0 || required <= 0) return;

        const ratio = available / required;
        const fitted = Math.max(ASCII_MIN_FONT, Math.min(ASCII_MAX_FONT, Math.floor(ASCII_MAX_FONT * ratio * 10) / 10));
        asciiFontPx = fitted;
        if (bannerShowCanvas) {
          queueMicrotask(() => {
            if (prepareBannerCanvas()) drawBannerGame();
          });
        }
      });
    });
  }

  function ensureBannerAudioCtx() {
    if (!browser) return null;
    if (!bannerAudioCtx) {
      try {
        bannerAudioCtx = new AudioContext();
      } catch {
        return null;
      }
    }
    if (bannerAudioCtx.state === "suspended") void bannerAudioCtx.resume();
    return bannerAudioCtx;
  }

  function playBannerTone(freq: number, duration = 0.06, gainMax = 0.04, type: OscillatorType = "triangle") {
    const ctx = ensureBannerAudioCtx();
    if (!ctx) return;

    const t0 = ctx.currentTime;
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.type = type;
    osc.frequency.setValueAtTime(freq, t0);
    gain.gain.setValueAtTime(0.0001, t0);
    gain.gain.exponentialRampToValueAtTime(gainMax, t0 + 0.01);
    gain.gain.exponentialRampToValueAtTime(0.0001, t0 + duration);
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.start(t0);
    osc.stop(t0 + duration + 0.01);
  }

  function playEatSound(intensity = 1) {
    const ctx = ensureBannerAudioCtx();
    if (!ctx) return;
    if (ctx.currentTime - bannerLastEatToneAt < 0.03) return;
    bannerLastEatToneAt = ctx.currentTime;
    const f = 340 + Math.min(420, intensity * 16) + Math.random() * 44;
    playBannerTone(f, 0.055, 0.032, "triangle");
  }

  function playWrapSound() {
    const ctx = ensureBannerAudioCtx();
    if (!ctx) return;
    if (ctx.currentTime - bannerLastWrapToneAt < 0.06) return;
    bannerLastWrapToneAt = ctx.currentTime;
    playBannerTone(220 + Math.random() * 34, 0.05, 0.028, "square");
  }

  function primeBannerAudio() {
    const ctx = ensureBannerAudioCtx();
    if (!ctx) return;
    if (ctx.state === "suspended") {
      void ctx.resume().then(() => {
        playBannerTone(300, 0.035, 0.012, "square");
      });
      return;
    }
    playBannerTone(300, 0.035, 0.012, "square");
  }

  function playWinSound() {
    const ctx = ensureBannerAudioCtx();
    if (!ctx) {
      println(`<span class="muted">\\a</span>`);
      return;
    }

    const base = ctx.currentTime;
    const notes = [523.25, 659.25, 783.99, 1046.5];
    notes.forEach((freq, i) => {
      const t = base + i * 0.085;
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();
      osc.type = "triangle";
      osc.frequency.setValueAtTime(freq, t);
      gain.gain.setValueAtTime(0.0001, t);
      gain.gain.exponentialRampToValueAtTime(0.045, t + 0.01);
      gain.gain.exponentialRampToValueAtTime(0.0001, t + 0.08);
      osc.connect(gain);
      gain.connect(ctx.destination);
      osc.start(t);
      osc.stop(t + 0.09);
    });
  }

  function activeBannerHeight() {
    const minH = Math.max(10, bannerSnake.radius * 2 + 4);
    const fromLines = Math.max(0, bannerActiveLines) * bannerLineHeight * bannerScaleY;
    if (fromLines <= 0) return minH;
    return Math.max(minH, Math.min(bannerViewH, fromLines));
  }

  function captureBannerSession(): BannerRuntimeSnapshot | null {
    if (!(bannerGameActive || bannerHasScars || bannerWon || bannerTotal > 0)) return null;

    return {
      gameActive: bannerGameActive,
      won: bannerWon,
      total: bannerTotal,
      remaining: bannerRemaining,
      chunks: bannerChunks.map((c) => ({ ...c })),
      lineAlive: [...bannerLineAlive],
      lineTotals: [...bannerLineTotals],
      activeLines: bannerActiveLines,
      step: bannerStep,
      lineHeight: bannerLineHeight,
      sourceW: bannerSourceW,
      sourceH: bannerSourceH,
      snake: {
        x: bannerSnake.x,
        y: bannerSnake.y,
        radius: bannerSnake.radius,
        speed: bannerSnake.speed,
        dirX: bannerSnake.dirX,
        dirY: bannerSnake.dirY,
        targetX: bannerSnake.targetX,
        targetY: bannerSnake.targetY,
        trail: bannerSnake.trail.map((p) => ({ ...p })),
        maxTrailLen: bannerSnake.maxTrailLen
      }
    };
  }

  function restoreBannerSession(snapshot?: BannerRuntimeSnapshot) {
    if (!snapshot) return false;
    if (!prepareBannerCanvas()) return false;

    bannerSourceW = Math.max(1, snapshot.sourceW || 1);
    bannerSourceH = Math.max(1, snapshot.sourceH || 1);
    bannerScaleX = bannerViewW / bannerSourceW;
    bannerScaleY = bannerViewH / bannerSourceH;

    bannerChunks = snapshot.chunks.map((c) => ({ ...c }));
    bannerLineAlive = [...snapshot.lineAlive];
    bannerLineTotals = [...snapshot.lineTotals];
    bannerActiveLines = snapshot.activeLines;
    bannerStep = snapshot.step;
    bannerLineHeight = snapshot.lineHeight;
    bannerTotal = snapshot.total;
    bannerRemaining = snapshot.remaining;
    bannerWon = snapshot.won;
    bannerGameActive = snapshot.gameActive;

    bannerSnake.x = snapshot.snake.x;
    bannerSnake.y = snapshot.snake.y;
    bannerSnake.radius = snapshot.snake.radius;
    bannerSnake.speed = snapshot.snake.speed;
    bannerSnake.dirX = snapshot.snake.dirX;
    bannerSnake.dirY = snapshot.snake.dirY;
    bannerSnake.targetX = snapshot.snake.targetX;
    bannerSnake.targetY = snapshot.snake.targetY;
    bannerSnake.trail = snapshot.snake.trail.map((p) => ({ ...p }));
    bannerSnake.maxTrailLen = snapshot.snake.maxTrailLen;

    bannerLastTs = 0;
    drawBannerGame();
    return true;
  }

  function persistBannerSessionNow() {
    const snap = captureBannerSession();
    if (snap) terminalBannerSessions.set(sessionId, snap);
    else terminalBannerSessions.delete(sessionId);
  }

  function collapseEatenLine(lineIdx: number) {
    if (lineIdx < 0 || lineIdx >= bannerLineAlive.length) return;
    const shift = bannerLineHeight;

    for (const c of bannerChunks) {
      if (c.life <= 0) continue;
      if (c.line > lineIdx) {
        c.line -= 1;
        c.y -= shift;
      }
    }

    bannerLineAlive.splice(lineIdx, 1);
    bannerLineTotals.splice(lineIdx, 1);
    bannerActiveLines = bannerLineAlive.length;
  }



  function prepareBannerCanvas() {
    if (!browser || !bannerCanvas || !asciiBase) return false;

    let rect = asciiBase.getBoundingClientRect();
    let cssW = Math.max(1, Math.ceil(rect.width));
    let cssH = Math.max(1, Math.ceil(rect.height));

    if ((cssW <= 1 || cssH <= 1) && asciiHost) {
      rect = asciiHost.getBoundingClientRect();
      cssW = Math.max(1, Math.ceil(rect.width));
      cssH = Math.max(1, Math.ceil(rect.height));
    }

    const ratio = Math.max(1, window.devicePixelRatio || 1);

    bannerCanvas.width = Math.max(1, Math.round(cssW * ratio));
    bannerCanvas.height = Math.max(1, Math.round(cssH * ratio));
    bannerCanvas.style.width = `${cssW}px`;
    bannerCanvas.style.height = `${cssH}px`;

    const ctx = bannerCanvas.getContext("2d");
    if (!ctx) return false;

    ctx.setTransform(ratio, 0, 0, ratio, 0, 0);
    bannerCtx = ctx;
    bannerViewW = cssW;
    bannerViewH = cssH;
    return true;
  }

  function buildBannerChunks() {
    if (!browser) return;

    const fontPx = asciiFontPx;
    const off = document.createElement("canvas");
    const octx = off.getContext("2d", { willReadFrequently: true });
    if (!octx) return;

    const computed = getComputedStyle(asciiBase!);
    const computedFontPx = parseFloat(computed.fontSize) || fontPx;
    const computedLineHeight = parseFloat(computed.lineHeight) || (computedFontPx * 1.06);
    octx.font = `${computed.fontWeight} ${computedFontPx}px ${computed.fontFamily}`;

    let maxWidth = 0;
    for (const line of KAGHA_LINES) {
      maxWidth = Math.max(maxWidth, Math.ceil(octx.measureText(line).width));
    }

    const padX = Math.max(2, Math.ceil(computedFontPx * 0.22));
    const textW = Math.max(1, maxWidth);
    const textH = Math.max(1, Math.ceil(computedLineHeight * KAGHA_LINES.length));

    off.width = textW + padX * 2;
    off.height = textH + 2;
    octx.clearRect(0, 0, off.width, off.height);
    octx.font = `${computed.fontWeight} ${computedFontPx}px ${computed.fontFamily}`;
    octx.textBaseline = "top";
    octx.fillStyle = "#fff";

    KAGHA_LINES.forEach((line, i) => {
      octx.fillText(line, padX, i * computedLineHeight);
    });

    bannerStep = Math.max(1, Math.round(computedFontPx * 0.14));
    bannerLineHeight = computedLineHeight;
    bannerLineAlive = Array.from({ length: KAGHA_LINES.length }, () => 0);
    bannerLineTotals = Array.from({ length: KAGHA_LINES.length }, () => 0);
    bannerActiveLines = KAGHA_LINES.length;
    const data = octx.getImageData(0, 0, off.width, off.height).data;
    const next: BannerChunk[] = [];

    for (let y = 0; y < off.height; y += bannerStep) {
      for (let x = 0; x < off.width; x += bannerStep) {
        const sx = Math.min(off.width - 1, x + (bannerStep >> 1));
        const sy = Math.min(off.height - 1, y + (bannerStep >> 1));
        const lineIdx = clamp(Math.floor(y / computedLineHeight), 0, KAGHA_LINES.length - 1);
        const alpha = data[(sy * off.width + sx) * 4 + 3];
        if (alpha > 16) {
          next.push({ x: x - padX, y, life: 1, line: lineIdx });
          bannerLineAlive[lineIdx] += 1;
          bannerLineTotals[lineIdx] += 1;
        }
      }
    }

    bannerSourceW = textW;
    bannerSourceH = textH;
    bannerScaleX = bannerViewW / bannerSourceW;
    bannerScaleY = bannerViewH / bannerSourceH;
    bannerChunks = next;
    bannerParticles = [];
    bannerTotal = next.length;
    bannerRemaining = next.length;
    bannerWon = false;
    bannerWinSoundPlayed = false;
  }

  function resetBannerSnake() {
    const r = Math.max(4, asciiFontPx * 0.42);
    bannerSnake.radius = r;
    bannerSnake.speed = Math.max(120, asciiFontPx * 16.5);
    bannerSnake.maxTrailLen = Math.max(100, bannerViewW * 0.46);
    bannerSnake.x = Math.max(r + 2, bannerViewW * 0.08);
    bannerSnake.y = bannerViewH * 0.55;
    bannerSnake.dirX = 1;
    bannerSnake.dirY = 0;
    bannerSnake.targetX = 1;
    bannerSnake.targetY = 0;
    bannerSnake.trail = [{ x: bannerSnake.x, y: bannerSnake.y }];
  }

  function restartBannerGame() {
    if (!(props.showBanner ?? true)) return;
    if (!prepareBannerCanvas()) return;
    buildBannerChunks();
    resetBannerSnake();
    bannerLastTs = 0;
    return bannerChunks.length > 0;
  }

   function startBannerGame(reset = false) {
    if (!(props.showBanner ?? true)) return false;
    primeBannerAudio();
    bannerGameActive = true;
    if (reset || !bannerChunks.length || bannerWon) {
      const ok = !!restartBannerGame();
      if (ok) {
        prepareBannerCanvas();
        drawBannerGame();
        // Force animation frame start
        bannerLastTs = 0;
        if (bannerRaf) cancelAnimationFrame(bannerRaf);
        bannerRaf = requestAnimationFrame(tickBannerGame);
      }
      persistBannerSessionNow();
      return ok;
    }
    if (!prepareBannerCanvas()) return false;
    if (!bannerSnake.trail.length) resetBannerSnake();
    bannerLastTs = 0;
    if (bannerRaf) cancelAnimationFrame(bannerRaf);
    bannerRaf = requestAnimationFrame(tickBannerGame);
    persistBannerSessionNow();
    return true;
  }

  function stopBannerGame() {
    if (!bannerGameActive) return false;
    bannerGameActive = false;
    bannerParticles = [];
    drawBannerGame();
    persistBannerSessionNow();
    return true;
  }

  function exitBannerGame() {
    bannerGameActive = false;
    bannerWon = false;
    bannerTotal = 0;
    bannerRemaining = 0;
    bannerChunks = [];
    bannerParticles = [];
    bannerLineAlive = [];
    bannerLineTotals = [];
    bannerActiveLines = 0;
    bannerWinSoundPlayed = false;
    persistBannerSessionNow();
    if (bannerCtx) bannerCtx.clearRect(0, 0, bannerViewW, bannerViewH);
    // Cancel animation frame on exit
    if (bannerRaf) cancelAnimationFrame(bannerRaf);
    bannerRaf = 0;
  }

  function emitBannerParticles(x: number, y: number, count: number) {
    for (let i = 0; i < count; i += 1) {
      const a = Math.random() * Math.PI * 2;
      const s = 22 + Math.random() * 72;
      bannerParticles.push({
        x,
        y,
        vx: Math.cos(a) * s,
        vy: Math.sin(a) * s,
        life: 0.26 + Math.random() * 0.22,
        size: 0.8 + Math.random() * 1.4
      });
    }
  }

  function updateBannerParticles(dt: number) {
    for (let i = bannerParticles.length - 1; i >= 0; i -= 1) {
      const p = bannerParticles[i];
      p.life -= dt;
      p.x += p.vx * dt;
      p.y += p.vy * dt;
      p.vx *= 0.96;
      p.vy *= 0.96;
      if (p.life <= 0) bannerParticles.splice(i, 1);
    }
  }

  function updateBannerGame(dt: number) {
    if (!bannerGameActive) return;
    if (!bannerCtx || !bannerChunks.length || !bannerTotal) return;

    const turn = clamp(17 * dt, 0, 1);
    bannerSnake.dirX += (bannerSnake.targetX - bannerSnake.dirX) * turn;
    bannerSnake.dirY += (bannerSnake.targetY - bannerSnake.dirY) * turn;
    const nd = normalize2d(bannerSnake.dirX, bannerSnake.dirY);
    bannerSnake.dirX = nd.x;
    bannerSnake.dirY = nd.y;

    bannerSnake.x += bannerSnake.dirX * bannerSnake.speed * dt;
    bannerSnake.y += bannerSnake.dirY * bannerSnake.speed * dt;

    const wrapPad = bannerSnake.radius + 1;
    let wrapped = false;
    if (bannerSnake.x > bannerViewW + wrapPad) {
      bannerSnake.x = -wrapPad;
      wrapped = true;
    } else if (bannerSnake.x < -wrapPad) {
      bannerSnake.x = bannerViewW + wrapPad;
      wrapped = true;
    }

    const activeH = activeBannerHeight();
    if (bannerSnake.y > activeH + wrapPad) {
      bannerSnake.y = -wrapPad;
      wrapped = true;
    } else if (bannerSnake.y < -wrapPad) {
      bannerSnake.y = activeH + wrapPad;
      wrapped = true;
    }
    if (wrapped) playWrapSound();

    bannerSnake.trail.unshift({ x: bannerSnake.x, y: bannerSnake.y });
    let len = 0;
    for (let i = 1; i < bannerSnake.trail.length; i += 1) {
      const a = bannerSnake.trail[i - 1];
      const b = bannerSnake.trail[i];
      len += Math.hypot(a.x - b.x, a.y - b.y);
      if (len > bannerSnake.maxTrailLen) {
        bannerSnake.trail.length = i + 1;
        break;
      }
    }

    if (!bannerWon && bannerRemaining > 0) {
      const hitRadius = bannerSnake.radius + Math.max(4, bannerStep * 1.05);
      const hitR2 = hitRadius * hitRadius;
      let removed = 0;
      const clearedLines = new Set<number>();
      for (const c of bannerChunks) {
        if (c.life <= 0) continue;

        const cx = c.x * bannerScaleX;
        const cy = c.y * bannerScaleY;
        const dx = bannerSnake.x - cx;
        const dy = bannerSnake.y - cy;

        if (dx * dx + dy * dy <= hitR2) {
          c.life = Math.max(0, c.life - dt * 5.2);
          if (c.life === 0) {
            bannerRemaining = Math.max(0, bannerRemaining - 1);
            if (c.line >= 0 && c.line < bannerLineAlive.length) {
              bannerLineAlive[c.line] = Math.max(0, bannerLineAlive[c.line] - 1);
              if (bannerLineAlive[c.line] === 0 && bannerLineTotals[c.line] > 0) {
                clearedLines.add(c.line);
              }
            }
            emitBannerParticles(cx, cy, 2);
            removed += 1;
          }
        }
      }

      if (clearedLines.size) {
        [...clearedLines]
          .sort((a, b) => b - a)
          .forEach((lineIdx) => collapseEatenLine(lineIdx));
      }

      if (removed > 0) playEatSound(removed);

      if (bannerRemaining <= 0) {
        bannerRemaining = 0;
        bannerWon = true;
        bannerGameActive = false;
        bannerParticles = [];
        if (!bannerWinSoundPlayed) {
          bannerWinSoundPlayed = true;
          playWinSound();
        }
      }
    }

    updateBannerParticles(dt);
  }

  function drawWinCelebration(ctx: CanvasRenderingContext2D) {
    const frame = WIN_ASCII_FRAMES[Math.floor((performance.now() / 170) % WIN_ASCII_FRAMES.length)];
    const line = frame[0] ?? "*** KAGHA CLEARED ***";
    const h = activeBannerHeight();
    const size = Math.max(8, Math.min(12, h * 0.62));
    const y = Math.max(0, h * 0.5 - size * 0.62);

    ctx.fillStyle = "rgba(15, 18, 28, 0.58)";
    ctx.fillRect(0, Math.max(0, y - 2), bannerViewW, size + 6);

    ctx.font = `700 ${size}px ui-monospace, Menlo, Consolas, monospace`;
    ctx.textBaseline = "top";
    ctx.textAlign = "center";
    ctx.fillStyle = "rgba(255, 226, 218, 0.94)";
    ctx.fillText(line, bannerViewW * 0.5, y);
  }

  function drawBannerGame() {
    if (!bannerCtx) return;
    const ctx = bannerCtx;

    ctx.clearRect(0, 0, bannerViewW, bannerViewH);
    if (!bannerShowCanvas) return;

    const clipH = activeBannerHeight();
    ctx.save();
    ctx.beginPath();
    ctx.rect(0, 0, bannerViewW, clipH);
    ctx.clip();

    for (const c of bannerChunks) {
      if (c.life <= 0) continue;
      const px = c.x * bannerScaleX;
      const py = c.y * bannerScaleY;
      const size = Math.max(1, Math.min(bannerScaleX, bannerScaleY) * bannerStep * (0.2 + c.life * 0.8));
      const alpha = 0.2 + c.life * 0.8;
      ctx.fillStyle = `rgba(246, 249, 255, ${alpha.toFixed(3)})`;
      ctx.fillRect(px - size / 2, py - size / 2, size, size);
    }

    if (bannerWon) {
      drawWinCelebration(ctx);
    }

    if (!bannerGameActive) {
      ctx.restore();
      return;
    }

    for (let i = bannerSnake.trail.length - 1; i >= 0; i -= 1) {
      const p = bannerSnake.trail[i];
      const t = 1 - i / Math.max(1, bannerSnake.trail.length - 1);
      const r = bannerSnake.radius * (0.16 + t * 0.84);
      const a = 0.1 + t * 0.82;
      const g = ctx.createRadialGradient(p.x - r * 0.3, p.y - r * 0.3, 1, p.x, p.y, r * 1.4);
      g.addColorStop(0, `rgba(255, 190, 182, ${a.toFixed(3)})`);
      g.addColorStop(1, `rgba(214, 74, 74, ${(a * 0.7).toFixed(3)})`);
      ctx.fillStyle = g;
      ctx.beginPath();
      ctx.arc(p.x, p.y, r, 0, Math.PI * 2);
      ctx.fill();
    }

    const hx0 = bannerSnake.x;
    const hy0 = bannerSnake.y;
    const hr = bannerSnake.radius + 1.4;
    const hg = ctx.createRadialGradient(hx0 - hr * 0.3, hy0 - hr * 0.3, 1, hx0, hy0, hr * 1.2);
    hg.addColorStop(0, "#ffd2c9");
    hg.addColorStop(1, "#e25a5a");
    ctx.fillStyle = hg;
    ctx.beginPath();
    ctx.arc(hx0, hy0, hr, 0, Math.PI * 2);
    ctx.fill();

    const n = normalize2d(bannerSnake.dirX, bannerSnake.dirY);
    const nx = -n.y;
    const ny = n.x;
    const fwd = Math.max(2.4, hr * 0.32);
    const spread = Math.max(2.2, hr * 0.26);
    ctx.fillStyle = "#07110e";
    ctx.beginPath();
    ctx.arc(hx0 + n.x * fwd + nx * spread, hy0 + n.y * fwd + ny * spread, 1.1, 0, Math.PI * 2);
    ctx.arc(hx0 + n.x * fwd - nx * spread, hy0 + n.y * fwd - ny * spread, 1.1, 0, Math.PI * 2);
    ctx.fill();

    for (const p of bannerParticles) {
      const alpha = clamp(p.life / 0.42, 0, 1);
      ctx.fillStyle = `rgba(255, 196, 174, ${alpha.toFixed(3)})`;
      ctx.beginPath();
      ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2);
      ctx.fill();
    }

    ctx.restore();
  }

  function tickBannerGame(ts: number) {
    if (!bannerGameActive && !bannerWon && bannerChunks.length === 0 && bannerParticles.length === 0) {
      if (bannerRaf) cancelAnimationFrame(bannerRaf);
      bannerRaf = 0;
      return;
    }
    const dt = clamp((ts - (bannerLastTs || ts)) / 1000, 0, 0.033);
    bannerLastTs = ts;
    if (bannerGameActive) {
      updateBannerGame(dt);
    }
    drawBannerGame();
    bannerRaf = requestAnimationFrame(tickBannerGame);
  }

  function handleBannerControlKey(e: KeyboardEvent, allowLetterControls: boolean) {
    if (!(props.showBanner ?? true)) return false;
    if (!bannerGameActive) return false;
    if (!bannerChunks.length) return false;
    if (e.altKey || e.ctrlKey || e.metaKey) return false;

    let dir: { x: number; y: number } | null = null;

    if (e.code === "ArrowUp") dir = { x: 0, y: -1 };
    else if (e.code === "ArrowLeft") dir = { x: -1, y: 0 };
    else if (e.code === "ArrowRight") dir = { x: 1, y: 0 };
    else if (e.code === "ArrowDown") dir = { x: 0, y: 1 };
    else if (allowLetterControls && e.code === "KeyZ") dir = { x: 0, y: -1 };
    else if (allowLetterControls && e.code === "KeyQ") dir = { x: -1, y: 0 };
    else if (allowLetterControls && e.code === "KeyD") dir = { x: 1, y: 0 };
    else if (allowLetterControls && e.code === "KeyS") dir = { x: 0, y: 1 };

     if (!dir) return false;
     ensureBannerAudioCtx();
     const n = normalize2d(dir.x, dir.y);
     bannerSnake.targetX = n.x;
     bannerSnake.targetY = n.y;
     e.preventDefault();
     focusInputSoon(true);
     return true;
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
    const restoreSavedBanner = () => {
      if (bannerGameActive || bannerWon || bannerHasScars || bannerTotal > 0) return;
      const saved = terminalBannerSessions.get(sessionId) as BannerRuntimeSnapshot | undefined;
      if (!saved) return;
      restoreBannerSession(saved);
    };
    const handleWindowKeys = (e: KeyboardEvent) => {
      if (!bannerGameActive) return;
      if (e.defaultPrevented) return;

      const target = e.target as HTMLElement | null;
      if (target) {
        const tag = target.tagName;
        if (tag === "INPUT" || tag === "TEXTAREA" || target.isContentEditable) return;
      }

      handleBannerControlKey(e, true);
    };
    const ro = new ResizeObserver(rerenderFit);
    if (asciiHost) ro.observe(asciiHost);

    // First render + after font settles.
    fitAscii();
    queueMicrotask(restoreSavedBanner);
    const t = setTimeout(() => {
      fitAscii();
      restoreSavedBanner();
    }, 180);
    window.addEventListener("resize", rerenderFit);
    window.addEventListener("keydown", handleWindowKeys);
    bannerRaf = requestAnimationFrame(tickBannerGame);

    return () => {
      ro.disconnect();
      window.removeEventListener("resize", rerenderFit);
      window.removeEventListener("keydown", handleWindowKeys);
      clearTimeout(t);
      cancelAnimationFrame(bannerRaf);
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

  const baseCmds = ["help", "clear", "ls", "exec", "split", "ps", "focus", "close", "dup", "ws", "game"] as const;

  function esc(s: string) {
    return s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");
  }

  const programAliases: Record<string, Program> = {
    terminal: "Terminal",
    termianl: "Terminal",
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

  function currentLangSuffix() {
    if (!browser) return "";
    const lang = (new URLSearchParams(window.location.search).get("lang") || "").toLowerCase();
    return lang.startsWith("fr") ? "?lang=fr" : "";
  }

  function kindFromCloseToken(raw: string): AppKind | undefined {
    const token = raw.trim();
    if (!token) return undefined;

    const prog = resolveProgram(token);
    if (prog) return programToKind[prog];

    const ws = workspaceForArg(token);
    if (ws) return ws.kind;

    const lowered = token.toLowerCase();
    switch (lowered) {
      case "terminal":
      case "term":
      case "termianl": return "terminal";
      case "about": return "about";
      case "experience": return "experience";
      case "work": return "work";
      case "skills": return "skills";
      case "contact": return "contact";
      default: return undefined;
    }
  }

  function closeAliasesForKind(kind: AppKind): string[] {
    switch (kind) {
      case "terminal": return ["Terminal"];
      case "about": return ["About", "AboutMe"];
      case "experience": return ["Experience"];
      case "work": return ["Work"];
      case "skills": return ["Skills"];
      case "contact": return ["Contact", "ContactMe"];
    }
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

  const helpSpecsEn: Record<string, HelpSpec> = {
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
      usage: "close [N|Name ...]",
      summary: "Close focused window or one/many targets by number or name.",
      examples: ["close", "close 2", "close AboutME", "close 2 Skills", "close --all Skills"]
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
    game: {
      usage: "game [start|pause|reset|exit]",
      summary: "Control the inline KAGHA snake game for this terminal.",
      examples: ["game", "game pause", "game reset", "game exit"]
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

  const helpSpecsFr: Record<string, HelpSpec> = {
    help: {
      usage: "help [command]",
      summary: "Affiche l'aide des commandes.",
      examples: ["help", "help close", "ls --help"]
    },
    clear: {
      usage: "clear",
      summary: "Nettoie la sortie du terminal (historique conservé).",
      examples: ["clear"]
    },
    ls: {
      usage: "ls [-a] [-l]",
      summary: "Liste les applications disponibles (format court ou long).",
      options: [
        { flag: "-a", desc: "Inclut les entrées . et ..." },
        { flag: "-l", desc: "Utilise le format d'affichage long." },
        { flag: "-h, --help", desc: "Affiche l'aide pour ls." }
      ],
      examples: ["ls", "ls -l", "ls -al"]
    },
    exec: {
      usage: "exec <Program>",
      summary: "Ouvre une application (split en WM, workspace hors WM).",
      examples: ["exec AboutMe", "exec Skills", "exec Terminal"]
    },
    split: {
      usage: "split <v|h> <Program>",
      summary: "Ouvre une application avec mode split force.",
      options: [
        { flag: "v", desc: "Split vertical (gauche/droite)." },
        { flag: "h", desc: "Split horizontal (haut/bas)." }
      ],
      examples: ["split v Work", "split h Terminal"]
    },
    ps: {
      usage: "ps",
      summary: "Liste les fenêtres ouvertes dans l'ordre.",
      examples: ["ps"]
    },
    focus: {
      usage: "focus <N>",
      summary: "Donne le focus à la fenêtre N.",
      examples: ["focus 2"]
    },
    close: {
      usage: "close [N|Name ...]",
      summary: "Ferme la fenêtre focus ou une/plusieurs cibles par numéro ou nom.",
      examples: ["close", "close 2", "close AboutME", "close 2 Skills", "close --all Skills"]
    },
    dup: {
      usage: "dup",
      summary: "Duplique la fenêtre focus.",
      examples: ["dup"]
    },
    ws: {
      usage: "ws <Name|N>",
      summary: "Navigue directement vers un workspace.",
      examples: ["ws About", "ws 3", "ws Skills"]
    },
    game: {
      usage: "game [start|pause|reset|exit]",
      summary: "Contrôle le mini-jeu snake KAGHA intégré à ce terminal.",
      examples: ["game", "game pause", "game reset", "game exit"]
    },
    "./X": {
      usage: "./<Program>",
      summary: "Alias de programme pour exec.",
      examples: ["./about", "./work", "./Terminal"]
    },
    Program: {
      usage: "<Program>",
      summary: "Ouvre directement une application par nom.",
      examples: ["AboutMe", "Experience", "Work", "Skills", "ContactMe"]
    }
  };

  function uiText() {
    if (isFr) {
      return {
        usage: "usage:",
        commands: "commandes:",
        options: "options:",
        examples: "exemples:",
        tip: "astuce:",
        noHelpTopic: "aucune aide pour:",
        unknownProgram: "programme inconnu",
        listAppsTip: "pour lister les applications disponibles",
        cmdHelpTip: "utilisez",
        cmdHelpTail: "pour les détails.",
        showHelpPages: "afficher les pages d'aide",
        clearOutput: "effacer la sortie",
        listApps: "lister les applications",
        openApp: "ouvrir une application",
        openSplit: "ouvrir une application en split",
        listWindows: "lister les fenêtres",
        focusWindow: "focus sur une fenêtre",
        closeWindows: "fermer une ou plusieurs fenêtres",
        duplicateFocused: "dupliquer la fenêtre focus",
        goWorkspace: "aller vers un workspace",
        snakeActions: "lancer/mettre en pause/réinitialiser/quitter le jeu snake",
        helpInFullTerminal: "l'aide est disponible dans le terminal complet",
        fullTerminalHelp: "utilisez le terminal complet pour l'aide des commandes",
        bannerRunning: "le jeu banner est déjà en cours",
        bannerNowRunning: "jeu banner en cours",
        bannerNotAvailable: "la bannière n'est pas disponible dans cette vue du terminal",
        bannerReset: "jeu banner réinitialisé",
        bannerPaused: "jeu banner en pause",
        bannerStopped: "le jeu banner est déjà arrêté",
        bannerExited: "jeu banner quitté ; ASCII original restauré",
        noWindows: "aucune fenêtre",
        duplicated: "duplique",
        pageHidden: "contenu de page masqué (Entrée x2 pour réinitialiser)",
        pageShown: "contenu de page affiché",
        layoutReset: "layout réinitialisé",
        closedFocused: "fenêtre focus fermée",
        closedCount: "fenêtre(s) fermée(s)",
        invalidCloseTarget: "cible close invalide:",
        focused: "focus sur",
        invalidFocus: "focus invalide",
        workspaceLabel: "workspace:",
        unknownWorkspace: "workspace inconnu",
        hint: "indice:",
        commandNotFound: "commande introuvable:"
      };
    }

    return {
      usage: "usage:",
      commands: "commands:",
      options: "options:",
      examples: "examples:",
      tip: "tip:",
      noHelpTopic: "no help topic for:",
      unknownProgram: "unknown program",
      listAppsTip: "to list available apps",
      cmdHelpTip: "use",
      cmdHelpTail: "for details.",
      showHelpPages: "show help pages",
      clearOutput: "clear output",
      listApps: "list apps",
      openApp: "open app",
      openSplit: "open app in split",
      listWindows: "list windows",
      focusWindow: "focus window",
      closeWindows: "close one or many windows",
      duplicateFocused: "duplicate focused window",
      goWorkspace: "go to workspace",
      snakeActions: "start/pause/reset/exit snake game",
      helpInFullTerminal: "help is available in the full terminal",
      fullTerminalHelp: "use full terminal for command help",
      bannerRunning: "banner game is already running",
      bannerNowRunning: "banner game running",
      bannerNotAvailable: "banner is not available in this terminal view",
      bannerReset: "banner game reset",
      bannerPaused: "banner game paused",
      bannerStopped: "banner game is already stopped",
      bannerExited: "banner game exited; original ASCII restored",
      noWindows: "no windows",
      duplicated: "duplicated",
      pageHidden: "page content hidden (double Enter to reset)",
      pageShown: "page content shown",
      layoutReset: "layout reset",
      closedFocused: "closed focused window",
      closedCount: "closed window(s)",
      invalidCloseTarget: "invalid close target:",
      focused: "focused",
      invalidFocus: "invalid focus",
      workspaceLabel: "workspace:",
      unknownWorkspace: "unknown workspace",
      hint: "hint:",
      commandNotFound: "command not found:"
    };
  }

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
          else return {
            ok: false as const,
            msg: isFr ? `ls: option non prise en charge -- '${f}'` : `ls: unsupported option -- '${f}'`
          };
        }
        continue;
      }

      if (token === "." || token === "~" || token === "/") continue;
      return {
        ok: false as const,
        msg: isFr
          ? `ls: impossible d'acceder a '${token}' : Aucun fichier ou dossier de ce type`
          : `ls: cannot access '${token}': No such file or directory`
      };
    }

    return { ok: true as const, showAll, longView };
  }

  function printUnknownProgram(target?: string) {
    const l = uiText();
    const prog = target?.trim();
    println(`<span class="err">${l.unknownProgram}${prog ? `: ${esc(prog)}` : ""}</span>`);
    println(`<span class="muted">${l.tip}</span> <span class="cmd">ls -al</span> <span class="muted">${l.listAppsTip}</span>`);
  }

  function escUsage(v: string) {
    return esc(v).replaceAll(" ", "&nbsp;");
  }

  type HelpKey = keyof typeof helpSpecsEn;

  function printHelp(name: HelpKey) {
    const l = uiText();
    const spec = (isFr ? helpSpecsFr : helpSpecsEn)[name];
    println(`<span class="muted">${l.usage}</span> <span class="cmd">${escUsage(spec.usage)}</span>`);
    println(`${esc(spec.summary)}`);

    if (spec.options?.length) {
      println(`<span class="muted">${l.options}</span>`);
      for (const opt of spec.options) {
        println(`<span class="cmd">${escUsage(opt.flag)}</span> <span class="muted">${esc(opt.desc)}</span>`);
      }
    }

    if (spec.examples?.length) {
      println(`<span class="muted">${l.examples}</span>`);
      for (const ex of spec.examples) {
        println(`<span class="cmd">${escUsage(ex)}</span>`);
      }
    }
  }

  function printGeneralHelp() {
    const l = uiText();
    println(`<span class="muted">${l.usage}</span> <span class="cmd">command [options]</span>`);
    println(`<span class="muted">${l.commands}</span>`);
    println(`<span class="cmd">help</span> <span class="muted">${l.showHelpPages}</span>`);
    println(`<span class="cmd">clear</span> <span class="muted">${l.clearOutput}</span>`);
    println(`<span class="cmd">ls</span> <span class="muted">${l.listApps}</span>`);
    println(`<span class="cmd">exec</span> <span class="muted">${l.openApp}</span>`);
    println(`<span class="cmd">split</span> <span class="muted">${l.openSplit}</span>`);
    println(`<span class="cmd">ps</span> <span class="muted">${l.listWindows}</span>`);
    println(`<span class="cmd">focus</span> <span class="muted">${l.focusWindow}</span>`);
    println(`<span class="cmd">close</span> <span class="muted">${l.closeWindows}</span>`);
    println(`<span class="cmd">dup</span> <span class="muted">${l.duplicateFocused}</span>`);
    println(`<span class="cmd">ws</span> <span class="muted">${l.goWorkspace}</span>`);
    println(`<span class="cmd">game</span> <span class="muted">${l.snakeActions}</span>`);
    println(`<span class="muted">${l.tip}</span> ${l.cmdHelpTip} <span class="cmd">&lt;command&gt; --help</span> ${l.cmdHelpTail}`);
  }

  function isHelpFlag(v: string) {
    return v === "-h" || v === "--help";
  }

  function normalizedHelpCommand(parts: string[]): HelpKey | undefined {
    if (!parts.length) return undefined;
    const cmd0 = parts[0];
    if (cmd0 in helpSpecsEn) return cmd0 as HelpKey;
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
        println(`<span class="err">${uiText().noHelpTopic} ${esc(parts[1])}</span>`);
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
    const snap = captureBannerSession();
    if (snap) terminalBannerSessions.set(sessionId, snap);
    else terminalBannerSessions.delete(sessionId);
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

  function focusInputSoon(onlyWhenActive = false) {
    const maybeFocus = () => {
      if (onlyWhenActive && !active) return;
      inputEl?.focus();
    };

    queueMicrotask(maybeFocus);

    if (!browser) return;
    requestAnimationFrame(maybeFocus);
    setTimeout(maybeFocus, 0);
  }

  $effect(() => {
    if (visible && active) focusInputSoon(true);
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
    if (wmMode) {
      const before = wmList().find((w) => w.focused)?.id;
      wmSplitOpen(kind, mode);
      if (before) wmFocus(before);
      focusInputSoon(true);
    }
    else {
      rememberWorkspaceSource();
      const path = routeForKind(kind);
      if (browser && window.location.pathname === path) return;
      void goto(`${path}${currentLangSuffix()}`, { keepFocus: true, noScroll: true, invalidateAll: false });
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

  function isCommandHead(head: string) {
    if (!head) return false;
    if ((baseCmds as readonly string[]).includes(head)) return true;
    if (head.startsWith("./")) return true;
    return resolveProgram(head) != null;
  }

  function withHelpFlags(cands: string[], tokenIndex: number, head: string) {
    if (tokenIndex < 1) return cands;
    if (!isCommandHead(head)) return cands;

    const out = [...cands];
    if (!out.includes("-h")) out.push("-h");
    if (!out.includes("--help")) out.push("--help");
    return out;
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
      return withHelpFlags([...programs], tokenIndex, head);
    }

    // split <v|h> <Program>
    if (head === "split") {
      if (tokenIndex === 1) return withHelpFlags(["v", "h"], tokenIndex, head);
      if (tokenIndex === 2) return withHelpFlags([...programs], tokenIndex, head);
      return withHelpFlags([], tokenIndex, head);
    }

    // focus <N>
    if (head === "focus") {
      const xs = wmList();
      return withHelpFlags(xs.map((x) => String(x.n)), tokenIndex, head);
    }

    // close [N|Name]
    if (head === "close") {
      const xs = wmList();
      const nums = xs.map((x) => String(x.n));
      const names = [...new Set([
        ...xs.flatMap((x) => [x.title, ...closeAliasesForKind(x.kind)])
      ])];
      return withHelpFlags([...nums, ...names, "--all"], tokenIndex, head);
    }

    // ws <name|N>
    if (head === "ws") {
      return withHelpFlags(["Terminal", "About", "Experience", "Work", "Skills", "Contact", "1", "2", "3", "4", "5", "6"], tokenIndex, head);
    }

    if (head === "ls") {
      return withHelpFlags(["-a", "-l", "-al", "-la", "-1"], tokenIndex, head);
    }

    if (head === "game") {
      if (tokenIndex === 1) return withHelpFlags(["start", "pause", "reset", "exit"], tokenIndex, head);
      return withHelpFlags([], tokenIndex, head);
    }

    if (head === "help") {
      return withHelpFlags([...baseCmds, ...programs], tokenIndex, head);
    }

    // allow program as second token in "./" not needed
    return withHelpFlags([], tokenIndex, head);
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
    println(`${promptHtml()}<span class="cmdText">&nbsp;&nbsp;${esc(c)}</span>`);
    if (!c) return;

    const l = uiText();

    history = [...history, c];
    historyIndex = history.length;

    if (maybeHandleHelp(c)) {
      return;
    }

    if (c === "help") {
      if (compact) {
        println(`<span class="muted">${l.helpInFullTerminal}</span>`);
        return;
      }
      printGeneralHelp();
      return;
    }

    if (c === "ls" || c.startsWith("ls ")) {
      const parsed = parseLsArgs(c);
      if (!parsed.ok) {
        if (parsed.msg === "__HELP__") {
          if (compact) println(`<span class="muted">${l.fullTerminalHelp}</span>`);
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

    if (c === "game" || c === "game start") {
      if (bannerGameActive && !bannerWon) {
        println(`<span class="muted">${l.bannerRunning}</span>`);
        return;
      }
      const started = startBannerGame(bannerWon);
      if (started) println(`<span class="muted">${l.bannerNowRunning}</span>`);
      else println(`<span class="err">${l.bannerNotAvailable}</span>`);
      return;
    }

    if (c === "game reset" || c === "game restart") {
      const started = startBannerGame(true);
      if (started) println(`<span class="muted">${l.bannerReset}</span>`);
      else println(`<span class="err">${l.bannerNotAvailable}</span>`);
      focusInputSoon(true);
      return;
    }

    if (c === "game pause" || c === "game stop") {
      const stopped = stopBannerGame();
      if (stopped) println(`<span class="muted">${l.bannerPaused}</span>`);
      else println(`<span class="muted">${l.bannerStopped}</span>`);
      return;
    }

    if (c === "game exit") {
      exitBannerGame();
      println(`<span class="muted">${l.bannerExited}</span>`);
      return;
    }

    if (c.startsWith("game ")) {
      println(`<span class="err">${l.usage}</span> <span class="cmd">game [start|pause|reset|exit]</span>`);
      return;
    }

    if (c === "ps") {
      const xs = wmList();
      const rows = xs
        .map((w) => `<div class="appRow"><span class="appIdx">${w.n}.</span>${appTokenHtml(w.kind, w.title)}${w.focused ? `<span class="focusMark">*</span>` : ""}<span class="muted">${w.kind}</span></div>`)
        .join("");
      println(rows || `<span class='muted'>${l.noWindows}</span>`);
      return;
    }

    if (c === "dup") {
      wmDuplicate(undefined, "auto");
      println(`<span class="muted">${l.duplicated}</span>`);
      return;
    }

    if (!wmMode && props.onLayoutAction && (c === "close" || c === "close page" || c === "close data")) {
      props.onLayoutAction("hide-content");
      println(`<span class="muted">${l.pageHidden}</span>`);
      return;
    }

    if (!wmMode && props.onLayoutAction && (c === "open" || c === "open page" || c === "show" || c === "show page")) {
      props.onLayoutAction("show-content");
      println(`<span class="muted">${l.pageShown}</span>`);
      return;
    }

    if (!wmMode && props.onLayoutAction && c === "reset") {
      props.onLayoutAction("reset-layout");
      out = [...starter];
      println(`<span class="muted">${l.layoutReset}</span>`);
      return;
    }

    if (c === "close") {
      wmClose();
      println(`<span class="muted">${l.closedFocused}</span>`);
      return;
    }

    if (c.startsWith("close ")) {
      const rawTargets = c.slice(6).split(/[\s,]+/).map((x) => x.trim()).filter(Boolean);
      const listed = wmList();
      const pool = [...listed];
      const ids: string[] = [];
      const invalid: string[] = [];

      const takeById = (id: string) => {
        const idx = pool.findIndex((w) => w.id === id);
        if (idx === -1) return false;
        pool.splice(idx, 1);
        ids.push(id);
        return true;
      };

      const takeByKindLatest = (kind: AppKind) => {
        for (let i = pool.length - 1; i >= 0; i -= 1) {
          const row = pool[i];
          if (row.kind === kind) {
            pool.splice(i, 1);
            ids.push(row.id);
            return true;
          }
        }
        return false;
      };

      const takeByNumber = (n: number) => {
        const row = pool.find((w) => w.n === n);
        if (!row) return false;
        return takeById(row.id);
      };

      const closeAllMode = rawTargets.length >= 1 && rawTargets[0] === "--all";

      if (closeAllMode) {
        const target = rawTargets[1];
        if (!target) {
          for (const row of [...pool].reverse()) {
            takeById(row.id);
          }
        } else {
          const kind = kindFromCloseToken(target);
          if (!kind) {
            invalid.push(target);
          } else {
            while (takeByKindLatest(kind)) {
              // Close all matching windows newest-first.
            }
            if (!ids.length) {
              invalid.push(target);
            }
          }
        }
      } else {
        for (const token of rawTargets) {
          const n = Number(token);
          if (Number.isFinite(n) && n > 0) {
            if (!takeByNumber(n)) invalid.push(token);
            continue;
          }

          const kind = kindFromCloseToken(token);
          if (!kind) {
            invalid.push(token);
            continue;
          }

          if (!takeByKindLatest(kind)) invalid.push(token);
        }
      }

      for (const id of ids) wmClose(id);

      if (ids.length) {
        const closeMsg = isFr
          ? `${ids.length} fenêtre${ids.length > 1 ? "s" : ""} fermée${ids.length > 1 ? "s" : ""}`
          : `closed ${ids.length} window${ids.length > 1 ? "s" : ""}`;
        println(`<span class="muted">${closeMsg}</span>`);
      }
      if (invalid.length) println(`<span class="err">${l.invalidCloseTarget} ${esc(invalid.join(", "))}</span>`);
      return;
    }

    if (c.startsWith("focus ")) {
      const n = Number(c.slice(6).trim());
      if (Number.isFinite(n) && n > 0) {
        wmFocusN(n);
        println(`<span class="muted">${l.focused} #${n}</span>`);
      } else {
        println(`<span class="err">${l.invalidFocus}</span>`);
      }
      return;
    }

    if (c.startsWith("ws ")) {
      const arg = c.slice(3).trim();
      const ws = workspaceForArg(arg);
      if (ws) {
        if (arg.toLowerCase() !== ws.name.toLowerCase() || arg !== ws.name) {
          println(`<span class="muted">${l.workspaceLabel}</span> <span class="cmd">ws ${ws.name}</span>`);
        }
        rememberWorkspaceSource();
        saveSession();
        if (browser && window.location.pathname !== ws.path) {
          void goto(`${ws.path}${currentLangSuffix()}`, { keepFocus: true, noScroll: true, invalidateAll: false });
        }
      }
      else println(`<span class="err">${l.unknownWorkspace}</span>`);
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
        if (wsCmd) println(`<span class="muted">${l.hint}</span> <span class="cmd">${esc(wsCmd)}</span>`);
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
        if (wsCmd) println(`<span class="muted">${l.hint}</span> <span class="cmd">${esc(wsCmd)}</span>`);
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
        if (wsCmd) println(`<span class="muted">${l.hint}</span> <span class="cmd">${esc(wsCmd)}</span>`);
      }
      open(programToKind[asProg], "auto");
      return;
    }

    println(`${l.commandNotFound} <span class="err">${esc(c)}</span>`);
  }

  function onKeydown(e: KeyboardEvent) {
    if (handleBannerControlKey(e, cmdline.trim().length === 0)) {
      queueMicrotask(() => inputEl?.setSelectionRange(cmdline.length, cmdline.length));
      return;
    }

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
          focusInputSoon(true);
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
      focusInputSoon(true);
      return;
    }

    if (e.key === "Tab") {
      e.preventDefault();
      applyCurrentMatch();
      focusInputSoon(true);
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

<div class={"term" + (compact ? " compact" : "") + ((props.embedded ?? false) ? " embedded" : "")} onpointerdown={() => inputEl?.focus()}>
  <div class="termBody" bind:this={bodyEl}>
    {#each out as row, i (i)}
      {#if row.kind === "banner"}
        <div class="line">
          <div class="banner">
            <div class="kaghaHero">
              <div
                class={"kaghaHit" + (bannerGameActive ? " playing" : "") + (bannerShowCanvas ? " scarred" : "")}
                bind:this={asciiHost}
                style={`--hx:${hx}; --hy:${hy}; --tx:${hx * 6}px; --ty:${hy * 4}px; --ascii-font:${asciiFontPx}px; --game-w:${Math.max(1, Math.ceil(bannerViewW))}px; --game-h:${Math.max(1, Math.ceil(bannerGameHeight))}px;`}
                onpointermove={onHeroMove}
                onpointerleave={onHeroLeave}
              >
                <div class="kaghaGameWrap">
                  <pre class="kaghaAscii base" bind:this={asciiBase}>{KAGHA_ASCII}</pre>
                  <pre class="kaghaAscii layer red" aria-hidden="true" style="--red-scale-y: {redScaleY};">{KAGHA_ASCII}</pre>
                  <pre class="kaghaAscii layer amber" aria-hidden="true" style="--amber-scale-y: {amberScaleY};">{KAGHA_ASCII}</pre>
                  <canvas class="kaghaGameCanvas" bind:this={bannerCanvas} aria-hidden="true"></canvas>
                  <div class="kaghaGameHud" aria-hidden="true">
                    <span>{bannerProgressPct}%</span>
                    {#if bannerWon}
                      <span class="kaghaWinTag">cleared</span>
                    {:else if bannerShowCanvas && !bannerGameActive}
                      <span class="kaghaPauseTag">paused</span>
                    {/if}
                  </div>
                </div>
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
          aria-label={isFr ? "Entree de commande du terminal" : "Terminal command input"}
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
