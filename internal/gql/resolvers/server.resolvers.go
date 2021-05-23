package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"cerberus-security-laboratories/des-wristband-ui/internal/gql"
	"cerberus-security-laboratories/des-wristband-ui/internal/gql/models"
	"context"
	"fmt"
	"log"
	"strconv"
)

func (r *mutationResolver) AddWristband(ctx context.Context, input models.AddWristbandInput) (*models.Wristband, error) {
	log.Print("[GQL] AddWristband")
	wb, err := r.Resolver.WristbandFactory.NewWristband(&input)
	if err != nil {
		log.Printf("Error creating new Wristband: %e\n", err)
		return nil, err
	}
	r.Resolver.wristband = append(r.Resolver.wristband, wb)
	// Get the model representation
	mWb, err := wb.Get()
	if err != nil {
		log.Printf("Error getting new Wristband state: %e\n", err)
		return nil, err
	}

	select {
	case r.Resolver.wbChan <- &mWb:
		// values are being read from r.Resolver.wbChan
		fmt.Println("r.Resolver.wbChan: inserted wb")
	default:
		// no subscribers, wb not in channel
		fmt.Println("r.Resolver.wbChan: wb created, not inserted")
	}

	// Trigger A Goroutine Which Listens on The Summary Channel After A Wristband is Added
	id, idErr := strconv.Atoi(mWb.ID)
	if idErr != nil {
		log.Printf("Error converting new Wristband ID to Integer: %e\n", idErr)
		return nil, err
	}
	updateWristband(id - 1)

	return &mWb, nil
}

func (r *mutationResolver) DeactivateWristband(ctx context.Context, input models.DeactivateWristbandInput) (*models.Wristband, error) {
	log.Print("[GQL] Deactivate A Wristband")

	var mWb models.Wristband
	var err error
	// Input for Wristband Deactivation is an ID of type string
	// Loop through array
	for _, wb := range r.Resolver.wristband {
		if wb.GetID() == input.ID {
			// Get the model representation
			mWb, err = wb.Get()
			if err != nil {
				log.Printf("Error getting a Wristband state: %e\n", err)
				return nil, err
			}
			wb.Deactivate()
		}
	}

	return &mWb, nil
}

func (r *mutationResolver) ReactivateWristband(ctx context.Context, input models.DeactivateWristbandInput) (*models.Wristband, error) {
	log.Print("[GQL] Reactivate A Wristband")

	var mWb models.Wristband
	var err error
	// Input for Wristband Deactivation is an ID of type string
	// Loop through array
	for _, wb := range r.Resolver.wristband {
		if wb.GetID() == input.ID {
			// Check if the wristband is NOT active
			if !wb.IsActive() {
				// Get the model representation
				mWb, err = wb.Get()
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
				wb.Activate()
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Active")
			}
		}
	}

	return &mWb, nil
}

func (r *mutationResolver) ReassignNewWristband(ctx context.Context, oldWristband models.DeactivateWristbandInput) (*models.Wristband, error) {
	log.Print("[GQL] Reassign A New Wristband")

	var mWbInput models.AddWristbandInput

	var mWb *models.Wristband
	var err error

	// Loop through array
	for _, wb := range r.Resolver.wristband {
		// String Match for Wristband IDs
		if wb.GetID() == oldWristband.ID {
			// Check if the wristband is active
			if wb.IsActive() {
				// Deactivate The Wristband
				deactivatedWristband, deactivatedErr := r.DeactivateWristband(ctx, oldWristband)
				if deactivatedErr != nil {
					log.Printf("Error getting a Wristband state: %e\n", deactivatedErr)
					return nil, deactivatedErr
				}

				// Extract Old Information from the Deactivated Wristband to a Placeholder
				mWbInput.Name = deactivatedWristband.Name
				mWbInput.OnOxygen = deactivatedWristband.OnOxygen
				mWbInput.DateOfBirth = deactivatedWristband.DateOfBirth
				mWbInput.Pregnant = deactivatedWristband.Pregnant
				mWbInput.Child = deactivatedWristband.Child
				mWbInput.Department = deactivatedWristband.Department

				// Add New Wristband with old information
				mWb, err = r.AddWristband(ctx, mWbInput)
				// mWb, err = wb.Get()
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Deactivated")
			}
		}
	}
	return mWb, nil
}

func (r *mutationResolver) ResetName(ctx context.Context, id string, value string) (*models.Wristband, error) {
	log.Print("[GQL] Reset A Wristband's Name")
	var mWb models.Wristband
	var err error

	for _, wb := range r.Resolver.wristband {
		// Find the selectec band
		if wb.GetID() == id {
			// Check if the wristband is active
			if wb.IsActive() {
				// Reset its name
				wb.SetName(value)
				// Get the model after the mutation and
				// return it with the updated value
				mWb, err = wb.Get()
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Deactivated")
			}
		}
	}
	return &mWb, nil
}

func (r *mutationResolver) ResetOnOxygen(ctx context.Context, id string, value bool) (*models.Wristband, error) {
	log.Print("[GQL] Reset A Wristband's OnOxygen Info")
	var mWb models.Wristband
	var err error

	for _, wb := range r.Resolver.wristband {
		// Find the selectec band
		if wb.GetID() == id {
			// Check if the wristband is active
			if wb.IsActive() {
				// Reset its onoxygen value
				wb.SetOnOxygen(value)
				// Get the model after the mutation and
				// return it with the updated value
				mWb, err = wb.Get()
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Deactivated")
			}
		}
	}
	return &mWb, nil
}

func (r *mutationResolver) ResetPregnant(ctx context.Context, id string, value bool) (*models.Wristband, error) {
	log.Print("[GQL] Reset A Wristband's Pregnant Info")
	var mWb models.Wristband
	var err error

	for _, wb := range r.Resolver.wristband {
		// Find the selectec band
		if wb.GetID() == id {
			// Check if the wristband is active
			if wb.IsActive() {
				// Reset its prenant value
				wb.SetPregnant(value)
				// Get the model after the mutation and
				// return it with the updated value
				mWb, err = wb.Get()
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Deactivated")
			}
		}
	}
	return &mWb, nil
}

func (r *mutationResolver) ResetChild(ctx context.Context, id string, value bool) (*models.Wristband, error) {
	log.Print("[GQL] Reset A Wristband's Child Info")
	var mWb models.Wristband
	var err error

	for _, wb := range r.Resolver.wristband {
		// Find the selectec band
		if wb.GetID() == id {
			// Check if the wristband is active
			if wb.IsActive() {
				// Reset its child value
				wb.SetChild(value)
				// Get the model after the mutation and
				// return it with the updated value
				mWb, err = wb.Get()
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Deactivated")
			}
		}
	}
	return &mWb, nil
}

func (r *mutationResolver) ResetDepartment(ctx context.Context, id string, value string) (*models.Wristband, error) {
	log.Print("[GQL] Reset A Wristband's Department Info")
	var mWb models.Wristband
	var err error

	for _, wb := range r.Resolver.wristband {
		// Find the selected band
		if wb.GetID() == id {
			// Check if the wristband is active
			if wb.IsActive() {
				// Reset its department value
				wb.SetDepartment(value)
				// Get the model after the mutation and
				// return it with the updated value
				mWb, err = wb.Get()
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Deactivated")
			}
		}
	}
	return &mWb, nil
}

func (r *mutationResolver) ResetMultipleFields(ctx context.Context, id string, options []string, values []string) (*models.Wristband, error) {
	log.Print("[GQL] Reset A Wristband's Info")
	var mWb models.Wristband
	var err error

	for _, wb := range r.Resolver.wristband {
		// Find the selectec band
		if wb.GetID() == id {
			// Check if the wristband is active
			if wb.IsActive() {
				for index, option := range options {

					// Check If The Value Is A Boolean In A String Format
					var booleanParsedValue bool
					var booleanParsedError error
					if values[index] == "true" || values[index] == "false" {
						// Convert The Stringified Boolean Values Into Proper Boolean Values
						booleanParsedValue, booleanParsedError = strconv.ParseBool(values[index])
						if booleanParsedError != nil {
							log.Printf("Error changing boolean fields: %e\n", err)
							return nil, err
						}
					}
					// Reset Fields
					// The Options Argument is to identify whic fields the client wants to make changes to
					// always starting from the Name field -> the Department field (1 -> 5)
					if option == "Name" {
						// Reset its name
						wb.SetName(values[index])
						// Get the model after the mutation and
						// return it with the updated value
						mWb, err = wb.Get()
						if err != nil {
							log.Printf("Error setting a Wristband state: %e\n", err)
							return nil, err
						}
					} else if option == "OnOxygen" {
						// Reset its onoxygen value
						wb.SetOnOxygen(booleanParsedValue)
						// Get the model after the mutation and
						// return it with the updated value
						mWb, err = wb.Get()
						if err != nil {
							log.Printf("Error setting a Wristband state: %e\n", err)
							return nil, err
						}
					} else if option == "Pregnant" {
						// Reset its prenant value
						wb.SetPregnant(booleanParsedValue)
						// Get the model after the mutation and
						// return it with the updated value
						mWb, err = wb.Get()
						if err != nil {
							log.Printf("Error setting a Wristband state: %e\n", err)
							return nil, err
						}
					} else if option == "Child" {
						// Reset its child value
						wb.SetChild(booleanParsedValue)
						// Get the model after the mutation and
						// return it with the updated value
						mWb, err = wb.Get()
						if err != nil {
							log.Printf("Error setting a Wristband state: %e\n", err)
							return nil, err
						}
					} else if option == "Department" {
						// Reset its child value
						wb.SetDepartment(values[index])
						// Get the model after the mutation and
						// return it with the updated value
						mWb, err = wb.Get()
						if err != nil {
							log.Printf("Error setting a Wristband state: %e\n", err)
							return nil, err
						}
					}
				}
			} else {
				return nil, fmt.Errorf("The Selected Wristband Has Already Been Deactivated")
			}
		}
	}

	return &mWb, nil
}

func (r *queryResolver) Node(ctx context.Context, id string) (models.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetWristband(ctx context.Context, id *string) (*models.Wristband, error) {
	log.Print("[GQL] Get A Wristband")
	var mWb models.Wristband
	var err error

	for _, wb := range r.Resolver.wristband {
		if wb.GetID() == *id {
			// Get the model representation
			mWb, err = wb.Get()
			if err != nil {
				log.Printf("Error getting a Wristband state: %e\n", err)
				return nil, err
			}
		}
	}

	return &mWb, nil
}

func (r *queryResolver) GetBridge(ctx context.Context, id *string) (*models.Bridge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetGateway(ctx context.Context, id *string) (*models.Gateway, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetWristbandData(ctx context.Context, id *string, how *string) (*models.WristbandData, error) {
	log.Print("[GQL] Get A Wristband Data")
	var wbData models.WristbandData

	// Loop through the wristband core array
	for _, wb := range r.Resolver.wristband {
		if wb.GetID() == *id {

			// If The Latest Wristband Data Is Needed
			if *how == "latest" {
				// Get the latest data from the selected wristband
				mWbData, err := wb.GetLatestData()
				if err != nil {
					log.Printf("Error getting a Wristband Data: %e\n", err)
					return nil, err
				}
				wbData = *mWbData
			} else if *how == "first" {
				// Get the first data from the selected wristband
				mWbData, err := wb.GetFirstData()
				if err != nil {
					log.Printf("Error getting a Wristband Data: %e\n", err)
					return nil, err
				}
				wbData = *mWbData
			}
		}
	}
	// r.dataChannel <- &wbData
	return &wbData, nil
}

func (r *queryResolver) GetWristbands(ctx context.Context) ([]*models.Wristband, error) {
	log.Print("[GQL] Get All Wristbands")
	var mWbs []*models.Wristband

	for _, wb := range r.Resolver.wristband {
		mWb, err := wb.Get()
		if err != nil {
			log.Printf("Error getting all Wristband states: %e\n", err)
			return nil, err
		}
		mWbs = append(mWbs, &mWb)
	}

	return mWbs, nil
}

func (r *queryResolver) GetImportantBands(ctx context.Context) ([]*models.Wristband, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetMultipleWristbandData(ctx context.Context, id *string, howMany *int, start *int, end *int) ([]*models.WristbandData, error) {
	log.Print("[GQL] Get Multiple Wristband Data")
	var mWbd []*models.WristbandData
	var mWbdChild []*models.WristbandData

	var err error
	var errChild error

	// Loop through the wristband core array
	for _, wb := range r.Resolver.wristband {
		// If ID is specified in a query
		if *id != "" {
			// Get Data of that specific wristband via its ID
			if wb.GetID() == *id {
				if *howMany > -1 {
					mWbd, err = wb.GetData(*howMany)
				} else {
					mWbd, err = wb.GetDataBlock(*start, *end)
				}
				if err != nil {
					log.Printf("Error getting a Wristband state: %e\n", err)
					return nil, err
				}
			}
			// Otherwise, get every last data of each wristband
		} else {
			if *howMany > -1 {
				mWbdChild, errChild = wb.GetData(*howMany)
			} else {
				mWbdChild, errChild = wb.GetDataBlock(*start, *end)
			}
			if errChild != nil {
				log.Printf("Error getting a Wristband state: %e\n", errChild)
				return nil, errChild
			}
			mWbd = append(mWbd, mWbdChild...)
		}
	}

	return mWbd, nil
}

func (r *queryResolver) GetSummary(ctx context.Context) (*models.Summary, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAlert(ctx context.Context) (*models.Alert, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetBridges(ctx context.Context) ([]*models.Bridge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetGateways(ctx context.Context) ([]*models.Gateway, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) UpdateWristbandAdded(ctx context.Context) (<-chan *models.Wristband, error) {
	log.Printf("UpdateWristbandAdded() is called")
	return r.Resolver.wbChan, nil
}

func (r *subscriptionResolver) UpdateWristbandData(ctx context.Context, id *string) (<-chan *models.WristbandData, error) {
	// Channel for WristbandData Subscription
	var dataChannel chan *models.WristbandData
	var dataError error = nil

	wristbands := r.Resolver.wristband
	if len(wristbands) > 0 {
		// Loop through the wristband core array
		for _, wb := range wristbands {
			if wb.GetID() == *id {
				// Get the latest data from the selected wristband
				dataChannel, dataError = wb.GetWristbandDataChan()
				if dataError != nil {
					log.Printf("Error getting wristband data from subscription: %e\n", dataError)
					return nil, dataError
				}
			}
		}
	} else {
		panic(fmt.Errorf("You need to create a wristband first before there is any wristband data available to query"))
	}
	return dataChannel, dataError
}

func (r *subscriptionResolver) UpdateWristbandDataAlert(ctx context.Context, id *string) (<-chan *models.Alert, error) {
	// Channel for Alert
	var alertChan chan *models.Alert
	var alertError error = nil

	wristbands := r.Resolver.wristband
	if len(wristbands) > 0 {
		// Loop through the wristband core array
		for _, wb := range wristbands {
			if wb.GetID() == *id {
				// Get the latest data from the selected wristband
				alertChan, alertError = wb.GetAlertChan()
				if alertError != nil {
					log.Printf("Error getting wristband data alert from subscription: %e\n", alertError)
					return nil, alertError
				}
			}
		}
	} else {
		panic(fmt.Errorf("You need to create a wristband first before there is any wristband data available to query"))
	}
	return alertChan, alertError
}

func (r *subscriptionResolver) UpdateLevel(ctx context.Context, id *string) (<-chan *models.Level, error) {
	var levelChan chan *models.Level
	var levelChanErr error = nil

	var level *models.Level
	var levelErr error = nil

	wristbands := r.Resolver.wristband
	if len(wristbands) > 0 {
		// Loop through the wristband core array
		for _, wb := range wristbands {
			if wb.GetID() == *id {
				// Get Level Channel
				levelChan, levelChanErr = wb.GetLevelChan()
				// Get Level model
				level, levelErr = wb.GetLevel()
				fmt.Println(level)

				if levelChanErr != nil {
					log.Printf("Error getting wristband data summary from subscription: %e\n", levelChanErr)
					return nil, levelChanErr
				}
				if levelErr != nil {
					log.Printf("Error getting wristband data summary from subscription: %e\n", levelErr)
					return nil, levelErr
				}

				// Do The Sum
				if level != nil {
					select {
					case levelChan <- level:
						fmt.Println("level created, not inserted")
					default:
						// no subscribers, wb not in channel
						fmt.Println("level created, not inserted")
					}
				}
			}
		}
	}

	return levelChan, nil
}

func (r *subscriptionResolver) UpdateSummary(ctx context.Context) (<-chan *models.Summary, error) {
	return r.Resolver.dataSumChan, nil
}

func (r *subscriptionResolver) UpdateImportantBands(ctx context.Context) (<-chan []*models.Wristband, error) {
	return r.Resolver.importantWristbandChan, nil
}

func (r *subscriptionResolver) UpdateWristbandActive(ctx context.Context, id *string) (<-chan *models.Active, error) {
	// Channel for WristbandData Subscription
	var activeChan chan *models.Active
	var activeError error = nil

	wristbands := r.Resolver.wristband
	if len(wristbands) > 0 {
		// Loop through the wristband core array
		for _, wb := range wristbands {
			if wb.GetID() == *id {
				// Get the latest data from the selected wristband
				activeChan, activeError = wb.GetActiveChan()
				if activeError != nil {
					log.Printf("Error getting wristband data from subscription: %e\n", activeError)
					return nil, activeError
				}
			}
		}
	} else {
		panic(fmt.Errorf("You need to create a wristband first before there is any wristband data available to query"))
	}
	return activeChan, activeError
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

// Subscription returns gql.SubscriptionResolver implementation.
func (r *Resolver) Subscription() gql.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
