package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"

	"github.com/aws/aws-sdk-go-v2/service/appstream"
)

type AppStreamImage struct {
	svc     *appstream.Client
	context context.Context

	name       *string
	visibility string

	sharedAccounts []*string
}

func init() {
	registerV2("AppStreamImage", ListAppStreamImages)
}

func ListAppStreamImages(cfg *aws.Config) ([]Resource, error) {
	svc := appstream.NewFromConfig(*cfg)
	ctx := context.TODO()

	resources := []Resource{}
	var nextToken *string

	for ok := true; ok; ok = (nextToken != nil) {
		params := &appstream.DescribeImagesInput{
			NextToken: nextToken,
		}

		output, err := svc.DescribeImages(ctx, params)
		if err != nil {
			return nil, err
		}
		nextToken = output.NextToken

		for _, image := range output.Images {
			sharedAccounts := []*string{}
			visibility := string(image.Visibility)

			// Filter out public images
			if strings.ToUpper(visibility) != "PUBLIC" {
				imagePerms, err := svc.DescribeImagePermissions(ctx, &appstream.DescribeImagePermissionsInput{
					Name: image.Name,
				})

				if err != nil {
					return nil, err
				}

				for _, permission := range imagePerms.SharedImagePermissionsList {
					sharedAccounts = append(sharedAccounts, permission.SharedAccountId)
				}

				resources = append(resources, &AppStreamImage{
					svc:            svc,
					context:        ctx,
					name:           image.Name,
					visibility:     visibility,
					sharedAccounts: sharedAccounts,
				})
			}
		}

	}

	return resources, nil
}

func (f *AppStreamImage) Remove() error {
	for _, account := range f.sharedAccounts {
		_, err := f.svc.DeleteImagePermissions(f.context, &appstream.DeleteImagePermissionsInput{
			Name:            f.name,
			SharedAccountId: account,
		})
		if err != nil {
			fmt.Println("Error deleting image permissions", err)
			return err
		}
	}

	_, err := f.svc.DeleteImage(f.context, &appstream.DeleteImageInput{
		Name: f.name,
	})
	if err != nil {
		fmt.Println("Error deleting image", err)
	}

	return err
}

func (f *AppStreamImage) String() string {
	return *f.name
}

func (f *AppStreamImage) Filter() error {
	if strings.ToUpper(f.visibility) == "PUBLIC" {
		return fmt.Errorf("cannot delete public AWS images")
	}
	return nil
}

func (f *AppStreamImage) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("Name", f.name)

	sharedAccounts := make([]string, len(f.sharedAccounts))
	for i, account := range f.sharedAccounts {
		sharedAccounts[i] = *account
	}
	if len(sharedAccounts) > 0 {
		properties.Set("Accounts with shared image", strings.Join(sharedAccounts, ", "))
	}

	return properties
}
