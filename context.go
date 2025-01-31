package yaml

import "context"

type ctxMergeKey struct{}

func withMerge(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxMergeKey{}, true)
}

func isMerge(ctx context.Context) bool {
	v, ok := ctx.Value(ctxMergeKey{}).(bool)
	if !ok {
		return false
	}
	return v
}
