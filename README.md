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
event:
  type: json
  fields:
    - name: customerId
      type: string
      required: true
    - name: orderNumber
      type: int
      required: true
    - name: orderStatus
      type: string
      required: true
    - name: nonRequiredField
      type: string
      
inputs:
  - name: webhook-orders
    type: webhook
    decoder:
      type: json
      mappers:
        - source: o_customer_id
          target: customerId
        - source: o_number
          target: orderNumber
        - source: o_status
          target: orderStatus
    params:
      endpoint: /webhook/orders

  - name: mqtt-orders
    type: mqtt
    decoder:
      type: json
      mappers:
        - source: customer_id
          target: customerId
        - source: number
          target: orderNumber
        - source: status
          target: orderStatus
    params:
      broker: mqtt.example.com
      user: my-user
      password: my-password
      topic: my-topic
      
outputs:
  - name: sse-orders
    type: sse
    params:
      endpoint: /sse/orders
      auth:
        type: jwt
        secret: my-secret-key
        encoder:
          type: json
          fields:
            - name: customer_id
              type: string
            - name: customer_email
              type: string
      conditions:
          operator: and
          conditions:
            - field: customerId
              operator: eq
              value: "customer_id"
            - field: orderStatus
              operator: eq
              value: "completed"
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