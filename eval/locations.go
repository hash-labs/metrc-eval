package eval

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hash-labs/metrc"
)

// LocationsResponse contains all data needed to verify the locations sheet.
type LocationsResponse struct {
	Create EvalRow `json:"create"` // first row
	Update EvalRow `json:"update"` // second row
	Get    EvalRow `json:"get"`    // third row
}

// EvalRow represents a row in the spreadsheet for evaluations.
type EvalRow struct {
	Code           int    `json:"code"`
	License        string `json:"license"`
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Request        string `json:"request"`
	BodyOrResponse string `json:"body_or_response"`
}

// Locations generates all data needed to verify the locations sheet.
func (e *EvalMetrc) Locations(license string) (LocationsResponse, error) {
	ts := time.Now().Format("2021.01.01 12:00:00")
	createName := fmt.Sprintf("Metrc Eval Location %s", ts)
	updateName := fmt.Sprintf("Metric Eval Location %s Updated", ts)

	clr, err := e.CreateLocation(license, createName)
	if err != nil {
		return LocationsResponse{}, fmt.Errorf("could not create location: %s", err)
	}

	id := clr.Id
	ulr, err := e.UpdateLocation(license, updateName, id)
	if err != nil {
		return LocationsResponse{}, fmt.Errorf("could not update location: %s", err)
	}

	glr, err := e.GetLocation(license, updateName, id)
	if err != nil {
		return LocationsResponse{}, fmt.Errorf("could not get final location: %s", err)
	}

	lr := LocationsResponse{
		Create: clr,
		Update: ulr,
		Get:    glr,
	}

	_, err = json.MarshalIndent(lr, "", "\t")
	if err != nil {
		return LocationsResponse{}, fmt.Errorf("could not marshal final struct: %s", err)
	}

	// Comment the below out for production / persistent results.
	// TODO: Make deletion configurable between testing and deploy.
	_, err = e.Metrc.DeleteLocationById(id, &license)
	if err != nil {
		return LocationsResponse{}, fmt.Errorf("could not delete location: %s", err)
	}

	return lr, nil
}

// CreateLocation creates a new location and returns its information.
// It corresponds to Step 1 in the Locations tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) CreateLocation(license string, name string) (EvalRow, error) {
	gotLocs, err := e.Metrc.GetLocationsActive(&license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not initially get active: %s", err)
	}

	// name := "Metrc Eval Location"
	inputLocs := []metrc.LocationPost{
		{
			Name:     name,
			TypeName: "Default Location type",
		},
	}

	_, err = e.Metrc.CreateLocations(inputLocs, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not create initial location: %s", err)
	}

	gotLocs, err = e.Metrc.GetLocationsActive(&license)
	var locId int
	var foundLocName bool
	for _, loc := range gotLocs {
		if loc.Name == name {
			locId = loc.Id
			foundLocName = true
			break
		}
	}

	if !foundLocName {
		return EvalRow{}, fmt.Errorf("could not get loc with matching name: %s", err)
	}

	endpoint := "locations/v1/create"
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)

	resp, err := json.MarshalIndent(inputLocs, "", "\t")
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not marshal input body")
	}
	response := string(resp)

	return EvalRow{
		Code:           200,
		License:        license,
		Id:             locId,
		Name:           name,
		Request:        request,
		BodyOrResponse: response,
	}, nil
}

// UpdateLocationResponse displays the information to verify updating a location.
// This corresponds to Step 2 in the Locations tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) UpdateLocation(license string, name string, id int) (EvalRow, error) {
	// name := "Metrc Eval Location Updated"
	updateLocs := []metrc.LocationPost{
		{
			Id:       id,
			Name:     name,
			TypeName: "Default Location type",
		},
	}

	_, err := e.Metrc.UpdateLocations(updateLocs, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not update locations: %s", err)
	}

	endpoint := "locations/v1/update"
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)
	resp, err := json.MarshalIndent(updateLocs, "", "\t")
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not marshal update body: %s", err)
	}
	response := string(resp)

	return EvalRow{
		Code:           200,
		License:        license,
		Id:             id,
		Name:           name,
		Request:        request,
		BodyOrResponse: response,
	}, nil
}

// GetLocationResponse displays the information to verify getting the final location.
// This corresponds to Step 3 in the Locations tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) GetLocation(license string, name string, id int) (EvalRow, error) {
	gotLoc, err := e.Metrc.GetLocationsById(id, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not get locs by id: %s", err)
	}

	endpoint := fmt.Sprintf("locations/v1/%d", id)
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)

	resp, err := json.MarshalIndent(gotLoc, "", "\t")
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not marshal got resp body: %s", err)
	}
	response := string(resp)

	return EvalRow{
		Code:           200,
		License:        license,
		Id:             id,
		Name:           name,
		Request:        request,
		BodyOrResponse: response,
	}, nil
}
