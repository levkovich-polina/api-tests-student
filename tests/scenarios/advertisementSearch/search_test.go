package advertisementSearch

import (
    "os"
	"strings"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"

    "api-tests-template/internal/assertions"
    "api-tests-template/internal/managers/advertisements"
    "api-tests-template/internal/managers/auth"
    "api-tests-template/tests"
)

const (
    existingAdTitle = "Змейка"      
    searchBeginning = "Зме"         
    searchMiddle    = "мей"         
    searchEnd       = "йка"         
)

type TestSuite struct {
    suite.Suite
    loginData struct {
        Token string
    }
}

func (s *TestSuite) SetupSuite() {
    tests.SetupSuite()

    s.T().Log("Проверка констант...")
    require.NotEmpty(s.T(), existingAdTitle, "задайте existingAdTitle")
    require.NotEmpty(s.T(), searchBeginning, "задайте searchBeginning")
    require.NotEmpty(s.T(), searchMiddle, "задайте searchMiddle")
    require.NotEmpty(s.T(), searchEnd, "задайте searchEnd")

    s.T().Log("Авторизация...")
    loginResp := auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
    s.loginData.Token = loginResp.Token
}

func (s *TestSuite) TestSearchByTitle() {
    list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)

    s.Run("В ответе есть хотя бы одно объявление", func() {
        require.NotEmpty(s.T(), list.Items, "список объявлений пуст")
    })

    s.Run("Среди найденных есть наше объявление", func() {
        assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
    })
}

func (s *TestSuite) TestSearchByTitle_Substrings() {
    cases := []struct {
        name   string
        search string
    }{
        {"начало", searchBeginning},
        {"середина", searchMiddle},
        {"конец", searchEnd},
        {"нижний регистр", strings.ToLower(searchMiddle)},
        {"верхний регистр", strings.ToUpper(searchMiddle)},
        {"смешанный регистр", mixCase(searchMiddle)},
    }

    for _, tc := range cases {
        s.Run(tc.name, func() {
            list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
            assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
        })
    }
}

func mixCase(s string) string {
    runes := []rune(s)
    for i := range runes {
        if i%2 == 0 {
            runes[i] = toLower(runes[i])
        } else {
            runes[i] = toUpper(runes[i])
        }
    }
    return string(runes)
}

func toLower(r rune) rune {
    if r >= 'A' && r <= 'Z' {
        return r + ('a' - 'A')
    }
    return r
}

func toUpper(r rune) rune {
    if r >= 'a' && r <= 'z' {
        return r - ('a' - 'A')
    }
    return r
}

func TestAdvertisementSearchSuite(t *testing.T) {
    suite.Run(t, new(TestSuite))
}