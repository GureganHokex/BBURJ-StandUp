(function () {
  var STORAGE_KEY = 'bburj_cookie_consent';

  function loadMetrika() {
    if (window.__bburjMetrikaLoaded) return;
    window.__bburjMetrikaLoaded = true;
    var script = document.createElement('script');
    script.src = '/static/js/yandex-metrika.js';
    script.defer = true;
    document.head.appendChild(script);
  }

  function hasConsent() {
    try {
      return localStorage.getItem(STORAGE_KEY) === '1';
    } catch (e) {
      return false;
    }
  }

  function acceptConsent() {
    try {
      localStorage.setItem(STORAGE_KEY, '1');
    } catch (e) {}
    var banner = document.getElementById('cookie-banner');
    if (banner) banner.hidden = true;
    loadMetrika();
  }

  function showBanner() {
    var banner = document.getElementById('cookie-banner');
    if (banner) banner.hidden = false;
  }

  document.addEventListener('DOMContentLoaded', function () {
    var acceptBtn = document.getElementById('cookie-banner-accept');
    if (acceptBtn) {
      acceptBtn.addEventListener('click', acceptConsent);
    }

    if (hasConsent()) {
      loadMetrika();
      return;
    }

    showBanner();
  });
})();
