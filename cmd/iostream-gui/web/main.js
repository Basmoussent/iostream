// Tiny vanilla controller — no framework. The whole app is three buttons that
// hit JSON endpoints and render the response into a few status spans.

const $ = (sel) => document.querySelector(sel);

async function callAPI(path) {
  const res = await fetch(path);
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error || `HTTP ${res.status}`);
  return body;
}

function setStatus(el, message, kind) {
  el.textContent = message;
  el.classList.toggle("ok", kind === "ok");
  el.classList.toggle("err", kind === "err");
}

async function withButton(btn, statusEl, busy, fn) {
  btn.disabled = true;
  setStatus(statusEl, busy, "");
  try {
    const out = await fn();
    return out;
  } catch (err) {
    setStatus(statusEl, err.message, "err");
    throw err;
  } finally {
    btn.disabled = false;
  }
}

function renderDevices(devs) {
  const tbody = $("#devices-table tbody");
  tbody.innerHTML = "";
  if (!devs || devs.length === 0) {
    setStatus($("#devices-status"), "No iOS devices found.", "");
    return;
  }
  for (const d of devs) {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${d.udid || "—"}</td>
      <td>${d.product_name || "—"}</td>
      <td>${d.usb_info || "—"}</td>
      <td>${d.quicktime_enabled ? "✓" : ""}</td>`;
    tbody.appendChild(tr);
  }
  setStatus($("#devices-status"), `${devs.length} device(s).`, "ok");
}

$("#setup-btn").addEventListener("click", async () => {
  await withButton($("#setup-btn"), $("#setup-status"), "Launching Zadig — accept the UAC prompt…", async () => {
    const res = await callAPI("/api/setup-driver");
    if (res.status === "skipped") {
      setStatus($("#setup-status"), `Skipped: ${res.reason}`, "");
    } else {
      setStatus($("#setup-status"), "Driver setup completed.", "ok");
    }
  });
});

$("#refresh-btn").addEventListener("click", () => {
  withButton($("#refresh-btn"), $("#devices-status"), "Refreshing…", async () => {
    const devs = await callAPI("/api/devices");
    renderDevices(devs);
  });
});

$("#stream-btn").addEventListener("click", () => {
  const udid = $("#stream-udid").value.trim();
  const url = udid ? `/api/stream?udid=${encodeURIComponent(udid)}` : "/api/stream";
  withButton($("#stream-btn"), $("#stream-status"), "Starting ffplay…", async () => {
    const res = await callAPI(url);
    const logHint = res.log ? ` Log: ${res.log}` : "";
    setStatus($("#stream-status"), `Streaming — close the ffplay window to stop.${logHint}`, "ok");
  });
});

// Auto-refresh device list on first load so the user sees something useful
// without having to click anything.
window.addEventListener("DOMContentLoaded", () => $("#refresh-btn").click());
