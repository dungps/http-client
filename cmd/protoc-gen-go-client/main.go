package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const release = "1.0.0"

var (
	showVersion = flag.Bool("version", false, "print version")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-client %v\n", release)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if f.Generate {
				generateFile(gen, f)
			}
		}

		return nil
	})
}
