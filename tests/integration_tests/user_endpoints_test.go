package integration_tests

import (
	"encoding/json"
	"github.com/krushev/go-footing/models"
	"github.com/stretchr/testify/assert"
	"net/http"
)

type UserResponse struct {
	Data []models.User `json:"data"`
	Total int          `json:"total"`
}

func (s *EndpointsTestSuite) TestFindAll() {

	//s.prepareTestData()

	req, _ := http.NewRequest(
		"GET",
		"http://localhost:3000/api/v0.0.1/users",
		nil,
	)

	// Perform the request plain with the app.
	// The -1 disables request latency.
	res, err := s.app.Test(req, -1)

	if err != nil {
		s.T().Errorf("Error while calling /api/v0.0.1/users %s", err)
	}

	// verify that no error occurred, that is not expected
	assert.Equal(s.T(), res.StatusCode, 200)

	var userResp *UserResponse
	err = json.NewDecoder(res.Body).Decode(&userResp)

	if err != nil {
		s.T().Errorf("Error while decoding response of /api/v0.0.1/users %s", err)
	}

	assert.Len(s.T(), userResp.Data, 2)
}
