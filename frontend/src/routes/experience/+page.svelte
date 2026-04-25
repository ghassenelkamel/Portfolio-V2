<script lang="ts">
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

  const jobs: Job[] = [
    {
      key: "bookbeo",
      title: "Full Stack Mobile Developer",
      org: "bookBeo",
      dates: "Mar 2023 – Present",
      location: "Rennes, France · On-site",
      summary: [
        "React Native mobile development with a full-stack mindset.",
        "Focus on scalability, performance, and production reliability."
      ],
      bullets: ["React Native", "AWS exploration", "Team delivery", "Feature ownership"],
      img: "/assets/p7.jpg"
    },
    {
      key: "st",
      title: "Security Project Team Member",
      org: "STMicroelectronics",
      dates: "Sep 2022 – Mar 2023",
      location: "Le Mans, France · On-site",
      summary: [
        "Study + implementation around PSA Crypto API and embedded crypto constraints.",
        "Hands-on C development with disciplined version control and Linux tooling."
      ],
      bullets: ["C", "Crypto APIs", "ARM Cortex-M context", "Git/Linux workflows"],
      img: "/assets/p6.jpg"
    },
    {
      key: "falcon",
      title: "Mobile Developer Trainee",
      org: "Falcon Inspection & Services",
      dates: "Jun 2022 – Aug 2022",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "Built a mobile app to verify and record multiple industrial tests.",
        "Each test had its own attributes + validation math."
      ],
      bullets: ["Flutter", "Node.js backend", "Docker", "Automation"],
      img: "/assets/p5.jpg"
    },
    {
      key: "dall-2021",
      title: "Monitoring & Storage of Meteo Parameters",
      org: "DigiArtLivingLab",
      dates: "Jul 2021 – Aug 2021",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "Automated collection + storage pipeline for weather data.",
        "Built user interfaces and learned JSON parsing + API usage."
      ],
      bullets: ["Java", "JavaFX/Swing", "OpenWeather API", "Server-side automation"],
      img: "/assets/p1.jpg"
    },
    {
      key: "sagemcom",
      title: "Electronic Components Validation",
      org: "SAGEMCOM",
      dates: "Feb 2020 – May 2020",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "LabView tooling for selection/detection of components.",
        "Trained a CNN and deployed a lightweight variant on Raspberry Pi."
      ],
      bullets: ["LabView", "CNN", "TensorFlow Lite", "Raspberry Pi", "SolidWorks"],
      img: "/assets/p2.jpg",
      link: "https://drive.google.com/file/d/1HEhe-tPGshCTeTujwacxx_Wc-qHJIF3j/view"
    },
    {
      key: "dall-2019",
      title: "Humidity/Temp Station + Scent Diffuser + Apps",
      org: "DigiArtLivingLab",
      dates: "Jan 2019 – Feb 2019",
      location: "Tunis, Tunisia · On-site",
      summary: [
        "Arduino builds with Bluetooth data transmission.",
        "Android apps under time constraints (MIT App Inventor) + real-time storage."
      ],
      bullets: ["Arduino", "Bluetooth", "MIT App Inventor", "Firebase"],
      img: "/assets/p3.jpg"
    },
    {
      key: "dall-2018",
      title: "Educational Game (Directions)",
      org: "DigiArtLivingLab",
      dates: "Jan 2018 – Feb 2018",
      location: "Nabeul, Tunisia · On-site",
      summary: [
        "Unity game pipeline + assets, small C++ exploration.",
        "Team coordination and building engaging kid-friendly experience."
      ],
      bullets: ["Unity", "C++ intro", "Team work", "Level design"],
      img: "/assets/p4.jpg"
    }
  ];

  let selected = $state<Job>(jobs[0]);

  function pick(k: string) {
    const j = jobs.find((x) => x.key === k);
    if (j) selected = j;
  }
</script>

<svelte:head>
  <title>Experience</title>
</svelte:head>

<div class="grid">
  <aside class="left">
    <div class="hdr">
      <div class="k">btop --jobs</div>
      <div class="t">My Internships</div>
      <div class="s">Select a job to view details.</div>
    </div>

    <div class="list" role="list">
      {#each jobs as j (j.key)}
        <button
          type="button"
          class={"item " + (j.key === selected.key ? "active" : "")}
          onclick={() => pick(j.key)}
        >
          <div class="iTitle">{j.title}</div>
          <div class="iMeta">{j.org} · {j.dates}</div>
        </button>
      {/each}
    </div>
  </aside>

  <section class="right">
    <div class="topStats">
      <div class="stat">
        <div class="k">ORG</div>
        <div class="v">{selected.org}</div>
      </div>
      <div class="stat">
        <div class="k">WHEN</div>
        <div class="v">{selected.dates}</div>
      </div>
      <div class="stat">
        <div class="k">WHERE</div>
        <div class="v">{selected.location}</div>
      </div>
    </div>

    <div class="card">
      <div class="cardHead">
        <div class="title">
          <div class="big">{selected.title}</div>
          <div class="small"><span class="org">@{selected.org}</span> · {selected.location}</div>
        </div>

        {#if selected.link}
          <a class="link" href={selected.link} target="_blank" rel="noreferrer">open proof</a>
        {/if}
      </div>

      <div class="body">
        <div class="text">
          {#each selected.summary as p, i (i)}
            <p style={"--d:" + (i * 70) + "ms"}>{p}</p>
          {/each}

          {#if selected.bullets?.length}
            <div class="bul">
              {#each selected.bullets as b (b)}
                <span>{b}</span>
              {/each}
            </div>
          {/if}
        </div>

        {#if selected.img}
          <div class="thumb" aria-label="Preview">
            <img src={selected.img} alt={selected.title} />
          </div>
        {/if}
      </div>
    </div>
  </section>
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 14px;
    min-height: 520px;
  }

  .left, .right {
    border: 1px solid rgba(204,102,102,0.18);
    border-radius: 14px;
    background: rgba(12,14,18,0.14);
    overflow: hidden;
  }

  .hdr {
    padding: 14px;
    border-bottom: 1px solid rgba(204,102,102,0.18);
    background: rgba(204,102,102,0.06);
  }
  .k {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.55);
  }
  .t {
    font-family: 'HACKED', ui-monospace, Menlo, Consolas, monospace;
    margin-top: 6px;
    font-size: 26px;
    letter-spacing: 1px;
    color: rgba(204,102,102,0.95);
  }
  .s {
    margin-top: 6px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12.5px;
    color: rgba(245,245,245,0.72);
  }

  .list {
    padding: 10px;
    display: grid;
    gap: 8px;
    max-height: 520px;
    overflow: auto;
  }

  .list::-webkit-scrollbar { width: 10px; }
  .list::-webkit-scrollbar-thumb { background: rgba(204,102,102,0.22); border-radius: 999px; }

  .item {
    text-align: left;
    border: 1px solid rgba(204,102,102,0.16);
    border-radius: 12px;
    background: rgba(12,14,18,0.10);
    padding: 10px 12px;
    cursor: pointer;
    transition: transform 0.15s ease, background 0.15s ease, border-color 0.15s ease;
  }

  .item:hover {
    transform: translateY(-1px);
    background: rgba(204,102,102,0.10);
  }

  .active {
    border-color: rgba(228,90,90,0.45);
    background: rgba(228,90,90,0.10);
  }

  .iTitle {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.86);
    font-size: 13px;
  }
  .iMeta {
    margin-top: 4px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.58);
    font-size: 12px;
  }

  .right { padding: 14px; }

  .topStats {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 10px;
    margin-bottom: 12px;
  }

  .stat {
    border: 1px solid rgba(204,102,102,0.16);
    border-radius: 12px;
    background: rgba(12,14,18,0.10);
    padding: 10px 12px;
  }

  .stat .k {
    font-size: 11px;
    letter-spacing: 0.5px;
    color: rgba(228,90,90,0.78);
  }

  .stat .v {
    margin-top: 6px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 13px;
    color: rgba(245,245,245,0.82);
  }

  .card {
    border: 1px solid rgba(204,102,102,0.18);
    border-radius: 14px;
    background: rgba(12,14,18,0.12);
    overflow: hidden;
  }

  .cardHead {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;

    padding: 14px;
    border-bottom: 1px solid rgba(204,102,102,0.18);
    background: rgba(204,102,102,0.06);
  }

  .big {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 15px;
    color: rgba(245,245,245,0.90);
  }

  .small {
    margin-top: 6px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12.5px;
    color: rgba(245,245,245,0.62);
  }

  .org { color: rgba(228,90,90,0.88); }

  .link {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.82);
    text-decoration: none;
    border-bottom: 1px dashed rgba(228,90,90,0.35);
    height: fit-content;
    padding-top: 2px;
  }
  .link:hover { border-bottom-color: rgba(228,90,90,0.75); }

  .body {
    display: grid;
    grid-template-columns: 1fr 220px;
    gap: 14px;
    padding: 14px;
  }

  .text p {
    margin: 0 0 10px 0;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 13px;
    color: rgba(245,245,245,0.78);
    opacity: 0;
    transform: translateY(2px);
    animation: in 380ms ease forwards;
    animation-delay: var(--d);
  }

  @keyframes in {
    to { opacity: 1; transform: translateY(0); }
  }

  .bul {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 10px;
  }

  .bul span {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.75);
    border: 1px solid rgba(204,102,102,0.16);
    background: rgba(12,14,18,0.10);
    padding: 6px 8px;
    border-radius: 999px;
  }

  .thumb {
    border: 1px solid rgba(204,102,102,0.18);
    border-radius: 14px;
    overflow: hidden;
    background: rgba(12,14,18,0.16);
    height: 190px;
    align-self: start;
  }

  .thumb img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    filter: saturate(0.95) contrast(1.05);
    opacity: 0.9;
    transition: transform 0.2s ease, opacity 0.2s ease;
  }
  .thumb:hover img { transform: scale(1.03); opacity: 1; }

  @media (max-width: 900px) {
    .grid { grid-template-columns: 1fr; }
    .body { grid-template-columns: 1fr; }
    .thumb { height: 220px; }
  }
</style>
