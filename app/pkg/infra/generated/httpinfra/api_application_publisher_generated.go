// Code generated by Signare OpenAPI generator. DO NOT EDIT

package httpinfra

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

// ApplicationAPIPublisherOptions options to create a ApplicationAPIRoutesPublished
type ApplicationAPIPublisherOptions struct {
	HTTPInfra httpinfra.HTTPRouter
	Handler   ApplicationAPIHTTPHandler
}

// ApplicationAPIRoutesPublished type for ApplicationAPI published routes
type ApplicationAPIRoutesPublished int

// ProvideApplicationAPIRoutes creates a new ApplicationAPIRoutesPublished
func ProvideApplicationAPIRoutes(options ApplicationAPIPublisherOptions) (ApplicationAPIRoutesPublished, error) {

	if options.HTTPInfra == nil {
		return 0, errors.New("missing mandatory HTTPInfra field to publish ApplicationAPI Routes")
	}

	if options.Handler == nil {
		return 0, errors.New("missing mandatory field Handler to publish ApplicationAPI Routes")
	}

	var err error

	err = PublishApplicationAccountsCreate(options.HTTPInfra, options.Handler)
	if err != nil {
		return 0, err
	}
	err = PublishApplicationAccountsRemove(options.HTTPInfra, options.Handler)
	if err != nil {
		return 0, err
	}
	err = PublishApplicationUsersCreate(options.HTTPInfra, options.Handler)
	if err != nil {
		return 0, err
	}
	err = PublishApplicationUsersDescribe(options.HTTPInfra, options.Handler)
	if err != nil {
		return 0, err
	}
	err = PublishApplicationUsersEdit(options.HTTPInfra, options.Handler)
	if err != nil {
		return 0, err
	}
	err = PublishApplicationUsersList(options.HTTPInfra, options.Handler)
	if err != nil {
		return 0, err
	}
	err = PublishApplicationUsersRemove(options.HTTPInfra, options.Handler)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

// PublishApplicationAccountsCreate publishes the ApplicationAccountsCreate endpoint
func PublishApplicationAccountsCreate(httpInfra httpinfra.HTTPRouter, handler ApplicationAPIHTTPHandler) error {
	opts := httpinfra.HandlerMatchOptions{Path: "/applications/{applicationId}/users/{userId}/accounts", Methods: []string{
		http.MethodPost,
	},
		Action: "application.accounts.create",
	}
	err := httpInfra.RegisterRawHandler(opts, handler.HandleHTTPApplicationAccountsCreate)
	if err != nil {
		return err
	}
	return nil
}

// PublishApplicationAccountsRemove publishes the ApplicationAccountsRemove endpoint
func PublishApplicationAccountsRemove(httpInfra httpinfra.HTTPRouter, handler ApplicationAPIHTTPHandler) error {
	opts := httpinfra.HandlerMatchOptions{Path: "/applications/{applicationId}/users/{userId}/accounts/{accountId}", Methods: []string{
		http.MethodDelete,
	},
		Action: "application.accounts.remove",
	}
	err := httpInfra.RegisterRawHandler(opts, handler.HandleHTTPApplicationAccountsRemove)
	if err != nil {
		return err
	}
	return nil
}

// PublishApplicationUsersCreate publishes the ApplicationUsersCreate endpoint
func PublishApplicationUsersCreate(httpInfra httpinfra.HTTPRouter, handler ApplicationAPIHTTPHandler) error {
	opts := httpinfra.HandlerMatchOptions{Path: "/applications/{applicationId}/users", Methods: []string{
		http.MethodPost,
	},
		Action: "application.users.create",
	}
	err := httpInfra.RegisterRawHandler(opts, handler.HandleHTTPApplicationUsersCreate)
	if err != nil {
		return err
	}
	return nil
}

// PublishApplicationUsersDescribe publishes the ApplicationUsersDescribe endpoint
func PublishApplicationUsersDescribe(httpInfra httpinfra.HTTPRouter, handler ApplicationAPIHTTPHandler) error {
	opts := httpinfra.HandlerMatchOptions{Path: "/applications/{applicationId}/users/{userId}", Methods: []string{
		http.MethodGet,
	},
		Action: "application.users.describe",
	}
	err := httpInfra.RegisterRawHandler(opts, handler.HandleHTTPApplicationUsersDescribe)
	if err != nil {
		return err
	}
	return nil
}

// PublishApplicationUsersEdit publishes the ApplicationUsersEdit endpoint
func PublishApplicationUsersEdit(httpInfra httpinfra.HTTPRouter, handler ApplicationAPIHTTPHandler) error {
	opts := httpinfra.HandlerMatchOptions{Path: "/applications/{applicationId}/users/{userId}", Methods: []string{
		http.MethodPut,
	},
		Action: "application.users.edit",
	}
	err := httpInfra.RegisterRawHandler(opts, handler.HandleHTTPApplicationUsersEdit)
	if err != nil {
		return err
	}
	return nil
}

// PublishApplicationUsersList publishes the ApplicationUsersList endpoint
func PublishApplicationUsersList(httpInfra httpinfra.HTTPRouter, handler ApplicationAPIHTTPHandler) error {
	opts := httpinfra.HandlerMatchOptions{Path: "/applications/{applicationId}/users", Methods: []string{
		http.MethodGet,
	},
		Action: "application.users.list",
	}
	err := httpInfra.RegisterRawHandler(opts, handler.HandleHTTPApplicationUsersList)
	if err != nil {
		return err
	}
	return nil
}

// PublishApplicationUsersRemove publishes the ApplicationUsersRemove endpoint
func PublishApplicationUsersRemove(httpInfra httpinfra.HTTPRouter, handler ApplicationAPIHTTPHandler) error {
	opts := httpinfra.HandlerMatchOptions{Path: "/applications/{applicationId}/users/{userId}", Methods: []string{
		http.MethodDelete,
	},
		Action: "application.users.remove",
	}
	err := httpInfra.RegisterRawHandler(opts, handler.HandleHTTPApplicationUsersRemove)
	if err != nil {
		return err
	}
	return nil
}
