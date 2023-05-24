#!/bin/sh

for i in {1..10}; do
    curl --location 'http://localhost:4000/webhook/orders' \
    --header 'Content-Type: application/json' \
    --data '{
        "o_customer_id": "titi",
        "o_number": 1122,
        "o_status": "completed"
    }'

    echo "call to webhook/orders $i"
    sleep 2
done
