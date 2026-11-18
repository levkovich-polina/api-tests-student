package advertisements

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// Формирует тело для создания объявления
func buildMultipartBody(title, description string, price float64, photoPath string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("title", title)
	writer.WriteField("description", description)
	writer.WriteField("price", fmt.Sprintf("%.0f", price))

	absPath, _ := filepath.Abs(photoPath)
	file, _ := os.Open(absPath)
	defer file.Close()

	part, _ := writer.CreateFormFile("photos", filepath.Base(absPath))
	io.Copy(part, file)

	writer.Close()
	return body, writer.FormDataContentType()
}

// Отправляет POST-запрос на создание объявления
func HttpPostAdvertisement(t *testing.T, token, title, description string, price float64, photoPath string) *http.Response {
	body, contentType := buildMultipartBody(title, description, price, photoPath)

	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		Headers(map[string]string{
			"Content-Type": contentType,
		}).
		Body(body.String()).
		Expect(t).End().Response
}

// Получает объявление по ID 
func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).End().Response
}

// Удаляет объявление по ID 
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).End().Response
}

// Получает фотографии объявления 
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get("/advertisements/" + id + "/photos").
		Expect(t).End().Response
}