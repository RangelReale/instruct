package instruct

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapTags(t *testing.T) {
	m := MapTags{
		"A": 1,
		"B": "2",
		"C": 3.3,
		"D": MapTags{
			"X": 66,
			"Y": "77",
			"H": MapTags{
				"U": 90,
				"H": 77,
			},
		},
		"ER": map[string]bool{
			"F": false,
		},
		"MP": map[string]any{
			"U": 321,
		},
	}

	tests := []struct {
		name  string
		path  []string
		want  interface{}
		found bool
	}{
		{
			name:  "item from root",
			path:  []string{"A"},
			want:  1,
			found: true,
		},
		{
			name:  "item from root not found",
			path:  []string{"Z"},
			found: false,
		},
		{
			name:  "item from second level",
			path:  []string{"D", "Y"},
			want:  "77",
			found: true,
		},
		{
			name:  "item from second level not found",
			path:  []string{"D", "J"},
			found: false,
		},
		{
			name:  "item from third level",
			path:  []string{"D", "H", "H"},
			want:  77,
			found: true,
		},
		{
			name:  "item from third level not found",
			path:  []string{"D", "H", "J"},
			found: false,
		},
		{
			name: "map item",
			path: []string{"D"},
			want: MapTags{
				"X": 66,
				"Y": "77",
				"H": MapTags{
					"U": 90,
					"H": 77,
				},
			},
			found: true,
		},
		{
			name:  "nil path",
			path:  nil,
			found: false,
		},
		{
			name:  "empty path",
			path:  []string{},
			found: false,
		},
		{
			name: "other map type",
			path: []string{"ER"},
			want: map[string]bool{
				"F": false,
			},
			found: true,
		},
		{
			name:  "other map type second level must not find",
			path:  []string{"ER", "F"},
			found: false,
		},
		{
			name:  "raw map type",
			path:  []string{"MP", "U"},
			want:  321,
			found: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			v, found := findPath(m, tt.path)
			require.Equal(t, tt.found, found)
			if found && tt.found {
				require.Equal(t, tt.want, v)
			}
		})
	}

}

func TestMapTagsUnusedFields(t *testing.T) {
	m := MapTags{
		"A": 1,
		"B": 2,
		"C": MapTags{
			"X": 99,
			"Y": MapTags{
				"Z": 12,
			},
		},
	}

	err := m.checkUnusedFields(map[string]bool{
		"A":     true,
		"B":     true,
		"C":     true,
		"C.X":   true,
		"C.Y":   true,
		"C.Y.Z": true,
	})
	require.NoError(t, err)

	// must have all
	err = m.checkUnusedFields(map[string]bool{
		"A":     true,
		"C.X":   true,
		"C.Y":   true,
		"C.Y.Z": true,
	})
	require.Error(t, err)

	// can have more
	err = m.checkUnusedFields(map[string]bool{
		"A":      true,
		"B":      true,
		"C":      true,
		"C.X":    true,
		"C.Y":    true,
		"C.Y.Z":  true,
		"C.Y.HH": true,
		"ZZ":     true,
	})
	require.NoError(t, err)

}
