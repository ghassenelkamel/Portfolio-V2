<script lang="ts">
  import { page } from "$app/stores";

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

  const fallbackContentEn: ExperienceContent = {
    eyebrow: "Experience",
    title: "Professional Journey",
    description: "Timeline highlights (company, role, dates, impact)."
  };

  const fallbackContentFr: ExperienceContent = {
    eyebrow: "Experience",
    title: "Parcours professionnel",
    description: "Points cles du parcours (entreprise, role, periodes, impact)."
  };

  let content = $state<ExperienceContent>(fallbackContentEn);

  const fallbackTimeline: Item[] = [
    {
      id: "spade-integrity",
      company: "Spade Integrity",
      role: "IoT Engineer",
      where: "On-site",
      dates: "May 2024 - Present · 2 yrs 1 mo",
      impact: [
        "Contribute to client platform development with focus on reliability, security, and operational continuity.",
        "Work on secure data recovery flows, NAS integration, and synchronization systems for alarms and critical data.",
        "Support end-to-end IoT/backend delivery by improving data pipelines, monitoring quality, and secure service integrations."
      ],
      stack: ["Go", "IoT platform", "Secure data recovery", "NAS integration", "Alarm synchronization", "Monitoring", "Infrastructure"]
    },
    {
      id: "bookbeo-2023",
      company: "bookBeo",
      role: "Full Stack Mobile Developer",
      where: "Rennes, Brittany, France · On-site",
      dates: "Mar 2023 - Sep 2023 · 7 mos",
      impact: [
        "Developed mobile product features and backend-connected workflows during internship delivery cycles.",
        "Collaborated on release quality and implementation speed with structured team practices and version control discipline."
      ],
      stack: ["GitLab", "Git", "React Native", "API integration", "Product delivery", "Performance", "Team collaboration"]
    },
    {
      id: "st-2022",
      company: "STMicroelectronics",
      role: "Security Project Team Member",
      where: "Le Mans, France · On-site",
      dates: "Sep 2022 – Mar 2023",
      impact: [
        "Worked on PSA Crypto API implementation under embedded constraints and strict security requirements.",
        "Delivered C code with structured Git workflows and Linux-based validation tooling."
      ],
      stack: ["C", "PSA Crypto", "Embedded security", "Linux tooling"]
    },
    {
      id: "falcon-2022",
      company: "Falcon Inspection & Services",
      role: "Mobile Developer Experience",
      where: "Tunis, Tunisia · On-site",
      dates: "Jun 2022 – Aug 2022",
      impact: [
        "Built a mobile workflow to record and validate multiple industrial tests in the field.",
        "Implemented per-test data models and rule validation to improve consistency of reports."
      ],
      stack: ["Flutter", "Node.js", "Validation flows", "Docker"]
    },
    {
      id: "dall-2021",
      company: "DigiArtLivingLab",
      role: "Weather Monitoring & Data Storage",
      where: "Tunis, Tunisia · On-site",
      dates: "Jul 2021 – Aug 2021",
      impact: [
        "Automated weather data collection and storage pipelines for continuous monitoring.",
        "Built UI modules and API integrations with JSON parsing and backend persistence."
      ],
      stack: ["Java", "JavaFX", "API integration", "Automation"]
    },
    {
      id: "sagemcom-2020",
      company: "SAGEMCOM",
      role: "Electronic Components Validation",
      where: "Tunis, Tunisia · On-site",
      dates: "Feb 2020 – May 2020",
      impact: [
        "Built LabView workflows for detection and qualification of electronic components.",
        "Trained a CNN model and deployed a lightweight inference variant on Raspberry Pi hardware."
      ],
      stack: ["LabView", "CNN", "TensorFlow Lite", "Raspberry Pi", "SolidWorks"]
    },
    {
      id: "dall-2019",
      company: "DigiArtLivingLab",
      role: "IoT Station + Mobile Apps",
      where: "Tunis, Tunisia · On-site",
      dates: "Jan 2019 – Feb 2019",
      impact: [
        "Developed Arduino-based sensing builds with Bluetooth communication.",
        "Created Android companion apps and live data sync under tight timeline constraints."
      ],
      stack: ["Arduino", "Bluetooth", "MIT App Inventor", "Firebase"]
    },
    {
      id: "dall-2018",
      company: "DigiArtLivingLab",
      role: "Educational Game Development",
      where: "Nabeul, Tunisia · On-site",
      dates: "Jan 2018 – Feb 2018",
      impact: [
        "Built a Unity educational game focused on directional learning for children.",
        "Coordinated assets, level pacing, and gameplay polish with a small collaborative team."
      ],
      stack: ["Unity", "Gameplay design", "Team collaboration", "C++ basics"]
    }
  ];

  let timeline = $state<Item[]>(fallbackTimeline);
  const langQuery = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr") ? "?lang=fr" : "");
  const isFr = $derived(langQuery === "?lang=fr");

  const ui = $derived.by(() => isFr
    ? {
        eyebrow: "Experience",
        title: "Parcours professionnel",
        desc: "Points cles du parcours (entreprise, role, periodes, impact).",
        impact: "Impact",
        stack: "Technologies"
      }
    : {
        eyebrow: "Experience",
        title: "Professional Journey",
        desc: "Timeline highlights (company, role, dates, impact).",
        impact: "Impact",
        stack: "Stack"
      });

  const fallbackTimelineFr: Item[] = [
    {
      id: "spade-integrity",
      company: "Spade Integrity",
      role: "Ingenieur IoT",
      where: "Sur site",
      dates: "Mai 2024 - Present · 2 ans 1 mois",
      impact: [
        "Contribue au developpement de la plateforme client avec un focus sur la fiabilite, la securite et la continuite de service.",
        "Travaille sur des flux de recuperation securisee des donnees, l'integration NAS et les mecanismes de synchronisation des alarmes et donnees critiques.",
        "Soutient les livraisons IoT/backend de bout en bout via l'amelioration des pipelines de donnees, du monitoring et des integrations securisees."
      ],
      stack: ["Go", "Plateforme IoT", "Recuperation securisee", "Integration NAS", "Synchronisation alarmes", "Monitoring", "Infrastructure"]
    },
    {
      id: "bookbeo-2023",
      company: "bookBeo",
      role: "Developpeur mobile full stack",
      where: "Rennes, Bretagne, France · Sur site",
      dates: "Mars 2023 - Sept 2023 · 7 mois",
      impact: [
        "Developpement de fonctionnalites mobiles et de workflows connectes au backend pendant les cycles de livraison de stage.",
        "Contribution a la qualite des releases et a la vitesse d'implementation grace a des pratiques d'equipe structurees et une discipline Git rigoureuse."
      ],
      stack: ["GitLab", "Git", "React Native", "Integration API", "Livraison produit", "Performance", "Collaboration equipe"]
    },
    {
      id: "st-2022",
      company: "STMicroelectronics",
      role: "Membre d'equipe projet securite",
      where: "Le Mans, France · Sur site",
      dates: "Sept 2022 - Mars 2023",
      impact: [
        "Travail sur l'implementation de l'API PSA Crypto sous contraintes embarquees et exigences de securite strictes.",
        "Livraison de code C avec des workflows Git structures et des outils de validation sous Linux."
      ],
      stack: ["C", "PSA Crypto", "Securite embarquee", "Outils Linux"]
    },
    {
      id: "falcon-2022",
      company: "Falcon Inspection & Services",
      role: "Experience developpement mobile",
      where: "Tunis, Tunisie · Sur site",
      dates: "Juin 2022 - Aout 2022",
      impact: [
        "Concu un workflow mobile pour enregistrer et valider plusieurs tests industriels sur le terrain.",
        "Implemente des modeles de donnees par test et des regles de validation pour ameliorer la coherence des rapports."
      ],
      stack: ["Flutter", "Node.js", "Flux de validation", "Docker"]
    },
    {
      id: "dall-2021",
      company: "DigiArtLivingLab",
      role: "Suivi meteo et stockage des donnees",
      where: "Tunis, Tunisie · Sur site",
      dates: "Juil 2021 - Aout 2021",
      impact: [
        "Automatisation de la collecte et du stockage des donnees meteo pour un suivi continu.",
        "Developpement de modules UI et d'integrations API avec parsing JSON et persistance backend."
      ],
      stack: ["Java", "JavaFX", "Integration API", "Automatisation"]
    },
    {
      id: "sagemcom-2020",
      company: "SAGEMCOM",
      role: "Validation de composants electroniques",
      where: "Tunis, Tunisie · Sur site",
      dates: "Fev 2020 - Mai 2020",
      impact: [
        "Conception de workflows LabView pour la detection et la qualification de composants electroniques.",
        "Entrainement d'un modele CNN et deploiement d'une version d'inference legere sur Raspberry Pi."
      ],
      stack: ["LabView", "CNN", "TensorFlow Lite", "Raspberry Pi", "SolidWorks"]
    },
    {
      id: "dall-2019",
      company: "DigiArtLivingLab",
      role: "Station IoT + applications mobiles",
      where: "Tunis, Tunisie · Sur site",
      dates: "Janv 2019 - Fev 2019",
      impact: [
        "Developpement de solutions de capteurs Arduino avec communication Bluetooth.",
        "Creation d'applications Android compagnon avec synchronisation temps reel des donnees sous contrainte de delai."
      ],
      stack: ["Arduino", "Bluetooth", "MIT App Inventor", "Firebase"]
    },
    {
      id: "dall-2018",
      company: "DigiArtLivingLab",
      role: "Developpement de jeu educatif",
      where: "Nabeul, Tunisie · Sur site",
      dates: "Janv 2018 - Fev 2018",
      impact: [
        "Creation d'un jeu educatif Unity axe sur l'apprentissage des directions pour enfants.",
        "Coordination des assets, du rythme des niveaux et du polish gameplay avec une petite equipe collaborative."
      ],
      stack: ["Unity", "Gameplay design", "Collaboration equipe", "Bases C++"]
    }
  ];

  const activeFallbackTimeline = $derived(isFr ? fallbackTimelineFr : fallbackTimeline);
  const activeFallbackContent = $derived(isFr ? fallbackContentFr : fallbackContentEn);

  $effect(() => {
    void (async () => {
      try {
        const res = await fetch(`/api/content/experience${langQuery}`);
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
          const byID = new Map<string, Item>(activeFallbackTimeline.map((x) => [x.id, x]));
          for (const m of mapped) byID.set(m.id, m);

          const merged = [...byID.values()];
          const fallbackOrder = new Map(activeFallbackTimeline.map((x, i) => [x.id, i]));
          merged.sort((a, b) => {
            const ia = fallbackOrder.get(a.id);
            const ib = fallbackOrder.get(b.id);
            if (ia != null && ib != null) return ia - ib;
            if (ia != null) return -1;
            if (ib != null) return 1;
            return a.id.localeCompare(b.id);
          });

          timeline = merged;
        }

        content = { ...activeFallbackContent, ...d };
      } catch {
        timeline = activeFallbackTimeline;
        content = activeFallbackContent;
      }
    })();
  });
</script>

<div class="surface">
  <header class="hdr">
    <div class="path">/experience</div>
    <h1>{content.title || ui.eyebrow}</h1>
    <p>{content.description || ui.desc}</p>
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
            <div class="label">{ui.impact}</div>
            <ul class="bul">
              {#each it.impact as x (x)}
                <li>{x}</li>
              {/each}
            </ul>
          </div>

          <div>
            <div class="label">{ui.stack}</div>
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
