# Helm Chart for Istio Ingress Gateways

This chart deploys and configures Istio ingress gateways with optional Microsoft Identity Sidecar Extension (MISE) for live token validation.

## Installation

Install the helm chart with the following commands:

```bash
# Create Namespace
NAMESPACE=istio-system

# Create a custom values file
cat > custom_values.yaml << EOF
# Enable MISE sidecar
mise:
  enabled: true
  image: mise-mock:latest
  imagePullPolicy: IfNotPresent

# Configure auth filter
authFilter:
  enabled: true
  namespace: istio-system
EOF

# Install the chart
helm upgrade --install istio-ingress . -n $NAMESPACE -f custom_values.yaml
```

## Microsoft Identity Sidecar Extension (MISE)

This chart includes support for deploying the MISE sidecar container with the Istio ingress gateway. MISE provides real-time token validation against Microsoft Entra ID, meeting security requirements for public-facing OSDU APIs.

### MISE Configuration

The MISE sidecar can be configured through the values file:

```yaml
mise:
  enabled: true                    # Enable the MISE sidecar
  image: mise-mock:latest          # Mock image for testing
  imagePullPolicy: IfNotPresent
  port: 9002                       # Port for the MISE service
  resources:                       # Resource limits and requests
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 256Mi
  envVars:                         # Environment variables for MISE
    LOG_LEVEL: "info"
    DEFAULT_LATENCY_MS: "10"

authFilter:
  enabled: true                    # Enable the ext_authz filter
  namespace: istio-system          # Namespace for the filter
  workloadSelector:                # Which workloads to apply to
    labels:
      istio: ingressgateway
  port: 80                         # Port to apply the filter
  authHost: localhost              # Host for auth service
  authPort: 9002                   # Port for auth service
  timeout: "0.2s"                  # Timeout for auth requests
  pathPrefix: "/ext_authz"         # Path for auth requests
```

For production deployments, replace `mise-mock` with the actual MISE implementation image that connects to Microsoft Entra ID.
