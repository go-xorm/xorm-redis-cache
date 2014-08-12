package xormrediscache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

type Point struct {
	X, Y int
}

func TestSerializationStruct(t *testing.T) {

	point := Point{X: 100, Y: -100}

	_, err := serialize(point)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestSerializationPtr(t *testing.T) {

	point := Point{X: 100, Y: -100}

	bytes, err := serialize(&point)

	if err != nil {
		t.Error(err)
	}

	ptr, err := deserialize(bytes)
	if err != nil {
		t.Error(err)
	}

	log.Println(ptr, "type:", reflect.TypeOf(ptr).Kind())

	if reflect.TypeOf(ptr).Kind() == reflect.Struct {
		t.Error(fmt.Errorf("deserialize func should return pointer of a struct"))
		t.FailNow()
	}
}

// func TestSerializationSliceOfPtr(t *testing.T) {

// 	points := []*Point{}
// 	for i := 1; i <= 3; i++ {
// 		points = append(points, &Point{3 * i, 4 * i})
// 	}

// 	bytes, err := serialize(points)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	ptr, err := deserialize(bytes)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
