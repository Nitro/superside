#!/usr/bin/env ruby

require 'json'
require 'excon'
require 'trollop'
require 'time'

opts = Trollop::options do
  opt :url,     'The URL to post updates to', default: 'http://localhost:7778/api/update'
  opt :fixture, 'The fixture to post', default: '../../fixtures/deployment-sample.json'
  opt :sleep,   'The number or decimal seconds to sleep between events', default: 0.3
end

connection = Excon.new(opts[:url])

fixture = File.expand_path(opts[:fixture], __FILE__)
data = JSON.parse(File.read(fixture))

%w{ nitro-dev nitro-prod }.shuffle.each do |clustername|
    data.each do |event|
      sleep opts[:sleep]
      event['Event']['Service']['Hostname'] = %w{ host1 host2 host3 host4 }.shuffle.first
      event['Event']['Service']['Name'] = %w{ awesome-svc good-svc }.shuffle.first

      time = if opts[:sleep] > 0
        Time.now.utc
      else
        event['Event']['Time']
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
