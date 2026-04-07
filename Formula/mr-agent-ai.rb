class MrAgentAi < Formula
  desc "Multi-Agent AI Skill Installer — installs agent skills into any project."
  homepage "https://github.com/SaidHernandez/mr-agent-ai"
  license "MIT"
  version "2.1.1"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/SaidHernandez/mr-agent-ai/releases/download/v2.1.1/mr-agent-ai_2.1.1_darwin_arm64.tar.gz"
      sha256 "fbf88d1ce19a7092d1841582dde65a9c1514de4cf22fbd8c2b96c2f468850bc5"
    else
      url "https://github.com/SaidHernandez/mr-agent-ai/releases/download/v2.1.1/mr-agent-ai_2.1.1_darwin_amd64.tar.gz"
      sha256 "ca35ddefa6b6c2919f3c22a71903755f83a70dc6913868ed613893763654ef20"
    end
  end

  def install
    bin.install "mr-agent-ai"
  end

  test do
    system "#{bin}/mr-agent-ai", "--help"
  end
end
