package cmd

import (
	"flag"
	"os"
	"time"

	konfigmanagerv1 "github.com/flanksource/konfig-manager/api/v1"
	"github.com/flanksource/konfig-manager/controllers"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var Operator = &cobra.Command{
	Use:   "operator",
	Short: "Start the kubernetes operator",
	Run:   run,
}

var enableLeaderElection, dev bool
var metricsAddr, probeAddr string
var syncPeriod time.Duration

func init() {
	Operator.Flags().StringVar(&metricsAddr, "metrics-bind-address", ":8081", "The address the metric endpoint binds to.")
	Operator.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":8082", "The address the probe endpoint binds to.")
	Operator.Flags().BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	Operator.Flags().BoolVar(&dev, "dev", false, "run operator in development mode")
	Operator.Flags().DurationVar(&syncPeriod, "sync-period", 10*time.Minute, "The time duration to run a full reconcile")
}

func run(cmd *cobra.Command, args []string) {
	scheme := runtime.NewScheme()
	setupLog := ctrl.Log.WithName("setup")

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(konfigmanagerv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme

	opts := zap.Options{
		Development: dev,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		SyncPeriod:             &syncPeriod,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "b4532c9b.flanksource.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.HierarchyConfigReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("canary"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Konfig")
		os.Exit(1)
	}
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
