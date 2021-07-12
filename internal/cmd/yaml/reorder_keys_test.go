package yaml

import (
	"reflect"
	"testing"

	yamlv2 "gopkg.in/yaml.v2"
)

func TestReorderKeys(t *testing.T) {
	tests := []struct {
		name         string
		yaml         map[interface{}]interface{}
		priorityKeys []string
		want         yamlv2.MapSlice
	}{
		{
			name: "one",
			yaml: map[interface{}]interface{}{"hello": "there"},
			want: yamlv2.MapSlice{yamlv2.MapItem{Key: "hello", Value: "there"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReorderKeys(tt.yaml, tt.priorityKeys)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReorderKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
