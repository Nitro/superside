;(function (angular) {
    'use strict';

    angular.module('superside.views')
        .controller('eventsController', eventsController);

    eventsController.$inject = ['$scope', 'stateService', 'websocketService'];

    function eventsController($scope, stateService, websocketService) {

		var self = this;
		self.events = stateService.events;
		self.deployments = stateService.deployments;
        self.services = stateService.services;
        self.clusters = stateService.clusters;
        self.filters = {
            cluster: '',
            service: ''
        };

		stateService.run();

		websocketService.setupSocket({
			onMessage: function() {
 				// WTF Angular... thanks for not updating the binding for us
				$scope.$apply();
			},
			onError: function(event) {
				console.log(event)
			}
		});

		// Compare versions. If they don't match, warning. If they don't match
		// and also are more than a day apart, then danger/red.
		self.getVersionMatchStatus = function(svcName) {
			var clusters = stateService.services[svcName];
			var groups = _.groupBy(clusters, function(cluster) {
				return cluster.Version
			});

			// All versions matched, good to go
			if (Object.keys(groups).length == 1) {
				return '';
			}

			var sorted = _.sortBy(clusters, 'EndTime')

			var oneDay = 3600*24*1000;
			var firstDate = new Date(sorted[0].EndTime).getTime();
			var lastDate = new Date(sorted[sorted.length-1].EndTime).getTime();

			if (lastDate > firstDate + oneDay) {
				return 'bg-danger';
			}

			return 'bg-warning';
		};

    }
})(angular);
