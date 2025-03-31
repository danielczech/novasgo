// ********spin
package novas

import "math"

var (
	//static double ang_last = -999.0
	//static double xx, yx, zx, xy, yy, zy, xz, yz, zz
	ang_last                           = -999.0
	xx, yx, zx, xy, yy, zy, xz, yz, zz float64
)

/*
------------------------------------------------------------------------

   PURPOSE:
      This function transforms a vector from one coordinate system
      to another with same origin and axes rotated about the z-axis.

   REFERENCES:
      Kaplan, G. H. et. al. (1989). Astron. Journ. 97, 1197-1210.

   INPUT
   ARGUMENTS:
      angle (double)
         Angle of coordinate system rotation, positive counterclockwise
         when viewed from +z, in degrees.
      pos1[3] (double)
         Position vector.

   OUTPUT
   ARGUMENTS:
      pos2[3] (double)
         Position vector expressed in new coordinate system rotated
         about z by 'angle'.

   RETURNED
   VALUE:
      None.

   GLOBALS
   USED:
      DEG2RAD            novascon.c

   FUNCTIONS
   CALLED:
      sin                math.h
      cos                math.h

   VER./DATE/
   PROGRAMMER:
      V1.0/08-93/WTH (USNO/AA) Translate Fortran.
      V2.0/10-03/JAB (USNO/AA) Update for IAU 2000 resolutions.
      V2.1/01-05/JAB (USNO/AA) Generalize the function.

   NOTES:
      1. This function is the C version of NOVAS Fortran routine 'spin'.

------------------------------------------------------------------------
*/
/*
void spin (double angle, double *pos1,
           double *pos2)
*/
func spin(angle float64, pos1, pos2 []float64) {

	//static double ang_last = -999.0
	//static double xx, yx, zx, xy, yy, zy, xz, yz, zz
	//double angr, cosang, sinang
	var angr, cosang, sinang float64

	if math.Abs(angle-ang_last) >= 1.0e-12 {
		angr = angle * DEG2RAD
		cosang = math.Cos(angr)
		sinang = math.Sin(angr)

		// Rotation matrix follows.
		xx = cosang
		yx = sinang
		zx = 0.0
		xy = -sinang
		yy = cosang
		zy = 0.0
		xz = 0.0
		yz = 0.0
		zz = 1.0

		ang_last = angle
	}

	// Perform rotation.
	pos2[0] = xx*pos1[0] + yx*pos1[1] + zx*pos1[2]
	pos2[1] = xy*pos1[0] + yy*pos1[1] + zy*pos1[2]
	pos2[2] = xz*pos1[0] + yz*pos1[1] + zz*pos1[2]
}
