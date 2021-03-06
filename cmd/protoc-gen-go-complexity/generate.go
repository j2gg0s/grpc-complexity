package main

import (
	"fmt"
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	contextPackage    = protogen.GoImportPath("context")
	grpcPackage       = protogen.GoImportPath("google.golang.org/grpc")
	codesPackage      = protogen.GoImportPath("google.golang.org/grpc/codes")
	statusPackage     = protogen.GoImportPath("google.golang.org/grpc/status")
	complexityPackage = protogen.GoImportPath("github.com/j2gg0s/grpc-complexity/complexity")
)

func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_complexity.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-complexity. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	for _, service := range file.Services {
		g.P("type " + service.GoName + "ComplexityServer interface{")
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				continue
			}
			g.P(methodSignature(g, method))
		}
		g.P("}")

		// Default Server
		g.P("type Default" + service.GoName + "ComplexityServer struct{}")
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				continue
			}
			g.P("func(s *Default" + service.GoName + "ComplexityServer) " + methodSignature(g, method) + "{")
			g.P(`return map[string]int{"default":1}`)
			g.P("}")
			g.P()
		}

		serverType := service.GoName + "ComplexityServer"
		serviceDescVar := service.GoName + "Complexity_ServiceDesc"
		// Server registration
		g.P("func Register", service.GoName, "ComplexityServer(s ", complexityPackage.Ident("ServiceRegistrar"), ", srv ", serverType, ") {")
		g.P("s.RegisterService(&", serviceDescVar, ", srv)")
		g.P("}")
		g.P()

		handlerNames := make([]string, len(service.Methods))
		for i, method := range service.Methods {
			handlerNames[i] = genServerMethod(gen, file, g, method)
		}

		// Service descriptor.
		g.P("var ", serviceDescVar, " = ", complexityPackage.Ident("ServiceDesc"), " {")
		g.P("ServiceName: ", strconv.Quote(string(service.Desc.FullName())), ",")
		g.P("HandlerType: (*", serverType, ")(nil),")
		g.P("Methods: []", complexityPackage.Ident("MethodDesc"), "{")
		for i, method := range service.Methods {
			if method.Desc.IsStreamingClient() && method.Desc.IsStreamingServer() {
				continue
			}
			g.P("{")
			g.P("MethodName: ", strconv.Quote(string(method.Desc.Name())), ",")
			g.P("ComplexityHandler: ", handlerNames[i], ",")
			g.P("},")
		}
		g.P("},")
		g.P("Metadata: \"", file.Desc.Path(), "\",")
		g.P("}")
	}

	return g
}

func methodSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	s := method.GoName + "(ctx " + g.QualifiedGoIdent(contextPackage.Ident("Context"))
	s += ", in *" + g.QualifiedGoIdent(method.Input.GoIdent)
	s += ", opts ..." + g.QualifiedGoIdent(grpcPackage.Ident("CallOption")) + ") map[string]int"
	return s
}

func genServerMethod(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, method *protogen.Method) string {
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		return ""
	}
	service := method.Parent
	hname := fmt.Sprintf("_%s_%s_ComplexityHandler", service.GoName, method.GoName)

	g.P("func ", hname, "(srv interface{}, ctx ", contextPackage.Ident("Context"), ", req interface{}) map[string]int {")
	g.P("return srv.(", service.GoName, "ComplexityServer).", method.GoName, "(ctx, req.(*", method.Input.GoIdent, "))")
	g.P("}")
	g.P()

	return hname
}
