
```console
root@pokt.network:~$ git clone https://github.com/pokt-foundation/gateway-settings-api.git
root@pokt.network:~$ cd gateway-settings-api
root@pokt.network:~$ go run main.go
```

1.Authentication (using JWTs)

Request
```console
root@pokt.network:~$ curl --location --request POST 'https://settings-api.portal.pokt.network/v1/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "your portal.pokt.network's account email",
    "password": "your portal.pokt.network's account password"
}'
```

Response (200 is successful, 401 if the email or password are wrong)
```json
{
    "data": "ey....4",
    "message": "Success login",
    "status": "success"
}
```

2. Adding a contract (use Bearer Token auth for the req)

`blockchain_id` can be found [here](https://docs.pokt.network/home/supported-blockchains#current-relaychains).

```console
curl --location --request POST 'https://settings-api.portal.pokt.network/v1/settings/add-contract' \
--header 'Authorization: Bearer ey....4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "your application ID",
    "gateway_settings": {
        "contracts_allowlist": [
            {
                "blockchain_id": "0053",
                "contracts": ["0x2....aecf652"]
            }
        ]
    }
}'
```
