package domain

import "context"

type MorzaCache interface {
	TryGetByServiceName(context context.Context, serviceName string) ([]byte, error)
	WarmByServiceName(context context.Context, serviceName string, data []byte) error
	CoolByServiceName(context context.Context, serviceName string) error
	CoolAll(context context.Context) error
}
