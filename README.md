# Gnat - Fast, lightweight load generator 🦟

Gnat is a lightweight, blazing-fast CLI tool for sending concurrent HTTP requests to test your APIs under load.

## 🔧 Installation

```shell
git clone https://github.com/<your-username>/gnat.git
cd gnat
go build -o gnat
```

## 🚀 Usage

```shell
./gnat run \
  --url https://your-api.com/ping \
  --duration 10s \
  --concurrency 12
```

### 🏁 Example Output

🚀 Sending requests to https://your-api.com/ping for 10s with 12 workers
✅ Load complete. Sent 601 requests in 10s

## License

Licensed under the [Apache License 2.0](LICENSE).