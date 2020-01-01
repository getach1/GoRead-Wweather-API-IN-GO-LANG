//GET https://api.darksky.net/forecast/0123456789abcdef9876543210fedcba/42.3601,-71.0589,255657600?exclude=currentl
//y,flags

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

type Day struct {
	//Time time.Time `json:"time"`
	Summary string `json:"summary"`
	TemperatureHigh float64`json:"temperatureMax"`
	TemperatureLow float64 `json:"temperatureMin"`
	Humidity float64 `json:"humidity"`
	Pressure float64 `json:"pressure"`
	PrecipType string `json:"precipType"`
}

type City struct {
	Latitude float32 `json:"latitude"`
	Longitude float32`json:"longitude"`
	Daily Daily `json:"daily"`
}

type  Daily struct {
	Summary string `json:"summary"`
	Data  []Day `json:"data"`
}

func search(latitudex, longitudey float64) City {
	latitude:=latitudex
	longitude:=longitudey
	exclude:="hourly,offset,flags"
	apikey:="e296602f412717faf216be94ed3728ef"
	url:=fmt.Sprintf("https://api.darksky.net/forecast/%s/%f,%f?exclude=%s",apikey,latitude,longitude,exclude)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("filed to request ",err)
		return City{}
	}
	defer res.Body.Close()
	var city  City

	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body,&city)

	return city
}


func request(w http.ResponseWriter,r *http.Request)  {
	tmpl:=template.Must(template.ParseFiles("form.html","display.html"))
	lat,_:=strconv.ParseFloat(r.FormValue("latitude"),64)
	long,_:=strconv.ParseFloat(r.FormValue("longitude"),64)
	 var daily Daily
	var city City=search(float64(lat), float64(long))
	 daily.Summary=city.Daily.Summary
	 daily.Data=city.Daily.Data
	tmpl.ExecuteTemplate(w,"display.html",daily)
	 for _,data:= range daily.Data{
	 	fmt.Println(data)
	 }
	fmt.Println(city)
}

func form(w http.ResponseWriter,r *http.Request)  {
	tmpl:=template.Must(template.ParseFiles("form.html","display.html"))
	tmpl.ExecuteTemplate(w,"form.html",nil)
}

func main() {
	//url := "https://community-open-weather-map.p.rapidapi.com/forecast/daily?q=san%20francisco%252Cus&lat=35&lon=139&cnt=10&units=metric%20or%20imperial"
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/get",request)
	http.HandleFunc("/",form)
	 http.ListenAndServe(":1000",nil)
}
