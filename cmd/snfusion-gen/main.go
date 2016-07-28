// snfusion-gen is a simple simulation program modelling fusion processes
// happening inside supernovae.
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/astrogo/snfusion/sim"
)

var (
	nCarbons = flag.Float64(
		"carbon-ratio", 60,
		"carbon ratio (0-100) giving the initial Carbon/Oxygen composition",
	)
	nIters = flag.Int(
		"n", 100000,
		"number of iterations to simulate",
	)
	seed = flag.Int64(
		"seed", 1234,
		"seed used for the MonteCarlo",
	)

	doprof = flag.Bool("cpu-prof", false, "enable CPU profiling")

	fname = flag.String("o", "output.csv", "output file name")
)

func main() {
	flag.Parse()

	if *doprof {
		f, err := os.Create(*fname + ".prof")
		if err != nil {
			log.Fatalf("error creating CPU-profiling file: %v\n", err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.SetFlags(0)
	log.SetPrefix("snfusion-gen: ")
	log.Printf("processing...\n")
	beg := time.Now()

	f, err := os.Create(*fname)
	if err != nil {
		log.Fatalf("error creating %s: %v\n", *fname, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	engine := sim.Engine{
		NumIters:   *nIters,
		NumCarbons: *nCarbons,
		Seed:       *seed,
	}

	err = engine.Run(w)
	delta := time.Now().Sub(beg)
	log.Printf("processing... [done]: %v\n", delta)

	if err != nil {
		log.Fatalf("error running engine: %v\n", err)
	}
}
