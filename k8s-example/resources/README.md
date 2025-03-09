## create namespace
```kubectl create namespace demo```

## apply resource quota to namespace
```kubectl apply --namespace demo -f k8s/resourcequota.yaml resourcequota "demo-resourcequota" created```


Now Kubernetes will block any API operations in the demo namespace that would exceed the quota. The example ResourceQuota limits the namespace to 100 Pods, so if there are 100 Pods already running and you try to start a new one, you will see an error message like this:
    Error from server (Forbidden): pods "demo" is forbidden: exceeded quota:
    demo-resourcequota, requested: pods=1, used: pods=100, limited: pods=100

Using ResourceQuotas is a good way to stop applications in one namespace from grabbing too many resources and starving those in other parts of the cluster.

## checking if resource quota is active
```kubectl get resourcequotas -n demo```


# Default Resource Requests and Limits
It’s not always easy to know what your container’s resource requirements are going to be in advance. You can set default requests and limits for all containers in a name‐ space using a LimitRange resource:


Any container in the namespace that doesn’t specify a resource limit or request will inherit the default value from the LimitRange. For example, a container with no cpu request specified will inherit the value of 200m from the LimitRange. Similarly, a con‐ tainer with no memory limit specified will inherit the value of 256Mi from the Limit‐ Range.

In theory, then, you could set the defaults in a LimitRange and not bother to specify requests or limits for individual containers. However, this isn’t good practice: it should be possible to look at a container spec and see what its requests and limits are, without having to know whether or not a LimitRange is in effect. Use the LimitRange only as a backstop to prevent problems with containers whose owners forgot to spec‐ ify requests and limits.


### Best Practice
> Use LimitRanges in each namespace to set default resource requests and limits for containers, but don’t rely on them; treat them as a backstop. Always specify explicit requests and limits in the container spec itself.


A good rule of thumb is that nodes should be big enough to run at least five of your typical Pods, keeping the proportion of stranded resources to around 10% or less. If the node can run 10 or more Pods, stranded resources will be below 5%.

The default limit in Kubernetes is 110 Pods per node. Although you can increase this limit by adjusting the --max-pods setting of the kubelet, this may not be possible with some managed services, and it’s a good idea to stick to the Kubernetes defaults unless there is a strong reason to change them.


The Pods-per-node limit means that you may not be able to take advantage of your cloud provider’s largest instance sizes. Instead, consider running a larger number of smaller nodes to get better utilization. For example, instead of 6 nodes with 8 vCPUs, run 12 nodes with 4 vCPUs.


### Best Practice
> Larger nodes tend to be more cost-effective, because less of their resources are consumed by system overhead. Size your nodes by looking at real-world  utilization figures for your cluster, aiming for between 10 and 100 Pods per node.


# Optimizing Storage
Best Practice
> Don’t use instance types with more storage than you need. Provi‐ sion the smallest, lowest-IOPS disk volumes you can, based on the throughput and space that you actually use.


# Using owner metadata

```yaml 
apiVersion: extensions/v1beta1
    kind: Deployment
    metadata:
      name: my-brilliant-app
      annotations:
        example.com/owner: "Customer Apps Team"
```

## Best Practice
Set owner annotations on all your resources, giving information about who to contact if there’s a problem with this resource, or if it seems abandoned and liable for termination.