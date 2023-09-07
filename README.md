# Check MIA Dates

This project is designed to fetch and update driving exam dates and times from the Georgia government API and store them in MongoDB.

It also provides endpoints to retrieve the saved data.

## Features

1. Pulls driving exam dates and times for various categories (theory, manual, and automatic) and cities in Georgia;
2. Uses randomized User-Agents when sending HTTP requests;
3. Periodically triggers data update every 10 minutes;
4. Provides a frontend static file server;
5. Logs events and errors using the `logrus` library.

## Endpoints

- `/api/update-dates` - Fetches data from the API and updates the MongoDB collections.
- `/api/get-theory` - Retrieves theory exam dates and times.
- `/api/get-manual` - Retrieves manual driving exam dates and times.
- `/api/get-auto` - Retrieves automatic driving exam dates and times.
- `/api/last-exec-time` - Fetches the last execution timestamp.
- `/` - Serves frontend static files.

## Prerequisites

1. Go language environment;
2. MongoDB instance;
3. Set the `MONGO_URI` environment variable to your MongoDB connection string.

## How to Run

1. Clone the repository;
2. Navigate to the project directory;
3. Run the command `go run main.go`.

> The application will start, and the logs will indicate that the application is listening on http://localhost:8080.

## Libraries Used

1. `logrus`: For logging events and errors;
2. `mongo-driver`: For connecting to and performing operations on MongoDB;
3. Standard Go libraries for HTTP server, JSON manipulation, and other utilities.

## Important Note

This project uses a hardcoded API endpoint.

Ensure you have the appropriate permissions or authorization to access the API before triggering the application. Also, avoid excessive requests to prevent potential rate-limiting or blacklisting.