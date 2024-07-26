package util

import "context"

type Handler[Req, Res any] func(context.Context, *Req) (*Res, error)
