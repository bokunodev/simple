package logger

import (
	"io"
	"log"
	"os"
)

var (
	defaultFlags = log.LstdFlags | log.LUTC | log.Lmsgprefix
	Erro         = log.New(os.Stderr, "[ERRO] ", defaultFlags)
	Info         = log.New(os.Stderr, "[INFO] ", defaultFlags)
	Warn         = log.New(os.Stderr, "[WARN] ", defaultFlags)
)

func ReInit() {
	log.SetFlags(Erro.Flags())
	log.SetOutput(Erro.Writer())
	log.SetPrefix(Erro.Prefix())
}

func init() { ReInit() }

func SetOutput(w io.Writer) {
	Erro.SetOutput(w)
	Info.SetOutput(w)
	Warn.SetOutput(w)
}
