package main

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// User model - GORM uses struct tags to map to database
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null;size:255"`
    Username  string    `gorm:"uniqueIndex;not null;size:50"`
    FullName  string    `gorm:"not null;size:100"`
    Age       *int      `gorm:"default:null"`
    IsActive  bool      `gorm:"default:true;index"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name
func (User) TableName() string {
    return "users"
}

func gormExamples() {
    // Connect to database
    dsn := "host=localhost user=user password=pass dbname=dbname port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Auto-migrate (creates table if it doesn't exist)
    db.AutoMigrate(&User{})

    // 1. CREATE - Insert a new user
    fmt.Println("=== CREATE USER ===")
    age := 30
    newUser := User{
        Email:    "john@example.com",
        Username: "johndoe",
        FullName: "John Doe",
        Age:      &age,
        IsActive: true,
    }
    result := db.Create(&newUser)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Created user: %+v (Rows affected: %d)\n", newUser, result.RowsAffected)

    // 2. READ - Get user by ID
    fmt.Println("\n=== GET USER BY ID ===")
    var user User
    result = db.First(&user, newUser.ID)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Found user: %+v\n", user)

    // 3. READ - Get user by email
    fmt.Println("\n=== GET USER BY EMAIL ===")
    var userByEmail User
    result = db.Where("email = ?", "john@example.com").First(&userByEmail)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Found user by email: %+v\n", userByEmail)

    // 4. LIST - Get paginated users
    fmt.Println("\n=== LIST USERS (PAGINATED) ===")
    var users []User
    result = db.Order("created_at DESC").Limit(10).Offset(0).Find(&users)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    for _, u := range users {
        fmt.Printf("User: %s (%s)\n", u.FullName, u.Email)
    }

    // 5. LIST - Get active users
    fmt.Println("\n=== LIST ACTIVE USERS ===")
    var activeUsers []User
    result = db.Where("is_active = ?", true).Order("created_at DESC").Find(&activeUsers)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Found %d active users\n", len(activeUsers))

    // 6. UPDATE - Update user details (using Save)
    fmt.Println("\n=== UPDATE USER (SAVE) ===")
    user.Email = "john.doe@example.com"
    user.FullName = "John Michael Doe"
    newAge := 31
    user.Age = &newAge
    result = db.Save(&user)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Updated user: %+v\n", user)

    // 7. UPDATE - Update specific fields only
    fmt.Println("\n=== UPDATE USER (UPDATES) ===")
    result = db.Model(&User{}).Where("id = ?", newUser.ID).Updates(map[string]interface{}{
        "is_active":  false,
        "updated_at": time.Now(),
    })
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("User status updated (Rows affected: %d)\n", result.RowsAffected)

    // 8. SEARCH - Search users by name (ILIKE)
    fmt.Println("\n=== SEARCH USERS BY NAME ===")
    var searchResults []User
    result = db.Where("full_name ILIKE ?", "%John%").Order("full_name").Find(&searchResults)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Found %d users matching 'John'\n", len(searchResults))

    // 9. FILTER - Get users by age range
    fmt.Println("\n=== GET USERS BY AGE RANGE ===")
    var usersInRange []User
    result = db.Where("age BETWEEN ? AND ?", 25, 35).Order("age").Find(&usersInRange)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Found %d users aged 25-35\n", len(usersInRange))

    // 10. COUNT - Count total users
    fmt.Println("\n=== COUNT USERS ===")
    var count int64
    result = db.Model(&User{}).Count(&count)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Total users: %d\n", count)

    // 11. DELETE - Delete user
    fmt.Println("\n=== DELETE USER ===")
    result = db.Delete(&User{}, newUser.ID)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("User deleted (Rows affected: %d)\n", result.RowsAffected)

    // 12. TRANSACTIONS - Multiple operations
    fmt.Println("\n=== TRANSACTION EXAMPLE ===")
    err = db.Transaction(func(tx *gorm.DB) error {
        // Create multiple users in a transaction
        age1 := 25
        user1 := User{
            Email:    "alice@example.com",
            Username: "alice",
            FullName: "Alice Smith",
            Age:      &age1,
            IsActive: true,
        }
        if err := tx.Create(&user1).Error; err != nil {
            return err
        }

        age2 := 28
        user2 := User{
            Email:    "bob@example.com",
            Username: "bob",
            FullName: "Bob Johnson",
            Age:      &age2,
            IsActive: true,
        }
        if err := tx.Create(&user2).Error; err != nil {
            return err
        }

        fmt.Printf("Created users in transaction: %s, %s\n", user1.Username, user2.Username)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }

    // 13. COMPLEX QUERY - Chain multiple conditions
    fmt.Println("\n=== COMPLEX QUERY ===")
    var complexResults []User
    result = db.Where("is_active = ?", true).
        Where("age >= ?", 25).
        Or("username LIKE ?", "admin%").
        Order("created_at DESC").
        Limit(5).
        Find(&complexResults)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Found %d users with complex filters\n", len(complexResults))

    // 14. SELECT SPECIFIC FIELDS
    fmt.Println("\n=== SELECT SPECIFIC FIELDS ===")
    var partialUsers []struct {
        ID       uint
        Username string
        Email    string
    }
    result = db.Model(&User{}).Select("id", "username", "email").Find(&partialUsers)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    for _, u := range partialUsers {
        fmt.Printf("User: %s (%s)\n", u.Username, u.Email)
    }

    // 15. PRELOADING (if you had relations)
    // Example: db.Preload("Posts").Preload("Profile").Find(&users)

    // 16. RAW SQL (when needed)
    fmt.Println("\n=== RAW SQL QUERY ===")
    var rawResults []User
    result = db.Raw("SELECT * FROM users WHERE age > ? ORDER BY age LIMIT ?", 20, 5).Scan(&rawResults)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Found %d users with raw SQL\n", len(rawResults))
}
