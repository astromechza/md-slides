package sliderenderer

const scriptHeader = `
<script>
var prevSlide = "/_slides?page=%d";
var nextSlide = "/_slides?page=%d";

document.onkeydown = function(evt) {
	evt = evt || window.event
	if ([13, 32, 39, 40].indexOf(evt.keyCode) >= 0) {
		window.location = nextSlide;
	}
	if ([8, 37, 38].indexOf(evt.keyCode) >= 0 ) {
		window.location = prevSlide;
	}
}

window.onresize = function(event) {
	var el = document.getElementById("body-inner");
	var m = 50;
	var wi = el.clientWidth + m;
	var hi = el.clientHeight + m;
	var ws = window.innerWidth / wi;
	var hs = window.innerHeight / hi;
	var ss = Math.min(ws, hs);
	el.style.transform = "scale(" + ss + ")";
};

document.addEventListener("DOMContentLoaded", function(event) {
    window.onresize(null);
});

</script>
`

const styleHeader = `
<style>
html {
	height: 100%%;
	font-size: 22px;
}

body {
	height: 100%%;
    display: flex;
    flex-flow: column;
	background-color: #444;
	justify-content: center;
}

.body-inner {
	align-self: center;
	display: grid;
	box-sizing: border-box;
	background: %s;
	padding: 1rem;
    border-radius: 0.3rem;
	box-shadow: 0px 0.2rem 0.6rem black;
	padding-left: 3rem;
    padding-right: 3rem;
	position: absolute;
	overflow: hidden;
	grid-auto-columns: 1fr;
	grid-auto-rows: 1fr;
}

.body-inner-halign-left {justify-items: start;}
.body-inner-halign-center {justify-items: center;}
.body-inner-halign-right {justify-items: end;}
.body-inner-valign-top {align-items: start;}
.body-inner-valign-center {align-items: center;}
.body-inner-valign-bottom {align-items: end;}
.body-inner-talign-left {text-align: left;}
.body-inner-talign-center {text-align: center;}
.body-inner-talign-right {text-align: right;}

.page-number {
	font-family: Palatino, "Palatino Linotype", "Palatino LT STD", "Book Antiqua", Georgia, serif;
	position: absolute;
	bottom: 0;
	right: 0;
	margin: 0.5rem;
	color: lightgrey;
}

</style>
`

const styleMultiHeader = `
<style>
body {
	justify-content: start;
	position: relative;
}

.body-inner {
	margin-top: 1.5rem;
	margin-bottom: 1.5rem;
	position: relative;
}
</style>

<style type="text/css" media="print">
.body-inner {
	margin-top: 20px;
	margin-bottom: 20px;
	page-break-before: always;
	page-break-inside: avoid;
	page-break-after: always;
}

@page {
	size: 1640px 940px;
}
</style>
`
