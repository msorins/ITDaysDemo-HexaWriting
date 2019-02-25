package RobotMiniProj

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/gorilla/mux"

	"mind/core/framework/skill"
	"mind/core/framework/log"
	"mind/core/framework/drivers/distance"
	"mind/core/framework/drivers/hexabody"
	"mind/core/framework/drivers/media"
	"mind/core/framework"
)

const (
	DISTANCE_TO_WALL   = 175 // millimeters
	MOVE_HEAD_DURATION = 500 // milliseconds
	WALK_SPEED         = 0.7 // cm per second
	SENSE_INTERVAL     = 250// four times per second
	FrameWidth  = 1280
	FrameHeight = 720
)



type RobotMiniProj struct {
	skill.Base
	degreeX            float64
	crtHeightMM        float64
	timeWalkingInFront int64
	lettersAnimation   map[string]string
	helloAnimation string
	clapAnimation string
	stop bool
}

func NewSkill() skill.Interface {
	// Use this method to create a new skill.
	return &RobotMiniProj{}
}


func (d *RobotMiniProj) OnStart() {
	// Use this method to do something when this skill is starting.
	log.Debug.Println("OnStart()")

	// Define all start up variables
	d.degreeX = 0
	d.timeWalkingInFront = 0
	d.lettersAnimation = map[string]string{
		"A": "V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A30V8A90V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A98V7A91V8A71V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A98V7A79V8A71V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A102V7A79V8A66V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A102V7A93V8A62V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A99V7A98V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200",
		"D": "V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A30V8A90V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A97V7A112V8A10V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A89V7A102V8A28V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A89V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200",
		"I": "V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A30V8A90V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A80V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A110V7A89V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A110V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A114V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A110V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A110V7A90V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200",
		"S": "V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A30V8A90V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A94V8A75V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A94V8A55V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A99V7A94V8A55V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A99V7A94V8A65V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A99V7A94V8A70V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A113V7A94V8A70V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A113V7A94V8A70V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A113V7A94V8A60V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A106V7A94V8A60V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A106V7A94V8A60V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200",
		"T": "V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A30V8A90V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A112V8A10V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A90V8A66V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A64V8A66V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200",
		"Y": "V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A30V8A90V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A90V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A105V7A102V8A43V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A110V7A109V8A27V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A110V7A89V8A27V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A111V7A73V8A64V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A111V7A83V8A64V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A104V7A98V8A39V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200|V0A90V1A81V2A133V3A90V4A81V5A133V6A104V7A84V8A39V9A90V10A81V11A133V12A90V13A81V14A133V15A90V16A81V17A133V18A0T200",

		" ": "",
	}
	d.helloAnimation = "V0A90V1A90V2A45V3A90V4A90V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A70V2A45V3A60V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A120V1A70V2A45V3A120V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A70V2A45V3A60V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A120V1A70V2A45V3A120V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A70V2A45V3A60V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A120V1A70V2A45V3A120V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A70V2A45V3A120V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A110V2A45V3A120V4A40V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A40V2A45V3A120V4A110V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A110V2A45V3A120V4A40V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A40V2A45V3A120V4A110V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A110V2A45V3A120V4A40V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300|V0A60V1A40V2A45V3A120V4A110V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T300"
	d.clapAnimation  = "V0A90V1A90V2A45V3A90V4A90V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T200|V0A60V1A70V2A45V3A120V4A70V5A45V6A120V7A110V8A135V9A90V10A70V11A135V12A90V13A70V14A135V15A60V16A110V17A135V18A0T200"
	d.stop = false

	// Start hexabody
	err := hexabody.Start()
	if err != nil {
		log.Error.Println("Hexabody start err:", err)
		return
	}

	// Start the distance
	err = distance.Start()
	if err != nil {
		log.Error.Println("Distance start err:", err)
		return
	}
	if !distance.Available() {
		log.Error.Println("Distance sensor is not available")
	}

	// Start media
	err = media.Start()
	if err != nil {
		log.Error.Println("Media start err:", err)
		return
	}

	// Execute the sequence of operations
	//d.executeSeqOfOperations()

	// Start the API
	go d.StartAPI()

	// Send the video feed to web
	go d.SendVideoFeedToWeb()
}

func (d *RobotMiniProj) OnClose() {
	// Use this method to do something when this skill is closing.
	log.Debug.Println("OnClose()")
	distance.Close()
	hexabody.Close()
	media.Close()
}

func (d *RobotMiniProj) OnConnect() {
	// Use this method to do something when the remote connected.
	log.Debug.Println("OnConnect()")
}

func (d *RobotMiniProj) OnDisconnect() {
	// Use this method to do something when the remote disconnected.
	log.Debug.Println("OnDisconnect()")
}

func (d *RobotMiniProj) OnRecvJSON(data []byte) {
	// Use this method to do something when skill receive json data from remote client.
}

func (d *RobotMiniProj) OnRecvString(data string) {
	// Use this method to do something when skill receive string from remote client.
	switch data {
	case "ReExec":
		d.executeSeqOfOperations()
	case "DrawHeadLeftRight":
		d.DrawHeadLeftRight()
	case "MoveFront":
		d.MoveFront()
	case "DrawLeft":
		d.DrawLeft()
	case "DrawRight":
		d.DrawRight()
	case "DrawUp":
		d.DrawUp(10)
	case "DrawDown":
		d.DrawDown(10)
	case "PitchFront":
		d.PitchFront(0.2)
	case "PitchBack":
		d.PitchBack(1)
	case "PointToBoard":
		d.PointToBoard(5, 5)
	}
}

func (d *RobotMiniProj) executeSeqOfOperations() {
	// Move head to front of robot
	d.moveHeadTo0()

	// Do the detection

	// Start walking until wall is reached
	d.walkFront()

	// Point user to robot's decision
	for i := 0; i <= 9; i++ {
		d.PointToBoard(i, 0)
		d.PointToBoard(i, 5)
		hexabody.StandWithHeight(100)
	}


	// Go back to starting position
	d.walkBack()
}

func (d *RobotMiniProj) moveHeadTo0() {
	hexabody.MoveHead(0, MOVE_HEAD_DURATION)
}

func (d *RobotMiniProj) getDistance() float64 {
	distanceVal, err := distance.Value()
	if err != nil {
		log.Error.Println(err)
	}

	log.Debug.Println("Distance: ", distanceVal)
	return distanceVal
}

func (d *RobotMiniProj) walkFront() {
	hexabody.SetStepLength(0.4)
	hexabody.WalkContinuously(0, WALK_SPEED)
	log.Debug.Println("walkFront()")
	start := time.Now()
	for {
		if d.getDistance() <= DISTANCE_TO_WALL {
			hexabody.StopWalkingContinuously()
			break
		}
		time.Sleep(SENSE_INTERVAL * time.Millisecond)
	}
	elapsed := time.Since(start)
	log.Debug.Println("Time walked: ", elapsed.Nanoseconds())
	d.timeWalkingInFront = elapsed.Nanoseconds()
	hexabody.SetStepLength(1.0)
}

func (d *RobotMiniProj) walkBack() {
	hexabody.SetStepLength(0.4)
	hexabody.WalkContinuously(180, WALK_SPEED)
	log.Debug.Println("walkBack)")
	start := time.Now()
	for {
		if time.Since(start).Nanoseconds() >= d.timeWalkingInFront {
			hexabody.StopWalkingContinuously()
			break
		}
		time.Sleep(SENSE_INTERVAL * time.Millisecond)
	}

	hexabody.SetStepLength(1.0)
}

func (d *RobotMiniProj) StandToHeight(heightMM float64) {
	hexabody.StandWithHeight(heightMM)
}

func (d *RobotMiniProj) DrawUp(mm float64) {
	d.crtHeightMM += mm
	hexabody.StandWithHeight(d.crtHeightMM)
}

func (d *RobotMiniProj) DrawDown(mm float64) {
	d.crtHeightMM -= mm
	hexabody.StandWithHeight(d.crtHeightMM)
}

func (d *RobotMiniProj) DrawLeft() {
	hexabody.SetStepLength(0.3)
	hexabody.Walk(90, 500 )
	hexabody.SetStepLength(1.0)
}

func (d *RobotMiniProj) DrawRight() {
	hexabody.SetStepLength(0.3)
	hexabody.Walk(270, 500)
	hexabody.SetStepLength(1.0)
}

func (d *RobotMiniProj) MoveFront() {
	log.Debug.Println("moveFront()")
	log.Debug.Println("startDistance: ", d.getDistance())
	hexabody.SetStepLength(0.3)
	hexabody.Walk(0, 500)
	hexabody.SetStepLength(1.0)
	log.Debug.Println("finishDistance: ", d.getDistance())
}

func (d *RobotMiniProj) DrawHeadLeftRight() {
	hexabody.MoveHead(0, MOVE_HEAD_DURATION)
	hexabody.MoveHead(20, MOVE_HEAD_DURATION)
	hexabody.MoveHead(-20, MOVE_HEAD_DURATION)
	hexabody.MoveHead(20, MOVE_HEAD_DURATION)
	hexabody.MoveHead(-20, MOVE_HEAD_DURATION)
	hexabody.MoveHead(0, MOVE_HEAD_DURATION)
}

func (d *RobotMiniProj) PitchFront(degree float64) {
	for startDegree := d.degreeX; d.degreeX >= startDegree - degree; d.degreeX -= 1 {
		hexabody.Pitch(d.degreeX, 100)
	}
}

func (d *RobotMiniProj) PitchBack(degree float64) {
	for startDegree := d.degreeX; d.degreeX <= startDegree + degree; d.degreeX += 1 {
		hexabody.Pitch(d.degreeX, 100)
	}
}

func (d *RobotMiniProj) PointToBoard(pos int, sleep int) {
	// Pos is a number [0,1, 2, 3, 4, 5, 6, 7, 8]
	log.Debug.Println("pointToBoard(", pos, ")")

	if pos == 0 {
		hexabody.MoveJoint(0, 1, 15, 200)
		hexabody.MoveJoint(0, 0, 75, 200)
		hexabody.MoveJoint(0, 2, 46, 200)
	}

	if pos == 1 {
		hexabody.MoveJoint(1, 1, 15, 200)
		hexabody.MoveJoint(1, 0, 130, 200)
		hexabody.MoveJoint(1, 2, 46, 200)
	}

	if pos == 2 {
		hexabody.MoveJoint(1, 1, 15, 200)
		hexabody.MoveJoint(1, 0, 100, 200)
		hexabody.MoveJoint(1, 2, 46, 200)
	}

	if pos == 3 {
		hexabody.MoveJoint(0, 1, 40, 200)
		hexabody.MoveJoint(0, 0, 75, 200)
		hexabody.MoveJoint(0, 2, 46, 200)
	}

	if pos == 4 {
		hexabody.MoveJoint(1, 1, 40, 200)
		hexabody.MoveJoint(1, 0, 130, 200)
		hexabody.MoveJoint(1, 2, 46, 200)
	}

	if pos == 5 {
		hexabody.MoveJoint(1, 1, 40, 200)
		hexabody.MoveJoint(1, 0, 100, 200)
		hexabody.MoveJoint(1, 2, 46, 200)
	}

	if pos == 6 {
		hexabody.MoveJoint(0, 1, 74, 200)
		hexabody.MoveJoint(0, 0, 75, 200)
		hexabody.MoveJoint(0, 2, 46, 200)
	}

	if pos == 7 {
		hexabody.MoveJoint(1, 1, 74, 200)
		hexabody.MoveJoint(1, 0, 130, 200)
		hexabody.MoveJoint(1, 2, 46, 200)
	}

	if pos == 8 {
		hexabody.MoveJoint(1, 1, 74, 200)
		hexabody.MoveJoint(1, 0, 100, 200)
		hexabody.MoveJoint(1, 2, 46, 200)
	}

	time.Sleep(time.Second * time.Duration(sleep))
}

func (d *RobotMiniProj) DrawHorizontalLine(left float64, right float64) {
	// Go left till end
	hexabody.WalkContinuously(90, WALK_SPEED)
	time.Sleep(time.Second * time.Duration(left / WALK_SPEED))
	hexabody.StopWalkingContinuously()

	//
	//// Go right till end
	//hexabody.WalkContinuously(270, WALK_SPEED / (left + right))
	//time.Sleep(time.Second)
	//hexabody.StopWalkingContinuously()
	//
	//
	//// go back
	//hexabody.WalkContinuously(90, WALK_SPEED / right)
	//time.Sleep(time.Second)
	//hexabody.StopWalkingContinuously()

}

func (d *RobotMiniProj) StartAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/moveFront",   d.MoveFrontAPI).Queries("dist", "{dist}")
	r.HandleFunc("/moveBack",    d.MoveBackAPI).Queries("dist", "{dist}")
	r.HandleFunc("/moveLeft",    d.MoveLeftAPI).Queries("dist", "{dist}")
	r.HandleFunc("/moveRight",   d.MoveRightAPI).Queries("dist", "{dist}")
	r.HandleFunc("/write",       d.writeLetter).Queries("letter", "{letter}")
	r.HandleFunc("/writeITDAYS", d.writeITDAYS)
	r.HandleFunc("/stop", d.stopAnimation)
	log.Error.Println(http.ListenAndServe(":8000", r))
}

func (d *RobotMiniProj) MoveFrontDistance(distInMM int) {
	log.Debug.Println("MoveFrontDistance(" + strconv.Itoa(distInMM) + ")")
}

func (d *RobotMiniProj) MoveBackDistance(distInMM int) {
	log.Debug.Println("MoveBackDistance(" + strconv.Itoa(distInMM) + ")")
}

func (d *RobotMiniProj) MoveLeftDistance(distInMM int) {
	log.Debug.Println("MoveLeftDistance(" + strconv.Itoa(distInMM) + ")")
}

func (d *RobotMiniProj) MoveRightDistance(distInMM int) {
	log.Debug.Println("MoveRightDistance(" + strconv.Itoa(distInMM) + ")")
}


func (d *RobotMiniProj) MoveFrontAPI(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println("/moveFront request received")

	// Get dist parameter
	vars := mux.Vars(r)
	log.Debug.Println(vars)
	distS := vars["dist"]
	distI, err := strconv.Atoi(distS)

	if distS != "" && err == nil {
		// Deal with the request
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `OK`)

		// Call move front on the robot
		d.MoveFrontDistance(distI)
	} else {
		log.Debug.Println("Error, not all parameters provided")
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `dist parameter is missing`)
	}
}

func (d *RobotMiniProj) RunCommand(command string) {
	var no, degree, duration int
	acts := map[int]int{}
	for {
		no, degree := -1, -1
		fmt.Sscanf(command, "V%dA%d", &no, &degree)

		if no == -1 {
			break
		}
		acts[no] = degree
		command = command[strings.IndexAny(command[1:], "VT")+1:]
	}

	fmt.Sscanf(command, "T%d", &duration)

	for no, degree = range acts {
		if no == 18 {
			go hexabody.MoveHead(float64(degree), duration)
		} else {
			go hexabody.MoveJoint(no/3, no%3, float64(degree), duration)
		}
	}

	time.Sleep(time.Millisecond * time.Duration(duration * 2))
}

func (d *RobotMiniProj) RunAnimaton(animation string) {
	cmds := strings.Split(animation, "|")
	for _, cmd := range cmds {
		d.RunCommand(cmd)
	}
}

func (d *RobotMiniProj) writeLetter(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println("/writeLetter request received")
	hexabody.Stand()
	time.Sleep(1 * time.Second)
	hexabody.Stand()

	vars := mux.Vars(r)
	distS, found := vars["letter"]

	var toWrite string
	if found == true {
		toWrite = string(distS[0])
	} else {
		return
	}

	d.RunAnimaton(d.lettersAnimation[toWrite])
}


func (d *RobotMiniProj) writeITDAYS(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println("/writeITDAYS request received")

	// Stand up && wait 1 second
	hexabody.Stand()
	time.Sleep(1 * time.Second)

	// Do Hello animation
	for i := 0; i < 2; i++ {
		if(d.stop) {
			d.stop = false
			return
		}
		d.RunAnimaton(d.helloAnimation)
	}

	// Walk in front 10 seconds
	//d.MoveLegsContinously(10000, 0.0)

	// Start writing
	toWrite := "IT DAYS"
	for _, chr := range(toWrite) {
		if(d.stop) {
			d.stop = false
			return
		}
		chrString := string(chr)

		// Draw the letter
		if chrString != " " {
			log.Debug.Println("drawLetter: ", chrString)
			d.RunAnimaton(d.lettersAnimation[chrString])
			d.RunAnimaton(d.lettersAnimation[chrString])
		}

		// Move to the right
		d.MoveLegsContinously(4000, 270.0)
	}

	// Go Left for x seconds
	d.MoveLegsContinously(len(toWrite) * 4000, 90)

	// Go Back for 10 seconds
	//d.MoveLegsContinously(10000, 180)

	// Do the clap animation
	for i:= 0; i < 100; i++ {
		if(d.stop) {
			d.stop = false
			return
		}
		d.RunAnimaton(d.clapAnimation)
	}

	// Deal with the request
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `OK`)
}

func (d *RobotMiniProj) stopAnimation(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println("/stopAnimation called")

	// Stop everything
	d.stop = true

	// Deal with the request
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `OK`)
}

func (d *RobotMiniProj) MoveLegsContinously(ms int, dir float64) {
	hexabody.SetStepLength(0.2)
	hexabody.WalkContinuously(dir, WALK_SPEED)
	log.Debug.Println("moveLegsContinously( ", ms, ", ", dir, " )")
	start := time.Now()
	for {
		duration := time.Since(start)

		if duration.Nanoseconds() > int64(ms) * 1000000{
			hexabody.StopWalkingContinuously()
			break
		}
		time.Sleep(SENSE_INTERVAL * time.Millisecond)
	}
	elapsed := time.Since(start)
	log.Debug.Println("Time walked: ", elapsed.Nanoseconds())
	d.timeWalkingInFront = elapsed.Nanoseconds()
	hexabody.SetStepLength(1.0)
}


func (d *RobotMiniProj) MoveBackAPI(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println("/moveBack request received")

	// Get dist parameter
	vars := mux.Vars(r)
	log.Debug.Println(vars)
	distS := vars["dist"]
	distI, err := strconv.Atoi(distS)

	if distS != "" && err == nil {
		// Deal with the request
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `OK`)

		// Call move front on the robot
		d.MoveBackDistance(distI)
	} else {
		log.Debug.Println("Error, not all parameters provided")
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `dist parameter is missing`)
	}
}

func (d *RobotMiniProj) MoveLeftAPI(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println("/moveLeft request received")

	// Get dist parameter
	vars := mux.Vars(r)
	log.Debug.Println(vars)
	distS := vars["dist"]
	distI, err := strconv.Atoi(distS)

	if distS != "" && err == nil {
		// Deal with the request
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `OK`)

		// Call move front on the robot
		d.MoveLeftDistance(distI)
	} else {
		log.Debug.Println("Error, not all parameters provided")
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `dist parameter is missing`)
	}
}

func (d *RobotMiniProj) MoveRightAPI(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println("/moveRight request received")

	// Get dist parameter
	vars := mux.Vars(r)
	log.Debug.Println(vars)
	distS := vars["dist"]
	distI, err := strconv.Atoi(distS)

	if distS != "" && err == nil {
		// Deal with the request
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `OK`)

		// Call move front on the robot
		d.MoveRightDistance(distI)
	} else {
		log.Debug.Println("Error, not all parameters provided")
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `dist parameter is missing`)
	}
}


func (d *RobotMiniProj) SendVideoFeedToWeb() {
	for true {
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, media.SnapshotYCbCr(), nil)
		str := base64.StdEncoding.EncodeToString(buf.Bytes())
		framework.SendString(str)
	}

}
