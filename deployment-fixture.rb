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
