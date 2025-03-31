package novas

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
)

const (
	TEST_ABS_ERR = 1.e-1
)

/*
   Function test status

      MakeCatEntry     novas.c PASS
      MakeObject        novas.c PASS
      tdb2tt             novas.c PASS
      ephemeris          novas.c PASS via light_time
      geo_posvel         novas.c PASS
      starvectors        novas.c PASS
      d_light            novas.c PASS
      proper_motion      novas.c PASS
      bary2obs           novas.c PASS
      light_time         novas.c TODO
      limb_angle         novas.c PASS
      grav_def           novas.c PASS
      aberration         novas.c PASS
      frame_tie          novas.c PASS
      precession         novas.c PASS
      nutation           novas.c PASS
      cio_location       novas.c
      cio_basis          novas.c
      rad_vel            novas.c PASS
      vector2radec       novas.c PASS
      nutation_angles    novas.c PASS
      ee_ct              novas.c PASS
      mean_obliq         novac.c PASS
      fund_args          novas.c PASS
      Etilt             novas.c PASS
*/

func checkString(s, es string) error {
	if s != es {
		fmt.Println("Got: ", s, " expected: ", es)
		return errors.New("")
	}
	return nil
}

func checkFloat(x, expX, tol float64) error {
	if math.Abs(x-expX) > tol {
		fmt.Println("Got: ", x, " expected: ", expX)
		return errors.New("")
	}
	return nil
}

func checkStar(star, expCatEntry CatEntry) error {
	if err := checkString(star.Starname, expCatEntry.Starname); err != nil {
		return errors.New("")
	}
	if err := checkString(star.Catalog, expCatEntry.Catalog); err != nil {
		return errors.New("")
	}
	if err := checkFloat(float64(star.Starnumber), float64(expCatEntry.Starnumber), ABS_ERR); err != nil {
		return errors.New("")
	}
	if err := checkFloat(star.Ra, expCatEntry.Ra, ABS_ERR); err != nil {
		return errors.New("")
	}
	if err := checkFloat(star.Dec, expCatEntry.Dec, ABS_ERR); err != nil {
		return errors.New("")
	}
	if err := checkFloat(star.Promora, expCatEntry.Promora, ABS_ERR); err != nil {
		return errors.New("")
	}
	if err := checkFloat(star.Promodec, expCatEntry.Promodec, ABS_ERR); err != nil {
		return errors.New("")
	}
	if err := checkFloat(star.Parallax, expCatEntry.Parallax, ABS_ERR); err != nil {
		return errors.New("")
	}
	if err := checkFloat(star.Radialvelocity, expCatEntry.Radialvelocity, ABS_ERR); err != nil {
		return errors.New("")
	}
	return nil
}

func TestMakeCatEntry(t *testing.T) {
	var star, expCatEntry CatEntry

	starName := "Star Name"
	catalog := "CAT"
	starNumber := int64(1)
	ra := 1.0
	dec := 2.0
	pm_ra := 3.0
	pm_dec := 4.0
	parallax := 5.0
	rad_vel := 6.0

	expCatEntry.Starname = starName
	expCatEntry.Catalog = catalog
	expCatEntry.Starnumber = starNumber
	expCatEntry.Ra = ra
	expCatEntry.Dec = dec
	expCatEntry.Promora = pm_ra
	expCatEntry.Promodec = pm_dec
	expCatEntry.Parallax = parallax
	expCatEntry.Radialvelocity = rad_vel

	if ec := MakeCatEntry(starName, catalog, starNumber, ra, dec, pm_ra, pm_dec, parallax, rad_vel, &star); ec != 0 {
		fmt.Println("Got errorCode: ", ec, " expected 0")
		t.Fail()
	}
	if ec := checkStar(star, expCatEntry); ec != nil {
		t.Fail()
	}
}

func TestMakeObject(t *testing.T) {
	objType := int16(2)
	number := int16(0) // for objType = 2
	name := "obn"
	var star CatEntry
	var cob Object

	starName := "Star"
	catalog := "CAT"
	starNumber := int64(1)
	ra := 1.0
	dec := 2.0
	pm_ra := 3.0
	pm_dec := 4.0
	parallax := 5.0
	rad_vel := 6.0

	if ec := MakeCatEntry(starName, catalog, starNumber, ra, dec, pm_ra, pm_dec, parallax, rad_vel, &star); ec != 0 {
		fmt.Println("Got errorCode: ", ec, " expected 0")
		t.Fail()
	}

	if ec := MakeObject(objType, number, name, &star, &cob); ec != 0 {
		fmt.Println("Got errorCode: ", ec, " expected 0")
		t.Fail()
	}
	if err := checkFloat(float64(cob.Type), float64(objType), ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(float64(cob.Number), float64(number), ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkString(cob.Name, strings.ToUpper(name)); err != nil {
		t.Fail()
	}
	if ec := checkStar(cob.Star, star); ec != nil {
		t.Fail()
	}
}

func TestMakeTdb2tt(t *testing.T) {

	var tdb_jd, tt_jd, secdiff float64
	var expTT_jd, expSecdiff float64

	tdb_jd = 2450300.500000000000
	expTT_jd = 2450300.500000009779
	expSecdiff = -0.000851349649
	tdb2tt(tdb_jd, &tt_jd, &secdiff)
	if err := checkFloat(tt_jd, expTT_jd, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(secdiff, expSecdiff, ABS_ERR); err != nil {
		t.Fail()
	}

	// 2nd Test
	tdb_jd = 2450417.500000000000
	expTT_jd = 2450417.500000011176
	expSecdiff = -0.000961169111
	tdb2tt(tdb_jd, &tt_jd, &secdiff)
	if err := checkFloat(tt_jd, expTT_jd, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(secdiff, expSecdiff, ABS_ERR); err != nil {
		t.Fail()
	}
}

/* TODO: Write Test
func TestEphemeris(t *testing.T) {
	fmt.Println("FAIL Ephemeris - TODO: Write Test")
	t.Fail()
}
*/

func TestGeoPosVel(t *testing.T) {

	var jd_tt, delta_t float64
	var accuracy int16
	var location Observer
	pos := make([]float64, 3)
	vel := make([]float64, 3)
	expPos := make([]float64, 3)
	expVel := make([]float64, 3)

	jd_tt = 2450300.500000000000
	delta_t = 60.000000000000
	accuracy = 0
	location.Where = 1
	location.On_surf.Latitude = 45.000000000000
	location.On_surf.Longitude = -75.000000000000
	location.On_surf.Height = 0.000000000000
	location.On_surf.Temperature = 10.000000000000
	location.On_surf.Pressure = 1010.000000000000
	location.Near_earth.Sc_pos[0] = 0.000000000000
	location.Near_earth.Sc_pos[1] = 0.000000000000
	location.Near_earth.Sc_pos[2] = 0.000000000000
	location.Near_earth.Sc_vel[0] = 0.000000000000
	location.Near_earth.Sc_vel[0] = 0.000000000000
	location.Near_earth.Sc_vel[0] = 0.000000000000
	expPos[0] = -0.000015736292
	expPos[1] = -0.000025781292
	expPos[2] = 0.000029989852
	expVel[0] = 0.000162423733
	expVel[1] = -0.000099083801
	expVel[2] = 0.000000047979

	if error := geo_posvel(jd_tt, delta_t, accuracy, &location, pos, vel); error != 0 {
		fmt.Println("Got error: ", error, " expected 0")
		t.Fail()
	}
	for idx, p := range pos {
		if err := checkFloat(p, expPos[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
	for idx, v := range vel {
		if err := checkFloat(v, expVel[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}

	// 2nd Test

	jd_tt = 2450417.500000000000
	delta_t = 60.000000000000
	accuracy = 0
	location.Where = 1
	location.On_surf.Latitude = 45.000000000000
	location.On_surf.Longitude = -75.000000000000
	location.On_surf.Height = 0.000000000000
	location.On_surf.Temperature = 10.000000000000
	location.On_surf.Pressure = 1010.000000000000
	location.Near_earth.Sc_pos[0] = 0.000000000000
	location.Near_earth.Sc_pos[1] = 0.000000000000
	location.Near_earth.Sc_pos[2] = 0.000000000000
	location.Near_earth.Sc_vel[0] = 0.000000000000
	location.Near_earth.Sc_vel[0] = 0.000000000000
	location.Near_earth.Sc_vel[0] = 0.000000000000
	expPos[0] = 0.000030020071
	expPos[1] = -0.000003193591
	expPos[2] = 0.000030004890
	expVel[0] = 0.000020111931
	expVel[1] = 0.000189194558
	expVel[2] = 0.000000014943
	if error := geo_posvel(jd_tt, delta_t, accuracy, &location, pos, vel); error != 0 {
		fmt.Println("Got error: ", error, " expected 0")
		t.Fail()
	}
	for idx, p := range pos {
		if err := checkFloat(p, expPos[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
	for idx, v := range vel {
		if err := checkFloat(v, expVel[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}

func TestStarVectors(t *testing.T) {

	var star CatEntry
	pos := make([]float64, 3)
	vel := make([]float64, 3)
	expPos := make([]float64, 3)
	expVel := make([]float64, 3)

	star.Starname = "Theta CAR"
	star.Catalog = "HIP"
	star.Starnumber = 2
	star.Ra = 10.715944806000
	star.Dec = -64.394450000000
	star.Promora = -18.870000000000
	star.Promodec = 12.060000000000
	star.Parallax = 7.430000000000
	star.Radialvelocity = 24.000000000000
	expPos[0] = -11326046.253463707864
	expPos[1] = 3957634.030029209331
	expPos[2] = -25034680.670301370323
	expVel[0] = -0.007145191487
	expVel[1] = 0.009862925083
	expVel[2] = -0.010580159896

	starvectors(&star, pos, vel)

	for idx, p := range pos {
		if err := checkFloat(p, expPos[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
	for idx, v := range vel {
		if err := checkFloat(v, expVel[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}

	// 2nd test
	star.Starname = "Delta ORI"
	star.Catalog = "HIP"
	star.Starnumber = 1
	star.Ra = 5.533444639000
	star.Dec = -0.299091944000
	star.Promora = 1.670000000000
	star.Promodec = 0.560000000000
	star.Parallax = 3.560000000000
	star.Radialvelocity = 16.000000000000
	expPos[0] = 7059283.251987532713
	expPos[1] = 57507101.873731084168
	expPos[2] = -302451.154011160368
	expVel[0] = -0.000148609903
	expVel[1] = 0.009331014346
	expVel[2] = 0.000382449798
	starvectors(&star, pos, vel)

	for idx, p := range pos {
		if err := checkFloat(p, expPos[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
	for idx, v := range vel {
		if err := checkFloat(v, expVel[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}

}

func TestDlight(t *testing.T) {

	pos1 := make([]float64, 3)
	pos_obs := make([]float64, 3)

	pos1[0] = -11326037.660523077473
	pos1[1] = 3957622.032589655370
	pos1[2] = -25034667.448240540922
	pos_obs[0] = 8.808874982543
	pos_obs[1] = 1.173739504753
	pos_obs[2] = 0.090145080957
	expDiflt := -0.020259558111

	diflt := d_light(pos1, pos_obs)
	if err := checkFloat(diflt, expDiflt, ABS_ERR); err != nil {
		t.Fail()
	}

	// 2nd test
	pos1[0] = -11326037.660490395501
	pos1[1] = 3957622.032455488108
	pos1[2] = -25034667.448276538402
	pos_obs[0] = 0.690697060183
	pos_obs[1] = -3.905584678288
	pos_obs[2] = -1.704455414592
	expDiflt = 0.004034133999
	diflt = d_light(pos1, pos_obs)
	if err := checkFloat(diflt, expDiflt, ABS_ERR); err != nil {
		t.Fail()
	}

}

func TestProperMotion(t *testing.T) {

	var jd_tdb1, jd_tdb2 float64
	pos := make([]float64, 3)
	vel := make([]float64, 3)
	pos2 := make([]float64, 3)
	expPos2 := make([]float64, 3)

	jd_tdb1 = 2451545.000000000000
	pos[0] = -11326046.253463707864
	pos[1] = 3957634.030029209331
	pos[2] = -25034680.670301370323
	vel[0] = -0.007145191487
	vel[1] = 0.009862925083
	vel[2] = -0.010580159896
	jd_tdb2 = 2450300.499353490304
	expPos2[0] = -11326037.361268281937
	expPos2[1] = 3957621.755612566601
	expPos2[2] = -25034667.503285538405

	proper_motion(jd_tdb1, pos, vel, jd_tdb2, pos2)
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}

	// 2nd test
	jd_tdb1 = 2451545.000000000000
	pos[0] = 7059283.251987532713
	pos[1] = 57507101.873731084168
	pos[2] = -302451.154011160368
	vel[0] = -0.000148609903
	vel[1] = 0.009331014346
	vel[2] = 0.000382449798
	jd_tdb2 = 2450300.496618329082
	expPos2[0] = 7059283.436933059245
	expPos2[1] = 57507090.261252179742
	expPos2[2] = -302451.629971226852
	proper_motion(jd_tdb1, pos, vel, jd_tdb2, pos2)
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}

func TestBary2obsAndLightime(t *testing.T) {

	pos := make([]float64, 3)
	pos_obs := make([]float64, 3)
	pos2 := make([]float64, 3)
	var lighttime float64

	pos[0] = 9.494239258762
	pos[1] = 0.498148482267
	pos[2] = -0.202596468352
	pos_obs[0] = 0.685364276219
	pos_obs[1] = -0.675591022487
	pos_obs[2] = -0.292741549309
	pos2[0] = 0.690697060183
	pos2[1] = -3.905584678288
	pos2[2] = -1.704455414592
	expLighttime := 0.051328103712

	bary2obs(pos, pos_obs, pos2, &lighttime)
	err := checkFloat(lighttime, expLighttime, ABS_ERR)
	if err != nil {
		t.Fail()
	}

	// 2nd test
	pos[0] = 1.376061336402
	pos[1] = -4.581175700775
	pos[2] = -1.997196963901
	pos_obs[0] = 0.685364276219
	pos_obs[1] = -0.675591022487
	pos_obs[2] = -0.292741549309
	pos2[0] = -0.690556491831
	pos2[1] = 0.681767468723
	pos2[2] = 0.295545347810
	expLighttime = 0.024932466703
	bary2obs(pos, pos_obs, pos2, &lighttime)
	err = checkFloat(lighttime, expLighttime, ABS_ERR)
	if err != nil {
		t.Fail()
	}
}

// Tested in bary2obs
/*
func TestLightTime(t *testing.T) {

		var jd_tdb, tlight, tlight0 float64
		var ss_object Object
		pos_obs := make([]float64, 3)
		accuracy := int16(0)
		pos := make([]float64, 3)

		errCode := light_time(jd_tdb, ss_object, pos_obs, tlight0, accuracy, pos, &tlight)

	fmt.Println("TODO: Finish test case")
	t.Fail()
}
*/

func TestLimbAngle(t *testing.T) {
	LIMBANG_ABS_ERR := 1.409e-7

	pos_obj := make([]float64, 3)
	pos_obs := make([]float64, 3)
	var limbAng, nadirAng float64

	pos_obj[0] = -11326038.046632558107
	pos_obj[1] = 3957622.431203589309
	pos_obj[2] = -25034667.210543990135
	pos_obs[0] = -0.000015736292
	pos_obs[1] = -0.000025781292
	pos_obs[2] = 0.000029989852

	expLimbAng := -34.813027594445
	expNadirAng := 0.613188582284

	limb_angle(pos_obj, pos_obs, &limbAng, &nadirAng)
	if err := checkFloat(limbAng, expLimbAng, LIMBANG_ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(nadirAng, expNadirAng, LIMBANG_ABS_ERR); err != nil {
		t.Fail()
	}

	// 2nd test
	pos_obj[0] = 276309.046717812540
	pos_obj[1] = 215499.199604222376
	pos_obj[2] = 27281467.053247429430
	pos_obs[0] = -0.000015736292
	pos_obs[1] = -0.000025781292
	pos_obs[2] = 0.000029989852
	expLimbAng = 44.106605680506
	expNadirAng = 1.490073396450

	limb_angle(pos_obj, pos_obs, &limbAng, &nadirAng)
	if err := checkFloat(limbAng, expLimbAng, LIMBANG_ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(nadirAng, expNadirAng, LIMBANG_ABS_ERR); err != nil {
		t.Fail()
	}
}

func TestAberration(t *testing.T) {

	pos := make([]float64, 3)
	ve := make([]float64, 3)
	var lighttime float64
	pos2 := make([]float64, 3)
	expPos2 := make([]float64, 3)

	pos[0] = -11326037.660533571616
	pos[1] = 3957622.032587273978
	pos[2] = -25034667.448236171156
	ve[0] = 0.012479921905
	ve[1] = 0.010578006607
	ve[2] = 0.004630147640
	lighttime = 160334.513329698297
	expPos2[0] = -11324544.167887844145
	expPos2[1] = 3959495.429304560646
	expPos2[2] = -25035046.850750934333

	aberration(pos, ve, lighttime, pos2)
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}

func TestFrameTie(t *testing.T) {

	pos1 := make([]float64, 3)
	var direction int16
	pos2 := make([]float64, 3)
	expPos2 := make([]float64, 3)

	pos1[0] = -11324544.167887844145298004150
	pos1[1] = 3959495.429304560646414756775
	pos1[2] = -25035046.850750934332609176636
	direction = 1
	expPos2[0] = -11324546.465012604370713233948
	expPos2[1] = 3959493.800052606035023927689
	expPos2[2] = -25035046.069331254810094833374

	frame_tie(pos1, direction, pos2)
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}

	// 2nd Test
	pos1[0] = -0.000015736287749355013
	pos1[1] = -0.000025781292172724189
	pos1[2] = 0.000029989854300406425
	direction = -1
	expPos2[0] = -0.000015736291990254234
	expPos2[1] = -0.000025781292050342656
	expPos2[2] = 0.000029989852180327329
	frame_tie(pos1, direction, pos2)
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}

func TestPrecession(t *testing.T) {

	var jd_tdb1, jd_tdb2 float64
	pos1 := make([]float64, 3)
	pos2 := make([]float64, 3)
	expPos2 := make([]float64, 3)

	jd_tdb1 = 2451545.000000000000
	pos1[0] = -11324546.465012604371
	pos1[1] = 3959493.800052606035
	pos1[2] = -25035046.069331254810
	jd_tdb2 = 2450300.499999990221
	expPos2[0] = -11329814.341698952019
	expPos2[1] = 3968123.718317992054
	expPos2[2] = -25031295.943741660565

	if errCode := precession(jd_tdb1, pos1, jd_tdb2, pos2); errCode != 0 {
		fmt.Println("Got error: ", errCode, " expected 0")
		t.Fail()
	}
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}

	// 2nd Test
	jd_tdb1 = 2450300.499999990221
	pos1[0] = 0.000162348211
	pos1[1] = -0.000099207506
	pos1[2] = -0.000000005793
	jd_tdb2 = 2451545.000000000000
	expPos2[0] = 0.000162423740
	expPos2[1] = -0.000099083790
	expPos2[2] = 0.000000047969
	if errCode := precession(jd_tdb1, pos1, jd_tdb2, pos2); errCode != 0 {
		fmt.Println("Got error: ", errCode, " expected 0")
		t.Fail()
	}
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}

func TestNutation(t *testing.T) {
	var jd_tdb float64
	var direction, accuracy int16
	pos := make([]float64, 3)
	pos2 := make([]float64, 3)
	expPos2 := make([]float64, 3)

	jd_tdb = 2450300.50
	direction = 0
	accuracy = 0
	pos[0] = -11329814.341698952019
	pos[1] = 3968123.718317992054
	pos[2] = -25031295.943741660565
	expPos2[0] = -11329676.905038006604
	expPos2[1] = 3966790.417636922095
	expPos2[2] = -25031569.477250423282

	nutation(jd_tdb, direction, accuracy, pos, pos2)
	for idx, p := range pos2 {
		if err := checkFloat(p, expPos2[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}

/* TODO: Write Test
func TestCioLocation(t *testing.T) {
	fmt.Println("FAIL CioLocaton. TODO: Write Test")
	t.Fail()
}
*/

/*
TODO: Write Test

	func TestCioBasis(t *testing.T) {
		fmt.Println("FAIL CioBasis. TODO: Write Test")
		t.Fail()
	}
*/
func TestRadVel(t *testing.T) {

	var cel_object Object

	cel_object.Type = 2
	cel_object.Number = 0
	cel_object.Name = "THETA CAR"
	cel_object.Star.Starname = "Theta CAR"
	cel_object.Star.Catalog = "HIP"
	cel_object.Star.Starnumber = 2
	cel_object.Star.Ra = 10.715944806000
	cel_object.Star.Dec = -64.394450000000
	cel_object.Star.Promora = -18.870000000000
	cel_object.Star.Promodec = 12.060000000000
	cel_object.Star.Parallax = 7.430000000000
	cel_object.Star.Radialvelocity = 24.000000000000

	pos := make([]float64, 3)
	pos[0] = -11326038.046632558107
	pos[1] = 3957622.431203589309
	pos[2] = -25034667.210543990135
	vel := make([]float64, 3)
	vel[0] = -0.007145191487
	vel[1] = 0.009862925083
	vel[2] = -0.010580159896
	vel_obs := make([]float64, 3)
	vel_obs[0] = 0.012479921905
	vel_obs[1] = 0.010578006607
	vel_obs[2] = 0.004630147640
	var d_obs_geo, d_obs_sun, d_obj_sun, rv float64
	d_obs_geo = 0.000042564036
	d_obs_sun = 1.014407315840
	d_obj_sun = 27761077.556339330971
	expRV := 37.431502600179

	rad_vel(&cel_object, pos, vel, vel_obs, d_obs_geo, d_obs_sun,
		d_obj_sun, &rv)

	if err := checkFloat(rv, expRV, ABS_ERR); err != nil {
		t.Fail()
	}

	// 2nd test
	cel_object.Type = 2
	cel_object.Number = 0
	cel_object.Name = "DELTA ORI"
	cel_object.Star.Starname = "Delta ORI"
	cel_object.Star.Catalog = "HIP"
	cel_object.Star.Starnumber = 1
	cel_object.Star.Ra = 5.533444639000
	cel_object.Star.Dec = -0.299091944000
	cel_object.Star.Promora = 1.670000000000
	cel_object.Star.Promodec = 0.560000000000
	cel_object.Star.Parallax = 3.560000000000
	cel_object.Star.Radialvelocity = 16.000000000000
	pos[0] = 7059282.751568783075
	pos[1] = 57507090.936843201518
	pos[2] = -302451.337229677534
	vel[0] = -0.000148609903
	vel[1] = 0.009331014346
	vel[2] = 0.000382449798
	vel_obs[0] = 0.012479921905
	vel_obs[1] = 0.010578006607
	vel_obs[2] = 0.004630147640
	d_obs_geo = 0.000042564036
	d_obs_sun = 1.014407315840
	d_obj_sun = 57939552.311117000878
	expRV = -4.773797553024
	rad_vel(&cel_object, pos, vel, vel_obs, d_obs_geo, d_obs_sun,
		d_obj_sun, &rv)

	if err := checkFloat(rv, expRV, ABS_ERR); err != nil {
		t.Fail()
	}
}

func TestVector2RaDec(t *testing.T) {
	pos := make([]float64, 3)
	var ra, dec float64

	// from novas.c
	pos[0] = -11329676.905038006604
	pos[1] = 3966790.417636922095
	pos[2] = -25031569.477250423282

	expRa := float64(10.713575393650)
	expDec := float64(-64.379669946942)

	vector2radec(pos, &ra, &dec)

	if err := checkFloat(ra, expRa, ABS_ERR); err != nil {
		t.Fail()
	}

	if err := checkFloat(dec, expDec, ABS_ERR); err != nil {
		t.Fail()
	}
}

func TestGravDef(t *testing.T) {
	jd_tdb := float64(2450203.500000)
	loc_code := int16(1)
	accuracy := int16(0)
	pos1 := make([]float64, 3)
	pos1[0] = 276311.148212
	pos1[1] = 215497.640561
	pos1[2] = 27281467.995022

	pos_obs := make([]float64, 3)
	pos_obs[0] = -0.776345
	pos_obs[1] = -0.587443
	pos_obs[2] = -0.254567
	pos2 := make([]float64, 3)

	exppos2 := [3]float64{276310.577695, 215497.220939, 27281468.004114}

	err := grav_def(jd_tdb, loc_code, accuracy, pos1, pos_obs, pos2)
	if err != 0 {
		fmt.Println("Got Error from grav_def= ", err, " expected 0")
		t.Fail()
	}
	for idx := 0; idx < 3; idx++ {
		if math.Abs(pos2[idx]-exppos2[idx]) > TEST_ABS_ERR {
			fmt.Println("Got pos2[", idx, "]= ", pos2[idx], " expected: ", exppos2[idx])
			t.Fail()
		}
	}

	// second call
	jd_tdb = 2450203.500000
	loc_code = 1
	accuracy = 0
	pos1[0] = 7059284.227693
	pos_obs[0] = -0.776345
	pos1[1] = 57507089.943582
	pos_obs[1] = -0.587443
	pos1[2] = -302451.412502
	pos_obs[2] = -0.254567

	exppos2[0] = 7059281.786355
	exppos2[1] = 57507090.238697
	exppos2[2] = -302452.281746

	err = grav_def(jd_tdb, loc_code, accuracy, pos1, pos_obs, pos2)
	if err != 0 {
		fmt.Println("2nd call: Got Error from grav_dev on second call= ", err, " expected 0")
		t.Fail()
	}
	for idx := 0; idx < 3; idx++ {
		if math.Abs(pos2[idx]-exppos2[idx]) > TEST_ABS_ERR {
			fmt.Println("2nd call: Got pos2[", idx, "]= ", pos2[idx], " expected: ", exppos2[idx])
			t.Fail()
		}
	}
}

func TestEtilt(t *testing.T) {

	var jd_tdb float64
	var accuracy int16

	var mobl, tobl, ee, dpsi, deps float64
	var expMobl, expTobl, expEe, expDpsi, expDeps float64

	jd_tdb = 2450300.5
	accuracy = 0
	expMobl = 23.439722735557
	expTobl = 23.437188578694
	expEe = 0.274478377631
	expDpsi = 4.487980745412
	expDeps = -9.122964707301

	Etilt(jd_tdb, accuracy, &mobl, &tobl, &ee, &dpsi, &deps)
	if err := checkFloat(mobl, expMobl, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(tobl, expTobl, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(ee, expEe, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(dpsi, expDpsi, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(deps, expDeps, ABS_ERR); err != nil {
		t.Fail()
	}
}

func TestEect(t *testing.T) {

	var jd_high, jd_low float64
	var accuracy int16

	var expRtn float64

	jd_high = 2450300.50
	jd_low = 0.000000000000
	accuracy = 0
	expRtn = -0.000000002195

	rtn := ee_ct(jd_high, jd_low, accuracy)
	if err := checkFloat(rtn, expRtn, ABS_ERR); err != nil {
		t.Fail()
	}

	// again to test statics(globals)
	rtn = ee_ct(jd_high, jd_low, accuracy)
	if err := checkFloat(rtn, expRtn, ABS_ERR); err != nil {
		t.Fail()
	}

	// 2nd Test
	jd_high = 2450417.50
	jd_low = 0.000000000000
	accuracy = 0
	expRtn = -0.000000001094
	rtn = ee_ct(jd_high, jd_low, accuracy)
	if err := checkFloat(rtn, expRtn, ABS_ERR); err != nil {
		t.Fail()
	}

	// 3rd Test
	jd_high = 2450300.500000000000
	jd_low = 0.0
	accuracy = 0
	expRtn = -0.000000002195
	rtn = ee_ct(jd_high, jd_low, accuracy)
	if err := checkFloat(rtn, expRtn, ABS_ERR); err != nil {
		t.Fail()
	}

}

func TestMeanObliq(t *testing.T) {

	var jd_tdb float64
	var expRtn float64

	jd_tdb = 2450300.499999990221
	expRtn = 84383.001848004453

	rtn := mean_obliq(jd_tdb)
	if err := checkFloat(rtn, expRtn, ABS_ERR); err != nil {
		t.Fail()
	}

	// 2nd Test
	jd_tdb = 2450417.499999988824
	expRtn = 84382.851816523558
	rtn = mean_obliq(jd_tdb)
	if err := checkFloat(rtn, expRtn, ABS_ERR); err != nil {
		t.Fail()
	}

}

func TestNutationAngles(t *testing.T) {
	var tdb float64
	var accuracy int16
	var dpsi, deps float64
	var expDpsi, expDeps float64

	tdb = -0.034072553046
	expDpsi = 4.487980745412
	expDeps = -9.122964707301

	nutation_angles(tdb, accuracy, &dpsi, &deps)
	if err := checkFloat(dpsi, expDpsi, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(deps, expDeps, ABS_ERR); err != nil {
		t.Fail()
	}

}

func TestIau2000a(t *testing.T) {

	var jd_high, jd_low float64
	var dpsi, deps float64
	var expDpsi, expDeps float64

	jd_high = 2451545.000000000000
	jd_low = -1244.500000009779
	expDpsi = 0.000021758345
	expDeps = -0.000044229381
	iau2000a(jd_high, jd_low, &dpsi, &deps)
	if err := checkFloat(dpsi, expDpsi, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(deps, expDeps, ABS_ERR); err != nil {
		t.Fail()
	}
}

func TestFundArgs(t *testing.T) {

	var tdb float64
	a := make([]float64, 5)
	expA := make([]float64, 5)

	tdb = -0.034072553046
	expA[0] = -4.964070968820
	expA[1] = -2.601420958394
	expA[2] = -2.978479161397
	expA[3] = -1.981595887543
	expA[4] = 3.332627977606
	fund_args(tdb, a)
	for idx, v := range a {
		if err := checkFloat(v, expA[idx], ABS_ERR); err != nil {
			t.Fail()
		}
	}
}

func TestTer2Cel(t *testing.T) {
	year := int16(2008)
	month := int16(4)
	day := int16(24)
	leap_sec := int16(33)
	accuracy := int16(0)
	de_num := int16(0)

	hour := 10.605
	ut1_utc := -0.387845
	latitude := 42.0
	longitude := -70.0
	height := 0.0
	temperature := 10.0
	pressure := 1010.0
	//xpole := -0.002
	//ypole := 0.529

	/*
		jd_ut1 := 2.460570411914789e+06
		delta_t := 69.35842
		var acc int16 = 0
		xp := 0.0
		yp := 0.0
	*/

	var geo_loc OnSurface
	//var obs_loc Observer
	var star, dummy_star CatEntry
	var moon, mars Object
	//var sky_pos t_place

	MakeOnSurface(latitude, longitude, height, temperature, pressure, &geo_loc)
	//make_observer_on_surface(latitude, longitude, height, temperature, pressure,
	//	&obs_loc)

	MakeCatEntry("GMB 1830", "FK6", 1307, 11.88299133, 37.71867646,
		4003.27, -5815.07, 109.21, -98.8, &star)

	MakeCatEntry("DUMMY", "xxx", 0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, &dummy_star)

	if error := MakeObject(0, 11, "Moon", &dummy_star, &moon); error != 0 {
		fmt.Printf("Error %d MakeCatEntry (Moon)\n", error)
	}

	if error := MakeObject(0, 4, "Mars", &dummy_star, &mars); error != 0 {
		fmt.Printf("Error %d MakeCatEntry (Mars)\n", error)
	}

	var jd_beg, jd_end float64
	if error := EphemOpen("JPLEPH", &jd_beg, &jd_end, &de_num); error != 0 {
		fmt.Printf("Error %d Opening JPLEPH\n", error)
	}

	jd_utc := JulianDate(year, month, day, hour)
	jd_tt := jd_utc + (float64(leap_sec)+32.184)/86400.0
	jd_ut1 := jd_utc + ut1_utc/86400.0
	delta_t := 32.184 + float64(leap_sec) - ut1_utc

	//jd_tdb := jd_tt // Approximation good to 0.0017 seconds.

	var rat, dect float64
	exprat := 11.8915479153
	expdect := 37.6586695456
	if err := TopoStar(jd_tt, delta_t, &star, &geo_loc, accuracy, &rat, &dect); err != nil {
		fmt.Printf("Error %v TopoStar\n", err)
	}
	if err := checkFloat(rat, exprat, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(dect, expdect, ABS_ERR); err != nil {
		t.Fail()
	}

	/*
		fmt.Printf("FK6 1307 geocentric and topocentric positions:\n")
		//fmt.Printf("%15.10f        %15.10f\n", ra, dec)
		fmt.Printf("%15.10f        %15.10f\n", rat, dect)
		fmt.Printf("\n")
	*/

	var dist float64
	exprat = 17.1031967646
	expdect = -28.2902502967
	expdist := 0.002703785126
	if error := TopoPlanet(jd_tt, &moon, delta_t, &geo_loc, accuracy,
		&rat, &dect, &dist); error != nil {
		fmt.Printf("Error %s from topo_planet.", error)
	}
	if err := checkFloat(rat, exprat, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(dect, expdect, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(dist, expdist, ABS_ERR); err != nil {
		t.Fail()
	}

	/*
		fmt.Printf("Moon geocentric and topocentric positions:\n")
		//fmt.Printf("%15.10f        %15.10f        %15.12f\n", ra, dec, dis)
		fmt.Printf("%15.10f        %15.10f        %15.12f\n", rat, dect, dist)
	*/

	var zd, az, rar, decr float64
	expzd := 81.6891016502
	expaz := 219.2708903405
	Equ2hor(jd_ut1, delta_t, accuracy, 0.0, 0.0, &geo_loc, rat, dect, 1,
		&zd, &az, &rar, &decr)

	if err := checkFloat(zd, expzd, ABS_ERR); err != nil {
		t.Fail()
	}
	if err := checkFloat(az, expaz, ABS_ERR); err != nil {
		t.Fail()
	}
}

/* TODO: Write test
func TestWobble(t *testing.T) {
	fmt.Println("FAIL Wobble: TODO: Write Test")
	t.Fail()
}
*/
