<script lang="ts">
  import { page } from "$app/stores";

  type Status = { kind: 'ok' | 'err' | 'info'; msg: string } | null;

  type ContactContent = {
    eyebrow?: string;
    title?: string;
    description?: string;
    email?: string;
    linkedin?: string;
    github?: string;
    portfolio?: string;
  };

  const { data } = $props<{ data: { content?: ContactContent } }>();
  const content = $derived.by(() => data?.content ?? {});
  const isFr = $derived(($page.url.searchParams.get("lang") || "").toLowerCase().startsWith("fr"));

  const ui = $derived.by(() => isFr
    ? {
        title: "Contact",
        header: "Canaux professionnels",
        description: "Choisissez un canal pour discuter du périmètre, de l'impact et d'une collaboration.",
        to: "Vers",
        copy: "copier",
        from: "De",
        fromPlaceholder: "votre@email.com",
        subject: "Objet",
        message: "Message",
        sending: "envoi...",
        send: "envoyer",
        emailCopied: "Email copié.",
        clipboardUnavailable: "Presse-papiers indisponible. Utilisez le lien mailto.",
        validEmail: "Entrez un email valide.",
        subjectShort: "Objet trop court.",
        messageShort: "Message trop court (min 10 caractères).",
        messageLong: "Message trop long (max 4000 caractères).",
        sendingStatus: "Envoi en cours...",
        sendFailed: "Echec d'envoi",
        sentThanks: "Message envoyé. Merci !",
        networkError: "Erreur réseau.",
        subjectPlaceholder: "Quel sujet souhaitez-vous aborder ?",
        messagePlaceholder: "Ecrivez votre message..."
      }
    : {
        title: "Contact",
        header: "Professional Reach",
        description: "Select a contact channel to discuss scope, impact, and collaboration.",
        to: "To",
        copy: "copy",
        from: "From",
        fromPlaceholder: "your@email.com",
        subject: "Subject",
        message: "Message",
        sending: "sending...",
        send: "send",
        emailCopied: "Email copied.",
        clipboardUnavailable: "Clipboard not available. Use the mailto link.",
        validEmail: "Enter a valid email.",
        subjectShort: "Subject is too short.",
        messageShort: "Message is too short (min 10 chars).",
        messageLong: "Message is too long (max 4000 chars).",
        sendingStatus: "Sending...",
        sendFailed: "Send failed",
        sentThanks: "Message sent. Thanks!",
        networkError: "Network error.",
        subjectPlaceholder: "What do you want to talk about?",
        messagePlaceholder: "Write your message..."
      });

  const toEmail = $derived(content.email || "Ghassenelkamel@live.fr");
  const linkedInUrl = $derived(content.linkedin || "https://linkedin.com/in/ghassenelkamel");
  const githubUrl = $derived(content.github || "https://github.com/ghassenelkamel");
  const portfolioUrl = $derived(content.portfolio || "https://ghassenelkamel.fr");

  let sending = $state(false);
  let status = $state<Status>(null);

  let from = $state("");
  let subject = $state("");
  let message = $state("");

  function setStatus(s: Status, ms = 2200) {
    status = s;
    if (s) setTimeout(() => (status = null), ms);
  }

  async function copyEmail() {
    try {
      await navigator.clipboard.writeText(toEmail);
      setStatus({ kind: "ok", msg: ui.emailCopied }, 1600);
    } catch {
      setStatus({ kind: "info", msg: ui.clipboardUnavailable }, 2200);
    }
  }

  function validate() {
    const e = from.trim();
    const s = subject.trim();
    const m = message.trim();

    if (!e || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e)) return ui.validEmail;
    if (!s || s.length < 3) return ui.subjectShort;
    if (!m || m.length < 10) return ui.messageShort;
    if (m.length > 4000) return ui.messageLong;
    return null;
  }

  async function submit() {
    const err = validate();
    if (err) return setStatus({ kind: "err", msg: err }, 2800);

    sending = true;
    setStatus({ kind: "info", msg: ui.sendingStatus }, 1200);

    try {
      const res = await fetch("/api/contact", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          name: "",
          email: from.trim(),
          subject: subject.trim(),
          message: message.trim()
        })
      });

      const txt = await res.text().catch(() => "");
      if (!res.ok) {
        let details = txt;
        try {
          const parsed = JSON.parse(txt) as { error?: string };
          if (parsed?.error) details = parsed.error;
        } catch {
          // Keep raw response text.
        }
        return setStatus({ kind: "err", msg: `${ui.sendFailed} (${res.status}). ${details}` }, 5000);
      }

      setStatus({ kind: "ok", msg: ui.sentThanks }, 2500);
      from = "";
      subject = "";
      message = "";
    } catch (e) {
      setStatus({ kind: "err", msg: e instanceof Error ? e.message : ui.networkError }, 3500);
    } finally {
      sending = false;
    }
  }
</script>

<svelte:head>
  <title>{ui.title}</title>
</svelte:head>

<div class="wrap">
  <div class="title">
    <div class="eyebrow">{content.eyebrow || ui.title}</div>
    <div class="t">{content.title || ui.header}</div>
    <div class="s">{content.description || ui.description}</div>
  </div>

  <div class="card">
    <div class="row">
      <div class="kv">
        <span>{ui.to}</span>
        <a href={"mailto:" + toEmail}>{toEmail}</a>
      </div>

      <button type="button" class="ghostBtn" onclick={copyEmail}>{ui.copy}</button>
    </div>

    <div class="form">
      <label>
        <span>{ui.from}</span>
        <input bind:value={from} type="email" placeholder={ui.fromPlaceholder} autocomplete="email" />
      </label>

      <label>
        <span>{ui.subject}</span>
        <input bind:value={subject} type="text" placeholder={ui.subjectPlaceholder} />
      </label>

      <label>
        <span>{ui.message}</span>
        <textarea bind:value={message} placeholder={ui.messagePlaceholder}></textarea>
      </label>

      <div class="actions">
        <button class="btn" type="button" disabled={sending} onclick={submit}>
          {sending ? ui.sending : ui.send}
        </button>
        <a class="link" href={linkedInUrl} target="_blank" rel="noreferrer">linkedin</a>
        <a class="link" href={githubUrl} target="_blank" rel="noreferrer">github</a>
        <a class="link" href={portfolioUrl} target="_blank" rel="noreferrer">portfolio</a>
      </div>

      {#if status}
        <div class={"status " + status.kind}>{status.msg}</div>
      {/if}
    </div>
  </div>
</div>

<style>
  .wrap { display: grid; gap: 14px; }

  .eyebrow {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: rgba(245,245,245,0.58);
  }

  .t {
    margin-top: 7px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 24px;
    line-height: 1.2;
    letter-spacing: 0.4px;
    color: rgba(228,90,90,0.95);
  }

  .s {
    margin-top: 8px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.72);
    font-size: 13px;
    line-height: 1.5;
    max-width: 66ch;
  }

  .card {
    border: 1px solid rgba(204,102,102,0.22);
    border-radius: 14px;
    background: rgba(12,14,18,0.18);
    padding: 14px;
  }

  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding-bottom: 12px;
    margin-bottom: 12px;
    border-bottom: 1px solid rgba(204,102,102,0.18);
  }

  .kv {
    display: flex;
    gap: 12px;
    align-items: baseline;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 13px;
  }
  .kv span { color: rgba(245,245,245,0.55); width: 40px; }
  .kv a { color: rgba(245,245,245,0.86); text-decoration: none; border-bottom: 1px dashed rgba(228,90,90,0.35); }
  .kv a:hover { border-bottom-color: rgba(228,90,90,0.75); }

  .ghostBtn {
    border: 1px solid rgba(204,102,102,0.22);
    background: rgba(12,14,18,0.14);
    color: rgba(245,245,245,0.82);
    border-radius: 999px;
    padding: 8px 12px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    cursor: pointer;
    transition: background 0.15s ease, transform 0.15s ease;
  }
  .ghostBtn:hover { background: rgba(204,102,102,0.12); transform: translateY(-1px); }

  .form { display: grid; gap: 12px; }

  label span {
    display: block;
    margin-bottom: 6px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(228,90,90,0.75);
  }

  input, textarea {
    width: 100%;
    border-radius: 12px;
    border: 1px solid rgba(204,102,102,0.18);
    background: rgba(12,14,18,0.12);
    color: rgba(245,245,245,0.86);
    padding: 10px 12px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    outline: none;
  }

  textarea { min-height: 170px; resize: vertical; }
  input:focus, textarea:focus { border-color: rgba(228,90,90,0.45); }

  .actions { display: flex; align-items: center; gap: 10px; margin-top: 6px; }

  .btn {
    border: 1px solid rgba(228,90,90,0.45);
    background: rgba(228,90,90,0.10);
    color: rgba(245,245,245,0.90);
    border-radius: 12px;
    padding: 10px 12px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    cursor: pointer;
    transition: transform 0.15s ease, background 0.15s ease;
  }
  .btn:hover { transform: translateY(-1px); background: rgba(204,102,102,0.12); }
  .btn:disabled { opacity: 0.6; cursor: default; transform: none; }

  .link {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.78);
    text-decoration: none;
    border-bottom: 1px dashed rgba(204,102,102,0.30);
    padding-bottom: 2px;
  }
  .link:hover { border-bottom-color: rgba(228,90,90,0.75); }

  .status { margin-top: 8px; font-family: ui-monospace, Menlo, Consolas, monospace; font-size: 12px; }
  .status.ok { color: rgba(140, 240, 170, 0.85); }
  .status.err { color: rgba(255, 140, 140, 0.90); }
  .status.info { color: rgba(245,245,245,0.68); }

  @media (max-width: 760px) {
    .t {
      font-size: 21px;
    }

    .row {
      flex-direction: column;
      align-items: flex-start;
    }

    .kv {
      width: 100%;
      flex-wrap: wrap;
      gap: 8px;
    }

    .actions {
      flex-wrap: wrap;
      align-items: stretch;
      gap: 8px;
    }

    .btn {
      width: 100%;
    }

    .link {
      display: inline-flex;
      align-items: center;
      min-height: 32px;
    }

    input,
    textarea {
      font-size: 16px;
    }
  }

  @media (max-width: 520px) {
    .card,
    .row {
      border-radius: 12px;
    }

    .kv a {
      overflow-wrap: anywhere;
    }

    .link {
      width: 100%;
      justify-content: flex-start;
    }
  }
</style>
