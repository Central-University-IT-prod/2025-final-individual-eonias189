package helpers

import (
	"context"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func ConfigureExpect(t *testing.T, ctx context.Context, baseUrl string) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		TestName: t.Name(),
		BaseURL:  baseUrl,
		Context:  ctx,
		Reporter: httpexpect.NewRequireReporter(t),
	})
}
