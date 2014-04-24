package main

/*
#include <girepository.h>
#include <gobject/gobject.h>
*/
import "C"
import (
	"strings"
)

func (self *Generator) GenFunction(info *C.GIFunctionInfo) {
	// name
	name := fromGStr(C.g_base_info_get_name(asBaseInfo(info)))
	self.currentFunctionName = name
	goName := convertFuncName(name)
	w(self.FuncsWriter, "func %s", goName)

	// collect args
	var inArgs, outArgs []*C.GIArgInfo
	nArgs := C.g_callable_info_get_n_args(asCallableInfo(info))
	for i := C.gint(0); i < nArgs; i++ {
		arg := C.g_callable_info_get_arg(asCallableInfo(info), i)
		argDirection := C.g_arg_info_get_direction(arg)
		switch argDirection {
		case C.GI_DIRECTION_IN:
			inArgs = append(inArgs, arg)
		case C.GI_DIRECTION_OUT:
			outArgs = append(outArgs, arg)
		case C.GI_DIRECTION_INOUT:
			inArgs = append(inArgs, arg)
			outArgs = append(outArgs, arg)
		}
	}

	// output in args
	w(self.FuncsWriter, "(")
	for _, arg := range inArgs {
		argName := convertArgName(fromGStr(C.g_base_info_get_name(asBaseInfo(arg))))
		goTypeStr := self.getGoTypeStr(C.g_arg_info_get_type(arg))
		/*
		argCallerAllocates := C.g_arg_info_is_caller_allocates(arg) == True
		if argCallerAllocates {
			goTypeStr = "*" + goTypeStr
		}
		*/
		w(self.FuncsWriter, "%s %s,", argName, goTypeStr)
	}
	w(self.FuncsWriter, ")")

	w(self.FuncsWriter, "{}\n") //TODO body

	// symbol
	symbol := fromGStr(C.g_function_info_get_symbol(info))
	_ = symbol

	// error
	throwsError := C.g_callable_info_can_throw_gerror(asCallableInfo(info)) == True
	_ = throwsError

	// args
	/*
		nArgs := C.g_callable_info_get_n_args(asCallableInfo(info))
		for i := C.gint(0); i < nArgs; i++ {
			argCallerAllocates := C.g_arg_info_is_caller_allocates(arg) == True
			if argCallerAllocates {
				p(" CALLER_ALLOCATES")
			}
			if argTypeTag == C.GI_TYPE_TAG_ARRAY {
				arrayType := C.g_type_info_get_array_type(argType)
				p(" ARRAY of %s", ArrayTypeToString(arrayType))
				if arrayType == C.GI_ARRAY_TYPE_C {
					elemType := C.g_type_info_get_param_type(argType, 0)
					elemTypeTag := C.g_type_info_get_tag(elemType)
					p(" %s", fromGStr(C.g_type_tag_to_string(elemTypeTag)))
				}
			}
			p("\n")
		}
	*/

	// return
	returnType := C.g_callable_info_get_return_type(asCallableInfo(info))
	_ = returnType
	//TODO same as arg type
	mayReturnNull := C.g_callable_info_may_return_null(asCallableInfo(info)) == True
	_ = mayReturnNull
	skipReturn := C.g_callable_info_skip_return(asCallableInfo(info)) == True
	_ = skipReturn

}

func convertFuncName(cName string) string {
	parts := strings.Split(cName, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func convertArgName(name string) string {
	if isGoKeyword(name) {
		return name + "_"
	}
	return name
}

func (self *Generator) getGoTypeStr(typeInfo *C.GITypeInfo) (ret string) {
	argTypeTag := C.g_type_info_get_tag(typeInfo)
	// map to basic type
	ret = tagMapToGoType(argTypeTag)
	if ret != "" {
		if C.g_type_info_is_pointer(typeInfo) == True && ret != "string" { // is pointer
			ret = "*" + ret
		}
		return
	}
	// complex type
	if argTypeTag == C.GI_TYPE_TAG_INTERFACE {
		ifaceInfo := C.g_type_info_get_interface(typeInfo) //GIBaseInfo
		if ifaceInfo != nil {
			ifaceType := C.g_base_info_get_type(ifaceInfo) //GIInfoType
			ifaceName := fromGStr(C.g_base_info_get_name(ifaceInfo))
			switch ifaceType {
			case C.GI_INFO_TYPE_STRUCT:
				/*
				ret = "*C." + strings.ToUpper(self.CPrefix) + ifaceName
				p("%s %s\n", self.currentFunctionName, ifaceName)
				return
				*/
				//TODO we need c:type here. girepository is not exposing.
			default:
				//p("INTERFACE %s %s\n", fromGStr(C.g_info_type_to_string(ifaceType)), ifaceName)
			}
		}
	}
	// fallback
	return "interface{}"
}

func tagMapToGoType(tag C.GITypeTag) string {
	switch tag {
	case C.GI_TYPE_TAG_BOOLEAN:
		return "bool"
	case C.GI_TYPE_TAG_INT8:
		return "int8"
	case C.GI_TYPE_TAG_UINT8:
		return "uint8"
	case C.GI_TYPE_TAG_INT16:
		return "int16"
	case C.GI_TYPE_TAG_UINT16:
		return "uint16"
	case C.GI_TYPE_TAG_INT32:
		return "int32"
	case C.GI_TYPE_TAG_UINT32:
		return "uint32"
	case C.GI_TYPE_TAG_INT64:
		return "int64"
	case C.GI_TYPE_TAG_UINT64:
		return "uint64"
	case C.GI_TYPE_TAG_FLOAT:
		return "float32"
	case C.GI_TYPE_TAG_DOUBLE:
		return "float64"
	case C.GI_TYPE_TAG_UTF8, C.GI_TYPE_TAG_FILENAME:
		return "string"
	case C.GI_TYPE_TAG_UNICHAR:
		return "rune"
	}
	return ""
}
