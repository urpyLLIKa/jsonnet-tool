package yaml

import (
	"math"
	"sort"

	yamlv2 "gopkg.in/yaml.v2"
)

func getKeysForMap(yaml map[interface{}]interface{}) []interface{} {
	keys := make([]interface{}, 0, len(yaml))
	for k := range yaml {
		keys = append(keys, k)
	}

	return keys
}

func recursivelyUpdateValue(v interface{}, priorityKeys []string) interface{} {
	switch v2 := v.(type) {
	case map[interface{}]interface{}:
		return recursivelyUpdateMap(v2, priorityKeys)
	case yamlv2.MapSlice:
		return recursivelyUpdateMapSlice(v2, priorityKeys)
	case []interface{}:
		return recursivelyUpdateArray(v2, priorityKeys)
	default:
		return v
	}
}

func recursivelyUpdateArray(yaml []interface{}, priorityKeys []string) []interface{} {
	r := make([]interface{}, 0, len(yaml))
	for _, v := range yaml {
		r = append(r, recursivelyUpdateValue(v, priorityKeys))
	}

	return r
}

func priorityOf(key interface{}, priorityKeys []string) int {
	w, ok := key.(string)
	if !ok {
		return math.MaxInt32
	}

	for i, k := range priorityKeys {
		if k == w {
			return i
		}
	}

	return math.MaxInt32
}

func comparator(i interface{}, j interface{}, priorityKeys []string) bool {
	keyA := priorityOf(i, priorityKeys)
	keyB := priorityOf(j, priorityKeys)

	if keyA == math.MaxInt32 && keyB == math.MaxInt32 {
		sa, oka := i.(string)
		sb, okb := j.(string)

		if oka && okb {
			return sa < sb
		}

		return false
	}

	return keyA < keyB
}

func recursivelyUpdateMapSlice(yaml yamlv2.MapSlice, priorityKeys []string) yamlv2.MapSlice {
	sort.Slice(yaml, func(i, j int) bool {
		return comparator(yaml[i].Key, yaml[j].Key, priorityKeys)
	})

	var r yamlv2.MapSlice

	for i := range yaml {
		k := yaml[i].Key
		v := yaml[i].Value

		w := recursivelyUpdateValue(v, priorityKeys)

		r = append(r, yamlv2.MapItem{Key: k, Value: w})
	}

	return r
}

func recursivelyUpdateMap(yaml map[interface{}]interface{}, priorityKeys []string) yamlv2.MapSlice {
	keys := getKeysForMap(yaml)

	sort.Slice(keys, func(i, j int) bool {
		return comparator(keys[i], keys[j], priorityKeys)
	})

	r := make(yamlv2.MapSlice, 0, len(keys))

	for _, k := range keys {
		v := yaml[k]
		w := recursivelyUpdateValue(v, priorityKeys)

		r = append(r, yamlv2.MapItem{Key: k, Value: w})
	}

	return r
}

// ReorderKeys reorders YAML prioritizing certain keys.
func ReorderKeys(yaml map[interface{}]interface{}, priorityKeys []string) yamlv2.MapSlice {
	return recursivelyUpdateMap(yaml, priorityKeys)
}
