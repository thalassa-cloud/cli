package contextstate

import (
	"context"
	"fmt"
	"net/url"
)

func Login(ctx context.Context, token string) error {
	context, err := GetContextConfiguration()
	if err != nil {
		return err
	}
	context.Users.User.Token = token
	if err := CombineConfigContext(context); err != nil {
		return err
	}
	return Save()
}

func LoginWithAPIEndpoint(ctx context.Context, token, apiEndpoint string) error {
	context, err := GetContextConfiguration()
	if err != nil {
		return err
	}
	// validate api endpoint
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return fmt.Errorf("invalid api endpoint: %w", err)
	}

	context.Users.User.Token = token
	context.Servers.API.Server = u.String()
	if err := CombineConfigContext(context); err != nil {
		return err
	}
	return Save()
}

func Clear() error {
	context, err := GetContextConfiguration()
	if err != nil {
		return err
	}
	return globalConfigManager.RemoveContext(context.Name)
}

func Logout() error {
	context, err := GetContextConfiguration()
	if err != nil {
		return err
	}
	context.Users.User.Token = ""
	if err := CombineConfigContext(context); err != nil {
		return err
	}
	return Save()
}
