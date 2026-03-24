# Flight Search API

A high-performance flight search aggregation API built with Go that queries multiple airline providers in parallel, normalizes their responses, and returns unified flight data with advanced filtering and caching capabilities.

## 🚀 Features

- **Multi-Provider Aggregation**: Simultaneously queries multiple airline APIs (Garuda Indonesia, Lion Air, Batik Air, AirAsia)
- **Concurrent Execution**: Uses goroutines and channels to fetch provider data in parallel for faster responses
- **Response Normalization**: Converts provider-specific formats into a unified domain model
- **Advanced Filtering**: Supports filtering by price range, departure/arrival time, airlines, duration, stops, seats, and passenger count
- **Sorting**: Enables sorting by price, duration, departure time, and arrival time
- **Smart Ranking**: Prioritizes optimal results based on price and convenience
- **Rate Limiting**: Prevents excessive requests to external providers
- **Intelligent Caching**: In-memory TTL caching to reduce API calls and improve performance
- **Extensible Architecture**: Hexagonal design for easy addition of new providers
- **Fault Tolerance**: Gracefully handles partial failures from providers
- **Docker Support**: Containerized deployment using docker-compose

## 🏗️ Architecture

### Data Flow

```
Client Request
      ↓
Cache Check
      ↓
Fetch Providers (Concurrent)
      ↓
Normalize Responses
      ↓
Aggregate Results
      ↓
Filter → Sort → Rank
      ↓
Cache Result
      ↓
Return Response
```

### Design Decisions

**Separation of Concerns**

- Providers, aggregation, filtering, sorting, and caching are handled independently

**Concurrent Providers**

- External APIs are queried in parallel using goroutines and channels

**Pluggable Providers**

- New providers can be added via a common interface

**Caching Strategy**

- In-memory cache with TTL
- Cache keys derived from request parameters

**Filtering & Sorting**

- Multi-dimensional filtering across price, time range, airline, duration, stops, and seat availability
- Flexible sorting by price, duration, and time (ascending/descending) across aggregated results

**Resilience**

- Partial failures are tolerated (graceful degradation)
- Rate limiting prevents provider overuse

### Providers

Integrates with:

- **Garuda Indonesia** - Mock delays: 50-100ms
- **Lion Air** - Mock delays: 100-200ms
- **Batik Air** - Mock delays: 200-400ms
- **AirAsia** - Mock delays: 50-150ms, 10% failure rate

Each provider adapter:

1. Loads mock provider responses from `internal/mock-api/` directory
2. Simulates realistic delays and potential failures
3. Returns results within `ProviderResult` domain model

## 📁 Project Structure

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
│   │   ├── aggregate.go               # Flight aggregation orchestrator
│   │   ├── filter.go                  # Flight filtering logic
│   │   ├── sort.go                    # Flight sorting logic
│   │   └── cache.go                   # Caching layer
│   ├── mapper/
│   │   └── response.go                # Domain → DTO conversion
│   ├── transport/
│   │   ├── http.go                    # HTTP handler
│   │   └── router.go                  # Route definitions
│   ├── mock-api/                      # Mock provider responses
│   ├── constant/                      # Application constants
│   └── utils/                         # Utility functions
├── Dockerfile                         # Docker build configuration
├── docker-compose.yml                 # Docker Compose setup
├── go.mod
└── README.md
```

### Package Responsibilities

- **domain**: Pure business domain models (no external dependencies)
- **provider**: Interfaces and implementations for external airline APIs
- **normalizer**: Converts provider-specific formats to unified domain models
- **aggregator**: Orchestrates parallel provider queries, aggregation, filtering, sorting, and caching
- **mapper**: Converts domain models to HTTP response DTOs
- **transport**: HTTP handler and request/response management
- **cmd/api**: Application bootstrap and dependency injection
- **constant**: Application-wide constants
- **utils**: Utility functions (formatter, etc.)

## 📡 API Endpoints

### POST /api/search-flights

Search for flights across all providers.

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

**Request with Additional Filters:**

```json
{
  "origin": "CGK",
  "destination": "DPS",
  "departureDate": "2025-12-15",
  "passengers": 1,
  "cabinClass": "economy",
  "priceRange": [1000000, 2000000],
  "numberOfStops": [0, 1, 2],
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
      "id": "JT742_Lion Air",
      "provider": "Lion Air",
      "airline": {
        "name": "Lion Air",
        "code": "JT"
      },
      "flight_number": "JT742",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T11:45:00+07:00",
        "timestamp": 1765773900
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T14:35:00+08:00",
        "timestamp": 1765780500
      },
      "duration": {
        "total_minutes": 110,
        "formatted": "1h 50m"
      },
      "stops": 0,
      "price": {
        "amount": 890000,
        "currency": "IDR",
        "formatted": "Rp 890.000"
      },
      "available_seats": 38,
      "cabin_class": "economy",
      "aircraft": "Boeing 737-800",
      "amenities": ["wiFi", "meal"],
      "baggage": {
        "carry_on": "7 kg",
        "checked": "20 kg"
      }
    }
  ]
}
```

### GET /health

Health check endpoint.

**Response:**

```json
{
  "status": "ok"
}
```

## ⚙️ Getting Started

### Prerequisites

- Go 1.25.0 or later
- Docker and Docker Compose (for containerized deployment)

### Method 1: Run with Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/irrangga/flight-search.git
cd flight-search

# Build and run with Docker Compose
docker-compose up --build

# Check if the API available at http://localhost:8080/health
```

### Method 2: Run Locally

```bash
# Clone the repository
git clone https://github.com/irrangga/flight-search.git
cd flight-search

# Install dependencies
go mod download

# Run the application
go run ./cmd/api

# Check if the API available at http://localhost:8080/health
```

### Method 3: Build and Run Binary

```bash
# Build the binary
go build -o flight-search ./cmd/api

# Run the binary
./flight-search
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/aggregator -v
```

## 📄 License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0).
