// snfusion-plot takes a CSV file created by snfusion-gen and
// creates a PNG plot out of it.
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/astrogo/snfusion/sim"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

func main() {
	ifname := flag.String("f", "output.csv", "input CSV file to analyze")
	ofname := flag.String("o", "output.png", "output PNG file")

	flag.Parse()

	log.SetPrefix("snfusion-plot: ")
	log.SetFlags(0)

	f, err := os.Open(*ifname)
	if err != nil {
		log.Fatalf("error opening %s: %v\n", *ifname, err)
	}
	defer f.Close()

	var engine sim.Engine
	var hdr []byte
	{
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			data := scanner.Bytes()
			if !bytes.HasPrefix(data, []byte("#")) {
				break
			}
			if !bytes.HasPrefix(data, sim.HeaderCSV) {
				continue
			}
			hdr = make([]byte, len(data)-len(sim.HeaderCSV))
			copy(hdr, data[len(sim.HeaderCSV):])
			break
		}
		err = scanner.Err()
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			log.Fatalf("error peeking at meta-data: %v\n", err)
		}
		if len(hdr) == 0 {
			log.Fatalf("could not find meta-data in file %s\n",
				*ifname,
			)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			log.Fatalf("error rewinding input file %s: %v\n",
				*ifname,
				err,
			)
		}
	}
	err = json.Unmarshal(hdr, &engine)
	if err != nil {
		log.Fatalf("error reading meta-data: %v\n", err)
	}

	log.Printf("plotting...\n")
	log.Printf("NumIters:   %d\n", engine.NumIters)
	log.Printf("NumCarbons: %d\n", engine.NumCarbons)
	log.Printf("Seed:       %d\n", engine.Seed)
	log.Printf("Nuclei:     %v\n", engine.Population)

	r := csv.NewReader(f)
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
		log.Fatalf("error reading data: %v\n", err)
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

	// Save the plot to a PNG file.
	if err := p.Save(25*vg.Centimeter, 15*vg.Centimeter, *ofname); err != nil {
		panic(err)
	}
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
