# Microsoft Identity Sidecar Extension (MISE) Mock

This is a mock implementation of the Microsoft Identity Sidecar Extension (MISE) for the OSDU platform. It provides a simulated authentication service that acts as a placeholder for the real MISE implementation, which will integrate with Microsoft Entra ID for live token validation.

## Purpose

As outlined in the MISE design requirements, Microsoft Entra now requires live (online) token validation for all public-facing OSDU APIs. To facilitate this transition, the MISE mock implementation:

1. Acts as a sidecar container in the istio-ingressgateway pods
2. Provides an ext_authz endpoint for Envoy to validate requests
3. Simulates the behavior of the real MISE service for testing and development

## Features

- Mock implementation of the MISE sidecar for Istio ingress gateway
- Support for ext_authz protocol integration with Envoy
- Configurable response headers
- Support for simulated latency
- Test-friendly interfaces

## Configuration

The mock MISE container is configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server listening port | `9002` |
| `LOG_LEVEL` | Logging verbosity (info, debug) | `info` |
| `DEFAULT_LATENCY_MS` | Default added latency in milliseconds | `0` |

## Testing Headers

For testing purposes, the mock understands specific headers:

| Header | Description |
|--------|-------------|
| `x-test-latency-ms` | Simulate specific response latency in milliseconds |
| `test-user-id` | Set a specific user ID in the response headers |
| `test-app-id` | Set a specific app ID in the response headers |
| `return-status-code` | Force a specific status code in the response |
| `return-*` | Any header prefixed with `return-` will be included in the response headers |

## Installation

The mock MISE container is automatically deployed as a sidecar with the istio-ingressgateway when using the provided Helm charts.

```bash
# Deploy using Helm
helm upgrade --install istio-ingress ./charts/istio-ingress --namespace istio-system
```

## Architecture

The MISE mock sidecar works as follows:

1. Requests arrive at the Istio ingress gateway.
2. The EnvoyFilter directs auth requests to the local MISE mock sidecar.
3. MISE mock validates the provided token (simulated).
4. Upon successful validation, MISE mock adds the necessary identity headers.
5. The request continues to its destination service.

This architecture enables testing of the full MISE workflow without requiring actual connectivity to Microsoft Entra ID.