# Keeprice

General purpose, scalable e-shop price scraper storage api

## Usage

```bash
keeprice api --storage=<influxdb | cassandra>
```

## Api 

Add new product 

```
POST 
/product 
{
    "site_id": "tannico",
    "product_id": "<product_url>",
    "product_price": "<product_price>",
    "scrape_date": <unix epoch sec>
}
```

Get product history 
```
GET
/product?site_id=<site_id>&product_id=<product_id>




```


## Supported backend storage

* InfluxDB

* Cassandra 