$(document).ready(function() {
    // If mobile or desktop
    if ($(window).width() < 480) {
        const tilt = $('.tilt').tilt();
        tilt.tilt.destroy.call(tilt);

        // Add dynamic colors based on covers
        var rgb = new ColorFinder(function favorDark(r, g, b) {
            return r > 245 && g > 245 && b > 245 ? 0 : 1;
        }).getMostProminentColor(document.getElementById('coverImg'));
        $('.cover-fill').css(
            'background-color',
            'rgb(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ')'
        );
        $('.cover-outline').css(
            'border',
            '1px solid rgb(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ')'
        );
        $('.cover-color').css(
            'color',
            'rgb(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ')'
        );
    } else {
        // Show Sponsor Popup
        $('.popup-toggle').click(function() {
            $('.hero').toggleClass('hero-hide');
            $('.popup-close').toggleClass('popup-close-show');
            $('.input-name').focus();
        });
        $('.popup-close').click(function() {
            $(this).toggleClass('popup-close-show');
            $('.hero').toggleClass('hero-hide');
            $('.input-name').blur();
        });

        // Sponsor form
        $('.form-sponsor').submit(function(e) {
            e.preventDefault();
            var $form = $(this);
            $.post($form.attr('action'), $form.serialize()).then(function() {
                $('.form-sponsor').addClass('form-hide');
                $('.form-success').addClass('form-success-show');
            });
        });

        // Tilted covers
        $('.tilt').tilt({
            perspective: 1000,
            maxGlare: 0.4,
            glare: true,
            scale: 1.1
        });

        // Add dynamic colors based on covers
        var rgb = new ColorFinder(function favorDark(r, g, b) {
            return r > 245 && g > 245 && b > 245 ? 0 : 1;
        }).getMostProminentColor(document.getElementById('coverImg'));
        $('.cover-fill').css(
            'background-color',
            'rgb(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ')'
        );
        $('.cover-fill').css(
            'box-shadow',
            '0 40px 60px rgb(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ',0.2)'
        );
        $('.cover-outline').css(
            'border',
            '1px solid rgb(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ')'
        );
        $('.cover-color').css(
            'color',
            'rgb(' + rgb.r + ',' + rgb.g + ',' + rgb.b + ')'
        );
    }

    // Audio player
    var castPlayers = document.querySelectorAll('.hero');
    var speeds = [1, 1.5, 2];
    for (i = 0; i < castPlayers.length; i++) {
        var player = castPlayers[i];
        var audio = player.querySelector('audio');
        var play = player.querySelector('.play');
        var pause = player.querySelector('.pause');
        var rewind = player.querySelector('.rewind');
        var progress = player.querySelector('.progress');
        var forward = player.querySelector('.forward');
        var speed = player.querySelector('.speed');
        var mute = player.querySelector('.mute');
        var currentTime = player.querySelector('.currenttime');
        var duration = player.querySelector('.duration');
        var currentSpeedIdx = 0;
        pause.style.display = 'none';
        var toHHMMSS = function(totalsecs) {
            var sec_num = parseInt(totalsecs, 10);
            var hours = Math.floor(sec_num / 3600);
            var minutes = Math.floor((sec_num - hours * 3600) / 60);
            var seconds = sec_num - hours * 3600 - minutes * 60;
            if (hours < 10) {
                hours = '0' + hours;
            }
            if (minutes < 10) {
                minutes = '0' + minutes;
            }
            if (seconds < 10) {
                seconds = '0' + seconds;
            }
            var time = hours + ':' + minutes + ':' + seconds;
            return time;
        };
        audio.addEventListener('loadedmetadata', function() {
            progress.setAttribute('max', Math.floor(audio.duration));
            duration.textContent = toHHMMSS(audio.duration);
        });
        audio.addEventListener('timeupdate', function() {
            progress.setAttribute('value', audio.currentTime);
            currentTime.textContent = toHHMMSS(audio.currentTime);
        });
        play.addEventListener(
            'click',
            function() {
                $('.cover').addClass('cover-playing');
                this.style.display = 'none';
                pause.style.display = 'inline-block';
                pause.focus();
                audio.play();
            },
            false
        );
        pause.addEventListener(
            'click',
            function() {
                $('.cover').removeClass('cover-playing');
                this.style.display = 'none';
                play.style.display = 'inline-block';
                play.focus();
                audio.pause();
            },
            false
        );
        rewind.addEventListener(
            'click',
            function() {
                audio.currentTime -= 15;
            },
            false
        );
        progress.addEventListener(
            'click',
            function(e) {
                audio.currentTime =
                    Math.floor(audio.duration) *
                    (e.offsetX / e.target.offsetWidth);
            },
            false
        );
        forward.addEventListener(
            'click',
            function() {
                audio.currentTime += 15;
            },
            false
        );
        speed.addEventListener(
            'click',
            function() {
                currentSpeedIdx =
                    currentSpeedIdx + 1 < speeds.length
                        ? currentSpeedIdx + 1
                        : 0;
                audio.playbackRate = speeds[currentSpeedIdx];
                this.textContent = speeds[currentSpeedIdx] + 'x';
                return true;
            },
            false
        );
        mute.addEventListener(
            'click',
            function() {
                if (audio.muted) {
                    audio.muted = false;
                    $('.icon-mute').addClass('hide');
                    $('.icon-unmute').removeClass('hide');
                } else {
                    audio.muted = true;
                    $('.icon-mute').removeClass('hide');
                    $('.icon-unmute').addClass('hide');
                }
            },
            false
        );
    }

    // Feather icons
    feather.replace();

    // Page transitions
    $('.animsition').animsition({
        linkElement: 'a:not([target="_blank"]):not([href^="#"])',
        linkElement: '.animsition-link',
        outDuration: 600,
        inDuration: 600
    });
});
