package welshpowell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func everyVertexHasSameColor(coloredVertices ColoredVertices, vertices []string) bool {
	if len(vertices) == 0 {
		return true
	}

	color := coloredVertices[vertices[0]]

	for i := 1; i < len(vertices); i++ {
		if coloredVertices[vertices[i]] != color {
			return false
		}
	}

	return true
}

func TestWelshPowell(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		graph                   map[string][]string
		expectedChromaticNumber int
		expectedGroups          [][]string
	}{
		{
			name:                    "an empty graph has no colors",
			graph:                   map[string][]string{},
			expectedChromaticNumber: 0,
			expectedGroups:          [][]string{},
		},
		{
			name: "when graph has 4 connected nodes, should use 2 colors",
			graph: map[string][]string{
				/*
					a -- b
					|    |
					c -- d
				*/
				"a": {"b", "c"},
				"b": {"a", "d"},
				"c": {"a", "d"},
				"d": {"b", "c"},
			},
			expectedChromaticNumber: 2,
			expectedGroups:          [][]string{{"a", "d"}, {"b", "c"}},
		},
		{
			name: "should use 3 colors",
			graph: map[string][]string{
				/*
					a -- b ---
					|    |    |
					c -- d -- e
				*/
				"a": {"b", "c"},
				"b": {"a", "d", "e"},
				"c": {"a", "d"},
				"d": {"b", "c", "e"},
				"e": {"b", "d"},
			},
			expectedChromaticNumber: 3,
			expectedGroups:          [][]string{{"a", "d"}, {"b", "c"}, {"e"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := WelshPowell(Graph(tt.graph))

			assert.Equal(t, tt.expectedChromaticNumber, actual.ChromaticNumber)

			for _, verticesThatShouldHaveSameColor := range tt.expectedGroups {
				assert.True(t, everyVertexHasSameColor(actual.Colors, verticesThatShouldHaveSameColor))
			}
		})
	}
}
