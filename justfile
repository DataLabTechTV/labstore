set shell := ["bash", "-uc"]
set dotenv-load

backend-serve:
    cd backend/ && go run main.go serve
