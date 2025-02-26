# Build the frontend 
FROM node:18-alpine AS build-frontend
WORKDIR /app
COPY Frontend/package.json Frontend/package-lock.json* ./
RUN npm ci
COPY Frontend/ ./ 
RUN npm run build

# Build stage for Go server
FROM golang:1.24-alpine AS build-backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine:latest AS final
WORKDIR /app
COPY --from=build-frontend /app/dist ./dist
COPY --from=build-backend /app/server ./server

# Expose the port the server runs on
EXPOSE 8080
# Command to run the server
CMD ["./server"]


