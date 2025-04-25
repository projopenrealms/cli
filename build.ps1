$version = "v1"
$name = "customrealms-cli"
$main = "./cmd/"

$targets = @(
  @{ OS = "windows"; ARCH = "amd64" },
  @{ OS = "windows"; ARCH = "386" },
  @{ OS = "windows"; ARCH = "arm64" },
  @{ OS = "linux"; ARCH = "amd64" },
  @{ OS = "linux"; ARCH = "386" },
  @{ OS = "linux"; ARCH = "arm64" },
  @{ OS = "darwin"; ARCH = "amd64" },
  @{ OS = "darwin"; ARCH = "arm64" }
)

foreach ($target in $targets) {
  $os = $target.OS
  $arch = $target.ARCH
  $folder = "dist/${name}_${os}_${arch}_${version}"
  $binary = if ($os -eq "windows") { "crx.exe" } else { "crx" }

  New-Item -ItemType Directory -Force -Path $folder | Out-Null
  Write-Host "→ Building $os/$arch → $folder\$binary"

  $env:GOOS = $os
  $env:GOARCH = $arch
  go build -o "$folder\$binary" $main
}