# TSIGKey deployment

## Specification

The `TSIGKey` specification contains the following fields:

| Field | Type | Required | Description |
| ----- | ---- |:--------:| ----------- |
| algorithm | string | Y | Algorithm for the TSIG key, one of "hmac-md5", "hmac-sha1", "hmac-sha224", "hmac-sha256", "hmac-sha384", "hmac-sha512" |
| secretRef | SecretReference | N* | Reference to a Kubernetes Secret containing the TSIG key value (recommended) |
| key | string | N* | Base64-encoded TSIG key value directly in the spec (see security warning below) |

*Either `secretRef` or `key` must be provided, but not both.

The `SecretReference` contains:

| Field | Type | Required | Description |
| ----- | ---- |:--------:| ----------- |
| name | string | Y | Name of the secret in the same namespace |
| namespace | string | N | Namespace of the secret (only used for ClusterTSIGKey) |

## Security Warning

⚠️ **WARNING**: Using the `key` field to store the TSIG key directly in the CRD is **insecure** as it will be visible in plain text in the Kubernetes API and stored unencrypted in etcd. This approach should only be used for testing or non-production environments. For production use, always use `secretRef` to store the key in a Kubernetes Secret.

## Examples

### Method 1: Using a Secret (Recommended)

First, create a Kubernetes Secret containing the TSIG key value:

```bash
kubectl create secret generic my-tsig-key \
  --from-literal=key='your-base64-encoded-key-here' \
  --namespace default
```

Then create a TSIGKey resource:

```yaml
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: TSIGKey
metadata:
  name: update-key
  namespace: default
spec:
  algorithm: hmac-sha256
  secretRef:
    name: my-tsig-key
```

### Method 2: Inline Key (Not Recommended for Production)

⚠️ **This method stores the key in plain text and is not secure. Use only for testing.**

```yaml
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: TSIGKey
metadata:
  name: update-key
  namespace: default
spec:
  algorithm: hmac-sha256
  key: "your-base64-encoded-key-here"
```

### Using the TSIG Key in a Zone

Once created, reference the TSIG key in your Zone definition:

```yaml
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: Zone
metadata:
  name: example.com
  namespace: default
spec:
  kind: Master
  nameservers:
    - ns1.example.com
    - ns2.example.com
  tsigKeyIds:
    - update-key
```

## Notes

- The secret must contain a key named `key` with the base64-encoded TSIG key value
- The TSIGKey resource name will be used as the key ID in PowerDNS
- For namespace-scoped TSIGKey resources, the secret must be in the same namespace
- For ClusterTSIGKey resources, you can specify the namespace in the secretRef

## Reconciliation Flow

The following diagram illustrates the reconciliation flow for TSIGKey resources:

```mermaid
sequenceDiagram
    participant U as User
    participant K as Kubernetes API
    participant C as Controller
    participant S as Secret
    participant P as PowerDNS API
    
    U->>K: kubectl apply tsigkey.yaml
    K->>C: TSIGKey Created Event
    
    Note over C: Reconciliation Loop Starts
    C->>C: Check Deletion Timestamp
    
    alt Resource is being deleted
        C->>P: DELETE /api/v1/servers/localhost/tsigkeys/{id}
        P-->>C: TSIG Key Deleted
        C->>C: Remove Finalizers
        C->>K: Update TSIGKey
        Note over C: Deletion Complete
    else Resource is being created/updated
        C->>C: Add Finalizers if missing
        C->>S: GET Secret
        
        alt Secret not found
            C->>K: Set Failed Status
            Note over C: Reconciliation Failed
        else Secret found
            C->>P: GET /api/v1/servers/localhost/tsigkeys/{id}
            
            alt TSIG Key doesn't exist in PowerDNS
                C->>P: POST /api/v1/servers/localhost/tsigkeys
                Note over P: Create TSIG Key
                P-->>C: TSIG Key Created Successfully
            else TSIG Key exists in PowerDNS
                C->>C: Compare desired vs actual state
                alt Differences found
                    C->>P: PUT /api/v1/servers/localhost/tsigkeys/{id}
                    P-->>C: TSIG Key Updated Successfully
                end
            end
            
            C->>K: Update TSIGKey Status
            C->>K: Set Available Condition
            Note over C: Reconciliation Succeeded
        end
    end
    
    K-->>U: TSIGKey Status Updated
```
