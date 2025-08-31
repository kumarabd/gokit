#!/usr/bin/env ruby

require 'net/http'
require 'json'

# Configuration
VERSION = ARGV[0] || "0.2.4"
REPO = "kumarabd/gokit"
FORMULA_PATH = "Formula/gokit.rb"

# GitHub API endpoint for releases
GITHUB_API = "https://api.github.com/repos/#{REPO}/releases/tags/v#{VERSION}"

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

# Main execution
puts "Updating Homebrew formula for GoKit v#{VERSION}"

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
  puts "\n✅ Successfully updated Homebrew formula!"
  puts "Next steps:"
  puts "1. Commit and push the updated formula to the homebrew branch"
  puts "2. Users can now install with: brew tap kumarabd/gokit && brew install gokit"
else
  puts "\n❌ Failed to calculate all SHA256 hashes"
  exit 1
end
