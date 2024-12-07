# argocd-applicationset-namespaces-generator-plugin

Namespaces Generator that discovers namespaces in a target cluster.

It can be used as ArgoCD ApplicationSet plugin https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/Generators-Plugin/.

It can discover existing namespaces in the cluster to produce an app per each namespace.

## Assumptions and prerequisites

- You are using JWT authentication to your clusters (i.e. Downward API tokens mounted to pods)
- If using external clusters, you must populate cluster annotation with its Certificate Authority

## Usage

1. Deploy the argocd-applicationset-namespaces-generator-plugin from `testdata/manifest.yaml`.

2. Deploy the ApplicationSet YAML from `testdata/appset.yaml`.

Here's an example to use together with clusters generator via matrix generator:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: project-hsiaoairplane-namespaces-generator
spec:
  goTemplate: true
  goTemplateOptions: ["missingkey=error"]
  generators:
  - matrix:
      generators:
      - clusters: {}
      - plugin:
          configMapRef:
            name: argocd-applicationset-namespaces-generator-plugin
          input:
            parameters:
              clusterName: "{{ .name }}"
              clusterEndpoint: "{{ .server }}"
              # Optional, if not set means all namespaces
              labelSelector:
                project: hsiaoairplane
          # OPTIONAL: Checks for changes every 30 seconds
          requeueAfterSeconds: 30
  template:
    metadata:
      name: 'project-hsiaoairplane-{{ .namespace }}-namespaces-generator'
      namespace: '{{ .namespace }}'
    spec:
      project: "default"
      source:
        repoURL: https://github.com/argoproj/argocd-example-apps
        targetRevision: master
        path: guestbook
      destination:
        server: '{{ .server }}'
        namespace: '{{ .namespace }}'
      syncPolicy:
        automated:
          prune: true
        syncOptions:
        - CreateNamespace=false
        - FailOnSharedResource=true
        - PruneLast=true
        - PrunePropagationPolicy=foreground
```

# Testing

```bash
go run ./... -v=4 --log-format=json server --local
curl -X POST -H "Content-Type: application/json" -d @testdata/request.json http://localhost:8080/api/v1/getparams.execute
```
