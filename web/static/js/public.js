(function () {
  var topBar = document.getElementById('site-top');
  var btn = document.querySelector('.site-menu-btn');
  var nav = document.getElementById('site-nav');
  var backdrop = document.getElementById('site-nav-backdrop');
  var closeBtn = document.querySelector('.site-nav-close');
  if (!topBar || !btn || !nav) return;

  var mq = window.matchMedia('(max-width: 768px)');

  function resetScrollX() {
    window.scrollTo(0, window.scrollY);
    document.documentElement.scrollLeft = 0;
    document.body.scrollLeft = 0;
  }

  function setOpen(open) {
    topBar.classList.toggle('is-menu-open', open);
    document.body.classList.toggle('is-menu-open', open);
    btn.setAttribute('aria-expanded', open ? 'true' : 'false');
    btn.setAttribute('aria-label', open ? 'Закрыть меню' : 'Открыть меню');
    nav.setAttribute('aria-hidden', open ? 'false' : mq.matches ? 'true' : 'false');
    document.body.classList.toggle('menu-open', open);
    if (backdrop) backdrop.hidden = !open;
  }

  function toggle() {
    setOpen(!topBar.classList.contains('is-menu-open'));
  }

  btn.addEventListener('click', toggle);
  if (closeBtn) closeBtn.addEventListener('click', function () { setOpen(false); });
  if (backdrop) backdrop.addEventListener('click', function () { setOpen(false); });

  nav.querySelectorAll('.site-nav-links a').forEach(function (link) {
    link.addEventListener('click', function () { setOpen(false); });
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' && topBar.classList.contains('is-menu-open')) setOpen(false);
  });

  function onBreakpoint() {
    if (!mq.matches) {
      setOpen(false);
      nav.setAttribute('aria-hidden', 'false');
    } else if (!topBar.classList.contains('is-menu-open')) {
      nav.setAttribute('aria-hidden', 'true');
    }
  }

  mq.addEventListener('change', onBreakpoint);
  onBreakpoint();
  resetScrollX();
  window.addEventListener('load', resetScrollX);
  window.addEventListener('pageshow', resetScrollX);
})();

(function () {
  var modal = document.getElementById('event-modal');
  if (!modal) return;

  var posterWrap = document.getElementById('event-modal-poster-wrap');
  var posterImg = document.getElementById('event-modal-poster');
  var dateEl = document.getElementById('event-modal-date');
  var titleEl = document.getElementById('event-modal-title');
  var cityEl = document.getElementById('event-modal-city');
  var descEl = document.getElementById('event-modal-desc');
  var buyEl = document.getElementById('event-modal-buy');

  function closeModal() {
    modal.hidden = true;
    document.body.classList.remove('event-modal-open');
  }

  function openModal(card) {
    var title = card.dataset.eventTitle || '';
    var city = card.dataset.eventCity || '';
    var desc = card.dataset.eventDesc || '';
    var date = card.dataset.eventDate || '';
    var poster = card.dataset.eventPoster || '';
    var ticket = card.dataset.eventTicket || '';

    if (titleEl) titleEl.textContent = title;
    if (cityEl) cityEl.textContent = city ? city.toUpperCase() : '';
    if (dateEl) dateEl.textContent = date;
    if (descEl) {
      descEl.textContent = desc;
      descEl.hidden = !desc;
    }

    if (posterWrap && posterImg) {
      if (poster) {
        posterImg.src = poster;
        posterWrap.classList.remove('hidden');
      } else {
        posterImg.removeAttribute('src');
        posterWrap.classList.add('hidden');
      }
    }

    if (buyEl) {
      if (ticket) {
        buyEl.href = ticket;
        buyEl.classList.remove('hidden');
      } else {
        buyEl.href = '#';
        buyEl.classList.add('hidden');
      }
    }

    modal.hidden = false;
    document.body.classList.add('event-modal-open');
  }

  document.addEventListener('click', function (e) {
    var detailsBtn = e.target.closest('.event-details-btn');
    if (detailsBtn) {
      var card = detailsBtn.closest('.event-card');
      if (card) openModal(card);
      return;
    }
    if (e.target.closest('[data-event-modal-close]')) closeModal();
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' && !modal.hidden) closeModal();
  });
})();
