package advertisementCreate

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	adv "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/tests"
)

type CreateSuite struct {
	suite.Suite
	token     string
	adID      string
	photoPath string
}

func (s *CreateSuite) SetupSuite() {
	tests.SetupSuite()
	s.token = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD")).Token
	s.photoPath = filepath.Join("..", "..", "..", "internal", "helpers", "testdata", "test.png")
}

func (s *CreateSuite) TearDownTest() {
	if s.adID != "" {
		adv.DeleteAd(s.T(), s.token, s.adID)
		s.adID = ""
	}
}

func uniqueTitle() string {
	return fmt.Sprintf("Test_%d", time.Now().UnixMilli())
}

// Создание объявления со всеми параметрами
func (s *CreateSuite) TestCreatePositive() {
	title := uniqueTitle()
	s.T().Logf(" Создаём объявление с заголовком: %s", title)

	ad := adv.CreateAd(s.T(), s.token, title, "desc", 100, s.photoPath)
	s.adID = ad.ID

	s.Require().NotEmpty(ad.ID, "ID не должен быть пустым")
	s.Require().Equal(title, ad.Title, "Заголовок не совпадает")
	s.Require().NotEmpty(ad.Photos, "Фотографии не загружены")

	s.T().Log("Проверяем GET по ID")
	fetched := adv.GetAd(s.T(), s.token, ad.ID)
	s.Require().Equal(ad.ID, fetched.ID, "ID не совпадает при GET")

	s.T().Log("Проверяем поиск по заголовку")
	list := adv.GetAdvertisementsList(s.T(), s.token, title)
	s.Require().NotEmpty(list.Items, "Объявление не найдено в поиске")

	s.T().Log("Проверяем фотографии")
	photos := adv.GetAdPhotos(s.T(), s.token, ad.ID)
	s.Require().NotEmpty(photos, "Фотографии не найдены")
}

func TestCreateSuite(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}