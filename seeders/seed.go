package seeders

import (
	"order-management/entity"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Clean() error {
	log.Info("Cleaning database...")

	// Delete in order to respect foreign key constraints
	if err := s.db.Where("1 = 1").Delete(&entity.OrderProduct{}).Error; err != nil {
		log.Error("Failed to clean OrderProduct table:", err)
		return err
	}
	log.Info("Cleaned OrderProduct table")

	if err := s.db.Where("1 = 1").Delete(&entity.Order{}).Error; err != nil {
		log.Error("Failed to clean Order table:", err)
		return err
	}
	log.Info("Cleaned Order table")

	if err := s.db.Where("1 = 1").Delete(&entity.Product{}).Error; err != nil {
		log.Error("Failed to clean Product table:", err)
		return err
	}
	log.Info("Cleaned Product table")

	if err := s.db.Where("1 = 1").Delete(&entity.Shop{}).Error; err != nil {
		log.Error("Failed to clean Shop table:", err)
		return err
	}
	log.Info("Cleaned Shop table")

	if err := s.db.Where("1 = 1").Delete(&entity.User{}).Error; err != nil {
		log.Error("Failed to clean User table:", err)
		return err
	}
	log.Info("Cleaned User table")

	log.Info("Database cleaning completed")
	return nil
}

func (s *Seeder) Seed() error {
	log.Info("Starting database seeding...")

	// Clean existing data
	if err := s.Clean(); err != nil {
		return err
	}

	// Seed users
	userIDs, err := s.seedUsers()
	if err != nil {
		return err
	}

	// Seed shops
	shopIDs, err := s.seedShops()
	if err != nil {
		return err
	}

	// Seed products
	productIDs, err := s.seedProducts(shopIDs)
	if err != nil {
		return err
	}

	// Seed orders
	if err := s.seedOrders(userIDs, productIDs); err != nil {
		return err
	}

	log.Info("Database seeding completed successfully")
	return nil
}

func (s *Seeder) seedUsers() ([]uint32, error) {
	log.Info("Seeding users...")
	users := []entity.User{
		{
			Email:    "user1@example.com",
			Password: "password1",
			Address:  "123 Main St, New York, NY 10001",
		},
		{
			Email:    "user2@example.com",
			Password: "password2",
			Address:  "456 Oak Ave, Los Angeles, CA 90001",
		},
		{
			Email:    "user3@example.com",
			Password: "password3",
			Address:  "789 Pine St, Chicago, IL 60601",
		},
		{
			Email:    "user4@example.com",
			Password: "password4",
			Address:  "321 Elm St, Houston, TX 77001",
		},
	}

	userIDs := make([]uint32, 0, len(users))
	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("Failed to hash password for user:", user.Email, err)
			return nil, err
		}
		user.Password = string(hashedPassword)
		if err := s.db.Create(&user).Error; err != nil {
			log.Error("Failed to create user:", user.Email, err)
			return nil, err
		}
		log.Info("Created user:", user.Email)
		userIDs = append(userIDs, user.ID)
	}

	log.Info("User seeding completed")
	return userIDs, nil
}

func (s *Seeder) seedShops() ([]uint32, error) {
	log.Info("Seeding shops...")
	shops := []entity.Shop{
		{
			Name:        "Tech Gadgets",
			Description: "Your one-stop shop for the latest tech gadgets and accessories",
			Password:    "tech123",
		},
		{
			Name:        "Fashion Boutique",
			Description: "Trendy clothing and accessories for men and women",
			Password:    "fashion123",
		},
		{
			Name:        "Home Decor",
			Description: "Beautiful home decor items to make your space special",
			Password:    "home123",
		},
		{
			Name:        "Sports Equipment",
			Description: "High-quality sports equipment for all your fitness needs",
			Password:    "sports123",
		},
	}

	shopIDs := make([]uint32, 0, len(shops))
	for _, shop := range shops {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(shop.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("Failed to hash password for shop:", shop.Name, err)
			return nil, err
		}
		shop.Password = string(hashedPassword)
		if err := s.db.Create(&shop).Error; err != nil {
			log.Error("Failed to create shop:", shop.Name, err)
			return nil, err
		}
		log.Info("Created shop:", shop.Name)
		shopIDs = append(shopIDs, shop.ID)
	}

	log.Info("Shop seeding completed")
	return shopIDs, nil
}

func (s *Seeder) seedProducts(shopIDs []uint32) ([]uint32, error) {
	log.Info("Seeding products...")
	products := []struct {
		Name        string
		Description string
		Price       uint32
		ShopIndex   int
	}{
		// Tech Gadgets products
		{
			Name:        "Wireless Earbuds",
			Description: "Premium wireless earbuds with noise cancellation",
			Price:       2999,
			ShopIndex:   0,
		},
		{
			Name:        "Smart Watch",
			Description: "Feature-rich smartwatch with health tracking",
			Price:       4999,
			ShopIndex:   0,
		},
		{
			Name:        "Portable Charger",
			Description: "High-capacity portable power bank",
			Price:       1999,
			ShopIndex:   0,
		},
		// Fashion Boutique products
		{
			Name:        "Leather Jacket",
			Description: "Classic black leather jacket",
			Price:       8999,
			ShopIndex:   1,
		},
		{
			Name:        "Designer Handbag",
			Description: "Elegant designer handbag",
			Price:       12999,
			ShopIndex:   1,
		},
		{
			Name:        "Silk Scarf",
			Description: "Luxurious silk scarf",
			Price:       2999,
			ShopIndex:   1,
		},
		// Home Decor products
		{
			Name:        "Modern Lamp",
			Description: "Contemporary table lamp",
			Price:       3999,
			ShopIndex:   2,
		},
		{
			Name:        "Wall Art",
			Description: "Abstract wall art painting",
			Price:       5999,
			ShopIndex:   2,
		},
		{
			Name:        "Throw Pillow",
			Description: "Decorative throw pillow",
			Price:       1999,
			ShopIndex:   2,
		},
		// Sports Equipment products
		{
			Name:        "Yoga Mat",
			Description: "Premium non-slip yoga mat",
			Price:       2499,
			ShopIndex:   3,
		},
		{
			Name:        "Dumbbell Set",
			Description: "Adjustable dumbbell set",
			Price:       7999,
			ShopIndex:   3,
		},
		{
			Name:        "Running Shoes",
			Description: "High-performance running shoes",
			Price:       8999,
			ShopIndex:   3,
		},
	}

	productIDs := make([]uint32, 0, len(products))
	for _, p := range products {
		product := entity.Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			ShopID:      shopIDs[p.ShopIndex],
		}
		if err := s.db.Create(&product).Error; err != nil {
			log.Error("Failed to create product:", product.Name, err)
			return nil, err
		}
		log.Info("Created product:", product.Name)
		productIDs = append(productIDs, product.ID)
	}

	log.Info("Product seeding completed")
	return productIDs, nil
}

func (s *Seeder) seedOrders(userIDs, productIDs []uint32) error {
	log.Info("Seeding orders...")
	// Create orders
	orders := []struct {
		Status    entity.Status
		Courier   string
		UserIndex int
	}{
		{
			Status:    entity.PENDING,
			Courier:   "J&T Express",
			UserIndex: 0,
		},
		{
			Status:    entity.SHIPPING,
			Courier:   "Kerry Express",
			UserIndex: 1,
		},
		{
			Status:    entity.COMPLETED,
			Courier:   "DHL",
			UserIndex: 2,
		},
		{
			Status:    entity.CANCELLED,
			Courier:   "FedEx",
			UserIndex: 3,
		},
		{
			Status:    entity.PENDING,
			Courier:   "UPS",
			UserIndex: 0,
		},
	}

	// Create order products first to calculate totals
	orderProducts := []struct {
		OrderIndex   int
		ProductIndex int
		Amount       uint32
	}{
		// Order 1: Tech items
		{
			OrderIndex:   0,
			ProductIndex: 0, // Wireless Earbuds (2999)
			Amount:       1,
		},
		{
			OrderIndex:   0,
			ProductIndex: 2, // Portable Charger (1999)
			Amount:       2,
		},
		// Order 2: Fashion items
		{
			OrderIndex:   1,
			ProductIndex: 3, // Leather Jacket (8999)
			Amount:       1,
		},
		{
			OrderIndex:   1,
			ProductIndex: 5, // Silk Scarf (2999)
			Amount:       3,
		},
		// Order 3: Home decor items
		{
			OrderIndex:   2,
			ProductIndex: 6, // Modern Lamp (3999)
			Amount:       2,
		},
		{
			OrderIndex:   2,
			ProductIndex: 7, // Wall Art (5999)
			Amount:       1,
		},
		// Order 4: Sports items
		{
			OrderIndex:   3,
			ProductIndex: 9, // Dumbbell Set (7999)
			Amount:       1,
		},
		{
			OrderIndex:   3,
			ProductIndex: 10, // Running Shoes (8999)
			Amount:       1,
		},
		// Order 5: Mixed items
		{
			OrderIndex:   4,
			ProductIndex: 1, // Smart Watch (4999)
			Amount:       1,
		},
		{
			OrderIndex:   4,
			ProductIndex: 4, // Designer Handbag (12999)
			Amount:       1,
		},
		{
			OrderIndex:   4,
			ProductIndex: 8, // Throw Pillow (1999)
			Amount:       2,
		},
	}

	// Calculate order totals
	orderTotals := make([]float32, len(orders))
	for _, op := range orderProducts {
		// Get product price
		var product entity.Product
		if err := s.db.First(&product, productIDs[op.ProductIndex]).Error; err != nil {
			log.Error("Failed to get product price:", err)
			return err
		}
		// Add to order total
		orderTotals[op.OrderIndex] += float32(product.Price) * float32(op.Amount)
	}

	// Create orders with correct totals
	orderIDs := make([]uint32, 0, len(orders))
	for i, o := range orders {
		order := entity.Order{
			Status:  o.Status,
			Total:   orderTotals[i],
			Courier: o.Courier,
			UserID:  userIDs[o.UserIndex],
		}
		if err := s.db.Create(&order).Error; err != nil {
			log.Error("Failed to create order for user:", order.UserID, err)
			return err
		}
		log.Info("Created order for user:", order.UserID, "with total:", order.Total)
		orderIDs = append(orderIDs, order.ID)
	}

	// Create order products
	for _, op := range orderProducts {
		orderProduct := entity.OrderProduct{
			OrderID:   orderIDs[op.OrderIndex],
			ProductID: productIDs[op.ProductIndex],
			Amount:    op.Amount,
		}
		if err := s.db.Create(&orderProduct).Error; err != nil {
			log.Error("Failed to create order product for order:", orderProduct.OrderID, err)
			return err
		}
		log.Info("Created order product for order:", orderProduct.OrderID)
	}

	log.Info("Order seeding completed")
	return nil
}
