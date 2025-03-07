package work_packages

import (
	"strconv"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type FilterOption int

const (
	Assignee FilterOption = iota
	Version
	Project
)

func Lookup(id uint64) (*models.WorkPackage, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	return workPackage.Convert(), nil
}

func All(filterOptions *map[FilterOption]string) ([]*models.WorkPackage, error) {
	var filters []requests.Filter
	var projectId *uint64

	for updateOpt, value := range *filterOptions {
		switch updateOpt {
		case Assignee:
			filters = append(filters, AssigneeFilter(value))
		case Version:
			filters = append(filters, VersionFilter(value))
		case Project:
			n, _ := strconv.ParseUint(value, 10, 64)
			projectId = &n
		}
	}

	query := requests.NewQuery(filters)

	requestUrl := paths.WorkPackages()

	if projectId != nil {
		requestUrl = paths.ProjectWorkPackages(*projectId)
	}

	response, err := requests.Get(requestUrl, &query)
	if err != nil {
		return nil, err
	}

	workPackageCollection := parser.Parse[dtos.WorkPackageCollectionDto](response)
	return workPackageCollection.Convert(), nil
}

func fetch(id uint64) (*dtos.WorkPackageDto, error) {
	response, err := requests.Get(paths.WorkPackage(id), nil)
	if err != nil {
		return nil, err
	}

	workPackage := parser.Parse[dtos.WorkPackageDto](response)
	return &workPackage, nil
}
