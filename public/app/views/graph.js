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
			// Count of total # of events
			var grouped = _.groupBy(self.events, 'Name');
			self.eventCounts = _.map(grouped, function(evts) { 
				return { key: evts[0].Name.slice(0, 20), y: evts.length }
			}).slice(0, 5);

			// Events where the new status was Unhealthy
			self.unhealthyEvents = [ 
				{
					key: 'Unhealthy Event Count',
					values: filterByStatus(grouped, 'Unhealthy').slice(0, 5)
				}
			];

			// Events where the new status was Tombstone
			self.tombstoneEvents = [
				{
					key: 'Tombstone Event Count',
					values: filterByStatus(grouped, 'Tombstone').slice(0, 5)
				}
			];

			// Line chart events showing status changes
			self.statusChangeEvents = aggregateStatusEvents(self.events, 'Status Change Events Over Time');
		});

        self.pieOptions = {
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

		self.barOptions = {
			chart: {
                x: function(d) { return d.label },
                y: function(d) { return d.value },
                type: 'discreteBarChart',
                height: height,
				width: width,
				margin: {
					top: 20,
					right: 20,
					bottom: 50,
					left: 55
				},
                duration: 500,
				xAxis: {
					rotateLabels: -90
				}
            }
		};

		self.lineOptions = {
			chart: {
                x: function(d) { return d[0] },
                y: function(d) { return d[1] },
				type: "lineChart",
				height: 450,
				margin: {
					top: 20,
					right: 20,
					bottom: 50,
					left: 65
				},
				color: [
					"#1f77b4",
					"#ff7f0e",
					"#2ca02c",
					"#d62728",
					"#9467bd",
					"#8c564b",
					"#e377c2",
					"#7f7f7f",
					"#bcbd22",
					"#17becf"
				],
				useInteractiveGuideline: true,
				clipVoronoi: false,
				xAxis: {
					axisLabel: "Time",
					showMaxMin: true,
					staggerLabels: true,
					tickFormat: function(d) {
						// How TF does this happen? Anyway, work around it
						if (_.isString(d)) {
							d = parseInt(d);
						}
            			return d3.time.format('%Y-%m-%d %H:%M:%S')(new Date(d));
        			},
					rotateLabels: -90
				},
				yAxis: {
					axisLabel: "Count",
					axisLabelDistance: 0
				}
			}
		};
	};

	// Expectes events grouped by Name and returns data formatted
	// for bar charts.
	function filterByStatus(events, status) {
		var filtered = _.map(events, function(evts) {
			return { 
					label: evts[0].Name,
					status: _.groupBy(evts, function(evt) { return evt.Status })
			};
		});

		filtered = _.map(filtered, function(evt) {
			var length = 0;
			if(evt.status[status]) {
				length = evt.status[status].length;
			}

			return { label: evt.label.slice(0, 20), value: length }
		});

		filtered = _.filter(filtered, function(evt) { return evt.value > 0; } );

 		return _.sortBy(filtered, function(evt) { return evt.value }).reverse();
	};

	// Aggregate status change events by time bucket
	function aggregateStatusEvents(events) {
		var filteredByEnv = _.groupBy(events, function(evt) { return evt.ClusterName });
		var result = [];

		for(var env in filteredByEnv) {
			var bucketed = _.reduce(filteredByEnv[env], function(memo, evt) {
				// Bucket by 5 minute intervals
				var bucket = Math.floor((Date.parse(evt.Time) / 300000)) * 300000
				if(memo[bucket] == null) {
					memo[bucket] = 0;
				}
				memo[bucket]++;

				return memo;
			}, {});

			var roughValues = _.map(bucketed, function(k, v) {
				return [ v, k ];
			});

			result.push({
				key: env,
				// Make sure the values are ordered or the chart flips out
				values: _.sortBy(roughValues, function(item) { item[0] })
			});

		}

		return result;
	};
})(angular);
