package xormrediscache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"reflect"
	"testing"
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

	log.Println(ptr, "type:", reflect.TypeOf(ptr))

	if reflect.TypeOf(ptr).Kind() == reflect.Struct {

		if ptr != point {
			t.Error(fmt.Errorf("decoded value:%v not identical to value:%v", ptr, point))
			t.FailNow()
		}
		t.Error(fmt.Errorf("deserialize func should return pointer of a struct"))
		t.FailNow()
	}

	ptrElem := reflect.ValueOf(ptr).Elem().Interface()

	if ptrElem != point {
		t.Error(fmt.Errorf("decoded value:%v not identical to value:%v", ptrElem, point))
		t.FailNow()
	}
}

func TestSerializationPtr2(t *testing.T) {

	point := Point{X: 100, Y: -100}

	bytes, err := serialize(&point)

	if err != nil {
		t.Error(err)
	}

	pointDec := Point{}
	err = deserialize2(bytes, &pointDec)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	log.Println(pointDec, "type:", reflect.TypeOf(pointDec))

	if pointDec != point {
		t.Error(fmt.Errorf("decoded value:%v not identical to value:%v", pointDec, point))
		t.FailNow()
	}

	// if reflect.TypeOf(ptr).Kind() == reflect.Struct {
	// 	t.Error(fmt.Errorf("deserialize func should return pointer of a struct"))
	// 	t.FailNow()
	// }
}

func TestSerializationPtr3(t *testing.T) {

	point := &Point{X: 100, Y: -100}

	bytes, err := serialize3(point)

	if err != nil {
		t.Error(err)
	}

	ptr := &Point{}
	log.Println("b4:", ptr, "type:", reflect.TypeOf(ptr))
	err = deserialize3(bytes, ptr)
	if err != nil {
		t.Error(err)
	}

	log.Println(ptr, "type:", reflect.TypeOf(ptr))

	if reflect.TypeOf(ptr).Kind() == reflect.Struct {
		t.Error(fmt.Errorf("deserialize func should return pointer of a struct"))
		t.FailNow()
	}
}

func serialize3(value *Point) ([]byte, error) {

	err := RegisterGobConcreteType(value)
	if err != nil {
		return nil, err
	}

	// if reflect.TypeOf(value).Kind() == reflect.Struct {
	// 	return nil, fmt.Errorf("serialize func only take pointer of a struct")
	// }

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := interfaceEncode3(encoder, value); err != nil {
		log.Fatalf("[xorm/redis_cacher] gob encoding '%s' failed: %s", value, err)
		return nil, err
	}
	return b.Bytes(), nil
}

func interfaceEncode3(enc *gob.Encoder, p *Point) error {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.

	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type.  If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.

	log.Printf("[xorm/redis_cacher] interfaceEncode type:%v", reflect.TypeOf(p))
	err := enc.Encode(p)
	if err != nil {
		log.Fatal("[xorm/redis_cacher] encode:", err)
	}
	return err
}

func deserialize3(byt []byte, ptr *Point) (err error) {
	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)

	if err = interfaceDecode3(decoder, ptr); err != nil {
		log.Fatalf("[xorm/redis_cacher] gob decoding failed: %s", err)
		return
	}

	return
}

func interfaceDecode3(dec *gob.Decoder, ptr *Point) error {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.

	log.Printf("[xorm/redis_cacher] interfaceDecode3 type b4 decode:%v", reflect.TypeOf(ptr))

	err := dec.Decode(ptr)
	if err != nil {
		log.Fatal("[xorm/redis_cacher] ", err)
	}
	// log.Printf("[xorm/redis_cacher] interfaceDecode3 type:%v", reflect.TypeOf(p))
	// *ptr = p
	log.Printf("[xorm/redis_cacher] interfaceDecode3 type:%v", reflect.TypeOf(ptr))

	return err
}

func TestSerializationPtr4(t *testing.T) {

	point := &Point{X: 100, Y: -100}

	bytes, err := serialize(point)

	if err != nil {
		t.Error(err)
	}

	ptr := &Point{}
	log.Println("b4:", ptr, "type:", reflect.TypeOf(ptr))
	err = deserialize2(bytes, ptr)
	if err != nil {
		t.Error(err)
	}

	log.Println(ptr, "type:", reflect.TypeOf(ptr))

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
