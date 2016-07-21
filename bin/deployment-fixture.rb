#!/usr/bin/env ruby

require 'json'
require 'excon'
require 'trollop'

opts = Trollop::options do
  opt :url,     'The URL to post updates to', default: 'http://localhost:7778/api/update'
  opt :fixture, 'The fixture to post', default: '../../fixtures/deployment-sample.json'
  opt :sleep,   'The number or decimal seconds to sleep between events', default: 0.3
end

connection = Excon.new(opts[:url])

state_blob = {
  'State' => {
    'ClusterName' => 'nitro-dev',
    'Hostname' => 'awesome-host'
  }
}

fixture = File.expand_path(opts[:fixture], __FILE__)
data = JSON.parse(File.read(fixture))

data.each do |event|
  sleep opts[:sleep]
  event['Event']['Service']['Hostname'] = %w{ host1 host2 host3 host4 }.shuffle.first

  data =  JSON.pretty_generate(state_blob.merge({
    'ChangeEvent' => {
      'Service' => event['Event']['Service'],
      'PreviousStatus' => event['Event']['PreviousStatus'],
      'Time' => event['Event']['Time']
    }
  }))

  connection.post(
    :body => data,
    :headers => { 'Content-type' => 'application/json' },
    :persistent => true
  )
end
