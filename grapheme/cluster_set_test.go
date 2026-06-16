package grapheme

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	helloString           = "Hello"
	graphemeClusterString = "👩🏽\u200d💻"
)

var (
	hello           = New('H', 'e', 'l', 'l', 'o')
	space           = New(' ')
	world           = New('W', 'o', 'r', 'l', 'd')
	graphemeCluster = New('👩', '🏽', '‍', '💻')
)

// Tests for [ClusterSet.Bytes] functions.
func Test_ClusterSet_Bytes(t *testing.T) {
	t.Run("empty ClusterSet", func(t *testing.T) {
		cs := ClusterSet{}
		var expected []byte
		output := cs.Bytes()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with multiple Clusters", func(t *testing.T) {
		cs := ClusterSet{
			hello,
			space,
			world,
		}
		expected := []byte("Hello World")
		output := cs.Bytes()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with complex grapheme cluster", func(t *testing.T) {
		cs := ClusterSet{
			graphemeCluster,
		}
		expected := []byte(graphemeClusterString)
		output := cs.Bytes()
		require.Equal(t, expected, output)
	})
}

// Tests for [ClusterSet.Runes] functions.
func Test_ClusterSet_Runes(t *testing.T) {
	t.Run("empty ClusterSet", func(t *testing.T) {
		cs := ClusterSet{}
		var expected []rune
		output := cs.Runes()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with multiple Clusters", func(t *testing.T) {
		cs := ClusterSet{
			hello,
			space,
			world,
		}
		expected := []rune("Hello World")
		output := cs.Runes()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with complex grapheme cluster", func(t *testing.T) {
		cs := ClusterSet{
			graphemeCluster,
		}
		expected := []rune(graphemeClusterString)
		output := cs.Runes()
		require.Equal(t, expected, output)
	})
}

// Tests for [ClusterSet.String] functions.
func Test_ClusterSet_String(t *testing.T) {
	t.Run("empty ClusterSet", func(t *testing.T) {
		cs := ClusterSet{}
		expected := ""
		output := cs.String()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with multiple Clusters", func(t *testing.T) {
		cs := ClusterSet{
			hello,
			space,
			world,
		}
		expected := "Hello World"
		output := cs.String()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with complex grapheme cluster", func(t *testing.T) {
		cs := ClusterSet{
			graphemeCluster,
		}
		expected := graphemeClusterString
		output := cs.String()
		require.Equal(t, expected, output)
	})
}
