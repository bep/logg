module github.com/bep/logg/benchmarks

go 1.24.0

replace github.com/bep/logg => ../

require (
	github.com/apex/log v1.9.0
	github.com/go-kit/log v0.2.0
	github.com/rs/zerolog v1.26.0
	github.com/sirupsen/logrus v1.8.1
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.19.1
	gopkg.in/inconshreveable/log15.v2 v2.0.0-20200109203555-b30bc20e4fd1
)

require (
	github.com/benbjohnson/clock v1.2.0 // indirect
	github.com/bep/clocks v0.5.0 // indirect
	github.com/bep/logg v0.0.0-20220809094309-f3eda2566f97
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	golang.org/x/sys v0.40.0 // indirect
)
