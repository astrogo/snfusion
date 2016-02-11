// snfusion-web is a simple command serving fusion processes analyses over http.
package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/astrogo/snfusion/sim"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
	"github.com/gonum/plot/vg/vgsvg"

	"golang.org/x/net/websocket"
)

var (
	rootfs = ""
)

func main() {
	srv := newServer()
	http.HandleFunc("/download", download)
	http.Handle("/", srv)
	http.Handle("/data", websocket.Handler(srv.dataHandler))
	log.Printf("listening on http://%s ...\n", srv.Addr)
	err := http.ListenAndServe(srv.Addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	Addr string
	tmpl *template.Template

	clients    map[*client]bool // registered clients
	register   chan *client
	unregister chan *client

	datac chan []byte
}

func newServer() *server {
	srv := &server{
		Addr:       getHostIP() + ":7071",
		tmpl:       template.Must(template.New("snfusion").Parse(string(MustAsset("index.html")))),
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		datac:      make(chan []byte),
	}
	go srv.run()
	return srv
}

func (srv *server) run() {
	for {
		select {
		case c := <-srv.register:
			srv.clients[c] = true
			log.Printf(">>> new-client: %v\n", c)
		case c := <-srv.unregister:
			if _, ok := srv.clients[c]; ok {
				delete(srv.clients, c)
				close(c.datac)
				log.Printf("client disconnected [%v]\n", c.ws.LocalAddr())
			}

		case data := <-srv.datac:
			log.Printf("data: %v\n", len(data))
		}
	}
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("accepting new connection from %v...\n", r.Host)
	srv.tmpl.Execute(w, srv)
	// srv.rootfs.ServeHTTP(w, r)
}

func (srv *server) dataHandler(ws *websocket.Conn) {
	c := &client{
		srv:   srv,
		datac: make(chan []byte, 256),
		ws:    ws,
	}
	srv.register <- c
	c.run()
}

type client struct {
	srv   *server
	ws    *websocket.Conn
	datac chan []byte
}

func (c *client) run() {
	var err error
	dir, err := ioutil.TempDir("", "snfusion-web-")
	if err != nil {
		log.Printf("error creating temporary directory: %v\n", err)
		return
	}

	defer func() {
		c.srv.unregister <- c
		c.ws.Close()
		c.srv = nil
		os.RemoveAll(dir)
	}()

	type params struct {
		ID         int   `json:"id"`
		NumIters   int   `json:"num_iters"`
		NumCarbons int   `json:"num_carbons"`
		Seed       int64 `json:"seed"`
	}

	type genReply struct {
		ID     int        `json:"id"`
		Stage  string     `json:"stage"`
		Err    error      `json:"err"`
		Msg    string     `json:"msg"`
		Engine sim.Engine `json:"engine"`
	}

	type plotReply struct {
		ID    int    `json:"id"`
		Stage string `json:"stage"`
		Err   error  `json:"err"`
		SVG   string `json:"svg"`
	}

	type zipReply struct {
		ID    int    `json:"id"`
		Stage string `json:"stage"`
		Err   error  `json:"err"`
		Href  string `json:"href"`
	}

	for {
		param := params{
			NumIters:   100000,
			NumCarbons: 60,
			Seed:       1234,
		}

		log.Printf("waiting for simulation parameters...\n")
		err = websocket.JSON.Receive(c.ws, &param)
		if err != nil {
			log.Printf("error rcv: %v\n", err)
			return
		}
		id := param.ID

		msgbuf := new(bytes.Buffer)
		msg := log.New(msgbuf, "snfusion-sim: ", 0)
		engine := sim.Engine{
			NumIters:   param.NumIters,
			NumCarbons: param.NumCarbons,
			Seed:       param.Seed,
		}
		engine.SetLogger(msg)

		log.Printf("processing... %#v\n", engine)
		csvbuf := new(bytes.Buffer)
		errc := make(chan error)
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			errc <- engine.Run(csvbuf)
			ticker.Stop()
		}()

		err = <-errc
		if err != nil {
			log.Printf("error: %v\n", err)
			_ = websocket.JSON.Send(c.ws, genReply{
				ID: id, Err: err, Engine: engine, Stage: "gen-done", Msg: msgbuf.String(),
			})
			return
		}

		err = websocket.JSON.Send(c.ws, genReply{
			ID: id, Err: err, Engine: engine, Stage: "gen-done", Msg: msgbuf.String(),
		})
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}

		csvdata := make([]byte, len(csvbuf.Bytes()))
		copy(csvdata, csvbuf.Bytes())

		log.Printf("running post-processing...\n")
		r := csv.NewReader(csvbuf)
		r.Comma = ';'
		r.Comment = '#'

		table := make([]plotter.XYs, len(engine.Population))
		for i := range table {
			table[i] = make(plotter.XYs, engine.NumIters+1)
		}

		for ix := 0; ix < engine.NumIters+1; ix++ {
			var text []string
			text, err = r.Read()
			if err != nil {
				break
			}
			for i := range engine.Population {
				table[i][ix].X = float64(ix)
				table[i][ix].Y = float64(atoi(text[i]))
			}
		}
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			log.Printf("error reading data: %v\n", err)
			return
		}

		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		p.Title.Text = fmt.Sprintf(
			"Time evolution of nuclei C%d-O%d (seed=%d)",
			engine.NumCarbons,
			100-engine.NumCarbons,
			engine.Seed,
		)
		p.X.Label.Text = "Iteration number"
		p.Y.Label.Text = "Atomic mass of nuclei"

		for i, n := range engine.Population {

			line, err := plotter.NewLine(table[i])
			if err != nil {
				log.Fatalf(
					"error adding data points for nucleus %v: %v\n",
					n, err,
				)
			}
			line.LineStyle.Color = col(n)
			line.LineStyle.Width = vg.Points(1)
			p.Add(line)
			p.Legend.Add(label(n), line)
		}

		p.Add(plotter.NewGrid())
		p.Legend.Top = true
		p.Legend.XOffs = -1 * vg.Centimeter

		figX := 25 * vg.Centimeter
		figY := figX / vg.Length(math.Phi)

		// Create a Canvas for writing SVG images.
		canvas := vgsvg.New(figX, figY)

		// Draw to the Canvas.
		p.Draw(draw.New(canvas))

		outsvg := new(bytes.Buffer)
		_, err = canvas.WriteTo(outsvg)
		if err != nil {
			log.Printf("error svg: %v\n", err)
			return
		}

		err = websocket.JSON.Send(c.ws, plotReply{
			ID: id, Err: err, SVG: outsvg.String(), Stage: "plot-done",
		})
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}

		pngcanvas := vgimg.PngCanvas{Canvas: vgimg.New(figX, figY)}
		p.Draw(draw.New(pngcanvas))
		outpng := new(bytes.Buffer)
		_, err = pngcanvas.WriteTo(outpng)
		if err != nil {
			log.Printf("error png: %v\n", err)
			return
		}

		href := filepath.Join(dir, fmt.Sprintf("output-%d.zip", id))
		zipf, err := os.Create(href)
		if err != nil {
			log.Printf("error creating zip file: %v\n", err)
		}
		defer zipf.Close()

		zipw := zip.NewWriter(zipf)
		defer zipw.Close()

		for _, file := range []struct {
			Name string
			Body []byte
		}{
			{"output.csv", csvdata},
			{"output.png", outpng.Bytes()},
		} {
			ff, err := zipw.Create(file.Name)
			if err != nil {
				log.Printf("error creating zip content %v: %v\n", file.Name, err)
				return
			}
			_, err = ff.Write(file.Body)
			if err != nil {
				log.Printf("error writing zip content %v: %v\n", file.Name, err)
				return
			}
		}
		err = zipw.Close()
		if err != nil {
			log.Printf("error closing zip-writer: %v\n", err)
			return
		}
		err = zipf.Close()
		if err != nil {
			log.Printf("error closing zip-file: %v\n", err)
			return
		}

		err = websocket.JSON.Send(c.ws, zipReply{
			ID: id, Err: err, Href: href, Stage: "zip-done",
		})
		if err != nil {
			log.Printf("error sending zip: %v\n", err)
			return
		}
		log.Printf("saved report under %v\n", href)
	}

}

func download(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleErr(w, "error parsing form", err, http.StatusInternalServerError)
	}

	log.Printf("download: %#v\n", *r)

	f, err := os.Open(r.Form.Get("file"))
	if err != nil {
		handleErr(w, "error opening report file", err, http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(w, f)

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename=output.zip")
	w.Header().Set("Content-Type", "application/force-download")
}

func handleErr(w http.ResponseWriter, stage string, err error, code int) {
	log.Printf(stage+": %v\n", err)
	fmt.Fprintf(w, stage+": %v\n", err)
	http.Error(w, err.Error(), code)
}

func init() {
	// FIXME(sbinet) makes this more reliable (multiple $GOPATH entries)
	gopath := os.Getenv("GOPATH")
	rootfs = filepath.Join(gopath, "src/github.com/astrogo/snfusion/cmd/snfusion-web/rootfs")
}

func atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func label(n sim.Nucleus) string {
	switch n {
	case sim.Nucleus{A: 12, Z: 6}:
		return "12-C"
	case sim.Nucleus{A: 16, Z: 8}:
		return "16-O"
	case sim.Nucleus{A: 24, Z: 12}:
		return "24-Mg"
	case sim.Nucleus{A: 28, Z: 14}:
		return "28-Si"
	case sim.Nucleus{A: 32, Z: 16}:
		return "32-S"
	case sim.Nucleus{A: 36, Z: 18}:
		return "36-Ar"
	case sim.Nucleus{A: 40, Z: 20}:
		return "40-Ca"
	case sim.Nucleus{A: 44, Z: 22}:
		return "44-Ti"
	case sim.Nucleus{A: 48, Z: 24}:
		return "48-Cr"
	case sim.Nucleus{A: 52, Z: 26}:
		return "52-Fe"
	case sim.Nucleus{A: 56, Z: 28}:
		return "56-Ni"
	}
	return n.String()
}

func rgb(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
}

func col(n sim.Nucleus) color.Color {
	switch n {
	case sim.Nucleus{A: 12, Z: 6}:
		return rgb(0, 0, 0)
	case sim.Nucleus{A: 16, Z: 8}:
		return rgb(0, 0, 255)
	case sim.Nucleus{A: 24, Z: 12}:
		return rgb(0, 255, 0)
	case sim.Nucleus{A: 28, Z: 14}:
		return rgb(0, 128, 255)
	case sim.Nucleus{A: 32, Z: 16}:
		return rgb(255, 255, 51)
	case sim.Nucleus{A: 36, Z: 18}:
		return rgb(128, 128, 128)
	case sim.Nucleus{A: 40, Z: 20}:
		return rgb(192, 192, 192)
	case sim.Nucleus{A: 44, Z: 22}:
		return rgb(255, 0, 255)
	case sim.Nucleus{A: 48, Z: 24}:
		return rgb(51, 255, 255)
	case sim.Nucleus{A: 52, Z: 26}:
		return rgb(255, 165, 0)
	case sim.Nucleus{A: 56, Z: 28}:
		return rgb(255, 0, 0)
	}
	return plotutil.Color(n.A)
}

func getHostIP() string {
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("could not retrieve hostname: %v\n", err)
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		log.Fatalf("could not lookup hostname IP: %v\n", err)
	}

	for _, addr := range addrs {
		ipv4 := addr.To4()
		if ipv4 == nil {
			continue
		}
		return ipv4.String()
	}

	log.Fatalf("could not infer host IP")
	return ""
}
