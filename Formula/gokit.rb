class Gokit < Formula
  desc "GoKit - Microservice development toolkit"
  homepage "https://github.com/kumarabd/gokit"
  version "0.2.4"
  head "https://github.com/kumarabd/gokit.git", branch: "homebrew"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/kumarabd/gokit/releases/download/v#{version}/gokit-darwin-arm64"
      sha256 "PLACEHOLDER_SHA256_ARM64"
    else
      url "https://github.com/kumarabd/gokit/releases/download/v#{version}/gokit-darwin-amd64"
      sha256 "PLACEHOLDER_SHA256_AMD64"
    end
  end

  on_linux do
    url "https://github.com/kumarabd/gokit/releases/download/v#{version}/gokit-linux-amd64"
    sha256 "PLACEHOLDER_SHA256_LINUX"
  end

  def install
    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "gokit-darwin-arm64" => "gokit"
      else
        bin.install "gokit-darwin-amd64" => "gokit"
      end
    else
      bin.install "gokit-linux-amd64" => "gokit"
    end
  end

  test do
    system "#{bin}/gokit", "version"
  end
end
