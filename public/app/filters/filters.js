;(function (angular) {
    'use strict';

    angular.module('superside.services')

    .filter('portsStr', function() {
        return function(svcPorts) {
            var ports = [];
            for (var i in svcPorts) {
                var port = svcPorts[i];

                if (port.Port == null) {
                    continue;
                }

                if (port.ServicePort != null && port.ServicePort != 0) {
                    ports.push(port.ServicePort.toString() + "->" + port.Port.toString());
                } else {
                    ports.push(port.Port.toString());
                }
            }

            return ports.join(", ");
        };
    })

    .filter('statusStr', function() {
        return function(status) {

            switch (status) {
                case 0: return "Alive";
                case 2: return "Unhealthy";
                case 3: return "Unknown";
                default: return "Tombstone";
            }

            return statusStr
        }
    })

    .filter('timeAgo', function() {
        return function(textDate) {
            if (textDate == null || textDate == "" || textDate == "1970-01-01T01:00:00+01:00") {
                return "never";
            }

            var date = Date.parse(textDate);
            var seconds = Math.floor((new Date() - date) / 1000);

            var interval = Math.floor(seconds / 31536000);

            if (interval > 1) {
                return interval + " years ago";
            }
            interval = Math.floor(seconds / 2592000);
            if (interval > 1) {
                return interval + " months ago";
            }
            interval = Math.floor(seconds / 86400);
            if (interval > 1) {
                return interval + " days ago";
            }
            interval = Math.floor(seconds / 3600);
            if (interval > 1) {
                return interval + " hours ago";
            }
            interval = Math.floor(seconds / 60);
            if (interval > 1) {
                return interval + " mins ago";
            }
            return Math.floor(seconds) + " secs ago";
        }
    })

    .filter('imageStr', function() {
        return function(imageName) {
            if (imageName.length < 25) {
                return imageName;
            }

            return imageName.substr(imageName.length - 25, imageName.length);
        }
    })

    .filter('extractTag', function() {
        return function(imageName) {
            return imageName.split(':')[1];
        }
    })

    .filter('uiEvent', function($filter) {

        return function(incident) {

            var cleanServiceEvent = {
                Name: null,
                Version: null,
                ClusterName: null,
				Hostnames: null,
                Type: null,
                Time: null,
                StartTime: null,
                EndTime: null,
                DeploymentId: null,
				StatusCode: null
            };

            if (incident.hasOwnProperty('Event')) {
                var service = incident.Event.Service;

                cleanServiceEvent.Type = 'ServiceEvent';
                cleanServiceEvent.ClusterName = incident.ClusterName;
                cleanServiceEvent.Name = service.Name;
                cleanServiceEvent.Version = $filter('extractTag')(service.Image);
                cleanServiceEvent.PreviousStatus = $filter('statusStr')(incident.Event.PreviousStatus);
				cleanServiceEvent.Status = $filter('statusStr')(service.Status);
				cleanServiceEvent.StatusCode = service.Status;
                cleanServiceEvent.Time = incident.Event.Time;
                cleanServiceEvent.Hostnames = [service.Hostname];
            } else {
                cleanServiceEvent.Type = 'Deployment';
                cleanServiceEvent.ClusterName = incident.ClusterName;
                cleanServiceEvent.Name = incident.Name;
                cleanServiceEvent.Version = incident.Version;
                cleanServiceEvent.Time = incident.StartTime;
                cleanServiceEvent.EndTime = incident.EndTime;
                cleanServiceEvent.StartTime = incident.StartTime;
                cleanServiceEvent.DeploymentId = incident.ID;
                cleanServiceEvent.Hostnames = incident.Hostnames;

            }

            return cleanServiceEvent;

        }

    })

    .filter('eventsFilter', function (){

        return function(events, filters) {
            var filteredEvents = [];

            if (filters.cluster === filters.service && filters.cluster === '') {
                return events;
            }

            _.forEach(events, function (event) {

                if ((filters.cluster === '' && filters.service === event.Name)
                    || (filters.service === '' && filters.cluster === event.ClusterName)
                    || (filters.service === event.Name && filters.cluster === event.ClusterName)) {

                    filteredEvents.push(event);
                }


            });

            return filteredEvents;
        }

    });

    if ( ! Array.prototype.groupBy) {
        Array.prototype.groupBy = function (f)
        {
            var groups = {};
            this.forEach(function(o) {
                var group = JSON.stringify(f(o));
                groups[group] = groups[group] || [];
                groups[group].push(o);
            });

            return Object.keys(groups).map(function (group) {
                return groups[group];
            });
        };
    }

})(angular);
