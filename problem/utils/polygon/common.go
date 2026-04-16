package polygon

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"problem/utils"
	"time"
)

/*
polygonApiCall(): make Polygon API calls
- Remember to do .Body.Close() the response
*/
func polygonApiCall(method string, params map[string]string) (*http.Response, error) {
	apiSecret := os.Getenv("POLYGON_API_SECRET")

	apiSig := ""

	// Create apiSig
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand_header := make([]byte, 6)
	for i := range 6 {
		rand_header[i] = charset[rand.Intn(len(charset))]
	}
	apiSig += fmt.Sprintf("%s/%s?", rand_header, method)

	keys := utils.GetSortedKeys(&params, func(a, b string) bool { return a < b })
	for k, v := range keys {
		if k > 0 {
			apiSig += "&"
		}
		apiSig += fmt.Sprintf("%s=%s", v, params[v])
	}
	apiSig += fmt.Sprintf("#%s", apiSecret)

	sha512 := sha512.New()
	sha512.Write([]byte(apiSig))
	apiSig = hex.EncodeToString(sha512.Sum(nil))

	apiSig = string(rand_header) + apiSig

	address := fmt.Sprintf("https://polygon.codeforces.com/api/%s?", method)
	requestParams := url.Values{}
	requestParams.Add("apiSig", apiSig)
	for key, value := range params {
		requestParams.Add(key, value)
	}

	address += requestParams.Encode()

	resp, err := http.Get(address)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
