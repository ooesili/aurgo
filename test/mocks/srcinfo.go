package mocks

import (
	"github.com/ooesili/aurgo/internal/cache"
)

type SrcInfo struct {
	ParseCall struct {
		Recieved struct {
			Input []byte
		}
		Returns struct {
			Package cache.Package
			Err     error
		}
	}
}

func (s *SrcInfo) Parse(input []byte) (cache.Package, error) {
	s.ParseCall.Recieved.Input = input
	returns := s.ParseCall.Returns
	return returns.Package, returns.Err
}
