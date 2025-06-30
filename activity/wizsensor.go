package activity

import (
	"context"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

type WizSensorParams struct {
	KubeconfigPath          string `json:"kubeconfigPath"`
	KubeconfigContext       string `json:"kubeconfigContext"`
	ImagePullSecretUsername string `json:"imagePullSecretUsername"`
	ImagePullSecretPassword string `json:"imagePullSecretPassword"`
	WizApiTokenClientId     string `json:"wizApiTokenClientId"`
	WizApiTokenClientToken  string `json:"wizApiTokenClientToken"`
}

func InstallWizSensorActivity(ctx context.Context, params WizSensorParams) error {
	namespace := "wiz"
	releaseName := "wiz-sensor"
	chartName := "wiz-sensor"
	repoUrl := "https://charts.wiz.io"

	//check for kubeconfig
	if _, err := os.Stat(params.KubeconfigPath); os.IsNotExist(err) {
		return fmt.Errorf("kubeconfig file does not exist at path: %s", params.KubeconfigPath)
	}

	// create the wiz namespace
	if err := ensureNamespace(params.KubeconfigPath, namespace); err != nil {
		return fmt.Errorf("failed to check or create namespace %s: %w", namespace, err)
	}

	//helm cli settings
	settings := cli.New()
	settings.KubeConfig = params.KubeconfigPath
	settings.KubeContext = params.KubeconfigContext
	settings.SetNamespace(namespace)

	// set up Helm
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secrets", nil); err != nil {
		return fmt.Errorf("failed to initialize Helm action configuration: %w", err)
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = releaseName
	install.Namespace = namespace
	install.ChartPathOptions.RepoURL = repoUrl

	// Load the chart
	chartPath, err := install.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return fmt.Errorf("failed to locate chart %s: %w", chartName, err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart %s: %w", chartName, err)
	}

	// Set the values for the chart
	values := map[string]interface{}{
		"imagePullSecret": map[string]interface{}{
			"username": params.ImagePullSecretUsername,
			"password": params.ImagePullSecretPassword,
		},
		"wizApiToken": map[string]interface{}{
			"clientId":    params.WizApiTokenClientId,
			"clientToken": params.WizApiTokenClientToken,
		},
	}

	// run the Helm install
	if _, err := install.RunWithContext(ctx, chart, values); err != nil {
		return fmt.Errorf("failed to install chart %s: %w", chartName, err)
	}

	// Ensure the pod is running
	if err := ensureRunningPod(params.KubeconfigPath, namespace, "wiz-sensor"); err != nil {
		return fmt.Errorf("failed to ensure pod is running: %w", err)
	}

	return nil
}

func ensureNamespace(kubeconfigPath, namespace string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	_, err = client.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err == nil {
		return nil // already exists
	}

	_, err = client.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: namespace},
	}, metav1.CreateOptions{})
	return err
}

func ensureRunningPod(kubeconfigPath, namespace, podName string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	timeout := 5 * time.Minute
	interval := 5 * time.Second
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get pod %s in namespace %s: %w", podName, namespace, err)
		}

		if pod.Status.Phase == corev1.PodRunning {
			return nil
		}

		time.Sleep(interval)
	}
	return fmt.Errorf("pod %s in namespace %s is not running", podName, namespace)
}
