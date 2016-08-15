#!/usr/bin/env ruby

require 'json'
require 'excon'
require 'trollop'
require 'time'

opts = Trollop::options do
  opt :url,     'The URL to post updates to', default: 'http://localhost:7779/api/update'
  opt :fixture, 'The fixture to post', default: '../../fixtures/deployment-sample.json'
  opt :sleep,   'The number or decimal seconds to sleep between events', default: 0.3
  opt :service, 'The name of the service to deploy', type: String
  opt :cluster, 'The cluster to deploy the service into', type: String
  opt :svc_version, 'The version to deploy', type: String
end

connection = Excon.new(opts[:url])

fixture = File.expand_path(opts[:fixture], __FILE__)

# If we specified a cluster name, use that, otherwise defaults
clusternames = if opts[:cluster]
  [opts[:cluster]]
else
  %w{ nitro-dev nitro-prod }
end

# If we specified a service name, use that, otherwise defaults
servicenames = if opts[:service]
  [opts[:service]]
else
  %w{ awesome-svc good-svc }
end

data = JSON.parse(File.read(fixture))

puts opts[:svc_version]
unless opts[:svc_version].nil? || opts[:svc_version].empty?
  data = data.select { |e| e['Event']['Service']['Image'] =~ /#{opts[:svc_version]}/ }
end

clusternames.shuffle.each do |clustername|
  data.each do |event|
    # We want prod to end up 1 version behind dev
    next if clustername == 'nitro-prod' && event['Event']['Service']['Image'] =~ /0.3/

    sleep opts[:sleep]

    event['Event']['Service']['Hostname'] = %w{ host1 host2 host3 host4 }.shuffle.first
    event['Event']['Service']['Name'] = servicenames.shuffle.first

    time = if opts[:sleep] > 0
      Time.now.utc
    else
      Time.parse(event['Event']['Time'])
    end

    json_data = JSON.pretty_generate(
      'State' => {
        'ClusterName' => clustername
      },
      'ChangeEvent' => {
        'Service' => event['Event']['Service'],
        'PreviousStatus' => event['Event']['PreviousStatus'],
        'Time' => time.xmlschema
      }
    )

    connection.post(
      :body => json_data,
      :headers => { 'Content-type' => 'application/json' },
      :persistent => true
    )
  end
end
