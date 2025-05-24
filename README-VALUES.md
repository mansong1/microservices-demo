# Image Configuration with values.yaml (Helm Template Format)

This document explains how to use the values.yaml file for managing container images in the microservices demo application using Helm template syntax.

## Overview

The Kubernetes manifests have been updated to reference image values from a centralized `values.yaml` file using Helm template syntax instead of hardcoded image references. This approach enables better integration with Helm charts, CI/CD pipelines, and easier image management.

## File Structure

- `values.yaml` - Contains image configurations and pipeline variables
- `kubernetes-manifests/*.yaml` - Updated Kubernetes manifests using Helm template variables

## values.yaml Structure

```yaml
images:
  adservice: <+artifacts.primary.image>
  cartservice: <+artifacts.primary.image>
  checkoutservice: <+artifacts.primary.image>
  currencyservice: <+artifacts.primary.image>
  emailservice: <+artifacts.primary.image>
  frontend: <+artifacts.primary.image>
  loadgenerator: <+artifacts.primary.image>
  paymentservice: <+artifacts.primary.image>
  productcatalogservice: <+artifacts.primary.image>
  recommendationservice: <+artifacts.primary.image>
  shippingservice: <+artifacts.primary.image>
  
defaultImages:
  adservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/adservice:v0.10.2
  cartservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/cartservice:v0.10.2
  checkoutservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/checkoutservice:v0.10.2
  currencyservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/currencyservice:v0.10.2
  emailservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/emailservice:v0.10.2
  frontend: us-central1-docker.pkg.dev/google-samples/microservices-demo/frontend:v0.10.2
  loadgenerator: us-central1-docker.pkg.dev/google-samples/microservices-demo/loadgenerator:v0.10.2
  paymentservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/paymentservice:v0.10.2
  productcatalogservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/productcatalogservice:v0.10.2
  recommendationservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/recommendationservice:v0.10.2
  shippingservice: us-central1-docker.pkg.dev/google-samples/microservices-demo/shippingservice:v0.10.2

thirdPartyImages:
  redis: redis:alpine
  busybox: busybox:latest
```

## Usage with Helm

### Helm Chart Deployment

```bash
# Deploy using Helm with values.yaml
helm install microservices-demo ./helm-chart -f values.yaml

# Deploy with custom image values
helm install microservices-demo ./helm-chart \
  --set images.frontend=my-registry/frontend:v1.0.0 \
  --set images.adservice=my-registry/adservice:v1.0.0
```

### Helm Template Rendering

```bash
# Render templates to see resolved values
helm template microservices-demo ./helm-chart -f values.yaml

# Render specific service
helm template microservices-demo ./helm-chart -f values.yaml --show-only templates/frontend.yaml
```

## Usage with CI/CD Pipelines

### Harness CI/CD

The `<+artifacts.primary.image>` syntax is compatible with Harness pipelines:

1. **Pipeline Configuration**: Set up artifact sources in your Harness pipeline
2. **Image Resolution**: Harness will automatically substitute `<+artifacts.primary.image>` with the actual image reference
3. **Deployment**: Use Helm or kubectl with resolved image references

### Template Substitution (Non-Helm)

For non-Helm deployments, substitute the Helm template variables:

```bash
# Using helm template command to render manifests
helm template microservices-demo . -f values.yaml | kubectl apply -f -

# Using yq to substitute values manually
yq eval '.images.frontend = "my-registry/frontend:v1.0.0"' values.yaml > custom-values.yaml
helm template microservices-demo . -f custom-values.yaml | kubectl apply -f -
```

## Updated Services (Helm Template Syntax)

All the following services now use Helm template syntax:

- adservice: `{{ .Values.images.adservice }}`
- cartservice: `{{ .Values.images.cartservice }}`
- checkoutservice: `{{ .Values.images.checkoutservice }}`
- currencyservice: `{{ .Values.images.currencyservice }}`
- emailservice: `{{ .Values.images.emailservice }}`
- frontend: `{{ .Values.images.frontend }}`
- loadgenerator: `{{ .Values.images.loadgenerator }}`
- paymentservice: `{{ .Values.images.paymentservice }}`
- productcatalogservice: `{{ .Values.images.productcatalogservice }}`
- recommendationservice: `{{ .Values.images.recommendationservice }}`
- shippingservice: `{{ .Values.images.shippingservice }}`

## Benefits

1. **Helm Compatibility**: Native Helm template syntax for better chart integration
2. **Centralized Management**: All image references in one place
3. **CI/CD Integration**: Easy integration with pipeline variables and Helm
4. **Environment Flexibility**: Different images for different environments using values files
5. **Version Control**: Track image changes through values.yaml
6. **Consistency**: Ensures all services use the same image versioning approach
7. **Template Validation**: Helm can validate template syntax and values

## Example Deployment Commands

```bash
# Using Helm (recommended)
helm install microservices-demo ./helm-chart -f values.yaml

# Using helm template + kubectl
helm template microservices-demo ./helm-chart -f values.yaml | kubectl apply -f -

# With custom values file for different environment
helm install microservices-demo-staging ./helm-chart -f values-staging.yaml

# Upgrade with new image versions
helm upgrade microservices-demo ./helm-chart \
  --set images.frontend=my-registry/frontend:v2.0.0 \
  --reuse-values
```

## Multiple Environment Support

Create environment-specific values files:

```bash
# values-dev.yaml
images:
  frontend: dev-registry/frontend:latest
  adservice: dev-registry/adservice:latest

# values-prod.yaml  
images:
  frontend: prod-registry/frontend:v1.2.3
  adservice: prod-registry/adservice:v1.2.3
```

Deploy to different environments:

```bash
# Development
helm install microservices-demo-dev ./helm-chart -f values-dev.yaml

# Production
helm install microservices-demo-prod ./helm-chart -f values-prod.yaml
```

## Troubleshooting

- Ensure Helm is installed and the chart structure is correct
- Use `helm template` to validate template rendering before deployment
- Check that image references are accessible from your Kubernetes cluster
- Verify that CI/CD pipeline has proper permissions to resolve artifact references
- Use `helm lint` to validate chart syntax and structure
- Check values.yaml syntax with `helm template --validate`
