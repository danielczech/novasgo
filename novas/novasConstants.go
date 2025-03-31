/*
  Naval Observatory Vector Astrometry Software (NOVAS)
  C Edition, Version 3.1

  novascon.h: Header file for novascon.c
  novascon.c: Constants for use with NOVAS

  U. S. Naval Observatory
  Astronomical Applications Dept.
  Washington, DC
  http://www.usno.navy.mil/USNO/astronomical-applications
*/
package novas

import "math"

const (
	// FN1 int16
	// FN0 int16
	NUM_TARGETS      = 12 // ie. number of bodies
	SIZE_OF_OBJ_NAME = 51
	SIZE_OF_CAT_NAME = 4

	// Define "origin" constants.
	BARYC  = 0
	HELIOC = 1

	// Approximate scale height of atmosphere in meters.
	AtmospherScaleHeight = float64(9.1e3)

	//   TDB Julian date of epoch J2000.0.
	T0 = float64(2451545.00000000)

	//Speed of light in meters/second is a defining physical constant.
	C = float64(299792458.0)

	//Light-time for one astronomical unit (AU) in seconds, from DE-405.
	AU_SEC = float64(499.0047838061)

	//Speed of light in AU/day.  Value is 86400 / AU_SEC.
	C_AUDAY = float64(173.1446326846693)

	// Astronomical unit in meters.  Value is AU_SEC * C.
	AU = float64(1.4959787069098932e+11)

	//Astronomical Unit in kilometers.
	AU_KM = float64(1.4959787069098932e+8)

	// Heliocentric gravitational constant in meters^3 / second^2, from
	// DE-405.
	GS = float64(1.32712440017987e+20)

	// Geocentric gravitational constant in meters^3 / second^2, from
	// DE-405.
	GE = float64(3.98600433e+14)

	// Radius of Earth in meters from IERS Conventions (2003).
	ERAD = float64(6378136.6)

	RADE = ERAD / AU

	// Earth ellipsoid flattening from IERS Conventions (2003).
	// Value is 1 / 298.25642.
	F = float64(0.003352819697896)

	// Rotational angular velocity of Earth in radians/sec from IERS
	// Conventions (2003).
	ANGVEL = float64(7.2921150e-5)

	// Reciprocal masses of solar system bodies, from DE-405
	// (Sun mass / body mass).
	// MASS[0] = Earth/Moon barycenter, MASS[1] = Mercury, ...,
	// MASS[9] = Pluto, MASS[10] = Sun, MASS[11] = Moon.

	// Value of 2 * pi in radians.
	TWOPI  = float64(6.283185307179586476925287)
	HALFPI = math.Pi / 2.0

	// Number of arcseconds in 360 degrees.
	ASEC360 = float64(1296000.0)

	// Angle conversion constants.
	ASEC2RAD = float64(4.848136811095359935899141e-6)
	DEG2RAD  = float64(0.017453292519943296)
	RAD2DEG  = float64(57.295779513082321)
)

// workaround GO not allowing constant arrays
func RMASS() [NUM_TARGETS]float64 {
	return [12]float64{328900.561400, 6023600.0, 408523.71,
		332946.050895, 3098708.0, 1047.3486, 3497.898, 22902.98,
		19412.24, 135200000.0, 1.0, 27068700.387534}
}
