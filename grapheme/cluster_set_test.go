package grapheme

import (
	"testing"

	"github.com/engmtcdrm/go-ansi"
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

// Tests for [ClusterSet.Bytes] function.
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

// TODO: Tests for [ClusterSet.Len] function.
func Test_ClusterSet_Len(t *testing.T) {
	t.Skip("Need to implement")
}

// Tests for [ClusterSet.Runes] function.
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

// Tests for [ClusterSet.String] function.
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

// Tests for [ClusterSet.VisualLen] function.
func Test_ClusterSet_VisualLen(t *testing.T) {
	t.Run("ClusterSet with all visible characters", func(t *testing.T) {
		cs := ClusterSet{
			hello,
			space,
			world,
		}
		expected := len("Hello World")
		output := cs.VisualLen()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with complex grapheme cluster characters", func(t *testing.T) {
		cs := ClusterSet{
			graphemeCluster,
		}
		expected := 4 // The visual length of a complex grapheme cluster is 1
		output := cs.VisualLen()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with ANSI escape sequence characters", func(t *testing.T) {
		cs := ClusterSet{
			New([]rune(ansi.Red)...), // ANSI escape sequence for red text
			hello,
			space,
			world,
			New([]rune(ansi.Reset)...), // ANSI escape sequence to reset
		}
		expected := len("Hello World") // The visual length should ignore ANSI escape sequences
		output := cs.VisualLen()
		require.Equal(t, expected, output)
	})

	t.Run("ClusterSet with complex ANSI escape sequence characters", func(t *testing.T) {
		cs := ClusterSet{
			New('\x1b', '[', '1', ';', '3', '1', 'm'), // ANSI escape sequence for red text
			hello,
			space,
			world,
			New([]rune(ansi.Reset)...), // ANSI escape sequence to reset
		}
		expected := len("Hello World") // The visual length should ignore ANSI escape sequences
		output := cs.VisualLen()
		require.Equal(t, expected, output)
	})
}
