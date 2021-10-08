package define

import (
	"github.com/spf13/afero"
)

var FS afero.Fs

func init() {
	FS = afero.NewOsFs()
}
