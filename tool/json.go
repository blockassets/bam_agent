package tool

import (
	"log"
	"math/rand"
	"time"
	"unsafe"

	"github.com/json-iterator/go"
)

/*
	The code in this file is all tested indirectly in the rest of the
	application. For some reason, tests can't be easily written. They
	just don't seem to work correctly with jsoniter. Maybe something
	to do with package scope and the registration.
*/

type TimeDuration struct {
	time.Duration
}

func RegisterJsonTypes() {
	RegisterTimeDuration()
	RegisterRandomDuration()
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

// Type insurance
var _ jsoniter.ValEncoder = &RandomDuration{}
var _ jsoniter.ValDecoder = &RandomDuration{}

type RandomDuration struct {
	time.Duration
}

func RegisterRandomDuration() {
	jsoniter.RegisterTypeEncoder("tool.RandomDuration", &RandomDuration{})
	jsoniter.RegisterTypeDecoder("tool.RandomDuration", &RandomDuration{})
}

func (codec *RandomDuration) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	val, err := time.ParseDuration(iter.ReadString())
	if err != nil {
		log.Fatal(err)
	}
	*((*time.Duration)(ptr)) = getRandomizedDuration(val)
}

func (codec *RandomDuration) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Duration)(ptr))
	stream.WriteString(ts.String())
}

func (codec *RandomDuration) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Duration)(ptr))
	return ts.Nanoseconds() == 0
}

/*
	Randomly add between 1-3600 seconds to a duration.
*/
func getRandomizedDuration(duration time.Duration) time.Duration {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return duration + time.Duration(r1.Intn(3600))*time.Second
}

/*
	Merges source into/over destination
*/
func Merge(src []byte, dst []byte) ([]byte, error) {
	var mergeMap map[string]interface{}
	err := jsoniter.Unmarshal(dst, &mergeMap)
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal(src, &mergeMap)
	if err != nil {
		return nil, err
	}
	return jsoniter.Marshal(mergeMap)
}