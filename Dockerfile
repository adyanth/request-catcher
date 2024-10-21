FROM node:23 AS frontend
WORKDIR /app
COPY frontend/ .
RUN npm run build

FROM golang:1.23 AS backend
WORKDIR /app
COPY go.* .
RUN go mod download
COPY catcher catcher
COPY main.go .
RUN CGO_ENABLED=0 go build -o /requestcatcher main.go

FROM gcr.io/distroless/base-debian12 AS requestcatcher
WORKDIR /frontend
COPY --from=frontend  /app/dist/ .
COPY frontend/favicon.ico .
WORKDIR /
COPY --from=backend /requestcatcher .
ENV FRONTEND_DIR=/frontend/
ENV FAVICON=/frontend/favicon.ico
ENTRYPOINT ["/requestcatcher"]
