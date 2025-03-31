// ********terra
package novas

import "math"

var (
	//static short int first_entry = 1
	//static double erad_km, ht_km
	first_entry    = int16(1)
	erad_km, ht_km float64
)

/*
------------------------------------------------------------------------

   PURPOSE:
      Computes the position and velocity vectors of a terrestrial
      observer with respect to the center of the Earth.

   REFERENCES:
      Kaplan, G. H. et. al. (1989). Astron. Journ. 97, 1197-1210.

   INPUT
   ARGUMENTS:
      *location (struct OnSurface)
         Pointer to structure containing observer's location (defined
         in novas.h).
      st (double)
         Local apparent sidereal time at reference meridian in hours.

   OUTPUT
   ARGUMENTS:
      pos[3] (double)
         Position vector of observer with respect to center of Earth,
         equatorial rectangular coordinates, referred to true equator
         and equinox of date, components in AU.
      vel[3] (double)
         Velocity vector of observer with respect to center of Earth,
         equatorial rectangular coordinates, referred to true equator
         and equinox of date, components in AU/day.

   RETURNED
   VALUE:
      None.

   GLOBALS
   USED:
      AU_KM, ERAD, F     novascon.c
      ANGVEL, DEG2RAD    novascon.c

   FUNCTIONS
   CALLED:
      sin                math.h
      cos                math.h
      sqrt               math.h

   VER./DATE/
   PROGRAMMER:
      V1.0/04-93/WTH (USNO/AA):  Translate Fortran.
      V1.1/06-98/JAB (USNO/AA):  Move constants 'f' and 'omega' to
                                 file 'novascon.c'.
      V1.2/10-03/JAB (USNO/AA):  Updates Notes; removed call to 'pow'.
      V1.3/12-04/JAB (USNO/AA):  Update to use 'OnSurface" structure.
      V1.4/09-09/WKP (USNO/AA):  Moved ht_km calculation from first_entry
                                 block.

   NOTES:
      1. If reference meridian is Greenwich and st=0, 'pos' is
      effectively referred to equator and Greenwich.
      2. This function ignores polar motion, unless the
      observer's longitude and latitude have been corrected for it,
      and variation in the length of day (angular velocity of earth).
      3. The true equator and equinox of date do not form an
      inertial system.  Therefore, with respect to an inertial system,
      the very small velocity component (several meters/day) due to
      the precession and nutation of the Earth's axis is not accounted
      for here.
      4. This function is the C version of NOVAS Fortran routine
      'terra'.

------------------------------------------------------------------------
*/
/*
void terra (OnSurface *location, double st,
            double *pos, double *vel)
*/
func terra(location *OnSurface, st float64, pos, vel []float64) {

	//double df, df2, phi, sinphi, cosphi, c, s, ach, ash, stlocl, sinst,
	//   cosst

	if first_entry != 0 {
		erad_km = ERAD / 1000.0
		first_entry = 0
	}

	// Compute parameters relating to geodetic to geocentric conversion.
	df := 1.0 - F
	df2 := df * df

	phi := location.Latitude * DEG2RAD
	sinphi := math.Sin(phi)
	cosphi := math.Cos(phi)
	c := 1.0 / math.Sqrt(cosphi*cosphi+df2*sinphi*sinphi)
	s := df2 * c
	ht_km := location.Height / 1000.0
	ach := erad_km*c + ht_km
	ash := erad_km*s + ht_km

	// Compute local sidereal time factors at the observer's longitude.
	stlocl := (st*15.0 + location.Longitude) * DEG2RAD
	sinst := math.Sin(stlocl)
	cosst := math.Cos(stlocl)

	// Compute position vector components in kilometers.
	pos[0] = ach * cosphi * cosst
	pos[1] = ach * cosphi * sinst
	pos[2] = ash * sinphi

	// Compute velocity vector components in kilometers/sec.
	vel[0] = -ANGVEL * ach * cosphi * sinst
	vel[1] = ANGVEL * ach * cosphi * cosst
	vel[2] = 0.0

	// Convert position and velocity components to AU and AU/DAY.
	for jdx := 0; jdx < 3; jdx++ {
		pos[jdx] /= AU_KM
		vel[jdx] /= AU_KM
		vel[jdx] *= 86400.0
	}
}
