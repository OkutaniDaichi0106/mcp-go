package mcp

type Capabilities map[string]map[string]bool

func (c Capabilities) AddCapability(key string, value string) Capabilities {
	if _, ok := c[key]; !ok {
		c[key] = map[string]bool{}
	}
	c[key][value] = true
	return c
}

func (c Capabilities) Merge(other Capabilities) Capabilities {
	for k, v := range other {
		if _, ok := c[k]; !ok {
			c[k] = v
		}
	}
	return c
}

func (c Capabilities) HasFeature(key string) bool {
	if _, ok := c[key]; !ok {
		return false
	}
	return true
}

func (c Capabilities) HasCapability(key string, value string) bool {
	if _, ok := c[key]; !ok {
		return false
	}
	if _, ok := c[key][value]; !ok {
		return false
	}
	return true
}
