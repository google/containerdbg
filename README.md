containerdbg: Automate container debugging tasks
============

containerdbg is an all-in-one command-line tool to help debug Kubernetes
containers with common issues that arise when moving to containers as part of
legacy application modernization.

Currently the tool looks for the following common issues:

* Files missing in the container image - by tracking failed open-file requests
  the tool can find files that weren't added to the container image (either at
  build-time or via mount). There is also special logic to specifically handle
  missing library support files (`*.so`, `*.py`, `*.rb`) 

* `EX_DEV` move failures - many legacy applications rely on a `move` or `rename` operation to be atomic - 
   but when moving files that were originally part of the base image - the
   overlay filesystem that supports containers needs to perform a copy and
   delete operation - this can cause strange hard to debug errors

* Failed network connections - sometimes a legacy application depends on a
  network service which is not available in the Kubernetes deployment - and not
  always is this easy to recognize - any failed network connection will be
  logged. The tool also recongizes failed DNS queries - which sometimes mean the
  network service is available but under a different name or undiscoverable by
  the workload.

* Static IP address usage - although this is not recommended - some legacy
  applications have configurations with static IP addresses which will definitly
  not work with a Kuberenetes based deployment. The tool uses the DNS and IP
  tracking - and if it finds a request to a network resource that wasn't
  initiated based on a DNS request it will warn you to update your
  configuration.

## Installation

Download the pre-compiled binaries:

* [Linux (amd64)](https://github.com/google/containerdbg/releases/download/v0.0.8/containerdbg_0.0.8_linux_amd64.tar.gz)
* [Linux (arm64)](https://github.com/google/containerdbg/releases/download/v0.0.8/containerdbg_0.0.8_linux_arm64.tar.gz)
* [MacOS (amd64)](https://github.com/google/containerdbg/releases/download/v0.0.8/containerdbg_0.0.8_darwin_amd64.tar.gz)
* [MacOS (arm64)](https://github.com/google/containerdbg/releases/download/v0.0.8/containerdbg_0.0.8_darwin_arm64.tar.gz)

```bash
VERSION=0.0.8
OS=linux
ARCH=amd64
tar xf containerdbg_${VERSION}_${OS}_${ARCH}.tar.gz
chmod +x containerdbg
sudo mv containerdbg /usr/local/bin
```

## containerdbg components

* containerdbg CLI: The containerdbg CLI provides you with simplified commands to debug and analyze debugging traces from your containers.
* containerdbg daemon: A Daemonset which is deployed by the CLI tool and manages all the eBPF filters needed for containerdbg debugging capabilities.
* entrypoint: A binary which will replace the original entrypoint in the container and will register the container for analysis.
* dnsproxy: A sidecar which acts as a proxy for all dns requests coming from the debugged container.

Building
--------

Check out [BUILDING](BUILDING.md)


Usage
-----
In this section we will show 2 main usage scenarios for containerdbg.

For a step-by-step guide with example application please refer to the [guide](examples/petclinic/).

### Analyzing a deployment yaml
This is use case is for when you have a kubernetes yaml which contains a Deployment resource. If this deployment yaml contains more than one Deployment resource please consider splitting it for simplicity.

1. the first stage is to run the following command `containerdbg debug -f <yaml file> -o record.pb`.
this will apply all the resources in the yaml file and modify the deployment inside the yaml to by debugged by containerdbg.
2. The output should look something like:

```bash
Installing containerdbg node daemon
NAMESPACE   RESOURCE                                  ACTION        STATUS      RECONCILED  CONDITIONS                                AGE     MESSAGE
            Namespace/containerdbg-system             Unchanged     Current                 <None>                                    4s      Resource is current
containerd  DaemonSet/containerdbg-daemonset          Created       Current                 <None>                                    2s      All replicas scheduled as expected. Repl

Press Ctrl-C to finish the debugging session and download the collected report
```
At this point, you can work with your deployment for a while until some errors occur. once you are done press Ctrl-C on the terminal in which you ran containerdbg to finish collecting information.

3. Now you can take record.pb and try to get a summary for the issues discovered by running `containerdbg analyze -f record.pb` this will print a short summary of what could have went wrong during the execution of your container.
```

While executing the container the following files were missing:
===============================================================
/var/lib/dpkg/arch is missing
/var/lib/dpkg/triggers/File is missing

While executing the container the library type files were missing:
==================================================================

While executing the container the following files where attempted to be moved but failed to docker limitation:
==============================================================================================================
```

### Analyzing a container image
In case you don't have a kubernetes yaml and you simply want to test an image you could run the following command `containerdbg debug <image> -o record.pb`.
As in the previous section the output should look like:

```bash
Installing containerdbg node daemon
NAMESPACE   RESOURCE                                  ACTION        STATUS      RECONCILED  CONDITIONS                                AGE     MESSAGE
            Namespace/containerdbg-system             Unchanged     Current                 <None>                                    4s      Resource is current
containerd  DaemonSet/containerdbg-daemonset          Created       Current                 <None>                                    2s      All replicas scheduled as expected. Repl

Press Ctrl-C to finish the debugging session and download the collected report
```
At this point, you can work with your deployment for a while until some errors occur. once you are done press Ctrl-C on the terminal in which you ran containerdbg to finish collecting information.

Now you can take record.pb and try to get a summary for the issues discovered by running `containerdbg analyze -f record.pb` this will print a short summary of what could have went wrong during the execution of your container.
```
While executing the container the following files were missing:
===============================================================
/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser is missing

While executing the container the library type files were missing:
==================================================================

While executing the container the following files where attempted to be moved but failed to docker limitation:
==============================================================================================================

While executing the container the following connections failed:
==============================================================================================================
10.108.0.105:5432
```

In this example we can see some missing configration files that have no real importance to the application and a failed connection to a posgres DB which might be the reason the application is failing.

Troubleshooting
---------------
In case you see the following errors
```
program sys_enter_open: apply CO-RE relocations: no BTF found for kernel version <version>: not supported
```

It means your cluster does not have btf support, in order to resolve this issue you can download the corresponding btf file from https://github.com/aquasecurity/btfhub-archive/ with the matching `<version>` and then do the following:
1. extract the downloaded file into `btf-install/` using `tar xf <filename>`
2. copy the resulting .btf file into btf-install folder.
3. run `export TARGET_REPO=<repo name>` where repo name is an image registry accesible to your cluster.
4. run `make install-btf`

The program should run succefully now.

## Technical background

The tools works by first deploying a workload (either from a YAML file or a container image) - it can also connect to an existing deployment.

Then it utilizes eBPF and a sidecar to instrument the workload while running -
so you can try and use the workload as usual.

When you are done you finish the analyze stage and the tool collects the
recorded data - and analyzes it displaying a report with issues found.

Contributing
------------

Contributions are welcome, see [CONTRIBUTING](./CONTRIBUTING.md)

Community
--------

**Come and ask us questions**
* [containerdbg-users mailing list](https://groups.google.com/g/containerdbg-users)
* For issues and feature requests please open a GitHub [issue](https://github.com/google/containerdbg/issues/new).

Thanks
------

* [miekg/dns](https://github.com/miekg/dns): containerdbg uses miekg's DNS library to trace all DNS requests coming from workloads.
* [cilium/ebpf](https://github.com/cilium/ebpf): containerdbg uses the Cilium eBPF library to inspect workloads.
