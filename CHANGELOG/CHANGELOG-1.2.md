

<hr>


## [v1.2.2](https://github.com/aws/aws-k8s-tester/releases/tag/v1.2.2) (2020-05)

See [code changes](https://github.com/aws/aws-k8s-tester/compare/v1.2.1...v1.2.2).

### `eks`

TODO: update secret, config map writes

### Dependency

- TODO


<hr>


## [v1.2.1](https://github.com/aws/aws-k8s-tester/releases/tag/v1.2.1) (2020-05-15)

See [code changes](https://github.com/aws/aws-k8s-tester/compare/v1.2.0...v1.2.1).

### `eks`

- Fix [`ecriface.ECRAPI.DescribeRepositories` calls](https://github.com/aws/aws-k8s-tester/commit/cc418cc3e8c01727c94c3b8fa8099775106020f5).

### `eksconfig`

- Add [`RepositoryAccountID` fields to `AddOnFargate`, `AddOnIRSA`, `AddOnIRSAFargate`, `AddOnHollowNodesRemote`, `AddOnClusterLoaderRemote`](https://github.com/aws/aws-k8s-tester/commit/cc418cc3e8c01727c94c3b8fa8099775106020f5).

### Dependency

- Upgrade [`github.com/aws/aws-sdk-go`](https://github.com/aws/aws-sdk-go/releases) from [`v1.30.28`](https://github.com/aws/aws-sdk-go/releases/tag/v1.30.28) to [`v1.30.29`](https://github.com/aws/aws-sdk-go/releases/tag/v1.30.29).


<hr>


## [v1.2.0](https://github.com/aws/aws-k8s-tester/releases/tag/v1.2.0) (2020-05-15)

See [code changes](https://github.com/aws/aws-k8s-tester/compare/v1.1.9...v1.2.0).

### `ec2config`

- Improve [README](https://github.com/aws/aws-k8s-tester/commit/4a15ae1d61cf58d286263c16e6074f8e3745077a).

### `eksconfig`

- Improve [README](https://github.com/aws/aws-k8s-tester/commit/4a15ae1d61cf58d286263c16e6074f8e3745077a).
- Remove [unnecessary fields from `AddOnIRSA`](https://github.com/aws/aws-k8s-tester/commit/52666165f7564922deb2e6e304c1f1c73412d691).
- Remove [unnecessary fields from `AddOnIRSAFargate`](https://github.com/aws/aws-k8s-tester/commit/52666165f7564922deb2e6e304c1f1c73412d691).
- Now [`AddOnFargate` optionally takes remote ECR image](https://github.com/aws/aws-k8s-tester/commit/afc73f3a7e77d817b953c5e4fe76be82f30fb6ff).
  - `AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_REPOSITORY_NAME` is optional.
  - `AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_REPOSITORY_URI` is optional.
  - `AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_REPOSITORY_IMAGE_TAG` is optional.
  - See https://github.com/aws/aws-k8s-tester/blob/master/Dockerfile and https://github.com/aws/aws-k8s-tester/blob/master/Makefile for container image build instructions.
- Now [`AddOnIRSA` requires remote ECR image](https://github.com/aws/aws-k8s-tester/commit/afc73f3a7e77d817b953c5e4fe76be82f30fb6ff).
  - `AWS_K8S_TESTER_EKS_ADD_ON_IRSA_REPOSITORY_NAME` is now required.
  - `AWS_K8S_TESTER_EKS_ADD_ON_IRSA_REPOSITORY_URI` is now required.
  - `AWS_K8S_TESTER_EKS_ADD_ON_IRSA_REPOSITORY_IMAGE_TAG` is now required.
  - See https://github.com/aws/aws-k8s-tester/blob/master/Dockerfile and https://github.com/aws/aws-k8s-tester/blob/master/Makefile for container image build instructions.
- Now [`AddOnIRSAFargate` requires remote ECR image](https://github.com/aws/aws-k8s-tester/commit/afc73f3a7e77d817b953c5e4fe76be82f30fb6ff).
  - `AWS_K8S_TESTER_EKS_ADD_ON_IRSA_FARGATE_REPOSITORY_NAME` is now required.
  - `AWS_K8S_TESTER_EKS_ADD_ON_IRSA_FARGATE_REPOSITORY_URI` is now required.
  - `AWS_K8S_TESTER_EKS_ADD_ON_IRSA_FARGATE_REPOSITORY_IMAGE_TAG` is now required.
  - See https://github.com/aws/aws-k8s-tester/blob/master/Dockerfile and https://github.com/aws/aws-k8s-tester/blob/master/Makefile for container image build instructions.

### `eks`

- Fix [`eks/hollow-nodes` with RBAC and clean up node labels](https://github.com/aws/aws-k8s-tester/commit/0f2c940680a8b1d430bf5726d6217d246cfa1ca2).
  - Previously, it did not work because of `"NodeRestriction"` from `"kube-apiserver --enable-admission-plugins"`. Now it works with `"NodeRestriction"`.
  - Add [`nodes/status` resource group](https://github.com/aws/aws-k8s-tester/commit/0aff1fb25565fc94d8fcadfe84c1f97c9ad6325d).
  - Add [`pods/status` resource group](https://github.com/aws/aws-k8s-tester/commit/0aff1fb25565fc94d8fcadfe84c1f97c9ad6325d).
- Fix and improve [`eks/irsa` configmap tests](https://github.com/aws/aws-k8s-tester/commit/52666165f7564922deb2e6e304c1f1c73412d691).
- Fix and improve [`eks/irsa-fargate` configmap tests](https://github.com/aws/aws-k8s-tester/commit/52666165f7564922deb2e6e304c1f1c73412d691).
- Improve [`eks/cluster-loader` `RequestSummary` output and separate results for reads](https://github.com/aws/aws-k8s-tester/commit/968fa2a18001112ca6c952439fe0a45b0dfd2b85).

### `pkg/aws/ssm`

- Check [`ssm.ListCommandInvocationsInput` batch limit](https://github.com/aws/aws-k8s-tester/commit/23d21857342930ceb0e165628ba8c124fb99198d).

### `pkg/metrics`

- Add [`HistogramBuckets.Table` method](https://github.com/aws/aws-k8s-tester/commit/52666165f7564922deb2e6e304c1f1c73412d691).

### Dependency

- Upgrade [`github.com/aws/aws-sdk-go`](https://github.com/aws/aws-sdk-go/releases) from [`v1.30.26`](https://github.com/aws/aws-sdk-go/releases/tag/v1.30.26) to [`v1.30.28`](https://github.com/aws/aws-sdk-go/releases/tag/v1.30.28).


<hr>
