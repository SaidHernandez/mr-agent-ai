class MrAgentAi < Formula
  desc "Multi-Agent AI Skill Installer — installs agent skills into any project."
  homepage "https://github.com/SaidHernandez/mr-agent-ai"
  license "MIT"
  version "2.0.1"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/SaidHernandez/mr-agent-ai/releases/download/v2.0.1/mr-agent-ai_2.0.1_darwin_arm64.tar.gz"
      sha256 "4fb5ef379014bfd88e5b1fdc6b98a5eb3a7cb7ac20cf2338a3c98b173d2e29e2"
    else
      url "https://github.com/SaidHernandez/mr-agent-ai/releases/download/v2.0.1/mr-agent-ai_2.0.1_darwin_amd64.tar.gz"
      sha256 "e543999d4e3041e43a6280964e8e8b32d44863cf525c539d959c19d71f1a0868"
    end
  end

  def install
    bin.install "mr-agent-ai"
  end

  test do
    system "#{bin}/mr-agent-ai", "--help"
  end
end
