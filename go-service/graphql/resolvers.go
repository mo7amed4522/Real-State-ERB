package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"my-property/go-service/database"
	"my-property/go-service/models"
	"my-property/go-service/pubsub"
	"my-property/go-service/services"
	"my-property/go-service/utils"
	"net/http"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

var propertyType *graphql.Object
var rootQuery *graphql.Object
var rootMutation *graphql.Object
var commentType *graphql.Object
var buildingType *graphql.Object
var financialResolvers *FinancialResolvers

// InitializeResolvers initializes all resolvers with their dependencies
func InitializeResolvers(financialService *services.FinancialService) {
	financialResolvers = NewFinancialResolvers(financialService)
	init()
}

func init() {
	commentType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.ID},
			"content":    &graphql.Field{Type: graphql.String},
			"userId":     &graphql.Field{Type: graphql.String},
			"likesCount": &graphql.Field{Type: graphql.Int},
			"parentId":   &graphql.Field{Type: graphql.ID},
			"replies": &graphql.Field{
				Type: graphql.NewList(commentType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if comment, ok := p.Source.(models.Comment); ok {
						var replies []models.Comment
						database.DB.Model(&comment).Association("Replies").Find(&replies)
						return replies, nil
					}
					return nil, nil
				},
			},
			"createdAt":  &graphql.Field{Type: graphql.String},
			"updatedAt":  &graphql.Field{Type: graphql.String},
		},
	})

	propertyType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Property",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.ID},
			"title":       &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"type":        &graphql.Field{Type: graphql.String},
			"status":      &graphql.Field{Type: graphql.String},
			"price":       &graphql.Field{Type: graphql.Float},
			"currency":    &graphql.Field{Type: graphql.String},
			"bedrooms":    &graphql.Field{Type: graphql.Int},
			"bathrooms":   &graphql.Field{Type: graphql.Int},
			"area":        &graphql.Field{Type: graphql.Float},
			"furnished":   &graphql.Field{Type: graphql.Boolean},
			"amenities": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if property, ok := p.Source.(models.Property); ok {
						return strings.Split(property.Amenities, ","), nil
					}
					return nil, nil
				},
			},
			"images": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if property, ok := p.Source.(models.Property); ok {
						return strings.Split(property.Images, ","), nil
					}
					return nil, nil
				},
			},
			"isFeatured": &graphql.Field{Type: graphql.Boolean},
			"views":      &graphql.Field{Type: graphql.Int},
			"favoritesCount": &graphql.Field{Type: graphql.Int},
			"comments": &graphql.Field{
				Type: graphql.NewList(commentType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if property, ok := p.Source.(models.Property); ok {
						var comments []models.Comment
						database.DB.Model(&property).Association("Comments").Find(&comments)
						return comments, nil
					}
					return nil, nil
				},
			},
			"createdAt":  &graphql.Field{Type: graphql.String},
			"updatedAt":  &graphql.Field{Type: graphql.String},
			"landlordId": &graphql.Field{Type: graphql.String},
		},
	})

	propertyInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "PropertyInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"title":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"description": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"type":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			"status":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"price":       &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"currency":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"bedrooms":    &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"bathrooms":   &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"area":        &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"furnished":   &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"amenities":   &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
			"images":      &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
			"isFeatured":  &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"landlordId":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	buildingType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Building",
			Fields: graphql.Fields{
				"id":             &graphql.Field{Type: graphql.ID},
				"title":          &graphql.Field{Type: graphql.String},
				"description":    &graphql.Field{Type: graphql.String},
				"latitude":       &graphql.Field{Type: graphql.Float},
				"longitude":      &graphql.Field{Type: graphql.Float},
				"address":        &graphql.Field{Type: graphql.String},
				"city":           &graphql.Field{Type: graphql.String},
				"region":         &graphql.Field{Type: graphql.String},
				"price":          &graphql.Field{Type: graphql.Float},
				"status":         &graphql.Field{Type: graphql.String},
				"sold_at":        &graphql.Field{Type: graphql.String},
				"company_id":     &graphql.Field{Type: graphql.ID},
				"developer_id":   &graphql.Field{Type: graphql.ID},
				"total_likes":    &graphql.Field{Type: graphql.Int},
				"total_comments": &graphql.Field{Type: graphql.Int},
				"total_views":    &graphql.Field{Type: graphql.Int},
				"created_at":     &graphql.Field{Type: graphql.String},
				"updated_at":     &graphql.Field{Type: graphql.String},
			},
		},
	)

	rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"properties": &graphql.Field{
				Type: graphql.NewList(propertyType),
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{Type: graphql.NewInputObject(graphql.InputObjectConfig{
						Name: "PropertyFilterInputArg",
						Fields: graphql.InputObjectConfigFieldMap{
							"isVerified": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
							"listedBy": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"ownershipType": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"rentalPeriod": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"tags": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
							"neighborhood": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"buildingName": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"minPrice": &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"maxPrice": &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"minBedrooms": &graphql.InputObjectFieldConfig{Type: graphql.Int},
							"maxBedrooms": &graphql.InputObjectFieldConfig{Type: graphql.Int},
							"minBathrooms": &graphql.InputObjectFieldConfig{Type: graphql.Int},
							"maxBathrooms": &graphql.InputObjectFieldConfig{Type: graphql.Int},
							"minArea": &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"maxArea": &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"furnished": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
							"status": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"type": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"ratingMin": &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"ratingMax": &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"isFeatured": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
							"boosted": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
						},
					})},
					"sortBy": &graphql.ArgumentConfig{Type: graphql.String},
					"sortOrder": &graphql.ArgumentConfig{Type: graphql.String},
					"limit": &graphql.ArgumentConfig{Type: graphql.Int},
					"offset": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var properties []models.Property
					db := database.DB

					if filter, ok := p.Args["filter"].(map[string]interface{}); ok && filter != nil {
						db = applyPropertyFilters(db, filter)
					}
					if sortBy, ok := p.Args["sortBy"].(string); ok && sortBy != "" {
						order := "asc"
						if sortOrder, ok := p.Args["sortOrder"].(string); ok && (sortOrder == "desc" || sortOrder == "DESC") {
							order = "desc"
						}
						db = db.Order(sortBy + " " + order)
					}
					if limit, ok := p.Args["limit"].(int); ok && limit > 0 {
						db = db.Limit(limit)
					}
					if offset, ok := p.Args["offset"].(int); ok && offset > 0 {
						db = db.Offset(offset)
					}
					db.Find(&properties)
					return properties, nil
				},
			},
			"property": &graphql.Field{
				Type: propertyType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					var property models.Property
					if err := database.DB.First(&property, id).Error; err != nil {
						return nil, err
					}
					return property, nil
				},
			},
			"building": &graphql.Field{
				Type: buildingType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(string)
					var building models.Building
					database.DB.First(&building, id)
					return building, nil
				},
			},
			"buildings": &graphql.Field{
				Type: graphql.NewList(buildingType),
				Args: graphql.FieldConfigArgument{
					"city":   &graphql.ArgumentConfig{Type: graphql.String},
					"region": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var buildings []models.Building
					query := database.DB
					if city, ok := params.Args["city"].(string); ok {
						query = query.Where("city = ?", city)
					}
					if region, ok := params.Args["region"].(string); ok {
						query = query.Where("region = ?", region)
					}
					query.Find(&buildings)
					return buildings, nil
				},
			},
			// Financial Queries
			"getSaleTransactions": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: "SaleTransaction",
					Fields: graphql.Fields{
						"id":            &graphql.Field{Type: graphql.ID},
						"buildingId":    &graphql.Field{Type: graphql.ID},
						"buyerId":       &graphql.Field{Type: graphql.ID},
						"sellerId":      &graphql.Field{Type: graphql.ID},
						"agentId":       &graphql.Field{Type: graphql.ID},
						"price":         &graphql.Field{Type: graphql.Float},
						"paymentMethod": &graphql.Field{Type: graphql.String},
						"status":        &graphql.Field{Type: graphql.String},
						"commission":    &graphql.Field{Type: graphql.Float},
						"taxAmount":     &graphql.Field{Type: graphql.Float},
						"fees":          &graphql.Field{Type: graphql.Float},
						"totalAmount":   &graphql.Field{Type: graphql.Float},
						"notes":         &graphql.Field{Type: graphql.String},
						"createdAt":     &graphql.Field{Type: graphql.String},
						"updatedAt":     &graphql.Field{Type: graphql.String},
						"completedAt":   &graphql.Field{Type: graphql.String},
					},
				})),
				Args: graphql.FieldConfigArgument{
					"buildingId": &graphql.ArgumentConfig{Type: graphql.ID},
					"buyerId":    &graphql.ArgumentConfig{Type: graphql.ID},
					"sellerId":   &graphql.ArgumentConfig{Type: graphql.ID},
					"agentId":    &graphql.ArgumentConfig{Type: graphql.ID},
					"status":     &graphql.ArgumentConfig{Type: graphql.String},
					"startDate":  &graphql.ArgumentConfig{Type: graphql.String},
					"endDate":    &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.GetSaleTransactions(p)
				},
			},
			"getLeaseContracts": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: "LeaseContract",
					Fields: graphql.Fields{
						"id":                &graphql.Field{Type: graphql.ID},
						"propertyId":        &graphql.Field{Type: graphql.ID},
						"tenantId":          &graphql.Field{Type: graphql.ID},
						"landlordId":        &graphql.Field{Type: graphql.ID},
						"agentId":           &graphql.Field{Type: graphql.ID},
						"durationMonths":    &graphql.Field{Type: graphql.Int},
						"startDate":         &graphql.Field{Type: graphql.String},
						"endDate":           &graphql.Field{Type: graphql.String},
						"monthlyRent":       &graphql.Field{Type: graphql.Float},
						"depositAmount":     &graphql.Field{Type: graphql.Float},
						"paymentFrequency":  &graphql.Field{Type: graphql.String},
						"status":            &graphql.Field{Type: graphql.String},
						"contractFileUrl":   &graphql.Field{Type: graphql.String},
						"contractFileKey":   &graphql.Field{Type: graphql.String},
						"utilitiesIncluded": &graphql.Field{Type: graphql.Boolean},
						"petAllowed":        &graphql.Field{Type: graphql.Boolean},
						"furnished":         &graphql.Field{Type: graphql.Boolean},
						"lateFeeAmount":     &graphql.Field{Type: graphql.Float},
						"gracePeriodDays":   &graphql.Field{Type: graphql.Int},
						"notes":             &graphql.Field{Type: graphql.String},
						"createdAt":         &graphql.Field{Type: graphql.String},
						"updatedAt":         &graphql.Field{Type: graphql.String},
						"signedAt":          &graphql.Field{Type: graphql.String},
						"terminatedAt":      &graphql.Field{Type: graphql.String},
					},
				})),
				Args: graphql.FieldConfigArgument{
					"propertyId": &graphql.ArgumentConfig{Type: graphql.ID},
					"tenantId":   &graphql.ArgumentConfig{Type: graphql.ID},
					"landlordId": &graphql.ArgumentConfig{Type: graphql.ID},
					"status":     &graphql.ArgumentConfig{Type: graphql.String},
					"activeOnly": &graphql.ArgumentConfig{Type: graphql.Boolean},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.GetLeaseContracts(p)
				},
			},
			"getOffers": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: "Offer",
					Fields: graphql.Fields{
						"id":              &graphql.Field{Type: graphql.ID},
						"title":           &graphql.Field{Type: graphql.String},
						"description":     &graphql.Field{Type: graphql.String},
						"discountPercent": &graphql.Field{Type: graphql.Int},
						"discountAmount":  &graphql.Field{Type: graphql.Float},
						"startDate":       &graphql.Field{Type: graphql.String},
						"endDate":         &graphql.Field{Type: graphql.String},
						"buildingId":      &graphql.Field{Type: graphql.ID},
						"companyId":       &graphql.Field{Type: graphql.ID},
						"imageUrl":        &graphql.Field{Type: graphql.String},
						"termsConditions": &graphql.Field{Type: graphql.String},
						"isActive":        &graphql.Field{Type: graphql.Boolean},
						"maxUses":         &graphql.Field{Type: graphql.Int},
						"currentUses":     &graphql.Field{Type: graphql.Int},
						"minAmount":       &graphql.Field{Type: graphql.Float},
						"maxAmount":       &graphql.Field{Type: graphql.Float},
						"code":            &graphql.Field{Type: graphql.String},
						"createdAt":       &graphql.Field{Type: graphql.String},
						"updatedAt":       &graphql.Field{Type: graphql.String},
					},
				})),
				Args: graphql.FieldConfigArgument{
					"buildingId": &graphql.ArgumentConfig{Type: graphql.ID},
					"companyId":  &graphql.ArgumentConfig{Type: graphql.ID},
					"activeOnly": &graphql.ArgumentConfig{Type: graphql.Boolean},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.GetOffers(p)
				},
			},
		},
	})

	rootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createProperty": &graphql.Field{
				Type: propertyType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(propertyInputType)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})

					langs := []string{"en", "ar", "fr", "de", "hi", "ru", "fil"}
					title := input["title"].(string)
					description := input["description"].(string)
					titleTranslations, err := utils.TranslateText(title, "en", langs)
					if err != nil {
						return nil, fmt.Errorf("translation error: %v", err)
					}
					descriptionTranslations, err := utils.TranslateText(description, "en", langs)
					if err != nil {
						return nil, fmt.Errorf("translation error: %v", err)
					}

					amenities, _ := input["amenities"].([]interface{})
					amenitiesStr := make([]string, len(amenities))
					for i, v := range amenities {
						amenitiesStr[i] = v.(string)
					}

					images, _ := input["images"].([]interface{})
					imagesStr := make([]string, len(images))
					for i, v := range images {
						imagesStr[i] = v.(string)
					}

					property := models.Property{
						Title:       titleTranslations,
						Description: descriptionTranslations,
						Type:        input["type"].(string),
						Status:      input["status"].(string),
						Price:       input["price"].(float64),
						Currency:    input["currency"].(string),
						Bedrooms:    int(input["bedrooms"].(float64)),
						Bathrooms:   int(input["bathrooms"].(float64)),
						Area:        input["area"].(float64),
						Furnished:   input["furnished"].(bool),
						Amenities:   strings.Join(amenitiesStr, ","),
						Images:      strings.Join(imagesStr, ","),
						IsFeatured:  input["isFeatured"].(bool),
						LandlordId:  input["landlordId"].(string),
						// Ownership & Availability
						ListedBy:         getString(input, "listedBy"),
						AvailabilityDate: parseTimePtr(getString(input, "availabilityDate")),
						IsVerified:       getBool(input, "isVerified"),
						OwnershipType:    getString(input, "ownershipType"),
						RentalPeriod:     getString(input, "rentalPeriod"),
						DepositRequired:  getBool(input, "depositRequired"),
						DepositAmount:    getFloat(input, "depositAmount"),
						CommissionAmount: getFloat(input, "commissionAmount"),
						// Digital Location Enhancements
						GoogleMapsLink:  getString(input, "googleMapsLink"),
						Neighborhood:    getString(input, "neighborhood"),
						NearbyLandmarks: joinStringArray(input, "nearbyLandmarks"),
						FloorNumber:     getInt(input, "floorNumber"),
						BuildingName:    getString(input, "buildingName"),
						// Building & Facility Info
						YearBuilt:      getInt(input, "yearBuilt"),
						ParkingSpaces:  getInt(input, "parkingSpaces"),
						Balcony:        getBool(input, "balcony"),
						Elevator:       getBool(input, "elevator"),
						MaintenanceFee: getFloat(input, "maintenanceFee"),
						FloorPlan:      getString(input, "floorPlan"),
						// Media Attachments
						VideoTourUrl:      getString(input, "videoTourUrl"),
						VirtualTour360Url: getString(input, "virtualTour360Url"),
						FloorPlanUrl:      getString(input, "floorPlanUrl"),
						Documents:         joinStringArray(input, "documents"),
						// Analytics & Engagement
						InquiryCount:  getInt(input, "inquiryCount"),
						LastViewedAt:  parseTimePtr(getString(input, "lastViewedAt")),
						BoostedUntil:  parseTimePtr(getString(input, "boostedUntil")),
						Rating:        getFloat(input, "rating"),
						Tags:          joinStringArray(input, "tags"),
					}
					database.DB.Create(&property)
					return property, nil
				},
			},
			"updateProperty": &graphql.Field{
				Type: propertyType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(propertyInputType)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					input := p.Args["input"].(map[string]interface{})
					var property models.Property
					if err := database.DB.First(&property, "id = ?", id).Error; err != nil {
						return nil, err
					}

					langs := []string{"en", "ar", "fr", "de", "hi", "ru", "fil"}
					if title, ok := input["title"]; ok {
						titleTranslations, err := utils.TranslateText(title.(string), "en", langs)
						if err != nil {
							return nil, fmt.Errorf("translation error: %v", err)
						}
						property.Title = titleTranslations
					}
					if description, ok := input["description"]; ok {
						descriptionTranslations, err := utils.TranslateText(description.(string), "en", langs)
						if err != nil {
							return nil, fmt.Errorf("translation error: %v", err)
						}
						property.Description = descriptionTranslations
					}
					if pType, ok := input["type"]; ok {
						property.Type = pType.(string)
					}
					if status, ok := input["status"]; ok {
						property.Status = status.(string)
					}
					if price, ok := input["price"]; ok {
						property.Price = price.(float64)
					}
					if currency, ok := input["currency"]; ok {
						property.Currency = currency.(string)
					}
					if bedrooms, ok := input["bedrooms"]; ok {
						property.Bedrooms = int(bedrooms.(float64))
					}
					if bathrooms, ok := input["bathrooms"]; ok {
						property.Bathrooms = int(bathrooms.(float64))
					}
					if area, ok := input["area"]; ok {
						property.Area = area.(float64)
					}
					if furnished, ok := input["furnished"]; ok {
						property.Furnished = furnished.(bool)
					}
					if amenities, ok := input["amenities"]; ok {
						amenitiesList := amenities.([]interface{})
						amenitiesStr := make([]string, len(amenitiesList))
						for i, v := range amenitiesList {
							amenitiesStr[i] = v.(string)
						}
						property.Amenities = strings.Join(amenitiesStr, ",")
					}
					if images, ok := input["images"]; ok {
						imagesList := images.([]interface{})
						imagesStr := make([]string, len(imagesList))
						for i, v := range imagesList {
							imagesStr[i] = v.(string)
						}
						property.Images = strings.Join(imagesStr, ",")
					}
					if isFeatured, ok := input["isFeatured"]; ok {
						property.IsFeatured = isFeatured.(bool)
					}
					if landlordId, ok := input["landlordId"]; ok {
						property.LandlordId = landlordId.(string)
					}
					// Ownership & Availability
					if listedBy, ok := input["listedBy"]; ok {
						property.ListedBy = listedBy.(string)
					}
					if availabilityDate, ok := input["availabilityDate"]; ok {
						property.AvailabilityDate = parseTimePtr(availabilityDate.(string))
					}
					if isVerified, ok := input["isVerified"]; ok {
						property.IsVerified = isVerified.(bool)
					}
					if ownershipType, ok := input["ownershipType"]; ok {
						property.OwnershipType = ownershipType.(string)
					}
					if rentalPeriod, ok := input["rentalPeriod"]; ok {
						property.RentalPeriod = rentalPeriod.(string)
					}
					if depositRequired, ok := input["depositRequired"]; ok {
						property.DepositRequired = depositRequired.(bool)
					}
					if depositAmount, ok := input["depositAmount"]; ok {
						property.DepositAmount = depositAmount.(float64)
					}
					if commissionAmount, ok := input["commissionAmount"]; ok {
						property.CommissionAmount = commissionAmount.(float64)
					}
					// Digital Location Enhancements
					if googleMapsLink, ok := input["googleMapsLink"]; ok {
						property.GoogleMapsLink = googleMapsLink.(string)
					}
					if neighborhood, ok := input["neighborhood"]; ok {
						property.Neighborhood = neighborhood.(string)
					}
					if nearbyLandmarks, ok := input["nearbyLandmarks"]; ok {
						property.NearbyLandmarks = joinStringArray(input, "nearbyLandmarks")
					}
					if floorNumber, ok := input["floorNumber"]; ok {
						property.FloorNumber = int(floorNumber.(float64))
					}
					if buildingName, ok := input["buildingName"]; ok {
						property.BuildingName = buildingName.(string)
					}
					// Building & Facility Info
					if yearBuilt, ok := input["yearBuilt"]; ok {
						property.YearBuilt = int(yearBuilt.(float64))
					}
					if parkingSpaces, ok := input["parkingSpaces"]; ok {
						property.ParkingSpaces = int(parkingSpaces.(float64))
					}
					if balcony, ok := input["balcony"]; ok {
						property.Balcony = balcony.(bool)
					}
					if elevator, ok := input["elevator"]; ok {
						property.Elevator = elevator.(bool)
					}
					if maintenanceFee, ok := input["maintenanceFee"]; ok {
						property.MaintenanceFee = maintenanceFee.(float64)
					}
					if floorPlan, ok := input["floorPlan"]; ok {
						property.FloorPlan = floorPlan.(string)
					}
					// Media Attachments
					if videoTourUrl, ok := input["videoTourUrl"]; ok {
						property.VideoTourUrl = videoTourUrl.(string)
					}
					if virtualTour360Url, ok := input["virtualTour360Url"]; ok {
						property.VirtualTour360Url = virtualTour360Url.(string)
					}
					if floorPlanUrl, ok := input["floorPlanUrl"]; ok {
						property.FloorPlanUrl = floorPlanUrl.(string)
					}
					if documents, ok := input["documents"]; ok {
						property.Documents = joinStringArray(input, "documents")
					}
					// Analytics & Engagement
					if inquiryCount, ok := input["inquiryCount"]; ok {
						property.InquiryCount = int(inquiryCount.(float64))
					}
					if lastViewedAt, ok := input["lastViewedAt"]; ok {
						property.LastViewedAt = parseTimePtr(lastViewedAt.(string))
					}
					if boostedUntil, ok := input["boostedUntil"]; ok {
						property.BoostedUntil = parseTimePtr(boostedUntil.(string))
					}
					if rating, ok := input["rating"]; ok {
						property.Rating = rating.(float64)
					}
					if tags, ok := input["tags"]; ok {
						property.Tags = joinStringArray(input, "tags")
					}

					database.DB.Save(&property)
					return property, nil
				},
			},
			"deleteProperty": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					var property models.Property
					if err := database.DB.First(&property, "id = ?", id).Error; err != nil {
						return nil, err
					}
					database.DB.Delete(&property)
					return "Property deleted successfully", nil
				},
			},
			"incrementFavoriteCount": &graphql.Field{
				Type: propertyType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					var property models.Property
					if err := database.DB.First(&property, "id = ?", id).Error; err != nil {
						return nil, err
					}
					property.FavoritesCount++
					database.DB.Save(&property)
					return property, nil
				},
			},
			"decrementFavoriteCount": &graphql.Field{
				Type: propertyType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					var property models.Property
					if err := database.DB.First(&property, "id = ?", id).Error; err != nil {
						return nil, err
					}
					if property.FavoritesCount > 0 {
						property.FavoritesCount--
					}
					database.DB.Save(&property)
					return property, nil
				},
			},
			"createComment": &graphql.Field{
				Type: commentType,
				Args: graphql.FieldConfigArgument{
					"propertyId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"content":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"userId":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"parentId":   &graphql.ArgumentConfig{Type: graphql.ID},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					propertyId, _ := p.Args["propertyId"].(string)
					content, _ := p.Args["content"].(string)
					userId, _ := p.Args["userId"].(string)
					var parentIdPtr *uint
					if parentIdRaw, ok := p.Args["parentId"]; ok && parentIdRaw != nil {
						parentIdStr := parentIdRaw.(string)
						var parentId uint
						_, err := fmt.Sscanf(parentIdStr, "%d", &parentId)
						if err == nil {
							parentIdPtr = &parentId
						}
					}

					var property models.Property
					if err := database.DB.First(&property, "id = ?", propertyId).Error; err != nil {
						return nil, fmt.Errorf("property not found")
					}

					comment := models.Comment{
						PropertyID: property.ID,
						Content:    content,
						UserID:     userId,
						ParentID:   parentIdPtr,
					}
					database.DB.Create(&comment)
					return comment, nil
				},
			},
			"updateComment": &graphql.Field{
				Type: commentType,
				Args: graphql.FieldConfigArgument{
					"commentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"content":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					commentId, _ := p.Args["commentId"].(string)
					content, _ := p.Args["content"].(string)

					var comment models.Comment
					if err := database.DB.First(&comment, "id = ?", commentId).Error; err != nil {
						return nil, fmt.Errorf("comment not found")
					}

					comment.Content = content
					database.DB.Save(&comment)
					return comment, nil
				},
			},
			"deleteComment": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"commentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					commentId, _ := p.Args["commentId"].(string)
					var comment models.Comment
					if err := database.DB.First(&comment, "id = ?", commentId).Error; err != nil {
						return nil, fmt.Errorf("comment not found")
					}
					database.DB.Delete(&comment)
					return "Comment deleted successfully", nil
				},
			},
			"incrementCommentLike": &graphql.Field{
				Type: commentType,
				Args: graphql.FieldConfigArgument{
					"commentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					commentId, _ := p.Args["commentId"].(string)
					var comment models.Comment
					if err := database.DB.First(&comment, "id = ?", commentId).Error; err != nil {
						return nil, fmt.Errorf("comment not found")
					}
					comment.LikesCount++
					database.DB.Save(&comment)
					return comment, nil
				},
			},
			"decrementCommentLike": &graphql.Field{
				Type: commentType,
				Args: graphql.FieldConfigArgument{
					"commentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					commentId, _ := p.Args["commentId"].(string)
					var comment models.Comment
					if err := database.DB.First(&comment, "id = ?", commentId).Error; err != nil {
						return nil, fmt.Errorf("comment not found")
					}
					if comment.LikesCount > 0 {
						comment.LikesCount--
					}
					database.DB.Save(&comment)
					return comment, nil
				},
			},
			"createBuilding": &graphql.Field{
				Type: buildingType,
				Args: graphql.FieldConfigArgument{
					"title":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"description":  &graphql.ArgumentConfig{Type: graphql.String},
					"latitude":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
					"longitude":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
					"address":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"city":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"region":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"price":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
					"company_id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"developer_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					building := models.Building{
						Title:       params.Args["title"].(string),
						Description: params.Args["description"].(string),
						Latitude:    params.Args["latitude"].(float64),
						Longitude:   params.Args["longitude"].(float64),
						Address:     params.Args["address"].(string),
						City:        params.Args["city"].(string),
						Region:      params.Args["region"].(string),
						Price:       params.Args["price"].(float64),
						Status:      "AVAILABLE",
					}
					database.DB.Create(&building)
					pubsub.GetInstance().Publish("BUILDING_ADDED", building)
					return building, nil
				},
			},
			// Financial Mutations
			"createSaleTransaction": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "SaleTransaction",
					Fields: graphql.Fields{
						"id":            &graphql.Field{Type: graphql.ID},
						"buildingId":    &graphql.Field{Type: graphql.ID},
						"buyerId":       &graphql.Field{Type: graphql.ID},
						"sellerId":      &graphql.Field{Type: graphql.ID},
						"agentId":       &graphql.Field{Type: graphql.ID},
						"price":         &graphql.Field{Type: graphql.Float},
						"paymentMethod": &graphql.Field{Type: graphql.String},
						"status":        &graphql.Field{Type: graphql.String},
						"commission":    &graphql.Field{Type: graphql.Float},
						"taxAmount":     &graphql.Field{Type: graphql.Float},
						"fees":          &graphql.Field{Type: graphql.Float},
						"totalAmount":   &graphql.Field{Type: graphql.Float},
						"notes":         &graphql.Field{Type: graphql.String},
						"createdAt":     &graphql.Field{Type: graphql.String},
						"updatedAt":     &graphql.Field{Type: graphql.String},
						"completedAt":   &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewInputObject(graphql.InputObjectConfig{
						Name: "CreateSaleTransactionInput",
						Fields: graphql.InputObjectConfigFieldMap{
							"buildingId":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
							"buyerId":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
							"sellerId":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
							"agentId":       &graphql.InputObjectFieldConfig{Type: graphql.ID},
							"price":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
							"paymentMethod": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"commission":    &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"taxAmount":     &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"fees":          &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"notes":         &graphql.InputObjectFieldConfig{Type: graphql.String},
						},
					})},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.CreateSaleTransaction(p)
				},
			},
			"updateSaleTransactionStatus": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "SaleTransaction",
					Fields: graphql.Fields{
						"id":            &graphql.Field{Type: graphql.ID},
						"buildingId":    &graphql.Field{Type: graphql.ID},
						"buyerId":       &graphql.Field{Type: graphql.ID},
						"sellerId":      &graphql.Field{Type: graphql.ID},
						"agentId":       &graphql.Field{Type: graphql.ID},
						"price":         &graphql.Field{Type: graphql.Float},
						"paymentMethod": &graphql.Field{Type: graphql.String},
						"status":        &graphql.Field{Type: graphql.String},
						"commission":    &graphql.Field{Type: graphql.Float},
						"taxAmount":     &graphql.Field{Type: graphql.Float},
						"fees":          &graphql.Field{Type: graphql.Float},
						"totalAmount":   &graphql.Field{Type: graphql.Float},
						"notes":         &graphql.Field{Type: graphql.String},
						"createdAt":     &graphql.Field{Type: graphql.String},
						"updatedAt":     &graphql.Field{Type: graphql.String},
						"completedAt":   &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"status": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.UpdateSaleTransactionStatus(p)
				},
			},
			"createLeaseContract": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "LeaseContract",
					Fields: graphql.Fields{
						"id":                &graphql.Field{Type: graphql.ID},
						"propertyId":        &graphql.Field{Type: graphql.ID},
						"tenantId":          &graphql.Field{Type: graphql.ID},
						"landlordId":        &graphql.Field{Type: graphql.ID},
						"agentId":           &graphql.Field{Type: graphql.ID},
						"durationMonths":    &graphql.Field{Type: graphql.Int},
						"startDate":         &graphql.Field{Type: graphql.String},
						"endDate":           &graphql.Field{Type: graphql.String},
						"monthlyRent":       &graphql.Field{Type: graphql.Float},
						"depositAmount":     &graphql.Field{Type: graphql.Float},
						"paymentFrequency":  &graphql.Field{Type: graphql.String},
						"status":            &graphql.Field{Type: graphql.String},
						"contractFileUrl":   &graphql.Field{Type: graphql.String},
						"contractFileKey":   &graphql.Field{Type: graphql.String},
						"utilitiesIncluded": &graphql.Field{Type: graphql.Boolean},
						"petAllowed":        &graphql.Field{Type: graphql.Boolean},
						"furnished":         &graphql.Field{Type: graphql.Boolean},
						"lateFeeAmount":     &graphql.Field{Type: graphql.Float},
						"gracePeriodDays":   &graphql.Field{Type: graphql.Int},
						"notes":             &graphql.Field{Type: graphql.String},
						"createdAt":         &graphql.Field{Type: graphql.String},
						"updatedAt":         &graphql.Field{Type: graphql.String},
						"signedAt":          &graphql.Field{Type: graphql.String},
						"terminatedAt":      &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewInputObject(graphql.InputObjectConfig{
						Name: "CreateLeaseContractInput",
						Fields: graphql.InputObjectConfigFieldMap{
							"propertyId":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
							"tenantId":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
							"landlordId":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
							"agentId":          &graphql.InputObjectFieldConfig{Type: graphql.ID},
							"durationMonths":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
							"startDate":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"monthlyRent":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
							"depositAmount":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
							"paymentFrequency": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"utilitiesIncluded": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
							"petAllowed":       &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
							"furnished":        &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
							"lateFeeAmount":    &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"gracePeriodDays":  &graphql.InputObjectFieldConfig{Type: graphql.Int},
							"notes":            &graphql.InputObjectFieldConfig{Type: graphql.String},
						},
					})},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.CreateLeaseContract(p)
				},
			},
			"updateLeaseStatus": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "LeaseContract",
					Fields: graphql.Fields{
						"id":                &graphql.Field{Type: graphql.ID},
						"propertyId":        &graphql.Field{Type: graphql.ID},
						"tenantId":          &graphql.Field{Type: graphql.ID},
						"landlordId":        &graphql.Field{Type: graphql.ID},
						"agentId":           &graphql.Field{Type: graphql.ID},
						"durationMonths":    &graphql.Field{Type: graphql.Int},
						"startDate":         &graphql.Field{Type: graphql.String},
						"endDate":           &graphql.Field{Type: graphql.String},
						"monthlyRent":       &graphql.Field{Type: graphql.Float},
						"depositAmount":     &graphql.Field{Type: graphql.Float},
						"paymentFrequency":  &graphql.Field{Type: graphql.String},
						"status":            &graphql.Field{Type: graphql.String},
						"contractFileUrl":   &graphql.Field{Type: graphql.String},
						"contractFileKey":   &graphql.Field{Type: graphql.String},
						"utilitiesIncluded": &graphql.Field{Type: graphql.Boolean},
						"petAllowed":        &graphql.Field{Type: graphql.Boolean},
						"furnished":         &graphql.Field{Type: graphql.Boolean},
						"lateFeeAmount":     &graphql.Field{Type: graphql.Float},
						"gracePeriodDays":   &graphql.Field{Type: graphql.Int},
						"notes":             &graphql.Field{Type: graphql.String},
						"createdAt":         &graphql.Field{Type: graphql.String},
						"updatedAt":         &graphql.Field{Type: graphql.String},
						"signedAt":          &graphql.Field{Type: graphql.String},
						"terminatedAt":      &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"status": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.UpdateLeaseStatus(p)
				},
			},
			"createLeasePayment": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "LeasePayment",
					Fields: graphql.Fields{
						"id":            &graphql.Field{Type: graphql.ID},
						"leaseId":       &graphql.Field{Type: graphql.ID},
						"amount":        &graphql.Field{Type: graphql.Float},
						"paymentDate":   &graphql.Field{Type: graphql.String},
						"dueDate":       &graphql.Field{Type: graphql.String},
						"status":        &graphql.Field{Type: graphql.String},
						"paymentMethod": &graphql.Field{Type: graphql.String},
						"lateFee":       &graphql.Field{Type: graphql.Float},
						"notes":         &graphql.Field{Type: graphql.String},
						"createdAt":     &graphql.Field{Type: graphql.String},
						"updatedAt":     &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewInputObject(graphql.InputObjectConfig{
						Name: "CreateLeasePaymentInput",
						Fields: graphql.InputObjectConfigFieldMap{
							"leaseId":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
							"amount":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
							"paymentDate":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"dueDate":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"paymentMethod": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"notes":         &graphql.InputObjectFieldConfig{Type: graphql.String},
						},
					})},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.CreateLeasePayment(p)
				},
			},
			"updatePaymentStatus": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "LeasePayment",
					Fields: graphql.Fields{
						"id":            &graphql.Field{Type: graphql.ID},
						"leaseId":       &graphql.Field{Type: graphql.ID},
						"amount":        &graphql.Field{Type: graphql.Float},
						"paymentDate":   &graphql.Field{Type: graphql.String},
						"dueDate":       &graphql.Field{Type: graphql.String},
						"status":        &graphql.Field{Type: graphql.String},
						"paymentMethod": &graphql.Field{Type: graphql.String},
						"lateFee":       &graphql.Field{Type: graphql.Float},
						"notes":         &graphql.Field{Type: graphql.String},
						"createdAt":     &graphql.Field{Type: graphql.String},
						"updatedAt":     &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"status": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.UpdatePaymentStatus(p)
				},
			},
			"createOffer": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Offer",
					Fields: graphql.Fields{
						"id":              &graphql.Field{Type: graphql.ID},
						"title":           &graphql.Field{Type: graphql.String},
						"description":     &graphql.Field{Type: graphql.String},
						"discountPercent": &graphql.Field{Type: graphql.Int},
						"discountAmount":  &graphql.Field{Type: graphql.Float},
						"startDate":       &graphql.Field{Type: graphql.String},
						"endDate":         &graphql.Field{Type: graphql.String},
						"buildingId":      &graphql.Field{Type: graphql.ID},
						"companyId":       &graphql.Field{Type: graphql.ID},
						"imageUrl":        &graphql.Field{Type: graphql.String},
						"termsConditions": &graphql.Field{Type: graphql.String},
						"isActive":        &graphql.Field{Type: graphql.Boolean},
						"maxUses":         &graphql.Field{Type: graphql.Int},
						"currentUses":     &graphql.Field{Type: graphql.Int},
						"minAmount":       &graphql.Field{Type: graphql.Float},
						"maxAmount":       &graphql.Field{Type: graphql.Float},
						"code":            &graphql.Field{Type: graphql.String},
						"createdAt":       &graphql.Field{Type: graphql.String},
						"updatedAt":       &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewInputObject(graphql.InputObjectConfig{
						Name: "CreateOfferInput",
						Fields: graphql.InputObjectConfigFieldMap{
							"title":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"description":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"discountPercent": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
							"discountAmount":  &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"startDate":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"endDate":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
							"buildingId":      &graphql.InputObjectFieldConfig{Type: graphql.ID},
							"companyId":       &graphql.InputObjectFieldConfig{Type: graphql.ID},
							"imageUrl":        &graphql.InputObjectFieldConfig{Type: graphql.String},
							"termsConditions": &graphql.InputObjectFieldConfig{Type: graphql.String},
							"maxUses":         &graphql.InputObjectFieldConfig{Type: graphql.Int},
							"minAmount":       &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"maxAmount":       &graphql.InputObjectFieldConfig{Type: graphql.Float},
							"code":            &graphql.InputObjectFieldConfig{Type: graphql.String},
						},
					})},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.CreateOffer(p)
				},
			},
			"useOffer": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "OfferUse",
					Fields: graphql.Fields{
						"id":       &graphql.Field{Type: graphql.ID},
						"offerId":  &graphql.Field{Type: graphql.ID},
						"userId":   &graphql.Field{Type: graphql.ID},
						"amount":   &graphql.Field{Type: graphql.Float},
						"discount": &graphql.Field{Type: graphql.Float},
						"usedAt":   &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"code":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"amount": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.UseOffer(p)
				},
			},
			"generateFinancialReport": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "FinancialReport",
					Fields: graphql.Fields{
						"id":              &graphql.Field{Type: graphql.ID},
						"reportType":      &graphql.Field{Type: graphql.String},
						"periodStart":     &graphql.Field{Type: graphql.String},
						"periodEnd":       &graphql.Field{Type: graphql.String},
						"totalSales":      &graphql.Field{Type: graphql.Float},
						"totalRentals":    &graphql.Field{Type: graphql.Float},
						"totalRevenue":    &graphql.Field{Type: graphql.Float},
						"totalCommissions": &graphql.Field{Type: graphql.Float},
						"totalTaxes":      &graphql.Field{Type: graphql.Float},
						"totalFees":       &graphql.Field{Type: graphql.Float},
						"reportData":      &graphql.Field{Type: graphql.String},
						"generatedAt":     &graphql.Field{Type: graphql.String},
						"generatedBy":     &graphql.Field{Type: graphql.ID},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"reportType":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"periodStart": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"periodEnd":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"generatedBy": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if financialResolvers == nil {
						return nil, fmt.Errorf("financial service not initialized")
					}
					return financialResolvers.GenerateFinancialReport(p)
				},
			},
		},
	})

	rootSubscription := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootSubscription",
		Fields: graphql.Fields{
			"buildingAdded": &graphql.Field{
				Type: buildingType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if building, ok := p.Source.(models.Building); ok {
						return building, nil
					}
					return nil, nil
				},
				Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
					return pubsub.GetInstance().Subscribe("BUILDING_ADDED"), nil
				},
			},
		},
	})

	Schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:        rootQuery,
			Mutation:     rootMutation,
			Subscription: rootSubscription,
		},
	)
}

var Schema graphql.Schema

func ExecuteQuery(query string, variables map[string]interface{}) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         Schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        context.Background(),
	})
	return result
}

func ExecuteSubscription(query string, variables map[string]interface{}) <-chan *graphql.Result {
	params := graphql.Params{
		Schema:         Schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        context.Background(),
	}
	// Note: In a real-world app, you'd handle multiple subscribers gracefully.
	// This is a simplified example.
	return graphql.Subscribe(params)
}

// --- utility functions ---
func getString(input map[string]interface{}, key string) string {
	if v, ok := input[key]; ok && v != nil {
		return v.(string)
	}
	return ""
}
func getBool(input map[string]interface{}, key string) bool {
	if v, ok := input[key]; ok && v != nil {
		return v.(bool)
	}
	return false
}
func getFloat(input map[string]interface{}, key string) float64 {
	if v, ok := input[key]; ok && v != nil {
		return v.(float64)
	}
	return 0
}
func getInt(input map[string]interface{}, key string) int {
	if v, ok := input[key]; ok && v != nil {
		return int(v.(float64))
	}
	return 0
}
func joinStringArray(input map[string]interface{}, key string) string {
	if v, ok := input[key]; ok && v != nil {
		arr := v.([]interface{})
		strs := make([]string, len(arr))
		for i, s := range arr {
			strs[i] = s.(string)
		}
		return strings.Join(strs, ",")
	}
	return ""
}
func parseTimePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func applyPropertyFilters(db *gorm.DB, filter map[string]interface{}) *gorm.DB {
	// OR logic
	if orFilters, ok := filter["or"]; ok && orFilters != nil {
		ors := orFilters.([]interface{})
		orDBs := make([]*gorm.DB, 0, len(ors))
		for _, orF := range ors {
			orDBs = append(orDBs, applyPropertyFilters(database.DB, orF.(map[string]interface{})))
		}
		// Combine OR queries
		if len(orDBs) > 0 {
			combined := orDBs[0]
			for i := 1; i < len(orDBs); i++ {
				combined = combined.Or(orDBs[i])
			}
			return db.Where(combined)
		}
	}
	// AND logic (default)
	if v, ok := filter["isVerified"]; ok {
		db = db.Where("is_verified = ?", v)
	}
	if v, ok := filter["listedBy"]; ok {
		db = db.Where("listed_by = ?", v)
	}
	if v, ok := filter["ownershipType"]; ok {
		db = db.Where("ownership_type = ?", v)
	}
	if v, ok := filter["rentalPeriod"]; ok {
		db = db.Where("rental_period = ?", v)
	}
	if v, ok := filter["neighborhood"]; ok {
		db = db.Where("neighborhood = ?", v)
	}
	if v, ok := filter["neighborhoods"]; ok {
		arr := v.([]interface{})
		if len(arr) > 0 {
			vals := make([]string, len(arr))
			for i, s := range arr {
				vals[i] = s.(string)
			}
			db = db.Where("neighborhood IN ?", vals)
		}
	}
	if v, ok := filter["buildingName"]; ok {
		db = db.Where("building_name = ?", v)
	}
	if v, ok := filter["buildingNames"]; ok {
		arr := v.([]interface{})
		if len(arr) > 0 {
			vals := make([]string, len(arr))
			for i, s := range arr {
				vals[i] = s.(string)
			}
			db = db.Where("building_name IN ?", vals)
		}
	}
	if v, ok := filter["minPrice"]; ok {
		db = db.Where("price >= ?", v)
	}
	if v, ok := filter["maxPrice"]; ok {
		db = db.Where("price <= ?", v)
	}
	if v, ok := filter["minBedrooms"]; ok {
		db = db.Where("bedrooms >= ?", v)
	}
	if v, ok := filter["maxBedrooms"]; ok {
		db = db.Where("bedrooms <= ?", v)
	}
	if v, ok := filter["minBathrooms"]; ok {
		db = db.Where("bathrooms >= ?", v)
	}
	if v, ok := filter["maxBathrooms"]; ok {
		db = db.Where("bathrooms <= ?", v)
	}
	if v, ok := filter["minArea"]; ok {
		db = db.Where("area >= ?", v)
	}
	if v, ok := filter["maxArea"]; ok {
		db = db.Where("area <= ?", v)
	}
	if v, ok := filter["furnished"]; ok {
		db = db.Where("furnished = ?", v)
	}
	if v, ok := filter["status"]; ok {
		db = db.Where("status = ?", v)
	}
	if v, ok := filter["statuses"]; ok {
		arr := v.([]interface{})
		if len(arr) > 0 {
			vals := make([]string, len(arr))
			for i, s := range arr {
				vals[i] = s.(string)
			}
			db = db.Where("status IN ?", vals)
		}
	}
	if v, ok := filter["type"]; ok {
		db = db.Where("type = ?", v)
	}
	if v, ok := filter["types"]; ok {
		arr := v.([]interface{})
		if len(arr) > 0 {
			vals := make([]string, len(arr))
			for i, s := range arr {
				vals[i] = s.(string)
			}
			db = db.Where("type IN ?", vals)
		}
	}
	if v, ok := filter["ratingMin"]; ok {
		db = db.Where("rating >= ?", v)
	}
	if v, ok := filter["ratingMax"]; ok {
		db = db.Where("rating <= ?", v)
	}
	if v, ok := filter["isFeatured"]; ok {
		db = db.Where("is_featured = ?", v)
	}
	if v, ok := filter["boosted"]; ok {
		db = db.Where("boosted_until > ?", time.Now())
	}
	if v, ok := filter["availabilityDate"]; ok {
		db = db.Where("availability_date <= ?", v)
	}
	if v, ok := filter["minYearBuilt"]; ok {
		db = db.Where("year_built >= ?", v)
	}
	if v, ok := filter["maxYearBuilt"]; ok {
		db = db.Where("year_built <= ?", v)
	}
	if v, ok := filter["minParkingSpaces"]; ok {
		db = db.Where("parking_spaces >= ?", v)
	}
	if v, ok := filter["maxParkingSpaces"]; ok {
		db = db.Where("parking_spaces <= ?", v)
	}
	if v, ok := filter["balcony"]; ok {
		db = db.Where("balcony = ?", v)
	}
	if v, ok := filter["elevator"]; ok {
		db = db.Where("elevator = ?", v)
	}
	if v, ok := filter["minMaintenanceFee"]; ok {
		db = db.Where("maintenance_fee >= ?", v)
	}
	if v, ok := filter["maxMaintenanceFee"]; ok {
		db = db.Where("maintenance_fee <= ?", v)
	}
	if v, ok := filter["videoTourUrl"]; ok {
		db = db.Where("video_tour_url = ?", v)
	}
	if v, ok := filter["virtualTour360Url"]; ok {
		db = db.Where("virtual_tour360_url = ?", v)
	}
	if v, ok := filter["documents"]; ok {
		arr := v.([]interface{})
		for _, doc := range arr {
			db = db.Where("documents LIKE ?", "%"+doc.(string)+"%")
		}
	}
	if v, ok := filter["minViewCount"]; ok {
		db = db.Where("views >= ?", v)
	}
	if v, ok := filter["maxViewCount"]; ok {
		db = db.Where("views <= ?", v)
	}
	if v, ok := filter["minInquiryCount"]; ok {
		db = db.Where("inquiry_count >= ?", v)
	}
	if v, ok := filter["maxInquiryCount"]; ok {
		db = db.Where("inquiry_count <= ?", v)
	}
	if v, ok := filter["lastViewedAt"]; ok {
		db = db.Where("last_viewed_at >= ?", v)
	}
	if v, ok := filter["boostedUntil"]; ok {
		db = db.Where("boosted_until >= ?", v)
	}
	if v, ok := filter["tags"]; ok {
		tags := v.([]interface{})
		for _, tag := range tags {
			db = db.Where("tags LIKE ?", "%"+tag.(string)+"%")
		}
	}
	return db
} 