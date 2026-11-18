package favorites

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	clientFav "api-tests-template/internal/client/http/favorites"
	adv "api-tests-template/internal/managers/advertisements"
	fav "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests"
)

const nonexistentAdID = "11111111-1111-1111-1111-111111111111"

type NegativeSuite struct {
	suite.Suite
	userToken  string
	ownerToken string
	photoPath  string
}

func (s *NegativeSuite) SetupSuite() {
	tests.SetupSuite()
	s.userToken, s.ownerToken = setupUsers(s.T())
	s.photoPath = filepath.Join("..", "..", "..", "internal", "helpers", "testdata", "test.png")
}

func (s *NegativeSuite) TestFavoritesWithoutAuth() {
	adID := CreateAdvertisementAsOwner(s.T(), s.ownerToken, s.photoPath)

	endpoints := []struct {
		name    string
		request func(t *testing.T) *http.Response
	}{
		{"POST /favorites", func(t *testing.T) *http.Response {
			return clientFav.HttpPostFavorite(t, "", adID)
		}},
		{"GET /favorites", func(t *testing.T) *http.Response {
			return clientFav.HttpGetFavorites(t, "", 1, 15)
		}},
		{"DELETE /favorites/{ad_id}", func(t *testing.T) *http.Response {
			return clientFav.HttpDeleteFavorite(t, "", adID)
		}},
	}

	for _, ep := range endpoints {
		s.T().Run(ep.name, func(t *testing.T) {
			resp := ep.request(t)
			resp.Body.Close()
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode,
				"%s без токена должен вернуть 401", ep.name)
		})
	}
}

func (s *NegativeSuite) TestNonExistentAd() {
	cases := []struct {
		name    string
		request func(t *testing.T) *http.Response
	}{
		{"добавление", func(t *testing.T) *http.Response {
			return fav.AddFavoriteRaw(t, s.userToken, nonexistentAdID)
		}},
		{"удаление", func(t *testing.T) *http.Response {
			return fav.DeleteFavoriteRaw(t, s.userToken, nonexistentAdID)
		}},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			resp := c.request(t)
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			require.Containsf(t, []int{http.StatusBadRequest, http.StatusNotFound}, resp.StatusCode,
				"%s несуществующего объявления должно вернуть 400 или 404, тело: %s", c.name, string(body))
		})
	}
}

func (s *NegativeSuite) TestAddDuplicateFavorite() {
	adID := CreateAdvertisementAsOwner(s.T(), s.ownerToken, s.photoPath)

	fav.AddFavorite(s.T(), s.userToken, adID)

	resp := fav.AddFavoriteRaw(s.T(), s.userToken, adID)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	s.Require().Equal(http.StatusConflict, resp.StatusCode,
		"Повторное добавление должно вернуть 409, тело: %s", string(body))
}

func (s *NegativeSuite) TestAddOwnAdvertisement() {
	title := fmt.Sprintf("Own_%d", time.Now().UnixNano())
	ad := adv.CreateAd(s.T(), s.userToken, title, "Описание своего объявления", 100, s.photoPath)

	s.T().Cleanup(func() {
		adv.DeleteAd(s.T(), s.userToken, ad.ID)
	})

	resp := fav.AddFavoriteRaw(s.T(), s.userToken, ad.ID)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	s.Require().Equal(http.StatusForbidden, resp.StatusCode,
		"Добавление своего объявления должно вернуть 403, тело: %s", string(body))
}

func TestNegativeSuite(t *testing.T) {
	suite.Run(t, new(NegativeSuite))
}