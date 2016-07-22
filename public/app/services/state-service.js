;(function (angular) {
	'use strict';

	angular.module('superside.services')
		.factory('stateService', stateService);

	stateService.$inject = ['$http', '$filter'];

	function stateService($http, $filter) {
		var services = {};
		var deployments = {};
		var events = [];
        var clusters = {};

        var addDeployment = function(deploy) {
			deployments[name] = deployments[name] || {};
            deployments[name][deploy.ID] = deploy;
            updateServiceVersion(deploy);
        };

        var updateServiceVersion = function(deploy) {
            services[deploy.Name] = services[deploy.Name] || {};
            services[deploy.Name][deploy.ClusterName] = services[deploy.Name][deploy.ClusterName] || {};
            services[deploy.Name][deploy.ClusterName].Version = deploy.Version;
            services[deploy.Name][deploy.ClusterName].EndTime = deploy.EndTime;
        };

        var addClusterName = function(event) {

            if (!clusters.hasOwnProperty(event.ClusterName)) {
                clusters[event.ClusterName] = '';
            }

        };

		return {
			events: events,
			services: services,
			deployments: deployments,
            clusters: clusters,

			addDeployment: addDeployment,
            addClusterName: addClusterName,

			run: function() {
				$http({
					method: 'GET',
					url: '/api/state/services',
					dataType: 'json'
				}).then(function(response) {

					_.each(response.data, function(event) {
                        var filteredEvent = $filter('uiEvent')(event);
						events.push(filteredEvent);
                        addClusterName(filteredEvent);
					});

					$http({
						method: 'GET',
						url: '/api/state/deployments',
						dataType: 'json'
					}).then(function(response) {
						_.each(response.data, function(service) {
							_.each(service, function(deploy) {
                                addDeployment(deploy);
								events.push($filter('uiEvent')(deploy));
                                addClusterName(deploy);
							});
						});

					}, function (error) {
						console.log('ERROR: ' + error);
					});
				}, function (error) {
					console.log('ERROR: ' + error);
				});
			}
		}
	}

})(angular);
