package core

// SizingType defines how an element sizes itself within its parent
type SizingType int

const (
	// SizingFit wraps tightly to content size (default)
	SizingFit SizingType = iota
	// SizingGrow expands to fill available space, sharing with other grow elements
	SizingGrow
	// SizingFixed uses exact pixel size
	SizingFixed
	// SizingPercent uses percentage of parent size
	SizingPercent
)

// SizingAxis defines sizing behavior for a single axis
type SizingAxis struct {
	Type    SizingType
	Value   float64 // For Fixed: pixels, For Percent: 0-1
	Min     float64 // Minimum size in pixels
	Max     float64 // Maximum size in pixels
}

// Sizing defines sizing behavior for both axes
type Sizing struct {
	Width  SizingAxis
	Height SizingAxis
}

// Default sizing configurations
var (
	SizingFitContent = Sizing{
		Width:  SizingAxis{Type: SizingFit, Min: 0, Max: 999999},
		Height: SizingAxis{Type: SizingFit, Min: 0, Max: 999999},
	}
	
	SizingFillWidth = Sizing{
		Width:  SizingAxis{Type: SizingGrow, Min: 0, Max: 999999},
		Height: SizingAxis{Type: SizingFit, Min: 0, Max: 999999},
	}
	
	SizingFillHeight = Sizing{
		Width:  SizingAxis{Type: SizingFit, Min: 0, Max: 999999},
		Height: SizingAxis{Type: SizingGrow, Min: 0, Max: 999999},
	}
	
	SizingFillBoth = Sizing{
		Width:  SizingAxis{Type: SizingGrow, Min: 0, Max: 999999},
		Height: SizingAxis{Type: SizingGrow, Min: 0, Max: 999999},
	}
)