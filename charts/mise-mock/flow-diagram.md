# MISE Integration Flow Diagram

The following diagram illustrates the logical flow of requests through the MISE sidecar:

```mermaid
graph TD
    A[Client] -->|1. HTTP Request| B[Istio Ingress Gateway]
    B -->|2. ext_authz request| C[MISE Sidecar]
    C -->|3. Check Auth Header| D{Auth Header Valid?}
    D -->|No| E[Return 403]
    D -->|Yes| F[Process User Identity]
    F -->|4. Set User Headers| G[Return 200 OK]
    G -->|5. Headers added| B
    B -->|6. Forward Request| H[Backend Service]
    E -->|Reject| B
    B -->|Rejected| A
```

## Request Flow

1. Client sends a request to an OSDU API endpoint
2. Istio ingress gateway receives the request and routes it to the MISE sidecar
3. MISE sidecar validates the authorization header
4. If the header is valid, MISE identifies the user and sets identity headers
5. Request is forwarded to the backend service with added headers
6. If validation fails, the request is rejected with a 403 error

## Header Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant E as Envoy Proxy
    participant M as MISE Sidecar
    participant S as Service
    
    C->>E: HTTP Request with Authorization Header
    E->>M: ext_authz request
    Note over M: Validate token
    alt Valid Token
        M->>E: 200 OK + x-user-id + x-app-id
        E->>S: Original request + identity headers
        S->>E: Service response
        E->>C: Response
    else Invalid Token
        M->>E: 403 Forbidden + WWW-Authenticate
        E->>C: 403 Forbidden + WWW-Authenticate
    end
```