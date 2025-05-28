# Chaos Engineering as a Service (CEaaS)

A self-service platform for running chaos engineering experiments on your Kubernetes infrastructure. This project provides a complete solution for implementing chaos engineering practices in your organization, allowing you to proactively identify weaknesses in your systems before they cause real outages.

## Features

- 🚀 Automated failure injection in Kubernetes clusters
- 🔒 Safe experimentation with automatic rollback
- 📊 Real-time monitoring and visualization
- 🛠️ Easy integration with existing CI/CD pipelines
- 🧪 Multiple chaos experiment types (pod failure, network latency, CPU/memory hogs)
- 🌐 Modern web dashboard for experiment management

## Architecture

The system consists of the following components:

- **Chaos Controller**: Kubernetes controller that manages the lifecycle of chaos experiments
- **Chaos Experiments**: Implementations of different chaos scenarios (pod failure, network latency, etc.)
- **API Server**: RESTful API for managing experiments programmatically
- **Web Dashboard**: React-based UI for monitoring and controlling experiments

## Prerequisites

- Kubernetes cluster (Minikube, Docker Desktop, or cloud-based)
- kubectl configured to communicate with your cluster
- Go 1.22+ (for development)
- Docker (for building container images)
- Helm (for deployment)
- Node.js and npm (for dashboard development)

## Quick Start

### Using Helm (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/chaos-engineering-service.git
   cd chaos-engineering-service
   ```

2. Install using Helm:
   ```bash
   helm install chaos-engineering ./charts/chaos-engineering-service
   ```

3. Access the dashboard:
   ```bash
   kubectl port-forward svc/chaos-api-server 8080:80 -n chaos-engineering
   ```
   Then open http://localhost:8080 in your browser.

### Manual Installation

1. Apply the CRDs:
   ```bash
   kubectl apply -f deploy/kubernetes/crds/
   ```

2. Deploy the controller and API server:
   ```bash
   kubectl apply -f deploy/kubernetes/deployment.yaml
   ```

3. Access the dashboard:
   ```bash
   kubectl port-forward svc/chaos-api-server 8080:80 -n chaos-engineering
   ```

## Running Chaos Experiments

### Using the Dashboard

1. Open the dashboard at http://localhost:8080
2. Click on "New Experiment"
3. Fill in the experiment details:
   - Name: A unique name for your experiment
   - Namespace: The Kubernetes namespace where the target resource is located
   - Target Kind: The kind of resource to target (Pod, Deployment, etc.)
   - Target Name: The name of the resource to target
   - Experiment Type: The type of chaos to inject
   - Duration: How long the experiment should run
   - Parameters: Specific parameters for the chosen experiment type
4. Click "Create Experiment" to start the chaos experiment

### Using kubectl

1. Create a YAML file for your experiment (see examples in the `examples/` directory)
2. Apply it using kubectl:
   ```bash
   kubectl apply -f examples/pod-failure-experiment.yaml
   ```

3. Monitor the experiment status:
   ```bash
   kubectl get chaosexperiments -n chaos-test
   ```

## Available Chaos Experiments

| Experiment Type | Description | Parameters |
|----------------|-------------|------------|
| pod-failure | Kills a pod to test resilience to pod failures | None |
| network-latency | Adds latency to network traffic | latency, jitter |
| cpu-hog | Consumes CPU resources | cpuCores |
| memory-hog | Consumes memory resources | memoryMB |

## Development

### Building the Project

```bash
# Download dependencies
go mod tidy

# Generate client code
make generate

# Build all components
make build

# Build the dashboard
make dashboard-build

# Build Docker images
make docker-build
```

### Running Locally

```bash
# Run the API server (which serves the dashboard)
make run-api
```

### Project Structure

- `cmd/controller/`: Controller entry point
- `pkg/chaos/apis/`: API definitions for CRDs
- `pkg/chaos/experiments/`: Chaos experiment implementations
- `pkg/controller/`: Controller implementation
- `api/`: API server implementation
- `dashboard/`: React dashboard
- `deploy/`: Kubernetes deployment manifests
- `charts/`: Helm charts
- `examples/`: Example chaos experiments
- `hack/`: Development scripts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Apache License 2.0
