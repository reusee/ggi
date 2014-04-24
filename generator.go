package main

/*
#include <girepository.h>
*/
import "C"
import (
	"bytes"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Generator struct {
	Lib     string
	OutDir  string
	CPrefix string

	FuncsWriter io.Writer

	currentFunctionName string
}

func Gen(lib, outDir string, includes []string) {
	generator := &Generator{
		Lib:    lib,
		OutDir: outDir,
	}
	filePrefix := strings.ToLower(lib) + "_"

	// open girepository
	repo := C.g_irepository_get_default()
	var gerr *C.GError
	C.g_irepository_require(repo, toGStr(lib), nil, 0, &gerr)
	if gerr != nil {
		log.Fatal(fromGStr(gerr.message))
	}
	generator.CPrefix = fromGStr(C.g_irepository_get_c_prefix(repo, toGStr(lib)))

	// functions output file
	funcsWriter := new(bytes.Buffer)
	generator.FuncsWriter = funcsWriter
	w(funcsWriter, "package %s\n", strings.ToLower(lib))
	for _, include := range includes {
		w(funcsWriter, "//#include <%s>\n", include)
	}
	w(funcsWriter, "import \"C\"\n")
	defer func() {
		out, err := os.Create(filepath.Join(outDir, filePrefix+"functions.go"))
		if err != nil {
			log.Fatal(err)
		}
		src, err := format.Source(funcsWriter.Bytes())
		if err != nil {
			p("%s\n", funcsWriter.Bytes())
			out.Write(funcsWriter.Bytes())
			out.Close()
			log.Fatal(err)
		}
		out.Write(src)
		out.Close()
	}()

	// iterate infos
	nInfos := C.g_irepository_get_n_infos(repo, toGStr(lib))
	for i := C.gint(0); i < nInfos; i++ {
		baseInfo := C.g_irepository_get_info(repo, toGStr(lib), i)
		t := C.g_base_info_get_type(baseInfo)
		switch t {

		// function
		case C.GI_INFO_TYPE_FUNCTION:
			info := asFunctionInfo(baseInfo)
			generator.GenFunction(info)
		}
	}
}
