package main

//go:generate sh -c "curl -sL https://raw.githubusercontent.com/LadybugDB/go-ladybug/refs/heads/master/download_lbug.sh | bash -s -- -out lib-ladybug"

/*
#cgo darwin LDFLAGS: -L${SRCDIR}/lib-ladybug -llbug -Wl,-rpath,${SRCDIR}/lib-ladybug
#cgo linux LDFLAGS: -L${SRCDIR}/lib-ladybug -llbug -Wl,-rpath,${SRCDIR}/lib-ladybug
#cgo windows LDFLAGS: -L${SRCDIR}/lib-ladybug -llbug_shared
#include <stdlib.h>
*/
import "C"

import (
	"os"

	"github.com/soundprediction/predicato/cmd/predicato"
)

func main() {
	if err := predicato.Execute(); err != nil {
		os.Exit(1)
	}
}
