package main

import "sync"

type ProductStore struct {
	mu sync.RWMutex // use RWmutex because our workload is read-heavy
	// ? how this variable would be used latter
	products map[int]Product 
}

// NewProductStore creates a store and seeds it with sample products.
func NewProductStore() *ProductStore {
	store := &ProductStore{
		products: make(map[int]Product),
	}

	// Seed with sample products so GET works immediately
	seedProducts := []Product{
		{
			ProductID:    1,
			SKU:          "ABC-123-XYZ",
			Manufacturer: "Acme Corporation",
			CategoryID:   10,
			Weight:       1250,
			SomeOtherID:  100,
		},
		{
			ProductID:    2,
			SKU:          "DEF-456-UVW",
			Manufacturer: "Globex Industries",
			CategoryID:   20,
			Weight:       500,
			SomeOtherID:  200,
		},
		{
			ProductID:    3,
			SKU:          "GHI-789-RST",
			Manufacturer: "Initech LLC",
			CategoryID:   10,
			Weight:       3000,
			SomeOtherID:  300,
		},
	}

	for _, p := range seedProducts {
		store.products[p.ProductID] = p
	}

	return store
}

// Get retrieves a product by ID.
// receiver function, method-like function 
func (s *ProductStore) Get(id int) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock() // using defer to trigger unlock at the end of function
	product, exists := s.products[id]
	return product, exists
}

// Update replaces the data for an existing product.
func (s *ProductStore) Update(id int, p Product) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.products[id]
	if !exists {
		return false
	}
	s.products[id] = p
	return exists
}