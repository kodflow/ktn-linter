package slice001

// Cas corrects - accès sécurisés

func GoodCheckedAccess(items []int, idx int) int {
	if idx < len(items) {
		return items[idx]
	}
	return -1
}

func GoodRangeLoop(data []string) {
	for i, v := range data {
		_ = i
		_ = v
	}
}

func GoodLiteralIndex() {
	arr := [5]int{1, 2, 3, 4, 5}
	_ = arr[2]
}

func GoodMultipleChecks(nums []int, i, j int) {
	if i < len(nums) {
		v1 := nums[i]
		_ = v1
	}
	if j < len(nums) {
		v2 := nums[j]
		_ = v2
	}
}
