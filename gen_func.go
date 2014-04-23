package main

/*
#include <girepository.h>
*/
import "C"

func (self *Generator) GenFunction(info *C.GIFunctionInfo) {

	// symbol
	symbol := fromGStr(C.g_function_info_get_symbol(info))
	p("SYMBOL %s\n", symbol)

	// error
	throwsError := C.g_callable_info_can_throw_gerror(asCallableInfo(info)) == True
	if throwsError {
		p("THROWS ERROR\n")
	}

	// args
	nArgs := C.g_callable_info_get_n_args(asCallableInfo(info))
	for i := C.gint(0); i < nArgs; i++ {
		arg := C.g_callable_info_get_arg(asCallableInfo(info), i)
		argName := fromGStr(C.g_base_info_get_name(asBaseInfo(arg)))
		argType := C.g_arg_info_get_type(arg)
		argTypeTag := C.g_type_info_get_tag(argType)
		argTypeTagRepr := fromGStr(C.g_type_tag_to_string(argTypeTag))
		argDirection := C.g_arg_info_get_direction(arg)
		p("%s %s %s", argName, argTypeTagRepr, argDirectionToString(argDirection))
		argIsPointer := C.g_type_info_is_pointer(argType) == True
		if argIsPointer {
			p(" POINTER")
		}
		argMayBeNull := C.g_arg_info_may_be_null(arg) == True
		if argMayBeNull {
			p(" NULLABLE")
		}
		argCallerAllocates := C.g_arg_info_is_caller_allocates(arg) == True
		if argCallerAllocates {
			p(" CALLER_ALLOCATES")
		}
		argIsOptional := C.g_arg_info_is_optional(arg) == True
		if argIsOptional { //TODO will it occur?
			p(" OPTIONAL")
		}
		argSkip := C.g_arg_info_is_skip(arg) == True
		if argSkip { //TODO will it occur?
			p(" SKIP")
		}
		if argTypeTag == C.GI_TYPE_TAG_INTERFACE {
			ifaceInfo := C.g_type_info_get_interface(argType)
			if ifaceInfo != nil {
				ifaceType := C.g_base_info_get_type(ifaceInfo)
				ifaceName := fromGStr(C.g_base_info_get_name(ifaceInfo))
				p(" INTERFACE %s %s", fromGStr(C.g_info_type_to_string(ifaceType)), ifaceName)
			}
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

	// return
	returnType := C.g_callable_info_get_return_type(asCallableInfo(info))
	_ = returnType
	//TODO same as arg type
	mayReturnNull := C.g_callable_info_may_return_null(asCallableInfo(info)) == True
	if mayReturnNull {
		p("MAY_RETURN_NULL\n")
	}
	skipReturn := C.g_callable_info_skip_return(asCallableInfo(info)) == True
	if skipReturn {
		p("SKIP_RETURN\n")
	}

	p("\n")
}
