# LiveTemplate
Code live template


## 浏览器变身编辑器

```html
data:text/html,
<html>
<head>
    <style type="text/css">

    @ font-face { font-family: 'Open Sans'; font-style: normal; font-weight: 40; }
    html { font-family: "Open Sans"; }
    * { -webkit-transition: all linear 1s; }
    </style >
    <script >
        window.onload = function () {
            var e = false ;
            var t = 0;
            setInterval( function () {
                if (!e) {
                    t = Math.round(Math.max(0, t - Math.max(t / 3, 1)))
                }
                var n = (255 - t * 2).toString(16); document.body.style.backgroundColor = "#ff" + n + "" + n
            }, 1e3);
            var n = null ;
            document.onkeydown = function () {
                t = Math.min(128, t + 2); e = true; clearTimeout(n);
                n = setTimeout( function () { e = false }, 1500)
            }
}</script >
</head>
<body contenteditable style="font-size :1rem;line-height: 1.4;max-width :50rem;margin: 0 auto;padding :0rem;">

```