# Knative Application Platform (under construction)

This project aims to bring a simple CI/CD experience to end users from building to deploying a Cloud application on Kubernets (or Knative). We offer some CRDs/Controllers and a cli (may be changed to a plugin of [Knative cli](https://github.com/knative/client) in future) which talk to [Tekton](https://github.com/tektoncd/pipeline) (a underlying pipeline system) to track and manage the lifecycle of buiding and deplying the application.

![knap_workflow](https://github.com/bluebosh/knap/blob/master/doc/knap.png)

### Use Case 1 ###
Jordan is a developer of a startup company, he is trying to develop a new cloud application and plan to deploy on Kubernetes. He has a Kubernets cluster sitting on IBM Cloud. He created 3 name spaces `dev`, `staing` and `production`, and he followed the installation guide to enabled Knative and installed Tekton and this project `knap`. At the very early phase, he quickly rolled out the code and just wanted to build and deploy the application on the Kubernete cluster automatically for each new git commit. To archive this, with `knap`, he just needed to run few commands to get knap to watch his git project and help build and run the application in `dev` namespace when a new commit shown up.

```
export KUBECONFIG=xxxxxx
knap create --watch https://github.com/bluebosh/jordan_app --template push --docker-token xxxxxx
```

After that, when a new commit is pushed in that git project, it's acknowedged by knap and knap will build the source code to an image and deploy into `dev` name space of that cluster a new revision.

### Use Case 2 ###
After some development interation cycles, the application got mature, Jordan would like to onboard the 1st production version and hoped to continuous integration and delivery. He needs a CI/CD pipeline to help on that...


### Prerequisites

What things you need to install the software and how to install them

```
Give examples
```

### Installing

A step by step series of examples that tell you how to get a development env running

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo

## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Jordan Zhang** - *cheif committer*
* **Grace Zhang** - *major contributor*
* **Edward Xiao** - *major contributor*
* **Matt Cui** - *project Manager*

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
