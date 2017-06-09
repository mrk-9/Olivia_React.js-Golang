export default function (name, email, created_at, immediately) {
  window.intercomSettings = {
    app_id: "ppiuwugb"
  };

  if (name) {
    window.intercomSettings.name = name;
  }

  if (email) {
    window.intercomSettings.email = email;
  }

  if (created_at) {
    if (typeof created_at == 'string') {
      created_at = parseInt((new Date(created_at)).getTime() / 1000, 10);
    }

    window.intercomSettings.created_at = created_at;
  }

  (function () {
    var w = window;
    var ic = w.Intercom;
    if (typeof ic === "function") {
      ic('reattach_activator');
      ic('update', intercomSettings);
    } else {
      var d = document;
      var i = function () {
        i.c(arguments)
      };
      i.q = [];
      i.c = function (args) {
        i.q.push(args)
      };
      w.Intercom = i;
      function l() {
        var s = d.createElement('script');
        s.type = 'text/javascript';
        s.async = true;
        s.src = 'https://widget.intercom.io/widget/APP_ID';
        var x = d.getElementsByTagName('script')[0];
        x.parentNode.insertBefore(s, x);
      }

      if (immediately === true) {
        l();
      } else {
        if (w.attachEvent) {
          w.attachEvent('onload', l);
        } else {
          w.addEventListener('load', l, false);
        }
      }
    }
  })()
}