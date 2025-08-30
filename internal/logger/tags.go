package logger

import "strings"

// Debug tags for filtering log output
const (
	// Core reactive system
	TagSignal   = "SIGNAL"
	TagEffect   = "EFFECT"
	TagMemo     = "MEMO"
	TagComputed = "COMPUTED"
	
	// Rendering system
	TagPaint    = "PAINT"
	TagRender   = "RENDER"
	TagPipeline = "PIPELINE"
	TagUpdate   = "UPDATE"
	TagDOM      = "DOM"
	TagCanvas   = "CANVAS"
	
	// Workflow engine
	TagEngine   = "ENGINE"
	TagWorkflow = "WORKFLOW"
	
	// Widget system
	TagWidget   = "WIDGET"
	TagText     = "TEXT"
	TagButton   = "BUTTON"
	
	// Application
	TagApp      = "APP"
	TagUI       = "UI"
	TagInit     = "INIT"
	
	// Maya framework
	TagMaya     = "MAYA"
	TagReactive = "REACTIVE"
)

// Common debug groups for convenience
var (
	// All reactive system tags
	ReactiveGroup = []string{TagSignal, TagEffect, TagMemo, TagComputed}
	
	// All rendering tags
	RenderGroup = []string{TagPaint, TagRender, TagPipeline, TagUpdate, TagDOM, TagCanvas}
	
	// All widget tags
	WidgetGroup = []string{TagWidget, TagText, TagButton}
	
	// Minimal debugging (just app and errors)
	MinimalGroup = []string{TagApp}
	
	// Common debugging scenario
	CommonGroup = []string{TagApp, TagUI, TagSignal, TagEffect}
)

// EnableGroup enables all tags in a group
func EnableGroup(group []string) {
	for _, tag := range group {
		EnableCategory(tag)
	}
}

// DisableGroup disables all tags in a group
func DisableGroup(group []string) {
	for _, tag := range group {
		DisableCategory(tag)
	}
}

// ParseDebugTags parses debug tags from string like "signal,effect,memo"
func ParseDebugTags(tags string) []string {
	if tags == "" {
		return nil
	}
	
	// Handle special group names
	switch tags {
	case "reactive":
		return ReactiveGroup
	case "render":
		return RenderGroup
	case "widget":
		return WidgetGroup
	case "minimal":
		return MinimalGroup
	case "common":
		return CommonGroup
	case "all":
		return append(append(ReactiveGroup, RenderGroup...), WidgetGroup...)
	}
	
	// Parse comma-separated tags
	result := []string{}
	for _, tag := range strings.Split(strings.ToUpper(tags), ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}