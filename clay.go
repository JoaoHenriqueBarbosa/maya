package clay

import "unsafe"

// VERSION: 0.14
// Direct transpilation from Clay C library to Go
// Original: https://github.com/nicbarker/clay

// -----------------------------------------
// HEADER DECLARATIONS ---------------------
// -----------------------------------------

// Public Helper Functions ------------------------

func CLAY__MAX(x, y float32) float32 {
    if x > y {
        return x
    }
    return y
}

func CLAY__MIN(x, y float32) float32 {
    if x < y {
        return x
    }
    return y
}

// CLAY_TEXT_CONFIG macro converted to function
func CLAY_TEXT_CONFIG(config Clay_TextElementConfig) *Clay_TextElementConfig {
    return Clay__StoreTextElementConfig(config)
}

// CLAY_BORDER_OUTSIDE macro converted to function
func CLAY_BORDER_OUTSIDE(widthValue float32) Clay_BorderWidth {
    return Clay_BorderWidth{uint16(widthValue), uint16(widthValue), uint16(widthValue), uint16(widthValue), 0}
}

// CLAY_BORDER_ALL macro converted to function
func CLAY_BORDER_ALL(widthValue float32) Clay_BorderWidth {
    return Clay_BorderWidth{uint16(widthValue), uint16(widthValue), uint16(widthValue), uint16(widthValue), uint16(widthValue)}
}

// CLAY_CORNER_RADIUS macro converted to function
func CLAY_CORNER_RADIUS(radius float32) Clay_CornerRadius {
    return Clay_CornerRadius{radius, radius, radius, radius}
}

// CLAY_PADDING_ALL macro converted to function
func CLAY_PADDING_ALL(padding float32) Clay_Padding {
    return Clay_Padding{uint16(padding), uint16(padding), uint16(padding), uint16(padding)}
}

// CLAY_SIZING_FIT macro converted to function
func CLAY_SIZING_FIT(min, max float32) Clay_SizingAxis {
    return Clay_SizingAxis{
        size: Clay_SizingAxisSize{
            minMax: Clay_SizingMinMax{min, max},
        },
        type_: CLAY__SIZING_TYPE_FIT,
    }
}

// CLAY_SIZING_GROW macro converted to function
func CLAY_SIZING_GROW(min, max float32) Clay_SizingAxis {
    return Clay_SizingAxis{
        size: Clay_SizingAxisSize{
            minMax: Clay_SizingMinMax{min, max},
        },
        type_: CLAY__SIZING_TYPE_GROW,
    }
}

// CLAY_SIZING_FIXED macro converted to function
func CLAY_SIZING_FIXED(fixedSize float32) Clay_SizingAxis {
    return Clay_SizingAxis{
        size: Clay_SizingAxisSize{
            minMax: Clay_SizingMinMax{fixedSize, fixedSize},
        },
        type_: CLAY__SIZING_TYPE_FIXED,
    }
}

// CLAY_SIZING_PERCENT macro converted to function
func CLAY_SIZING_PERCENT(percentOfParent float32) Clay_SizingAxis {
    return Clay_SizingAxis{
        size: Clay_SizingAxisSize{
            percent: percentOfParent,
        },
        type_: CLAY__SIZING_TYPE_PERCENT,
    }
}

// Note: If a compile error led you here, you might be trying to use CLAY_ID with something other than a string literal. To construct an ID with a dynamic string, use CLAY_SID instead.
// CLAY_ID macro converted to function
func CLAY_ID(label string) Clay_ElementId {
    return CLAY_SID(CLAY_STRING(label))
}

// CLAY_SID macro converted to function
func CLAY_SID(label Clay_String) Clay_ElementId {
    return Clay__HashString(label, 0)
}

// Note: If a compile error led you here, you might be trying to use CLAY_IDI with something other than a string literal. To construct an ID with a dynamic string, use CLAY_SIDI instead.
// CLAY_IDI macro converted to function
func CLAY_IDI(label string, index uint32) Clay_ElementId {
    return CLAY_SIDI(CLAY_STRING(label), index)
}

// CLAY_SIDI macro converted to function
func CLAY_SIDI(label Clay_String, index uint32) Clay_ElementId {
    return Clay__HashStringWithOffset(label, index, 0)
}

// Note: If a compile error led you here, you might be trying to use CLAY_ID_LOCAL with something other than a string literal. To construct an ID with a dynamic string, use CLAY_SID_LOCAL instead.
// CLAY_ID_LOCAL macro converted to function
func CLAY_ID_LOCAL(label string) Clay_ElementId {
    return CLAY_SID_LOCAL(CLAY_STRING(label), 0)
}

// CLAY_SID_LOCAL macro converted to function
func CLAY_SID_LOCAL(label Clay_String, index uint32) Clay_ElementId {
    _ = index // index parameter exists in macro but not used
    return Clay__HashString(label, Clay__GetParentElementId())
}

// Note: If a compile error led you here, you might be trying to use CLAY_IDI_LOCAL with something other than a string literal. To construct an ID with a dynamic string, use CLAY_SIDI_LOCAL instead.
// CLAY_IDI_LOCAL macro converted to function
func CLAY_IDI_LOCAL(label string, index uint32) Clay_ElementId {
    return CLAY_SIDI_LOCAL(CLAY_STRING(label), index)
}

// CLAY_SIDI_LOCAL macro converted to function
func CLAY_SIDI_LOCAL(label Clay_String, index uint32) Clay_ElementId {
    return Clay__HashStringWithOffset(label, index, Clay__GetParentElementId())
}

// CLAY__STRING_LENGTH macro converted to function
func CLAY__STRING_LENGTH(s string) int32 {
    return int32(len(s))
}

// CLAY__ENSURE_STRING_LITERAL macro - not needed in Go
// In C this ensures x is a string literal at compile time
// Go handles strings differently

// Note: If an error led you here, it's because CLAY_STRING can only be used with string literals, i.e. CLAY_STRING("SomeString") and not CLAY_STRING(yourString)
// CLAY_STRING macro converted to function
func CLAY_STRING(s string) Clay_String {
    return Clay_String{
        isStaticallyAllocated: true,
        length:                int32(len(s)),
        chars:                 unsafe.Pointer(unsafe.StringData(s)),
    }
}

// CLAY_STRING_CONST - used for compile-time string constants
// In Go, we can just use CLAY_STRING function instead

var CLAY__ELEMENT_DEFINITION_LATCH uint8

// GCC marks the above CLAY__ELEMENT_DEFINITION_LATCH as an unused variable for files that include clay.h but don't declare any layout
// This is to suppress that warning
func Clay__SuppressUnusedLatchDefinitionVariableWarning() {
    _ = CLAY__ELEMENT_DEFINITION_LATCH
}

// Publicly visible layout element macros -----------------------------------------------------

/* This macro looks scary on the surface, but is actually quite simple.
  It turns a macro call like this:

  CLAY({
    .id = CLAY_ID("Container"),
    .backgroundColor = { 255, 200, 200, 255 }
  }) {
      ...children declared here
  }

  Into calls like this:

  Clay_OpenElement();
  Clay_ConfigureOpenElement((Clay_ElementDeclaration) {
    .id = CLAY_ID("Container"),
    .backgroundColor = { 255, 200, 200, 255 }
  });
  ...children declared here
  Clay_CloseElement();

  The for loop will only ever run a single iteration, putting Clay__CloseElement() in the increment of the loop
  means that it will run after the body - where the children are declared. It just exists to make sure you don't forget
  to call Clay_CloseElement().
*/
// CLAY macro converted to function
// In Go, we need to use a function with a callback for children
func CLAY(config Clay_ElementDeclaration, children func()) {
    Clay__OpenElement()
    Clay__ConfigureOpenElement(config)
    if children != nil {
        children()
    }
    Clay__CloseElement()
}

// These macros exist to allow the CLAY() macro to be called both with an inline struct definition, such as
// CLAY({ .id = something... });
// As well as by passing a predefined declaration struct
// Clay_ElementDeclaration declarationStruct = ...
// CLAY(declarationStruct);
// CLAY__WRAPPER macros - not needed in Go, Go handles struct initialization differently

// CLAY_TEXT macro converted to function
func CLAY_TEXT(text Clay_String, textConfig *Clay_TextElementConfig) {
    Clay__OpenTextElement(text, textConfig)
}

// C++ and C conditional compilation macros - not needed in Go
// CLAY__INIT - in Go we just use type{}
// CLAY_PACKED_ENUM - Go doesn't have packed enums, we'll use uint8 for enum types
// CLAY__DEFAULT_STRUCT - Go uses zero values by default

// Utility Structs -------------------------

// Note: Clay_String is not guaranteed to be null terminated. It may be if created from a literal C string,
// but it is also used to represent slices.
type Clay_String struct {
    // Set this boolean to true if the char* data underlying this string will live for the entire lifetime of the program.
    // This will automatically be set for strings created with CLAY_STRING, as the macro requires a string literal.
    isStaticallyAllocated bool
    length                int32
    // The underlying character memory. Note: this will not be copied and will not extend the lifetime of the underlying memory.
    chars unsafe.Pointer
}

// Clay_StringSlice is used to represent non owning string slices, and includes
// a baseChars field which points to the string this slice is derived from.
type Clay_StringSlice struct {
    length    int32
    chars     unsafe.Pointer
    baseChars unsafe.Pointer // The source string / char* that this slice was derived from
}

// Forward declaration of Clay_Context - not needed in Go

// Clay_Arena is a memory arena structure that is used by clay to manage its internal allocations.
// Rather than creating it by hand, it's easier to use Clay_CreateArenaWithCapacityAndMemory()
type Clay_Arena struct {
    nextAllocation uintptr
    capacity       uint64
    memory         []byte
}

type Clay_Dimensions struct {
    width  float32
    height float32
}

type Clay_Vector2 struct {
    x float32
    y float32
}

// Internally clay conventionally represents colors as 0-255, but interpretation is up to the renderer.
type Clay_Color struct {
    r float32
    g float32
    b float32
    a float32
}

type Clay_BoundingBox struct {
    x      float32
    y      float32
    width  float32
    height float32
}

// Primarily created via the CLAY_ID(), CLAY_IDI(), CLAY_ID_LOCAL() and CLAY_IDI_LOCAL() macros.
// Represents a hashed string ID used for identifying and finding specific clay UI elements, required
// by functions such as Clay_PointerOver() and Clay_GetElementData().
type Clay_ElementId struct {
    id       uint32      // The resulting hash generated from the other fields.
    offset   uint32      // A numerical offset applied after computing the hash from stringId.
    baseId   uint32      // A base hash value to start from, for example the parent element ID is used when calculating CLAY_ID_LOCAL().
    stringId Clay_String // The string id to hash.
}

// A sized array of Clay_ElementId.
type Clay_ElementIdArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_ElementId
}

// Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
type Clay_CornerRadius struct {
    topLeft     float32
    topRight    float32
    bottomLeft  float32
    bottomRight float32
}

// Element Configs ---------------------------

// Controls the direction in which child elements will be automatically laid out.
type Clay_LayoutDirection uint8

const (
    // (Default) Lays out child elements from left to right with increasing x.
    CLAY_LEFT_TO_RIGHT Clay_LayoutDirection = iota
    // Lays out child elements from top to bottom with increasing y.
    CLAY_TOP_TO_BOTTOM
)

// Controls the alignment along the x axis (horizontal) of child elements.
type Clay_LayoutAlignmentX uint8

const (
    // (Default) Aligns child elements to the left hand side of this element, offset by padding.width.left
    CLAY_ALIGN_X_LEFT Clay_LayoutAlignmentX = iota
    // Aligns child elements to the right hand side of this element, offset by padding.width.right
    CLAY_ALIGN_X_RIGHT
    // Aligns child elements horizontally to the center of this element
    CLAY_ALIGN_X_CENTER
)

// Controls the alignment along the y axis (vertical) of child elements.
type Clay_LayoutAlignmentY uint8

const (
    // (Default) Aligns child elements to the top of this element, offset by padding.width.top
    CLAY_ALIGN_Y_TOP Clay_LayoutAlignmentY = iota
    // Aligns child elements to the bottom of this element, offset by padding.width.bottom
    CLAY_ALIGN_Y_BOTTOM
    // Aligns child elements vertically to the center of this element
    CLAY_ALIGN_Y_CENTER
)

// Controls how the element takes up space inside its parent container.
type Clay__SizingType uint8

const (
    // (default) Wraps tightly to the size of the element's contents.
    CLAY__SIZING_TYPE_FIT Clay__SizingType = iota
    // Expands along this axis to fill available space in the parent element, sharing it with other GROW elements.
    CLAY__SIZING_TYPE_GROW
    // Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
    CLAY__SIZING_TYPE_PERCENT
    // Clamps the axis size to an exact size in pixels.
    CLAY__SIZING_TYPE_FIXED
)

// Controls how child elements are aligned on each axis.
type Clay_ChildAlignment struct {
    x Clay_LayoutAlignmentX // Controls alignment of children along the x axis.
    y Clay_LayoutAlignmentY // Controls alignment of children along the y axis.
}

// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to,
// overriding sizing types such as FIT or GROW.
type Clay_SizingMinMax struct {
    min float32 // The smallest final size of the element on this axis will be this value in pixels.
    max float32 // The largest final size of the element on this axis will be this value in pixels.
}

// Controls the sizing of this element along one axis inside its parent container.
type Clay_SizingAxisSize struct {
    minMax  Clay_SizingMinMax // Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to, overriding sizing types such as FIT or GROW.
    percent float32           // Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
}

type Clay_SizingAxis struct {
    size  Clay_SizingAxisSize
    type_ Clay__SizingType // Controls how the element takes up space inside its parent container.
}

// Controls the sizing of this element along one axis inside its parent container.
type Clay_Sizing struct {
    width  Clay_SizingAxis // Controls the width sizing of the element, along the x axis.
    height Clay_SizingAxis // Controls the height sizing of the element, along the y axis.
}

// Controls "padding" in pixels, which is a gap between the bounding box of this element and where its children
// will be placed.
type Clay_Padding struct {
    left   uint16
    right  uint16
    top    uint16
    bottom uint16
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Controls various settings that affect the size and position of an element, as well as the sizes and positions
// of any child elements.
type Clay_LayoutConfig struct {
    sizing          Clay_Sizing          // Controls the sizing of this element inside it's parent container, including FIT, GROW, PERCENT and FIXED sizing.
    padding         Clay_Padding         // Controls "padding" in pixels, which is a gap between the bounding box of this element and where its children will be placed.
    childGap        uint16               // Controls the gap in pixels between child elements along the layout axis (horizontal gap for LEFT_TO_RIGHT, vertical gap for TOP_TO_BOTTOM).
    childAlignment  Clay_ChildAlignment  // Controls how child elements are aligned on each axis.
    layoutDirection Clay_LayoutDirection // Controls the direction in which child elements will be automatically laid out.
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Default layout configuration
var CLAY_LAYOUT_DEFAULT Clay_LayoutConfig

// Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
type Clay_TextElementConfigWrapMode uint8

const (
    // (default) breaks on whitespace characters.
    CLAY_TEXT_WRAP_WORDS Clay_TextElementConfigWrapMode = iota
    // Don't break on space characters, only on newlines.
    CLAY_TEXT_WRAP_NEWLINES
    // Disable text wrapping entirely.
    CLAY_TEXT_WRAP_NONE
)

// Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
type Clay_TextAlignment uint8

const (
    // (default) Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
    CLAY_TEXT_ALIGN_LEFT Clay_TextAlignment = iota
    // Horizontally aligns wrapped lines of text to the center of their bounding box.
    CLAY_TEXT_ALIGN_CENTER
    // Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
    CLAY_TEXT_ALIGN_RIGHT
)

// Controls various functionality related to text elements.
type Clay_TextElementConfig struct {
    // A pointer that will be transparently passed through to the resulting render command.
    userData interface{}
    // The RGBA color of the font to render, conventionally specified as 0-255.
    textColor Clay_Color
    // An integer transparently passed to Clay_MeasureText to identify the font to use.
    // The debug view will pass fontId = 0 for its internal text.
    fontId uint16
    // Controls the size of the font. Handled by the function provided to Clay_MeasureText.
    fontSize uint16
    // Controls extra horizontal spacing between characters. Handled by the function provided to Clay_MeasureText.
    letterSpacing uint16
    // Controls additional vertical space between wrapped lines of text.
    lineHeight uint16
    // Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
    // CLAY_TEXT_WRAP_WORDS (default) breaks on whitespace characters.
    // CLAY_TEXT_WRAP_NEWLINES doesn't break on space characters, only on newlines.
    // CLAY_TEXT_WRAP_NONE disables wrapping entirely.
    wrapMode Clay_TextElementConfigWrapMode
    // Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
    // CLAY_TEXT_ALIGN_LEFT (default) - Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
    // CLAY_TEXT_ALIGN_CENTER - Horizontally aligns wrapped lines of text to the center of their bounding box.
    // CLAY_TEXT_ALIGN_RIGHT - Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
    textAlignment Clay_TextAlignment
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Aspect Ratio --------------------------------

// Controls various settings related to aspect ratio scaling element.
type Clay_AspectRatioElementConfig struct {
    aspectRatio float32 // A float representing the target "Aspect ratio" for an element, which is its final width divided by its final height.
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Image --------------------------------

// Controls various settings related to image elements.
type Clay_ImageElementConfig struct {
    imageData interface{} // A transparent pointer used to pass image data through to the renderer.
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Floating -----------------------------

// Controls where a floating element is offset relative to its parent element.
// Note: see https://github.com/user-attachments/assets/b8c6dfaa-c1b1-41a4-be55-013473e4a6ce for a visual explanation.
type Clay_FloatingAttachPointType uint8

const (
    CLAY_ATTACH_POINT_LEFT_TOP Clay_FloatingAttachPointType = iota
    CLAY_ATTACH_POINT_LEFT_CENTER
    CLAY_ATTACH_POINT_LEFT_BOTTOM
    CLAY_ATTACH_POINT_CENTER_TOP
    CLAY_ATTACH_POINT_CENTER_CENTER
    CLAY_ATTACH_POINT_CENTER_BOTTOM
    CLAY_ATTACH_POINT_RIGHT_TOP
    CLAY_ATTACH_POINT_RIGHT_CENTER
    CLAY_ATTACH_POINT_RIGHT_BOTTOM
)

// Controls where a floating element is offset relative to its parent element.
type Clay_FloatingAttachPoints struct {
    element Clay_FloatingAttachPointType // Controls the origin point on a floating element that attaches to its parent.
    parent  Clay_FloatingAttachPointType // Controls the origin point on the parent element that the floating element attaches to.
}

// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath a floating element.
type Clay_PointerCaptureMode uint8

const (
    // (default) "Capture" the pointer event and don't allow events like hover and click to pass through to elements underneath.
    CLAY_POINTER_CAPTURE_MODE_CAPTURE Clay_PointerCaptureMode = iota
    //    CLAY_POINTER_CAPTURE_MODE_PARENT, TODO pass pointer through to attached parent

    // Transparently pass through pointer events like hover and click to elements underneath the floating element.
    CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH
)

// Controls which element a floating element is "attached" to (i.e. relative offset from).
type Clay_FloatingAttachToElement uint8

const (
    // (default) Disables floating for this element.
    CLAY_ATTACH_TO_NONE Clay_FloatingAttachToElement = iota
    // Attaches this floating element to its parent, positioned based on the .attachPoints and .offset fields.
    CLAY_ATTACH_TO_PARENT
    // Attaches this floating element to an element with a specific ID, specified with the .parentId field. positioned based on the .attachPoints and .offset fields.
    CLAY_ATTACH_TO_ELEMENT_WITH_ID
    // Attaches this floating element to the root of the layout, which combined with the .offset field provides functionality similar to "absolute positioning".
    CLAY_ATTACH_TO_ROOT
)

// Controls whether or not a floating element is clipped to the same clipping rectangle as the element it's attached to.
type Clay_FloatingClipToElement uint8

const (
    // (default) - The floating element does not inherit clipping.
    CLAY_CLIP_TO_NONE Clay_FloatingClipToElement = iota
    // The floating element is clipped to the same clipping rectangle as the element it's attached to.
    CLAY_CLIP_TO_ATTACHED_PARENT
)

// Controls various settings related to "floating" elements, which are elements that "float" above other elements, potentially overlapping their boundaries,
// and not affecting the layout of sibling or parent elements.
type Clay_FloatingElementConfig struct {
    // Offsets this floating element by the provided x,y coordinates from its attachPoints.
    offset Clay_Vector2
    // Expands the boundaries of the outer floating element without affecting its children.
    expand Clay_Dimensions
    // When used in conjunction with .attachTo = CLAY_ATTACH_TO_ELEMENT_WITH_ID, attaches this floating element to the element in the hierarchy with the provided ID.
    // Hint: attach the ID to the other element with .id = CLAY_ID("yourId"), and specify the id the same way, with .parentId = CLAY_ID("yourId").id
    parentId uint32
    // Controls the z index of this floating element and all its children. Floating elements are sorted in ascending z order before output.
    // zIndex is also passed to the renderer for all elements contained within this floating element.
    zIndex int16
    // Controls how mouse pointer events like hover and click are captured or passed through to elements underneath / behind a floating element.
    // Enum is of the form CLAY_ATTACH_POINT_foo_bar. See Clay_FloatingAttachPoints for more details.
    // Note: see <img src="https://github.com/user-attachments/assets/b8c6dfaa-c1b1-41a4-be55-013473e4a6ce />
    // and <img src="https://github.com/user-attachments/assets/ebe75e0d-1904-46b0-982d-418f929d1516 /> for a visual explanation.
    attachPoints Clay_FloatingAttachPoints
    // Controls how mouse pointer events like hover and click are captured or passed through to elements underneath a floating element.
    // CLAY_POINTER_CAPTURE_MODE_CAPTURE (default) - "Capture" the pointer event and don't allow events like hover and click to pass through to elements underneath.
    // CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH - Transparently pass through pointer events like hover and click to elements underneath the floating element.
    pointerCaptureMode Clay_PointerCaptureMode
    // Controls which element a floating element is "attached" to (i.e. relative offset from).
    // CLAY_ATTACH_TO_NONE (default) - Disables floating for this element.
    // CLAY_ATTACH_TO_PARENT - Attaches this floating element to its parent, positioned based on the .attachPoints and .offset fields.
    // CLAY_ATTACH_TO_ELEMENT_WITH_ID - Attaches this floating element to an element with a specific ID, specified with the .parentId field. positioned based on the .attachPoints and .offset fields.
    // CLAY_ATTACH_TO_ROOT - Attaches this floating element to the root of the layout, which combined with the .offset field provides functionality similar to "absolute positioning".
    attachTo Clay_FloatingAttachToElement
    // Controls whether or not a floating element is clipped to the same clipping rectangle as the element it's attached to.
    // CLAY_CLIP_TO_NONE (default) - The floating element does not inherit clipping.
    // CLAY_CLIP_TO_ATTACHED_PARENT - The floating element is clipped to the same clipping rectangle as the element it's attached to.
    clipTo Clay_FloatingClipToElement
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Custom -----------------------------

// Controls various settings related to custom elements.
type Clay_CustomElementConfig struct {
    // A transparent pointer through which you can pass custom data to the renderer.
    // Generates CUSTOM render commands.
    customData interface{}
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Scroll -----------------------------

// Controls the axis on which an element switches to "scrolling", which clips the contents and allows scrolling in that direction.
type Clay_ClipElementConfig struct {
    horizontal   bool         // Clip overflowing elements on the X axis.
    vertical     bool         // Clip overflowing elements on the Y axis.
    childOffset  Clay_Vector2 // Offsets the x,y positions of all child elements. Used primarily for scrolling containers.
}

// CLAY__WRAPPER_STRUCT - not needed in Go

// Border -----------------------------

// Controls the widths of individual element borders.
type Clay_BorderWidth struct {
    left   uint16
    right  uint16
    top    uint16
    bottom uint16
    // Creates borders between each child element, depending on the .layoutDirection.
    // e.g. for LEFT_TO_RIGHT, borders will be vertical lines, and for TOP_TO_BOTTOM borders will be horizontal lines.
    // .betweenChildren borders will result in individual RECTANGLE render commands being generated.
    betweenChildren uint16
}

// Controls settings related to element borders.
type Clay_BorderElementConfig struct {
    color Clay_Color        // Controls the color of all borders with width > 0. Conventionally represented as 0-255, but interpretation is up to the renderer.
    width Clay_BorderWidth // Controls the widths of individual borders. At least one of these should be > 0 for a BORDER render command to be generated.
}

// Render Command Data -----------------------------

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_TEXT
type Clay_TextRenderData struct {
    // A string slice containing the text to be rendered.
    // Note: this is not guaranteed to be null terminated.
    stringContents Clay_StringSlice
    // Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
    textColor Clay_Color
    // An integer representing the font to use to render this text, transparently passed through from the text declaration.
    fontId uint16
    fontSize uint16
    // Specifies the extra whitespace gap in pixels between each character.
    letterSpacing uint16
    // The height of the bounding box for this line of text.
    lineHeight uint16
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_RECTANGLE
type Clay_RectangleRenderData struct {
    // The solid background color to fill this rectangle with. Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
    backgroundColor Clay_Color
    // Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
    // The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
    cornerRadius Clay_CornerRadius
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_IMAGE
type Clay_ImageRenderData struct {
    // The tint color for this image. Note that the default value is 0,0,0,0 and should likely be interpreted
    // as "untinted".
    // Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
    backgroundColor Clay_Color
    // Controls the "radius", or corner rounding of this image.
    // The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
    cornerRadius Clay_CornerRadius
    // A pointer transparently passed through from the original element definition, typically used to represent image data.
    imageData interface{}
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_CUSTOM
type Clay_CustomRenderData struct {
    // Passed through from .backgroundColor in the original element declaration.
    // Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
    backgroundColor Clay_Color
    // Controls the "radius", or corner rounding of this custom element.
    // The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
    cornerRadius Clay_CornerRadius
    // A pointer transparently passed through from the original element definition.
    customData interface{}
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_START || commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
type Clay_ClipRenderData struct {
    horizontal bool
    vertical   bool
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_BORDER
type Clay_BorderRenderData struct {
    // Controls a shared color for all this element's borders.
    // Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
    color Clay_Color
    // Specifies the "radius", or corner rounding of this border element.
    // The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
    cornerRadius Clay_CornerRadius
    // Controls individual border side widths.
    width Clay_BorderWidth
}

// A struct union containing data specific to this command's .commandType
type Clay_RenderData struct {
    // Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_RECTANGLE
    rectangle Clay_RectangleRenderData
    // Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_TEXT
    text Clay_TextRenderData
    // Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_IMAGE
    image Clay_ImageRenderData
    // Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_CUSTOM
    custom Clay_CustomRenderData
    // Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_BORDER
    border Clay_BorderRenderData
    // Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_START|END
    clip Clay_ClipRenderData
}

// Miscellaneous Structs & Enums ---------------------------------

// Data representing the current internal state of a scrolling element.
type Clay_ScrollContainerData struct {
    // Note: This is a pointer to the real internal scroll position, mutating it may cause a change in final layout.
    // Intended for use with external functionality that modifies scroll position, such as scroll bars or auto scrolling.
    scrollPosition *Clay_Vector2
    // The bounding box of the scroll element.
    scrollContainerDimensions Clay_Dimensions
    // The outer dimensions of the inner scroll container content, including the padding of the parent scroll container.
    contentDimensions Clay_Dimensions
    // The config that was originally passed to the clip element.
    config Clay_ClipElementConfig
    // Indicates whether an actual scroll container matched the provided ID or if the default struct was returned.
    found bool
}

// Bounding box and other data for a specific UI element.
type Clay_ElementData struct {
    // The rectangle that encloses this UI element, with the position relative to the root of the layout.
    boundingBox Clay_BoundingBox
    // Indicates whether an actual Element matched the provided ID or if the default struct was returned.
    found bool
}

// Used by renderers to determine specific handling for each render command.
type Clay_RenderCommandType uint8

const (
    // This command type should be skipped.
    CLAY_RENDER_COMMAND_TYPE_NONE Clay_RenderCommandType = iota
    // The renderer should draw a solid color rectangle.
    CLAY_RENDER_COMMAND_TYPE_RECTANGLE
    // The renderer should draw a colored border inset into the bounding box.
    CLAY_RENDER_COMMAND_TYPE_BORDER
    // The renderer should draw text.
    CLAY_RENDER_COMMAND_TYPE_TEXT
    // The renderer should draw an image.
    CLAY_RENDER_COMMAND_TYPE_IMAGE
    // The renderer should begin clipping all future draw commands, only rendering content that falls within the provided boundingBox.
    CLAY_RENDER_COMMAND_TYPE_SCISSOR_START
    // The renderer should finish any previously active clipping, and begin rendering elements in full again.
    CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
    // The renderer should provide a custom implementation for handling this render command based on its .customData
    CLAY_RENDER_COMMAND_TYPE_CUSTOM
)

type Clay_RenderCommand struct {
    // A rectangular box that fully encloses this UI element, with the position relative to the root of the layout.
    boundingBox Clay_BoundingBox
    // A struct union containing data specific to this command's commandType.
    renderData Clay_RenderData
    // A pointer transparently passed through from the original element declaration.
    userData interface{}
    // The id of this element, transparently passed through from the original element declaration.
    id uint32
    // The z order required for drawing this command correctly.
    // Note: the render command array is already sorted in ascending order, and will produce correct results if drawn in naive order.
    // This field is intended for use in batching renderers for improved performance.
    zIndex int16
    // Specifies how to handle rendering of this command.
    // CLAY_RENDER_COMMAND_TYPE_RECTANGLE - The renderer should draw a solid color rectangle.
    // CLAY_RENDER_COMMAND_TYPE_BORDER - The renderer should draw a colored border inset into the bounding box.
    // CLAY_RENDER_COMMAND_TYPE_TEXT - The renderer should draw text.
    // CLAY_RENDER_COMMAND_TYPE_IMAGE - The renderer should draw an image.
    // CLAY_RENDER_COMMAND_TYPE_SCISSOR_START - The renderer should begin clipping all future draw commands, only rendering content that falls within the provided boundingBox.
    // CLAY_RENDER_COMMAND_TYPE_SCISSOR_END - The renderer should finish any previously active clipping, and begin rendering elements in full again.
    // CLAY_RENDER_COMMAND_TYPE_CUSTOM - The renderer should provide a custom implementation for handling this render command based on its .customData
    commandType Clay_RenderCommandType
}

// A sized array of render commands.
type Clay_RenderCommandArray struct {
    // The underlying max capacity of the array, not necessarily all initialized.
    capacity int32
    // The number of initialized elements in this array. Used for loops and iteration.
    length int32
    // A pointer to the first element in the internal array.
    internalArray *Clay_RenderCommand
}

// Represents the current state of interaction with clay this frame.
type Clay_PointerDataInteractionState uint8

const (
    // A left mouse click, or touch occurred this frame.
    CLAY_POINTER_DATA_PRESSED_THIS_FRAME Clay_PointerDataInteractionState = iota
    // The left mouse button click or touch happened at some point in the past, and is still currently held down this frame.
    CLAY_POINTER_DATA_PRESSED
    // The left mouse button click or touch was released this frame.
    CLAY_POINTER_DATA_RELEASED_THIS_FRAME
    // The left mouse button click or touch is not currently down / was released at some point in the past.
    CLAY_POINTER_DATA_RELEASED
)

// Information on the current state of pointer interactions this frame.
type Clay_PointerData struct {
    // The position of the mouse / touch / pointer relative to the root of the layout.
    position Clay_Vector2
    // Represents the current state of interaction with clay this frame.
    // CLAY_POINTER_DATA_PRESSED_THIS_FRAME - A left mouse click, or touch occurred this frame.
    // CLAY_POINTER_DATA_PRESSED - The left mouse button click or touch happened at some point in the past, and is still currently held down this frame.
    // CLAY_POINTER_DATA_RELEASED_THIS_FRAME - The left mouse button click or touch was released this frame.
    // CLAY_POINTER_DATA_RELEASED - The left mouse button click or touch is not currently down / was released at some point in the past.
    state Clay_PointerDataInteractionState
}

type Clay_ElementDeclaration struct {
    // Primarily created via the CLAY_ID(), CLAY_IDI(), CLAY_ID_LOCAL() and CLAY_IDI_LOCAL() macros.
    // Represents a hashed string ID used for identifying and finding specific clay UI elements, required by functions such as Clay_PointerOver() and Clay_GetElementData().
    id Clay_ElementId
    // Controls various settings that affect the size and position of an element, as well as the sizes and positions of any child elements.
    layout Clay_LayoutConfig
    // Controls the background color of the resulting element.
    // By convention specified as 0-255, but interpretation is up to the renderer.
    // If no other config is specified, .backgroundColor will generate a RECTANGLE render command, otherwise it will be passed as a property to IMAGE or CUSTOM render commands.
    backgroundColor Clay_Color
    // Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
    cornerRadius Clay_CornerRadius
    // Controls settings related to aspect ratio scaling.
    aspectRatio Clay_AspectRatioElementConfig
    // Controls settings related to image elements.
    image Clay_ImageElementConfig
    // Controls whether and how an element "floats", which means it layers over the top of other elements in z order, and doesn't affect the position and size of siblings or parent elements.
    // Note: in order to activate floating, .floating.attachTo must be set to something other than the default value.
    floating Clay_FloatingElementConfig
    // Used to create CUSTOM render commands, usually to render element types not supported by Clay.
    custom Clay_CustomElementConfig
    // Controls whether an element should clip its contents, as well as providing child x,y offset configuration for scrolling.
    clip Clay_ClipElementConfig
    // Controls settings related to element borders, and will generate BORDER render commands.
    border Clay_BorderElementConfig
    // A pointer that will be transparently passed through to resulting render commands.
    userData interface{}
}

// Represents the type of error clay encountered while computing layout.
type Clay_ErrorType uint8

const (
    // A text measurement function wasn't provided using Clay_SetMeasureTextFunction(), or the provided function was null.
    CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED Clay_ErrorType = iota
    // Clay attempted to allocate its internal data structures but ran out of space.
    // The arena passed to Clay_Initialize was created with a capacity smaller than that required by Clay_MinMemorySize().
    CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED
    // Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxElementCount().
    CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED
    // Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxMeasureTextCacheWordCount().
    CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED
    // Two elements were declared with exactly the same ID within one layout.
    CLAY_ERROR_TYPE_DUPLICATE_ID
    // A floating element was declared using CLAY_ATTACH_TO_ELEMENT_ID and either an invalid .parentId was provided or no element with the provided .parentId was found.
    CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND
    // An element was declared that using CLAY_SIZING_PERCENT but the percentage value was over 1. Percentage values are expected to be in the 0-1 range.
    CLAY_ERROR_TYPE_PERCENTAGE_OVER_1
    // Clay encountered an internal error. It would be wonderful if you could report this so we can fix it!
    CLAY_ERROR_TYPE_INTERNAL_ERROR
)

// Data to identify the error that clay has encountered.
type Clay_ErrorData struct {
    // Represents the type of error clay encountered while computing layout.
    // CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED - A text measurement function wasn't provided using Clay_SetMeasureTextFunction(), or the provided function was null.
    // CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED - Clay attempted to allocate its internal data structures but ran out of space. The arena passed to Clay_Initialize was created with a capacity smaller than that required by Clay_MinMemorySize().
    // CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED - Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxElementCount().
    // CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED - Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxMeasureTextCacheWordCount().
    // CLAY_ERROR_TYPE_DUPLICATE_ID - Two elements were declared with exactly the same ID within one layout.
    // CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND - A floating element was declared using CLAY_ATTACH_TO_ELEMENT_ID and either an invalid .parentId was provided or no element with the provided .parentId was found.
    // CLAY_ERROR_TYPE_PERCENTAGE_OVER_1 - An element was declared that using CLAY_SIZING_PERCENT but the percentage value was over 1. Percentage values are expected to be in the 0-1 range.
    // CLAY_ERROR_TYPE_INTERNAL_ERROR - Clay encountered an internal error. It would be wonderful if you could report this so we can fix it!
    errorType Clay_ErrorType
    // A string containing human-readable error text that explains the error in more detail.
    errorText Clay_String
    // A transparent pointer passed through from when the error handler was first provided.
    userData interface{}
}

// A wrapper struct around Clay's error handler function.
type Clay_ErrorHandler struct {
    // A user provided function to call when Clay encounters an error during layout.
    errorHandlerFunction func(errorText Clay_ErrorData)
    // A pointer that will be transparently passed through to the error handler when it is called.
    userData interface{}
}

// End of header section

// -----------------------------------------
// IMPLEMENTATION --------------------------
// -----------------------------------------

const CLAY__NULL = 0
const CLAY__MAXFLOAT = 3.40282346638528859812e+38

// The below comment describes the array macros that were converted to Go
// Array types and functions will be generated as needed when encountering their usage

var Clay__currentContext *Clay_Context
var Clay__defaultMaxElementCount int32 = 8192
var Clay__defaultMaxMeasureTextWordCacheCount int32 = 16384

func Clay__ErrorHandlerFunctionDefault(errorText Clay_ErrorData) {
    _ = errorText
}

var CLAY__SPACECHAR = Clay_String{length: 1, chars: unsafe.Pointer(unsafe.StringData(" "))}
var CLAY__STRING_DEFAULT = Clay_String{length: 0, chars: unsafe.Pointer(unsafe.StringData(""))}

type Clay_BooleanWarnings struct {
    maxElementsExceeded           bool
    maxRenderCommandsExceeded     bool
    maxTextMeasureCacheExceeded   bool
    textMeasurementFunctionNotSet bool
}

type Clay__Warning struct {
    baseMessage    Clay_String
    dynamicMessage Clay_String
}

var CLAY__WARNING_DEFAULT Clay__Warning

type Clay__WarningArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__Warning
}

type Clay_SharedElementConfig struct {
    backgroundColor Clay_Color
    cornerRadius    Clay_CornerRadius
    userData        interface{}
}

// Forward declarations removed - implementations are below

// Array types generated from CLAY__ARRAY_DEFINE macros
type Clay__boolArray struct {
    capacity      int32
    length        int32
    internalArray *bool
}

type Clay__int32_tArray struct {
    capacity      int32
    length        int32
    internalArray *int32
}

type Clay__charArray struct {
    capacity      int32
    length        int32
    internalArray *byte
}

type Clay__LayoutConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_LayoutConfig
}

type Clay__TextElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_TextElementConfig
}

type Clay__AspectRatioElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_AspectRatioElementConfig
}

type Clay__ImageElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_ImageElementConfig
}

type Clay__FloatingElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_FloatingElementConfig
}

type Clay__CustomElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_CustomElementConfig
}

type Clay__ClipElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_ClipElementConfig
}

type Clay__BorderElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_BorderElementConfig
}

type Clay__StringArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_String
}

type Clay__SharedElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_SharedElementConfig
}

// RenderCommandArray functions will be generated later

type Clay__ElementConfigType uint8

const (
    CLAY__ELEMENT_CONFIG_TYPE_NONE Clay__ElementConfigType = iota
    CLAY__ELEMENT_CONFIG_TYPE_BORDER
    CLAY__ELEMENT_CONFIG_TYPE_FLOATING
    CLAY__ELEMENT_CONFIG_TYPE_CLIP
    CLAY__ELEMENT_CONFIG_TYPE_ASPECT
    CLAY__ELEMENT_CONFIG_TYPE_IMAGE
    CLAY__ELEMENT_CONFIG_TYPE_TEXT
    CLAY__ELEMENT_CONFIG_TYPE_CUSTOM
    CLAY__ELEMENT_CONFIG_TYPE_SHARED
)

type Clay_ElementConfigUnion struct {
    textElementConfig        *Clay_TextElementConfig
    aspectRatioElementConfig *Clay_AspectRatioElementConfig
    imageElementConfig       *Clay_ImageElementConfig
    floatingElementConfig    *Clay_FloatingElementConfig
    customElementConfig      *Clay_CustomElementConfig
    clipElementConfig        *Clay_ClipElementConfig
    borderElementConfig      *Clay_BorderElementConfig
    sharedElementConfig      *Clay_SharedElementConfig
}

type Clay_ElementConfig struct {
    _type  Clay__ElementConfigType
    config Clay_ElementConfigUnion
}

type Clay__ElementConfigArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_ElementConfig
}

type Clay__WrappedTextLine struct {
    dimensions Clay_Dimensions
    line       Clay_String
}

type Clay__WrappedTextLineArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__WrappedTextLine
}

type Clay__WrappedTextLineArraySlice struct {
    length        int32
    internalArray *Clay__WrappedTextLine
}

type Clay__TextElementData struct {
    text                Clay_String
    preferredDimensions Clay_Dimensions
    elementIndex        int32
    wrappedLines        Clay__WrappedTextLineArraySlice
}

type Clay__TextElementDataArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__TextElementData
}

type Clay__LayoutElementChildren struct {
    elements *int32
    length   uint16
}

type Clay__ElementConfigArraySlice struct {
    length        int32
    internalArray *Clay_ElementConfig
}

type Clay_LayoutElement struct {
    childrenOrTextContent struct {
        children        Clay__LayoutElementChildren
        textElementData *Clay__TextElementData
    }
    dimensions     Clay_Dimensions
    minDimensions  Clay_Dimensions
    layoutConfig   *Clay_LayoutConfig
    elementConfigs Clay__ElementConfigArraySlice
    id             uint32
}

type Clay_LayoutElementArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_LayoutElement
}

type Clay__ScrollContainerDataInternal struct {
    layoutElement       *Clay_LayoutElement
    boundingBox         Clay_BoundingBox
    contentSize         Clay_Dimensions
    scrollOrigin        Clay_Vector2
    pointerOrigin       Clay_Vector2
    scrollMomentum      Clay_Vector2
    scrollPosition      Clay_Vector2
    previousDelta       Clay_Vector2
    momentumTime        float32
    elementId           uint32
    openThisFrame       bool
    pointerScrollActive bool
}

type Clay__ScrollContainerDataInternalArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__ScrollContainerDataInternal
}

type Clay__DebugElementData struct {
    collision bool
    collapsed bool
}

type Clay__DebugElementDataArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__DebugElementData
}

type Clay_LayoutElementHashMapItem struct { // todo get this struct into a single cache line
    boundingBox           Clay_BoundingBox
    elementId             Clay_ElementId
    layoutElement         *Clay_LayoutElement
    onHoverFunction       func(elementId Clay_ElementId, pointerInfo Clay_PointerData, userData uintptr)
    hoverFunctionUserData uintptr
    nextIndex             int32
    generation            uint32
    idAlias               uint32
    debugData             *Clay__DebugElementData
}

type Clay__LayoutElementHashMapItemArray struct {
    capacity      int32
    length        int32
    internalArray *Clay_LayoutElementHashMapItem
}

type Clay__MeasuredWord struct {
    startOffset int32
    length      int32
    width       float32
    next        int32
}

type Clay__MeasuredWordArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__MeasuredWord
}

type Clay__MeasureTextCacheItem struct {
    unwrappedDimensions     Clay_Dimensions
    measuredWordsStartIndex int32
    minWidth                float32
    containsNewlines        bool
    // Hash map data
    id         uint32
    nextIndex  int32
    generation uint32
}

type Clay__MeasureTextCacheItemArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__MeasureTextCacheItem
}

type Clay__LayoutElementTreeNode struct {
    layoutElement   *Clay_LayoutElement
    position        Clay_Vector2
    nextChildOffset Clay_Vector2
}

type Clay__LayoutElementTreeNodeArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__LayoutElementTreeNode
}

type Clay__LayoutElementTreeRoot struct {
    layoutElementIndex int32
    parentId           uint32         // This can be zero in the case of the root layout tree
    clipElementId      uint32         // This can be zero if there is no clip element
    zIndex             int16
    pointerOffset      Clay_Vector2 // Only used when scroll containers are managed externally
}

type Clay__LayoutElementTreeRootArray struct {
    capacity      int32
    length        int32
    internalArray *Clay__LayoutElementTreeRoot
}

type Clay_Context struct {
    maxElementCount              int32
    maxMeasureTextCacheWordCount int32
    warningsEnabled              bool
    errorHandler                 Clay_ErrorHandler
    booleanWarnings              Clay_BooleanWarnings
    warnings                     Clay__WarningArray

    pointerInfo                    Clay_PointerData
    layoutDimensions               Clay_Dimensions
    dynamicElementIndexBaseHash    Clay_ElementId
    dynamicElementIndex            uint32
    debugModeEnabled               bool
    disableCulling                 bool
    externalScrollHandlingEnabled  bool
    debugSelectedElementId         uint32
    generation                     uint32
    arenaResetOffset               uintptr
    measureTextUserData            interface{}
    queryScrollOffsetUserData      interface{}
    internalArena                  Clay_Arena
    // Layout Elements / Render Commands
    layoutElements              Clay_LayoutElementArray
    renderCommands              Clay_RenderCommandArray
    openLayoutElementStack      Clay__int32_tArray
    layoutElementChildren       Clay__int32_tArray
    layoutElementChildrenBuffer Clay__int32_tArray
    textElementData             Clay__TextElementDataArray
    aspectRatioElementIndexes   Clay__int32_tArray
    reusableElementIndexBuffer  Clay__int32_tArray
    layoutElementClipElementIds Clay__int32_tArray
    // Configs
    layoutConfigs                Clay__LayoutConfigArray
    elementConfigs               Clay__ElementConfigArray
    textElementConfigs           Clay__TextElementConfigArray
    aspectRatioElementConfigs    Clay__AspectRatioElementConfigArray
    imageElementConfigs          Clay__ImageElementConfigArray
    floatingElementConfigs       Clay__FloatingElementConfigArray
    clipElementConfigs           Clay__ClipElementConfigArray
    customElementConfigs         Clay__CustomElementConfigArray
    borderElementConfigs         Clay__BorderElementConfigArray
    sharedElementConfigs         Clay__SharedElementConfigArray
    // Misc Data Structures
    layoutElementIdStrings         Clay__StringArray
    wrappedTextLines               Clay__WrappedTextLineArray
    layoutElementTreeNodeArray1    Clay__LayoutElementTreeNodeArray
    layoutElementTreeRoots         Clay__LayoutElementTreeRootArray
    layoutElementsHashMapInternal  Clay__LayoutElementHashMapItemArray
    layoutElementsHashMap          Clay__int32_tArray
    measureTextHashMapInternal     Clay__MeasureTextCacheItemArray
    measureTextHashMapInternalFreeList Clay__int32_tArray
    measureTextHashMap             Clay__int32_tArray
    measuredWords                  Clay__MeasuredWordArray
    measuredWordsFreeList          Clay__int32_tArray
    openClipElementStack           Clay__int32_tArray
    pointerOverIds                 Clay_ElementIdArray
    scrollContainerDatas           Clay__ScrollContainerDataInternalArray
    treeNodeVisited                Clay__boolArray
    dynamicStringData              Clay__charArray
    debugElementData               Clay__DebugElementDataArray
}

// Array helper functions and default constants
var Clay_TextElementConfig_DEFAULT Clay_TextElementConfig
var Clay_AspectRatioElementConfig_DEFAULT Clay_AspectRatioElementConfig
var Clay_ImageElementConfig_DEFAULT Clay_ImageElementConfig
var Clay_FloatingElementConfig_DEFAULT Clay_FloatingElementConfig
var Clay_CustomElementConfig_DEFAULT Clay_CustomElementConfig
var Clay_ClipElementConfig_DEFAULT Clay_ClipElementConfig
var Clay_BorderElementConfig_DEFAULT Clay_BorderElementConfig
var Clay_SharedElementConfig_DEFAULT Clay_SharedElementConfig
var Clay_LayoutElementHashMapItem_DEFAULT Clay_LayoutElementHashMapItem
var Clay__MeasureTextCacheItem_DEFAULT Clay__MeasureTextCacheItem
var Clay_ElementId_DEFAULT Clay_ElementId
var Clay__CornerRadius_DEFAULT Clay_CornerRadius
var Clay__BorderWidth_DEFAULT Clay_BorderWidth
var Clay_RenderCommandData Clay_RenderCommand

// Array Get functions
func Clay_LayoutElementArray_Get(array *Clay_LayoutElementArray, index int32) *Clay_LayoutElement {
    if index >= 0 && index < array.length {
        return (*Clay_LayoutElement)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay_LayoutElement{})))
    }
    return nil
}

func Clay__int32_tArray_GetValue(array *Clay__int32_tArray, index int32) int32 {
    if index >= 0 && index < array.length {
        return *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(int32(0))))
    }
    return 0
}

func Clay__int32_tArray_Get(array *Clay__int32_tArray, index int32) *int32 {
    if index >= 0 && index < array.length {
        return (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(int32(0))))
    }
    return nil
}

func Clay__int32_tArray_Set(array *Clay__int32_tArray, index int32, value int32) {
    if index >= 0 && index < array.capacity {
        *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(int32(0)))) = value
    }
}

func Clay__int32_tArray_Add(array *Clay__int32_tArray, value int32) *int32 {
    if array.length < array.capacity {
        ptr := (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(int32(0))))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__int32_tArray_RemoveSwapback(array *Clay__int32_tArray, index int32) {
    if index >= 0 && index < array.length {
        array.length--
        if index != array.length {
            *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(int32(0)))) = 
                *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(int32(0))))
        }
    }
}

// Add functions for other array types
func Clay__LayoutConfigArray_Add(array *Clay__LayoutConfigArray, value Clay_LayoutConfig) *Clay_LayoutConfig {
    if array.length < array.capacity {
        ptr := (*Clay_LayoutConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_LayoutConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__TextElementConfigArray_Add(array *Clay__TextElementConfigArray, value Clay_TextElementConfig) *Clay_TextElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_TextElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_TextElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__AspectRatioElementConfigArray_Add(array *Clay__AspectRatioElementConfigArray, value Clay_AspectRatioElementConfig) *Clay_AspectRatioElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_AspectRatioElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_AspectRatioElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__ImageElementConfigArray_Add(array *Clay__ImageElementConfigArray, value Clay_ImageElementConfig) *Clay_ImageElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_ImageElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_ImageElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__FloatingElementConfigArray_Add(array *Clay__FloatingElementConfigArray, value Clay_FloatingElementConfig) *Clay_FloatingElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_FloatingElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_FloatingElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__CustomElementConfigArray_Add(array *Clay__CustomElementConfigArray, value Clay_CustomElementConfig) *Clay_CustomElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_CustomElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_CustomElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__ClipElementConfigArray_Add(array *Clay__ClipElementConfigArray, value Clay_ClipElementConfig) *Clay_ClipElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_ClipElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_ClipElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__BorderElementConfigArray_Add(array *Clay__BorderElementConfigArray, value Clay_BorderElementConfig) *Clay_BorderElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_BorderElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_BorderElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__SharedElementConfigArray_Add(array *Clay__SharedElementConfigArray, value Clay_SharedElementConfig) *Clay_SharedElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_SharedElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_SharedElementConfig{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__ElementConfigArray_Add(array *Clay__ElementConfigArray, value Clay_ElementConfig) *Clay_ElementConfig {
    if array.length < array.capacity {
        ptr := (*Clay_ElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_ElementDeclaration{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__ElementConfigArraySlice_Get(slice *Clay__ElementConfigArraySlice, index int32) *Clay_ElementConfig {
    if index >= 0 && index < slice.length {
        return (*Clay_ElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(slice.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay_ElementDeclaration{})))
    }
    return nil
}

// More array functions
func Clay__MeasuredWordArray_Get(array *Clay__MeasuredWordArray, index int32) *Clay__MeasuredWord {
    if index >= 0 && index < array.length {
        return (*Clay__MeasuredWord)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__MeasuredWord{})))
    }
    return nil
}

func Clay__MeasuredWordArray_Set(array *Clay__MeasuredWordArray, index int32, value Clay__MeasuredWord) {
    if index >= 0 && index < array.capacity {
        *(*Clay__MeasuredWord)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__MeasuredWord{}))) = value
    }
}

func Clay__MeasuredWordArray_Add(array *Clay__MeasuredWordArray, value Clay__MeasuredWord) *Clay__MeasuredWord {
    if array.length < array.capacity {
        ptr := (*Clay__MeasuredWord)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__MeasuredWord{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__MeasureTextCacheItemArray_Get(array *Clay__MeasureTextCacheItemArray, index int32) *Clay__MeasureTextCacheItem {
    if index >= 0 && index < array.length {
        return (*Clay__MeasureTextCacheItem)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__MeasureTextCacheItem{})))
    }
    return nil
}

func Clay__MeasureTextCacheItemArray_Set(array *Clay__MeasureTextCacheItemArray, index int32, value *Clay__MeasureTextCacheItem) {
    if index >= 0 && index < array.capacity {
        ptr := (*Clay__MeasureTextCacheItem)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__MeasureTextCacheItem{})))
        *ptr = *value
        if index >= array.length {
            array.length = index + 1
        }
    }
}

func Clay__MeasureTextCacheItemArray_Add(array *Clay__MeasureTextCacheItemArray, value Clay__MeasureTextCacheItem) *Clay__MeasureTextCacheItem {
    if array.length < array.capacity {
        ptr := (*Clay__MeasureTextCacheItem)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__MeasureTextCacheItem{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__LayoutElementHashMapItemArray_Get(array *Clay__LayoutElementHashMapItemArray, index int32) *Clay_LayoutElementHashMapItem {
    if index >= 0 && index < array.length {
        return (*Clay_LayoutElementHashMapItem)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay_LayoutElementHashMapItem{})))
    }
    return nil
}

func Clay__LayoutElementHashMapItemArray_Add(array *Clay__LayoutElementHashMapItemArray, value Clay_LayoutElementHashMapItem) *Clay_LayoutElementHashMapItem {
    if array.length < array.capacity {
        ptr := (*Clay_LayoutElementHashMapItem)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_LayoutElementHashMapItem{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__DebugElementDataArray_Add(array *Clay__DebugElementDataArray, value Clay__DebugElementData) *Clay__DebugElementData {
    if array.length < array.capacity {
        ptr := (*Clay__DebugElementData)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__DebugElementData{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__StringArray_Add(array *Clay__StringArray, value Clay_String) *Clay_String {
    if array.length < array.capacity {
        ptr := (*Clay_String)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_String{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__StringArray_Set(array *Clay__StringArray, index int32, value *Clay_String) {
    if index >= 0 && index < array.capacity {
        ptr := (*Clay_String)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay_String{})))
        *ptr = *value
        if index >= array.length {
            array.length = index + 1
        }
    }
}

func Clay_LayoutElementArray_Add(array *Clay_LayoutElementArray, value Clay_LayoutElement) *Clay_LayoutElement {
    if array.length < array.capacity {
        ptr := (*Clay_LayoutElement)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_LayoutElement{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__TextElementDataArray_Add(array *Clay__TextElementDataArray, value Clay__TextElementData) *Clay__TextElementData {
    if array.length < array.capacity {
        ptr := (*Clay__TextElementData)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__TextElementData{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__TextElementDataArray_Get(array *Clay__TextElementDataArray, index int32) *Clay__TextElementData {
    if index >= 0 && index < array.length {
        return (*Clay__TextElementData)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__TextElementData{})))
    }
    return nil
}

// Clay__LayoutElementTreeRootArray functions
func Clay__LayoutElementTreeRootArray_Get(array *Clay__LayoutElementTreeRootArray, index int32) *Clay__LayoutElementTreeRoot {
    if index >= 0 && index < array.length {
        return (*Clay__LayoutElementTreeRoot)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__LayoutElementTreeRoot{})))
    }
    return nil
}

func Clay__LayoutElementTreeRootArray_Set(array *Clay__LayoutElementTreeRootArray, index int32, value *Clay__LayoutElementTreeRoot) {
    if index >= 0 && index < array.capacity {
        *(*Clay__LayoutElementTreeRoot)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__LayoutElementTreeRoot{}))) = *value
    }
}

func Clay__LayoutElementTreeRootArray_Add(array *Clay__LayoutElementTreeRootArray, value Clay__LayoutElementTreeRoot) *Clay__LayoutElementTreeRoot {
    if array.length < array.capacity {
        ptr := (*Clay__LayoutElementTreeRoot)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__LayoutElementTreeRoot{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

// Clay__ScrollContainerDataInternalArray functions
func Clay__ScrollContainerDataInternalArray_Get(array *Clay__ScrollContainerDataInternalArray, index int32) *Clay__ScrollContainerDataInternal {
    if index >= 0 && index < array.length {
        return (*Clay__ScrollContainerDataInternal)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__ScrollContainerDataInternal{})))
    }
    return nil
}

func Clay__ScrollContainerDataInternalArray_Add(array *Clay__ScrollContainerDataInternalArray, value Clay__ScrollContainerDataInternal) *Clay__ScrollContainerDataInternal {
    if array.length < array.capacity {
        ptr := (*Clay__ScrollContainerDataInternal)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__ScrollContainerDataInternal{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__ScrollContainerDataInternalArray_RemoveSwapback(array *Clay__ScrollContainerDataInternalArray, index int32) {
    if index >= 0 && index < array.length {
        array.length--
        if index != array.length {
            *(*Clay__ScrollContainerDataInternal)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__ScrollContainerDataInternal{}))) = 
                *(*Clay__ScrollContainerDataInternal)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__ScrollContainerDataInternal{})))
        }
    }
}

// Clay__LayoutElementTreeNodeArray functions
func Clay__LayoutElementTreeNodeArray_Get(array *Clay__LayoutElementTreeNodeArray, index int32) *Clay__LayoutElementTreeNode {
    if index >= 0 && index < array.length {
        return (*Clay__LayoutElementTreeNode)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__LayoutElementTreeNode{})))
    }
    return nil
}

func Clay__LayoutElementTreeNodeArray_Add(array *Clay__LayoutElementTreeNodeArray, value Clay__LayoutElementTreeNode) *Clay__LayoutElementTreeNode {
    if array.length < array.capacity {
        ptr := (*Clay__LayoutElementTreeNode)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__LayoutElementTreeNode{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

// Clay_RenderCommandArray functions
func Clay_RenderCommandArray_Add(array *Clay_RenderCommandArray, value Clay_RenderCommand) *Clay_RenderCommand {
    if array.length < array.capacity {
        ptr := (*Clay_RenderCommand)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_RenderCommand{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

// Clay__WrappedTextLineArray functions
func Clay__WrappedTextLineArray_Add(array *Clay__WrappedTextLineArray, value Clay__WrappedTextLine) *Clay__WrappedTextLine {
    if array.length < array.capacity {
        ptr := (*Clay__WrappedTextLine)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__WrappedTextLine{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay__WrappedTextLineArraySlice_Get(slice *Clay__WrappedTextLineArraySlice, index int32) *Clay__WrappedTextLine {
    if index >= 0 && index < slice.length {
        return (*Clay__WrappedTextLine)(unsafe.Pointer(uintptr(unsafe.Pointer(slice.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay__WrappedTextLine{})))
    }
    return nil
}

// Clay_ElementIdArray functions
func Clay_ElementIdArray_Add(array *Clay_ElementIdArray, value Clay_ElementId) *Clay_ElementId {
    if array.length < array.capacity {
        ptr := (*Clay_ElementId)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay_ElementId{})))
        *ptr = value
        array.length++
        return ptr
    }
    return nil
}

func Clay_ElementIdArray_Get(array *Clay_ElementIdArray, index int32) *Clay_ElementId {
    if index >= 0 && index < array.length {
        return (*Clay_ElementId)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(index)*unsafe.Sizeof(Clay_ElementId{})))
    }
    return nil
}

func Clay__Context_Allocate_Arena(arena *Clay_Arena) *Clay_Context {
    var totalSizeBytes uint64 = uint64(unsafe.Sizeof(Clay_Context{}))
    if totalSizeBytes > arena.capacity {
        return nil
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return (*Clay_Context)(unsafe.Pointer(&arena.memory[0]))
}

func Clay__int32_tArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__int32_tArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(int32(0)))
    if totalSizeBytes > arena.capacity {
        return Clay__int32_tArray{}
    }
    result := Clay__int32_tArray{
        internalArray: (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay_LayoutElementArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay_LayoutElementArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_LayoutElement{}))
    if totalSizeBytes > arena.capacity {
        return Clay_LayoutElementArray{}
    }
    result := Clay_LayoutElementArray{
        internalArray: (*Clay_LayoutElement)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__LayoutConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__LayoutConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_LayoutConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__LayoutConfigArray{}
    }
    result := Clay__LayoutConfigArray{
        internalArray: (*Clay_LayoutConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__ElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__ElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_ElementDeclaration{}))
    if totalSizeBytes > arena.capacity {
        return Clay__ElementConfigArray{}
    }
    result := Clay__ElementConfigArray{
        internalArray: (*Clay_ElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__TextElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__TextElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_TextElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__TextElementConfigArray{}
    }
    result := Clay__TextElementConfigArray{
        internalArray: (*Clay_TextElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__AspectRatioElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__AspectRatioElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_AspectRatioElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__AspectRatioElementConfigArray{}
    }
    result := Clay__AspectRatioElementConfigArray{
        internalArray: (*Clay_AspectRatioElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__ImageElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__ImageElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_ImageElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__ImageElementConfigArray{}
    }
    result := Clay__ImageElementConfigArray{
        internalArray: (*Clay_ImageElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__FloatingElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__FloatingElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_FloatingElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__FloatingElementConfigArray{}
    }
    result := Clay__FloatingElementConfigArray{
        internalArray: (*Clay_FloatingElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__ClipElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__ClipElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_ClipElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__ClipElementConfigArray{}
    }
    result := Clay__ClipElementConfigArray{
        internalArray: (*Clay_ClipElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__CustomElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__CustomElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_CustomElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__CustomElementConfigArray{}
    }
    result := Clay__CustomElementConfigArray{
        internalArray: (*Clay_CustomElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__BorderElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__BorderElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_BorderElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__BorderElementConfigArray{}
    }
    result := Clay__BorderElementConfigArray{
        internalArray: (*Clay_BorderElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__SharedElementConfigArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__SharedElementConfigArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_SharedElementConfig{}))
    if totalSizeBytes > arena.capacity {
        return Clay__SharedElementConfigArray{}
    }
    result := Clay__SharedElementConfigArray{
        internalArray: (*Clay_SharedElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__StringArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__StringArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_String{}))
    if totalSizeBytes > arena.capacity {
        return Clay__StringArray{}
    }
    result := Clay__StringArray{
        internalArray: (*Clay_String)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__WrappedTextLineArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__WrappedTextLineArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__WrappedTextLine{}))
    if totalSizeBytes > arena.capacity {
        return Clay__WrappedTextLineArray{}
    }
    result := Clay__WrappedTextLineArray{
        internalArray: (*Clay__WrappedTextLine)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__LayoutElementTreeNodeArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__LayoutElementTreeNodeArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__LayoutElementTreeNode{}))
    if totalSizeBytes > arena.capacity {
        return Clay__LayoutElementTreeNodeArray{}
    }
    result := Clay__LayoutElementTreeNodeArray{
        internalArray: (*Clay__LayoutElementTreeNode)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__LayoutElementTreeRootArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__LayoutElementTreeRootArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__LayoutElementTreeRoot{}))
    if totalSizeBytes > arena.capacity {
        return Clay__LayoutElementTreeRootArray{}
    }
    result := Clay__LayoutElementTreeRootArray{
        internalArray: (*Clay__LayoutElementTreeRoot)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__TextElementDataArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__TextElementDataArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__TextElementData{}))
    if totalSizeBytes > arena.capacity {
        return Clay__TextElementDataArray{}
    }
    result := Clay__TextElementDataArray{
        internalArray: (*Clay__TextElementData)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay_RenderCommandArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay_RenderCommandArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_RenderCommand{}))
    if totalSizeBytes > arena.capacity {
        return Clay_RenderCommandArray{}
    }
    result := Clay_RenderCommandArray{
        internalArray: (*Clay_RenderCommand)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__boolArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__boolArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(bool(false)))
    if totalSizeBytes > arena.capacity {
        return Clay__boolArray{}
    }
    result := Clay__boolArray{
        internalArray: (*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__charArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__charArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(byte(0)))
    if totalSizeBytes > arena.capacity {
        return Clay__charArray{}
    }
    result := Clay__charArray{
        internalArray: (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__ScrollContainerDataInternalArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__ScrollContainerDataInternalArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__ScrollContainerDataInternal{}))
    if totalSizeBytes > arena.capacity {
        return Clay__ScrollContainerDataInternalArray{}
    }
    result := Clay__ScrollContainerDataInternalArray{
        internalArray: (*Clay__ScrollContainerDataInternal)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__LayoutElementHashMapItemArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__LayoutElementHashMapItemArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_LayoutElementHashMapItem{}))
    if totalSizeBytes > arena.capacity {
        return Clay__LayoutElementHashMapItemArray{}
    }
    result := Clay__LayoutElementHashMapItemArray{
        internalArray: (*Clay_LayoutElementHashMapItem)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__MeasureTextCacheItemArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__MeasureTextCacheItemArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__MeasureTextCacheItem{}))
    if totalSizeBytes > arena.capacity {
        return Clay__MeasureTextCacheItemArray{}
    }
    result := Clay__MeasureTextCacheItemArray{
        internalArray: (*Clay__MeasureTextCacheItem)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__MeasuredWordArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__MeasuredWordArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__MeasuredWord{}))
    if totalSizeBytes > arena.capacity {
        return Clay__MeasuredWordArray{}
    }
    result := Clay__MeasuredWordArray{
        internalArray: (*Clay__MeasuredWord)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay_ElementIdArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay_ElementIdArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay_ElementId{}))
    if totalSizeBytes > arena.capacity {
        return Clay_ElementIdArray{}
    }
    result := Clay_ElementIdArray{
        internalArray: (*Clay_ElementId)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__DebugElementDataArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__DebugElementDataArray {
    totalSizeBytes := uint64(capacity) * uint64(unsafe.Sizeof(Clay__DebugElementData{}))
    if totalSizeBytes > arena.capacity {
        return Clay__DebugElementDataArray{}
    }
    result := Clay__DebugElementDataArray{
        internalArray: (*Clay__DebugElementData)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + arena.nextAllocation)),
        capacity: capacity,
        length: 0,
    }
    arena.nextAllocation += uintptr(totalSizeBytes)
    return result
}

func Clay__WriteStringToCharBuffer(buffer *Clay__charArray, string Clay_String) Clay_String {
    for i := int32(0); i < string.length; i++ {
        *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(buffer.internalArray)) + uintptr(buffer.length+i))) = *(*byte)(unsafe.Pointer(uintptr(string.chars) + uintptr(i)))
    }
    buffer.length += string.length
    return Clay_String{
        length: string.length,
        chars:  string.chars, // Will need proper pointer arithmetic
    }
}

// Function pointers for text measurement and scroll offset
var Clay__MeasureText func(text Clay_StringSlice, config *Clay_TextElementConfig, userData unsafe.Pointer) Clay_Dimensions
var Clay__QueryScrollOffset func(elementId uint32, userData unsafe.Pointer) Clay_Vector2

func Clay__GetOpenLayoutElement() *Clay_LayoutElement {
    context := Clay_GetCurrentContext()
    return Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length-1))
}

func Clay__GetParentElementId() uint32 {
    context := Clay_GetCurrentContext()
    return Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length-2)).id
}

func Clay__StoreLayoutConfig(config Clay_LayoutConfig) *Clay_LayoutConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &CLAY_LAYOUT_DEFAULT
    }
    return Clay__LayoutConfigArray_Add(&context.layoutConfigs, config)
}
func Clay__StoreTextElementConfig(config Clay_TextElementConfig) *Clay_TextElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_TextElementConfig_DEFAULT
    }
    return Clay__TextElementConfigArray_Add(&context.textElementConfigs, config)
}
func Clay__StoreAspectRatioElementConfig(config Clay_AspectRatioElementConfig) *Clay_AspectRatioElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_AspectRatioElementConfig_DEFAULT
    }
    return Clay__AspectRatioElementConfigArray_Add(&context.aspectRatioElementConfigs, config)
}
func Clay__StoreImageElementConfig(config Clay_ImageElementConfig) *Clay_ImageElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_ImageElementConfig_DEFAULT
    }
    return Clay__ImageElementConfigArray_Add(&context.imageElementConfigs, config)
}
func Clay__StoreFloatingElementConfig(config Clay_FloatingElementConfig) *Clay_FloatingElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_FloatingElementConfig_DEFAULT
    }
    return Clay__FloatingElementConfigArray_Add(&context.floatingElementConfigs, config)
}
func Clay__StoreCustomElementConfig(config Clay_CustomElementConfig) *Clay_CustomElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_CustomElementConfig_DEFAULT
    }
    return Clay__CustomElementConfigArray_Add(&context.customElementConfigs, config)
}
func Clay__StoreClipElementConfig(config Clay_ClipElementConfig) *Clay_ClipElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_ClipElementConfig_DEFAULT
    }
    return Clay__ClipElementConfigArray_Add(&context.clipElementConfigs, config)
}
func Clay__StoreBorderElementConfig(config Clay_BorderElementConfig) *Clay_BorderElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_BorderElementConfig_DEFAULT
    }
    return Clay__BorderElementConfigArray_Add(&context.borderElementConfigs, config)
}
func Clay__StoreSharedElementConfig(config Clay_SharedElementConfig) *Clay_SharedElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return &Clay_SharedElementConfig_DEFAULT
    }
    return Clay__SharedElementConfigArray_Add(&context.sharedElementConfigs, config)
}

func Clay__AttachElementConfig(config Clay_ElementConfigUnion, _type Clay__ElementConfigType) Clay_ElementConfig {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return Clay_ElementConfig{}
    }
    openLayoutElement := Clay__GetOpenLayoutElement()
    openLayoutElement.elementConfigs.length++
    return *Clay__ElementConfigArray_Add(&context.elementConfigs, Clay_ElementConfig{_type: _type, config: config})
}

func Clay__FindElementConfigWithType(element *Clay_LayoutElement, _type Clay__ElementConfigType) Clay_ElementConfigUnion {
    for i := int32(0); i < element.elementConfigs.length; i++ {
        config := Clay__ElementConfigArraySlice_Get(&element.elementConfigs, i)
        if config._type == _type {
            return config.config
        }
    }
    return Clay_ElementConfigUnion{}
}

func Clay__HashNumber(offset uint32, seed uint32) Clay_ElementId {
    hash := seed
    hash += (offset + 48)
    hash += (hash << 10)
    hash ^= (hash >> 6)

    hash += (hash << 3)
    hash ^= (hash >> 11)
    hash += (hash << 15)
    return Clay_ElementId{id: hash + 1, offset: offset, baseId: seed, stringId: CLAY__STRING_DEFAULT} // Reserve the hash result of zero as "null id"
}

func Clay__HashString(key Clay_String, seed uint32) Clay_ElementId {
    hash := seed

    for i := int32(0); i < key.length; i++ {
        hash += uint32(*(*byte)(unsafe.Pointer(uintptr(key.chars) + uintptr(i))))
        hash += (hash << 10)
        hash ^= (hash >> 6)
    }

    hash += (hash << 3)
    hash ^= (hash >> 11)
    hash += (hash << 15)
    return Clay_ElementId{id: hash + 1, offset: 0, baseId: hash + 1, stringId: key} // Reserve the hash result of zero as "null id"
}

func Clay__HashStringWithOffset(key Clay_String, offset uint32, seed uint32) Clay_ElementId {
    var hash uint32 = 0
    base := seed

    for i := int32(0); i < key.length; i++ {
        base += uint32(*(*byte)(unsafe.Pointer(uintptr(key.chars) + uintptr(i))))
        base += (base << 10)
        base ^= (base >> 6)
    }
    hash = base
    hash += offset
    hash += (hash << 10)
    hash ^= (hash >> 6)

    hash += (hash << 3)
    base += (base << 3)
    hash ^= (hash >> 11)
    base ^= (base >> 11)
    hash += (hash << 15)
    base += (base << 15)
    return Clay_ElementId{id: hash + 1, offset: offset, baseId: base + 1, stringId: key} // Reserve the hash result of zero as "null id"
}


// Non-SIMD version of hash function
func Clay__HashData(data *uint8, length uintptr) uint64 {
    var hash uint64 = 0

    for i := uintptr(0); i < length; i++ {
        hash += uint64(*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(data)) + i)))
        hash += (hash << 10)
        hash ^= (hash >> 6)
    }
    return hash
}


func Clay__HashStringContentsWithConfig(text *Clay_String, config *Clay_TextElementConfig) uint32 {
    var hash uint32 = 0
    if text.isStaticallyAllocated {
        hash += uint32(uintptr(text.chars))
        hash += (hash << 10)
        hash ^= (hash >> 6)
        hash += uint32(text.length)
        hash += (hash << 10)
        hash ^= (hash >> 6)
    } else {
        hash = uint32(Clay__HashData((*uint8)(text.chars), uintptr(text.length)) % 0xFFFFFFFF)
    }

    hash += uint32(config.fontId)
    hash += (hash << 10)
    hash ^= (hash >> 6)

    hash += uint32(config.fontSize)
    hash += (hash << 10)
    hash ^= (hash >> 6)

    hash += uint32(config.letterSpacing)
    hash += (hash << 10)
    hash ^= (hash >> 6)

    hash += (hash << 3)
    hash ^= (hash >> 11)
    hash += (hash << 15)
    return hash + 1 // Reserve the hash result of zero as "null id"
}

func Clay__AddMeasuredWord(word Clay__MeasuredWord, previousWord *Clay__MeasuredWord) *Clay__MeasuredWord {
    context := Clay_GetCurrentContext()
    if context.measuredWordsFreeList.length > 0 {
        newItemIndex := Clay__int32_tArray_GetValue(&context.measuredWordsFreeList, context.measuredWordsFreeList.length-1)
        context.measuredWordsFreeList.length--
        Clay__MeasuredWordArray_Set(&context.measuredWords, newItemIndex, word)
        previousWord.next = int32(newItemIndex)
        return Clay__MeasuredWordArray_Get(&context.measuredWords, newItemIndex)
    } else {
        previousWord.next = int32(context.measuredWords.length)
        return Clay__MeasuredWordArray_Add(&context.measuredWords, word)
    }
}

func Clay__MeasureTextCached(text *Clay_String, config *Clay_TextElementConfig) *Clay__MeasureTextCacheItem {
    context := Clay_GetCurrentContext()
    if Clay__MeasureText == nil {
        if !context.booleanWarnings.textMeasurementFunctionNotSet {
            context.booleanWarnings.textMeasurementFunctionNotSet = true
            context.errorHandler.errorHandlerFunction(Clay_ErrorData{
                    errorType: CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED,
                    errorText: CLAY_STRING("Clay's internal MeasureText function is null. You may have forgotten to call Clay_SetMeasureTextFunction(), or passed a NULL function pointer by mistake."),
                    userData: context.errorHandler.userData})
        }
        return &Clay__MeasureTextCacheItem_DEFAULT
    }
    id := Clay__HashStringContentsWithConfig(text, config)
    hashBucket := id % uint32(context.maxMeasureTextCacheWordCount/32)
    elementIndexPrevious := int32(0)
    elementIndex := *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.measureTextHashMap.internalArray)) + uintptr(hashBucket)*unsafe.Sizeof(int32(0))))
    for elementIndex != 0 {
        hashEntry := Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, elementIndex)
        if hashEntry.id == id {
            hashEntry.generation = context.generation
            return hashEntry
        }
        // This element hasn't been seen in a few frames, delete the hash map item
        if context.generation-hashEntry.generation > 2 {
            // Add all the measured words that were included in this measurement to the freelist
            nextWordIndex := hashEntry.measuredWordsStartIndex
            for nextWordIndex != -1 {
                measuredWord := Clay__MeasuredWordArray_Get(&context.measuredWords, nextWordIndex)
                Clay__int32_tArray_Add(&context.measuredWordsFreeList, nextWordIndex)
                nextWordIndex = measuredWord.next
            }
            nextIndex := hashEntry.nextIndex
            Clay__MeasureTextCacheItemArray_Set(&context.measureTextHashMapInternal, elementIndex, &Clay__MeasureTextCacheItem{measuredWordsStartIndex: -1})
            Clay__int32_tArray_Add(&context.measureTextHashMapInternalFreeList, elementIndex)
            if elementIndexPrevious == 0 {
                *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.measureTextHashMap.internalArray)) + uintptr(hashBucket)*unsafe.Sizeof(int32(0)))) = nextIndex
            } else {
                previousHashEntry := Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, elementIndexPrevious)
                previousHashEntry.nextIndex = nextIndex
            }
            elementIndex = nextIndex
        } else {
            elementIndexPrevious = elementIndex
            elementIndex = hashEntry.nextIndex
        }
    }
    newItemIndex := int32(0)
    newCacheItem := Clay__MeasureTextCacheItem{measuredWordsStartIndex: -1, id: id, generation: context.generation}
    var measured *Clay__MeasureTextCacheItem
    if context.measureTextHashMapInternalFreeList.length > 0 {
        newItemIndex = Clay__int32_tArray_GetValue(&context.measureTextHashMapInternalFreeList, context.measureTextHashMapInternalFreeList.length-1)
        context.measureTextHashMapInternalFreeList.length--
        Clay__MeasureTextCacheItemArray_Set(&context.measureTextHashMapInternal, newItemIndex, &newCacheItem)
        measured = Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, newItemIndex)
    } else {
        if context.measureTextHashMapInternal.length == context.measureTextHashMapInternal.capacity-1 {
            if !context.booleanWarnings.maxTextMeasureCacheExceeded {
                context.errorHandler.errorHandlerFunction(Clay_ErrorData{
                    errorType: CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
                    errorText: CLAY_STRING("Clay ran out of capacity while attempting to measure text elements. Try using Clay_SetMaxElementCount() with a higher value."),
                    userData: context.errorHandler.userData})
                context.booleanWarnings.maxTextMeasureCacheExceeded = true
            }
            return &Clay__MeasureTextCacheItem_DEFAULT
        }
        measured = Clay__MeasureTextCacheItemArray_Add(&context.measureTextHashMapInternal, newCacheItem)
        newItemIndex = context.measureTextHashMapInternal.length - 1
    }
    start := int32(0)
    end := int32(0)
    lineWidth := float32(0)
    measuredWidth := float32(0)
    measuredHeight := float32(0)
    spaceWidth := Clay__MeasureText(Clay_StringSlice{length: 1, chars: CLAY__SPACECHAR.chars, baseChars: CLAY__SPACECHAR.chars}, config, context.measureTextUserData.(unsafe.Pointer)).width
    tempWord := Clay__MeasuredWord{next: -1}
    previousWord := &tempWord
    for end < text.length {
        if context.measuredWords.length == context.measuredWords.capacity-1 {
            if !context.booleanWarnings.maxTextMeasureCacheExceeded {
                context.errorHandler.errorHandlerFunction(Clay_ErrorData{
                    errorType: CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED,
                    errorText: CLAY_STRING("Clay has run out of space in it's internal text measurement cache. Try using Clay_SetMaxMeasureTextCacheWordCount() (default 16384, with 1 unit storing 1 measured word)."),
                    userData: context.errorHandler.userData})
                context.booleanWarnings.maxTextMeasureCacheExceeded = true
            }
            return &Clay__MeasureTextCacheItem_DEFAULT
        }
        current := *(*byte)(unsafe.Pointer(uintptr(text.chars) + uintptr(end)))
        if current == ' ' || current == '\n' {
            length := end - start
            dimensions := Clay_Dimensions{}
            if length > 0 {
                dimensions = Clay__MeasureText(Clay_StringSlice{length: length, chars: unsafe.Pointer(uintptr(text.chars) + uintptr(start)), baseChars: text.chars}, config, context.measureTextUserData.(unsafe.Pointer))
            }
            measured.minWidth = CLAY__MAX(dimensions.width, measured.minWidth)
            measuredHeight = CLAY__MAX(measuredHeight, dimensions.height)
            if current == ' ' {
                dimensions.width += spaceWidth
                previousWord = Clay__AddMeasuredWord(Clay__MeasuredWord{startOffset: start, length: length + 1, width: dimensions.width, next: -1}, previousWord)
                lineWidth += dimensions.width
            }
            if current == '\n' {
                if length > 0 {
                    previousWord = Clay__AddMeasuredWord(Clay__MeasuredWord{startOffset: start, length: length, width: dimensions.width, next: -1}, previousWord)
                }
                previousWord = Clay__AddMeasuredWord(Clay__MeasuredWord{startOffset: end + 1, length: 0, width: 0, next: -1}, previousWord)
                lineWidth += dimensions.width
                measuredWidth = CLAY__MAX(lineWidth, measuredWidth)
                measured.containsNewlines = true
                lineWidth = 0
            }
            start = end + 1
        }
        end++
    }
    if end-start > 0 {
        dimensions := Clay__MeasureText(Clay_StringSlice{length: end - start, chars: unsafe.Pointer(uintptr(text.chars) + uintptr(start)), baseChars: text.chars}, config, context.measureTextUserData.(unsafe.Pointer))
        Clay__AddMeasuredWord(Clay__MeasuredWord{startOffset: start, length: end - start, width: dimensions.width, next: -1}, previousWord)
        lineWidth += dimensions.width
        measuredHeight = CLAY__MAX(measuredHeight, dimensions.height)
        measured.minWidth = CLAY__MAX(dimensions.width, measured.minWidth)
    }
    measuredWidth = CLAY__MAX(lineWidth, measuredWidth) - float32(config.letterSpacing)
    measured.measuredWordsStartIndex = tempWord.next
    measured.unwrappedDimensions.width = measuredWidth
    measured.unwrappedDimensions.height = measuredHeight
    if elementIndexPrevious != 0 {
        Clay__MeasureTextCacheItemArray_Get(&context.measureTextHashMapInternal, elementIndexPrevious).nextIndex = newItemIndex
    } else {
        *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.measureTextHashMap.internalArray)) + uintptr(hashBucket)*unsafe.Sizeof(int32(0)))) = newItemIndex
    }
    return measured
}

func Clay__PointIsInsideRect(point Clay_Vector2, rect Clay_BoundingBox) bool {
    return point.x >= rect.x && point.x <= rect.x+rect.width && point.y >= rect.y && point.y <= rect.y+rect.height
}

func Clay__AddHashMapItem(elementId Clay_ElementId, layoutElement *Clay_LayoutElement, idAlias uint32) *Clay_LayoutElementHashMapItem {
    context := Clay_GetCurrentContext()
    if context.layoutElementsHashMapInternal.length == context.layoutElementsHashMapInternal.capacity-1 {
        return nil
    }
    item := Clay_LayoutElementHashMapItem{
        elementId:    elementId,
        layoutElement: layoutElement,
        nextIndex:    -1,
        generation:   context.generation + 1,
        idAlias:      idAlias,
    }
    hashBucket := int32(elementId.id) % context.layoutElementsHashMap.capacity
    hashItemPrevious := int32(-1)
    hashItemIndex := *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.layoutElementsHashMap.internalArray)) + uintptr(hashBucket)*unsafe.Sizeof(int32(0))))
    for hashItemIndex != -1 { // Just replace collision, not a big deal - leave it up to the end user
        hashItem := Clay__LayoutElementHashMapItemArray_Get(&context.layoutElementsHashMapInternal, hashItemIndex)
        if hashItem.elementId.id == elementId.id { // Collision - resolve based on generation
            item.nextIndex = hashItem.nextIndex
            if hashItem.generation <= context.generation { // First collision - assume this is the "same" element
                hashItem.elementId = elementId // Make sure to copy this across. If the stringId reference has changed, we should update the hash item to use the new one.
                hashItem.idAlias = idAlias
                hashItem.generation = context.generation + 1
                hashItem.layoutElement = layoutElement
                hashItem.debugData.collision = false
                hashItem.onHoverFunction = nil
                hashItem.hoverFunctionUserData = 0
            } else { // Multiple collisions this frame - two elements have the same ID
                context.errorHandler.errorHandlerFunction(Clay_ErrorData{
                    errorType: CLAY_ERROR_TYPE_DUPLICATE_ID,
                    errorText: CLAY_STRING("An element with this ID was already previously declared during this layout."),
                    userData:  context.errorHandler.userData,
                })
                if context.debugModeEnabled {
                    hashItem.debugData.collision = true
                }
            }
            return hashItem
        }
        hashItemPrevious = hashItemIndex
        hashItemIndex = hashItem.nextIndex
    }
    hashItem := Clay__LayoutElementHashMapItemArray_Add(&context.layoutElementsHashMapInternal, item)
    hashItem.debugData = Clay__DebugElementDataArray_Add(&context.debugElementData, Clay__DebugElementData{})
    if hashItemPrevious != -1 {
        Clay__LayoutElementHashMapItemArray_Get(&context.layoutElementsHashMapInternal, hashItemPrevious).nextIndex = int32(context.layoutElementsHashMapInternal.length - 1)
    } else {
        *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.layoutElementsHashMap.internalArray)) + uintptr(hashBucket)*unsafe.Sizeof(int32(0)))) = int32(context.layoutElementsHashMapInternal.length - 1)
    }
    return hashItem
}

func Clay__GetHashMapItem(id uint32) *Clay_LayoutElementHashMapItem {
    context := Clay_GetCurrentContext()
    hashBucket := int32(id) % context.layoutElementsHashMap.capacity
    elementIndex := *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.layoutElementsHashMap.internalArray)) + uintptr(hashBucket)*unsafe.Sizeof(int32(0))))
    for elementIndex != -1 {
        hashEntry := Clay__LayoutElementHashMapItemArray_Get(&context.layoutElementsHashMapInternal, elementIndex)
        if hashEntry.elementId.id == id {
            return hashEntry
        }
        elementIndex = hashEntry.nextIndex
    }
    return &Clay_LayoutElementHashMapItem_DEFAULT
}

func Clay__GenerateIdForAnonymousElement(openLayoutElement *Clay_LayoutElement) Clay_ElementId {
    context := Clay_GetCurrentContext()
    parentElement := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length-2))
    elementId := Clay__HashNumber(uint32(parentElement.childrenOrTextContent.children.length), parentElement.id)
    openLayoutElement.id = elementId.id
    Clay__AddHashMapItem(elementId, openLayoutElement, 0)
    Clay__StringArray_Add(&context.layoutElementIdStrings, elementId.stringId)
    return elementId
}

func Clay__ElementHasConfig(layoutElement *Clay_LayoutElement, _type Clay__ElementConfigType) bool {
    for i := int32(0); i < layoutElement.elementConfigs.length; i++ {
        if Clay__ElementConfigArraySlice_Get(&layoutElement.elementConfigs, i)._type == _type {
            return true;
        }
    }
    return false;
}

func Clay__UpdateAspectRatioBox(layoutElement *Clay_LayoutElement) {
    for j := int32(0); j < layoutElement.elementConfigs.length; j++ {
        config := Clay__ElementConfigArraySlice_Get(&layoutElement.elementConfigs, j)
        if config._type == CLAY__ELEMENT_CONFIG_TYPE_ASPECT {
            aspectConfig := config.config.aspectRatioElementConfig
            if aspectConfig.aspectRatio == 0 {
                break
            }
            if layoutElement.dimensions.width == 0 && layoutElement.dimensions.height != 0 {
                layoutElement.dimensions.width = layoutElement.dimensions.height * aspectConfig.aspectRatio
            } else if layoutElement.dimensions.width != 0 && layoutElement.dimensions.height == 0 {
                layoutElement.dimensions.height = layoutElement.dimensions.width * (1 / aspectConfig.aspectRatio)
            }
            break
        }
    }
}

func Clay__CloseElement() {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return
    }
    openLayoutElement := Clay__GetOpenLayoutElement()
    layoutConfig := openLayoutElement.layoutConfig
    elementHasClipHorizontal := false
    elementHasClipVertical := false
    for i := int32(0); i < openLayoutElement.elementConfigs.length; i++ {
        config := Clay__ElementConfigArraySlice_Get(&openLayoutElement.elementConfigs, i)
        if config._type == CLAY__ELEMENT_CONFIG_TYPE_CLIP {
            elementHasClipHorizontal = config.config.clipElementConfig.horizontal
            elementHasClipVertical = config.config.clipElementConfig.vertical
            context.openClipElementStack.length--
            break
        } else if config._type == CLAY__ELEMENT_CONFIG_TYPE_FLOATING {
            context.openClipElementStack.length--
        }
    }

    leftRightPadding := float32(layoutConfig.padding.left + layoutConfig.padding.right)
    topBottomPadding := float32(layoutConfig.padding.top + layoutConfig.padding.bottom)

    // Attach children to the current open element
    openLayoutElement.childrenOrTextContent.children.elements = (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.layoutElementChildren.internalArray)) + uintptr(context.layoutElementChildren.length)*unsafe.Sizeof(int32(0))))
    if layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT {
        openLayoutElement.dimensions.width = leftRightPadding
        openLayoutElement.minDimensions.width = leftRightPadding
        for i := int32(0); i < int32(openLayoutElement.childrenOrTextContent.children.length); i++ {
            childIndex := Clay__int32_tArray_GetValue(&context.layoutElementChildrenBuffer, context.layoutElementChildrenBuffer.length-int32(openLayoutElement.childrenOrTextContent.children.length)+i)
            child := Clay_LayoutElementArray_Get(&context.layoutElements, childIndex)
            openLayoutElement.dimensions.width += child.dimensions.width
            openLayoutElement.dimensions.height = CLAY__MAX(openLayoutElement.dimensions.height, child.dimensions.height+topBottomPadding)
            // Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
            if !elementHasClipHorizontal {
                openLayoutElement.minDimensions.width += child.minDimensions.width
            }
            if !elementHasClipVertical {
                openLayoutElement.minDimensions.height = CLAY__MAX(openLayoutElement.minDimensions.height, child.minDimensions.height+topBottomPadding)
            }
            Clay__int32_tArray_Add(&context.layoutElementChildren, childIndex)
        }
        childGap := CLAY__MAX(float32(int32(openLayoutElement.childrenOrTextContent.children.length)-1), 0) * float32(layoutConfig.childGap)
        openLayoutElement.dimensions.width += childGap
        if !elementHasClipHorizontal {
            openLayoutElement.minDimensions.width += childGap
        }
    } else if layoutConfig.layoutDirection == CLAY_TOP_TO_BOTTOM {
        openLayoutElement.dimensions.height = topBottomPadding
        openLayoutElement.minDimensions.height = topBottomPadding
        for i := int32(0); i < int32(openLayoutElement.childrenOrTextContent.children.length); i++ {
            childIndex := Clay__int32_tArray_GetValue(&context.layoutElementChildrenBuffer, context.layoutElementChildrenBuffer.length-int32(openLayoutElement.childrenOrTextContent.children.length)+i)
            child := Clay_LayoutElementArray_Get(&context.layoutElements, childIndex)
            openLayoutElement.dimensions.height += child.dimensions.height
            openLayoutElement.dimensions.width = CLAY__MAX(openLayoutElement.dimensions.width, child.dimensions.width+leftRightPadding)
            // Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
            if !elementHasClipVertical {
                openLayoutElement.minDimensions.height += child.minDimensions.height
            }
            if !elementHasClipHorizontal {
                openLayoutElement.minDimensions.width = CLAY__MAX(openLayoutElement.minDimensions.width, child.minDimensions.width+leftRightPadding)
            }
            Clay__int32_tArray_Add(&context.layoutElementChildren, childIndex)
        }
        childGap := CLAY__MAX(float32(int32(openLayoutElement.childrenOrTextContent.children.length)-1), 0) * float32(layoutConfig.childGap)
        openLayoutElement.dimensions.height += childGap
        if !elementHasClipVertical {
            openLayoutElement.minDimensions.height += childGap
        }
    }

    context.layoutElementChildrenBuffer.length -= int32(openLayoutElement.childrenOrTextContent.children.length)

    // Clamp element min and max width to the values configured in the layout
    if layoutConfig.sizing.width.type_ != CLAY__SIZING_TYPE_PERCENT {
        if layoutConfig.sizing.width.size.minMax.max <= 0 { // Set the max size if the user didn't specify, makes calculations easier
            layoutConfig.sizing.width.size.minMax.max = CLAY__MAXFLOAT
        }
        openLayoutElement.dimensions.width = CLAY__MIN(CLAY__MAX(openLayoutElement.dimensions.width, layoutConfig.sizing.width.size.minMax.min), layoutConfig.sizing.width.size.minMax.max)
        openLayoutElement.minDimensions.width = CLAY__MIN(CLAY__MAX(openLayoutElement.minDimensions.width, layoutConfig.sizing.width.size.minMax.min), layoutConfig.sizing.width.size.minMax.max)
    } else {
        openLayoutElement.dimensions.width = 0
    }

    // Clamp element min and max height to the values configured in the layout
    if layoutConfig.sizing.height.type_ != CLAY__SIZING_TYPE_PERCENT {
        if layoutConfig.sizing.height.size.minMax.max <= 0 { // Set the max size if the user didn't specify, makes calculations easier
            layoutConfig.sizing.height.size.minMax.max = CLAY__MAXFLOAT
        }
        openLayoutElement.dimensions.height = CLAY__MIN(CLAY__MAX(openLayoutElement.dimensions.height, layoutConfig.sizing.height.size.minMax.min), layoutConfig.sizing.height.size.minMax.max)
        openLayoutElement.minDimensions.height = CLAY__MIN(CLAY__MAX(openLayoutElement.minDimensions.height, layoutConfig.sizing.height.size.minMax.min), layoutConfig.sizing.height.size.minMax.max)
    } else {
        openLayoutElement.dimensions.height = 0
    }

    Clay__UpdateAspectRatioBox(openLayoutElement)

    elementIsFloating := Clay__ElementHasConfig(openLayoutElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)

    // Close the currently open element
    closingElementIndex := Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length-1)
    Clay__int32_tArray_RemoveSwapback(&context.openLayoutElementStack, context.openLayoutElementStack.length-1)
    openLayoutElement = Clay__GetOpenLayoutElement()

    if !elementIsFloating && context.openLayoutElementStack.length > 1 {
        openLayoutElement.childrenOrTextContent.children.length++
        Clay__int32_tArray_Add(&context.layoutElementChildrenBuffer, closingElementIndex)
    }
}

// Clay__MemCmp will be implemented with SIMD support based on build tags
func Clay__MemCmp(s1, s2 string, length int32) bool {
    // Simple byte comparison for Go
    if int32(len(s1)) < length || int32(len(s2)) < length {
        return false
    }
    for i := int32(0); i < length; i++ {
        if s1[i] != s2[i] {
            return false
        }
    }
    return true
}

func Clay__OpenElement() {
    context := Clay_GetCurrentContext()
    if context.layoutElements.length == context.layoutElements.capacity-1 || context.booleanWarnings.maxElementsExceeded {
        context.booleanWarnings.maxElementsExceeded = true
        return
    }
    layoutElement := Clay_LayoutElement{}
    Clay_LayoutElementArray_Add(&context.layoutElements, layoutElement)
    Clay__int32_tArray_Add(&context.openLayoutElementStack, context.layoutElements.length-1)
    if context.openClipElementStack.length > 0 {
        Clay__int32_tArray_Set(&context.layoutElementClipElementIds, context.layoutElements.length-1, Clay__int32_tArray_GetValue(&context.openClipElementStack, context.openClipElementStack.length-1))
    } else {
        Clay__int32_tArray_Set(&context.layoutElementClipElementIds, context.layoutElements.length-1, 0)
    }
}

func Clay__OpenTextElement(text Clay_String, textConfig *Clay_TextElementConfig) {
    context := Clay_GetCurrentContext()
    if context.layoutElements.length == context.layoutElements.capacity-1 || context.booleanWarnings.maxElementsExceeded {
        context.booleanWarnings.maxElementsExceeded = true
        return
    }
    parentElement := Clay__GetOpenLayoutElement()

    layoutElement := Clay_LayoutElement{}
    textElement := Clay_LayoutElementArray_Add(&context.layoutElements, layoutElement)
    if context.openClipElementStack.length > 0 {
        Clay__int32_tArray_Set(&context.layoutElementClipElementIds, context.layoutElements.length-1, Clay__int32_tArray_GetValue(&context.openClipElementStack, context.openClipElementStack.length-1))
    } else {
        Clay__int32_tArray_Set(&context.layoutElementClipElementIds, context.layoutElements.length-1, 0)
    }

    Clay__int32_tArray_Add(&context.layoutElementChildrenBuffer, context.layoutElements.length-1)
    textMeasured := Clay__MeasureTextCached(&text, textConfig)
    elementId := Clay__HashNumber(uint32(parentElement.childrenOrTextContent.children.length), parentElement.id)
    textElement.id = elementId.id
    Clay__AddHashMapItem(elementId, textElement, 0)
    Clay__StringArray_Add(&context.layoutElementIdStrings, elementId.stringId)
    var textDimensions Clay_Dimensions
    if textConfig.lineHeight > 0 {
        textDimensions = Clay_Dimensions{width: textMeasured.unwrappedDimensions.width, height: float32(textConfig.lineHeight)}
    } else {
        textDimensions = Clay_Dimensions{width: textMeasured.unwrappedDimensions.width, height: textMeasured.unwrappedDimensions.height}
    }
    textElement.dimensions = textDimensions
    textElement.minDimensions = Clay_Dimensions{width: textMeasured.minWidth, height: textDimensions.height}
    textElement.childrenOrTextContent.textElementData = Clay__TextElementDataArray_Add(&context.textElementData, Clay__TextElementData{text: text, preferredDimensions: textMeasured.unwrappedDimensions, elementIndex: context.layoutElements.length - 1})
    textElement.elementConfigs = Clay__ElementConfigArraySlice{
        length: 1,
        internalArray: Clay__ElementConfigArray_Add(&context.elementConfigs, Clay_ElementConfig{_type: CLAY__ELEMENT_CONFIG_TYPE_TEXT, config: Clay_ElementConfigUnion{textElementConfig: textConfig}}),
    }
    textElement.layoutConfig = &CLAY_LAYOUT_DEFAULT
    parentElement.childrenOrTextContent.children.length++
}

func Clay__AttachId(elementId Clay_ElementId) Clay_ElementId {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return Clay_ElementId_DEFAULT
    }
    openLayoutElement := Clay__GetOpenLayoutElement()
    idAlias := openLayoutElement.id
    openLayoutElement.id = elementId.id
    Clay__AddHashMapItem(elementId, openLayoutElement, idAlias)
    Clay__StringArray_Set(&context.layoutElementIdStrings, context.layoutElements.length-1, &elementId.stringId)
    return elementId
}

func Clay__ConfigureOpenElementPtr(declaration *Clay_ElementDeclaration) {
    context := Clay_GetCurrentContext()
    openLayoutElement := Clay__GetOpenLayoutElement()
    openLayoutElement.layoutConfig = Clay__StoreLayoutConfig(declaration.layout)
    if (declaration.layout.sizing.width.type_ == CLAY__SIZING_TYPE_PERCENT && declaration.layout.sizing.width.size.percent > 1) || (declaration.layout.sizing.height.type_ == CLAY__SIZING_TYPE_PERCENT && declaration.layout.sizing.height.size.percent > 1) {
        context.errorHandler.errorHandlerFunction(Clay_ErrorData{
            errorType: CLAY_ERROR_TYPE_PERCENTAGE_OVER_1,
            errorText: CLAY_STRING("An element was configured with CLAY_SIZING_PERCENT, but the provided percentage value was over 1.0. Clay expects a value between 0 and 1, i.e. 20% is 0.2."),
            userData: context.errorHandler.userData,
        })
    }

    openLayoutElementId := declaration.id

    openLayoutElement.elementConfigs.internalArray = (*Clay_ElementConfig)(unsafe.Pointer(uintptr(unsafe.Pointer(context.elementConfigs.internalArray)) + uintptr(context.elementConfigs.length)*unsafe.Sizeof(Clay_ElementDeclaration{})))
    var sharedConfig *Clay_SharedElementConfig
    if declaration.backgroundColor.a > 0 {
        sharedConfig = Clay__StoreSharedElementConfig(Clay_SharedElementConfig{backgroundColor: declaration.backgroundColor})
        Clay__AttachElementConfig(Clay_ElementConfigUnion{sharedElementConfig: sharedConfig}, CLAY__ELEMENT_CONFIG_TYPE_SHARED)
    }
    if !Clay__MemCmp(string(unsafe.Slice((*byte)(unsafe.Pointer(&declaration.cornerRadius)), 16)), string(unsafe.Slice((*byte)(unsafe.Pointer(&Clay__CornerRadius_DEFAULT)), 16)), 16) {
        if sharedConfig != nil {
            sharedConfig.cornerRadius = declaration.cornerRadius
        } else {
            sharedConfig = Clay__StoreSharedElementConfig(Clay_SharedElementConfig{cornerRadius: declaration.cornerRadius})
            Clay__AttachElementConfig(Clay_ElementConfigUnion{sharedElementConfig: sharedConfig}, CLAY__ELEMENT_CONFIG_TYPE_SHARED)
        }
    }
    if declaration.userData != 0 {
        if sharedConfig != nil {
            sharedConfig.userData = declaration.userData
        } else {
            sharedConfig = Clay__StoreSharedElementConfig(Clay_SharedElementConfig{userData: declaration.userData})
            Clay__AttachElementConfig(Clay_ElementConfigUnion{sharedElementConfig: sharedConfig}, CLAY__ELEMENT_CONFIG_TYPE_SHARED)
        }
    }
    if declaration.image.imageData != nil {
        Clay__AttachElementConfig(Clay_ElementConfigUnion{imageElementConfig: Clay__StoreImageElementConfig(declaration.image)}, CLAY__ELEMENT_CONFIG_TYPE_IMAGE)
    }
    if declaration.aspectRatio.aspectRatio > 0 {
        Clay__AttachElementConfig(Clay_ElementConfigUnion{aspectRatioElementConfig: Clay__StoreAspectRatioElementConfig(declaration.aspectRatio)}, CLAY__ELEMENT_CONFIG_TYPE_ASPECT)
        Clay__int32_tArray_Add(&context.aspectRatioElementIndexes, context.layoutElements.length-1)
    }
    if declaration.floating.attachTo != CLAY_ATTACH_TO_NONE {
        floatingConfig := declaration.floating
        // This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here
        hierarchicalParent := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length-2))
        if hierarchicalParent != nil {
            clipElementId := uint32(0)
            if declaration.floating.attachTo == CLAY_ATTACH_TO_PARENT {
                // Attach to the element's direct hierarchical parent
                floatingConfig.parentId = hierarchicalParent.id
                if context.openClipElementStack.length > 0 {
                    clipElementId = uint32(Clay__int32_tArray_GetValue(&context.openClipElementStack, context.openClipElementStack.length-1))
                }
            } else if declaration.floating.attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID {
                parentItem := Clay__GetHashMapItem(floatingConfig.parentId)
                if parentItem == &Clay_LayoutElementHashMapItem_DEFAULT {
                    context.errorHandler.errorHandlerFunction(Clay_ErrorData{
                        errorType: CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND,
                        errorText: CLAY_STRING("A floating element was declared with a parentId, but no element with that ID was found."),
                        userData: context.errorHandler.userData,
                    })
                } else {
                    clipElementId = uint32(Clay__int32_tArray_GetValue(&context.layoutElementClipElementIds, int32(uintptr(unsafe.Pointer(parentItem.layoutElement))-uintptr(unsafe.Pointer(context.layoutElements.internalArray)))))
                }
            } else if declaration.floating.attachTo == CLAY_ATTACH_TO_ROOT {
                floatingConfig.parentId = Clay__HashString(CLAY_STRING("Clay__RootContainer"), 0).id
            }
            if openLayoutElementId.id == 0 {
                openLayoutElementId = Clay__HashStringWithOffset(CLAY_STRING("Clay__FloatingContainer"), uint32(context.layoutElementTreeRoots.length), 0)
            }
            if declaration.floating.clipTo == CLAY_CLIP_TO_NONE {
                clipElementId = 0
            }
            currentElementIndex := Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length-1)
            Clay__int32_tArray_Set(&context.layoutElementClipElementIds, currentElementIndex, int32(clipElementId))
            Clay__int32_tArray_Add(&context.openClipElementStack, int32(clipElementId))
            Clay__LayoutElementTreeRootArray_Add(&context.layoutElementTreeRoots, Clay__LayoutElementTreeRoot{
                layoutElementIndex: Clay__int32_tArray_GetValue(&context.openLayoutElementStack, context.openLayoutElementStack.length-1),
                parentId: floatingConfig.parentId,
                clipElementId: clipElementId,
                zIndex: floatingConfig.zIndex,
            })
            Clay__AttachElementConfig(Clay_ElementConfigUnion{floatingElementConfig: Clay__StoreFloatingElementConfig(floatingConfig)}, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)
        }
    }
    if declaration.custom.customData != nil {
        Clay__AttachElementConfig(Clay_ElementConfigUnion{customElementConfig: Clay__StoreCustomElementConfig(declaration.custom)}, CLAY__ELEMENT_CONFIG_TYPE_CUSTOM)
    }

    if openLayoutElementId.id != 0 {
        Clay__AttachId(openLayoutElementId)
    } else if openLayoutElement.id == 0 {
        openLayoutElementId = Clay__GenerateIdForAnonymousElement(openLayoutElement)
    }

    if declaration.clip.horizontal || declaration.clip.vertical {
        Clay__AttachElementConfig(Clay_ElementConfigUnion{clipElementConfig: Clay__StoreClipElementConfig(declaration.clip)}, CLAY__ELEMENT_CONFIG_TYPE_CLIP)
        Clay__int32_tArray_Add(&context.openClipElementStack, int32(openLayoutElement.id))
        // Retrieve or create cached data to track scroll position across frames
        var scrollOffset *Clay__ScrollContainerDataInternal
        for i := int32(0); i < context.scrollContainerDatas.length; i++ {
            mapping := Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i)
            if openLayoutElement.id == mapping.elementId {
                scrollOffset = mapping
                scrollOffset.layoutElement = openLayoutElement
                scrollOffset.openThisFrame = true
            }
        }
        if scrollOffset == nil {
            scrollOffset = Clay__ScrollContainerDataInternalArray_Add(&context.scrollContainerDatas, Clay__ScrollContainerDataInternal{layoutElement: openLayoutElement, scrollOrigin: Clay_Vector2{x: -1, y: -1}, elementId: openLayoutElement.id, openThisFrame: true})
        }
        if context.externalScrollHandlingEnabled {
            scrollOffset.scrollPosition = Clay__QueryScrollOffset(scrollOffset.elementId, context.queryScrollOffsetUserData.(unsafe.Pointer))
        }
    }
    if !Clay__MemCmp(string(unsafe.Slice((*byte)(unsafe.Pointer(&declaration.border.width)), 10)), string(unsafe.Slice((*byte)(unsafe.Pointer(&Clay__BorderWidth_DEFAULT)), 10)), 10) {
        Clay__AttachElementConfig(Clay_ElementConfigUnion{borderElementConfig: Clay__StoreBorderElementConfig(declaration.border)}, CLAY__ELEMENT_CONFIG_TYPE_BORDER)
    }
}

func Clay__ConfigureOpenElement(declaration Clay_ElementDeclaration) {
    Clay__ConfigureOpenElementPtr(&declaration)
}

func Clay__InitializeEphemeralMemory(context *Clay_Context) {
    maxElementCount := context.maxElementCount
    // Ephemeral Memory - reset every frame
    arena := &context.internalArena
    arena.nextAllocation = context.arenaResetOffset

    context.layoutElementChildrenBuffer = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.layoutElements = Clay_LayoutElementArray_Allocate_Arena(arena, maxElementCount)
    context.warnings = Clay__WarningArray_Allocate_Arena(arena, 100)

    context.layoutConfigs = Clay__LayoutConfigArray_Allocate_Arena(arena, maxElementCount)
    context.elementConfigs = Clay__ElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.textElementConfigs = Clay__TextElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.aspectRatioElementConfigs = Clay__AspectRatioElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.imageElementConfigs = Clay__ImageElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.floatingElementConfigs = Clay__FloatingElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.clipElementConfigs = Clay__ClipElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.customElementConfigs = Clay__CustomElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.borderElementConfigs = Clay__BorderElementConfigArray_Allocate_Arena(arena, maxElementCount)
    context.sharedElementConfigs = Clay__SharedElementConfigArray_Allocate_Arena(arena, maxElementCount)

    context.layoutElementIdStrings = Clay__StringArray_Allocate_Arena(arena, maxElementCount)
    context.wrappedTextLines = Clay__WrappedTextLineArray_Allocate_Arena(arena, maxElementCount)
    context.layoutElementTreeNodeArray1 = Clay__LayoutElementTreeNodeArray_Allocate_Arena(arena, maxElementCount)
    context.layoutElementTreeRoots = Clay__LayoutElementTreeRootArray_Allocate_Arena(arena, maxElementCount)
    context.layoutElementChildren = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.openLayoutElementStack = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.textElementData = Clay__TextElementDataArray_Allocate_Arena(arena, maxElementCount)
    context.aspectRatioElementIndexes = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.renderCommands = Clay_RenderCommandArray_Allocate_Arena(arena, maxElementCount)
    context.treeNodeVisited = Clay__boolArray_Allocate_Arena(arena, maxElementCount)
    context.treeNodeVisited.length = context.treeNodeVisited.capacity // This array is accessed directly rather than behaving as a list
    context.openClipElementStack = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.reusableElementIndexBuffer = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.layoutElementClipElementIds = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.dynamicStringData = Clay__charArray_Allocate_Arena(arena, maxElementCount)
}

func Clay__InitializePersistentMemory(context *Clay_Context) {
    // Persistent memory - initialized once and not reset
    maxElementCount := context.maxElementCount
    maxMeasureTextCacheWordCount := context.maxMeasureTextCacheWordCount
    arena := &context.internalArena

    context.scrollContainerDatas = Clay__ScrollContainerDataInternalArray_Allocate_Arena(arena, 100)
    context.layoutElementsHashMapInternal = Clay__LayoutElementHashMapItemArray_Allocate_Arena(arena, maxElementCount)
    context.layoutElementsHashMap = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.measureTextHashMapInternal = Clay__MeasureTextCacheItemArray_Allocate_Arena(arena, maxElementCount)
    context.measureTextHashMapInternalFreeList = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.measuredWordsFreeList = Clay__int32_tArray_Allocate_Arena(arena, maxMeasureTextCacheWordCount)
    context.measureTextHashMap = Clay__int32_tArray_Allocate_Arena(arena, maxElementCount)
    context.measuredWords = Clay__MeasuredWordArray_Allocate_Arena(arena, maxMeasureTextCacheWordCount)
    context.pointerOverIds = Clay_ElementIdArray_Allocate_Arena(arena, maxElementCount)
    context.debugElementData = Clay__DebugElementDataArray_Allocate_Arena(arena, maxElementCount)
    context.arenaResetOffset = arena.nextAllocation
}

const CLAY__EPSILON float32 = 0.01

func Clay__FloatEqual(left, right float32) bool {
    subtracted := left - right
    return subtracted < CLAY__EPSILON && subtracted > -CLAY__EPSILON
}

func Clay__SizeContainersAlongAxis(xAxis bool) {
    context := Clay_GetCurrentContext()
    bfsBuffer := context.layoutElementChildrenBuffer
    resizableContainerBuffer := context.openLayoutElementStack
    for rootIndex := int32(0); rootIndex < context.layoutElementTreeRoots.length; rootIndex++ {
        bfsBuffer.length = 0
        root := Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex)
        rootElement := Clay_LayoutElementArray_Get(&context.layoutElements, root.layoutElementIndex)
        Clay__int32_tArray_Add(&bfsBuffer, root.layoutElementIndex)

        // Size floating containers to their parents
        if Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) {
            floatingElementConfig := Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig
            parentItem := Clay__GetHashMapItem(floatingElementConfig.parentId)
            if parentItem != nil && parentItem != &Clay_LayoutElementHashMapItem_DEFAULT {
                parentLayoutElement := parentItem.layoutElement
                switch rootElement.layoutConfig.sizing.width.type_ {
                case CLAY__SIZING_TYPE_GROW:
                    rootElement.dimensions.width = parentLayoutElement.dimensions.width
                case CLAY__SIZING_TYPE_PERCENT:
                    rootElement.dimensions.width = parentLayoutElement.dimensions.width * rootElement.layoutConfig.sizing.width.size.percent
                }
                switch rootElement.layoutConfig.sizing.height.type_ {
                case CLAY__SIZING_TYPE_GROW:
                    rootElement.dimensions.height = parentLayoutElement.dimensions.height
                case CLAY__SIZING_TYPE_PERCENT:
                    rootElement.dimensions.height = parentLayoutElement.dimensions.height * rootElement.layoutConfig.sizing.height.size.percent
                }
            }
        }

        if rootElement.layoutConfig.sizing.width.type_ != CLAY__SIZING_TYPE_PERCENT {
            rootElement.dimensions.width = CLAY__MIN(CLAY__MAX(rootElement.dimensions.width, rootElement.layoutConfig.sizing.width.size.minMax.min), rootElement.layoutConfig.sizing.width.size.minMax.max)
        }
        if rootElement.layoutConfig.sizing.height.type_ != CLAY__SIZING_TYPE_PERCENT {
            rootElement.dimensions.height = CLAY__MIN(CLAY__MAX(rootElement.dimensions.height, rootElement.layoutConfig.sizing.height.size.minMax.min), rootElement.layoutConfig.sizing.height.size.minMax.max)
        }

        for i := int32(0); i < bfsBuffer.length; i++ {
            parentIndex := Clay__int32_tArray_GetValue(&bfsBuffer, i)
            parent := Clay_LayoutElementArray_Get(&context.layoutElements, parentIndex)
            parentStyleConfig := parent.layoutConfig
            growContainerCount := int32(0)
            var parentSize float32
            if xAxis {
                parentSize = parent.dimensions.width
            } else {
                parentSize = parent.dimensions.height
            }
            var parentPadding float32
            if xAxis {
                parentPadding = float32(parent.layoutConfig.padding.left + parent.layoutConfig.padding.right)
            } else {
                parentPadding = float32(parent.layoutConfig.padding.top + parent.layoutConfig.padding.bottom)
            }
            innerContentSize := float32(0)
            totalPaddingAndChildGaps := parentPadding
            sizingAlongAxis := (xAxis && parentStyleConfig.layoutDirection == CLAY_LEFT_TO_RIGHT) || (!xAxis && parentStyleConfig.layoutDirection == CLAY_TOP_TO_BOTTOM)
            resizableContainerBuffer.length = 0
            parentChildGap := parentStyleConfig.childGap

            for childOffset := int32(0); childOffset < int32(parent.childrenOrTextContent.children.length); childOffset++ {
                childElementIndex := *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(parent.childrenOrTextContent.children.elements)) + uintptr(childOffset)*unsafe.Sizeof(int32(0))))
                childElement := Clay_LayoutElementArray_Get(&context.layoutElements, childElementIndex)
                var childSizing Clay_SizingAxis
                if xAxis {
                    childSizing = childElement.layoutConfig.sizing.width
                } else {
                    childSizing = childElement.layoutConfig.sizing.height
                }
                var childSize float32
                if xAxis {
                    childSize = childElement.dimensions.width
                } else {
                    childSize = childElement.dimensions.height
                }

                if !Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && childElement.childrenOrTextContent.children.length > 0 {
                    Clay__int32_tArray_Add(&bfsBuffer, childElementIndex)
                }

                if childSizing.type_ != CLAY__SIZING_TYPE_PERCENT &&
                    childSizing.type_ != CLAY__SIZING_TYPE_FIXED &&
                    (!Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (Clay__FindElementConfigWithType(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig.wrapMode == CLAY_TEXT_WRAP_WORDS)) {
                    Clay__int32_tArray_Add(&resizableContainerBuffer, childElementIndex)
                }

                if sizingAlongAxis {
                    if childSizing.type_ == CLAY__SIZING_TYPE_PERCENT {
                        innerContentSize += 0
                    } else {
                        innerContentSize += childSize
                    }
                    if childSizing.type_ == CLAY__SIZING_TYPE_GROW {
                        growContainerCount++
                    }
                    if childOffset > 0 {
                        innerContentSize += float32(parentChildGap) // For children after index 0, the childAxisOffset is the gap from the previous child
                        totalPaddingAndChildGaps += float32(parentChildGap)
                    }
                } else {
                    innerContentSize = CLAY__MAX(childSize, innerContentSize)
                }
            }

            // Expand percentage containers to size
            for childOffset := int32(0); childOffset < int32(parent.childrenOrTextContent.children.length); childOffset++ {
                childElementIndex := *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(parent.childrenOrTextContent.children.elements)) + uintptr(childOffset)*unsafe.Sizeof(int32(0))))
                childElement := Clay_LayoutElementArray_Get(&context.layoutElements, childElementIndex)
                var childSizing Clay_SizingAxis
                if xAxis {
                    childSizing = childElement.layoutConfig.sizing.width
                } else {
                    childSizing = childElement.layoutConfig.sizing.height
                }
                childSize := &childElement.dimensions.width
                if !xAxis {
                    childSize = &childElement.dimensions.height
                }
                if childSizing.type_ == CLAY__SIZING_TYPE_PERCENT {
                    *childSize = (parentSize - totalPaddingAndChildGaps) * childSizing.size.percent
                    if sizingAlongAxis {
                        innerContentSize += *childSize
                    }
                    Clay__UpdateAspectRatioBox(childElement)
                }
            }

            if sizingAlongAxis {
                sizeToDistribute := parentSize - parentPadding - innerContentSize
                // The content is too large, compress the children as much as possible
                if sizeToDistribute < 0 {
                    // If the parent clips content in this axis direction, don't compress children, just leave them alone
                    clipElementConfig := Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig
                    if clipElementConfig != nil {
                        if (xAxis && clipElementConfig.horizontal) || (!xAxis && clipElementConfig.vertical) {
                            continue
                        }
                    }
                    // Scrolling containers preferentially compress before others
                    for sizeToDistribute < -CLAY__EPSILON && resizableContainerBuffer.length > 0 {
                        largest := float32(0)
                        secondLargest := float32(0)
                        widthToAdd := sizeToDistribute
                        for childIndex := int32(0); childIndex < resizableContainerBuffer.length; childIndex++ {
                            child := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex))
                            var childSize float32
                            if xAxis {
                                childSize = child.dimensions.width
                            } else {
                                childSize = child.dimensions.height
                            }
                            if Clay__FloatEqual(childSize, largest) {
                                continue
                            }
                            if childSize > largest {
                                secondLargest = largest
                                largest = childSize
                            }
                            if childSize < largest {
                                secondLargest = CLAY__MAX(secondLargest, childSize)
                                widthToAdd = secondLargest - largest
                            }
                        }

                        widthToAdd = CLAY__MAX(widthToAdd, sizeToDistribute/float32(resizableContainerBuffer.length))

                        for childIndex := int32(0); childIndex < resizableContainerBuffer.length; childIndex++ {
                            child := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex))
                            childSize := &child.dimensions.width
                            if !xAxis {
                                childSize = &child.dimensions.height
                            }
                            var minSize float32
                            if xAxis {
                                minSize = child.minDimensions.width
                            } else {
                                minSize = child.minDimensions.height
                            }
                            previousWidth := *childSize
                            if Clay__FloatEqual(*childSize, largest) {
                                *childSize += widthToAdd
                                if *childSize <= minSize {
                                    *childSize = minSize
                                    Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex)
                                    childIndex--
                                }
                                sizeToDistribute -= (*childSize - previousWidth)
                            }
                        }
                    }
                // The content is too small, allow SIZING_GROW containers to expand
                } else if sizeToDistribute > 0 && growContainerCount > 0 {
                    for childIndex := int32(0); childIndex < resizableContainerBuffer.length; childIndex++ {
                        child := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex))
                        var childSizing Clay__SizingType
                        if xAxis {
                            childSizing = child.layoutConfig.sizing.width.type_
                        } else {
                            childSizing = child.layoutConfig.sizing.height.type_
                        }
                        if childSizing != CLAY__SIZING_TYPE_GROW {
                            Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex)
                            childIndex--
                        }
                    }
                    for sizeToDistribute > CLAY__EPSILON && resizableContainerBuffer.length > 0 {
                        smallest := float32(CLAY__MAXFLOAT)
                        secondSmallest := float32(CLAY__MAXFLOAT)
                        widthToAdd := sizeToDistribute
                        for childIndex := int32(0); childIndex < resizableContainerBuffer.length; childIndex++ {
                            child := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex))
                            var childSize float32
                            if xAxis {
                                childSize = child.dimensions.width
                            } else {
                                childSize = child.dimensions.height
                            }
                            if Clay__FloatEqual(childSize, smallest) {
                                continue
                            }
                            if childSize < smallest {
                                secondSmallest = smallest
                                smallest = childSize
                            }
                            if childSize > smallest {
                                secondSmallest = CLAY__MIN(secondSmallest, childSize)
                                widthToAdd = secondSmallest - smallest
                            }
                        }

                        widthToAdd = CLAY__MIN(widthToAdd, sizeToDistribute/float32(resizableContainerBuffer.length))

                        for childIndex := int32(0); childIndex < resizableContainerBuffer.length; childIndex++ {
                            child := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex))
                            childSize := &child.dimensions.width
                            if !xAxis {
                                childSize = &child.dimensions.height
                            }
                            var maxSize float32
                            if xAxis {
                                maxSize = child.layoutConfig.sizing.width.size.minMax.max
                            } else {
                                maxSize = child.layoutConfig.sizing.height.size.minMax.max
                            }
                            previousWidth := *childSize
                            if Clay__FloatEqual(*childSize, smallest) {
                                *childSize += widthToAdd
                                if *childSize >= maxSize {
                                    *childSize = maxSize
                                    Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex)
                                    childIndex--
                                }
                                sizeToDistribute -= (*childSize - previousWidth)
                            }
                        }
                    }
                }
            // Sizing along the non layout axis ("off axis")
            } else {
                for childOffset := int32(0); childOffset < resizableContainerBuffer.length; childOffset++ {
                    childElement := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childOffset))
                    var childSizing Clay_SizingAxis
                    if xAxis {
                        childSizing = childElement.layoutConfig.sizing.width
                    } else {
                        childSizing = childElement.layoutConfig.sizing.height
                    }
                    var minSize float32
                    if xAxis {
                        minSize = childElement.minDimensions.width
                    } else {
                        minSize = childElement.minDimensions.height
                    }
                    childSize := &childElement.dimensions.width
                    if !xAxis {
                        childSize = &childElement.dimensions.height
                    }

                    maxSize := parentSize - parentPadding
                    // If we're laying out the children of a scroll panel, grow containers expand to the size of the inner content, not the outer container
                    if Clay__ElementHasConfig(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP) {
                        clipElementConfig := Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig
                        if (xAxis && clipElementConfig.horizontal) || (!xAxis && clipElementConfig.vertical) {
                            maxSize = CLAY__MAX(maxSize, innerContentSize)
                        }
                    }
                    if childSizing.type_ == CLAY__SIZING_TYPE_GROW {
                        *childSize = CLAY__MIN(maxSize, childSizing.size.minMax.max)
                    }
                    *childSize = CLAY__MAX(minSize, CLAY__MIN(*childSize, maxSize))
                }
            }
        }
    }
}

func Clay__IntToString(integer int32) Clay_String {
    if integer == 0 {
        return CLAY_STRING("0")
    }
    context := Clay_GetCurrentContext()
    chars := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(context.dynamicStringData.internalArray)) + uintptr(context.dynamicStringData.length)))
    length := int32(0)
    sign := integer

    if integer < 0 {
        integer = -integer
    }
    for integer > 0 {
        *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(chars)) + uintptr(length))) = byte(integer%10 + '0')
        length++
        integer /= 10
    }

    if sign < 0 {
        *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(chars)) + uintptr(length))) = '-'
        length++
    }

    // Reverse the string to get the correct order
    for j, k := int32(0), length-1; j < k; j, k = j+1, k-1 {
        temp := *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(chars)) + uintptr(j)))
        *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(chars)) + uintptr(j))) = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(chars)) + uintptr(k)))
        *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(chars)) + uintptr(k))) = temp
    }
    context.dynamicStringData.length += length
    return Clay_String{length: length, chars: unsafe.Pointer(chars)}
}

func Clay__AddRenderCommand(renderCommand Clay_RenderCommand) {
    context := Clay_GetCurrentContext()
    if context.renderCommands.length < context.renderCommands.capacity-1 {
        Clay_RenderCommandArray_Add(&context.renderCommands, renderCommand)
    } else {
        if !context.booleanWarnings.maxRenderCommandsExceeded {
            context.booleanWarnings.maxRenderCommandsExceeded = true
            context.errorHandler.errorHandlerFunction(Clay_ErrorData{
                errorType: CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
                errorText: CLAY_STRING("Clay ran out of capacity while attempting to create render commands. This is usually caused by a large amount of wrapping text elements while close to the max element capacity. Try using Clay_SetMaxElementCount() with a higher value."),
                userData: context.errorHandler.userData,
            })
        }
    }
}

func Clay__ElementIsOffscreen(boundingBox *Clay_BoundingBox) bool {
    context := Clay_GetCurrentContext()
    if context.disableCulling {
        return false
    }

    return (boundingBox.x > float32(context.layoutDimensions.width)) ||
           (boundingBox.y > float32(context.layoutDimensions.height)) ||
           (boundingBox.x+boundingBox.width < 0) ||
           (boundingBox.y+boundingBox.height < 0)
}

func Clay__CalculateFinalLayout() {
    context := Clay_GetCurrentContext()
    // Calculate sizing along the X axis
    Clay__SizeContainersAlongAxis(true)

    // Wrap text
    for textElementIndex := int32(0); textElementIndex < context.textElementData.length; textElementIndex++ {
        textElementData := Clay__TextElementDataArray_Get(&context.textElementData, textElementIndex)
        textElementData.wrappedLines = Clay__WrappedTextLineArraySlice{length: 0, internalArray: (*Clay__WrappedTextLine)(unsafe.Pointer(uintptr(unsafe.Pointer(context.wrappedTextLines.internalArray)) + uintptr(context.wrappedTextLines.length)*unsafe.Sizeof(Clay__WrappedTextLine{})))}
        containerElement := Clay_LayoutElementArray_Get(&context.layoutElements, textElementData.elementIndex)
        textConfig := Clay__FindElementConfigWithType(containerElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig
        measureTextCacheItem := Clay__MeasureTextCached(&textElementData.text, textConfig)
        lineWidth := float32(0)
        var lineHeight float32
        if textConfig.lineHeight > 0 {
            lineHeight = float32(textConfig.lineHeight)
        } else {
            lineHeight = textElementData.preferredDimensions.height
        }
        lineLengthChars := int32(0)
        lineStartOffset := int32(0)
        if !measureTextCacheItem.containsNewlines && textElementData.preferredDimensions.width <= containerElement.dimensions.width {
            Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, Clay__WrappedTextLine{dimensions: containerElement.dimensions, line: textElementData.text})
            textElementData.wrappedLines.length++
            continue
        }
        spaceWidth := Clay__MeasureText(Clay_StringSlice{length: 1, chars: CLAY__SPACECHAR.chars, baseChars: CLAY__SPACECHAR.chars}, textConfig, context.measureTextUserData.(unsafe.Pointer)).width
        wordIndex := measureTextCacheItem.measuredWordsStartIndex
        for wordIndex != -1 {
            if context.wrappedTextLines.length > context.wrappedTextLines.capacity-1 {
                break
            }
            measuredWord := Clay__MeasuredWordArray_Get(&context.measuredWords, wordIndex)
            // Only word on the line is too large, just render it anyway
            if lineLengthChars == 0 && lineWidth+measuredWord.width > containerElement.dimensions.width {
                Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, Clay__WrappedTextLine{dimensions: Clay_Dimensions{width: measuredWord.width, height: lineHeight}, line: Clay_String{length: measuredWord.length, chars: unsafe.Pointer(uintptr(textElementData.text.chars) + uintptr(measuredWord.startOffset))}})
                textElementData.wrappedLines.length++
                wordIndex = measuredWord.next
                lineStartOffset = measuredWord.startOffset + measuredWord.length
            } else if measuredWord.length == 0 || lineWidth+measuredWord.width > containerElement.dimensions.width {
                // Wrapped text lines list has overflowed, just render out the line
                finalCharIsSpace := *(*byte)(unsafe.Pointer(uintptr(textElementData.text.chars) + uintptr(CLAY__MAX(float32(lineStartOffset+lineLengthChars-1), 0)))) == ' '
                var widthAdjustment float32
                var lengthAdjustment int32
                if finalCharIsSpace {
                    widthAdjustment = -spaceWidth
                    lengthAdjustment = -1
                } else {
                    widthAdjustment = 0
                    lengthAdjustment = 0
                }
                Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, Clay__WrappedTextLine{dimensions: Clay_Dimensions{width: lineWidth + widthAdjustment, height: lineHeight}, line: Clay_String{length: lineLengthChars + lengthAdjustment, chars: unsafe.Pointer(uintptr(textElementData.text.chars) + uintptr(lineStartOffset))}})
                textElementData.wrappedLines.length++
                if lineLengthChars == 0 || measuredWord.length == 0 {
                    wordIndex = measuredWord.next
                }
                lineWidth = 0
                lineLengthChars = 0
                lineStartOffset = measuredWord.startOffset
            } else {
                lineWidth += measuredWord.width + float32(textConfig.letterSpacing)
                lineLengthChars += measuredWord.length
                wordIndex = measuredWord.next
            }
        }
        if lineLengthChars > 0 {
            Clay__WrappedTextLineArray_Add(&context.wrappedTextLines, Clay__WrappedTextLine{dimensions: Clay_Dimensions{width: lineWidth - float32(textConfig.letterSpacing), height: lineHeight}, line: Clay_String{length: lineLengthChars, chars: unsafe.Pointer(uintptr(textElementData.text.chars) + uintptr(lineStartOffset))}})
            textElementData.wrappedLines.length++
        }
        containerElement.dimensions.height = lineHeight * float32(textElementData.wrappedLines.length)
    }

    // Scale vertical heights according to aspect ratio
    for i := int32(0); i < context.aspectRatioElementIndexes.length; i++ {
        aspectElement := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.aspectRatioElementIndexes, i))
        config := Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig
        aspectElement.dimensions.height = (1 / config.aspectRatio) * aspectElement.dimensions.width
        aspectElement.layoutConfig.sizing.height.size.minMax.max = aspectElement.dimensions.height
    }

    // Propagate effect of text wrapping, aspect scaling etc. on height of parents
    dfsBuffer := context.layoutElementTreeNodeArray1
    dfsBuffer.length = 0
    for i := int32(0); i < context.layoutElementTreeRoots.length; i++ {
        root := Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, i)
        *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length)*unsafe.Sizeof(bool(false)))) = false
        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, Clay__LayoutElementTreeNode{layoutElement: Clay_LayoutElementArray_Get(&context.layoutElements, root.layoutElementIndex)})
    }
    for dfsBuffer.length > 0 {
        currentElementTreeNode := Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, dfsBuffer.length-1)
        currentElement := currentElementTreeNode.layoutElement
        if !*(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length-1)*unsafe.Sizeof(bool(false)))) {
            *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length-1)*unsafe.Sizeof(bool(false)))) = true
            // If the element has no children or is the container for a text element, don't bother inspecting it
            if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement.childrenOrTextContent.children.length == 0 {
                dfsBuffer.length--
                continue
            }
            // Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
            for i := int32(0); i < int32(currentElement.childrenOrTextContent.children.length); i++ {
                *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length)*unsafe.Sizeof(bool(false)))) = false
                Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, Clay__LayoutElementTreeNode{layoutElement: Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))})
            }
            continue
        }
        dfsBuffer.length--

        // DFS node has been visited, this is on the way back up to the root
        layoutConfig := currentElement.layoutConfig
        if layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT {
            // Resize any parent containers that have grown in height along their non layout axis
            for j := int32(0); j < int32(currentElement.childrenOrTextContent.children.length); j++ {
                childElement := Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(j)*unsafe.Sizeof(int32(0)))))
                childHeightWithPadding := CLAY__MAX(childElement.dimensions.height+float32(layoutConfig.padding.top+layoutConfig.padding.bottom), currentElement.dimensions.height)
                currentElement.dimensions.height = CLAY__MIN(CLAY__MAX(childHeightWithPadding, layoutConfig.sizing.height.size.minMax.min), layoutConfig.sizing.height.size.minMax.max)
            }
        } else if layoutConfig.layoutDirection == CLAY_TOP_TO_BOTTOM {
            // Resizing along the layout axis
            contentHeight := float32(layoutConfig.padding.top + layoutConfig.padding.bottom)
            for j := int32(0); j < int32(currentElement.childrenOrTextContent.children.length); j++ {
                childElement := Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(j)*unsafe.Sizeof(int32(0)))))
                contentHeight += childElement.dimensions.height
            }
            contentHeight += CLAY__MAX(float32(int32(currentElement.childrenOrTextContent.children.length)-1), 0) * float32(layoutConfig.childGap)
            currentElement.dimensions.height = CLAY__MIN(CLAY__MAX(contentHeight, layoutConfig.sizing.height.size.minMax.min), layoutConfig.sizing.height.size.minMax.max)
        }
    }

    // Calculate sizing along the Y axis
    Clay__SizeContainersAlongAxis(false)

    // Scale horizontal widths according to aspect ratio
    for i := int32(0); i < context.aspectRatioElementIndexes.length; i++ {
        aspectElement := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.aspectRatioElementIndexes, i))
        config := Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig
        aspectElement.dimensions.width = config.aspectRatio * aspectElement.dimensions.height
    }

    // Sort tree roots by z-index
    sortMax := context.layoutElementTreeRoots.length - 1
    for sortMax > 0 { // todo dumb bubble sort
        for i := int32(0); i < sortMax; i++ {
            current := *Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, i)
            next := *Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, i+1)
            if next.zIndex < current.zIndex {
                Clay__LayoutElementTreeRootArray_Set(&context.layoutElementTreeRoots, i, &next)
                Clay__LayoutElementTreeRootArray_Set(&context.layoutElementTreeRoots, i+1, &current)
            }
        }
        sortMax--
    }

    // Calculate final positions and generate render commands
    context.renderCommands.length = 0
    dfsBuffer.length = 0
    for rootIndex := int32(0); rootIndex < context.layoutElementTreeRoots.length; rootIndex++ {
        dfsBuffer.length = 0
        root := Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex)
        rootElement := Clay_LayoutElementArray_Get(&context.layoutElements, root.layoutElementIndex)
        rootPosition := Clay_Vector2{}
        parentHashMapItem := Clay__GetHashMapItem(root.parentId)
        // Position root floating containers
        if Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) && parentHashMapItem != nil {
            config := Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig
            rootDimensions := rootElement.dimensions
            parentBoundingBox := parentHashMapItem.boundingBox
            // Set X position
            targetAttachPosition := Clay_Vector2{}
            switch config.attachPoints.parent {
            case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_LEFT_BOTTOM:
                targetAttachPosition.x = parentBoundingBox.x
            case CLAY_ATTACH_POINT_CENTER_TOP, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_CENTER_BOTTOM:
                targetAttachPosition.x = parentBoundingBox.x + (parentBoundingBox.width / 2)
            case CLAY_ATTACH_POINT_RIGHT_TOP, CLAY_ATTACH_POINT_RIGHT_CENTER, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
                targetAttachPosition.x = parentBoundingBox.x + parentBoundingBox.width
            }
            switch config.attachPoints.element {
            case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_LEFT_BOTTOM:
                // No adjustment needed
            case CLAY_ATTACH_POINT_CENTER_TOP, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_CENTER_BOTTOM:
                targetAttachPosition.x -= (rootDimensions.width / 2)
            case CLAY_ATTACH_POINT_RIGHT_TOP, CLAY_ATTACH_POINT_RIGHT_CENTER, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
                targetAttachPosition.x -= rootDimensions.width
            }
            switch config.attachPoints.parent { // I know I could merge the x and y switch statements, but this is easier to read
            case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_RIGHT_TOP, CLAY_ATTACH_POINT_CENTER_TOP:
                targetAttachPosition.y = parentBoundingBox.y
            case CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_RIGHT_CENTER:
                targetAttachPosition.y = parentBoundingBox.y + (parentBoundingBox.height / 2)
            case CLAY_ATTACH_POINT_LEFT_BOTTOM, CLAY_ATTACH_POINT_CENTER_BOTTOM, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
                targetAttachPosition.y = parentBoundingBox.y + parentBoundingBox.height
            }
            switch config.attachPoints.element {
            case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_RIGHT_TOP, CLAY_ATTACH_POINT_CENTER_TOP:
                // No adjustment needed
            case CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_RIGHT_CENTER:
                targetAttachPosition.y -= (rootDimensions.height / 2)
            case CLAY_ATTACH_POINT_LEFT_BOTTOM, CLAY_ATTACH_POINT_CENTER_BOTTOM, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
                targetAttachPosition.y -= rootDimensions.height
            }
            targetAttachPosition.x += config.offset.x
            targetAttachPosition.y += config.offset.y
            rootPosition = targetAttachPosition
        }
        if root.clipElementId != 0 {
            clipHashMapItem := Clay__GetHashMapItem(root.clipElementId)
            if clipHashMapItem != nil {
                // Floating elements that are attached to scrolling contents won't be correctly positioned if external scroll handling is enabled, fix here
                if context.externalScrollHandlingEnabled {
                    clipConfig := Clay__FindElementConfigWithType(clipHashMapItem.layoutElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig
                    if clipConfig.horizontal {
                        rootPosition.x += clipConfig.childOffset.x
                    }
                    if clipConfig.vertical {
                        rootPosition.y += clipConfig.childOffset.y
                    }
                }
                Clay__AddRenderCommand(Clay_RenderCommand{
                    boundingBox: clipHashMapItem.boundingBox,
                    userData: 0,
                    id: Clay__HashNumber(rootElement.id, uint32(rootElement.childrenOrTextContent.children.length+10)).id, // TODO need a better strategy for managing derived ids
                    zIndex: root.zIndex,
                    commandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_START,
                })
            }
        }
        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, Clay__LayoutElementTreeNode{layoutElement: rootElement, position: rootPosition, nextChildOffset: Clay_Vector2{x: float32(rootElement.layoutConfig.padding.left), y: float32(rootElement.layoutConfig.padding.top)}})

        *(*bool)(unsafe.Pointer(context.treeNodeVisited.internalArray)) = false
        for dfsBuffer.length > 0 {
            currentElementTreeNode := Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, dfsBuffer.length-1)
            currentElement := currentElementTreeNode.layoutElement
            layoutConfig := currentElement.layoutConfig
            scrollOffset := Clay_Vector2{}

            // This will only be run a single time for each element in downwards DFS order
            if !*(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length-1)*unsafe.Sizeof(bool(false)))) {
                *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length-1)*unsafe.Sizeof(bool(false)))) = true

                currentElementBoundingBox := Clay_BoundingBox{x: currentElementTreeNode.position.x, y: currentElementTreeNode.position.y, width: currentElement.dimensions.width, height: currentElement.dimensions.height}
                if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) {
                    floatingElementConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig
                    expand := floatingElementConfig.expand
                    currentElementBoundingBox.x -= expand.width
                    currentElementBoundingBox.width += expand.width * 2
                    currentElementBoundingBox.y -= expand.height
                    currentElementBoundingBox.height += expand.height * 2
                }

                var scrollContainerData *Clay__ScrollContainerDataInternal
                // Apply scroll offsets to container
                if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP) {
                    clipConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig

                    // This linear scan could theoretically be slow under very strange conditions, but I can't imagine a real UI with more than a few 10's of scroll containers
                    for i := int32(0); i < context.scrollContainerDatas.length; i++ {
                        mapping := Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i)
                        if mapping.layoutElement == currentElement {
                            scrollContainerData = mapping
                            mapping.boundingBox = currentElementBoundingBox
                            scrollOffset = clipConfig.childOffset
                            if context.externalScrollHandlingEnabled {
                                scrollOffset = Clay_Vector2{}
                            }
                            break
                        }
                    }
                }

                hashMapItem := Clay__GetHashMapItem(currentElement.id)
                if hashMapItem != nil {
                    hashMapItem.boundingBox = currentElementBoundingBox
                    if hashMapItem.idAlias != 0 {
                        hashMapItemAlias := Clay__GetHashMapItem(hashMapItem.idAlias)
                        if hashMapItemAlias != nil {
                            hashMapItemAlias.boundingBox = currentElementBoundingBox
                        }
                    }
                }

                var sortedConfigIndexes [20]int32
                for elementConfigIndex := int32(0); elementConfigIndex < currentElement.elementConfigs.length; elementConfigIndex++ {
                    sortedConfigIndexes[elementConfigIndex] = elementConfigIndex
                }
                sortMax := currentElement.elementConfigs.length - 1
                for sortMax > 0 { // todo dumb bubble sort
                    for i := int32(0); i < sortMax; i++ {
                        current := sortedConfigIndexes[i]
                        next := sortedConfigIndexes[i+1]
                        currentType := Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, current)._type
                        nextType := Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, next)._type
                        if nextType == CLAY__ELEMENT_CONFIG_TYPE_CLIP || currentType == CLAY__ELEMENT_CONFIG_TYPE_BORDER {
                            sortedConfigIndexes[i] = next
                            sortedConfigIndexes[i+1] = current
                        }
                    }
                    sortMax--
                }

                emitRectangle := false
                // Create the render commands for this element
                sharedConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig
                if sharedConfig != nil && sharedConfig.backgroundColor.a > 0 {
                   emitRectangle = true
                } else if sharedConfig == nil {
                    emitRectangle = false
                    sharedConfig = &Clay_SharedElementConfig_DEFAULT
                }
                for elementConfigIndex := int32(0); elementConfigIndex < currentElement.elementConfigs.length; elementConfigIndex++ {
                    elementConfig := Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, sortedConfigIndexes[elementConfigIndex])
                    renderCommand := Clay_RenderCommand{
                        boundingBox: currentElementBoundingBox,
                        userData: sharedConfig.userData,
                        id: currentElement.id,
                    }

                    offscreen := Clay__ElementIsOffscreen(&currentElementBoundingBox)
                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
                    shouldRender := !offscreen
                    switch elementConfig._type {
                    case CLAY__ELEMENT_CONFIG_TYPE_ASPECT, CLAY__ELEMENT_CONFIG_TYPE_FLOATING, CLAY__ELEMENT_CONFIG_TYPE_SHARED, CLAY__ELEMENT_CONFIG_TYPE_BORDER:
                        shouldRender = false
                    case CLAY__ELEMENT_CONFIG_TYPE_CLIP:
                        renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START
                        renderCommand.renderData = Clay_RenderData{
                            clip: Clay_ClipRenderData{
                                horizontal: elementConfig.config.clipElementConfig.horizontal,
                                vertical: elementConfig.config.clipElementConfig.vertical,
                            },
                        }
                    case CLAY__ELEMENT_CONFIG_TYPE_IMAGE:
                        renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_IMAGE
                        renderCommand.renderData = Clay_RenderData{
                            image: Clay_ImageRenderData{
                                backgroundColor: sharedConfig.backgroundColor,
                                cornerRadius: sharedConfig.cornerRadius,
                                imageData: elementConfig.config.imageElementConfig.imageData,
                            },
                        }
                        emitRectangle = false
                    case CLAY__ELEMENT_CONFIG_TYPE_TEXT:
                        if !shouldRender {
                            break
                        }
                        shouldRender = false
                        configUnion := elementConfig.config
                        textElementConfig := configUnion.textElementConfig
                        naturalLineHeight := currentElement.childrenOrTextContent.textElementData.preferredDimensions.height
                        var finalLineHeight float32
                        if textElementConfig.lineHeight > 0 {
                            finalLineHeight = float32(textElementConfig.lineHeight)
                        } else {
                            finalLineHeight = naturalLineHeight
                        }
                        lineHeightOffset := (finalLineHeight - naturalLineHeight) / 2
                        yPosition := lineHeightOffset
                        for lineIndex := int32(0); lineIndex < currentElement.childrenOrTextContent.textElementData.wrappedLines.length; lineIndex++ {
                            wrappedLine := Clay__WrappedTextLineArraySlice_Get(&currentElement.childrenOrTextContent.textElementData.wrappedLines, lineIndex)
                            if wrappedLine.line.length == 0 {
                                yPosition += finalLineHeight
                                continue
                            }
                            offset := currentElementBoundingBox.width - wrappedLine.dimensions.width
                            if textElementConfig.textAlignment == CLAY_TEXT_ALIGN_LEFT {
                                offset = 0
                            }
                            if textElementConfig.textAlignment == CLAY_TEXT_ALIGN_CENTER {
                                offset /= 2
                            }
                            Clay__AddRenderCommand(Clay_RenderCommand{
                                boundingBox: Clay_BoundingBox{x: currentElementBoundingBox.x + offset, y: currentElementBoundingBox.y + yPosition, width: wrappedLine.dimensions.width, height: wrappedLine.dimensions.height},
                                renderData: Clay_RenderData{
                                    text: Clay_TextRenderData{
                                        stringContents: Clay_StringSlice{length: wrappedLine.line.length, chars: wrappedLine.line.chars, baseChars: currentElement.childrenOrTextContent.textElementData.text.chars},
                                        textColor: textElementConfig.textColor,
                                        fontId: textElementConfig.fontId,
                                        fontSize: textElementConfig.fontSize,
                                        letterSpacing: textElementConfig.letterSpacing,
                                        lineHeight: textElementConfig.lineHeight,
                                    },
                                },
                                userData: textElementConfig.userData,
                                id: Clay__HashNumber(uint32(lineIndex), currentElement.id).id,
                                zIndex: root.zIndex,
                                commandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
                            })
                            yPosition += finalLineHeight

                            if !context.disableCulling && (currentElementBoundingBox.y+yPosition > float32(context.layoutDimensions.height)) {
                                break
                            }
                        }
                    case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM:
                        renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_CUSTOM
                        renderCommand.renderData = Clay_RenderData{
                            custom: Clay_CustomRenderData{
                                backgroundColor: sharedConfig.backgroundColor,
                                cornerRadius: sharedConfig.cornerRadius,
                                customData: elementConfig.config.customElementConfig.customData,
                            },
                        }
                        emitRectangle = false
                    }
                    if shouldRender {
                        Clay__AddRenderCommand(renderCommand)
                    }
                    if offscreen {
                        // NOTE: You may be tempted to try an early return / continue if an element is off screen. Why bother calculating layout for its children, right?
                        // Unfortunately, a FLOATING_CONTAINER may be defined that attaches to a child or grandchild of this element, which is large enough to still
                        // be on screen, even if this element isn't. That depends on this element and it's children being laid out correctly (even if they are entirely off screen)
                    }
                }

                if emitRectangle {
                    Clay__AddRenderCommand(Clay_RenderCommand{
                        boundingBox: currentElementBoundingBox,
                        renderData: Clay_RenderData{
                            rectangle: Clay_RectangleRenderData{
                                backgroundColor: sharedConfig.backgroundColor,
                                cornerRadius: sharedConfig.cornerRadius,
                            },
                        },
                        userData: sharedConfig.userData,
                        id: currentElement.id,
                        zIndex: root.zIndex,
                        commandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                    })
                }

                // Setup initial on-axis alignment
                if !Clay__ElementHasConfig(currentElementTreeNode.layoutElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) {
                    contentSize := Clay_Dimensions{width: 0, height: 0}
                    if layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT {
                        for i := int32(0); i < int32(currentElement.childrenOrTextContent.children.length); i++ {
                            childElement := Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))
                            contentSize.width += childElement.dimensions.width
                            contentSize.height = CLAY__MAX(contentSize.height, childElement.dimensions.height)
                        }
                        contentSize.width += CLAY__MAX(float32(int32(currentElement.childrenOrTextContent.children.length)-1), 0) * float32(layoutConfig.childGap)
                        extraSpace := currentElement.dimensions.width - float32(layoutConfig.padding.left + layoutConfig.padding.right) - contentSize.width
                        switch layoutConfig.childAlignment.x {
                            case CLAY_ALIGN_X_LEFT: extraSpace = 0
                            case CLAY_ALIGN_X_CENTER: extraSpace /= 2
                            default:
                        }
                        currentElementTreeNode.nextChildOffset.x += extraSpace
                        extraSpace = CLAY__MAX(0, extraSpace)
                    } else {
                        for i := int32(0); i < int32(currentElement.childrenOrTextContent.children.length); i++ {
                            childElement := Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))
                            contentSize.width = CLAY__MAX(contentSize.width, childElement.dimensions.width)
                            contentSize.height += childElement.dimensions.height
                        }
                        contentSize.height += CLAY__MAX(float32(int32(currentElement.childrenOrTextContent.children.length)-1), 0) * float32(layoutConfig.childGap)
                        extraSpace := currentElement.dimensions.height - float32(layoutConfig.padding.top + layoutConfig.padding.bottom) - contentSize.height
                        switch layoutConfig.childAlignment.y {
                            case CLAY_ALIGN_Y_TOP: extraSpace = 0
                            case CLAY_ALIGN_Y_CENTER: extraSpace /= 2
                            default:
                        }
                        extraSpace = CLAY__MAX(0, extraSpace)
                        currentElementTreeNode.nextChildOffset.y += extraSpace
                    }

                    if scrollContainerData != nil {
                        scrollContainerData.contentSize = Clay_Dimensions{ width: contentSize.width + float32(layoutConfig.padding.left + layoutConfig.padding.right), height: contentSize.height + float32(layoutConfig.padding.top + layoutConfig.padding.bottom) }
                    }
                }
            } else {
                // DFS is returning upwards backwards
                closeClipElement := false
                clipConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig
                if clipConfig != nil {
                    closeClipElement = true
                    for i := int32(0); i < context.scrollContainerDatas.length; i++ {
                        mapping := Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i)
                        if mapping.layoutElement == currentElement {
                            scrollOffset = clipConfig.childOffset
                            if context.externalScrollHandlingEnabled {
                                scrollOffset = Clay_Vector2{x: 0, y: 0}
                            }
                            break
                        }
                    }
                }

                if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER) {
                    currentElementData := Clay__GetHashMapItem(currentElement.id)
                    currentElementBoundingBox := currentElementData.boundingBox

                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
                    if !Clay__ElementIsOffscreen(&currentElementBoundingBox) {
                        var sharedConfig *Clay_SharedElementConfig
                        if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED) {
                            sharedConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig
                        } else {
                            sharedConfig = &Clay_SharedElementConfig_DEFAULT
                        }
                        borderConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER).borderElementConfig
                        renderCommand := Clay_RenderCommand{
                            boundingBox: currentElementBoundingBox,
                            renderData: Clay_RenderData{
                                border: Clay_BorderRenderData{
                                    color: borderConfig.color,
                                    cornerRadius: sharedConfig.cornerRadius,
                                    width: borderConfig.width,
                                },
                            },
                            userData: sharedConfig.userData,
                            id: Clay__HashNumber(currentElement.id, uint32(currentElement.childrenOrTextContent.children.length)).id,
                            commandType: CLAY_RENDER_COMMAND_TYPE_BORDER,
                        }
                        Clay__AddRenderCommand(renderCommand)
                        if borderConfig.width.betweenChildren > 0 && borderConfig.color.a > 0 {
                            halfGap := float32(layoutConfig.childGap) / 2
                            borderOffset := Clay_Vector2{ x: float32(layoutConfig.padding.left) - halfGap, y: float32(layoutConfig.padding.top) - halfGap }
                            if layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT {
                                for i := int32(0); i < int32(currentElement.childrenOrTextContent.children.length); i++ {
                                    childElement := Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))
                                    if i > 0 {
                                        Clay__AddRenderCommand(Clay_RenderCommand{
                                            boundingBox: Clay_BoundingBox{ x: currentElementBoundingBox.x + borderOffset.x + scrollOffset.x, y: currentElementBoundingBox.y + scrollOffset.y, width: float32(borderConfig.width.betweenChildren), height: currentElement.dimensions.height },
                                            renderData: Clay_RenderData{
                                                rectangle: Clay_RectangleRenderData{
                                                    backgroundColor: borderConfig.color,
                                                },
                                            },
                                            userData: sharedConfig.userData,
                                            id: Clay__HashNumber(currentElement.id, uint32(currentElement.childrenOrTextContent.children.length) + 1 + uint32(i)).id,
                                            commandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                                        })
                                    }
                                    borderOffset.x += childElement.dimensions.width + float32(layoutConfig.childGap)
                                }
                            } else {
                                for i := int32(0); i < int32(currentElement.childrenOrTextContent.children.length); i++ {
                                    childElement := Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))
                                    if i > 0 {
                                        Clay__AddRenderCommand(Clay_RenderCommand{
                                            boundingBox: Clay_BoundingBox{ x: currentElementBoundingBox.x + scrollOffset.x, y: currentElementBoundingBox.y + borderOffset.y + scrollOffset.y, width: currentElement.dimensions.width, height: float32(borderConfig.width.betweenChildren) },
                                            renderData: Clay_RenderData{
                                                rectangle: Clay_RectangleRenderData{
                                                    backgroundColor: borderConfig.color,
                                                },
                                            },
                                            userData: sharedConfig.userData,
                                            id: Clay__HashNumber(currentElement.id, uint32(currentElement.childrenOrTextContent.children.length) + 1 + uint32(i)).id,
                                            commandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                                        })
                                    }
                                    borderOffset.y += childElement.dimensions.height + float32(layoutConfig.childGap)
                                }
                            }
                        }
                    }
                }
                // This exists because the scissor needs to end _after_ borders between elements
                if closeClipElement {
                    Clay__AddRenderCommand(Clay_RenderCommand{
                        id: Clay__HashNumber(currentElement.id, uint32(rootElement.childrenOrTextContent.children.length + 11)).id,
                        commandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
                    })
                }

                dfsBuffer.length--
                continue
            }

            // Add children to the DFS buffer
            if !Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) {
                dfsBuffer.length += int32(currentElement.childrenOrTextContent.children.length)
                for i := int32(0); i < int32(currentElement.childrenOrTextContent.children.length); i++ {
                    childElement := Clay_LayoutElementArray_Get(&context.layoutElements, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))
                    // Alignment along non layout axis
                    if layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT {
                        currentElementTreeNode.nextChildOffset.y = float32(currentElement.layoutConfig.padding.top)
                        whiteSpaceAroundChild := currentElement.dimensions.height - float32(layoutConfig.padding.top + layoutConfig.padding.bottom) - childElement.dimensions.height
                        switch layoutConfig.childAlignment.y {
                            case CLAY_ALIGN_Y_TOP:
                            case CLAY_ALIGN_Y_CENTER: currentElementTreeNode.nextChildOffset.y += whiteSpaceAroundChild / 2
                            case CLAY_ALIGN_Y_BOTTOM: currentElementTreeNode.nextChildOffset.y += whiteSpaceAroundChild
                        }
                    } else {
                        currentElementTreeNode.nextChildOffset.x = float32(currentElement.layoutConfig.padding.left)
                        whiteSpaceAroundChild := currentElement.dimensions.width - float32(layoutConfig.padding.left + layoutConfig.padding.right) - childElement.dimensions.width
                        switch layoutConfig.childAlignment.x {
                            case CLAY_ALIGN_X_LEFT:
                            case CLAY_ALIGN_X_CENTER: currentElementTreeNode.nextChildOffset.x += whiteSpaceAroundChild / 2
                            case CLAY_ALIGN_X_RIGHT: currentElementTreeNode.nextChildOffset.x += whiteSpaceAroundChild
                        }
                    }

                    childPosition := Clay_Vector2{
                        x: currentElementTreeNode.position.x + currentElementTreeNode.nextChildOffset.x + scrollOffset.x,
                        y: currentElementTreeNode.position.y + currentElementTreeNode.nextChildOffset.y + scrollOffset.y,
                    }

                    // DFS buffer elements need to be added in reverse because stack traversal happens backwards
                    newNodeIndex := uint32(dfsBuffer.length - 1 - i)
                    *(*Clay__LayoutElementTreeNode)(unsafe.Pointer(uintptr(unsafe.Pointer(dfsBuffer.internalArray)) + uintptr(newNodeIndex)*unsafe.Sizeof(Clay__LayoutElementTreeNode{}))) = Clay__LayoutElementTreeNode{
                        layoutElement: childElement,
                        position: Clay_Vector2{ x: childPosition.x, y: childPosition.y },
                        nextChildOffset: Clay_Vector2{ x: float32(childElement.layoutConfig.padding.left), y: float32(childElement.layoutConfig.padding.top) },
                    }
                    *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(newNodeIndex)*unsafe.Sizeof(bool(false)))) = false

                    // Update parent offsets
                    if layoutConfig.layoutDirection == CLAY_LEFT_TO_RIGHT {
                        currentElementTreeNode.nextChildOffset.x += childElement.dimensions.width + float32(layoutConfig.childGap)
                    } else {
                        currentElementTreeNode.nextChildOffset.y += childElement.dimensions.height + float32(layoutConfig.childGap)
                    }
                }
            }
        }

        if root.clipElementId != 0 {
            Clay__AddRenderCommand(Clay_RenderCommand{ id: Clay__HashNumber(rootElement.id, uint32(rootElement.childrenOrTextContent.children.length + 11)).id, commandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_END })
        }
    }
}

func Clay_GetPointerOverIds() Clay_ElementIdArray {
    return Clay_GetCurrentContext().pointerOverIds
}

// DebugTools
var CLAY__DEBUGVIEW_COLOR_1 = Clay_Color{r: 58, g: 56, b: 52, a: 255}
var CLAY__DEBUGVIEW_COLOR_2 = Clay_Color{r: 62, g: 60, b: 58, a: 255}
var CLAY__DEBUGVIEW_COLOR_3 = Clay_Color{r: 141, g: 133, b: 135, a: 255}
var CLAY__DEBUGVIEW_COLOR_4 = Clay_Color{r: 238, g: 226, b: 231, a: 255}
var CLAY__DEBUGVIEW_COLOR_SELECTED_ROW = Clay_Color{r: 102, g: 80, b: 78, a: 255}
const CLAY__DEBUGVIEW_ROW_HEIGHT int32 = 30
const CLAY__DEBUGVIEW_OUTER_PADDING int32 = 10
var Clay__debugViewWidth uint32 = 400
var Clay__debugViewHighlightColor = Clay_Color{ r: 168, g: 66, b: 28, a: 100 }
const CLAY__DEBUGVIEW_INDENT_WIDTH int32 = 16
var Clay__DebugView_TextNameConfig = Clay_TextElementConfig{textColor: Clay_Color{r: 238, g: 226, b: 231, a: 255}, fontSize: 16, wrapMode: CLAY_TEXT_WRAP_NONE}
var Clay__DebugView_ScrollViewItemLayoutConfig Clay_LayoutConfig

type Clay__DebugElementConfigTypeLabelConfig struct {
    label Clay_String
    color Clay_Color
}

func Clay__DebugGetElementConfigTypeLabel(_type Clay__ElementConfigType) Clay__DebugElementConfigTypeLabelConfig {
    switch _type {
        case CLAY__ELEMENT_CONFIG_TYPE_SHARED: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Shared"), color: Clay_Color{r: 243, g: 134, b: 48, a: 255} }
        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Text"), color: Clay_Color{r: 105, g: 210, b: 231, a: 255} }
        case CLAY__ELEMENT_CONFIG_TYPE_ASPECT: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Aspect"), color: Clay_Color{r: 101, g: 149, b: 194, a: 255} }
        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Image"), color: Clay_Color{r: 121, g: 189, b: 154, a: 255} }
        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Floating"), color: Clay_Color{r: 250, g: 105, b: 0, a: 255} }
        case CLAY__ELEMENT_CONFIG_TYPE_CLIP: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Scroll"), color: Clay_Color{r: 242, g: 196, b: 90, a: 255} }
        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Border"), color: Clay_Color{r: 108, g: 91, b: 123, a: 255} }
        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Custom"), color: Clay_Color{r: 11, g: 72, b: 107, a: 255} }
        default:
    }
    return Clay__DebugElementConfigTypeLabelConfig{ label: CLAY_STRING("Error"), color: Clay_Color{r: 0, g: 0, b: 0, a: 255} }
}

type Clay__RenderDebugLayoutData struct {
    rowCount int32
    selectedElementRowIndex int32
}

// Returns row count
func Clay__RenderDebugLayoutElementsList(initialRootsLength int32, highlightedRowIndex int32) Clay__RenderDebugLayoutData {
    context := Clay_GetCurrentContext()
    dfsBuffer := context.reusableElementIndexBuffer
    Clay__DebugView_ScrollViewItemLayoutConfig = Clay_LayoutConfig{ sizing: Clay_Sizing{ height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT)) }, childGap: 6, childAlignment: Clay_ChildAlignment{ y: CLAY_ALIGN_Y_CENTER }}
    layoutData := Clay__RenderDebugLayoutData{}

    highlightedElementId := uint32(0)

    for rootIndex := int32(0); rootIndex < initialRootsLength; rootIndex++ {
        dfsBuffer.length = 0
        root := Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex)
        Clay__int32_tArray_Add(&dfsBuffer, int32(root.layoutElementIndex))
        *(*bool)(unsafe.Pointer(context.treeNodeVisited.internalArray)) = false
        if rootIndex > 0 {
            CLAY(Clay_ElementDeclaration{ id: CLAY_IDI("Clay__DebugView_EmptyRowOuter", uint32(rootIndex)), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0)}, padding: Clay_Padding{left: uint16(CLAY__DEBUGVIEW_INDENT_WIDTH / 2), right: 0, top: 0, bottom: 0} } }, func() {
                CLAY(Clay_ElementDeclaration{ id: CLAY_IDI("Clay__DebugView_EmptyRow", uint32(rootIndex)), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT)) }}, border: Clay_BorderElementConfig{ color: CLAY__DEBUGVIEW_COLOR_3, width: Clay_BorderWidth{ top: 1 } } }, func() {})
            })
            layoutData.rowCount++
        }
        for dfsBuffer.length > 0 {
            currentElementIndex := Clay__int32_tArray_GetValue(&dfsBuffer, int32(dfsBuffer.length - 1))
            currentElement := Clay_LayoutElementArray_Get(&context.layoutElements, int32(currentElementIndex))
            if *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length - 1)*unsafe.Sizeof(bool(false)))) {
                if !Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && currentElement.childrenOrTextContent.children.length > 0 {
                    Clay__CloseElement()
                    Clay__CloseElement()
                    Clay__CloseElement()
                }
                dfsBuffer.length--
                continue
            }

            if highlightedRowIndex == layoutData.rowCount {
                if context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
                    context.debugSelectedElementId = currentElement.id
                }
                highlightedElementId = currentElement.id
            }

            *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length - 1)*unsafe.Sizeof(bool(false)))) = true
            currentElementData := Clay__GetHashMapItem(currentElement.id)
            offscreen := Clay__ElementIsOffscreen(&currentElementData.boundingBox)
            if context.debugSelectedElementId == currentElement.id {
                layoutData.selectedElementRowIndex = layoutData.rowCount
            }
            CLAY(Clay_ElementDeclaration{ id: CLAY_IDI("Clay__DebugView_ElementOuter", currentElement.id), layout: Clay__DebugView_ScrollViewItemLayoutConfig }, func() {
                // Collapse icon / button
                if !(Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement.childrenOrTextContent.children.length == 0) {
                    CLAY(Clay_ElementDeclaration{
                        id: CLAY_IDI("Clay__DebugView_CollapseElement", currentElement.id),
                        layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_FIXED(16), height: CLAY_SIZING_FIXED(16)}, childAlignment: Clay_ChildAlignment{ x: CLAY_ALIGN_X_CENTER, y: CLAY_ALIGN_Y_CENTER} },
                        cornerRadius: CLAY_CORNER_RADIUS(4),
                        border: Clay_BorderElementConfig{ color: CLAY__DEBUGVIEW_COLOR_3, width: Clay_BorderWidth{left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0} },
                    }, func() {
                        var text Clay_String
                        if currentElementData != nil && currentElementData.debugData.collapsed {
                            text = CLAY_STRING("+")
                        } else {
                            text = CLAY_STRING("-")
                        }
                        CLAY_TEXT(text, CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_4, fontSize: 16 }))
                    })
                } else { // Square dot for empty containers
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_FIXED(16), height: CLAY_SIZING_FIXED(16)}, childAlignment: Clay_ChildAlignment{ x: CLAY_ALIGN_X_CENTER, y: CLAY_ALIGN_Y_CENTER } } }, func() {
                        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_FIXED(8), height: CLAY_SIZING_FIXED(8)} }, backgroundColor: CLAY__DEBUGVIEW_COLOR_3, cornerRadius: CLAY_CORNER_RADIUS(2) }, func() {})
                    })
                }
                // Collisions and offscreen info
                if currentElementData != nil {
                    if currentElementData.debugData.collision {
                        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: 8, right: 8, top: 2, bottom: 2 }}, border: Clay_BorderElementConfig{ color: Clay_Color{r: 177, g: 147, b: 8, a: 255}, width: Clay_BorderWidth{left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0} } }, func() {
                            CLAY_TEXT(CLAY_STRING("Duplicate ID"), CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_3, fontSize: 16 }))
                        })
                    }
                    if offscreen {
                        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: 8, right: 8, top: 2, bottom: 2 } }, border: Clay_BorderElementConfig{  color: CLAY__DEBUGVIEW_COLOR_3, width: Clay_BorderWidth{ left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0} } }, func() {
                            CLAY_TEXT(CLAY_STRING("Offscreen"), CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_3, fontSize: 16 }))
                        })
                    }
                }
                idString := *(*Clay_String)(unsafe.Pointer(uintptr(unsafe.Pointer(context.layoutElementIdStrings.internalArray)) + uintptr(currentElementIndex)*unsafe.Sizeof(Clay_String{})))
                if idString.length > 0 {
                    var textConfig *Clay_TextElementConfig
                    if offscreen {
                        textConfig = &Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_3, fontSize: 16 }
                    } else {
                        textConfig = &Clay__DebugView_TextNameConfig
                    }
                    CLAY_TEXT(idString, textConfig)
                }
                for elementConfigIndex := int32(0); elementConfigIndex < currentElement.elementConfigs.length; elementConfigIndex++ {
                    elementConfig := Clay__ElementConfigArraySlice_Get(&currentElement.elementConfigs, elementConfigIndex)
                    if elementConfig._type == CLAY__ELEMENT_CONFIG_TYPE_SHARED {
                        labelColor := Clay_Color{r: 243, g: 134, b: 48, a: 90}
                        labelColor.a = 90
                        backgroundColor := elementConfig.config.sharedElementConfig.backgroundColor
                        radius := elementConfig.config.sharedElementConfig.cornerRadius
                        if backgroundColor.a > 0 {
                            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: 8, right: 8, top: 2, bottom: 2 } }, backgroundColor: labelColor, cornerRadius: CLAY_CORNER_RADIUS(4), border: Clay_BorderElementConfig{ color: labelColor, width: Clay_BorderWidth{ left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0} } }, func() {
                                var textColor Clay_Color
                                if offscreen {
                                    textColor = CLAY__DEBUGVIEW_COLOR_3
                                } else {
                                    textColor = CLAY__DEBUGVIEW_COLOR_4
                                }
                                CLAY_TEXT(CLAY_STRING("Color"), CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: textColor, fontSize: 16 }))
                            })
                        }
                        if radius.bottomLeft > 0 {
                            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: 8, right: 8, top: 2, bottom: 2 } }, backgroundColor: labelColor, cornerRadius: CLAY_CORNER_RADIUS(4), border: Clay_BorderElementConfig{ color: labelColor, width: Clay_BorderWidth{ left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0 } } }, func() {
                                var textColor Clay_Color
                                if offscreen {
                                    textColor = CLAY__DEBUGVIEW_COLOR_3
                                } else {
                                    textColor = CLAY__DEBUGVIEW_COLOR_4
                                }
                                CLAY_TEXT(CLAY_STRING("Radius"), CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: textColor, fontSize: 16 }))
                            })
                        }
                        continue
                    }
                    config := Clay__DebugGetElementConfigTypeLabel(elementConfig._type)
                    backgroundColor := config.color
                    backgroundColor.a = 90
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: 8, right: 8, top: 2, bottom: 2 } }, backgroundColor: backgroundColor, cornerRadius: CLAY_CORNER_RADIUS(4), border: Clay_BorderElementConfig{ color: config.color, width: Clay_BorderWidth{ left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0 } } }, func() {
                        var textColor Clay_Color
                        if offscreen {
                            textColor = CLAY__DEBUGVIEW_COLOR_3
                        } else {
                            textColor = CLAY__DEBUGVIEW_COLOR_4
                        }
                        CLAY_TEXT(config.label, CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: textColor, fontSize: 16 }))
                    })
                }
            })

            // Render the text contents below the element as a non-interactive row
            if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) {
                layoutData.rowCount++
                textElementData := currentElement.childrenOrTextContent.textElementData
                var rawTextConfig *Clay_TextElementConfig
                if offscreen {
                    rawTextConfig = &Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_3, fontSize: 16 }
                } else {
                    rawTextConfig = &Clay__DebugView_TextNameConfig
                }
                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT))}, childAlignment: Clay_ChildAlignment{ y: CLAY_ALIGN_Y_CENTER } } }, func() {
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_INDENT_WIDTH + 16)) } } }, func() {})
                    CLAY_TEXT(CLAY_STRING("\""), rawTextConfig)
                    var displayText Clay_String
                    if textElementData.text.length > 40 {
                        displayText = Clay_String{ length: 40, chars: textElementData.text.chars }
                    } else {
                        displayText = textElementData.text
                    }
                    CLAY_TEXT(displayText, rawTextConfig)
                    if textElementData.text.length > 40 {
                        CLAY_TEXT(CLAY_STRING("..."), rawTextConfig)
                    }
                    CLAY_TEXT(CLAY_STRING("\""), rawTextConfig)
                })
            } else if currentElement.childrenOrTextContent.children.length > 0 {
                Clay__OpenElement()
                Clay__ConfigureOpenElement(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: 8 } } })
                Clay__OpenElement()
                Clay__ConfigureOpenElement(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: uint16(CLAY__DEBUGVIEW_INDENT_WIDTH) }}, border: Clay_BorderElementConfig{ color: CLAY__DEBUGVIEW_COLOR_3, width: Clay_BorderWidth{ left: 1 } }})
                Clay__OpenElement()
                Clay__ConfigureOpenElement(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_TOP_TO_BOTTOM } })
            }

            layoutData.rowCount++
            if !(Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (currentElementData != nil && currentElementData.debugData.collapsed)) {
                for i := currentElement.childrenOrTextContent.children.length - 1; i >= 0; i-- {
                    Clay__int32_tArray_Add(&dfsBuffer, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))
                    *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length - 1)*unsafe.Sizeof(bool(false)))) = false // TODO needs to be ranged checked
                }
            }
        }
    }

    if context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
        collapseButtonId := Clay__HashString(CLAY_STRING("Clay__DebugView_CollapseElement"), 0)
        for i := int(context.pointerOverIds.length) - 1; i >= 0; i-- {
            elementId := Clay_ElementIdArray_Get(&context.pointerOverIds, int32(i))
            if elementId.baseId == collapseButtonId.baseId {
                highlightedItem := Clay__GetHashMapItem(elementId.offset)
                highlightedItem.debugData.collapsed = !highlightedItem.debugData.collapsed
                break
            }
        }
    }

    if highlightedElementId != 0 {
        CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugView_ElementHighlight"), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_GROW(0, 0)} }, floating: Clay_FloatingElementConfig{ parentId: highlightedElementId, zIndex: 32767, pointerCaptureMode: CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH, attachTo: CLAY_ATTACH_TO_ELEMENT_WITH_ID } }, func() {
            CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugView_ElementHighlightRectangle"), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_GROW(0, 0)} }, backgroundColor: Clay__debugViewHighlightColor }, func() {})
        })
    }
    return layoutData
}

func Clay__RenderDebugLayoutSizing(sizing Clay_SizingAxis, infoTextConfig *Clay_TextElementConfig) {
    sizingLabel := CLAY_STRING("GROW")
    if sizing.type_ == CLAY__SIZING_TYPE_FIT {
        sizingLabel = CLAY_STRING("FIT")
    } else if sizing.type_ == CLAY__SIZING_TYPE_PERCENT {
        sizingLabel = CLAY_STRING("PERCENT")
    } else if sizing.type_ == CLAY__SIZING_TYPE_FIXED {
        sizingLabel = CLAY_STRING("FIXED")
    }
    CLAY_TEXT(sizingLabel, infoTextConfig)
    if sizing.type_ == CLAY__SIZING_TYPE_GROW || sizing.type_ == CLAY__SIZING_TYPE_FIT || sizing.type_ == CLAY__SIZING_TYPE_FIXED {
        CLAY_TEXT(CLAY_STRING("("), infoTextConfig)
        if sizing.size.minMax.min != 0 {
            CLAY_TEXT(CLAY_STRING("min: "), infoTextConfig)
            CLAY_TEXT(Clay__IntToString(int32(sizing.size.minMax.min)), infoTextConfig)
            if sizing.size.minMax.max != CLAY__MAXFLOAT {
                CLAY_TEXT(CLAY_STRING(", "), infoTextConfig)
            }
        }
        if sizing.size.minMax.max != CLAY__MAXFLOAT {
            CLAY_TEXT(CLAY_STRING("max: "), infoTextConfig)
            CLAY_TEXT(Clay__IntToString(int32(sizing.size.minMax.max)), infoTextConfig)
        }
        CLAY_TEXT(CLAY_STRING(")"), infoTextConfig)
    } else if sizing.type_ == CLAY__SIZING_TYPE_PERCENT {
        CLAY_TEXT(CLAY_STRING("("), infoTextConfig)
        CLAY_TEXT(Clay__IntToString(int32(sizing.size.percent * 100)), infoTextConfig)
        CLAY_TEXT(CLAY_STRING("%)"), infoTextConfig)
    }
}

func Clay__RenderDebugViewElementConfigHeader(elementId Clay_String, _type Clay__ElementConfigType) {
    config := Clay__DebugGetElementConfigTypeLabel(_type)
    backgroundColor := config.color
    backgroundColor.a = 90
    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_GROW(0, 0) }, padding: CLAY_PADDING_ALL(float32(CLAY__DEBUGVIEW_OUTER_PADDING)), childAlignment: Clay_ChildAlignment{ y: CLAY_ALIGN_Y_CENTER } } }, func() {
        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: Clay_Padding{ left: 8, right: 8, top: 2, bottom: 2 } }, backgroundColor: backgroundColor, cornerRadius: CLAY_CORNER_RADIUS(4), border: Clay_BorderElementConfig{ color: config.color, width: Clay_BorderWidth{ left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0 } } }, func() {
            CLAY_TEXT(config.label, CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_4, fontSize: 16 }))
        })
        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_GROW(0, 0) } } }, func() {})
        CLAY_TEXT(elementId, CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_3, fontSize: 16, wrapMode: CLAY_TEXT_WRAP_NONE }))
    })
}

func Clay__RenderDebugViewColor(color Clay_Color, textConfig *Clay_TextElementConfig) {
    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ childAlignment: Clay_ChildAlignment{y: CLAY_ALIGN_Y_CENTER} } }, func() {
        CLAY_TEXT(CLAY_STRING("{ r: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(color.r)), textConfig)
        CLAY_TEXT(CLAY_STRING(", g: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(color.g)), textConfig)
        CLAY_TEXT(CLAY_STRING(", b: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(color.b)), textConfig)
        CLAY_TEXT(CLAY_STRING(", a: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(color.a)), textConfig)
        CLAY_TEXT(CLAY_STRING(" }"), textConfig)
        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_FIXED(10) } } }, func() {})
        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT - 8)), height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT - 8))} }, backgroundColor: color, cornerRadius: CLAY_CORNER_RADIUS(4), border: Clay_BorderElementConfig{ color: CLAY__DEBUGVIEW_COLOR_4, width: Clay_BorderWidth{ left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0 } } }, func() {})
    })
}

func Clay__RenderDebugViewCornerRadius(cornerRadius Clay_CornerRadius, textConfig *Clay_TextElementConfig) {
    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ childAlignment: Clay_ChildAlignment{y: CLAY_ALIGN_Y_CENTER} } }, func() {
        CLAY_TEXT(CLAY_STRING("{ topLeft: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(cornerRadius.topLeft)), textConfig)
        CLAY_TEXT(CLAY_STRING(", topRight: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(cornerRadius.topRight)), textConfig)
        CLAY_TEXT(CLAY_STRING(", bottomLeft: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(cornerRadius.bottomLeft)), textConfig)
        CLAY_TEXT(CLAY_STRING(", bottomRight: "), textConfig)
        CLAY_TEXT(Clay__IntToString(int32(cornerRadius.bottomRight)), textConfig)
        CLAY_TEXT(CLAY_STRING(" }"), textConfig)
    })
}

func HandleDebugViewCloseButtonInteraction(elementId Clay_ElementId, pointerInfo Clay_PointerData, userData uintptr) {
    context := Clay_GetCurrentContext()
    _ = elementId
    _ = userData
    if pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
        context.debugModeEnabled = false
    }
}

func Clay__RenderDebugView() {
    context := Clay_GetCurrentContext()
    closeButtonId := Clay__HashString(CLAY_STRING("Clay__DebugViewTopHeaderCloseButtonOuter"), 0)
    if context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
        for i := int32(0); i < context.pointerOverIds.length; i++ {
            elementId := Clay_ElementIdArray_Get(&context.pointerOverIds, i)
            if elementId.id == closeButtonId.id {
                context.debugModeEnabled = false
                return
            }
        }
    }

    initialRootsLength := uint32(context.layoutElementTreeRoots.length)
    initialElementsLength := uint32(context.layoutElements.length)
    infoTextConfig := &Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_4, fontSize: 16, wrapMode: CLAY_TEXT_WRAP_NONE }
    infoTitleConfig := &Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_3, fontSize: 16, wrapMode: CLAY_TEXT_WRAP_NONE }
    scrollId := Clay__HashString(CLAY_STRING("Clay__DebugViewOuterScrollPane"), 0)
    scrollYOffset := float32(0)
    pointerInDebugView := context.pointerInfo.position.y < context.layoutDimensions.height - 300
    for i := int32(0); i < context.scrollContainerDatas.length; i++ {
        scrollContainerData := Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i)
        if scrollContainerData.elementId == scrollId.id {
            if !context.externalScrollHandlingEnabled {
                scrollYOffset = scrollContainerData.scrollPosition.y
            } else {
                pointerInDebugView = context.pointerInfo.position.y + scrollContainerData.scrollPosition.y < context.layoutDimensions.height - 300
            }
            break
        }
    }
    var highlightedRow int32
    if pointerInDebugView {
        highlightedRow = int32((context.pointerInfo.position.y - scrollYOffset) / float32(CLAY__DEBUGVIEW_ROW_HEIGHT)) - 1
    } else {
        highlightedRow = -1
    }
    if context.pointerInfo.position.x < context.layoutDimensions.width - float32(Clay__debugViewWidth) {
        highlightedRow = -1
    }
    layoutData := Clay__RenderDebugLayoutData{}
    CLAY(Clay_ElementDeclaration{ 
        id: CLAY_ID("Clay__DebugView"),
        layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_FIXED(float32(Clay__debugViewWidth)), height: CLAY_SIZING_FIXED(context.layoutDimensions.height) }, layoutDirection: CLAY_TOP_TO_BOTTOM },
        floating: Clay_FloatingElementConfig{ zIndex: 32765, attachPoints: Clay_FloatingAttachPoints{ element: CLAY_ATTACH_POINT_LEFT_CENTER, parent: CLAY_ATTACH_POINT_RIGHT_CENTER }, attachTo: CLAY_ATTACH_TO_ROOT, clipTo: CLAY_CLIP_TO_ATTACHED_PARENT },
        border: Clay_BorderElementConfig{ color: CLAY__DEBUGVIEW_COLOR_3, width: Clay_BorderWidth{ bottom: 1 } },
    }, func() {
        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT))}, padding: Clay_Padding{left: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), right: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), top: 0, bottom: 0 }, childAlignment: Clay_ChildAlignment{y: CLAY_ALIGN_Y_CENTER} }, backgroundColor: CLAY__DEBUGVIEW_COLOR_2 }, func() {
            CLAY_TEXT(CLAY_STRING("Clay Debug Tools"), infoTextConfig)
            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_GROW(0, 0) } } }, func() {})
            // Close button
            CLAY(Clay_ElementDeclaration{
                layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT - 10)), height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT - 10))}, childAlignment: Clay_ChildAlignment{x: CLAY_ALIGN_X_CENTER, y: CLAY_ALIGN_Y_CENTER} },
                backgroundColor: Clay_Color{r: 217, g: 91, b: 67, a: 80},
                cornerRadius: CLAY_CORNER_RADIUS(4),
                border: Clay_BorderElementConfig{ color: Clay_Color{ r: 217, g: 91, b: 67, a: 255 }, width: Clay_BorderWidth{ left: 1, right: 1, top: 1, bottom: 1, betweenChildren: 0 } },
            }, func() {
                Clay_OnHover(HandleDebugViewCloseButtonInteraction, 0)
                CLAY_TEXT(CLAY_STRING("x"), CLAY_TEXT_CONFIG(Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_4, fontSize: 16 }))
            })
        })
        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(1)} }, backgroundColor: CLAY__DEBUGVIEW_COLOR_3 }, func() {})
        CLAY(Clay_ElementDeclaration{ id: scrollId, layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_GROW(0, 0)} }, clip: Clay_ClipElementConfig{ horizontal: true, vertical: true, childOffset: Clay_GetScrollOffset() } }, func() {
            var bgColor Clay_Color
            if ((initialElementsLength + initialRootsLength) & 1) == 0 {
                bgColor = CLAY__DEBUGVIEW_COLOR_2
            } else {
                bgColor = CLAY__DEBUGVIEW_COLOR_1
            }
            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_GROW(0, 0)}, layoutDirection: CLAY_TOP_TO_BOTTOM }, backgroundColor: bgColor }, func() {
                panelContentsId := Clay__HashString(CLAY_STRING("Clay__DebugViewPaneOuter"), 0)
                // Element list
                CLAY(Clay_ElementDeclaration{ id: panelContentsId, layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_GROW(0, 0)} }, floating: Clay_FloatingElementConfig{ zIndex: 32766, pointerCaptureMode: CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH, attachTo: CLAY_ATTACH_TO_PARENT, clipTo: CLAY_CLIP_TO_ATTACHED_PARENT } }, func() {
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_GROW(0, 0)}, padding: Clay_Padding{ left: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), right: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), top: 0, bottom: 0 }, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                        layoutData = Clay__RenderDebugLayoutElementsList(int32(initialRootsLength), highlightedRow)
                    })
                })
                contentWidth := Clay__GetHashMapItem(panelContentsId.id).layoutElement.dimensions.width
                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_FIXED(contentWidth) }, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {})
                for i := int32(0); i < layoutData.rowCount; i++ {
                    var rowColor Clay_Color
                    if (i & 1) == 0 {
                        rowColor = CLAY__DEBUGVIEW_COLOR_2
                    } else {
                        rowColor = CLAY__DEBUGVIEW_COLOR_1
                    }
                    if i == layoutData.selectedElementRowIndex {
                        rowColor = CLAY__DEBUGVIEW_COLOR_SELECTED_ROW
                    }
                    if i == highlightedRow {
                        rowColor.r = min(1.0, rowColor.r * 1.25)
                        rowColor.g = min(1.0, rowColor.g * 1.25)
                        rowColor.b = min(1.0, rowColor.b * 1.25)
                    }
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT))}, layoutDirection: CLAY_TOP_TO_BOTTOM }, backgroundColor: rowColor }, func() {})
                }
            })
        })
        CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(1)} }, backgroundColor: CLAY__DEBUGVIEW_COLOR_3 }, func() {})
        if context.debugSelectedElementId != 0 {
            selectedItem := Clay__GetHashMapItem(context.debugSelectedElementId)
            CLAY(Clay_ElementDeclaration{
                layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(300)}, layoutDirection: CLAY_TOP_TO_BOTTOM },
                backgroundColor: CLAY__DEBUGVIEW_COLOR_2,
                clip: Clay_ClipElementConfig{ vertical: true, childOffset: Clay_GetScrollOffset() },
                border: Clay_BorderElementConfig{ color: CLAY__DEBUGVIEW_COLOR_3, width: Clay_BorderWidth{ betweenChildren: 1 } },
            }, func() {
                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT + 8))}, padding: Clay_Padding{left: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), right: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), top: 0, bottom: 0 }, childAlignment: Clay_ChildAlignment{y: CLAY_ALIGN_Y_CENTER} } }, func() {
                    CLAY_TEXT(CLAY_STRING("Layout Config"), infoTextConfig)
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_GROW(0, 0) } } }, func() {})
                    if selectedItem.elementId.stringId.length != 0 {
                        CLAY_TEXT(selectedItem.elementId.stringId, infoTitleConfig)
                        if selectedItem.elementId.offset != 0 {
                            CLAY_TEXT(CLAY_STRING(" ("), infoTitleConfig)
                            CLAY_TEXT(Clay__IntToString(int32(selectedItem.elementId.offset)), infoTitleConfig)
                            CLAY_TEXT(CLAY_STRING(")"), infoTitleConfig)
                        }
                    }
                })
                attributeConfigPadding := Clay_Padding{left: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), right: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), top: 8, bottom: 8}
                // Clay_LayoutConfig debug info
                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                    // .boundingBox
                    CLAY_TEXT(CLAY_STRING("Bounding Box"), infoTitleConfig)
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                        CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(selectedItem.boundingBox.x)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(selectedItem.boundingBox.y)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(", width: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(selectedItem.boundingBox.width)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(selectedItem.boundingBox.height)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig)
                    })
                    // .layoutDirection
                    CLAY_TEXT(CLAY_STRING("Layout Direction"), infoTitleConfig)
                    layoutConfig := selectedItem.layoutElement.layoutConfig
                    var directionText Clay_String
                    if layoutConfig.layoutDirection == CLAY_TOP_TO_BOTTOM {
                        directionText = CLAY_STRING("TOP_TO_BOTTOM")
                    } else {
                        directionText = CLAY_STRING("LEFT_TO_RIGHT")
                    }
                    CLAY_TEXT(directionText, infoTextConfig)
                    // .sizing
                    CLAY_TEXT(CLAY_STRING("Sizing"), infoTitleConfig)
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                        CLAY_TEXT(CLAY_STRING("width: "), infoTextConfig)
                        Clay__RenderDebugLayoutSizing(layoutConfig.sizing.width, infoTextConfig)
                    })
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                        CLAY_TEXT(CLAY_STRING("height: "), infoTextConfig)
                        Clay__RenderDebugLayoutSizing(layoutConfig.sizing.height, infoTextConfig)
                    })
                    // .padding
                    CLAY_TEXT(CLAY_STRING("Padding"), infoTitleConfig)
                    CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewElementInfoPadding") }, func() {
                        CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(layoutConfig.padding.left)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(layoutConfig.padding.right)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(layoutConfig.padding.top)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig)
                        CLAY_TEXT(Clay__IntToString(int32(layoutConfig.padding.bottom)), infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig)
                    })
                    // .childGap
                    CLAY_TEXT(CLAY_STRING("Child Gap"), infoTitleConfig)
                    CLAY_TEXT(Clay__IntToString(int32(layoutConfig.childGap)), infoTextConfig)
                    // .childAlignment
                    CLAY_TEXT(CLAY_STRING("Child Alignment"), infoTitleConfig)
                    CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                        CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig)
                        alignX := CLAY_STRING("LEFT")
                        if layoutConfig.childAlignment.x == CLAY_ALIGN_X_CENTER {
                            alignX = CLAY_STRING("CENTER")
                        } else if layoutConfig.childAlignment.x == CLAY_ALIGN_X_RIGHT {
                            alignX = CLAY_STRING("RIGHT")
                        }
                        CLAY_TEXT(alignX, infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig)
                        alignY := CLAY_STRING("TOP")
                        if layoutConfig.childAlignment.y == CLAY_ALIGN_Y_CENTER {
                            alignY = CLAY_STRING("CENTER")
                        } else if layoutConfig.childAlignment.y == CLAY_ALIGN_Y_BOTTOM {
                            alignY = CLAY_STRING("BOTTOM")
                        }
                        CLAY_TEXT(alignY, infoTextConfig)
                        CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig)
                    })
                })
                for elementConfigIndex := int32(0); elementConfigIndex < selectedItem.layoutElement.elementConfigs.length; elementConfigIndex++ {
                    elementConfig := Clay__ElementConfigArraySlice_Get(&selectedItem.layoutElement.elementConfigs, elementConfigIndex)
                    Clay__RenderDebugViewElementConfigHeader(selectedItem.elementId.stringId, elementConfig._type)
                    switch elementConfig._type {
                        case CLAY__ELEMENT_CONFIG_TYPE_SHARED:
                            sharedConfig := elementConfig.config.sharedElementConfig
                            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM }}, func() {
                                // .backgroundColor
                                CLAY_TEXT(CLAY_STRING("Background Color"), infoTitleConfig)
                                Clay__RenderDebugViewColor(sharedConfig.backgroundColor, infoTextConfig)
                                // .cornerRadius
                                CLAY_TEXT(CLAY_STRING("Corner Radius"), infoTitleConfig)
                                Clay__RenderDebugViewCornerRadius(sharedConfig.cornerRadius, infoTextConfig)
                            })
                        case CLAY__ELEMENT_CONFIG_TYPE_TEXT:
                            textConfig := elementConfig.config.textElementConfig
                            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                                // .fontSize
                                CLAY_TEXT(CLAY_STRING("Font Size"), infoTitleConfig)
                                CLAY_TEXT(Clay__IntToString(int32(textConfig.fontSize)), infoTextConfig)
                                // .fontId
                                CLAY_TEXT(CLAY_STRING("Font ID"), infoTitleConfig)
                                CLAY_TEXT(Clay__IntToString(int32(textConfig.fontId)), infoTextConfig)
                                // .lineHeight
                                CLAY_TEXT(CLAY_STRING("Line Height"), infoTitleConfig)
                                var lineHeightText Clay_String
                                if textConfig.lineHeight == 0 {
                                    lineHeightText = CLAY_STRING("auto")
                                } else {
                                    lineHeightText = Clay__IntToString(int32(textConfig.lineHeight))
                                }
                                CLAY_TEXT(lineHeightText, infoTextConfig)
                                // .letterSpacing
                                CLAY_TEXT(CLAY_STRING("Letter Spacing"), infoTitleConfig)
                                CLAY_TEXT(Clay__IntToString(int32(textConfig.letterSpacing)), infoTextConfig)
                                // .wrapMode
                                CLAY_TEXT(CLAY_STRING("Wrap Mode"), infoTitleConfig)
                                wrapMode := CLAY_STRING("WORDS")
                                if textConfig.wrapMode == CLAY_TEXT_WRAP_NONE {
                                    wrapMode = CLAY_STRING("NONE")
                                } else if textConfig.wrapMode == CLAY_TEXT_WRAP_NEWLINES {
                                    wrapMode = CLAY_STRING("NEWLINES")
                                }
                                CLAY_TEXT(wrapMode, infoTextConfig)
                                // .textAlignment
                                CLAY_TEXT(CLAY_STRING("Text Alignment"), infoTitleConfig)
                                textAlignment := CLAY_STRING("LEFT")
                                if textConfig.textAlignment == CLAY_TEXT_ALIGN_CENTER {
                                    textAlignment = CLAY_STRING("CENTER")
                                } else if textConfig.textAlignment == CLAY_TEXT_ALIGN_RIGHT {
                                    textAlignment = CLAY_STRING("RIGHT")
                                }
                                CLAY_TEXT(textAlignment, infoTextConfig)
                                // .textColor
                                CLAY_TEXT(CLAY_STRING("Text Color"), infoTitleConfig)
                                Clay__RenderDebugViewColor(textConfig.textColor, infoTextConfig)
                            })
                        case CLAY__ELEMENT_CONFIG_TYPE_ASPECT:
                            aspectRatioConfig := elementConfig.config.aspectRatioElementConfig
                            CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewElementInfoAspectRatioBody"), layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                                CLAY_TEXT(CLAY_STRING("Aspect Ratio"), infoTitleConfig)
                                // Aspect Ratio
                                CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewElementInfoAspectRatio") }, func() {
                                    CLAY_TEXT(Clay__IntToString(int32(aspectRatioConfig.aspectRatio)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING("."), infoTextConfig)
                                    frac := aspectRatioConfig.aspectRatio - float32(int32(aspectRatioConfig.aspectRatio))
                                    frac *= 100
                                    if int32(frac) < 10 {
                                        CLAY_TEXT(CLAY_STRING("0"), infoTextConfig)
                                    }
                                    CLAY_TEXT(Clay__IntToString(int32(frac)), infoTextConfig)
                                })
                            })
                        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE:
                            imageConfig := elementConfig.config.imageElementConfig
                            aspectConfig := Clay_AspectRatioElementConfig{ aspectRatio: 1 }
                            if Clay__ElementHasConfig(selectedItem.layoutElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT) {
                                aspectConfig = *Clay__FindElementConfigWithType(selectedItem.layoutElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig
                            }
                            CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewElementInfoImageBody"), layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                                // Image Preview
                                CLAY_TEXT(CLAY_STRING("Preview"), infoTitleConfig)
                                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_GROW(64, 128), height: CLAY_SIZING_GROW(64, 128) }}, aspectRatio: aspectConfig, image: *imageConfig }, func() {})
                            })
                        case CLAY__ELEMENT_CONFIG_TYPE_CLIP:
                            clipConfig := elementConfig.config.clipElementConfig
                            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                                // .vertical
                                CLAY_TEXT(CLAY_STRING("Vertical"), infoTitleConfig)
                                var verticalText Clay_String
                                if clipConfig.vertical {
                                    verticalText = CLAY_STRING("true")
                                } else {
                                    verticalText = CLAY_STRING("false")
                                }
                                CLAY_TEXT(verticalText, infoTextConfig)
                                // .horizontal
                                CLAY_TEXT(CLAY_STRING("Horizontal"), infoTitleConfig)
                                var horizontalText Clay_String
                                if clipConfig.horizontal {
                                    horizontalText = CLAY_STRING("true")
                                } else {
                                    horizontalText = CLAY_STRING("false")
                                }
                                CLAY_TEXT(horizontalText, infoTextConfig)
                            })
                        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING:
                            floatingConfig := elementConfig.config.floatingElementConfig
                            CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                                // .offset
                                CLAY_TEXT(CLAY_STRING("Offset"), infoTitleConfig)
                                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                                    CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(floatingConfig.offset.x)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(floatingConfig.offset.y)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig)
                                })
                                // .expand
                                CLAY_TEXT(CLAY_STRING("Expand"), infoTitleConfig)
                                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                                    CLAY_TEXT(CLAY_STRING("{ width: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(floatingConfig.expand.width)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(floatingConfig.expand.height)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig)
                                })
                                // .zIndex
                                CLAY_TEXT(CLAY_STRING("z-index"), infoTitleConfig)
                                CLAY_TEXT(Clay__IntToString(int32(floatingConfig.zIndex)), infoTextConfig)
                                // .parentId
                                CLAY_TEXT(CLAY_STRING("Parent"), infoTitleConfig)
                                hashItem := Clay__GetHashMapItem(floatingConfig.parentId)
                                CLAY_TEXT(hashItem.elementId.stringId, infoTextConfig)
                                // .attachPoints
                                CLAY_TEXT(CLAY_STRING("Attach Points"), infoTitleConfig)
                                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                                    CLAY_TEXT(CLAY_STRING("{ element: "), infoTextConfig)
                                    attachPointElement := CLAY_STRING("LEFT_TOP")
                                    if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_LEFT_CENTER {
                                        attachPointElement = CLAY_STRING("LEFT_CENTER")
                                    } else if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_LEFT_BOTTOM {
                                        attachPointElement = CLAY_STRING("LEFT_BOTTOM")
                                    } else if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_CENTER_TOP {
                                        attachPointElement = CLAY_STRING("CENTER_TOP")
                                    } else if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_CENTER_CENTER {
                                        attachPointElement = CLAY_STRING("CENTER_CENTER")
                                    } else if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_CENTER_BOTTOM {
                                        attachPointElement = CLAY_STRING("CENTER_BOTTOM")
                                    } else if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_RIGHT_TOP {
                                        attachPointElement = CLAY_STRING("RIGHT_TOP")
                                    } else if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_RIGHT_CENTER {
                                        attachPointElement = CLAY_STRING("RIGHT_CENTER")
                                    } else if floatingConfig.attachPoints.element == CLAY_ATTACH_POINT_RIGHT_BOTTOM {
                                        attachPointElement = CLAY_STRING("RIGHT_BOTTOM")
                                    }
                                    CLAY_TEXT(attachPointElement, infoTextConfig)
                                    attachPointParent := CLAY_STRING("LEFT_TOP")
                                    if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_LEFT_CENTER {
                                        attachPointParent = CLAY_STRING("LEFT_CENTER")
                                    } else if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_LEFT_BOTTOM {
                                        attachPointParent = CLAY_STRING("LEFT_BOTTOM")
                                    } else if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_CENTER_TOP {
                                        attachPointParent = CLAY_STRING("CENTER_TOP")
                                    } else if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_CENTER_CENTER {
                                        attachPointParent = CLAY_STRING("CENTER_CENTER")
                                    } else if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_CENTER_BOTTOM {
                                        attachPointParent = CLAY_STRING("CENTER_BOTTOM")
                                    } else if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_TOP {
                                        attachPointParent = CLAY_STRING("RIGHT_TOP")
                                    } else if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_CENTER {
                                        attachPointParent = CLAY_STRING("RIGHT_CENTER")
                                    } else if floatingConfig.attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_BOTTOM {
                                        attachPointParent = CLAY_STRING("RIGHT_BOTTOM")
                                    }
                                    CLAY_TEXT(CLAY_STRING(", parent: "), infoTextConfig)
                                    CLAY_TEXT(attachPointParent, infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig)
                                })
                                // .pointerCaptureMode
                                CLAY_TEXT(CLAY_STRING("Pointer Capture Mode"), infoTitleConfig)
                                pointerCaptureMode := CLAY_STRING("NONE")
                                if floatingConfig.pointerCaptureMode == CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH {
                                    pointerCaptureMode = CLAY_STRING("PASSTHROUGH")
                                }
                                CLAY_TEXT(pointerCaptureMode, infoTextConfig)
                                // .attachTo
                                CLAY_TEXT(CLAY_STRING("Attach To"), infoTitleConfig)
                                attachTo := CLAY_STRING("NONE")
                                if floatingConfig.attachTo == CLAY_ATTACH_TO_PARENT {
                                    attachTo = CLAY_STRING("PARENT")
                                } else if floatingConfig.attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID {
                                    attachTo = CLAY_STRING("ELEMENT_WITH_ID")
                                } else if floatingConfig.attachTo == CLAY_ATTACH_TO_ROOT {
                                    attachTo = CLAY_STRING("ROOT")
                                }
                                CLAY_TEXT(attachTo, infoTextConfig)
                                // .clipTo
                                CLAY_TEXT(CLAY_STRING("Clip To"), infoTitleConfig)
                                clipTo := CLAY_STRING("ATTACHED_PARENT")
                                if floatingConfig.clipTo == CLAY_CLIP_TO_NONE {
                                    clipTo = CLAY_STRING("NONE")
                                }
                                CLAY_TEXT(clipTo, infoTextConfig)
                            })
                        case CLAY__ELEMENT_CONFIG_TYPE_BORDER:
                            borderConfig := elementConfig.config.borderElementConfig
                            CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewElementInfoBorderBody"), layout: Clay_LayoutConfig{ padding: attributeConfigPadding, childGap: 8, layoutDirection: CLAY_TOP_TO_BOTTOM } }, func() {
                                CLAY_TEXT(CLAY_STRING("Border Widths"), infoTitleConfig)
                                CLAY(Clay_ElementDeclaration{ layout: Clay_LayoutConfig{ layoutDirection: CLAY_LEFT_TO_RIGHT } }, func() {
                                    CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(borderConfig.width.left)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(borderConfig.width.right)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(borderConfig.width.top)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig)
                                    CLAY_TEXT(Clay__IntToString(int32(borderConfig.width.bottom)), infoTextConfig)
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig)
                                })
                                // .textColor
                                CLAY_TEXT(CLAY_STRING("Border Color"), infoTitleConfig)
                                Clay__RenderDebugViewColor(borderConfig.color, infoTextConfig)
                            })
                        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM:
                        default:
                    }
                }
            })
        } else {
            CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewWarningsScrollPane"), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(300)}, childGap: 6, layoutDirection: CLAY_TOP_TO_BOTTOM }, backgroundColor: CLAY__DEBUGVIEW_COLOR_2, clip: Clay_ClipElementConfig{ horizontal: true, vertical: true, childOffset: Clay_GetScrollOffset() } }, func() {
                warningConfig := &Clay_TextElementConfig{ textColor: CLAY__DEBUGVIEW_COLOR_4, fontSize: 16, wrapMode: CLAY_TEXT_WRAP_NONE }
                CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewWarningItemHeader"), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT))}, padding: Clay_Padding{left: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), right: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), top: 0, bottom: 0 }, childGap: 8, childAlignment: Clay_ChildAlignment{y: CLAY_ALIGN_Y_CENTER} } }, func() {
                    CLAY_TEXT(CLAY_STRING("Warnings"), warningConfig)
                })
                CLAY(Clay_ElementDeclaration{ id: CLAY_ID("Clay__DebugViewWarningsTopBorder"), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{ width: CLAY_SIZING_GROW(0, 0), height: CLAY_SIZING_FIXED(1)} }, backgroundColor: Clay_Color{r: 200, g: 200, b: 200, a: 255} }, func() {})
                previousWarningsLength := int32(context.warnings.length)
                for i := int32(0); i < previousWarningsLength; i++ {
                    warning := *(*Clay__Warning)(unsafe.Pointer(uintptr(unsafe.Pointer(context.warnings.internalArray)) + uintptr(i)*unsafe.Sizeof(Clay__Warning{})))
                    CLAY(Clay_ElementDeclaration{ id: CLAY_IDI("Clay__DebugViewWarningItem", uint32(i)), layout: Clay_LayoutConfig{ sizing: Clay_Sizing{height: CLAY_SIZING_FIXED(float32(CLAY__DEBUGVIEW_ROW_HEIGHT))}, padding: Clay_Padding{left: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), right: uint16(CLAY__DEBUGVIEW_OUTER_PADDING), top: 0, bottom: 0 }, childGap: 8, childAlignment: Clay_ChildAlignment{y: CLAY_ALIGN_Y_CENTER} } }, func() {
                        CLAY_TEXT(warning.baseMessage, warningConfig)
                        if warning.dynamicMessage.length > 0 {
                            CLAY_TEXT(warning.dynamicMessage, warningConfig)
                        }
                    })
                }
            })
        }
    })
}
// endregion

func Clay__WarningArray_Allocate_Arena(arena *Clay_Arena, capacity int32) Clay__WarningArray {
    totalSizeBytes := uint(capacity) * uint(unsafe.Sizeof(Clay_String{}))
    array := Clay__WarningArray{capacity: capacity, length: 0}
    nextAllocOffset := uintptr(arena.nextAllocation + (64 - (arena.nextAllocation % 64)))
    if uint64(nextAllocOffset) + uint64(totalSizeBytes) <= arena.capacity {
        array.internalArray = (*Clay__Warning)(unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + nextAllocOffset))
        arena.nextAllocation = uintptr(nextAllocOffset) + uintptr(totalSizeBytes)
    } else {
        Clay__currentContext.errorHandler.errorHandlerFunction(Clay_ErrorData{
            errorType: CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
            errorText: CLAY_STRING("Clay attempted to allocate memory in its arena, but ran out of capacity. Try increasing the capacity of the arena passed to Clay_Initialize()"),
            userData: Clay__currentContext.errorHandler.userData,
        })
    }
    return array
}

func Clay__WarningArray_Add(array *Clay__WarningArray, item Clay__Warning) *Clay__Warning {
    if array.length < array.capacity {
        *(*Clay__Warning)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length)*unsafe.Sizeof(Clay__Warning{}))) = item
        array.length++
        return (*Clay__Warning)(unsafe.Pointer(uintptr(unsafe.Pointer(array.internalArray)) + uintptr(array.length - 1)*unsafe.Sizeof(Clay__Warning{})))
    }
    return &CLAY__WARNING_DEFAULT
}

func Clay__Array_Allocate_Arena(capacity int32, itemSize uint32, arena *Clay_Arena) unsafe.Pointer {
    totalSizeBytes := uint(capacity) * uint(itemSize)
    nextAllocOffset := uintptr(arena.nextAllocation + ((64 - (arena.nextAllocation % 64)) & 63))
    if uint64(nextAllocOffset) + uint64(totalSizeBytes) <= arena.capacity {
        arena.nextAllocation = uintptr(nextAllocOffset) + uintptr(totalSizeBytes)
        return unsafe.Pointer(uintptr(unsafe.Pointer(&arena.memory[0])) + nextAllocOffset)
    } else {
        Clay__currentContext.errorHandler.errorHandlerFunction(Clay_ErrorData{
            errorType: CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
            errorText: CLAY_STRING("Clay attempted to allocate memory in its arena, but ran out of capacity. Try increasing the capacity of the arena passed to Clay_Initialize()"),
            userData: Clay__currentContext.errorHandler.userData,
        })
    }
    return nil
}

func Clay__Array_RangeCheck(index int32, length int32) bool {
    if index < length && index >= 0 {
        return true
    }
    context := Clay_GetCurrentContext()
    context.errorHandler.errorHandlerFunction(Clay_ErrorData{
        errorType: CLAY_ERROR_TYPE_INTERNAL_ERROR,
        errorText: CLAY_STRING("Clay attempted to make an out of bounds array access. This is an internal error and is likely a bug."),
        userData: context.errorHandler.userData,
    })
    return false
}

func Clay__Array_AddCapacityCheck(length int32, capacity int32) bool {
    if length < capacity {
        return true
    }
    context := Clay_GetCurrentContext()
    context.errorHandler.errorHandlerFunction(Clay_ErrorData{
        errorType: CLAY_ERROR_TYPE_INTERNAL_ERROR,
        errorText: CLAY_STRING("Clay attempted to make an out of bounds array access. This is an internal error and is likely a bug."),
        userData: context.errorHandler.userData,
    })
    return false
}

// PUBLIC API FROM HERE ---------------------------------------

func Clay_MinMemorySize() uint32 {
    fakeContext := Clay_Context{
        maxElementCount: Clay__defaultMaxElementCount,
        maxMeasureTextCacheWordCount: Clay__defaultMaxMeasureTextWordCacheCount,
        internalArena: Clay_Arena{
            capacity: ^uint64(0),
            memory: nil,
        },
    }
    currentContext := Clay_GetCurrentContext()
    if currentContext != nil {
        fakeContext.maxElementCount = currentContext.maxElementCount
        fakeContext.maxMeasureTextCacheWordCount = currentContext.maxMeasureTextCacheWordCount
    }
    // Reserve space in the arena for the context, important for calculating min memory size correctly
    Clay__Context_Allocate_Arena(&fakeContext.internalArena)
    Clay__InitializePersistentMemory(&fakeContext)
    Clay__InitializeEphemeralMemory(&fakeContext)
    return uint32(fakeContext.internalArena.nextAllocation) + 128
}

func Clay_CreateArenaWithCapacityAndMemory(capacity uint, memory unsafe.Pointer) Clay_Arena {
    arena := Clay_Arena{
        capacity: uint64(capacity),
        memory: (*[1 << 30]byte)(memory)[:capacity],
    }
    return arena
}

func Clay_SetMeasureTextFunction(measureTextFunction func(text Clay_StringSlice, config *Clay_TextElementConfig, userData unsafe.Pointer) Clay_Dimensions, userData unsafe.Pointer) {
    context := Clay_GetCurrentContext()
    Clay__MeasureText = measureTextFunction
    context.measureTextUserData = userData
}

func Clay_SetQueryScrollOffsetFunction(queryScrollOffsetFunction func(elementId uint32, userData unsafe.Pointer) Clay_Vector2, userData unsafe.Pointer) {
    context := Clay_GetCurrentContext()
    Clay__QueryScrollOffset = queryScrollOffsetFunction
    context.queryScrollOffsetUserData = userData
}

func Clay_SetLayoutDimensions(dimensions Clay_Dimensions) {
    Clay_GetCurrentContext().layoutDimensions = dimensions
}

func Clay_SetPointerState(position Clay_Vector2, isPointerDown bool) {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return
    }
    context.pointerInfo.position = position
    context.pointerOverIds.length = 0
    dfsBuffer := context.layoutElementChildrenBuffer
    for rootIndex := context.layoutElementTreeRoots.length - 1; rootIndex >= 0; rootIndex-- {
        dfsBuffer.length = 0
        root := Clay__LayoutElementTreeRootArray_Get(&context.layoutElementTreeRoots, rootIndex)
        Clay__int32_tArray_Add(&dfsBuffer, int32(root.layoutElementIndex))
        *(*bool)(unsafe.Pointer(context.treeNodeVisited.internalArray)) = false
        found := false
        for dfsBuffer.length > 0 {
            if *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length - 1)*unsafe.Sizeof(bool(false)))) {
                dfsBuffer.length--
                continue
            }
            *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length - 1)*unsafe.Sizeof(bool(false)))) = true
            currentElement := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&dfsBuffer, int32(dfsBuffer.length - 1)))
            mapItem := Clay__GetHashMapItem(currentElement.id) // TODO think of a way around this, maybe the fact that it's essentially a binary tree limits the cost, but the worst case is not great
            clipElementId := Clay__int32_tArray_GetValue(&context.layoutElementClipElementIds, int32(uintptr(unsafe.Pointer(currentElement)) - uintptr(unsafe.Pointer(context.layoutElements.internalArray))))
            clipItem := Clay__GetHashMapItem(uint32(clipElementId))
            if mapItem != nil {
                elementBox := mapItem.boundingBox
                elementBox.x -= root.pointerOffset.x
                elementBox.y -= root.pointerOffset.y
                if Clay__PointIsInsideRect(position, elementBox) && (clipElementId == 0 || Clay__PointIsInsideRect(position, clipItem.boundingBox) || context.externalScrollHandlingEnabled) {
                    if mapItem.onHoverFunction != nil {
                        mapItem.onHoverFunction(mapItem.elementId, context.pointerInfo, mapItem.hoverFunctionUserData)
                    }
                    Clay_ElementIdArray_Add(&context.pointerOverIds, mapItem.elementId)
                    found = true

                    if mapItem.idAlias != 0 {
                        Clay_ElementIdArray_Add(&context.pointerOverIds, Clay_ElementId{ id: mapItem.idAlias })
                    }
                }
                if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) {
                    dfsBuffer.length--
                    continue
                }
                for i := currentElement.childrenOrTextContent.children.length - 1; i >= 0; i-- {
                    Clay__int32_tArray_Add(&dfsBuffer, *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(currentElement.childrenOrTextContent.children.elements)) + uintptr(i)*unsafe.Sizeof(int32(0)))))
                    *(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(context.treeNodeVisited.internalArray)) + uintptr(dfsBuffer.length - 1)*unsafe.Sizeof(bool(false)))) = false // TODO needs to be ranged checked
                }
            } else {
                dfsBuffer.length--
            }
        }

        rootElement := Clay_LayoutElementArray_Get(&context.layoutElements, root.layoutElementIndex)
        if found && Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) &&
                Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig.pointerCaptureMode == CLAY_POINTER_CAPTURE_MODE_CAPTURE {
            break
        }
    }

    if isPointerDown {
        if context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
            context.pointerInfo.state = CLAY_POINTER_DATA_PRESSED
        } else if context.pointerInfo.state != CLAY_POINTER_DATA_PRESSED {
            context.pointerInfo.state = CLAY_POINTER_DATA_PRESSED_THIS_FRAME
        }
    } else {
        if context.pointerInfo.state == CLAY_POINTER_DATA_RELEASED_THIS_FRAME {
            context.pointerInfo.state = CLAY_POINTER_DATA_RELEASED
        } else if context.pointerInfo.state != CLAY_POINTER_DATA_RELEASED {
            context.pointerInfo.state = CLAY_POINTER_DATA_RELEASED_THIS_FRAME
        }
    }
}

func Clay_Initialize(arena Clay_Arena, layoutDimensions Clay_Dimensions, errorHandler Clay_ErrorHandler) *Clay_Context {
    // Cacheline align memory passed in
    baseOffset := uintptr(64 - (uintptr(unsafe.Pointer(&arena.memory[0])) % 64))
    if baseOffset == 64 {
        baseOffset = 0
    }
    arena.memory = arena.memory[baseOffset:]
    context := Clay__Context_Allocate_Arena(&arena)
    if context == nil {
        return nil
    }
    // DEFAULTS
    oldContext := Clay_GetCurrentContext()
    var maxElementCount int32
    var maxMeasureTextCacheWordCount int32
    if oldContext != nil {
        maxElementCount = oldContext.maxElementCount
        maxMeasureTextCacheWordCount = oldContext.maxMeasureTextCacheWordCount
    } else {
        maxElementCount = Clay__defaultMaxElementCount
        maxMeasureTextCacheWordCount = Clay__defaultMaxMeasureTextWordCacheCount
    }
    var errorHandlerToUse Clay_ErrorHandler
    if errorHandler.errorHandlerFunction != nil {
        errorHandlerToUse = errorHandler
    } else {
        errorHandlerToUse = Clay_ErrorHandler{ errorHandlerFunction: Clay__ErrorHandlerFunctionDefault, userData: nil }
    }
    *context = Clay_Context{
        maxElementCount: maxElementCount,
        maxMeasureTextCacheWordCount: maxMeasureTextCacheWordCount,
        errorHandler: errorHandlerToUse,
        layoutDimensions: layoutDimensions,
        internalArena: arena,
    }
    Clay_SetCurrentContext(context)
    Clay__InitializePersistentMemory(context)
    Clay__InitializeEphemeralMemory(context)
    for i := int32(0); i < context.layoutElementsHashMap.capacity; i++ {
        *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.layoutElementsHashMap.internalArray)) + uintptr(i)*unsafe.Sizeof(int32(0)))) = -1
    }
    for i := int32(0); i < context.measureTextHashMap.capacity; i++ {
        *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.measureTextHashMap.internalArray)) + uintptr(i)*unsafe.Sizeof(int32(0)))) = 0
    }
    context.measureTextHashMapInternal.length = 1 // Reserve the 0 value to mean "no next element"
    context.layoutDimensions = layoutDimensions
    return context
}

func Clay_GetCurrentContext() *Clay_Context {
    return Clay__currentContext
}

func Clay_SetCurrentContext(context *Clay_Context) {
    Clay__currentContext = context
}

func Clay_GetScrollOffset() Clay_Vector2 {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return Clay_Vector2{}
    }
    openLayoutElement := Clay__GetOpenLayoutElement()
    // If the element has no id attached at this point, we need to generate one
    if openLayoutElement.id == 0 {
        Clay__GenerateIdForAnonymousElement(openLayoutElement)
    }
    for i := int32(0); i < context.scrollContainerDatas.length; i++ {
        mapping := Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i)
        if mapping.layoutElement == openLayoutElement {
            return mapping.scrollPosition
        }
    }
    return Clay_Vector2{}
}

func Clay_UpdateScrollContainers(enableDragScrolling bool, scrollDelta Clay_Vector2, deltaTime float32) {
    context := Clay_GetCurrentContext()
    isPointerActive := enableDragScrolling && (context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED || context.pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME)
    // Don't apply scroll events to ancestors of the inner element
    highestPriorityElementIndex := int32(-1)
    var highestPriorityScrollData *Clay__ScrollContainerDataInternal
    for i := int32(0); i < context.scrollContainerDatas.length; i++ {
        scrollData := Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i)
        if !scrollData.openThisFrame {
            Clay__ScrollContainerDataInternalArray_RemoveSwapback(&context.scrollContainerDatas, i)
            continue
        }
        scrollData.openThisFrame = false
        hashMapItem := Clay__GetHashMapItem(scrollData.elementId)
        // Element isn't rendered this frame but scroll offset has been retained
        if hashMapItem == nil {
            Clay__ScrollContainerDataInternalArray_RemoveSwapback(&context.scrollContainerDatas, i)
            continue
        }

        // Touch / click is released
        if !isPointerActive && scrollData.pointerScrollActive {
            xDiff := scrollData.scrollPosition.x - scrollData.scrollOrigin.x
            if xDiff < -10 || xDiff > 10 {
                scrollData.scrollMomentum.x = (scrollData.scrollPosition.x - scrollData.scrollOrigin.x) / (scrollData.momentumTime * 25)
            }
            yDiff := scrollData.scrollPosition.y - scrollData.scrollOrigin.y
            if yDiff < -10 || yDiff > 10 {
                scrollData.scrollMomentum.y = (scrollData.scrollPosition.y - scrollData.scrollOrigin.y) / (scrollData.momentumTime * 25)
            }
            scrollData.pointerScrollActive = false

            scrollData.pointerOrigin = Clay_Vector2{x: 0, y: 0}
            scrollData.scrollOrigin = Clay_Vector2{x: 0, y: 0}
            scrollData.momentumTime = 0
        }

        // Apply existing momentum
        scrollData.scrollPosition.x += scrollData.scrollMomentum.x
        scrollData.scrollMomentum.x *= 0.95
        scrollOccurred := scrollDelta.x != 0 || scrollDelta.y != 0
        if (scrollData.scrollMomentum.x > -0.1 && scrollData.scrollMomentum.x < 0.1) || scrollOccurred {
            scrollData.scrollMomentum.x = 0
        }
        scrollData.scrollPosition.x = CLAY__MIN(CLAY__MAX(scrollData.scrollPosition.x, -(CLAY__MAX(scrollData.contentSize.width - scrollData.layoutElement.dimensions.width, 0))), 0)

        scrollData.scrollPosition.y += scrollData.scrollMomentum.y
        scrollData.scrollMomentum.y *= 0.95
        if (scrollData.scrollMomentum.y > -0.1 && scrollData.scrollMomentum.y < 0.1) || scrollOccurred {
            scrollData.scrollMomentum.y = 0
        }
        scrollData.scrollPosition.y = CLAY__MIN(CLAY__MAX(scrollData.scrollPosition.y, -(CLAY__MAX(scrollData.contentSize.height - scrollData.layoutElement.dimensions.height, 0))), 0)

        for j := int32(0); j < context.pointerOverIds.length; j++ { // TODO n & m are small here but this being n*m gives me the creeps
            if scrollData.layoutElement.id == Clay_ElementIdArray_Get(&context.pointerOverIds, j).id {
                highestPriorityElementIndex = j
                highestPriorityScrollData = scrollData
            }
        }
    }

    if highestPriorityElementIndex > -1 && highestPriorityScrollData != nil {
        scrollElement := highestPriorityScrollData.layoutElement
        clipConfig := Clay__FindElementConfigWithType(scrollElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig
        canScrollVertically := clipConfig.vertical && highestPriorityScrollData.contentSize.height > scrollElement.dimensions.height
        canScrollHorizontally := clipConfig.horizontal && highestPriorityScrollData.contentSize.width > scrollElement.dimensions.width
        // Handle wheel scroll
        if canScrollVertically {
            highestPriorityScrollData.scrollPosition.y = highestPriorityScrollData.scrollPosition.y + scrollDelta.y * 10
        }
        if canScrollHorizontally {
            highestPriorityScrollData.scrollPosition.x = highestPriorityScrollData.scrollPosition.x + scrollDelta.x * 10
        }
        // Handle click / touch scroll
        if isPointerActive {
            highestPriorityScrollData.scrollMomentum = Clay_Vector2{}
            if !highestPriorityScrollData.pointerScrollActive {
                highestPriorityScrollData.pointerOrigin = context.pointerInfo.position
                highestPriorityScrollData.scrollOrigin = highestPriorityScrollData.scrollPosition
                highestPriorityScrollData.pointerScrollActive = true
            } else {
                scrollDeltaX := float32(0)
                scrollDeltaY := float32(0)
                if canScrollHorizontally {
                    oldXScrollPosition := highestPriorityScrollData.scrollPosition.x
                    highestPriorityScrollData.scrollPosition.x = highestPriorityScrollData.scrollOrigin.x + (context.pointerInfo.position.x - highestPriorityScrollData.pointerOrigin.x)
                    highestPriorityScrollData.scrollPosition.x = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.x, 0), -(highestPriorityScrollData.contentSize.width - highestPriorityScrollData.boundingBox.width))
                    scrollDeltaX = highestPriorityScrollData.scrollPosition.x - oldXScrollPosition
                }
                if canScrollVertically {
                    oldYScrollPosition := highestPriorityScrollData.scrollPosition.y
                    highestPriorityScrollData.scrollPosition.y = highestPriorityScrollData.scrollOrigin.y + (context.pointerInfo.position.y - highestPriorityScrollData.pointerOrigin.y)
                    highestPriorityScrollData.scrollPosition.y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.y, 0), -(highestPriorityScrollData.contentSize.height - highestPriorityScrollData.boundingBox.height))
                    scrollDeltaY = highestPriorityScrollData.scrollPosition.y - oldYScrollPosition
                }
                if scrollDeltaX > -0.1 && scrollDeltaX < 0.1 && scrollDeltaY > -0.1 && scrollDeltaY < 0.1 && highestPriorityScrollData.momentumTime > 0.15 {
                    highestPriorityScrollData.momentumTime = 0
                    highestPriorityScrollData.pointerOrigin = context.pointerInfo.position
                    highestPriorityScrollData.scrollOrigin = highestPriorityScrollData.scrollPosition
                } else {
                     highestPriorityScrollData.momentumTime += deltaTime
                }
            }
        }
        // Clamp any changes to scroll position to the maximum size of the contents
        if canScrollVertically {
            highestPriorityScrollData.scrollPosition.y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.y, 0), -(highestPriorityScrollData.contentSize.height - scrollElement.dimensions.height))
        }
        if canScrollHorizontally {
            highestPriorityScrollData.scrollPosition.x = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.scrollPosition.x, 0), -(highestPriorityScrollData.contentSize.width - scrollElement.dimensions.width))
        }
    }
}

func Clay_BeginLayout() {
    context := Clay_GetCurrentContext()
    Clay__InitializeEphemeralMemory(context)
    context.generation++
    context.dynamicElementIndex = 0
    // Set up the root container that covers the entire window
    rootDimensions := Clay_Dimensions{width: context.layoutDimensions.width, height: context.layoutDimensions.height}
    if context.debugModeEnabled {
        rootDimensions.width -= float32(Clay__debugViewWidth)
    }
    context.booleanWarnings = Clay_BooleanWarnings{}
    Clay__OpenElement()
    Clay__ConfigureOpenElement(Clay_ElementDeclaration{
        id: CLAY_ID("Clay__RootContainer"),
        layout: Clay_LayoutConfig{ sizing: Clay_Sizing{width: CLAY_SIZING_FIXED(rootDimensions.width), height: CLAY_SIZING_FIXED(rootDimensions.height)} },
    })
    Clay__int32_tArray_Add(&context.openLayoutElementStack, 0)
    Clay__LayoutElementTreeRootArray_Add(&context.layoutElementTreeRoots, Clay__LayoutElementTreeRoot{ layoutElementIndex: 0 })
}

func Clay_EndLayout() Clay_RenderCommandArray {
    context := Clay_GetCurrentContext()
    Clay__CloseElement()
    elementsExceededBeforeDebugView := context.booleanWarnings.maxElementsExceeded
    if context.debugModeEnabled && !elementsExceededBeforeDebugView {
        context.warningsEnabled = false
        Clay__RenderDebugView()
        context.warningsEnabled = true
    }
    if context.booleanWarnings.maxElementsExceeded {
        var message Clay_String
        if !elementsExceededBeforeDebugView {
            message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount after adding the debug-view to the layout.")
        } else {
            message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount")
        }
        Clay__AddRenderCommand(Clay_RenderCommand{
            boundingBox: Clay_BoundingBox{ x: context.layoutDimensions.width / 2 - 59 * 4, y: context.layoutDimensions.height / 2, width: 0, height: 0 },
            renderData: Clay_RenderData{ text: Clay_TextRenderData{ stringContents: Clay_StringSlice{ length: message.length, chars: message.chars, baseChars: message.chars }, textColor: Clay_Color{r: 255, g: 0, b: 0, a: 255}, fontSize: 16 } },
            commandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
        })
    } else {
        Clay__CalculateFinalLayout()
    }
    return context.renderCommands
}

func Clay_GetElementId(idString Clay_String) Clay_ElementId {
    return Clay__HashString(idString, 0)
}

func Clay_GetElementIdWithIndex(idString Clay_String, index uint32) Clay_ElementId {
    return Clay__HashStringWithOffset(idString, index, 0)
}

func Clay_Hovered() bool {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return false
    }
    openLayoutElement := Clay__GetOpenLayoutElement()
    // If the element has no id attached at this point, we need to generate one
    if openLayoutElement.id == 0 {
        Clay__GenerateIdForAnonymousElement(openLayoutElement)
    }
    for i := int32(0); i < context.pointerOverIds.length; i++ {
        if Clay_ElementIdArray_Get(&context.pointerOverIds, i).id == openLayoutElement.id {
            return true
        }
    }
    return false
}

func Clay_OnHover(onHoverFunction func(elementId Clay_ElementId, pointerInfo Clay_PointerData, userData uintptr), userData uintptr) {
    context := Clay_GetCurrentContext()
    if context.booleanWarnings.maxElementsExceeded {
        return
    }
    openLayoutElement := Clay__GetOpenLayoutElement()
    if openLayoutElement.id == 0 {
        Clay__GenerateIdForAnonymousElement(openLayoutElement)
    }
    hashMapItem := Clay__GetHashMapItem(openLayoutElement.id)
    hashMapItem.onHoverFunction = onHoverFunction
    hashMapItem.hoverFunctionUserData = userData
}

func Clay_PointerOver(elementId Clay_ElementId) bool { // TODO return priority for separating multiple results
    context := Clay_GetCurrentContext()
    for i := int32(0); i < context.pointerOverIds.length; i++ {
        if Clay_ElementIdArray_Get(&context.pointerOverIds, i).id == elementId.id {
            return true
        }
    }
    return false;
}

//export Clay_GetScrollContainerData
func Clay_GetScrollContainerData(id Clay_ElementId) Clay_ScrollContainerData {
    context := Clay_GetCurrentContext()
    for i := int32(0); i < context.scrollContainerDatas.length; i++ {
        scrollContainerData := Clay__ScrollContainerDataInternalArray_Get(&context.scrollContainerDatas, i)
        if scrollContainerData.elementId == id.id {
            clipElementConfig := Clay__FindElementConfigWithType(scrollContainerData.layoutElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig
            if clipElementConfig == nil { // This can happen on the first frame before a scroll container is declared
                return Clay_ScrollContainerData{}
            }
            return Clay_ScrollContainerData{
                scrollPosition: &scrollContainerData.scrollPosition,
                scrollContainerDimensions: Clay_Dimensions{ scrollContainerData.boundingBox.width, scrollContainerData.boundingBox.height },
                contentDimensions: scrollContainerData.contentSize,
                config: *clipElementConfig,
                found: true,
            }
        }
    }
    return Clay_ScrollContainerData{}
}

//export Clay_GetElementData
func Clay_GetElementData(id Clay_ElementId) Clay_ElementData {
    item := Clay__GetHashMapItem(id.id)
    if item == &Clay_LayoutElementHashMapItem_DEFAULT {
        return Clay_ElementData{}
    }

    return Clay_ElementData{
        boundingBox: item.boundingBox,
        found: true,
    }
}

func Clay_SetDebugModeEnabled(enabled bool) {
    context := Clay_GetCurrentContext()
    context.debugModeEnabled = enabled
}

func Clay_IsDebugModeEnabled() bool {
    context := Clay_GetCurrentContext()
    return context.debugModeEnabled
}

func Clay_SetCullingEnabled(enabled bool) {
    context := Clay_GetCurrentContext()
    context.disableCulling = !enabled
}

func Clay_SetExternalScrollHandlingEnabled(enabled bool) {
    context := Clay_GetCurrentContext()
    context.externalScrollHandlingEnabled = enabled
}

func Clay_GetMaxElementCount() int32 {
    context := Clay_GetCurrentContext()
    return context.maxElementCount
}

func Clay_SetMaxElementCount(maxElementCount int32) {
    context := Clay_GetCurrentContext()
    if context != nil {
        context.maxElementCount = maxElementCount
    } else {
        Clay__defaultMaxElementCount = maxElementCount // TODO: Fix this
        Clay__defaultMaxMeasureTextWordCacheCount = maxElementCount * 2
    }
}

func Clay_GetMaxMeasureTextCacheWordCount() int32 {
    context := Clay_GetCurrentContext()
    return context.maxMeasureTextCacheWordCount
}

func Clay_SetMaxMeasureTextCacheWordCount(maxMeasureTextCacheWordCount int32) {
    context := Clay_GetCurrentContext()
    if context != nil {
        Clay__currentContext.maxMeasureTextCacheWordCount = maxMeasureTextCacheWordCount
    } else {
        Clay__defaultMaxMeasureTextWordCacheCount = maxMeasureTextCacheWordCount // TODO: Fix this
    }
}

func Clay_ResetMeasureTextCache() {
    context := Clay_GetCurrentContext()
    context.measureTextHashMapInternal.length = 0
    context.measureTextHashMapInternalFreeList.length = 0
    context.measureTextHashMap.length = 0
    context.measuredWords.length = 0
    context.measuredWordsFreeList.length = 0
    
    for i := int32(0); i < context.measureTextHashMap.capacity; i++ {
        *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(context.measureTextHashMap.internalArray)) + uintptr(i)*unsafe.Sizeof(int32(0)))) = 0
    }
    context.measureTextHashMapInternal.length = 1 // Reserve the 0 value to mean "no next element"
}
