"use strict";

// Catalog admin panel. Talks to the same server's JSON API:
//   POST /login, DELETE /signout, GET/POST /sarees, PUT/DELETE /sarees/:id, POST /uploads

const TOKEN_KEY = "saree-admin-token";
const USER_KEY = "saree-admin-user";
const MAX_PHOTO_EDGE = 1600; // px; phone photos are shrunk before upload

const $ = (id) => document.getElementById(id);

// sessionStorage can throw (private mode, blocked storage); treat that as empty.
const storage = {
  get(key) {
    try { return window.sessionStorage.getItem(key); } catch { return null; }
  },
  set(key, value) {
    try {
      if (value == null) window.sessionStorage.removeItem(key);
      else window.sessionStorage.setItem(key, value);
    } catch { /* not persisted */ }
  },
};

const state = {
  token: storage.get(TOKEN_KEY),
  user: parseJson(storage.get(USER_KEY)),
  sarees: [],
  editingId: null, // null = adding a new saree
  photos: [],
  dirty: false,
  uploads: 0,
  resumeEdit: false,
};

// ---------- helpers ----------
class SessionError extends Error {}

function parseJson(text) {
  try { return JSON.parse(text || "null"); } catch { return null; }
}

function show(view) {
  for (const id of ["login-view", "list-view", "edit-view"]) $(id).hidden = id !== view;
  $("session-actions").hidden = view === "login-view";
  window.scrollTo(0, 0);
}

let toastTimer;
function toast(message, bad = false) {
  const el = $("toast");
  el.textContent = message;
  el.classList.toggle("bad", bad);
  el.hidden = false;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => { el.hidden = true; }, bad ? 6000 : 3000);
}

function showError(id, message) {
  const el = $(id);
  el.textContent = message || "";
  el.hidden = !message;
}

const rupees = new Intl.NumberFormat("en-IN", { style: "currency", currency: "INR", maximumFractionDigits: 0 });

function safeImageUrl(url) {
  return typeof url === "string" && /^(https?:\/\/|\/)/.test(url) ? url : "";
}

// ---------- API ----------

async function api(path, { method = "GET", body, form } = {}) {
  const headers = {};
  if (state.token) headers.Authorization = `Bearer ${state.token}`;
  if (body !== undefined) headers["Content-Type"] = "application/json";

  const response = await fetch(path, {
    method,
    headers,
    body: form ?? (body !== undefined ? JSON.stringify(body) : undefined),
  });
  const data = await response.json().catch(() => null);

  if (response.status === 401) {
    requireLogin("Your session expired. Sign in again to continue — nothing you typed is lost.");
    throw new SessionError("session expired");
  }
  if (response.status === 403) {
    requireLogin("This account isn't an admin. Sign in with an email listed in ADMIN_EMAILS.");
    throw new SessionError("not an admin");
  }
  if (!response.ok) {
    throw new Error((data && data.error) || `Request failed (${response.status})`);
  }
  return data;
}

function report(error) {
  if (error instanceof SessionError) return;
  console.error(error);
  toast(error.message || "Something went wrong", true);
}

// ---------- session ----------

function setSession(token, user) {
  state.token = token;
  state.user = user;
  storage.set(TOKEN_KEY, token);
  storage.set(USER_KEY, user ? JSON.stringify(user) : null);
  $("who").textContent = user ? user.email : "";
}

function requireLogin(reason) {
  state.resumeEdit = !$("edit-view").hidden;
  setSession(null, null);
  $("login-reason").textContent = reason || "Use an admin account to add and edit sarees.";
  showError("login-error", "");
  show("login-view");
  $("login-form").email.focus();
}

$("login-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = event.currentTarget;
  const button = form.querySelector("button[type=submit]");
  showError("login-error", "");
  button.disabled = true;
  try {
    const response = await fetch("/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: form.email.value.trim(), password: form.password.value }),
    });
    const data = await response.json().catch(() => null);
    if (!response.ok) throw new Error((data && data.error) || "Sign in failed");
    if (!data.user || data.user.role !== "admin") {
      throw new Error("This account isn't an admin. Ask the site owner to add your email to ADMIN_EMAILS.");
    }
    setSession(data.tokens.access_token, data.user);
    form.password.value = "";
    if (state.resumeEdit) {
      state.resumeEdit = false;
      show("edit-view");
      toast("Signed in — you can save now");
    } else {
      await loadList();
    }
  } catch (error) {
    showError("login-error", error.message);
  } finally {
    button.disabled = false;
  }
});

$("logout").addEventListener("click", async () => {
  if (state.dirty && !confirm("Discard your unsaved changes and sign out?")) return;
  try { await api("/signout", { method: "DELETE" }); } catch { /* signing out locally anyway */ }
  state.dirty = false;
  state.resumeEdit = false;
  setSession(null, null);
  requireLogin("You're signed out.");
  state.resumeEdit = false;
});

// ---------- list ----------

async function loadList() {
  try {
    state.sarees = (await api("/sarees")) || [];
    renderList();
    show("list-view");
  } catch (error) {
    report(error);
  }
}

function renderList() {
  const query = $("search").value.trim().toLowerCase();
  const list = $("saree-list");
  list.replaceChildren();

  const matches = state.sarees.filter((s) =>
    !query || [s.name, s.category, s.fabric_type, s.color, s.design_pattern, s.description]
      .some((v) => (v || "").toLowerCase().includes(query)));

  const total = state.sarees.length;
  const live = state.sarees.filter((s) => s.availability).length;
  $("count").textContent = total ? `${total} saree${total === 1 ? "" : "s"} · ${live} on the website` : "";
  $("empty").hidden = total > 0;

  const template = $("saree-item");
  for (const saree of matches) {
    const item = template.content.firstElementChild.cloneNode(true);
    const cover = safeImageUrl((saree.images_url || [])[0]);
    if (cover) item.querySelector(".thumb").src = cover;
    item.querySelector(".name").textContent = saree.name || "(untitled)";
    item.querySelector(".meta").textContent =
      [saree.category, saree.fabric_type, saree.color].filter(Boolean).join(" · ") || "No details yet";

    const meta2 = item.querySelector(".meta2");
    const badge = document.createElement("span");
    badge.className = `badge ${saree.availability ? "badge-on" : "badge-off"}`;
    badge.textContent = saree.availability ? "On website" : "Hidden";
    meta2.append(`${rupees.format(saree.price || 0)} · ${saree.stock || 0} in stock `, badge);

    item.querySelector('[data-action="edit"]').addEventListener("click", () => openEditor(saree));
    item.querySelector('[data-action="delete"]').addEventListener("click", () => deleteSaree(saree));
    list.append(item);
  }

  fillDatalist("category-options", state.sarees.map((s) => s.category));
  fillDatalist("fabric-options", state.sarees.map((s) => s.fabric_type));
}

function fillDatalist(id, values) {
  const options = [...new Set(values.filter(Boolean))].sort().map((value) => {
    const option = document.createElement("option");
    option.value = value;
    return option;
  });
  $(id).replaceChildren(...options);
}

$("search").addEventListener("input", renderList);
$("add-saree").addEventListener("click", () => openEditor(null));

async function deleteSaree(saree) {
  if (!confirm(`Delete "${saree.name}"? It will disappear from the website.`)) return;
  try {
    await api(`/sarees/${encodeURIComponent(saree.id)}`, { method: "DELETE" });
    state.sarees = state.sarees.filter((s) => s.id !== saree.id);
    renderList();
    toast("Saree deleted");
  } catch (error) {
    report(error);
  }
}

// ---------- editor ----------

const form = $("saree-form");

function openEditor(saree) {
  state.editingId = saree ? saree.id : null;
  $("edit-title").textContent = saree ? "Edit saree" : "Add saree";
  $("save").textContent = saree ? "Save changes" : "Save saree";
  form.reset();
  showError("form-error", "");
  $("upload-status").textContent = "";

  const s = saree || {};
  for (const field of ["name", "description", "category", "fabric_type", "color", "design_pattern"]) {
    form[field].value = s[field] || "";
  }
  form.price.value = saree ? s.price ?? 0 : "";
  form.stock.value = saree ? s.stock ?? 0 : "";
  form.availability.checked = saree ? Boolean(s.availability) : true;
  form.saree_length.value = s.saree_dimensions?.length || "";
  form.saree_breadth.value = s.saree_dimensions?.breadth || "";
  form.blouse_length.value = s.blouse_dimensions?.length || "";
  form.blouse_breadth.value = s.blouse_dimensions?.breadth || "";

  state.photos = (s.images_url || []).filter(safeImageUrl);
  renderPhotos();
  state.dirty = false;
  show("edit-view");
  form.name.focus();
}

function closeEditor() {
  if (state.uploads > 0 && !confirm("Photos are still uploading. Leave anyway?")) return;
  if (state.dirty && !confirm("Discard your unsaved changes?")) return;
  state.dirty = false;
  loadList();
}

for (const button of document.querySelectorAll('[data-action="cancel"]')) {
  button.addEventListener("click", closeEditor);
}
form.addEventListener("input", () => { state.dirty = true; });
window.addEventListener("beforeunload", (event) => {
  if (state.dirty || state.uploads > 0) event.preventDefault();
});

function renderPhotos() {
  const list = $("photos");
  list.replaceChildren();
  const template = $("photo-item");
  state.photos.forEach((url, index) => {
    const item = template.content.firstElementChild.cloneNode(true);
    item.querySelector("img").src = url;
    item.querySelector('[data-action="left"]').addEventListener("click", () => movePhoto(index, -1));
    item.querySelector('[data-action="right"]').addEventListener("click", () => movePhoto(index, 1));
    item.querySelector('[data-action="remove"]').addEventListener("click", () => {
      state.photos.splice(index, 1);
      state.dirty = true;
      renderPhotos();
    });
    list.append(item);
  });
}

function movePhoto(index, delta) {
  const target = index + delta;
  if (target < 0 || target >= state.photos.length) return;
  [state.photos[index], state.photos[target]] = [state.photos[target], state.photos[index]];
  state.dirty = true;
  renderPhotos();
}

// Shrinks large photos so the website loads quickly. Falls back to the
// original file if the browser can't decode it.
async function shrink(file) {
  if (file.type === "image/gif" || !window.createImageBitmap) return file;
  try {
    const bitmap = await createImageBitmap(file);
    const scale = Math.min(1, MAX_PHOTO_EDGE / Math.max(bitmap.width, bitmap.height));
    if (scale === 1 && file.size < 1.5 * 1024 * 1024) return file;
    const canvas = document.createElement("canvas");
    canvas.width = Math.round(bitmap.width * scale);
    canvas.height = Math.round(bitmap.height * scale);
    canvas.getContext("2d").drawImage(bitmap, 0, 0, canvas.width, canvas.height);
    const blob = await new Promise((resolve) => canvas.toBlob(resolve, "image/jpeg", 0.85));
    return blob || file;
  } catch {
    return file;
  }
}

async function uploadFiles(files) {
  const images = [...files].filter((f) => f.type.startsWith("image/"));
  if (!images.length) return;
  state.uploads += images.length;
  $("save").disabled = true;
  let done = 0;
  let failed = 0;
  const status = () => {
    $("upload-status").textContent = state.uploads > 0
      ? `Uploading ${done + 1} of ${images.length}…`
      : failed ? `${failed} photo${failed === 1 ? "" : "s"} could not be uploaded.` : "Photos added.";
  };
  status();

  for (const file of images) {
    try {
      const body = new FormData();
      body.append("file", await shrink(file), file.name.replace(/\.\w+$/, "") + ".jpg");
      const { url } = await api("/uploads", { method: "POST", form: body });
      state.photos.push(url);
      state.dirty = true;
      renderPhotos();
    } catch (error) {
      failed += 1;
      report(error);
      if (error instanceof SessionError) {
        state.uploads -= images.length - done;
        break;
      }
    }
    done += 1;
    state.uploads -= 1;
    if (state.uploads > 0) status();
  }
  status();
  $("save").disabled = state.uploads > 0;
}

$("photo-input").addEventListener("change", (event) => {
  uploadFiles(event.target.files);
  event.target.value = "";
});

const dropzone = $("dropzone");
dropzone.addEventListener("dragover", (event) => { event.preventDefault(); dropzone.classList.add("dragging"); });
dropzone.addEventListener("dragleave", () => dropzone.classList.remove("dragging"));
dropzone.addEventListener("drop", (event) => {
  event.preventDefault();
  dropzone.classList.remove("dragging");
  uploadFiles(event.dataTransfer.files);
});

function wholeNumber(value) {
  const n = Number.parseInt(value, 10);
  return Number.isFinite(n) && n > 0 ? n : 0;
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  showError("form-error", "");
  if (!form.name.value.trim()) {
    showError("form-error", "Give the saree a name.");
    form.name.focus();
    return;
  }
  if (state.uploads > 0) {
    showError("form-error", "Wait for the photos to finish uploading.");
    return;
  }

  const saree = {
    name: form.name.value.trim(),
    description: form.description.value.trim(),
    category: form.category.value.trim(),
    fabric_type: form.fabric_type.value.trim(),
    color: form.color.value.trim(),
    design_pattern: form.design_pattern.value.trim(),
    price: wholeNumber(form.price.value),
    stock: wholeNumber(form.stock.value),
    availability: form.availability.checked,
    images_url: state.photos,
    saree_dimensions: { length: form.saree_length.value.trim(), breadth: form.saree_breadth.value.trim() },
    blouse_dimensions: { length: form.blouse_length.value.trim(), breadth: form.blouse_breadth.value.trim() },
  };

  const button = $("save");
  button.disabled = true;
  try {
    if (state.editingId) {
      await api(`/sarees/${encodeURIComponent(state.editingId)}`, { method: "PUT", body: saree });
    } else {
      await api("/sarees", { method: "POST", body: saree });
    }
    state.dirty = false;
    toast(state.editingId ? "Changes saved" : "Saree added");
    await loadList();
  } catch (error) {
    if (!(error instanceof SessionError)) showError("form-error", error.message);
  } finally {
    button.disabled = false;
  }
});

// ---------- start ----------

if (state.token && state.user && state.user.role === "admin") {
  $("who").textContent = state.user.email;
  loadList();
} else {
  requireLogin();
}
