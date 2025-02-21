package gofiberswagger

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type RouteInfo = openapi3.Operation

var AcquiredRoutesInfo map[string]*RouteInfo

func RegisterPath(method string, path string, info *RouteInfo) {
	if AcquiredRoutesInfo == nil {
		AcquiredRoutesInfo = make(map[string]*RouteInfo)
	}
	if info == nil {
		info = &RouteInfo{}
	}
	AcquiredRoutesInfo[getAcquiredRoutesInfoId(method, path)] = info
}

func getAcquiredRoutesInfo(method string, path string) *RouteInfo {
	if AcquiredRoutesInfo == nil {
		return nil
	}
	return AcquiredRoutesInfo[getAcquiredRoutesInfoId(method, path)]
}

func getAcquiredRoutesInfoId(method string, path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(method)+path, " ", ""), "//", "/")
}
