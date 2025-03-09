# Demo

## Kubernetes manifest File

Run kubectl apply -f ${deplument yaml file}

After a few seconds, a demo Pod should be running:

kubectl get pods --selector app=demo


## Service resources

Suppose you want to make a network connection to a Pod (such as our example application)

You can see that it looks somewhat similar to the Deployment resource we showed earlier. However, the kind is Service, instead of Deployment, and the spec just includes a list of ports, plus a selector and a type.

If you zoom in a little, you can see that the Service is forwarding its port 9999 to the Pod’s port 8888:

The selector is the part that tells the Service how to route requests to particular Pods. Requests will be forwarded to any Pods matching the specified set of labels; in this case, just app: demo


## Go ahead and apply the manifest now, to create the Service:

kubectl apply -f k8s/service.yaml

kubectl port-forward service/demo 9999:8888

As before, kubectl port-forward will connect the demo service to a port on your local machine, so that you can connect to http://localhost:9999/ with your web browser.

Once you’re satisfied that everything is working correctly, run the following com‐ mand to clean up before moving on to the next section:

kubectl delete -f k8s/


To see comprehensive information about an individual Pod (or any other resource), use kubectl describe:

kubectl describe pod/demo-dev-6c96484c48-69vss