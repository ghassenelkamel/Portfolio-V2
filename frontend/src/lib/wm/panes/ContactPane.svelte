<script lang="ts">
  let from = $state("");
  let subject = $state("");
  let message = $state("");
  let status = $state<string | null>(null);
  let sending = $state(false);

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
      if (!r.ok) status = j?.error ?? `failed (${r.status})`;
      else status = "sent";
    } catch (e) {
      status = e instanceof Error ? e.message : "unknown error";
    } finally {
      sending = false;
      setTimeout(() => (status = null), 2200);
    }
  }
</script>

<div class="pane">
  <h2>Contact</h2>

  <div class="form">
    <label>
      <span>From</span>
      <input bind:value={from} type="email" placeholder="your@email.com" />
    </label>

    <label>
      <span>Subject</span>
      <input bind:value={subject} type="text" placeholder="subject" />
    </label>

    <label>
      <span>Message</span>
      <textarea bind:value={message} placeholder="write your message…"></textarea>
    </label>

    <button type="button" class="btn" disabled={sending} onclick={submit}>
      {sending ? "sending…" : "send"}
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
