# Online Food Ordering API Server

API Server serves as backend to Online food ordering app.

## Endpoints

### 1. Get List of Products
- **GET** `/products`
- Returns all available products

### 2. Find Product by ID
- **GET** `/products/{id}`
- Returns a specific product by its ID

### 3. Place an Order
- **POST** `/orders`
- Places a new order
- Coupon code: Optional

## Authentication

Authentication is optional. Use api_key header to pass auth-key value.

## Data Sources

- Product catalogue is static JSON file available on the server: [products.json](./data/products.json)
- Coupon data is available in S3 bucket, which we download during server bootup and perform validation

### Configuration

Configure coupon data source by setting these environment variables in App Runner (requires restart):

- `COUPON_REGION` - AWS region (default: `ap-southeast-2`)
- `COUPON_BUCKET` - S3 bucket name (default: `orderfoodonline-files`)
- `COUPON_FILES` - Comma-separated filenames (default: `couponbase1,couponbase2,couponbase3`)

## Verification
For testing, we already deployed in AWS App Runner. Endpoints are at:
https://d5ppvd5edn.ap-southeast-2.awsapprunner.com/api/

Swagger doc link:
https://d5ppvd5edn.ap-southeast-2.awsapprunner.com/apidoc/index.html

## Improvements which can be done:

### Coupon data to redis-cache
As per requirement, we have coupon data available as S3 blobs. Due to this constraint, max optimized respone time we could able to achieve while validatin coupons is 2-4 seconds. Which is not at all acceptable time for production grade.

Without this constraint, we can load this data in to redis-cache and let the server consume from there, which would reduce respone time considerable to milli-seconds.

### Making coupon/product config hot-reload
Current implementation, when we need to change coupon data, we need to change app-runner environment values and restart the service.

Instead, we can have hot-reload, which could be a go-routine doing downloads from s3 in timely manner when there is a change in cloudconfig var.

Same implementation can be done for product.json, sourcing from different location and doing a hot reload.

OR

We can move to DB.

