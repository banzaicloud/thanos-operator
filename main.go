// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"os"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/controllers"
	thanosv1alpha1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	thanosconfig "github.com/banzaicloud/thanos-operator/controllers/config"
	// +kubebuilder:scaffold:imports
)

const DefaultLeaderElectionID = "banzaicloud-thanos-operator"

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = thanosv1alpha1.AddToScheme(scheme)
	_ = apiextensions.AddToScheme(scheme)
	_ = prometheus.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection, enablePromCRDWatches bool
	var leaderElectionId string
	var leaderElectionNamespace string
	var flagThanosImage string
	var flagThanosImageTag string

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&enablePromCRDWatches, "enable-prom-crd-watches", true, "Enable dynamic watches of Prometheus CRDs")
	flag.StringVar(&leaderElectionId, "leader-election-id", "", "The ID of the leader election")
	flag.StringVar(&leaderElectionNamespace, "leader-election-ns", "", "The NS  of the leader election")
	flag.StringVar(&flagThanosImage, "thanos-image", "", "Overrides the default Thanos image")
	flag.StringVar(&flagThanosImageTag, "thanos-image-tag", "", "Overrides the default Thanos image tag")
	flag.Parse()

	ctrl.SetLogger(utils.Log)

	if leaderElectionId == "" {
		leaderElectionId = DefaultLeaderElectionID
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      metricsAddr,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        leaderElectionId,
		LeaderElectionNamespace: leaderElectionNamespace,
		Port:                    9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	controllerConfig := thanosconfig.GetControllerConfig()
	controllerConfig.AddConfigItem("ThanosImage", flagThanosImage)
	controllerConfig.AddConfigItem("ThanosImageTag", flagThanosImageTag)

	var ThanosController, ObjectStoreController, ReceiverController controller.Controller

	if ObjectStoreController, err = (&controllers.ObjectStoreReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ObjectStore"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ObjectStore")
		os.Exit(1)
	}
	if ThanosController, err = (&controllers.ThanosReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Thanos"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Thanos")
		os.Exit(1)
	}
	if err = (&controllers.StoreEndpointReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("StoreEndpoint"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "StoreEndpoint")
		os.Exit(1)
	}
	if ReceiverController, err = (&controllers.ReceiverReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Receiver"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Receiver")
		os.Exit(1)
	}
	if err = (&controllers.ThanosEndpointReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ThanosEndpoint"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ThanosEndpoint")
		os.Exit(1)
	}
	if err = (&controllers.ThanosPeerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ThanosPeer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ThanosPeer")
		os.Exit(1)
	}

	if enablePromCRDWatches {
		if err = (&controllers.ServiceMonitorWatchReconciler{
			Log: ctrl.Log.WithName("controllers").WithName("ServiceMonitorWatch"),
			Controllers: map[string]controllers.ControllerWithSource{
				"receiver": {
					Controller: ReceiverController,
					Source:     &thanosv1alpha1.Receiver{},
				},
				"thanos": {
					Controller: ThanosController,
					Source:     &thanosv1alpha1.Thanos{},
				},
				"objectstore": {
					Controller: ObjectStoreController,
					Source:     &thanosv1alpha1.ObjectStore{},
				},
			},
			Client: mgr.GetClient(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ServiceMonitorWatch")
			os.Exit(1)
		}
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
