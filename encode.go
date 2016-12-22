package tags

import (
	"fmt"
	"reflect"
	"strings"
)

type URLValue struct {
	TagName string
}

func (u *URLValue) Encode(a interface{}) string {
	v := reflect.ValueOf(a)

	fmt.Println(v.Kind())

	p := &printer{tagname: u.TagName, visited: make(map[visit]int), depth: 0}
	p.printValue(v)
	fmt.Printf("***** %#v", p)

	return strings.Join(p.kvs, "&")
}

// printValue must keep track of already-printed pointer values to avoid
// infinite recursion.
type visit struct {
	v   uintptr
	typ reflect.Type
}

type printer struct {
	key     []string
	kvs     []string
	tagname string
	visited map[visit]int
	depth   int
}

func (p *printer) append(kv string) {
	p.kvs = append(p.kvs, kv)
}

func (p *printer) printValue(v reflect.Value) {
	fmt.Println(p.depth, v.Kind(), v.Interface(), p.kvs, (p.key))
	if p.depth > 10 {
		p.append("!%v(DEPTH EXCEEDED)")
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.String:
		fmt.Println("==", p.key, p.kvs, v.Interface())
		p.append(strings.Join(p.key, ".") + "=" + fmt.Sprint(v.Interface()))
	case reflect.Map:
		//t := v.Type()
		//if showType {
		//	io.WriteString(p, t.String())
		//}
		//writeByte(p, '{')
		//if nonzero(v) {
		//	expand := !canInline(v.Type())
		//	pp := p
		//	if expand {
		//		writeByte(p, '\n')
		//		pp = p.indent()
		//	}
		//	keys := v.MapKeys()
		//	for i := 0; i < v.Len(); i++ {
		//		showTypeInStruct := true
		//		k := keys[i]
		//		mv := v.MapIndex(k)
		//		pp.printValue(k, false, true)
		//		writeByte(pp, ':')
		//		if expand {
		//			writeByte(pp, '\t')
		//		}
		//		showTypeInStruct = t.Elem().Kind() == reflect.Interface
		//		pp.printValue(mv, showTypeInStruct, true)
		//		if expand {
		//			io.WriteString(pp, ",\n")
		//		} else if i < v.Len()-1 {
		//			io.WriteString(pp, ", ")
		//		}
		//	}
		//	if expand {
		//		pp.tw.Flush()
		//	}
		//}
	case reflect.Struct:
		t := v.Type()
		if v.CanAddr() {
			addr := v.UnsafeAddr()
			vis := visit{addr, t}
			if vd, ok := p.visited[vis]; ok && vd < p.depth {
				p.append(t.String() + "{(CYCLIC REFERENCE)}")
				break // don't print v again
			}
			p.visited[vis] = p.depth
		}

		fmt.Println("---", t.String())

		if nonzero(v) {
			for i := 0; i < v.NumField(); i++ {
				f := t.Field(i)
				if tag, ok := f.Tag.Lookup(p.tagname); ok && tag != "" {
					n := len(p.key)
					p.key = append(p.key, tag)
					p.printValue(getField(v, i))
					p.key = p.key[0:n]
				}
			}
		}
	case reflect.Array, reflect.Slice:
		fmt.Println("array, slice")
	case reflect.Ptr:
		e := v.Elem()
		if e.IsValid() {
			pp := *p
			pp.depth++
			pp.printValue(e)
		}
	}
}

//
//
//func (p *printer) printInline(v reflect.Value, x interface{}) {
//	fmt.Fprintf(p, "(%v)", x)
//}

func getField(v reflect.Value, i int) reflect.Value {
	val := v.Field(i)
	if val.Kind() == reflect.Interface && !val.IsNil() {
		val = val.Elem()
	}
	return val
}

func nonzero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() != complex(0, 0)
	case reflect.String:
		return v.String() != ""
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if nonzero(getField(v, i)) {
				return true
			}
		}
		return false
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if nonzero(v.Index(i)) {
				return true
			}
		}
		return false
	case reflect.Map, reflect.Interface, reflect.Slice, reflect.Ptr, reflect.Chan, reflect.Func:
		return !v.IsNil()
	case reflect.UnsafePointer:
		return v.Pointer() != 0
	}
	return true
}

func labelType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface, reflect.Struct:
		return true
	}
	return false
}
func canInline(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map:
		return !canExpand(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if canExpand(t.Field(i).Type) {
				return false
			}
		}
		return true
	case reflect.Interface:
		return false
	case reflect.Array, reflect.Slice:
		return !canExpand(t.Elem())
	case reflect.Ptr:
		return false
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return false
	}
	return true
}

func canExpand(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map, reflect.Struct,
		reflect.Interface, reflect.Array, reflect.Slice,
		reflect.Ptr:
		return true
	}
	return false
}
