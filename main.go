package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"os"
	ctrlRuntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	// since we invoke tests with -ginkgo.junit-report we need to import ginkgo.
	_ "github.com/onsi/ginkgo/v2"
)

var (
	globalName                             = "ca-controller-for-strimzi"
	scheme                                 = runtime.NewScheme()
	log                                    = ctrlRuntime.Log.WithName(globalName)
	reconcileAnnotationKey                 = "sebastian.gaiser.bayern/tls-strimzi-ca"
	reconcileAnnotationValue               = "reconcile"
	managedByLabelKey                      = "sebastian.gaiser.bayern/managed-by"
	hashLabelKey                           = "sebastian.gaiser.bayern/hash"
	managedByLabelValue                    = "ca-controller-for-strimzi"
	targetClusterNameKey                   = "sebastian.gaiser.bayern/target-cluster-name"
	targetClusterNameValue                 = ""
	targetSecretAnnotationNameKey          = "sebastian.gaiser.bayern/target-secret-name"
	targetSecretAnnotationKeyNameKey       = "sebastian.gaiser.bayern/target-secret-key-name"
	strimziClusterLabel                    = "strimzi.io/cluster"
	strimziKindLabel                       = "strimzi.io/kind"
	strimziKindValue                       = "Kafka"
	strimziCaCertGeneration                = "strimzi.io/ca-cert-generation"
	strimziCaKeyGeneration                 = "strimzi.io/ca-key-generation"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	logOpts := zap.Options{
		Development: os.Getenv("DEBUG") != "",
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	logOpts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrlRuntime.SetLogger(zap.New(zap.UseFlagOptions(&logOpts)))
	ctrlRuntime.Log.Info("Logger initialized")

	healthProbeBindAddress := os.Getenv("HEALTH_BIND_ADDR")
	if healthProbeBindAddress == "" {
		healthProbeBindAddress = ":8081"
	}
	controllerNamespace := os.Getenv("CONTROLLER_NAMESPACE")
	if controllerNamespace == "" {
		controllerNamespace = "kafka"
	}

	options := ctrlRuntime.Options{Scheme: scheme, Cache: cache.Options{DefaultNamespaces: map[string]cache.Config{controllerNamespace: {}}}}
	options.HealthProbeBindAddress = healthProbeBindAddress

	mgr, err := ctrlRuntime.NewManager(ctrlRuntime.GetConfigOrDie(), options)
	if err != nil {
		ctrlRuntime.Log.Error(err, "could not create manager")
		os.Exit(1)
	}

	if err = (&Reconciler{
		client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		ctrlRuntime.Log.Error(err, "unable to create controller", "controller", globalName)
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		ctrlRuntime.Log.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		ctrlRuntime.Log.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// Setup metrics serving
	metricsBindAddress := os.Getenv("METRICS_BIND_ADDR")
	if metricsBindAddress == "" {
		metricsBindAddress = ":9000"
	}
	metricsPath := os.Getenv("METRICS_PATH")
	if metricsPath == "" {
		metricsPath = "/metrics"
	}
	http.Handle(metricsPath, promhttp.Handler())

	go func() {
		if err := http.ListenAndServe(metricsBindAddress, nil); err != nil {
			ctrlRuntime.Log.Error(err, "failed running metrics server.")
			os.Exit(1)
		}
	}()

	if err := mgr.Start(ctrlRuntime.SetupSignalHandler()); err != nil {
		ctrlRuntime.Log.Error(err, "could not start manager")
		os.Exit(1)
	}
}

type Reconciler struct {
	client client.Client
}

func (r *Reconciler) SetupWithManager(mgr ctrlRuntime.Manager) error {
	return ctrlRuntime.
		NewControllerManagedBy(mgr).
		For(&corev1.Secret{}, builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		WithEventFilter(predicate.NewPredicateFuncs(func(obj client.Object) bool {
			annotations := obj.GetAnnotations()
			if annotations != nil {
				return annotations[reconcileAnnotationKey] == reconcileAnnotationValue
			}
			return false
		})).
		Complete(r)
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrlRuntime.Request) (ctrlRuntime.Result, error) {
	// verify secret should be reconciled
	tlsSecret := &corev1.Secret{}
	err := r.client.Get(ctx, req.NamespacedName, tlsSecret)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to lookup secret: %w", err)
	}
	tlsSecretAnnotations := tlsSecret.GetAnnotations()
	if tlsSecret.Type != corev1.SecretTypeTLS {
		ctrlRuntime.Log.Info(fmt.Sprintf("Skipping secret with name: %s because it is not of secret type: %s", tlsSecret.Name, string(corev1.SecretTypeTLS)))
		return ctrlRuntime.Result{}, nil
	}
	if value, exists := tlsSecretAnnotations[reconcileAnnotationKey]; !exists {
		ctrlRuntime.Log.Info(fmt.Sprintf("Secret with name '%s' should be reconciled but is missing annotation '%s'", tlsSecret.Name, reconcileAnnotationKey))
		return ctrlRuntime.Result{}, nil
	} else {
		if value != reconcileAnnotationValue {
			ctrlRuntime.Log.Info(fmt.Sprintf("Found secret with name '%s' but value '%s' is not matching required reconcile annotation '%s'", tlsSecret.Name, value, reconcileAnnotationKey))
			return ctrlRuntime.Result{}, nil
		} else {
			ctrlRuntime.Log.Info(fmt.Sprintf("Found secret with name '%s' containing '%s' annotation", tlsSecret.Name, value))
		}
	}

	// verify targetSecretAnnotationNameKey is set
	if value, exists := tlsSecretAnnotations[targetSecretAnnotationNameKey]; !exists {
		ctrlRuntime.Log.Info(fmt.Sprintf("Secret with name '%s' should be reconciled but is missing annotation '%s'", tlsSecret.Name, targetSecretAnnotationNameKey))
		return ctrlRuntime.Result{}, nil
	} else {
		if value == "" {
			ctrlRuntime.Log.Info("Secret with name: " + tlsSecret.Name + " should be reconciled but the target secret name annotation value is empty...")
			return ctrlRuntime.Result{}, nil
		} else {
			ctrlRuntime.Log.Info("Secret with name: " + tlsSecret.Name + " will get reconciled with target secret name: " + tlsSecretAnnotations[targetSecretAnnotationNameKey])
		}
	}

	// verify targetSecretAnnotationKeyNameKey is set
	if value, exists := tlsSecretAnnotations[targetSecretAnnotationKeyNameKey]; !exists {
		ctrlRuntime.Log.Info("Secret with name: " + tlsSecret.Name + " should be reconciled but is missing annotation: " + targetSecretAnnotationKeyNameKey)
		return ctrlRuntime.Result{}, nil
	} else {
		if value == "" {
			ctrlRuntime.Log.Info("Secret with name: " + tlsSecret.Name + " should be reconciled but the target secret name key annotation value is empty...")
			return ctrlRuntime.Result{}, nil
		} else {
			ctrlRuntime.Log.Info("Secret with name: " + tlsSecret.Name + " will get reconciled with target secret name: " + tlsSecretAnnotations[targetSecretAnnotationKeyNameKey])
		}
	}

	// check if target secrets are existing
	targetSecret := &corev1.Secret{}
	targetSecretExists := false
	targetSecretNeedsUpdate := false
	targetSecretKey := &corev1.Secret{}
	targetSecretKeyExists := false
	targetSecretKeyNeedsUpdate := false

	// process tlsSecretHash for comparison with target secrets label
	tlsSecretHashRaw := sha256.Sum256([]byte(string(tlsSecret.Data["ca.crt"]) + string(tlsSecret.Data["tls.crt"]) + string(tlsSecret.Data["tls.key"])))
	tlsSecretHash := truncateString(hex.EncodeToString(tlsSecretHashRaw[:]), 63)

	// TODO export to function
	err = r.client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: tlsSecretAnnotations[targetSecretAnnotationNameKey]}, targetSecret)
	if err == nil {
		ctrlRuntime.Log.Info(fmt.Sprintf("Secret %s already exists", tlsSecretAnnotations[targetSecretAnnotationNameKey]))
		targetSecretExists = true

		// check if target secret need an update
		// but first check if the hash label exists...
		targetSecretHash, ok := targetSecret.Labels[hashLabelKey]
		if !ok {
			return ctrlRuntime.Result{}, errors.New(fmt.Sprintf("label '%s' not found for target secret '%s'", targetSecret, hashLabelKey))
		}
		// now check if it needs an update
		if tlsSecretHash != targetSecretHash {
			targetSecretNeedsUpdate = true
		}
	}

	err = r.client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: tlsSecretAnnotations[targetSecretAnnotationKeyNameKey]}, targetSecretKey)
	if err == nil {
		ctrlRuntime.Log.Info(fmt.Sprintf("Secret %s already exists", tlsSecretAnnotations[targetSecretAnnotationKeyNameKey]))
		targetSecretKeyExists = true

		// check if target secret need an update
		// but first check if the hash label exists...
		targetSecretKeyHash, ok := targetSecretKey.Labels[hashLabelKey]
		if !ok {
			return ctrlRuntime.Result{}, errors.New(fmt.Sprintf("label '%s' not found for target secret '%s'", targetSecretKey, hashLabelKey))
		}
		// now check if it needs an update
		if tlsSecretHash != targetSecretKeyHash {
			targetSecretKeyNeedsUpdate = true
		}
	}

	if targetSecretExists && !targetSecretNeedsUpdate && targetSecretKeyExists && !targetSecretKeyNeedsUpdate {
		ctrlRuntime.Log.Info("All target secrets are up-to-date")
		return ctrlRuntime.Result{}, nil
	}

	caCrt := string(tlsSecret.Data["ca.crt"])
	tlsCrt := string(tlsSecret.Data["tls.crt"])
	tlsKey := string(tlsSecret.Data["tls.key"])
	combinedCrt := tlsCrt + caCrt

	clusterName := tlsSecretAnnotations[targetClusterNameKey]

	targetSecretToApply := buildTargetSecret(targetSecretExists, targetSecret, map[string]string{"ca.crt": combinedCrt}, tlsSecretAnnotations[targetSecretAnnotationNameKey], req.Namespace, tlsSecretHash, strimziCaCertGeneration, clusterName)
	targetSecretKeyToApply := buildTargetSecret(targetSecretKeyExists, targetSecretKey, map[string]string{"ca.key": tlsKey}, tlsSecretAnnotations[targetSecretAnnotationKeyNameKey], req.Namespace, tlsSecretHash, strimziCaKeyGeneration, clusterName)

	if !targetSecretExists {
		ctrlRuntime.Log.Info(fmt.Sprintf("Creating target secret %s/%s", targetSecretKeyToApply.Name, targetSecretKeyToApply.Namespace))
		err := r.client.Create(ctx, targetSecretToApply)
		if err != nil {
			return ctrlRuntime.Result{}, err
		} else {
			ctrlRuntime.Log.Info(fmt.Sprintf("Target secret %s successfully reconciled", tlsSecretAnnotations[targetSecretAnnotationNameKey]))
		}
	} else {
		ctrlRuntime.Log.Info(fmt.Sprintf("Updating target secret %s/%s", targetSecretKeyToApply.Name, targetSecretKeyToApply.Namespace))
		err := r.client.Update(ctx, targetSecretToApply)
		if err != nil {
			return ctrlRuntime.Result{}, err
		} else {
			ctrlRuntime.Log.Info(fmt.Sprintf("Target secret %s successfully reconciled", tlsSecretAnnotations[targetSecretAnnotationNameKey]))
		}
	}
	if !targetSecretKeyExists {
		ctrlRuntime.Log.Info(fmt.Sprintf("Creating target secret %s/%s", targetSecretKeyToApply.Name, targetSecretKeyToApply.Namespace))
		err := r.client.Create(ctx, targetSecretKeyToApply)
		if err != nil {
			return ctrlRuntime.Result{}, err
		} else {
			ctrlRuntime.Log.Info(fmt.Sprintf("Target secret %s successfully reconciled", tlsSecretAnnotations[targetSecretAnnotationKeyNameKey]))
		}
	} else {
		ctrlRuntime.Log.Info(fmt.Sprintf("Updating target secret %s/%s", targetSecretKeyToApply.Name, targetSecretKeyToApply.Namespace))
		err := r.client.Update(ctx, targetSecretKeyToApply)
		if err != nil {
			return ctrlRuntime.Result{}, err
		} else {
			ctrlRuntime.Log.Info(fmt.Sprintf("Target secret %s successfully reconciled", tlsSecretAnnotations[targetSecretAnnotationKeyNameKey]))
		}
	}

	return ctrlRuntime.Result{}, nil
}
