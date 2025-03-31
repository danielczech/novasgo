/*
  Naval Observatory Vector Astrometry Software (NOVAS)
  C Edition, Version 3.1

  Datatypes for novas

  U. S. Naval Observatory
  Astronomical Applications Dept.
  Washington, DC
  http://www.usno.navy.mil/USNO/astronomical-applications
*/
package novas

/*
   Structures
*/

/*
   struct CatEntry:  basic astrometric data for any celestial object
                      located outside the solar system; the catalog
                      data for a star

   starname[SIZE_OF_OBJ_NAME] = name of celestial object
   catalog[SIZE_OF_CAT_NAME]  = catalog designator (e.g., HIP)
   starnumber                 = integer identifier assigned to object
   ra                         = ICRS right ascension (hours)
   dec                        = ICRS declination (degrees)
   promora                    = ICRS proper motion in right ascension
                                (milliarcseconds/year)
   promodec                   = ICRS proper motion in declination
                                (milliarcseconds/year)
   parallax                   = parallax (milliarcseconds)
   radialvelocity             = radial velocity (km/s)

   SIZE_OF_OBJ_NAME and SIZE_OF_CAT_NAME are defined below.  Each is the
   number of characters in the string (the string length) plus the null
   terminator.
*/

type CatEntry struct {
	//char starname[SIZE_OF_OBJ_NAME];
	Starname string
	// char catalog[SIZE_OF_CAT_NAME];
	Catalog        string
	Starnumber     int64
	Ra             float64
	Dec            float64
	Promora        float64
	Promodec       float64
	Parallax       float64
	Radialvelocity float64
}

/*
   struct object:    specifies the celestial object of interest

   type              = type of object
                     = 0 ... major planet, Pluto, Sun, or Moon
                     = 1 ... minor planet
                     = 2 ... object located outside the solar system
                             (star, nebula, galaxy, etc.)
   number            = object number
                       For 'type' = 0: Mercury = 1, ..., Pluto = 9,
                                       Sun = 10, Moon = 11
                       For 'type' = 1: minor planet number
                       For 'type' = 2: set to 0 (object is
                       fully specified in 'struct CatEntry')
   name              = name of the object (limited to
                       (SIZE_OF_OBJ_NAME - 1) characters)
   star              = basic astrometric data for any celestial object
                       located outside the solar system; the catalog
                       data for a star
*/
type Object struct {
	Type   int16
	Number int16
	//char name[SIZE_OF_OBJ_NAME];
	Name string
	Star CatEntry
}

/*
   struct OnSurface: data for an observer's location on the surface of
                      the Earth.  The atmospheric parameters are used
                      only by the refraction function called from
                      function 'equ2hor'. Additional parameters can be
                      added to this structure if a more sophisticated
                      refraction model is employed.

   latitude           = geodetic (ITRS) latitude; north positive (degrees)
   longitude          = geodetic (ITRS) longitude; east positive (degrees)
   height             = height of the observer (meters)
   temperature        = temperature (degrees Celsius)
   pressure           = atmospheric pressure (millibars)
*/

type OnSurface struct {
	Latitude    float64
	Longitude   float64
	Height      float64
	Temperature float64
	Pressure    float64
}

/*
   struct in_space:   data for an observer's location on a near-Earth
                      spacecraft

   sc_pos[3]          = geocentric position vector (x, y, z), components
                        in km
   sc_vel[3]          = geocentric velocity vector (x_dot, y_dot,
                        z_dot), components in km/s

                        Both vectors with respect to true equator and
                        equinox of date
*/

type InSpace struct {
	Sc_pos [3]float64
	Sc_vel [3]float64
}

/*
   struct observer:   data specifying the location of the observer

   where              = integer code specifying location of observer
                        = 0: observer at geocenter
                        = 1: observer on surface of earth
                        = 2: observer on near-earth spacecraft
   OnSurface         = structure containing data for an observer's
                        location on the surface of the Earth (where = 1)
   near_earth         = data for an observer's location on a near-Earth
                        spacecraft (where = 2)
*/

type Observer struct {
	Where      int16
	On_surf    OnSurface
	Near_earth InSpace
}

/*
   struct sky_pos:    data specifying a celestial object's place on the
                      sky; contains the output from function 'place'

   r_hat[3]           = unit vector toward object (dimensionless)
   ra                 = apparent, topocentric, or astrometric
                        right ascension (hours)
   dec                = apparent, topocentric, or astrometric
                        declination (degrees)
   dis                = true (geometric, Euclidian) distance to solar
                        system body or 0.0 for star (AU)
   rv                 = radial velocity (km/s)
*/

type SkyPos struct {
	R_hat [3]float64
	Ra    float64
	Dec   float64
	Dis   float64
	Rv    float64
}

/*
   struct ra_of_cio:  right ascension of the Celestial Intermediate
                      Origin (CIO) with respect to the GCRS

   jd_tdb             = TDB Julian date
   ra_cio             = right ascension of the CIO with respect
                        to the GCRS (arcseconds)
*/

type ra_of_cio struct {
	Jd_tdb float64
	Ra_cio float64
}

type CioHeader struct {
	Jd_beg float64
	Jd_end float64
	T_int  float64
	N_recs int64
}

type CioData struct {
	T  float64
	Ra float64
}
