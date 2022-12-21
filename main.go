package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	appName string = "eks-shared-usage-calculator"
	version string = "1.0.0"

	kubecostServiceFQDN    = "kubecost-cost-analyzer.kubecost.svc.cluster.local"
	kubecostServicePort    = "9090"
	kubecostApiEndpoint    = "/model/allocation"
	kubecostDataWindow     = "1d"
	kubecostDataAggregate  = "deployment"
	kubecostDataAccumulate = true
)

var (
	KubecostApiUrl = fmt.Sprintf("http://%s:%s%s?window=%s&aggregate=%s&accumulate=%v",
		kubecostServiceFQDN, kubecostServicePort, kubecostApiEndpoint,
		kubecostDataWindow, kubecostDataAggregate, kubecostDataAccumulate)
)

// KubernetesPodId kubernetes deployment name
type KubernetesDeploymentId string

// KubernetesPod kubernetes pod
type KubernetesDeployment struct {
	Name      string  `json:"name"`
	TotalCost float32 `json:"totalCost"`
}

// kubernetesDeployment kubernetes deployment
type kubernetesDeployment struct {
	Code int                                               `json:"code"`
	Data []map[KubernetesDeploymentId]KubernetesDeployment `json:"data"`
}

func main() {
	log.Printf("Starting %s...\n", appName)

	// Get API response JSON data
	log.Printf("Scraping kubecost API endpoint: %v\n", KubecostApiUrl)
	kubecostRawData := getKubecostRawData()

	rawData := []byte(kubecostRawData)

	var deployments kubernetesDeployment
	err := json.Unmarshal(rawData, &deployments)
	if err != nil {
		log.Fatalln(err)
	}

	// From scraped kubecost API data, calculate percentages per deployment of TotalCost in range
	fmt.Printf("---\n")
	percentageCalculatedDeployments := calculateDeploymentSharePercentages(iterateThrough(deployments))
	for deploymentId, deploymentSharePercentage := range percentageCalculatedDeployments {
		fmt.Printf("%s = %.5f\n", deploymentId, deploymentSharePercentage/100)
	}

	fmt.Printf("---\n")
	log.Printf("DONE!\n")
	log.Printf("Scrape job completed.\n")
}

// Calculates the share percentages of a given map with String key and Float32 values; returns map[string]float32 :: deploymentSharePercentages
func calculateDeploymentSharePercentages(m map[string]float32) map[string]float32 {
	var total float32
	for _, v := range m {
		total += v
	}

	deploymentSharePercentages := make(map[string]float32)
	for k, v := range m {
		deploymentSharePercentages[k] = v / total * 100
	}
	return deploymentSharePercentages
}

// Iterates through deployment JSON objects and returns a map of kubernets deployments based on scraped JSON data; returns map[string]float32 :: deployments
func iterateThrough(deployment kubernetesDeployment) map[string]float32 {
	deployments := map[string]float32{}
	for _, deploymentList := range deployment.Data {
		for deploymentId, deployment := range deploymentList {
			deployments[string(deploymentId)] = float32(deployment.TotalCost)
		}
	}
	return deployments
}

// Obtain raw http response body content from kubecost API; returns String :: kubecostResponse
func getKubecostRawData() string {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", KubecostApiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "kubecost-datadog-exporter")
	req.Header.Add("Contet-Type", "text/plain; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	kubecostApiResponse, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(kubecostApiResponse)
}
