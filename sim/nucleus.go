package sim

import "fmt"

// Nucleus models a standard model nucleus.
// It holds the mass number A and the atomic number Z
// of this nucleus.
type Nucleus struct {
	A int // mass number
	Z int // atomic number
}

// N returns the number of nucleons
func (n Nucleus) N() int {
	return n.A - n.Z
}

func (n Nucleus) String() string {
	return fmt.Sprintf("Nucleus{A: %2d, Z:%2d}", n.A, n.Z)
}

// Fuse returns the product of the fusion of two nuclei n1 and n2,
// or false if the fusion is not physically possible.
func Fuse(n1, n2 Nucleus) (Nucleus, bool) {
	o := Nucleus{
		A: n1.A + n2.A,
		Z: n1.Z + n2.Z,
	}
	if o.A <= 56 {
		return o, true
	}
	return Nucleus{}, false
}

// Nuclei is a slice of nuclei, sortable by increasing mass number.
type Nuclei []Nucleus

func (p Nuclei) Len() int           { return len(p) }
func (p Nuclei) Less(i, j int) bool { return p[i].A < p[j].A }
func (p Nuclei) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var (
	nC  = Nucleus{A: 12, Z: 6}
	nO  = Nucleus{A: 16, Z: 8}
	nMg = Nucleus{A: 24, Z: 12}
	nSi = Nucleus{A: 28, Z: 14}
	nS  = Nucleus{A: 32, Z: 16}
	nAr = Nucleus{A: 36, Z: 18}
	nCa = Nucleus{A: 40, Z: 20}
	nTi = Nucleus{A: 44, Z: 22}
	nCr = Nucleus{A: 48, Z: 24}
	nFe = Nucleus{A: 52, Z: 26}
	nNi = Nucleus{A: 56, Z: 28}

	// xsects is all the cross-sections that a sim.Engine can handle
	xsects = map[pair]float64{
		pair{nC, nC}:   0.8315672884,
		pair{nO, nC}:   1,
		pair{nO, nO}:   0.9872126376,
		pair{nMg, nC}:  0.97267664,
		pair{nMg, nO}:  0.965386457,
		pair{nMg, nMg}: 0.8924961757,
		pair{nSi, nC}:  0.7969537454,
		pair{nSi, nO}:  0.6755141681,
		pair{nSi, nMg}: 0.7702788102,
		pair{nSi, nSi}: 0.6517919696,
		pair{nS, nC}:   0.6883015304,
		pair{nS, nO}:   0.7202932946,
		pair{nS, nMg}:  0.8330120436,
		pair{nAr, nC}:  0.7513868241,
		pair{nAr, nO}:  0.7976021702,
		pair{nCa, nC}:  0.8048923532,
		pair{nCa, nO}:  0.8548382242,
		pair{nTi, nC}:  0.9762778016,
	}
)

type pair [2]Nucleus

func init() {
	for k, v := range xsects {
		if k[0] == k[1] {
			continue
		}
		xsects[pair{k[1], k[0]}] = v
	}
}
