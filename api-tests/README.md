# API Testing with Postman

This directory contains Postman collections and environments for comprehensive API testing of the Smart Fit Girl application.

## 📁 Structure

```
api-tests/
├── collections/
│   └── Smart_Fit_Girl_API.postman_collection.json
├── environments/
│   └── Local_Development.postman_environment.json
├── data/
│   └── (test data files - to be added)
└── README.md
```

## 🚀 Quick Start

### 1. Import Collection and Environment

1. **Import Collection**: 
   - Open Postman
   - Click "Import" → "File" 
   - Select `collections/Smart_Fit_Girl_API.postman_collection.json`

2. **Import Environment**:
   - Click "Import" → "File"
   - Select `environments/Local_Development.postman_environment.json`
   - Set "Local Development" as your active environment

### 2. Start Services

Ensure all services are running:
```bash
cd /path/to/smart-fit-girl
docker-compose up -d
# OR
make up
```

Verify services are healthy:
```bash
curl http://localhost:8080/health  # API Service
curl http://localhost:8081/health  # User Service
```

### 3. Run Tests

**Recommended Test Sequence:**
1. **Health Checks** → Verify all services are running
2. **Authentication** → Login and get JWT token
3. **Protected Routes** → Test JWT authentication
4. **Goals & Surveys** → Test survey functionality
5. **Error Handling** → Verify error scenarios

## 📋 Collection Details

### Health Checks
- **API Service Health**: `GET /health`
- **User Service Health**: `GET /health` (direct to user-service)

### Authentication
- **Login - Sophia**: Uses seed user `sophia.woytowitz@gmail.com`
- **Login - Alex Johnson**: Uses seed user `alex.johnson@example.com`
- **Invalid Login**: Tests error handling for wrong credentials

**Auto-Generated Variables:**
- `jwt_token`: Automatically captured from successful login
- `user_id`: Automatically captured from login response

### Protected Routes
- **Get User Profile**: Tests JWT-protected endpoint
- **Access Protected Without Token**: Tests unauthorized access

### User Management (Direct)
- **Get All Users**: Direct call to user-service
- **Get User by ID**: Uses dynamic `{{user_id}}` variable
- **Verify User Credentials**: Tests password verification

### Goals & Surveys
- **Get All Goals**: Returns goals grouped by category for UI
- **Create Survey for User**: Creates fitness survey with goals
- **Get Latest Survey**: Retrieves most recent survey with goals
- **Create Second Survey**: Tests progress tracking with multiple surveys

**Survey Data Structure:**
```json
{
  "currentWeight": 165.5,
  "targetWeight": 150.0,
  "activityLevel": 4,
  "goalIds": [1, 5, 7, 9]
}
```

**Response Structure (Optimized for React UI):**
```json
{
  "id": 1,
  "userId": 1,
  "currentWeight": 165.5,
  "targetWeight": 150.0,
  "activityLevel": 4,
  "createdAt": "2025-08-02T19:17:57.858252Z",
  "goals": {
    "weight": [
      {"id": 1, "name": "Lose", "description": "...", "selected": true},
      {"id": 2, "name": "Maintain", "description": "...", "selected": false}
    ],
    "appearance": [...],
    "strength": [...],
    "endurance": [...]
  }
}
```

### Error Handling
- **Invalid User ID**: Tests 404 responses
- **Invalid Survey Data**: Tests validation errors
- **No Surveys for User**: Tests empty result handling

## 🧪 Test Automation

Each request includes automated tests that verify:

### Status Code Tests
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
```

### Response Structure Tests
```javascript
pm.test("Response has required properties", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('email');
});
```

### Business Logic Tests
```javascript
pm.test("Goals are properly selected", function () {
    var jsonData = pm.response.json();
    var weightGoals = jsonData.goals.weight;
    var selectedCount = weightGoals.filter(g => g.selected).length;
    pm.expect(selectedCount).to.be.greaterThan(0);
});
```

## 🔧 Environment Configuration

### Local Development Environment

| Variable | Value | Description |
|----------|-------|-------------|
| `api_base_url` | `http://localhost:8080` | API Gateway endpoint |
| `user_service_url` | `http://localhost:8081` | User service direct endpoint |
| `meal_service_url` | `http://localhost:8083` | Meal service endpoint (future) |
| `tracking_service_url` | `http://localhost:8084` | Tracking service endpoint (future) |
| `jwt_token` | (auto-populated) | JWT token from login |
| `user_id` | `1` | Default test user ID |

### Custom Environments

For different environments (staging, production), create new environment files:

```json
{
  "name": "Staging",
  "values": [
    {
      "key": "api_base_url",
      "value": "https://api-staging.smartfitgirl.com"
    }
  ]
}
```

## 📊 Running Collection Tests

### Via Postman UI
1. Select collection "Smart Fit Girl API"
2. Click "Run" button
3. Select "Local Development" environment
4. Choose folders/requests to run
5. Click "Run Smart Fit Girl API"

### Via Newman (CLI)
```bash
# Install Newman globally
npm install -g newman

# Run entire collection
newman run collections/Smart_Fit_Girl_API.postman_collection.json \
  -e environments/Local_Development.postman_environment.json

# Run specific folder
newman run collections/Smart_Fit_Girl_API.postman_collection.json \
  -e environments/Local_Development.postman_environment.json \
  --folder "Authentication"

# Generate HTML report
newman run collections/Smart_Fit_Girl_API.postman_collection.json \
  -e environments/Local_Development.postman_environment.json \
  -r htmlextra --reporter-htmlextra-export reports/
```

## 🔍 Test Data

### Seed Users Available
- **Sophia Woytowitz**: `sophia.woytowitz@gmail.com` / `password`
- **Alex Johnson**: `alex.johnson@example.com` / `password`  
- **Emma Thompson**: `emma.thompson@example.co.uk` / `password`

### Test Goals (by Category)
- **Weight**: Lose, Maintain, Gain
- **Appearance**: Lean, Bulk, Definition/Cut
- **Strength**: Gain, Maintain
- **Endurance**: Gain, Maintain

## 🚦 Best Practices

### Test Order
1. Always run **Health Checks** first
2. Run **Authentication** to get JWT token
3. Test **Protected Routes** to verify JWT works
4. Run business logic tests (**Goals & Surveys**)
5. Test **Error Handling** scenarios

### Variable Management
- Use `{{variable_name}}` syntax for dynamic values
- Let authentication tests auto-populate JWT tokens
- Use pre-request scripts for data setup if needed

### Assertions
- Test both success and error scenarios
- Verify response structure matches expected format
- Check business logic (e.g., goal selection)
- Validate data relationships (e.g., survey-goal associations)

## 🔄 CI/CD Integration

To integrate with CI/CD pipelines:

```yaml
# GitHub Actions example
- name: Run API Tests
  run: |
    newman run api-tests/collections/Smart_Fit_Girl_API.postman_collection.json \
      -e api-tests/environments/Local_Development.postman_environment.json \
      --reporters cli,junit \
      --reporter-junit-export test-results.xml
```

## 📈 Monitoring & Reporting

### Test Results
- Green: All tests passed ✅
- Red: Test failures (check response/assertions) ❌
- Yellow: Request failed (check service availability) ⚠️

### Common Issues
- **401 Unauthorized**: Run authentication request first
- **Connection Refused**: Ensure services are running
- **404 Not Found**: Check endpoint URLs and user data

---

## 📝 Adding New Tests

When adding new endpoints:

1. **Add Request**: Create new request in appropriate folder
2. **Add Tests**: Include status code and business logic assertions
3. **Update Variables**: Add any new environment variables needed
4. **Document**: Update this README with new test descriptions

Example test script:
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has expected structure", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('expectedField');
});
```

# API Tests

This directory contains Postman collections and environments for testing the Smart Fit Girl APIs independently of the React frontend.

## API-First Development Approach

We follow an API-first approach where:
1. Design and implement each API endpoint
2. Test thoroughly with Postman
3. Then build the corresponding React UI components

## Directory Structure

```
api-tests/
├── collections/     # Postman collection files (.json)
├── environments/    # Postman environment files (.json)
├── data/           # Test data files (JSON, CSV, etc.)
└── scripts/        # Custom test scripts and utilities
```

## API Endpoints Structure

### User Service (Port 8082)
- `POST /users/register` - User registration
- `POST /users/login` - User authentication
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile
- `GET /users/settings` - Get user settings
- `PUT /users/settings` - Update user settings

### Meal Service (Port 8083)
- `GET /meals/plans` - Get meal plans
- `POST /meals/plans` - Create meal plan
- `PUT /meals/plans/:id` - Update meal plan
- `DELETE /meals/plans/:id` - Delete meal plan

### Tracking Service (Port 8084)
- `POST /tracking/checkin` - Daily check-in (weight, mood, sleep)
- `GET /tracking/checkins` - Get check-in history
- `POST /tracking/evening` - Evening check-in (steps, activities)
- `GET /tracking/stats` - Get tracking statistics

### API Service (Port 8080)
All endpoints are accessible through the main API service which routes to appropriate microservices.

## Testing Workflow

1. Start backend services: `docker-compose up -d`
2. Import Postman collections from `collections/` directory
3. Set up environment variables in Postman
4. Test each endpoint thoroughly
5. Document any issues or changes needed
6. Once API is stable, implement the corresponding UI component

## Environment Variables for Postman

```json
{
  "api_base_url": "http://localhost:8080",
  "user_service_url": "http://localhost:8082",
  "meal_service_url": "http://localhost:8083",
  "tracking_service_url": "http://localhost:8084"
}
```
