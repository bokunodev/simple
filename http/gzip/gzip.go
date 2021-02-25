package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"sync"
)

var gzipWriterPool = sync.Pool{New: func() interface{} { return &gzip.Writer{} }}

func Get(w io.Writer) *gzip.Writer {
	g := gzipWriterPool.Get().(*gzip.Writer)
	g.Reset(w)
	return g
}

func Put(g *gzip.Writer) {
	_ = g.Close()
	gzipWriterPool.Put(g)
}

/*
Accept-Encoding: gzip
Accept-Encoding: compress
Accept-Encoding: deflate
Accept-Encoding: br
Accept-Encoding: identity
Accept-Encoding: *
*/
type FileServer struct{ http.FileSystem }

func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
