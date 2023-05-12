package rest

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Kin-dza-dzaa/task/internal/transport/http/v1/rest/servicemock"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slog"
)

func Test_encodeResopnse(t *testing.T) {
	type args struct {
		w        *httptest.ResponseRecorder
		response httpResponse
		status   int
	}
	tests := []struct {
		wantRes    httpResponse
		wantStatus int
		name       string
		args       args
	}{
		{
			name: "Valid_response",
			args: args{
				w:      httptest.NewRecorder(),
				status: 200,
				response: httpResponse{
					Path:    "/some_path",
					Message: "test",
				},
			},
			wantStatus: 200,
			wantRes: httpResponse{
				Path:    "/some_path",
				Message: "test",
			},
		},
	}
	for _, tt := range tests {
		currencyServiceMock := servicemock.NewCurrencyService(t)
		h := New(slog.Default(), currencyServiceMock)

		t.Run(tt.name, func(t *testing.T) {
			h.encode(tt.args.w, tt.args.status, tt.args.response)
			var gotResponse httpResponse
			err := json.Unmarshal(tt.args.w.Body.Bytes(), &gotResponse)
			if err != nil {
				t.Fatalf("%v - json.Unmarshal: %v", tt.name, err)
			}
			if diff := cmp.Diff(tt.wantRes, gotResponse); diff != "" {
				t.Fatalf("wanted: %v got: %v dif: %v", tt.wantRes, gotResponse, diff)
			}
			if diff := cmp.Diff(tt.wantStatus, tt.args.w.Code); diff != "" {
				t.Fatalf("wanted: %v got: %v dif: %v", tt.wantStatus, tt.args.w.Code, diff)
			}
		})
	}
}
