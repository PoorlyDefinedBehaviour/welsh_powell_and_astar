package astar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAStar(t *testing.T) {
	t.Parallel()

	cities := map[string]Coordinates{
		"Cascavel":          {Latitude: 24.9578, Longitude: 53.4595},
		"Curitiba":          {Latitude: 25.429, Longitude: 49.2671},
		"Foz do Iguaçu":     {Latitude: 25.5163, Longitude: 54.5854},
		"Francisco Beltrão": {Latitude: 26.0779, Longitude: 53.052},
		"Guarapuava":        {Latitude: 25.3907, Longitude: 51.4628},
		"Londrina":          {Latitude: 23.3045, Longitude: 51.1696},
		"Maringá":           {Latitude: 23.421, Longitude: 51.9331},
		"Paranaguá":         {Latitude: 25.5149, Longitude: 48.5226},
		"Ponta Grossa":      {Latitude: 25.0994, Longitude: 50.1583},
		"São Mateus do Sul": {Latitude: 25.8682, Longitude: 50.3842},
		"Toledo":            {Latitude: 24.7251, Longitude: 53.7417},
		"Umuarama":          {Latitude: 23.7661, Longitude: 53.3206},
	}

	graph := map[string][]string{
		"Cascavel":          {"Toledo", "Foz do Iguaçu", "Francisco Beltrão", "Guarapuava"},
		"Curitiba":          {"São Mateus do Sul", "Ponta Grossa", "Paranaguá"},
		"Foz do Iguaçu":     {"Cascavel"},
		"Francisco Beltrão": {"Cascavel", "São Mateus do Sul"},
		"Guarapuava":        {"Cascavel", "Ponta Grossa"},
		"Londrina":          {"Maringá", "Ponta Grossa"},
		"Maringá":           {"Ponta Grossa", "Umuarama", "Londrina"},
		"Paranaguá":         {"Curitiba"},
		"Ponta Grossa":      {"Curitiba", "Guarapuava", "Maringá", "Londrina"},
		"São Mateus do Sul": {"Francisco Beltrão", "Curitiba"},
		"Toledo":            {"Cascavel", "Umuarama"},
		"Umuarama":          {"Toledo", "Maringá"},
	}

	tests := []struct {
		name     string
		from     string
		to       string
		expected []Vertex
	}{
		{
			name: "Cascavel to Curitiba",
			from: "Cascavel",
			to:   "Curitiba",
			expected: []Vertex{
				{
					Value:     "Cascavel",
					Latitude:  24.9578,
					Longitude: 53.4595,
				},
				{
					Value:     "Guarapuava",
					Latitude:  25.3907,
					Longitude: 51.4628,
				},
				{
					Value:     "Ponta Grossa",
					Latitude:  25.0994,
					Longitude: 50.1583,
				},
				{
					Value:     "Curitiba",
					Latitude:  25.429,
					Longitude: 49.2671,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			graphInput := GraphInput{
				Cities: cities,
				Graph:  graph,
				Path: Path{
					From: tt.from,
					To:   tt.to,
				},
			}

			path, err := AStar(&graphInput)

			assert.Nil(t, err)

			assert.Equal(t, tt.expected, path)
		})
	}
}
