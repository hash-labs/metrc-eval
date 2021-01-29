package eval

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hash-labs/metrc"
)

// ItemsResponse contains the verification data for the items sheet.
type ItemsResponse struct {
	Create EvalRow `json:"create"` // first row
	Update EvalRow `json:"update"` // second row
	Get    EvalRow `json:"get"`    // third row
}

var itemName string = fmt.Sprintf("Metrc Item Name %s", time.Now().Format(timeLayoutFmt))

var defaultItem metrc.ItemPost = metrc.ItemPost{
	ItemCategory:                    "Capsule (weight)",
	Name:                            itemName,
	UnitOfMeasure:                   "Ounces", // this changes in the Update function
	Strain:                          "Spring Hill Kush",
	UnitThcContent:                  10.0,
	UnitThcContentUnitOfMeasure:     "Milligrams",
	UnitThcContentDose:              5.0,
	UnitThcContentDoseUnitOfMeasure: "Milligrams",
	UnitWeight:                      100.0,
	UnitWeightUnitOfMeasure:         "Milligrams",
	NumberOfDoses:                   2,
}

// Items constructs the information needed to verify the Items tab.
func (e *EvalMetrc) Items(license string) (ItemsResponse, error) {
	ci, err := e.CreateItem(license)
	if err != nil {
		return ItemsResponse{}, fmt.Errorf("could not create item: %s", err)
	}

	id := ci.Id
	ui, err := e.UpdateItem(license, id)
	if err != nil {
		return ItemsResponse{}, fmt.Errorf("could not update item: %s", err)
	}

	gi, err := e.GetItem(license, id)
	if err != nil {
		return ItemsResponse{}, fmt.Errorf("could not get item: %s", err)
	}

	ir := ItemsResponse{
		Create: ci,
		Update: ui,
		Get:    gi,
	}

	_, err = json.MarshalIndent(ir, "", "\t")
	if err != nil {
		return ItemsResponse{}, fmt.Errorf("could not marshal final struct: %s", err)
	}

	return ir, nil
}

// CreateItem creates a new Item and returns its information.
// It corresponds to Step 1 in the Items tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) CreateItem(license string) (EvalRow, error) {
	gotItems, err := e.Metrc.GetItemsActive(&license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not initially get active items: %s", err)
	}

	inputItems := []metrc.ItemPost{defaultItem}
	_, err = e.Metrc.CreateItems(inputItems, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not create initial items: %s", err)
	}

	gotItems, err = e.Metrc.GetItemsActive(&license)
	var itemId int
	var foundItemName bool
	for _, item := range gotItems {
		if item.Name == itemName {
			itemId = item.Id
			foundItemName = true
			break
		}
	}

	if !foundItemName {
		return EvalRow{}, fmt.Errorf("could not get item with matching name: %s", err)
	}

	endpoint := "items/v1/create"
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)

	resp, err := json.MarshalIndent(inputItems, "", "\t")
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not marshal input body")
	}
	response := string(resp)

	return EvalRow{
		Code:           200,
		License:        license,
		Id:             itemId,
		Name:           itemName,
		Request:        request,
		BodyOrResponse: response,
	}, nil
}

// UpdateItem displays the information to verify updating an item.
// This corresponds to Step 2 in the Items tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) UpdateItem(license string, id int) (EvalRow, error) {
	// The sheet requests that the Unit of Measure Type is changed.
	updateItem := defaultItem
	updateItem.Id = id
	updateItem.UnitOfMeasure = "Milligrams"
	updateItems := []metrc.ItemPost{updateItem}

	_, err := e.Metrc.UpdateItems(updateItems, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not update items: %s", err)
	}

	endpoint := "items/v1/update"
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)

	bodyBytes, err := json.MarshalIndent(updateItems, "", "\t")
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not marshap update body: %s", err)
	}
	body := string(bodyBytes)

	return EvalRow{
		Code:           200,
		License:        license,
		Id:             id,
		Name:           itemName,
		Request:        request,
		BodyOrResponse: body,
	}, nil
}

// GetItem displays the information to verify getting the final item.
// This corresponds to Step 3 in the Items tab of the Metrc Evaluation spreadsheet.
func (e *EvalMetrc) GetItem(license string, id int) (EvalRow, error) {
	gotItem, err := e.Metrc.GetItemsById(id, &license)
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not get items by id: %s", err)
	}

	endpoint := fmt.Sprintf("items/v1/%d", id)
	queryParam := fmt.Sprintf("?licenseNumber=%s", license)
	request := fmt.Sprintf("%s/%s%s", metrcUrl, endpoint, queryParam)

	resp, err := json.MarshalIndent(gotItem, "", "\t")
	if err != nil {
		return EvalRow{}, fmt.Errorf("could not marshal got item body: %s", err)
	}
	response := string(resp)

	return EvalRow{
		Code:           200,
		License:        license,
		Id:             id,
		Name:           itemName,
		Request:        request,
		BodyOrResponse: response,
	}, nil

}
