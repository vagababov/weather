// Fetcher is a program that allows fetching historical data from the NOAA site.
// NOAA site: http://www.wrh.noaa.gov/otx/climate/lcd/lcd.php
// (C) Victor Agababov (vagababov@gmail.com)
// Author maintains no responsibility for the functioning of this
// program. Use at your own risk.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// Flags.
var (
	startDate = flag.String("start_date", "", "Starting date in YYYYmmdd format. "+
		"Must be the last day of the month")
	endDate = flag.String("end_date", "", "Ending date in YYYYmmdd format. Must be "+
		"the last day of the month")
	station     = flag.String("station", "sew", "Meteo station to fetch data for.")
	noaaStation = flag.String("noaa_station", "sew", "NOAA center station.")
	outputDir   = flag.String("output_dir", "./", "Output directory")
)

// January is 1.
var daysInMonth = []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

const df = "20060102"

func fetchPost(date, noaaStation, station string) {
	uri := fmt.Sprintf("http://www.weather.gov/climate/getclimate.php?wfo=%s", noaaStation)
	resp, err := http.PostForm(uri, url.Values{
		"product":          {"CF6"},
		"printer_friendly": {"yes"},
		"station":          {station},
		"recent":           {"no"},
		"date":             {date},
	})
	if err != nil {
		fmt.Printf("Error fetching data for station: %s/%s date: %s: %v",
			noaaStation, station, date, err)
		return
	}

	fn := fmt.Sprintf("%s-%s", station, date)
	p := path.Join(*outputDir, fn)
	f, err := os.Create(p)
	if err != nil {
		fmt.Printf("Error creating file %s: %v", p, err)
		return
	}
	defer f.Close()

	io.Copy(f, resp.Body)
}

func validateFlags() bool {
	// Date format.
	sd, err := time.Parse(df, *startDate)
	if err != nil {
		fmt.Printf("time.Parse(%s): %v\n", *startDate, err)
		return false
	}
	ed, err := time.Parse(df, *endDate)
	if err != nil {
		fmt.Printf("time.Parse(%s): %v\n", *endDate, err)
		return false
	}
	if ed.Before(sd) {
		fmt.Printf("StartDate(%v) must not be larger than EndDate(%v)\n", sd, ed)
		return false
	}
	if len(*noaaStation) != 3 {
		fmt.Printf("NOAA station must be present and be 3 characters long; was: %s\n", *noaaStation)
		return false
	}
	if len(*station) != 3 {
		fmt.Printf("station must be present and be 3 characters long; was: %s\n", *station)
		return false
	}

	return true
}

func fetch(year int, month string) bool {
	fmt.Printf("Fetching Year: %d Month: %s\n", year, month)
	_, err := http.Get("http://www.wrh.noaa.gov/otx/climate/lcd/lcd.php?yr=15&mon=jan&stn=eat")
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return false
	}
	return true
}

func nextDate(d time.Time) time.Time {
	y, m := d.Year(), d.Month()+1
	if m == 13 {
		m = 1
		y++
	}
	return time.Date(y, m, daysInMonth[m], 0, 0, 0, 0, time.Local)
}

func main() {
	flag.Parse()
	if !validateFlags() {
		return
	}

	// We know those dates parse successfully.
	sd, _ := time.Parse(df, *startDate)
	ed, _ := time.Parse(df, *endDate)
	for d := sd; d.Before(ed) || d.Equal(ed); {
		ds := d.Format("20060102")
		fmt.Printf("Fetching for date: %s\n", ds)
		fetchPost(ds, *noaaStation, *station)
		d = nextDate(d)
	}
	fmt.Printf("\nDone!\n")
}
