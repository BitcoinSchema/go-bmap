module github.com/bitcoinschema/go-bmap

go 1.22

toolchain go1.22.5

require (
	github.com/bitcoin-sv/go-sdk v1.1.7
	github.com/bitcoinschema/go-aip v0.2.3
	github.com/bitcoinschema/go-b v0.1.1
	github.com/bitcoinschema/go-bap v0.3.3
	github.com/bitcoinschema/go-bob v0.4.3
	github.com/bitcoinschema/go-boost v0.1.0
	github.com/bitcoinschema/go-bpu v0.1.3
	github.com/bitcoinschema/go-map v0.1.1
	github.com/bitcoinschema/go-sigma v0.0.2
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.24.0 // indirect
)

replace github.com/bitcoinschema/go-aip => ../go-aip

replace github.com/bitcoinschema/go-b => ../go-b

replace github.com/bitcoinschema/go-bap => ../go-bap

replace github.com/bitcoinschema/go-bob => ../go-bob

replace github.com/bitcoinschema/go-boost => ../go-boost

replace github.com/bitcoinschema/go-bpu => ../go-bpu

replace github.com/bitcoinschema/go-map => ../go-map

replace github.com/bitcoinschema/go-sigma => ../go-sigma
