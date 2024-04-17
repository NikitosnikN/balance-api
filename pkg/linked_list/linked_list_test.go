package linked_list

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLinkedList(t *testing.T) {
	t.Run("integer linked list", func(t *testing.T) {
		t.Parallel()
		list := NewLinkedList[int]()

		list.Insert(1)
		list.Insert(2)
		list.Insert(3)

		expected := []int{1, 2, 3, 1, 2, 3}

		for i := 0; i < len(expected); i++ {
			node := list.Next()

			require.NotNil(t, node)
			require.Equal(t, expected[i], node.data)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		list := NewLinkedList[int]()

		require.Nil(t, list.Next())

		elements := list.GetElements()

		require.Equal(t, []int{}, elements)
	})

	t.Run("get elements", func(t *testing.T) {
		t.Parallel()

		list := NewLinkedList[int]()

		list.Insert(1)
		list.Insert(2)
		list.Insert(3)
		list.Insert(4)

		result := list.GetElements()
		expected := []int{1, 2, 3, 4}

		require.Equal(t, expected, result)
	})
}
