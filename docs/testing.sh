curl -X 'POST' \
  'http://localhost:8080/api/v1/subscriptions' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 10,
  "service_name": "string",
  "start_date": "06-2026",
  "stop_date": "07-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}' -v

curl -X 'POST' \
  'http://localhost:8080/api/v1/subscriptions' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 11,
  "service_name": "YouTube",
  "start_date": "04-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbc"
}' -v

curl -X 'POST' \
  'http://localhost:8080/api/v1/subscriptions' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 15,
  "service_name": "Netflix",
  "start_date": "04-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbc"
}' -v

curl -X 'POST' \
  'http://localhost:8080/api/v1/subscriptions' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 10,
  "service_name": "Netflix",
  "start_date": "06-2026",
  "stop_date": "07-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}' -v

curl -X 'POST' \
  'http://localhost:8080/api/v1/subscriptions' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 0,
  "service_name": "Netflix",
  "start_date": "06-2026",
  "stop_date": "07-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}' -v
{"success":false,"message":"Validation error: Key: 'CreateUserSubscriptionRequest.Price' Error:Field validation for 'Price' failed on the 'required' tag\n"}

curl -X 'POST' \
  'http://localhost:8080/api/v1/subscriptions' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 1,
  "service_name": "",
  "start_date": "06-2026",
  "stop_date": "07-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}' -v
{"success":false,"message":"Validation error: ServiceName=Validation error in tag: 'min'"}

curl -X 'POST' \
  'http://localhost:8080/api/v1/subscriptions' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 1,
  "service_name": "Netflix",
  "start_date": "06-2026",
  "stop_date": "07-2026",
  "user_id": "606"
}' -v

# /api/v1/subscriptions/{id} ###########################################################################################
curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/9' \
  -H 'accept: application/json' -v

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/0' \
  -H 'accept: application/json' -v

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/abc' \
  -H 'accept: application/json' -v

# /api/v1/user/{userID}/subscriptions ##################################################################################
curl -X 'GET' \
  'http://localhost:8080/api/v1/user/{userID}/subscriptions' \
  -H 'accept: application/json' -v

curl -X 'GET' \
  'http://localhost:8080/api/v1/user/60601fee-2bf1-4721-ae6f-7636e79a0cba/subscriptions' \
  -H 'accept: application/json' -v

curl -X 'GET' \
  'http://localhost:8080/api/v1/user/60601fee-2bf1-4721-ae6f-7636e79a0cbb/subscriptions' \
  -H 'accept: application/json' -v

# /api/v1/subscriptions/{id} [put] #####################################################################################
curl -X 'PUT' \
  'http://localhost:8080/api/v1/subscriptions/{id}' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 0,
  "service_name": "string",
  "start_date": "07-2026",
  "stop_date": "07-2026",
  "user_id": "string"
}' -v

curl -X 'PUT' \
  'http://localhost:8080/api/v1/subscriptions/1' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 14,
  "service_name": "MeGoo",
  "start_date": "04-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}' -v

curl -X 'PUT' \
  'http://localhost:8080/api/v1/subscriptions/11' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 0,
  "service_name": "",
  "start_date": "07-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbd"
}' -v

curl -X 'PUT' \
  'http://localhost:8080/api/v1/subscriptions/11' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 14,
  "service_name": "NetflixEx",
  "start_date": "04-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbd"
}' -v

curl -X 'PUT' \
  'http://localhost:8080/api/v1/subscriptions/111' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 14,
  "service_name": "NetflixEx",
  "start_date": "04-2026",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbd"
}' -v

# /api/v1/subscriptions/{id} [patch] ###################################################################################
curl -X 'PATCH' \
  'http://localhost:8080/api/v1/subscriptions/111' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 1,
  "start_date": "07-2026"
}' -v

curl -X 'PATCH' \
  'http://localhost:8080/api/v1/subscriptions/11' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 15,
  "start_date": "05-2026"
}' -v
# /api/v1/subscriptions/{id} [delete] ##################################################################################
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/subscriptions/{id}' \
  -H 'accept: application/json' -v

curl -X 'DELETE' \
  'http://localhost:8080/api/v1/subscriptions/111' \
  -H 'accept: application/json' -v

curl -X 'DELETE' \
  'http://localhost:8080/api/v1/subscriptions/12' \
  -H 'accept: application/json' -v

# /api/v1/subscriptions/total? #########################################################################################
curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=05-2026&date_to=07-2026' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=05-2026&date_to=07-2026&service_name=NetflixEx' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=05-2026&date_to=07-2026&service_name=Netflix' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=05-2026&date_to=07-2026&service_name=string' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=05-2026&date_to=07-2026&service_name=YouTube' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=05-2026&date_to=12-2026&service_name=YouTube' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=04-2026&date_to=07-2026&service_name=YouTube' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=03-2026&date_to=07-2026&service_name=YouTube' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=03-2026&date_to=07-2026&service_name=YouTube' \
  -H 'accept: application/json'

curl -X 'GET' \
  'http://localhost:8080/api/v1/subscriptions/total?date_from=05-2026&date_to=07-2026&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba' \
  -H 'accept: application/json'