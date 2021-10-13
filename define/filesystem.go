package define

import (
	"github.com/spf13/afero"
)

var FS afero.Fs

func init() {
	FS = afero.NewOsFs()
}

func GetSymbolLinker() afero.Linker {
	linker, ok := FS.(afero.Linker)
	if !ok {
		return nil
	}

	return linker
}

func GetSymbolLinkReader() afero.LinkReader {
	reader, ok := FS.(afero.LinkReader)
	if !ok {
		return nil
	}

	return reader
}
