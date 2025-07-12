module github.com/amirkode/go-mongr8

go 1.24.1

replace github.com/ONSdigital/dp-mongodb-in-memory => github.com/amirkode/dp-mongodb-in-memory v1.8.0

// internal module
// TODO: enable again with proper usecase
// require github.com/amirkode/go-mongr8/internal v0.0.0
// replace github.com/amirkode/go-mongr8/internal => ./internal

// external depedencies
require (
	github.com/golang/snappy v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.17.4
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)

require (
	github.com/ONSdigital/dp-mongodb-in-memory v1.8.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/smartystreets/goconvey v1.8.1
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.260.0 // indirect
	github.com/ONSdigital/dp-net/v2 v2.11.2 // indirect
	github.com/ONSdigital/log.go/v2 v2.4.3 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/smarty/assertions v1.15.1 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
)
