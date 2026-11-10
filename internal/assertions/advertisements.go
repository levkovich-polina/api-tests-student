package assertions

import (
    "testing"

    "github.com/stretchr/testify/require"
    "api-tests-template/internal/managers/advertisements"
)

func ContainsAdvertisementByTitle(t *testing.T, items []advertisements.Advertisement, title string) {
    for _, item := range items {
        if item.Title == title {
            return 
        }
    }
    require.Fail(t, "Объявление не найдено", title)
}