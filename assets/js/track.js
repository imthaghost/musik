var _gauges = _gauges || [];
(function() {
    var h = _gauges['slice'] ? _gauges.slice(0) : [];
    _gauges = {
        track_referrer: true,
        image: new Image(),
        track: function() {
            this.setCookie('_gauges_cookie', 1, 1);
            var a = this.url();
            if (a) {
                this.image.src = a;
                var b = 60 * 60,
                    f = b * 24,
                    c = f * 31,
                    d = f * 365,
                    j = d * 10;
                if (!this.getCookie('_gauges_unique_hour')) {
                    this.setCookie('_gauges_unique_hour', 1, b);
                }
                if (!this.getCookie('_gauges_unique_day')) {
                    this.setCookie('_gauges_unique_day', 1, f);
                }
                if (!this.getCookie('_gauges_unique_month')) {
                    this.setCookie('_gauges_unique_month', 1, c);
                }
                if (!this.getCookie('_gauges_unique_year')) {
                    this.setCookie('_gauges_unique_year', 1, d);
                }
                this.setCookie('_gauges_unique', 1, d);
            }
        },
        push: function(a) {
            var b = a.shift();
            if (b == 'track') {
                _gauges.track();
            }
        },
        url: function() {
            var a,
                b,
                f,
                c = this.$('gauges-tracker');
            if (c) {
                b = c.getAttribute('data-site-id');
                f = c.getAttribute('data-track-path');
                if (!f) {
                    f = c.src.replace('/track.js', '/track.gif');
                }
                a = String(f);
                a += '?h[site_id]=' + b;
                a += '&h[resource]=' + this.resource();
                a += '&h[referrer]=' + this.referrer();
                a += '&h[title]=' + this.title();
                a += '&h[user_agent]=' + this.agent();
                a += '&h[unique]=' + this.unique();
                a += '&h[unique_hour]=' + this.uniqueHour();
                a += '&h[unique_day]=' + this.uniqueDay();
                a += '&h[unique_month]=' + this.uniqueMonth();
                a += '&h[unique_year]=' + this.uniqueYear();
                a += '&h[screenx]=' + this.screenWidth();
                a += '&h[browserx]=' + this.browserWidth();
                a += '&h[browsery]=' + this.browserHeight();
                a += '&timestamp=' + this.timestamp();
            }
            return a;
        },
        domain: function() {
            return window.location.hostname;
        },
        referrer: function() {
            var a = '';
            if (!this.track_referrer) {
                return a;
            }
            this.track_referrer = false;
            try {
                a = top.document.referrer;
            } catch (e1) {
                try {
                    a = parent.document.referrer;
                } catch (e2) {
                    a = '';
                }
            }
            if (a == '') {
                a = document.referrer;
            }
            return this.escape(a);
        },
        agent: function() {
            return this.escape(navigator.userAgent);
        },
        escape: function(a) {
            return typeof encodeURIComponent == 'function'
                ? encodeURIComponent(a)
                : escape(a);
        },
        resource: function() {
            return this.escape(document.location.href);
        },
        timestamp: function() {
            return new Date().getTime();
        },
        title: function() {
            return document.title && document.title != ''
                ? this.escape(document.title)
                : '';
        },
        uniqueHour: function() {
            if (!this.getCookie('_gauges_cookie')) {
                return 0;
            }
            return this.getCookie('_gauges_unique_hour') ? 0 : 1;
        },
        uniqueDay: function() {
            if (!this.getCookie('_gauges_cookie')) {
                return 0;
            }
            return this.getCookie('_gauges_unique_day') ? 0 : 1;
        },
        uniqueMonth: function() {
            if (!this.getCookie('_gauges_cookie')) {
                return 0;
            }
            return this.getCookie('_gauges_unique_month') ? 0 : 1;
        },
        uniqueYear: function() {
            if (!this.getCookie('_gauges_cookie')) {
                return 0;
            }
            return this.getCookie('_gauges_unique_year') ? 0 : 1;
        },
        unique: function() {
            if (!this.getCookie('_gauges_cookie')) {
                return 0;
            }
            return this.getCookie('_gauges_unique') ? 0 : 1;
        },
        screenWidth: function() {
            try {
                return screen.width;
            } catch (e) {
                return 0;
            }
        },
        browserDimensions: function() {
            var a = 0,
                b = 0;
            try {
                if (typeof window.innerWidth == 'number') {
                    a = window.innerWidth;
                    b = window.innerHeight;
                } else if (
                    document.documentElement &&
                    document.documentElement.clientWidth
                ) {
                    a = document.documentElement.clientWidth;
                    b = document.documentElement.clientHeight;
                } else if (document.body && document.body.clientWidth) {
                    a = document.body.clientWidth;
                    b = document.body.clientHeight;
                }
            } catch (e) {}
            return { width: a, height: b };
        },
        browserWidth: function() {
            return this.browserDimensions().width;
        },
        browserHeight: function() {
            return this.browserDimensions().height;
        },
        $: function(a) {
            if (document.getElementById) {
                return document.getElementById(a);
            }
            return null;
        },
        setCookie: function(a, b, f) {
            var c, d;
            b = escape(b);
            if (f) {
                c = new Date();
                c.setTime(c.getTime() + f * 1000);
                d = '; expires=' + c.toGMTString();
            } else {
                d = '';
            }
            document.cookie = a + '=' + b + d + '; path=/';
        },
        getCookie: function(a) {
            var b = a + '=',
                f = document.cookie.split(';');
            for (var c = 0; c < f.length; c++) {
                var d = f[c];
                while (d.charAt(0) == ' ') {
                    d = d.substring(1, d.length);
                }
                if (d.indexOf(b) == 0) {
                    return unescape(d.substring(b.length, d.length));
                }
            }
            return null;
        }
    };
    _gauges.track();
    for (var g = 0, i = h.length; g < i; g++) {
        _gauges.push(h[g]);
    }
})();
