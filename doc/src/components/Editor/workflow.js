export const workflow = `name: incoming-orders

event_type:
  format: json
  fields:
    - name: customerId
      data_type: string
      required: true
    - name: orderNumber
      data_type: int
      required: true
    - name: orderStatus
      data_type: string
      required: true
    - name: nonRequiredField
      data_type: string

middlewares:
  cors:
    default:
      allowed_origins:
        - http://localhost:3000
        - http://localhost:4000
      allowed_methods:
        - POST
      allowed_headers:
        - Content-Type
        - Authorization
      allow_credentials: true
      max_age: 3600

  auth:
    default:
      type: token
      params:
        token: 2dc7929e5b589cb7861bcae19e13ad96
        query_param: token

inputs:
  - name: webhook-orders
    source: webhook
    params:
      endpoint: /webhook/orders
      auth: default
    decoder:
      format: json
      mappers:
        - source: o_customer_id
          target: customerId
        - source: o_number
          target: orderNumber
        - source: o_status
          target: orderStatus

outputs:
  - name: sse-orders-completed
    destination: sse
    params:
      endpoint: /sse/orders/completed
      auth: default
      cors: default
    encoder:
      format: json
    condition:
      operator: or
      conditions:
        - field: orderStatus
          operator: eq
          value: "completed"
        - field: orderStatus
          operator: eq
          value: "pending"

  - name: smtp-orders-failed
    destination: smtp
    params:
      from: hi@valensto.com
      to: v.e.brochard@gmail.com
      subject: Order failed
      base_url: http://localhost:4000
      host: smtp.gmail.com
      port: "587"
      username: hi@valensto.com
      password: bbfrhiypmitrixyw
      enable_starttls: true
    encoder:
      format: html
    condition:
      field: orderStatus
      operator: eq
      value: "failed"`