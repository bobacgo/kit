package tag

// TODO 设计一个默认赋值的标签
func Default[T any](v T) T {
	return v
}