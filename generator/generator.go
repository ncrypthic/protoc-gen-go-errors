package generator

import (
	"flag"

	pb "github.com/ncrypthic/protoc-gen-go-errors/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	*protogen.GeneratedFile
	indent []string
}

func NewGenerator(g *protogen.GeneratedFile) *generator {
	return &generator{g, []string{}}
}

func (g *generator) Enter() {
	g.indent = append(g.indent, "\t")
}

func (g *generator) Exit() {
	g.indent = g.indent[1:]
}

func (g *generator) P(args ...interface{}) {
	payload := make([]interface{}, len(g.indent)+len(args))
	counter := 0
	for _, i := range g.indent {
		payload[counter] = i
		counter++
	}
	for _, a := range args {
		payload[counter] = a
		counter++
	}
	g.GeneratedFile.P(payload...)
}

func (g *generator) generateErr(err *pb.Error) {
	statusNew := protogen.GoIdent{
		GoImportPath: "google.golang.org/grpc/status",
		GoName:       "New",
	}
	grpcCode := "Unknown"
	switch err.GetStatus() {
	case pb.Status_OK:
		grpcCode = codes.OK.String()
	case pb.Status_CANCELED:
		grpcCode = codes.Canceled.String()
	case pb.Status_UNKNOWN:
		grpcCode = codes.Unknown.String()
	case pb.Status_INVALID_ARGUMENT:
		grpcCode = codes.InvalidArgument.String()
	case pb.Status_DEADLINE_EXCEEDED:
		grpcCode = codes.DeadlineExceeded.String()
	case pb.Status_NOT_FOUND:
		grpcCode = codes.NotFound.String()
	case pb.Status_ALREADY_EXISTS:
		grpcCode = codes.AlreadyExists.String()
	case pb.Status_PERMISSION_DENIED:
		grpcCode = codes.PermissionDenied.String()
	case pb.Status_RESOURCE_EXHAUSTED:
		grpcCode = codes.ResourceExhausted.String()
	case pb.Status_FAILED_PRECONDITION:
		grpcCode = codes.FailedPrecondition.String()
	case pb.Status_ABORTED:
		grpcCode = codes.Aborted.String()
	case pb.Status_OUT_OF_RANGE:
		grpcCode = codes.OutOfRange.String()
	case pb.Status_UNIMPLEMENTED:
		grpcCode = codes.Unimplemented.String()
	case pb.Status_INTERNAL:
		grpcCode = codes.Internal.String()
	case pb.Status_UNAVAILABLE:
		grpcCode = codes.Unavailable.String()
	default:
		grpcCode = codes.Unknown.String()
	}
	codes := protogen.GoIdent{
		GoImportPath: "google.golang.org/grpc/codes",
		GoName:       grpcCode,
	}
	g.P("ERR_", err.GetCode(), " error", " = ", statusNew, "(", codes, ", \"", err.GetCode(), "\").Err()")
}

func (g *generator) generateErrors(f *protogen.File) {
	g.P("var (")
	g.Enter()
	for _, svc := range f.Services {
		for _, m := range svc.Methods {
			opt := proto.GetExtension(m.Desc.Options(), pb.E_Types)
			if opt == nil {
				continue
			}
			defs, ok := opt.([]*pb.Error)
			if !ok {
				continue
			}
			for _, e := range defs {
				g.generateErr(e)
			}
		}
	}
	g.Exit()
	g.P(")")
}

func Generate() {
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			filename := f.GeneratedFilenamePrefix + ".errors.pb.go"
			file := gen.NewGeneratedFile(filename, f.GoImportPath)
			file.P("package ", f.GoPackageName)
			g := NewGenerator(file)
			g.generateErrors(f)
		}
		return nil
	})
}
