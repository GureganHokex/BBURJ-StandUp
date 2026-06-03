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
    if (e.key === 'Escape') setOpen(false);
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
