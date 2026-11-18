package advertisementCreate

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	clientAdv "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/tests"
)

type CreateNegativeSuite struct {
	suite.Suite
	token     string
	photoPath string
}

func (s *CreateNegativeSuite) SetupSuite() {
	tests.SetupSuite()
	s.token = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD")).Token
	s.photoPath = filepath.Join("..", "..", "..", "internal", "helpers", "testdata", "test.png")
}

// Создание без авторизации.
func (s *CreateNegativeSuite) TestCreateNoAuth() {
	s.T().Log("Создание объявления без токена")
	raw := clientAdv.HttpPostAdvertisement(s.T(), "", "t", "d", 100, s.photoPath)
	raw.Body.Close()
	s.Require().Equal(http.StatusUnauthorized, raw.StatusCode, "Без токена должен быть 401")
}

// Создание с пустым title.
func (s *CreateNegativeSuite) TestCreateEmptyTitle() {
	s.T().Log("Создание объявления с пустым title")
	raw := clientAdv.HttpPostAdvertisement(s.T(), s.token, "", "d", 100, s.photoPath)
	raw.Body.Close()
	s.Require().Equal(http.StatusBadRequest, raw.StatusCode, "Пустой title должен дать 400")
}

func TestCreateNegativeSuite(t *testing.T) {
	suite.Run(t, new(CreateNegativeSuite))
}