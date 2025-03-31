/*
Naval Observatory Vector Astrometry Software (NOVAS)
C Edition, Version 3.1

checkout-stars-full.c: Checkout program for use with solsys1 in full-accuracy mode

U. S. Naval Observatory
Astronomical Applications Dept.
Washington, DC
http://www.usno.navy.mil/USNO/astronomical-applications
*/
package novas

import (
	"fmt"
	"math"
	"testing"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include "eph_manager.h"
#include "novas.h"

#define N_STARS 3
#define N_TIMES 4
*/

const (
	N_STARS = 3
	N_TIMES = 4
	ABS_ERR = 1.e-7
)

// int main (void)
// {
func TestStarsFull(t *testing.T) {
	/*
	   Main function to check out many parts of NOVAS-C by calling
	   function 'topo_star' with version 1 of function 'solarsystem'.

	   For use with NOVAS-C Version 3.1.
	*/
	var expected_ra = [N_TIMES][N_STARS]float64{
		{2.446988922, 5.530110723, 10.714525513},
		{2.446988922, 5.530110723, 10.714525513},
		{2.509480139, 5.531195904, 10.714444761},
		{2.481177533, 5.530372288, 10.713575394}}
	var expected_dec = [N_TIMES][N_STARS]float64{
		{89.24635149, -0.30571737, -64.38130590},
		{89.24635149, -0.30571737, -64.38130590},
		{89.25196813, -0.30301961, -64.37366514},
		{89.24254404, -0.30231606, -64.3796699}}

	//short int error = 0;
	//short int accuracy = 0;
	//short int i, j, de_num;
	error := int16(0)
	accuracy := int16(0)
	var de_num int16

	/*
	   'deltat' is the difference in time scales, TT - UT1.

	    The array 'tjd' contains four selected Julian dates at which the
	    star positions will be evaluated.
	*/

	//double deltat = 60.0;
	//double tjd[N_TIMES] = {2450203.5, 2450203.5, 2450417.5, 2450300.5};
	//double jd_beg, jd_end, ra, dec;
	deltat := float64(60.0)
	tjd := []float64{2450203.5, 2450203.5, 2450417.5, 2450300.5}
	var jd_beg, jd_end, ra, dec float64

	/*
	   Hipparcos (ICRS) catalog data for three selected stars.
	*/

	stars := [N_STARS]CatEntry{
		{"POLARIS", "HIP", 0, 2.530301028, 89.264109444,
			44.22, -11.75, 7.56, -17.4},
		{"Delta ORI", "HIP", 1, 5.533444639, -0.299091944,
			1.67, 0.56, 3.56, 16.0},
		{"Theta CAR", "HIP", 2, 10.715944806, -64.394450000,
			-18.87, 12.06, 7.43, 24.0}}

	/*
	   The observer's terrestrial coordinates (latitude, longitude, height).
	*/

	geo_loc := OnSurface{45.0, -75.0, 0.0, 10.0, 1010.0}

	/*
	   Open the JPL ephemeris file.
	*/
	if error = EphemOpen("JPLEPH", &jd_beg, &jd_end, &de_num); error != 0 {
		fmt.Println("Error opending JPLEPH")
		t.Fail()
	}

	/*
	   Compute the topocentric places of the three stars at the four
	   selected Julian dates.
	*/

	for idx := 0; idx < N_TIMES; idx++ {
		for jdx := 0; jdx < N_STARS; jdx++ {
			if err := TopoStar(tjd[idx], deltat, &stars[jdx], &geo_loc,
				accuracy, &ra, &dec); err != nil {
				fmt.Printf("Error %v from topo_star. Star %d  Time %d\n",
					err, jdx, idx)
				t.Fail()
			} else {
				if math.Abs(ra-expected_ra[idx][jdx]) > ABS_ERR {
					fmt.Printf("Got ra= %18.12f, Expected: %18.12f\n", ra, expected_ra[idx][jdx])
					fmt.Printf("JD = %f  Star = %s\n", tjd[idx], stars[jdx].Starname)
					fmt.Printf("RA = %18.12f  Dec = %18.12f\n", ra, dec)
					fmt.Printf("\n")
					t.Fail()
				}
				if math.Abs(dec-expected_dec[idx][jdx]) > ABS_ERR {
					fmt.Printf("Got dec= %18.12f, Expected: %18.12f\n", dec, expected_dec[idx][jdx])
					fmt.Printf("JD = %f  Star = %s\n", tjd[idx], stars[jdx].Starname)
					fmt.Printf("RA = %18.12f  Dec = %18.12f\n", ra, dec)
					fmt.Printf("\n")
					t.Fail()
				}
			}
		}
		fmt.Printf("\n")
	}
	EPHFILE.Close()
}
