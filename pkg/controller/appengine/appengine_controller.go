package appengine

import (
	"context"
	"fmt"
	"time"

	//"github.com/pkg/errors"
	knapv1alpha1 "github.com/bluebosh/knap/pkg/apis/knap/v1alpha1"
	pipelinev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	servingv1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_appengine")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Appengine Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAppengine{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("appengine-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Appengine
	err = c.Watch(&source.Kind{Type: &knapv1alpha1.Appengine{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Appengine
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &knapv1alpha1.Appengine{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAppengine implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAppengine{}

// ReconcileAppengine reconciles a Appengine object
type ReconcileAppengine struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Appengine object and makes changes based on the state read
// and what is in the Appengine.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example+
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAppengine) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Appengine")

	// Fetch the Appengine instance
	app := &knapv1alpha1.Appengine{}
	err := r.client.Get(context.TODO(), request.NamespacedName, app)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	trygitresource := &pipelinev1alpha1.PipelineResource{}
	err = r.client.Get(context.TODO(), client.ObjectKey{Namespace: app.Namespace, Name: app.Spec.AppName + "-git"}, trygitresource)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			pipelinerun, err := runProcess(r, app)

			if err != nil {
				reqLogger.Error(err, "Failed to create new app", "app name", app.Spec.AppName)
				return reconcile.Result{}, err
			}
			reqLogger.Info("Finish build new app", "pipelinerun name", pipelinerun.Name, "pipelinerun status", pipelinerun.Status.Results)
			return reconcile.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
		} else {
			reqLogger.Error(err, "Failed to get git resource", "git resource", app.Spec.AppName+"-git")
			return reconcile.Result{}, err
		}
	} else {
		if app.Spec.GitRepo != trygitresource.Spec.Params[0].Value || app.Spec.GitRevision != trygitresource.Spec.Params[1].Value {
			pipelinerun, err := runProcess(r, app)

			if err != nil {
				reqLogger.Error(err, "Failed to update the app", "app name", app.Spec.AppName)
				return reconcile.Result{}, err
			}
			reqLogger.Info("Finish re-build the app, waiting for result", "pipelinerun name", pipelinerun.Name, "pipelinerun status", pipelinerun.Status.Results)
			return reconcile.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
		} else {
			pipelinerunresult := &pipelinev1alpha1.PipelineRun{}
			err = r.client.Get(context.TODO(), client.ObjectKey{Namespace: app.Namespace, Name: app.Spec.AppName + "-pr-" + fmt.Sprintf("%d", app.ObjectMeta.Generation)}, pipelinerunresult)
			if err != nil {
				if apiErrors.IsNotFound(err) {
					pipelinerun, err := runProcess(r, app)

					if err != nil {
						reqLogger.Error(err, "Failed to re-run the pipeline run", "pipelinerun name", app.Spec.AppName+"-pr"+fmt.Sprintf("%d", app.ObjectMeta.Generation))
						return reconcile.Result{}, err
					}
					reqLogger.Info("Finish to re-run the pipelinerun, waiting for result", "pipelinerun name", pipelinerun.Name, "pipelinerun status", pipelinerun.Status.Results)
					return reconcile.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
				} else {
					reqLogger.Error(err, "Failed to get the pipelinerun", "pipelinerun name", app.Spec.AppName+"-pr-"+fmt.Sprintf("%d", app.ObjectMeta.Generation))
					return reconcile.Result{}, err
				}
			} else {
				if string(pipelinerunresult.Status.Conditions[0].Type) != "Succeeded" || string(pipelinerunresult.Status.Conditions[0].Status) != "True" {
					app.Status.Status = string(pipelinerunresult.Status.Conditions[0].Type)
					app.Status.Ready = "Pending"
					//app.Status.PipelineRun = pipelinerunresult
					app.Status.Domain = "NotReady"
					err := r.client.Status().Update(context.TODO(), app)
					if err != nil {
						reqLogger.Error(err, "Failed to update app status during build", "app name", app.Spec.AppName)
						return reconcile.Result{}, err
					}
					reqLogger.Info("Finish the whole process, waiting for result", "pipelinerun name", pipelinerunresult.Name, "pipelinerun type", pipelinerunresult.Status.Conditions[0].Type, "pipelinerun status", pipelinerunresult.Status.Conditions[0].Status)
					return reconcile.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
				} else {
					reqLogger.Info("The pipeline is done, wait for service startup", "pipelinerun name", pipelinerunresult.Name, "pipelinerun type", pipelinerunresult.Status.Conditions[0].Type, "pipelinerun status", pipelinerunresult.Status.Conditions[0].Status)
					approute := &servingv1.Route{}
					err = r.client.Get(context.TODO(), client.ObjectKey{Namespace: app.Namespace, Name: app.Spec.AppName}, approute)
					if err != nil {
						reqLogger.Error(err, "Failed to get route for app", "app name", app.Spec.AppName)
						return reconcile.Result{}, err
					} else {
						if approute.Status.Conditions[2].Status != "True" {
							app.Status.Status = string(pipelinerunresult.Status.Conditions[0].Type)
							app.Status.Ready = "Deployed"
							//app.Status.PipelineRun = pipelinerunresult
							app.Status.Domain = "Preparing"
							app.Status.Instance = 0
							err := r.client.Status().Update(context.TODO(), app)
							if err != nil {
								reqLogger.Error(err, "Failed to update app status after process", "app name", app.Spec.AppName)
								return reconcile.Result{}, err
							}
							reqLogger.Info("The route is not ready", "route status", approute.Status, "route domain", approute.Spec)
							return reconcile.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
						} else {
							app.Status.Status = string(pipelinerunresult.Status.Conditions[0].Type)
							app.Status.Ready = "Running"
							app.Status.Domain = approute.Status.DeprecatedDomain
							app.Status.Instance = 1
							err := r.client.Status().Update(context.TODO(), app)
							if err != nil {
								reqLogger.Error(err, "Failed to update app status when route is ready", "app name", app.Spec.AppName)
								return reconcile.Result{}, err
							}
							reqLogger.Info("The route is ready", "route status", approute.Status, "route domain", approute.Spec)
						}
					}

					return reconcile.Result{}, nil
				}
			}
		}
	}
}

func appSpecr(spec pipelinev1alpha1.PipelineResourceSpec) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		app := obj.(*pipelinev1alpha1.PipelineResource)
		app.Spec = spec
		return nil
	}
}

func createOrUpdateGitResource(r *ReconcileAppengine, app *knapv1alpha1.Appengine) error {
	// Step 1 Create or update git resource for application
	gitresource := &pipelinev1alpha1.PipelineResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Spec.AppName + "-git",
			Namespace: app.Namespace,
		},
		Spec: pipelinev1alpha1.PipelineResourceSpec{
			Type: pipelinev1alpha1.PipelineResourceTypeGit,
			Params: []pipelinev1alpha1.Param{{
				Name:  "url",
				Value: app.Spec.GitRepo,
			}, {
				Name:  "revision",
				Value: app.Spec.GitRevision,
			}},
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, gitresource, appSpecr(gitresource.Spec))
	// err = r.client.Create(context.TODO(), gitresource)
	return err
}

func runPipeline(r *ReconcileAppengine, app *knapv1alpha1.Appengine) (*pipelinev1alpha1.PipelineRun, error) {
	// Step 2 Run build and deploy pipeline
	pipelinerun := &pipelinev1alpha1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Spec.AppName + "-pr-" + fmt.Sprintf("%d", app.ObjectMeta.Generation),
			Namespace: app.Namespace,
		},
		Spec: pipelinev1alpha1.PipelineRunSpec{
			PipelineRef: pipelinev1alpha1.PipelineRef{
				Name: "build-and-deploy-pipeline",
			},
			ServiceAccount: "pipeline-account",
			Resources: []pipelinev1alpha1.PipelineResourceBinding{{
				Name: "git-source",
				ResourceRef: pipelinev1alpha1.PipelineResourceRef{
					Name: app.Spec.AppName + "-git",
				},
			},
			},
			Params: []pipelinev1alpha1.Param{{
				Name:  "pathToYamlFile",
				Value: "knative/" + app.Spec.AppName + ".yaml",
			}, {
				Name:  "imageUrl",
				Value: "us.icr.io/knative_jordan/" + app.Spec.AppName,
			}, {
				Name:  "imageTag",
				Value: fmt.Sprintf("%d", app.ObjectMeta.Generation) + ".0",
			}},
		},
	}
	err := r.client.Create(context.TODO(), pipelinerun)
	return pipelinerun, err
}

func runProcess(r *ReconcileAppengine, app *knapv1alpha1.Appengine) (*pipelinev1alpha1.PipelineRun, error) {
	err := createOrUpdateGitResource(r, app)
	if err != nil {
		return nil, err
	}
	pipelinerun, err := runPipeline(r, app)
	if err != nil {
		return nil, err
	}
	return pipelinerun, nil
}
