package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"recommendation/internal/model"

	"github.com/stretchr/testify/assert"
)

func Test_userInteractionsHandler_Collector(t *testing.T) {
	type (
		fields struct {
			userInteractionsService UserInteractionService
		}
		args struct {
			w *httptest.ResponseRecorder
			r *http.Request
		}
		response struct {
			status int
			body   []byte
		}
	)
	var (
		//TODO move mock data
		collectorInteractionJSONRequestBody = strings.NewReader(`{
			"user_id":"user123",
			"interactions":[
				{
					"product_sku":"3",
					"action":"add_to_cart",
					"interaction_timestamp":"2024-09-19T22:00:00Z"
				}
			]
			}`)
		collectorInteractionMissingFieldsJSONRequestBody = strings.NewReader(`{
			"user_id":"user123"	}`)
		malformedJSONRequestBody = strings.NewReader(`---`)

		malformedBodyErrorResponse, _    = json.Marshal(model.NewApiError("Malformed Body", 400))
		missingFieldBodyErrorResponse, _ = json.Marshal(model.NewApiError("Missing user interactions", 400))
	)

	// req := httptest.NewRequest("POST", "/collector/interactions", collectorInteractionReequestBody)
	// rr := httptest.NewRecorder()

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected response
	}{
		{
			name: "When user collect interactions sucessfully then return 201",
			fields: fields{
				userInteractionsService: NewUserInteractionServiceMock(nil),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/collector/interactions", collectorInteractionJSONRequestBody),
			},
			expected: response{
				status: http.StatusCreated,
				body:   nil,
			},
		},
		{
			name: "When user send malformed JSON then return 400",
			fields: fields{
				userInteractionsService: NewUserInteractionServiceMock(nil),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/collector/interactions", malformedJSONRequestBody),
			},
			expected: response{
				status: http.StatusBadRequest,
				body:   malformedBodyErrorResponse,
			},
		},
		{
			name: "When user miss required field then return 400",
			fields: fields{
				userInteractionsService: NewUserInteractionServiceMock(nil),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/collector/interactions", collectorInteractionMissingFieldsJSONRequestBody),
			},
			expected: response{
				status: http.StatusBadRequest,
				body:   missingFieldBodyErrorResponse,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserInteractionHandler(tt.fields.userInteractionsService)
			h.Collector(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			resBody, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.expected.status, res.StatusCode)
			if tt.expected.body == nil {
				assert.Empty(t, resBody)
			} else {
				assert.JSONEq(t, string(tt.expected.body), string(resBody))
			}

		})
	}
}

type userInteractionServiceMock struct {
	responseError error
}

func NewUserInteractionServiceMock(responseError error) *userInteractionServiceMock {
	return &userInteractionServiceMock{
		responseError: responseError,
	}
}

func (s *userInteractionServiceMock) Collect(userInteraction model.UserInteraction) error {
	return s.responseError
}
