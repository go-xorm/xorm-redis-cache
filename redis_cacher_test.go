package xormrediscache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"reflect"
	"testing"
	// "unsafe"
)

type Point struct {
	X, Y int
}

func TestSerializationGob(t *testing.T) {
	var network bytes.Buffer // Stand-in for the network.

	// Create an encoder and send a value.
	enc := gob.NewEncoder(&network)
	err := enc.Encode(&Point{3, 4})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&network)
	var v Point
	err = dec.Decode(&v)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	log.Println(v)
}

func TestSerializationStruct(t *testing.T) {

	point := Point{X: 100, Y: -100}

	_, err := serialize(point)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}

var (
	TEST_POINT = Point{X: 100, Y: -100}
)

func serializePoint() ([]byte, error) {
	var point Point = TEST_POINT
	return serialize(&point)
}

func TestSerializationPtr1(t *testing.T) {

	bytes, err := serializePoint()

	if err != nil {
		t.Error(err)
	}

	ptr, err := deserialize(bytes)
	if err != nil {
		t.Error(err)
	}

	log.Println(ptr, "ptr type:", reflect.TypeOf(ptr))

	if reflect.TypeOf(ptr).Kind() == reflect.Struct {

		if ptr != TEST_POINT {
			t.Error(fmt.Errorf("decoded value:%v not identical to value:%v", ptr, TEST_POINT))
			t.FailNow()
		}
		t.Error(fmt.Errorf("deserialize func should return pointer of a struct"))
		t.FailNow()
	}

	ptrElem := reflect.ValueOf(ptr).Elem().Interface()
	log.Println(ptrElem, "elem type:", reflect.TypeOf(ptrElem), "can addr", reflect.ValueOf(ptrElem).CanAddr())
	// if ptrElem != point {
	if ptrElem != TEST_POINT {
		t.Error(fmt.Errorf("decoded value:%v not identical to value:%v", ptrElem, TEST_POINT))
		t.FailNow()
	}

	points := []Point{}
	points = append(points, ptrElem.(Point))

	log.Println("points:%v", points)

	// !nashtsai! how to make following compile?
	// pointPtrSlice := []*Point{}
	// pointPtrSlice = append(pointPtrSlice, ptr.(*Point))

	// datas := reflect.ValueOf(ptr).Elem().InterfaceData()
	// fmt.Println("data:", datas[0], datas[1])

	// sp := reflect.NewAt(reflect.TypeOf(ptrElem),
	// 	unsafe.Pointer(datas[1])).Interface()
	// fmt.Println("sp:", sp, sp.(*Point))

	pointPtrSlice := []*Point{}
	pointPtrSlice = append(pointPtrSlice, ptr.(*Point))
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
