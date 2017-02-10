package mocks

type Cache struct {
	SyncCall struct {
		Received struct {
			Package string
		}
		Returns struct {
			Err error
		}
	}
}

func (c *Cache) Sync(pkg string) error {
	c.SyncCall.Received.Package = pkg
	return c.SyncCall.Returns.Err
}
