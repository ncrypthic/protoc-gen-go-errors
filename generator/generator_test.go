package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	fds   *descriptorpb.FileDescriptorSet = &descriptorpb.FileDescriptorSet{}
	files *protoregistry.Files            = &protoregistry.Files{}
	p     *protogen.Plugin                = &protogen.Plugin{}
)

func TestMain(m *testing.M) {
	os.RemoveAll("./test/test.pb.descriptor")
	cmd := exec.Command("protoc", "-o", "test/test.pb.descriptor", "--include_imports", "-I", "../../", "-I", "test", "test.proto")
	err := cmd.Start()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(cmd.Args, " ")))
	}
	err = cmd.Wait()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(cmd.Args, " ")))
	}
	gengocmd := exec.Command("protoc", "--go_out=:test", "--go_opt=paths=source_relative", "-I", "../../", "-I", "test", "test.proto")
	err = gengocmd.Start()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(gengocmd.Args, " ")))
	}
	err = gengocmd.Wait()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(gengocmd.Args, " ")))
	}
	b, err := ioutil.ReadFile("./test/test.pb.descriptor")
	if err != nil {
		panic(fmt.Errorf("failed to parse proto descriptor: %w", err))
	}
	err = proto.Unmarshal(b, fds)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshall proto descriptor: %w", err))
	}
	files, err = protodesc.NewFiles(fds)
	if err != nil {
		panic(fmt.Errorf("failed to read FileDescriptorSet: %w", err))
	}
	p, err = protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{
		ProtoFile: fds.File,
	})
	if err != nil {
		panic(fmt.Errorf("failed to generate plugin: %w", err))
	}

	os.Exit(m.Run())
}

func TestGenerateError(t *testing.T) {
	g := p.NewGeneratedFile("t", "")
	NewGenerator(g).generateErrors(p.FilesByPath["test.proto"])
	txt, err := g.Content()
	require.Nil(t, err)
	fmt.Println(string(txt))
}
