# Chirpy
**Chirpy** - a social media platform to share quick thoughts and quips

## Requirements
[Postgresql](https://www.postgresql.org/) and [Go](https://go.dev/)

## Installation
Install by executing go install github.com/mike-moseley/chirpy@latest

## Setup
Have a running postgresql server, and setup a .env file containing:
- DB_URL: URL to your postgresql server
- PLATFORM: Whether or not you are working on a development instance, only forbids resetting the db via the /admin/reset endpoint
- TOKEN_STRING: Secret token for generating and validating authentication tokens, a uses HMAC-SHA256 hash
- POLKA_KEY: An API key for a fake payment processor to emulate receiving webhooks

## Running the Server
- Run by entering `./chirpy` into your shell

## API
- API endpoints documented in [docs/API.md](docs/API.md)
