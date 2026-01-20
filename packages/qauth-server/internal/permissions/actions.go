package permissions

// CURD 操作枚举
type Action int8

const (
	Create Action = iota
	Read
	Update
	Delete
)
