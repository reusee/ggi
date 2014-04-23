package main

/*
#include <girepository.h>
#cgo pkg-config: gobject-introspection-1.0
gboolean T = TRUE;
gboolean F = FALSE;
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

var (
	p     = fmt.Printf
	True  = C.T
	False = C.F
)

func toGStr(s string) *C.gchar {
	return (*C.gchar)(unsafe.Pointer(C.CString(s)))
}

func fromGStr(s *C.gchar) string {
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

func asBaseInfo(p interface{}) *C.GIBaseInfo {
	return (*C.GIBaseInfo)(unsafe.Pointer(reflect.ValueOf(p).Pointer()))
}

func asCallableInfo(p interface{}) *C.GICallableInfo {
	return (*C.GICallableInfo)(unsafe.Pointer(reflect.ValueOf(p).Pointer()))
}

func asFunctionInfo(p interface{}) *C.GIFunctionInfo {
	return (*C.GIFunctionInfo)(unsafe.Pointer(reflect.ValueOf(p).Pointer()))
}

func argDirectionToString(d C.GIDirection) (ret string) {
	switch d {
	case C.GI_DIRECTION_IN:
		ret = "IN"
	case C.GI_DIRECTION_OUT:
		ret = "OUT"
	case C.GI_DIRECTION_INOUT:
		ret = "INOUT"
	default:
		panic("error")
	}
	return
}

func ArrayTypeToString(t C.GIArrayType) (ret string) {
	switch t {
	case C.GI_ARRAY_TYPE_C:
		ret = "C"
	case C.GI_ARRAY_TYPE_ARRAY:
		ret = "GArray"
	case C.GI_ARRAY_TYPE_PTR_ARRAY:
		ret = "GPtrArray"
	case C.GI_ARRAY_TYPE_BYTE_ARRAY:
		ret = "GByteArray"
	default:
		panic("error")
	}
	return
}
