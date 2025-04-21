package test

import "os"

// TestData provides test JSON queries for conversion tests
type TestData struct {
	Name        string
	Description string
	JSON        string
}

// GetTestData returns a slice of test cases for JSON to SQL conversion
func GetTestData() []TestData {
	return []TestData{
		{
			Name:        "EcommerceQuery",
			Description: "E-commerce order with multiple nested relations",
			JSON: `{
  "orders": {
    "select": ["id", "order_date", "total_amount", "status"],
    "where": {
      "and": [
        { "status": {"in": ["shipped", "delivered"]} },
        { "order_date": {">=": "2023-06-01"} },
        { "order_date": {"<=": "2023-12-31"} }
      ],
      "total_amount": {">": 50}
    },
    "order": ["-order_date", "total_amount"],
    "limit": 100,
    "customer": {
      "select": ["id", "name", "email", "phone"],
      "join": "customer_id:id",
      "addresses": {
        "select": ["id", "street", "city", "postal_code", "country"],
        "join": "customer_id:id",
        "where": {
          "is_primary": true
        }
      }
    },
    "items": {
      "select": ["id", "product_id", "quantity", "unit_price", "subtotal"],
      "join": "order_id:id",
      "order": ["-subtotal"],
      "product": {
        "select": ["id", "name", "sku", "category_id"],
        "join": "id:product_id",
        "category": {
          "select": ["id", "name", "parent_id"],
          "join": "id:category_id"
        }
      }
    },
    "payments": {
      "select": ["id", "amount", "payment_method", "status", "transaction_id"],
      "join": "order_id:id",
      "where": {
        "status": "completed"
      }
    },
    "shipments": {
      "select": ["id", "tracking_number", "carrier", "shipping_date", "status"],
      "join": "order_id:id"
    }
  }
}`,
		},
		{
			Name:        "BlogQuery",
			Description: "Blog with advanced filtering",
			JSON: `{
  "posts": {
    "select": ["id", "title", "content", "published_at", "view_count"],
    "where": {
      "and": [
        { "status": "published" },
        { "published_at": {">=": "2023-01-01"} },
        { 
          "or": [
            { "category_id": {"in": [1, 3, 5]} },
            { "is_featured": true }
          ]
        }
      ]
    },
    "order": ["-published_at"],
    "limit": 20,
    "author": {
      "select": ["id", "name", "bio", "avatar_url"],
      "join": "author_id:id"
    },
    "categories": {
      "select": ["id", "name", "slug"],
      "join": "post_id:id"
    },
    "tags": {
      "select": ["id", "name", "slug"],
      "join": "post_id:id"
    },
    "comments": {
      "select": ["id", "content", "created_at", "status"],
      "join": "post_id:id",
      "where": {
        "status": "approved"
      },
      "order": ["created_at"],
      "limit": 10,
      "user": {
        "select": ["id", "name", "avatar_url"],
        "join": "user_id:id"
      }
    }
  }
}`,
		},
		{
			Name:        "AnalyticsQuery",
			Description: "Analytics dashboard with complex aggregations",
			JSON: `{
  "page_views": {
    "select": ["date", "page_path", "view_count", "visitor_count"],
    "where": {
      "and": [
        { "date": {">=": "2023-09-01"} },
        { "date": {"<=": "2023-09-30"} }
      ],
      "site_id": 123
    },
    "order": ["-view_count"],
    "limit": 50,
    "referrers": {
      "select": ["source", "medium", "count"],
      "join": "page_view_id:id",
      "order": ["-count"],
      "limit": 5
    },
    "user_metrics": {
      "select": ["device_type", "browser", "operating_system", "user_count"],
      "join": "page_view_id:id",
      "where": {
        "user_count": {">": 10}
      }
    }
  }
}`,
		},
		{
			Name:        "HRQuery",
			Description: "HR system with multiple entity relationships",
			JSON: `{
  "employees": {
    "select": ["id", "first_name", "last_name", "email", "hire_date", "status"],
    "where": {
      "and": [
        { "status": "active" },
        { 
          "or": [
            { "department_id": {"in": [2, 3, 7]} },
            { "position": {"in": ["manager", "director", "vp"]} }
          ]
        }
      ]
    },
    "order": ["last_name", "first_name"],
    "department": {
      "select": ["id", "name", "location_id"],
      "join": "department_id:id",
      "location": {
        "select": ["id", "city", "country", "address"],
        "join": "id:location_id"
      }
    },
    "salary_history": {
      "select": ["id", "amount", "effective_date", "reason"],
      "join": "employee_id:id",
      "order": ["-effective_date"],
      "limit": 3
    },
    "performance_reviews": {
      "select": ["id", "review_date", "rating", "notes"],
      "join": "employee_id:id",
      "where": {
        "review_date": {">=": "2022-01-01"}
      },
      "order": ["-review_date"]
    },
    "projects": {
      "select": ["id", "name", "start_date", "end_date", "status"],
      "join": "employee_id:id",
      "where": {
        "status": {"in": ["in_progress", "completed"]}
      }
    }
  }
}`,
		},
		{
			Name:        "LibraryQuery",
			Description: "Library management system",
			JSON: `{
  "books": {
    "select": ["id", "title", "isbn", "publication_year", "language", "copies_available"],
    "where": {
      "and": [
        { "status": "active" },
        { "copies_available": {">": 0} }
      ],
      "or": [
        { "category_id": {"in": [3, 5, 8]} },
        { "is_bestseller": true }
      ]
    },
    "order": ["title"],
    "limit": 50,
    "authors": {
      "select": ["id", "name", "birth_year", "nationality"],
      "join": "book_id:id"
    },
    "publisher": {
      "select": ["id", "name", "location"],
      "join": "publisher_id:id"
    },
    "categories": {
      "select": ["id", "name", "parent_id"],
      "join": "book_id:id"
    },
    "reviews": {
      "select": ["id", "rating", "comment", "created_at"],
      "join": "book_id:id",
      "where": {
        "rating": {">=": 4}
      },
      "order": ["-created_at"],
      "limit": 5,
      "user": {
        "select": ["id", "username", "membership_level"],
        "join": "user_id:id"
      }
    },
    "borrowing_history": {
      "select": ["id", "checkout_date", "due_date", "return_date", "status"],
      "join": "book_id:id",
      "order": ["-checkout_date"],
      "limit": 10
    }
  }
}`,
		},
		{
			Name:        "RealEstateQuery",
			Description: "Real estate property listings",
			JSON: `{
  "properties": {
    "select": ["id", "title", "description", "price", "built_year", "listing_date", "status"],
    "where": {
      "and": [
        { "status": "active" },
        { "price": {">=": 300000} },
        { "price": {"<=": 700000} },
        { "bedrooms": {">=": 3} },
        { "bathrooms": {">=": 2} }
      ],
      "or": [
        { "property_type": {"in": ["single_family", "townhouse", "condo"]} },
        { "has_pool": true }
      ]
    },
    "order": ["-listing_date", "price"],
    "limit": 30,
    "location": {
      "select": ["id", "address", "city", "state", "zip_code", "latitude", "longitude"],
      "join": "location_id:id",
      "where": {
        "state": {"in": ["CA", "NY", "FL", "TX"]}
      }
    },
    "features": {
      "select": ["id", "name", "category"],
      "join": "property_id:id"
    },
    "photos": {
      "select": ["id", "url", "caption", "is_primary"],
      "join": "property_id:id",
      "order": ["display_order"]
    },
    "agent": {
      "select": ["id", "name", "email", "phone", "license_number"],
      "join": "agent_id:id",
      "agency": {
        "select": ["id", "name", "logo_url", "website"],
        "join": "agency_id:id"
      }
    },
    "open_houses": {
      "select": ["id", "start_time", "end_time", "date"],
      "join": "property_id:id",
      "where": {
        "date": {">=": "2023-01-01"}
      },
      "order": ["date", "start_time"]
    }
  }
}`,
		},
		{
			Name:        "HealthcareQuery",
			Description: "Healthcare patient records",
			JSON: `{
  "patients": {
    "select": ["id", "first_name", "last_name", "dob", "gender", "blood_type", "insurance_id"],
    "where": {
      "and": [
        { "status": "active" },
        { 
          "or": [
            { "age": {">=": 65} },
            { "has_chronic_condition": true }
          ]
        }
      ]
    },
    "order": ["last_name", "first_name"],
    "limit": 100,
    "appointments": {
      "select": ["id", "date", "time", "type", "status", "notes"],
      "join": "patient_id:id",
      "where": {
        "date": {">=": "2023-01-01"}
      },
      "order": ["-date", "-time"],
      "doctor": {
        "select": ["id", "name", "specialty", "license_number"],
        "join": "doctor_id:id"
      }
    },
    "medical_records": {
      "select": ["id", "record_date", "diagnosis", "treatment", "notes"],
      "join": "patient_id:id",
      "order": ["-record_date"],
      "prescriptions": {
        "select": ["id", "medication", "dosage", "frequency", "start_date", "end_date"],
        "join": "medical_record_id:id"
      }
    },
    "lab_results": {
      "select": ["id", "test_name", "result", "reference_range", "date", "is_abnormal"],
      "join": "patient_id:id",
      "where": {
        "date": {">=": "2023-06-01"}
      },
      "order": ["-date"]
    },
    "emergency_contacts": {
      "select": ["id", "name", "relationship", "phone", "email"],
      "join": "patient_id:id"
    }
  }
}`,
		},
	}
}

// SaveTestDataToFiles saves all test data to JSON files in the specified directory
func SaveTestDataToFiles(dir string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Get test data
	testData := GetTestData()

	// Save each test case to a file
	for _, td := range testData {
		filename := dir + "/" + td.Name + ".json"
		if err := os.WriteFile(filename, []byte(td.JSON), 0644); err != nil {
			return err
		}
	}

	return nil
}
