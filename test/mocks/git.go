package mocks

type Git struct {
	CloneCall struct {
		Received struct {
			URL  string
			Path string
		}
		Returns struct {
			Err error
		}
	}
}

func (g *Git) Clone(url, path string) error {
	g.CloneCall.Received.URL = url
	g.CloneCall.Received.Path = path
	return g.CloneCall.Returns.Err
}
