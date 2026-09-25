// Package models defines domain entities, value objects, data transfer objects (DTOs),
// and error definitions for the Cinema Food & Beverage Recommender service.
//
// File organization:
//   - cinema.go         : Movie theaters, regional locations, and coordinates.
//   - concession.go     : Concession food/beverage items, categories, and dynamic live pricing.
//   - promotion.go      : Scope Privilege Club tiers, discount types, and savings breakdowns.
//   - recommendation.go : Recommendation request/result DTOs and cart evaluation payloads.
//   - payment.go        : Supported Sri Lankan payment methods (LankaQR, FriMi, eZ Cash, etc.) and checkout DTOs.
//   - ad.go             : In-theater lobby screen promotional advertisements and combo deals.
//   - errors.go         : Sentinel domain error values.
package models
