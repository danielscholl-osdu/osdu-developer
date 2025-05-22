
## Install Process

Either manually modify the values.yaml for the chart or generate a custom_values yaml to use.

_The following commands can help generate a prepopulated custom_values file._
```bash
# Setup Variables
GROUP=<your_group>

# Translate Values File
cat > custom_values.yaml << EOF
################################################################################
# Specify the azure environment specific values
#
azure:
  tenantId: $(az account show --query tenantId -otsv)
  clientId: $(az identity list --resource-group $GROUP --query "[?contains(name, 'service')].clientId" -otsv)
  keyvaultName: $(az keyvault list --resource-group $GROUP --query '[].name' -otsv)

################################################################################
# Specify the resource limits
#
resourceLimits:
  defaultCpuRequests: "0.5"
  defaultMemoryRequests: "4Gi"
  defaultCpuLimits: "1"
  defaultMemoryLimits: "4Gi"
EOF

NAMESPACE=osdu-azure
helm template osdu-developer-base . -f custom_values.yaml

helm upgrade --install osdu-developer-base -f custom_values.yaml . -n $NAMESPACE
```

## Microsoft Identity Secure Enclave (MISE) Integration

This chart includes support for MISE validation as an additional security layer on top of standard Istio JWT validation. MISE validation is disabled by default but can be enabled by setting `mise.enabled: true` in your values.yaml file.

### Prerequisites

- Istio version 1.16.0 or higher
- AKS with Istio Service Mesh enabled
- Recommended Istio revision: `asm-1-23`

### Enabling MISE Validation

To enable MISE validation, add the following configuration to your custom_values.yaml file:

```yaml
################################################################################
# MISE Configuration
#
mise:
  enabled: true
  image:
    repository: mcr.microsoft.com/mise
    tag: latest
  # Use EnvoyFilter by default (set to true if you prefer using IstioOperator mesh config)
  useMeshConfig: false 
```

### Architecture

MISE validation extends the existing JWT validation flow:

1. Istio RequestAuthentication validates JWT signature
2. Envoy Lua filter extracts claims and sets x-user-id and x-app-id headers
3. MISE validation provides additional authorization checks via Envoy ext_authz filter
4. Request only proceeds if MISE validation succeeds

For more details on the MISE integration, see [mise-integration.md](mise-integration.md).

### Mesh Configuration

If you choose to use the IstioOperator mesh configuration (by setting `mise.useMeshConfig: true`), you'll need to apply the configuration separately:

```bash
# Extract the mesh config
helm template osdu-developer-base . -f custom_values.yaml --show-only templates/mise-mesh-config.yaml > mise-mesh-config.yaml

# Apply using istioctl
istioctl apply -f mise-mesh-config.yaml
```