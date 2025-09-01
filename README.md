# SQL GORM for Tinh Tinh

<div align="center">
<img alt="GitHub Release" src="https://img.shields.io/github/v/release/tinh-tinh/sqlorm">
<img alt="GitHub License" src="https://img.shields.io/github/license/tinh-tinh/sqlorm">
<a href="https://codecov.io/gh/tinh-tinh/sqlorm">
    <img src="https://codecov.io/gh/tinh-tinh/sqlorm/graph/badge.svg?token=TS4B5QAO3T"/>
</a>
<a href="https://pkg.go.dev/github.com/tinh-tinh/sqlorm"><img src="https://pkg.go.dev/badge/github.com/tinh-tinh/sqlorm.svg" alt="Go Reference"></a>
</div>

<div align="center">
    <img src="https://avatars.githubusercontent.com/u/178628733?s=400&u=2a8230486a43595a03a6f9f204e54a0046ce0cc4&v=4" width="200" alt="Tinh Tinh Logo">
</div>

## Overview

SQL GORM for Tinh Tinh is a powerful database toolkit designed to work seamlessly with the Tinh Tinh framework. It provides an elegant and efficient way to interact with SQL databases using GORM, the fantastic ORM library for Golang.

## Features

- 🚀 Full GORM integration with Tinh Tinh
- 📦 Easy-to-use database operations
- 🔄 Auto Migration support
- 🎯 Type-safe query building
- 🛠️ Advanced features like:
  - Associations handling
  - Hooks
  - Transactions
  - Custom data types
  - And more!

## Installation

To install the package, use:

```bash
go get -u github.com/tinh-tinh/sqlorm/v2
```

## Quick Start

```go
package main

import (
    "github.com/tinh-tinh/sqlorm/v2"
)

// User represents your database model
type User struct {
    ID    uint   `gorm:"primarykey"`
    Name  string
    Email string
}

func main() {
    // Initialize your database connection
    db := sqlorm.New(&sqlorm.Config{
        Driver:   "postgres",
        Host:     "localhost",
        Port:     5432,
        Database: "mydb",
        Username: "user",
        Password: "password",
    })

    // Auto migrate your models
    db.AutoMigrate(&User{})

    // Create a new user
    user := User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    db.Create(&user)
}
```

## Configuration

The package supports various database configurations:

```go
type Config struct {
    Driver   string // "postgres", "mysql", "sqlite"
    Host     string
    Port     int
    Database string
    Username string
    Password string
    SSLMode  string
    TimeZone string
}
```

## Supported Databases

- PostgreSQL
- MySQL
- SQLite
- Microsoft SQL Server

## Documentation

For detailed documentation and examples, please visit:
- [Go Package Documentation](https://pkg.go.dev/github.com/tinh-tinh/sqlorm)
- [GORM Official Documentation](https://gorm.io/docs/)

## Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

If you encounter any issues or need help, you can:
- Open an issue in the GitHub repository
- Check our documentation
- Join our community discussions
