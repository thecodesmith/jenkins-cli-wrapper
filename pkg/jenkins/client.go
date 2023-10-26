package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/bndr/gojenkins"

	config "github.com/thecodesmith/jenkinsw/pkg/config"
	"github.com/thecodesmith/jenkinsw/pkg/utils"
)

type Client struct {
	ioStreams *utils.IOStreams
	api       *gojenkins.Jenkins
	ctx       context.Context
}

func NewClient(ctx *config.Context, streams *utils.IOStreams) (*Client, error) {
	httpClient := http.DefaultClient
	httpClient.Transport = http.DefaultTransport
	httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}

	httpCtx := context.TODO()
	jenkins := gojenkins.CreateJenkins(nil, ctx.Host, ctx.Username, ctx.ApiToken)
	_, err := jenkins.Init(httpCtx)

	if err != nil {
		return nil, err
	}

	return &Client{api: jenkins, ctx: httpCtx, ioStreams: streams}, nil
}

func (c *Client) Version() string {
	return c.api.Version
}

// ListJobs details
func (c *Client) ListJobs(depth int) ([]gojenkins.InnerJob, error) {
	immediateJobs, err := c.api.GetAllJobNames(c.ctx)
	if err != nil {
		return nil, err
	}
	if depth == 0 {
		return immediateJobs, nil
	}
	finalJobs := make([]gojenkins.InnerJob, 0)
	for _, jobRaw := range immediateJobs {
		if jobRaw.Class == "com.cloudbees.hudson.plugins.folder.Folder" {
			job, err := c.api.GetJob(c.ctx, jobRaw.Name)
			if err != nil {
				return nil, err
			}
			receivedJobs, err := c.getInnerJobs(job.GetName(), job, 0, depth)
			if err != nil {
				return nil, err
			}
			if len(receivedJobs) > 0 {
				finalJobs = append(finalJobs, receivedJobs...)
			}
		} else {
			finalJobs = append(finalJobs, jobRaw)
		}
	}
	return finalJobs, nil
}

func (jc *Client) getInnerJobs(parent string, job *gojenkins.Job, depth, limit int) ([]gojenkins.InnerJob, error) {
	if depth == limit {
		if job == nil {
			return nil, nil
		}
		return []gojenkins.InnerJob{{
			Color: job.GetDetails().Color, Name: parent,
			Class: job.GetDetails().Class, Url: job.GetDetails().URL}}, nil
	}
	finalJobs := make([]gojenkins.InnerJob, 0)

	jobsList, err := job.GetInnerJobs(jc.ctx)
	if err != nil {
		return nil, err
	}
	for _, jobInner := range jobsList {
		jobName := fmt.Sprintf("%s/job/%s", parent, jobInner.GetName())
		if jobInner.GetDetails().Class == "com.cloudbees.hudson.plugins.folder.Folder" {
			receivedJobs, err := jc.getInnerJobs(jobName, jobInner, depth+1, limit)
			if err != nil {
				return nil, err
			}
			if len(receivedJobs) > 0 {
				finalJobs = append(finalJobs, receivedJobs...)
			}
		} else {
			finalJobs = append(finalJobs, gojenkins.InnerJob{
				Color: jobInner.GetDetails().Color, Name: jobName,
				Class: jobInner.GetDetails().Class, Url: jobInner.GetDetails().URL})
		}
	}
	return finalJobs, nil
}
