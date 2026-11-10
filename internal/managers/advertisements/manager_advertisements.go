package advertisements

import (
    "encoding/json"
    "io"
    "testing"

    "github.com/stretchr/testify/require"
    "api-tests-template/internal/client/http/advertisements"
)

func GetAdvertisementsList(t *testing.T, token, text string) AdvertisementListResponse {
    raw := advertisements.HttpGetAdvertisements(t, token, text)
    defer raw.Body.Close()

    body, err := io.ReadAll(raw.Body)
    require.NoError(t, err)

    var resp AdvertisementListResponse
    require.NoError(t, json.Unmarshal(body, &resp))

    return resp
}