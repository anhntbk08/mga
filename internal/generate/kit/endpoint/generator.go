package endpoint

import (
	"bytes"
	"fmt"

	"github.com/dave/jennifer/jen"
)

// Generate generates a go-kit endpoint struct.
func Generate(pkg string, spec ServiceSpec, withOc bool) (string, error) {
	file := jen.NewFilePath(pkg)

	file.PackageComment("Code generated by mga tool. DO NOT EDIT.")

	file.ImportName("github.com/go-kit/kit/endpoint", "endpoint")
	file.ImportAlias("github.com/sagikazarmark/kitx/endpoint", "kitxendpoint")
	file.ImportName(spec.Package.Path, spec.Package.Name)

	endpoints := make([]jen.Code, 0, len(spec.Endpoints))
	endpointDict := jen.Dict{}
	for _, endpoint := range spec.Endpoints {
		endpoints = append(endpoints, jen.Id(endpoint.Name).Qual("github.com/go-kit/kit/endpoint", "Endpoint"))
		endpointDict[jen.Id(endpoint.Name)] = jen.Id("mw").Call(
			jen.Id(fmt.Sprintf("Make%sEndpoint", endpoint.Name)).Call(jen.Id("service")),
		)
	}

	file.Comment("Endpoints collects all of the endpoints that compose the underlying service. It's")
	file.Comment("meant to be used as a helper struct, to collect all of the endpoints into a")
	file.Comment("single parameter.")
	file.Type().Id("Endpoints").Struct(endpoints...)

	file.Comment("MakeEndpoints returns an Endpoints struct where each endpoint invokes")
	file.Comment("the corresponding method on the provided service.")
	file.Func().Id("MakeEndpoints").
		Params(
			jen.Id("service").Qual(spec.Package.Path, spec.Name),
			jen.Id("middleware").Op("...").Qual("github.com/go-kit/kit/endpoint", "Middleware"),
		).
		Params(jen.Id("Endpoints")).
		Block(
			jen.Id("mw").
				Op(":=").
				Qual("github.com/sagikazarmark/kitx/endpoint", "Chain").
				Call(jen.Id("middleware").Op("...")),
			jen.Line(),
			jen.Return(
				jen.Id("Endpoints").Values(endpointDict),
			),
		)

	if withOc {
		file.ImportAlias("github.com/go-kit/kit/tracing/opencensus", "kitoc")

		endpointDict := jen.Dict{}
		for _, endpoint := range spec.Endpoints {
			endpointDict[jen.Id(endpoint.Name)] = jen.Qual(
				"github.com/go-kit/kit/tracing/opencensus",
				"TraceEndpoint",
			).
				Call(jen.Lit(fmt.Sprintf("%s.%s", spec.Package.Name, endpoint.Name))).
				Call(jen.Id("endpoints").Dot(endpoint.Name))
		}

		file.Comment("TraceEndpoints returns an Endpoints struct where each endpoint is wrapped with a tracing middleware.")
		file.Func().Id("TraceEndpoints").
			Params(jen.Id("endpoints").Id("Endpoints")).
			Params(jen.Id("Endpoints")).
			Block(jen.Return(jen.Id("Endpoints").Values(endpointDict)))
	}

	var buf bytes.Buffer

	err := file.Render(&buf)

	return buf.String(), err
}