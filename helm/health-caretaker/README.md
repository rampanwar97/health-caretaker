# Health Caretaker Helm Chart

A Helm chart for deploying the Health Caretaker application - a modern, real-time health monitoring application for HTTP/HTTPS endpoints.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- Prometheus Operator (optional, for ServiceMonitor and PrometheusRule)

## Installing the Chart

To install the chart with the release name `my-health-caretaker`:

```bash
helm install my-health-caretaker ./helm/health-caretaker
```

The command deploys Health Caretaker on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-health-caretaker` deployment:

```bash
helm uninstall my-health-caretaker
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Parameters

### Global parameters

| Name                      | Description                                     | Value |
| ------------------------- | ----------------------------------------------- | ----- |
| `nameOverride`            | String to partially override common.names.fullname | `""` |
| `fullnameOverride`        | String to fully override common.names.fullname | `""` |

### Image parameters

| Name                | Description                                                                 | Value           |
| ------------------- | --------------------------------------------------------------------------- | --------------- |
| `image.repository`  | Health Caretaker image repository                                           | `your-username/health-caretaker` |
| `image.pullPolicy`  | Image pull policy                                                            | `IfNotPresent`  |
| `image.tag`         | Health Caretaker image tag (immutable tags are recommended)                 | `""`            |

### Image pull secrets

| Name                        | Description                                                                 | Value  |
| --------------------------- | --------------------------------------------------------------------------- | ------ |
| `imagePullSecrets`          | Specify docker-registry secret names as an array                           | `[]`   |

### Service Account

| Name                                 | Description                                                                 | Value   |
| ------------------------------------ | --------------------------------------------------------------------------- | ------- |
| `serviceAccount.create`              | Specifies whether a service account should be created                       | `true`  |
| `serviceAccount.annotations`         | Annotations to add to the service account                                   | `{}`    |
| `serviceAccount.name`                | The name of the service account to use                                      | `""`    |

### Pod Security Context

| Name                | Description                                                                 | Value   |
| ------------------- | --------------------------------------------------------------------------- | ------- |
| `podSecurityContext` | Pod security context                                                        | `{}`    |

### Container Security Context

| Name                | Description                                                                 | Value   |
| ------------------- | --------------------------------------------------------------------------- | ------- |
| `securityContext`   | Container security context                                                  | `{}`    |

### Service parameters

| Name                    | Description                                                                 | Value       |
| ----------------------- | --------------------------------------------------------------------------- | ----------- |
| `service.type`          | Kubernetes Service type                                                     | `ClusterIP` |
| `service.webPort`       | Service web port                                                             | `8080`      |
| `service.metricsPort`   | Service metrics port                                                         | `9091`      |
| `service.annotations`   | Additional custom annotations for the service                               | `{}`        |

### Ingress parameters

| Name                       | Description                                                                 | Value                    |
| -------------------------- | --------------------------------------------------------------------------- | ------------------------ |
| `ingress.enabled`          | Enable ingress record generation for Health Caretaker                      | `false`                  |
| `ingress.className`        | IngressClass resource name                                                  | `""`                     |
| `ingress.annotations`      | Additional annotations for the Ingress resource                             | `{}`                     |
| `ingress.hosts`            | An array of hosts to be covered with this ingress record                   | `[]`                     |
| `ingress.tls`              | TLS configuration for additional hostname(s) to be covered with this ingress record | `[]` |

### Resource parameters

| Name                       | Description                                                                 | Value   |
| -------------------------- | --------------------------------------------------------------------------- | ------- |
| `resources.limits`         | The resources limits for the Health Caretaker containers                    | `{}`    |
| `resources.requests`       | The requested resources for the Health Caretaker containers                 | `{}`    |

### Autoscaling parameters

| Name                                       | Description                                                                 | Value   |
| ------------------------------------------ | --------------------------------------------------------------------------- | ------- |
| `autoscaling.enabled`                      | Enable Horizontal POD autoscaling for Health Caretaker                     | `false` |
| `autoscaling.minReplicas`                  | Minimum number of Health Caretaker replicas                                | `1`     |
| `autoscaling.maxReplicas`                  | Maximum number of Health Caretaker replicas                                | `10`    |
| `autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization percentage                                          | `80`    |
| `autoscaling.targetMemoryUtilizationPercentage` | Target Memory utilization percentage                                    | `80`    |

### Other parameters

| Name                       | Description                                                                 | Value   |
| -------------------------- | --------------------------------------------------------------------------- | ------- |
| `replicaCount`             | Number of Health Caretaker replicas to deploy                              | `1`     |
| `podAnnotations`           | Annotations for Health Caretaker pods                                       | `{}`    |
| `nodeSelector`             | Node labels for pod assignment                                              | `{}`    |
| `tolerations`              | Tolerations for pod assignment                                              | `[]`    |
| `affinity`                 | Affinity for pod assignment                                                 | `{}`    |

### Health Caretaker Configuration

| Name                       | Description                                                                 | Value   |
| -------------------------- | --------------------------------------------------------------------------- | ------- |
| `config.endpoints`         | List of endpoints to monitor                                                | `[]`    |
| `env.WEB_PORT`             | Web UI and API port                                                         | `8080`  |
| `env.METRICS_PORT`         | Prometheus metrics port                                                     | `9091`  |
| `env.METRICS_ENABLED`      | Enable/disable metrics endpoint                                             | `true`  |
| `env.METRICS_PATH`         | Path for metrics endpoint                                                   | `/metrics` |

### Pod Disruption Budget

| Name                       | Description                                                                 | Value   |
| -------------------------- | --------------------------------------------------------------------------- | ------- |
| `podDisruptionBudget.enabled` | Enable Pod Disruption Budget                                               | `false` |
| `podDisruptionBudget.minAvailable` | Minimum number of pods that must be available                             | `1`     |

### Network Policy

| Name                       | Description                                                                 | Value   |
| -------------------------- | --------------------------------------------------------------------------- | ------- |
| `networkPolicy.enabled`    | Enable Network Policy                                                       | `false` |
| `networkPolicy.ingress`    | Ingress rules for Network Policy                                           | `[]`    |
| `networkPolicy.egress`     | Egress rules for Network Policy                                            | `[]`    |

### ServiceMonitor (Prometheus Operator)

| Name                       | Description                                                                 | Value   |
| -------------------------- | --------------------------------------------------------------------------- | ------- |
| `serviceMonitor.enabled`   | Enable ServiceMonitor for Prometheus Operator                              | `false` |
| `serviceMonitor.namespace` | Namespace for ServiceMonitor                                                | `""`    |
| `serviceMonitor.interval`  | Scrape interval                                                             | `30s`   |
| `serviceMonitor.scrapeTimeout` | Scrape timeout                                                          | `10s`   |
| `serviceMonitor.labels`    | Labels for ServiceMonitor                                                   | `{}`    |
| `serviceMonitor.annotations` | Annotations for ServiceMonitor                                           | `{}`    |


## Configuration and installation details

### ConfigMap

The application configuration is managed via a ConfigMap that contains the `config.json` file. The ConfigMap is automatically generated from the values in `values.yaml` under the `config.endpoints` section.

### Endpoint Configuration

Endpoints are configured in the `values.yaml` file under `config.endpoints`. Each endpoint can have the following properties:

- `name`: Friendly name for the endpoint
- `url`: HTTP/HTTPS URL to monitor
- `method`: HTTP method (GET, POST, PUT, DELETE, HEAD)
- `interval`: Check interval in seconds
- `timeout`: Request timeout in seconds
- `labels`: Custom labels for metrics (optional)
- `probe_type`: Type of probe (optional)

### Example Configuration

```yaml
config:
  endpoints:
    - name: "Google"
      url: "https://www.google.com"
      method: "GET"
      interval: 30
      timeout: 10
      labels:
        service: "search"
        environment: "production"
        team: "platform"
        criticality: "high"
```

### Updating Configuration

To update the endpoint configuration:

1. Edit the `values.yaml` file
2. Run: `helm upgrade my-health-caretaker ./helm/health-caretaker`

## Examples

### Basic Installation

```bash
helm install my-health-caretaker ./helm/health-caretaker
```

### Installation with Custom Values

```bash
helm install my-health-caretaker ./helm/health-caretaker \
  --set image.tag=v1.0.0 \
  --set replicaCount=3 \
  --set service.type=LoadBalancer
```

### Installation with Custom Configuration

```bash
helm install my-health-caretaker ./helm/health-caretaker \
  --set-file config.custom=./my-config.json
```

### Installation with Ingress

```bash
helm install my-health-caretaker ./helm/health-caretaker \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=health-caretaker.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

### Installation with Prometheus Monitoring

```bash
helm install my-health-caretaker ./helm/health-caretaker \
  --set serviceMonitor.enabled=true
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -l app.kubernetes.io/name=health-caretaker
```

### View Logs

```bash
kubectl logs -l app.kubernetes.io/name=health-caretaker
```

### Check ConfigMap

```bash
kubectl get configmap -l app.kubernetes.io/name=health-caretaker -o yaml
```

### Port Forward for Testing

```bash
kubectl port-forward svc/my-health-caretaker 8080:8080
kubectl port-forward svc/my-health-caretaker 9091:9091
```

## License

This chart is licensed under the MIT License.
