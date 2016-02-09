package main

const index = `
<!doctype html>

<html>

	<head>
		<meta charset="utf-8"/>
		<title>SuperNovae Fusion</title>
		<meta name="viewport" content="width=device-width, minimum-scale=1.0, initial-scale=1.0, user-scalable=yes" />
		<script src="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/webcomponentsjs/webcomponents.js"></script>


		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/polymer/polymer.html"/>
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/iron-icons/iron-icons.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/iron-input/iron-input.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-button/paper-button.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-input/paper-input.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-spinner/paper-spinner.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-toast/paper-toast.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-toolbar/paper-toolbar.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-scroll-header-panel/paper-scroll-header-panel.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-icon-button/paper-icon-button.html" />
		<link rel="import" href="https://cdn.rawgit.com/download/polymer-cdn/1.2.3.2/lib/paper-styles/color.html" />

		<style>
	paper-scroll-header-panel {
      position: absolute;
      top: 0;
      right: 0;
      bottom: 0;
      left: 0;
      background-color: var(--paper-grey-200, #eee);
    }

    paper-toolbar {
      background-color: var(--google-blue-500, #4285f4);
    }

    paper-toolbar .title {
      margin: 0 8px;
    }

    paper-scroll-header-panel .content {
      padding: 8px;
    }

    paper-icon-button {
      --paper-icon-button-ink-color: white;
    }

    .spacer {
      @apply(--layout-flex);
    }

	paper-input {
		display: block;
	}

	body {
		padding: 40px;
	}

	div.content {
		width: 60%;
	}

	.center {
		margin: auto;
		width: 60%;
		border: 1px solid;
		padding: 10px;
	}

	paper-button[raised].colorful {
		background-color: #4285f4;
		color: #fff;
	}
		</style>

<script type="text/javascript">

var sock = null;
var wsuri = "ws://{{.Addr}}/data";

window.onload = function() {
	console.log("onload");

	sock = new WebSocket(wsuri);

	sock.onopen = function() {
		console.log("connected to " + wsuri);
	}

	sock.onclose = function(e) {
		console.log("connection closed (" + e.code + ")");
	}

	sock.onmessage = function(e) {
		var obj = JSON.parse(e.data);
		console.log("got: "+JSON.stringify(obj));
		switch (obj["stage"]) {
			case "gen-done":
				document.getElementById("sim-spinner").active = false;
				if (obj["err"] != null) {
					document.getElementById("snfusion-gen-output").innerHTML = JSON.stringify(obj["err"]);
				}
				break;
			case "plot-done":
				document.getElementById("snfusion-plot").innerHTML = obj["svg"];
		}
	}
};

function snfusionGen() {
	document.getElementById("sim-spinner").active = true;
	var data = {
		"num_iters": Number(document.getElementById("num-iters").value),
		"num_carbons": Number(document.getElementById("num-carbons").value),
		"seed": Number(document.getElementById("seed").value)
	}
	console.log("data: "+JSON.stringify(data));

	sock.send(JSON.stringify(data));
}
</script>

	</head>

	<body unresolved>

		<paper-scroll-header-panel fixed>

			<paper-toolbar>
				<paper-icon-button icon="arrow-back"></paper-icon-button>
				<div class="spacer title">sn-fusion</div>
				<paper-icon-button icon="search"></paper-icon-button>
				<paper-icon-button icon="more-vert"></paper-icon-button>
			</paper-toolbar>


			<div class="content snfusion-gen-params">
				<div class="center">
					Please specify the simulation parameters...
					<br>
					<div class="center">
						<paper-input id="num-iters" label="# iters" value="10000"></paper-input>
						<paper-input id="num-carbons" label="% carbon atoms" value="60"></paper-input>
						<paper-input id="seed" label="seed" value="1234"></paper-input>
					</div>
					<br>
					<center>
						<paper-button raised class="colorful" onclick="snfusionGen()">Launch simulation</paper-button>
						<br>
						<paper-spinner alt="Running simulation..." id="sim-spinner"></paper-spinner>
						<p id="snfusion-gen-output"></p>
						<p id="snfusion-plot"></p>
					</center>
				</div>
			</div>

		</paper-scroll-header-panel>


	</body>

</html>`
