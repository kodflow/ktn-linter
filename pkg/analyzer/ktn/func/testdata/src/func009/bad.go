package func009

type Counter struct {
	count int
	cache int
}

// Bad: Getter that modifies state
func (c *Counter) GetCount() int {
	c.count++ // want "KTN-FUNC-009"
	return c.count
}

// Bad: Getter that assigns to field
func (c *Counter) GetCachedValue() int {
	c.cache = c.count * 2 // want "KTN-FUNC-009"
	return c.cache
}

// Bad: IsValid with side effect
func (c *Counter) IsReady() bool {
	c.cache = 100 // want "KTN-FUNC-009"
	return c.count > 0
}

// Bad: HasData with increment
func (c *Counter) HasData() bool {
	c.count++ // want "KTN-FUNC-009"
	return c.count > 0
}
