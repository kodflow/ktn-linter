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
