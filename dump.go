package main

/*
#include <girepository.h>
*/
import "C"
import "unsafe"

func DumpObjectInfo(info *C.GIObjectInfo) {
	isAbstract := C.g_object_info_get_abstract(info) == C.gboolean(1)
	p("is abstract %v\n", isAbstract)
	parent := C.g_object_info_get_parent(info)
	if parent != nil {
		parentName := fromGStr(C.g_base_info_get_name(asBaseInfo(parent)))
		p("parent %s\n", parentName)
	}
	typeName := fromGStr(C.g_object_info_get_type_name(info))
	p("type name %s\n", typeName)
	typeInit := fromGStr(C.g_object_info_get_type_init(info))
	p("type init %s\n", typeInit)
	nConsts := C.g_object_info_get_n_constants(info)
	p("%d consts\n", nConsts)
	for i := C.gint(0); i < nConsts; i++ {
		constInfo := C.g_object_info_get_constant(info, i)
		DumpConstantInfo(constInfo)
	}
	nFields := C.g_object_info_get_n_fields(info)
	p("%d fields\n", nFields)
	for i := C.gint(0); i < nFields; i++ {
		field := C.g_object_info_get_field(info, i)
		DumpFieldInfo(field)
	}
	nInterfaces := C.g_object_info_get_n_interfaces(info)
	p("%d interfaces\n", nInterfaces)
	for i := C.gint(0); i < nInterfaces; i++ {
		interf := C.g_object_info_get_interface(info, i)
		DumpInterfaceInfo(interf)
	}
	nMethods := C.g_object_info_get_n_methods(info)
	p("%d methods\n", nMethods)
	for i := C.gint(0); i < nMethods; i++ {
		f := C.g_object_info_get_method(info, i)
		DumpFunctionInfo(f)
	}
	nProperties := C.g_object_info_get_n_properties(info)
	p("%d properties\n", nProperties)
	for i := C.gint(0); i < nProperties; i++ {
		property := C.g_object_info_get_property(info, i)
		DumpPropertyInfo(property)
	}
	nSignals := C.g_object_info_get_n_signals(info)
	p("%d signals\n", nSignals)
	for i := C.gint(0); i < nSignals; i++ {
		signal := C.g_object_info_get_signal(info, i)
		DumpSignalInfo(signal)
	}
	nVFuncs := C.g_object_info_get_n_vfuncs(info)
	p("%d vfuncs\n", nVFuncs)
	for i := C.gint(0); i < nVFuncs; i++ {
		vfunc := C.g_object_info_get_vfunc(info, i)
		DumpVFuncInfo(vfunc)
	}
}

func DumpStructInfo(info *C.GIStructInfo) {
	align := C.g_struct_info_get_alignment(info)
	p("alignment %d bytes\n", align)
	size := C.g_struct_info_get_size(info)
	p("size %d bytes\n", size)
	isGtypeStruct := C.g_struct_info_is_gtype_struct(info) == C.gboolean(1)
	p("is gtype struct %v\n", isGtypeStruct)
	isForeign := C.g_struct_info_is_foreign(info) == C.gboolean(1)
	p("is foreign %v\n", isForeign)
	nFields := C.g_struct_info_get_n_fields(info)
	p("%d fields\n", nFields)
	for i := C.gint(0); i < nFields; i++ {
		field := C.g_struct_info_get_field(info, i)
		DumpFieldInfo(field)
	}
	nMethods := C.g_struct_info_get_n_methods(info)
	p("%d methods\n", nMethods)
	for i := C.gint(0); i < nMethods; i++ {
		f := C.g_struct_info_get_method(info, i)
		DumpFunctionInfo(f)
	}
}

func DumpEnumInfo(info *C.GIEnumInfo) {
	nValues := C.g_enum_info_get_n_values(info)
	p("%d values\n", nValues)
	for i := C.gint(0); i < nValues; i++ {
		valueInfo := C.g_enum_info_get_value(info, i)
		DumpValueInfo(valueInfo)
	}
	nMethods := C.g_enum_info_get_n_methods(info)
	p("%d methods\n", nMethods)
	for i := C.gint(0); i < nMethods; i++ {
		f := C.g_enum_info_get_method(info, i)
		DumpFunctionInfo(f)
	}
	storageType := C.g_enum_info_get_storage_type(info)
	p("%s\n", TypeTagGetName(storageType))
	errorDomain := fromGStr(C.g_enum_info_get_error_domain(info))
	p("error domain %s\n", errorDomain)
}

func DumpValueInfo(valueInfo *C.GIValueInfo) {
	name := fromGStr(C.g_base_info_get_name(asBaseInfo(valueInfo)))
	value := C.g_value_info_get_value(valueInfo)
	p("%s = %v\n", name, value)
}

func DumpFunctionInfo(info *C.GIFunctionInfo) {
	DumpCallableInfo(asCallableInfo(info))
	flags := C.g_function_info_get_flags(info)
	if flags&C.GI_FUNCTION_IS_METHOD > 0 {
		p("is method\n")
	}
	if flags&C.GI_FUNCTION_IS_CONSTRUCTOR > 0 {
		p("is constructor\n")
	}
	if flags&C.GI_FUNCTION_IS_GETTER > 0 {
		p("is getter\n")
		property := C.g_function_info_get_property(info)
		DumpPropertyInfo(property)
	}
	if flags&C.GI_FUNCTION_IS_SETTER > 0 {
		p("is setter\n")
		property := C.g_function_info_get_property(info)
		DumpPropertyInfo(property)
	}
	if flags&C.GI_FUNCTION_WRAPS_VFUNC > 0 {
		p("wraps vfunc\n")
		_ = C.g_function_info_get_vfunc(info)
	}
	if flags&C.GI_FUNCTION_THROWS > 0 {
		p("throws error\n")
	}
	symbol := fromGStr(C.g_function_info_get_symbol(info))
	p("symbol %s\n", symbol)
}

func DumpCallableInfo(info *C.GICallableInfo) {
	throwsError := C.g_callable_info_can_throw_gerror(info) == C.gboolean(1)
	p("can throws error %v\n", throwsError)
	nArgs := C.g_callable_info_get_n_args(info)
	for i := C.gint(0); i < nArgs; i++ {
		argInfo := C.g_callable_info_get_arg(info, i)
		DumpArgInfo(argInfo)
	}
	returnOwnership := C.g_callable_info_get_caller_owns(info)
	p("return value ownership %s\n", TransferGetName(returnOwnership))
	returnType := C.g_callable_info_get_return_type(info)
	defer C.g_base_info_unref(asBaseInfo(returnType))
	p("return type %v\n", returnType)
	DumpTypeInfo(returnType)
	isMethod := C.g_callable_info_is_method(info) == C.gboolean(1)
	p("is method %v\n", isMethod)
	var iter C.GIAttributeIter
	var key, value *C.char
	for C.g_callable_info_iterate_return_attributes(info, &iter, &key, &value) == C.gboolean(1) {
		p("Attr %s = %s\n", C.GoString(key), C.GoString(value))
	}
	mayReturnNull := C.g_callable_info_may_return_null(info) == C.gboolean(1)
	p("may return null %v\n", mayReturnNull)
	skipReturn := C.g_callable_info_skip_return(info) == C.gboolean(1)
	p("skip return %v\n", skipReturn)
}

func DumpInterfaceInfo(info *C.GIInterfaceInfo) {
	nBase := C.g_interface_info_get_n_prerequisites(info)
	p("%d prerequisites\n", nBase)
	for i := C.gint(0); i < nBase; i++ {
		base := (*C.GIInterfaceInfo)(unsafe.Pointer(C.g_interface_info_get_prerequisite(info, i)))
		DumpInterfaceInfo(base)
	}
	nProperties := C.g_interface_info_get_n_properties(info)
	p("%d properties\n", nProperties)
	for i := C.gint(0); i < nProperties; i++ {
		property := C.g_interface_info_get_property(info, i)
		DumpPropertyInfo(property)
	}
	nMethods := C.g_interface_info_get_n_methods(info)
	p("%d methods\n", nMethods)
	for i := C.gint(0); i < nMethods; i++ {
		method := C.g_interface_info_get_method(info, i)
		DumpFunctionInfo(method)
	}
	nSignals := C.g_interface_info_get_n_signals(info)
	p("%d signals\n", nSignals)
	for i := C.gint(0); i < nSignals; i++ {
		signal := C.g_interface_info_get_signal(info, i)
		DumpSignalInfo(signal)
	}
	nVFuncs := C.g_interface_info_get_n_vfuncs(info)
	p("%d vfuncs\n", nVFuncs)
	for i := C.gint(0); i < nVFuncs; i++ {
		vfunc := C.g_interface_info_get_vfunc(info, i)
		DumpVFuncInfo(vfunc)
	}
	nConstants := C.g_interface_info_get_n_constants(info)
	p("%d constants\n", nConstants)
	for i := C.gint(0); i < nConstants; i++ {
		constant := C.g_interface_info_get_constant(info, i)
		DumpConstantInfo(constant)
	}
	structInfo := C.g_interface_info_get_iface_struct(info)
	if structInfo != nil {
		DumpStructInfo(structInfo)
	}
}

func DumpFieldInfo(info *C.GIFieldInfo) {
	flags := C.g_field_info_get_flags(info)
	if flags&C.GI_FIELD_IS_READABLE > 0 {
		p("readable\n")
	}
	if flags&C.GI_FIELD_IS_WRITABLE > 0 {
		p("writable\n")
	}
	offset := C.g_field_info_get_offset(info)
	p("offset %d\n", offset)
	size := C.g_field_info_get_size(info)
	p("size %d\n", size)
	t := C.g_field_info_get_type(info)
	DumpTypeInfo(t)
}

func DumpPropertyInfo(info *C.GIPropertyInfo) {
	//TODO
}

func DumpSignalInfo(info *C.GISignalInfo) {
	//TODO
}

func DumpVFuncInfo(info *C.GIVFuncInfo) {
	//TODO
}

func DumpArgInfo(info *C.GIArgInfo) {
	//TODO
}

func DumpTypeInfo(info *C.GITypeInfo) {
	//TODO
}

func DumpUnionInfo(info *C.GIUnionInfo) {
	nFields := C.g_union_info_get_n_fields(info)
	p("%d fields\n", nFields)
	for i := C.gint(0); i < nFields; i++ {
		field := C.g_union_info_get_field(info, i)
		DumpFieldInfo(field)
	}
	nMethods := C.g_union_info_get_n_methods(info)
	p("%d methods\n", nMethods)
	for i := C.gint(0); i < nMethods; i++ {
		method := C.g_union_info_get_method(info, i)
		DumpFunctionInfo(method)
	}
	isDiscriminated := C.g_union_info_is_discriminated(info) == C.gboolean(1)
	p("is discriminated %v\n", isDiscriminated)
	if isDiscriminated {
		offset := C.g_union_info_get_discriminator_offset(info)
		p("discriminated offset %d\n", offset)
		discriminatedType := C.g_union_info_get_discriminator_type(info)
		p("discriminated type %d\n", discriminatedType)
		DumpTypeInfo(discriminatedType)
		for i := C.gint(0); i < nFields; i++ {
			discriminator := C.g_union_info_get_discriminator(info, i)
			DumpConstantInfo(discriminator)
		}
	}
	size := C.g_union_info_get_size(info)
	p("size %d bytes\n", size)
	align := C.g_union_info_get_alignment(info)
	p("alignment %d bytes\n", align)
}

func DumpConstantInfo(info *C.GIConstantInfo) {
	var value C.GIArgument
	C.g_constant_info_get_value(info, &value)
	defer C.g_constant_info_free_value(info, &value)
	p("value %v\n", value)
	t := C.g_constant_info_get_type(info)
	DumpTypeInfo(t)
}

func TypeTagGetName(t C.GITypeTag) (ret string) {
	return map[C.GITypeTag]string{
		C.GI_TYPE_TAG_VOID:      "void",
		C.GI_TYPE_TAG_BOOLEAN:   "bool",
		C.GI_TYPE_TAG_INT8:      "int8",
		C.GI_TYPE_TAG_UINT8:     "uint8",
		C.GI_TYPE_TAG_INT16:     "int16",
		C.GI_TYPE_TAG_UINT16:    "uint16",
		C.GI_TYPE_TAG_INT32:     "int32",
		C.GI_TYPE_TAG_UINT32:    "uint32",
		C.GI_TYPE_TAG_INT64:     "int64",
		C.GI_TYPE_TAG_UINT64:    "uint64",
		C.GI_TYPE_TAG_FLOAT:     "float",
		C.GI_TYPE_TAG_DOUBLE:    "double",
		C.GI_TYPE_TAG_GTYPE:     "gtype",
		C.GI_TYPE_TAG_UTF8:      "utf8",
		C.GI_TYPE_TAG_FILENAME:  "filename",
		C.GI_TYPE_TAG_ARRAY:     "array",
		C.GI_TYPE_TAG_INTERFACE: "interface",
		C.GI_TYPE_TAG_GLIST:     "glist",
		C.GI_TYPE_TAG_GSLIST:    "gslist",
		C.GI_TYPE_TAG_GHASH:     "ghash",
		C.GI_TYPE_TAG_ERROR:     "error",
		C.GI_TYPE_TAG_UNICHAR:   "unichar",
	}[t]
}

func TransferGetName(t C.GITransfer) string {
	return map[C.GITransfer]string{
		C.GI_TRANSFER_NOTHING:    "nothing",
		C.GI_TRANSFER_CONTAINER:  "container",
		C.GI_TRANSFER_EVERYTHING: "everything",
	}[t]
}
