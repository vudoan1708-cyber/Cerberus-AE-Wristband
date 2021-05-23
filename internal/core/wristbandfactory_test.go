package core

// import (
// 	"cerberus-security-laboratories/des-wristband-ui/internal/gql/models"
// 	"fmt"
// 	"math/rand"
// 	"strings"
// 	"testing"
// )

// // ErrorContains checks if the error message in out contains the text in
// // want.
// //
// // This is safe when out is nil. Use an empty string for want if you want to
// // test that err is nil.
// func ErrorContains(out error, want string) bool {
// 	if out == nil {
// 		return want == ""
// 	}
// 	if want == "" {
// 		return false
// 	}
// 	return strings.Contains(out.Error(), want)
// }

// func TestWristbandFactory_NewWristband(t *testing.T) {
// 	type fields struct {
// 		tickPeriod int
// 		wbCount    int
// 		dataPath   string
// 		filePrefix string
// 	}
// 	name1 := "Vlbzg Aicmrajwwh"
// 	tests := []struct {
// 		name    string
// 		seed    int64
// 		fields  fields
// 		want    *models.Wristband
// 		wantErr string
// 	}{
// 		{
// 			"BasicWristbandProxy creation",
// 			1,
// 			fields{
// 				1,
// 				0,
// 				"",
// 				"",
// 			},
// 			&models.Wristband{
// 				ID:          "1",
// 				Msid:        196609,
// 				Type:        1214578689,
// 				TypeVer:     1214578689,
// 				Key:         "9IU9e2Xk01IB1ye8KZJBHA6seBDzWPv8zA1jTZw7olQ=",
// 				Tic:         "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
// 				Active:      true,
// 				Activated:   nil,
// 				Deactivated: nil,
// 				Name:        nil,
// 				OnOxygen:    false,
// 				Pregnant:    false,
// 				Child:       false,
// 			},
// 			"",
// 		},
// 		{
// 			"FileWristbandProxy creation with valid file",
// 			1,
// 			fields{
// 				1,
// 				0,
// 				"./unittest_data/TestWristbandFactory_NewWristband",
// 				"wbData_",
// 			},
// 			&models.Wristband{
// 				ID:          "1",
// 				Msid:        196609,
// 				Type:        1214578689,
// 				TypeVer:     1214578689,
// 				Key:         "9IU9e2Xk01IB1ye8KZJBHA6seBDzWPv8zA1jTZw7olQ=",
// 				Tic:         "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
// 				Active:      true,
// 				Activated:   nil,
// 				Deactivated: nil,
// 				Name:        &name1,
// 				OnOxygen:    false,
// 				Pregnant:    false,
// 				Child:       false,
// 			},
// 			"",
// 		},
// 		{
// 			"FileWristbandProxy creation with incorrect data directory",
// 			1,
// 			fields{
// 				1,
// 				0,
// 				"./blah/TestWristbandFactory_NewWristband",
// 				"wbData_",
// 			},
// 			nil,
// 			"Tried to open wristband data file",
// 		},
// 		{
// 			"FileWristbandProxy creation with incorrect file prefix",
// 			1,
// 			fields{
// 				1,
// 				0,
// 				"./unittest_data/TestWristbandFactory_NewWristband",
// 				"blah_",
// 			},
// 			nil,
// 			"Tried to open wristband data file",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rand.Seed(tt.seed)
// 			wf := &WristbandFactory{
// 				tickPeriod: tt.fields.tickPeriod,
// 				wbCount:    tt.fields.wbCount,
// 				dataPath:   tt.fields.dataPath,
// 				filePrefix: tt.fields.filePrefix,
// 			}

// 			// Create a wristband
// 			got, err := wf.NewWristband()

// 			if tt.wantErr != "" {
// 				if !ErrorContains(err, tt.wantErr) {
// 					// Not the expected error
// 					t.Errorf("WristbandFactory.NewWristband() error = %v, wantErr %v", err, tt.wantErr)
// 				}
// 			} else {
// 				wb, err := got.Get()
// 				if err != nil {
// 					t.Errorf("WristbandProxy.Get() error: %v", err)
// 				}
// 				err = wb.EqualValues(tt.want)
// 				if err != nil {
// 					t.Errorf("WristbandFactory.NewWristband() Wristband mismatch: %v", err)
// 					fmt.Printf("%s\n", wb.PrettyPrint())
// 				}
// 			}
// 		})
// 	}
// }
