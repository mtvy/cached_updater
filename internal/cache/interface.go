package cache

import "context"

type responser interface {
	Get(ctx context.Context) (interface{}, error)
}
