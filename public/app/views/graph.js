;(function (angular) {
    'use strict';

    angular.module('superside.views')
        .controller('graphController', graphController);

    graphController.$inject = ['$scope', 'stateService', 'websocketService'];

    function graphController($scope, stateService, websocketService) {

		var self = this;

		self.events = stateService.events;
		self.deployments = stateService.deployments;
        self.clusters = stateService.clusters;

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

        var height = 350;
        var width = 350;

		// We have to wait on the data to come back from the stateService
		stateService.onSuccess.push(function() {
			// Percentage of total events
			var grouped = _.groupBy(self.events, 'Name');
			self.eventPercentages = _.map(grouped, function(evts) { 
				return { key: evts[0].Name.slice(0, 20), y: evts.length }
			}).slice(0, 5);

			// Events where the new status was Unhealthy
			var flaps = _.map(grouped, function(evts) {
				return { 
						key: evts[0].Name,
						flaps: _.groupBy(evts, function(evt) { return evt.Status })
				};
			}).slice(0, 5);

			self.eventFlaps = _.map(flaps, function(flap) {
				return { key: flap.key, y: flap.flaps.Unhealthy.length }
			});

		});

        self.options = {
            chart: {
                x: function(d) { return d.key },
                y: function(d) { return d.y },
                type: 'pieChart',
                height: height,
				width: width,
				donut: true,
                showLabels: true,
                duration: 500,
                labelThreshold: 0.1,
				labelsOutside: true,
                labelSunbeamLayout: true,
				callback: function() {
					d3.selectAll('.nv-pieLabels text').style('fill', "white");
				},
                legend: {
                    margin: {
                        top: 5,
                        right: 35,
                        bottom: 5,
                        left: 0
                    }
                }
            }
		};
    }
})(angular);
