<script lang="ts">
  import { page } from "$app/stores";
  import { onDestroy } from "svelte";

  type Status = { kind: 'ok' | 'err' | 'info'; msg: string } | null;
  type FieldErrors = { from?: string; subject?: string; message?: string };

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
        messagePlaceholder: "Ecrivez votre message...",
        allFieldsRequired: "Tous les champs ci-dessous sont obligatoires.",
        emailRequired: "L'email est obligatoire.",
        subjectRequired: "L'objet est obligatoire.",
        messageRequired: "Le message est obligatoire.",
        fixErrors: "Veuillez corriger les champs en erreur puis réessayer.",
        dismiss: "fermer"
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
        messagePlaceholder: "Write your message...",
        allFieldsRequired: "All fields below are required.",
        emailRequired: "Email is required.",
        subjectRequired: "Subject is required.",
        messageRequired: "Message is required.",
        fixErrors: "Please fix the highlighted fields and try again.",
        dismiss: "dismiss"
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
  let fieldErrors = $state<FieldErrors>({});
  let statusTimer: ReturnType<typeof setTimeout> | null = null;

  function setStatus(s: Status, autoHideMs = 0) {
    if (statusTimer) {
      clearTimeout(statusTimer);
      statusTimer = null;
    }
    status = s;
    if (s && autoHideMs > 0) {
      statusTimer = setTimeout(() => {
        status = null;
        statusTimer = null;
      }, autoHideMs);
    }
  }

  function clearStatus() {
    if (statusTimer) {
      clearTimeout(statusTimer);
      statusTimer = null;
    }
    status = null;
  }

  function clearFieldError(field: keyof FieldErrors) {
    if (!fieldErrors[field]) return;
    fieldErrors = { ...fieldErrors, [field]: undefined };
  }

  onDestroy(() => {
    if (statusTimer) clearTimeout(statusTimer);
  });

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

    const errors: FieldErrors = {};

    if (!e) errors.from = ui.emailRequired;
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e)) errors.from = ui.validEmail;

    if (!s) errors.subject = ui.subjectRequired;
    else if (s.length < 3) errors.subject = ui.subjectShort;

    if (!m) errors.message = ui.messageRequired;
    else if (m.length < 10) errors.message = ui.messageShort;
    else if (m.length > 4000) errors.message = ui.messageLong;

    return errors;
  }

  async function submit() {
    const errors = validate();
    fieldErrors = errors;

    if (errors.from || errors.subject || errors.message) {
      setStatus({ kind: "err", msg: ui.fixErrors });
      return;
    }

    sending = true;
    setStatus({ kind: "info", msg: ui.sendingStatus });

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
        return setStatus({ kind: "err", msg: `${ui.sendFailed} (${res.status}). ${details}` });
      }

      setStatus({ kind: "ok", msg: ui.sentThanks });
      from = "";
      subject = "";
      message = "";
      fieldErrors = {};
    } catch (e) {
      setStatus({ kind: "err", msg: e instanceof Error ? e.message : ui.networkError });
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
      <p class="formHint">{ui.allFieldsRequired}</p>

      <label>
        <span>{ui.from}</span>
        <input
          bind:value={from}
          type="email"
          placeholder={ui.fromPlaceholder}
          autocomplete="email"
          required
          aria-invalid={fieldErrors.from ? "true" : undefined}
          aria-describedby={fieldErrors.from ? "fromError" : undefined}
          oninput={() => clearFieldError("from")}
        />
        {#if fieldErrors.from}
          <span id="fromError" class="fieldError">{fieldErrors.from}</span>
        {/if}
      </label>

      <label>
        <span>{ui.subject}</span>
        <input
          bind:value={subject}
          type="text"
          placeholder={ui.subjectPlaceholder}
          required
          maxlength="180"
          aria-invalid={fieldErrors.subject ? "true" : undefined}
          aria-describedby={fieldErrors.subject ? "subjectError" : undefined}
          oninput={() => clearFieldError("subject")}
        />
        {#if fieldErrors.subject}
          <span id="subjectError" class="fieldError">{fieldErrors.subject}</span>
        {/if}
      </label>

      <label>
        <span>{ui.message}</span>
        <textarea
          bind:value={message}
          placeholder={ui.messagePlaceholder}
          required
          maxlength="4000"
          aria-invalid={fieldErrors.message ? "true" : undefined}
          aria-describedby={fieldErrors.message ? "messageError" : undefined}
          oninput={() => clearFieldError("message")}
        ></textarea>
        {#if fieldErrors.message}
          <span id="messageError" class="fieldError">{fieldErrors.message}</span>
        {/if}
      </label>

      <div class="actions">
        <button class="btn" type="button" disabled={sending} aria-busy={sending} onclick={submit}>
          {sending ? ui.sending : ui.send}
        </button>
        <a class="link" href={linkedInUrl} target="_blank" rel="noreferrer">linkedin</a>
        <a class="link" href={githubUrl} target="_blank" rel="noreferrer">github</a>
        <a class="link" href={portfolioUrl} target="_blank" rel="noreferrer">portfolio</a>
      </div>

      {#if status}
        <div class={"status " + status.kind} role={status.kind === "err" ? "alert" : "status"} aria-live="polite">
          <span>{status.msg}</span>
          <button type="button" class="statusClose" onclick={clearStatus}>{ui.dismiss}</button>
        </div>
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
  .formHint {
    margin: 0;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.66);
  }

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
  input[aria-invalid="true"], textarea[aria-invalid="true"] {
    border-color: rgba(255, 140, 140, 0.66);
  }

  .fieldError {
    display: block;
    margin-top: 6px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(255, 140, 140, 0.92);
  }

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

  .status {
    margin-top: 8px;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 10px;
    border-radius: 10px;
    border: 1px solid rgba(255,255,255,0.14);
    padding: 10px 12px;
    background: rgba(12,14,18,0.22);
  }
  .status.ok {
    color: rgba(140, 240, 170, 0.9);
    border-color: rgba(140, 240, 170, 0.34);
  }
  .status.err {
    color: rgba(255, 140, 140, 0.95);
    border-color: rgba(255, 140, 140, 0.4);
  }
  .status.info {
    color: rgba(245,245,245,0.8);
    border-color: rgba(245,245,245,0.24);
  }

  .statusClose {
    border: 1px solid rgba(245,245,245,0.20);
    border-radius: 999px;
    background: rgba(255,255,255,0.06);
    color: inherit;
    font-family: inherit;
    font-size: 11px;
    padding: 5px 9px;
    cursor: pointer;
    flex-shrink: 0;
  }

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
