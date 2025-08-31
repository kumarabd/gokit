#!/usr/bin/env ruby

require 'net/http'
require 'json'
require 'fileutils'

# Configuration
VERSION = ARGV[0] || "0.2.4"
MAIN_REPO = "kumarabd/gokit"
HOMEBREW_REPO = "kumarabd/homebrew-gokit"
HOMEBREW_DIR = "../homebrew-gokit"
FORMULA_PATH = "#{HOMEBREW_DIR}/Formula/gokit.rb"

# GitHub API endpoint for releases
GITHUB_API = "https://api.github.com/repos/#{MAIN_REPO}/releases/tags/v#{VERSION}"

def get_release_assets
  uri = URI(GITHUB_API)
  response = Net::HTTP.get_response(uri)
  
  if response.code != "200"
    puts "Error: Could not fetch release information for v#{VERSION}"
    puts "Response: #{response.code} - #{response.body}"
    exit 1
  end
  
  data = JSON.parse(response.body)
  data['assets']
end

def download_and_calculate_sha256(url, filename)
  puts "Downloading #{filename}..."
  
  uri = URI(url)
  response = Net::HTTP.get_response(uri)
  
  if response.code != "200"
    puts "Error: Could not download #{filename}"
    return nil
  end
  
  # Calculate SHA256
  require 'digest'
  sha256 = Digest::SHA256.hexdigest(response.body)
  
  puts "SHA256 for #{filename}: #{sha256}"
  sha256
end

def update_formula(version, hashes)
  formula_content = File.read(FORMULA_PATH)
  
  # Update version
  formula_content.gsub!(/version "[\d.]+"/, "version \"#{version}\"")
  
  # Update SHA256 hashes
  hashes.each do |platform, hash|
    placeholder = "PLACEHOLDER_SHA256_#{platform.upcase}"
    formula_content.gsub!(placeholder, hash)
  end
  
  File.write(FORMULA_PATH, formula_content)
  puts "Updated #{FORMULA_PATH} with version #{version}"
end

def commit_and_push_changes(version)
  Dir.chdir(HOMEBREW_DIR) do
    # Check if there are changes
    status = `git status --porcelain`
    if status.empty?
      puts "No changes to commit"
      return
    end
    
    # Commit changes
    system("git add Formula/gokit.rb")
    system("git commit -m \"Update GoKit to v#{version}\"")
    
    # Push changes
    if system("git push origin main")
      puts "✅ Successfully pushed changes to homebrew-gokit repository"
    else
      puts "❌ Failed to push changes"
      exit 1
    end
  end
end

# Main execution
puts "Updating Homebrew tap for GoKit v#{VERSION}"

# Check if homebrew directory exists
unless Dir.exist?(HOMEBREW_DIR)
  puts "Error: Homebrew directory not found at #{HOMEBREW_DIR}"
  puts "Please clone the homebrew-gokit repository first:"
  puts "git clone https://github.com/#{HOMEBREW_REPO}.git #{HOMEBREW_DIR}"
  exit 1
end

assets = get_release_assets()
hashes = {}

assets.each do |asset|
  name = asset['name']
  url = asset['browser_download_url']
  
  case name
  when /darwin-arm64/
    hashes['ARM64'] = download_and_calculate_sha256(url, name)
  when /darwin-amd64/
    hashes['AMD64'] = download_and_calculate_sha256(url, name)
  when /linux-amd64/
    hashes['LINUX'] = download_and_calculate_sha256(url, name)
  end
end

if hashes.values.all?
  update_formula(VERSION, hashes)
  commit_and_push_changes(VERSION)
  puts "\n✅ Successfully updated Homebrew tap!"
  puts "Users can now install with: brew tap kumarabd/gokit && brew install gokit"
else
  puts "\n❌ Failed to calculate all SHA256 hashes"
  exit 1
end
