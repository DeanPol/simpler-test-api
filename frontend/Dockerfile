# Use the official Node.js image from the Docker Hub
FROM node:18 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json
COPY frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY frontend/ ./

# Build the React app
RUN npm run build

# Use a simple web server to serve the static files
FROM nginx:alpine

# Copy the React build from the previous stage
COPY --from=build /app/build /usr/share/nginx/html

# Expose port 80 to the outside world
EXPOSE 80

# Command to run nginx
CMD ["nginx", "-g", "daemon off;"]
