package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Ballerina Version. Value is set during the build process.
var balVersion string
var balPOS string
var balPArch string
var balZip string
var balDist string

//const resourceStore = "https://staging-cdn-updates.private.wso2.com/wum-build-resources/"
var resourceStore = "http://10.100.7.90:8000/"

var windowsData = map[string]string{
	"java.wis": `<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi"
xmlns:util="http://schemas.microsoft.com/wix/UtilExtension">
	</Wix>`,
	"installer.wxs": `<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi"
xmlns:util="http://schemas.microsoft.com/wix/UtilExtension">
<?if $(var.Arch) = 386 ?>
  <?define ProdId = {6147400c-24be-4f94-ba82-5a1c76320f10} ?>
  <?define UpgradeCode = {1008193c-0d2f-4c55-b7e2-5b3342ef042b} ?>
  <?define SysFolder=SystemFolder ?>
  <?define PlatformArch=x86 ?>
  <?define ProgramFilesDir=ProgramFilesFolder ?>
<?else?>
  <?define ProdId = {8ca1298f-2d74-4ca1-8f56-1d1147df5034} ?>
  <?define UpgradeCode = {7b6926c8-c11f-426a-94f9-66fdf933a682} ?>
  <?define SysFolder=System64Folder ?>
  <?define PlatformArch=x64 ?>
  <?define ProgramFilesDir=ProgramFiles64Folder ?>
<?endif?>

<Product
    Id="*"
    Name="Ballerina $(var.balVersion)"
    Language="1033"
    Version="$(var.WixbalVersion)"
    Manufacturer="https://wso2.com"
    UpgradeCode="$(var.UpgradeCode)" >
<Package
    Id='*'
    Keywords='Installer'
    Description="The Ballerina Installer"
    Comments="The Ballerina Installer."
    InstallerVersion="300"
    Compressed="yes"
    InstallScope="perMachine"
    Languages="1033"
    Platform="$(var.PlatformArch)" />

	<UI>
		<UIRef Id="WixUI_InstallDir"/>
		<Publish Dialog="WelcomeDlg" Control="Next" Event="NewDialog" Value="InstallDirDlg">NOT Installed</Publish>
		<Publish Dialog="InstallDirDlg" Control="Back" Event="NewDialog" Value="WelcomeDlg" Order="2">1</Publish>
		<Publish Dialog="InstallDirDlg" Control="Next" Event="NewDialog" Value="PrepareDlg" Order="5">WIXUI_DONTVALIDATEPATH OR WIXUI_INSTALLDIR_VALID="1"</Publish>
	</UI>

<Property Id="ARPCOMMENTS" Value="Ballerina is a general purpose, concurrent and strongly typed programming language with both textual and graphical syntaxes, optimized for integration." />
<Property Id="ARPCONTACT" Value="https://ballerinalang.org/" />
<Property Id="ARPHELPLINK" Value="https://ballerinalang.org/" />
<Property Id="ARPREADME" Value="https://ballerinalang.org/" />
<Property Id="ARPURLINFOABOUT" Value="https://ballerinalang.org/" />
<Media Id='1' Cabinet="bal.cab" EmbedCab="yes" CompressionLevel="high" />
<Condition Message="Windows XP or greater required."> VersionNT >= 500</Condition>
<MajorUpgrade AllowDowngrades="yes" />
<SetDirectory Id="INSTALLDIRROOT" Value="[WindowsVolume]Ballerina"/>
<CustomAction
    Id="SetApplicationRootDirectory"
    Property="ARPINSTALLLOCATION"
    Value="[INSTALLDIR]" />
<!-- Define the directory structure and environment variables -->
<Directory Id="TARGETDIR" Name="SourceDir">
  <Directory Id="INSTALLDIRROOT">
			<Directory Id="INSTALLDIR" Name="ballerina-$(var.balVersion)">
		</Directory>
  </Directory>
  <Directory Id="ProgramMenuFolder">
    <Directory Id="BallerinaProgramShortcutsDir" Name="Ballerina-$(var.balVersion)"/>
  </Directory>
  <Directory Id="EnvironmentEntries">
    <Directory Id="BallerinaEnvironmentEntries" Name="Ballerina-$(var.balVersion)"/>
  </Directory>
</Directory>
<!-- Programs Menu Shortcuts -->
<DirectoryRef Id="BallerinaProgramShortcutsDir">
  <Component Id="Component_BallerinaProgramShortCuts" Guid="{764ee6d4-917f-422c-87cb-cc0fff389765}">
    <Shortcut
        Id="UninstallShortcut"
        Name="Uninstall Ballerina-$(var.balVersion)"
        Description="Uninstalls Ballerina-$(var.balVersion)"
        Target="[$(var.SysFolder)]msiexec.exe"
        Arguments="/x [ProductCode]" />
    <RemoveFolder
        Id="BallerinaProgramShortcutsDir"
        On="uninstall" />
    <RegistryValue
        Root="HKCU"
        Key="Software\Ballerina-$(var.balVersion)"
        Name="ShortCuts"
        Type="integer"
        Value="1"
        KeyPath="yes" />
  </Component>
</DirectoryRef>
<!-- Registry & Environment Settings -->
<DirectoryRef Id="BallerinaEnvironmentEntries">
  <Component Id="Component_BallerinaEnvironment" Guid="{f9f2e5e9-d6fb-4ef3-8faf-38b5fd283237}">
    <RegistryKey
        Root="HKCU"
        Key="Software\Ballerina-$(var.balVersion)"
        Action="create" >
            <RegistryValue
                Name="installed"
                Type="integer"
                Value="1"
                KeyPath="yes" />
            <RegistryValue
                Name="installLocation"
                Type="string"
                Value="[INSTALLDIR]" />
    </RegistryKey>
    <Environment
        Id="BallerinaHome"
        Action="set"
        Part="all"
        Name="BALLERINA_HOME"
        Permanent="no"
        System="yes"
        Value="[INSTALLDIR]" />
	<Environment
        Id="BallerinaPathEntry"
        Action="set"
        Part="last"
        Name="PATH"
        Permanent="no"
        System="yes"
        Value="%BALLERINA_HOME%\bin" />
    <RemoveFolder
        Id="BallerinaEnvironmentEntries"
        On="uninstall" />
  </Component>
</DirectoryRef>
<!-- Install the files -->
<Feature
    Id="BallerinaTools"
    Title="Ballerina"
    Level="1">
      <ComponentRef Id="Component_BallerinaEnvironment" />
      <ComponentGroupRef Id="AppFiles" />
      <ComponentRef Id="Component_BallerinaProgramShortCuts" />
</Feature>
<!-- Update the environment -->
<InstallExecuteSequence>
    <Custom Action="SetApplicationRootDirectory" Before="InstallFinalize" />
</InstallExecuteSequence>
<!-- Include the user interface -->
<WixVariable Id="WixUIBannerBmp" Value="Banner.jpg" />
<WixVariable Id="WixUIDialogBmp" Value="Dialog.jpg" />
<Property Id="WIXUI_INSTALLDIR" Value="INSTALLDIR" />
</Product>

</Wix>
`,
	"Banner.jpg":  resourceStore + "windows/ballerina-banner.jpg",
	"Dialog.jpg":  resourceStore + "windows/ballerina-dialog.jpg",
}

var versionRe = regexp.MustCompile(`^Ballerina(\d+(\.\d+)*)`)

func main() {
	targetDir := ""
	unzip(balZip, targetDir)
	os.Rename(balDist, "ballerina-" + balVersion)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	switch runtime.GOOS {
	case "windows":
		err = windowsMSI()
	case "darwin":
	//Todo: for macOSX
	}
	if err != nil {
		log.Fatal(err)
	}
}

func windowsMSI() error {
	cwd, version, err := environmentInfo()
	if err != nil {
		return err
	}

	//Install Wix tools.
	wix := filepath.Join(cwd, "wix")
	defer os.RemoveAll(wix)
	if err := installWix(wix); err != nil {
		return err
	}

	// Write out windows data that is used by the packaging process.
	win := filepath.Join(cwd, "windows")

	defer os.RemoveAll(win)
	if err := writeDataFiles(windowsData, win); err != nil {
		return err
	}

	// Gather files.
	balDir := filepath.Join(cwd, "ballerina-" + balVersion)
	appfiles := filepath.Join(win, "AppFiles.wxs")
	if err := runDir(win, filepath.Join(wix, "heat"),
		"dir", balDir,
		"-nologo",
		"-gg", "-g1", "-srd", "-sfrag", "-sreg",
		"-cg", "AppFiles",
		"-template", "fragment",
		"-dr", "INSTALLDIR",
		"-var", "var.SourceDir",
		"-out", appfiles,
	); err != nil {
		return err
	}

	// Build package.
	if err := runDir(win, filepath.Join(wix, "candle"),
		"-nologo",
		"-dbalVersion=" + version,
		"-dWixbalVersion=" + wixVersion(version),
		"-dArch=" + runtime.GOARCH,
		"-dSourceDir=" + balDir,
		"-ext", "WixUtilExtension",
		filepath.Join(win, "installer.wxs"),
		appfiles,
	); err != nil {
		return err
	}

	msi := filepath.Join(cwd, "msi")

	if err := createDirIfNotExist(msi); err != nil {
		return err
	}

	return runDir(win, filepath.Join(wix, "light"),
		"-nologo",
		"-dcl:high",
		"-sice:ICE60",
		"-ext", "WixUIExtension",
		"-ext", "WixUtilExtension",
		"-loc", "../resources/en-us.wxl",
		"AppFiles.wixobj",
		"installer.wixobj",
		"-o", filepath.Join(msi, balDist + "-" + balPOS + "-" + balPArch + ".msi"),
	)
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func environmentInfo() (cwd, version string, err error) {
	cwd, err = os.Getwd()
	if err != nil {
		return
	}
	version = strings.TrimSpace(balVersion)
	return
}

func installWix(path string) error {
	// Fetch wix binary zip file.
	body, err := httpGet(resourceStore + "windows/wix310-binaries.zip")
	if err != nil {
		return err
	}

	// Unzip to path.
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	if len(zr.File) <= 0 {
		fmt.Println("No zip")
	}
	for _, f := range zr.File {
		name := filepath.FromSlash(f.Name)
		err := os.MkdirAll(filepath.Join(path, filepath.Dir(name)), 0755)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(path, name), b, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func httpGet(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.New(r.Status)
	}
	return body, nil
}

func writeDataFiles(data map[string]string, base string) error {
	for name, body := range data {
		dst := filepath.Join(base, name)
		err := os.MkdirAll(filepath.Dir(dst), 0755)
		if err != nil {
			return err
		}
		b := []byte(body)
		if strings.HasPrefix(body, resourceStore) {
			b, err = httpGet(body)
			if err != nil {
				return err
			}
		}
		// (We really mean 0755 on the next line; some of these files
		// are executable, and there's no harm in making them all so.)
		if err := ioutil.WriteFile(dst, b, 0755); err != nil {
			return err
		}
	}
	return nil
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func runDir(dir, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func wixVersion(v string) string {
	m := versionRe.FindStringSubmatch(v)
	if m == nil {
		return "0.0.0"
	}
	return m[1]
}

func cpDir(dst, src string) error {
	walk := func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, srcPath[len(src):])
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}
		return cp(dstPath, srcPath)
	}
	return filepath.Walk(src, walk)
}

func cp(dst, src string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	fi, err := sf.Stat()
	if err != nil {
		return err
	}
	tmpDst := dst + ".tmp"
	df, err := os.Create(tmpDst)
	if err != nil {
		return err
	}
	defer df.Close()

	if runtime.GOOS != "windows" {
		if err := df.Chmod(fi.Mode()); err != nil {
			return err
		}
	}
	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}
	if err := df.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpDst, dst); err != nil {
		return err
	}
	// Ensure the destination has the same mtime as the source.
	return os.Chtimes(dst, fi.ModTime(), fi.ModTime())
}