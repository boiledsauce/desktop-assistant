# Go Personal Automation Service

This is a Go application that utilizes clean architecture principles to provide automation services for personal needs. The application is designed to be extensible, with the ability to add more services and endpoints as needed.

The whole purpose and goal is to reduce the amount of repetetive tasks you have to do - which may be very individual.

## Features

- **Server-Sent Events (SSE)**: The application provides a `/events` endpoint that clients can listen to for real-time updates. This utilizes the Server-Sent Events (SSE) protocol.

- **File Download Completion Notification**: When a file download is complete, the application sends an event to the client who initiated connection with the `/events` endpoint with the path of the downloaded file.

# Getting Started

## Recommended path (Kubernetes)

1. Make sure you have kubernetes installed
2. Modify all the k8s files inside k8s/ folder according to your needs (e.g set the proper mounted volume, as the filesystem will be watched by a file.
3. Run the jenkins job in any way you desire (It is currently a local filesystem pipeline)

Port forward the kubernetes pod by

```
kubectl get pods
kubectl port-forward <pod name> 8081:8080
The pod should now be ready for outwards connections
```

## Go native path

To run the application, you will need to have Go installed on your machine. Once you have Go installed, you can run the application with the following command:

```go
docker run
```

This will start the application and it will begin listening for events on the hardcoded URL.

## Endpoints

- `/events`: This endpoint streams events to the client in real-time through Server-Sent Events (SSE). Clients can listen to this endpoint to receive continuous updates about file downloads.

## Pipeline

The pipeline:

- Creates a kubrenetes cluster if it does not exist
- Prepares ENV with variables
- Builds a Docker Image of the project and tags it accordingly
- Updates the Kubernetes deployment by replacing the docker image and writing its new build-tag to the deployment.yaml and applies all Kubernetes files.

## Future Plans

We plan to add more endpoints and services to this application to automate more personal needs.
Desirable features are:

- Replace duplicates of files at a path
- Sort files at a certain folder path
- Automate mundane sequential steps of manual labour. (Lets be lazy together)

## Contributing

Contributions are welcome! Please feel free to submit a pull request.
