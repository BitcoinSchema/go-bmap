## bmap-go

A Go version of the [bmap js library](https://github.com/rohenaz/bmap/)

## Usage

Transform from [BOB](https://github.com/rohenaz/go-bob) Tx to BMAP format

```go
import "github.com/rohenaz/go-bmap"

// Transform from BOB to BMAP
bmapData := bmap.New()
bmapData.FromBob(bobData)
```

## Supported protocols

- [AIP](https://github.com/rohenaz/go-aip)
- [BAP](https://github.com/rohenaz/go-bap)
- [MAP](https://github.com/rohenaz/go-map)

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
