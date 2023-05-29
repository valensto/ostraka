#!/bin/sh

echo "Starting to call webhook/orders every 2 seconds... Press Ctrl+C to stop."
echo "Open $PWD/index.html in your browser to see the events."

while true; do
    o_customer_id=$(uuidgen | tr -d '-')
    o_number=$(od -An -N2 -i /dev/urandom | awk '{ print $1 }')
    o_number=$((o_number % 9000 + 1000))
    options=("completed" "pending" "failed")
    random_index=$(awk -v min=0 -v max=$((${#options[@]} - 1)) 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    o_status=${options[$random_index]}

    curl --location 'http://localhost:4000/webhook/orders' \
    --header 'Content-Type: application/json' \
    --data '{
        "o_customer_id": "'"$o_customer_id"'",
        "o_number": '"$o_number"',
        "o_status": "'"$o_status"'"
    }'

    echo "Call to webhook/orders order number: $o_number order status: $o_status"
    sleep 3
done
