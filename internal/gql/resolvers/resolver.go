//go:generate go run ../../../scripts/gqlgen.go

package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"cerberus-security-laboratories/des-wristband-ui/internal/core"
	"cerberus-security-laboratories/des-wristband-ui/internal/gql/models"

	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/apex/log"
)

var resolver *Resolver

// Goroutines
func updateWristband(id int) {
	wristbands := resolver.GetWristbandCore()

	// This will get triggered whenever a wristband level channel got sent a data
	// This will lead to a more up-to-date all bands summary
	go updateLevelSummary(wristbands[id].GetLevelChan())

	// This will get triggered when an overall news 2 score channel got sent a data
	// This will then run a goroutine that sort important bands
	go updateImportantBands(wristbands[id].GetNews2Chan())
}

// SUMMARY CHANNEL
func updateLevelSummary(ch chan *models.Level, err error) {
	if err != nil {
		return
	}
	for range ch {
		// Get The Wristband Core Array
		wristbands := resolver.GetWristbandCore()

		high := 0
		medium := 0
		lowMedium := 0
		low := 0
		other := 0

		// The Reason to loop through all bands again is to get ALL bands summary correct
		for i, wb := range wristbands {
			level, levelErr := wb.GetLevel()
			fmt.Println(i, level)

			if levelErr != nil {
				return
			}

			if level != nil {
				// Do The Sum
				// values are being read from wb.GetLevelChan()
				if level.Text == "high" {
					fmt.Println("Individual Wristband's levels: HIGH")
					high++
				} else if level.Text == "medium" {
					fmt.Println("Individual Wristband's levels: MEDIUM")
					medium++
				} else if level.Text == "low-medium" {
					fmt.Println("Individual Wristband's levels: LOW-MEDIUM")
					lowMedium++
				} else if level.Text == "low" {
					fmt.Println("Individual Wristband's levels: LOW")
					low++
				} else if level.Text == "error" {
					fmt.Println("Individual Wristband's levels: OTHER")
					other++
				}
			}
		}

		// Summary
		var summary models.Summary

		dataSumChan := resolver.GetWristbandDataSumChan()

		// Assign to The Summary Model
		summary.High = high
		summary.Medium = medium
		summary.LowMedium = lowMedium
		summary.Low = low
		summary.Other = other

		// Send The Summary Data to The Channel
		select {
		case dataSumChan <- &summary:
			// values are being  read from r.Resolver.dataSumChan
			fmt.Println("r.Resolver.dataSumChan: inserted summary")
		default:
			// no subscribers, wb not in channel
			fmt.Println("r.Resolver.dataSumChan: summary created, not inserted")
		}
	}
}

// IMPORTANT CHANNEL
func updateImportantBands(ch chan *models.News2, err error) {
	if err != nil {
		return
	}
	for range ch {
		var importantBands []*models.Wristband
		// Get The Wristband Core Array
		wristbands := resolver.GetWristbandCore()

		for i, wb := range wristbands {
			// Get Each Individual Band
			mwb, err := wb.Get()
			if err != nil {
				return
			}

			// Get Their News 2 Score
			wnews, newsErr := wb.GetNews2()
			if newsErr != nil {
				return
			}

			// Get Their Latest Data for Proximity (Skin Contact) && Get Their Battery Level
			wbd_latest, latestErr := wb.GetLatestData()
			if latestErr != nil {
				return
			}

			// Check A Screnario Where A Wristband Has News 2 Score, but A Recently Added Doesn't Have One Yet
			if wnews != nil && wbd_latest != nil {
				// Filter out the low news score (0 - 4)
				if wnews.Overall >= 0 && wnews.Overall <= 4 {
					fmt.Println("LOW OVERALL SCORE ON WRISTBAND OF INDEX ", i)
					// Otherwise, if overall score is larger than 4
					// or, baterry level is below 20
					// or, proximity (skin contact) is false
				} else if wnews.Overall > 4 || wbd_latest.BatteryLevel < 20 || !wbd_latest.SensorData.Proximity {
					importantBands = append(importantBands, &mwb)
				}
			}
		}

		// After Collecting Important Bands
		importantChan := resolver.GetImportantBandsChan()
		// Send The Important Bands Array to The Channel
		select {
		case importantChan <- importantBands:
			// values are being  read from importantChan
			fmt.Println("importantChan: inserted important bands")
		default:
			// no subscribers, wb not in channel
			fmt.Println("importantChan: important bands created, not inserted")
		}
	}
}

///////////////////////////////////////////////
// SHOWCASE ONLY
// Wristband Input JSON files
type wristbandInputs struct {
	Tic         string `json:"tic"`
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
	OnOxygen    bool   `json:"onOxygen"`
	Pregnant    bool   `json:"pregnant"`
	Child       bool   `json:"child"`
	Department  string `json:"department"`
}

func showcaseOnly(finalFolder string) {
	// Get Current Working Directory
	curr_wd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	curr_wd = filepath.Join(curr_wd, "internal", "core", "unittest_data", "Test_Showcase_AddWristbandInputs", finalFolder)

	fileName := filepath.Join(curr_wd, "input.json")
	fmt.Println(fileName)
	fd, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	// defer fd.Close()
	// Load file into a JSON decoder
	decoder := json.NewDecoder(fd)
	// Read the array open bracket
	decoder.Token()

	var numFiles int = 10

	if finalFolder == "Nurse" {
		numFiles = 11
	}

	for i := 0; i < numFiles; i++ {
		if !decoder.More() {
			return
		}
		// Assign Input Fields to models.AddWristbandInput obj
		var wbInputs wristbandInputs
		decoder.Decode(&wbInputs)
		var input models.AddWristbandInput
		input.Name = &wbInputs.Name
		input.DateOfBirth = wbInputs.DateOfBirth
		input.OnOxygen = wbInputs.OnOxygen
		input.Pregnant = wbInputs.Pregnant
		input.Child = wbInputs.Child
		input.Department = wbInputs.Department
		// Call The Mutation Resolver 10 times
		resolver.Mutation().AddWristband(context.TODO(), input)
	}
}

func ResolverInitialisation(dataPath *string, dataPrefix *string, tickPeriod *int, finalFolder string) *Resolver {
	// Instantiate Wristband Factory
	wf, err := core.NewWristbandFactory(*dataPath, *dataPrefix, *tickPeriod)
	if err != nil {
		panic(fmt.Errorf("Could not create Wristband Factory: %v", err))
	}
	// Instantiate Wristband Channel
	wbChan := make(chan *models.Wristband, 1)

	// Instantiate Wristband Data Summary Channel
	dataSumChan := make(chan *models.Summary, 1)

	// Instantiate Important Wristband Channel
	importantWristbandChan := make(chan []*models.Wristband, 1)

	resolver = &Resolver{
		wbChan:                 wbChan,
		dataSumChan:            dataSumChan,
		importantWristbandChan: importantWristbandChan,
		WristbandFactory:       *wf,
	}

	// Trigger A Goroutine Which Listens on The Summary Channel After A Wristband is Added
	// updateWristband()
	// Showcase only
	showcaseOnly(finalFolder)

	return resolver
}

// Resolver contains all the state of the GQL resolvers
type Resolver struct {
	// Logging
	logger *log.Entry
	// sync
	mu sync.Mutex
	// Wristband Factory instance
	WristbandFactory core.WristbandFactory
	// Wristband Channel
	wbChan chan *models.Wristband
	// Wristband Data Summary Channel
	dataSumChan chan *models.Summary
	// Important Bands Channel
	importantWristbandChan chan []*models.Wristband
	// Array of Wristbands
	wristband []core.WristbandProxy
}

func (r *Resolver) GetWristbandCore() []core.WristbandProxy {
	return r.wristband
}

func (r *Resolver) GetImportantBandsChan() chan []*models.Wristband {
	return r.importantWristbandChan
}

func (r *Resolver) GetWristbandDataSumChan() chan *models.Summary {
	return r.dataSumChan
}
