<script lang="ts">
  type Item = {
    id: string;
    company: string;
    role: string;
    where: string;
    dates: string;
    impact: string[];
    stack: string[];
    image?: string;
    proofUrl?: string;
    sortOrder?: number;
  };

  type ExperienceContent = {
    eyebrow?: string;
    title?: string;
    description?: string;
    items?: Array<{
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
    }>;
  };

  let content = $state<ExperienceContent>({
    eyebrow: "Experience",
    title: "Professional Journey",
    description: "Timeline highlights (company, role, dates, impact)."
  });

  const fallbackTimeline: Item[] = [
    {
      id: "spade-integrity",
      company: "Spade Integrity",
      role: "IoT & Systems Engineer",
      where: "Paris, FR",
      dates: "2024 — Present",
      impact: [
        "Designed and shipped IoT platform components spanning embedded firmware, secure networking, and backend services.",
        "Built Go services and APIs for device telemetry, alarms, and operational tooling with an emphasis on reliability.",
        "Delivered secure connectivity patterns (VPN + MQTTS/TLS) and hardened remote operations for field devices.",
        "Implemented signal-processing pipelines and data flows to convert raw sensor events into actionable alerts."
      ],
      stack: ["Go", "PostgreSQL", "MQTT", "TLS", "OpenVPN/WireGuard", "Yocto", "Linux", "Docker", "Svelte"]
    },
    {
      id: "selected-projects",
      company: "Selected engineering projects",
      role: "Backend / Infra / Embedded",
      where: "Open-source & personal",
      dates: "2019 — Present",
      impact: [
        "Developed production-style backend patterns (REST APIs, caching, pagination, rate-limit aware integrations).",
        "Maintained reproducible Linux environments and tooling for fast iteration (dev shells, containers, automation).",
        "Prototyped IoT workflows end-to-end: device → broker → backend → dashboard."
      ],
      stack: ["Linux", "Go", "Networking", "Docker", "MQTT", "TypeScript", "Svelte"]
    }
  ];

  let timeline = $state<Item[]>(fallbackTimeline);

  $effect(() => {
    void (async () => {
      try {
        const res = await fetch("/api/content/experience");
        const j = await res.json();
        const d = (j?.data ?? {}) as ExperienceContent;

        const mapped: Item[] = (Array.isArray(d.items) ? d.items : [])
          .slice()
          .sort((a, b) => (a.sortOrder ?? 0) - (b.sortOrder ?? 0))
          .map((it) => ({
            id: it.id,
            company: it.org,
            role: it.role,
            where: it.location,
            dates: it.dates,
            impact: Array.isArray(it.summary) ? it.summary : [],
            stack: Array.isArray(it.bullets) ? it.bullets : [],
            image: it.image,
            proofUrl: it.proofUrl,
            sortOrder: it.sortOrder
          }))
          .filter((it) => it.id && it.company && it.role && it.dates && it.impact.length);

        if (mapped.length) {
          timeline = mapped;
        }

        content = {
          eyebrow: d.eyebrow || "Experience",
          title: d.title || "Professional Journey",
          description: d.description || "Timeline highlights (company, role, dates, impact)."
        };
      } catch {
        // Keep fallback timeline.
      }
    })();
  });
</script>

<div class="surface">
  <header class="hdr">
    <div class="path">/experience</div>
    <h1>{content.title || "Experience"}</h1>
    <p>{content.description || "Timeline highlights (company, role, dates, impact)."}</p>
  </header>

  <div class="list">
    {#each timeline as it (it.id)}
      <article class="card">
        <div class="top">
          <div class="who">
            <div class="company">{it.company}</div>
            <div class="role">{it.role}</div>
          </div>

          <div class="meta">
            <div class="where">{it.where}</div>
            <div class="dates">{it.dates}</div>
          </div>
        </div>

        <div class="grid">
          <div>
            <div class="label">Impact</div>
            <ul class="bul">
              {#each it.impact as x (x)}
                <li>{x}</li>
              {/each}
            </ul>
          </div>

          <div>
            <div class="label">Stack</div>
            <div class="tags">
              {#each it.stack as t (t)}
                <span class="tag">{t}</span>
              {/each}
            </div>
          </div>
        </div>
      </article>
    {/each}
  </div>
</div>

<style>
  .surface {
    border: 1px solid var(--border, rgba(255, 255, 255, 0.10));
    background: var(--surface, rgba(15, 20, 25, 0.35));
    backdrop-filter: blur(var(--blur, 18px));
    -webkit-backdrop-filter: blur(var(--blur, 18px));
    border-radius: var(--radius, 14px);
    padding: 16px;
  }

  .hdr {
    display: grid;
    gap: 8px;
    margin-bottom: 16px;
  }

  .path {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    color: rgba(255, 255, 255, 0.55);
  }

  h1 {
    margin: 0;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 22px;
    font-weight: 650;
    color: rgba(255, 255, 255, 0.92);
    letter-spacing: -0.25px;
  }

  p {
    margin: 0;
    font-family: var(--font-sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell,
      "Open Sans", "Helvetica Neue", sans-serif);
    color: rgba(255, 255, 255, 0.72);
    font-size: 14px;
    line-height: 1.55;
  }

  .list {
    display: grid;
    gap: 16px;
  }

  .card {
    border: 1px solid rgba(255, 255, 255, 0.08);
    background: var(--surface-2, rgba(15, 20, 25, 0.25));
    border-radius: 12px;
    padding: 16px;

    transition: transform 0.15s ease, background 0.15s ease, border-color 0.15s ease;
    position: relative;
    overflow: hidden;
  }

  .card::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 3px;
    height: 100%;
    background: linear-gradient(to bottom, var(--accent, rgba(228, 90, 90, 1)), rgba(228, 90, 90, 0));
  }

  .card:hover {
    background: rgba(15, 20, 25, 0.35);
    border-color: rgba(255, 255, 255, 0.12);
    transform: translateY(-2px);
  }

  .top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;

    padding-bottom: 12px;
    margin-bottom: 12px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  }

  .company {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    color: var(--accent, rgba(228, 90, 90, 0.95));
    font-size: 15px;
    font-weight: 600;
  }

  .role {
    margin-top: 2px;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    color: rgba(255, 255, 255, 0.9);
    font-size: 14px;
    font-weight: 550;
  }

  .meta {
    text-align: right;
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 13px;
    color: rgba(255, 255, 255, 0.55);
  }

  .where {
    color: rgba(255, 255, 255, 0.6);
  }

  .dates {
    margin-top: 2px;
    color: rgba(255, 255, 255, 0.78);
    font-weight: 550;
  }

  .grid {
    display: grid;
    grid-template-columns: 1.6fr 1fr;
    gap: 16px;
  }

  .label {
    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: rgba(255, 255, 255, 0.5);
    margin-bottom: 10px;
  }

  .bul {
    margin: 0;
    padding: 0;
    list-style: none;

    font-family: var(--font-sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell,
      "Open Sans", "Helvetica Neue", sans-serif);
    font-size: 13px;
    line-height: 1.7;
    color: rgba(255, 255, 255, 0.86);
  }

  .bul li {
    padding: 4px 0;
    margin-left: 10px;
    padding-left: 12px;
    border-left: 2px solid rgba(228, 90, 90, 0.22);
    position: relative;
  }

  .bul li::before {
    content: "→";
    position: absolute;
    left: 0;
    color: var(--accent, rgba(228, 90, 90, 0.95));
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .tag {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(228, 90, 90, 0.08);
    border-radius: 8px;
    padding: 6px 10px;

    font-family: var(--font-mono, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Menlo, monospace);
    font-size: 12px;
    color: var(--accent, rgba(228, 90, 90, 0.95));

    transition: transform 0.15s ease, background 0.15s ease;
  }

  .tag:hover {
    transform: translateY(-1px);
    background: rgba(228, 90, 90, 0.15);
  }

  @media (max-width: 900px) {
    .top {
      flex-direction: column;
      align-items: stretch;
      gap: 10px;
    }

    .meta {
      text-align: left;
    }

    .grid {
      grid-template-columns: 1fr;
    }
  }
</style>
