# Go Personal Automation Service

This is a Go application that utilizes clean architecture principles to provide automation services for personal needs. The application is designed to be extensible, with the ability to add more services and endpoints as needed.

The whole purpose and goal is to reduce the amount of repetetive tasks you have to do - which may be very individual.

## Features

- **Server-Sent Events (SSE)**: The application provides a `/events` endpoint that clients can listen to for real-time updates. This utilizes the Server-Sent Events (SSE) protocol.

- **File Download Completion Notification**: When a file download is complete, the application sends an event to the client who initiated connection with the `/events` endpoint with the path of the downloaded file.

## Getting Started

To run the application, you will need to have Go installed on your machine. Once you have Go installed, you can run the application with the following command:

```go
docker run
```

This will start the application and it will begin listening for events on the hardcoded URL.

## Endpoints

- `/events`: This endpoint streams events to the client in real-time through Server-Sent Events (SSE). Clients can listen to this endpoint to receive continuous updates about file downloads.

## Future Plans

We plan to add more endpoints and services to this application to automate more personal needs.

# Roadmap UseCases

- Replace duplicates of files at a path
- Sort at a path
- Bundle multiple usecases in one

## Contributing

Contributions are welcome! Please feel free to submit a pull request.
