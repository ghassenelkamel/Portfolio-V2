<script lang="ts">
  import { onMount } from "svelte";

  const props = $props<{
    profile: {
      name: string;
      tagline: string;
      about: string;
      working: string;
      goals: string;
      email: string;
      portfolio: string;
      github: string;
      skills: string[];
    };
  }>();

  const P = "[kagha@archbios about]$ ";

  let out = $state<string[]>([]);
  let cmd = $state("");
  let inputEl = $state<HTMLInputElement | null>(null);

  function esc(s: string){
    return s.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;");
  }

  function println(s: string){
    out = [...out, s];
  }

  async function typeBlock(lines: string[], baseDelay = 70){
    for (let i = 0; i < lines.length; i++){
      await new Promise(r => setTimeout(r, baseDelay));
      println(esc(lines[i]));
    }
  }

  function toLines(block: string){
    if (!block) return [];
    return block
      .split("\n")
      .map(l => l.trimEnd())
      .filter(l => l.trim().length > 0)
      // remove markdown bullets formatting a bit
      .map(l => l.replace(/^\-\s+/, "• "));
  }

  function help(){
    return [
      "help                         # show commands",
      "ls                           # list available docs",
      "cat about | now | goals      # print sections (from GitHub README)",
      "cat skills                   # show skill set",
      "open github|portfolio|email  # open link",
      "clear                        # clear screen",
      "refresh                      # reload (pull newest README)"
    ];
  }

  function ls(){
    return [
      "about.md",
      "now.md",
      "goals.md",
      "skills.md",
      "contact.md"
    ];
  }

  async function run(raw: string){
    const c = raw.trim();
    println(esc(P + c));

    if (!c) return;

    if (c === "clear"){
      out = [];
      return;
    }

    if (c === "help"){
      await typeBlock(help(), 45);
      return;
    }

    if (c === "ls"){
      await typeBlock(ls().map(x => "  " + x), 35);
      return;
    }

    if (c === "refresh"){
      location.reload();
      return;
    }

    if (c.startsWith("open ")){
      const x = c.slice(5).trim();
      if (x === "github") window.open(props.profile.github, "_blank", "noreferrer");
      else if (x === "portfolio") window.open(props.profile.portfolio, "_blank", "noreferrer");
      else if (x === "email") window.location.href = `mailto:${props.profile.email}`;
      else await typeBlock([`Unknown target: ${x}`], 40);
      return;
    }

    if (c.startsWith("cat ")){
      const what = c.slice(4).trim();
      if (what === "about"){
        await typeBlock([`# ${props.profile.name}`, props.profile.tagline, ""], 40);
        await typeBlock(toLines(props.profile.about), 35);
        return;
      }
      if (what === "now"){
        await typeBlock(["# NOW", ""], 40);
        await typeBlock(toLines(props.profile.working), 35);
        return;
      }
      if (what === "goals"){
        await typeBlock(["# GOALS", ""], 40);
        await typeBlock(toLines(props.profile.goals), 35);
        return;
      }
      if (what === "skills"){
        await typeBlock(["# SKILLS", ""], 40);
        const skills = (props.profile.skills?.length ? props.profile.skills : [
          "Go","Linux","Svelte","MQTT","OpenVPN","WireGuard","Docker","PostgreSQL"
        ]);
        await typeBlock(skills.map((s: string) => "• " + s), 20);
        return;
      }
      if (what === "contact"){
        await typeBlock(["# CONTACT",""], 40);
        await typeBlock([
          `Email: ${props.profile.email}`,
          `Portfolio: ${props.profile.portfolio}`,
          `GitHub: ${props.profile.github}`
        ], 35);
        return;
      }

      await typeBlock([`Unknown doc: ${what}`], 40);
      return;
    }

    await typeBlock([`Unknown command: ${c}`, `Type "help"`], 40);
  }

  function onKeyDown(e: KeyboardEvent){
    if (e.key === "Enter"){
      e.preventDefault();
      const raw = cmd;
      cmd = "";
      void run(raw);
      queueMicrotask(() => inputEl?.focus());
    }
    if (e.key === "Tab"){
      // tiny “smart hint” (not full completion, but feels terminal)
      e.preventDefault();
      if (!cmd) cmd = "help";
      else if (cmd === "cat") cmd = "cat about";
      else if (cmd === "open") cmd = "open github";
    }
  }

  onMount(async () => {
    // Boot sequence: feels like a terminal
    await typeBlock([
      "cat ~/about.md",
      ""
    ], 35);
    await typeBlock([
      props.profile.name,
      props.profile.tagline,
      "",
      'Type "help" then explore with "ls" and "cat about/now/goals".'
    ], 40);
    inputEl?.focus();
  });
</script>

<div class="aboutTerm">
  <div class="screen">
    {#each out as l, i (i)}
      <div class="l">{@html l}</div>
    {/each}
  </div>

  <div class="in">
    <span class="p">{P}</span>
    <input bind:this={inputEl} value={cmd} oninput={(e)=>cmd=(e.target as HTMLInputElement).value} onkeydown={onKeyDown} autocomplete="off" spellcheck="false" />
  </div>

  <div class="mini">
    <span>Hints:</span>
    <span class="k">help</span>
    <span class="k">ls</span>
    <span class="k">cat about</span>
    <span class="k">cat now</span>
    <span class="k">open portfolio</span>
  </div>
</div>

<style>
  .aboutTerm{
    border: 1px solid var(--border, rgba(255, 255, 255, 0.10));
    border-radius: var(--radius, 14px);

    background: linear-gradient(
      180deg,
      var(--term-surface, rgba(15, 20, 25, 0.16)),
      var(--term-surface-2, rgba(15, 20, 25, 0.09))
    );

    backdrop-filter: blur(var(--term-blur, 14px)) saturate(145%);
    -webkit-backdrop-filter: blur(var(--term-blur, 14px)) saturate(145%);

    padding: 14px;
  }

  .screen{
    min-height: 420px;
    max-height: 540px;
    overflow: auto;
    padding: 6px 4px 10px 4px;
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 13.5px;
    color: rgba(245,245,245,0.86);
  }

  .screen::-webkit-scrollbar{ width: 10px; }
  .screen::-webkit-scrollbar-thumb{ background: rgba(204,102,102,0.22); border-radius: 999px; }

  .l{ margin-bottom: 8px; white-space: pre-wrap; }

  .in{
    display: flex;
    gap: 10px;
    align-items: center;
    padding-top: 10px;
    border-top: 1px solid rgba(204,102,102,0.16);
    margin-top: 10px;

    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 13.5px;
  }

  .p{ color: var(--accent, rgba(228, 90, 90, 0.92)); }
  input{
    width: 100%;
    border: 0;
    outline: none;
    background: transparent;
    color: rgba(245,245,245,0.90);
    font-family: inherit;
    font-size: inherit;
    caret-color: var(--accent, rgba(228, 90, 90, 0.90));
  }

  .mini{
    margin-top: 10px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;

    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    font-size: 12px;
    color: rgba(245,245,245,0.55);
  }
  .k{
    padding: 4px 8px;
    border-radius: 999px;
    border: 1px solid rgba(228,90,90,0.28);
    background: rgba(228,90,90,0.06);
    color: rgba(228,90,90,0.88);
  }
</style>
