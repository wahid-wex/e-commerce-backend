package migrations

import (
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/models"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.GetConfig())

// Up runs all migrations
func Up() {
	database := db.GetDb()
	createTables(database)
	createRolesAndPermissions(database)
	createDefaultUsers(database)
	createCategories(database)
	createProducts(database)
	createProductAttributes(database)
}

// createTables creates all tables in the database
func createTables(db *gorm.DB) {
	tables := []interface{}{}

	// User related tables
	tables = addNewTable(db, models.Customer{}, tables)
	tables = addNewTable(db, models.Seller{}, tables)
	tables = addNewTable(db, models.Admin{}, tables)

	// Role and Permission tables
	tables = addNewTable(db, models.Role{}, tables)
	tables = addNewTable(db, models.Permission{}, tables)
	tables = addNewTable(db, models.RolePermission{}, tables)
	tables = addNewTable(db, models.UserRole{}, tables)

	// Product related tables
	tables = addNewTable(db, models.Category{}, tables)
	tables = addNewTable(db, models.Product{}, tables)
	tables = addNewTable(db, models.ProductAttribute{}, tables)
	tables = addNewTable(db, models.ProductStock{}, tables)
	tables = addNewTable(db, models.Favorite{}, tables)
	tables = addNewTable(db, models.Review{}, tables)

	// Cart and Order tables
	tables = addNewTable(db, models.Cart{}, tables)
	tables = addNewTable(db, models.CartItem{}, tables)
	tables = addNewTable(db, models.Order{}, tables)
	tables = addNewTable(db, models.OrderItem{}, tables)
	tables = addNewTable(db, models.Payment{}, tables)

	err := db.Migrator().CreateTable(tables...)
	if err != nil {
		logger.Error(logging.Postgres, logging.Migration, err.Error(), nil)
	}
}

// addNewTable adds a new table to the list if it doesn't exist
func addNewTable(db *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !db.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

// createRolesAndPermissions creates default roles and permissions
func createRolesAndPermissions(db *gorm.DB) {
	// Create roles
	roles := []models.Role{
		{Name: "admin", Description: "مدیر سیستم"},
		{Name: "seller", Description: "فروشنده"},
		{Name: "customer", Description: "خریدار"},
	}

	for _, role := range roles {
		var count int64
		db.Model(&models.Role{}).Where("name = ?", role.Name).Count(&count)
		if count == 0 {
			db.Create(&role)
		}
	}

	// Create permissions
	permissions := []models.Permission{
		{Name: "user_management", Description: "مدیریت کاربران"},
		{Name: "product_management", Description: "مدیریت محصولات"},
		{Name: "order_management", Description: "مدیریت سفارشات"},
		{Name: "category_management", Description: "مدیریت دسته‌بندی‌ها"},
		{Name: "review_management", Description: "مدیریت نظرات"},
		{Name: "payment_management", Description: "مدیریت پرداخت‌ها"},
		{Name: "view_products", Description: "مشاهده محصولات"},
		{Name: "add_to_cart", Description: "افزودن به سبد خرید"},
		{Name: "checkout", Description: "تکمیل خرید"},
		{Name: "view_orders", Description: "مشاهده سفارشات"},
		{Name: "add_product", Description: "افزودن محصول"},
		{Name: "edit_product", Description: "ویرایش محصول"},
		{Name: "delete_product", Description: "حذف محصول"},
		{Name: "manage_stock", Description: "مدیریت موجودی"},
		{Name: "answer_questions", Description: "پاسخ به سوالات"},
	}

	for _, permission := range permissions {
		var count int64
		db.Model(&models.Permission{}).Where("name = ?", permission.Name).Count(&count)
		if count == 0 {
			db.Create(&permission)
		}
	}

	// Assign permissions to roles
	assignPermissionsToRoles(db)
}

// assignPermissionsToRoles assigns permissions to roles
func assignPermissionsToRoles(db *gorm.DB) {
	// Get roles
	var adminRole, sellerRole, customerRole models.Role
	db.Where("name = ?", "admin").First(&adminRole)
	db.Where("name = ?", "seller").First(&sellerRole)
	db.Where("name = ?", "customer").First(&customerRole)

	// Get permissions
	var permissions []models.Permission
	db.Find(&permissions)

	// Create a map for easier access
	permMap := make(map[string]uint)
	for _, perm := range permissions {
		permMap[perm.Name] = perm.ID
	}

	// Admin permissions (all permissions)
	for _, perm := range permissions {
		var count int64
		db.Model(&models.RolePermission{}).
			Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).
			Count(&count)

		if count == 0 {
			db.Create(&models.RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: perm.ID,
			})
		}
	}

	// Seller permissions
	sellerPermissions := []string{
		"product_management", "view_products", "add_product",
		"edit_product", "delete_product", "manage_stock",
		"answer_questions", "view_orders", "order_management",
	}

	for _, permName := range sellerPermissions {
		if permID, ok := permMap[permName]; ok {
			var count int64
			db.Model(&models.RolePermission{}).
				Where("role_id = ? AND permission_id = ?", sellerRole.ID, permID).
				Count(&count)

			if count == 0 {
				db.Create(&models.RolePermission{
					RoleID:       sellerRole.ID,
					PermissionID: permID,
				})
			}
		}
	}

	// Customer permissions
	customerPermissions := []string{
		"view_products", "add_to_cart", "checkout", "view_orders",
	}

	for _, permName := range customerPermissions {
		if permID, ok := permMap[permName]; ok {
			var count int64
			db.Model(&models.RolePermission{}).
				Where("role_id = ? AND permission_id = ?", customerRole.ID, permID).
				Count(&count)

			if count == 0 {
				db.Create(&models.RolePermission{
					RoleID:       customerRole.ID,
					PermissionID: permID,
				})
			}
		}
	}
}

// createDefaultUsers creates default users for testing
func createDefaultUsers(db *gorm.DB) {
	// Create admin user
	createAdminUser(db)

	// Create seller users
	createSellerUsers(db)

	// Create customer users
	createCustomerUsers(db)
}

// createAdminUser creates a default admin user
func createAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.Admin{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		// Create base user
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		// Create admin
		admin := models.Admin{
			FirstName:  "مدیر",
			LastName:   "سیستم",
			Department: "مدیریت",
			Username:   "admin",
			Email:      "admin@example.com",
			Password:   string(hashedPassword),
		}

		db.Create(&admin)

		// Assign admin role
		var adminRole models.Role
		db.Where("name = ?", "admin").First(&adminRole)

		userRole := models.UserRole{
			AdminID: &admin.ID,
			RoleID:  adminRole.ID,
		}
		db.Create(&userRole)
	}
}

// createSellerUsers creates default seller users
func createSellerUsers(db *gorm.DB) {
	sellers := []struct {
		Username        string
		Email           string
		Password        string
		StoreName       string
		BusinessLicense string
		NationalID      string
		Address         string
		Phone           string
		Description     string
	}{
		{
			Username:        "seller1",
			Email:           "seller1@example.com",
			Password:        "seller123",
			StoreName:       "فروشگاه دیجیتال آلفا",
			BusinessLicense: "12345678",
			NationalID:      "0123456789",
			Address:         "تهران، خیابان ولیعصر",
			Phone:           "09121234567",
			Description:     "فروشگاه محصولات الکترونیکی با بهترین قیمت",
		},
		{
			Username:        "seller2",
			Email:           "seller2@example.com",
			Password:        "seller123",
			StoreName:       "فروشگاه گجت برتر",
			BusinessLicense: "87654321",
			NationalID:      "9876543210",
			Address:         "تهران، خیابان شریعتی",
			Phone:           "09129876543",
			Description:     "ارائه دهنده جدیدترین محصولات الکترونیکی",
		},
	}

	var sellerRole models.Role
	db.Where("name = ?", "seller").First(&sellerRole)

	for _, s := range sellers {
		var count int64
		db.Model(&models.Seller{}).Where("username = ?", s.Username).Count(&count)

		if count == 0 {
			// Create base user
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
			// Create seller
			seller := models.Seller{
				Username:         s.Username,
				Email:            s.Email,
				Password:         string(hashedPassword),
				StoreName:        s.StoreName,
				BusinessLicense:  s.BusinessLicense,
				NationalID:       s.NationalID,
				SatisfactionRate: 0,
				RegistrationDate: time.Now(),
				IsVerified:       true,
				Address:          s.Address,
				Phone:            s.Phone,
				Description:      s.Description,
				Logo:             "",
			}
			db.Create(&seller)

			// Assign seller role
			userRole := models.UserRole{
				SellerID: &seller.ID,
				RoleID:   sellerRole.ID,
			}
			db.Create(&userRole)
		}
	}
}

// createCustomerUsers creates default customer users
func createCustomerUsers(db *gorm.DB) {
	customers := []struct {
		Username   string
		Email      string
		Password   string
		FirstName  string
		LastName   string
		Address    string
		PostalCode string
		Phone      string
		CardNumber string
	}{
		{
			Username:   "customer1",
			Email:      "customer1@example.com",
			Password:   "customer123",
			FirstName:  "علی",
			LastName:   "محمدی",
			Address:    "تهران، خیابان انقلاب",
			PostalCode: "1234567890",
			Phone:      "09131234567",
			CardNumber: "6037991234567890",
		},
		{
			Username:   "customer2",
			Email:      "customer2@example.com",
			Password:   "customer123",
			FirstName:  "مریم",
			LastName:   "احمدی",
			Address:    "اصفهان، خیابان چهارباغ",
			PostalCode: "9876543210",
			Phone:      "09139876543",
			CardNumber: "6037999876543210",
		},
	}

	var customerRole models.Role
	db.Where("name = ?", "customer").First(&customerRole)

	for _, c := range customers {
		var count int64
		db.Model(&models.Customer{}).Where("username = ?", c.Username).Count(&count)

		if count == 0 {
			// Create base user
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
			// Create customer
			customer := models.Customer{
				Username:        c.Username,
				Email:           c.Email,
				Password:        string(hashedPassword),
				FirstName:       c.FirstName,
				LastName:        c.LastName,
				Address:         c.Address,
				PostalCode:      c.PostalCode,
				Phone:           c.Phone,
				CardNumber:      c.CardNumber,
				ShippingAddress: c.Address,
			}
			db.Create(&customer)

			// Create cart for customer
			cart := models.Cart{
				CustomerID: customer.ID,
			}
			db.Create(&cart)

			// Assign customer role
			userRole := models.UserRole{
				SellerID: &customer.ID,
				RoleID:   customerRole.ID,
			}
			db.Create(&userRole)
		}
	}
}

// createCategories creates default product categories
func createCategories(db *gorm.DB) {
	categories := []struct {
		Name        string
		Description string
		ImageURL    string
	}{
		{
			Name:        "لپ تاپ",
			Description: "انواع لپ تاپ‌های گیمینگ و اداری",
			ImageURL:    "/images/categories/laptop.jpg",
		},
		{
			Name:        "موبایل",
			Description: "گوشی‌های هوشمند با برندهای مختلف",
			ImageURL:    "/images/categories/mobile.jpg",
		},
		{
			Name:        "تبلت",
			Description: "تبلت‌های اندرویدی و iOS",
			ImageURL:    "/images/categories/tablet.jpg",
		},
		{
			Name:        "لوازم جانبی",
			Description: "لوازم جانبی برای انواع دستگاه‌های الکترونیکی",
			ImageURL:    "/images/categories/accessories.jpg",
		},
		{
			Name:        "کامپیوتر",
			Description: "کامپیوترهای رومیزی و قطعات",
			ImageURL:    "/images/categories/computer.jpg",
		},
		{
			Name:        "صوتی و تصویری",
			Description: "تجهیزات صوتی و تصویری",
			ImageURL:    "/images/categories/audio-video.jpg",
		},
	}

	for _, c := range categories {
		var count int64
		db.Model(&models.Category{}).Where("name = ?", c.Name).Count(&count)

		if count == 0 {
			category := models.Category{
				Name:        c.Name,
				Description: c.Description,
				ImageURL:    c.ImageURL,
			}
			db.Create(&category)
		}
	}
}

// createProducts creates default products
func createProducts(db *gorm.DB) {
	// Get categories
	var laptopCategory, mobileCategory, tabletCategory, accessoriesCategory models.Category
	db.Where("name = ?", "لپ تاپ").First(&laptopCategory)
	db.Where("name = ?", "موبایل").First(&mobileCategory)
	db.Where("name = ?", "تبلت").First(&tabletCategory)
	db.Where("name = ?", "لوازم جانبی").First(&accessoriesCategory)

	// Get sellers
	var seller1, seller2 models.Seller
	db.Where("username = ?", "seller1").First(&seller1)
	db.Where("username = ?", "seller2").First(&seller2)

	// Create products
	products := []struct {
		CategoryID  uint
		Name        string
		Description string
		Price       float64
		ImageURL    string
		IsActive    bool
	}{
		{
			CategoryID:  laptopCategory.ID,
			Name:        "لپ تاپ ایسوس ROG Strix G15",
			Description: "لپ تاپ گیمینگ با پردازنده قدرتمند و کارت گرافیک مناسب برای بازی‌های سنگین",
			Price:       45000000,
			ImageURL:    "/images/products/asus-rog.jpg",
			IsActive:    true,
		},
		{
			CategoryID:  laptopCategory.ID,
			Name:        "لپ تاپ اپل MacBook Pro",
			Description: "لپ تاپ حرفه‌ای اپل با پردازنده M1 و صفحه نمایش Retina",
			Price:       75000000,
			ImageURL:    "/images/products/macbook-pro.jpg",
			IsActive:    true,
		},
		{
			CategoryID:  mobileCategory.ID,
			Name:        "گوشی سامسونگ Galaxy S21",
			Description: "گوشی هوشمند سامسونگ با دوربین حرفه‌ای و صفحه نمایش AMOLED",
			Price:       25000000,
			ImageURL:    "/images/products/samsung-s21.jpg",
			IsActive:    true,
		},
		{
			CategoryID:  mobileCategory.ID,
			Name:        "گوشی اپل iPhone 13",
			Description: "گوشی هوشمند اپل با پردازنده A15 Bionic و دوربین دوگانه",
			Price:       35000000,
			ImageURL:    "/images/products/iphone-13.jpg",
			IsActive:    true,
		},
		{
			CategoryID:  tabletCategory.ID,
			Name:        "تبلت اپل iPad Pro",
			Description: "تبلت حرفه‌ای اپل با پردازنده M1 و صفحه نمایش Liquid Retina",
			Price:       30000000,
			ImageURL:    "/images/products/ipad-pro.jpg",
			IsActive:    true,
		},
		{
			CategoryID:  accessoriesCategory.ID,
			Name:        "هدفون بی‌سیم سونی WH-1000XM4",
			Description: "هدفون بی‌سیم با قابلیت حذف نویز و کیفیت صدای فوق‌العاده",
			Price:       8000000,
			ImageURL:    "/images/products/sony-headphone.jpg",
			IsActive:    true,
		},
	}

	for _, p := range products {
		var count int64
		db.Model(&models.Product{}).Where("name = ? AND seller_id = ?", p.Name).Count(&count)

		if count == 0 {
			product := models.Product{
				CategoryID:       p.CategoryID,
				Name:             p.Name,
				Description:      p.Description,
				Price:            p.Price,
				ImageURL:         p.ImageURL,
				IsActive:         p.IsActive,
				SatisfactionRate: 0,
			}
			db.Create(&product)

			// Create product stock
			stock := models.ProductStock{
				ProductID: product.ID,
				Quantity:  100,
			}
			db.Create(&stock)
		}
	}
}

// createProductAttributes creates attributes for products
func createProductAttributes(db *gorm.DB) {
	// Get products
	var products []models.Product
	db.Find(&products)

	for _, product := range products {
		var count int64
		db.Model(&models.ProductAttribute{}).Where("product_id = ?", product.ID).Count(&count)

		if count == 0 {
			// Create attributes based on product category
			switch product.CategoryID {
			case 1: // Laptop
				createLaptopAttributes(db, product.ID)
			case 2: // Mobile
				createMobileAttributes(db, product.ID)
			case 3: // Tablet
				createTabletAttributes(db, product.ID)
			case 4: // Accessories
				createAccessoryAttributes(db, product.ID)
			}
		}
	}
}

// createLaptopAttributes creates attributes for laptop products
func createLaptopAttributes(db *gorm.DB, productID uint) {
	attributes := []struct {
		Key   string
		Value string
	}{
		{Key: "پردازنده", Value: "Intel Core i7-11800H"},
		{Key: "حافظه رم", Value: "16GB DDR4"},
		{Key: "حافظه داخلی", Value: "1TB SSD"},
		{Key: "کارت گرافیک", Value: "NVIDIA GeForce RTX 3060"},
		{Key: "اندازه صفحه نمایش", Value: "15.6 اینچ"},
		{Key: "رزولوشن", Value: "1920x1080"},
		{Key: "سیستم عامل", Value: "Windows 11"},
		{Key: "وزن", Value: "2.3 کیلوگرم"},
		{Key: "رنگ", Value: "مشکی"},
	}

	for _, attr := range attributes {
		productAttr := models.ProductAttribute{
			ProductID: productID,
			Key:       attr.Key,
			Value:     attr.Value,
		}
		db.Create(&productAttr)
	}
}

// createMobileAttributes creates attributes for mobile products
func createMobileAttributes(db *gorm.DB, productID uint) {
	attributes := []struct {
		Key   string
		Value string
	}{
		{Key: "پردازنده", Value: "Snapdragon 888"},
		{Key: "حافظه رم", Value: "8GB"},
		{Key: "حافظه داخلی", Value: "128GB"},
		{Key: "دوربین اصلی", Value: "64 مگاپیکسل"},
		{Key: "دوربین سلفی", Value: "12 مگاپیکسل"},
		{Key: "اندازه صفحه نمایش", Value: "6.2 اینچ"},
		{Key: "رزولوشن", Value: "2400x1080"},
		{Key: "باتری", Value: "4500mAh"},
		{Key: "سیستم عامل", Value: "Android 12"},
		{Key: "رنگ", Value: "آبی"},
	}

	for _, attr := range attributes {
		productAttr := models.ProductAttribute{
			ProductID: productID,
			Key:       attr.Key,
			Value:     attr.Value,
		}
		db.Create(&productAttr)
	}
}

// createTabletAttributes creates attributes for tablet products
func createTabletAttributes(db *gorm.DB, productID uint) {
	attributes := []struct {
		Key   string
		Value string
	}{
		{Key: "پردازنده", Value: "Apple M1"},
		{Key: "حافظه رم", Value: "8GB"},
		{Key: "حافظه داخلی", Value: "256GB"},
		{Key: "دوربین اصلی", Value: "12 مگاپیکسل"},
		{Key: "دوربین سلفی", Value: "7 مگاپیکسل"},
		{Key: "اندازه صفحه نمایش", Value: "11 اینچ"},
		{Key: "رزولوشن", Value: "2388x1668"},
		{Key: "باتری", Value: "7538mAh"},
		{Key: "سیستم عامل", Value: "iPadOS 15"},
		{Key: "رنگ", Value: "نقره‌ای"},
	}

	for _, attr := range attributes {
		productAttr := models.ProductAttribute{
			ProductID: productID,
			Key:       attr.Key,
			Value:     attr.Value,
		}
		db.Create(&productAttr)
	}
}

// createAccessoryAttributes creates attributes for accessory products
func createAccessoryAttributes(db *gorm.DB, productID uint) {
	attributes := []struct {
		Key   string
		Value string
	}{
		{Key: "نوع اتصال", Value: "بی‌سیم"},
		{Key: "قابلیت حذف نویز", Value: "دارد"},
		{Key: "عمر باتری", Value: "30 ساعت"},
		{Key: "وزن", Value: "254 گرم"},
		{Key: "رنگ", Value: "مشکی"},
	}

	for _, attr := range attributes {
		productAttr := models.ProductAttribute{
			ProductID: productID,
			Key:       attr.Key,
			Value:     attr.Value,
		}
		db.Create(&productAttr)
	}
}
