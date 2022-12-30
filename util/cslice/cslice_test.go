package cslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSlice(t *testing.T) {
	assert := assert.New(t)

	// 测试 Add 方法
	t.Run("Test Add", func(t *testing.T) {
		c := New[int](10)
		c.Add(1)
		assert.Equal(1, c.Length(), "添加元素后长度应为 1")
		c.Add(2)
		assert.Equal(2, c.Length(), "再添加元素后长度应为 2")
	})

	// 测试 RemoveLast 方法
	t.Run("Test RemoveLast", func(t *testing.T) {
		c := New[int](10)
		assert.Nil(c.RemoveLast(), "空切片删除元素应返回 nil")

		c.Add(1)
		r := c.RemoveLast()
		assert.Equal(1, *r, "删除的元素应为 1")
		assert.Nil(c.RemoveLast(), "再次删除元素应返回 nil")

		c.Add(1)
		c.Add(2)
		r = c.RemoveLast()
		assert.Equal(2, *r, "删除的元素应为 2")
	})

	// 测试 Length 方法
	t.Run("Test Length", func(t *testing.T) {
		c := New[int](10)

		assert.Equal(0, c.Length(), "新建的切片长度应为 0")

		c.Add(1)
		assert.Equal(1, c.Length(), "添加元素后长度应为 1")
	})

	// 测试 Reset 方法
	t.Run("Test Reset", func(t *testing.T) {
		c := New[int](10)

		c.Reset()
		assert.Equal(0, c.Length(), "空切片应保持为空")

		c.Add(1)
		c.Reset()
		assert.Equal(0, c.Length(), "重置后切片应为空")
	})

	// 测试 FillTail 方法
	t.Run("Test FillTail", func(t *testing.T) {
		c := New[int](10)
		items := []int{1, 2, 3}

		c.FillTail(items)
		assert.Equal(3, c.Length(), "元素添加后长度应为 3")

		c.FillTail(items)
		assert.Equal(6, c.Length(), "再次添加元素后长度应为 6")
	})
}
