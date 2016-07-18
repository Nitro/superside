#!/usr/bin/env ruby

require 'excon'
require 'json'

raw = Excon.get('http://localhost:7778/state').body
data = JSON.parse(raw)

def status_str_for(status)
  case status
  when 0 then 'Alive'
  when 1 then 'Tombstone'
  when 2 then 'Unhealthy'
  when 3 then 'Unknown'
  end
end

puts 'Events'
puts '-' * 80
data.each do |notification|
  svc = notification['Event']['Service']
  puts "#{'%-30s' % svc['Updated']} #{'%15s' % notification['ClusterName']} #{'%20s' % svc['Hostname']} " +
    "#{'%25s' % svc['Name']}  " +
    "#{status_str_for(notification['Event']['PreviousStatus'])} --> " +
    "#{status_str_for(svc['Status'])}"
end
