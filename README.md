# wbtask

## Running the Application

To run the application and its dependencies:

```bash
docker-compose up --build
```

## Running the integration tests

To run the integration tests:

```bash
docker-compose -f docker-compose-test.yaml up --abort-on-container-exit
```

## API Endpoints
### Save User
To create a new user:
```bash
curl -X POST http://localhost:8080/save \
-H "Content-Type: application/json" \
-d '{
    "id": "<uuid>",
    "name": "John Doe",
    "email": "johndoe@example.com",
    "date_of_birth": "1990-01-01T00:00:00Z"
}'
```

### Get User
To get a user by id:
```bash
curl http://localhost:8080/<uuid>
```
