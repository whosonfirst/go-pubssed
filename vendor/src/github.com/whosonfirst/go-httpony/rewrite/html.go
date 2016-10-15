package rewrite

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
)

type HTMLRewriter interface {
	Rewrite(node *html.Node, writer io.Writer) error
	SetKey(key string, value interface{}) error
}

type HTMLRewriteHandler struct {
	writer HTMLRewriter
}

func NewHTMLRewriterHandler(writer HTMLRewriter) (*HTMLRewriteHandler, error) {

	h := HTMLRewriteHandler{
		writer: writer,
	}

	return &h, nil
}

func (h HTMLRewriteHandler) Handler(reader io.Reader) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		doc, err := html.Parse(reader)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		// This bit is still problematic and triggers all kinds of weird
		// memory-pointer errors. I have no idea why (20160627/thisisaaronland)

		h.writer.SetKey("request", req)
		h.writer.Rewrite(doc, rsp)
		return
	}

	return http.HandlerFunc(fn)
}
