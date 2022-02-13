package handler

import (
	"testing"
	//"github.com/stretchr/testify/assert"
)

type MockService struct {
}

func TestFileUploadHandler(t *testing.T) {
	//response, err := e(context.Background(), domain.StoreReadings{})
	//expectedResponse := domain.StoreReadings{}

	//assert.NoError(t, err)
	//assert.Equal(t, expectedResponse, response)

	/* w := httptest.NewRecorder()
	r := gin.Default()
	ms := MockService()
	r.POST("/",FileUploadHandler())

	registrationPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	} */
}
