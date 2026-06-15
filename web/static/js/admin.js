(function () {
  const csrf = document.querySelector('meta[name="csrf-token"]')?.content || '';
  const page = document.getElementById('admin-page');
  const pageData = page ? page.dataset : document.body.dataset;

  function apiHeaders(json) {
    const h = { 'X-CSRF-Token': csrf };
    if (json) h['Content-Type'] = 'application/json';
    return h;
  }

  function escapeHtml(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function formatKopecksAsRubles(kopecks) {
    const total = Number(kopecks) || 0;
    const rub = Math.floor(total / 100);
    const kop = total % 100;
    if (kop === 0) return String(rub);
    return `${rub},${String(kop).padStart(2, '0')}`;
  }

  function formatKopecksDisplay(kopecks) {
    const total = Number(kopecks) || 0;
    const rub = Math.floor(total / 100);
    const kop = total % 100;
    if (kop === 0) return `${rub} ₽`;
    return `${rub},${String(kop).padStart(2, '0')} ₽`;
  }

  function parseRublesToKopecks(raw) {
    const normalized = String(raw).trim().replace(/\s/g, '').replace(',', '.');
    if (!normalized) return 0;
    const match = normalized.match(/^(\d+)(?:\.(\d{0,2}))?$/);
    if (!match) return NaN;
    const rub = parseInt(match[1], 10) || 0;
    const frac = (match[2] || '').padEnd(2, '0').slice(0, 2);
    const kop = frac ? parseInt(frac, 10) || 0 : 0;
    return rub * 100 + kop;
  }

  /* ——— Image upload (drag & drop) ——— */
  function initDropzones(root) {
    root.querySelectorAll('.upload-dropzone:not([data-bound])').forEach((zone) => {
      zone.dataset.bound = '1';
      const inputName = zone.dataset.targetInput;
      const fileInput = zone.querySelector('.upload-file-input');
      const hidden = root.querySelector(`input[name="${inputName}"]`) || zone.parentElement.querySelector('.upload-hidden-url');
      const previewWrap = zone.parentElement.querySelector('.upload-preview-wrap');
      const previewImg = previewWrap?.querySelector('.upload-preview');
      const errBox = zone.parentElement.querySelector('.upload-error');

      if (!hidden || !fileInput) return;

      function setPreview(url) {
        if (!url || !previewWrap || !previewImg) return;
        previewImg.src = url;
        previewWrap.classList.remove('hidden');
      }

      function setError(msg) {
        if (!errBox) return;
        if (msg) {
          errBox.textContent = msg;
          errBox.classList.remove('hidden');
        } else {
          errBox.classList.add('hidden');
          errBox.textContent = '';
        }
      }

      async function uploadFile(file) {
        if (!file.type.startsWith('image/')) {
          setError('Выберите файл изображения');
          return;
        }
        setError('');
        zone.classList.add('uploading');
        const fd = new FormData();
        fd.append('file', file);
        try {
          const r = await fetch('/api/upload', {
            method: 'POST',
            headers: { 'X-CSRF-Token': csrf },
            credentials: 'same-origin',
            body: fd,
          });
          const data = await r.json().catch(() => ({}));
          if (!r.ok) {
            setError(data.error || 'Ошибка загрузки');
            return;
          }
          const url = data.data?.url;
          if (url) {
            hidden.value = url;
            setPreview(url);
            const rootForm = zone.closest('form');
            if (rootForm) updateAdminEventPreview(rootForm);
          }
        } catch {
          setError('Сеть недоступна');
        } finally {
          zone.classList.remove('uploading');
        }
      }

      zone.addEventListener('click', () => fileInput.click());
      zone.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          fileInput.click();
        }
      });
      fileInput.addEventListener('change', () => {
        if (fileInput.files[0]) uploadFile(fileInput.files[0]);
      });
      zone.addEventListener('dragover', (e) => {
        e.preventDefault();
        zone.classList.add('dragover');
      });
      zone.addEventListener('dragleave', () => zone.classList.remove('dragover'));
      zone.addEventListener('drop', (e) => {
        e.preventDefault();
        zone.classList.remove('dragover');
        const file = e.dataTransfer?.files?.[0];
        if (file) uploadFile(file);
      });

      if (hidden.value) setPreview(hidden.value);
    });
  }

  function formatEventPreviewDay(dateStr) {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    if (Number.isNaN(d.getTime())) return '—';
    return String(d.getDate()).padStart(2, '0');
  }

  function formatEventPreviewMeta(dateStr) {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    if (Number.isNaN(d.getTime())) return '';
    const months = ['ЯНВ', 'ФЕВ', 'МАР', 'АПР', 'МАЙ', 'ИЮН', 'ИЮЛ', 'АВГ', 'СЕН', 'ОКТ', 'НОЯ', 'ДЕК'];
    const days = ['ВС', 'ПН', 'ВТ', 'СР', 'ЧТ', 'ПТ', 'СБ'];
    return `${months[d.getMonth()]} · ${days[d.getDay()]}`;
  }

  function updateAdminEventPreview(form) {
    const bg = document.getElementById('admin-event-preview-bg');
    const day = document.getElementById('admin-event-preview-day');
    const meta = document.getElementById('admin-event-preview-meta');
    const city = document.getElementById('admin-event-preview-city');
    const title = document.getElementById('admin-event-preview-title');
    if (!bg || !form) return;

    const dateVal = form.querySelector('[name="date"]')?.value || '';
    const cityVal = (form.querySelector('[name="city"]')?.value || '').trim();
    const titleVal = (form.querySelector('[name="title"]')?.value || '').trim();
    const posterVal = (form.querySelector('[name="poster_image_url"]')?.value || '').trim();

    if (day) day.textContent = formatEventPreviewDay(dateVal);
    if (meta) meta.textContent = formatEventPreviewMeta(dateVal);
    if (city) city.textContent = cityVal ? cityVal.toUpperCase() : 'ГОРОД';
    if (title) title.textContent = titleVal || 'Название события';

    if (posterVal) {
      bg.style.backgroundImage = `linear-gradient(180deg, rgba(10, 10, 11, 0.2) 0%, rgba(10, 10, 11, 0.85) 100%), radial-gradient(ellipse at 85% 15%, rgba(217, 173, 42, 0.15) 0%, transparent 45%), url("${posterVal.replace(/"/g, '\\"')}")`;
    } else {
      bg.style.backgroundImage = '';
    }
  }

  function setPosterURL(form, url) {
    const hidden = form.querySelector('[name="poster_image_url"]');
    const manual = document.getElementById('poster-url-manual');
    const previewWrap = hidden?.closest('.form-field')?.querySelector('.upload-preview-wrap');
    const previewImg = previewWrap?.querySelector('.upload-preview');
    if (hidden) hidden.value = url || '';
    if (manual) manual.value = url || '';
    if (url && previewWrap && previewImg) {
      previewImg.src = url;
      previewWrap.classList.remove('hidden');
    }
    updateAdminEventPreview(form);
  }

  function initEventPosterPreview(form) {
    const previewFields = ['title', 'city', 'date', 'poster_image_url'];
    previewFields.forEach((name) => {
      const el = form.querySelector(`[name="${name}"]`);
      if (el) el.addEventListener('input', () => updateAdminEventPreview(form));
      if (el && name === 'date') el.addEventListener('change', () => updateAdminEventPreview(form));
    });

    const manual = document.getElementById('poster-url-manual');
    if (manual) {
      manual.addEventListener('change', () => {
        setPosterURL(form, manual.value.trim());
      });
    }

    const previewBtn = document.getElementById('ticket-url-preview-btn');
    const previewStatus = document.getElementById('ticket-url-preview-status');
    if (previewBtn) {
      previewBtn.addEventListener('click', () => fetchTicketPagePreview(form, previewStatus));
    }

    const ticketURL = form.querySelector('[name="ticket_url"]');
    if (ticketURL) {
      ticketURL.addEventListener('blur', () => {
        const url = ticketURL.value.trim();
        if (!url || form.querySelector('[name="poster_image_url"]')?.value) return;
        fetchTicketPagePreview(form, previewStatus, true);
      });
    }

    updateAdminEventPreview(form);
  }

  async function fetchTicketPagePreview(form, statusEl, silent) {
    const ticketURL = form.querySelector('[name="ticket_url"]')?.value?.trim();
    if (!ticketURL) {
      if (!silent && statusEl) statusEl.textContent = 'Сначала укажите ссылку на билеты';
      return;
    }
    if (statusEl) statusEl.textContent = 'Загрузка предпросмотра…';
    try {
      const r = await fetch(`/api/events/preview-ticket?url=${encodeURIComponent(ticketURL)}`, {
        headers: apiHeaders(false),
        credentials: 'same-origin',
      });
      const payload = await r.json().catch(() => ({}));
      const data = payload.data || {};
      if (!r.ok) {
        if (statusEl) statusEl.textContent = payload.error || 'Не удалось получить данные со страницы';
        return;
      }
      if (data.poster_image_url) setPosterURL(form, data.poster_image_url);
      const titleEl = form.querySelector('[name="title"]');
      const descEl = form.querySelector('[name="description"]');
      if (data.title && titleEl && !titleEl.value.trim()) titleEl.value = data.title;
      if (data.description && descEl && !descEl.value.trim()) descEl.value = data.description;
      updateAdminEventPreview(form);
      if (statusEl) {
        statusEl.textContent = data.poster_image_url
          ? 'Данные подтянуты — проверьте и при необходимости исправьте'
          : 'Текст подтянут, постер на странице не найден — загрузите вручную';
      }
    } catch {
      if (statusEl) statusEl.textContent = 'Сеть недоступна';
    }
  }

  function showFieldErrors(errBox, errs) {
    errBox.classList.remove('hidden');
    errBox.innerHTML = Object.entries(errs)
      .map(([k, v]) => `<div><strong>${escapeHtml(k)}</strong>: ${escapeHtml(String(v))}</div>`)
      .join('');
  }

  const tableBody = document.getElementById('admin-table-body');
  if (tableBody && pageData.pageMode === 'list' && pageData.apiPath) {
    fetch(pageData.apiPath, { headers: apiHeaders(false), credentials: 'same-origin' })
      .then((r) => {
        if (!r.ok) throw new Error('http ' + r.status);
        return r.json();
      })
      .then((payload) => {
        const cols = (pageData.columns || '').split(',').filter(Boolean);
        const slug = pageData.modelSlug;
        tableBody.innerHTML = (payload.data || [])
          .map((row) => renderRow(slug, row, cols))
          .join('') || '<tr><td colspan="10" class="px-4 py-6 text-slate-500">Нет записей</td></tr>';
      })
      .catch(() => {
        tableBody.innerHTML = '<tr><td colspan="10" class="px-4 py-6 text-red-600">Ошибка загрузки</td></tr>';
      });
  }

  function renderRow(slug, row, cols) {
    const cells = cols
      .map((key) => {
        let val = row[key];
        if (key === 'date' && val) val = new Date(val).toLocaleString('ru-RU');
        if (key === 'price' && val != null) val = formatKopecksDisplay(val);
        if (key === 'ticket_source' && val) {
          const labels = { manual: 'вручную', timepad: 'Timepad', ticketscloud: 'TicketsCloud' };
          val = labels[val] || val;
        }
        if (key === 'image_url' && val) {
          return `<td class="px-4 py-3"><img src="${escapeHtml(String(val))}" alt="" class="h-10 w-10 object-cover rounded border"></td>`;
        }
        if (key === 'url' && val) val = String(val).slice(0, 40) + '…';
        return `<td class="px-4 py-3">${escapeHtml(String(val ?? ''))}</td>`;
      })
      .join('');
    return `<tr class="hover:bg-slate-50">
      <td class="px-4 py-3 text-slate-500">${row.id}</td>
      ${cells}
      <td class="px-4 py-3 text-right admin-row-actions whitespace-nowrap">
        <a href="/admin/${slug}/${row.id}/edit" class="text-blue-700 hover:underline text-sm">Edit</a>
        <button type="button" data-delete-id="${row.id}" class="text-red-700 hover:underline text-sm">Delete</button>
      </td>
    </tr>`;
  }

  tableBody?.addEventListener('click', function (e) {
    const btn = e.target.closest('[data-delete-id]');
    if (!btn) return;
    const id = btn.dataset.deleteId;
    if (!confirm('Удалить запись #' + id + '?')) return;
    fetch(`${pageData.apiPath}/${id}`, {
      method: 'DELETE',
      headers: apiHeaders(false),
      credentials: 'same-origin',
    }).then(() => location.reload());
  });

  const settingsForm = document.getElementById('settings-form');
  if (settingsForm && pageData.pageMode === 'settings') {
    fetch('/api/settings', { headers: apiHeaders(false), credentials: 'same-origin' })
      .then((r) => r.json())
      .then((payload) => {
        const data = payload.data || {};
        settingsForm.querySelectorAll('[name]').forEach((el) => {
          if (data[el.name] != null) el.value = data[el.name];
        });
        initDropzones(settingsForm);
      });

    settingsForm.addEventListener('submit', function (e) {
      e.preventDefault();
      const body = {};
      settingsForm.querySelectorAll('[name]').forEach((el) => {
        body[el.name] = el.value;
      });
      fetch('/api/settings', {
        method: 'PUT',
        headers: apiHeaders(true),
        credentials: 'same-origin',
        body: JSON.stringify(body),
      }).then(async (r) => {
        const data = await r.json().catch(() => ({}));
        const errBox = document.getElementById('form-errors');
        if (!r.ok) {
          const errs = data.errors || { error: data.error || 'Ошибка сохранения' };
          showFieldErrors(errBox, errs);
          return;
        }
        errBox.innerHTML = '<div class="text-green-800">Сохранено</div>';
        errBox.classList.remove('hidden');
      });
    });
  }

  const form = document.getElementById('admin-form');
  if (form && pageData.apiPath && pageData.pageMode !== 'settings' && pageData.pageMode !== 'account') {
    const recordId = parseInt(pageData.recordId || '0', 10);
    if (pageData.pageMode === 'event-form') {
      initTicketImporter(form);
      initEventPosterPreview(form);
    }
    if (recordId > 0) {
      fetch(`${pageData.apiPath}/${recordId}`, {
        headers: apiHeaders(false),
        credentials: 'same-origin',
      })
        .then((r) => r.json())
        .then((payload) => {
          const data = payload.data || {};
          form.querySelectorAll('[name]').forEach((el) => {
            const name = el.name;
            let val = data[name];
            if (name === 'date' && val) {
              val = new Date(val).toISOString().slice(0, 16);
            }
            if (name === 'price' && val != null && el.closest('.form-field')?.dataset.fieldType === 'price_rub') {
              val = formatKopecksAsRubles(val);
            }
            if (val != null) el.value = val;
          });
          if (data.ticket_source) {
            const src = form.querySelector('[name="ticket_source"]');
            if (src) src.value = data.ticket_source;
          }
          if (data.external_id) {
            const ext = form.querySelector('[name="external_id"]');
            if (ext) ext.value = data.external_id;
          }
          initDropzones(form);
          if (pageData.pageMode === 'event-form') {
            const manual = document.getElementById('poster-url-manual');
            if (manual && data.poster_image_url) manual.value = data.poster_image_url;
            updateAdminEventPreview(form);
          }
        });
    } else {
      initDropzones(form);
    }

    form.addEventListener('submit', function (e) {
      e.preventDefault();
      const errBox = document.getElementById('form-errors');
      const body = {};
      for (const el of form.querySelectorAll('[name]')) {
        const fieldType = el.closest('.form-field')?.dataset.fieldType;
        let v = el.value;
        if (fieldType === 'price_rub') {
          v = parseRublesToKopecks(v);
          if (Number.isNaN(v)) {
            errBox.classList.remove('hidden');
            errBox.innerHTML = '<div><strong>price</strong>: укажите цену в рублях, например 1990 или 1990,50</div>';
            return;
          }
        } else if (el.type === 'number') {
          v = parseInt(v, 10) || 0;
        }
        body[el.name] = v;
      }
      const srcEl = form.querySelector('[name="ticket_source"]');
      const extEl = form.querySelector('[name="external_id"]');
      if (srcEl) body.ticket_source = srcEl.value || 'manual';
      if (extEl) body.external_id = extEl.value || '';
      const method = recordId > 0 ? 'PUT' : 'POST';
      const url = recordId > 0 ? `${pageData.apiPath}/${recordId}` : pageData.apiPath;
      fetch(url, {
        method,
        headers: apiHeaders(true),
        credentials: 'same-origin',
        body: JSON.stringify(body),
      })
        .then(async (r) => {
          const data = await r.json().catch(() => ({}));
          if (!r.ok) {
            errBox.classList.remove('hidden');
            const errs = data.errors || { error: data.error || 'Ошибка сохранения' };
            showFieldErrors(errBox, errs);
            return;
          }
          window.location.href = `/admin/${pageData.modelSlug}`;
        });
    });
  }

  function initTicketImporter(form) {
    const select = document.getElementById('ticket-source-select');
    const loadBtn = document.getElementById('ticket-source-load');
    const list = document.getElementById('ticket-source-list');
    const status = document.getElementById('ticket-source-status');
    if (!select || !loadBtn || !list) return;

    fetch('/api/ticket-catalog/providers', { headers: apiHeaders(false), credentials: 'same-origin' })
      .then((r) => r.json())
      .then((payload) => {
        const providers = (payload.data || []).filter((p) => p.configured);
        select.innerHTML = providers.length
          ? providers.map((p) => `<option value="${escapeHtml(p.id)}">${escapeHtml(p.name)}</option>`).join('')
          : '<option value="">Сначала настройте агрегатор в Settings</option>';
        if (providers.length) loadTicketEvents();
      })
      .catch(() => {
        select.innerHTML = '<option value="">Ошибка загрузки</option>';
      });

    let cachedItems = [];

    loadBtn.addEventListener('click', loadTicketEvents);

    function loadTicketEvents() {
      const source = select.value;
      if (!source) return;
      status.textContent = 'Загрузка…';
      list.innerHTML = '';
      fetch(`/api/ticket-catalog/${encodeURIComponent(source)}/events`, {
        headers: apiHeaders(false),
        credentials: 'same-origin',
      })
        .then(async (r) => {
          const payload = await r.json().catch(() => ({}));
          if (!r.ok) {
            status.textContent = payload.error || 'Не удалось загрузить события';
            return;
          }
          cachedItems = payload.data || [];
          status.textContent = cachedItems.length ? `Найдено: ${cachedItems.length}` : 'Событий не найдено';
          list.innerHTML = cachedItems
            .map((item, index) => {
              const date = item.date ? new Date(item.date).toLocaleString('ru-RU') : '';
              const added = item.already_added ? ' · уже в афише' : '';
              const thumb = item.poster_image_url
                ? `<img src="${escapeHtml(item.poster_image_url)}" alt="" class="ticket-pick-thumb">`
                : '<div class="ticket-pick-thumb"></div>';
              return `<button type="button" class="ticket-pick w-full text-left px-3 py-2 hover:bg-slate-50 text-sm" data-index="${index}">
              <div class="ticket-pick-row">${thumb}<div>
              <div class="font-medium text-slate-900">${escapeHtml(item.title || '')}${added}</div>
              <div class="text-slate-500 text-xs">${escapeHtml(date)} · ${escapeHtml(item.city || '')}</div>
              </div></div>
            </button>`;
            })
            .join('');
        })
        .catch(() => {
          status.textContent = 'Не удалось загрузить события';
        });
    }

    list.addEventListener('click', (e) => {
      const btn = e.target.closest('.ticket-pick');
      if (!btn) return;
      const item = cachedItems[Number(btn.dataset.index)];
      if (!item) return;
      fillEventFromTicket(form, item);
      status.textContent = 'Поля заполнены — проверьте и нажмите «Сохранить»';
    });
  }

  function fillEventFromTicket(form, item) {
    const set = (name, value) => {
      const el = form.querySelector(`[name="${name}"]`);
      if (el && value != null) el.value = value;
    };
    set('title', item.title || '');
    set('city', item.city || '');
    set('description', item.description || '');
    set('ticket_url', item.ticket_url || '');
    if (item.date) {
      set('date', new Date(item.date).toISOString().slice(0, 16));
    }
    set('ticket_source', item.source || 'manual');
    set('external_id', item.external_id || '');
    if (item.poster_image_url) {
      setPosterURL(form, item.poster_image_url);
    } else {
      updateAdminEventPreview(form);
    }
  }

  const accountForm = document.getElementById('account-form');
  if (accountForm && pageData.pageMode === 'account') {
    accountForm.addEventListener('submit', function (e) {
      e.preventDefault();
      const errBox = document.getElementById('form-errors');
      const body = {
        current_password: accountForm.querySelector('[name="current_password"]').value,
        new_password: accountForm.querySelector('[name="new_password"]').value,
      };
      fetch('/api/account/password', {
        method: 'PUT',
        headers: apiHeaders(true),
        credentials: 'same-origin',
        body: JSON.stringify(body),
      }).then(async (r) => {
        const data = await r.json().catch(() => ({}));
        if (!r.ok) {
          const errs = data.errors || { error: data.error || 'Ошибка смены пароля' };
          showFieldErrors(errBox, errs);
          return;
        }
        window.location.href = '/admin/login';
      });
    });
  }
})();
