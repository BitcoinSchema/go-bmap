## bmap-go

A Go version of the [bmap js library](https://github.com/bitcoinschema/bmap/)

## Usage

Transform from [BOB](https://github.com/bitcoinschema/go-bob) Tx to BMAP format

```go
import "github.com/bitcoinschema/go-bmap"

// Transform from BOB to BMAP
bmapTx, err := bmapData.NewFromString(bobData)
```

## Supported protocols

- [AIP](https://github.com/bitcoinschema/go-aip)
- [BAP](https://github.com/bitcoinschema/go-bap)
- [MAP](https://github.com/bitcoinschema/go-map)

## Example

Preparing for mongo insert:

```go

bsonData := bson.M{
  "tx":  bobData.Tx,
  "in":  bobData.In,
  "out": bobData.Out,
  "blk": bobData.Blk,
}

if bmapData.AIP != nil {
  bsonData["AIP"] = bmapData.AIP
}

if bmapData.BAP != nil {
  bsonData["BAP"] = bmapData.BAP
}

if bmapData.MAP != nil {
  bsonData["MAP"] = bmapData.MAP
}

_, err := conn.InsertOne(collectionName, bsonData)
```

## ToDo

- B support
