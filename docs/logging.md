# Maya Framework Logging System

## Overview

Maya Framework includes a configurable logging system that allows you to control the verbosity of logs during development and debugging.

## Log Levels

The framework supports the following log levels (from least to most verbose):

1. **silent** (default) - No logs output
2. **error** - Only error messages
3. **warn** - Warnings and errors
4. **info** - Important information, warnings, and errors
5. **debug** - Debug messages and all above
6. **trace** - Detailed trace logs including all reactive system operations

## Configuration

### Environment Variables

Control logging through environment variables:

```bash
# Set log level
export MAYA_LOG_LEVEL=debug

# Set specific categories to log (comma-separated)
export MAYA_LOG_CATEGORIES=SIGNAL,EFFECT,MEMO
```

### Available Categories

- **APP** - Application lifecycle events
- **SIGNAL** - Signal read/write operations
- **EFFECT** - Effect execution and dependencies
- **MEMO** - Memoized computation caching
- **COMPUTED** - Computed value derivations
- **UI** - User interface events (clicks, etc.)
- **PIPELINE** - Rendering pipeline stages
- **RENDER** - Rendering operations
- **UPDATE** - DOM updates

## Usage Examples

### Build with Different Log Levels

```bash
# No logs (default)
./build.sh

# Show only errors
MAYA_LOG_LEVEL=error ./build.sh

# Show info messages
MAYA_LOG_LEVEL=info ./build.sh

# Debug mode
MAYA_LOG_LEVEL=debug ./build.sh

# Full trace (very verbose)
MAYA_LOG_LEVEL=trace ./build.sh
```

### Filter by Category

```bash
# Show only Signal and Effect logs
MAYA_LOG_LEVEL=trace MAYA_LOG_CATEGORIES=SIGNAL,EFFECT ./build.sh

# Debug only UI interactions
MAYA_LOG_LEVEL=debug MAYA_LOG_CATEGORIES=UI ./build.sh

# Track rendering pipeline
MAYA_LOG_LEVEL=debug MAYA_LOG_CATEGORIES=PIPELINE,RENDER ./build.sh
```

## In Your Code

You can also control logging programmatically:

```go
import "github.com/maya-framework/maya/internal/logger"

// Set log level
logger.SetLevel(logger.LevelDebug)

// Enable specific categories
logger.EnableCategory("SIGNAL")
logger.EnableCategory("EFFECT")

// Disable a category
logger.DisableCategory("MEMO")

// Use in your code
logger.Debug("UI", "Button clicked: %s", buttonName)
logger.Info("APP", "Application started")
logger.Error("RENDER", "Failed to render: %v", err)
```

## Performance Considerations

- **Production**: Always use `silent` or `error` level
- **Development**: Use `debug` for general development
- **Debugging**: Use `trace` only when investigating specific issues
- **Category Filtering**: Use categories to reduce noise when debugging specific features

## Troubleshooting

If you're seeing too many logs:
1. Check your MAYA_LOG_LEVEL environment variable
2. Use MAYA_LOG_CATEGORIES to filter to specific systems
3. Ensure you're not running a debug build in production

If you're not seeing expected logs:
1. Verify the log level is set appropriately
2. Check that the category is enabled (if using categories)
3. Ensure the build was compiled with the correct environment variables