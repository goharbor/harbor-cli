package selection

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/list"
)

func TestSelectionAsync(t *testing.T) {
	items := []list.Item{Item("Option 1")}
	
	// 1. Synchronous A/B Test
	syncModel := NewModel(items, "Item")
	syncView := syncModel.View().Content
	
	if strings.Contains(syncView, "Loading selections...") {
		t.Errorf("Synchronous model should NOT show loading state")
	}
	if !strings.Contains(syncView, "Select a Item") {
		t.Errorf("Synchronous model should show list title")
	}

	// 2. Asynchronous A/B Test
	fetchFunc := func() ([]list.Item, error) {
		return items, nil
	}
	
	asyncModel := NewModelWithFetch(fetchFunc, "Item")
	asyncView := asyncModel.View().Content
	
	if !strings.Contains(asyncView, "Loading selections...") {
		t.Errorf("Async model should show loading state initially")
	}

	// Simulate Bubbletea fetching data
	cmd := asyncModel.Init()
	if cmd == nil {
		t.Errorf("Init() should return a cmd when loading")
	}

	msg := asyncModel.fetchCmd()
	dataMsg, ok := msg.(DataLoadedMsg)
	if !ok {
		t.Fatalf("fetchCmd did not return DataLoadedMsg")
	}

	updatedModel, _ := asyncModel.Update(dataMsg)
	asyncModel = updatedModel.(Model)

	finalView := asyncModel.View().Content
	if strings.Contains(finalView, "Loading selections...") {
		t.Errorf("Async model should NOT show loading state after data arrives")
	}
	if !strings.Contains(finalView, "Select a Item") {
		t.Errorf("Async model should show list title after data arrives")
	}
}
