# kube-pushgateway-test

Prometheus Push Gateway tested.

Testing how to collect metrics of cronjobs using Prometheus.

For test you need a Kubernetes Cluster or Minikube.

If you are using Minikube:
```sh
minikube -p pushgateway start
```

On the folder project (Skip this if you are using Docker for Mac with Kubernetes):
```sh
eval $(minikube -p pushgateway docker-env)
```

Build test Image inside the minikube cluster
```sh
docker build -t test-cronjob:1.0 .
```

Now deploy the Prometheus Push Gateway:
```sh
kubectl --context pushgateway apply -f pushgateway.k8s.yaml
```
replace the context with your context. (docker-desktop if you are using Docker for Mac)

And now deploy the CronJob
```sh
kubectl --context pushgateway apply -f cronjob.k8s.yaml
```

To see if all is working type:
```sh
kubectl --context pushgateway -n kube-pushgateway-test get all
```

That return something like:
```
NAME                                     READY   STATUS      RESTARTS   AGE
pod/kube-cronjob-test-1570811520-88hqb   0/1     Completed   0          2m29s
pod/kube-cronjob-test-1570811580-z6s7n   0/1     Completed   0          89s
pod/kube-cronjob-test-1570811640-kw2mw   0/1     Completed   0          29s
pod/pushgateway-56d8d86bc6-d9nx5         1/1     Running     1          19h

NAME                  TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
service/pushgateway   ClusterIP   10.99.182.29   <none>        9091/TCP   19h

NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/pushgateway   1/1     1            1           19h

NAME                                     DESIRED   CURRENT   READY   AGE
replicaset.apps/pushgateway-56d8d86bc6   1         1         1       19h

NAME                                     COMPLETIONS   DURATION   AGE
job.batch/kube-cronjob-test-1570811520   1/1           3s         2m29s
job.batch/kube-cronjob-test-1570811580   1/1           4s         89s
job.batch/kube-cronjob-test-1570811640   1/1           4s         29s

NAME                              SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
cronjob.batch/kube-cronjob-test   */1 * * * *   False     0        34s             143m
```

To test if the metrics are fine you can port-foward to the Push Gateway service:
```sh
kubectl --context pushgateway -n kube-pushgateway-test port-forward svc/pushgateway 9091:9091
```

```sh
curl http://localhost:9091/metrics
```

```
# HELP test_duration_seconds The duration of the last in seconds.
# TYPE test_duration_seconds gauge
test_duration_seconds{instance="",job="test-job"} 2.000900967
# HELP test_last_completion_timestamp_seconds The timestamp of the last completion, successful or not.
# TYPE test_last_completion_timestamp_seconds gauge
test_last_completion_timestamp_seconds{instance="",job="test-job"} 1.570811829237434e+09
# HELP test_last_success_timestamp_seconds The timestamp of the last successful completion.
# TYPE test_last_success_timestamp_seconds gauge
test_last_success_timestamp_seconds{instance="",job="test-job"} 1.5708118292374587e+09
```

Some alarms that you can set is the duration, difference between success & completion
and completion with current time.

Last step is to configure Prometheus Operator to scrap the metrics.
