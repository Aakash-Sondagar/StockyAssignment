# Stocky - Stock Reward System Backend

## Overview

Stocky is a backend service designed for a fintech platform where users earn **Indian stocks** (e.g., Reliance, TCS) as rewards. The system manages user portfolios, tracks real-time stock valuations, and maintains a **Double-Entry Ledger** to ensure financial accuracy and auditability.

This project is built using **Golang** for high performance and **PostgreSQL** for transactional integrity.

## Tech Stack

- **Language:** Golang (1.20+)
- **Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (HTTP Web Framework)
- **Database:** PostgreSQL (Relational DB)
- **Logging:** [Logrus](https://github.com/sirupsen/logrus) (Structured Logging)
- **Containerization:** Docker (for Database)

## Project Structure

```bash
├── main.go           # Application entry point, routes, and logic
├── schema.sql        # Database schema definitions
├── collection.json   # Postman collection for API testing
├── go.mod            # Go module definition
├── go.sum            # Checksums for dependencies
└── .env              # Environment variables (optional)