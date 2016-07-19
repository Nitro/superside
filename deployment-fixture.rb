#!/usr/bin/env ruby

require 'json'
require 'excon'

connection = Excon.new('http://localhost:7778/update')

state_blob = {
  'State' => {
    'ClusterName' => 'nitro-dev',
    'Hostname' => 'awesome-host'
  }
}

data = JSON.parse(File.read('deployment-sample.json'))

data.each do |event|
  connection.post(
    :body => state_blob.merge({
        'ChangeEvent' => {
          'Service' => event['Event']['Service']
        }
      }).to_json,
    :headers => { 'Content-type' => 'application/json' },
    :persistent => true
  )
end
