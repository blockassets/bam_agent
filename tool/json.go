package tool

import (
	"log"
	"time"
	"unsafe"

	"github.com/json-iterator/go"
)

type TimeDuration struct {
	time.Duration
}

func RegisterTimeDuration() {
	jsoniter.RegisterTypeEncoder("time.Duration", &TimeDuration{})
	jsoniter.RegisterTypeDecoder("time.Duration", &TimeDuration{})
}

func (codec *TimeDuration) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	val, err := time.ParseDuration(iter.ReadString())
	if err != nil {
		log.Fatal(err)
	}
	*((*time.Duration)(ptr)) = val
}

func (codec *TimeDuration) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Duration)(ptr))
	stream.WriteString(ts.String())
}

func (codec *TimeDuration) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Duration)(ptr))
	return ts.Nanoseconds() == 0
}
