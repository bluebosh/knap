// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/bluebosh/knap/pkg/apis/knap/v1alpha1"
	scheme "github.com/bluebosh/knap/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AppenginesGetter has a method to return a AppengineInterface.
// A group's client should implement this interface.
type AppenginesGetter interface {
	Appengines(namespace string) AppengineInterface
}

// AppengineInterface has methods to work with Appengine resources.
type AppengineInterface interface {
	Create(*v1alpha1.Appengine) (*v1alpha1.Appengine, error)
	Update(*v1alpha1.Appengine) (*v1alpha1.Appengine, error)
	UpdateStatus(*v1alpha1.Appengine) (*v1alpha1.Appengine, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Appengine, error)
	List(opts v1.ListOptions) (*v1alpha1.AppengineList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Appengine, err error)
	AppengineExpansion
}

// appengines implements AppengineInterface
type appengines struct {
	client rest.Interface
	ns     string
}

// newAppengines returns a Appengines
func newAppengines(c *KnapV1alpha1Client, namespace string) *appengines {
	return &appengines{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the appengine, and returns the corresponding appengine object, and an error if there is any.
func (c *appengines) Get(name string, options v1.GetOptions) (result *v1alpha1.Appengine, err error) {
	result = &v1alpha1.Appengine{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("appengines").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Appengines that match those selectors.
func (c *appengines) List(opts v1.ListOptions) (result *v1alpha1.AppengineList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AppengineList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("appengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested appengines.
func (c *appengines) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("appengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a appengine and creates it.  Returns the server's representation of the appengine, and an error, if there is any.
func (c *appengines) Create(appengine *v1alpha1.Appengine) (result *v1alpha1.Appengine, err error) {
	result = &v1alpha1.Appengine{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("appengines").
		Body(appengine).
		Do().
		Into(result)
	return
}

// Update takes the representation of a appengine and updates it. Returns the server's representation of the appengine, and an error, if there is any.
func (c *appengines) Update(appengine *v1alpha1.Appengine) (result *v1alpha1.Appengine, err error) {
	result = &v1alpha1.Appengine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("appengines").
		Name(appengine.Name).
		Body(appengine).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *appengines) UpdateStatus(appengine *v1alpha1.Appengine) (result *v1alpha1.Appengine, err error) {
	result = &v1alpha1.Appengine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("appengines").
		Name(appengine.Name).
		SubResource("status").
		Body(appengine).
		Do().
		Into(result)
	return
}

// Delete takes name of the appengine and deletes it. Returns an error if one occurs.
func (c *appengines) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("appengines").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *appengines) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("appengines").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched appengine.
func (c *appengines) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Appengine, err error) {
	result = &v1alpha1.Appengine{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("appengines").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
