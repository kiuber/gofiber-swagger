package gofiberswagger

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type RouteInfo = openapi3.Operation

var acquiredRoutesInfo map[string]*RouteInfo

func RegisterPath(method string, path string, info *RouteInfo) {
	if acquiredRoutesInfo == nil {
		acquiredRoutesInfo = make(map[string]*RouteInfo)
	}
	if info == nil {
		info = &RouteInfo{}
	}
	acquiredRoutesInfo[getAcquiredRoutesInfoId(method, path)] = info
}

func getAcquiredRoutesInfo(method string, path string) *RouteInfo {
	if acquiredRoutesInfo == nil {
		return nil
	}
	return acquiredRoutesInfo[getAcquiredRoutesInfoId(method, path)]
}

func getAcquiredRoutesInfoId(method string, path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(method)+path, " ", ""), "//", "/")
}
