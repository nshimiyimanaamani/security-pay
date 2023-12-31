package platypus

import (
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const kparams ctxKey = "params"

const isleaf = "isleaf"

var _ Handler = (*HandlerFunc)(nil)
var _ Handler = (*Mux)(nil)

// Handler ...
type Handler interface {
	Process(context.Context, *Command) (Result, error)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as a Handler.
type HandlerFunc func(context.Context, *Command) (Result, error)

// Process calls fn(ctx, cmd)
func (fn HandlerFunc) Process(ctx context.Context, cmd *Command) (Result, error) {
	return fn(ctx, cmd)
}

// Mux ...
type Mux struct {
	tree     *node
	notFound Handler
}

// New ...
func New(prefix string, notFound Handler) *Mux {
	prefix = strings.TrimSuffix(prefix, "#")
	node := node{key: prefix, isParam: false}
	return &Mux{tree: &node, notFound: notFound}
}

// Process dispatches a command sequense to the handler whose
// pattern most closely matches the cmd pattern.
func (mux *Mux) Process(ctx context.Context, cmd *Command) (Result, error) {

	params := ParamsFromContext(ctx)

	node, key := mux.tree.traverse(strings.Split(cmd.Request, "*")[1:], params)

	if node != nil {
		params.Add(isleaf, node.isLeaf)
	}
	ctx = ContextWithParams(ctx, params)

	//for visualization purposes
	logrus.WithFields(logrus.Fields{
		"request": cmd.Request,
		"params":  params,
		"key":     key,
		"phone":   cmd.Phone,
	}).Infof("dispatching command to handler %v", time.Now().Format("2006-01-02 15:04:05"))

	if node.action != nil {
		return node.action.Process(ctx, cmd)
	}
	return mux.notFound.Process(ctx, cmd)
}

// Handle ...
func (mux *Mux) Handle(pattern string, handler Handler, ts Transformer) {
	if pattern[0] != '*' {
		panic("Path has to start with a *.")
	}
	if handler == nil {
		panic("mux: nil handler")
	}
	mux.tree.insertNode(pattern, handler, ts)

}

// HandlerFunc registers the handler function for the given pattern.
func (mux *Mux) HandlerFunc(pattern string, handler func(context.Context, *Command) (Result, error), ts Transformer) {
	mux.Handle(pattern, HandlerFunc(handler), ts)
}

// NotFound returns an error indicating that the handler was not found for the given task.
func NotFound(ctx context.Context, cmd *Command) (Result, error) {
	return Result{Out: "undefined"}, nil
}

// NotFoundHandler returns a simple task handler that returns a “not found“ error.
func NotFoundHandler() Handler { return HandlerFunc(NotFound) }
