class MrAgentAi < Formula
  desc "Multi-Agent AI Skill Installer — installs agent skills into any project."
  homepage "https://github.com/SaidHernandez/mr-agent-ai"
  license "MIT"
  version "3.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/SaidHernandez/mr-agent-ai/releases/download/v3.0.0/mr-agent-ai_3.0.0_darwin_arm64.tar.gz"
      sha256 ""
    else
      url "https://github.com/SaidHernandez/mr-agent-ai/releases/download/v3.0.0/mr-agent-ai_3.0.0_darwin_amd64.tar.gz"
      sha256 ""
    end
  end

  def install
    bin.install "mr-agent-ai"
  end

  test do
    system "#{bin}/mr-agent-ai", "--help"
  end
end
