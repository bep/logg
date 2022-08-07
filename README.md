
[![Tests on Linux, MacOS and Windows](https://github.com/bep/log/workflows/Test/badge.svg)](https://github.com/bep/log/actions?query=workflow:Test)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/log)](https://goreportcard.com/report/github.com/bep/log)
[![GoDoc](https://godoc.org/github.com/bep/log?status.svg)](https://godoc.org/github.com/bep/log)

This is a fork of the exellent [Apex Log](https://github.com/apex/log) library.

Main changes:

* Trim unneeded dependencies.
* Make `Fields` into a slice to preserve log order.
* Split `Entry` into `Entry` and `EntryFields`. This is easier to reason about and more effective.
* Rework the logger interface to allow lazy creation of messages, e.g. `Info(fmt.Stringer)`.