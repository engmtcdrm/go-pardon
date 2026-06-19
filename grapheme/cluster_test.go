package grapheme

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for [New] function.
func Test_New(t *testing.T) {
	t.Run("create a Cluster with a single rune", func(t *testing.T) {
		r := 'a'
		expected := Cluster{r}
		output := New(r)
		require.Equal(t, expected, output)
	})

	t.Run("create a Cluster with multiple runes", func(t *testing.T) {
		expected := Cluster{}
		expected = append(expected, hello...)

		output := New(hello...)
		require.Equal(t, expected, output)
	})

	t.Run("create a Cluster with a complex grapheme cluster", func(t *testing.T) {
		expected := Cluster{}
		expected = append(expected, graphemeCluster...)

		output := New(graphemeCluster...)
		require.Equal(t, expected, output)
	})
}

// Tests for [Cluster.Bytes] function.
func Test_Cluster_Bytes(t *testing.T) {
	t.Run("Cluster with a single rune", func(t *testing.T) {
		c := New('a')
		expected := []byte("a")
		output := c.Bytes()
		require.Equal(t, expected, output)
	})

	t.Run("Cluster with multiple runes", func(t *testing.T) {
		c := New(hello...)
		expected := []byte(helloString)
		output := c.Bytes()
		require.Equal(t, expected, output)
	})

	t.Run("Cluster with a complex grapheme cluster", func(t *testing.T) {
		c := New(graphemeCluster...)
		expected := []byte(graphemeClusterString)
		output := c.Bytes()
		require.Equal(t, expected, output)
	})
}

// Tests for [Cluster.IsANSIEscapeSequence] function.
func Test_Cluster_IsANSIEscapeSequence(t *testing.T) {
	t.Run("Cluster with a single rune", func(t *testing.T) {
		c := New('a')
		expected := false
		output := c.IsANSIEscapeSequence()
		require.Equal(t, expected, output)
	})

	t.Run("Cluster with an ANSI escape sequence", func(t *testing.T) {
		c := New('\x1b', '[', '0', 'm')
		expected := true
		output := c.IsANSIEscapeSequence()
		require.Equal(t, expected, output)
	})
}

// TODO: Tests for [Cluster.Len] function.
func Test_Cluster_Len(t *testing.T) {
	t.Skip("Need to finish")
}

// Tests for [Cluster.String] function.
func Test_Cluster_String(t *testing.T) {
	t.Run("Cluster with a single rune", func(t *testing.T) {
		c := New('a')
		expected := "a"
		output := c.String()
		require.Equal(t, expected, output)
	})

	t.Run("Cluster with multiple runes", func(t *testing.T) {
		c := New(hello...)
		expected := helloString
		output := c.String()
		require.Equal(t, expected, output)
	})

	t.Run("Cluster with a complex grapheme cluster", func(t *testing.T) {
		c := New(graphemeCluster...)
		expected := graphemeClusterString
		output := c.String()
		require.Equal(t, expected, output)
	})
}

// TODO: Tests for [Cluster.VisualLen] function.
func Test_Cluster_VisualLen(t *testing.T) {
	t.Skip("Need to finish")
}
