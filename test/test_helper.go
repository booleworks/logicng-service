package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
	"github.com/booleworks/logicng-service/srv"
)

const (
	host    = "localhost"
	port    = "9999"
	timeout = 2 * time.Second
)

func runServer(t *testing.T) context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	cfg := &config.Config{
		Host:                  host,
		Port:                  port,
		SyncComputationTimout: timeout,
	}
	go srv.Run(ctx, cfg)
	return ctx
}

func endpoint(path string) string {
	return fmt.Sprintf(`http://%s:%s/%s`, host, port, path)
}

func callServiceJSON(
	ctx context.Context,
	method string,
	endpoint string,
	body string,
) (*http.Response, error) {
	return callService(ctx, method, endpoint, []byte(body), "application/json")
}

func callServiceProtoBuf(
	ctx context.Context,
	method string,
	endpoint string,
	body []byte,
) (*http.Response, error) {
	return callService(ctx, method, endpoint, body, "application/protobuf")
}

// Taken from the great article at:
// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
func callService(
	ctx context.Context,
	method string,
	endpoint string,
	body []byte,
	content string,
) (*http.Response, error) {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, _ := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(body))
		req.Header.Set("accept", content)
		req.Header.Set("Content-Type", content)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			timeout := 5 * time.Second
			if time.Since(startTime) >= timeout {
				return nil, fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func jsonFormulaInput(formula string) string {
	return fmt.Sprintf(`{"formula": "%s"}`, formula)
}

func pbFormulaInput(formula string) []byte {
	input := sio.FormulaInput{Formula: formula}
	pb, _ := input.ProtoBuf()
	return pb
}

func validateProtoBufFormulaResult(t *testing.T, response *http.Response, formula string) {
	validateSuccess(t, response, "application/protobuf")
	data, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	object, err := sio.FormulaResult{}.DeserProtoBuf(data)
	assert.Nil(t, err)
	assert.True(t, object.State.Success)
	assert.Empty(t, object.State.Error)
	assert.Equal(t, formula, object.Formula)
}

func validateJSONFormulaResult(t *testing.T, response *http.Response, formula string) {
	validateSuccess(t, response, "application/json")
	var object sio.FormulaResult
	err := json.NewDecoder(response.Body).Decode(&object)
	assert.Nil(t, err)
	assert.True(t, object.State.Success)
	assert.Empty(t, object.State.Error)
	assert.Equal(t, formula, object.Formula)
}

func validateJSONBoolResult(t *testing.T, response *http.Response, value bool) {
	validateSuccess(t, response, "application/json")
	var object sio.BoolResult
	err := json.NewDecoder(response.Body).Decode(&object)
	assert.Nil(t, err)
	assert.True(t, object.State.Success)
	assert.Empty(t, object.State.Error)
	assert.Equal(t, value, object.Value)
}

func validateSuccess(t *testing.T, response *http.Response, content string) {
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, content, response.Header.Get("Content-Type"))
	assert.NotNil(t, response.Header.Get("X-Correlation-ID"))
}

func extractJSONBody(response *http.Response) string {
	var bytes []byte
	bytes, _ = io.ReadAll(response.Body)
	defer response.Body.Close()
	return string(bytes)
}
