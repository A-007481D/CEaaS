#!/bin/bash

# Test script for Chaos Engineering as a Service
set -e

echo "===== Testing Chaos Engineering as a Service ====="

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "kubectl is not installed. Please install it first."
    exit 1
fi

# Check if we're connected to a Kubernetes cluster
if ! kubectl get nodes &> /dev/null; then
    echo "Not connected to a Kubernetes cluster. Please configure kubectl."
    exit 1
fi

echo "✅ Kubernetes connection verified"

# Create the test namespace and deploy the test application
echo "📦 Deploying test application..."
kubectl apply -f examples/test-app.yaml
echo "✅ Test application deployed"

# Wait for the test application to be ready
echo "⏳ Waiting for test application to be ready..."
kubectl wait --for=condition=available --timeout=60s deployment/nginx-test -n chaos-test
echo "✅ Test application is ready"

# Apply the CRDs
echo "📦 Applying Custom Resource Definitions..."
kubectl apply -f deploy/kubernetes/crds/chaosexperiment.yaml
echo "✅ CRDs applied"

# Deploy the controller
echo "📦 Deploying Chaos Controller..."
kubectl apply -f deploy/kubernetes/deployment.yaml
echo "✅ Chaos Controller deployed"

# Wait for the controller to be ready
echo "⏳ Waiting for Chaos Controller to be ready..."
kubectl wait --for=condition=available --timeout=60s deployment/chaos-controller -n chaos-engineering
echo "✅ Chaos Controller is ready"

# Run a pod failure experiment
echo "🧪 Running pod failure experiment..."
kubectl apply -f examples/pod-failure-experiment.yaml
echo "✅ Pod failure experiment started"

# Wait for a moment
sleep 5

# Check the experiment status
echo "📊 Checking experiment status..."
kubectl get chaosexperiments -n chaos-test

echo ""
echo "===== Test Complete ====="
echo "You can access the dashboard by running:"
echo "kubectl port-forward svc/chaos-api-server 8080:80 -n chaos-engineering"
echo "Then open http://localhost:8080 in your browser"
