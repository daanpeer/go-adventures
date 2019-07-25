package main

func has(subject map[string]interface{}, key string) bool {
	_, ok := subject[key]
	return ok
}
