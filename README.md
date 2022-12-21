# EKS Shared Usage Calculator

EKS Shared Usage Calculator extracts and calculates share usage cost metrics from kubecost (OpenCost) cost analyzer deployments on EKS clusters to Datadog.

This tool is a custom API scraper that interacts with Kubecost Cost Analyzer deployments on Kubernetes clusters. This provides the capability to scrape and calculate cost data metrics and usage percentages to be sourced for the Datadog Cost-per-service Dashboard to assist with accurate cost attribution per service. It is an improvement over [kubectl-cost](https://github.com/kubecost/kubectl-cost) plugin for a specific use case.

# Dependencies
- This tool requires that Kubecost Cost Analyzer has already been succesfully deployed to a Kubernetes cluster and has been given 24 hours to obtain and store cost data before polling with this tool.
- It is required that this tool be run within a Kubernetes cluster environment. The tool "could" be run standalone via kube-admin however port forwarding would be required on the cluster and API endpoint contants will need to be modified in this source code. Future iterations of this tool may include the functionality to pass in arguments to this tool.

## Usage
This tool is primarily designed to be deployed as a pod/container to a kubernetes cluster in order to interact with the local cluster DNS names of the `kubecost cost analyzer` service ([kubecost-cost-analyzer.kubecost.svc.cluster.local](kubecost-cost-analyzer.kubecost.svc.cluster.local)).

## kube-admin
This tool can be launched directly from `kube-admin` container when working on EKS clusters. just run the following from `kube-admin` tool in order to gernerate the shared cost tag data for an EKS cluster, you can find this script on the kube-admin container tool here: []()
```bash
/usr/local/bin/eks-shared-usage-calculator
```
