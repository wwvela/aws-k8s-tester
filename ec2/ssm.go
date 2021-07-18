package ec2

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-k8s-tester/ec2config"
	"github.com/aws/aws-k8s-tester/pkg/randutil"
	aws_v2 "github.com/aws/aws-sdk-go-v2/aws"
	aws_ssm_v2 "github.com/aws/aws-sdk-go-v2/service/ssm"
	aws_ssm_v2_types "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	smithy "github.com/aws/smithy-go"
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
)

// NOT WORKING...
// invalid document content

func (ts *Tester) createSSM() error {
	if err := ts.createSSMDocument(); err != nil {
		return err
	}
	if err := ts.sendSSMDocumentCommand(); err != nil {
		return err
	}
	return nil
}

func (ts *Tester) deleteSSM() error {
	if err := ts.deleteSSMDocument(); err != nil {
		return err
	}
	return nil
}

func (ts *Tester) createSSMDocument() error {
	createStart := time.Now()

	for asgName, cur := range ts.cfg.ASGs {
		if cur.SSM == nil {
			continue
		}

		if !cur.SSM.DocumentCreate {
			ts.lg.Info("skipping SSM document create",
				zap.String("asg-name", asgName),
				zap.String("ssm-document-name", cur.SSM.DocumentName),
			)
			continue
		}

		// ref. https://docs.aws.amazon.com/systems-manager/latest/userguide/create-ssm-document-api.html
		content := `schemaVersion: '2.2'
description: SSM document
parameters:
  region:
    type: String
	description: 'AWS Region'
  executionTimeoutSeconds:
    type: String
	description: 'timeout for script, in seconds'
  moreCommands:
    type: String
	description: 'more commands'
mainSteps:
  - action: aws:runShellScript
    name: %s
    inputs:
      timeoutSeconds: '{{ executionTimeoutSeconds }}'
      runCommand:
        - |
          AWS_DEFAULT_REGION={{region}}
          echo "running SSM with AWS_DEFAULT_REGION: ${AWS_DEFAULT_REGION}"
          echo "running more SSM command"
          {{ moreCommands }}
`

		ts.lg.Info("creating SSM document",
			zap.String("asg-name", asgName),
			zap.String("ssm-document-name", cur.SSM.DocumentName),
		)
		_, err := ts.ssmAPIV2.CreateDocument(
			context.Background(),
			&aws_ssm_v2.CreateDocumentInput{
				Name:           aws_v2.String(cur.SSM.DocumentName),
				DocumentFormat: aws_ssm_v2_types.DocumentFormatYaml,
				DocumentType:   aws_ssm_v2_types.DocumentTypeCommand,
				VersionName:    aws_v2.String("v1"),
				Tags: []aws_ssm_v2_types.Tag{
					{
						Key:   aws_v2.String("Name"),
						Value: aws_v2.String(ts.cfg.Name),
					},
					{
						Key:   aws_v2.String("DocumentName"),
						Value: aws_v2.String(cur.SSM.DocumentName),
					},
					{
						Key:   aws_v2.String("DocumentVersion"),
						Value: aws_v2.String("v1"),
					},
				},
				// ref. https://docs.aws.amazon.com/systems-manager/latest/userguide/create-ssm-document-api.html
				Content: aws_v2.String(fmt.Sprintf(content, cur.SSM.DocumentName)),
			},
		)
		if err != nil {
			return err
		}

		ts.lg.Info("created SSM Document",
			zap.String("asg-name", cur.Name),
			zap.String("ssm-document-name", cur.SSM.DocumentName),
			zap.String("started", humanize.RelTime(createStart, time.Now(), "ago", "from now")),
		)
	}

	ts.cfg.Sync()
	return nil
}

func (ts *Tester) deleteSSMDocument() (err error) {
	for asgName, cur := range ts.cfg.ASGs {
		if cur.SSM == nil {
			continue
		}

		if !cur.SSM.DocumentCreate {
			ts.lg.Info("skipping SSM document delete",
				zap.String("asg-name", asgName),
				zap.String("ssm-document-name", cur.SSM.DocumentName),
			)
			continue
		}
		ts.lg.Info("deleting SSM document",
			zap.String("asg-name", cur.Name),
			zap.String("ssm-document-name", cur.SSM.DocumentName),
		)
		_, err = ts.ssmAPIV2.DeleteDocument(
			context.Background(),
			&aws_ssm_v2.DeleteDocumentInput{
				Name:  aws_v2.String(cur.SSM.DocumentName),
				Force: true,
			},
		)
		if err != nil {
			ts.lg.Warn("failed to delete SSM document", zap.Error(err))
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				if strings.Contains(apiErr.ErrorCode(), "NotFound") {
					ts.cfg.DeletedResources[cur.SSM.DocumentName] = "SSM.DocumentName"
					ts.cfg.Sync()
					err = nil
				}
			}
			// InvalidDocument: Document eks2021071804awseyzymhjfdInstallBottlerocket does not exist in your account
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					ts.cfg.DeletedResources[cur.SSM.DocumentName] = "SSM.DocumentName"
					ts.cfg.Sync()
					err = nil
				}
			}
		} else {
			ts.cfg.DeletedResources[cur.SSM.DocumentName] = "SSM.DocumentName"
			ts.cfg.Sync()
		}
		if err == nil {
			ts.cfg.RecordStatus(fmt.Sprintf("%q/%s", cur.SSM.DocumentName, ec2config.StatusDELETEDORNOTEXIST))
		}

		ts.lg.Info("deleted SSM document",
			zap.String("asg-name", cur.Name),
			zap.String("ssm-document-name", cur.SSM.DocumentName),
		)
	}

	ts.cfg.Sync()
	return err
}

func (ts *Tester) sendSSMDocumentCommand() error {
	for asgName, cur := range ts.cfg.ASGs {
		if cur.SSM == nil {
			continue
		}

		if cur.SSM.DocumentName == "" {
			ts.lg.Info("skipping SSM document send",
				zap.String("asg-name", asgName),
				zap.String("ssm-document-name", cur.SSM.DocumentName),
			)
			continue
		}
		if len(cur.Instances) == 0 {
			return fmt.Errorf("no instance found for SSM document %q", cur.SSM.DocumentName)
		}
		ids := make([]string, 0)
		for id := range cur.Instances {
			ids = append(ids, id)
		}

		// batch by 50
		// e.g. 'instanceIds' failed to satisfy constraint: Member must have length less than or equal to 50
		ts.lg.Info("sending SSM document",
			zap.String("asg-name", asgName),
			zap.String("ssm-document-name", cur.SSM.DocumentName),
			zap.Int("instance-ids", len(ids)),
		)
		left := make([]string, len(ids))
		copy(left, ids)
		for len(left) > 0 {
			batch := make([]string, 0)
			switch {
			case len(left) <= 50:
				batch = append(batch, left...)
				left = left[:0:0]
			case len(left) > 50:
				batch = append(batch, left[:50]...)
				left = left[50:]
			}
			ssmInput := &aws_ssm_v2.SendCommandInput{
				DocumentName:   aws_v2.String(cur.SSM.DocumentName),
				Comment:        aws_v2.String(cur.SSM.DocumentName + "-" + randutil.String(10)),
				InstanceIds:    batch,
				MaxConcurrency: aws_v2.String(fmt.Sprintf("%d", len(batch))),
				Parameters: map[string][]string{
					"region":                  {ts.cfg.Region},
					"executionTimeoutSeconds": {fmt.Sprintf("%d", cur.SSM.DocumentExecutionTimeoutSeconds)},
				},
				OutputS3BucketName: aws_v2.String(ts.cfg.S3.BucketName),
				OutputS3KeyPrefix:  aws_v2.String(path.Join(ts.cfg.Name, "ssm-outputs")),
			}
			if len(cur.SSM.DocumentCommands) > 0 {
				ssmInput.Parameters["moreCommands"] = []string{cur.SSM.DocumentCommands}
			}
			cmd, err := ts.ssmAPIV2.SendCommand(
				context.Background(),
				ssmInput,
			)
			if err != nil {
				return err
			}
			docName := aws_v2.ToString(cmd.Command.DocumentName)
			if docName != cur.SSM.DocumentName {
				return fmt.Errorf("SSM Document Name expected %q, got %q", cur.SSM.DocumentName, docName)
			}
			cmdID := aws_v2.ToString(cmd.Command.CommandId)
			cur.SSM.DocumentCommandIDs = append(cur.SSM.DocumentCommandIDs, cmdID)

			ts.lg.Info("sent SSM document",
				zap.String("asg-name", asgName),
				zap.String("ssm-document-name", cur.SSM.DocumentName),
				zap.String("ssm-command-id", cmdID),
				zap.Int("sent-instance-ids", len(batch)),
				zap.Int("left-instance-ids", len(left)),
			)
			if len(left) == 0 {
				break
			}

			ts.lg.Info("waiting for next SSM run batch", zap.Int("left", len(left)))
			time.Sleep(15 * time.Second)
		}

		ts.cfg.ASGs[asgName] = cur
		ts.cfg.Sync()
	}

	ts.cfg.Sync()
	return nil
}
