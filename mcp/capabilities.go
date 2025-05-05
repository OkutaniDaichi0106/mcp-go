package mcp

type Capabilities map[string]map[string]bool

func (c Capabilities) EnableCapability(feature string, method string) Capabilities {
	if _, ok := c[feature]; !ok {
		c[feature] = map[string]bool{}
	}
	c[feature][method] = true
	return c
}

func (c Capabilities) DisableCapability(feature string, method string) Capabilities {
	if _, ok := c[feature]; !ok {
		c[feature] = map[string]bool{}
	}
	c[feature][method] = false
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

func (c Capabilities) HasFeature(feature string) bool {
	if _, ok := c[feature]; !ok {
		return false
	}
	return true
}

func (c Capabilities) HasCapability(feature string, method string) bool {
	if _, ok := c[feature]; !ok {
		return false
	}
	if _, ok := c[feature][method]; !ok {
		return false
	}
	return true
}
