package polygon

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"problem/models"
	"strconv"
	"time"
)

// type packages struct {
// 	pkgs []models.Package
// }

type returnStruct struct {
	Status string           `json:"status"`
	Result []models.Package `json:"result"`
}

func GetLastestPackage(problemId uint64) (uint64, error) {
	resp, err := polygonApiCall("problem.packages", map[string]string{
		"problemId": strconv.Itoa(int(problemId)),
		"apiKey":    os.Getenv("POLYGON_API_KEY"),
		"time":      fmt.Sprintf("%d", time.Now().Unix()),
	})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("get error while getting latest packageId: %v", resp)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// fmt.Println(string(body))

	var return_json returnStruct
	if err := json.Unmarshal(body, &return_json); err != nil {
		return 0, err
	}

	for _, pkg := range return_json.Result {
		if pkg.State != "READY" {
			continue

		}
		if pkg.Type == "linux" {
			continue
		}

		return pkg.ID, nil
	}

	return 0, fmt.Errorf("no valid packages found")
}
