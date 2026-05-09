<script lang="ts">
  import { page } from "$app/stores";

  type AboutContent = {
    headline?: string;
    subtitle?: string;
    summary?: string[];
    keywords?: string[];
    specialties?: string[];
    currently?: string[];
    goals?: string[];
    email?: string;
    site?: string;
    github?: string;
  };

  const { data } = $props<{ data: { content?: AboutContent } }>();
  const content = $derived.by(() => data?.content ?? {});

  const isFr = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr"));

  const defaultKeywords = $derived.by(() => isFr
    ? ["Architecture IoT", "Backend Go", "Linux embarque", "MQTT / TLS", "Reseaux VPN", "APIs d'automatisation"]
    : ["IoT architecture", "Go backend", "Embedded Linux", "MQTT / TLS", "VPN networking", "Automation APIs"]);

  const ui = $derived.by(() => isFr
    ? {
        title: "À propos",
        kicker: "à propos",
        copy: "copier",
        copied: "copié",
        clipboard: "Presse-papiers indisponible",
        about: "À propos",
        specialize: "Je me spécialise en",
        noSpecialties: "Aucune spécialité configurée.",
        currently: "Actuellement",
        noCurrently: "Aucun element 'Actuellement'.",
        goals: "Objectifs",
        noGoals: "Aucun objectif renseigné.",
        defaultIntro: "Je construis des systèmes fiables pour les produits connectés : connectivité sécurisée, flux de données propres et services backend robustes en production.",
        emailLabel: "Email",
        linkLabel: "Lien",
        rowEmail: "courriel",
        rowSite: "site",
        rowGithub: "github"
      }
    : {
        title: "About",
        kicker: "about",
        copy: "copy",
        copied: "copied",
        clipboard: "Clipboard not available",
        about: "About",
        specialize: "I specialize in",
        noSpecialties: "No specialties configured.",
        currently: "Currently",
        noCurrently: "No 'Currently' items yet.",
        goals: "Goals",
        noGoals: "No 'Goals' items yet.",
        defaultIntro: "I build reliable systems for connected products: secure connectivity, clean data flow, and backend services that hold up in production.",
        emailLabel: "Email",
        linkLabel: "Link",
        rowEmail: "email",
        rowSite: "site",
        rowGithub: "github"
      });

  const headline = $derived.by(() => content.headline || "Ghassen El Kamel");
  const subtitle = $derived.by(() => content.subtitle || "IoT & Systems Engineer · Backend · Infrastructure");
  const summary = $derived.by(() => Array.isArray(content.summary) ? content.summary.filter(Boolean) : []);
  const keywords = $derived.by(() => Array.isArray(content.keywords) && content.keywords.length
    ? content.keywords.filter(Boolean)
    : defaultKeywords);
  const specialties = $derived.by(() => Array.isArray(content.specialties) ? content.specialties.filter(Boolean) : []);
  const now = $derived.by(() => Array.isArray(content.currently) ? content.currently.filter(Boolean) : []);
  const goals = $derived.by(() => Array.isArray(content.goals) ? content.goals.filter(Boolean) : []);
  const email = $derived.by(() => content.email || "Ghassenelkamel@live.fr");
  const site = $derived.by(() => content.site || "https://ghassenelkamel.fr");
  const github = $derived.by(() => content.github || "https://github.com/ghassenelkamel");

  let copied = $state<string | null>(null);

  async function copy(text: string, label: string) {
    try {
      await navigator.clipboard.writeText(text);
      copied = `${label} ${ui.copied}`;
      setTimeout(() => (copied = null), 1200);
    } catch {
      copied = ui.clipboard;
      setTimeout(() => (copied = null), 1400);
    }
  }
</script>

<svelte:head>
  <title>{ui.title}</title>
</svelte:head>

<div class="about">
  <section class="hero">
    <img class="avatar" src="/assets/avatar.png" alt="Ghassen El Kamel" />

    <div class="heroText">
      <div class="kicker">{ui.kicker}</div>
      <h1>{headline}</h1>
      <div class="sub">{subtitle}</div>

      <div class="intro">
        {#if summary.length}
          {#each summary as p}
            <p>{p}</p>
          {/each}
        {:else}
          <p>{ui.defaultIntro}</p>
        {/if}
      </div>

      <div class="links">
        <div class="linkRow">
          <span class="label">{ui.rowEmail}</span>
          <a class="link" href={"mailto:" + email}>{email}</a>
          <button class="mini" type="button" onclick={() => copy(email, ui.emailLabel)}>{ui.copy}</button>
        </div>

        <div class="linkRow">
          <span class="label">{ui.rowSite}</span>
          <a class="link" href={site} target="_blank" rel="noreferrer">{site}</a>
          <button class="mini" type="button" onclick={() => copy(site, ui.linkLabel)}>{ui.copy}</button>
        </div>

        <div class="linkRow">
          <span class="label">{ui.rowGithub}</span>
          <a class="link" href={github} target="_blank" rel="noreferrer">{github.replace(/^https?:\/\//, "")}</a>
        </div>

        {#if copied}
          <div class="toast">{copied}</div>
        {/if}
      </div>
    </div>
  </section>

  <section class="grid">
    <div class="card">
      <div class="cardTitle">{ui.about}</div>
      <div class="chips">
        {#each keywords.slice(0, 6) as s}
          <span class="chip">{s}</span>
        {/each}
      </div>
      <div class="hr"></div>
      <div class="listTitle">{ui.specialize}</div>
      {#if specialties.length}
        <ul class="list">
          {#each specialties as it}
            <li>{it}</li>
          {/each}
        </ul>
      {:else}
        <div class="empty">{ui.noSpecialties}</div>
      {/if}
    </div>

    <div class="card">
      <div class="cardTitle">{ui.currently}</div>
      {#if now.length}
        <ul class="list">
          {#each now as it}
            <li>{it}</li>
          {/each}
        </ul>
      {:else}
        <div class="empty">{ui.noCurrently}</div>
      {/if}
    </div>

    <div class="card">
      <div class="cardTitle">{ui.goals}</div>
      {#if goals.length}
        <ul class="list">
          {#each goals as it}
            <li>{it}</li>
          {/each}
        </ul>
      {:else}
        <div class="empty">{ui.noGoals}</div>
      {/if}
    </div>
  </section>
</div>

<style>
  .about { padding: 18px; display: grid; gap: 16px; }

  .hero {
    display: grid;
    grid-template-columns: 84px 1fr;
    gap: 14px;
    align-items: start;

    padding: 14px;
    border-radius: 14px;
    border: 1px solid rgba(245, 245, 245, 0.10);
    background: rgba(12, 14, 18, 0.34);
    backdrop-filter: blur(22px) saturate(150%);
    -webkit-backdrop-filter: blur(22px) saturate(150%);
  }

  .avatar {
    width: 84px;
    height: 84px;
    border-radius: 999px;
    object-fit: cover;
    border: 1px solid rgba(204, 102, 102, 0.30);
    box-shadow: 0 10px 30px rgba(0,0,0,0.35);
  }

  .kicker {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    letter-spacing: 1.2px;
    color: rgba(228, 90, 90, 0.75);
    text-transform: uppercase;
    margin-bottom: 6px;
  }

  h1 {
    margin: 0;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 22px;
    letter-spacing: 0.6px;
    color: rgba(245, 245, 245, 0.94);
  }

  .sub {
    margin-top: 6px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 13px;
    color: rgba(245, 245, 245, 0.70);
  }

  .intro { margin-top: 12px; display: grid; gap: 10px; }
  .intro p { margin: 0; line-height: 1.55; color: rgba(245, 245, 245, 0.88); font-size: 14px; }

  .links { margin-top: 12px; display: grid; gap: 8px; }
  .linkRow { display: grid; grid-template-columns: 70px 1fr auto; gap: 10px; align-items: center; }

  .label {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.55);
  }

  .link {
    text-decoration: none;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.86);
    border-bottom: 1px dashed rgba(228, 90, 90, 0.35);
    padding-bottom: 2px;
    overflow-wrap: anywhere;
  }
  .link:hover { border-bottom-color: rgba(228, 90, 90, 0.75); }

  .mini {
    border: 1px solid rgba(204, 102, 102, 0.25);
    background: rgba(12, 14, 18, 0.20);
    color: rgba(245, 245, 245, 0.82);
    border-radius: 999px;
    padding: 6px 10px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    cursor: pointer;
  }
  .mini:hover { background: rgba(204, 102, 102, 0.14); }

  .toast {
    margin-top: 6px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.70);
  }

  .grid { display: grid; grid-template-columns: 1.2fr 1fr 1fr; gap: 14px; }

  .card {
    padding: 14px;
    border-radius: 14px;
    border: 1px solid rgba(245, 245, 245, 0.10);
    background: rgba(12, 14, 18, 0.28);
    backdrop-filter: blur(22px) saturate(150%);
    -webkit-backdrop-filter: blur(22px) saturate(150%);
    overflow: hidden;
  }

  .cardTitle {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    letter-spacing: 1.1px;
    color: rgba(228, 90, 90, 0.75);
    text-transform: uppercase;
    margin-bottom: 10px;
  }

  .chips { display: flex; flex-wrap: wrap; gap: 8px; }
  .chip {
    padding: 7px 10px;
    border-radius: 999px;
    border: 1px solid rgba(204, 102, 102, 0.22);
    background: rgba(204, 102, 102, 0.10);
    color: rgba(245, 245, 245, 0.86);
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    white-space: nowrap;
  }

  .hr { height: 1px; background: rgba(204, 102, 102, 0.18); margin: 12px 0; }

  .listTitle {
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245, 245, 245, 0.72);
    margin-bottom: 8px;
  }

  .list {
    margin: 0;
    padding-left: 18px;
    display: grid;
    gap: 8px;
    color: rgba(245, 245, 245, 0.80);
    font-size: 13px;
    line-height: 1.45;
  }

  .empty { color: rgba(245, 245, 245, 0.62); font-size: 13px; line-height: 1.45; }

  @media (max-width: 950px) {
    .grid { grid-template-columns: 1fr; }
    .hero { grid-template-columns: 72px 1fr; }
    .avatar { width: 72px; height: 72px; }
  }

  @media (max-width: 760px) {
    .about {
      padding: 12px;
      gap: 12px;
    }

    .hero {
      grid-template-columns: 1fr;
      gap: 12px;
      padding: 12px;
    }

    .avatar {
      width: 68px;
      height: 68px;
      border-radius: 999px;
    }

    h1 {
      font-size: 20px;
    }

    .sub {
      font-size: 12px;
    }

    .linkRow {
      grid-template-columns: 56px minmax(0, 1fr);
      align-items: start;
    }

    .mini {
      grid-column: 2;
      justify-self: start;
      margin-top: 4px;
    }
  }

  @media (max-width: 520px) {
    .card,
    .hero {
      border-radius: 12px;
    }

    .chip {
      white-space: normal;
      line-height: 1.3;
    }
  }
</style>
