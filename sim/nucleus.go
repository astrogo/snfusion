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
