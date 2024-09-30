# Examify
An online IELTS exam platform powered by microservices communicating via gRPC.

## About
This repository contains the codebase for taking IELTS exams online. The system is built using a **microservices architecture**, with services communicating through **gRPC** for efficient and scalable interactions between the components.

## Project Structure

```
~/GolandProjects/examify git:[master]
tree
.
├── api-gateway
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── Dockerfile
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── grpc_clients
│   │   │   ├── auth_client.go
│   │   │   ├── ielts_client.go
│   │   │   └── user_client.go
│   │   ├── handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── ielts_book_handler.go
│   │   │   ├── ielts_exam_handler.go
│   │   │   ├── init_clients.go
│   │   │   └── user_handler.go
│   │   ├── middleware
│   │   │   └── auth_middleware.go
│   │   ├── models
│   │   │   └── ielts.go
│   │   └── routes
│   │       ├── ielts_routes.go
│   │       ├── routes.go
│   │       └── user_routes.go
│   ├── proto
│   │   ├── auth.proto
│   │   ├── common.proto
│   │   ├── ielts.proto
│   │   ├── pb
│   │   │   ├── auth_grpc.pb.go
│   │   │   ├── auth.pb.go
│   │   │   ├── common.pb.go
│   │   │   ├── ielts_grpc.pb.go
│   │   │   ├── ielts.pb.go
│   │   │   ├── user_grpc.pb.go
│   │   │   └── user.pb.go
│   │   └── user.proto
│   └── utils
│       └── utils.go
├── auth-service
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── grpc_clients
│   │   │   └── UserClient.go
│   │   ├── repository
│   │   │   └── repository.go
│   │   ├── security
│   │   │   └── jwt.go
│   │   ├── server
│   │   │   └── server.go
│   │   ├── service
│   │   │   └── service.go
│   │   └── telegram
│   │       ├── bot.go
│   │       ├── handlers.go
│   │       └── utils.go
│   └── proto
│       ├── auth.proto
│       ├── common.proto
│       ├── pb
│       │   ├── auth_grpc.pb.go
│       │   ├── auth.pb.go
│       │   ├── common.pb.go
│       │   ├── user_grpc.pb.go
│       │   └── user.pb.go
│       └── user.proto
├── ielts-service
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── grpc_clients
│   │   │   ├── integration_client.go
│   │   │   └── user_client.go
│   │   ├── models
│   │   │   └── models.go
│   │   ├── repository
│   │   │   ├── db.go
│   │   │   ├── helper.go
│   │   │   └── repository.go
│   │   ├── server
│   │   │   └── server.go
│   │   ├── service
│   │   │   ├── ielts_service.go
│   │   │   └── ielts_service_test.go
│   │   └── utils
│   │       └── utils.go
│   ├── migrations
│   │   ├── ielts_service_down.sql
│   │   └── ielts_service_up.sql
│   └── proto
│       ├── common.proto
│       ├── ielts.proto
│       ├── integration.proto
│       ├── pb
│       │   ├── common.pb.go
│       │   ├── ielts_grpc.pb.go
│       │   ├── ielts.pb.go
│       │   ├── integration_grpc.pb.go
│       │   ├── integration.pb.go
│       │   ├── user_grpc.pb.go
│       │   └── user.pb.go
│       └── user.proto
├── integration-service
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── server
│   │   │   └── server.go
│   │   ├── service
│   │   │   ├── service.go
│   │   │   └── writing_service.go
│   │   └── utils
│   │       └── utils.go
│   ├── main.go
│   └── proto
│       ├── integration.proto
│       └── pb
│           ├── integration_grpc.pb.go
│           └── integration.pb.go
├── Makefile
├── scripts
│   ├── drop_databases.sh
│   └── setup_databases.sh
└── user-service
    ├── cmd
    │   └── main.go
    ├── config
    │   ├── config.go
    │   └── config.yaml
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── internal
    │   ├── repository
    │   │   ├── db.go
    │   │   └── user_repository.go
    │   ├── server
    │   │   └── server.go
    │   └── service
    │       └── user_service.go
    ├── migrations
    │   ├── user_service_down.sql
    │   └── user_service_up.sql
    └── proto
        ├── common.proto
        ├── pb
        │   ├── common.pb.go
        │   ├── user_grpc.pb.go
        │   └── user.pb.go
        └── user.proto

59 directories, 119 files
```

## Services
The project is composed of the following microservices:

- **API Gateway**: Routes requests to the appropriate microservices.
- **Authentication Service**: Manages user authentication, login, and session handling.
- **IELTS Service**: Core service that facilitates the exam process, including listening, reading, writing, and speaking modules.
- **Integration Service**: Handles third-party integrations, including AI services for automated scoring.
- **User Service**: Manages user data, profile information, and exam history.

## Technology Stack
- **Primary Language**: Go (98.9%)
- **Communication Protocol**: gRPC for efficient microservice-to-microservice communication.
