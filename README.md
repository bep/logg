
[![Tests on Linux, MacOS and Windows](https://github.com/bep/log/workflows/Test/badge.svg)](https://github.com/bep/log/actions?query=workflow:Test)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/log)](https://goreportcard.com/report/github.com/bep/log)
[![GoDoc](https://godoc.org/github.com/bep/log?status.svg)](https://godoc.org/github.com/bep/log)

This is a fork of the exellent [Apex Log](https://github.com/apex/log) library.

Main changes:

* Trim unneeded dependencies.
* Make `Fields` into a slice to preserve log order.
* Split `Entry` into `Entry` and `EntryFields`. This is easier to reason about and more effective (see benchmarks below)

The existing benchmarks compared to the `7e0ed94172ec33f01a921811762c99e93d531cf6`:

```bash
name              old time/op    new time/op    delta
Logger_small-10      108ns ± 0%      75ns ± 0%  -30.54%  (p=0.029 n=4+4)
Logger_medium-10     338ns ± 0%     175ns ± 0%  -48.17%  (p=0.029 n=4+4)
Logger_large-10     1.11µs ± 0%    0.41µs ± 0%  -62.92%  (p=0.029 n=4+4)

name              old alloc/op   new alloc/op   delta
Logger_small-10       272B ± 0%       96B ± 0%  -64.71%  (p=0.029 n=4+4)
Logger_medium-10      904B ± 0%      312B ± 0%  -65.49%  (p=0.029 n=4+4)
Logger_large-10     2.47kB ± 0%    1.26kB ± 0%  -49.15%  (p=0.029 n=4+4)

name              old allocs/op  new allocs/op  delta
Logger_small-10       3.00 ± 0%      2.00 ± 0%  -33.33%  (p=0.029 n=4+4)
Logger_medium-10      7.00 ± 0%      5.00 ± 0%  -28.57%  (p=0.029 n=4+4)
Logger_large-10       19.0 ± 0%      10.0 ± 0%  -47.37%  (p=0.029 n=4+4)
```