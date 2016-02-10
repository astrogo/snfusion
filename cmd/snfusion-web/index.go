package main

import (
	"fmt"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _index_html = []byte(`<!doctype html>

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

</html>
`)

func index_html_bytes() ([]byte, error) {
	return _index_html, nil
}

func index_html() (*asset, error) {
	bytes, err := index_html_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "index.html", size: 4708, mode: os.FileMode(420), modTime: time.Unix(1455095779, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"index.html": index_html,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"index.html": &_bintree_t{index_html, map[string]*_bintree_t{
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

