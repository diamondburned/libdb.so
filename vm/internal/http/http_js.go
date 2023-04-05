package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"syscall/js"
)

type fetchDoer struct {
	fetch js.Value
}

func init() {
	DefaultClient = &Client{Doer: fetchDoer{fetch: js.Global().Get("fetch")}}
}

type promiseResult[T any] struct {
	value T
	err   error
}

func (d fetchDoer) Do(ctx context.Context, req *Request) (*Response, error) {
	ch := make(chan promiseResult[*Response], 1)

	go func() {
		jsHeaders := js.Global().Get("Headers").New()
		for k, v := range req.Header {
			jsHeaders.Call("append", k, v[0])
		}

		promise := d.fetch.Invoke(req.URL, js.ValueOf(map[string]any{
			"method":  req.Method,
			"headers": jsHeaders,
			"body":    req.Body,
		}))

		promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resp := args[0]

			headers := make(Header)
			jsHeaders := resp.Get("headers")
			jsHeaders.Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				header := headers[args[0].String()]
				header = append(header, args[1].String())
				headers[args[0].String()] = header
				return nil
			}))

			var body io.ReadCloser
			jsBody := resp.Get("body")
			if jsBody.Truthy() {
				body = jsReadCloser{
					jsBody:   jsBody,
					jsReader: jsBody.Call("getReader"),
				}
			}

			ch <- promiseResult[*Response]{
				value: &Response{
					Status: resp.Get("status").Int(),
					Header: resp.Get("headers"),
					Body:   resp.Get("body"),
				},
			}
			return nil
		}))
	}()
	return promise, nil
}

type jsReadCloser struct {
	jsReader    js.Value
	buffer      bytes.Buffer
	bufferMu    sync.Mutex
	bufferErr   error
	doneReading chan struct{}
	keepReading chan struct{}
}

func newJSReadCloser(jsReader js.Value) *jsReadCloser {
	r := jsReadCloser{
		jsReader:    jsReader,
		doneReading: make(chan struct{}),
		keepReading: make(chan struct{}),
	}

	go func() {
		for range r.keepReading {
			promise := r.jsReader.Call("read")
			promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				jsResult := args[0]
				if jsResult.IsUndefined() || jsResult.Get("done").Truthy() {
					close(r.done)
					return nil
				}

				jsValue := jsResult.Get("value") // Uint8Array
				if jsValue.IsUndefined() {
					return nil
				}

				r.bufferMu.Lock()
				// Hopefully this is less expensive than allocating a new byte
				// slice per call.
				r.buffer.Grow(jsValue.Length())
				for i := 0; i < jsValue.Length(); i++ {
					r.buffer.WriteByte(byte(jsValue.Index(i).Int()))
				}
				r.bufferMu.Unlock()

				return nil
			}))
			promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				r.bufferMu.Lock()
				r.bufferErr = errors.New(args[0].Get("message").String())
				r.bufferMu.Unlock()
				close(r.done)
				return nil
			}))
		}
	}()
}

func (r *jsReadCloser) Read(p []byte) (int, error) {
	for {
		select {
		case <-r.doneReading:
			return 0, io.EOF
		default:

		}
	}
	if r.buffer.Len() == 0 {
		ch := make(chan error, 1)

		go func() {
			promise := r.jsReader.Call("read")
			promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- promiseResult[js.Value]{value: args[0]}
				return nil
			}))
		}()

	}

	jsResult := r.jsBody.Call("read", js.ValueOf(p))
	if jsResult.IsUndefined() {
		return 0, io.EOF
	}
	return jsResult.Get("value").Int(), nil
}
