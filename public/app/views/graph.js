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

		self.testdata2 = [
			{ key: 'One', y: 1 },
			{ key: 'Two', y: 1 },
			{ key: 'Three', y: 1 },
			{ key: 'Four', y: 1 },
			{ key: 'Five', y: 1 },
			{ key: 'Six', y: 1 },
			{ key: 'Seven', y: 1 }
		];

        var arcRadius2 = [
			{ inner: 0.9, outer: 1 },
			{ inner: 0.8, outer: 1 },
			{ inner: 0.7, outer: 1 },
			{ inner: 0.6, outer: 1 },
			{ inner: 0.5, outer: 1 },
			{ inner: 0.4, outer: 1 },
			{ inner: 0.3, outer: 1 }
		];

        var height = 350;
        var width = 350;

        var colors = ['green', 'gray'];

		// We have to wait on the data to come back from the stateService
		stateService.onSuccess.push(function() {
			var grouped = _.groupBy(self.events, 'Name');
			self.percentages = _.map(grouped, function(evts) { return { key: evts[0].Name, y: evts.length } });

			// Scale the bars by percentage as well
			self.arcRadii = _.map(grouped, function(evts) { 
				return { 
					inner: evts.length / self.events.length,
					outer: 1
				};
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
				arcsRadius: self.arcRadii,
                labelThreshold: 0.01,
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
