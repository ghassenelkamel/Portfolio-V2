<script lang="ts">
  type Status = { kind: 'ok' | 'err' | 'info'; msg: string } | null;

  const toEmail = "Ghassenelkamel@live.fr";

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
      setStatus({ kind: "ok", msg: "Email copied." }, 1600);
    } catch {
      setStatus({ kind: "info", msg: "Clipboard not available. Use the mailto link." }, 2200);
    }
  }

  function validate() {
    const e = from.trim();
    const s = subject.trim();
    const m = message.trim();

    if (!e || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e)) return "Enter a valid email.";
    if (!s || s.length < 3) return "Subject is too short.";
    if (!m || m.length < 10) return "Message is too short (min 10 chars).";
    if (m.length > 4000) return "Message is too long (max 4000 chars).";
    return null;
  }

  async function submit() {
    const err = validate();
    if (err) return setStatus({ kind: "err", msg: err }, 2800);

    sending = true;
    setStatus({ kind: "info", msg: "Sending..." }, 1200);

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
        return setStatus({ kind: "err", msg: `Send failed (${res.status}). ${txt.slice(0, 160)}` }, 3500);
      }

      setStatus({ kind: "ok", msg: "Message sent. Thanks!" }, 2500);
      from = "";
      subject = "";
      message = "";
    } catch (e) {
      setStatus({ kind: "err", msg: e instanceof Error ? e.message : "Network error." }, 3500);
    } finally {
      sending = false;
    }
  }
</script>

<svelte:head>
  <title>Contact</title>
</svelte:head>

<div class="wrap">
  <div class="title">
    <div class="k">mail --compose</div>
    <h1>Contact</h1>
    <p>Direct API send (no client secrets). Server can forward later if you configure it.</p>
  </div>

  <div class="card">
    <div class="row">
      <div class="kv">
        <span>To</span>
        <a href={"mailto:" + toEmail}>{toEmail}</a>
      </div>

      <button type="button" class="ghostBtn" onclick={copyEmail}>copy</button>
    </div>

    <div class="form">
      <label>
        <span>From</span>
        <input bind:value={from} type="email" placeholder="your@email.com" autocomplete="email" />
      </label>

      <label>
        <span>Subject</span>
        <input bind:value={subject} type="text" placeholder="What do you want to talk about?" />
      </label>

      <label>
        <span>Message</span>
        <textarea bind:value={message} placeholder="Write your message..."></textarea>
      </label>

      <div class="actions">
        <button class="btn" type="button" disabled={sending} onclick={submit}>
          {sending ? "sending..." : "send"}
        </button>
        <a class="link" href="https://linkedin.com/in/ghassenelkamel" target="_blank" rel="noreferrer">linkedin</a>
        <a class="link" href="https://github.com/ghassenelkamel" target="_blank" rel="noreferrer">github</a>
      </div>

      {#if status}
        <div class={"status " + status.kind}>{status.msg}</div>
      {/if}
    </div>
  </div>
</div>

<style>
  .wrap { display: grid; gap: 14px; }

  .k {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
    color: rgba(245,245,245,0.55);
  }

  h1 {
    margin: 6px 0 6px 0;
    font-family: 'HACKED', ui-monospace, Menlo, Consolas, monospace;
    color: rgba(204,102,102,0.95);
    letter-spacing: 1px;
  }

  p {
    margin: 0;
    font-family: ui-monospace, Menlo, Consolas, monospace;
    color: rgba(245,245,245,0.72);
    font-size: 13px;
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
</style>
