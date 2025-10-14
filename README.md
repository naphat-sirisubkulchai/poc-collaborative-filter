# POC Collaborative Filter

A production-ready recommendation engine POC for wealth management platforms. This service implements intelligent **multi-strategy recommendation algorithms** including Collaborative Filtering, Content-Based, Hybrid, and Popularity-based approaches. The system analyzes customer similarity across personal information, behavioral patterns, and transaction history to deliver personalized investment recommendations.

## Architecture

This service implements **Clean Architecture** principles with a modular, scalable design:

### Core Technologies
- **GraphQL API** - Type-safe API with gqlgen code generation
- **PostgreSQL** - Relational database with JSONB support for flexible data modeling
- **GORM** - Modern ORM with migrations and type-safe queries
- **Uber FX** - Dependency injection for clean separation of concerns
- **Gin** - High-performance HTTP router and middleware
- **Docker** - Containerized deployment with hot-reload development environment

### Key Features
- **Multi-Strategy Recommendation Engine** with automatic strategy selection
- **Real-time similarity calculation** across 4 weighted dimensions
- **Flexible GraphQL API** for complex queries and mutations
- **Type-safe domain models** with PostgreSQL-specific types (JSONB, text arrays)
- **Extensible architecture** ready for caching, background jobs, and microservices integration

## Project Structure

```
poc-collaborative-filter/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── app/
│   │   └── server.go                  # HTTP server setup (Gin + GraphQL)
│   ├── di/
│   │   └── container.go               # Uber FX dependency injection
│   ├── domain/                        # Domain models (entities)
│   │   ├── customer.go                # Customer with JSONB traits
│   │   ├── recommendation.go          # Recommendation entity
│   │   ├── recommendation_config.go   # Strategy configuration
│   │   ├── strategy.go                # Recommendation strategies
│   │   ├── constants.go               # Domain constants
│   │   ├── lead.go, team.go, task.go # Additional entities
│   │   ├── compliance.go, campaign.go # Compliance & campaigns
│   ├── repository/                    # Data access layer
│   │   ├── interfaces.go              # Repository interfaces
│   │   ├── customer_repository.go     # Customer data access
│   │   └── recommendation_repository.go # Recommendation storage
│   ├── service/                       # Business logic layer
│   │   ├── collaborative_filter_service.go # CF algorithm
│   │   ├── recommendation_service.go       # Strategy orchestration
│   │   ├── customer_service.go             # Customer operations
│   │   └── diversification.go              # Diversification logic
│   ├── utils/
│   │   └── similarity.go              # Similarity metric calculations
│   ├── config/
│   │   └── config.go                  # Environment configuration
│   ├── database/
│   │   └── database.go                # GORM connection setup
│   └── middleware/                    # HTTP middleware
│       ├── cors.go                    # CORS configuration
│       └── logger.go                  # Request logging
├── graph/                             # GraphQL layer
│   ├── schema.graphqls                # GraphQL schema (edit this)
│   ├── schema.resolvers.go            # Resolver implementations
│   ├── resolver.go                    # Resolver setup
│   ├── model/                         # GraphQL models
│   │   └── scalars.go                 # Custom scalars (Time, JSON, Decimal)
│   └── generated/                     # Auto-generated code (don't edit)
├── migrations/                        # Database migrations
├── docker-compose.yml                 # Production Docker setup
├── docker-compose.dev.yml             # Development with hot-reload
├── Dockerfile                         # Container build instructions
├── Makefile                           # Build and dev commands
├── gqlgen.yml                         # GraphQL codegen config
├── mockdata.json                      # Seed data for testing
└── seed-data.sh                       # Database seeding script
```

## Multi-Strategy Recommendation Engine

### Strategy Selection Algorithm

The system automatically selects the best recommendation strategy based on customer profile and activity:

| Strategy | Trigger Conditions | Use Case |
|----------|-------------------|----------|
| **Collaborative Filtering** | Activity Score ≥ 20 AND Transactions ≥ 3 | Active customers with sufficient behavior data |
| **Hybrid** | Activity Score ≥ 5 AND Profile Complete ≥ 10 | Moderate activity with good profile data |
| **Content-Based** | Profile Complete ≥ 5 | New customers with stated preferences |
| **Popularity** | Cold Start (minimal data) | Brand new customers or prospects |

**Activity Score Calculation:**
```
ActivityScore = (Searches × 0.1) + (Views × 0.1) + (Factsheet × 0.2) +
                (Buy × 1.0) + (Sell × 0.8) + (Switch × 0.6)
```

### Collaborative Filtering Algorithm

The system calculates customer similarity using **5 different similarity metrics** across four weighted dimensions:

**Dimension Weights:**
- Personal Information: 25%
- Survey Data: 35% (highest weight - explicit preferences)
- Activity Tracking: 20%
- Transaction History: 20%

#### 1. Personal Similarity (25%)

| Metric Type | Property | Formula | Example |
|------------|----------|---------|---------|
| **Distance** | Age | `max(0, 1 - \|age1-age2\|/50)` | 7 years apart → 0.86 |
| **Ratio** | Annual Income | `min/max` | 200k vs 250k → 0.80 |
| **Distance** | Experience | `max(0, 1 - \|exp1-exp2\|/20)` | 5 years apart → 0.75 |
| **Binary** | Occupation | `1.0 if same else 0.0` | Match → 1.0 |
| **Binary** | Risk Profile | `1.0 if same else 0.0` | Match → 1.0 |

#### 2. Survey Similarity (35%) - Highest Weight

| Metric Type | Property | Formula | Why This Metric? |
|------------|----------|---------|------------------|
| **Jaccard** | Investment Types | `\|A∩B\| / \|A∪B\|` | Conservative - good for similar-sized sets |
| **Jaccard** | Industry Sectors | `\|A∩B\| / \|A∪B\|` | Precise matching for sector interests |
| **Dice** | Regional Preferences | `2×\|A∩B\| / (\|A\|+\|B\|)` | Forgiving - handles size differences |
| **Binary** | ESG Interest | `1.0 if same else 0.0` | Boolean preference |
| **Binary** | Crypto Interest | `1.0 if same else 0.0` | Boolean preference |

**Why Dice for Regions?** Some customers select 1-2 regions, others select 5+. Dice Coefficient is more forgiving for size differences than Jaccard, preventing harsh penalties for new customers.

#### 3. Activity Similarity (20%) - With Cold Start Handling

| Metric Type | Property | Formula | Cold Start Handling |
|------------|----------|---------|---------------------|
| **Ratio** | Total Searches | `min/max` | Both 0 → 0.5, One 0 → 0.2 |
| **Ratio** | Total Views | `min/max` | Both 0 → 0.5, One 0 → 0.2 |
| **Ratio** | Factsheet Reads | `min/max` | Both 0 → 0.5, One 0 → 0.2 |
| **Ratio** | Share Count | `min/max` | Both 0 → 0.5, One 0 → 0.2 |

**Cold Start Strategy:** Instead of returning 0.0 (which kills similarity), the system returns:
- **0.5** (neutral) if both customers have no activity
- **0.2** (low but not zero) if only one customer has activity

#### 4. Transaction Similarity (20%)

| Metric Type | Property | Formula | Note |
|------------|----------|---------|------|
| **Ratio** | Buy Count | `min/max` | Requires both > 0 |
| **Ratio** | Sell Count | `min/max` | Requires both > 0 |
| **Ratio** | Switch Count | `min/max` | Requires both > 0 |
| **Ratio** | Fund Detail Access | `min/max` | Requires both > 0 |

**Final Similarity Score:**
```
TotalSimilarity = (PersonalSim × 0.25) +
                  (SurveySim × 0.35) +
                  (ActivitySim × 0.20) +
                  (TransactionSim × 0.20)
```

### Similarity Metrics Explained

#### Jaccard Similarity
```
Jaccard(A, B) = |A ∩ B| / |A ∪ B|
```
**Example:**
```
A = ["Equity", "Mutual Funds", "ETF"]
B = ["Equity", "Mutual Funds", "Bonds"]
Intersection = 2, Union = 4
Score = 2/4 = 0.5
```

#### Dice Coefficient (Sørensen-Dice)
```
Dice(A, B) = 2 × |A ∩ B| / (|A| + |B|)
```
**Example:**
```
A = ["US", "Asia"]
B = ["US", "Europe", "Asia", "Global"]
Intersection = 2, |A| = 2, |B| = 4
Score = 2×2 / (2+4) = 0.67

Jaccard would give: 2/4 = 0.5 (more harsh)
Dice gives: 0.67 (more forgiving)
```

#### Distance-Based Similarity
```
Similarity = max(0, 1 - (|value1 - value2| / normalizationFactor))
```
**Example (Age):**
```
Age1 = 43, Age2 = 50, Norm = 50
Difference = 7
Score = 1 - (7/50) = 0.86
```

#### Ratio Similarity
```
Similarity = min(value1, value2) / max(value1, value2)
```
**Example (Income):**
```
Income1 = 185,000, Income2 = 250,000
Score = 185,000 / 250,000 = 0.74
```

### Content-Based Filtering

For customers with low activity (Activity Score < 20), uses their **stated preferences** directly:
- Investment Types from survey
- Industry Sectors from profile
- Risk Profile alignment

No similarity calculation needed - directly recommends based on explicit preferences.

### Popularity-Based Filtering

For cold start customers (new prospects with minimal data):
- Returns predefined popular items
- Descending score: 0.6, 0.5, 0.4, 0.3, 0.2
- No personalization - same for all cold start users

## Domain Models

### Customer
Uses **PostgreSQL-specific types**:
- `pq.StringArray` for text[] columns (AccountType, Holdings)
- `datatypes.JSON` for JSONB columns (Traits)
- `decimal.Decimal` for financial amounts
- UUID primary keys as strings

### Traits Structure (JSONB)
```json
{
  "personality": ["Risk-Averse", "Conservative"],
  "preferences": ["Long-term Growth", "Dividend Income"],
  "notes": "Prefers ESG investments",
  "tags": ["High-Value", "Active"],
  "properties": [
    {"key": "Age", "value": "45", "type": "number"},
    {"key": "Occupation", "value": "Engineer", "type": "text"},
    {"key": "Total Searches", "value": "145", "type": "number"}
  ],
  "relatedPersons": [...]
}
```

## Getting Started

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 16+ (or use Docker)
- Redis (optional, for caching)

### Development Setup

1. **Install dependencies:**
```bash
make deps
```

2. **Start development environment:**
```bash
make dev
```

This starts:
- App container with hot reload (Air)
- PostgreSQL database
- Redis cache

3. **Generate GraphQL code** (after schema changes):
```bash
make generate
```

4. **Run migrations:**
```bash
make cli-migrate
```

5. **Seed data:**
```bash
make cli-seed
```

### Environment Variables

Copy `.env.example` to `.env` and configure:

```env
# Database
DB_DRIVER=postgres
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=collaborative_filter_dev

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=debug

# Collaborative Filtering
SIMILARITY_THRESHOLD=0.5
MAX_RECOMMENDATIONS=10
CACHE_TTL=3600
```

## GraphQL API

### Queries

**Get Customer:**
```graphql
query {
  customer(id: "uuid-here") {
    id
    name
    segment
    riskProfile
    traits {
      properties {
        key
        value
        type
      }
    }
  }
}
```

**Find Similar Customers:**
```graphql
query {
  similarCustomers(customerId: "uuid-here", limit: 5) {
    customer {
      id
      name
      segment
    }
    similarityScore
    weights {
      personal
      survey
      activityAccess
      transaction
    }
  }
}
```

**Generate Recommendations:**
```graphql
query {
  generateRecommendations(
    customerId: "uuid-here"
    type: FUND
  ) {
    id
    itemName
    score
    reason
    metadata {
      algorithm
      confidence
      basedOnUsers
      similarityAvg
    }
  }
}
```

### Mutations

**Create Customer:**
```graphql
mutation {
  createCustomer(input: {
    name: "John Doe"
    email: "john@example.com"
    segment: HNI
    status: ACTIVE
    riskProfile: MODERATE
  }) {
    id
    name
  }
}
```

**Accept Recommendation:**
```graphql
mutation {
  acceptRecommendation(id: "rec-uuid") {
    id
    status
  }
}
```

## Development Commands

```bash
make help               # Show all available commands
make deps               # Install Go dependencies
make build              # Build binary
make run                # Run locally
make dev                # Start with hot reload (Docker)
make dev-stop           # Stop development containers
make dev-logs           # View logs
make generate           # Generate GraphQL code
make test               # Run tests
make clean              # Clean build artifacts
```

## Production Deployment

```bash
make docker-build       # Build production image
make docker-up          # Start production containers
make docker-down        # Stop production containers
make docker-logs        # View production logs
```

## Key Technologies

- **Go 1.25** - Programming language
- **gqlgen** - GraphQL server code generation
- **GORM** - ORM with PostgreSQL support
- **Gin** - HTTP web framework
- **Uber FX** - Dependency injection
- **Zap** - Structured logging
- **lib/pq** - PostgreSQL array support
- **shopspring/decimal** - Decimal precision for finances
- **Air** - Hot reload for development

## PostgreSQL-Specific Features

### String Arrays (pq.StringArray)
Used for: AccountType, Holdings, Specialization, CustomerType
```go
type Customer struct {
    Holdings pq.StringArray `gorm:"type:text[]"`
}
```

### JSONB (datatypes.JSON)
Used for: Traits, Metadata, MeetingLink, Location
```go
type Customer struct {
    Traits datatypes.JSON `gorm:"type:jsonb"`
}
```

## Next Steps for CRM Integration

1. **Add Authentication & Authorization**
   - JWT tokens
   - Role-based access control

2. **Implement Caching Layer**
   - Redis for similarity scores
   - Cache invalidation strategy

3. **Add More Recommendation Types**
   - Product recommendations
   - Sector recommendations
   - Customer-to-customer matching

4. **Performance Optimization**
   - Batch similarity calculations
   - Pre-compute similarity matrices
   - Background job processing

5. **Integrate with CRM Service**
   - Event-driven updates
   - Shared database or API calls
   - Message queue (NATS/RabbitMQ)

## License

Proprietary - Vahalla Backend
