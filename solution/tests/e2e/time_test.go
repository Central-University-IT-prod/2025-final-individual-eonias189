package e2e

import (
	"advertising/tests/helpers"
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestAdvanceDay(t *testing.T) {
	ctx := context.Background()
	advertisingServerUrl := "http://localhost:8080"
	// advertisingServerUrl := helpers.SetUpInfrastructure(ctx, t, "../../advertising-service/migrations")

	e := helpers.ConfigureExpect(t, ctx, advertisingServerUrl)

	// check set day
	advanceDaySuccess(e, pointer(42)).
		JSON().
		IsObject().
		Object().
		HasValue("current_date", 42)

	// check increment
	advanceDaySuccess(e, nil).
		JSON().Object().
		ContainsKey("current_date").
		HasValue("current_date", 43)

}

func advanceDaySuccess(e *httpexpect.Expect, day *int) *httpexpect.Response {
	req := e.POST("/time/advance")
	if day != nil {
		req = req.WithJSON(helpers.JSON{"current_date": *day})
	}
	return req.
		Expect().
		Status(http.StatusOK)
}

func pointer[T any](v T) *T {
	return &v
}
