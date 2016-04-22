snfusion
========

[![GoDoc](https://godoc.org/github.com/astrogo/snfusion?status.svg)](https://godoc.org/github.com/astrogo/snfusion)

`snfusion` is a simple simulation program modelling fusion processes happening in supernovae.

## Installation

```sh
$> go get github.com/astrogo/snfusion/...
```

## Contribute

`astrogo/snfusion` is released under the `BSD-3` license.
Please send a pull request to [astrogo/license](https://github.com/astrogo/license), adding
yourself to the `AUTHORS` and/or `CONTRIBUTORS` files.

## Documentation

Documentation is available on [godoc.org](https://godoc.org/github.com/astrogo/snfusion).

## Example

```sh
$> snfusion-gen -h
Usage of snfusion-gen:
  -carbon-ratio int
    	carbon ratio (0-100) giving the initial Carbon/Oxygen composition (default 60)
  -cpu-prof
    	enable CPU profiling
  -n int
    	number of iterations to simulate (default 100000)
  -o string
    	output file name (default "output.csv")
  -seed int
    	seed used for the MonteCarlo (default 1234)

$> snfusion-gen -n 30000
snfusion-gen: processing...
snfusion-gen: composition of 10000 nuclei:
Nucleus{A: 12, Z: 6}: 6127
Nucleus{A: 16, Z: 8}: 3873
snfusion-gen: iter #3000/30000...
snfusion-gen: iter #6000/30000...
snfusion-gen: iter #9000/30000...
snfusion-gen: iter #12000/30000...
snfusion-gen: iter #15000/30000...
snfusion-gen: iter #18000/30000...
snfusion-gen: iter #21000/30000...
snfusion-gen: iter #24000/30000...
snfusion-gen: iter #27000/30000...
snfusion-gen: iter #30000/30000...
snfusion-gen: composition of 3066 nuclei:
Nucleus{A: 12, Z: 6}: 71
Nucleus{A: 16, Z: 8}: 63
Nucleus{A: 24, Z:12}: 126
Nucleus{A: 28, Z:14}: 170
Nucleus{A: 32, Z:16}: 124
Nucleus{A: 36, Z:18}: 165
Nucleus{A: 40, Z:20}: 377
Nucleus{A: 44, Z:22}: 343
Nucleus{A: 48, Z:24}: 348
Nucleus{A: 52, Z:26}: 640
Nucleus{A: 56, Z:28}: 639
snfusion-gen: processing... [done]: 10.52320492s

$> ll output.csv
-rw-r--r-- 1 binet binet 1.7M Jan 14 21:32 output.csv

$> head output.csv
# snfusion-gen={"NumIters":30000,"NumCarbons":60,"Seed":1234,"Population":[{"A":12,"Z":6},{"A":16,"Z":8},{"A":24,"Z":12},{"A":28,"Z":14},{"A":32,"Z":16},{"A":36,"Z":18},{"A":40,"Z":20},{"A":44,"Z":22},{"A":48,"Z":24},{"A":52,"Z":26},{"A":56,"Z":28}]}
73524;61968;0;0;0;0;0;0;0;0;0
73524;61968;0;0;0;0;0;0;0;0;0
73512;61952;0;28;0;0;0;0;0;0;0
73512;61952;0;28;0;0;0;0;0;0;0

$> snfusion-plot -f output.csv -o output.png
snfusion-plot: plotting...
snfusion-plot: NumIters:   30000
snfusion-plot: NumCarbons: 60
snfusion-plot: Seed:       1234
snfusion-plot: Nuclei:     [Nucleus{A: 12, Z: 6} Nucleus{A: 16, Z: 8} Nucleus{A: 24, Z:12} Nucleus{A: 28, Z:14} Nucleus{A: 32, Z:16} Nucleus{A: 36, Z:18} Nucleus{A: 40, Z:20} Nucleus{A: 44, Z:22} Nucleus{A: 48, Z:24} Nucleus{A: 52, Z:26} Nucleus{A: 56, Z:28}]
```

![60 Carbon-12, 40 Oxygen-16](/doc/output.png)
