package middleware

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerStatic(prefix string, embedFs embed.FS) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get sub file system
		subFs, err := fs.Sub(embedFs, prefix)
		if err != nil {
			panic(err)
		}

		fsys := http.FS(subFs)
		// open static resource, if it fails, it will be skipped and handed over to other routes for processing.
		f, err := fsys.Open(ctx.Request.URL.Path)
		if err != nil {
			ctx.Next()
			return
		}
		defer f.Close()

		// file server handle request
		http.FileServer(fsys).ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	}
}
