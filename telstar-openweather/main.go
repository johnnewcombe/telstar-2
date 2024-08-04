package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar-library/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	MAXFORECASTPAGES        = 12
	ENV_OPENWEATHER_API_KEY = "OPENWEATHER_API_KEY"
	SAMPLE_API_KEY          = "demo"
	// TODO define this outside of the binary?
	BASEURL = "https://api.openweathermap.org/data/2.5/%s?q=%s,uk&units=metric&appid=%s"
)

//go:embed templates/weather.json
var templateWeather string

//go:embed templates/forecast.json
var templateForecast string

//go:embed sample-data/weather-sample.json
var sampleWeather string

//go:embed sample-data/forecast-sample.json
var sampleForecast string

func main() {

	var (
		weather    WeatherResponse
		forecast   ForecastResponse
		path       string
		output     strings.Builder
		content    string
		err        error
		city       string
		outputText []byte
	)

	// syntax [api_key] city
	// the api open weather api can be specified as an argument or as an environment variable
	// if the environment variable OPENWEATHER_API_KEY is detected, the utilty will expect only
	// one parameter, that of the city. If the environment variable is not set then two arguments
	// are expected, the api key and the city.

	apiKey := os.Getenv(ENV_OPENWEATHER_API_KEY)
	argCount := len(os.Args)

	if len(apiKey) > 0 {
		// we have environment variable set so looking fpr two arguments (including the cmd)
		if len(os.Args) != 2 {
			log.Fatalf("incorrect number of arguments api key is set in an Environment Variable therefore 1 argument was expected (%d received).", argCount-1) //this could be due to the api key being set as an Environment variable and passed as a parameter.")
		}
		city = os.Args[1]
	} else {
		if len(os.Args) != 3 {
			log.Fatalf("incorrect number of arguments api key is not set in an Environment Variable therefore 2 arguments were expected (%d received).", argCount-1) //this could be due to the api key being set as an Environment variable and passed as a parameter.")
		}
		apiKey = os.Args[1]
		city = os.Args[2]
	}

	path, err = os.Executable()
	if err != nil {
		log.Fatalf("get pathname, %s, %v", path, err)
	}

	if weather, err = getWeather(city, apiKey); err != nil {
		log.Fatalf("get weather, %v", err)
	}

	if len(weather.Name) > 0 {
		content = createWeatherPage(weather)
	} else {
		content = createNoDataPage()
	}

	// load the template into a frame and allows us to get the required frame id for the results
	wf := types.Frame{}
	if err = wf.Load([]byte(templateWeather)); err != nil {
		log.Fatal(err)
	}

	wf.Content.Data = strings.ReplaceAll(wf.Content.Data, "[CONTENT]", content)

	// update the content for the weather page and add to the output
	if outputText, err = wf.Dump(); err != nil {
		log.Fatalf("write output, %v", err)
	}
	if _, err = output.Write(outputText); err != nil {
		log.Fatal(err)
	}
	if _, err = output.WriteString(","); err != nil {
		log.Fatal(err)
	}

	if len(weather.Name) > 0 {

		// get forecast
		if forecast, err = getForecast(city, apiKey); err != nil {
			log.Fatalf("get forecast, %v", err)
		}

		// add the forecast pages
		forecastPages := createForecastPages(forecast)
		frameId := rune(wf.PID.FrameId[0])
		pageNumber := wf.PID.PageNumber

		// iterate through the first 15 forcast pages
		for f := 0; f < len(forecastPages) && f < MAXFORECASTPAGES; f++ {

			forecastPage := forecastPages[f]

			// calculate the next PID
			if pageNumber, frameId, err = utils.GetFollowOnPID(pageNumber, frameId); err != nil {
				log.Fatal(err)
			}
			// load the template into a frame and allows us to get the required frame id for the results
			ff := types.Frame{}
			ff.Load([]byte(templateForecast))
			ff.Content.Data = strings.ReplaceAll(ff.Content.Data, "[CONTENT]", forecastPage)
			ff.PID.PageNumber = pageNumber
			ff.PID.FrameId = string(frameId)

			// update the content for the weather page and add to the output
			if outputText, err = ff.Dump(); err != nil {
				log.Fatalf("write output, %v", err)
			}
			if _, err = output.Write(outputText); err != nil {
				log.Fatal(err)
			}
			output.WriteString(",")
		}
	}

	outs := output.String()
	fmt.Printf("[%s]", outs[:len(outs)-1])

	// FIXME return error frame
	//  could this be the template frame with the content overwritten?
}

func getWeather(city string, apiKey string) (WeatherResponse, error) {

	var (
		weatherContent string
		err            error
		url            string
	)

	url = fmt.Sprintf(BASEURL, "weather", city, apiKey)

	useSampleData := strings.ToLower(apiKey) == SAMPLE_API_KEY

	weather := WeatherResponse{}
	weatherError := ErrorResponse{}

	if useSampleData {
		weatherContent = sampleWeather
	} else {
		if weatherContent, err = getData(url); err != nil {
			return weather, err
		}
	}

	if len(weatherContent) == 0 {
		return weather, nil
	}

	//file.WriteFile("weather-sample.json",[]byte(weatherContent))

	if err = json.Unmarshal([]byte(weatherContent), &weather); err != nil {
		// we will get this error if the search term wasn't found
		if err = json.Unmarshal([]byte(weatherContent), &weatherError); err != nil {
			return weather, fmt.Errorf("%v, %s", err, weatherError.Message)
		}
		return weather, err
	}

	if useSampleData {
		if len(city) > 2 {
			//weather.Name = utils.ToUpperCamelCase(city)
			weather.Name = strings.ToUpper(city)
		}
		weather.Dt = time.Now().Unix()
	}

	return weather, nil
}

func getForecast(city string, apiKey string) (ForecastResponse, error) {

	var (
		forecastContent string
		err             error
		url             string
	)

	url = fmt.Sprintf(BASEURL, "forecast", city, apiKey)

	useSampleData := strings.ToLower(apiKey) == SAMPLE_API_KEY

	forecast := ForecastResponse{}
	forecastError := ErrorResponse{}

	if strings.ToLower(apiKey) == SAMPLE_API_KEY {
		forecastContent = sampleForecast
	} else {
		if forecastContent, err = getData(url); err != nil {
			return forecast, err
		}
	}

	if len(forecastContent) == 0 {
		return forecast, nil
	}

	//file.WriteFile("forecast-sample.json",[]byte(forecastContent))

	if err = json.Unmarshal([]byte(forecastContent), &forecast); err != nil {
		// we will get this error if the search term wasn't found
		if err = json.Unmarshal([]byte(forecastContent), &forecastError); err != nil {
			return forecast, errors.New(forecastError.Message)
		}
		return forecast, err
	}

	if useSampleData {
		if len(city) > 2 {
			forecast.City.Name = strings.ToUpper(city)
		}
		for i := 0; i < len(forecast.List); i++ {
			// fixme need to make time in three hour chunks 15:00, 18:00 from the current time onwards
			forecast.List[i].Dt = utils.TruncateToStartOfHour(time.Now().Local().Add(time.Hour * time.Duration(3*(i+1)))).Unix()
			//forecast.List[i].Dt = utils.TruncateToStartOfHour(time.Now()).Unix() // + time.Hour*index
			forecast.List[i].DtTxt = ""
		}
	}
	return forecast, nil
}

func getData(url string) (string, error) {

	var (
		client     http.Client
		bodyString string
		bodyBytes  []byte
	)

	resp, err := client.Get(url)
	if err != nil {
		return bodyString, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return bodyString, err
		}
		bodyString = string(bodyBytes)
	}
	return bodyString, nil
}
