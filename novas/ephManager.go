/*
  Naval Observatory Vector Astrometry Software (NOVAS)
  C Edition, Version 3.1

  eph_manager.c: C version of JPL Ephemeris Manager for use with solsys1

  U. S. Naval Observatory
  Astronomical Applications Dept.
  Washington, DC
  http://www.usno.navy.mil/USNO/astronomical-applications

  Golang implementation
*/
package novas

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

// Order matters.
type EphemEpoch struct {
	JDbegin    float64
	JDend      float64
	TimeFactor float64
}

type EphemHeader struct {
	TTL     [252]byte
	Cnam    [2400]byte
	Epoch   EphemEpoch
	Ncon    int32
	Jplau   float64
	EmRatio float64
	IPT     [NUM_TARGETS][3]int32
	DEnum   int32
	LPT     [3]int32
}

const (
	// NOTES:
	// KM...flag defining physical units of the output states.
	// = 1, km and km/sec
	// = 0, AU and AU/day
	// Default value is 0 (KM determines time unit for nutations.
	//                        Angle unit is always radians.)
	EARTH     = 2
	MOON      = 9
	KM        = int16(0)
	SecPerDay = 86400.0
	SEEKSET   = 0 // Seek into file from beginning. Seem man fseek and os.File
	SEEKCUR   = 1 // Seek into file from current position
	SEEKEND   = 2 // Seek into file from end
)

var (
	EphemFilename string

	eh EphemHeader

	//LPT [3]int32 // LPT[3] int[]

	//IPT [12][3]int32 // IPT[3][12]

	NRL           int64
	NP            int64
	NV            int64
	RECORD_LENGTH int64

	PC       [18]float64
	VC       [18]float64
	TWOT     float64
	EM_RATIO float64
	Buffer   []float64
	EPHFILE  *(os.File)
)

func init() {
	var jd_beg, jd_end float64
	var de_num int16
	if error := EphemOpen("JPLEPH", &jd_beg, &jd_end, &de_num); error != 0 {
		fmt.Println("ephManager.go: init(). Error from EphemEpen. error= ", error)
		os.Exit(int(error))
	}
}

// readFloqt64Slice is a helper function to read 4 bytes from file and cast as int
func readFloat64Slice(fp *os.File, n int) ([]float64, error) {
	b := make([]byte, n*8)
	f := make([]float64, n)
	_, err := fp.Read(b)
	if err != nil {
		return nil, err
	}
	for idx := 0; idx < n; idx++ {
		f[idx] = float64(b[idx*n])
	}
	return f, nil
}

// Helper to read binary data into a []float64. endian is of:
// binary.LittleEndian or binary.BigEndian. n is the number of elements in
// the return slice.
func readBinary2Float64Slice(r io.Reader, n int, endian binary.ByteOrder) ([]float64, error) {
	sl := make([]float64, n)

	// binary.LittleEndian
	err := binary.Read(r, endian, &sl)
	if err != nil {
		return nil, err
	}
	return sl, nil
}

// Helper to read binary data into a []int. endian is of:
// binary.LittleEndian or binary.BigEndian. n is the number of elements in
// the return slice.
func readBinary2IntSlice(r io.ByteReader, n int, endian binary.ByteOrder) ([]int, error) {
	sl := make([]int, n)

	for idx, _ := range sl {
		val, err := binary.ReadVarint(r)
		if err != nil {
			return nil, err
		}
		sl[idx] = int(val)
	}
	return sl, nil
}

func test() {
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}
	r := bytes.NewReader(b)

	type D struct {
		PI   float64
		Uate uint8
		Mine [3]byte
		Too  uint16
	}
	var data D

	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	/*
		fmt.Println(data.PI)
		fmt.Println(data.Uate)
		fmt.Printf("% x\n", data.Mine)
		fmt.Println(data.Too)
	*/
}

func ReadBinary2EphemHeader(r io.Reader) error {

	err := binary.Read(r, binary.LittleEndian, &eh)
	if err != nil {
		return err
	}
	return nil
}

func ReadBinary2Buffer(offset int64) error {
	// Open JPL file ephem_name readonly
	EPHFILE, err := os.Open(EphemFilename)
	if err != nil {
		fmt.Println("Error opening JPL file")
		return err
	}
	defer EPHFILE.Close()
	_, serr := EPHFILE.Seek(offset, os.SEEK_SET)
	if serr != nil {
		fmt.Println("Error seeking using offset: ", offset)
		return serr
	}
	//fmt.Println("sought to: ", of)
	err = binary.Read(EPHFILE, binary.LittleEndian, Buffer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// readInt is a helper function to read 4 bytes from file and cast as int
/*
func readInt(fp *os.File) (int, error) {
	b := make([]byte, 4)
	_, err := fp.Read(b)
	if err != nil {
		return 0, err
	}
	return int(b), nil
}
*/

/*
   Define global variables


short int KM;


//   IPT and LPT defined as int to support 64 bit systems.
int IPT[3][12], LPT[3];

long int NRL, NP, NV;
long int RECORD_LENGTH;

double SS[3], JPLAU, PC[18], VC[18], TWOT, EM_RATIO;
double *BUFFER;

FILE *EPHFILE = NULL;
*/

/*
------------------------------------------------------------------------

   PURPOSE:
      This function opens a JPL planetary ephemeris file and
      sets initial values.  This function must be called
      prior to calls to the other JPL ephemeris functions.

   REFERENCES:
      Standish, E.M. and Newhall, X X (1988). "The JPL Export
         Planetary Ephemeris"; JPL document dated 17 June 1988.

   INPUT
   ARGUMENTS:
      *ephem_name (char)
         Name of the direct-access ephemeris file.

   OUTPUT
   ARGUMENTS:
      *jd_begin (double)
         Beginning Julian date of the ephemeris file.
      *jd_end (double)
         Ending Julian date of the ephemeris file.
      *de_number (short int)
         DE number of the ephemeris file opened.

   RETURNED
   VALUE:
      (short int)
          0   ...file exists and is opened correctly.
          1   ...file does not exist/not found.
          2-10...error reading from file header.
          11  ...unable to set record length; ephemeris (DE number)
                 not in look-up table.

   GLOBALS
   USED:
      SS                eph_manager.h
      JPLAU             eph_manager.h
      PC                eph_manager.h
      VC                eph_manager.h
      TWOT              eph_manager.h
      EM_RATIO          eph_manager.h
      BUFFER            eph_manager.h
      IPT               eph_manager.h
      LPT               eph_manager.h
      NRL               eph_manager.h
      KM                eph_manager.h
      NP                eph_manager.h
      NV                eph_manager.h
      RECORD_LENGTH     eph_manager.h
      EPHFILE           eph_manager.h

   FUNCTIONS
   CALLED:
      fclose            stdio.h
      free              stdlib.h
      fopen             stdio.h
      fread             stdio.h
      calloc            stdlib.h

   VER./DATE/
   PROGRAMMER:
      V1.0/06-90/JAB (USNO/NA)
      V1.1/06-92/JAB (USNO/AA): Restructure and add initializations.
      V1.2/07-98/WTH (USNO/AA): Modified to open files for different
                                ephemeris types. (200,403,404,405,406)
      V1.3/11-07/WKP (USNO/AA): Updated prolog.
      V1.4/09-10/WKP (USNO/AA): Changed ncon and denum variables and
                                sizeof ipt array to type 'int' for
                                64-bit system compatibility.
      V1.5/09-10/WTH (USNO/AA): Added support for DE421, default case
                                for switch, close file on error.
      V1.6/10-10/WKP (USNO/AA): Renamed function to lowercase to
                                comply with coding standards.

   NOTES:
      KM...flag defining physical units of the output states.
         = 1, km and km/sec
         = 0, AU and AU/day
      Default value is 0 (KM determines time unit for nutations.
                          Angle unit is always radians.)

------------------------------------------------------------------------
*/
func EphemOpen(ephem_name string,
	jd_begin *float64, jd_end *float64,
	de_number *int16) int16 {

	// save for later when reading Buffer contents
	EphemFilename = ephem_name

	var i int16

	if EPHFILE != nil {
		EPHFILE.Close()
	}

	// Open JPL file ephem_name readonly
	EPHFILE, err := os.Open(ephem_name)
	if err != nil {
		return 1 // remove magic number, use a constant.
	}
	defer EPHFILE.Close()

	//	 File found. Set initializations and default values.

	NRL = 0

	NP = 2
	NV = 3
	TWOT = 0.0

	//PC := make([]float64, 18)
	//VC := make([]float64, 18)

	for idx := 0; i < 18; i++ {
		PC[idx] = 0.0
		VC[idx] = 0.0
	}

	PC[0] = 1.0
	VC[1] = 1.0

	//   Read in values from the first record, aka the header.
	ReadBinary2EphemHeader(EPHFILE)

	//   Set the value of the record length according to what JPL ephemeris is
	//   being opened.

	//fmt.Println("DEnum: ", eh.DEnum)
	switch eh.DEnum {
	case 200:
		RECORD_LENGTH = 6608
	case 403, 405, 421:
		RECORD_LENGTH = 8144
	case 404, 406:
		RECORD_LENGTH = 5824

		//   An unknown DE file was opened. Close the file and return an error
		//   code.
	default:
		*jd_begin = 0.0
		*jd_end = 0.0
		*de_number = 0
		return 11
	}
	//fmt.Println("RECORD_LENGTH= ", RECORD_LENGTH)

	//BUFFER = (double *) calloc (RECORD_LENGTH / 8, sizeof(double))
	Buffer = make([]float64, RECORD_LENGTH/8)
	//fmt.Println("sizeof Buffer: ", len(Buffer))

	*de_number = int16(eh.DEnum)
	*jd_begin = eh.Epoch.JDbegin
	*jd_end = eh.Epoch.JDend

	return 0
}

/*
------------------------------------------------------------------------

   PURPOSE:
      This function closes a JPL planetary ephemeris file and frees the
      memory.

   REFERENCES:
      Standish, E.M. and Newhall, X X (1988). "The JPL Export
         Planetary Ephemeris"; JPL document dated 17 June 1988.

   INPUT
   ARGUMENTS:
      None.

   OUTPUT
   ARGUMENTS:
      None.

   RETURNED
   VALUE:
      (short int)
          0  ...file was already closed or closed correctly.
          EOF...error closing file.

   GLOBALS
   USED:
      BUFFER            eph_manager.h
      EPHFILE           eph_manager.h

   FUNCTIONS
   CALLED:
      fclose            stdio.h
      free              stdlib.h

   VER./DATE/
   PROGRAMMER:
      V1.0/11-07/WKP (USNO/AA)
      V1.1/09-10/WKP (USNO/AA): Explicitly cast fclose return value to
                                type 'short int'.
      V1.2/10-10/WKP (USNO/AA): Renamed function to lowercase to
                                comply with coding standards.

   NOTES:
      None.

------------------------------------------------------------------------
*/
/*
func EphemClose () int16 {

	var error int16 // short int error = 0

   if EPHFILE {

      error := EPHFILE.Close()
		 if error != nil {
			 return 999
     }
	 }
   return  0
}
*/

/*
------------------------------------------------------------------------

   PURPOSE:
      This function accesses the JPL planetary ephemeris to give the
      position and velocity of the target object with respect to the
      center object.

   REFERENCES:
      Standish, E.M. and Newhall, X X (1988). "The JPL Export
         Planetary Ephemeris"; JPL document dated 17 June 1988.

   INPUT
   ARGUMENTS:
      tjd[2] (double)
         Two-element array containing the Julian date, which may be
         split any way (although the first element is usually the
         "integer" part, and the second element is the "fractional"
         part).  Julian date is in the TDB or "T_eph" time scale.
      target (short int)
         Number of 'target' point.
      center (short int)
         Number of 'center' (origin) point.
         The numbering convention for 'target' and'center' is:
            0 = Mercury           7 = Neptune
            1 = Venus             8 = Pluto
            2 = Earth             9 = Moon
            3 = Mars             10 = Sun
            4 = Jupiter          11 = Solar system bary.
            5 = Saturn           12 = Earth-Moon bary.
            6 = Uranus           13 = Nutations (long int. and obliq.)
            (If nutations are desired, set 'target' = 13;
             'center' will be ignored on that call.)

   OUTPUT
   ARGUMENTS:
      *position (double)
         Position vector array of target relative to center, measured
         in AU.
      *velocity (double)
         Velocity vector array of target relative to center, measured
         in AU/day.

   RETURNED
   VALUE:
      (short int)
         0  ...everything OK.
         1,2...error returned from State.

   GLOBALS
   USED:
      EM_RATIO          eph_manager.h

   FUNCTIONS
   CALLED:
      state             eph_manager.h

   VER./DATE/
   PROGRAMMER:
      V1.0/03-93/WTH (USNO/AA): Convert FORTRAN to C.
      V1.1/07-93/WTH (USNO/AA): Update to C standards.
      V2.0/07-98/WTH (USNO/AA): Modified for ease of use and linearity.
      V3.0/11-06/JAB (USNO/AA): Allowed for use of input 'split' Julian
                                date for higher precision.
      V3.1/11-07/WKP (USNO/AA): Updated prolog and error codes.
      V3.1/12-07/WKP (USNO/AA): Removed unreferenced variables.
      V3.2/10-10/WKP (USNO/AA): Renamed function to lowercase to
                                comply with coding standards.

   NOTES:
      None.

------------------------------------------------------------------------
*/
/*
short int planet_ephemeris (double tjd[2], short int target,
                            short int center,

                            double *position, double *velocity)
{
*/
func PlanetEphemeris(tjd [2]float64, target, center int16,
	position, velocity []float64) int16 {

	/*
		fmt.Println("Inside PlanetEphemeris")
		fmt.Println("   tjd= ", tjd)
		fmt.Println("   target= ", target)
		fmt.Println("   center= ", center)
		fmt.Println("PlanetEphemeris: eh.Epoch= ", eh.Epoch)
	*/
	var error int16
	earth := int16(EARTH)
	moon := int16(MOON)

	//short int i, error = 0, earth = 2, moon = 9
	//short int do_earth = 0, do_moon = 0
	var do_earth, do_moon int16

	/*
			 double jed[2]
		   double pos_moon[3] = {0.0,0.0,0.0}, vel_moon[3] = {0.0,0.0,0.0},
		          pos_earth[3] = {0.0,0.0,0.0}, vel_earth[3] = {0.0,0.0,0.0}
		   double target_pos[3] = {0.0,0.0,0.0}, target_vel[3] = {0.0,0.0,0.0},
		          center_pos[3] = {0.0,0.0,0.0}, center_vel[3] = {0.0,0.0,0.0}
	*/
	jed := make([]float64, 2)
	pos_moon := make([]float64, 3)
	vel_moon := make([]float64, 3)
	pos_earth := make([]float64, 3)
	vel_earth := make([]float64, 3)
	target_pos := make([]float64, 3)
	target_vel := make([]float64, 3)
	center_pos := make([]float64, 3)
	center_vel := make([]float64, 3)

	//   Initialize 'jed' for 'state' and set up component count.

	jed[0] = tjd[0]
	jed[1] = tjd[1]

	//   Check for target point = center point.

	if target == center {
		for idx := 0; idx < 3; idx++ {
			position[idx] = 0.0
			velocity[idx] = 0.0
		}
		return 0
	}

	//   Check for instances of target or center being Earth or Moon,
	//   and for target or center being the Earth-Moon barycenter.

	if (target == earth) || (center == earth) {
		do_moon = 1
	}
	if (target == moon) || (center == moon) {
		do_earth = 1
	}
	if (target == 12) || (center == 12) {
		do_earth = 1
	}

	if do_earth == 1 {
		error := state(jed, 2, pos_earth, vel_earth)
		if error != 0 {
			fmt.Println("PlanetEphemeris: do_earth: error calling state(). err= ", error)
			return error
		}
	}

	if do_moon == 1 {
		error := state(jed, 9, pos_moon, vel_moon)
		if error != 0 {
			fmt.Println("PlanetEphemeris: do_moon: error calling state(). err= ", error)
			return error
		}
	}

	//   Make call to State for target object.

	if target == 11 {
		for idx := 0; idx < 3; idx++ {
			target_pos[idx] = 0.0
			target_vel[idx] = 0.0
		}
	} else if target == 12 {
		for idx := 0; idx < 3; idx++ {
			target_pos[idx] = pos_earth[idx]
			target_vel[idx] = vel_earth[idx]
		}
	} else {
		error = state(jed, target, target_pos, target_vel)
		if error != 0 {
			fmt.Println("PlanetEphemeris: target: error calling state(). err= ", error)
			fmt.Println("  called state with target= ", target)
			return error
		}
	}

	//   Make call to State for center object.

	//   If the requested center is the Solar System barycenter,
	//   then don't bother with a second call to State.

	if center == 11 {
		for idx := 0; idx < 3; idx++ {
			center_pos[idx] = 0.0
			center_vel[idx] = 0.0
		}
	} else if center == 12 {

		//   Center is Earth-Moon barycenter, which was already computed above.

		for idx := 0; idx < 3; idx++ {
			center_pos[idx] = pos_earth[idx]
			center_vel[idx] = vel_earth[idx]
		}
	} else {
		if error = state(jed, center, center_pos, center_vel); error != 0 {
			fmt.Println("PlanetEphemeris: center obj: error calling state(). err= ", error)
			return error
		}
	}

	//   Check for cases of Earth as target and Moon as center or vice versa.

	if (target == earth) && (center == moon) {
		for idx := 0; idx < 3; idx++ {
			position[idx] = -center_pos[idx]
			velocity[idx] = -center_vel[idx]
		}
		return 0
	} else if (target == moon) && (center == earth) {
		for idx := 0; idx < 3; idx++ {
			position[idx] = target_pos[idx]
			velocity[idx] = target_vel[idx]
		}
		return 0
	} else if target == earth {

		//   Check for Earth as target, or as center.

		for idx := 0; idx < 3; idx++ {
			target_pos[idx] = target_pos[idx] - (pos_moon[idx] /
				(1.0 + eh.EmRatio))
			target_vel[idx] = target_vel[idx] - (vel_moon[idx] /
				(1.0 + eh.EmRatio))
		}
	} else if center == earth {
		for idx := 0; idx < 3; idx++ {
			center_pos[idx] = center_pos[idx] - (pos_moon[idx] /
				(1.0 + eh.EmRatio))
			center_vel[idx] = center_vel[idx] - (vel_moon[idx] /
				(1.0 + eh.EmRatio))
		}
	} else if target == moon {

		//   Check for Moon as target, or as center.

		for idx := 0; idx < 3; idx++ {
			target_pos[idx] = (pos_earth[idx] - (target_pos[idx] /
				(1.0 + eh.EmRatio))) + target_pos[idx]
			target_vel[idx] = (vel_earth[idx] - (target_vel[idx] /
				(1.0 + eh.EmRatio))) + target_vel[idx]
		}
	} else if center == moon {
		for idx := 0; idx < 3; idx++ {
			center_pos[idx] = (pos_earth[idx] - (center_pos[idx] /
				(1.0 + eh.EmRatio))) + center_pos[idx]
			center_vel[idx] = (vel_earth[idx] - (center_vel[idx] /
				(1.0 + eh.EmRatio))) + center_vel[idx]
		}
	}

	//   Compute position and velocity vectors.

	for idx := 0; idx < 3; idx++ {
		position[idx] = target_pos[idx] - center_pos[idx]
		velocity[idx] = target_vel[idx] - center_vel[idx]
	}

	return 0
}

/*
------------------------------------------------------------------------

   PURPOSE:
      This function reads and interpolates the JPL planetary
      ephemeris file.

   REFERENCES:
      Standish, E.M. and Newhall, X X (1988). "The JPL Export
         Planetary Ephemeris"; JPL document dated 17 June 1988.

   INPUT
   ARGUMENTS:
      *jed (double)
         2-element Julian date (TDB) at which interpolation is wanted.
         Any combination of jed[0]+jed[1] which falls within the time
         span on the file is a permissible epoch.  See Note 1 below.
      target (short int)
         The requested body to get data for from the ephemeris file.
         The designation of the astronomical bodies is:
                 0 = Mercury                    6 = Uranus
                 1 = Venus                      7 = Neptune
                 2 = Earth-Moon barycenter      8 = Pluto
                 3 = Mars                       9 = geocentric Moon
                 4 = Jupiter                   10 = Sun
                 5 = Saturn

   OUTPUT
   ARGUMENTS:
      *target_pos (double)
         The barycentric position vector array of the requested object,
         in AU.
         (If target object is the Moon, then the vector is geocentric.)
      *target_vel (double)
         The barycentric velocity vector array of the requested object,
         in AU/Day.

         Both vectors are referenced to the Earth mean equator and
         equinox of epoch.

   RETURNED
   VALUE:
      (short int)
         0...everything OK.
         1...error reading ephemeris file.
         2...epoch out of range.

   GLOBALS
   USED:
      KM                eph_manager.h
      EPHFILE           eph_manager.h
      IPT               eph_manager.h
      BUFFER            eph_manager.h
      NRL               eph_manager.h
      RECORD_LENGTH     eph_manager.h
      SS                eph_manager.h
      JPLAU             eph_manager.h

   FUNCTIONS
   CALLED:
      split             eph_manager.h
      fseek             stdio.h
      fread             stdio.h
      interpolate       eph_manager.h
      ephem_close       eph_manager.h

   VER./DATE/
   PROGRAMMER:
      V1.0/03-93/WTH (USNO/AA): Convert FORTRAN to C.
      V1.1/07-93/WTH (USNO/AA): Update to C standards.
      V2.0/07-98/WTH (USNO/AA): Modify to make position and velocity
                                two distinct vector arrays.  Routine set
                                to compute one state per call.
      V2.1/11-07/WKP (USNO/AA): Updated prolog.
      V2.2/10-10/WKP (USNO/AA): Renamed function to lowercase to
                                comply with coding standards.

   NOTES:
      1. For ease in programming, the user may put the entire epoch in
         jed[0] and set jed[1] = 0. For maximum interpolation accuracy,
         set jed[0] = the most recent midnight at or before
         interpolation epoch, and set jed[1] = fractional part of a day
         elapsed between jed[0] and epoch. As an alternative, it may
         prove convenient to set jed[0] = some fixed epoch, such as
         start of the integration and jed[1] = elapsed interval between
         then and epoch.

------------------------------------------------------------------------
*/
func state(jed []float64, target int16, target_pos, target_vel []float64) int16 {

	// long int nr, rec
	var nr, rec int64

	//double t[2], aufac = 1.0, jd[4], s

	t := make([]float64, 2)
	aufac := float64(1.0)
	jd := make([]float64, 4)

	//   Set units based on value of the 'KM' flag.

	if KM != 0 {
		t[1] = eh.Epoch.TimeFactor * SecPerDay
	} else {
		t[1] = eh.Epoch.TimeFactor
		aufac = 1.0 / eh.Jplau
	}

	//   Check epoch.

	s := jed[0] - 0.5
	split(s, jd)
	split(jed[1], jd[2:])
	jd[0] += jd[2] + 0.5
	jd[1] += jd[3]
	split(jd[1], jd[2:])
	jd[0] += jd[2]

	//   Return error code if date is out of range.

	if (jd[0] < eh.Epoch.JDbegin) || ((jd[0] + jd[3]) > eh.Epoch.JDend) {
		fmt.Println("state: epoc out of range. jd[0]= ", jd[0], ", jd[3]= ", jd[3], " begin= ", eh.Epoch.JDbegin, " end= ", eh.Epoch.JDend)
		return 2
	}

	//   Calculate record number and relative time interval.

	nr = int64((jd[0]-eh.Epoch.JDbegin)/eh.Epoch.TimeFactor) + 3
	if jd[0] == eh.Epoch.JDend {
		nr -= 2
	}
	t[0] = ((jd[0] - (float64((nr-3))*eh.Epoch.TimeFactor + eh.Epoch.JDbegin)) + jd[3]) / eh.Epoch.TimeFactor

	//   Read correct record if it is not already in memory.
	//fmt.Println("nr: ", nr)
	if nr != NRL {
		NRL = nr
		rec = (nr - 1) * RECORD_LENGTH // offset seek
		err := ReadBinary2Buffer(rec)
		if err != nil {
			fmt.Println("Error reading buffer: ", err.Error())
			EPHFILE.Close()
			return 1
		}
	}

	//   Check and interpolate for requested body.
	/*
		fmt.Println("calling interpolate with: ")
		fmt.Println("   target: ", target)
		fmt.Println("   eh.IPT[target]: ", eh.IPT[target])
		fmt.Println("   eh.IPT[target][0]-1: ", eh.IPT[target][0]-1)
		fmt.Println("   Buffer= ", Buffer[eh.IPT[target][0]-1])
		fmt.Println("   t= ", t)
		fmt.Println("   eh.IPT[target][1]= ", eh.IPT[target][1])
		fmt.Println("   eh.IPT[target][2]= ", eh.IPT[target][2])
	*/
	interpolate(Buffer[eh.IPT[target][0]-1:], t, int64(eh.IPT[target][1]),
		int64(eh.IPT[target][2]), target_pos, target_vel)

	for idx := 0; idx < 3; idx++ {
		target_pos[idx] *= aufac
		target_vel[idx] *= aufac
		//fmt.Println("   target_pos[", idx, "]= ", target_pos[idx])
		//fmt.Println("   target_vel[", idx, "]= ", target_vel[idx])
	}
	return 0
}

/*
------------------------------------------------------------------------

   PURPOSE:
      This function breaks up a double number into a double integer
      part and a fractional part.

   REFERENCES:
      Standish, E.M. and Newhall, X X (1988). "The JPL Export
         Planetary Ephemeris"; JPL document dated 17 June 1988.

   INPUT
   ARGUMENTS:
      tt (double)
         Input number.

   OUTPUT
   ARGUMENTS:
      *fr (double)
         2-element output array;
            fr[0] contains integer part,
            fr[1] contains fractional part.
         For negative input numbers,
            fr[0] contains the next more negative integer;
            fr[1] contains a positive fraction.

   RETURNED
   VALUE:
      None.

   GLOBALS
   USED:
      None.

   FUNCTIONS
   CALLED:
      None.

   VER./DATE/
   PROGRAMMER:
      V1.0/06-90/JAB (USNO/NA): CA coding standards
      V1.1/03-93/WTH (USNO/AA): Convert to C.
      V1.2/07-93/WTH (USNO/AA): Update to C standards.
      V1.3/10-10/WKP (USNO/AA): Renamed function to lowercase to
                                comply with coding standards.

   NOTES:
      None.

------------------------------------------------------------------------
*/
func split(tt float64, fr []float64) {

	//   Get integer and fractional parts.

	fr[0] = float64((int64(tt)))
	fr[1] = tt - fr[0]

	//   Make adjustments for negative input number.

	if (tt >= 0.0) || (fr[1] == 0.0) {
		return
	} else {
		fr[0] -= 1.0
		fr[1] += 1.0
	}
}

/*
------------------------------------------------------------------------

   PURPOSE:
      This function differentiates and interpolates a set of
      Chebyshev coefficients to give position and velocity.

   REFERENCES:
      Standish, E.M. and Newhall, X X (1988). "The JPL Export
         Planetary Ephemeris"; JPL document dated 17 June 1988.

   INPUT
   ARGUMENTS:
      *buf (double)
         Array of Chebyshev coefficients of position.
      *t (double)
         t[0] is fractional time interval covered by coefficients at
         which interpolation is desired (0 <= t[0] <= 1).
         t[1] is length of whole interval in input time units.
      ncf (long int)
         Number of coefficients per component.
      na (long int)
         Number of sets of coefficients in full array
         (i.e., number of sub-intervals in full interval).

   OUTPUT
   ARGUMENTS:
      *position (double)
         Position array of requested object.
      *velocity (double)
         Velocity array of requested object.

   RETURNED
   VALUE:
      None.

   GLOBALS
   USED:
      NP                eph_manager.h
      NV                eph_manager.h
      PC                eph_manager.h
      VC                eph_manager.h
      TWOT              eph_manager.h

   FUNCTIONS
   CALLED:
      fmod              math.h

   VER./DATE/
   PROGRAMMER:
      V1.0/03-93/WTH (USNO/AA): Convert FORTRAN to C.
      V1.1/07-93/WTH (USNO/AA): Update to C standards.
      V1.2/07-98/WTH (USNO/AA): Modify to make position and velocity
                                two distinct vector arrays.
      V1.3/11-07/WKP (USNO/AA): Updated prolog.
      V1.4/12-07/WKP (USNO/AA): Changed ncf and na arguments from short
                                int to long int.
      V1.5/10-10/WKP (USNO/AA): Renamed function to lowercase to
                                comply with coding standards.

   NOTES:
      None.

------------------------------------------------------------------------
*/
func interpolate(buf, t []float64, ncf, na int64,
	position, velocity []float64) {

	var dna, dt1, temp, tc, vfac float64

	/*
	   Get correct sub-interval number for this set of coefficients and
	   then get normalized Chebyshev time within that subinterval.
	*/

	dna = float64(na)
	dt1 = float64(int64(t[0]))
	temp = dna * t[0]
	l := int64(temp - dt1)

	/*
	   'tc' is the normalized Chebyshev time (-1 <= tc <= 1).
	*/

	tc = 2.0*(math.Mod(temp, 1.0)+dt1) - 1.0
	/*
		fmt.Println("temp= ", temp)
		fmt.Println("math.Mod(temp, 1)= ", math.Mod(temp, 1.0))
		fmt.Println("dt1= ", dt1)
		fmt.Println("tc= ", tc)
	*/
	/*
	   Check to see whether Chebyshev time has changed, and compute new
	   polynomial values if it has.  (The element PC[1] is the value of
	   t1[tc] and hence contains the value of 'tc' on the previous call.)
	*/

	if tc != PC[1] {
		NP = 2
		NV = 3
		PC[1] = tc
		TWOT = tc + tc
	}

	/*
	   Be sure that at least 'ncf' polynomials have been evaluated and
	   are stored in the array 'PC'.
	*/

	if NP < ncf {
		for idx := NP; idx < ncf; idx++ {
			PC[idx] = TWOT*PC[idx-1] - PC[idx-2]
		}
		NP = ncf
	}

	/*
	   Interpolate to get position for each component.
	*/

	for idx := 0; idx < 3; idx++ {
		position[idx] = 0.0
		//fmt.Println("i= ", i)
		for jdx := ncf - 1; jdx >= 0; jdx-- {
			//fmt.Println("  jdx= ", jdx)
			kdx := jdx + (int64(idx) * ncf) + (l * (3 * ncf))
			/*
				fmt.Println("  ncf= ", ncf)
				fmt.Println("  l= ", l)
				fmt.Println("  kdx= ", kdx)
				fmt.Println("  buf[kdx]= ", buf[kdx])
				fmt.Println("  PC[jdx]= ", PC[jdx])
			*/
			position[idx] += PC[jdx] * buf[kdx]
			//fmt.Println("  position[idx]= ", position[idx])
		}
	}

	/*
	   If velocity interpolation is desired, be sure enough derivative
	   polynomials have been generated and stored.
	*/

	vfac = (2.0 * dna) / t[1]
	VC[2] = 2.0 * TWOT
	if NV < ncf {
		for idx := NV; idx < ncf; idx++ {
			VC[idx] = TWOT*VC[idx-1] + PC[idx-1] + PC[idx-1] - VC[idx-2]
		}
		NV = ncf
	}

	/*
	   Interpolate to get velocity for each component.
	*/

	for idx := 0; idx < 3; idx++ {
		velocity[idx] = 0.0
		for jdx := ncf - 1; jdx > 0; jdx-- {
			kdx := jdx + (int64(idx) * ncf) + (l * (3 * ncf))
			velocity[idx] += VC[jdx] * buf[kdx]
		}
		velocity[idx] *= vfac
	}
}
