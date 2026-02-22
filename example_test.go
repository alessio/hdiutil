package hdiutil_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"al.essio.dev/pkg/hdiutil"
)

func ExampleConfig_Validate() {
	cfg := hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
		VolumeName: "My App",
	}

	err := cfg.Validate()
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleConfig_Validate_withFormat() {
	cfg := hdiutil.Config{
		SourceDir:   "./dist",
		OutputPath:  "MyApp.dmg",
		ImageFormat: "UDBZ",
		FileSystem:  "APFS",
	}

	err := cfg.Validate()
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleConfig_Validate_sandboxSafeAPFS() {
	cfg := hdiutil.Config{
		SourceDir:   "./dist",
		OutputPath:  "MyApp.dmg",
		SandboxSafe: true,
		FileSystem:  "APFS",
	}

	err := cfg.Validate()
	fmt.Println(errors.Is(err, hdiutil.ErrSandboxAPFS))
	// Output:
	// true
}

func ExampleConfig_Validate_invalidFormat() {
	cfg := hdiutil.Config{
		SourceDir:   "./dist",
		OutputPath:  "MyApp.dmg",
		ImageFormat: "INVALID",
	}

	err := cfg.Validate()
	fmt.Println(errors.Is(err, hdiutil.ErrInvFormatOpt))
	// Output:
	// true
}

func ExampleConfig_Validate_missingExtension() {
	cfg := hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.iso",
	}

	err := cfg.Validate()
	fmt.Println(errors.Is(err, hdiutil.ErrImageFileExt))
	// Output:
	// true
}

func ExampleConfig_Validate_unsafeArg() {
	cfg := hdiutil.Config{
		SourceDir:  "src\x00evil",
		OutputPath: "test.dmg",
	}

	err := cfg.Validate()
	fmt.Println(errors.Is(err, hdiutil.ErrUnsafeArg))
	// Output:
	// true
}

func ExampleConfig_VolumeNameOpt() {
	cfg := hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.VolumeNameOpt())
	// Output:
	// MyApp
}

func ExampleConfig_VolumeNameOpt_explicit() {
	cfg := hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
		VolumeName: "My Application",
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.VolumeNameOpt())
	// Output:
	// My Application
}

func ExampleConfig_FilesystemOpts() {
	cfg := hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
		FileSystem: "APFS",
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.FilesystemOpts())
	// Output:
	// [-fs APFS]
}

func ExampleConfig_ImageFormatOpts() {
	cfg := hdiutil.Config{
		SourceDir:   "./dist",
		OutputPath:  "MyApp.dmg",
		ImageFormat: "UDBZ",
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.ImageFormatOpts())
	// Output:
	// [-format UDBZ -imagekey bzip2-level=9]
}

func ExampleConfig_VolumeSizeOpts() {
	cfg := hdiutil.Config{
		SourceDir:    "./dist",
		OutputPath:   "MyApp.dmg",
		VolumeSizeMb: 256,
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.VolumeSizeOpts())
	// Output:
	// [-size 256m]
}

func ExampleConfig_VolumeSizeOpts_auto() {
	cfg := hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
	}

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.VolumeSizeOpts())
	// Output:
	// []
}

func ExampleConfig_ToJSON() {
	cfg := &hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
		VolumeName: "My App",
		FileSystem: "HFS+",
	}

	var buf bytes.Buffer
	if err := cfg.ToJSON(&buf); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(buf.String())
	// Output:
	// {
	//   "volume_name": "My App",
	//   "filesystem": "HFS+",
	//   "output_path": "MyApp.dmg",
	//   "source_dir": "./dist"
	// }
}

func ExampleConfig_FromJSON() {
	jsonStr := `{
		"source_dir": "./dist",
		"output_path": "MyApp.dmg",
		"volume_name": "My App",
		"image_format": "ULFO"
	}`

	cfg := &hdiutil.Config{}
	if err := cfg.FromJSON(strings.NewReader(jsonStr)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.SourceDir)
	fmt.Println(cfg.OutputPath)
	fmt.Println(cfg.VolumeName)
	fmt.Println(cfg.ImageFormat)
	// Output:
	// ./dist
	// MyApp.dmg
	// My App
	// ULFO
}

func ExampleLoadConfig() {
	dir, err := os.MkdirTemp("", "example-*")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "dmg.json")
	data := []byte(`{"source_dir":"./dist","output_path":"App.dmg","volume_name":"App"}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		fmt.Println(err)
		return
	}

	cfg, err := hdiutil.LoadConfig(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.VolumeName)
	fmt.Println(cfg.OutputPath)
	// Output:
	// App
	// App.dmg
}

func ExampleNew() {
	cfg := &hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
		VolumeName: "My App",
		Simulate:   true,
	}

	runner := hdiutil.New(cfg)
	defer runner.Cleanup()

	if err := runner.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if err := runner.Start(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok")
	// Output:
	// ok
}

func ExampleNew_sandboxSafe() {
	cfg := &hdiutil.Config{
		SourceDir:   "./dist",
		OutputPath:  "MyApp.dmg",
		SandboxSafe: true,
		Simulate:    true,
	}

	runner := hdiutil.New(cfg)
	defer runner.Cleanup()

	if err := runner.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if err := runner.Start(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok")
	// Output:
	// ok
}

func ExampleNew_fullWorkflow() {
	cfg := &hdiutil.Config{
		SourceDir:   "./dist",
		OutputPath:  "MyApp.dmg",
		VolumeName:  "My App",
		ImageFormat: "UDBZ",
		FileSystem:  "HFS+",
		Simulate:    true,
	}

	runner := hdiutil.New(cfg)
	defer runner.Cleanup()

	if err := runner.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if err := runner.Start(); err != nil {
		fmt.Println(err)
		return
	}

	if err := runner.FinalizeDMG(); err != nil {
		fmt.Println(err)
		return
	}

	// Codesign and Notarize are no-ops without credentials.
	if err := runner.Codesign(); err != nil {
		fmt.Println(err)
		return
	}

	if err := runner.Notarize(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok")
	// Output:
	// ok
}

func ExampleSetLogWriter() {
	var buf bytes.Buffer
	hdiutil.SetLogWriter(&buf)
	defer hdiutil.SetLogWriter(os.Stderr)

	cfg := &hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
		Simulate:   true,
	}

	runner := hdiutil.New(cfg)
	defer runner.Cleanup()

	_ = runner.Setup()
	_ = runner.Start()

	// Log output was captured in buf.
	fmt.Println(buf.Len() > 0)
	// Output:
	// true
}

func ExampleWithExecutor() {
	mock := &noopExecutor{}

	cfg := &hdiutil.Config{
		SourceDir:  "./dist",
		OutputPath: "MyApp.dmg",
		VolumeName: "My App",
	}

	runner := hdiutil.New(cfg, hdiutil.WithExecutor(mock))
	defer runner.Cleanup()

	if err := runner.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if err := runner.Start(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok")
	// Output:
	// ok
}