package core

// import (
// 	"cerberus-security-laboratories/des-wristband-ui/internal/gql/models"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"testing"
// 	"time"
// )

// func Test_readSingleJSONDataEntry(t *testing.T) {
// 	type args struct {
// 		fileName string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *models.WristbandData
// 		wantErr bool
// 	}{
// 		{
// 			"Read correct data",
// 			args{
// 				"./unittest_data/Test_readSingleJSONDataEntry/valid.json",
// 			},
// 			&models.WristbandData{
// 				ID: "",
// 				// Wristband: &models.Wristband{
// 				// 	ID:          "aWbyyzBaYH",
// 				// 	Msid:        0x0499141551,
// 				// 	Type:        0x7164083171,
// 				// 	TypeVer:     0x1,
// 				// 	Key:         "thxOC",
// 				// 	Tic:         "o8vI5irq2s",
// 				// 	Active:      true,
// 				// 	Activated:   nil,
// 				// 	Deactivated: nil,
// 				// 	Data:        nil,
// 				// 	Name:        nil,
// 				// 	OnOxygen:    false,
// 				// 	Pregnant:    false,
// 				// 	Child:       false,
// 				// },
// 				Time: "",
// 				SensorData: &models.WristbandSensorData{
// 					Respiration:   82,
// 					Sp02:          1,
// 					Pulse:         407,
// 					Temperature:   53,
// 					BloodPressure: nil,
// 					Motion:        true,
// 					Proximity:     true,
// 				},
// 				News2: &models.News2{
// 					Respiration:   82,
// 					Sp02:          1,
// 					OnOxygen:      false,
// 					Pulse:         407,
// 					Temperature:   53,
// 					BloodPressure: nil,
// 					Motion:        true,
// 					Overall:       1499,
// 				},
// 				// BatteryLevel: 69,
// 				Location: []*models.BridgeSignal{
// 					&models.BridgeSignal{
// 						Bridge: &models.Bridge{
// 							ID:      "aWbyyzBaYH",
// 							Msid:    499141551,
// 							Type:    7164083171,
// 							TypeVer: 0x1,
// 							Key:     "thxOC",
// 							Tic:     "o8vI5irq2s",
// 						},
// 						Signal: -35,
// 					},
// 					&models.BridgeSignal{
// 						Bridge: &models.Bridge{
// 							ID:      "HwbTHwLhd",
// 							Msid:    783652924,
// 							Type:    7164083171,
// 							TypeVer: 0x1,
// 							Key:     "sfkhj",
// 							Tic:     "jkfsjkfjxd",
// 						},
// 						Signal: -40,
// 					},
// 				},
// 			},
// 			false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			file, err := os.Open(tt.args.fileName)
// 			if err != nil {
// 				t.Errorf("could not open test data file %s: %v", tt.args.fileName, err)
// 			}
// 			defer file.Close()
// 			decoder := json.NewDecoder(file)
// 			// Read the array open bracket
// 			decoder.Token()

// 			got, err := readSingleJSONDataEntry(decoder)

// 			fmt.Printf("Got: %s", got.PrettyPrint())
// 			fmt.Printf("Exp: %s", tt.want.PrettyPrint())

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("readSingleJSONDataEntry() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			err = got.EqualValues(tt.want)
// 			if err != nil {
// 				t.Errorf("readSingleJSONDataEntry() = %v", err)
// 			}
// 		})
// 	}
// }

// func TestNewWristbandProxyFile(t *testing.T) {
// 	type args struct {
// 		id          int
// 		tickPeriod  int
// 		name        string
// 		dateOfBirth string
// 		onOxygen    bool
// 		pregnant    bool
// 		child       bool
// 		department 	string
// 		dataDecoder *json.Decoder
// 		fileName    string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    WristbandProxy
// 		wantErr bool
// 	}{
// 		{
// 			"Basic creation and file reading",
// 			args{
// 				id:          946482362,
// 				tickPeriod:  1,
// 				name:        "Andy",
// 				onOxygen:    false,
// 				pregnant:    false,
// 				child:       false,
// 				dataDecoder: nil,
// 				// fileName:    "../../tests/unit/core/TestNewWristbandProxyFile/basic.json",
// 				fileName: "./unittest_data/TestNewWristbandProxyFile/basic.json",
// 			},
// 			&wristbandProxyFile{
// 				id:          946482362,
// 				tickPeriod:  1,
// 				name:        "Andy",
// 				onOxygen:    false,
// 				pregnant:    false,
// 				child:       false,
// 				dataDecoder: nil,
// 				data: []*models.WristbandData{
// 					{
// 						Time: "",
// 						SensorData: &models.WristbandSensorData{
// 							Respiration:   82,
// 							Sp02:          1,
// 							Pulse:         407,
// 							Temperature:   53,
// 							BloodPressure: nil,
// 							Motion:        true,
// 							Proximity:     true,
// 						},
// 						News2: &models.News2{
// 							Respiration:   82,
// 							Sp02:          1,
// 							OnOxygen:      false,
// 							Pulse:         407,
// 							Temperature:   53,
// 							BloodPressure: nil,
// 							Motion:        true,
// 							Overall:       1499,
// 						},
// 						// BatteryLevel: 69,
// 						Location: []*models.BridgeSignal{
// 							&models.BridgeSignal{
// 								Bridge: &models.Bridge{
// 									ID:      "aWbyyzBaYH",
// 									Msid:    499141551,
// 									Type:    7164083171,
// 									TypeVer: 0x1,
// 									Key:     "thxOC",
// 									Tic:     "o8vI5irq2s",
// 								},
// 								Signal: -35,
// 							},
// 							&models.BridgeSignal{
// 								Bridge: &models.Bridge{
// 									ID:      "HwbTHwLhd",
// 									Msid:    783652924,
// 									Type:    7164083171,
// 									TypeVer: 0x1,
// 									Key:     "sfkhj",
// 									Tic:     "jkfsjkfjxd",
// 								},
// 								Signal: -40,
// 							},
// 						},
// 					},
// 					{

// 						Time: "",
// 						SensorData: &models.WristbandSensorData{
// 							Respiration:   79,
// 							Sp02:          5,
// 							Pulse:         432,
// 							Temperature:   43,
// 							BloodPressure: nil,
// 							Motion:        true,
// 							Proximity:     true,
// 						},
// 						News2: &models.News2{
// 							Respiration:   72,
// 							Sp02:          13,
// 							OnOxygen:      false,
// 							Pulse:         65,
// 							Temperature:   36,
// 							BloodPressure: nil,
// 							Motion:        true,
// 							Overall:       96,
// 						},
// 						// BatteryLevel: 69,
// 						Location: []*models.BridgeSignal{
// 							&models.BridgeSignal{
// 								Bridge: &models.Bridge{
// 									ID:      "yAJifjeHwnJ",
// 									Msid:    943573498,
// 									Type:    7164083171,
// 									TypeVer: 0x1,
// 									Key:     "thxOC",
// 									Tic:     "o8vI5irq2s",
// 								},
// 								Signal: -35,
// 							},
// 							&models.BridgeSignal{
// 								Bridge: &models.Bridge{
// 									ID:      "wALFJdkjWhW",
// 									Msid:    783652924,
// 									Type:    7164083171,
// 									TypeVer: 0x1,
// 									Key:     "sfkhj",
// 									Tic:     "jkfsjkfjxd",
// 								},
// 								Signal: -40,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			file, err := os.Open(tt.args.fileName)
// 			if err != nil {
// 				t.Errorf("could not open test data file %s: %v", tt.args.fileName, err)
// 			}
// 			defer file.Close()
// 			decoder := json.NewDecoder(file)
// 			// Read the array open bracket
// 			decoder.Token()

// 			got, err := NewWristbandProxyFile(tt.args.id, tt.args.tickPeriod, tt.args.name, tt.args.onOxygen, tt.args.dateOfBirth, tt.args.pregnant, tt.args.child, tt.args.department, decoder)
// 			// Give some time for go routine and ticker to run
// 			time.Sleep(100 * time.Millisecond)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("NewWristbandProxyFile() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				// Collect expected data
// 				expData, err := tt.want.GetData(0)
// 				if err != nil {
// 					t.Errorf("Want data could not be accessed")
// 				}
// 				// Collect generated data
// 				gotData, err := got.GetData(0)
// 				if err != nil {
// 					t.Errorf("Got data could not be accessed")
// 				}
// 				// Check there is the same amount as generated
// 				if len(expData) != len(gotData) {
// 					t.Errorf("Want data and got data lengths do not match")
// 				}

// 				for i, d := range expData {
// 					err = d.EqualValues(gotData[i])
// 					if err != nil {
// 						t.Errorf("WristbandData[%d] does not match: %v", i, err)
// 					}
// 				}
// 			}
// 		})
// 	}
// }
