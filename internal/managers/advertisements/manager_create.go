package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"api-tests-template/internal/client/http/advertisements"
)

// Создает новое объявление через API и возвращает его модель
func CreateAd(t *testing.T, token, title, description string, price float64, photoPath string) Advertisement {
	t.Helper()
	raw := advertisements.HttpPostAdvertisement(t, token, title, description, price, photoPath)
	defer raw.Body.Close()

	require.Equal(t, http.StatusCreated, raw.StatusCode, "Создание должно вернуть 201")

	body, _ := io.ReadAll(raw.Body)
	var ad Advertisement
	require.NoError(t, json.Unmarshal(body, &ad))
	return ad
}

// Получает данные объявления по его ID
func GetAd(t *testing.T, token, id string) Advertisement {
	t.Helper()
	raw := advertisements.HttpGetAdvertisement(t, token, id)
	defer raw.Body.Close()

	require.Equal(t, http.StatusOK, raw.StatusCode, "GET объявления должен вернуть 200")

	body, _ := io.ReadAll(raw.Body)
	var ad Advertisement
	require.NoError(t, json.Unmarshal(body, &ad))
	return ad
}

// Удаляет объявление по его ID
func DeleteAd(t *testing.T, token, id string) {
	t.Helper()
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()
}

// Получает список фотографий объявления
func GetAdPhotos(t *testing.T, token, id string) []Photo {
	t.Helper()
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	require.Equal(t, http.StatusOK, raw.StatusCode, "GET фотографий должен вернуть 200")

	body, _ := io.ReadAll(raw.Body)
	var photos []Photo
	require.NoError(t, json.Unmarshal(body, &photos))
	return photos
}