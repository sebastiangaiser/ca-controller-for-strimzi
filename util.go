package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlRuntime "sigs.k8s.io/controller-runtime"
	"strconv"
)

const (
	defaultTargetSecretGeneration = "0"
)

func incrementString(input string) string {
	i, err := strconv.Atoi(input)
	if err != nil {
		log.Error(err, "failed to increment string")
	}
	i++
	s := strconv.Itoa(i)
	return s
}

func buildTargetSecret(exists bool, sourceSecret *corev1.Secret, data map[string]string, secretName, secretNamespace, tlsSecretHash, certOrKeyLabel, clusterName string) *corev1.Secret {
	var generation string
	if !exists {
		generation = defaultTargetSecretGeneration
	} else {
		oldGeneration := sourceSecret.GetAnnotations()[certOrKeyLabel]
		generation = incrementString(oldGeneration)
		ctrlRuntime.Log.Info(fmt.Sprintf("Secret %s was updated from %s to %s", sourceSecret.Name, oldGeneration, generation))
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: secretNamespace,
			Labels: map[string]string{
				managedByLabelKey:   managedByLabelValue,
				hashLabelKey:        tlsSecretHash,
				strimziClusterLabel: clusterName,
				strimziKindLabel:    strimziKindValue,
			},
			Annotations: map[string]string{
				certOrKeyLabel: generation,
			},
		},
		StringData: data,
		Type:       corev1.SecretTypeOpaque,
	}
	return secret
}

func truncateString(s string, max int) string {
	return s[:max]
}

func buildHistoricalSecret(sourceSecret *corev1.Secret, certOrKeyLabel string) *corev1.Secret {
	generation := sourceSecret.GetAnnotations()[certOrKeyLabel]
	historicalName := fmt.Sprintf("%s-generation-%s", sourceSecret.Name, generation)
	ctrlRuntime.Log.Info(fmt.Sprintf("Creating historical secret %s from %s", historicalName, sourceSecret.Name))

	// Copy labels from source secret, excluding the hash label since this is a historical snapshot
	labels := make(map[string]string)
	for k, v := range sourceSecret.Labels {
		labels[k] = v
	}
	// Mark this as a historical secret
	labels["sebastian.gaiser.bayern/historical"] = "true"
	labels["sebastian.gaiser.bayern/historical-generation"] = generation

	// Copy annotations from source secret
	annotations := make(map[string]string)
	for k, v := range sourceSecret.Annotations {
		annotations[k] = v
	}

	// Copy data from source secret
	data := make(map[string][]byte)
	for k, v := range sourceSecret.Data {
		data[k] = v
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        historicalName,
			Namespace:   sourceSecret.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Data: data,
		Type: sourceSecret.Type,
	}
	return secret
}
