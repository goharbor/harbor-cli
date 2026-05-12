package tablelist

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/table"
)

func TestTableListAsync(t *testing.T) {
	// 1. Synchronous A/B Test
	cols := []table.Column{{Title: "ID", Width: 10}}
	rows := []table.Row{{"1"}}
	
	syncModel := NewModel(cols, rows, 5)
	syncView := syncModel.View().Content
	
	if strings.Contains(syncView, "Loading data...") {
		t.Errorf("Synchronous model should NOT show loading state")
	}
	if !strings.Contains(syncView, "ID") {
		t.Errorf("Synchronous model should show table headers")
	}

	// 2. Asynchronous A/B Test
	fetchFunc := func() ([]table.Row, error) {
		return []table.Row{{"2"}}, nil
	}
	
	asyncModel := NewModelWithFetch(cols, fetchFunc, 5)
	asyncView := asyncModel.View().Content
	
	if !strings.Contains(asyncView, "Loading data...") {
		t.Errorf("Async model should show loading state initially")
	}

	// Simulate Init (should return batch cmd with fetchCmd)
	cmd := asyncModel.Init()
	if cmd == nil {
		t.Errorf("Init() should return a cmd when loading")
	}

	// Execute fetchCmd to simulate Bubbletea runtime
	// Since tea.Batch returns a batchMsg, we execute the fetch directly
	msg := asyncModel.fetchCmd()
	
	// Ensure the returned msg is DataLoadedMsg
	dataMsg, ok := msg.(DataLoadedMsg)
	if !ok {
		t.Fatalf("fetchCmd did not return DataLoadedMsg")
	}

	// Update the model with the loaded data
	updatedModel, _ := asyncModel.Update(dataMsg)
	asyncModel = updatedModel.(Model)

	// Now the view should NOT contain loading, and SHOULD contain the table
	finalView := asyncModel.View().Content
	if strings.Contains(finalView, "Loading data...") {
		t.Errorf("Async model should NOT show loading state after data arrives")
	}
	if !strings.Contains(finalView, "ID") {
		t.Errorf("Async model should show table headers after data arrives")
	}
}
