package main

/*
#include <girepository.h>
*/
import "C"
import "log"

type Generator struct {
	Lib    string
	OutDir string
}

func Gen(lib, outDir string) {
	generator := &Generator{
		Lib:    lib,
		OutDir: outDir,
	}

	repo := C.g_irepository_get_default()
	var gerr *C.GError
	C.g_irepository_require(repo, toGStr(lib), nil, 0, &gerr)
	if gerr != nil {
		log.Fatal(fromGStr(gerr.message))
	}

	nInfos := C.g_irepository_get_n_infos(repo, toGStr(lib))
	for i := C.gint(0); i < nInfos; i++ {
		baseInfo := C.g_irepository_get_info(repo, toGStr(lib), i)
		if C.g_base_info_is_deprecated(baseInfo) == True { // skip deprecated
			continue
		}
		t := C.g_base_info_get_type(baseInfo)
		name := fromGStr(C.g_base_info_get_name(baseInfo)) // name
		switch t {

		// function
		case C.GI_INFO_TYPE_FUNCTION:
			p("=> %s\n", name)
			info := asFunctionInfo(baseInfo)
			generator.GenFunction(info)
		}
	}
}
