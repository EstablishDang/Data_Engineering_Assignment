# app_mssql_server

# Build binary
1. Clone code from git:
```
git clone https://github.com/EstablishDang/Data_Engineering_Assignment.git
```

2. How to build:
```
make
```

# How to run 
1. Example:
```
   ./bin/app_mssql_server
```

# Define APIs
1. GetVendors
- Get with specific name of vendor:
```
   curl --location --request GET 'http://localhost:20000/data_azure/v1/vendors/Australia Bike Retailer'
```

- List all name of Vendors:
```
   curl --location --request GET 'http://localhost:20000/data_azure/v1/vendors/Name'
```

- List all Info of Vendors:
```
   curl --location --request GET 'http://localhost:20000/data_azure/v1/vendors/All'
```

2. GetShipMethod
```
   curl --location --request GET 'http://localhost:20000/data_azure/v1/ship_method'
```

3. GetPurchaseInfo by Id
```
   curl --location --request GET 'http://localhost:20000/data_azure/v1/purchase_order_info/100'
```

4. AddNewVendor
curl --location --request POST 'http://localhost:20000/data_azure/v1/vendor' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "name":"Nhom7",
    "credit_rating":"2",
    "url_web":"nhom7.test.com"
}'