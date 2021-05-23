package core

import (
	"cerberus-security-laboratories/des-wristband-ui/internal/gql/models"
	// "cerberus-security-laboratories/des-wristband-ui/internal/gql/resolvers"
	"cerberus-security-laboratories/des-wristband-ui/internal/titan"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

// createRandomString creates a letter string of length n
func createRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return strings.Title(string(b))
}

// createRandomNumberInRange generates random number in a pre-determined range
func createRandomNumberInRange(min int, max int) int {
	return rand.Intn(max-min) + min
}
func createRandomFloatNumberInRange(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// WristbandProxy interface gives a common interface to all the different
// underlying implementations of the WristbandProxy
type WristbandProxy interface {
	Get() (models.Wristband, error)
	Activate()
	Deactivate()
	IsActive() bool
	GetFirstData() (*models.WristbandData, error)
	GetLatestData() (*models.WristbandData, error)
	GetData(int) ([]*models.WristbandData, error)
	GetDataBlock(int, int) ([]*models.WristbandData, error)
	GetWristbandDataChan() (chan *models.WristbandData, error)
	GetAlertChan() (chan *models.Alert, error)
	GetLevelChan() (chan *models.Level, error)
	GetNews2Chan() (chan *models.News2, error)
	GetActiveChan() (chan *models.Active, error)
	GetNews2() (*models.News2, error)
	GetLevel() (*models.Level, error)
	// GetSpecificData(string) (*models.WristbandData, error)
	GetID() string
	// GetMSID() int
	// GetType() int
	// GetTypeVer() int
	// GetKey() string
	// GetTIC() string
	// GetActivatedTime() string
	// GetDeactivatedTime() string
	// GetName() string
	SetName(string)
	// IsOnOxygen() bool
	SetOnOxygen(bool)
	// IsPregnant() bool
	SetPregnant(bool)
	// IsChild() bool
	SetChild(bool)
	SetDepartment(string)
}

// wristbandProxyBasic is a wristband model that generates normal in-range
// sensor readings and requires no configuration
type wristbandProxyBasic struct {
	id              int
	tickPeriod      int
	data            []*models.WristbandData
	dataChan        chan *models.WristbandData
	alertChan       chan *models.Alert
	levelChan       chan *models.Level
	news2Chan       chan *models.News2
	activeChan      chan *models.Active
	prevAlert       *models.Alert
	prevNews2       *models.News2
	isInAlert       bool
	titanIdentity   *titan.Identity
	active          bool
	activatedTime   string
	deactivatedTime string
	name            string
	dateOfBirth     string
	onOxygen        bool
	pregnant        bool
	child           bool
	department      string
	level           *models.Level
}

////////////////////////////
// UTILS
// addData to this Wristband
func (wb *wristbandProxyBasic) addData(wbd *models.WristbandData) error {
	wb.data = append(wb.data, wbd)
	return nil
}

////////////////////////////
// PROTOTYPING DUMMY DATA
func createAlert(wsd models.WristbandSensorData, overall int, level string, target string, trend string, overallLevel string) (*models.Alert, error) {
	var wAlert models.Alert
	wAlert.SensorData = &wsd
	wAlert.Level = level
	wAlert.Target = target
	wAlert.OverallLevel = overallLevel
	wAlert.Overall = overall
	wAlert.Trend = trend
	return &wAlert, nil
}

func createAlertTrends(latest int, prev int) string {
	var alertTrend string
	if latest < prev {
		alertTrend = "falling"
	} else if latest > prev {
		alertTrend = "rising"
	} else {
		alertTrend = "unchanged"
	}

	return alertTrend
}
func createAlertTrendsTemp(latest float64, prev float64) string {
	var alertTrend string
	if latest < prev {
		alertTrend = "falling"
	} else if latest > prev {
		alertTrend = "rising"
	} else {
		alertTrend = "unchanged"
	}

	return alertTrend
}

// Check for Sensor Data Alert
func (wb *wristbandProxyBasic) checkDataForAlert(wsd models.WristbandSensorData, overall int, level string, target string, overallLevel string) error {
	var wAlert *models.Alert
	var alertErr error = nil

	// Check for Alert Trend
	var alertTrend string
	if wb.prevAlert != nil {
		// Check If They are The Same Parameter
		// if target == wb.prevAlert.Target {
		// bloodPressure
		if target == "bloodPressure" {
			alertTrend = createAlertTrends(wsd.BloodPressure, wb.prevAlert.SensorData.BloodPressure)
			// sp02
		} else if target == "sp02" {
			alertTrend = createAlertTrends(wsd.Sp02, wb.prevAlert.SensorData.Sp02)
			// pulse
		} else if target == "pulse" {
			alertTrend = createAlertTrends(wsd.Pulse, wb.prevAlert.SensorData.Pulse)
			// respiration
		} else if target == "respiration" {
			// temperature
			alertTrend = createAlertTrends(wsd.Respiration, wb.prevAlert.SensorData.Respiration)
		} else if target == "temperature" {
			alertTrend = createAlertTrendsTemp(wsd.Temperature, wb.prevAlert.SensorData.Temperature)
			// motion
			// proximity
		} else if target == "motion" || target == "proximity" {
			alertTrend = "none"
		}

		// Check if Overall Level is not Low
		// And make sure this is not the 1st alert, and it is unchanged
		if overallLevel != "low" && alertTrend == "unchanged" {
			// Then make sure to use the last critical trend
			alertTrend = wb.prevAlert.Trend
		}
	} else {
		alertTrend = "none"
	}

	// Send Alert Data to Alert Channel
	wAlert, alertErr = createAlert(wsd, overall, level, target, alertTrend, overallLevel)

	// The Current Alert will be Considered as Previous Data after Using it
	wb.prevAlert = wAlert

	if alertErr != nil {
		panic(fmt.Errorf("Error getting Alert Channel for subscription: %e\n", alertErr))
	} else {
		select {
		case wb.alertChan <- wAlert:
			// values are being read from wb.alertChan
			fmt.Println("wb.alertChan: inserted data")
		default:
			// no subscribers, alert not in channel
			fmt.Println("wb.alertChan: data created, not inserted")
		}
	}
	return nil
}
func (wb *wristbandProxyBasic) createLevel(level string) (*models.Level, error) {
	var wLevel models.Level
	wLevel.ID = fmt.Sprint(wb.id)
	wLevel.Text = level
	return &wLevel, nil
}

// Check for Each Wristband's Level
func (wb *wristbandProxyBasic) checkForLevel(overallLevel string) error {
	var wLevel *models.Level
	var levelErr error = nil

	// Send Level Data to Level Channel
	wLevel, levelErr = wb.createLevel(overallLevel)

	if levelErr != nil {
		panic(fmt.Errorf("Error getting Alert Channel for subscription: %e\n", levelErr))
	} else {
		// resolvers.UpdateLevelSummary(wLevel)
		wb.level = wLevel
		log.Println(wb.level.Text)
		// Assign the Level to the Wristband Level State
		select {
		case wb.levelChan <- wb.level:
			fmt.Println("level inserted")
		default:
			fmt.Println("level not inserted")
		}
	}
	return nil
}

// Sensor Data
func (wb *wristbandProxyBasic) createSensorData(respiration int, sp02 int, pulse int, temperature float64, bloodPressure int, motion bool, proximity bool, overall int, level string, target string, overallLevel string) *models.WristbandSensorData {
	var wsd models.WristbandSensorData

	// populate data to the sensorData object
	wsd.Respiration = respiration
	wsd.Sp02 = sp02
	wsd.Pulse = pulse
	wsd.Temperature = temperature
	wsd.BloodPressure = bloodPressure
	wsd.Motion = motion
	wsd.Proximity = proximity

	// log.Println(level)
	// Check for Alert (none and low level is not included)
	if level != "" {
		// if level != "low" {
		// }
		wb.checkDataForAlert(wsd, overall, level, target, overallLevel)
		// All Level is Included for Summary
		wb.checkForLevel(overallLevel)
	}

	return &wsd
}

// News2 Score
func (wb *wristbandProxyBasic) createNews2(respiration int, sp02 int, airNotOxygen bool, pulse int, temperature int, bloodPressure int, motion bool, overall int) *models.News2 {
	var wnews models.News2

	// populate data to the sensorData object
	wnews.Respiration = respiration
	wnews.Sp02 = sp02
	wnews.OnOxygen = airNotOxygen
	wnews.Pulse = pulse
	wnews.Temperature = temperature
	wnews.BloodPressure = &bloodPressure
	wnews.Motion = motion
	wnews.Overall = overall

	wb.prevNews2 = &wnews

	select {
	case wb.news2Chan <- wb.prevNews2:
		fmt.Println("wb.news2Chan inserted")
	default:
		fmt.Println("wb.news2Chan not inserted")
	}

	return &wnews
}

// VITAL SCORE (RANGE(PULSE), RANGE(RESPIRATION), RANGE(BLOODPRESSURE), RANGE(TEMPERATURE), RANGE(SP02), RANGE(MOTION), RANGE(BATTERY))
// HEALTHY VITAL SCORE
// NEED TO KNOW ABOUT THE LOW AND HIGH RANGE FOR MOTION SCORES (NORMAL: 6 - 10?)
func healthyScores() (int, int, int, int, int, int, float64, float64, int, int, int, int) {
	return 51, 90, 12, 20, 111, 219, 37, 38, 96, 250, 80, 100
}

// LOW VITAL SCORE
func lowScores() (int, int, int, int, int, int, float64, float64, int, int, int, int) {
	return 0, 50, 0, 11, 0, 110, 0, 36, 0, 87, 0, 79
} // HIGH VITAL SCORE
func highScores() (int, int, int, int, int, int, float64, float64, int, int, int, int) {
	return 91, 200, 21, 40, 220, 400, 39, 40, 91, 99, 90, 100
}

// RANDOM VITAL SCORE
func randomScores() (int, int, int, int, int, int, float64, float64, int, int, int, int) {
	return 40, 131, 8, 25, 90, 220, 35, 40, 91, 99, 0, 100
}

// createWristbandData generates dummy vital score data for testing and returns it as a models.WristbandData
func (wb *wristbandProxyBasic) createWristbandData(condition string) (*models.WristbandData, error) {
	var wbData models.WristbandData

	var pulseRange [2]int
	var respirationRange [2]int
	var bloodPressureRange [2]int
	var temperatureRange [2]float64
	var sp02Range [2]int
	var batteryRange [2]int

	if condition == "healthy" {
		pulseRange[0], pulseRange[1], respirationRange[0], respirationRange[1],
			bloodPressureRange[0], bloodPressureRange[1], temperatureRange[0],
			temperatureRange[1], sp02Range[0], sp02Range[1], batteryRange[0], batteryRange[1] = healthyScores()
	} else if condition == "low" {
		pulseRange[0], pulseRange[1], respirationRange[0], respirationRange[1],
			bloodPressureRange[0], bloodPressureRange[1], temperatureRange[0],
			temperatureRange[1], sp02Range[0], sp02Range[1], batteryRange[0], batteryRange[1] = lowScores()
	} else if condition == "high" {
		pulseRange[0], pulseRange[1], respirationRange[0], respirationRange[1],
			bloodPressureRange[0], bloodPressureRange[1], temperatureRange[0],
			temperatureRange[1], sp02Range[0], sp02Range[1], batteryRange[0], batteryRange[1] = highScores()
	} else if condition == "random" {
		pulseRange[0], pulseRange[1], respirationRange[0], respirationRange[1],
			bloodPressureRange[0], bloodPressureRange[1], temperatureRange[0],
			temperatureRange[1], sp02Range[0], sp02Range[1], batteryRange[0], batteryRange[1] = randomScores()
	}

	// Vital Scores:

	// Variables Used For Alert
	var level string = ""
	var target string

	// Now Check for Individual Parameter to Get The Target One that Causes Abnormal Vital Score
	// If Two or More Scores Are As High As One Another, Then Prioritise Descending Alphabetical Order
	// Blood Pressure (B)
	// blood pressure (int)
	bloodPressure := createRandomNumberInRange(bloodPressureRange[0], bloodPressureRange[1])
	var bpAggregateScore int
	if bloodPressure <= 90 || bloodPressure >= 220 {
		level = "high"
		target = "bloodPressure"
		bpAggregateScore = 3
	} else if bloodPressure >= 91 && bloodPressure <= 100 {
		if level == "medium" || level == "low-medium" || level == "low" || level == "" {
			level = "medium"
			target = "bloodPressure"
		}
		bpAggregateScore = 2
	} else if bloodPressure >= 101 && bloodPressure <= 110 {
		if level == "low-medium" || level == "low" || level == "" {
			level = "low-medium"
			target = "bloodPressure"
		}
		bpAggregateScore = 1
	} else if bloodPressure >= 111 && bloodPressure <= 219 {
		if level == "low" || level == "" {
			level = "low"
			target = "bloodPressure"
		}
		bpAggregateScore = 0
	}

	// Motion (M)
	// motion (bool)
	var motion bool
	if condition == "healthy" {
		motion = true
	} else {
		if rand.Float64() > 0.5 {
			motion = true
		} else {
			motion = false
		}
	}

	// Proximity (P, r)
	// proximity (bool)
	var proximity bool
	if condition == "healthy" {
		proximity = true
	} else {
		if rand.Float64() > 0.2 {
			proximity = true
		} else {
			proximity = false
		}
	}

	// Pulse (P, u)
	// pulse (int)
	pulse := createRandomNumberInRange(pulseRange[0], pulseRange[1])
	var pulseAggregateScore int
	if pulse <= 40 || pulse >= 131 {
		level = "high"
		target = "pulse"
		pulseAggregateScore = 3
	} else if pulse >= 111 && pulse <= 130 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "pulse"
		}
		pulseAggregateScore = 2
	} else if pulse >= 41 && pulse <= 50 || pulse >= 91 && pulse <= 110 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "pulse"
		}
		pulseAggregateScore = 1
	} else if pulse >= 51 && pulse <= 90 {
		if level == "low" {
			level = "low"
			target = "pulse"
		}
		pulseAggregateScore = 0
	}

	// Respiration (R)
	// respiration (int)
	respiration := createRandomNumberInRange(respirationRange[0], respirationRange[1])
	var resAggregateScore int
	if respiration <= 8 || respiration >= 25 {
		level = "high"
		target = "respiration"
		resAggregateScore = 3
	} else if respiration >= 21 && respiration <= 24 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "respiration"
		}
		resAggregateScore = 2
	} else if respiration >= 9 && respiration <= 11 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "respiration"
		}
		resAggregateScore = 1
	} else if respiration >= 12 && respiration <= 20 {
		if level == "low" {
			level = "low"
			target = "respiration"
		}
		resAggregateScore = 0
	}

	// airNotOxygen (bool)
	var airNotOxygen bool
	if condition == "healthy" {
		airNotOxygen = true
	} else {
		if rand.Float64() > 0.5 {
			airNotOxygen = true
		} else {
			airNotOxygen = false
		}
	}

	// Sp02 (S)
	// sp02 (int) (Scale 2)
	sp02 := createRandomNumberInRange(sp02Range[0], sp02Range[1])
	// var sp02AggregateScore2 int
	// if (sp02 <= 83 || sp02 >= 97 && !airNotOxygen) {
	// 	sp02AggregateScore2 = 3
	// } else if (sp02 >= 84 && sp02 <= 85 || sp02 >= 95 && sp02 <= 96 && !airNotOxygen) {
	// 	sp02AggregateScore2 = 2
	// } else if (sp02 >= 86 && sp02 <= 87 || sp02 >= 93 && sp02 <= 94 && !airNotOxygen) {
	// 	sp02AggregateScore2 = 1
	// } else if (sp02 >= 88 && sp02 <= 92 || sp02 >= 93 && airNotOxygen) {
	// 	sp02AggregateScore2 = 0
	// }

	// sp02 (Scale 1)
	var sp02AggregateScore1 int
	if sp02 <= 91 {
		level = "high"
		target = "sp02"
		sp02AggregateScore1 = 3
	} else if sp02 >= 92 && sp02 <= 93 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "sp02"
		}
		sp02AggregateScore1 = 2
	} else if sp02 >= 94 && sp02 <= 95 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "sp02"
		}
		sp02AggregateScore1 = 1
	} else {
		if level == "low" {
			level = "low"
			target = "sp02"
		}
		sp02AggregateScore1 = 0
	}

	// Temperature (T)
	// temperature (int)
	temperature := createRandomFloatNumberInRange(temperatureRange[0], temperatureRange[1])
	var tempAggregateScore int
	if temperature <= 35 {
		level = "high"
		target = "temperature"
		tempAggregateScore = 3
	} else if float64(temperature) >= 39.1 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "temperature"
		}
		tempAggregateScore = 2
	} else if float64(temperature) >= 35.1 && temperature <= 36 || float64(temperature) >= 38.1 && temperature <= 39 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "temperature"
		}
		tempAggregateScore = 1
	} else if float64(temperature) >= 36.1 && temperature <= 38 {
		if level == "low" {
			level = "low"
			target = "temperature"
		}
		tempAggregateScore = 0
	}

	// THE ONLY REASON WHY MOTION AND PROXIMITY IS NOT ORDERED ALPHABETICALLY WITH THE ABOVE PARAMETERS IS BECAUSE THEY ARE BOOLEAN VARIABLES
	// AND THEY WILL THROW ERROR INSTEAD OF AN ALERT STATE, ERROR STATE WILL BE PRIORITISED BEFORE ALERT
	// SO THEY WILL BE ORDERED ALPHABETICALLY DESCENDING SEPARATELY
	if !motion {
		level = "high"
		target = "motion"
	}
	if !proximity {
		level = "error"
		target = "proximity"

		// Deactivate This Wristband Immediately to Stop Generating Data
		// wb.Deactivate()
	}

	// overall (int)
	overall := tempAggregateScore + pulseAggregateScore + resAggregateScore + bpAggregateScore + sp02AggregateScore1

	// Baterry Level
	batteryLevel := createRandomNumberInRange(batteryRange[0], batteryRange[1])

	// Overall Level
	var overallLevel string = ""
	if overall >= 0 && overall <= 4 {
		overallLevel = "low"
	}
	if level == "high" && overall <= 4 {
		overallLevel = "low-medium"
	} else if overall >= 5 && overall <= 6 {
		overallLevel = "medium"
	} else if overall >= 7 {
		overallLevel = "high"
	} else if !proximity {
		overallLevel = "error"
	}

	//check for low battery or no connection
	if batteryLevel < 20 {
		level = "error"
		target = "batteryLevel"
		if batteryLevel < 1 {
			level = "error"
			target = "noConnection"
			overallLevel = "error"
		}
	}

	// check for child or pregnant
	if wb.child || wb.pregnant {
		// level = "error"
		overallLevel = "error"
	}

	// populate all the generated data to the object
	wbData.Time = ""
	wbData.SensorData = wb.createSensorData(respiration, sp02, pulse, temperature, bloodPressure, motion, proximity, overall, level, target, overallLevel)
	wbData.News2 = wb.createNews2(resAggregateScore, sp02AggregateScore1, airNotOxygen, pulseAggregateScore, tempAggregateScore, bpAggregateScore, motion, overall)
	// wbData.Location = append(wbd.Location, data.Location...)
	// wbData.Location = []
	wbData.BatteryLevel = batteryLevel

	return &wbData, nil
}

////////////////////////////

// Wristband Data Alert Channel
func (wb *wristbandProxyBasic) GetAlertChan() (chan *models.Alert, error) {
	log.Println("wristbandProxyBasic.GetAlertChan() called")
	return wb.alertChan, nil
}

// Wristband Data Channel
func (wb *wristbandProxyBasic) GetWristbandDataChan() (chan *models.WristbandData, error) {
	log.Println("wristbandProxyBasic.GetWristbandDataChan() called")
	return wb.dataChan, nil
}

// Summary of All Wristband's Emergency Levels
func (wb *wristbandProxyBasic) GetLevelChan() (chan *models.Level, error) {
	log.Println("wristbandProxyBasic.GetLevelChan() called")
	return wb.levelChan, nil
}

func (wb *wristbandProxyBasic) GetLevel() (*models.Level, error) {
	log.Println("wristbandProxyBasic.GetLevel() called")
	return wb.level, nil
}

func (wb *wristbandProxyBasic) GetNews2Chan() (chan *models.News2, error) {
	log.Println("wristbandProxyBasic.GetNews2Chan() called")
	return wb.news2Chan, nil
}

func (wb *wristbandProxyBasic) GetNews2() (*models.News2, error) {
	return wb.prevNews2, nil
}

func (wb *wristbandProxyBasic) GetActiveChan() (chan *models.Active, error) {
	return wb.activeChan, nil
}

// NewWristbandProxyBasic creates a new instance of a wristbandProxyBasic
func NewWristbandProxyBasic(id int, tickPeriod int, name string, onOxygen bool, dateOfBirth string, pregnant bool, child bool, department string) (WristbandProxy, error) {
	wb := wristbandProxyBasic{id: id, name: name, dateOfBirth: dateOfBirth, department: department, tickPeriod: tickPeriod}

	// Create the Titan Identity
	var err error
	wb.titanIdentity, err = titan.NewIdentity(0x00030001, 0x48650001, 0x00000000)
	if err != nil {
		return nil, err
	}

	// Record the activation
	wb.active = true
	wb.activatedTime = base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02 15:04:05")))

	// Invent a name if required
	if name == "" {
		firstName := createRandomString(4 + rand.Intn(10))
		lastName := createRandomString(5 + rand.Intn(10))
		wb.name = firstName + " " + lastName
	}

	// Other fields
	wb.onOxygen = onOxygen
	wb.pregnant = pregnant
	wb.child = child
	wb.dateOfBirth = dateOfBirth
	wb.department = department

	// Instantiate Data Channel for Each Wristband
	wb.dataChan = make(chan *models.WristbandData, 1)

	// Instantiate Alert Channel for Each Wristband
	wb.alertChan = make(chan *models.Alert, 1)

	// Instantiate Level Channel (String) for Each Wristband
	wb.levelChan = make(chan *models.Level, 1)

	// Instantiate News2 Channel for Each Wristband
	wb.news2Chan = make(chan *models.News2, 1)

	// Instatiate Active Channel for Each Wristband
	wb.activeChan = make(chan *models.Active, 1)

	// Set up a ticker channel to trigger reading in of new data
	tick := time.NewTicker(time.Duration(tickPeriod) * time.Millisecond).C
	// Create a go routine to create wristband data
	go func(ch <-chan time.Time) {
		for range ch {
			if wb.active {
				wbd, err := wb.createWristbandData("random")
				if err != nil {
					return
				}
				wb.addData(wbd)

				select {
				case wb.dataChan <- wbd:
					// values are being read from r.Resolver.wbChan
					fmt.Println("wb.dataChan: inserted data")
				default:
					// no subscribers, wb not in channel
					fmt.Println("wb.dataChan: data created, not inserted")
				}
			}
		}
	}(tick)

	return &wb, nil
}

func (wb *wristbandProxyBasic) Get() (models.Wristband, error) {
	log.Println("wristbandProxyBasic.Get() called")

	// Create and populate structure from object state
	mWb := models.Wristband{
		ID:          fmt.Sprintf("%d", wb.id),
		Msid:        wb.titanIdentity.GetMsid(),
		Type:        wb.titanIdentity.GetType(),
		TypeVer:     wb.titanIdentity.GetTypeVer(),
		Key:         wb.titanIdentity.GetKeyAsString(),
		Tic:         wb.titanIdentity.GetCertificateAsString(),
		Active:      wb.active,
		Activated:   &wb.activatedTime,
		Deactivated: &wb.deactivatedTime,
		Data:        wb.data,
		Name:        &wb.name,
		DateOfBirth: wb.dateOfBirth,
		OnOxygen:    wb.onOxygen,
		Pregnant:    wb.pregnant,
		Child:       wb.child,
		Department:  wb.department,
	}

	return mWb, nil
}

func (wb *wristbandProxyBasic) GetID() string {
	log.Println("wristbandProxyBasic.GetID() called")

	return fmt.Sprintf("%d", wb.id)
}

// GET WRISTBAND DATA
////////////////////////////////
// Latest Data
func (wb *wristbandProxyBasic) GetLatestData() (*models.WristbandData, error) {
	if len(wb.data) == 0 {
		return nil, fmt.Errorf("Wristband has no data")
	}
	return wb.data[len(wb.data)-1], nil
}

// First Data
func (wb *wristbandProxyBasic) GetFirstData() (*models.WristbandData, error) {
	if len(wb.data) == 0 {
		return nil, fmt.Errorf("Wristband has no data")
	}
	return wb.data[0], nil
}

// Multiple Data (from a specified element to the total length of the array)
func (wb *wristbandProxyBasic) GetData(num int) ([]*models.WristbandData, error) {
	if len(wb.data) == 0 {
		return nil, fmt.Errorf("Wristband has no data")
	}
	// if num is 0 then return all
	if num == 0 {
		return wb.data, nil
	} else {
		if len(wb.data) < num {
			// Return The Entire Array
			return wb.data, nil
			// return nil, fmt.Errorf("Wristband does not have the amount of data requested")
		}
		return wb.data[len(wb.data)-num:], nil
	}
}

// Multiple Data (In Block)
func (wb *wristbandProxyBasic) GetDataBlock(start int, end int) ([]*models.WristbandData, error) {
	var reversed []*models.WristbandData

	for i := range wb.data {
		// Reverse The Array
		rData := wb.data[len(wb.data)-1-i]
		// Append To A Temp Array
		reversed = append(reversed, rData)
	}

	// index 0 is now the last element of the array
	// The Returned Array Will Be Backwards In Time
	// End Value Must Be Less Than zWristband Data's Total Size
	// Otherwise
	if len(wb.data) < end {
		// Return The Data Value To The Very Last Data It Has
		return reversed[start:(end - (end - len(wb.data)))], nil
	} else {
		if start < end {
			return reversed[start:end], nil
		}
		return nil, fmt.Errorf("The Start Value Must Be Less Than The End Value")
	}
}

////////////////////////////////

// SET NEW INFO TO EXISTING WRISTBAND DATA
////////////////////////////////
// Set New Name
func (wb *wristbandProxyBasic) SetName(newName string) {
	// Invent a name if required
	if newName == "" {
		firstName := createRandomString(4 + rand.Intn(10))
		lastName := createRandomString(5 + rand.Intn(10))
		wb.name = firstName + " " + lastName
	} else {
		wb.name = newName
	}
}

// Set OnOxygen Info
func (wb *wristbandProxyBasic) SetOnOxygen(isOnOxygen bool) {
	wb.onOxygen = isOnOxygen
}

// Set New Pregnant Info
func (wb *wristbandProxyBasic) SetPregnant(isPregnant bool) {
	wb.pregnant = isPregnant
}

// Set New Child Info
func (wb *wristbandProxyBasic) SetChild(isChild bool) {
	wb.child = isChild
}

// Set New Department
func (wb *wristbandProxyBasic) SetDepartment(newDepartment string) {
	wb.department = newDepartment
}

/////////////////////////////////
// Deactivate / Activate
// Deactivate A Wristband
func (wb *wristbandProxyBasic) Deactivate() {
	var activeWb models.Active

	// reset the active key to false
	wb.active = false
	// Record the activation
	wb.deactivatedTime = base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02 15:04:05")))

	activeWb.ID = fmt.Sprint(wb.id)
	activeWb.Active = wb.active
	select {
	case wb.activeChan <- &activeWb:
		// values are being read from r.Resolver.wbChan
		fmt.Println("wb.activeChan: inserted active state")
	default:
		// no subscribers, wb not in channel
		fmt.Println("wb.activeChan: active chan created, not inserted")
	}
}

// Activate A Wristband
func (wb *wristbandProxyBasic) Activate() {
	var activeWb models.Active

	// reset the active key to true
	wb.active = true
	// Record the activation
	wb.activatedTime = base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02 15:04:05")))

	activeWb.ID = fmt.Sprint(wb.id)
	activeWb.Active = wb.active
	select {
	case wb.activeChan <- &activeWb:
		// values are being read from r.Resolver.wbChan
		fmt.Println("wb.activeChan: inserted active state")
	default:
		// no subscribers, wb not in channel
		fmt.Println("wb.activeChan: active chan created, not inserted")
	}
}

///////////////////////////////////
// IsActive
func (wb *wristbandProxyBasic) IsActive() bool {
	return wb.active
}

// wristbandDataFile matches what is in the Wristband data JSON files
type wristbandDataFile struct {
	Time         string                     `json:"time"`
	SensorData   models.WristbandSensorData `json:"sensorData"`
	News2        models.News2               `json:"news2"`
	BatteryLevel int                        `json:"batteryLevel"`
	Location     []*models.BridgeSignal     `json:"location"`
}

// addData to this Wristband
func (wb *wristbandProxyFile) addData(wbd *models.WristbandData) error {
	wb.data = append(wb.data, wbd)
	return nil
}

// Create Alert
// func createAlert(wsd models.WristbandSensorData, overall int, level string, target string, trend string, overallLevel string) (*models.Alert, error) {
// 	var wAlert models.Alert
// 	wAlert.SensorData = &wsd
// 	wAlert.Level = level
// 	wAlert.Target = target
// 	wAlert.OverallLevel = overallLevel
// 	wAlert.Overall = overall
// 	wAlert.Trend = trend
// 	return &wAlert, nil
// }

// func createAlertTrends(latest int, prev int) string {
// 	var alertTrend string
// 	if latest < prev {
// 		alertTrend = "falling"
// 	} else if latest > prev {
// 		alertTrend = "rising"
// 	} else {
// 		alertTrend = "unchanged"
// 	}

// 	return alertTrend
// }

// Check for Sensor Data Alert
func (wb *wristbandProxyFile) checkDataForAlert(wsd models.WristbandSensorData, overall int, level string, target string, overallLevel string) error {
	var wAlert *models.Alert
	var alertErr error = nil

	// Check for Alert Trend
	var alertTrend string
	if wb.prevAlert != nil {
		// Check If They are The Same Parameter
		// if target == wb.prevAlert.Target {
		// bloodPressure
		if target == "bloodPressure" {
			alertTrend = createAlertTrends(wsd.BloodPressure, wb.prevAlert.SensorData.BloodPressure)
			// sp02
		} else if target == "sp02" {
			alertTrend = createAlertTrends(wsd.Sp02, wb.prevAlert.SensorData.Sp02)
			// pulse
		} else if target == "pulse" {
			alertTrend = createAlertTrends(wsd.Pulse, wb.prevAlert.SensorData.Pulse)
			// respiration
		} else if target == "respiration" {
			alertTrend = createAlertTrends(wsd.Respiration, wb.prevAlert.SensorData.Respiration)
			// temperature
		} else if target == "temperature" {
			alertTrend = createAlertTrendsTemp(wsd.Temperature, wb.prevAlert.SensorData.Temperature)
			// motion
			// proximity
		} else if target == "motion" || target == "proximity" {
			alertTrend = "none"
		}

		// Check if Overall Level is not Low
		// And make sure this is not the 1st alert, and it is unchanged
		if overallLevel != "low" && alertTrend == "unchanged" {
			// Then make sure to use the last critical trend
			alertTrend = wb.prevAlert.Trend
		}
	} else {
		alertTrend = "none"
	}

	// Send Alert Data to Alert Channel
	wAlert, alertErr = createAlert(wsd, overall, level, target, alertTrend, overallLevel)

	// The Current Alert will be Considered as Previous Data after Using it
	wb.prevAlert = wAlert

	if alertErr != nil {
		panic(fmt.Errorf("Error getting Alert Channel for subscription: %e\n", alertErr))
	} else {
		select {
		case wb.alertChan <- wAlert:
			// values are being read from wb.alertChan
			fmt.Println("wb.alertChan: inserted data")
		default:
			// no subscribers, alert not in channel
			fmt.Println("wb.alertChan: data created, not inserted")
		}
	}
	return nil
}
func (wb *wristbandProxyFile) createLevel(level string) (*models.Level, error) {
	var wLevel models.Level
	wLevel.ID = fmt.Sprint(wb.id)
	wLevel.Text = level
	return &wLevel, nil
}

// Check for Each Wristband's Level
func (wb *wristbandProxyFile) checkForLevel(overallLevel string) error {
	var wLevel *models.Level
	var levelErr error = nil

	// Send Level Data to Level Channel
	wLevel, levelErr = wb.createLevel(overallLevel)

	if levelErr != nil {
		panic(fmt.Errorf("Error getting Alert Channel for subscription: %e\n", levelErr))
	} else {
		// resolvers.UpdateLevelSummary(wLevel)
		wb.level = wLevel
		log.Println(wb.level.Text)
		// Assign the Level to the Wristband Level State
		select {
		case wb.levelChan <- wb.level:
			fmt.Println("level inserted")
		default:
			fmt.Println("level not inserted")
		}
	}
	return nil
}

// Check For Alert and Level
func (wb *wristbandProxyFile) checkParameters(wsd *models.WristbandSensorData, overall int, level string, target string, overallLevel string) {
	// Check for Alert (none and low level is not included)
	if level != "" {
		// if level != "low" {
		// }
		wb.checkDataForAlert(*wsd, overall, level, target, overallLevel)
		// All Level is Included for Summary
		wb.checkForLevel(overallLevel)
	}
}

func (wb *wristbandProxyFile) readSingleJSONDataEntry(dec *json.Decoder) (*models.WristbandData, error) {
	if !dec.More() {
		return nil, fmt.Errorf("Finished file")
	}
	var data wristbandDataFile
	dec.Decode(&data)
	var wbd models.WristbandData
	wbd.Time = data.Time
	wbd.SensorData = &data.SensorData
	wbd.News2 = &data.News2
	wbd.Location = append(wbd.Location, data.Location...)
	wbd.BatteryLevel = data.BatteryLevel

	// Create Alert Parameters (level, target, overall level)
	// Vital Scores:
	var wbs models.WristbandSensorData = *wbd.SensorData
	// var wnews models.News2 = *wbd.News2
	// log.Println(&data.SensorData)

	// Variables Used For Alert
	var level string = ""
	var target string

	// Now Check for Individual Parameter to Get The Target One that Causes Abnormal Vital Score
	// If Two or More Scores Are As High As One Another, Then Prioritise Descending Alphabetical Order
	// Blood Pressure (B)
	// blood pressure (int)
	var bpAggregateScore int
	if wbs.BloodPressure <= 90 || wbs.BloodPressure >= 220 {
		level = "high"
		target = "bloodPressure"
		bpAggregateScore = 3
	} else if wbs.BloodPressure >= 91 && wbs.BloodPressure <= 100 {
		if level == "medium" || level == "low-medium" || level == "low" || level == "" {
			level = "medium"
			target = "bloodPressure"
		}
		bpAggregateScore = 2
	} else if wbs.BloodPressure >= 101 && wbs.BloodPressure <= 110 {
		if level == "low-medium" || level == "low" || level == "" {
			level = "low-medium"
			target = "bloodPressure"
		}
		bpAggregateScore = 1
	} else if wbs.BloodPressure >= 111 && wbs.BloodPressure <= 219 {
		if level == "low" || level == "" {
			level = "low"
			target = "bloodPressure"
		}
		bpAggregateScore = 0
	}

	// Motion (M)
	// motion (bool)
	var motion bool = wbs.Motion

	// Proximity (P, r)
	// proximity (bool)
	var proximity bool = wbs.Proximity

	// Pulse (P, u)
	// pulse (int)
	var pulseAggregateScore int
	if wbs.Pulse <= 40 || wbs.Pulse >= 131 {
		level = "high"
		target = "pulse"
		pulseAggregateScore = 3
	} else if wbs.Pulse >= 111 && wbs.Pulse <= 130 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "pulse"
		}
		pulseAggregateScore = 2
	} else if wbs.Pulse >= 41 && wbs.Pulse <= 50 || wbs.Pulse >= 91 && wbs.Pulse <= 110 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "pulse"
		}
		pulseAggregateScore = 1
	} else if wbs.Pulse >= 51 && wbs.Pulse <= 90 {
		if level == "low" {
			level = "low"
			target = "pulse"
		}
		pulseAggregateScore = 0
	}

	// Respiration (R)
	// respiration (int)
	var resAggregateScore int
	if wbs.Respiration <= 8 || wbs.Respiration >= 25 {
		level = "high"
		target = "respiration"
		resAggregateScore = 3
	} else if wbs.Respiration >= 21 && wbs.Respiration <= 24 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "respiration"
		}
		resAggregateScore = 2
	} else if wbs.Respiration >= 9 && wbs.Respiration <= 11 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "respiration"
		}
		resAggregateScore = 1
	} else if wbs.Respiration >= 12 && wbs.Respiration <= 20 {
		if level == "low" {
			level = "low"
			target = "respiration"
		}
		resAggregateScore = 0
	}

	// airNotOxygen (bool)
	// var airNotOxygen bool = wnews.OnOxygen

	// Sp02 (S)
	// sp02 (Scale 1)
	var sp02AggregateScore1 int
	if wbs.Sp02 <= 91 {
		level = "high"
		target = "sp02"
		sp02AggregateScore1 = 3
	} else if wbs.Sp02 >= 92 && wbs.Sp02 <= 93 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "sp02"
		}
		sp02AggregateScore1 = 2
	} else if wbs.Sp02 >= 94 && wbs.Sp02 <= 95 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "sp02"
		}
		sp02AggregateScore1 = 1
	} else {
		if level == "low" {
			level = "low"
			target = "sp02"
		}
		sp02AggregateScore1 = 0
	}

	// Temperature (T)
	var tempAggregateScore int
	if wbs.Temperature <= 35 {
		level = "high"
		target = "temperature"
		tempAggregateScore = 3
	} else if float64(wbs.Temperature) >= 39.1 {
		if level == "medium" || level == "low-medium" || level == "low" {
			level = "medium"
			target = "temperature"
		}
		tempAggregateScore = 2
	} else if float64(wbs.Temperature) >= 35.1 && wbs.Temperature <= 36 || float64(wbs.Temperature) >= 38.1 && wbs.Temperature <= 39 {
		if level == "low-medium" || level == "low" {
			level = "low-medium"
			target = "temperature"
		}
		tempAggregateScore = 1
	} else if float64(wbs.Temperature) >= 36.1 && wbs.Temperature <= 38 {
		if level == "low" {
			level = "low"
			target = "temperature"
		}
		tempAggregateScore = 0
	}

	// THE ONLY REASON WHY MOTION AND PROXIMITY IS NOT ORDERED ALPHABETICALLY WITH THE ABOVE PARAMETERS IS BECAUSE THEY ARE BOOLEAN VARIABLES
	// AND THEY WILL THROW ERROR INSTEAD OF AN ALERT STATE, ERROR STATE WILL BE PRIORITISED BEFORE ALERT
	// SO THEY WILL BE ORDERED ALPHABETICALLY DESCENDING SEPARATELY
	if !motion {
		level = "high"
		target = "motion"
	}
	if !proximity {
		level = "error"
		target = "proximity"

		// Deactivate This Wristband Immediately to Stop Generating Data
		// wb.Deactivate()
	}

	// overall (int)
	overall := tempAggregateScore + pulseAggregateScore + resAggregateScore + bpAggregateScore + sp02AggregateScore1

	// Overall Level
	var overallLevel string = ""
	if overall >= 0 && overall <= 4 {
		overallLevel = "low"
	}
	if level == "high" && overall <= 4 {
		overallLevel = "low-medium"
	} else if overall >= 5 && overall <= 6 {
		overallLevel = "medium"
	} else if overall >= 7 {
		overallLevel = "high"
	} else if !proximity {
		overallLevel = "error"
	}

	if wbd.BatteryLevel < 20 {
		level = "error"
		target = "batteryLevel"
		if wbd.BatteryLevel < 1 {
			level = "error"
			target = "noConnection"
			overallLevel = "error"
		}
	}

	// check for child or pregnant
	if wb.child || wb.pregnant {
		// level = "error"
		overallLevel = "error"
	}

	// Check For Alert, Level
	wb.checkParameters(&wbs, overall, level, target, overallLevel)

	return &wbd, nil
}

type wristbandProxyFile struct {
	id              int
	tickPeriod      int
	data            []*models.WristbandData
	dataChan        chan *models.WristbandData
	alertChan       chan *models.Alert
	levelChan       chan *models.Level
	news2Chan       chan *models.News2
	activeChan      chan *models.Active
	prevAlert       *models.Alert
	prevNews2       *models.News2
	isInAlert       bool
	titanIdentity   *titan.Identity
	active          bool
	activatedTime   string
	deactivatedTime string
	name            string
	dateOfBirth     string
	onOxygen        bool
	pregnant        bool
	child           bool
	department      string
	level           *models.Level
	dataDecoder     *json.Decoder
}

// Wristband Data Alert Channel
func (wb *wristbandProxyFile) GetAlertChan() (chan *models.Alert, error) {
	return wb.alertChan, nil
}

// Wristband Data Channel
func (wb *wristbandProxyFile) GetWristbandDataChan() (chan *models.WristbandData, error) {
	return wb.dataChan, nil
}

// Summary of All Wristband's Emergency Levels
func (wb *wristbandProxyFile) GetLevelChan() (chan *models.Level, error) {
	return wb.levelChan, nil
}

func (wb *wristbandProxyFile) GetLevel() (*models.Level, error) {
	log.Println("wristbandProxyBasic.GetLevel() called")
	return wb.level, nil
}

func (wb *wristbandProxyFile) GetNews2Chan() (chan *models.News2, error) {
	log.Println("wristbandProxyBasic.GetNews2Chan() called")
	return wb.news2Chan, nil
}

func (wb *wristbandProxyFile) GetNews2() (*models.News2, error) {
	return wb.prevNews2, nil
}

func (wb *wristbandProxyFile) GetActiveChan() (chan *models.Active, error) {
	return wb.activeChan, nil
}

// NewWristbandProxyFile creates a new instance of a wristbandProxyFile
func NewWristbandProxyFile(
	id int,
	tickPeriod int,
	name string,
	onOxygen bool,
	dateOfBirth string,
	pregnant bool,
	child bool,
	department string,
	dataDecoder *json.Decoder,
) (WristbandProxy, error) {
	wb := wristbandProxyFile{id: id, name: name, dateOfBirth: dateOfBirth, department: department, tickPeriod: tickPeriod, dataDecoder: dataDecoder}

	// Create the Titan Identity
	var err error
	wb.titanIdentity, err = titan.NewIdentity(0x00030001, 0x48650001, 0x00000000)
	if err != nil {
		return nil, err
	}

	// Record the activation
	wb.active = true
	wb.activatedTime = base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02 15:04:05")))

	// Invent a name if required
	if name == "" {
		firstName := createRandomString(4 + rand.Intn(10))
		lastName := createRandomString(5 + rand.Intn(10))
		wb.name = firstName + " " + lastName
	}

	// Other fields
	wb.onOxygen = onOxygen
	wb.pregnant = pregnant
	wb.child = child
	wb.dateOfBirth = dateOfBirth
	wb.department = department

	// Instantiate Wristband Data Channel
	wb.dataChan = make(chan *models.WristbandData, 1)

	// Instantiate Alert Channel for Each Wristband
	wb.alertChan = make(chan *models.Alert, 1)

	// Instantiate Level Channel (String) for Each Wristband
	wb.levelChan = make(chan *models.Level, 1)

	// Instantiate News2 Channel for Each Wristband
	wb.news2Chan = make(chan *models.News2, 1)

	// Instatiate Active Channel for Each Wristband
	wb.activeChan = make(chan *models.Active, 1)

	// Set up a ticker channel to trigger reading in of new data
	tick := time.NewTicker(time.Duration(tickPeriod) * time.Millisecond).C
	// Create a go routine to read an entry from the JSON file and add to the data array
	go func(ch <-chan time.Time) {
		for range ch {
			if wb.active {
				wbd, err := wb.readSingleJSONDataEntry(wb.dataDecoder)
				if err != nil {
					// File finished
					return
				}
				wb.addData(wbd)

				select {
				case wb.dataChan <- wbd:
					// values are being read from r.Resolver.wbChan
					fmt.Println("wb.dataChan: inserted data")
				default:
					// no subscribers, wb not in channel
					fmt.Println("wb.dataChan: data created, not inserted")
				}
			}
		}
	}(tick)

	return &wb, nil
}

func (wb *wristbandProxyFile) Get() (models.Wristband, error) {
	log.Println("wristbandProxyFile.Get() called")

	// Create and populate structure from object state
	mWb := models.Wristband{
		ID:          fmt.Sprintf("%d", wb.id),
		Msid:        wb.titanIdentity.GetMsid(),
		Type:        wb.titanIdentity.GetType(),
		TypeVer:     wb.titanIdentity.GetTypeVer(),
		Key:         wb.titanIdentity.GetKeyAsString(),
		Tic:         wb.titanIdentity.GetCertificateAsString(),
		Active:      wb.active,
		Activated:   &wb.activatedTime,
		Deactivated: &wb.deactivatedTime,
		Data:        wb.data,
		Name:        &wb.name,
		DateOfBirth: wb.dateOfBirth,
		OnOxygen:    wb.onOxygen,
		Pregnant:    wb.pregnant,
		Child:       wb.child,
		Department:  wb.department,
	}

	return mWb, nil
}

//////////////////////////////////
func (wb *wristbandProxyFile) GetFirstData() (*models.WristbandData, error) {
	if len(wb.data) == 0 {
		return nil, fmt.Errorf("Wristband has no data")
	}
	return wb.data[0], nil
}

func (wb *wristbandProxyFile) GetLatestData() (*models.WristbandData, error) {
	log.Println("wristbandProxyFile.GetLatestData() called")
	if len(wb.data) == 0 {
		return nil, fmt.Errorf("Wristband has no data")
	}
	return wb.data[len(wb.data)-1], nil
}

func (wb *wristbandProxyFile) GetData(num int) ([]*models.WristbandData, error) {
	if len(wb.data) == 0 {
		return nil, fmt.Errorf("Wristband has no data")
	}
	// if num is 0 then return all
	if num == 0 {
		return wb.data, nil
	} else {
		if len(wb.data) < num {
			return nil, fmt.Errorf("Wristband does not have the amount of data requested")
		}
		// get the latest wristbands
		return wb.data[len(wb.data)-num:], nil
	}
}

// Multiple Data (In Block)
func (wb *wristbandProxyFile) GetDataBlock(start int, end int) ([]*models.WristbandData, error) {
	return wb.data[len(wb.data):], nil
}

///////////////////////////////

// SET NEW INFO TO EXISTING WRISTBAND DATA
////////////////////////////////
// Set New Name
func (wb *wristbandProxyFile) SetName(newName string) {
	// Invent a name if required
	if newName == "" {
		firstName := createRandomString(4 + rand.Intn(10))
		lastName := createRandomString(5 + rand.Intn(10))
		wb.name = firstName + " " + lastName
	} else {
		wb.name = newName
	}
}

// Set OnOxygen Info
func (wb *wristbandProxyFile) SetOnOxygen(isOnOxygen bool) {
	wb.onOxygen = isOnOxygen
}

// Set Pregnant Info
func (wb *wristbandProxyFile) SetPregnant(isPregnant bool) {
	wb.pregnant = isPregnant
}

// Set New Child Info
func (wb *wristbandProxyFile) SetChild(isChild bool) {
	wb.child = isChild
}

// Set New Department
func (wb *wristbandProxyFile) SetDepartment(newDepartment string) {
	wb.department = newDepartment
}

/////////////////////////////////
// Deactivate / Activate
// Deactivate A Wristband
func (wb *wristbandProxyFile) Deactivate() {
	var activeWb models.Active

	// reset the active key to false
	wb.active = false
	// Record the activation
	wb.deactivatedTime = base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02 15:04:05")))

	activeWb.ID = fmt.Sprint(wb.id)
	activeWb.Active = wb.active
	select {
	case wb.activeChan <- &activeWb:
		// values are being read from r.Resolver.wbChan
		fmt.Println("wb.activeChan: inserted active state")
	default:
		// no subscribers, wb not in channel
		fmt.Println("wb.activeChan: active chan created, not inserted")
	}
}

// Activate A Wristband
func (wb *wristbandProxyFile) Activate() {
	var activeWb models.Active

	// reset the active key to true
	wb.active = true
	// Record the activation
	wb.activatedTime = base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02 15:04:05")))

	activeWb.ID = fmt.Sprint(wb.id)
	activeWb.Active = wb.active
	select {
	case wb.activeChan <- &activeWb:
		// values are being read from r.Resolver.wbChan
		fmt.Println("wb.activeChan: inserted active state")
	default:
		// no subscribers, wb not in channel
		fmt.Println("wb.activeChan: active chan created, not inserted")
	}
}

///////////////////////////////////
// IsActive
func (wb *wristbandProxyFile) IsActive() bool {
	return wb.active
}

func (wb *wristbandProxyFile) GetID() string {
	return fmt.Sprintf("%d", wb.id)
}
