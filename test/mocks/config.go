package mocks

type Config struct {
	PackagesCall struct {
		Returns struct {
			Packages []string
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
		}
	}
	SourceBaseCall struct {
		SourceBase string
	}
}

func (c *Config) Packages() []string {
	return c.PackagesCall.Returns.Packages
}

func (c *Config) AurRepoURL(pkg string) string {
	c.AurRepoURLCall.Received.Package = pkg
	return c.AurRepoURLCall.Returns.URL
}

func (c *Config) SourcePath(pkg string) string {
	c.SourcePathCall.Received.Package = pkg
	return c.SourcePathCall.Returns.Path
}

func (c *Config) SourceBase() string {
	return c.SourceBaseCall.SourceBase
}
