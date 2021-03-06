package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var runCmd = &cobra.Command{
	Use:   "run [cronjob name]",
	Short: "Run a cronjob immediately",
	Long:  "Run a cronjob immediately",
	Run:   run,
}

func run(_ *cobra.Command, args []string) {

	cronjobName := args[0]

	clientset := mustGetClientset()

	cronjob := getCronjob(namespace, cronjobName)

	createJob(clientset, cronjobName, cronjob)
}

func createJob(clientset *kubernetes.Clientset, cronjobName string, cronjob *v1beta1.CronJob) {
	jobsClient := clientset.BatchV1().Jobs(namespace)

	suffix := "-" + strconv.Itoa(int(time.Now().Unix()))

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cronjobName + suffix,
			Namespace: namespace,
		},
		Spec: cronjob.Spec.JobTemplate.Spec,
	}

	result, err := jobsClient.Create(job)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created job %q.\n", result.GetObjectMeta().GetName())
}
