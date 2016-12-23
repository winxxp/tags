package tags

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type URLValue struct {
	TagName string
}

func (u *URLValue) Encode(a interface{}) string {
	v := reflect.ValueOf(a)
	p := &printer{
		tagname: u.TagName,
		kvs:     make(url.Values),
		visited: make(map[visit]int),
		depth:   0,
	}
	p.printValue(v)

	return p.encode()
}

type visit struct {
	v   uintptr
	typ reflect.Type
}

type printer struct {
	key     bytes.Buffer
	kvs     url.Values
	tagname string
	visited map[visit]int
	depth   int
}

func (p *printer) encode() string {
	return p.kvs.Encode()
}

func (p *printer) printValue(v reflect.Value) {
	if p.depth > 10 {
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		p.kvs.Set(p.key.String(), strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.kvs.Set(p.key.String(), strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.kvs.Set(p.key.String(), strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		p.kvs.Set(p.key.String(), strconv.FormatFloat(v.Float(), 'G', -1, 64))
	case reflect.Complex64, reflect.Complex128:
		p.kvs.Set(p.key.String(), fmt.Sprint(v.Interface()))
	case reflect.String:
		p.kvs.Set(p.key.String(), v.String())
	case reflect.Map:
		keys := v.MapKeys()
		for i := 0; i < v.Len(); i++ {
			k := keys[i]
			n := p.key.Len()
			p.printKey(fmt.Sprint(k), n)
			p.printValue(v.MapIndex(k))
			p.key.Truncate(n)
		}
	case reflect.Struct:
		t := v.Type()
		if v.CanAddr() {
			addr := v.UnsafeAddr()
			vis := visit{addr, t}
			if vd, ok := p.visited[vis]; ok && vd < p.depth {
				break // don't print v again
			}
			p.visited[vis] = p.depth
		}

		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			if tag, ok := f.Tag.Lookup(p.tagname); ok && tag != "" {
				n := p.key.Len()
				p.printKey(tag, n)
				p.printValue(getField(v, i))
				p.key.Truncate(n)
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			n := p.key.Len()
			p.key.WriteByte('[')
			p.key.WriteString(strconv.Itoa(i))
			p.key.WriteByte(']')
			p.printValue(v.Index(i))
			p.key.Truncate(n)
		}

	case reflect.Ptr:
		e := v.Elem()
		if e.IsValid() {
			pp := *p
			pp.depth++
			pp.printValue(e)
		}
	}
}

func (p *printer) printKey(key string, n int) {
	if n > 0 {
		p.key.WriteByte('.')
	}
	p.key.WriteString(key)
}

func getField(v reflect.Value, i int) reflect.Value {
	val := v.Field(i)

	if val.Kind() == reflect.Interface && !val.IsNil() {
		val = val.Elem()
	}

	return val
}
