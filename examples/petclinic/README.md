# Petclinic debugging example

In this example we will show how you can use containerdbg in order to debug petclinic, an application which allows you to manage a clinic for pets.

The application contains:
1. A tomcat server with the web application written in Java.
2. A posgres DB which contains the information about the pets.

In addition the normal petclinic Java web app, the app has been modified to throw exception by trying to open non existant files on specific actions, an issue which we will debug.

The code for the petclinic application can be found in examples/spring-framework-petclinic/ and is a modified version of the opensource [project](https://github.com/spring-petclinic/spring-framework-petclinic)

## Preperation

In order to run this tutorial you need to have the following resources:
1. a kubernetes cluster `kind` cluster can be used for this.
2. have `skaffold` installed on your machine https://skaffold.dev/docs/install/

Next we will have to build and prepare our yamls for the tutorial, this can be done in the following steps:
1. Set up your image repository (this is not needed for kind clusters)
```
export REGISTRY=<my repository>
```
2. Run the following command to generate the tomcat yaml
```
cd tomcat
skaffold build -d $REGISTRY -o artifacts.json
skaffold render -a artifacts.json > ../tomcat.yaml
cd ..
```
3. Run the following command to generate the db yaml
```
cd db
skaffold build -d $REGISTRY -o artifacts.json
skaffold render -a artifacts.json > ../db.yaml
cd ..
```

## Deploying a working application
First we will deploy this application and show you where it fails.

In order to start the app, please run the following while being connected to a cluster:
```
kubectl apply -f ./db.yaml
kubectl apply -f ./tomcat.yaml
```

Now retrieve the external IP address of the tomcat site by running:

```
export IP=$(kubectl get services petclinic-tomcat -ojsonpath='{.status.loadBalancer.ingress[0].ip}')
```

Now open your browser and go to `http://$IP:8080/petclinic` try the following:
1. Open find owners page.
2. Add new owner.
3. Add new pet for the owner.

At this point you should see an error that something is missing.

In the next section we will see how to use containerdbg in order to discover what is missing.


## Using containerdbg to detect what is not working

### Connection issues
In this example we will demonstrate a situation where one container was moved to another location but the DB wasn't made to be accesible in the new location.

To demonstrate this, we will delete the previously created db pod by running the following command:
```
kubectl delete statefulsets.apps petclinic-postgres
```

Now by traversing the site (for example go to the find owners page) you will see errors are starting to occur.

Now in order to test what is broken in the tomcat container we will remove the currently installed tomcat container:
```
kubectl delete -f tomcat.yaml
```

Next we will run tomcat using containerdbg:
```
containerdbg debug -f tomcat.yaml -o tomcat.pb
```

Then open your browser and traverse to the petclinic application.
```
export IP=$(kubectl get services petclinic-tomcat -ojsonpath='{.status.loadBalancer.ingress[0].ip}')
```

open your browser and go to `http://$IP:8080/petclinic`.

You should see an error when entering the site, after seeing the error press Ctrl-C in the terminal from which you ran containerdbg.


Now run the analysis tool:
```
containerdbg analyze -f tomcat.pb
```

this should produce the following output:
```
While executing the container the following files were missing:
===============================================================
/usr/local/openjdk-11/conf/jndi.properties is missing
/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser is missing

While executing the container the library type files were missing:
==================================================================

While executing the container the following files where attempted to be moved but failed to docker limitation:
==============================================================================================================

While executing the container the following connections failed:
==============================================================================================================
10.108.12.9:5432
```

The last connection using port 5432 is the postgres DB cluster IP which can be viewed from `kubectl get services`.

Now restore the db by running `kubectl apply -f db.yaml` and proceed to the next example.

### Missing files
In this step we will once again containerdbg to run tomcat
```
containerdbg debug -f tomcat.yaml -o tomcat2.pb
```

And repeat the steps to reproduce the error as described in [Deploying a working application](#deploying-a-working-application).


After seeing the error stop the debugging session by pressing Ctrl-C in the same terminal where you ran containerdbg.

By running `containerdbg analyze -f tomcat2.pb` you will see the following output:
```
While executing the container the following files were missing:
===============================================================
 is missing
/etc/group is missing
/etc/passwd is missing
/etc/selinux/config is missing
/etc/selinux/contexts/lxc_contexts is missing
/usr/lib/locale/locale-archive is missing
/usr/local/openjdk-11/conf/jndi.properties is missing
/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser is missing
/usr/share/containers/selinux/contexts is missing
file.txt is missing

While executing the container the library type files were missing:
==================================================================

While executing the container the following files where attempted to be moved but failed to docker limitation:
==============================================================================================================

While executing the container the following connections failed:
==============================================================================================================
```

By comparing to the old output or by seeing that most files are not relevant (selinux related or /etc/group /etc/passwd files) we can see that the relevant error is caused by file.txt .

This should lead you to edit the tomcat Dockerfile and add the missing file, luckily we can see that there is a commented out line in tomcat/Dockerfile adding this file so simply uncomment it.

Now repeat the steps in [Deploying a working application](#deploying-a-working-application) and the application should run without issues.


# Summary
In this example we saw how, by running containerdbg multiple times we could detect issues in the application and fix them this is an iterative approach that should be repeated multiple times until issues are fixed.


