// Package sim provides utilities to simulate fusion processes
// happening in a supernova.
//
// A typical use case would look like:
//
//	var w io.Writer
//	engine := sim.Engine{
//		NumIters:   *nIters,
//		NumCarbons: *nCarbons,
//		Seed:       *seed,
//	}
//	err = engine.Run(w)
//
package sim

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

var (
	// Population is the default list of nuclei we want to study.
	Population = []Nucleus{
		{A: 12, Z: 6},  // 12-C
		{A: 16, Z: 8},  // 16-O
		{A: 24, Z: 12}, // 24-Mg
		{A: 28, Z: 14}, // 28-Si
		{A: 32, Z: 16}, // 32-S
		{A: 36, Z: 18}, // 36-Ar
		{A: 40, Z: 20}, // 40-Ca
		{A: 44, Z: 22}, // 44-Ti
		{A: 48, Z: 24}, // 48-Cr
		{A: 52, Z: 26}, // 52-Fe
		{A: 56, Z: 28}, // 56-Ni
	}

	// HeaderCSV identifies the start of meta-data
	HeaderCSV = []byte("# snfusion-gen=")
)

func itoa(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// Engine controls the time evolution of an SN-Fusion simulation.
type Engine struct {
	NumIters   int
	NumCarbons int
	nuclei     []Nucleus
	Seed       int64
	Population []Nucleus
	rng        *rand.Rand
	w          io.Writer
	wcsv       *csv.Writer
}

// Run runs the whole simulation and writes data (as well as
// metadata) into w.
// The data is written as a CSV file with '#' comments and ';' separators.
func (e *Engine) Run(w io.Writer) error {
	err := e.init(w)
	if err != nil {
		return err
	}

	defer e.wcsv.Flush()

	log.Printf("%v\n", e.stats())

	for i := 0; i < e.NumIters; i++ {
		err = e.process()
		if err != nil {
			return err
		}
		if (i+1)%int(float64(e.NumIters)*0.1) == 0 {
			log.Printf("iter #%d/%d...\n", i+1, e.NumIters)
		}
	}

	log.Printf("%v\n", e.stats())

	e.wcsv.Flush()
	err = e.wcsv.Error()
	if err != nil {
		return err
	}

	return err
}

func (e *Engine) init(w io.Writer) error {
	e.rng = rand.New(rand.NewSource(e.Seed))
	e.w = w

	var err error
	const nmax = 10000
	e.nuclei = make([]Nucleus, 0, nmax)
	for i := 0; i < nmax; i++ {
		var n Nucleus
		v := e.rng.Intn(100)
		switch {
		case v <= e.NumCarbons:
			n.Z = 6 // carbon-12
			n.A = 12
		case v > e.NumCarbons:
			n.Z = 8 // oxygen-16
			n.A = 16
		}
		e.nuclei = append(e.nuclei, n)
	}

	if e.Population == nil {
		e.Population = make([]Nucleus, len(Population))
		copy(e.Population, Population)
	}

	hdr, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(e.w, "%v%v\n", string(HeaderCSV), string(hdr))
	if err != nil {
		return err
	}

	e.wcsv = csv.NewWriter(e.w)
	e.wcsv.Comma = ';'

	err = e.writeRecord()
	if err != nil {
		return err
	}

	return err
}

func (e *Engine) process() error {
	var err error
	defer func() {
		if err == nil {
			err = e.writeRecord()
		}
	}()

	i := e.rng.Intn(len(e.nuclei))
	j := e.rng.Intn(len(e.nuclei))
	if i == j {
		return err
	}
	ni := e.nuclei[i]
	nj := e.nuclei[j]
	o, ok := Fuse(ni, nj)
	if !ok {
		// can't fuse nuclei
		return err
	}
	fuse := e.rng.Float64() < xsects[pair{ni, nj}]
	switch fuse {
	case true:
		e.nuclei[i] = o
		e.delete(j)
	case false:
		return err
	}

	return err
}

func (e *Engine) delete(i int) {
	e.nuclei[i] = e.nuclei[len(e.nuclei)-1]
	e.nuclei = e.nuclei[:len(e.nuclei)-1]
}

func (e *Engine) stats() stats {
	histo := make(map[Nucleus]int, len(e.nuclei)/2)
	for _, n := range e.nuclei {
		histo[n]++
	}
	nuclei := make(Nuclei, 0, len(histo))
	for n := range histo {
		nuclei = append(nuclei, n)
	}

	sort.Sort(nuclei)

	stats := stats{
		n:      len(e.nuclei),
		nuclei: nuclei,
		histo:  histo,
	}

	return stats
}

func (e *Engine) writeRecord() error {
	data := make([]string, len(e.Population))
	stats := e.stats()
	for i, n := range e.Population {
		data[i] = itoa(stats.histo[n] * n.A)
	}
	return e.wcsv.Write(data)
}

type stats struct {
	n      int
	nuclei Nuclei
	histo  map[Nucleus]int
}

func (s stats) String() string {
	o := []string{}
	o = append(o, fmt.Sprintf("composition of %d nuclei:", s.n))
	for _, n := range s.nuclei {
		o = append(o, fmt.Sprintf("%v: %d", n, s.histo[n]))
	}
	return strings.Join(o, "\n")
}
