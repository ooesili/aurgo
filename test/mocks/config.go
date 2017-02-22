package mocks

type Config struct {
	PackagesCall struct {
		Returns struct {
			Packages []string
			Err      error
		}
	}
	AurRepoURLCall struct {
		Received struct {
			Package string
		}
		Returns struct {
			URL string
		}
	}
	SourcePathCall struct {
		Received struct {
			Package string
		}
		Returns struct {
			Path string
			Err  error
		}
	}
	SourceBaseCall struct {
		SourceBase string
	}
}

func (c *Config) Packages() ([]string, error) {
	returns := c.PackagesCall.Returns
	return returns.Packages, returns.Err
}

func (c *Config) AurRepoURL(pkg string) string {
	c.AurRepoURLCall.Received.Package = pkg
	return c.AurRepoURLCall.Returns.URL
}

func (c *Config) SourcePath(pkg string) (string, error) {
	c.SourcePathCall.Received.Package = pkg
	returns := c.SourcePathCall.Returns
	return returns.Path, returns.Err
}

func (c *Config) SourceBase() string {
	return c.SourceBaseCall.SourceBase
}
