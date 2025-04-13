erDiagram
User {
uint ID PK
string Username
string Email
string Password
string FirstName
string LastName
string Address
string Phone
enum Role
datetime CreatedAt
datetime UpdatedAt
}

    Category {
        uint ID PK
        string Name
        string Description
        string ImageURL
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    Product {
        uint ID PK
        uint SellerID FK
        uint CategoryID FK
        string Name
        string Description
        float Price
        string ImageURL
        boolean IsActive
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    ProductStock {
        uint ID PK
        uint ProductID FK
        uint SellerID FK
        int Quantity
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    Cart {
        uint ID PK
        uint UserID FK
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    CartItem {
        uint ID PK
        uint CartID FK
        uint ProductID FK
        int Quantity
        float Price
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    Order {
        uint ID PK
        uint UserID FK
        uint SellerID FK
        float TotalAmount
        enum Status
        string ShippingAddress
        string TrackingNumber
        datetime OrderDate
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    OrderItem {
        uint ID PK
        uint OrderID FK
        uint ProductID FK
        int Quantity
        float Price
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    Payment {
        uint ID PK
        uint OrderID FK
        float Amount
        enum Status
        string TransactionID
        string PaymentMethod
        datetime PaymentDate
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    Review {
        uint ID PK
        uint UserID FK
        uint ProductID FK
        string Content
        int Rating
        datetime CreatedAt
        datetime UpdatedAt
    }
    
    User ||--o{ Product : "sells"
    User ||--o{ Cart : "has"
    User ||--o{ Order : "places"
    User ||--o{ Review : "writes"
    
    Category ||--o{ Product : "contains"
    
    Product ||--o{ CartItem : "added_to"
    Product ||--o{ OrderItem : "ordered_in"
    Product ||--o{ Review : "receives"
    Product ||--o{ ProductStock : "has"
    
    Cart ||--o{ CartItem : "contains"
    
    Order ||--o{ OrderItem : "contains"
    Order ||--|| Payment : "has"
