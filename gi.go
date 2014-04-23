package main

/*
#include <girepository.h>
#cgo pkg-config: gobject-introspection-1.0
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

var p = fmt.Printf

func main() {
	repo := C.g_irepository_get_default()
	lib := os.Args[1]

	var err *C.GError
	C.g_irepository_require(repo, toGStr(lib), nil, 0, &err)
	if err != nil {
		log.Fatal(fromGStr(err.message))
	}

	nInfos := C.g_irepository_get_n_infos(repo, toGStr(lib))
	var nObjects int
	for i := C.gint(0); i < nInfos; i++ {
		baseInfo := C.g_irepository_get_info(repo, toGStr(lib), i)
		if C.g_base_info_is_deprecated(baseInfo) == C.gboolean(1) { // skip deprecated
			continue
		}
		name := fromGStr(C.g_base_info_get_name(baseInfo)) // name
		p("=> %s\n", name)
		var iter C.GIAttributeIter
		var key, value *C.char
		for C.g_base_info_iterate_attributes(baseInfo, &iter, &key, &value) == C.gboolean(1) {
			p("Attr %s = %s\n", C.GoString(key), C.GoString(value))
		}
		t := C.g_base_info_get_type(baseInfo) // type
		switch t {                            // dump
		case C.GI_INFO_TYPE_OBJECT:
			info := (*C.GIObjectInfo)(unsafe.Pointer(baseInfo))
			p("Object\n")
			DumpObjectInfo(info)
			nObjects++
		case C.GI_INFO_TYPE_STRUCT:
			p("Struct\n")
			info := (*C.GIStructInfo)(unsafe.Pointer(baseInfo))
			DumpStructInfo(info)
		case C.GI_INFO_TYPE_FLAGS:
			p("Flags\n")
			info := (*C.GIEnumInfo)(unsafe.Pointer(baseInfo))
			DumpEnumInfo(info)
		case C.GI_INFO_TYPE_CALLBACK:
			p("Callback\n")
			info := (*C.GIFunctionInfo)(unsafe.Pointer(baseInfo))
			DumpFunctionInfo(info)
		case C.GI_INFO_TYPE_INTERFACE:
			p("Interface\n")
			info := (*C.GIInterfaceInfo)(unsafe.Pointer(baseInfo))
			DumpInterfaceInfo(info)
		case C.GI_INFO_TYPE_UNION:
			p("Union\n")
			info := (*C.GIUnionInfo)(unsafe.Pointer(baseInfo))
			DumpUnionInfo(info)
		case C.GI_INFO_TYPE_ENUM:
			p("Enum\n")
			info := (*C.GIEnumInfo)(unsafe.Pointer(baseInfo))
			DumpEnumInfo(info)
		case C.GI_INFO_TYPE_FUNCTION:
			p("Function\n")
			info := (*C.GIFunctionInfo)(unsafe.Pointer(baseInfo))
			DumpFunctionInfo(info)
		case C.GI_INFO_TYPE_CONSTANT:
			p("Constant\n")
			info := (*C.GIConstantInfo)(unsafe.Pointer(baseInfo))
			DumpConstantInfo(info)
		default:
			panic(fmt.Sprintf("unknown type %d", t))
		}
		C.g_base_info_unref(baseInfo)
		p(strings.Repeat("-", 64))
		p("\n")
	}
	p("%d\n", nInfos)
	p("%d object types\n", nObjects)
}

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