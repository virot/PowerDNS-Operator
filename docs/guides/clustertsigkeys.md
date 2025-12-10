# ClusterTSIGKey deployment

## Specification

The `ClusterTSIGKey` specification is identical to the TSIGKey specification but is cluster-scoped rather than namespace-scoped.

| Field | Type | Required | Description |
| ----- | ---- |:--------:| ----------- |
| algorithm | string | Y | Algorithm for the TSIG key, one of "hmac-md5", "hmac-sha1", "hmac-sha224", "hmac-sha256", "hmac-sha384", "hmac-sha512" |
| secretRef | SecretReference | Y | Reference to a Kubernetes Secret containing the TSIG key value |

The `SecretReference` contains:

| Field | Type | Required | Description |
| ----- | ---- |:--------:| ----------- |
| name | string | Y | Name of the secret |
| namespace | string | Y | Namespace of the secret (required for ClusterTSIGKey) |

## Example

### Creating a Secret with the TSIG Key

First, create a Kubernetes Secret containing the TSIG key value:

```bash
kubectl create secret generic global-tsig-key \
  --from-literal=key='your-base64-encoded-key-here' \
  --namespace powerdns-operator-system
```

### Creating a ClusterTSIGKey Resource

```yaml
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: ClusterTSIGKey
metadata:
  name: global-update-key
spec:
  algorithm: hmac-sha256
  secretRef:
    name: global-tsig-key
    namespace: powerdns-operator-system
```

### Using the ClusterTSIGKey in a ClusterZone

Once created, reference the TSIG key in your ClusterZone definition:

```yaml
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: ClusterZone
metadata:
  name: example.org
spec:
  kind: Master
  nameservers:
    - ns1.example.org
    - ns2.example.org
  tsigKeyIds:
    - global-update-key
```

## Notes

- ClusterTSIGKey is a cluster-scoped resource and can be referenced by any Zone or ClusterZone
- The secret namespace must be explicitly specified in the secretRef for ClusterTSIGKey
- The secret must contain a key named `key` with the base64-encoded TSIG key value
- The ClusterTSIGKey resource name will be used as the key ID in PowerDNS
- Useful for shared TSIG keys that are used across multiple zones in different namespaces

## Use Cases

### Cluster-Wide TSIG Keys

Use ClusterTSIGKey when you have a TSIG key that should be available to all zones in the cluster:

```yaml
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: ClusterTSIGKey
metadata:
  name: primary-transfer-key
spec:
  algorithm: hmac-sha512
  secretRef:
    name: primary-transfer-key
    namespace: powerdns-operator-system
```

### Namespace-Specific Zones Using Cluster Keys

Zones in any namespace can reference the cluster-wide TSIG key:

```yaml
# Zone in namespace "team-a"
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: Zone
metadata:
  name: team-a.example.com
  namespace: team-a
spec:
  kind: Master
  nameservers:
    - ns1.example.com
  tsigKeyIds:
    - primary-transfer-key
---
# Zone in namespace "team-b"
apiVersion: dns.cav.enablers.ob/v1alpha2
kind: Zone
metadata:
  name: team-b.example.com
  namespace: team-b
spec:
  kind: Master
  nameservers:
    - ns1.example.com
  tsigKeyIds:
    - primary-transfer-key
```
