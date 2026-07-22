package app

import (
	"DailyDashBoard/internal/tui/features"
	"testing"
)

func TestDailyMenuMatchesRegistry(t *testing.T) {
	m := InitialModel()

	items := m.menus[dailyMenu].items
	registry := features.Daily()

	if len(items) != len(registry)+1 {
		t.Fatalf("expected %d menu items, got %d", len(registry)+1, len(items))
	}

	for i, feature := range registry {
		if items[i].label != feature.MenuLabel {
			t.Fatalf("expected item %d to be %q, got %q", i, feature.MenuLabel, items[i].label)
		}

		if items[i].feature == nil {
			t.Fatalf("expected item %d to reference its feature", i)
		}
	}

	if items[len(items)-1].label != "back" {
		t.Fatalf("expected final item to be back, got %q", items[len(items)-1].label)
	}
}
