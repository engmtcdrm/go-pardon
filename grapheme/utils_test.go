package grapheme

import (
	"testing"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/stretchr/testify/require"
)

// Tests for [ClusterSetFromString] function.
// func Test_ClusterSetFromString(t *testing.T) {
// 	t.Run("empty string", func(t *testing.T) {
// 		var expected ClusterSet
// 		result := ClusterSetFromString("")
// 		require.Equal(t, expected, result)
// 	})

// 	t.Run("simple string", func(t *testing.T) {
// 		expected := "Hello"
// 		result := ClusterSetFromString("Hello")
// 		require.Equal(t, expected, result.String())
// 	})
// }

// Tests for [Equal] function.
func Test_Equal(t *testing.T) {
	t.Run("Empty clusters", func(t *testing.T) {
		var a, b Cluster
		require.True(t, Equal(a, b), "Expected empty clusters to be equal")
	})

	t.Run("Identical clusters", func(t *testing.T) {
		a := NewFromString("Hello")
		b := NewFromString("Hello")
		require.True(t, Equal(a, b), "Expected identical clusters to be equal")
	})

	t.Run("Different clusters", func(t *testing.T) {
		a := NewFromString("Hello")
		b := NewFromString("World")
		require.False(t, Equal(a, b), "Expected different clusters to not be equal")
	})

	t.Run("Clusters with different lengths", func(t *testing.T) {
		a := NewFromString("Hello")
		b := NewFromString("Hello!")
		require.False(t, Equal(a, b), "Expected clusters of different lengths to not be equal")
	})
}

// Tests for [EqualFold] function.
func Test_EqualFold(t *testing.T) {
	t.Run("Empty clusters", func(t *testing.T) {
		var a, b Cluster
		require.True(t, EqualFold(a, b), "Expected empty clusters to be equal")
	})

	t.Run("Identical clusters", func(t *testing.T) {
		a := NewFromString("Hello")
		b := NewFromString("Hello")
		require.True(t, EqualFold(a, b), "Expected identical clusters to be equal")
	})

	t.Run("Case-insensitive match", func(t *testing.T) {
		a := NewFromString("Hello")
		b := NewFromString("hello")
		require.True(t, EqualFold(a, b), "Expected clusters to be equal ignoring case")
	})

	t.Run("Different clusters", func(t *testing.T) {
		a := NewFromString("Hello")
		b := NewFromString("World")
		require.False(t, EqualFold(a, b), "Expected different clusters to not be equal")
	})

	t.Run("Clusters with different lengths", func(t *testing.T) {
		a := NewFromString("Hello")
		b := NewFromString("Hello!")
		require.False(t, EqualFold(a, b), "Expected clusters of different lengths to not be equal")
	})
}

// Tests for [HasSequenceEnd] function.
func Test_HasSequenceEnd(t *testing.T) {
	t.Run("Empty cluster", func(t *testing.T) {
		var cluster Cluster
		require.False(t, HasSequenceEnd(cluster), "Expected empty cluster to not have sequence end")
	})

	t.Run("Cluster without sequence end", func(t *testing.T) {
		cluster := NewFromString("-----")
		require.False(t, HasSequenceEnd(cluster), "Expected cluster without sequence end to return false")
	})

	t.Run("Cluster with sequence end", func(t *testing.T) {
		cluster := NewFromString("\x1b[31mHello\x1b[0m")
		require.True(t, HasSequenceEnd(cluster), "Expected cluster with sequence end to return true")
	})
}

// Tests for [IndexOfSequenceEnd] function.
func Test_IndexOfSequenceEnd(t *testing.T) {
	t.Run("Empty cluster", func(t *testing.T) {
		var cluster Cluster
		result := IndexOfSequenceEnd(cluster)
		require.Equal(t, -1, result, "Expected empty cluster to return -1")
	})

	t.Run("Cluster without sequence end", func(t *testing.T) {
		cluster := NewFromString("-----")
		result := IndexOfSequenceEnd(cluster)
		require.Equal(t, -1, result, "Expected cluster without sequence end to return -1")
	})

	t.Run("Cluster with sequence end", func(t *testing.T) {
		cluster := NewFromString("\x1b[31mHello\x1b[0m")
		result := IndexOfSequenceEnd(cluster)
		require.NotEqual(t, -1, result, "Expected cluster with sequence end to return index of sequence end")
	})
}

// Tests for [IsFeEscapeSequence] function.
func Test_IsFeEscapeSequence(t *testing.T) {
	t.Run("Empty cluster", func(t *testing.T) {
		var cluster Cluster
		require.False(t, IsFeEscapeSequence(cluster), "Expected empty cluster to not be a FE escape sequence")
	})

	t.Run("Cluster that is a FE escape sequence", func(t *testing.T) {
		cluster := NewFromString("\x1b[31m")
		require.True(t, IsFeEscapeSequence(cluster), "Expected cluster that is a FE escape sequence to return true")
	})

	t.Run("Cluster has incorrect second character", func(t *testing.T) {
		cluster := NewFromString("\x1b-")
		require.False(t, IsFeEscapeSequence(cluster), "Expected cluster that is not a FE escape sequence to return false")
	})

	t.Run("Cluster that is not a FE escape sequence", func(t *testing.T) {
		cluster := NewFromString("Hello")
		require.False(t, IsFeEscapeSequence(cluster), "Expected cluster that is not a FE escape sequence to return false")
	})
}

// Tests for [IsSequenceEnd] function.
func Test_IsSequenceEnd(t *testing.T) {
	t.Run("Empty cluster", func(t *testing.T) {
		var cluster Cluster
		require.False(t, IsSequenceEnd(cluster), "Expected empty cluster to not be a sequence end")
	})

	t.Run("Cluster that is too long", func(t *testing.T) {
		cluster := NewFromString("\x1b[0m")
		require.False(t, IsSequenceEnd(cluster), "Expected cluster without sequence end to return false")
	})

	t.Run("Cluster that is a sequence end", func(t *testing.T) {
		cluster := NewFromString("m")
		require.True(t, IsSequenceEnd(cluster), "Expected cluster that is a sequence end to return true")
	})

	t.Run("Cluster that is not a sequence end", func(t *testing.T) {
		cluster := NewFromString("-")
		require.False(t, IsSequenceEnd(cluster), "Expected cluster that is not a sequence end to return false")
	})
}

// TODO: Tests for [clusterANSIEscape] function.
func Test_clusterANSIEscape(t *testing.T) {
	t.Run("Empty ClusterSet", func(t *testing.T) {
		var clusterSet ClusterSet
		result := clusterANSIEscape(clusterSet)
		require.Equal(t, clusterSet, result)
	})

	t.Run("One character", func(t *testing.T) {
		clusterSet := ClusterSet{New('A')}
		result := clusterANSIEscape(clusterSet)
		require.Equal(t, clusterSet, result)
	})

	t.Run("Has escape, but not ANSI escape sequence", func(t *testing.T) {
		clusterSet := testClusterSetFromString(t, "This \x1b is an escape a test")
		result := clusterANSIEscape(clusterSet)
		require.Equal(t, clusterSet, result)
	})

	t.Run("Only an escape character", func(t *testing.T) {
		clusterSet := ClusterSetFromString("\x1b")
		result := clusterANSIEscape(clusterSet)
		require.Equal(t, clusterSet, result)
	})

	t.Run("Has CSI escape sequence, but no sequence end", func(t *testing.T) {
		clusterSet := ClusterSetFromString("This \x1b[ is red but no sequence end")
		result := clusterANSIEscape(clusterSet)
		require.NotEqual(t, clusterSet, result)
	})

	t.Run("Has ANSI escape sequence", func(t *testing.T) {
		clusterSet := ClusterSetFromString("This \x1b[31m is red")
		result := clusterANSIEscape(clusterSet)
		require.NotEqual(t, clusterSet, result)
	})
}

// testClusterSetFromString is a helper function to create a ClusterSet from a
// string, splitting it into grapheme clusters.
func testClusterSetFromString(t *testing.T, s string) ClusterSet {
	t.Helper()
	var pendingSet ClusterSet
	gr := graphemes.FromString(s)

	for gr.Next() {
		cluster := New([]rune(gr.Value())...)
		pendingSet = append(pendingSet, cluster)
	}
	return pendingSet
}
