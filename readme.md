# go app with helm charts

This repo consists of a simple go web application.
It also has a helm chart that deploys it to a kubernetes cluster along with a config map.
The objective of the repository is to understand how to use configmaps along with kubernetes apps

## Running the helm chart

### Pre-requisites

1. Docker
2. minikube
3. kubernetes
4. helm3
5. go

### Running the code in kubernetes

1. Start minikube
   ```shell
    minikube start
   ```
2. Run the following command to be able to build the docker image in the minikube docker env
   ```shell
    eval $(minikube docker-env)
   ```
3. Build the docker image in the same terminal where the above command was executed
   ```shell
    docker build -t gok8 .
   ```
4. Install the application using helm
   The syntax for the following command is `helm install <release-name> <chart-location>`
   ```shell
    helm install v1 ./go-helm
   ```

### Verifying the installed code

1. Run `kubectl get pods` to check the running pods. The output should be something like
   ```shell
    kubectl get pods
    NAME                          READY   STATUS    RESTARTS   AGE
    go-web-app-54b544476d-mgktm   1/1     Running   0          5s
    go-web-app-54b544476d-vcvjj   1/1     Running   0          5s
   ```
   
2. Since the app exposes 2 rest endpoints namely `/` and `/config` one can access it by running
   ```shell
    minikube service go-web-service
   ```
   where `go-web-service` is the name of the kubernetes service we deployed
   The above command establishes a ssh tunnel with our kubernetes app and opens a browser
   
   The output of the above command is something like
   ```shell
    minikube service go-web-service
    |-----------|----------------|-------------|---------------------------|
    | NAMESPACE |      NAME      | TARGET PORT |            URL            |
    |-----------|----------------|-------------|---------------------------|
    | default   | go-web-service | http/80     | http://192.168.49.2:32029 |
    |-----------|----------------|-------------|---------------------------|
    🏃  Starting tunnel for service go-web-service.
    |-----------|----------------|-------------|------------------------|
    | NAMESPACE |      NAME      | TARGET PORT |          URL           |
    |-----------|----------------|-------------|------------------------|
    | default   | go-web-service |             | http://127.0.0.1:57922 |
    |-----------|----------------|-------------|------------------------|
    🎉  Opening service default/go-web-service in default browser...
    ❗  Because you are using a Docker driver on darwin, the terminal needs to be open to run it.
   ```
    
3. One can now navigate to `http://127.0.0.1:57922/config` to view the config map that was loaded
   Note that the port number will be different every time
   
### Changing the config map

The `deployment.yaml` has an annotation that is tied to the `sha256sum` of the `client_config.json` file.
This is so that the pods restart automatically if there is any change to the `client_config.json` file

```yaml
template:
    metadata:
      labels:
        name: go-web-app
      annotations:
        checksum/config: {{ .Files.Get "client_config.json" | sha256sum }}
```

So now modify the `client_config.json` and upgrade the release by running

```shell
helm upgrade v1 ./go-helm
```
Immediately run `kubectl get pods` and you should see that the old pods are being terminated and new pods get created due to change in config map
```shell
kubectl get pods
NAME                          READY   STATUS        RESTARTS   AGE
go-web-app-54b544476d-mgktm   0/1     Terminating   0          12m
go-web-app-54b544476d-vcvjj   0/1     Terminating   0          12m
go-web-app-696994cfc4-76lm5   1/1     Running       0          6s
go-web-app-696994cfc4-pzp97   1/1     Running       0          5s
```

Verify the updated configs by again navigating to the `/config` end point of the app

### Clean up

To clean up the installation run the following helm command
```shell
helm delete <release-name>
```

In this case it is
```shell
helm delete v1
```
