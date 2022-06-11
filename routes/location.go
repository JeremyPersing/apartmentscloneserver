package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kataras/iris/v12"
)

func Autocomplete(ctx iris.Context) {
	limit := "10"
	location := ctx.URLParam("location")
	limitQuery := ctx.URLParam("limit")
	if limitQuery != "" {
		limit = limitQuery
	}

	apiKey := os.Getenv("LOCATION_TOKEN")
	url :=
		"https://api.locationiq.com/v1/autocomplete.php?key=" + apiKey + "&q=" + location + "&limit=" + limit

	fetchLocations(url, ctx)
}

func Search(ctx iris.Context) {
	location := ctx.URLParam("location")

	apiKey := os.Getenv("LOCATION_TOKEN")
	url :=
		"https://api.locationiq.com/v1/search.php?key=" + apiKey + "&q=" + location + "&format=json"

	fetchLocations(url, ctx)
}

func fetchLocations(url string, ctx iris.Context) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	res, locationErr := client.Do(req)
	if locationErr != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail("Internal Server Error"))
		return
	}

	defer res.Body.Close()

	body, bodyErr := ioutil.ReadAll(res.Body)
	if bodyErr != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail("Internal Server Error"))
		return
	}

	var objMap []map[string]interface{}
	jsonErr := json.Unmarshal(body, &objMap)
	if jsonErr != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail("Internal Server Error"))
		return
	}

	ctx.JSON(objMap)
}
