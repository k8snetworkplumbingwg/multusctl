package commands

import (
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
	"k8s.io/klog/v2"
)

var logTo string
var verbose int
var maxWidth int
var colorize string
var kubeconfigPath string
var masterUrl string
var context string

func init() {
	var defaultKubeconfigPath string
	if u, err := user.Current(); err == nil {
		defaultKubeconfigPath = filepath.Join(u.HomeDir, ".kube", "config")
	}

	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().IntVarP(&maxWidth, "width", "j", 0, "maximum output width (0 to use terminal width, -1 for no maximum)")
	rootCommand.PersistentFlags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	rootCommand.PersistentFlags().StringVarP(&masterUrl, "master", "m", "", "address of the Kubernetes API server")
	rootCommand.PersistentFlags().StringVarP(&kubeconfigPath, "kubeconfig", "k", defaultKubeconfigPath, "path to Kubernetes configuration")
	rootCommand.PersistentFlags().StringVarP(&context, "context", "x", "", "name of context in Kubernetes configuration")
}

var rootCommand = &cobra.Command{
	Use:   toolName,
	Short: "Control Multus CNI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		util.InitializeColorization(colorize)
		commonlog.Initialize(verbose, logTo)
		if writer := commonlog.GetWriter(); writer != nil {
			klog.SetOutput(writer)
		}
	},
}

func Execute() {
	err := rootCommand.Execute()
	util.FailOnError(err)
}
