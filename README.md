# Longshore

`longshore` is a small tool for providing a healthcheck against the Docker daemon in a Kubernetes cluster. It performs this by exposing an HTTP endpoint `/readyz` which returns the JSON-ified equivalent of the command-line execution of `docker info`.

## Endpoints

- `/`, `/readyz`, `/healthz`  
  The JSON-ified output equivalent of `docker info`. If the daemon is non-responsive, the HTTP response code will be `500`. Otherwise, it will be `200`.

- `/livez`  
  Returns `200` whenever the `longshore` service is running.

## Suggested deployment

Deploy `longshore` as a DaemonSet, which should run on all nodes which have a Docker daemon process. For example:

```yaml
---
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: longshore
  namespace: kube-system
spec:
  selector:
    matchLabels:
      group: monitoring
      k8s-app: longshore
      kubernetes.io/cluster-service: "true"
  template:
    metadata:
      labels:
        group: monitoring
        k8s-app: longshore
        kubernetes.io/cluster-service: "true"
    spec:
      containers:
        - name: longshore
          image: docker.io/moinknet/longshore:0.1.0
          ports:
            - containerPort: 8080
          livenessProbe:
            httpGet:
              path: /livez
              port: 8080
            initialDelaySeconds: 3
            periodSeconds: 3
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 30
          resources:
            limits:
              cpu: 100m
              memory: 32Mi
          securityContext:
            runAsUser: 0
          volumeMounts:
            - name: docker-socket
              readOnly: true
              mountPath: /var/run/docker.sock
      volumes:
        - name: docker-socket
          hostPath:
            path: /var/run/docker.sock
```
