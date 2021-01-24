package eval

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hash-labs/metrc"
)

// StrainsResponse contains all data needed to verify the `strains` sheet.
type StrainsResponse struct {
	Create EvalRow `json:"create"` // first row
	Get    EvalRow `json:"get"`    // second row
}

func (e *EvalMetrc) Strains(license string) (StrainsResponse, error) {
	ts := time.Now().Format("2021.01.01 12:00:00")
	name := fmt.Sprintf("Metrc Strain Name %s", ts)

	cs, err := e.CreateStrain(license, name)
	if err != nil {
		return StrainsResponse{}, fmt.Errorf("could not create strain: %s", err)
	}

	id := cs.Id
	gs, err := e.GetStrain(license, name, id)
	if err != nil {
		return StrainsResponse{}, fmt.Errorf("could not get strain: %s", err)
	}

	sr := StrainsResponse{
		Create: cs,
		Get:    gs,
	}

	// Comment deletion out for production / persistent results.
	// TODO: Make deletion configurable between testing and deploy.
	_, err = e.Metrc.DeleteStrainById(id, &license)
	if err != nil {
		return StrainsResponse{}, fmt.Errorf("could not delete strain: %s", err)
	}
	return sr, nil
}

// CreateStrain creates a new strain and returns its information.
// It corresponds to Step 1 in the Strains tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) CreateStrain(license string, name string) (EvalRow, error) {
	// TODO: Understand why ThcLevel and CbdLevel aren't persisted.
	inputStrains := []metrc.Strain{
		{
			Name:             name,
			TestingStatus:    "None",
			ThcLevel:         .185,
			CbdLevel:         .275,
			IndicaPercentage: 25.0,
			SativaPercentage: 75.0,
		},
	}

	_, err := e.Metrc.CreateStrains(inputStrains, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not create strains: %s", err)
	}

	gotStrains, err := e.Metrc.GetStrainsActive(&license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not get active strains: %s", err)
	}

	var foundStrainName bool
	var gotStrain metrc.Strain
	for _, strain := range gotStrains {
		if strain.Name == name {
			foundStrainName = true
			gotStrain = strain
			break
		}
	}

	if !foundStrainName {
		return EvalRow{}, fmt.Errorf("could not get strain with matching name: %s", err)
	}

	endpoint := "strains/v1/create"
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)

	bodyBytes, err := json.MarshalIndent(gotStrain, "", "\t")
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not marshal input body")
	}
	body := string(bodyBytes)

	return EvalRow{
		Code:           200,
		License:        license,
		Id:             gotStrain.Id,
		Name:           name,
		Request:        request,
		BodyOrResponse: body,
	}, nil
}

// GetStrain gets the created strain and returns its information.
// It corresponds to Step 2 in the Strains tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) GetStrain(license string, name string, id int) (EvalRow, error) {
	// TODO: Implement
	gotStrain, err := e.Metrc.GetStrainsById(id, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not get strain by id: %s", err)
	}

	endpoint := fmt.Sprintf("strains/v1/%d", id)
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)

	resp, err := json.MarshalIndent(gotStrain, "", "\t")
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
