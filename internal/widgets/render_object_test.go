package widgets

import (
	"testing"

	"github.com/maya-framework/maya/internal/core"
)

func TestRenderBox(t *testing.T) {
	box := &RenderBox{
		Size: Size{Width: 100, Height: 50},
	}
	
	if box.Size.Width != 100 {
		t.Errorf("Expected width 100, got %f", box.Size.Width)
	}
	
	if box.Size.Height != 50 {
		t.Errorf("Expected height 50, got %f", box.Size.Height)
	}
}

func TestRenderParagraph(t *testing.T) {
	style := TextStyle{
		FontFamily: "Arial",
		FontSize:   16,
		FontWeight: FontWeightBold,
		Color:      core.Color{R: 255, G: 0, B: 0, A: 255},
		LineHeight: 1.5,
	}
	
	paragraph := &RenderParagraph{
		Text:  "Hello World",
		Style: style,
	}
	
	if paragraph.Text != "Hello World" {
		t.Errorf("Expected text 'Hello World', got %s", paragraph.Text)
	}
	
	if paragraph.Style.FontFamily != "Arial" {
		t.Errorf("Expected font family 'Arial', got %s", paragraph.Style.FontFamily)
	}
	
	if paragraph.Style.FontSize != 16 {
		t.Errorf("Expected font size 16, got %f", paragraph.Style.FontSize)
	}
	
	if paragraph.Style.FontWeight != FontWeightBold {
		t.Errorf("Expected bold font weight, got %d", paragraph.Style.FontWeight)
	}
	
	if paragraph.Style.Color != style.Color {
		t.Errorf("Expected color %v, got %v", style.Color, paragraph.Style.Color)
	}
	
	if paragraph.Style.LineHeight != 1.5 {
		t.Errorf("Expected line height 1.5, got %f", paragraph.Style.LineHeight)
	}
}

func TestRenderButton(t *testing.T) {
	var pressed bool
	onPressed := func() {
		pressed = true
	}
	
	button := &RenderButton{
		Label:     "Click Me",
		OnPressed: onPressed,
		Disabled:  false,
	}
	
	if button.Label != "Click Me" {
		t.Errorf("Expected label 'Click Me', got %s", button.Label)
	}
	
	if button.Disabled {
		t.Error("Button should not be disabled")
	}
	
	// Test callback
	button.OnPressed()
	if !pressed {
		t.Error("OnPressed callback should be called")
	}
}

func TestRenderButton_DisabledState(t *testing.T) {
	button := &RenderButton{
		Label:    "Disabled",
		Disabled: true,
	}
	
	if !button.Disabled {
		t.Error("Button should be disabled")
	}
}


func TestRenderButton_NilCallback(t *testing.T) {
	button := &RenderButton{
		Label:     "No Callback",
		OnPressed: nil,
	}
	
	// Should not panic
	if button.OnPressed != nil {
		button.OnPressed()
	}
}

func TestEdgeInsets(t *testing.T) {
	insets := EdgeInsets{
		Top:    10,
		Right:  20,
		Bottom: 30,
		Left:   40,
	}
	
	if insets.Top != 10 {
		t.Errorf("Expected top 10, got %f", insets.Top)
	}
	
	if insets.Right != 20 {
		t.Errorf("Expected right 20, got %f", insets.Right)
	}
	
	if insets.Bottom != 30 {
		t.Errorf("Expected bottom 30, got %f", insets.Bottom)
	}
	
	if insets.Left != 40 {
		t.Errorf("Expected left 40, got %f", insets.Left)
	}
}

func TestSize(t *testing.T) {
	size := Size{
		Width:  100,
		Height: 50,
	}
	
	if size.Width != 100 {
		t.Errorf("Expected width 100, got %f", size.Width)
	}
	
	if size.Height != 50 {
		t.Errorf("Expected height 50, got %f", size.Height)
	}
}

func TestTextStyle_FontWeights(t *testing.T) {
	testCases := []struct {
		name   string
		weight FontWeight
		value  int
	}{
		{"Normal", FontWeightNormal, 400},
		{"Bold", FontWeightBold, 700},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if int(tc.weight) != tc.value {
				t.Errorf("Expected %s to be %d, got %d", tc.name, tc.value, int(tc.weight))
			}
		})
	}
}

func TestColors(t *testing.T) {
	testCases := []struct {
		name  string
		color core.Color
		r, g, b, a uint8
	}{
		{"Black", ColorBlack, 0, 0, 0, 255},
		{"White", ColorWhite, 255, 255, 255, 255},
		{"Red", ColorRed, 255, 0, 0, 255},
		{"Green", ColorGreen, 0, 255, 0, 255},
		{"Blue", ColorBlue, 0, 0, 255, 255},
		{"Transparent", ColorTransparent, 0, 0, 0, 0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.color.R != tc.r || tc.color.G != tc.g || tc.color.B != tc.b || tc.color.A != tc.a {
				t.Errorf("Expected %s to be (%d,%d,%d,%d), got (%d,%d,%d,%d)",
					tc.name, tc.r, tc.g, tc.b, tc.a,
					tc.color.R, tc.color.G, tc.color.B, tc.color.A)
			}
		})
	}
}

func TestAlignment(t *testing.T) {
	testCases := []struct {
		name      string
		alignment Alignment
		x, y      float64
	}{
		{"TopLeft", AlignmentTopLeft, -1, -1},
		{"TopCenter", AlignmentTopCenter, 0, -1},
		{"TopRight", AlignmentTopRight, 1, -1},
		{"CenterLeft", AlignmentCenterLeft, -1, 0},
		{"Center", AlignmentCenter, 0, 0},
		{"CenterRight", AlignmentCenterRight, 1, 0},
		{"BottomLeft", AlignmentBottomLeft, -1, 1},
		{"BottomCenter", AlignmentBottomCenter, 0, 1},
		{"BottomRight", AlignmentBottomRight, 1, 1},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.alignment.X != tc.x || tc.alignment.Y != tc.y {
				t.Errorf("Expected %s to be (%f,%f), got (%f,%f)",
					tc.name, tc.x, tc.y, tc.alignment.X, tc.alignment.Y)
			}
		})
	}
}

func TestBoxShadow(t *testing.T) {
	shadow := BoxShadow{
		Color:        core.Color{R: 0, G: 0, B: 0, A: 128},
		Offset:       core.Offset{X: 5, Y: 5},
		BlurRadius:   10,
		SpreadRadius: 2,
	}
	
	if shadow.Color.A != 128 {
		t.Errorf("Expected shadow alpha 128, got %d", shadow.Color.A)
	}
	
	if shadow.Offset.X != 5 || shadow.Offset.Y != 5 {
		t.Errorf("Expected shadow offset (5,5), got (%f,%f)", shadow.Offset.X, shadow.Offset.Y)
	}
	
	if shadow.BlurRadius != 10 {
		t.Errorf("Expected blur radius 10, got %f", shadow.BlurRadius)
	}
	
	if shadow.SpreadRadius != 2 {
		t.Errorf("Expected spread radius 2, got %f", shadow.SpreadRadius)
	}
}