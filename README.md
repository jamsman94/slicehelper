## slicehelper
A simple go tool to make json Marshalling less painful for slices

## install
```
go get -u github.com/jamsman94/slicehelper
```

## Usage

Replacing All nil slices with initialized empty slices for json marshal

This currently supports custom structs, pointers, slices and maps, and the combination of above types.

```go
import "github.com/jamsman94/slicehelper"

var myslice []string                               // nil

slicehelper.ReplaceNilWithEmptySlice(myslice)      // []

type SupportedType struct {
    PointerToSlice   *[]string           `json:"pointer_to_slice"`
    OnlySlice        []string            `json:"only_string"`
    SliceOfPointer   []*string           `json:"slice_of_pointer"`
    MapofSlice       map[string][]string `json:"map_of_slice"`
}
}
```