package sliderenderer

const scriptHeader = `
<script>
document.onkeydown = function(evt) {
	evt = evt || window.event
	if ([13, 32, 39, 40].indexOf(evt.keyCode) >= 0) {
		window.location = "{{ .URLPath }}?page={{ .PageNext }}";
	}
	if ([8, 37, 38].indexOf(evt.keyCode) >= 0 ) {
		window.location = "{{ .URLPath }}?page={{ .PagePrev }}";
	}
}

window.onresize = function(event) {
	var el = document.getElementsByClassName("slide-wrap")[0];
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
	height: 100%;
	font-size: {{ .FontSize }}px;
}

body {
	height: 100%;
    display: flex;
    flex-flow: column;
	background: #444444;
	justify-content: center;
}

.slide-wrap {
	align-self: center;
	display: grid;
	box-sizing: border-box;
	background: white;
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

.slide-wrap-halign-left {justify-items: start;}
.slide-wrap-halign-center {justify-items: center;}
.slide-wrap-halign-right {justify-items: end;}
.slide-wrap-valign-top {align-items: start;}
.slide-wrap-valign-center {align-items: center;}
.slide-wrap-valign-bottom {align-items: end;}
.slide-wrap-talign-left {text-align: left;}
.slide-wrap-talign-center {text-align: center;}
.slide-wrap-talign-right {text-align: right;}

.page-number {
	font-family: Palatino, "Palatino Linotype", "Palatino LT STD", "Book Antiqua", Georgia, serif;
	position: absolute;
	bottom: 0;
	right: 0;
	margin: 0.5rem;
	color: lightgrey;
}

pre.chroma {
	padding: 1rem;
	border-radius: 0.5rem;
	border: 1px solid lightgrey;
}
</style>
`

const styleMultiHeader = `
<style>
html {
	height: auto;
	font-size: {{ .FontSize }}px;
}

body {
	justify-content: start;
	position: relative;
}

.slide-wrap {
	margin-top: 1.5rem;
	margin-bottom: 1.5rem;
	position: relative;
}

@media print {
	.slide-wrap {
		margin-top: 20px;
		margin-bottom: 20px;
		page-break-before: always;
		page-break-inside: avoid;
	}
}

@page {
	size: {{ add .XRes 40 }}px {{ add .YRes 40 }}px;
}
</style>
`
