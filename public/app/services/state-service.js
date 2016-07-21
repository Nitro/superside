;(function (angular) {
	'use strict';

	angular.module('superside.services')
		.factory('stateService', stateService);

	stateService.$inject = ['$http', '$q', '$filter'];

	function stateService($http, $q, $filter) {
		var services = {};
		var deployments = {};
		var events = [];

        var addDeployment = function(deploy) {
            if (deployments[name] == null) {
                deployments[name] = {};
            }
            deployments[name][deploy.ID] = deploy;
            updateServiceVersion(deploy);
        };

        var updateServiceVersion = function(event) {
            services[event.Name] = services[event.Name] || {};
            services[event.Name][event.ClusterName] = services[event.Name][event.ClusterName] || {};

            services[event.Name][event.ClusterName].Version = event.Version;
            services[event.Name][event.ClusterName].Time = event.Time;
        };

		return {
			events: events,
			services: services,
			deployments: deployments,

			addDeployment: addDeployment,

			run: function() {
				$http({
					method: 'GET',
					url: '/api/state/services',
					dataType: 'json'
				}).then(function(response) {

					_.each(response.data, function(event) {

                        var filteredEvent = $filter('uiEvent')(event);

						events.push(filteredEvent);
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
