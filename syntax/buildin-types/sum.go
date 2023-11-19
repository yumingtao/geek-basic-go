package main

func Sum(vals []int) int {
	var res int
	for _, v := range vals {
		res += v
	}
	return res
}

func SumInt64(vals []int64) int64 {
	var res int64
	for _, v := range vals {
		res += v
	}
	return res
}

func Keys(m map[string]string) []string {
	keys := make([]string, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func KeysAny(m map[any]any) []any {
	keys := make([]any, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
