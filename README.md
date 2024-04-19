# Balance API

## Endpoints

### 1. GET /balance

Fetches the balance of the provided Ethereum address.

#### Query Parameters:

- `address` (required): The Ethereum address for which to retrieve the balance.
    - Example: `0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5`
- `blockTag` (optional): The state block number. If not provided, defaults to `latest`.

#### Request Example:

```bash
curl --request GET \
  --url 'http://localhost:8000/balance?address=0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5&blockTag=latest'
```

### 2. GET /ht

Performs a health check on the service.

#### Response:

- Returns a 200 status code if at least one provider is alive.
- Returns a 500 status code if no providers are alive.

#### Request Example:

```bash
curl --request GET \
  --url 'http://localhost:8000/ht'
```

### 3. GET /metrics

Retrieves Prometheus metrics for the service.

#### Request Example: 

```bash
curl --request GET \
  --url 'http://localhost:8000/metrics'
```
