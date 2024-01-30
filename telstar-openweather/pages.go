package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"fmt"
	"strings"
)

func createWeatherPage(weather WeatherResponse) string {

	var outputWeather strings.Builder

	outputWeather.WriteString("[W]") // start with white text just below the title
	outputWeather.WriteString(weather.Name)
	outputWeather.WriteString(", ")
	outputWeather.WriteString(weather.Sys.Country)
	outputWeather.WriteString("\r\n")
	outputWeather.WriteString("[W]") // White text
	outputWeather.WriteString(utils.GetDateTimeFromUnixData(weather.Dt))
	outputWeather.WriteString("\r\n[b][m.]\r\n") // dots horizontal line
	outputWeather.WriteString("[D]")                 // double height
	outputWeather.WriteString(weather.Weather[0].Main)
	outputWeather.WriteString(": ")
	outputWeather.WriteString(weather.Weather[0].Description)
	outputWeather.WriteString("\n\r\n[b][m.]\r\n") // dots horizontal line
	outputWeather.WriteString("[Y]Temp             :[W]")
	outputWeather.WriteString(fmt.Sprintf("%02.2f C\r\n", weather.Main.Temp))
	outputWeather.WriteString("[Y]Humidity         :[W]")
	outputWeather.WriteString(fmt.Sprintf("%d %%\r\n", weather.Main.Humidity))
	outputWeather.WriteString("[Y]Pressure         :[W]")
	outputWeather.WriteString(fmt.Sprintf("%d mb\r\n", weather.Main.Pressure))
	outputWeather.WriteString("[Y]Sunrise/Sunset   :[W]")
	outputWeather.WriteString(fmt.Sprintf("%s/%s\r\n", utils.GetTimeFromUnixData(weather.Sys.Sunrise), utils.GetTimeFromUnixData(weather.Sys.Sunset)))
	outputWeather.WriteString("[Y]Wind Speed       :[W]")
	outputWeather.WriteString(fmt.Sprintf("%2.2f mph\r\n", weather.Wind.Speed))
	outputWeather.WriteString("[Y]Wind Direction   :[W]")
	outputWeather.WriteString(fmt.Sprintf("%d deg\r\n", weather.Wind.Deg))
	//	outputWeather.WriteString("")
	//	outputWeather.WriteString("")
	//	outputWeather.WriteString("")
	//	outputWeather.WriteString("")
	//	outputWeather.WriteString("")

	return outputWeather.String()

}

func createForecastPages(forecast ForecastResponse) []string {

	var output []string

	for l := 0; l < len(forecast.List); l++ {
		var outputForecast strings.Builder

		outputForecast.WriteString("[W]Forecast for ")
		outputForecast.WriteString(utils.GetDateTimeFromUnixData(forecast.List[l].Dt))
		outputForecast.WriteString("\n\r[b][m.]\r\n")
		outputForecast.WriteString("[D]") // double height
		outputForecast.WriteString(forecast.List[l].Weather[0].Main)
		outputForecast.WriteString(": ")
		outputForecast.WriteString(forecast.List[l].Weather[0].Description)
		outputForecast.WriteString("\n\r\n[b][m.]\r\n")
		outputForecast.WriteString("[Y]Temp             :[W]")
		outputForecast.WriteString(fmt.Sprintf("%02.2f C\r\n", forecast.List[l].Main.Temp))
		outputForecast.WriteString("[Y]Humidity         :[W]")
		outputForecast.WriteString(fmt.Sprintf("%d %%\r\n", forecast.List[l].Main.Humidity))
		outputForecast.WriteString("[Y]Pressure         :[W]")
		outputForecast.WriteString(fmt.Sprintf("%d mb\r\n", forecast.List[l].Main.Pressure))
		outputForecast.WriteString("[Y]Wind Speed       :[W]")
		outputForecast.WriteString(fmt.Sprintf("%2.2f mph\r\n", forecast.List[l].Wind.Speed))
		outputForecast.WriteString("[Y]Wind Direction   :[W]")
		outputForecast.WriteString(fmt.Sprintf("%d deg\r\n", forecast.List[l].Wind.Deg))

		output = append(output, outputForecast.String())
	}

	return output
}

func createNoDataPage() string {
	return "\\r\\n\\r\\n\\r\\n  [D]Nothing found press _ to try again"
}
