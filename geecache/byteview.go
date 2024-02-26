package geecache

// 存储的数据
type ByteView struct {
	b []byte
}

// 查看占用内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// 因为原有存储为可读的，所以返回一个拷贝以免被修改缓存
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 将缓存数据转换为字符串形式
func (v ByteView) String() string {
	return string(v.b)
}

// 拷贝原有缓存
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
