package util

// RemoveElement 去除某个元素
func RemoveElement(arr []string, ele string) (newArr []string) {
	for _, v := range arr {
		if v != ele {
			newArr = append(newArr, v)
		}
	}
	return newArr
}

// RemoveDuplicates 数组去重
func RemoveDuplicates[T comparable](arr []T) []T {
	encountered := map[T]bool{}
	var result []T
	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}
