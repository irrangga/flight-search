# Flight Search API

A clean architecture flight search aggregation system that queries multiple airline providers in parallel, normalizes their responses, and returns unified flight data.

## Architecture

This project follows **Clean Code** and **Hexagonal Architecture** principles:

```
flight-search/
├── cmd/
│   └── api/
│       └── main.go                    # Application bootstrap
├── internal/                          # Private packages (Go convention)
│   ├── domain/
│   │   └── flight.go                  # Unified domain models
│   ├── provider/
│   │   ├── provider.go                # Provider interface
│   │   └── garuda.go, lion.go, ...    # Provider implementations
│   ├── normalizer/
│   │   ├── normalizer.go              # Normalization router
│   │   └── garuda.go, lion.go, ...    # Provider-specific normalizers
│   ├── aggregator/
│   │   └── aggregate.go               # Flight aggregation orchestrator
│   ├── mapper/
│   │   └── response.go                # Domain → DTO conversion
│   └── transport/
│       └── http.go                    # HTTP handler & routes
├── mock-api/                          # Mock provider responses
├── go.mod
└── README.md
```

## Package Responsibilities

- **domain**: Pure business domain models (no external dependencies)
- **provider**: Interfaces and implementations for external airline APIs
- **normalizer**: Converts provider-specific formats to unified domain models
- **aggregator**: Orchestrates parallel provider queries and aggregation
- **mapper**: Converts domain models to HTTP response DTOs
- **transport**: HTTP handler and request/response management
- **cmd/api**: Application bootstrap and dependency injection

## Running

### Build & Run

```bash
# Run directly
go run ./cmd/api

# Or build and run binary
go build -o flight-search ./cmd/api
./flight-search
```

Server starts on `http://localhost:8080`

## API Endpoint

### POST `/api/search-flights`

**Request:**

```json
{
  "origin": "CGK",
  "destination": "DPS",
  "departureDate": "2025-12-15",
  "passengers": 1,
  "cabinClass": "economy"
}
```

**Request with Filters:**

```json
{
  "origin": "CGK",
  "destination": "DPS",
  "departureDate": "2025-12-15",
  "passengers": 1,
  "cabinClass": "economy",
  "priceRange": [1000000, 2000000],
  "stopsRange": [0, 1],
  "departureTimeRange": ["06:00", "12:00"],
  "arrivalTimeRange": ["10:00", "18:00"],
  "airlines": ["AirAsia", "Garuda"],
  "durationRange": [60, 360],
  "sortBy": "price_asc"
}
```

**Response:**

```json
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "cabin_class": "economy"
  },
  "metadata": {
    "total_results": 12,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 450,
    "cache_hit": false
  },
  "flights": [
    {
      "id": "GA001_Garuda Indonesia",
      "provider": "Garuda Indonesia",
      "airline": { "name": "Garuda Indonesia", "code": "GA" },
      "flight_number": "GA001",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T08:00:00Z",
        "timestamp": 1702614000
      },
      "arrival": {
        "airport": "DPS",
        "city": "Bali",
        "datetime": "2025-12-15T11:00:00Z",
        "timestamp": 1702624800
      },
      "duration": {
        "total_minutes": 180,
        "formatted": "3h 0m"
      },
      "stops": 0,
      "price": { "amount": 1500000, "currency": "IDR" },
      "available_seats": 100,
      "cabin_class": "economy",
      "aircraft": "Boeing 737",
      "amenities": ["WiFi", "Meal"],
      "baggage": { "carry_on": "1 piece", "checked": "1 piece" }
    }
  ]
}
```

## Testing

```bash
# Run all tests
go test -v ./...

# Run specific package tests
go test -v ./internal/aggregator
go test -v ./internal/transport
```

## Providers

Integrates with:

- **Garuda Indonesia** - Mock delays: 50-100ms
- **Lion Air** - Mock delays: 100-200ms
- **Batik Air** - Mock delays: 200-400ms
- **AirAsia** - Mock delays: 50-150ms, 10% failure rate

Each provider adapter:

1. Loads mock provider responses from `mock-api/` directory
2. Simulates realistic delays and potential failures
3. Returns results within `ProviderResult` domain model

## Design Decisions

1. **Separation of Concerns**: Each package has a single responsibility
2. **Dependency Inversion**: Providers and normalizers follow interface patterns
3. **No Framework Dependencies**: Uses Go standard library HTTP only
4. **Parallel Processing**: Goroutines for concurrent provider queries
5. **Type Safety**: Domain models ensure compile-time type correctness
6. **Go Conventions**: `internal/` package for private implementations, `cmd/` for entry points
7. **Layer Pattern**: Clear separation between domain, application, and transport layers

## Data Flow

```
HTTP Request
    ↓
transport/Handler
    ↓
aggregator/Aggregator.Search()
    ├─→ provider/Provider.SearchFlights() × 4 (parallel)
    │       └─→ normalizer/NormalizeFlight()
    │
    ├─→ Collect & sort results
    └─→ domain/SearchResult
        ↓
    mapper/ToSearchResponse()
        ↓
    HTTP Response (JSON)
```

## Future Enhancements (Skipped for Now)

- [ ] Caching layer (`internal/cache/`)
- [ ] Advanced filtering (`internal/service/filter.go`)
- [ ] Result sorting strategies (`internal/service/sort.go`)
- [ ] Flight ranking algorithm (`internal/service/ranking.go`)
- [ ] GraphQL endpoint
- [ ] Database persistence
- [ ] Authentication & authorization
- [ ] Rate limiting

```bash
# Run as web server
go run ./cmd/api

# Server starts on http://localhost:8080
# API Endpoint: POST /api/search-flights
```

**Example API Request:**

```bash
curl -X POST http://localhost:8080/api/search-flights \
  -H "Content-Type: application/json" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "economy"
  }'
```

## Building

```bash
go build -o flight-search ./cmd/api
./flight-search
```

````

## Testing

```bash
go test ./...
````

## Future Enhancements

- Add caching layer (Redis/memory)
- Implement "best value" scoring algorithm
- Support round-trip searches
- Add rate limiting
- Implement retry logic with exponential backoff
- Add currency formatting
- Support multi-city searches
