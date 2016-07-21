;(function (angular) {
	'use strict';

	angular.module('superside.services')
		.factory('stateService', stateService);

	stateService.$inject = ['$http', '$q', '$filter'];

	function stateService($http, $q, $filter) {
		var services = [];
		var deployments = {};
		var events = [];

		var self = this

		return {
			events: events,
			services: services,
			deployments: deployments,

			addDeployment: function(deploy) {
				if (deployments[name] == null) {
					deployments[name] = {};
				}
				deployments[name][deploy.ID] = deploy;
			},

			run: function() {
				$http({
					method: 'GET',
					url: '/api/state/services',
					dataType: 'json'
				}).then(function(response) {
					services = response.data;
					_.each(response.data, function(event) {
						events.push($filter('uiEvent')(event));
					});

					$http({
						method: 'GET',
						url: '/api/state/deployments',
						dataType: 'json'
					}).then(function(response) {
						_.each(response.data, function(service, name) {
							_.each(service, function(deploy) {
								// TODO this should call addDeployment but scope...
								if (deployments[name] == null) {
									deployments[name] = {};
								}
								deployments[name][deploy.ID] = deploy;
								events.push($filter('uiEvent')(deploy));
							});
						});

						events.sort(function(evt1, evt2) {
							if (evt1.sortCode() < evt2.sortCode()) return -1;
							return 1;
						});

					}, function (error) {
						console.log('ERROR: ' + error);
					});
				}, function (error) {
					console.log('ERROR: ' + error);
				});
			},
		}
	}

})(angular);
