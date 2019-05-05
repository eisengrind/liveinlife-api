package endpoint_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/51st-state/api/pkg/encode"
	"go.uber.org/zap"

	"github.com/51st-state/api/pkg/api/endpoint"
)

func TestNew(t *testing.T) {
	ep := endpoint.New(encode.NewJSONEncoder(), func(_ context.Context, _ *http.Request) (interface{}, error) { return nil, nil })
	ep.HandlerFunc(nil)
}

func TestEndpointWithBefore(t *testing.T) {
	mwFunc := func(ctx context.Context, _ *http.Request) (context.Context, error) {
		t.Log("middleware called")
		return ctx, nil
	}

	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	ep := endpoint.New(encode.NewJSONEncoder(), func(_ context.Context, _ *http.Request) (interface{}, error) { return nil, nil })
	ep.WithBefore(mwFunc)
	ep.WithBefore(mwFunc)

	hndl := ep.HandlerFunc(logger)
	hndl(httptest.NewRecorder(), httptest.NewRequest("get", "/test/path", bytes.NewBufferString("test")))
}
