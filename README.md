# GUtils

A comprehensive collection of production-ready Go utilities for building modern applications. GUtils provides type-safe abstractions for common tasks including database operations, caching, authentication, messaging, and more.

## Features

- **Database Abstraction** - Generic CRUD operations for SQL (PostgreSQL, MySQL, SQLite) and NoSQL (MongoDB)
- **Caching** - Pluggable cache with in-memory and Redis backends
- **Authentication** - JWT token generation and validation
- **Message Broker** - RabbitMQ pub/sub with delayed message support
- **HTTP Client** - Unified interface with multiple backend implementations
- **USSD Framework** - Session management for USSD menu navigation
- **Image Processing** - Compression and format conversion
- **Error Handling** - Structured errors with HTTP status codes
- **Type-Safe Utilities** - Environment variables, logging, UUID generation, and more

## Installation

```bash
go get github.com/ochom/gutils
```

## Quick Start

```go
package main

import (
    "github.com/ochom/gutils/env"
    "github.com/ochom/gutils/logs"
    "github.com/ochom/gutils/sqlr"
)

func main() {
    // Load environment variables
    dbUrl := env.MustGet("DATABASE_URL")

    // Initialize database
    sqlr.Init(sqlr.Config{Url: dbUrl})

    // Use structured logging
    logs.Info("Application started successfully")
}
```

## Package Documentation

### 1. sqlr - SQL Database Abstraction

Generic GORM wrapper supporting PostgreSQL, MySQL, and SQLite with type-safe CRUD operations.

#### Installation
```go
import "github.com/ochom/gutils/sqlr"
```

#### Configuration
```go
config := sqlr.Config{
    Url:      "postgresql://user:pass@localhost:5432/db",
    LogLevel: 4,  // GORM log level (1=Silent, 2=Error, 3=Warn, 4=Info)
    MaxOpenConns: 10,
    MaxIdleConns: 5,
    MaxConnIdleTime: 5 * time.Minute,
    MaxConnLifeTime: 1 * time.Hour,
}

sqlr.Init(config)
```

#### Usage Examples
```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Email string
}

// Create
user := User{Name: "John Doe", Email: "john@example.com"}
if err := sqlr.Create(&user); err != nil {
    logs.Error("Failed to create user: %v", err)
}

// Find by ID
user, err := sqlr.FindOneById[User](1)

// Find with conditions
db := sqlr.GORM()
users, err := sqlr.FindAll[User](db.Where("email LIKE ?", "%@example.com"))

// Update entire record
user.Name = "Jane Doe"
sqlr.Update(&user)

// Update specific fields
sqlr.UpdateOne[User](
    db.Where("id = ?", 1),
    map[string]interface{}{"name": "New Name"},
)

// Delete
sqlr.DeleteById[User](1)

// Count
count, _ := sqlr.Count[User](db.Where("email LIKE ?", "%@example.com"))

// Pagination
users, _ := sqlr.FindWithLimit[User](db, 1, 20) // page 1, 20 per page
```

### 2. nosql - MongoDB Abstraction

Generic MongoDB operations with automatic collection naming.

#### Installation
```go
import "github.com/ochom/gutils/nosql"
```

#### Configuration
```go
nosql.Init("mongodb://localhost:27017", "mydb")
```

#### Usage Examples
```go
type Product struct {
    ID    string `bson:"_id"`
    Name  string `bson:"name"`
    Price float64 `bson:"price"`
}

// Create
product := Product{ID: uuid.New(), Name: "Laptop", Price: 999.99}
nosql.Create(&product)

// Find one
filter := bson.M{"name": "Laptop"}
product, err := nosql.FindOne[Product](filter)

// Find all
products, _ := nosql.FindAll[Product](bson.M{"price": bson.M{"$gt": 500}})

// Update
product.Price = 899.99
nosql.Update(&product)

// Aggregation
pipeline := []bson.M{
    {"$match": bson.M{"price": bson.M{"$gt": 100}}},
    {"$group": bson.M{"_id": "$category", "total": bson.M{"$sum": 1}}},
}
results, _ := nosql.Pipe[Product](pipeline)

// Delete
nosql.DeleteByID[Product]("product-id")
```

### 3. cache - Caching Layer

Pluggable caching with in-memory and Redis support.

#### Installation
```go
import "github.com/ochom/gutils/cache"
```

#### Configuration
Set environment variable:
```bash
CACHE_DRIVER=0  # 0 = Memory, 1 = Redis
REDIS_URL=redis://localhost:6379  # Required if using Redis
```

#### Usage Examples
```go
import "time"

// Set with expiration
cache.Set("user:1", userData, 5*time.Minute)

// Get
if data, found := cache.Get("user:1"); found {
    // Use cached data
}

// Delete
cache.Delete("user:1")

// Access Redis client directly (when using Redis driver)
client := cache.Client()
```

### 4. auth - JWT Authentication

Generate and validate JWT tokens for authentication.

#### Installation
```go
import "github.com/ochom/gutils/auth"
```

#### Configuration
Set environment variable:
```bash
JWT_SECRET=your-secret-key
```

#### Usage Examples
```go
// Generate tokens
userData := map[string]string{
    "user_id": "123",
    "email": "user@example.com",
    "role": "admin",
}

tokens, err := auth.GenerateAuthTokens(userData, 0, 0)
// tokens["access_token"]  - expires in 3 hours (default)
// tokens["refresh_token"] - expires in 7 days (default)

// Custom expiry
tokens, _ := auth.GenerateAuthTokens(userData, 1*time.Hour, 30*24*time.Hour)

// Validate and extract claims
claims, err := auth.GetAuthClaims(tokenString)
if err != nil {
    logs.Error("Invalid token: %v", err)
}

userID := claims["user_id"]
role := claims["role"]
```

### 5. gttp - HTTP Client

Unified HTTP client interface with multiple backends.

#### Installation
```go
import "github.com/ochom/gutils/gttp"
```

#### Configuration
```bash
HTTP_CLIENT=0  # 0 = Default Go HTTP, 1 = GoFiber
```

#### Usage Examples
```go
// GET request
headers := gttp.M{"Authorization": "Bearer token"}
response := gttp.Get("https://api.example.com/users", headers, 10*time.Second)

if response.StatusCode == 200 {
    var users []User
    json.Unmarshal(response.Body, &users)
}

// POST request
payload := map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
}
body, _ := json.Marshal(payload)

headers := gttp.M{
    "Content-Type": "application/json",
    "Authorization": "Bearer token",
}
response := gttp.Post("https://api.example.com/users", headers, body, 10*time.Second)

// Generic request
response := gttp.SendRequest(
    "https://api.example.com/data",
    "PUT",
    headers,
    body,
    15*time.Second,
)
```

### 6. pubsub - RabbitMQ Messaging

Publish/subscribe messaging with delayed message support.

#### Installation
```go
import "github.com/ochom/gutils/pubsub"
```

#### Publisher Example
```go
publisher := pubsub.NewPublisher(
    "amqp://guest:guest@localhost:5672/",
    "my-exchange",
    "my-queue",
)

publisher.SetExchangeType(pubsub.Direct)
publisher.SetRoutingKey("notifications")

// Publish immediately
message := map[string]string{"user_id": "123", "event": "signup"}
publisher.Publish(message)

// Publish with delay
publisher.PublishWithDelay(message, 5*time.Minute)
```

#### Consumer Example
```go
consumer := pubsub.NewConsumer(
    "amqp://guest:guest@localhost:5672/",
    "my-queue",
)

consumer.Consume(func(body []byte) {
    var message map[string]string
    json.Unmarshal(body, &message)

    // Process message
    logs.Info("Received: %v", message)
})
```

### 7. ussd - USSD Session Management

Framework for building USSD menu-driven applications.

#### Installation
```go
import "github.com/ochom/gutils/ussd"
```

#### Usage Example
```go
// Define root menu
rootMenu := ussd.NewStep("Welcome to MyApp\n1. Check Balance\n2. Send Money\n3. Exit")

// Add sub-menus
balanceMenu := ussd.NewStep(func(p ussd.Params) string {
    // Fetch balance from database
    return fmt.Sprintf("Your balance is: KES 1,000")
})
rootMenu.AddStep("1", balanceMenu)

// Send money flow
sendMenu := ussd.NewStep("Enter amount:")
sendMenu.Run = func(p ussd.Params) string {
    amount := p.Text
    // Process transaction
    return fmt.Sprintf("Sent KES %s successfully", amount)
}
rootMenu.AddStep("2", sendMenu)

// Initialize parser
parser := ussd.New(rootMenu)

// Handle USSD request
params := ussd.Params{
    SessionId:   "session-123",
    PhoneNumber: "254712345678",
    Text:        "1*500",  // User selected 1, then entered 500
}

response := parser.Parse(params)
// response contains CON/END prefix and menu text
```

### 8. env - Environment Variables

Type-safe environment variable access with defaults.

#### Installation
```go
import "github.com/ochom/gutils/env"
```

#### Usage Examples
```go
// With defaults (safe)
dbHost := env.Get("DB_HOST", "localhost")
dbPort := env.Int("DB_PORT", 5432)
debugMode := env.Bool("DEBUG", false)
timeout := env.Float("TIMEOUT", 30.0)

// Must exist (panics if missing)
apiKey := env.MustGet("API_KEY")
maxConns := env.MustInt("MAX_CONNECTIONS")
enableCache := env.MustBool("ENABLE_CACHE")
```

### 9. logs - Structured Logging

Color-coded logging with file/line information.

#### Installation
```go
import "github.com/ochom/gutils/logs"
```

#### Usage Examples
```go
logs.Debug("User input: %v", userInput)          // Blue
logs.Info("Server started on port %d", 8080)     // Green
logs.Warn("High memory usage: %.2f%%", 85.5)     // Yellow
logs.Error("Database connection failed: %v", err) // Red
logs.Fatal("Critical error: %v", err)            // Red + exit(1)

// Redirect output
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
logs.SetOutput(file)
```

### 10. errors - Custom Error Types

Structured errors with HTTP status codes.

#### Installation
```go
import "github.com/ochom/gutils/errors"
```

#### Usage Examples
```go
// Create errors with status codes
err := errors.NotFound("User with id %d not found", userID)
err := errors.BadRequest("Invalid email format: %s", email)
err := errors.Unauthorized("Invalid credentials")
err := errors.Forbidden("Access denied to resource %s", resourceID)
err := errors.Conflict("Email %s already exists", email)
err := errors.New("Unexpected error: %v", originalErr)

// Use in HTTP handlers
if err != nil {
    customErr := err.(*errors.CustomError)
    return c.Status(customErr.Code).JSON(fiber.Map{
        "error": customErr.Message,
    })
}
```

### 11. uuid - Unique ID Generation

Generate MongoDB-style object IDs.

#### Installation
```go
import "github.com/ochom/gutils/uuid"
```

#### Usage Example
```go
id := uuid.New()  // Returns hex string like "507f1f77bcf86cd799439011"

user := User{
    ID:   uuid.New(),
    Name: "John Doe",
}
```

### 12. images - Image Processing

Compress and convert images to JPEG format.

#### Installation
```go
import "github.com/ochom/gutils/images"
```

#### Usage Example
```go
// Read image file
imageData, _ := ioutil.ReadFile("photo.png")

// Compress and convert to JPEG
compressedData, newFilename, err := images.CompressImage(imageData, 85, "photo.png")
if err != nil {
    logs.Error("Compression failed: %v", err)
}

// Save compressed image
ioutil.WriteFile(newFilename, compressedData, 0644)
```

### 13. arrays - Array Utilities

JavaScript-style array operations using Go generics.

#### Installation
```go
import "github.com/ochom/gutils/arrays"
```

#### Usage Examples
```go
numbers := []int{1, 2, 3, 4, 5, 6}

// Filter
evens := arrays.Filter(numbers, func(n int) bool {
    return n%2 == 0
}) // [2, 4, 6]

// Map
doubled := arrays.Map(numbers, func(n int) int {
    return n * 2
}) // [2, 4, 6, 8, 10, 12]

// Find
first := arrays.Find(numbers, func(n int) bool {
    return n > 3
}) // 4

// Reduce
sum := arrays.Reduce(numbers, 0, func(acc, n int) int {
    return acc + n
}) // 21

// Chunk
chunks := arrays.Chunk(numbers, 2) // [[1,2], [3,4], [5,6]]

// ForEach
arrays.ForEach(numbers, func(n int) {
    fmt.Println(n)
})
```

### 14. helpers - Utility Functions

Collection of helper functions for common tasks.

#### Installation
```go
import "github.com/ochom/gutils/helpers"
```

#### Usage Examples
```go
// Password hashing
hash, _ := helpers.HashPassword("mypassword")
isValid := helpers.ComparePassword(hash, "mypassword")

// Generate OTP
otp := helpers.GenerateOTP(6) // Returns 6-digit numeric OTP

// Phone number parsing (Kenya format)
normalized, err := helpers.ParseMobile("0712345678")
// Returns "254712345678"

hashedPhone := helpers.HashPhone("254712345678")

// Email parsing
username, domain, err := helpers.ParseEmail("user@example.com")
// username: "user", domain: "example.com"

// Data conversion
bytes := helpers.ToBytes(myStruct)
var result MyStruct
helpers.FromBytes(bytes, &result)

// CSV parsing
file, _ := os.Open("data.csv")
rowsChan := helpers.ParseCSV(file)
for row := range rowsChan {
    // Process each CSV row
}

// Read entire file
data, _ := helpers.ReadFile("/path/to/file.txt")
```

## Configuration

GUtils uses environment variables for configuration. Create a `.env` file in your project root:

```bash
# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/mydb
DB_LOG_LEVEL=4

# Cache
CACHE_DRIVER=1
REDIS_URL=redis://localhost:6379

# Authentication
JWT_SECRET=your-super-secret-key

# HTTP Client
HTTP_CLIENT=0

# MongoDB
MONGO_URL=mongodb://localhost:27017
MONGO_DB=mydb

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

## Testing

Run tests for all packages:

```bash
make test
```

Run linter:

```bash
make lint
```

Clean dependencies:

```bash
make tidy
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
- All tests pass
- Code follows Go conventions
- New features include tests and documentation

## Requirements

- Go 1.22 or higher
- Optional dependencies based on usage:
  - PostgreSQL/MySQL/SQLite for `sqlr`
  - Redis for `cache` (when using Redis driver)
  - MongoDB for `nosql`
  - RabbitMQ for `pubsub`

## License

[Add your license here]

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/ochom/gutils).

## Acknowledgments

This library provides utilities commonly needed in Go applications, with inspiration from various frameworks and libraries in the Go ecosystem.
