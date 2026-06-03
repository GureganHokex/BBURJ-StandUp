(function () {
  const csrf = document.querySelector('meta[name="csrf-token"]')?.content || '';

  function apiHeaders(json) {
    const h = { 'X-CSRF-Token': csrf };
    if (json) h['Content-Type'] = 'application/json';
    return h;
  }

  function escapeHtml(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
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

  document.body.addEventListener('htmx:configRequest', function (evt) {
    evt.detail.headers['X-CSRF-Token'] = csrf;
  });

  const tableBody = document.getElementById('admin-table-body');
  if (tableBody && document.body.dataset.pageMode === 'list' && document.body.dataset.apiPath) {
    fetch(document.body.dataset.apiPath, { headers: apiHeaders(false), credentials: 'same-origin' })
      .then((r) => r.json())
      .then((payload) => {
        const cols = (document.body.dataset.columns || '').split(',').filter(Boolean);
        const slug = document.body.dataset.modelSlug;
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
        if (key === 'price' && val != null) val = (val / 100).toFixed(0) + ' ₽';
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
    fetch(`${document.body.dataset.apiPath}/${id}`, {
      method: 'DELETE',
      headers: apiHeaders(false),
      credentials: 'same-origin',
    }).then(() => location.reload());
  });

  const settingsForm = document.getElementById('settings-form');
  if (settingsForm && document.body.dataset.pageMode === 'settings') {
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
          errBox.classList.remove('hidden');
          const errs = data.errors || { error: data.error || 'Ошибка сохранения' };
          errBox.innerHTML = Object.entries(errs)
            .map(([k, v]) => `<div><strong>${k}</strong>: ${v}</div>`)
            .join('');
          return;
        }
        errBox.innerHTML = '<div class="text-green-800">Сохранено</div>';
        errBox.classList.remove('hidden');
      });
    });
  }

  const form = document.getElementById('admin-form');
  if (form && document.body.dataset.apiPath && document.body.dataset.pageMode !== 'settings') {
    const recordId = parseInt(document.body.dataset.recordId || '0', 10);
    if (recordId > 0) {
      fetch(`${document.body.dataset.apiPath}/${recordId}`, {
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
            if (val != null) el.value = val;
          });
          initDropzones(form);
        });
    } else {
      initDropzones(form);
    }

    form.addEventListener('submit', function (e) {
      e.preventDefault();
      const body = {};
      form.querySelectorAll('[name]').forEach((el) => {
        let v = el.value;
        if (el.type === 'number') v = parseInt(v, 10) || 0;
        body[el.name] = v;
      });
      const method = recordId > 0 ? 'PUT' : 'POST';
      const url = recordId > 0 ? `${document.body.dataset.apiPath}/${recordId}` : document.body.dataset.apiPath;
      fetch(url, {
        method,
        headers: apiHeaders(true),
        credentials: 'same-origin',
        body: JSON.stringify(body),
      })
        .then(async (r) => {
          const data = await r.json().catch(() => ({}));
          const errBox = document.getElementById('form-errors');
          if (!r.ok) {
            errBox.classList.remove('hidden');
            const errs = data.errors || { error: data.error || 'Ошибка сохранения' };
            errBox.innerHTML = Object.entries(errs)
              .map(([k, v]) => `<div><strong>${k}</strong>: ${v}</div>`)
              .join('');
            return;
          }
          window.location.href = `/admin/${document.body.dataset.modelSlug}`;
        });
    });
  }
})();
