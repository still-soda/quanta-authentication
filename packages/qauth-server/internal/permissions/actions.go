package permissions

// CURD 操作枚举
type Action int8

const (
	Create Action = 1
	Read   Action = 2
	Update Action = 3
	Delete Action = 4
)
