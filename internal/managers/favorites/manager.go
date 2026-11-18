package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	clientFav "api-tests-template/internal/client/http/favorites"
)

// Добавляет объявление в избранное и возвращает распарсенную запись
func AddFavorite(t *testing.T, token string, adID string) FavoriteResponse {
	t.Helper()
	t.Logf("Добавляем объявление %s в избранное", adID)

	raw := clientFav.HttpPostFavorite(t, token, adID)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ добавления в избранное")
	require.Equal(t, http.StatusCreated, raw.StatusCode,
		"Добавление в избранное должно вернуть 201, тело: %s", string(body))

	var resp FavoriteResponse
	require.NoError(t, json.Unmarshal(body, &resp), "Не удалось распарсить FavoriteResponse")
	return resp
}

// Добавляет объявление в избранное и возвращает HTTP-ответ без проверок
func AddFavoriteRaw(t *testing.T, token string, adID string) *http.Response {
	t.Helper()
	t.Logf("Пытаемся добавить объявление %s в избранное", adID)
	return clientFav.HttpPostFavorite(t, token, adID)
}

// Получает список избранного и парсит ответ
func GetFavorites(t *testing.T, token string, page, pageSize int) FavoriteListResponse {
	t.Helper()
	t.Logf("Получаем список избранного (page=%d, page_size=%d)", page, pageSize)

	raw := clientFav.HttpGetFavorites(t, token, page, pageSize)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа избранного")
	require.Equal(t, http.StatusOK, raw.StatusCode,
		"GET /favorites должен вернуть 200, тело: %s", string(body))

	var resp FavoriteListResponse
	require.NoError(t, json.Unmarshal(body, &resp), "Не удалось распарсить список избранного")
	return resp
}

// Удаляет объявление из избранного и возвращает HTTP-ответ
func DeleteFavoriteRaw(t *testing.T, token string, adID string) *http.Response {
	t.Helper()
	t.Logf("Удаляем объявление %s из избранного", adID)
	return clientFav.HttpDeleteFavorite(t, token, adID)
}

// Удаляет объявление из избранного и проверяет код 204
func DeleteFavorite(t *testing.T, token string, adID string) {
	t.Helper()
	raw := DeleteFavoriteRaw(t, token, adID)
	defer raw.Body.Close()
	require.Equal(t, http.StatusNoContent, raw.StatusCode,
		"Удаление из избранного должно вернуть 204")
}

// Удаляет объявление из избранного, игнорируя код ответа
func DeleteFavoriteIgnoreNotFound(t *testing.T, token string, adID string) {
	t.Helper()
	raw := clientFav.HttpDeleteFavorite(t, token, adID)
	raw.Body.Close()
}

// Проверяет, есть ли объявление с указанным ID в списке избранного
func ContainsAd(list FavoriteListResponse, adID string) bool {
	for _, item := range list.Items {
		if item.Advertisement.ID == adID {
			return true
		}
	}
	return false
}