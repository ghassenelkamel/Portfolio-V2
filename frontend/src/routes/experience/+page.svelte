<script lang="ts">
  import { fade } from "svelte/transition";
  import { page } from "$app/stores";

  type Job = {
    key: string;
    title: string;
    org: string;
    dates: string;
    location: string;
    summary: string[];
    bullets?: string[];
    img?: string;
    link?: string;
  };

  type ExperienceItem = {
    id: string;
    role: string;
    org: string;
    dates: string;
    location: string;
    summary: string[];
    bullets?: string[];
    image?: string;
    proofUrl?: string;
    sortOrder?: number;
  };

  type ExperienceContent = {
    eyebrow?: string;
    title?: string;
    description?: string;
    items?: ExperienceItem[];
  };

  const { data } = $props<{ data: { content?: ExperienceContent } }>();
  const content = $derived.by(() => data?.content ?? {});
  const isFr = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr"));

  const ui = $derived.by(() => isFr
    ? {
        title: "Experience",
        eyebrow: "Experience",
        fallbackTitle: "Parcours professionnel",
        fallbackDesc: "Sélectionnez une expérience pour explorer le périmètre, l'impact et les outils utilisés.",
        company: "Entreprise",
        period: "Période",
        location: "Lieu",
        highlights: "Points cles",
        openProof: "ouvrir la preuve",
        clickToExpand: "cliquer pour agrandir",
        openPreviewFor: "Ouvrir l'aperçu image pour",
        close: "fermer"
      }
    : {
        title: "Experience",
        eyebrow: "Experience",
        fallbackTitle: "Professional Journey",
        fallbackDesc: "Select an experience to explore scope, impact, and tools used.",
        company: "Company",
        period: "Period",
        location: "Location",
        highlights: "Highlights",
        openProof: "open proof",
        clickToExpand: "click to expand",
        openPreviewFor: "Open image preview for",
        close: "close"
      });

  const fallbackJobs: Job[] = [
    {
      key: "spade-integrity-2024",
      title: "IoT Engineer",
      org: "Spade Integrity",
      dates: "May 2024 - Present · 2 yrs 1 mo",
      location: "On-site",
      summary: [
        "Contribute to client platform development with focus on reliability, security, and operational continuity.",
        "Work on secure data recovery flows, NAS integration, and synchronization systems for alarms and critical data.",
        "Support end-to-end IoT/backend delivery by improving data pipelines, monitoring quality, and secure service integrations."
      ],
      bullets: ["Go", "IoT platform", "Secure data recovery", "NAS integration", "Alarm synchronization", "Monitoring", "Infrastructure"],
      img: "/assets/spade-integrity.svg"
    },
    {
      key: "bookbeo-2023",
      title: "Full Stack Mobile Developer",
      org: "bookBeo",
      dates: "Mar 2023 - Sep 2023 · 7 mos",
      location: "Rennes, Brittany, France · On-site",
      summary: [
        "Developed mobile product features and backend-connected workflows during internship delivery cycles.",
        "Collaborated on release quality and implementation speed with structured team practices and version control discipline."
      ],
      bullets: ["GitLab", "Git", "React Native", "API integration", "Product delivery", "Performance", "Team collaboration"],
      img: "/assets/p7.jpg"
    },
    {
      key: "st-2022",
      title: "Security Project Team Member",
      org: "STMicroelectronics",
      dates: "Sep 2022 – Mar 2023",
      location: "Le Mans, France · On-site",
      summary: [
        "Worked on PSA Crypto API implementation under embedded constraints and strict security requirements.",
        "Delivered C code with structured Git workflows and Linux-based validation tooling."
      ],
      bullets: ["C", "PSA Crypto", "Embedded security", "Linux tooling"],
      img: "/assets/p6.jpg"
    },
    {
      key: "falcon-2022",
      title: "Mobile Developer Experience",
      org: "Falcon Inspection & Services",
      dates: "Jun 2022 – Aug 2022",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "Built a mobile workflow to record and validate multiple industrial tests in the field.",
        "Implemented per-test data models and rule validation to improve consistency of reports."
      ],
      bullets: ["Flutter", "Node.js", "Validation flows", "Docker"],
      img: "/assets/p5.jpg"
    },
    {
      key: "dall-2021",
      title: "Weather Monitoring & Data Storage",
      org: "DigiArtLivingLab",
      dates: "Jul 2021 – Aug 2021",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "Automated weather data collection and storage pipelines for continuous monitoring.",
        "Built UI modules and API integrations with JSON parsing and backend persistence."
      ],
      bullets: ["Java", "JavaFX", "API integration", "Automation"],
      img: "/assets/p1.jpg"
    },
    {
      key: "sagemcom-2020",
      title: "Electronic Components Validation",
      org: "SAGEMCOM",
      dates: "Feb 2020 – May 2020",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "Built LabView workflows for detection and qualification of electronic components.",
        "Trained a CNN model and deployed a lightweight inference variant on Raspberry Pi hardware."
      ],
      bullets: ["LabView", "CNN", "TensorFlow Lite", "Raspberry Pi", "SolidWorks"],
      img: "/assets/p2.jpg",
      link: "https://drive.google.com/file/d/1HEhe-tPGshCTeTujwacxx_Wc-qHJIF3j/view"
    },
    {
      key: "dall-2019",
      title: "IoT Station + Mobile Apps",
      org: "DigiArtLivingLab",
      dates: "Jan 2019 – Feb 2019",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "Developed Arduino-based sensing builds with Bluetooth communication.",
        "Created Android companion apps and live data sync under tight timeline constraints."
      ],
      bullets: ["Arduino", "Bluetooth", "MIT App Inventor", "Firebase"],
      img: "/assets/p3.jpg"
    },
    {
      key: "dall-2018",
      title: "Educational Game Development",
      org: "DigiArtLivingLab",
      dates: "Jan 2018 – Feb 2018",
      location: "Nabeul, Tunisia · On-site",
      summary: [
        "Built a Unity educational game focused on directional learning for children.",
        "Coordinated assets, level pacing, and gameplay polish with a small collaborative team."
      ],
      bullets: ["Unity", "Gameplay design", "Team collaboration", "C++ basics"],
      img: "/assets/p4.jpg"
    }
  ];

  const jobs = $derived.by<Job[]>(() => {
    const items: ExperienceItem[] = Array.isArray(content.items) ? content.items : [];
    const jobsFromContent = items
      .slice()
      .sort((a, b) => (a.sortOrder ?? 0) - (b.sortOrder ?? 0))
      .map((item) => ({
        key: item.id,
        title: item.role,
        org: item.org,
        dates: item.dates,
        location: item.location,
        summary: Array.isArray(item.summary) ? item.summary : [],
        bullets: Array.isArray(item.bullets) ? item.bullets : [],
        img: item.image,
        link: item.proofUrl
      }))
      .filter((item) => item.key && item.title && item.org && item.dates && item.location && item.summary.length);

    if (!jobsFromContent.length) return fallbackJobs;

    const byKey = new Map<string, Job>(fallbackJobs.map((j) => [j.key, j]));
    for (const j of jobsFromContent) byKey.set(j.key, j);

    const merged = [...byKey.values()];

    const orderByFallback = new Map(fallbackJobs.map((j, i) => [j.key, i]));
    merged.sort((a, b) => {
      const ia = orderByFallback.get(a.key);
      const ib = orderByFallback.get(b.key);
      if (ia != null && ib != null) return ia - ib;
      if (ia != null) return -1;
      if (ib != null) return 1;
      return a.key.localeCompare(b.key);
    });

    return merged;
  });

  let selectedKey = $state<string>(fallbackJobs[0].key);

  $effect(() => {
    if (!jobs.length) return;
    if (!jobs.some((j) => j.key === selectedKey)) {
      selectedKey = jobs[0].key;
    }
  });

  const selected = $derived.by(() => {
    if (!jobs.length) return fallbackJobs[0];
    return jobs.find((j) => j.key === selectedKey) ?? jobs[0];
  });
  let modalOpen = $state(false);
  let modalSrc = $state("");
  let modalAlt = $state("");
  let modalCaption = $state("");

  let thumbRx = $state(0);
  let thumbRy = $state(0);
  let thumbTx = $state(0);
  let thumbTy = $state(0);
  let thumbLx = $state(50);
  let thumbLy = $state(50);

  function resetThumbMotion() {
    thumbRx = 0;
    thumbRy = 0;
    thumbTx = 0;
    thumbTy = 0;
    thumbLx = 50;
    thumbLy = 50;
  }

  function onThumbMove(e: PointerEvent) {
    if (typeof window !== "undefined" && window.matchMedia("(hover: none), (pointer: coarse)").matches) {
      return;
    }

    const el = e.currentTarget as HTMLElement;
    const r = el.getBoundingClientRect();
    if (r.width <= 0 || r.height <= 0) return;

    const px = (e.clientX - r.left) / r.width;
    const py = (e.clientY - r.top) / r.height;
    const cx = px - 0.5;
    const cy = py - 0.5;

    thumbRy = cx * 10;
    thumbRx = -cy * 8;
    thumbTx = cx * 8;
    thumbTy = cy * 6;
    thumbLx = Math.max(0, Math.min(100, px * 100));
    thumbLy = Math.max(0, Math.min(100, py * 100));
  }

  function onThumbLeave() {
    resetThumbMotion();
  }

  function openImageModal(job: Job) {
    if (!job.img) return;
    modalSrc = job.img;
    modalAlt = job.title;
    modalCaption = `${job.org} - ${job.title}`;
    modalOpen = true;
  }

  function closeImageModal() {
    modalOpen = false;
  }

  function onLightboxBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) closeImageModal();
  }

  function onWindowKeydown(e: KeyboardEvent) {
    if (e.key === "Escape" && modalOpen) closeImageModal();
  }

  function pick(k: string) {
    selectedKey = k;
    closeImageModal();
    resetThumbMotion();
  }
</script>

<svelte:head>
  <title>{ui.title}</title>
</svelte:head>

<svelte:window onkeydown={onWindowKeydown} />

<div class="grid">
  <aside class="left">
    <div class="hdr">
      <div class="eyebrow">{content.eyebrow || ui.eyebrow}</div>
      <div class="t">{content.title || ui.fallbackTitle}</div>
      <div class="s">{content.description || ui.fallbackDesc}</div>
    </div>

    <div class="list" role="list">
      {#each jobs as j, idx (j.key)}
        <button
          type="button"
          class={"item " + (j.key === selected.key ? "active" : "")}
          onclick={() => pick(j.key)}
          style={`--i:${idx};`}
        >
          <div class="iTitle">{j.title}</div>
          <div class="iMeta">{j.org}</div>
          <div class="iSub">{j.dates} · {j.location}</div>
        </button>
      {/each}
    </div>
  </aside>

  <section class="right">
    <div class="topStats">
      <div class="stat">
        <div class="k">{ui.company}</div>
        <div class="v">{selected.org}</div>
      </div>
      <div class="stat">
        <div class="k">{ui.period}</div>
        <div class="v">{selected.dates}</div>
      </div>
      <div class="stat">
        <div class="k">{ui.location}</div>
        <div class="v">{selected.location}</div>
      </div>
    </div>

    {#key selected.key}
      <article class="card">
        <div class="cardGlow" aria-hidden="true"></div>

        <div class="cardHead">
          <div class="title">
            <div class="big">{selected.title}</div>
            <div class="small"><span class="org">@{selected.org}</span> · {selected.location}</div>
          </div>

          {#if selected.link}
            <a class="link" href={selected.link} target="_blank" rel="noreferrer">{ui.openProof}</a>
          {/if}
        </div>

        <div class="body">
          <div class="text">
            <div class="label">{ui.highlights}</div>

            {#each selected.summary as p, i (i)}
              <p>{p}</p>
            {/each}

            {#if selected.bullets?.length}
              <div class="bul">
                {#each selected.bullets as b, i (b)}
                  <span>{b}</span>
                {/each}
              </div>
            {/if}
          </div>

          {#if selected.img}
            <figure class="thumb" aria-label="Experience preview">
              <button
                type="button"
                class="thumbBtn"
                aria-label={`${ui.openPreviewFor} ${selected.title}`}
                title={ui.clickToExpand}
                onclick={() => openImageModal(selected)}
                onpointermove={onThumbMove}
                onpointerleave={onThumbLeave}
                style={`--rx:${thumbRx}deg; --ry:${thumbRy}deg; --tx:${thumbTx}px; --ty:${thumbTy}px; --lx:${thumbLx}%; --ly:${thumbLy}%;`}
              >
                <img class="thumbImg" src={selected.img} alt={selected.title} loading="eager" decoding="async" />
                <span class="zoomHint">{ui.clickToExpand}</span>
              </button>
              <figcaption>{selected.org}</figcaption>
            </figure>
          {/if}
        </div>
      </article>
    {/key}
  </section>
</div>

{#if modalOpen}
  <div
    class="lightbox"
    role="dialog"
    aria-modal="true"
    aria-label={modalCaption}
    tabindex="0"
    onclick={onLightboxBackdrop}
    onkeydown={(e) => {
      if (e.key === "Escape") closeImageModal();
    }}
    in:fade={{ duration: 150 }}
    out:fade={{ duration: 100 }}
  >
    <div class="lightboxInner">
      <button type="button" class="lightboxClose" onclick={closeImageModal}>{ui.close}</button>
      <img class="lightboxImg" src={modalSrc} alt={modalAlt} loading="eager" decoding="async" />
      <div class="lightboxCap">{modalCaption}</div>
    </div>
  </div>
{/if}

<style>
  .grid {
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 14px;
    min-height: 560px;
  }

  .left,
  .right {
    border: 1px solid rgba(204, 102, 102, 0.18);
    border-radius: 14px;
    background: rgba(12, 14, 18, 0.16);
    backdrop-filter: blur(22px) saturate(140%);
    -webkit-backdrop-filter: blur(22px) saturate(140%);
    overflow: hidden;
  }

  .hdr {
    padding: 14px;
    border-bottom: 1px solid rgba(204, 102, 102, 0.18);
    background:
      radial-gradient(180px 90px at 10% -20%, rgba(228, 90, 90, 0.18), transparent 70%),
      rgba(204, 102, 102, 0.06);
  }

  .eyebrow {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: rgba(245, 245, 245, 0.58);
  }

  .t {
    margin-top: 7px;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 24px;
    line-height: 1.2;
    letter-spacing: 0.4px;
    color: rgba(228, 90, 90, 0.94);
  }

  .s {
    margin-top: 8px;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12.5px;
    line-height: 1.5;
    color: rgba(245, 245, 245, 0.74);
    max-width: 34ch;
  }

  .list {
    padding: 10px;
    display: grid;
    gap: 8px;
    max-height: 560px;
    overflow: auto;
  }

  .list::-webkit-scrollbar {
    width: 10px;
  }

  .list::-webkit-scrollbar-thumb {
    background: rgba(204, 102, 102, 0.24);
    border-radius: 999px;
  }

  .item {
    text-align: left;
    border: 1px solid rgba(204, 102, 102, 0.16);
    border-radius: 12px;
    background: linear-gradient(180deg, rgba(20, 24, 30, 0.45), rgba(12, 14, 18, 0.22));
    padding: 11px 12px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    opacity: 0;
    animation: itemIn 380ms ease forwards;
    animation-delay: calc(var(--i) * 42ms);
    transition: transform 0.18s ease, background 0.18s ease, border-color 0.18s ease;
  }

  .item::before {
    content: "";
    position: absolute;
    left: 0;
    top: 10px;
    bottom: 10px;
    width: 2px;
    background: rgba(228, 90, 90, 0);
    transition: background 0.18s ease;
  }

  .item:hover {
    transform: translateY(-1px);
    background: linear-gradient(180deg, rgba(228, 90, 90, 0.14), rgba(12, 14, 18, 0.24));
    border-color: rgba(228, 90, 90, 0.3);
  }

  .active {
    border-color: rgba(228, 90, 90, 0.44);
    background: linear-gradient(180deg, rgba(228, 90, 90, 0.16), rgba(12, 14, 18, 0.24));
    box-shadow: 0 0 0 1px rgba(228, 90, 90, 0.16) inset;
  }

  .active::before {
    background: rgba(228, 90, 90, 0.9);
  }

  .iTitle {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    color: rgba(245, 245, 245, 0.9);
    font-size: 13.2px;
    line-height: 1.35;
  }

  .iMeta,
  .iSub {
    margin-top: 4px;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
  }

  .iMeta {
    color: rgba(245, 245, 245, 0.68);
  }

  .iSub {
    color: rgba(245, 245, 245, 0.52);
  }

  .right {
    padding: 14px;
    display: grid;
    grid-template-rows: auto 1fr;
    gap: 12px;
  }

  .topStats {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 10px;
  }

  .stat {
    border: 1px solid rgba(204, 102, 102, 0.16);
    border-radius: 12px;
    background: rgba(12, 14, 18, 0.14);
    padding: 10px 12px;
  }

  .stat .k {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 11px;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: rgba(228, 90, 90, 0.82);
  }

  .stat .v {
    margin-top: 6px;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 13px;
    color: rgba(245, 245, 245, 0.84);
    line-height: 1.35;
  }

  .card {
    position: relative;
    border: 1px solid rgba(204, 102, 102, 0.22);
    border-radius: 14px;
    background: linear-gradient(180deg, rgba(16, 20, 24, 0.5), rgba(12, 14, 18, 0.26));
    overflow: hidden;
    min-height: 0;
  }

  .cardGlow {
    position: absolute;
    inset: -120px auto auto -40px;
    width: 260px;
    height: 260px;
    border-radius: 999px;
    background: radial-gradient(circle, rgba(228, 90, 90, 0.22), rgba(228, 90, 90, 0));
    pointer-events: none;
    animation: drift 8s ease-in-out infinite;
  }

  .cardHead {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
    padding: 14px;
    border-bottom: 1px solid rgba(204, 102, 102, 0.2);
    background: rgba(228, 90, 90, 0.06);
  }

  .big {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 15.5px;
    color: rgba(245, 245, 245, 0.94);
    line-height: 1.3;
  }

  .small {
    margin-top: 6px;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12.5px;
    color: rgba(245, 245, 245, 0.64);
  }

  .org {
    color: rgba(228, 90, 90, 0.9);
  }

  .link {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    color: rgba(245, 245, 245, 0.86);
    text-decoration: none;
    border: 1px solid rgba(228, 90, 90, 0.28);
    background: rgba(228, 90, 90, 0.1);
    border-radius: 999px;
    padding: 7px 10px;
    height: fit-content;
    transition: transform 0.15s ease, background 0.15s ease;
  }

  .link:hover {
    transform: translateY(-1px);
    background: rgba(228, 90, 90, 0.2);
  }

  .body {
    position: relative;
    z-index: 1;
    display: grid;
    grid-template-columns: minmax(0, 1fr) 340px;
    gap: 14px;
    padding: 14px;
  }

  .label {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    letter-spacing: 0.6px;
    text-transform: uppercase;
    color: rgba(228, 90, 90, 0.86);
    margin-bottom: 10px;
  }

  .text p {
    margin: 0 0 10px;
    font-family: var(--font-sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell,
      "Open Sans", "Helvetica Neue", sans-serif);
    font-size: 14px;
    line-height: 1.62;
    color: rgba(245, 245, 245, 0.9);
    max-width: 68ch;
  }

  .bul {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 12px;
  }

  .bul span {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    color: rgba(245, 245, 245, 0.84);
    border: 1px solid rgba(228, 90, 90, 0.26);
    background: rgba(228, 90, 90, 0.1);
    padding: 6px 10px;
    border-radius: 999px;
    transition: transform 0.15s ease, background 0.15s ease;
  }

  .bul span:hover {
    transform: translateY(-1px);
    background: rgba(228, 90, 90, 0.17);
  }

  .thumb {
    margin: 0;
    border: 1px solid rgba(228, 90, 90, 0.28);
    border-radius: 14px;
    overflow: hidden;
    background: rgba(12, 14, 18, 0.26);
    min-height: 300px;
    align-self: start;
    position: relative;
    box-shadow: 0 14px 32px rgba(0, 0, 0, 0.35);
    perspective: 1100px;
  }

  .thumbBtn {
    border: 0;
    padding: 0;
    margin: 0;
    background: transparent;
    width: 100%;
    height: 100%;
    min-height: 300px;
    position: relative;
    cursor: zoom-in;
    overflow: hidden;
    display: block;
  }

  .thumbBtn::before {
    content: "";
    position: absolute;
    inset: 0;
    z-index: 2;
    pointer-events: none;
    background:
      linear-gradient(to top, rgba(6, 9, 12, 0.5), rgba(6, 9, 12, 0.12) 45%, rgba(6, 9, 12, 0.02)),
      radial-gradient(130px 90px at 78% 10%, rgba(228, 90, 90, 0.2), rgba(228, 90, 90, 0));
  }

  .thumbBtn::after {
    content: "";
    position: absolute;
    inset: 0;
    z-index: 3;
    pointer-events: none;
    background: radial-gradient(180px 120px at var(--lx) var(--ly), rgba(255, 255, 255, 0.14), rgba(255, 255, 255, 0));
    opacity: 0.8;
    transition: opacity 0.2s ease;
  }

  .thumbImg {
    width: 100%;
    height: 100%;
    min-height: 300px;
    object-fit: cover;
    object-position: center;
    filter: saturate(1.14) contrast(1.14) brightness(1.06);
    transform: translate3d(var(--tx), var(--ty), 0) scale(1.05) rotateX(var(--rx)) rotateY(var(--ry));
    transition: transform 0.16s linear, filter 0.22s ease;
  }

  .thumbBtn:hover .thumbImg {
    filter: saturate(1.22) contrast(1.16) brightness(1.08);
  }

  .zoomHint {
    position: absolute;
    right: 10px;
    top: 10px;
    z-index: 4;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 11px;
    color: rgba(245, 245, 245, 0.9);
    background: rgba(10, 12, 16, 0.54);
    border: 1px solid rgba(228, 90, 90, 0.3);
    border-radius: 999px;
    padding: 5px 8px;
    opacity: 0;
    transform: translateY(-2px);
    transition: opacity 0.2s ease, transform 0.2s ease;
  }

  .thumbBtn:hover .zoomHint {
    opacity: 1;
    transform: translateY(0);
  }

  figcaption {
    position: absolute;
    z-index: 5;
    left: 10px;
    bottom: 10px;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 11.5px;
    letter-spacing: 0.4px;
    color: rgba(245, 245, 245, 0.9);
    background: rgba(6, 9, 12, 0.54);
    border: 1px solid rgba(228, 90, 90, 0.26);
    border-radius: 999px;
    padding: 5px 8px;
  }

  .lightbox {
    position: fixed;
    inset: 0;
    z-index: 140;
    background: rgba(3, 5, 8, 0.82);
    backdrop-filter: blur(7px);
    -webkit-backdrop-filter: blur(7px);
    display: grid;
    place-items: center;
    padding: 20px;
  }

  .lightboxInner {
    width: min(1040px, 100%);
    max-height: min(90vh, 900px);
    display: grid;
    grid-template-rows: auto minmax(0, 1fr) auto;
    gap: 10px;
    border: 1px solid rgba(228, 90, 90, 0.35);
    border-radius: 16px;
    background: linear-gradient(180deg, rgba(15, 18, 24, 0.88), rgba(8, 10, 14, 0.92));
    box-shadow: 0 26px 64px rgba(0, 0, 0, 0.5);
    overflow: hidden;
  }

  .lightboxClose {
    justify-self: end;
    margin: 10px 10px 0 0;
    border: 1px solid rgba(228, 90, 90, 0.32);
    background: rgba(228, 90, 90, 0.12);
    color: rgba(245, 245, 245, 0.92);
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    border-radius: 999px;
    padding: 7px 11px;
    cursor: pointer;
  }

  .lightboxClose:hover {
    background: rgba(228, 90, 90, 0.22);
  }

  .lightboxImg {
    width: 100%;
    height: 100%;
    max-height: min(74vh, 760px);
    object-fit: contain;
    object-position: center;
    padding: 0 14px;
    box-sizing: border-box;
    filter: saturate(1.16) contrast(1.12) brightness(1.05);
  }

  .lightboxCap {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    color: rgba(245, 245, 245, 0.82);
    padding: 0 14px 12px;
  }

  @keyframes itemIn {
    from {
      opacity: 0;
      transform: translateY(6px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes drift {
    0%,
    100% {
      transform: translate3d(0, 0, 0);
    }
    50% {
      transform: translate3d(18px, 8px, 0);
    }
  }

  @media (max-width: 980px) {
    .grid {
      grid-template-columns: 1fr;
      min-height: auto;
    }

    .list {
      max-height: 340px;
    }

    .topStats {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .cardHead {
      flex-direction: column;
      align-items: flex-start;
    }

    .link {
      align-self: flex-start;
    }

    .body {
      grid-template-columns: 1fr;
    }

    .text p {
      max-width: none;
    }

    .thumb {
      min-height: 260px;
    }

    .thumbBtn {
      min-height: 260px;
    }

    .thumbImg {
      min-height: 260px;
    }

    .lightbox {
      padding: 10px;
    }

    .lightboxInner {
      width: 100%;
      max-height: 92vh;
    }
  }

  @media (max-width: 760px) {
    .left,
    .right,
    .card,
    .stat,
    .item,
    .thumb {
      border-radius: 12px;
    }

    .right {
      padding: 10px;
      gap: 10px;
    }

    .topStats {
      grid-template-columns: 1fr;
      gap: 8px;
    }

    .cardHead,
    .body {
      padding: 12px;
    }

    .big {
      font-size: 14px;
    }

    .small,
    .stat .v {
      font-size: 12px;
    }

    .thumb,
    .thumbBtn,
    .thumbImg {
      min-height: 220px;
    }

    .lightboxClose {
      padding: 8px 12px;
    }

    .lightboxCap {
      font-size: 11px;
    }
  }

  @media (hover: none), (pointer: coarse) {
    .thumbImg {
      transform: none;
      transition: filter 0.22s ease;
    }

    .zoomHint {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .item,
    .cardGlow,
    .thumbImg,
    .thumbBtn::after,
    .link,
    .bul span {
      animation: none;
      transition: none;
    }
  }
</style>
