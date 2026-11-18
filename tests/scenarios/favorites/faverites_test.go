package favorites

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	apiRunner "api-tests-template/internal/helpers/api-runner"
	adv "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
)

func setupUsers(t *testing.T) (buyerToken, sellerToken string) {
	t.Helper()
	buyerToken = auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD")).Token
	require.NotEmpty(t, buyerToken, "Токен покупателя пустой")

	sellerToken, _ = RegisterSecondUser(t)
	require.NotEmpty(t, sellerToken, "Токен продавца пустой")

	return buyerToken, sellerToken
}

func RegisterSecondUser(t *testing.T) (token, userID string) {
	t.Helper()
	email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
	registerReq := map[string]string{
		"email":      email,
		"first_name": "Test",
		"last_name":  "User",
		"password":   "TestPass123!",
	}

	resp := apiRunner.GetRunner().Create().
		Post("/auth/register").
		JSON(registerReq).
		Expect(t).End().Response
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Регистрация должна вернуть 201")

	var parsed struct {
		Token string `json:"token"`
		User  struct {
			ID string `json:"id"`
		} `json:"user"`
	}
	body, _ := io.ReadAll(resp.Body)
	require.NoError(t, json.Unmarshal(body, &parsed), "Не удалось распарсить ответ регистрации")
	require.NotEmpty(t, parsed.Token, "Токен второго пользователя пустой")
	require.NotEmpty(t, parsed.User.ID, "ID второго пользователя пустой")

	return parsed.Token, parsed.User.ID
}

func CreateAdvertisementAsOwner(t *testing.T, sellerToken, photoPath string) string {
	t.Helper()
	title := fmt.Sprintf("TestAd_%d", time.Now().UnixNano())
	ad := adv.CreateAd(t, sellerToken, title, "Описание для теста избранного", 100, photoPath)

	t.Cleanup(func() {
		adv.DeleteAd(t, sellerToken, ad.ID)
	})

	return ad.ID
}

func LoginCurrentUser(t *testing.T) string {
	t.Helper()
	loginResp := auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	require.NotEmpty(t, loginResp.Token, "Токен текущего пользователя пустой")
	return loginResp.Token
}