# FeedbackBot — Backend Dockerfile (placeholder)
# To be completed by Sardor during development

FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o feedbackbot-api .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/feedbackbot-api .
EXPOSE 8080
CMD ["./feedbackbot-api"]
