package func009

type MyStruct struct {
	value int
	name  string
}

// Good: Simple getter with no side effects
func (m *MyStruct) GetValue() int {
	return m.value
}

// Good: IsValid with no side effects
func (m *MyStruct) IsValid() bool {
	return m.value > 0
}

// Good: HasName with no side effects
func (m *MyStruct) HasName() bool {
	return m.name != ""
}

// Good: Getter with local variable (no side effect)
func (m *MyStruct) GetDoubleValue() int {
	result := m.value * 2
	return result
}

// Good: Not a getter, can have side effects
func (m *MyStruct) SetValue(v int) {
	m.value = v
}

// Good: Not a getter (Update prefix)
func (m *MyStruct) UpdateValue(v int) {
	m.value = v
}

// Good: Getter with local slice/map assignment (no side effect on struct)
func (m *MyStruct) GetProcessed() []int {
	local := make([]int, 10)
	local[0] = m.value
	return local
}

// Good: Getter with local map assignment (no side effect on struct)
func (m *MyStruct) GetMap() map[string]int {
	result := make(map[string]int)
	result["value"] = m.value
	return result
}

// Good: Test function that is a getter with side effects (exempt)
func TestGetValue(m *MyStruct) int {
	m.value++
	return m.value
}

// Good: Benchmark function that is a getter with side effects (exempt)
func BenchmarkGetValue(m *MyStruct) int {
	m.value++
	return m.value
}

// Good: Getter that assigns to local index (array/map/slice) - no side effect on struct
func (m *MyStruct) GetLocalArray() []int {
	arr := make([]int, 5)
	arr[0] = m.value
	arr[1] = m.value * 2
	return arr
}

// Good: Getter that assigns to local variable index - no side effect
func (m *MyStruct) GetLocalMapValue() int {
	localMap := make(map[string]int)
	localMap["key"] = m.value
	return localMap["key"]
}

// Good: External function declaration (no body) - should be skipped
func GetExternalData() int

// Good: Interface method declaration (no body) - should be skipped
type DataReader interface {
	GetData() string
	IsReady() bool
	HasItems() bool
}

// Good: Getter with local variable increment (no side effect on struct)
func (m *MyStruct) GetIncrementedValue() int {
	local := m.value
	local++ // This is OK - incrementing local variable, not struct field
	return local
}

// Good: Getter with local variable decrement (no side effect on struct)
func (m *MyStruct) GetDecrementedValue() int {
	count := 10
	count-- // This is OK - decrementing local variable
	return count
}
// GetExternalFunc est une fonction externe - ignorée.
func GetExternalFunc() int
