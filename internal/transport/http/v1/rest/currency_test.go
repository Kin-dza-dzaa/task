package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kin-dza-dzaa/task/internal/transport/http/v1/rest/servicemock"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"
)

var errTest = errors.New("test error")

func Test_CurrencyHandler_Rates(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		want      httpResponse
		setupMock func(currencyServiceMock *servicemock.CurrencyService)
	}{
		{
			name:      "Invalid_date_format",
			setupMock: func(currencyServiceMock *servicemock.CurrencyService) {},
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/rates/invalid_date", http.NoBody)
					chiCtx := chi.NewRouteContext()
					chiCtx.URLParams.Add(datePathKey, "invalid_date")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
				}(),
				w: httptest.NewRecorder(),
			},
			want: httpResponse{
				Message: http.StatusText(http.StatusBadRequest),
				Path:    "/rates/invalid_date",
			},
		},
		{
			name: "Internal_server_error",
			setupMock: func(currencyServiceMock *servicemock.CurrencyService) {
				currencyServiceMock.On("Rates", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, errTest)
			},
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/rates/11.11.2011", http.NoBody)
					chiCtx := chi.NewRouteContext()
					chiCtx.URLParams.Add(datePathKey, "11.11.2011")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
				}(),
				w: httptest.NewRecorder(),
			},
			want: httpResponse{
				Message: http.StatusText(http.StatusInternalServerError),
				Path:    "/rates/11.11.2011",
			},
		},
	}
	for _, tt := range tests {
		currencyServiceMock := servicemock.NewCurrencyService(t)
		h := New(slog.Default(), currencyServiceMock)
		tt.setupMock(currencyServiceMock)

		t.Run(tt.name, func(t *testing.T) {
			h.Rates(tt.args.w, tt.args.r)
			var got httpResponse
			if err := json.NewDecoder(tt.args.w.Body).Decode(&got); err != nil {
				t.Fatalf("json.NewDecoder.Decode: %v", err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("got %v, but wanted %v diff: %v", got, tt.want, diff)
			}
		})
	}
}

func Test_CurrencyHandler_CreateRates(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		want      httpResponse
		setupMock func(currencyServiceMock *servicemock.CurrencyService)
	}{
		{
			name: "Accepted",
			setupMock: func(currencyServiceMock *servicemock.CurrencyService) {
				currencyServiceMock.On("CreateRates", mock.Anything, mock.Anything).Once()
			},
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/rates/11.11.2011", http.NoBody)
					chiCtx := chi.NewRouteContext()
					chiCtx.URLParams.Add(datePathKey, "11.11.2011")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
				}(),
				w: httptest.NewRecorder(),
			},
			want: httpResponse{
				Message: http.StatusText(http.StatusAccepted),
				Path:    "/rates/11.11.2011",
			},
		},
		{
			name:      "Invalid_date_format",
			setupMock: func(currencyServiceMock *servicemock.CurrencyService) {},
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/rates/invalid_date", http.NoBody)
					chiCtx := chi.NewRouteContext()
					chiCtx.URLParams.Add(datePathKey, "invalid_date")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
				}(),
				w: httptest.NewRecorder(),
			},
			want: httpResponse{
				Message: http.StatusText(http.StatusBadRequest),
				Path:    "/rates/invalid_date",
			},
		},
	}
	for _, tt := range tests {
		currencyServiceMock := servicemock.NewCurrencyService(t)
		h := New(slog.Default(), currencyServiceMock)
		tt.setupMock(currencyServiceMock)

		t.Run(tt.name, func(t *testing.T) {
			h.CreateRates(tt.args.w, tt.args.r)
			var got httpResponse
			if err := json.NewDecoder(tt.args.w.Body).Decode(&got); err != nil {
				t.Fatalf("json.NewDecoder.Decode: %v", err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("got %v, but wanted %v diff: %v", got, tt.want, diff)
			}
		})
	}
}
