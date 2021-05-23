package core

import (
	"cerberus-security-laboratories/des-wristband-ui/internal/gql/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// WristbandFactory is used to create new WristbandProxy objects
type WristbandFactory struct {
	tickPeriod int    // period between wristband updates
	wbCount    int    // total number of created wristbands
	dataPath   string // Directory path to data files
	filePrefix string // file prefix for JSON config files
}

// NewWristbandFactory creates a new instance of this struct
func NewWristbandFactory(dataPath string, filePrefix string, tickPeriod int) (*WristbandFactory, error) {
	log.Printf("Creating new Basic Wristband Factory")

	// Check dataPath is valid
	_, err := os.Stat(dataPath)
	if err != nil {
		return nil, fmt.Errorf("dataPath [%s] does not exist: %v", dataPath, err)
	}

	wf := WristbandFactory{
		tickPeriod: tickPeriod,
		wbCount:    0,
		dataPath:   dataPath,
		filePrefix: filePrefix,
	}

	log.Printf("Created new Basic Wristband Factory")
	return &wf, nil
}

// NewWristband creates a new wristband instance and genrates all the values
func (wf *WristbandFactory) NewWristband(input *models.AddWristbandInput) (WristbandProxy, error) {
	log.Println("WristbandFactory.NewWristband called")

	wf.wbCount++

	var name string = "N/A"
	if input.Name != nil {
		name = *input.Name
	}

	// If dataPath and filePrefix are blank then create a
	// basic wristband that randomly generates "in range" data.
	if wf.dataPath == "" && wf.filePrefix == "" {
		wp, err := NewWristbandProxyBasic(
			wf.wbCount,
			wf.tickPeriod,
			name,
			input.OnOxygen,
			input.DateOfBirth,
			input.Pregnant,
			input.Child,
			input.Department,
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Created new Basic Wristband Proxy #%d", wf.wbCount)
		return wp, nil
	}

	// Otherwise create a wristband that reads data from a file
	fileName := filepath.Join(wf.dataPath, wf.filePrefix+fmt.Sprintf("%03d", wf.wbCount)+".json")
	fd, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, fmt.Errorf("Tried to open wristband data file: %e", err)
	}
	// defer fd.Close()
	// Load file into a JSON decoder
	decoder := json.NewDecoder(fd)
	// Read the array open bracket
	decoder.Token()

	wp, err := NewWristbandProxyFile(
		wf.wbCount,
		wf.tickPeriod,
		name,
		input.OnOxygen,
		input.DateOfBirth,
		input.Pregnant,
		input.Child,
		input.Department,
		decoder,
	)
	if err != nil {
		return nil, err
	}

	log.Printf("Created new Wristband Proxy #%d with data source %s", wf.wbCount, fileName)

	return wp, nil
}

// AddWristband creates a new wristband instance and imports all the values
//TODO currently this ignores the TIC and creates a new one...
func (wf *WristbandFactory) AddWristband(input *models.AddWristbandInput) (WristbandProxy, error) {
	wf.wbCount++
	var name string = "N/A"
	if input.Name != nil {
		name = *input.Name
	}
	wp, err := NewWristbandProxyBasic(
		wf.wbCount,
		wf.tickPeriod,
		name,
		input.OnOxygen,
		input.DateOfBirth,
		input.Pregnant,
		input.Child,
		input.Department,
	)
	// wp, err := wf.NewWristband()
	if err != nil {
		return nil, err
	}
	log.Printf("Added new Basic Wristband Proxy #%d", wf.wbCount)

	return wp, nil
}
