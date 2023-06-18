# Ostraka

Ostraka is an open-source project that simplifies event synchronization in microservice systems. It allows you to transform data from various sources, such as gRPC, Kafka, MQTT, and webhooks, into Server-Sent Events (SSE) or email notifications for easy consumption in the user interface. The project aims to provide a flexible and configurable solution for event synchronization between microservices.

## Features

- Configuration Management: Ostraka allows you to configure and manage different types of inputs, including gRPC, Kafka, MQTT, and webhooks, to collect events from diverse systems.
- Event Transformation: Incoming events from inputs are transformed to the SSE format, enabling continuous streaming to the frontend.
- Flexible Configuration: Ostraka provides a configuration system based on YAML files, allowing users to define inputs and outputs according to their specific needs.
- Modularity: The project is designed to be extensible through modules, providing the ability to add middleware for authentication, logging, and other functionalities to inputs and outputs.

## Getting Started

### Prerequisites

- Go programming language (version > 1.18)
- Docker
- Make

### Installation

1. Clone the repository:

```shell
git clone git@github.com:valensto/ostraka.git
```

### Configuration

Create a YAML configuration file for your contexts, inputs, and outputs. See the Configuration section below for more details.

Place the YAML configuration file in the config directory `config/resources`.

### Usage

Run the Ostraka microservice using Docker:
    
```shell
make dev
```

To run test

```shell
make test
```

To run sse example

```shell
make sse-example
```

Ostraka will start and load the configuration from the YAML file.

Access the Ostraka SSE endpoint in your browser or make API calls to interact with the microservice.

### Configuration

Ostraka uses YAML files for configuration. You can define the inputs and outputs based on your requirements.

Example YAML configuration:

```yaml
name: incoming-orders

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
    - name: default
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
    - name: default
      type: jwt
      params:
        header: Authorization
        secret: secret
        algorithm: HS256
        verify_expiration: true
        max_age: 3600
        payload:
          - name: accountId
            data_type: string
            required: true

    - name: webhook
      type: token
      params:
        token: 2dc7929e5b589cb7861bcae19e13ad96
        query_param: token

inputs:
  - name: webhook-orders
    source: webhook
    params:
      endpoint: /webhook/orders
      auth: webhook
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
    condition:
      operator: or
      conditions:
        - field: orderStatus
          operator: eq
          value: "completed"
        - field: orderStatus
          operator: eq
          value: "pending"
  - name: sse-orders-failed
    destination: sse
    params:
      endpoint: /sse/orders/failed
      cors: default
      auth: default
    condition:
      field: orderStatus
      operator: eq
      value: "failed"
```

Modify the configuration file to match your specific inputs, outputs, and other settings.

### Contributing

Contributions to Ostraka are welcome! If you would like to contribute to the project, please follow these steps:

1. Fork the repository on GitHub.
2. Create a new branch for your feature or bug fix:
```shell
git checkout -b my-feature
```
3. Make the necessary changes and commit them:
4. Push your changes to your forked repository:
5. Open a pull request on the original repository. Provide a clear description of your changes and any relevant information.
6. Your pull request will be reviewed by the project maintainers. They may provide feedback or request further changes.
7. Once your pull request is approved, it will be merged into the main repository. Congratulations on your contribution!

Please ensure that your contributions adhere to the project's coding conventions and standards. Also, make sure to include tests for any new functionality or bug fixes.

### License

Ostraka is licensed under the MIT License. See the LICENSE file for more information.

### Contact

If you have any questions, suggestions, or feedback, you can reach out to the project maintainers at hi@valensto.com.

---

Thank you for your interest in Ostraka! We hope you find the project useful and look forward to your contributions.