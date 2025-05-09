# Utility Package

This package contains utility functions for common tasks across the application.

## Mask Utility

The `mask.go` module provides functions for masking sensitive data in logging and error messages.

### Available Functions

- `MaskDsn(dsn string)`: Masks password in database connection strings
- `MaskSensitive(value string, visiblePrefixChars, visibleSuffixChars int)`: General masking function
- `MaskCredential(credential string)`: Masks credentials showing first/last 2 chars
- `MaskEmail(email string)`: Masks email addresses (local part only)
- `MaskJWT(token string)`: Masks JWT tokens while preserving structure
- `MaskURL(url string)`: Masks sensitive parts of URLs (auth, tokens, etc.)

### Usage Examples

```go
package main

import (
    "fmt"
    "github.com/linkeunid/go-api/pkg/util"
)

func main() {
    // Mask a database connection string
    dsn := "user:password123@tcp(localhost:3306)/mydb"
    fmt.Println(util.MaskDsn(dsn)) 
    // Output: user:******@tcp(localhost:3306)/mydb

    // Mask a credential like an API key
    apiKey := "abcdef1234567890"
    fmt.Println(util.MaskCredential(apiKey))
    // Output: ab************90

    // Mask an email address
    email := "user@example.com"
    fmt.Println(util.MaskEmail(email))
    // Output: us**@example.com

    // Mask a URL with sensitive params
    url := "https://api.example.com/v1?api_key=secret&user=bob"
    fmt.Println(util.MaskURL(url))
    // Output: https://api.example.com/v1?api_key=******&user=bob

    // Custom masking with specific visible parts
    customData := "sensitive-information-here"
    fmt.Println(util.MaskSensitive(customData, 3, 4))
    // Output: sen**************here
}
```

### Integration with Logger

The mask utility integrates well with structured loggers like zap:

```go
import (
    "github.com/linkeunid/go-api/pkg/util"
    "go.uber.org/zap"
)

func logExample(logger *zap.Logger) {
    apiKey := "abcdef1234567890"
    
    // Log with masked sensitive data
    logger.Info("API request",
        zap.String("endpoint", "/api/users"),
        zap.String("api_key", util.MaskCredential(apiKey)),
    )
}
```

### Additional Helper Function

A general `MaskSensitiveData` function is available in the bootstrap package that automatically selects the appropriate masking function based on data type:

```go
import "github.com/linkeunid/go-api/internal/bootstrap"

// Examples
bootstrap.MaskSensitiveData("email", "user@example.com")  // Masks email
bootstrap.MaskSensitiveData("password", "secret123")      // Masks password
bootstrap.MaskSensitiveData("token", "eyJhbGciOi...")     // Masks JWT token
bootstrap.MaskSensitiveData("url", "https://user:pass@...") // Masks URL
``` 