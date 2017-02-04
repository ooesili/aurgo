package aur_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/ooesili/aurgo/internal/aur"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("New", func() {
	Context("when given an invalid URL", func() {
		It("returns an error", func() {
			_, err := New("://oopsIForgotTheProtocol")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("API", func() {
	var (
		api     API
		server  *httptest.Server
		query   url.Values
		path    string
		fixture string
	)

	BeforeEach(func() {
		query = nil
		path = ""
		fixture = ""

		hfunc := func(resp http.ResponseWriter, req *http.Request) {
			path = req.URL.Path
			query = req.URL.Query()
			fmt.Fprint(resp, fixture)
		}
		server = httptest.NewServer(http.HandlerFunc(hfunc))
	})

	JustBeforeEach(func() {
		var err error
		api, err = New(server.URL)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Version", func() {
		var (
			version string
			err     error
		)

		JustBeforeEach(func() {
			version, err = api.Version("xcape")
		})

		Context("when the package exists", func() {
			BeforeEach(func() {
				fixture = fixtures["info-xcape"]
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns the version from the response", func() {
				Expect(version).To(Equal("1.2-1"))
			})

			It("hits the right endpoint", func() {
				Expect(path).To(Equal("/rpc"))
				Expect(query["v"]).To(ConsistOf("5"))
				Expect(query["type"]).To(ConsistOf("info"))
				Expect(query["arg[]"]).To(ConsistOf("xcape"))
			})
		})

		Context("when the package is not found", func() {
			BeforeEach(func() {
				fixture = fixtures["info-not-found"]
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("package not found: xcape"))
			})
		})

		Context("when the server returns invalid JSON", func() {
			BeforeEach(func() {
				fixture = ""
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when a transport error occurs", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
