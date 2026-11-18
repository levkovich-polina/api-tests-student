package favorites

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	adv "api-tests-template/internal/managers/advertisements"
	fav "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests"
)

type PositiveSuite struct {
	suite.Suite
	userToken  string
	ownerToken string
	photoPath  string
}

func (s *PositiveSuite) SetupSuite() {
	tests.SetupSuite()
	s.userToken, s.ownerToken = setupUsers(s.T())
	s.photoPath = filepath.Join("..", "..", "..", "internal", "helpers", "testdata", "test.png")
}

func (s *PositiveSuite) TestFavoritesFullCycle() {
	adID := CreateAdvertisementAsOwner(s.T(), s.ownerToken, s.photoPath)

	beforeTotal := fav.GetFavorites(s.T(), s.userToken, 1, 50).Total

	added := fav.AddFavorite(s.T(), s.userToken, adID)
	s.Require().Equal(adID, added.Advertisement.ID, "В избранное добавлено другое объявление")
	s.Require().NotEmpty(added.FavoriteID, "Не заполнен favorite_id")
	s.Require().NotEmpty(added.AddedAt, "Не заполнено время добавления")

	list := fav.GetFavorites(s.T(), s.userToken, 1, 50)
	s.Require().True(fav.ContainsAd(list, adID), "Объявление %s отсутствует в избранном", adID)
	s.Require().Equal(beforeTotal+1, list.Total, "После добавления total должен увеличиться на 1")

	fav.DeleteFavorite(s.T(), s.userToken, adID)

	listAfter := fav.GetFavorites(s.T(), s.userToken, 1, 50)
	s.Require().False(fav.ContainsAd(listAfter, adID), "Объявление %s всё ещё в избранном", adID)
	s.Require().Equal(beforeTotal, listAfter.Total, "После удаления total должен вернуться к исходному")
}

func (s *PositiveSuite) TestDeleteAdvertisementRemovesFromFavorites() {
	adID := CreateAdvertisementAsOwner(s.T(), s.ownerToken, s.photoPath)

	fav.AddFavorite(s.T(), s.userToken, adID)
	s.Require().True(fav.ContainsAd(fav.GetFavorites(s.T(), s.userToken, 1, 50), adID),
		"Объявление должно быть в избранном до удаления")

	adv.DeleteAd(s.T(), s.ownerToken, adID)

	list := fav.GetFavorites(s.T(), s.userToken, 1, 50)
	s.Require().False(fav.ContainsAd(list, adID),
		"Удалённое объявление %s не должно оставаться в избранном", adID)
}

func TestPositiveSuite(t *testing.T) {
	suite.Run(t, new(PositiveSuite))
}