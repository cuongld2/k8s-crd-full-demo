# How to Access Kubernetes CRDs from client-go

## Introduction

Kubernetes is a popular container orchestration platform responsible for creating and managing containers that allow software applications to scale out to handle the growing workload of users. Besides built-in resources such as pods or deployments, Kubernetes provides
custom resource definition (CRD) support so that you can define your resources with the exact format you need. Kubernetes CRD provides you the following benefits:

- You can use the powerful command line utility `kubectl` with a number of functionalities like creating or updating resources
- The custom resources are managed by Kubernetes directly so that they can scale out or scale in when needed
- Kubernetes also provides a client tool that allows you to interact with Kubernetes resources programmatically.
- Kubernetes supports several popular programming languages for client tool such as Python, Java, Javascript, or Go. 

This article will show you how to access and manipulate Kubernetes CRDS using the [client-go](https://github.com/kubernetes/client-go).

## Demonstration

Let's say your software department relies on Kubernetes to build applications and tools for both production and internal purposes. When implementing a new application, you wonder whether the existing Kubernetes cluster already provided the database you need to store data for the new app. To resolve that problem, you create a custom resource to manage databases inside the Kubernetes cluster. You can search for more information about the new `database` resource, such as the currently supported database in Kubernetes, the total number of database instances, or the available database instances for each database.

### Prerequisites

To follow along with the article, you would need to prepare the following:
1. A ready-to-use Linux machine. This demo uses an Ubuntu machine with version 20.04.
2. A ready-to-use Kubernetes cluster. You can set up your Kubernetes cluster using virtual machines or a Kubernetes cloud provider solution. This demo uses the Vultr Kubernetes Engine version 1.24.4.
3. A ready-to-use Go environment since you will build the tool using Kubernetes client-go to interact with the `database` custom resource. This demo uses Go version 1.19.
4. A ready-to-use `kubectl` tool to interact with the Kubernetes cluster from the command line.

### Access the Vultr Kubernetes cluster using kubectl

If you choose to use Vultr Kubernetes cluster for the demo (Vultr offers 100$ credits for you to try out Vultr products when you create your new account), please download the Kubernetes config file for your Kubernetes server from Vultr page. You need that config file to access the Vultr Kubernetes cluster.

1. Download the config file
From the Kubernetes overview page in Vultr, click the "Download Configuration" button to download the config file.
![Vultr Kubernetes page](https://i.imgur.com/zBq2tLn.png)
The downloaded file will have a name like "vke-2adb84ff-6b5a-4e5e-a92e-0cafa3b13486.yaml". You should rename it to "vke.yaml" and move it to the home directory for convenience. Open up your terminal and type the following commands:

    $ cd ~/Downloads
    $ mv ${your_config_file.yaml} ~/vke.yaml



2. Export the config file to environment variable

You need to export the config file as an environment variable for the `kubectl` command line tool to access the Kubernetes cluster. Run the commands below:


    $ cd ~

    // Get your current home directory path 
    $ echo $HOME 
    $ export KUBECONFIG='${path_to_your_home_directory}/vke.yaml'
    $ kubectl get node

You should be able to see the nodes that the Kubernetes cluster has, similar to the below:

    NAME                   STATUS   ROLES    AGE     VERSION
    k8s-crd-ba11fd0aaa9b   Ready    <none>   6d20h   v1.24.4
    k8s-crd-e29c4afea916   Ready    <none>   6d20h   v1.24.4

Now that you can successfully access the Kubernetes cluster using `kubectl`, let's move on to the next section to see how to create the `database` custom resource definition using `kubectl`.

### Create the Kubernetes custom resource definition using kubectl

Kubernetes use `yaml` files as instructions to allow clients to interact with the Kubernetes server. The custom resource definition file (CRD file) is in `yaml` format. The CRD file provides information such as `apiVersion`, `metadata`, `spec`, and `scope` of the resource. Check out [Kubernetes guide for creating custom resource definitions](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) for more details about how CRD file works. 

1. Add a new `database` custom resource definition

Run the following commands to create the custom resource definition file:

    $ mkdir k8s-crd-demo
    $ cd k8s-crd-demo
    $ nano dbs_crd.k8s.yaml

Then copy the following yaml definitions into the `dbs_crd.k8s.yaml`, save the file.

    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
    name: databases.resource.donald.com
    spec:
    group: resource.donald.com
    versions:
        - name: v1
        served: true
        storage: true
        schema:
            openAPIV3Schema:
            type: object
            properties:
                spec:
                type: object
                properties:
                    dbName:
                    type: string
                    nullable: false
                    description:
                    type: string
                    nullable: false
                    total:
                    type: integer
                    default: 10
                    minimum: 1
                    maximum: 100
                    available:
                    type: integer
                    default: 10
                    minimum: 1
                    maximum: 100
                    dbType:
                    type: string
                    enum:
                    - sql
                    - noSQL
                    - timeSeries
                    - messageQueue
                    - caching
                    nullable: false
                    tags:
                    type: string
                    nullable: true
                required: ["dbName", "total", "available", "dbType"]
            required: ["spec"]
    scope: Cluster
    names:
        plural: databases
        singular: database
        kind: Database
        shortNames:
        - db

Here, you define the `apiVersion` for the custom resource with `apiextensions.k8s.io/v1`, which is version 1 for API extensions of Kubernetes. The name of the CRD is `databases.resource.donald.com`. The name of the resource group is `resource.donald.com`. You need to use these names when interacting with the Kubernetes custom resources using the Kubernetes go-client tool. The `scope` of the custom resource by default is `Cluster`, which means you can access the custom from anywhere inside the Kubernetes cluster. You can also set the `scope` value to `Namespace` to restrict access to the custom resource inside a particular namespace.

The `database` custom resource has information about `dbName`, `description`, `total`, `available`, `dbType`, and `tags`. The `total` and `available` fields you restrict to be `integer` data types and have values in the range from 1 to 100 instances. The `dbType` must be `string` and can only be one of the values as `sql`, `noSQL`, `timeSeries`, `messageQueue`, or `caching`.

To create this `database` custom resource in the Kubernetes cluster, run the following command:


    $ kubectl apply -f dbs_crd.k8s.yaml


Using the `apply` option with `kubectl` tells the Kubernetes cluster to create or update the target resource. The `-f` option indicates that you are using a file to apply the action. You should be able to see similar output like:

    customresourcedefinition.apiextensions.k8s.io/databases.resource.donald.com created

Now you successfully create the custom resource definition.

2. Add a new database to `database` custom resource definition

Let's add a new database resource item into the `database` custom resource definition. To do it, run the below commands:

    $ nano mysql_resource_object.yaml

Copy the following content into the `mysql_resource_object.yaml`:

    apiVersion: "resource.donald.com/v1"
    kind: Database
    metadata:
    name: mysql
    spec:
    dbName: mysql
    description: Used for storing relation structured data.
    total: 50
    available: 50
    dbType: sql
    tags: Web Development, Data Engineering, Embedded software

You set the `apiVersion` for the resource definition with the value `resource.donald.com/v1`. The `apiVersion` must be in the format of resourceGroup.version. The `kind` of resource is `Database` and must match the `kind` of the custom resource definition you already created earlier. The name of the `database` item is "mysql" with `dbType` as "sql" and `available` instances are 50.

Run the following command to add the `mysql` database item to the `database` resource definition.

    $ kubectl apply -f mysql_resource_object.yaml

Similar to creating the resource definition, we use `kubectl` with the `apply` option to add a new resource. You should be able to see similar output like:

    database.resource.donald.com/mysql created

You now successfully add the "mysql" resource to the `database` custom resource definition. To check the available databases in the Kubernetes cluster, run the following:

    $ kubectl get db


You should be able to see the output like:

    NAME    AGE
    mysql   2m58s

Or you can get detailed information for the `database` custom resource definition using the following command:

    $ kubectl get db -o yaml

The output should look like this:

    apiVersion: v1
    items:
    - apiVersion: resource.donald.com/v1
    kind: Database
    metadata:
        annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
            {"apiVersion":"resource.donald.com/v1","kind":"Database","metadata":{"annotations":{},"name":"mysql"},"spec":{"available":50,"dbName":"mysql","dbType":"sql","description":"Used for storing relation structured data.","tags":"Web Development, Data Engineering, Embedded software","total":50}}
        creationTimestamp: "2022-11-17T17:58:30Z"
        generation: 1
        name: mysql
        resourceVersion: "1419745"
        uid: 40ed6d7e-a372-4f64-8400-20376fd8fdba
    spec:
        available: 50
        dbName: mysql
        dbType: sql
        description: Used for storing relation structured data.
        tags: Web Development, Data Engineering, Embedded software
        total: 50
    kind: List
    metadata:
    resourceVersion: ""


At this step, you successfully create the `database` custom resource definition and added `mysql` database. Let's move on to see how you can programmatically access the `database` custom resource definition using Go with the help of [Kubernetes go-client tool](https://github.com/kubernetes/client-go).

### Interact with Kubernetes custom resources using go-client

You must initiate a go module environment and install the needed dependencies to build an app that interacts with the Kubernetes custom resources.

#### Install needed dependencies

Open the terminal and type the following `go mod` command to initialize the go module environment.

    $ go mod init k8s-resource.com/m

The go module will automatically create a `go.mod` file. Add the following dependencies into your app's `go.mod` file to connect with the Kubernetes cluster.


    require k8s.io/client-go v0.24.4

    require (
        github.com/google/go-cmp v0.5.9 // indirect
        github.com/kr/pretty v0.3.0 // indirect
        github.com/rogpeppe/go-internal v1.8.0 // indirect
        github.com/stretchr/testify v1.7.1 // indirect
        gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
        sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
        sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
        sigs.k8s.io/yaml v1.2.0 // indirect
    )

    require (
        k8s.io/api v0.24.4 // indirect
        k8s.io/apimachinery v0.24.4
    )

    require (
        github.com/davecgh/go-spew v1.1.1 // indirect
        github.com/go-logr/logr v1.2.3 // indirect
        github.com/gogo/protobuf v1.3.2 // indirect
        github.com/golang/protobuf v1.5.2 // indirect
        github.com/google/gofuzz v1.2.0 // indirect
        github.com/imdario/mergo v0.3.13 // indirect; indirectap
        github.com/json-iterator/go v1.1.12 // indirect
        github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
        github.com/modern-go/reflect2 v1.0.2 // indirect
        github.com/spf13/pflag v1.0.5 // indirect
        golang.org/x/net v0.2.0 // indirect
        golang.org/x/oauth2 v0.2.0 // indirect
        golang.org/x/sys v0.2.0 // indirect
        golang.org/x/term v0.2.0 // indirect
        golang.org/x/text v0.4.0 // indirect
        golang.org/x/time v0.2.0 // indirect
        google.golang.org/appengine v1.6.7 // indirect
        google.golang.org/protobuf v1.28.1 // indirect
        gopkg.in/inf.v0 v0.9.1 // indirect
        gopkg.in/yaml.v2 v2.4.0 // indirect
        k8s.io/klog/v2 v2.80.1 // indirect
        k8s.io/utils v0.0.0-20221108210102-8e77b1f39fe2 // indirect
    )

**NOTE**:
The version of the go-client library (currently, you are using version v24.4) should match the Kubernetes cluster version (v24.4) to prevent incompatible issues. Check out [this guide](https://github.com/kubernetes/client-go#compatibility-matrix) for compatibility matrix details.


Then run `go mod tidy` to install these dependencies:

    $ go mod tidy

Now that you install the dependencies, let's write some codes to interact with the Kubernetes `database` custom resources.

#### Write the code to interact with the Kubernetes custom resources

Let's write the code that allows the app to:
- Create a new custom resource
- Remove an existing one 
- Get all the current custom resources 
- Get the custom resource by the resource name 

To do it, you use several built-in methods from Kubernetes go-client:

    type Interface interface {
        GetRateLimiter() flowcontrol.RateLimiter
        Verb(verb string) *Request
        Post() *Request
        Put() *Request
        Patch(pt types.PatchType) *Request
        Get() *Request
        Delete() *Request
        APIVersion() schema.GroupVersion
    }

You use the `Post` method to create a new resource, `Get` to retrieve all the resources or a specific resource by its name, and `Delete` to remove an existing resource.

##### Implemented Database structs and methods to interact with Kubernetes runtime

1. Create Database structs

You must create structs for `DatabaseSpec`, `Database`, and `DatabaseList` to interact with the existing `database` custom resource definition. Run the following commands to create a new `database.go` file.


    $ mkdir api
    $ cd api
    $ nano database.go

Copy the following codes into the `database.go` file:

    package api

    import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    type DatabaseSpec struct {
        DbName      string `json:"dbName"`
        Description string `json:"description,omitempty"`
        Total       int    `json:"total"`
        Available   int    `json:"available"`
        DbType      string `json:"dbType"`
        Tags        string `json:"tags,omitempty"`
    }

    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
    type Database struct {
        metav1.TypeMeta   `json:",inline"`
        metav1.ObjectMeta `json:"metadata,omitempty"`

        Spec DatabaseSpec `json:"spec"`
    }

    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
    type DatabaseList struct {
        metav1.TypeMeta `json:",inline"`
        metav1.ListMeta `json:"metadata,omitempty"`

        Items []Database `json:"items"`
    }


The `DatabaseSpec` have fields that match with the current spec `database` resource definition are `dbName`, `description`, `total`, `available`,`dbType`, and `tags`. Similarly, the `Database` and `DatabaseList` structs consist of fields that match with `database` resource definition metadata information.

2. Create DeepCopy methods

You create a `deepcopy.go` file to define methods so your app can interact with the Kubernetes runtime.

    $ nano deepcopy.go

Copy the following code into the `deepcopy.go` file.

    package api

    import "k8s.io/apimachinery/pkg/runtime"

    func (in *Database) DeepCopyInto(out *Database) {
        out.TypeMeta = in.TypeMeta
        out.ObjectMeta = in.ObjectMeta
        out.Spec = DatabaseSpec{
            DbName:      in.Spec.DbName,
            Description: in.Spec.Description,
            Total:       in.Spec.Total,
            Available:   in.Spec.Available,
            DbType:      in.Spec.DbType,
            Tags:        in.Spec.Tags,
        }
    }

    func (in *Database) DeepCopyObject() runtime.Object {
        out := Database{}
        in.DeepCopyInto(&out)

        return &out
    }

    func (in *DatabaseList) DeepCopyObject() runtime.Object {
        out := DatabaseList{}
        out.TypeMeta = in.TypeMeta
        out.ListMeta = in.ListMeta

        if in.Items != nil {
            out.Items = make([]Database, len(in.Items))
            for i := range in.Items {
                in.Items[i].DeepCopyInto(&out.Items[i])
            }
        }

        return &out
    }

Here you define the `DeepCopyInto` method for the `Database` struct, the `DeepCopyObject` method for the `Database` struct, and another `DeepCopyObject` method for the `DatabaseList` struct so that Kubernetes runtime can interact with these defined structs.


3. Adding schema types to work with Kubernetes runtime

Create the `register.go` file for adding schema types to work with Kubernetes runtime.


    $ nano register.go

Then copy the following code into `register.go` file:

    package api

    import (
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/apimachinery/pkg/runtime"
        "k8s.io/apimachinery/pkg/runtime/schema"
    )

    const GroupName = "resource.donald.com"
    const GroupVersion = "v1"

    var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

    var (
        SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
        AddToScheme   = SchemeBuilder.AddToScheme
    )

    func addKnownTypes(scheme *runtime.Scheme) error {
        scheme.AddKnownTypes(SchemeGroupVersion,
            &Database{},
            &DatabaseList{},
        )

        metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
        return nil
    }


You set the `GroupName` and `GroupVersion` that match the group name and group version of the `database` custom resource definition. Then inside the `addKnownTypes` function, you add the type for `Database` and `DatabaseList` to Kubernetes runtime.

You have just implemented the Go structs, functions, and methods to interact with Kubernetes runtime at this step. The next part of the article is about defining the Kubernetes client and methods to:
- Create a new resource
- Get existing resources
- Delete an existing one.

##### Implementing Kubernetes client and methods for interacting with Kubernetes custom resources

1. Define the configuration for Kubernetes Rest client
You need to define the configuration for the Kubernetes Rest client. Run the following commands to create a new `api.go` file.

    $ cd ..
    $ mkdir clientset
    $ cd clientset
    $ nano api.go

Copy the following code into `api.go`:

    package clientset

    import (
        "context"

        "k8s-resource.com/m/api"
        "k8s.io/apimachinery/pkg/runtime/schema"
        "k8s.io/client-go/kubernetes/scheme"
        "k8s.io/client-go/rest"
    )

    type ExampleInterface interface {
        Databases(ctx context.Context) DatabaseInterface
    }

    type ExampleClient struct {
        restClient rest.Interface
    }

    func NewForConfig(c *rest.Config) (*ExampleClient, error) {
        config := *c
        config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: api.GroupName, Version: api.GroupVersion}
        config.APIPath = "/apis"
        config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
        config.UserAgent = rest.DefaultKubernetesUserAgent()

        client, err := rest.RESTClientFor(&config)
        if err != nil {
            return nil, err
        }

        return &ExampleClient{restClient: client}, nil
    }

    func (c *ExampleClient) Databases(ctx context.Context) DatabaseInterface {
        return &databaseClient{
            restClient: c.restClient,
            ctx:        ctx,
        }
    }

Here you add the rest configuration for the Kubernetes client to connect with `database` custom resources.

2. Add methods for creating, deleting, and getting custom resources

You need to create a new file named `databases.go`.

    $ nano databases.go

Copy the following code into the `databases.go` file.


    package clientset

    import (
        "context"

        "k8s-resource.com/m/api"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/kubernetes/scheme"
        "k8s.io/client-go/rest"
    )

    type DatabaseInterface interface {
        List(opts metav1.ListOptions) (*api.DatabaseList, error)
        Get(name string, options metav1.GetOptions) (*api.Database, error)
        Create(*api.Database) (*api.Database, error)
        Delete(name string, options metav1.DeleteOptions) (*api.Database, error)
        // ...
    }

    type databaseClient struct {
        restClient rest.Interface
        ctx        context.Context
    }

    func (c *databaseClient) List(opts metav1.ListOptions) (*api.DatabaseList, error) {
        result := api.DatabaseList{}
        err := c.restClient.
            Get().
            AbsPath("/apis/resource.donald.com/v1/databases").
            Do(c.ctx).
            Into(&result)

        return &result, err
    }

    func (c *databaseClient) Get(name string, opts metav1.GetOptions) (*api.Database, error) {
        result := api.Database{}
        err := c.restClient.
            Get().
            AbsPath("/apis/resource.donald.com/v1/databases").
            Name(name).
            VersionedParams(&opts, scheme.ParameterCodec).
            Do(c.ctx).
            Into(&result)

        return &result, err
    }

    func (c *databaseClient) Create(database *api.Database) (*api.Database, error) {
        result := api.Database{}
        err := c.restClient.
            Post().
            AbsPath("/apis/resource.donald.com/v1/databases").
            Body(database).
            Do(c.ctx).
            Into(&result)

        return &result, err
    }

    func (c *databaseClient) Delete(name string, opts metav1.DeleteOptions) (*api.Database, error) {

        result := api.Database{}

        err := c.restClient.
            Delete().
            AbsPath("/apis/resource.donald.com/v1/databases").
            Name(name).
            VersionedParams(&opts, scheme.ParameterCodec).
            Do(c.ctx).Into(&result)
        return &result, err
    }

Here you define the `Create` method to create a new resource, the `Get` method to get a resource by name, the `List` to get all current resources, and the `Delete` to remove an existing resource no longer needed.

Now you've added the codes for defining the Kubernetes client and methods to interact with Kubernetes custom resources. Let's move on to create a `main.go` file.

##### Creating a `main.go` file to interact with the Kubernetes resources.

Suppose that in your next software project, you need to use MongoDB to store data for your app. To add the "mongodb" database into the `database` custom resource definition, you need to do these below steps:

1. Copy the `vke.yaml` config file into the current directory.

    $ cd ..
    $ cp ~/vke.yaml .

2. Create a `main.go` file.

    $ cd ..
    $ nano main.go


3. Add the following code to the `main.go` file:

    package main

    import (
        "context"
        "flag"
        "fmt"
        "log"
        "os"

        "k8s-resource.com/m/api"
        client "k8s-resource.com/m/clientset"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/tools/clientcmd"

        "k8s.io/client-go/kubernetes/scheme"
        "k8s.io/client-go/rest"
    )

    var kubeconfig string

    func init() {
        path, err := os.Getwd()
        if err != nil {
            log.Println(err)
        }
        flag.StringVar(&kubeconfig, "kubeconfig", path+"/vke.yaml", "path to Kubernetes config file")
        flag.Parse()
    }

    func main() {
        var config *rest.Config
        var err error

        if kubeconfig == "" {
            log.Printf("using in-cluster configuration")
            config, err = rest.InClusterConfig()
        } else {
            log.Printf("using configuration from '%s'", kubeconfig)
            config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        }

        if err != nil {
            panic(err)
        }

        api.AddToScheme(scheme.Scheme)

        clientSet, err := client.NewForConfig(config)
        if err != nil {
            panic(err)
        }

        context := context.TODO()

        newDatabase := new(api.Database) // pa == &Student{"", 0}
        newDatabase.Name = "mongodb"
        newDatabase.Kind = "Database" // pa == &Student{"Alice", 0}
        newDatabase.APIVersion = "resource.donald.com/v1"
        newDatabase.Spec.DbName = "mongodb"
        newDatabase.Spec.Description = "Used storing unstructured data"
        newDatabase.Spec.Total = 100
        newDatabase.Spec.Available = 50
        newDatabase.Spec.DbType = "noSQL"
        newDatabase.Spec.Tags = "Web Development, nosql data"
        newDatabase.Spec.Available = 70

        projectCreated, err := clientSet.Databases(context).Create(newDatabase)
        if err != nil {
            panic(err)
        }

        fmt.Println(projectCreated)
        }

Here you call the `Create` method to add `mongodb` database to the `database` custom resource definition. 

4. Execute the action

Let's run the `main.go` file.

    $ go run main.go

After running this command, you should see the similar output below:

    2022/11/18 02:14:55 using configuration from '/home/donaldle/Projects/Personal/vultr/k8s-crd/k8s-crd-full-demo/vke.yaml'
    &{{ } {mongodb    f8ba273e-fd1f-4b40-b036-cf13b8c72366 1430720 1 2022-11-18 02:14:55 +0700 +07 <nil> <nil> map[] map[] [] []  [{main Update resource.donald.com/v1 2022-11-18 02:14:55 +0700 +07 FieldsV1 {"f:spec":{".":{},"f:available":{},"f:dbName":{},"f:dbType":{},"f:description":{},"f:tags":{},"f:total":{}}} }]} {mongodb Used storing unstructured data 100 70 noSQL Web Development, nosql data}}

You just added the "mongodb" database. Let's try to get detailed information about the "mongodb" database using the `Get` method. 

5. Get detailed information for "mongodb" database

To do this, replace the `main.go` code with the below code.

    package main

    import (
        "context"
        "flag"
        "fmt"
        "log"
        "os"

        "k8s-resource.com/m/api"
        client "k8s-resource.com/m/clientset"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/tools/clientcmd"

        "k8s.io/client-go/kubernetes/scheme"
        "k8s.io/client-go/rest"
    )

    var kubeconfig string

    func init() {
        path, err := os.Getwd()
        if err != nil {
            log.Println(err)
        }
        flag.StringVar(&kubeconfig, "kubeconfig", path+"/vke.yaml", "path to Kubernetes config file")
        flag.Parse()
    }

    func main() {
        var config *rest.Config
        var err error

        if kubeconfig == "" {
            log.Printf("using in-cluster configuration")
            config, err = rest.InClusterConfig()
        } else {
            log.Printf("using configuration from '%s'", kubeconfig)
            config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        }

        if err != nil {
            panic(err)
        }

        api.AddToScheme(scheme.Scheme)

        clientSet, err := client.NewForConfig(config)
        if err != nil {
            panic(err)
        }

        context := context.TODO()

        projectGet, err := clientSet.Databases(context).Get("mongodb", metav1.GetOptions{})
        if err != nil {
            panic(err)
        }

        fmt.Println(projectGet)

    }

Then run the command:

    $ go run main.go

You should see a similar output as below:

    2022/11/18 02:18:20 using configuration from '/home/donaldle/Projects/Personal/vultr/k8s-crd/k8s-crd-full-demo/vke.yaml'
    &{{ } {mongodb    f8ba273e-fd1f-4b40-b036-cf13b8c72366 1430720 1 2022-11-18 02:14:55 +0700 +07 <nil> <nil> map[] map[] [] []  [{main Update resource.donald.com/v1 2022-11-18 02:14:55 +0700 +07 FieldsV1 {"f:spec":{".":{},"f:available":{},"f:dbName":{},"f:dbType":{},"f:description":{},"f:tags":{},"f:total":{}}} }]} {mongodb Used storing unstructured data 100 70 noSQL Web Development, nosql data}}


6. Remove "mysql" database from Kubernetes cluster

Let's say you no longer need the `mysql` database in the Kubernetes cluster. To remove the `mysql` resource from the Kubernetes cluster, replace the code in `main.go` with the following code:

    package main

    import (
        "context"
        "flag"
        "log"
        "os"

        "k8s-resource.com/m/api"
        client "k8s-resource.com/m/clientset"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/tools/clientcmd"

        "k8s.io/client-go/kubernetes/scheme"
        "k8s.io/client-go/rest"
    )

    var kubeconfig string

    func init() {
        path, err := os.Getwd()
        if err != nil {
            log.Println(err)
        }
        flag.StringVar(&kubeconfig, "kubeconfig", path+"/vke.yaml", "path to Kubernetes config file")
        flag.Parse()
    }

    func main() {
        var config *rest.Config
        var err error

        if kubeconfig == "" {
            log.Printf("using in-cluster configuration")
            config, err = rest.InClusterConfig()
        } else {
            log.Printf("using configuration from '%s'", kubeconfig)
            config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        }

        if err != nil {
            panic(err)
        }

        api.AddToScheme(scheme.Scheme)

        clientSet, err := client.NewForConfig(config)
        if err != nil {
            panic(err)
        }

        context := context.TODO()

        _, err = clientSet.Databases(context).Delete("mysql", metav1.DeleteOptions{})
        if err != nil {
            panic(err)
        }

    }


Then run :

    $ go run main.go


7. Check if "mysql" database is actually removed

Now, let's try to get all the current custom resources to see whether you successfully removed the "mysql" database. Replace the existing code in the `main.go` file with the following content:

package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"

    "k8s-resource.com/m/api"
    client "k8s-resource.com/m/clientset"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/clientcmd"

    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/rest"
)

var kubeconfig string

func init() {
    path, err := os.Getwd()
    if err != nil {
        log.Println(err)
    }
    flag.StringVar(&kubeconfig, "kubeconfig", path+"/vke.yaml", "path to Kubernetes config file")
    flag.Parse()
}

func main() {
    var config *rest.Config
    var err error

    if kubeconfig == "" {
        log.Printf("using in-cluster configuration")
        config, err = rest.InClusterConfig()
    } else {
        log.Printf("using configuration from '%s'", kubeconfig)
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
    }

    if err != nil {
        panic(err)
    }

    api.AddToScheme(scheme.Scheme)

    clientSet, err := client.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    context := context.TODO()

    projects, err := clientSet.Databases(context).List(metav1.ListOptions{})
    if err != nil {
        panic(err)
    }

    for _, k := range projects.Items {

        fmt.Println(k.Name)

    }

}

Let's run the `main.go` file:

    $ go run main.go


You should only see the `mongodb` database displayed in the output.

    2022/11/18 02:24:08 using configuration from '/home/donaldle/Projects/Personal/vultr/k8s-crd/k8s-crd-full-demo/vke.yaml'
    mongodb

And that's how you can interact with Kubernetes custom resources using Kubernetes go-client tool.

## Conclusion

The article presented you what Kubernetes CRD is, why you would want to use Kubernetes CRD in your current working project, and you can use Kubernetes go-client tool to interact with Kubernetes CRD programmatically. Working with Kubernetes is both fun and challenging, so get ready to face new obstacles when working with it. If you want to learn more about other use cases of using Kubernetes go-client, checkout [how to write Kubernetes infrastructure in Go with cdk8s](https://www.vultr.com/docs/write-your-kubernetes-infrastructure-in-go-with-cdk8s/) or [creating and managing Kubernetes jobs using go-client](https://www.vultr.com/docs/create-and-manage-kubernetes-jobs-in-go-with-client-go-api/).