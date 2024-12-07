# argocd-applicationset-namespaces-generator-plugin

Namespaces Generator that discovers namespaces in a target cluster.

It can be used as ArgoCD ApplicationSet plugin https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/Generators-Plugin/.

It can discover existing namespaces in the cluster to produce an app per each namespace.

# Assumptions and prerequisites

- You are using JWT authentication to your clusters (i.e. Downward API tokens mounted to pods)
- If using external clusters, you must populate cluster annotation with its Certificate Authority

# Usage

## Local

```bash
go run ./... -v=4 --log-format=json server --local
curl -X POST -H "Content-Type: application/json" -d @testdata/request.json http://localhost:8080/api/v1/getparams.execute
```

## In Cluster

1. Deploy the argocd-applicationset-namespaces-generator-plugin Deployment.

   ```console
   kubectl apply -f testdata/manifest.yaml
   ```

2. Deploy the ApplicationSet YAMLs.

   ```console
   kubectl apply -f testdata/project-airplanehsiao-appset.yaml
   kubectl apply -f testdata/project-hsiaoairplane-appset.yaml
   ```

3. Create ArgoCD ApplicationProjects.

   ```console
   argocd proj create project-hsiaoairplane
   argocd proj create project-airplanehsiao
   ```

4. Create the project "hsiaoairplane" namespaces.

   ```console
   kubectl create namespace foo
   kubectl label namespace foo project=hsiaoairplane --overwrite=true
   kubectl create namespace foobar
   kubectl label namespace foobar project=hsiaoairplane --overwrite=true
   ```

5. Create the project "airplanehsiao" namespaces.

   ```console
   kubectl create namespace bar
   kubectl label namespace bar project=airplanehsiao --overwrite=true
   kubectl create namespace barfoo
   kubectl label namespace barfoo project=airplanehsiao --overwrite=true
   ```

6. Access the ArgoCD GUI.

   ```console
   argocd admin dashboard -n argocd
   ```

7. Open the browser http://localhost:8080/
