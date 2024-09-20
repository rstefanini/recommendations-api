# Project recommendations

Products recommendations api

**This is an incomplete version of the challenge as it lacks many improvements I couldn't add**
    -Swagger Api documentation
    -Integration Tests
    -Higher coverage
    -Dynamic Caching to prevent stress the DB
    -Context handling to fail fast situations
    -Request tracing
    -Thorough Logging and environment dependant
    -Optimizations for db large datasets using concurrency



## Getting Started

This Api exposes the following Endpoints(by default localhost:8080):
    GET /recommendations/users/{userID}     It return the recommended products for the User
                Sample Response:{
                                    "products": [
                                        "2",
                                        "3"
                                    ]
                                }
    POST /collector/interaction             To process products interactions by the user
                Sample body:{
                                "user_id":"user123",
                                "interactions":[
                                    {
                                        "product_sku":"3",
                                        "action":"add_to_cart",
                                        "interaction_timestamp":"2024-09-20T14:00:00Z"
                                    }
                                ]
                            }
    

# Dependencies
In order to get this api working in local you need:
    Docker Daemon running (for mongoDB)


# Steps
    
    make docker-run    
    make run

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```