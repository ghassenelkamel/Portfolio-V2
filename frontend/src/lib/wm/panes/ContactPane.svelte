<script lang="ts">
  import { page } from "$app/stores";

  type ContactContent = {
    eyebrow?: string;
    title?: string;
    description?: string;
    email?: string;
    linkedin?: string;
    github?: string;
    portfolio?: string;
  };

  const fallbackContentEn: ContactContent = {
    eyebrow: "Contact",
    title: "Professional Reach",
    description: "Select a contact channel to discuss scope, impact, and collaboration.",
    email: "Ghassenelkamel@live.fr",
    linkedin: "https://linkedin.com/in/ghassenelkamel",
    github: "https://github.com/ghassenelkamel",
    portfolio: "https://ghassenelkamel.fr"
  };

  const fallbackContentFr: ContactContent = {
    eyebrow: "Contact",
    title: "Canaux professionnels",
    description: "Choisissez un canal pour discuter du périmètre, de l'impact et d'une collaboration.",
    email: "Ghassenelkamel@live.fr",
    linkedin: "https://linkedin.com/in/ghassenelkamel",
    github: "https://github.com/ghassenelkamel",
    portfolio: "https://ghassenelkamel.fr"
  };

  let content = $state<ContactContent>(fallbackContentEn);
  const langQuery = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr") ? "?lang=fr" : "");
  const isFr = $derived(langQuery === "?lang=fr");

  const ui = $derived.by(() => isFr
    ? {
        title: "Contact",
        desc: "Choisissez un canal pour discuter du périmètre, de l'impact et d'une collaboration.",
        from: "De",
        subject: "Objet",
        message: "Message",
        fromPlaceholder: "votre@email.com",
        subjectPlaceholder: "objet",
        messagePlaceholder: "ecrivez votre message...",
        sending: "envoi...",
        send: "envoyer",
        sent: "envoye",
        failed: "echec",
        unknown: "erreur inconnue"
      }
    : {
        title: "Contact",
        desc: "Select a contact channel to discuss scope, impact, and collaboration.",
        from: "From",
        subject: "Subject",
        message: "Message",
        fromPlaceholder: "your@email.com",
        subjectPlaceholder: "subject",
        messagePlaceholder: "write your message...",
        sending: "sending...",
        send: "send",
        sent: "sent",
        failed: "failed",
        unknown: "unknown error"
      });

  const activeFallbackContent = $derived(isFr ? fallbackContentFr : fallbackContentEn);

  let from = $state("");
  let subject = $state("");
  let message = $state("");
  let status = $state<string | null>(null);
  let sending = $state(false);

  $effect(() => {
    void (async () => {
      try {
        const res = await fetch(`/api/content/contact${langQuery}`);
        const j = await res.json();
        const d = (j?.data ?? {}) as ContactContent;
        content = { ...activeFallbackContent, ...d };
      } catch {
        content = activeFallbackContent;
      }
    })();
  });

  async function submit() {
    status = null;
    sending = true;
    try {
      const r = await fetch("/api/contact", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ name: "", email: from, subject, message })
      });
      const j = await r.json().catch(() => ({}));
      if (!r.ok) status = j?.error ?? `${ui.failed} (${r.status})`;
      else status = ui.sent;
    } catch (e) {
      status = e instanceof Error ? e.message : ui.unknown;
    } finally {
      sending = false;
      setTimeout(() => (status = null), 2200);
    }
  }
</script>

<div class="pane">
  <h2>{content.title || ui.title}</h2>
  <p class="desc">{content.description || ui.desc}</p>

  <div class="links">
    <a href={"mailto:" + (content.email || "Ghassenelkamel@live.fr")}>{content.email || "Ghassenelkamel@live.fr"}</a>
    <a href={content.linkedin || "https://linkedin.com/in/ghassenelkamel"} target="_blank" rel="noreferrer">linkedin</a>
    <a href={content.github || "https://github.com/ghassenelkamel"} target="_blank" rel="noreferrer">github</a>
  </div>

  <div class="form">
    <label>
      <span>{ui.from}</span>
      <input bind:value={from} type="email" placeholder={ui.fromPlaceholder} />
    </label>

    <label>
      <span>{ui.subject}</span>
      <input bind:value={subject} type="text" placeholder={ui.subjectPlaceholder} />
    </label>

    <label>
      <span>{ui.message}</span>
      <textarea bind:value={message} placeholder={ui.messagePlaceholder}></textarea>
    </label>

    <button type="button" class="btn" disabled={sending} onclick={submit}>
      {sending ? ui.sending : ui.send}
    </button>

    {#if status}
      <div class="status">{status}</div>
    {/if}
  </div>
</div>

<style>
  .pane {
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    color: rgba(255, 255, 255, 0.88);
  }

  h2 {
    margin: 0 0 10px;
    color: var(--accent, #e45a5a);
    font-size: 18px;
    letter-spacing: -0.25px;
  }

  .desc {
    margin: 0 0 10px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.65);
    line-height: 1.45;
  }

  .links {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
    margin: 0 0 12px;
  }

  .links a {
    color: rgba(255, 255, 255, 0.78);
    text-decoration: none;
    font-size: 12px;
    border-bottom: 1px dashed rgba(228, 90, 90, 0.3);
    padding-bottom: 2px;
  }

  .form {
    display: grid;
    gap: 12px;
  }

  label span {
    display: block;
    margin-bottom: 6px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.60);
  }

  input,
  textarea {
    width: 100%;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.10);
    background: rgba(255, 255, 255, 0.04);
    color: rgba(255, 255, 255, 0.90);
    padding: 10px 12px;
    outline: none;
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
  }

  input:focus,
  textarea:focus {
    border-color: rgba(228, 90, 90, 0.35);
    box-shadow: 0 0 0 3px rgba(228, 90, 90, 0.10);
  }

  textarea {
    min-height: 150px;
    resize: vertical;
  }

  .btn {
    border: 1px solid rgba(228, 90, 90, 0.35);
    background: rgba(228, 90, 90, 0.10);
    color: rgba(255, 255, 255, 0.92);
    padding: 10px 12px;
    border-radius: 12px;
    cursor: pointer;
    font-family: var(--font-mono, ui-monospace, Menlo, Consolas, monospace);
    transition: transform 0.15s ease, background 0.15s ease;
  }

  .btn:hover {
    transform: translateY(-1px);
    background: rgba(228, 90, 90, 0.15);
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: default;
    transform: none;
  }

  .status {
    color: rgba(255, 255, 255, 0.65);
    font-size: 12px;
  }
</style>
