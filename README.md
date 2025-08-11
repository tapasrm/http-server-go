# Go HTTP Server

This is a simple, barebones HTTP/1.1 server built with Go. It's capable of handling multiple clients simultaneously and serves as a basic example of building web servers in Go.

## Features

- **Concurrent Connections:** Handles multiple clients at once using Goroutines.
- **Basic Routing:** Includes a simple router to handle different URL paths.
- **Static & Dynamic Responses:** Serves both static content and dynamically generated responses based on the request.

## Getting Started

### Prerequisites

- Go (version 1.18 or later)

### Running the Server

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-http-server.git
    cd go-http-server
    ```

2.  **Run the server:**
    ```bash
    ./run.sh
    ```

The server will start on `localhost:4221`.

## Endpoints

The server exposes the following endpoints:

- `GET /`: Returns a simple "Hello, World!" message.
- `GET /echo/{message}`: Echoes back the `{message}` provided in the URL path.

## Project Structure

```
.
├── app/
│   ├── main.go       # Main application entry point
│   ├── request.go    # Request parsing logic
│   ├── response.go   # Response building logic
│   └── router.go     # Routing logic
├── go.mod
├── go.sum
├── run.sh # Script to run the server
└── README.md
```

## Technologies Used

- **Go:** The server is built entirely using the standard Go library, with no external frameworks.
