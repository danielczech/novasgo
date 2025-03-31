package main

import (
	"fmt"
	"log"
	"time"

	nov "github.com/rh-codebase/novasgo/novas"
)

var (
	// Version is the git version at build time
	Version string
	// Build is git hash at buid time
	Build string
)

func main() {
	log.Println("Version: ", Version)
	log.Println("Build: ", Build)

	var ut1Utc float64 = -0.17442
	//var xpole float64 = 0.0824
	//var ypole float64 = 0.4126

	var catEntry nov.CatEntry //nov
	var earth nov.Object      //nov
	nov.MakeObject(0, 3, "earth", &catEntry, &earth)

	var source nov.Object //nov
	nov.MakeObject(0, 10, "sun", &catEntry, &source)
	printBody(source)
	// Other object examples
	//nov.MakeObject(0, 11, "moon", &catEntry, &source)
	//nov.MakeObject(0, 4, "mars", &catEntry, &source)
	//nov.MakeObject(0, 5, "jupiter", &catEntry, &source)
	//nov.MakeObject(0, 1, "mercury", &catEntry, &source)

	var year int16 = 2021
	var month int16 = 4
	var day int16 = 2
	var numberLeapSec int16 = 37
	var hour float64 = 0.0

	year = 2024
	month = 9
	day = 16
	hour = 15.7

	var si nov.OnSurface // nov
	nov.MakeOnSurface(37.2339, -118.282, 1222., 0.0, 0.0, &si)

	for {
		ti := time.Now().UTC()
		year = int16(ti.Year())
		month = int16(ti.Month())
		day = int16(ti.Day())
		hr := ti.Hour()
		min := ti.Minute()
		sec := ti.Second()
		ns := ti.Nanosecond()
		hour = float64(hr) + float64(min)/60. + (float64(sec)+float64(ns)/1e9)/3600.
		fmt.Println("year,month,day,hr,min,sec,ns,dec_hour: ",
			year, month, day, hr, min, sec, ns, hour)
		jdUTC := nov.JulianDate(year, month, day, hour) //nov
		fmt.Println(jdUTC)
		jdTT := jdUTC + (float64(numberLeapSec)+32.184)/86400.
		jdUT1 := jdUTC + ut1Utc/86400.0
		deltaT := 32.184 + float64(numberLeapSec) - ut1Utc
		var gast float64
		nov.SiderealTime(jdUT1, 0.0, deltaT, 1, 1, 0, &gast) //nov
		fmt.Println("gast: ", gast)
		last := gast + si.Longitude/15.0
		fmt.Println("last: ", last)

		/*
			  var mobl, tobl, ee, dpsi, deps float64
				nov.Etilt(jdTT, 10, &mobl, &tobl, &ee, &dpsi, &deps)
				printBody(earth)
				printSite(si)
		*/

		var ra, dec, dis float64
		var accuracy int16 = 0
		err := nov.TopoPlanet(jdTT, &source, deltaT, &si, accuracy, &ra, &dec, &dis) //nov
		if err != nil {
			fmt.Println("Error TopoPlanet: ", err)
		}
		//printSite(si)
		// convert ra,dec to az,el
		var zd, az, rar, decr float64
		doRefraction := int16(0)
		fmt.Println("jdUT1, deltaT, ra, dec: ", jdUT1, deltaT, ra, dec)
		nov.Equ2hor(jdUT1, deltaT, accuracy, 0.0, 0.0, &si, ra, dec, doRefraction, &zd, &az, &rar, &decr) //nov
		fmt.Println("zd,az: ", zd, az)
		el := 90.0 - zd
		azd := int(az)
		azm := int((az - float64(azd)) * 60.)
		azs := (((az - float64(azd)) * 60.) - float64(azm)) * 60.0
		eld := int(el)
		elm := int((el - float64(eld)) * 60.)
		els := (((el - float64(eld)) * 60.) - float64(elm)) * 60.0
		fmt.Println("jdUT1: ", jdUT1)
		fmt.Println("deltaT: ", deltaT)
		fmt.Printf("Topo: %d:%d:%06.3f ra: %6.3f, dec: %6.3f  (%6.3f)%d:%d:%6.3f,(%6.3f)%d:%d:%6.3f: ", hr, min, float64(sec)+float64(ns)/1e9, ra, dec, az, azd, azm, azs, el, eld, elm, els)
		time.Sleep(1 * time.Second)
	}
	// example using internal package. See import above.
	//rtn := example.Ex("Bears")
	//log.Println(rtn)
}

func printBody(b nov.Object) { //nov
	fmt.Println("Body: ")
	fmt.Println("  type[0=Major planet]= ", b.Type)
	fmt.Println("  number[Mer=1...] = ", b.Number)
	fmt.Println("  name= ", b.Name)
}

func printSite(si nov.OnSurface) { //nov
	fmt.Println("Site: ")
	fmt.Println("   geodetic si.latitide[deg] = ", si.Latitude)
	fmt.Println("   geodetic si.longitide[deg] = ", si.Longitude)
	fmt.Println("   si.height = ", si.Height)
	fmt.Println("   si.Temperature[C]= ", si.Temperature)
	fmt.Println("   si.Pressure[mb]= ", si.Pressure)
}
