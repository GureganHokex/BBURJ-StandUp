(function () {
  const box = document.getElementById('photo-lightbox');
  if (!box) return;

  const backdrop = box.querySelector('.lightbox-backdrop');
  const closeBtn = box.querySelector('.lightbox-close');
  const prevBtn = box.querySelector('.lightbox-prev');
  const nextBtn = box.querySelector('.lightbox-next');
  const img = box.querySelector('.lightbox-image');
  const caption = box.querySelector('.lightbox-caption');

  let items = [];
  let index = 0;

  function collectItems(trigger) {
    const gallery = trigger.closest('[data-lightbox-gallery]');
    if (!gallery) return [trigger];
    return Array.from(gallery.querySelectorAll('[data-lightbox-src]'));
  }

  function show(i) {
    if (!items.length) return;
    index = (i + items.length) % items.length;
    const el = items[index];
    const src = el.dataset.lightboxSrc;
    const cap = el.dataset.lightboxCaption || '';
    img.src = src;
    img.alt = cap || 'Фото';
    caption.textContent = cap;
    caption.hidden = !cap;
    const showNav = items.length > 1;
    prevBtn.hidden = !showNav;
    nextBtn.hidden = !showNav;
    box.classList.add('is-open');
    box.setAttribute('aria-hidden', 'false');
    document.body.classList.add('lightbox-open');
  }

  function close() {
    box.classList.remove('is-open');
    box.setAttribute('aria-hidden', 'true');
    document.body.classList.remove('lightbox-open');
    img.removeAttribute('src');
  }

  document.addEventListener('click', (e) => {
    const thumb = e.target.closest('[data-lightbox-src]');
    if (!thumb) return;
    e.preventDefault();
    items = collectItems(thumb);
    const idx = items.indexOf(thumb);
    show(idx >= 0 ? idx : 0);
  });

  closeBtn?.addEventListener('click', close);
  backdrop?.addEventListener('click', close);
  prevBtn?.addEventListener('click', (e) => {
    e.stopPropagation();
    show(index - 1);
  });
  nextBtn?.addEventListener('click', (e) => {
    e.stopPropagation();
    show(index + 1);
  });

  document.addEventListener('keydown', (e) => {
    if (!box.classList.contains('is-open')) return;
    if (e.key === 'Escape') close();
    if (e.key === 'ArrowLeft') show(index - 1);
    if (e.key === 'ArrowRight') show(index + 1);
  });
})();
