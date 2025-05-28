package novas

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"testing"
)

func TestReadBinaryFloa64Slice(t *testing.T) {
	n := 2
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40,
		0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	buf := bytes.NewReader(b)

	fs, err := readBinary2Float64Slice(buf, n, binary.LittleEndian)
	if err != nil {
		t.Fail()
	}
	expectedVal := 3.141592653589793
	for _, val := range fs {
		if val != expectedVal {
			fmt.Println("expected: ", expectedVal, " got: ", val)
			t.Fail()
		}
	}
}

func TestTest(t *testing.T) {
	test()
}

func TestReadEphemHeader(t *testing.T) {
	/*  From NOVAS eph_manager.c using printf statements
			ttl[0]: 74
			SS[2]: 32.000000
		  ncon: 156
			JPLAU: 149597870.691000
			EM_RATION: 81.300560
			IPT10: 14    NOVAS C has IPT[3][120], In Go, must be [120][3]IPT
			IPT01: 171   so IPT01 in NOVAS is IPT[1][0] in GO
		  denum: 405
	    LPT0: 899
	    LPT1: 10
	    LPT2: 4

	*/
	//var eh EphemHeader

	var expectedEh EphemHeader
	expectedEh.TTL = [252]byte{74, 80, 76, 32, 80, 108, 97, 110, 101, 116, 97, 114, 121, 32, 69, 112, 104, 101, 109, 101, 114, 105, 115, 32, 68, 69, 52, 48, 53, 47, 68, 69, 52, 48, 53, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 83, 116, 97, 114, 116, 32, 69, 112, 111, 99, 104, 58, 32, 74, 69, 68, 61, 32, 32, 50, 51, 48, 53, 52, 50, 52, 46, 53, 32, 49, 53, 57, 57, 32, 68, 69, 67, 32, 48, 57, 32, 48, 48, 58, 48, 48, 58, 48, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 70, 105, 110, 97, 108, 32, 69, 112, 111, 99, 104, 58, 32, 74, 69, 68, 61, 32, 32, 50, 53, 50, 53, 48, 48, 56, 46, 53, 32, 50, 50, 48, 49, 32, 70, 69, 66, 32, 50, 48, 32, 48, 48, 58, 48, 48, 58, 48, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
	expectedEh.Epoch.JDbegin = 2.3054245e+06
	expectedEh.Epoch.JDend = 2.5250085e+06
	expectedEh.Epoch.TimeFactor = 32.0
	expectedEh.Ncon = 156
	expectedEh.Jplau = 1.49597870691e+08
	expectedEh.EmRatio = 81.30056
	expectedEh.IPT[0][0] = 3
	expectedEh.IPT[0][1] = 14
	expectedEh.IPT[1][0] = 171
	expectedEh.LPT[0] = 899
	expectedEh.LPT[1] = 10
	expectedEh.LPT[2] = 4
	ephemFilename := GetEphemFilename()
	fp, err := os.Open(ephemFilename)
	if err != nil {
		fmt.Println("Error opening JPLEPH")
		t.Fail()
	}
	err = ReadBinary2EphemHeader(fp)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	//fmt.Println(eh)
	for idx, ttl := range eh.TTL {
		if ttl != expectedEh.TTL[idx] {
			fmt.Println("Got TTL[", idx, "]: ", ttl, " Expected: ", expectedEh.TTL[idx])
			t.Fail()
		}
	}
	if eh.IPT[0][0] != expectedEh.IPT[0][0] {
		fmt.Println("Got IPT[0][0]: ", eh.IPT[0][0], " Expected: ", expectedEh.IPT[0][0])
		t.Fail()
	}
	if eh.IPT[0][1] != expectedEh.IPT[0][1] {
		fmt.Println("Got IPT[0][1]: ", eh.IPT[0][1], " Expected: ", expectedEh.IPT[0][1])
		t.Fail()
	}
	if eh.IPT[1][0] != expectedEh.IPT[1][0] {
		fmt.Println("Got IPT[1][0]: ", eh.IPT[1][0], " Expected: ", expectedEh.IPT[1][0])
		t.Fail()
	}
	/*
		fmt.Println("TTL[0]: ", eh.TTL[0])
		fmt.Println("JDbegin: ", eh.Epoch.JDbegin)
		fmt.Println("JDend: ", eh.Epoch.JDend)
		fmt.Println("SS[2]: ", eh.Epoch.TimeFactor)
		fmt.Println("Ncon: ", eh.Ncon)
		fmt.Println("Jplau: ", eh.Jplau)
		fmt.Println("EmRatio: ", eh.EmRatio)
		fmt.Println("IPT00: ", eh.IPT[0][0])
		fmt.Println("IPT10: ", eh.IPT[1][0])
		fmt.Println("Dnum: ", eh.DEnum)
		fmt.Println("LPT0: ", eh.LPT[0])
		fmt.Println("LPT1: ", eh.LPT[1])
		fmt.Println("LPT2: ", eh.LPT[2])
	*/
}

/*
func TestReadBinaryIntSlice(t *testing.T) {
	n := 2
	b := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10}
	buf := bytes.NewReader(b)

	sl, err := readBinary2IntSlice(buf, n, binary.LittleEndian)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	expectedVal := 1
	for _, val := range sl {
		if val != expectedVal {
			fmt.Println("expected: ", expectedVal, " got: ", val)
			t.Fail()
		}
	}
}
*/

func TestEphemOpen(t *testing.T) {
	var jd_begin float64
	var jd_end float64
	var de_number int16

	ex_jd_begin := 2.3054245e+06
	ex_jd_end := 2.5250085e+06
	ex_de_number := int16(405)
	rtn := EphemOpen("", &jd_begin, &jd_end,
		&de_number)
	if rtn != 0 {
		fmt.Println("Got rtn= ", rtn, " expected 0")
		t.Fail()
	}
	if jd_begin != ex_jd_begin {
		fmt.Println("Got jd_begin= ", jd_begin, " expected: ", ex_jd_begin)
		t.Fail()
	}
	if jd_end != ex_jd_end {
		fmt.Println("Got jd_end= ", jd_end, " expected: ", ex_jd_end)
		t.Fail()
	}
	if de_number != ex_de_number {
		fmt.Println("Got de_number= ", de_number, " expected: ", ex_de_number)
		t.Fail()
	}
}

func TestSplit(t *testing.T) {
	val := 12345.50
	fr := make([]float64, 4)
	split(val, fr)
	if fr[0] != 12345. {
		fmt.Println("Got: ", fr[0], " Expected: ", 12345.0)
		t.Fail()
	}
	if fr[1] != 0.5 {
		fmt.Println("Got: ", fr[1], " Expected: ", .50)
		t.Fail()
	}
	//fmt.Println("fr0: ", fr[0])
	//fmt.Println("fr1: ", fr[1])

	val = -12345.75
	split(val, fr[2:])
	if fr[2] != -12346. {
		fmt.Println("Got: ", fr[0], " Expected: ", -12346.0)
		t.Fail()
	}
	if fr[3] != 0.25 {
		fmt.Println("Got: ", fr[1], " Expected: ", .50)
		t.Fail()
	}
	//fmt.Println("fr0: ", fr[0])
	//fmt.Println("fr1: ", fr[1])

}

func TestState(t *testing.T) {
	/*
		    Indirectly tests interpolate function.

										JPL Ephemeris DE405 open. jd_beg = 2305424.50  jd_end = 2525008.50

										inside functioin state called with:
										   jed[0]= 2450203.500000017229
								       jed[1]= 0.000000000000
										   target= 9
										   calling interpolate with:
										      target: 9
										      BUFFER= 144604.775240
										      t[0],t[1]= 0.343750, 32.000000
										      IPT[1][target]= 13
										      IPT[2][target]= 8
						           target_pos[0]= -0.00259811933318615
						           target_vel[0]= -0.000002720689436937
						           target_pos[1]=  0.00017511673539366
						           target_vel[1]= -0.000553317582229670
						           target_pos[2]= -0.00001538456962788
						           target_vel[2]= -0.000184376219255533
										inside functioin state called with:
										   jed[0]= 2450203.500000017229
								       jed[1]= 0.000000000000
										   target= 2
										   calling interpolate with:
										      target: 2
										      BUFFER= -120263458.253920
										      t[0],t[1]= 0.343750, 32.000000
										      IPT[1][target]= 13
										      IPT[2][target]= 2
						           target_pos[0]= -0.77635201791459341
						           target_vel[0]= 0.010767895467870266
						           target_pos[1]= -0.58745898514689132
						           target_vel[1]= -0.012155741363468383
						           target_pos[2]= -0.25459718897479006
						           target_vel[2]= -0.005270042318329569
	*/
	expTargetPos := make([]float64, 3)
	expTargetVel := make([]float64, 3)
	expTargetPos[0] = -0.00259811933318615
	expTargetPos[1] = 0.00017511673539366
	expTargetPos[2] = -0.00001538456962788
	expTargetVel[0] = -0.000002720689436937
	expTargetVel[1] = -0.000553317582229670
	expTargetVel[2] = -0.000184376219255533

	var jd_begin float64
	var jd_end float64
	var de_number int16

	rtn := EphemOpen("", &jd_begin, &jd_end,
		&de_number)

	if rtn != 0 {
		fmt.Println("Error opening JPL ephemeris file")
		t.Fail()
	}

	jed := make([]float64, 2)
	jed[0] = 2450203.500000
	jed[1] = 0.0

	target := int16(9)

	targetPos := make([]float64, 3)
	targetVel := make([]float64, 3)

	rtn = state(jed, target, targetPos, targetVel)
	if rtn != 0 {
		fmt.Println("state return error", rtn)
		t.Fail()
	}
	for idx, tp := range targetPos {
		if math.Abs(tp-expTargetPos[idx]) > ABS_ERR {
			fmt.Println("Got: targetPos[", idx, "]= ", tp, " Expected: ", expTargetPos[idx])
			t.Fail()
		}
		if math.Abs(targetVel[idx]-expTargetVel[idx]) > ABS_ERR {
			fmt.Println("Got: targetVel[", idx, "]= ", targetVel[idx], " Expected: ", expTargetVel[idx])
			t.Fail()
		}
	}

	// 2nd call
	jed[0], jed[1] = 2450203.500000, 0.000000
	target = 2
	expTargetPos[0] = -0.77635201791459341
	expTargetPos[1] = -0.58745898514689132
	expTargetPos[2] = -0.25459718897479006
	expTargetVel[0] = 0.010767895467870266
	expTargetVel[1] = -0.012155741363468383
	expTargetVel[2] = -0.005270042318329569

	rtn = state(jed, target, targetPos, targetVel)
	if rtn != 0 {
		fmt.Println("2nd call: state return error", rtn)
		t.Fail()
	}
	for idx, tp := range targetPos {
		if math.Abs(tp-expTargetPos[idx]) > ABS_ERR {
			fmt.Println("2nd call: Got: targetPos[", idx, "]= ", tp, " Expected: ", expTargetPos[idx])
			t.Fail()
		}
		if math.Abs(targetVel[idx]-expTargetVel[idx]) > ABS_ERR {
			fmt.Println("2nd call: Got: targetVel[", idx, "]= ", targetVel[idx], " Expected: ", expTargetVel[idx])
			t.Fail()
		}
	}

}

func checkPlanetEphemeris(errCode int16, pos, vel []float64, expPos, expVel [3]float64) error {
	if errCode != 0 {
		fmt.Println("PlanetEphemeris: Got: ", errCode, " expected: 0")
		return errors.New("")
	}
	for idx, p := range pos {
		if math.Abs(p-expPos[idx]) > ABS_ERR {
			fmt.Println("Got: pos[", idx, "]: ", p, " expected: ", expPos[idx])
			return errors.New("")
		}
		if math.Abs(vel[idx]-expVel[idx]) > ABS_ERR {
			fmt.Println("Got vel[", idx, "]: ", vel[idx], " expected: ", expVel[idx])
			errors.New("")
		}
	}
	return nil
}

func TestPlanetEphemeris(t *testing.T) {
	//fmt.Println("TODO: PlanetEphemeris")
	tjd := [2]float64{2450300.50, 0.0}
	target := int16(5)
	center := int16(11)

	// from eph_manager.c
	expPos := [3]float64{9.494239258762, 0.498148482267, -0.202596468352}
	expVel := [3]float64{-0.000519417132, 0.005134496687, 0.002142896333}

	pos := make([]float64, 3)
	vel := make([]float64, 3)

	errCode := PlanetEphemeris(tjd, target, center, pos, vel)
	err := checkPlanetEphemeris(errCode, pos, vel, expPos, expVel)
	if err != nil {
		t.Fail()
	}

	// 2nd test
	tjd[0] = 2450300.495965856127
	tjd[1] = 0.0
	target = int16(4)
	center = int16(11)

	// from eph_manager.c
	expPos[0] = 1.376032395735
	expPos[1] = -4.581184701051
	expPos[2] = -1.997200116756
	expVel[0] = 0.007173953577
	expVel[1] = 0.002231010868
	expVel[2] = 0.000781535866

	errCode = PlanetEphemeris(tjd, target, center, pos, vel)
	err = checkPlanetEphemeris(errCode, pos, vel, expPos, expVel)
	if err != nil {
		t.Fail()
	}

	// 3rd test
	tjd[0] = 2450300.499350775033
	tjd[1] = 0.000000000000
	target = int16(10)
	center = int16(11)

	// from eph_manager.c
	expPos[0] = -0.005192211046
	expPos[1] = 0.006176448717
	expPos[2] = 0.002803799441
	expVel[0] = -0.000007033372
	expVel[1] = -0.000003820588
	expVel[2] = -0.000001448105

	errCode = PlanetEphemeris(tjd, target, center, pos, vel)
	err = checkPlanetEphemeris(errCode, pos, vel, expPos, expVel)
	if err != nil {
		t.Fail()
	}

}
