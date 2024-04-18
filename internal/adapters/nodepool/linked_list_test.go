package nodepool

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLinkedList(t *testing.T) {
	t.Run("linked list Next()", func(t *testing.T) {
		t.Parallel()
		list := NewLinkedList()

		list.Insert(NewNode("1", "test.test"))
		list.Insert(NewNode("2", "test.test"))
		list.Insert(NewNode("3", "test.test"))

		expected := []string{"1", "2", "3", "1", "2", "3"}

		for i := 0; i < len(expected); i++ {
			el := list.Next()

			require.NotNil(t, el)
			require.Equal(t, expected[i], el.node.Name)
		}
	})

	t.Run("linked list NextAliveNode() with alive node", func(t *testing.T) {
		t.Parallel()
		list := NewLinkedList()

		nodes := []*Node{
			NewNode("1", "test.test"),
			NewNode("2", "test.test"),
			NewNode("3", "test.test"),
		}
		nodes[2].IsAlive = true

		for _, node := range nodes {
			list.Insert(node)
		}

		require.Equal(t, nodes[2], list.NextAliveNode())
		// test that second time works ok too
		require.Equal(t, nodes[2], list.NextAliveNode())
	})

	t.Run("linked list NextAliveNode() without alive node", func(t *testing.T) {
		t.Parallel()
		list := NewLinkedList()

		nodes := []*Node{
			NewNode("1", "test.test"),
			NewNode("2", "test.test"),
			NewNode("3", "test.test"),
		}

		for _, node := range nodes {
			list.Insert(node)
		}

		require.Nil(t, list.NextAliveNode())
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		var expected []*Node

		list := NewLinkedList()

		require.Nil(t, list.Next())

		elements := list.GetElements()

		require.Equal(t, expected, elements)
	})

	t.Run("get elements", func(t *testing.T) {
		t.Parallel()

		list := NewLinkedList()

		expected := []*Node{
			NewNode("1", "test.test"),
			NewNode("2", "test.test"),
			NewNode("3", "test.test"),
		}

		list.Insert(expected[0])
		list.Insert(expected[1])
		list.Insert(expected[2])

		result := list.GetElements()

		require.Equal(t, expected, result)
	})
}
