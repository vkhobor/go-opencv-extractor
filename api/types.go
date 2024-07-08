package api

import "context"

// TODO remove the structs leave the func type
type WithBody[T any] struct {
	Body T
}

type WithPathId struct {
	ID string `path:"id"`
}

type Empty struct {
}

type Handler[Req, Res any] func(context.Context, *Req) (*Res, error)
