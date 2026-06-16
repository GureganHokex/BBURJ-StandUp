(function () {
  // Debug overlay appears when ?_ym_debug=2 or _ym_debug cookie is set (e.g. after counter check).
  try {
    if (window._ym_debug) {
      delete window._ym_debug;
    }

    if (/\b_ym_debug=/.test(location.search)) {
      var params = new URLSearchParams(location.search);
      params.delete('_ym_debug');
      var query = params.toString();
      history.replaceState(null, '', location.pathname + (query ? '?' + query : '') + location.hash);
    }

    document.cookie.split(';').forEach(function (part) {
      var name = part.split('=')[0].trim();
      if (name === '_ym_debug') {
        document.cookie = name + '=; Max-Age=0; Path=/';
        document.cookie = name + '=; Max-Age=0; Path=/; Domain=' + location.hostname;
        document.cookie = name + '=; Max-Age=0; Path=/; Domain=.' + location.hostname.replace(/^www\./, '');
      }
    });

    Object.keys(localStorage).forEach(function (key) {
      if (key.indexOf('_ym_debug') === 0) {
        localStorage.removeItem(key);
      }
    });
  } catch (e) {}

  (function (m, e, t, r, i, k, a) {
    m[i] = m[i] || function () { (m[i].a = m[i].a || []).push(arguments); };
    m[i].l = 1 * new Date();
    for (var j = 0; j < document.scripts.length; j++) {
      if (document.scripts[j].src === r) return;
    }
    k = e.createElement(t);
    a = e.getElementsByTagName(t)[0];
    k.async = 1;
    k.src = r;
    a.parentNode.insertBefore(k, a);
  })(window, document, 'script', 'https://mc.yandex.ru/metrika/tag.js?id=109884189', 'ym');

  ym(109884189, 'init', {
    webvisor: true,
    clickmap: true,
    ecommerce: 'dataLayer',
    accurateTrackBounce: true,
    trackLinks: true,
  });
})();
