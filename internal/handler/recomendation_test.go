package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"recommendation/internal/model"

	"github.com/stretchr/testify/assert"
)

func Test_productRecommendationHandler_Get(t *testing.T) {
	type (
		fields struct {
			recommendationService ProductRecommendationService
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
		successResponseBody = []byte(`{
			"products": [
				"3",
				"2",
				"1"
			]
		}`)

		serviceErrorResponseBody, _       = json.Marshal(model.NewApiError("any error", 500))
		invalidUserIdErrorResponseBody, _ = json.Marshal(model.NewApiError("Invalid User ID", 400))
		req                               = httptest.NewRequest("GET", "/recommendations/users/user123", nil)
	)
	req.SetPathValue("user_id", "user123")

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected response
	}{
		{
			name: "When the backend return an error then return 500",
			fields: fields{
				recommendationService: NewProductRecommendationService(nil, errors.New("any error")),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: req,
			},
			expected: response{
				status: http.StatusInternalServerError,
				body:   serviceErrorResponseBody,
			},
		},
		{
			name: "When user ID is not valid then return 400",
			fields: fields{
				recommendationService: NewProductRecommendationService(nil, nil),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/recommendations/users/", nil),
			},
			expected: response{
				status: http.StatusBadRequest,
				body:   invalidUserIdErrorResponseBody,
			},
		},
		{
			name: "When user ask for recommendations then return a product list",
			fields: fields{
				recommendationService: NewProductRecommendationService(&model.ProductsRecommendation{Products: []model.ProductSKU{"3", "2", "1"}}, nil),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: req,
			},
			expected: response{
				status: http.StatusOK,
				body:   successResponseBody,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewProductRecommendationHandler(tt.fields.recommendationService)
			h.Get(tt.args.w, tt.args.r)
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

type productRecommendationServiceMock struct {
	productsRecommendation *model.ProductsRecommendation
	err                    error
}

func NewProductRecommendationService(productsRecommendation *model.ProductsRecommendation, err error) *productRecommendationServiceMock {
	return &productRecommendationServiceMock{
		productsRecommendation: productsRecommendation,
		err:                    err,
	}
}

func (s *productRecommendationServiceMock) GetProductRecommendation(u model.UserID) (*model.ProductsRecommendation, error) {

	return s.productsRecommendation, s.err
}
