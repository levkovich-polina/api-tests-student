package favorites

import (
	"fmt"
	"net/http"
	"testing"

	apiRunner "api-tests-template/internal/helpers/api-runner"
)

const favoritesPath = "/favorites"

// Тело запроса POST /favorites
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

// Возвращает заголовки с токеном авторизации
func buildAuthHeaders(token string) map[string]string {
	if token == "" {
		return nil
	}
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// Добавляет объявление в избранное
func HttpPostFavorite(t *testing.T, token string, adID string) *http.Response {
	req := AddFavoriteRequest{AdvertisementID: adID}
	return apiRunner.GetRunner().Create().
		Post(favoritesPath).
		Headers(buildAuthHeaders(token)).
		JSON(req).
		Expect(t).End().Response
}

// Получает список избранного с пагинацией
func HttpGetFavorites(t *testing.T, token string, page, pageSize int) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(favoritesPath).
		Headers(buildAuthHeaders(token)).
		Query("page", fmt.Sprintf("%d", page)).
		Query("page_size", fmt.Sprintf("%d", pageSize)).
		Expect(t).End().Response
}

// Удаляет объявление из избранного
func HttpDeleteFavorite(t *testing.T, token string, adID string) *http.Response {
	return apiRunner.GetRunner().Create().
		Delete(favoritesPath + "/" + adID).
		Headers(buildAuthHeaders(token)).
		Expect(t).End().Response
}
