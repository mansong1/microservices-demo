#!/bin/bash

# Array of service names from kubernetes-manifests
services=("adservice" "cartservice" "checkoutservice" "currencyservice" "emailservice" "frontend" "loadgenerator" "paymentservice" "productcatalogservice" "recommendationservice" "shippingservice")

# Harness API details
ACCOUNT_ID="BKB_Vic2RbWnsMSUAWTgXw"
ORG_ID="default"
PROJECT_ID="Resilience"
API_KEY="pat.BKB_Vic2RbWnsMSUAWTgXw.682f933cca2a0932c51ab1de.8k3EQKDplnMLyc5VQKvf"
BASE_URL="https://app.harness.io/gateway/ng/api/servicesV2"

# Function to get the correct container port for each service
get_container_port() {
    local service_name=$1
    case $service_name in
        "adservice") echo "9555" ;;
        "cartservice") echo "7070" ;;
        "checkoutservice") echo "5050" ;;
        "currencyservice") echo "7000" ;;
        "emailservice") echo "8080" ;;
        "frontend") echo "8080" ;;
        "loadgenerator") echo "8080" ;;
        "paymentservice") echo "50051" ;;
        "productcatalogservice") echo "3550" ;;
        "recommendationservice") echo "8080" ;;
        "shippingservice") echo "50051" ;;
        *) echo "8080" ;;
    esac
}

# Function to create a service with GitHub connector for manifests and GCP connector for artifacts
create_service() {
    local service_name=$1
    local container_port=$(get_container_port $service_name)
    
    echo "Creating service: $service_name in project ${PROJECT_ID}"
    echo "Container port: $container_port"
    echo "Repository Name: microservices-demo"
    echo "GitHub manifest path: /kubernetes-manifests/${service_name}.yaml"
    echo "GCP artifact path: google-samples/microservices-demo/${service_name}"
    
    curl -X POST "${BASE_URL}?accountIdentifier=${ACCOUNT_ID}&orgIdentifier=${ORG_ID}&projectIdentifier=${PROJECT_ID}" \
      -H "Content-Type: application/json" \
      -H "x-api-key: ${API_KEY}" \
      -d "{
        \"name\": \"${service_name}\",
        \"identifier\": \"${service_name}\",
        \"orgIdentifier\": \"${ORG_ID}\",
        \"projectIdentifier\": \"${PROJECT_ID}\",
        \"description\": \"${service_name} from microservices demo with GitHub manifests and GCP artifacts\",
        \"tags\": {},
        \"yaml\": \"service:\\n  name: ${service_name}\\n  identifier: ${service_name}\\n  orgIdentifier: ${ORG_ID}\\n  projectIdentifier: ${PROJECT_ID}\\n  serviceDefinition:\\n    type: Kubernetes\\n    spec:\\n      manifests:\\n        - manifest:\\n            identifier: ${service_name}_manifest\\n            type: K8sManifest\\n            spec:\\n              store:\\n                type: Github\\n                spec:\\n                  connectorRef: account.Github\\n                  gitFetchType: Branch\\n                  branch: main\\n                  repoName: microservices-demo\\n                  paths:\\n                    - /kubernetes-manifests/${service_name}.yaml\\n      artifacts:\\n        primary:\\n          primaryArtifactRef: <+input>\\n          sources:\\n            - spec:\\n                connectorRef: account.GCP\\n                imagePath: google-samples/microservices-demo/${service_name}\\n                tag: <+input>\\n                registryHostname: us-central1-docker.pkg.dev\\n              identifier: ${service_name}_artifact\\n              type: Gcr\"
      }"
    
    echo ""
    echo "Service $service_name creation request sent for project ${PROJECT_ID}"
    echo "---"
}

echo "Creating all Harness services in Resilience project with proper project scoping..."
echo "Project ID: ${PROJECT_ID}"
echo "Organization: ${ORG_ID}"
echo "Account: ${ACCOUNT_ID}"
echo "Repository Name: microservices-demo"
echo ""

# Create all services
for service in "${services[@]}"; do
    create_service "$service"
    sleep 2  # Small delay between requests
done

echo "All service creation requests completed!"
echo ""
echo "Services created with:"
echo "✅ Project: ${PROJECT_ID} (properly scoped)"
echo "✅ Organization: ${ORG_ID}"
echo "✅ GitHub Connector: account.Github"
echo "✅ Repository Name: microservices-demo"
echo "✅ Manifest Paths: /kubernetes-manifests/[service-name].yaml"
echo "✅ GCP Connector: account.GCP"
echo "✅ Artifact Registry: us-central1-docker.pkg.dev/google-samples/microservices-demo/[service-name]"
