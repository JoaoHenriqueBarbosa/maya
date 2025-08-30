package clay

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
    return Clay_BorderWidth{widthValue, widthValue, widthValue, widthValue, 0}
}

// CLAY_BORDER_ALL macro converted to function
func CLAY_BORDER_ALL(widthValue float32) Clay_BorderWidth {
    return Clay_BorderWidth{widthValue, widthValue, widthValue, widthValue, widthValue}
}

// CLAY_CORNER_RADIUS macro converted to function
func CLAY_CORNER_RADIUS(radius float32) Clay_CornerRadius {
    return Clay_CornerRadius{radius, radius, radius, radius}
}

// CLAY_PADDING_ALL macro converted to function
func CLAY_PADDING_ALL(padding float32) Clay_Padding {
    return Clay_Padding{padding, padding, padding, padding}
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
    return CLAY_SID_LOCAL(CLAY_STRING(label))
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
        chars:                 s,
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
    chars string
}

// Clay_StringSlice is used to represent non owning string slices, and includes
// a baseChars field which points to the string this slice is derived from.
type Clay_StringSlice struct {
    length    int32
    chars     string
    baseChars string // The source string / char* that this slice was derived from
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
    internalArray []Clay_ElementId
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

// Function Forward Declarations ---------------------------------

// Public API functions ------------------------------------------

// Returns the size, in bytes, of the minimum amount of memory Clay requires to operate at its current settings.
func Clay_MinMemorySize() uint32

// Creates an arena for clay to use for its internal allocations, given a certain capacity in bytes and a pointer to an allocation of at least that size.
// Intended to be used with Clay_MinMemorySize in the following way:
// uint32_t minMemoryRequired = Clay_MinMemorySize();
// Clay_Arena clayMemory = Clay_CreateArenaWithCapacityAndMemory(minMemoryRequired, malloc(minMemoryRequired));
func Clay_CreateArenaWithCapacityAndMemory(capacity uintptr, memory interface{}) Clay_Arena

// Sets the state of the "pointer" (i.e. the mouse or touch) in Clay's internal data. Used for detecting and responding to mouse events in the debug view,
// as well as for Clay_Hovered() and scroll element handling.
func Clay_SetPointerState(position Clay_Vector2, pointerDown bool)

// Initialize Clay's internal arena and setup required data before layout can begin. Only needs to be called once.
// - arena can be created using Clay_CreateArenaWithCapacityAndMemory()
// - layoutDimensions are the initial bounding dimensions of the layout (i.e. the screen width and height for a full screen layout)
// - errorHandler is used by Clay to inform you if something has gone wrong in configuration or layout.
func Clay_Initialize(arena Clay_Arena, layoutDimensions Clay_Dimensions, errorHandler Clay_ErrorHandler) *Clay_Context

// Returns the Context that clay is currently using. Used when using multiple instances of clay simultaneously.
func Clay_GetCurrentContext() *Clay_Context
// Sets the context that clay will use to compute the layout.
// Used to restore a context saved from Clay_GetCurrentContext when using multiple instances of clay simultaneously.
func Clay_SetCurrentContext(context *Clay_Context)
// Updates the state of Clay's internal scroll data, updating scroll content positions if scrollDelta is non zero, and progressing momentum scrolling.
// - enableDragScrolling when set to true will enable mobile device like "touch drag" scroll of scroll containers, including momentum scrolling after the touch has ended.
// - scrollDelta is the amount to scroll this frame on each axis in pixels.
// - deltaTime is the time in seconds since the last "frame" (scroll update)
func Clay_UpdateScrollContainers(enableDragScrolling bool, scrollDelta Clay_Vector2, deltaTime float32)
// Returns the internally stored scroll offset for the currently open element.
// Generally intended for use with clip elements to create scrolling containers.
func Clay_GetScrollOffset() Clay_Vector2
// Updates the layout dimensions in response to the window or outer container being resized.
func Clay_SetLayoutDimensions(dimensions Clay_Dimensions)
// Called before starting any layout declarations.
func Clay_BeginLayout()
// Called when all layout declarations are finished.
// Computes the layout and generates and returns the array of render commands to draw.
func Clay_EndLayout() Clay_RenderCommandArray
// Calculates a hash ID from the given idString.
// Generally only used for dynamic strings when CLAY_ID("stringLiteral") can't be used.
func Clay_GetElementId(idString Clay_String) Clay_ElementId
// Calculates a hash ID from the given idString and index.
// - index is used to avoid constructing dynamic ID strings in loops.
// Generally only used for dynamic strings when CLAY_IDI("stringLiteral", index) can't be used.
func Clay_GetElementIdWithIndex(idString Clay_String, index uint32) Clay_ElementId
// Returns layout data such as the final calculated bounding box for an element with a given ID.
// The returned Clay_ElementData contains a `found` bool that will be true if an element with the provided ID was found.
// This ID can be calculated either with CLAY_ID() for string literal IDs, or Clay_GetElementId for dynamic strings.
func Clay_GetElementData(id Clay_ElementId) Clay_ElementData
// Returns true if the pointer position provided by Clay_SetPointerState is within the current element's bounding box.
// Works during element declaration, e.g. CLAY({ .backgroundColor = Clay_Hovered() ? BLUE : RED });
func Clay_Hovered() bool
// Bind a callback that will be called when the pointer position provided by Clay_SetPointerState is within the current element's bounding box.
// - onHoverFunction is a function pointer to a user defined function.
// - userData is a pointer that will be transparently passed through when the onHoverFunction is called.
func Clay_OnHover(onHoverFunction func(elementId Clay_ElementId, pointerData Clay_PointerData, userData uintptr), userData uintptr)
// An imperative function that returns true if the pointer position provided by Clay_SetPointerState is within the element with the provided ID's bounding box.
// This ID can be calculated either with CLAY_ID() for string literal IDs, or Clay_GetElementId for dynamic strings.
func Clay_PointerOver(elementId Clay_ElementId) bool
// Returns the array of element IDs that the pointer is currently over.
func Clay_GetPointerOverIds() Clay_ElementIdArray
// Returns data representing the state of the scrolling element with the provided ID.
// The returned Clay_ScrollContainerData contains a `found` bool that will be true if a scroll element was found with the provided ID.
// An imperative function that returns true if the pointer position provided by Clay_SetPointerState is within the element with the provided ID's bounding box.
// This ID can be calculated either with CLAY_ID() for string literal IDs, or Clay_GetElementId for dynamic strings.
func Clay_GetScrollContainerData(id Clay_ElementId) Clay_ScrollContainerData
// Binds a callback function that Clay will call to determine the dimensions of a given string slice.
// - measureTextFunction is a user provided function that adheres to the interface Clay_Dimensions (Clay_StringSlice text, Clay_TextElementConfig *config, void *userData);
// - userData is a pointer that will be transparently passed through when the measureTextFunction is called.
func Clay_SetMeasureTextFunction(measureTextFunction func(text Clay_StringSlice, config *Clay_TextElementConfig, userData interface{}) Clay_Dimensions, userData interface{})
// Experimental - Used in cases where Clay needs to integrate with a system that manages its own scrolling containers externally.
// Please reach out if you plan to use this function, as it may be subject to change.
func Clay_SetQueryScrollOffsetFunction(queryScrollOffsetFunction func(elementId uint32, userData interface{}) Clay_Vector2, userData interface{})
// A bounds-checked "get" function for the Clay_RenderCommandArray returned from Clay_EndLayout().
func Clay_RenderCommandArray_Get(array *Clay_RenderCommandArray, index int32) *Clay_RenderCommand
// Enables and disables Clay's internal debug tools.
// This state is retained and does not need to be set each frame.
func Clay_SetDebugModeEnabled(enabled bool)
// Returns true if Clay's internal debug tools are currently enabled.
func Clay_IsDebugModeEnabled() bool
// Enables and disables visibility culling. By default, Clay will not generate render commands for elements whose bounding box is entirely outside the screen.
func Clay_SetCullingEnabled(enabled bool)
// Returns the maximum number of UI elements supported by Clay's current configuration.
func Clay_GetMaxElementCount() int32
// Modifies the maximum number of UI elements supported by Clay's current configuration.
// This may require reallocating additional memory, and re-calling Clay_Initialize();
func Clay_SetMaxElementCount(maxElementCount int32)
// Returns the maximum number of measured "words" (whitespace seperated runs of characters) that Clay can store in its internal text measurement cache.
func Clay_GetMaxMeasureTextCacheWordCount() int32
// Modifies the maximum number of measured "words" (whitespace seperated runs of characters) that Clay can store in its internal text measurement cache.
// This may require reallocating additional memory, and re-calling Clay_Initialize();
func Clay_SetMaxMeasureTextCacheWordCount(maxMeasureTextCacheWordCount int32)
// Resets Clay's internal text measurement cache. Useful if font mappings have changed or fonts have been reloaded.
func Clay_ResetMeasureTextCache()

// Internal API functions required by macros ----------------------

func Clay__OpenElement()
func Clay__ConfigureOpenElement(config Clay_ElementDeclaration)
func Clay__ConfigureOpenElementPtr(config *Clay_ElementDeclaration)
func Clay__CloseElement()
func Clay__HashString(key Clay_String, seed uint32) Clay_ElementId
func Clay__HashStringWithOffset(key Clay_String, offset uint32, seed uint32) Clay_ElementId
func Clay__OpenTextElement(text Clay_String, textConfig *Clay_TextElementConfig)
func Clay__StoreTextElementConfig(config Clay_TextElementConfig) *Clay_TextElementConfig
func Clay__GetParentElementId() uint32

var Clay__debugViewHighlightColor Clay_Color
var Clay__debugViewWidth uint32

// End of header section

// -----------------------------------------
// IMPLEMENTATION --------------------------
// -----------------------------------------

const CLAY__NULL = 0
const CLAY__MAXFLOAT = 3.40282346638528859812e+38

var CLAY_LAYOUT_DEFAULT Clay_LayoutConfig

var Clay__Color_DEFAULT Clay_Color
var Clay__CornerRadius_DEFAULT Clay_CornerRadius
var Clay__BorderWidth_DEFAULT Clay_BorderWidth

// The below comment describes the array macros that were converted to Go
// Array types and functions will be generated as needed when encountering their usage

var Clay__currentContext *Clay_Context
var Clay__defaultMaxElementCount int32 = 8192
var Clay__defaultMaxMeasureTextWordCacheCount int32 = 16384

func Clay__ErrorHandlerFunctionDefault(errorText Clay_ErrorData) {
    _ = errorText
}

var CLAY__SPACECHAR = Clay_String{length: 1, chars: " "}
var CLAY__STRING_DEFAULT = Clay_String{length: 0, chars: ""}

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

func Clay__WarningArray_Allocate_Arena(capacity int32, arena *Clay_Arena) Clay__WarningArray
func Clay__WarningArray_Add(array *Clay__WarningArray, item Clay__Warning) *Clay__Warning
func Clay__Array_Allocate_Arena(capacity int32, itemSize uint32, arena *Clay_Arena) interface{}
func Clay__Array_RangeCheck(index int32, length int32) bool
func Clay__Array_AddCapacityCheck(length int32, capacity int32) bool

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

func Clay__Context_Allocate_Arena(arena *Clay_Arena) *Clay_Context {
    var totalSizeBytes uintptr = unsafe.Sizeof(Clay_Context{})
    if totalSizeBytes > arena.capacity {
        return nil
    }
    arena.nextAllocation += totalSizeBytes
    return (*Clay_Context)(arena.memory)
}

func Clay__WriteStringToCharBuffer(buffer *Clay__charArray, string Clay_String) Clay_String {
    for i := int32(0); i < string.length; i++ {
        *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(buffer.internalArray)) + uintptr(buffer.length+i))) = string.chars[i]
    }
    buffer.length += string.length
    return Clay_String{
        length: string.length,
        chars:  string.chars, // Will need proper pointer arithmetic
    }
}

// Function pointers for text measurement and scroll offset
var Clay__MeasureText func(text Clay_StringSlice, config *Clay_TextElementConfig, userData interface{}) Clay_Dimensions
var Clay__QueryScrollOffset func(elementId uint32, userData interface{}) Clay_Vector2

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
        hash += uint32(key.chars[i])
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
        base += uint32(key.chars[i])
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
        hash += uint32(uintptr(unsafe.Pointer(&text.chars[0])))
        hash += (hash << 10)
        hash ^= (hash >> 6)
        hash += uint32(text.length)
        hash += (hash << 10)
        hash ^= (hash >> 6)
    } else {
        hash = uint32(Clay__HashData((*uint8)(unsafe.Pointer(&text.chars[0])), uintptr(text.length)) % 0xFFFFFFFF)
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
        newItemIndex := Clay__int32_tArray_GetValue(&context.measuredWordsFreeList, int(context.measuredWordsFreeList.length-1))
        context.measuredWordsFreeList.length--
        Clay__MeasuredWordArray_Set(&context.measuredWords, int(newItemIndex), word)
        previousWord.next = int32(newItemIndex)
        return Clay__MeasuredWordArray_Get(&context.measuredWords, int(newItemIndex))
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
                    errorText: Clay_String{chars: "Clay's internal MeasureText function is null. You may have forgotten to call Clay_SetMeasureTextFunction(), or passed a NULL function pointer by mistake.", length: 156},
                    userData: context.errorHandler.userData})
        }
        return &Clay__MeasureTextCacheItem_DEFAULT
    }
    // TODO: Complete conversion of this large function
    id := Clay__HashStringContentsWithConfig(text, config)
    _ = id
    // Temporarily return default to continue conversion
    return &Clay__MeasureTextCacheItem_DEFAULT
}

func Clay__PointIsInsideRect(point Clay_Vector2, rect Clay_BoundingBox) bool {
    return point.x >= rect.x && point.x <= rect.x+rect.width && point.y >= rect.y && point.y <= rect.y+rect.height
}

func Clay__AddHashMapItem(elementId Clay_ElementId, layoutElement *Clay_LayoutElement, idAlias uint32) *Clay_LayoutElementHashMapItem {
    // TODO: Complete conversion of hash map logic
    _ = elementId
    _ = layoutElement
    _ = idAlias
    return nil
}

func Clay__GetHashMapItem(id uint32) *Clay_LayoutElementHashMapItem {
    // TODO: Complete hash map implementation
    _ = id
    return &Clay_LayoutElementHashMapItem_DEFAULT
}

func Clay__GenerateIdForAnonymousElement(openLayoutElement *Clay_LayoutElement) Clay_ElementId {
    context := Clay_GetCurrentContext()
    parentElement := Clay_LayoutElementArray_Get(&context.layoutElements, Clay__int32_tArray_GetValue(&context.openLayoutElementStack, int(context.openLayoutElementStack.length-2)))
    elementId := Clay__HashNumber(uint32(parentElement.childrenOrTextContent.children.length), parentElement.id)
    openLayoutElement.id = elementId.id
    Clay__AddHashMapItem(elementId, openLayoutElement, 0)
    Clay__StringArray_Add(&context.layoutElementIdStrings, elementId.stringId)
    return elementId
}

func Clay__ElementHasConfig(layoutElement *Clay_LayoutElement, _type Clay__ElementConfigType) bool {
    for i := int32(0); i < layoutElement.elementConfigs.length; i++ {
        if Clay__ElementConfigArraySlice_Get(&layoutElement.elementConfigs, int(i))._type == _type {
            return true;
        }
    }
    return false;
}

func Clay__UpdateAspectRatioBox(layoutElement *Clay_LayoutElement) {
    for j := int32(0); j < layoutElement.elementConfigs.length; j++ {
        config := Clay__ElementConfigArraySlice_Get(&layoutElement.elementConfigs, int(j))
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

void Clay__CloseElement(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    Clay_LayoutConfig *layoutConfig = openLayoutElement->layoutConfig;
    bool elementHasClipHorizontal = false;
    bool elementHasClipVertical = false;
    for (int32_t i = 0; i < openLayoutElement->elementConfigs.length; i++) {
        Clay_ElementConfig *config = Clay__ElementConfigArraySlice_Get(&openLayoutElement->elementConfigs, i);
        if (config->type == CLAY__ELEMENT_CONFIG_TYPE_CLIP) {
            elementHasClipHorizontal = config->config.clipElementConfig->horizontal;
            elementHasClipVertical = config->config.clipElementConfig->vertical;
            context->openClipElementStack.length--;
            break;
        } else if (config->type == CLAY__ELEMENT_CONFIG_TYPE_FLOATING) {
            context->openClipElementStack.length--;
        }
    }

    float leftRightPadding = (float)(layoutConfig->padding.left + layoutConfig->padding.right);
    float topBottomPadding = (float)(layoutConfig->padding.top + layoutConfig->padding.bottom);

    // Attach children to the current open element
    openLayoutElement->childrenOrTextContent.children.elements = &context->layoutElementChildren.internalArray[context->layoutElementChildren.length];
    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
        openLayoutElement->dimensions.width = leftRightPadding;
        openLayoutElement->minDimensions.width = leftRightPadding;
        for (int32_t i = 0; i < openLayoutElement->childrenOrTextContent.children.length; i++) {
            int32_t childIndex = Clay__int32_tArray_GetValue(&context->layoutElementChildrenBuffer, (int)context->layoutElementChildrenBuffer.length - openLayoutElement->childrenOrTextContent.children.length + i);
            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, childIndex);
            openLayoutElement->dimensions.width += child->dimensions.width;
            openLayoutElement->dimensions.height = CLAY__MAX(openLayoutElement->dimensions.height, child->dimensions.height + topBottomPadding);
            // Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
            if (!elementHasClipHorizontal) {
                openLayoutElement->minDimensions.width += child->minDimensions.width;
            }
            if (!elementHasClipVertical) {
                openLayoutElement->minDimensions.height = CLAY__MAX(openLayoutElement->minDimensions.height, child->minDimensions.height + topBottomPadding);
            }
            Clay__int32_tArray_Add(&context->layoutElementChildren, childIndex);
        }
        float childGap = (float)(CLAY__MAX(openLayoutElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
        openLayoutElement->dimensions.width += childGap;
        if (!elementHasClipHorizontal) {
            openLayoutElement->minDimensions.width += childGap;
        }
    }
    else if (layoutConfig->layoutDirection == CLAY_TOP_TO_BOTTOM) {
        openLayoutElement->dimensions.height = topBottomPadding;
        openLayoutElement->minDimensions.height = topBottomPadding;
        for (int32_t i = 0; i < openLayoutElement->childrenOrTextContent.children.length; i++) {
            int32_t childIndex = Clay__int32_tArray_GetValue(&context->layoutElementChildrenBuffer, (int)context->layoutElementChildrenBuffer.length - openLayoutElement->childrenOrTextContent.children.length + i);
            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, childIndex);
            openLayoutElement->dimensions.height += child->dimensions.height;
            openLayoutElement->dimensions.width = CLAY__MAX(openLayoutElement->dimensions.width, child->dimensions.width + leftRightPadding);
            // Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
            if (!elementHasClipVertical) {
                openLayoutElement->minDimensions.height += child->minDimensions.height;
            }
            if (!elementHasClipHorizontal) {
                openLayoutElement->minDimensions.width = CLAY__MAX(openLayoutElement->minDimensions.width, child->minDimensions.width + leftRightPadding);
            }
            Clay__int32_tArray_Add(&context->layoutElementChildren, childIndex);
        }
        float childGap = (float)(CLAY__MAX(openLayoutElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
        openLayoutElement->dimensions.height += childGap;
        if (!elementHasClipVertical) {
            openLayoutElement->minDimensions.height += childGap;
        }
    }

    context->layoutElementChildrenBuffer.length -= openLayoutElement->childrenOrTextContent.children.length;

    // Clamp element min and max width to the values configured in the layout
    if (layoutConfig->sizing.width.type != CLAY__SIZING_TYPE_PERCENT) {
        if (layoutConfig->sizing.width.size.minMax.max <= 0) { // Set the max size if the user didn't specify, makes calculations easier
            layoutConfig->sizing.width.size.minMax.max = CLAY__MAXFLOAT;
        }
        openLayoutElement->dimensions.width = CLAY__MIN(CLAY__MAX(openLayoutElement->dimensions.width, layoutConfig->sizing.width.size.minMax.min), layoutConfig->sizing.width.size.minMax.max);
        openLayoutElement->minDimensions.width = CLAY__MIN(CLAY__MAX(openLayoutElement->minDimensions.width, layoutConfig->sizing.width.size.minMax.min), layoutConfig->sizing.width.size.minMax.max);
    } else {
        openLayoutElement->dimensions.width = 0;
    }

    // Clamp element min and max height to the values configured in the layout
    if (layoutConfig->sizing.height.type != CLAY__SIZING_TYPE_PERCENT) {
        if (layoutConfig->sizing.height.size.minMax.max <= 0) { // Set the max size if the user didn't specify, makes calculations easier
            layoutConfig->sizing.height.size.minMax.max = CLAY__MAXFLOAT;
        }
        openLayoutElement->dimensions.height = CLAY__MIN(CLAY__MAX(openLayoutElement->dimensions.height, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
        openLayoutElement->minDimensions.height = CLAY__MIN(CLAY__MAX(openLayoutElement->minDimensions.height, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
    } else {
        openLayoutElement->dimensions.height = 0;
    }

    Clay__UpdateAspectRatioBox(openLayoutElement);

    bool elementIsFloating = Clay__ElementHasConfig(openLayoutElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING);

    // Close the currently open element
    int32_t closingElementIndex = Clay__int32_tArray_RemoveSwapback(&context->openLayoutElementStack, (int)context->openLayoutElementStack.length - 1);
    openLayoutElement = Clay__GetOpenLayoutElement();

    if (!elementIsFloating && context->openLayoutElementStack.length > 1) {
        openLayoutElement->childrenOrTextContent.children.length++;
        Clay__int32_tArray_Add(&context->layoutElementChildrenBuffer, closingElementIndex);
    }
}

bool Clay__MemCmp(const char *s1, const char *s2, int32_t length);
#if !defined(CLAY_DISABLE_SIMD) && (defined(__x86_64__) || defined(_M_X64) || defined(_M_AMD64))
    bool Clay__MemCmp(const char *s1, const char *s2, int32_t length) {
        while (length >= 16) {
            __m128i v1 = _mm_loadu_si128((const __m128i *)s1);
            __m128i v2 = _mm_loadu_si128((const __m128i *)s2);

            if (_mm_movemask_epi8(_mm_cmpeq_epi8(v1, v2)) != 0xFFFF) { // If any byte differs
                return false;
            }

            s1 += 16;
            s2 += 16;
            length -= 16;
        }

        // Handle remaining bytes
        while (length--) {
            if (*s1 != *s2) {
                return false;
            }
            s1++;
            s2++;
        }

        return true;
    }
#elif !defined(CLAY_DISABLE_SIMD) && defined(__aarch64__)
    bool Clay__MemCmp(const char *s1, const char *s2, int32_t length) {
        while (length >= 16) {
            uint8x16_t v1 = vld1q_u8((const uint8_t *)s1);
            uint8x16_t v2 = vld1q_u8((const uint8_t *)s2);

            // Compare vectors
            if (vminvq_u32(vreinterpretq_u32_u8(vceqq_u8(v1, v2))) != 0xFFFFFFFF) { // If there's a difference
                return false;
            }

            s1 += 16;
            s2 += 16;
            length -= 16;
        }

        // Handle remaining bytes
        while (length--) {
            if (*s1 != *s2) {
                return false;
            }
            s1++;
            s2++;
        }

        return true;
    }
#else
    bool Clay__MemCmp(const char *s1, const char *s2, int32_t length) {
        for (int32_t i = 0; i < length; i++) {
            if (s1[i] != s2[i]) {
                return false;
            }
        }
        return true;
    }
#endif

void Clay__OpenElement(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->layoutElements.length == context->layoutElements.capacity - 1 || context->booleanWarnings.maxElementsExceeded) {
        context->booleanWarnings.maxElementsExceeded = true;
        return;
    }
    Clay_LayoutElement layoutElement = CLAY__DEFAULT_STRUCT;
    Clay_LayoutElementArray_Add(&context->layoutElements, layoutElement);
    Clay__int32_tArray_Add(&context->openLayoutElementStack, context->layoutElements.length - 1);
    if (context->openClipElementStack.length > 0) {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1));
    } else {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, 0);
    }
}

void Clay__OpenTextElement(Clay_String text, Clay_TextElementConfig *textConfig) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->layoutElements.length == context->layoutElements.capacity - 1 || context->booleanWarnings.maxElementsExceeded) {
        context->booleanWarnings.maxElementsExceeded = true;
        return;
    }
    Clay_LayoutElement *parentElement = Clay__GetOpenLayoutElement();

    Clay_LayoutElement layoutElement = CLAY__DEFAULT_STRUCT;
    Clay_LayoutElement *textElement = Clay_LayoutElementArray_Add(&context->layoutElements, layoutElement);
    if (context->openClipElementStack.length > 0) {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1));
    } else {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, 0);
    }

    Clay__int32_tArray_Add(&context->layoutElementChildrenBuffer, context->layoutElements.length - 1);
    Clay__MeasureTextCacheItem *textMeasured = Clay__MeasureTextCached(&text, textConfig);
    Clay_ElementId elementId = Clay__HashNumber(parentElement->childrenOrTextContent.children.length, parentElement->id);
    textElement->id = elementId.id;
    Clay__AddHashMapItem(elementId, textElement, 0);
    Clay__StringArray_Add(&context->layoutElementIdStrings, elementId.stringId);
    Clay_Dimensions textDimensions = { .width = textMeasured->unwrappedDimensions.width, .height = textConfig->lineHeight > 0 ? (float)textConfig->lineHeight : textMeasured->unwrappedDimensions.height };
    textElement->dimensions = textDimensions;
    textElement->minDimensions = CLAY__INIT(Clay_Dimensions) { .width = textMeasured->minWidth, .height = textDimensions.height };
    textElement->childrenOrTextContent.textElementData = Clay__TextElementDataArray_Add(&context->textElementData, CLAY__INIT(Clay__TextElementData) { .text = text, .preferredDimensions = textMeasured->unwrappedDimensions, .elementIndex = context->layoutElements.length - 1 });
    textElement->elementConfigs = CLAY__INIT(Clay__ElementConfigArraySlice) {
            .length = 1,
            .internalArray = Clay__ElementConfigArray_Add(&context->elementConfigs, CLAY__INIT(Clay_ElementConfig) { .type = CLAY__ELEMENT_CONFIG_TYPE_TEXT, .config = { .textElementConfig = textConfig }})
    };
    textElement->layoutConfig = &CLAY_LAYOUT_DEFAULT;
    parentElement->childrenOrTextContent.children.length++;
}

Clay_ElementId Clay__AttachId(Clay_ElementId elementId) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return Clay_ElementId_DEFAULT;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    uint32_t idAlias = openLayoutElement->id;
    openLayoutElement->id = elementId.id;
    Clay__AddHashMapItem(elementId, openLayoutElement, idAlias);
    Clay__StringArray_Set(&context->layoutElementIdStrings, context->layoutElements.length - 1, elementId.stringId);
    return elementId;
}

void Clay__ConfigureOpenElementPtr(const Clay_ElementDeclaration *declaration) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    openLayoutElement->layoutConfig = Clay__StoreLayoutConfig(declaration->layout);
    if ((declaration->layout.sizing.width.type == CLAY__SIZING_TYPE_PERCENT && declaration->layout.sizing.width.size.percent > 1) || (declaration->layout.sizing.height.type == CLAY__SIZING_TYPE_PERCENT && declaration->layout.sizing.height.size.percent > 1)) {
        context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                .errorType = CLAY_ERROR_TYPE_PERCENTAGE_OVER_1,
                .errorText = CLAY_STRING("An element was configured with CLAY_SIZING_PERCENT, but the provided percentage value was over 1.0. Clay expects a value between 0 and 1, i.e. 20% is 0.2."),
                .userData = context->errorHandler.userData });
    }

    Clay_ElementId openLayoutElementId = declaration->id;

    openLayoutElement->elementConfigs.internalArray = &context->elementConfigs.internalArray[context->elementConfigs.length];
    Clay_SharedElementConfig *sharedConfig = NULL;
    if (declaration->backgroundColor.a > 0) {
        sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .backgroundColor = declaration->backgroundColor });
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
    }
    if (!Clay__MemCmp((char *)(&declaration->cornerRadius), (char *)(&Clay__CornerRadius_DEFAULT), sizeof(Clay_CornerRadius))) {
        if (sharedConfig) {
            sharedConfig->cornerRadius = declaration->cornerRadius;
        } else {
            sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .cornerRadius = declaration->cornerRadius });
            Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
        }
    }
    if (declaration->userData != 0) {
        if (sharedConfig) {
            sharedConfig->userData = declaration->userData;
        } else {
            sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .userData = declaration->userData });
            Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
        }
    }
    if (declaration->image.imageData) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .imageElementConfig = Clay__StoreImageElementConfig(declaration->image) }, CLAY__ELEMENT_CONFIG_TYPE_IMAGE);
    }
    if (declaration->aspectRatio.aspectRatio > 0) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .aspectRatioElementConfig = Clay__StoreAspectRatioElementConfig(declaration->aspectRatio) }, CLAY__ELEMENT_CONFIG_TYPE_ASPECT);
        Clay__int32_tArray_Add(&context->aspectRatioElementIndexes, context->layoutElements.length - 1);
    }
    if (declaration->floating.attachTo != CLAY_ATTACH_TO_NONE) {
        Clay_FloatingElementConfig floatingConfig = declaration->floating;
        // This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here
        Clay_LayoutElement *hierarchicalParent = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 2));
        if (hierarchicalParent) {
            uint32_t clipElementId = 0;
            if (declaration->floating.attachTo == CLAY_ATTACH_TO_PARENT) {
                // Attach to the element's direct hierarchical parent
                floatingConfig.parentId = hierarchicalParent->id;
                if (context->openClipElementStack.length > 0) {
                    clipElementId = Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1);
                }
            } else if (declaration->floating.attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID) {
                Clay_LayoutElementHashMapItem *parentItem = Clay__GetHashMapItem(floatingConfig.parentId);
                if (parentItem == &Clay_LayoutElementHashMapItem_DEFAULT) {
                    context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                            .errorType = CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND,
                            .errorText = CLAY_STRING("A floating element was declared with a parentId, but no element with that ID was found."),
                            .userData = context->errorHandler.userData });
                } else {
                    clipElementId = Clay__int32_tArray_GetValue(&context->layoutElementClipElementIds, (int32_t)(parentItem->layoutElement - context->layoutElements.internalArray));
                }
            } else if (declaration->floating.attachTo == CLAY_ATTACH_TO_ROOT) {
                floatingConfig.parentId = Clay__HashString(CLAY_STRING("Clay__RootContainer"), 0).id;
            }
            if (!openLayoutElementId.id) {
                openLayoutElementId = Clay__HashStringWithOffset(CLAY_STRING("Clay__FloatingContainer"), context->layoutElementTreeRoots.length, 0);
            }
            if (declaration->floating.clipTo == CLAY_CLIP_TO_NONE) {
                clipElementId = 0;
            }
            int32_t currentElementIndex = Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1);
            Clay__int32_tArray_Set(&context->layoutElementClipElementIds, currentElementIndex, clipElementId);
            Clay__int32_tArray_Add(&context->openClipElementStack, clipElementId);
            Clay__LayoutElementTreeRootArray_Add(&context->layoutElementTreeRoots, CLAY__INIT(Clay__LayoutElementTreeRoot) {
                    .layoutElementIndex = Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1),
                    .parentId = floatingConfig.parentId,
                    .clipElementId = clipElementId,
                    .zIndex = floatingConfig.zIndex,
            });
            Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .floatingElementConfig = Clay__StoreFloatingElementConfig(floatingConfig) }, CLAY__ELEMENT_CONFIG_TYPE_FLOATING);
        }
    }
    if (declaration->custom.customData) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .customElementConfig = Clay__StoreCustomElementConfig(declaration->custom) }, CLAY__ELEMENT_CONFIG_TYPE_CUSTOM);
    }

    if (openLayoutElementId.id != 0) {
        Clay__AttachId(openLayoutElementId);
    } else if (openLayoutElement->id == 0) {
        openLayoutElementId = Clay__GenerateIdForAnonymousElement(openLayoutElement);
    }

    if (declaration->clip.horizontal | declaration->clip.vertical) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .clipElementConfig = Clay__StoreClipElementConfig(declaration->clip) }, CLAY__ELEMENT_CONFIG_TYPE_CLIP);
        Clay__int32_tArray_Add(&context->openClipElementStack, (int)openLayoutElement->id);
        // Retrieve or create cached data to track scroll position across frames
        Clay__ScrollContainerDataInternal *scrollOffset = CLAY__NULL;
        for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
            Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
            if (openLayoutElement->id == mapping->elementId) {
                scrollOffset = mapping;
                scrollOffset->layoutElement = openLayoutElement;
                scrollOffset->openThisFrame = true;
            }
        }
        if (!scrollOffset) {
            scrollOffset = Clay__ScrollContainerDataInternalArray_Add(&context->scrollContainerDatas, CLAY__INIT(Clay__ScrollContainerDataInternal){.layoutElement = openLayoutElement, .scrollOrigin = {-1,-1}, .elementId = openLayoutElement->id, .openThisFrame = true});
        }
        if (context->externalScrollHandlingEnabled) {
            scrollOffset->scrollPosition = Clay__QueryScrollOffset(scrollOffset->elementId, context->queryScrollOffsetUserData);
        }
    }
    if (!Clay__MemCmp((char *)(&declaration->border.width), (char *)(&Clay__BorderWidth_DEFAULT), sizeof(Clay_BorderWidth))) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .borderElementConfig = Clay__StoreBorderElementConfig(declaration->border) }, CLAY__ELEMENT_CONFIG_TYPE_BORDER);
    }
}

void Clay__ConfigureOpenElement(const Clay_ElementDeclaration declaration) {
    Clay__ConfigureOpenElementPtr(&declaration);
}

void Clay__InitializeEphemeralMemory(Clay_Context* context) {
    int32_t maxElementCount = context->maxElementCount;
    // Ephemeral Memory - reset every frame
    Clay_Arena *arena = &context->internalArena;
    arena->nextAllocation = context->arenaResetOffset;

    context->layoutElementChildrenBuffer = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElements = Clay_LayoutElementArray_Allocate_Arena(maxElementCount, arena);
    context->warnings = Clay__WarningArray_Allocate_Arena(100, arena);

    context->layoutConfigs = Clay__LayoutConfigArray_Allocate_Arena(maxElementCount, arena);
    context->elementConfigs = Clay__ElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->textElementConfigs = Clay__TextElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->aspectRatioElementConfigs = Clay__AspectRatioElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->imageElementConfigs = Clay__ImageElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->floatingElementConfigs = Clay__FloatingElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->clipElementConfigs = Clay__ClipElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->customElementConfigs = Clay__CustomElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->borderElementConfigs = Clay__BorderElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->sharedElementConfigs = Clay__SharedElementConfigArray_Allocate_Arena(maxElementCount, arena);

    context->layoutElementIdStrings = Clay__StringArray_Allocate_Arena(maxElementCount, arena);
    context->wrappedTextLines = Clay__WrappedTextLineArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementTreeNodeArray1 = Clay__LayoutElementTreeNodeArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementTreeRoots = Clay__LayoutElementTreeRootArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementChildren = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->openLayoutElementStack = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->textElementData = Clay__TextElementDataArray_Allocate_Arena(maxElementCount, arena);
    context->aspectRatioElementIndexes = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->renderCommands = Clay_RenderCommandArray_Allocate_Arena(maxElementCount, arena);
    context->treeNodeVisited = Clay__boolArray_Allocate_Arena(maxElementCount, arena);
    context->treeNodeVisited.length = context->treeNodeVisited.capacity; // This array is accessed directly rather than behaving as a list
    context->openClipElementStack = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->reusableElementIndexBuffer = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementClipElementIds = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->dynamicStringData = Clay__charArray_Allocate_Arena(maxElementCount, arena);
}

void Clay__InitializePersistentMemory(Clay_Context* context) {
    // Persistent memory - initialized once and not reset
    int32_t maxElementCount = context->maxElementCount;
    int32_t maxMeasureTextCacheWordCount = context->maxMeasureTextCacheWordCount;
    Clay_Arena *arena = &context->internalArena;

    context->scrollContainerDatas = Clay__ScrollContainerDataInternalArray_Allocate_Arena(100, arena);
    context->layoutElementsHashMapInternal = Clay__LayoutElementHashMapItemArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementsHashMap = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->measureTextHashMapInternal = Clay__MeasureTextCacheItemArray_Allocate_Arena(maxElementCount, arena);
    context->measureTextHashMapInternalFreeList = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->measuredWordsFreeList = Clay__int32_tArray_Allocate_Arena(maxMeasureTextCacheWordCount, arena);
    context->measureTextHashMap = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->measuredWords = Clay__MeasuredWordArray_Allocate_Arena(maxMeasureTextCacheWordCount, arena);
    context->pointerOverIds = Clay_ElementIdArray_Allocate_Arena(maxElementCount, arena);
    context->debugElementData = Clay__DebugElementDataArray_Allocate_Arena(maxElementCount, arena);
    context->arenaResetOffset = arena->nextAllocation;
}

const float CLAY__EPSILON = 0.01;

bool Clay__FloatEqual(float left, float right) {
    float subtracted = left - right;
    return subtracted < CLAY__EPSILON && subtracted > -CLAY__EPSILON;
}

void Clay__SizeContainersAlongAxis(bool xAxis) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__int32_tArray bfsBuffer = context->layoutElementChildrenBuffer;
    Clay__int32_tArray resizableContainerBuffer = context->openLayoutElementStack;
    for (int32_t rootIndex = 0; rootIndex < context->layoutElementTreeRoots.length; ++rootIndex) {
        bfsBuffer.length = 0;
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)root->layoutElementIndex);
        Clay__int32_tArray_Add(&bfsBuffer, (int32_t)root->layoutElementIndex);

        // Size floating containers to their parents
        if (Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)) {
            Clay_FloatingElementConfig *floatingElementConfig = Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
            Clay_LayoutElementHashMapItem *parentItem = Clay__GetHashMapItem(floatingElementConfig->parentId);
            if (parentItem && parentItem != &Clay_LayoutElementHashMapItem_DEFAULT) {
                Clay_LayoutElement *parentLayoutElement = parentItem->layoutElement;
                switch (rootElement->layoutConfig->sizing.width.type) {
                    case CLAY__SIZING_TYPE_GROW: {
                        rootElement->dimensions.width = parentLayoutElement->dimensions.width;
                        break;
                    }
                    case CLAY__SIZING_TYPE_PERCENT: {
                        rootElement->dimensions.width = parentLayoutElement->dimensions.width * rootElement->layoutConfig->sizing.width.size.percent;
                        break;
                    }
                    default: break;
                }
                switch (rootElement->layoutConfig->sizing.height.type) {
                    case CLAY__SIZING_TYPE_GROW: {
                        rootElement->dimensions.height = parentLayoutElement->dimensions.height;
                        break;
                    }
                    case CLAY__SIZING_TYPE_PERCENT: {
                        rootElement->dimensions.height = parentLayoutElement->dimensions.height * rootElement->layoutConfig->sizing.height.size.percent;
                        break;
                    }
                    default: break;
                }
            }
        }

        if (rootElement->layoutConfig->sizing.width.type != CLAY__SIZING_TYPE_PERCENT) {
            rootElement->dimensions.width = CLAY__MIN(CLAY__MAX(rootElement->dimensions.width, rootElement->layoutConfig->sizing.width.size.minMax.min), rootElement->layoutConfig->sizing.width.size.minMax.max);
        }
        if (rootElement->layoutConfig->sizing.height.type != CLAY__SIZING_TYPE_PERCENT) {
            rootElement->dimensions.height = CLAY__MIN(CLAY__MAX(rootElement->dimensions.height, rootElement->layoutConfig->sizing.height.size.minMax.min), rootElement->layoutConfig->sizing.height.size.minMax.max);
        }

        for (int32_t i = 0; i < bfsBuffer.length; ++i) {
            int32_t parentIndex = Clay__int32_tArray_GetValue(&bfsBuffer, i);
            Clay_LayoutElement *parent = Clay_LayoutElementArray_Get(&context->layoutElements, parentIndex);
            Clay_LayoutConfig *parentStyleConfig = parent->layoutConfig;
            int32_t growContainerCount = 0;
            float parentSize = xAxis ? parent->dimensions.width : parent->dimensions.height;
            float parentPadding = (float)(xAxis ? (parent->layoutConfig->padding.left + parent->layoutConfig->padding.right) : (parent->layoutConfig->padding.top + parent->layoutConfig->padding.bottom));
            float innerContentSize = 0, totalPaddingAndChildGaps = parentPadding;
            bool sizingAlongAxis = (xAxis && parentStyleConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) || (!xAxis && parentStyleConfig->layoutDirection == CLAY_TOP_TO_BOTTOM);
            resizableContainerBuffer.length = 0;
            float parentChildGap = parentStyleConfig->childGap;

            for (int32_t childOffset = 0; childOffset < parent->childrenOrTextContent.children.length; childOffset++) {
                int32_t childElementIndex = parent->childrenOrTextContent.children.elements[childOffset];
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, childElementIndex);
                Clay_SizingAxis childSizing = xAxis ? childElement->layoutConfig->sizing.width : childElement->layoutConfig->sizing.height;
                float childSize = xAxis ? childElement->dimensions.width : childElement->dimensions.height;

                if (!Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && childElement->childrenOrTextContent.children.length > 0) {
                    Clay__int32_tArray_Add(&bfsBuffer, childElementIndex);
                }

                if (childSizing.type != CLAY__SIZING_TYPE_PERCENT
                    && childSizing.type != CLAY__SIZING_TYPE_FIXED
                    && (!Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (Clay__FindElementConfigWithType(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig->wrapMode == CLAY_TEXT_WRAP_WORDS)) // todo too many loops
//                    && (xAxis || !Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT))
                ) {
                    Clay__int32_tArray_Add(&resizableContainerBuffer, childElementIndex);
                }

                if (sizingAlongAxis) {
                    innerContentSize += (childSizing.type == CLAY__SIZING_TYPE_PERCENT ? 0 : childSize);
                    if (childSizing.type == CLAY__SIZING_TYPE_GROW) {
                        growContainerCount++;
                    }
                    if (childOffset > 0) {
                        innerContentSize += parentChildGap; // For children after index 0, the childAxisOffset is the gap from the previous child
                        totalPaddingAndChildGaps += parentChildGap;
                    }
                } else {
                    innerContentSize = CLAY__MAX(childSize, innerContentSize);
                }
            }

            // Expand percentage containers to size
            for (int32_t childOffset = 0; childOffset < parent->childrenOrTextContent.children.length; childOffset++) {
                int32_t childElementIndex = parent->childrenOrTextContent.children.elements[childOffset];
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, childElementIndex);
                Clay_SizingAxis childSizing = xAxis ? childElement->layoutConfig->sizing.width : childElement->layoutConfig->sizing.height;
                float *childSize = xAxis ? &childElement->dimensions.width : &childElement->dimensions.height;
                if (childSizing.type == CLAY__SIZING_TYPE_PERCENT) {
                    *childSize = (parentSize - totalPaddingAndChildGaps) * childSizing.size.percent;
                    if (sizingAlongAxis) {
                        innerContentSize += *childSize;
                    }
                    Clay__UpdateAspectRatioBox(childElement);
                }
            }

            if (sizingAlongAxis) {
                float sizeToDistribute = parentSize - parentPadding - innerContentSize;
                // The content is too large, compress the children as much as possible
                if (sizeToDistribute < 0) {
                    // If the parent clips content in this axis direction, don't compress children, just leave them alone
                    Clay_ClipElementConfig *clipElementConfig = Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                    if (clipElementConfig) {
                        if (((xAxis && clipElementConfig->horizontal) || (!xAxis && clipElementConfig->vertical))) {
                            continue;
                        }
                    }
                    // Scrolling containers preferentially compress before others
                    while (sizeToDistribute < -CLAY__EPSILON && resizableContainerBuffer.length > 0) {
                        float largest = 0;
                        float secondLargest = 0;
                        float widthToAdd = sizeToDistribute;
                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float childSize = xAxis ? child->dimensions.width : child->dimensions.height;
                            if (Clay__FloatEqual(childSize, largest)) { continue; }
                            if (childSize > largest) {
                                secondLargest = largest;
                                largest = childSize;
                            }
                            if (childSize < largest) {
                                secondLargest = CLAY__MAX(secondLargest, childSize);
                                widthToAdd = secondLargest - largest;
                            }
                        }

                        widthToAdd = CLAY__MAX(widthToAdd, sizeToDistribute / resizableContainerBuffer.length);

                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float *childSize = xAxis ? &child->dimensions.width : &child->dimensions.height;
                            float minSize = xAxis ? child->minDimensions.width : child->minDimensions.height;
                            float previousWidth = *childSize;
                            if (Clay__FloatEqual(*childSize, largest)) {
                                *childSize += widthToAdd;
                                if (*childSize <= minSize) {
                                    *childSize = minSize;
                                    Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
                                }
                                sizeToDistribute -= (*childSize - previousWidth);
                            }
                        }
                    }
                // The content is too small, allow SIZING_GROW containers to expand
                } else if (sizeToDistribute > 0 && growContainerCount > 0) {
                    for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                        Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                        Clay__SizingType childSizing = xAxis ? child->layoutConfig->sizing.width.type : child->layoutConfig->sizing.height.type;
                        if (childSizing != CLAY__SIZING_TYPE_GROW) {
                            Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
                        }
                    }
                    while (sizeToDistribute > CLAY__EPSILON && resizableContainerBuffer.length > 0) {
                        float smallest = CLAY__MAXFLOAT;
                        float secondSmallest = CLAY__MAXFLOAT;
                        float widthToAdd = sizeToDistribute;
                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float childSize = xAxis ? child->dimensions.width : child->dimensions.height;
                            if (Clay__FloatEqual(childSize, smallest)) { continue; }
                            if (childSize < smallest) {
                                secondSmallest = smallest;
                                smallest = childSize;
                            }
                            if (childSize > smallest) {
                                secondSmallest = CLAY__MIN(secondSmallest, childSize);
                                widthToAdd = secondSmallest - smallest;
                            }
                        }

                        widthToAdd = CLAY__MIN(widthToAdd, sizeToDistribute / resizableContainerBuffer.length);

                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float *childSize = xAxis ? &child->dimensions.width : &child->dimensions.height;
                            float maxSize = xAxis ? child->layoutConfig->sizing.width.size.minMax.max : child->layoutConfig->sizing.height.size.minMax.max;
                            float previousWidth = *childSize;
                            if (Clay__FloatEqual(*childSize, smallest)) {
                                *childSize += widthToAdd;
                                if (*childSize >= maxSize) {
                                    *childSize = maxSize;
                                    Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
                                }
                                sizeToDistribute -= (*childSize - previousWidth);
                            }
                        }
                    }
                }
            // Sizing along the non layout axis ("off axis")
            } else {
                for (int32_t childOffset = 0; childOffset < resizableContainerBuffer.length; childOffset++) {
                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childOffset));
                    Clay_SizingAxis childSizing = xAxis ? childElement->layoutConfig->sizing.width : childElement->layoutConfig->sizing.height;
                    float minSize = xAxis ? childElement->minDimensions.width : childElement->minDimensions.height;
                    float *childSize = xAxis ? &childElement->dimensions.width : &childElement->dimensions.height;

                    float maxSize = parentSize - parentPadding;
                    // If we're laying out the children of a scroll panel, grow containers expand to the size of the inner content, not the outer container
                    if (Clay__ElementHasConfig(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP)) {
                        Clay_ClipElementConfig *clipElementConfig = Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                        if (((xAxis && clipElementConfig->horizontal) || (!xAxis && clipElementConfig->vertical))) {
                            maxSize = CLAY__MAX(maxSize, innerContentSize);
                        }
                    }
                    if (childSizing.type == CLAY__SIZING_TYPE_GROW) {
                        *childSize = CLAY__MIN(maxSize, childSizing.size.minMax.max);
                    }
                    *childSize = CLAY__MAX(minSize, CLAY__MIN(*childSize, maxSize));
                }
            }
        }
    }
}

Clay_String Clay__IntToString(int32_t integer) {
    if (integer == 0) {
        return CLAY__INIT(Clay_String) { .length = 1, .chars = "0" };
    }
    Clay_Context* context = Clay_GetCurrentContext();
    char *chars = (char *)(context->dynamicStringData.internalArray + context->dynamicStringData.length);
    int32_t length = 0;
    int32_t sign = integer;

    if (integer < 0) {
        integer = -integer;
    }
    while (integer > 0) {
        chars[length++] = (char)(integer % 10 + '0');
        integer /= 10;
    }

    if (sign < 0) {
        chars[length++] = '-';
    }

    // Reverse the string to get the correct order
    for (int32_t j = 0, k = length - 1; j < k; j++, k--) {
        char temp = chars[j];
        chars[j] = chars[k];
        chars[k] = temp;
    }
    context->dynamicStringData.length += length;
    return CLAY__INIT(Clay_String) { .length = length, .chars = chars };
}

void Clay__AddRenderCommand(Clay_RenderCommand renderCommand) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->renderCommands.length < context->renderCommands.capacity - 1) {
        Clay_RenderCommandArray_Add(&context->renderCommands, renderCommand);
    } else {
        if (!context->booleanWarnings.maxRenderCommandsExceeded) {
            context->booleanWarnings.maxRenderCommandsExceeded = true;
            context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                .errorType = CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
                .errorText = CLAY_STRING("Clay ran out of capacity while attempting to create render commands. This is usually caused by a large amount of wrapping text elements while close to the max element capacity. Try using Clay_SetMaxElementCount() with a higher value."),
                .userData = context->errorHandler.userData });
        }
    }
}

bool Clay__ElementIsOffscreen(Clay_BoundingBox *boundingBox) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->disableCulling) {
        return false;
    }

    return (boundingBox->x > (float)context->layoutDimensions.width) ||
           (boundingBox->y > (float)context->layoutDimensions.height) ||
           (boundingBox->x + boundingBox->width < 0) ||
           (boundingBox->y + boundingBox->height < 0);
}

void Clay__CalculateFinalLayout(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    // Calculate sizing along the X axis
    Clay__SizeContainersAlongAxis(true);

    // Wrap text
    for (int32_t textElementIndex = 0; textElementIndex < context->textElementData.length; ++textElementIndex) {
        Clay__TextElementData *textElementData = Clay__TextElementDataArray_Get(&context->textElementData, textElementIndex);
        textElementData->wrappedLines = CLAY__INIT(Clay__WrappedTextLineArraySlice) { .length = 0, .internalArray = &context->wrappedTextLines.internalArray[context->wrappedTextLines.length] };
        Clay_LayoutElement *containerElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)textElementData->elementIndex);
        Clay_TextElementConfig *textConfig = Clay__FindElementConfigWithType(containerElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig;
        Clay__MeasureTextCacheItem *measureTextCacheItem = Clay__MeasureTextCached(&textElementData->text, textConfig);
        float lineWidth = 0;
        float lineHeight = textConfig->lineHeight > 0 ? (float)textConfig->lineHeight : textElementData->preferredDimensions.height;
        int32_t lineLengthChars = 0;
        int32_t lineStartOffset = 0;
        if (!measureTextCacheItem->containsNewlines && textElementData->preferredDimensions.width <= containerElement->dimensions.width) {
            Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { containerElement->dimensions,  textElementData->text });
            textElementData->wrappedLines.length++;
            continue;
        }
        float spaceWidth = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = 1, .chars = CLAY__SPACECHAR.chars, .baseChars = CLAY__SPACECHAR.chars }, textConfig, context->measureTextUserData).width;
        int32_t wordIndex = measureTextCacheItem->measuredWordsStartIndex;
        while (wordIndex != -1) {
            if (context->wrappedTextLines.length > context->wrappedTextLines.capacity - 1) {
                break;
            }
            Clay__MeasuredWord *measuredWord = Clay__MeasuredWordArray_Get(&context->measuredWords, wordIndex);
            // Only word on the line is too large, just render it anyway
            if (lineLengthChars == 0 && lineWidth + measuredWord->width > containerElement->dimensions.width) {
                Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { measuredWord->width, lineHeight }, { .length = measuredWord->length, .chars = &textElementData->text.chars[measuredWord->startOffset] } });
                textElementData->wrappedLines.length++;
                wordIndex = measuredWord->next;
                lineStartOffset = measuredWord->startOffset + measuredWord->length;
            }
            // measuredWord->length == 0 means a newline character
            else if (measuredWord->length == 0 || lineWidth + measuredWord->width > containerElement->dimensions.width) {
                // Wrapped text lines list has overflowed, just render out the line
                bool finalCharIsSpace = textElementData->text.chars[CLAY__MAX(lineStartOffset + lineLengthChars - 1, 0)] == ' ';
                Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { lineWidth + (finalCharIsSpace ? -spaceWidth : 0), lineHeight }, { .length = lineLengthChars + (finalCharIsSpace ? -1 : 0), .chars = &textElementData->text.chars[lineStartOffset] } });
                textElementData->wrappedLines.length++;
                if (lineLengthChars == 0 || measuredWord->length == 0) {
                    wordIndex = measuredWord->next;
                }
                lineWidth = 0;
                lineLengthChars = 0;
                lineStartOffset = measuredWord->startOffset;
            } else {
                lineWidth += measuredWord->width + textConfig->letterSpacing;
                lineLengthChars += measuredWord->length;
                wordIndex = measuredWord->next;
            }
        }
        if (lineLengthChars > 0) {
            Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { lineWidth - textConfig->letterSpacing, lineHeight }, {.length = lineLengthChars, .chars = &textElementData->text.chars[lineStartOffset] } });
            textElementData->wrappedLines.length++;
        }
        containerElement->dimensions.height = lineHeight * (float)textElementData->wrappedLines.length;
    }

    // Scale vertical heights according to aspect ratio
    for (int32_t i = 0; i < context->aspectRatioElementIndexes.length; ++i) {
        Clay_LayoutElement* aspectElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->aspectRatioElementIndexes, i));
        Clay_AspectRatioElementConfig *config = Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig;
        aspectElement->dimensions.height = (1 / config->aspectRatio) * aspectElement->dimensions.width;
        aspectElement->layoutConfig->sizing.height.size.minMax.max = aspectElement->dimensions.height;
    }

    // Propagate effect of text wrapping, aspect scaling etc. on height of parents
    Clay__LayoutElementTreeNodeArray dfsBuffer = context->layoutElementTreeNodeArray1;
    dfsBuffer.length = 0;
    for (int32_t i = 0; i < context->layoutElementTreeRoots.length; ++i) {
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i);
        context->treeNodeVisited.internalArray[dfsBuffer.length] = false;
        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)root->layoutElementIndex) });
    }
    while (dfsBuffer.length > 0) {
        Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
        Clay_LayoutElement *currentElement = currentElementTreeNode->layoutElement;
        if (!context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
            context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
            // If the element has no children or is the container for a text element, don't bother inspecting it
            if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement->childrenOrTextContent.children.length == 0) {
                dfsBuffer.length--;
                continue;
            }
            // Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
            for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; i++) {
                context->treeNodeVisited.internalArray[dfsBuffer.length] = false;
                Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]) });
            }
            continue;
        }
        dfsBuffer.length--;

        // DFS node has been visited, this is on the way back up to the root
        Clay_LayoutConfig *layoutConfig = currentElement->layoutConfig;
        if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
            // Resize any parent containers that have grown in height along their non layout axis
            for (int32_t j = 0; j < currentElement->childrenOrTextContent.children.length; ++j) {
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[j]);
                float childHeightWithPadding = CLAY__MAX(childElement->dimensions.height + layoutConfig->padding.top + layoutConfig->padding.bottom, currentElement->dimensions.height);
                currentElement->dimensions.height = CLAY__MIN(CLAY__MAX(childHeightWithPadding, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
            }
        } else if (layoutConfig->layoutDirection == CLAY_TOP_TO_BOTTOM) {
            // Resizing along the layout axis
            float contentHeight = (float)(layoutConfig->padding.top + layoutConfig->padding.bottom);
            for (int32_t j = 0; j < currentElement->childrenOrTextContent.children.length; ++j) {
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[j]);
                contentHeight += childElement->dimensions.height;
            }
            contentHeight += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
            currentElement->dimensions.height = CLAY__MIN(CLAY__MAX(contentHeight, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
        }
    }

    // Calculate sizing along the Y axis
    Clay__SizeContainersAlongAxis(false);

    // Scale horizontal widths according to aspect ratio
    for (int32_t i = 0; i < context->aspectRatioElementIndexes.length; ++i) {
        Clay_LayoutElement* aspectElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->aspectRatioElementIndexes, i));
        Clay_AspectRatioElementConfig *config = Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig;
        aspectElement->dimensions.width = config->aspectRatio * aspectElement->dimensions.height;
    }

    // Sort tree roots by z-index
    int32_t sortMax = context->layoutElementTreeRoots.length - 1;
    while (sortMax > 0) { // todo dumb bubble sort
        for (int32_t i = 0; i < sortMax; ++i) {
            Clay__LayoutElementTreeRoot current = *Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i);
            Clay__LayoutElementTreeRoot next = *Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i + 1);
            if (next.zIndex < current.zIndex) {
                Clay__LayoutElementTreeRootArray_Set(&context->layoutElementTreeRoots, i, next);
                Clay__LayoutElementTreeRootArray_Set(&context->layoutElementTreeRoots, i + 1, current);
            }
        }
        sortMax--;
    }

    // Calculate final positions and generate render commands
    context->renderCommands.length = 0;
    dfsBuffer.length = 0;
    for (int32_t rootIndex = 0; rootIndex < context->layoutElementTreeRoots.length; ++rootIndex) {
        dfsBuffer.length = 0;
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)root->layoutElementIndex);
        Clay_Vector2 rootPosition = CLAY__DEFAULT_STRUCT;
        Clay_LayoutElementHashMapItem *parentHashMapItem = Clay__GetHashMapItem(root->parentId);
        // Position root floating containers
        if (Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) && parentHashMapItem) {
            Clay_FloatingElementConfig *config = Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
            Clay_Dimensions rootDimensions = rootElement->dimensions;
            Clay_BoundingBox parentBoundingBox = parentHashMapItem->boundingBox;
            // Set X position
            Clay_Vector2 targetAttachPosition = CLAY__DEFAULT_STRUCT;
            switch (config->attachPoints.parent) {
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_LEFT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x; break;
                case CLAY_ATTACH_POINT_CENTER_TOP:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + (parentBoundingBox.width / 2); break;
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + parentBoundingBox.width; break;
            }
            switch (config->attachPoints.element) {
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_LEFT_BOTTOM: break;
                case CLAY_ATTACH_POINT_CENTER_TOP:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x -= (rootDimensions.width / 2); break;
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x -= rootDimensions.width; break;
            }
            switch (config->attachPoints.parent) { // I know I could merge the x and y switch statements, but this is easier to read
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_CENTER_TOP: targetAttachPosition.y = parentBoundingBox.y; break;
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y = parentBoundingBox.y + (parentBoundingBox.height / 2); break;
                case CLAY_ATTACH_POINT_LEFT_BOTTOM:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y = parentBoundingBox.y + parentBoundingBox.height; break;
            }
            switch (config->attachPoints.element) {
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_CENTER_TOP: break;
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y -= (rootDimensions.height / 2); break;
                case CLAY_ATTACH_POINT_LEFT_BOTTOM:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y -= rootDimensions.height; break;
            }
            targetAttachPosition.x += config->offset.x;
            targetAttachPosition.y += config->offset.y;
            rootPosition = targetAttachPosition;
        }
        if (root->clipElementId) {
            Clay_LayoutElementHashMapItem *clipHashMapItem = Clay__GetHashMapItem(root->clipElementId);
            if (clipHashMapItem) {
                // Floating elements that are attached to scrolling contents won't be correctly positioned if external scroll handling is enabled, fix here
                if (context->externalScrollHandlingEnabled) {
                    Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(clipHashMapItem->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                    if (clipConfig->horizontal) {
                        rootPosition.x += clipConfig->childOffset.x;
                    }
                    if (clipConfig->vertical) {
                        rootPosition.y += clipConfig->childOffset.y;
                    }
                }
                Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                    .boundingBox = clipHashMapItem->boundingBox,
                    .userData = 0,
                    .id = Clay__HashNumber(rootElement->id, rootElement->childrenOrTextContent.children.length + 10).id, // TODO need a better strategy for managing derived ids
                    .zIndex = root->zIndex,
                    .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START,
                });
            }
        }
        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = rootElement, .position = rootPosition, .nextChildOffset = { .x = (float)rootElement->layoutConfig->padding.left, .y = (float)rootElement->layoutConfig->padding.top } });

        context->treeNodeVisited.internalArray[0] = false;
        while (dfsBuffer.length > 0) {
            Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
            Clay_LayoutElement *currentElement = currentElementTreeNode->layoutElement;
            Clay_LayoutConfig *layoutConfig = currentElement->layoutConfig;
            Clay_Vector2 scrollOffset = CLAY__DEFAULT_STRUCT;

            // This will only be run a single time for each element in downwards DFS order
            if (!context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
                context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;

                Clay_BoundingBox currentElementBoundingBox = { currentElementTreeNode->position.x, currentElementTreeNode->position.y, currentElement->dimensions.width, currentElement->dimensions.height };
                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)) {
                    Clay_FloatingElementConfig *floatingElementConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
                    Clay_Dimensions expand = floatingElementConfig->expand;
                    currentElementBoundingBox.x -= expand.width;
                    currentElementBoundingBox.width += expand.width * 2;
                    currentElementBoundingBox.y -= expand.height;
                    currentElementBoundingBox.height += expand.height * 2;
                }

                Clay__ScrollContainerDataInternal *scrollContainerData = CLAY__NULL;
                // Apply scroll offsets to container
                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP)) {
                    Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;

                    // This linear scan could theoretically be slow under very strange conditions, but I can't imagine a real UI with more than a few 10's of scroll containers
                    for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
                        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
                        if (mapping->layoutElement == currentElement) {
                            scrollContainerData = mapping;
                            mapping->boundingBox = currentElementBoundingBox;
                            scrollOffset = clipConfig->childOffset;
                            if (context->externalScrollHandlingEnabled) {
                                scrollOffset = CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
                            }
                            break;
                        }
                    }
                }

                Clay_LayoutElementHashMapItem *hashMapItem = Clay__GetHashMapItem(currentElement->id);
                if (hashMapItem) {
                    hashMapItem->boundingBox = currentElementBoundingBox;
                    if (hashMapItem->idAlias) {
                        Clay_LayoutElementHashMapItem *hashMapItemAlias = Clay__GetHashMapItem(hashMapItem->idAlias);
                        if (hashMapItemAlias) {
                            hashMapItemAlias->boundingBox = currentElementBoundingBox;
                        }
                    }
                }

                int32_t sortedConfigIndexes[20];
                for (int32_t elementConfigIndex = 0; elementConfigIndex < currentElement->elementConfigs.length; ++elementConfigIndex) {
                    sortedConfigIndexes[elementConfigIndex] = elementConfigIndex;
                }
                sortMax = currentElement->elementConfigs.length - 1;
                while (sortMax > 0) { // todo dumb bubble sort
                    for (int32_t i = 0; i < sortMax; ++i) {
                        int32_t current = sortedConfigIndexes[i];
                        int32_t next = sortedConfigIndexes[i + 1];
                        Clay__ElementConfigType currentType = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, current)->type;
                        Clay__ElementConfigType nextType = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, next)->type;
                        if (nextType == CLAY__ELEMENT_CONFIG_TYPE_CLIP || currentType == CLAY__ELEMENT_CONFIG_TYPE_BORDER) {
                            sortedConfigIndexes[i] = next;
                            sortedConfigIndexes[i + 1] = current;
                        }
                    }
                    sortMax--;
                }

                bool emitRectangle = false;
                // Create the render commands for this element
                Clay_SharedElementConfig *sharedConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig;
                if (sharedConfig && sharedConfig->backgroundColor.a > 0) {
                   emitRectangle = true;
                }
                else if (!sharedConfig) {
                    emitRectangle = false;
                    sharedConfig = &Clay_SharedElementConfig_DEFAULT;
                }
                for (int32_t elementConfigIndex = 0; elementConfigIndex < currentElement->elementConfigs.length; ++elementConfigIndex) {
                    Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, sortedConfigIndexes[elementConfigIndex]);
                    Clay_RenderCommand renderCommand = {
                        .boundingBox = currentElementBoundingBox,
                        .userData = sharedConfig->userData,
                        .id = currentElement->id,
                    };

                    bool offscreen = Clay__ElementIsOffscreen(&currentElementBoundingBox);
                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
                    bool shouldRender = !offscreen;
                    switch (elementConfig->type) {
                        case CLAY__ELEMENT_CONFIG_TYPE_ASPECT:
                        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING:
                        case CLAY__ELEMENT_CONFIG_TYPE_SHARED:
                        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: {
                            shouldRender = false;
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_CLIP: {
                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START;
                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
                                .clip = {
                                    .horizontal = elementConfig->config.clipElementConfig->horizontal,
                                    .vertical = elementConfig->config.clipElementConfig->vertical,
                                }
                            };
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: {
                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_IMAGE;
                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
                                .image = {
                                    .backgroundColor = sharedConfig->backgroundColor,
                                    .cornerRadius = sharedConfig->cornerRadius,
                                    .imageData = elementConfig->config.imageElementConfig->imageData,
                               }
                            };
                            emitRectangle = false;
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: {
                            if (!shouldRender) {
                                break;
                            }
                            shouldRender = false;
                            Clay_ElementConfigUnion configUnion = elementConfig->config;
                            Clay_TextElementConfig *textElementConfig = configUnion.textElementConfig;
                            float naturalLineHeight = currentElement->childrenOrTextContent.textElementData->preferredDimensions.height;
                            float finalLineHeight = textElementConfig->lineHeight > 0 ? (float)textElementConfig->lineHeight : naturalLineHeight;
                            float lineHeightOffset = (finalLineHeight - naturalLineHeight) / 2;
                            float yPosition = lineHeightOffset;
                            for (int32_t lineIndex = 0; lineIndex < currentElement->childrenOrTextContent.textElementData->wrappedLines.length; ++lineIndex) {
                                Clay__WrappedTextLine *wrappedLine = Clay__WrappedTextLineArraySlice_Get(&currentElement->childrenOrTextContent.textElementData->wrappedLines, lineIndex);
                                if (wrappedLine->line.length == 0) {
                                    yPosition += finalLineHeight;
                                    continue;
                                }
                                float offset = (currentElementBoundingBox.width - wrappedLine->dimensions.width);
                                if (textElementConfig->textAlignment == CLAY_TEXT_ALIGN_LEFT) {
                                    offset = 0;
                                }
                                if (textElementConfig->textAlignment == CLAY_TEXT_ALIGN_CENTER) {
                                    offset /= 2;
                                }
                                Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                                    .boundingBox = { currentElementBoundingBox.x + offset, currentElementBoundingBox.y + yPosition, wrappedLine->dimensions.width, wrappedLine->dimensions.height },
                                    .renderData = { .text = {
                                        .stringContents = CLAY__INIT(Clay_StringSlice) { .length = wrappedLine->line.length, .chars = wrappedLine->line.chars, .baseChars = currentElement->childrenOrTextContent.textElementData->text.chars },
                                        .textColor = textElementConfig->textColor,
                                        .fontId = textElementConfig->fontId,
                                        .fontSize = textElementConfig->fontSize,
                                        .letterSpacing = textElementConfig->letterSpacing,
                                        .lineHeight = textElementConfig->lineHeight,
                                    }},
                                    .userData = textElementConfig->userData,
                                    .id = Clay__HashNumber(lineIndex, currentElement->id).id,
                                    .zIndex = root->zIndex,
                                    .commandType = CLAY_RENDER_COMMAND_TYPE_TEXT,
                                });
                                yPosition += finalLineHeight;

                                if (!context->disableCulling && (currentElementBoundingBox.y + yPosition > context->layoutDimensions.height)) {
                                    break;
                                }
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: {
                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_CUSTOM;
                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
                                .custom = {
                                    .backgroundColor = sharedConfig->backgroundColor,
                                    .cornerRadius = sharedConfig->cornerRadius,
                                    .customData = elementConfig->config.customElementConfig->customData,
                                }
                            };
                            emitRectangle = false;
                            break;
                        }
                        default: break;
                    }
                    if (shouldRender) {
                        Clay__AddRenderCommand(renderCommand);
                    }
                    if (offscreen) {
                        // NOTE: You may be tempted to try an early return / continue if an element is off screen. Why bother calculating layout for its children, right?
                        // Unfortunately, a FLOATING_CONTAINER may be defined that attaches to a child or grandchild of this element, which is large enough to still
                        // be on screen, even if this element isn't. That depends on this element and it's children being laid out correctly (even if they are entirely off screen)
                    }
                }

                if (emitRectangle) {
                    Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                        .boundingBox = currentElementBoundingBox,
                        .renderData = { .rectangle = {
                                .backgroundColor = sharedConfig->backgroundColor,
                                .cornerRadius = sharedConfig->cornerRadius,
                        }},
                        .userData = sharedConfig->userData,
                        .id = currentElement->id,
                        .zIndex = root->zIndex,
                        .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                    });
                }

                // Setup initial on-axis alignment
                if (!Clay__ElementHasConfig(currentElementTreeNode->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
                    Clay_Dimensions contentSize = {0,0};
                    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                        for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                            Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                            contentSize.width += childElement->dimensions.width;
                            contentSize.height = CLAY__MAX(contentSize.height, childElement->dimensions.height);
                        }
                        contentSize.width += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
                        float extraSpace = currentElement->dimensions.width - (float)(layoutConfig->padding.left + layoutConfig->padding.right) - contentSize.width;
                        switch (layoutConfig->childAlignment.x) {
                            case CLAY_ALIGN_X_LEFT: extraSpace = 0; break;
                            case CLAY_ALIGN_X_CENTER: extraSpace /= 2; break;
                            default: break;
                        }
                        currentElementTreeNode->nextChildOffset.x += extraSpace;
                        extraSpace = CLAY__MAX(0, extraSpace);
                    } else {
                        for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                            Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                            contentSize.width = CLAY__MAX(contentSize.width, childElement->dimensions.width);
                            contentSize.height += childElement->dimensions.height;
                        }
                        contentSize.height += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
                        float extraSpace = currentElement->dimensions.height - (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) - contentSize.height;
                        switch (layoutConfig->childAlignment.y) {
                            case CLAY_ALIGN_Y_TOP: extraSpace = 0; break;
                            case CLAY_ALIGN_Y_CENTER: extraSpace /= 2; break;
                            default: break;
                        }
                        extraSpace = CLAY__MAX(0, extraSpace);
                        currentElementTreeNode->nextChildOffset.y += extraSpace;
                    }

                    if (scrollContainerData) {
                        scrollContainerData->contentSize = CLAY__INIT(Clay_Dimensions) { contentSize.width + (float)(layoutConfig->padding.left + layoutConfig->padding.right), contentSize.height + (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) };
                    }
                }
            }
            else {
                // DFS is returning upwards backwards
                bool closeClipElement = false;
                Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                if (clipConfig) {
                    closeClipElement = true;
                    for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
                        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
                        if (mapping->layoutElement == currentElement) {
                            scrollOffset = clipConfig->childOffset;
                            if (context->externalScrollHandlingEnabled) {
                                scrollOffset = CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
                            }
                            break;
                        }
                    }
                }

                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER)) {
                    Clay_LayoutElementHashMapItem *currentElementData = Clay__GetHashMapItem(currentElement->id);
                    Clay_BoundingBox currentElementBoundingBox = currentElementData->boundingBox;

                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
                    if (!Clay__ElementIsOffscreen(&currentElementBoundingBox)) {
                        Clay_SharedElementConfig *sharedConfig = Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED) ? Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig : &Clay_SharedElementConfig_DEFAULT;
                        Clay_BorderElementConfig *borderConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER).borderElementConfig;
                        Clay_RenderCommand renderCommand = {
                                .boundingBox = currentElementBoundingBox,
                                .renderData = { .border = {
                                    .color = borderConfig->color,
                                    .cornerRadius = sharedConfig->cornerRadius,
                                    .width = borderConfig->width
                                }},
                                .userData = sharedConfig->userData,
                                .id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length).id,
                                .commandType = CLAY_RENDER_COMMAND_TYPE_BORDER,
                        };
                        Clay__AddRenderCommand(renderCommand);
                        if (borderConfig->width.betweenChildren > 0 && borderConfig->color.a > 0) {
                            float halfGap = layoutConfig->childGap / 2;
                            Clay_Vector2 borderOffset = { (float)layoutConfig->padding.left - halfGap, (float)layoutConfig->padding.top - halfGap };
                            if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                                for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                                    if (i > 0) {
                                        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                                            .boundingBox = { currentElementBoundingBox.x + borderOffset.x + scrollOffset.x, currentElementBoundingBox.y + scrollOffset.y, (float)borderConfig->width.betweenChildren, currentElement->dimensions.height },
                                            .renderData = { .rectangle = {
                                                .backgroundColor = borderConfig->color,
                                            } },
                                            .userData = sharedConfig->userData,
                                            .id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length + 1 + i).id,
                                            .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                                        });
                                    }
                                    borderOffset.x += (childElement->dimensions.width + (float)layoutConfig->childGap);
                                }
                            } else {
                                for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                                    if (i > 0) {
                                        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                                            .boundingBox = { currentElementBoundingBox.x + scrollOffset.x, currentElementBoundingBox.y + borderOffset.y + scrollOffset.y, currentElement->dimensions.width, (float)borderConfig->width.betweenChildren },
                                            .renderData = { .rectangle = {
                                                    .backgroundColor = borderConfig->color,
                                            } },
                                            .userData = sharedConfig->userData,
                                            .id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length + 1 + i).id,
                                            .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                                        });
                                    }
                                    borderOffset.y += (childElement->dimensions.height + (float)layoutConfig->childGap);
                                }
                            }
                        }
                    }
                }
                // This exists because the scissor needs to end _after_ borders between elements
                if (closeClipElement) {
                    Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                        .id = Clay__HashNumber(currentElement->id, rootElement->childrenOrTextContent.children.length + 11).id,
                        .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
                    });
                }

                dfsBuffer.length--;
                continue;
            }

            // Add children to the DFS buffer
            if (!Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
                dfsBuffer.length += currentElement->childrenOrTextContent.children.length;
                for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                    // Alignment along non layout axis
                    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                        currentElementTreeNode->nextChildOffset.y = currentElement->layoutConfig->padding.top;
                        float whiteSpaceAroundChild = currentElement->dimensions.height - (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) - childElement->dimensions.height;
                        switch (layoutConfig->childAlignment.y) {
                            case CLAY_ALIGN_Y_TOP: break;
                            case CLAY_ALIGN_Y_CENTER: currentElementTreeNode->nextChildOffset.y += whiteSpaceAroundChild / 2; break;
                            case CLAY_ALIGN_Y_BOTTOM: currentElementTreeNode->nextChildOffset.y += whiteSpaceAroundChild; break;
                        }
                    } else {
                        currentElementTreeNode->nextChildOffset.x = currentElement->layoutConfig->padding.left;
                        float whiteSpaceAroundChild = currentElement->dimensions.width - (float)(layoutConfig->padding.left + layoutConfig->padding.right) - childElement->dimensions.width;
                        switch (layoutConfig->childAlignment.x) {
                            case CLAY_ALIGN_X_LEFT: break;
                            case CLAY_ALIGN_X_CENTER: currentElementTreeNode->nextChildOffset.x += whiteSpaceAroundChild / 2; break;
                            case CLAY_ALIGN_X_RIGHT: currentElementTreeNode->nextChildOffset.x += whiteSpaceAroundChild; break;
                        }
                    }

                    Clay_Vector2 childPosition = {
                        currentElementTreeNode->position.x + currentElementTreeNode->nextChildOffset.x + scrollOffset.x,
                        currentElementTreeNode->position.y + currentElementTreeNode->nextChildOffset.y + scrollOffset.y,
                    };

                    // DFS buffer elements need to be added in reverse because stack traversal happens backwards
                    uint32_t newNodeIndex = dfsBuffer.length - 1 - i;
                    dfsBuffer.internalArray[newNodeIndex] = CLAY__INIT(Clay__LayoutElementTreeNode) {
                        .layoutElement = childElement,
                        .position = { childPosition.x, childPosition.y },
                        .nextChildOffset = { .x = (float)childElement->layoutConfig->padding.left, .y = (float)childElement->layoutConfig->padding.top },
                    };
                    context->treeNodeVisited.internalArray[newNodeIndex] = false;

                    // Update parent offsets
                    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                        currentElementTreeNode->nextChildOffset.x += childElement->dimensions.width + (float)layoutConfig->childGap;
                    } else {
                        currentElementTreeNode->nextChildOffset.y += childElement->dimensions.height + (float)layoutConfig->childGap;
                    }
                }
            }
        }

        if (root->clipElementId) {
            Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) { .id = Clay__HashNumber(rootElement->id, rootElement->childrenOrTextContent.children.length + 11).id, .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END });
        }
    }
}

CLAY_WASM_EXPORT("Clay_GetPointerOverIds")
CLAY_DLL_EXPORT Clay_ElementIdArray Clay_GetPointerOverIds(void) {
    return Clay_GetCurrentContext()->pointerOverIds;
}

#pragma region DebugTools
Clay_Color CLAY__DEBUGVIEW_COLOR_1 = {58, 56, 52, 255};
Clay_Color CLAY__DEBUGVIEW_COLOR_2 = {62, 60, 58, 255};
Clay_Color CLAY__DEBUGVIEW_COLOR_3 = {141, 133, 135, 255};
Clay_Color CLAY__DEBUGVIEW_COLOR_4 = {238, 226, 231, 255};
Clay_Color CLAY__DEBUGVIEW_COLOR_SELECTED_ROW = {102, 80, 78, 255};
const int32_t CLAY__DEBUGVIEW_ROW_HEIGHT = 30;
const int32_t CLAY__DEBUGVIEW_OUTER_PADDING = 10;
const int32_t CLAY__DEBUGVIEW_INDENT_WIDTH = 16;
Clay_TextElementConfig Clay__DebugView_TextNameConfig = {.textColor = {238, 226, 231, 255}, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE };
Clay_LayoutConfig Clay__DebugView_ScrollViewItemLayoutConfig = CLAY__DEFAULT_STRUCT;

typedef struct {
    Clay_String label;
    Clay_Color color;
} Clay__DebugElementConfigTypeLabelConfig;

Clay__DebugElementConfigTypeLabelConfig Clay__DebugGetElementConfigTypeLabel(Clay__ElementConfigType type) {
    switch (type) {
        case CLAY__ELEMENT_CONFIG_TYPE_SHARED: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Shared"), {243,134,48,255} };
        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Text"), {105,210,231,255} };
        case CLAY__ELEMENT_CONFIG_TYPE_ASPECT: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Aspect"), {101,149,194,255} };
        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Image"), {121,189,154,255} };
        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Floating"), {250,105,0,255} };
        case CLAY__ELEMENT_CONFIG_TYPE_CLIP: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) {CLAY_STRING("Scroll"), {242, 196, 90, 255} };
        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) {CLAY_STRING("Border"), {108, 91, 123, 255} };
        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Custom"), {11,72,107,255} };
        default: break;
    }
    return CLAY__INIT(Clay__DebugElementConfigTypeLabelConfig) { CLAY_STRING("Error"), {0,0,0,255} };
}

typedef struct {
    int32_t rowCount;
    int32_t selectedElementRowIndex;
} Clay__RenderDebugLayoutData;

// Returns row count
Clay__RenderDebugLayoutData Clay__RenderDebugLayoutElementsList(int32_t initialRootsLength, int32_t highlightedRowIndex) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__int32_tArray dfsBuffer = context->reusableElementIndexBuffer;
    Clay__DebugView_ScrollViewItemLayoutConfig = CLAY__INIT(Clay_LayoutConfig) { .sizing = { .height = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT) }, .childGap = 6, .childAlignment = { .y = CLAY_ALIGN_Y_CENTER }};
    Clay__RenderDebugLayoutData layoutData = CLAY__DEFAULT_STRUCT;

    uint32_t highlightedElementId = 0;

    for (int32_t rootIndex = 0; rootIndex < initialRootsLength; ++rootIndex) {
        dfsBuffer.length = 0;
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
        Clay__int32_tArray_Add(&dfsBuffer, (int32_t)root->layoutElementIndex);
        context->treeNodeVisited.internalArray[0] = false;
        if (rootIndex > 0) {
            CLAY({ .id = CLAY_IDI("Clay__DebugView_EmptyRowOuter", rootIndex), .layout = { .sizing = {.width = CLAY_SIZING_GROW(0)}, .padding = {CLAY__DEBUGVIEW_INDENT_WIDTH / 2, 0, 0, 0} } }) {
                CLAY({ .id = CLAY_IDI("Clay__DebugView_EmptyRow", rootIndex), .layout = { .sizing = { .width = CLAY_SIZING_GROW(0), .height = CLAY_SIZING_FIXED((float)CLAY__DEBUGVIEW_ROW_HEIGHT) }}, .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .width = { .top = 1 } } }) {}
            }
            layoutData.rowCount++;
        }
        while (dfsBuffer.length > 0) {
            int32_t currentElementIndex = Clay__int32_tArray_GetValue(&dfsBuffer, (int)dfsBuffer.length - 1);
            Clay_LayoutElement *currentElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)currentElementIndex);
            if (context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
                if (!Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && currentElement->childrenOrTextContent.children.length > 0) {
                    Clay__CloseElement();
                    Clay__CloseElement();
                    Clay__CloseElement();
                }
                dfsBuffer.length--;
                continue;
            }

            if (highlightedRowIndex == layoutData.rowCount) {
                if (context->pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
                    context->debugSelectedElementId = currentElement->id;
                }
                highlightedElementId = currentElement->id;
            }

            context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
            Clay_LayoutElementHashMapItem *currentElementData = Clay__GetHashMapItem(currentElement->id);
            bool offscreen = Clay__ElementIsOffscreen(&currentElementData->boundingBox);
            if (context->debugSelectedElementId == currentElement->id) {
                layoutData.selectedElementRowIndex = layoutData.rowCount;
            }
            CLAY({ .id = CLAY_IDI("Clay__DebugView_ElementOuter", currentElement->id), .layout = Clay__DebugView_ScrollViewItemLayoutConfig }) {
                // Collapse icon / button
                if (!(Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement->childrenOrTextContent.children.length == 0)) {
                    CLAY({
                        .id = CLAY_IDI("Clay__DebugView_CollapseElement", currentElement->id),
                        .layout = { .sizing = {CLAY_SIZING_FIXED(16), CLAY_SIZING_FIXED(16)}, .childAlignment = { CLAY_ALIGN_X_CENTER, CLAY_ALIGN_Y_CENTER} },
                        .cornerRadius = CLAY_CORNER_RADIUS(4),
                        .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .width = {1, 1, 1, 1, 0} },
                    }) {
                        CLAY_TEXT((currentElementData && currentElementData->debugData->collapsed) ? CLAY_STRING("+") : CLAY_STRING("-"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
                    }
                } else { // Square dot for empty containers
                    CLAY({ .layout = { .sizing = {CLAY_SIZING_FIXED(16), CLAY_SIZING_FIXED(16)}, .childAlignment = { CLAY_ALIGN_X_CENTER, CLAY_ALIGN_Y_CENTER } } }) {
                        CLAY({ .layout = { .sizing = {CLAY_SIZING_FIXED(8), CLAY_SIZING_FIXED(8)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3, .cornerRadius = CLAY_CORNER_RADIUS(2) }) {}
                    }
                }
                // Collisions and offscreen info
                if (currentElementData) {
                    if (currentElementData->debugData->collision) {
                        CLAY({ .layout = { .padding = { 8, 8, 2, 2 }}, .border = { .color = {177, 147, 8, 255}, .width = {1, 1, 1, 1, 0} } }) {
                            CLAY_TEXT(CLAY_STRING("Duplicate ID"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }));
                        }
                    }
                    if (offscreen) {
                        CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .border = {  .color = CLAY__DEBUGVIEW_COLOR_3, .width = { 1, 1, 1, 1, 0} } }) {
                            CLAY_TEXT(CLAY_STRING("Offscreen"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }));
                        }
                    }
                }
                Clay_String idString = context->layoutElementIdStrings.internalArray[currentElementIndex];
                if (idString.length > 0) {
                    CLAY_TEXT(idString, offscreen ? CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }) : &Clay__DebugView_TextNameConfig);
                }
                for (int32_t elementConfigIndex = 0; elementConfigIndex < currentElement->elementConfigs.length; ++elementConfigIndex) {
                    Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, elementConfigIndex);
                    if (elementConfig->type == CLAY__ELEMENT_CONFIG_TYPE_SHARED) {
                        Clay_Color labelColor = {243,134,48,90};
                        labelColor.a = 90;
                        Clay_Color backgroundColor = elementConfig->config.sharedElementConfig->backgroundColor;
                        Clay_CornerRadius radius = elementConfig->config.sharedElementConfig->cornerRadius;
                        if (backgroundColor.a > 0) {
                            CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = labelColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = labelColor, .width = { 1, 1, 1, 1, 0} } }) {
                                CLAY_TEXT(CLAY_STRING("Color"), CLAY_TEXT_CONFIG({ .textColor = offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
                            }
                        }
                        if (radius.bottomLeft > 0) {
                            CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = labelColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = labelColor, .width = { 1, 1, 1, 1, 0 } } }) {
                                CLAY_TEXT(CLAY_STRING("Radius"), CLAY_TEXT_CONFIG({ .textColor = offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
                            }
                        }
                        continue;
                    }
                    Clay__DebugElementConfigTypeLabelConfig config = Clay__DebugGetElementConfigTypeLabel(elementConfig->type);
                    Clay_Color backgroundColor = config.color;
                    backgroundColor.a = 90;
                    CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = backgroundColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = config.color, .width = { 1, 1, 1, 1, 0 } } }) {
                        CLAY_TEXT(config.label, CLAY_TEXT_CONFIG({ .textColor = offscreen ? CLAY__DEBUGVIEW_COLOR_3 : CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
                    }
                }
            }

            // Render the text contents below the element as a non-interactive row
            if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
                layoutData.rowCount++;
                Clay__TextElementData *textElementData = currentElement->childrenOrTextContent.textElementData;
                Clay_TextElementConfig *rawTextConfig = offscreen ? CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16 }) : &Clay__DebugView_TextNameConfig;
                CLAY({ .layout = { .sizing = { .height = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .childAlignment = { .y = CLAY_ALIGN_Y_CENTER } } }) {
                    CLAY({ .layout = { .sizing = {.width = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_INDENT_WIDTH + 16) } } }) {}
                    CLAY_TEXT(CLAY_STRING("\""), rawTextConfig);
                    CLAY_TEXT(textElementData->text.length > 40 ? (CLAY__INIT(Clay_String) { .length = 40, .chars = textElementData->text.chars }) : textElementData->text, rawTextConfig);
                    if (textElementData->text.length > 40) {
                        CLAY_TEXT(CLAY_STRING("..."), rawTextConfig);
                    }
                    CLAY_TEXT(CLAY_STRING("\""), rawTextConfig);
                }
            } else if (currentElement->childrenOrTextContent.children.length > 0) {
                Clay__OpenElement();
                Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) { .layout = { .padding = { .left = 8 } } });
                Clay__OpenElement();
                Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) { .layout = { .padding = { .left = CLAY__DEBUGVIEW_INDENT_WIDTH }}, .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .width = { .left = 1 } }});
                Clay__OpenElement();
                Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) { .layout = { .layoutDirection = CLAY_TOP_TO_BOTTOM } });
            }

            layoutData.rowCount++;
            if (!(Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (currentElementData && currentElementData->debugData->collapsed))) {
                for (int32_t i = currentElement->childrenOrTextContent.children.length - 1; i >= 0; --i) {
                    Clay__int32_tArray_Add(&dfsBuffer, currentElement->childrenOrTextContent.children.elements[i]);
                    context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = false; // TODO needs to be ranged checked
                }
            }
        }
    }

    if (context->pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
        Clay_ElementId collapseButtonId = Clay__HashString(CLAY_STRING("Clay__DebugView_CollapseElement"), 0);
        for (int32_t i = (int)context->pointerOverIds.length - 1; i >= 0; i--) {
            Clay_ElementId *elementId = Clay_ElementIdArray_Get(&context->pointerOverIds, i);
            if (elementId->baseId == collapseButtonId.baseId) {
                Clay_LayoutElementHashMapItem *highlightedItem = Clay__GetHashMapItem(elementId->offset);
                highlightedItem->debugData->collapsed = !highlightedItem->debugData->collapsed;
                break;
            }
        }
    }

    if (highlightedElementId) {
        CLAY({ .id = CLAY_ID("Clay__DebugView_ElementHighlight"), .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .floating = { .parentId = highlightedElementId, .zIndex = 32767, .pointerCaptureMode = CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH, .attachTo = CLAY_ATTACH_TO_ELEMENT_WITH_ID } }) {
            CLAY({ .id = CLAY_ID("Clay__DebugView_ElementHighlightRectangle"), .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .backgroundColor = Clay__debugViewHighlightColor }) {}
        }
    }
    return layoutData;
}

void Clay__RenderDebugLayoutSizing(Clay_SizingAxis sizing, Clay_TextElementConfig *infoTextConfig) {
    Clay_String sizingLabel = CLAY_STRING("GROW");
    if (sizing.type == CLAY__SIZING_TYPE_FIT) {
        sizingLabel = CLAY_STRING("FIT");
    } else if (sizing.type == CLAY__SIZING_TYPE_PERCENT) {
        sizingLabel = CLAY_STRING("PERCENT");
    } else if (sizing.type == CLAY__SIZING_TYPE_FIXED) {
        sizingLabel = CLAY_STRING("FIXED");
    }
    CLAY_TEXT(sizingLabel, infoTextConfig);
    if (sizing.type == CLAY__SIZING_TYPE_GROW || sizing.type == CLAY__SIZING_TYPE_FIT || sizing.type == CLAY__SIZING_TYPE_FIXED) {
        CLAY_TEXT(CLAY_STRING("("), infoTextConfig);
        if (sizing.size.minMax.min != 0) {
            CLAY_TEXT(CLAY_STRING("min: "), infoTextConfig);
            CLAY_TEXT(Clay__IntToString(sizing.size.minMax.min), infoTextConfig);
            if (sizing.size.minMax.max != CLAY__MAXFLOAT) {
                CLAY_TEXT(CLAY_STRING(", "), infoTextConfig);
            }
        }
        if (sizing.size.minMax.max != CLAY__MAXFLOAT) {
            CLAY_TEXT(CLAY_STRING("max: "), infoTextConfig);
            CLAY_TEXT(Clay__IntToString(sizing.size.minMax.max), infoTextConfig);
        }
        CLAY_TEXT(CLAY_STRING(")"), infoTextConfig);
    } else if (sizing.type == CLAY__SIZING_TYPE_PERCENT) {
        CLAY_TEXT(CLAY_STRING("("), infoTextConfig);
        CLAY_TEXT(Clay__IntToString(sizing.size.percent * 100), infoTextConfig);
        CLAY_TEXT(CLAY_STRING("%)"), infoTextConfig);
    }
}

void Clay__RenderDebugViewElementConfigHeader(Clay_String elementId, Clay__ElementConfigType type) {
    Clay__DebugElementConfigTypeLabelConfig config = Clay__DebugGetElementConfigTypeLabel(type);
    Clay_Color backgroundColor = config.color;
    backgroundColor.a = 90;
    CLAY({ .layout = { .sizing = { .width = CLAY_SIZING_GROW(0) }, .padding = CLAY_PADDING_ALL(CLAY__DEBUGVIEW_OUTER_PADDING), .childAlignment = { .y = CLAY_ALIGN_Y_CENTER } } }) {
        CLAY({ .layout = { .padding = { 8, 8, 2, 2 } }, .backgroundColor = backgroundColor, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = config.color, .width = { 1, 1, 1, 1, 0 } } }) {
            CLAY_TEXT(config.label, CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
        }
        CLAY({ .layout = { .sizing = { .width = CLAY_SIZING_GROW(0) } } }) {}
        CLAY_TEXT(elementId, CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE }));
    }
}

void Clay__RenderDebugViewColor(Clay_Color color, Clay_TextElementConfig *textConfig) {
    CLAY({ .layout = { .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
        CLAY_TEXT(CLAY_STRING("{ r: "), textConfig);
        CLAY_TEXT(Clay__IntToString(color.r), textConfig);
        CLAY_TEXT(CLAY_STRING(", g: "), textConfig);
        CLAY_TEXT(Clay__IntToString(color.g), textConfig);
        CLAY_TEXT(CLAY_STRING(", b: "), textConfig);
        CLAY_TEXT(Clay__IntToString(color.b), textConfig);
        CLAY_TEXT(CLAY_STRING(", a: "), textConfig);
        CLAY_TEXT(Clay__IntToString(color.a), textConfig);
        CLAY_TEXT(CLAY_STRING(" }"), textConfig);
        CLAY({ .layout = { .sizing = { .width = CLAY_SIZING_FIXED(10) } } }) {}
        CLAY({ .layout = { .sizing = { CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 8), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 8)} }, .backgroundColor = color, .cornerRadius = CLAY_CORNER_RADIUS(4), .border = { .color = CLAY__DEBUGVIEW_COLOR_4, .width = { 1, 1, 1, 1, 0 } } }) {}
    }
}

void Clay__RenderDebugViewCornerRadius(Clay_CornerRadius cornerRadius, Clay_TextElementConfig *textConfig) {
    CLAY({ .layout = { .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
        CLAY_TEXT(CLAY_STRING("{ topLeft: "), textConfig);
        CLAY_TEXT(Clay__IntToString(cornerRadius.topLeft), textConfig);
        CLAY_TEXT(CLAY_STRING(", topRight: "), textConfig);
        CLAY_TEXT(Clay__IntToString(cornerRadius.topRight), textConfig);
        CLAY_TEXT(CLAY_STRING(", bottomLeft: "), textConfig);
        CLAY_TEXT(Clay__IntToString(cornerRadius.bottomLeft), textConfig);
        CLAY_TEXT(CLAY_STRING(", bottomRight: "), textConfig);
        CLAY_TEXT(Clay__IntToString(cornerRadius.bottomRight), textConfig);
        CLAY_TEXT(CLAY_STRING(" }"), textConfig);
    }
}

void HandleDebugViewCloseButtonInteraction(Clay_ElementId elementId, Clay_PointerData pointerInfo, intptr_t userData) {
    Clay_Context* context = Clay_GetCurrentContext();
    (void) elementId; (void) pointerInfo; (void) userData;
    if (pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
        context->debugModeEnabled = false;
    }
}

void Clay__RenderDebugView(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay_ElementId closeButtonId = Clay__HashString(CLAY_STRING("Clay__DebugViewTopHeaderCloseButtonOuter"), 0);
    if (context->pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
        for (int32_t i = 0; i < context->pointerOverIds.length; ++i) {
            Clay_ElementId *elementId = Clay_ElementIdArray_Get(&context->pointerOverIds, i);
            if (elementId->id == closeButtonId.id) {
                context->debugModeEnabled = false;
                return;
            }
        }
    }

    uint32_t initialRootsLength = context->layoutElementTreeRoots.length;
    uint32_t initialElementsLength = context->layoutElements.length;
    Clay_TextElementConfig *infoTextConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
    Clay_TextElementConfig *infoTitleConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_3, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
    Clay_ElementId scrollId = Clay__HashString(CLAY_STRING("Clay__DebugViewOuterScrollPane"), 0);
    float scrollYOffset = 0;
    bool pointerInDebugView = context->pointerInfo.position.y < context->layoutDimensions.height - 300;
    for (int32_t i = 0; i < context->scrollContainerDatas.length; ++i) {
        Clay__ScrollContainerDataInternal *scrollContainerData = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
        if (scrollContainerData->elementId == scrollId.id) {
            if (!context->externalScrollHandlingEnabled) {
                scrollYOffset = scrollContainerData->scrollPosition.y;
            } else {
                pointerInDebugView = context->pointerInfo.position.y + scrollContainerData->scrollPosition.y < context->layoutDimensions.height - 300;
            }
            break;
        }
    }
    int32_t highlightedRow = pointerInDebugView
            ? (int32_t)((context->pointerInfo.position.y - scrollYOffset) / (float)CLAY__DEBUGVIEW_ROW_HEIGHT) - 1
            : -1;
    if (context->pointerInfo.position.x < context->layoutDimensions.width - (float)Clay__debugViewWidth) {
        highlightedRow = -1;
    }
    Clay__RenderDebugLayoutData layoutData = CLAY__DEFAULT_STRUCT;
    CLAY({ .id = CLAY_ID("Clay__DebugView"),
         .layout = { .sizing = { CLAY_SIZING_FIXED((float)Clay__debugViewWidth) , CLAY_SIZING_FIXED(context->layoutDimensions.height) }, .layoutDirection = CLAY_TOP_TO_BOTTOM },
        .floating = { .zIndex = 32765, .attachPoints = { .element = CLAY_ATTACH_POINT_LEFT_CENTER, .parent = CLAY_ATTACH_POINT_RIGHT_CENTER }, .attachTo = CLAY_ATTACH_TO_ROOT, .clipTo = CLAY_CLIP_TO_ATTACHED_PARENT },
        .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .width = { .bottom = 1 } }
    }) {
        CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_2 }) {
            CLAY_TEXT(CLAY_STRING("Clay Debug Tools"), infoTextConfig);
            CLAY({ .layout = { .sizing = { .width = CLAY_SIZING_GROW(0) } } }) {}
            // Close button
            CLAY({
                .layout = { .sizing = {CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10)}, .childAlignment = {CLAY_ALIGN_X_CENTER, CLAY_ALIGN_Y_CENTER} },
                .backgroundColor = {217,91,67,80},
                .cornerRadius = CLAY_CORNER_RADIUS(4),
                .border = { .color = { 217,91,67,255 }, .width = { 1, 1, 1, 1, 0 } },
            }) {
                Clay_OnHover(HandleDebugViewCloseButtonInteraction, 0);
                CLAY_TEXT(CLAY_STRING("x"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
            }
        }
        CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(1)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3 } ) {}
        CLAY({ .id = scrollId, .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .clip = { .horizontal = true, .vertical = true, .childOffset = Clay_GetScrollOffset() } }) {
            CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)}, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = ((initialElementsLength + initialRootsLength) & 1) == 0 ? CLAY__DEBUGVIEW_COLOR_2 : CLAY__DEBUGVIEW_COLOR_1 }) {
                Clay_ElementId panelContentsId = Clay__HashString(CLAY_STRING("Clay__DebugViewPaneOuter"), 0);
                // Element list
                CLAY({ .id = panelContentsId, .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .floating = { .zIndex = 32766, .pointerCaptureMode = CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH, .attachTo = CLAY_ATTACH_TO_PARENT, .clipTo = CLAY_CLIP_TO_ATTACHED_PARENT } }) {
                    CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)}, .padding = { CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                        layoutData = Clay__RenderDebugLayoutElementsList((int32_t)initialRootsLength, highlightedRow);
                    }
                }
                float contentWidth = Clay__GetHashMapItem(panelContentsId.id)->layoutElement->dimensions.width;
                CLAY({ .layout = { .sizing = {.width = CLAY_SIZING_FIXED(contentWidth) }, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {}
                for (int32_t i = 0; i < layoutData.rowCount; i++) {
                    Clay_Color rowColor = (i & 1) == 0 ? CLAY__DEBUGVIEW_COLOR_2 : CLAY__DEBUGVIEW_COLOR_1;
                    if (i == layoutData.selectedElementRowIndex) {
                        rowColor = CLAY__DEBUGVIEW_COLOR_SELECTED_ROW;
                    }
                    if (i == highlightedRow) {
                        rowColor.r *= 1.25f;
                        rowColor.g *= 1.25f;
                        rowColor.b *= 1.25f;
                    }
                    CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = rowColor } ) {}
                }
            }
        }
        CLAY({ .layout = { .sizing = {.width = CLAY_SIZING_GROW(0), .height = CLAY_SIZING_FIXED(1)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3 }) {}
        if (context->debugSelectedElementId != 0) {
            Clay_LayoutElementHashMapItem *selectedItem = Clay__GetHashMapItem(context->debugSelectedElementId);
            CLAY({
                .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(300)}, .layoutDirection = CLAY_TOP_TO_BOTTOM },
                .backgroundColor = CLAY__DEBUGVIEW_COLOR_2 ,
                .clip = { .vertical = true, .childOffset = Clay_GetScrollOffset() },
                .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .width = { .betweenChildren = 1 } }
            }) {
                CLAY({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT + 8)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
                    CLAY_TEXT(CLAY_STRING("Layout Config"), infoTextConfig);
                    CLAY({ .layout = { .sizing = { .width = CLAY_SIZING_GROW(0) } } }) {}
                    if (selectedItem->elementId.stringId.length != 0) {
                        CLAY_TEXT(selectedItem->elementId.stringId, infoTitleConfig);
                        if (selectedItem->elementId.offset != 0) {
                            CLAY_TEXT(CLAY_STRING(" ("), infoTitleConfig);
                            CLAY_TEXT(Clay__IntToString(selectedItem->elementId.offset), infoTitleConfig);
                            CLAY_TEXT(CLAY_STRING(")"), infoTitleConfig);
                        }
                    }
                }
                Clay_Padding attributeConfigPadding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 8, 8};
                // Clay_LayoutConfig debug info
                CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                    // .boundingBox
                    CLAY_TEXT(CLAY_STRING("Bounding Box"), infoTitleConfig);
                    CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                        CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.x), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.y), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(", width: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.width), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.height), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
                    }
                    // .layoutDirection
                    CLAY_TEXT(CLAY_STRING("Layout Direction"), infoTitleConfig);
                    Clay_LayoutConfig *layoutConfig = selectedItem->layoutElement->layoutConfig;
                    CLAY_TEXT(layoutConfig->layoutDirection == CLAY_TOP_TO_BOTTOM ? CLAY_STRING("TOP_TO_BOTTOM") : CLAY_STRING("LEFT_TO_RIGHT"), infoTextConfig);
                    // .sizing
                    CLAY_TEXT(CLAY_STRING("Sizing"), infoTitleConfig);
                    CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                        CLAY_TEXT(CLAY_STRING("width: "), infoTextConfig);
                        Clay__RenderDebugLayoutSizing(layoutConfig->sizing.width, infoTextConfig);
                    }
                    CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                        CLAY_TEXT(CLAY_STRING("height: "), infoTextConfig);
                        Clay__RenderDebugLayoutSizing(layoutConfig->sizing.height, infoTextConfig);
                    }
                    // .padding
                    CLAY_TEXT(CLAY_STRING("Padding"), infoTitleConfig);
                    CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoPadding") }) {
                        CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(layoutConfig->padding.left), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(layoutConfig->padding.right), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(layoutConfig->padding.top), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig);
                        CLAY_TEXT(Clay__IntToString(layoutConfig->padding.bottom), infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
                    }
                    // .childGap
                    CLAY_TEXT(CLAY_STRING("Child Gap"), infoTitleConfig);
                    CLAY_TEXT(Clay__IntToString(layoutConfig->childGap), infoTextConfig);
                    // .childAlignment
                    CLAY_TEXT(CLAY_STRING("Child Alignment"), infoTitleConfig);
                    CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                        CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
                        Clay_String alignX = CLAY_STRING("LEFT");
                        if (layoutConfig->childAlignment.x == CLAY_ALIGN_X_CENTER) {
                            alignX = CLAY_STRING("CENTER");
                        } else if (layoutConfig->childAlignment.x == CLAY_ALIGN_X_RIGHT) {
                            alignX = CLAY_STRING("RIGHT");
                        }
                        CLAY_TEXT(alignX, infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
                        Clay_String alignY = CLAY_STRING("TOP");
                        if (layoutConfig->childAlignment.y == CLAY_ALIGN_Y_CENTER) {
                            alignY = CLAY_STRING("CENTER");
                        } else if (layoutConfig->childAlignment.y == CLAY_ALIGN_Y_BOTTOM) {
                            alignY = CLAY_STRING("BOTTOM");
                        }
                        CLAY_TEXT(alignY, infoTextConfig);
                        CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
                    }
                }
                for (int32_t elementConfigIndex = 0; elementConfigIndex < selectedItem->layoutElement->elementConfigs.length; ++elementConfigIndex) {
                    Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&selectedItem->layoutElement->elementConfigs, elementConfigIndex);
                    Clay__RenderDebugViewElementConfigHeader(selectedItem->elementId.stringId, elementConfig->type);
                    switch (elementConfig->type) {
                        case CLAY__ELEMENT_CONFIG_TYPE_SHARED: {
                            Clay_SharedElementConfig *sharedConfig = elementConfig->config.sharedElementConfig;
                            CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM }}) {
                                // .backgroundColor
                                CLAY_TEXT(CLAY_STRING("Background Color"), infoTitleConfig);
                                Clay__RenderDebugViewColor(sharedConfig->backgroundColor, infoTextConfig);
                                // .cornerRadius
                                CLAY_TEXT(CLAY_STRING("Corner Radius"), infoTitleConfig);
                                Clay__RenderDebugViewCornerRadius(sharedConfig->cornerRadius, infoTextConfig);
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: {
                            Clay_TextElementConfig *textConfig = elementConfig->config.textElementConfig;
                            CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                                // .fontSize
                                CLAY_TEXT(CLAY_STRING("Font Size"), infoTitleConfig);
                                CLAY_TEXT(Clay__IntToString(textConfig->fontSize), infoTextConfig);
                                // .fontId
                                CLAY_TEXT(CLAY_STRING("Font ID"), infoTitleConfig);
                                CLAY_TEXT(Clay__IntToString(textConfig->fontId), infoTextConfig);
                                // .lineHeight
                                CLAY_TEXT(CLAY_STRING("Line Height"), infoTitleConfig);
                                CLAY_TEXT(textConfig->lineHeight == 0 ? CLAY_STRING("auto") : Clay__IntToString(textConfig->lineHeight), infoTextConfig);
                                // .letterSpacing
                                CLAY_TEXT(CLAY_STRING("Letter Spacing"), infoTitleConfig);
                                CLAY_TEXT(Clay__IntToString(textConfig->letterSpacing), infoTextConfig);
                                // .wrapMode
                                CLAY_TEXT(CLAY_STRING("Wrap Mode"), infoTitleConfig);
                                Clay_String wrapMode = CLAY_STRING("WORDS");
                                if (textConfig->wrapMode == CLAY_TEXT_WRAP_NONE) {
                                    wrapMode = CLAY_STRING("NONE");
                                } else if (textConfig->wrapMode == CLAY_TEXT_WRAP_NEWLINES) {
                                    wrapMode = CLAY_STRING("NEWLINES");
                                }
                                CLAY_TEXT(wrapMode, infoTextConfig);
                                // .textAlignment
                                CLAY_TEXT(CLAY_STRING("Text Alignment"), infoTitleConfig);
                                Clay_String textAlignment = CLAY_STRING("LEFT");
                                if (textConfig->textAlignment == CLAY_TEXT_ALIGN_CENTER) {
                                    textAlignment = CLAY_STRING("CENTER");
                                } else if (textConfig->textAlignment == CLAY_TEXT_ALIGN_RIGHT) {
                                    textAlignment = CLAY_STRING("RIGHT");
                                }
                                CLAY_TEXT(textAlignment, infoTextConfig);
                                // .textColor
                                CLAY_TEXT(CLAY_STRING("Text Color"), infoTitleConfig);
                                Clay__RenderDebugViewColor(textConfig->textColor, infoTextConfig);
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_ASPECT: {
                            Clay_AspectRatioElementConfig *aspectRatioConfig = elementConfig->config.aspectRatioElementConfig;
                            CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoAspectRatioBody"), .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                                CLAY_TEXT(CLAY_STRING("Aspect Ratio"), infoTitleConfig);
                                // Aspect Ratio
                                CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoAspectRatio") }) {
                                    CLAY_TEXT(Clay__IntToString(aspectRatioConfig->aspectRatio), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING("."), infoTextConfig);
                                    float frac = aspectRatioConfig->aspectRatio - (int)(aspectRatioConfig->aspectRatio);
                                    frac *= 100;
                                    if ((int)frac < 10) {
                                        CLAY_TEXT(CLAY_STRING("0"), infoTextConfig);
                                    }
                                    CLAY_TEXT(Clay__IntToString(frac), infoTextConfig);
                                }
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: {
                            Clay_ImageElementConfig *imageConfig = elementConfig->config.imageElementConfig;
                            Clay_AspectRatioElementConfig aspectConfig = { 1 };
                            if (Clay__ElementHasConfig(selectedItem->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT)) {
                                aspectConfig = *Clay__FindElementConfigWithType(selectedItem->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig;
                            }
                            CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoImageBody"), .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                                // Image Preview
                                CLAY_TEXT(CLAY_STRING("Preview"), infoTitleConfig);
                                CLAY({ .layout = { .sizing = { .width = CLAY_SIZING_GROW(64, 128), .height = CLAY_SIZING_GROW(64, 128) }}, .aspectRatio = aspectConfig, .image = *imageConfig }) {}
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_CLIP: {
                            Clay_ClipElementConfig *clipConfig = elementConfig->config.clipElementConfig;
                            CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                                // .vertical
                                CLAY_TEXT(CLAY_STRING("Vertical"), infoTitleConfig);
                                CLAY_TEXT(clipConfig->vertical ? CLAY_STRING("true") : CLAY_STRING("false") , infoTextConfig);
                                // .horizontal
                                CLAY_TEXT(CLAY_STRING("Horizontal"), infoTitleConfig);
                                CLAY_TEXT(clipConfig->horizontal ? CLAY_STRING("true") : CLAY_STRING("false") , infoTextConfig);
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING: {
                            Clay_FloatingElementConfig *floatingConfig = elementConfig->config.floatingElementConfig;
                            CLAY({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                                // .offset
                                CLAY_TEXT(CLAY_STRING("Offset"), infoTitleConfig);
                                CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                                    CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(floatingConfig->offset.x), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(floatingConfig->offset.y), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
                                }
                                // .expand
                                CLAY_TEXT(CLAY_STRING("Expand"), infoTitleConfig);
                                CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                                    CLAY_TEXT(CLAY_STRING("{ width: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(floatingConfig->expand.width), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(floatingConfig->expand.height), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
                                }
                                // .zIndex
                                CLAY_TEXT(CLAY_STRING("z-index"), infoTitleConfig);
                                CLAY_TEXT(Clay__IntToString(floatingConfig->zIndex), infoTextConfig);
                                // .parentId
                                CLAY_TEXT(CLAY_STRING("Parent"), infoTitleConfig);
                                Clay_LayoutElementHashMapItem *hashItem = Clay__GetHashMapItem(floatingConfig->parentId);
                                CLAY_TEXT(hashItem->elementId.stringId, infoTextConfig);
                                // .attachPoints
                                CLAY_TEXT(CLAY_STRING("Attach Points"), infoTitleConfig);
                                CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                                    CLAY_TEXT(CLAY_STRING("{ element: "), infoTextConfig);
                                    Clay_String attachPointElement = CLAY_STRING("LEFT_TOP");
                                    if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_LEFT_CENTER) {
                                        attachPointElement = CLAY_STRING("LEFT_CENTER");
                                    } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_LEFT_BOTTOM) {
                                        attachPointElement = CLAY_STRING("LEFT_BOTTOM");
                                    } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_CENTER_TOP) {
                                        attachPointElement = CLAY_STRING("CENTER_TOP");
                                    } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_CENTER_CENTER) {
                                        attachPointElement = CLAY_STRING("CENTER_CENTER");
                                    } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_CENTER_BOTTOM) {
                                        attachPointElement = CLAY_STRING("CENTER_BOTTOM");
                                    } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_RIGHT_TOP) {
                                        attachPointElement = CLAY_STRING("RIGHT_TOP");
                                    } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_RIGHT_CENTER) {
                                        attachPointElement = CLAY_STRING("RIGHT_CENTER");
                                    } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_RIGHT_BOTTOM) {
                                        attachPointElement = CLAY_STRING("RIGHT_BOTTOM");
                                    }
                                    CLAY_TEXT(attachPointElement, infoTextConfig);
                                    Clay_String attachPointParent = CLAY_STRING("LEFT_TOP");
                                    if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_LEFT_CENTER) {
                                        attachPointParent = CLAY_STRING("LEFT_CENTER");
                                    } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_LEFT_BOTTOM) {
                                        attachPointParent = CLAY_STRING("LEFT_BOTTOM");
                                    } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_CENTER_TOP) {
                                        attachPointParent = CLAY_STRING("CENTER_TOP");
                                    } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_CENTER_CENTER) {
                                        attachPointParent = CLAY_STRING("CENTER_CENTER");
                                    } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_CENTER_BOTTOM) {
                                        attachPointParent = CLAY_STRING("CENTER_BOTTOM");
                                    } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_TOP) {
                                        attachPointParent = CLAY_STRING("RIGHT_TOP");
                                    } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_CENTER) {
                                        attachPointParent = CLAY_STRING("RIGHT_CENTER");
                                    } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_BOTTOM) {
                                        attachPointParent = CLAY_STRING("RIGHT_BOTTOM");
                                    }
                                    CLAY_TEXT(CLAY_STRING(", parent: "), infoTextConfig);
                                    CLAY_TEXT(attachPointParent, infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
                                }
                                // .pointerCaptureMode
                                CLAY_TEXT(CLAY_STRING("Pointer Capture Mode"), infoTitleConfig);
                                Clay_String pointerCaptureMode = CLAY_STRING("NONE");
                                if (floatingConfig->pointerCaptureMode == CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH) {
                                    pointerCaptureMode = CLAY_STRING("PASSTHROUGH");
                                }
                                CLAY_TEXT(pointerCaptureMode, infoTextConfig);
                                // .attachTo
                                CLAY_TEXT(CLAY_STRING("Attach To"), infoTitleConfig);
                                Clay_String attachTo = CLAY_STRING("NONE");
                                if (floatingConfig->attachTo == CLAY_ATTACH_TO_PARENT) {
                                    attachTo = CLAY_STRING("PARENT");
                                } else if (floatingConfig->attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID) {
                                    attachTo = CLAY_STRING("ELEMENT_WITH_ID");
                                } else if (floatingConfig->attachTo == CLAY_ATTACH_TO_ROOT) {
                                    attachTo = CLAY_STRING("ROOT");
                                }
                                CLAY_TEXT(attachTo, infoTextConfig);
                                // .clipTo
                                CLAY_TEXT(CLAY_STRING("Clip To"), infoTitleConfig);
                                Clay_String clipTo = CLAY_STRING("ATTACHED_PARENT");
                                if (floatingConfig->clipTo == CLAY_CLIP_TO_NONE) {
                                    clipTo = CLAY_STRING("NONE");
                                }
                                CLAY_TEXT(clipTo, infoTextConfig);
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: {
                            Clay_BorderElementConfig *borderConfig = elementConfig->config.borderElementConfig;
                            CLAY({ .id = CLAY_ID("Clay__DebugViewElementInfoBorderBody"), .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
                                CLAY_TEXT(CLAY_STRING("Border Widths"), infoTitleConfig);
                                CLAY({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
                                    CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(borderConfig->width.left), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(borderConfig->width.right), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(borderConfig->width.top), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig);
                                    CLAY_TEXT(Clay__IntToString(borderConfig->width.bottom), infoTextConfig);
                                    CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
                                }
                                // .textColor
                                CLAY_TEXT(CLAY_STRING("Border Color"), infoTitleConfig);
                                Clay__RenderDebugViewColor(borderConfig->color, infoTextConfig);
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM:
                        default: break;
                    }
                }
            }
        } else {
            CLAY({ .id = CLAY_ID("Clay__DebugViewWarningsScrollPane"), .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(300)}, .childGap = 6, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_2, .clip = { .horizontal = true, .vertical = true, .childOffset = Clay_GetScrollOffset() } }) {
                Clay_TextElementConfig *warningConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
                CLAY({ .id = CLAY_ID("Clay__DebugViewWarningItemHeader"), .layout = { .sizing = {.height = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childGap = 8, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
                    CLAY_TEXT(CLAY_STRING("Warnings"), warningConfig);
                }
                CLAY({ .id = CLAY_ID("Clay__DebugViewWarningsTopBorder"), .layout = { .sizing = { .width = CLAY_SIZING_GROW(0), .height = CLAY_SIZING_FIXED(1)} }, .backgroundColor = {200, 200, 200, 255} }) {}
                int32_t previousWarningsLength = context->warnings.length;
                for (int32_t i = 0; i < previousWarningsLength; i++) {
                    Clay__Warning warning = context->warnings.internalArray[i];
                    CLAY({ .id = CLAY_IDI("Clay__DebugViewWarningItem", i), .layout = { .sizing = {.height = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childGap = 8, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
                        CLAY_TEXT(warning.baseMessage, warningConfig);
                        if (warning.dynamicMessage.length > 0) {
                            CLAY_TEXT(warning.dynamicMessage, warningConfig);
                        }
                    }
                }
            }
        }
    }
}
#pragma endregion

uint32_t Clay__debugViewWidth = 400;
Clay_Color Clay__debugViewHighlightColor = { 168, 66, 28, 100 };

Clay__WarningArray Clay__WarningArray_Allocate_Arena(int32_t capacity, Clay_Arena *arena) {
    size_t totalSizeBytes = capacity * sizeof(Clay_String);
    Clay__WarningArray array = {.capacity = capacity, .length = 0};
    uintptr_t nextAllocOffset = arena->nextAllocation + (64 - (arena->nextAllocation % 64));
    if (nextAllocOffset + totalSizeBytes <= arena->capacity) {
        array.internalArray = (Clay__Warning*)((uintptr_t)arena->memory + (uintptr_t)nextAllocOffset);
        arena->nextAllocation = nextAllocOffset + totalSizeBytes;
    }
    else {
        Clay__currentContext->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
            .errorType = CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
            .errorText = CLAY_STRING("Clay attempted to allocate memory in its arena, but ran out of capacity. Try increasing the capacity of the arena passed to Clay_Initialize()"),
            .userData = Clay__currentContext->errorHandler.userData });
    }
    return array;
}

Clay__Warning *Clay__WarningArray_Add(Clay__WarningArray *array, Clay__Warning item)
{
    if (array->length < array->capacity) {
        array->internalArray[array->length++] = item;
        return &array->internalArray[array->length - 1];
    }
    return &CLAY__WARNING_DEFAULT;
}

void* Clay__Array_Allocate_Arena(int32_t capacity, uint32_t itemSize, Clay_Arena *arena)
{
    size_t totalSizeBytes = capacity * itemSize;
    uintptr_t nextAllocOffset = arena->nextAllocation + ((64 - (arena->nextAllocation % 64)) & 63);
    if (nextAllocOffset + totalSizeBytes <= arena->capacity) {
        arena->nextAllocation = nextAllocOffset + totalSizeBytes;
        return (void*)((uintptr_t)arena->memory + (uintptr_t)nextAllocOffset);
    }
    else {
        Clay__currentContext->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                .errorType = CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
                .errorText = CLAY_STRING("Clay attempted to allocate memory in its arena, but ran out of capacity. Try increasing the capacity of the arena passed to Clay_Initialize()"),
                .userData = Clay__currentContext->errorHandler.userData });
    }
    return CLAY__NULL;
}

bool Clay__Array_RangeCheck(int32_t index, int32_t length)
{
    if (index < length && index >= 0) {
        return true;
    }
    Clay_Context* context = Clay_GetCurrentContext();
    context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
            .errorType = CLAY_ERROR_TYPE_INTERNAL_ERROR,
            .errorText = CLAY_STRING("Clay attempted to make an out of bounds array access. This is an internal error and is likely a bug."),
            .userData = context->errorHandler.userData });
    return false;
}

bool Clay__Array_AddCapacityCheck(int32_t length, int32_t capacity)
{
    if (length < capacity) {
        return true;
    }
    Clay_Context* context = Clay_GetCurrentContext();
    context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
        .errorType = CLAY_ERROR_TYPE_INTERNAL_ERROR,
        .errorText = CLAY_STRING("Clay attempted to make an out of bounds array access. This is an internal error and is likely a bug."),
        .userData = context->errorHandler.userData });
    return false;
}

// PUBLIC API FROM HERE ---------------------------------------

CLAY_WASM_EXPORT("Clay_MinMemorySize")
uint32_t Clay_MinMemorySize(void) {
    Clay_Context fakeContext = {
        .maxElementCount = Clay__defaultMaxElementCount,
        .maxMeasureTextCacheWordCount = Clay__defaultMaxMeasureTextWordCacheCount,
        .internalArena = {
            .capacity = SIZE_MAX,
            .memory = NULL,
        }
    };
    Clay_Context* currentContext = Clay_GetCurrentContext();
    if (currentContext) {
        fakeContext.maxElementCount = currentContext->maxElementCount;
        fakeContext.maxMeasureTextCacheWordCount = currentContext->maxMeasureTextCacheWordCount;
    }
    // Reserve space in the arena for the context, important for calculating min memory size correctly
    Clay__Context_Allocate_Arena(&fakeContext.internalArena);
    Clay__InitializePersistentMemory(&fakeContext);
    Clay__InitializeEphemeralMemory(&fakeContext);
    return (uint32_t)fakeContext.internalArena.nextAllocation + 128;
}

CLAY_WASM_EXPORT("Clay_CreateArenaWithCapacityAndMemory")
Clay_Arena Clay_CreateArenaWithCapacityAndMemory(size_t capacity, void *memory) {
    Clay_Arena arena = {
        .capacity = capacity,
        .memory = (char *)memory
    };
    return arena;
}

#ifndef CLAY_WASM
void Clay_SetMeasureTextFunction(Clay_Dimensions (*measureTextFunction)(Clay_StringSlice text, Clay_TextElementConfig *config, void *userData), void *userData) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__MeasureText = measureTextFunction;
    context->measureTextUserData = userData;
}
void Clay_SetQueryScrollOffsetFunction(Clay_Vector2 (*queryScrollOffsetFunction)(uint32_t elementId, void *userData), void *userData) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__QueryScrollOffset = queryScrollOffsetFunction;
    context->queryScrollOffsetUserData = userData;
}
#endif

CLAY_WASM_EXPORT("Clay_SetLayoutDimensions")
void Clay_SetLayoutDimensions(Clay_Dimensions dimensions) {
    Clay_GetCurrentContext()->layoutDimensions = dimensions;
}

CLAY_WASM_EXPORT("Clay_SetPointerState")
void Clay_SetPointerState(Clay_Vector2 position, bool isPointerDown) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return;
    }
    context->pointerInfo.position = position;
    context->pointerOverIds.length = 0;
    Clay__int32_tArray dfsBuffer = context->layoutElementChildrenBuffer;
    for (int32_t rootIndex = context->layoutElementTreeRoots.length - 1; rootIndex >= 0; --rootIndex) {
        dfsBuffer.length = 0;
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
        Clay__int32_tArray_Add(&dfsBuffer, (int32_t)root->layoutElementIndex);
        context->treeNodeVisited.internalArray[0] = false;
        bool found = false;
        while (dfsBuffer.length > 0) {
            if (context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
                dfsBuffer.length--;
                continue;
            }
            context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
            Clay_LayoutElement *currentElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&dfsBuffer, (int)dfsBuffer.length - 1));
            Clay_LayoutElementHashMapItem *mapItem = Clay__GetHashMapItem(currentElement->id); // TODO think of a way around this, maybe the fact that it's essentially a binary tree limits the cost, but the worst case is not great
            int32_t clipElementId = Clay__int32_tArray_GetValue(&context->layoutElementClipElementIds, (int32_t)(currentElement - context->layoutElements.internalArray));
            Clay_LayoutElementHashMapItem *clipItem = Clay__GetHashMapItem(clipElementId);
            if (mapItem) {
                Clay_BoundingBox elementBox = mapItem->boundingBox;
                elementBox.x -= root->pointerOffset.x;
                elementBox.y -= root->pointerOffset.y;
                if ((Clay__PointIsInsideRect(position, elementBox)) && (clipElementId == 0 || (Clay__PointIsInsideRect(position, clipItem->boundingBox)) || context->externalScrollHandlingEnabled)) {
                    if (mapItem->onHoverFunction) {
                        mapItem->onHoverFunction(mapItem->elementId, context->pointerInfo, mapItem->hoverFunctionUserData);
                    }
                    Clay_ElementIdArray_Add(&context->pointerOverIds, mapItem->elementId);
                    found = true;

                    if (mapItem->idAlias != 0) {
                        Clay_ElementIdArray_Add(&context->pointerOverIds, CLAY__INIT(Clay_ElementId) { .id = mapItem->idAlias });
                    }
                }
                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
                    dfsBuffer.length--;
                    continue;
                }
                for (int32_t i = currentElement->childrenOrTextContent.children.length - 1; i >= 0; --i) {
                    Clay__int32_tArray_Add(&dfsBuffer, currentElement->childrenOrTextContent.children.elements[i]);
                    context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = false; // TODO needs to be ranged checked
                }
            } else {
                dfsBuffer.length--;
            }
        }

        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context->layoutElements, root->layoutElementIndex);
        if (found && Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) &&
                Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig->pointerCaptureMode == CLAY_POINTER_CAPTURE_MODE_CAPTURE) {
            break;
        }
    }

    if (isPointerDown) {
        if (context->pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
            context->pointerInfo.state = CLAY_POINTER_DATA_PRESSED;
        } else if (context->pointerInfo.state != CLAY_POINTER_DATA_PRESSED) {
            context->pointerInfo.state = CLAY_POINTER_DATA_PRESSED_THIS_FRAME;
        }
    } else {
        if (context->pointerInfo.state == CLAY_POINTER_DATA_RELEASED_THIS_FRAME) {
            context->pointerInfo.state = CLAY_POINTER_DATA_RELEASED;
        } else if (context->pointerInfo.state != CLAY_POINTER_DATA_RELEASED)  {
            context->pointerInfo.state = CLAY_POINTER_DATA_RELEASED_THIS_FRAME;
        }
    }
}

CLAY_WASM_EXPORT("Clay_Initialize")
Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler) {
    // Cacheline align memory passed in
    uintptr_t baseOffset = 64 - ((uintptr_t)arena.memory % 64);
    baseOffset = baseOffset == 64 ? 0 : baseOffset;
    arena.memory += baseOffset;
    Clay_Context *context = Clay__Context_Allocate_Arena(&arena);
    if (context == NULL) return NULL;
    // DEFAULTS
    Clay_Context *oldContext = Clay_GetCurrentContext();
    *context = CLAY__INIT(Clay_Context) {
        .maxElementCount = oldContext ? oldContext->maxElementCount : Clay__defaultMaxElementCount,
        .maxMeasureTextCacheWordCount = oldContext ? oldContext->maxMeasureTextCacheWordCount : Clay__defaultMaxMeasureTextWordCacheCount,
        .errorHandler = errorHandler.errorHandlerFunction ? errorHandler : CLAY__INIT(Clay_ErrorHandler) { Clay__ErrorHandlerFunctionDefault, 0 },
        .layoutDimensions = layoutDimensions,
        .internalArena = arena,
    };
    Clay_SetCurrentContext(context);
    Clay__InitializePersistentMemory(context);
    Clay__InitializeEphemeralMemory(context);
    for (int32_t i = 0; i < context->layoutElementsHashMap.capacity; ++i) {
        context->layoutElementsHashMap.internalArray[i] = -1;
    }
    for (int32_t i = 0; i < context->measureTextHashMap.capacity; ++i) {
        context->measureTextHashMap.internalArray[i] = 0;
    }
    context->measureTextHashMapInternal.length = 1; // Reserve the 0 value to mean "no next element"
    context->layoutDimensions = layoutDimensions;
    return context;
}

CLAY_WASM_EXPORT("Clay_GetCurrentContext")
Clay_Context* Clay_GetCurrentContext(void) {
    return Clay__currentContext;
}

CLAY_WASM_EXPORT("Clay_SetCurrentContext")
void Clay_SetCurrentContext(Clay_Context* context) {
    Clay__currentContext = context;
}

CLAY_WASM_EXPORT("Clay_GetScrollOffset")
Clay_Vector2 Clay_GetScrollOffset(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    // If the element has no id attached at this point, we need to generate one
    if (openLayoutElement->id == 0) {
        Clay__GenerateIdForAnonymousElement(openLayoutElement);
    }
    for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
        if (mapping->layoutElement == openLayoutElement) {
            return mapping->scrollPosition;
        }
    }
    return CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
}

CLAY_WASM_EXPORT("Clay_UpdateScrollContainers")
void Clay_UpdateScrollContainers(bool enableDragScrolling, Clay_Vector2 scrollDelta, float deltaTime) {
    Clay_Context* context = Clay_GetCurrentContext();
    bool isPointerActive = enableDragScrolling && (context->pointerInfo.state == CLAY_POINTER_DATA_PRESSED || context->pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME);
    // Don't apply scroll events to ancestors of the inner element
    int32_t highestPriorityElementIndex = -1;
    Clay__ScrollContainerDataInternal *highestPriorityScrollData = CLAY__NULL;
    for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
        Clay__ScrollContainerDataInternal *scrollData = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
        if (!scrollData->openThisFrame) {
            Clay__ScrollContainerDataInternalArray_RemoveSwapback(&context->scrollContainerDatas, i);
            continue;
        }
        scrollData->openThisFrame = false;
        Clay_LayoutElementHashMapItem *hashMapItem = Clay__GetHashMapItem(scrollData->elementId);
        // Element isn't rendered this frame but scroll offset has been retained
        if (!hashMapItem) {
            Clay__ScrollContainerDataInternalArray_RemoveSwapback(&context->scrollContainerDatas, i);
            continue;
        }

        // Touch / click is released
        if (!isPointerActive && scrollData->pointerScrollActive) {
            float xDiff = scrollData->scrollPosition.x - scrollData->scrollOrigin.x;
            if (xDiff < -10 || xDiff > 10) {
                scrollData->scrollMomentum.x = (scrollData->scrollPosition.x - scrollData->scrollOrigin.x) / (scrollData->momentumTime * 25);
            }
            float yDiff = scrollData->scrollPosition.y - scrollData->scrollOrigin.y;
            if (yDiff < -10 || yDiff > 10) {
                scrollData->scrollMomentum.y = (scrollData->scrollPosition.y - scrollData->scrollOrigin.y) / (scrollData->momentumTime * 25);
            }
            scrollData->pointerScrollActive = false;

            scrollData->pointerOrigin = CLAY__INIT(Clay_Vector2){0,0};
            scrollData->scrollOrigin = CLAY__INIT(Clay_Vector2){0,0};
            scrollData->momentumTime = 0;
        }

        // Apply existing momentum
        scrollData->scrollPosition.x += scrollData->scrollMomentum.x;
        scrollData->scrollMomentum.x *= 0.95f;
        bool scrollOccurred = scrollDelta.x != 0 || scrollDelta.y != 0;
        if ((scrollData->scrollMomentum.x > -0.1f && scrollData->scrollMomentum.x < 0.1f) || scrollOccurred) {
            scrollData->scrollMomentum.x = 0;
        }
        scrollData->scrollPosition.x = CLAY__MIN(CLAY__MAX(scrollData->scrollPosition.x, -(CLAY__MAX(scrollData->contentSize.width - scrollData->layoutElement->dimensions.width, 0))), 0);

        scrollData->scrollPosition.y += scrollData->scrollMomentum.y;
        scrollData->scrollMomentum.y *= 0.95f;
        if ((scrollData->scrollMomentum.y > -0.1f && scrollData->scrollMomentum.y < 0.1f) || scrollOccurred) {
            scrollData->scrollMomentum.y = 0;
        }
        scrollData->scrollPosition.y = CLAY__MIN(CLAY__MAX(scrollData->scrollPosition.y, -(CLAY__MAX(scrollData->contentSize.height - scrollData->layoutElement->dimensions.height, 0))), 0);

        for (int32_t j = 0; j < context->pointerOverIds.length; ++j) { // TODO n & m are small here but this being n*m gives me the creeps
            if (scrollData->layoutElement->id == Clay_ElementIdArray_Get(&context->pointerOverIds, j)->id) {
                highestPriorityElementIndex = j;
                highestPriorityScrollData = scrollData;
            }
        }
    }

    if (highestPriorityElementIndex > -1 && highestPriorityScrollData) {
        Clay_LayoutElement *scrollElement = highestPriorityScrollData->layoutElement;
        Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(scrollElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
        bool canScrollVertically = clipConfig->vertical && highestPriorityScrollData->contentSize.height > scrollElement->dimensions.height;
        bool canScrollHorizontally = clipConfig->horizontal && highestPriorityScrollData->contentSize.width > scrollElement->dimensions.width;
        // Handle wheel scroll
        if (canScrollVertically) {
            highestPriorityScrollData->scrollPosition.y = highestPriorityScrollData->scrollPosition.y + scrollDelta.y * 10;
        }
        if (canScrollHorizontally) {
            highestPriorityScrollData->scrollPosition.x = highestPriorityScrollData->scrollPosition.x + scrollDelta.x * 10;
        }
        // Handle click / touch scroll
        if (isPointerActive) {
            highestPriorityScrollData->scrollMomentum = CLAY__INIT(Clay_Vector2)CLAY__DEFAULT_STRUCT;
            if (!highestPriorityScrollData->pointerScrollActive) {
                highestPriorityScrollData->pointerOrigin = context->pointerInfo.position;
                highestPriorityScrollData->scrollOrigin = highestPriorityScrollData->scrollPosition;
                highestPriorityScrollData->pointerScrollActive = true;
            } else {
                float scrollDeltaX = 0, scrollDeltaY = 0;
                if (canScrollHorizontally) {
                    float oldXScrollPosition = highestPriorityScrollData->scrollPosition.x;
                    highestPriorityScrollData->scrollPosition.x = highestPriorityScrollData->scrollOrigin.x + (context->pointerInfo.position.x - highestPriorityScrollData->pointerOrigin.x);
                    highestPriorityScrollData->scrollPosition.x = CLAY__MAX(CLAY__MIN(highestPriorityScrollData->scrollPosition.x, 0), -(highestPriorityScrollData->contentSize.width - highestPriorityScrollData->boundingBox.width));
                    scrollDeltaX = highestPriorityScrollData->scrollPosition.x - oldXScrollPosition;
                }
                if (canScrollVertically) {
                    float oldYScrollPosition = highestPriorityScrollData->scrollPosition.y;
                    highestPriorityScrollData->scrollPosition.y = highestPriorityScrollData->scrollOrigin.y + (context->pointerInfo.position.y - highestPriorityScrollData->pointerOrigin.y);
                    highestPriorityScrollData->scrollPosition.y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData->scrollPosition.y, 0), -(highestPriorityScrollData->contentSize.height - highestPriorityScrollData->boundingBox.height));
                    scrollDeltaY = highestPriorityScrollData->scrollPosition.y - oldYScrollPosition;
                }
                if (scrollDeltaX > -0.1f && scrollDeltaX < 0.1f && scrollDeltaY > -0.1f && scrollDeltaY < 0.1f && highestPriorityScrollData->momentumTime > 0.15f) {
                    highestPriorityScrollData->momentumTime = 0;
                    highestPriorityScrollData->pointerOrigin = context->pointerInfo.position;
                    highestPriorityScrollData->scrollOrigin = highestPriorityScrollData->scrollPosition;
                } else {
                     highestPriorityScrollData->momentumTime += deltaTime;
                }
            }
        }
        // Clamp any changes to scroll position to the maximum size of the contents
        if (canScrollVertically) {
            highestPriorityScrollData->scrollPosition.y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData->scrollPosition.y, 0), -(highestPriorityScrollData->contentSize.height - scrollElement->dimensions.height));
        }
        if (canScrollHorizontally) {
            highestPriorityScrollData->scrollPosition.x = CLAY__MAX(CLAY__MIN(highestPriorityScrollData->scrollPosition.x, 0), -(highestPriorityScrollData->contentSize.width - scrollElement->dimensions.width));
        }
    }
}

CLAY_WASM_EXPORT("Clay_BeginLayout")
void Clay_BeginLayout(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__InitializeEphemeralMemory(context);
    context->generation++;
    context->dynamicElementIndex = 0;
    // Set up the root container that covers the entire window
    Clay_Dimensions rootDimensions = {context->layoutDimensions.width, context->layoutDimensions.height};
    if (context->debugModeEnabled) {
        rootDimensions.width -= (float)Clay__debugViewWidth;
    }
    context->booleanWarnings = CLAY__INIT(Clay_BooleanWarnings) CLAY__DEFAULT_STRUCT;
    Clay__OpenElement();
    Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) {
            .id = CLAY_ID("Clay__RootContainer"),
            .layout = { .sizing = {CLAY_SIZING_FIXED((rootDimensions.width)), CLAY_SIZING_FIXED(rootDimensions.height)} }
    });
    Clay__int32_tArray_Add(&context->openLayoutElementStack, 0);
    Clay__LayoutElementTreeRootArray_Add(&context->layoutElementTreeRoots, CLAY__INIT(Clay__LayoutElementTreeRoot) { .layoutElementIndex = 0 });
}

CLAY_WASM_EXPORT("Clay_EndLayout")
Clay_RenderCommandArray Clay_EndLayout(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__CloseElement();
    bool elementsExceededBeforeDebugView = context->booleanWarnings.maxElementsExceeded;
    if (context->debugModeEnabled && !elementsExceededBeforeDebugView) {
        context->warningsEnabled = false;
        Clay__RenderDebugView();
        context->warningsEnabled = true;
    }
    if (context->booleanWarnings.maxElementsExceeded) {
        Clay_String message;
        if (!elementsExceededBeforeDebugView) {
            message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount after adding the debug-view to the layout.");
        } else {
            message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount");
        }
        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand ) {
            .boundingBox = { context->layoutDimensions.width / 2 - 59 * 4, context->layoutDimensions.height / 2, 0, 0 },
            .renderData = { .text = { .stringContents = CLAY__INIT(Clay_StringSlice) { .length = message.length, .chars = message.chars, .baseChars = message.chars }, .textColor = {255, 0, 0, 255}, .fontSize = 16 } },
            .commandType = CLAY_RENDER_COMMAND_TYPE_TEXT
        });
    } else {
        Clay__CalculateFinalLayout();
    }
    return context->renderCommands;
}

CLAY_WASM_EXPORT("Clay_GetElementId")
Clay_ElementId Clay_GetElementId(Clay_String idString) {
    return Clay__HashString(idString, 0);
}

CLAY_WASM_EXPORT("Clay_GetElementIdWithIndex")
Clay_ElementId Clay_GetElementIdWithIndex(Clay_String idString, uint32_t index) {
    return Clay__HashStringWithOffset(idString, index, 0);
}

bool Clay_Hovered(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return false;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    // If the element has no id attached at this point, we need to generate one
    if (openLayoutElement->id == 0) {
        Clay__GenerateIdForAnonymousElement(openLayoutElement);
    }
    for (int32_t i = 0; i < context->pointerOverIds.length; ++i) {
        if (Clay_ElementIdArray_Get(&context->pointerOverIds, i)->id == openLayoutElement->id) {
            return true;
        }
    }
    return false;
}

void Clay_OnHover(void (*onHoverFunction)(Clay_ElementId elementId, Clay_PointerData pointerInfo, intptr_t userData), intptr_t userData) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    if (openLayoutElement->id == 0) {
        Clay__GenerateIdForAnonymousElement(openLayoutElement);
    }
    Clay_LayoutElementHashMapItem *hashMapItem = Clay__GetHashMapItem(openLayoutElement->id);
    hashMapItem->onHoverFunction = onHoverFunction;
    hashMapItem->hoverFunctionUserData = userData;
}

CLAY_WASM_EXPORT("Clay_PointerOver")
bool Clay_PointerOver(Clay_ElementId elementId) { // TODO return priority for separating multiple results
    Clay_Context* context = Clay_GetCurrentContext();
    for (int32_t i = 0; i < context->pointerOverIds.length; ++i) {
        if (Clay_ElementIdArray_Get(&context->pointerOverIds, i)->id == elementId.id) {
            return true;
        }
    }
    return false;
}

CLAY_WASM_EXPORT("Clay_GetScrollContainerData")
Clay_ScrollContainerData Clay_GetScrollContainerData(Clay_ElementId id) {
    Clay_Context* context = Clay_GetCurrentContext();
    for (int32_t i = 0; i < context->scrollContainerDatas.length; ++i) {
        Clay__ScrollContainerDataInternal *scrollContainerData = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
        if (scrollContainerData->elementId == id.id) {
            Clay_ClipElementConfig *clipElementConfig = Clay__FindElementConfigWithType(scrollContainerData->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
            if (!clipElementConfig) { // This can happen on the first frame before a scroll container is declared
                return CLAY__INIT(Clay_ScrollContainerData) CLAY__DEFAULT_STRUCT;
            }
            return CLAY__INIT(Clay_ScrollContainerData) {
                .scrollPosition = &scrollContainerData->scrollPosition,
                .scrollContainerDimensions = { scrollContainerData->boundingBox.width, scrollContainerData->boundingBox.height },
                .contentDimensions = scrollContainerData->contentSize,
                .config = *clipElementConfig,
                .found = true
            };
        }
    }
    return CLAY__INIT(Clay_ScrollContainerData) CLAY__DEFAULT_STRUCT;
}

CLAY_WASM_EXPORT("Clay_GetElementData")
Clay_ElementData Clay_GetElementData(Clay_ElementId id){
    Clay_LayoutElementHashMapItem * item = Clay__GetHashMapItem(id.id);
    if(item == &Clay_LayoutElementHashMapItem_DEFAULT) {
        return CLAY__INIT(Clay_ElementData) CLAY__DEFAULT_STRUCT;
    }

    return CLAY__INIT(Clay_ElementData){
        .boundingBox = item->boundingBox,
        .found = true
    };
}

CLAY_WASM_EXPORT("Clay_SetDebugModeEnabled")
void Clay_SetDebugModeEnabled(bool enabled) {
    Clay_Context* context = Clay_GetCurrentContext();
    context->debugModeEnabled = enabled;
}

CLAY_WASM_EXPORT("Clay_IsDebugModeEnabled")
bool Clay_IsDebugModeEnabled(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    return context->debugModeEnabled;
}

CLAY_WASM_EXPORT("Clay_SetCullingEnabled")
void Clay_SetCullingEnabled(bool enabled) {
    Clay_Context* context = Clay_GetCurrentContext();
    context->disableCulling = !enabled;
}

CLAY_WASM_EXPORT("Clay_SetExternalScrollHandlingEnabled")
void Clay_SetExternalScrollHandlingEnabled(bool enabled) {
    Clay_Context* context = Clay_GetCurrentContext();
    context->externalScrollHandlingEnabled = enabled;
}

CLAY_WASM_EXPORT("Clay_GetMaxElementCount")
int32_t Clay_GetMaxElementCount(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    return context->maxElementCount;
}

CLAY_WASM_EXPORT("Clay_SetMaxElementCount")
void Clay_SetMaxElementCount(int32_t maxElementCount) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context) {
        context->maxElementCount = maxElementCount;
    } else {
        Clay__defaultMaxElementCount = maxElementCount; // TODO: Fix this
        Clay__defaultMaxMeasureTextWordCacheCount = maxElementCount * 2;
    }
}

CLAY_WASM_EXPORT("Clay_GetMaxMeasureTextCacheWordCount")
int32_t Clay_GetMaxMeasureTextCacheWordCount(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    return context->maxMeasureTextCacheWordCount;
}

CLAY_WASM_EXPORT("Clay_SetMaxMeasureTextCacheWordCount")
void Clay_SetMaxMeasureTextCacheWordCount(int32_t maxMeasureTextCacheWordCount) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context) {
        Clay__currentContext->maxMeasureTextCacheWordCount = maxMeasureTextCacheWordCount;
    } else {
        Clay__defaultMaxMeasureTextWordCacheCount = maxMeasureTextCacheWordCount; // TODO: Fix this
    }
}

CLAY_WASM_EXPORT("Clay_ResetMeasureTextCache")
void Clay_ResetMeasureTextCache(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    context->measureTextHashMapInternal.length = 0;
    context->measureTextHashMapInternalFreeList.length = 0;
    context->measureTextHashMap.length = 0;
    context->measuredWords.length = 0;
    context->measuredWordsFreeList.length = 0;
    
    for (int32_t i = 0; i < context->measureTextHashMap.capacity; ++i) {
        context->measureTextHashMap.internalArray[i] = 0;
    }
    context->measureTextHashMapInternal.length = 1; // Reserve the 0 value to mean "no next element"
}

#endif // CLAY_IMPLEMENTATION

/*
LICENSE
zlib/libpng license

Copyright (c) 2024 Nic Barker

This software is provided 'as-is', without any express or implied warranty.
In no event will the authors be held liable for any damages arising from the
use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

    1. The origin of this software must not be misrepresented; you must not
    claim that you wrote the original software. If you use this software in a
    product, an acknowledgment in the product documentation would be
    appreciated but is not required.

    2. Altered source versions must be plainly marked as such, and must not
    be misrepresented as being the original software.

    3. This notice may not be removed or altered from any source
    distribution.
*/
