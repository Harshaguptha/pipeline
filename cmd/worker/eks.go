// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"time"

	"emperror.dev/errors"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"

	eksworkflow "github.com/banzaicloud/pipeline/internal/providers/amazon/eks/workflow"
	pkgEks "github.com/banzaicloud/pipeline/pkg/cluster/eks"
)

const asgWaitLoopSleepSeconds = 5
const asgFulfillmentTimeout = 2 * time.Minute

func registerEKSWorkflows(secretStore eksworkflow.SecretStore) error {

	vpcTemplate, err := pkgEks.GetVPCTemplate()
	if err != nil {
		return errors.WrapIf(err, "failed to get CloudFormation template for VPC")
	}

	subnetTemplate, err := pkgEks.GetSubnetTemplate()
	if err != nil {
		return errors.WrapIf(err, "failed to get CloudFormation template for Subnet")
	}

	iamRolesTemplate, err := pkgEks.GetIAMTemplate()
	if err != nil {
		return errors.WrapIf(err, "failed to get CloudFormation template for IAM roles")
	}

	nodePoolTemplate, err := pkgEks.GetNodePoolTemplate()
	if err != nil {
		return errors.WrapIf(err, "failed to get CloudFormation template for node pools")
	}

	workflow.RegisterWithOptions(eksworkflow.CreateClusterWorkflow, workflow.RegisterOptions{Name: eksworkflow.CreateClusterWorkflowName})
	workflow.RegisterWithOptions(eksworkflow.CreateInfrastructureWorkflow, workflow.RegisterOptions{Name: eksworkflow.CreateInfraWorkflowName})

	awsSessionFactory := eksworkflow.NewAWSSessionFactory(secretStore)

	createVPCActivity := eksworkflow.NewCreateVPCActivity(awsSessionFactory, vpcTemplate)
	activity.RegisterWithOptions(createVPCActivity.Execute, activity.RegisterOptions{Name: eksworkflow.CreateVpcActivityName})

	createSubnetActivity := eksworkflow.NewCreateSubnetActivity(awsSessionFactory, subnetTemplate)
	activity.RegisterWithOptions(createSubnetActivity.Execute, activity.RegisterOptions{Name: eksworkflow.CreateSubnetActivityName})

	getSubnetsDetailsActivity := eksworkflow.NewGetSubnetsDetailsActivity(awsSessionFactory)
	activity.RegisterWithOptions(getSubnetsDetailsActivity.Execute, activity.RegisterOptions{Name: eksworkflow.GetSubnetsDetailsActivityName})

	createIamRolesActivity := eksworkflow.NewCreateIamRolesActivity(awsSessionFactory, iamRolesTemplate)
	activity.RegisterWithOptions(createIamRolesActivity.Execute, activity.RegisterOptions{Name: eksworkflow.CreateIamRolesActivityName})

	uploadSSHActivityActivity := eksworkflow.NewUploadSSHKeyActivity(awsSessionFactory)
	activity.RegisterWithOptions(uploadSSHActivityActivity.Execute, activity.RegisterOptions{Name: eksworkflow.UploadSSHKeyActivityName})

	getVpcConfigActivity := eksworkflow.NewGetVpcConfigActivity(awsSessionFactory)
	activity.RegisterWithOptions(getVpcConfigActivity.Execute, activity.RegisterOptions{Name: eksworkflow.GetVpcConfigActivityName})

	createEksClusterActivity := eksworkflow.NewCreateEksClusterActivity(awsSessionFactory)
	activity.RegisterWithOptions(createEksClusterActivity.Execute, activity.RegisterOptions{Name: eksworkflow.CreateEksControlPlaneActivityName})

	waitAttempts := int(asgFulfillmentTimeout.Seconds() / asgWaitLoopSleepSeconds)
	waitInterval := asgWaitLoopSleepSeconds * time.Second
	createAsgActivity := eksworkflow.NewCreateAsgActivity(awsSessionFactory, nodePoolTemplate, waitAttempts, waitInterval)
	activity.RegisterWithOptions(createAsgActivity.Execute, activity.RegisterOptions{Name: eksworkflow.CreateAsgActivityName})

	createUserAccessKeyActivity := eksworkflow.NewCreateClusterUserAccessKeyActivity(awsSessionFactory)
	activity.RegisterWithOptions(createUserAccessKeyActivity.Execute, activity.RegisterOptions{Name: eksworkflow.CreateClusterUserAccessKeyActivityName})

	bootstrapActivity := eksworkflow.NewBootstrapActivity(awsSessionFactory)
	activity.RegisterWithOptions(bootstrapActivity.Execute, activity.RegisterOptions{Name: eksworkflow.BootstrapActivityName})

	return nil
}