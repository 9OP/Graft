package utils

func Min[K uint | uint8 | uint16 | uint32 | uint64 | int](value_0, value_1 K) K {
	if value_0 < value_1 {
		return value_0
	}
	return value_1
}

func Max[K uint | uint8 | uint16 | uint32 | uint64 | int](value_0, value_1 K) K {
	if value_0 > value_1 {
		return value_0
	}
	return value_1
}

func CopyMap[K interface{}, V comparable](originalMap map[V]K) map[V]K {
	newMap := make(map[V]K, len(originalMap))
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}
