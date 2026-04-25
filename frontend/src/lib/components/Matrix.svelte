<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { browser } from "$app/environment";

  let canvas: HTMLCanvasElement | null = null;
  let ctx: CanvasRenderingContext2D | null = null;

  let interval: number | null = null;

  let w = 0;
  let h = 0;
  let dpr = 1;

  // Matrix-ish charset: Katakana + numerics + symbols
  const KATAKANA =
    "アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン";
  const EXTRA = "0123456789@#$%&*+-/<>{}[]";
  const CHARS = (KATAKANA + EXTRA).split("");

  let fontSize = 14;
  let columns = 0;
  let drops: { y: number; speed: number }[] = [];

  function rand(min: number, max: number) {
    return Math.random() * (max - min) + min;
  }

  function resize() {
    if (!browser || !canvas) return;

    dpr = window.devicePixelRatio || 1;
    w = window.innerWidth;
    h = window.innerHeight;

    canvas.width = Math.floor(w * dpr);
    canvas.height = Math.floor(h * dpr);
    canvas.style.width = `${w}px`;
    canvas.style.height = `${h}px`;

    ctx = canvas.getContext("2d");
    if (!ctx) return;
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    fontSize = w < 700 ? 12 : 14;
    columns = Math.floor(w / fontSize);
    drops = Array.from({ length: columns }, () => ({
      y: rand(-h, 0),
      speed: rand(0.8, 2.6)
    }));
  }

function step() {
     if (!browser || !ctx) return;

     // Dark trail (slight red tint)
     ctx.fillStyle = "rgba(10, 0, 0, 0.07)";
     ctx.fillRect(0, 0, w, h);

     ctx.font = `${fontSize}px 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Menlo, monospace`;

     for (let i = 0; i < drops.length; i++) {
       const x = i * fontSize;
       const y = drops[i].y;

       // Matrix color: RED (requested)
       const head = "rgba(255, 80, 80, 0.92)";
       const tail = "rgba(255, 80, 80, 0.36)";

       const ch = CHARS[(Math.random() * CHARS.length) | 0];

       // Subtle glow effect
       ctx.shadowColor = "rgba(255, 80, 80, 0.22)";
       ctx.shadowBlur = 6;

       ctx.fillStyle = head;
       ctx.fillText(ch, x, y);

       ctx.shadowBlur = 0;
       ctx.fillStyle = tail;
       ctx.fillText(CHARS[(Math.random() * CHARS.length) | 0], x, y - fontSize * 2);

       drops[i].y += drops[i].speed * fontSize;

       if (drops[i].y > h + rand(0, 200) && Math.random() > 0.975) {
         drops[i].y = rand(-200, 0);
         drops[i].speed = rand(0.8, 2.6);
       }
     }
   }

  onMount(() => {
    if (!browser) return;
    resize();
    window.addEventListener("resize", resize);
    interval = window.setInterval(step, 35);
  });

  onDestroy(() => {
    if (!browser) return;
    window.removeEventListener("resize", resize);
    if (interval !== null) window.clearInterval(interval);
    interval = null;
  });
</script>

<canvas bind:this={canvas} class="matrix" aria-hidden="true"></canvas>

<style>
  .matrix {
    position: fixed;
    inset: 0;
    z-index: 0;
    background: #000;
    pointer-events: none;
  }
</style>
