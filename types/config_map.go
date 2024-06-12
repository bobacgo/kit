package types

const _configMapDef = "default"

type ConfigMap[T any] map[string]*T

func (c ConfigMap[T]) Get(key string) *T {
	return c[key]
}

func (c ConfigMap[T]) Default() *T {
	return c[_configMapDef]
}
