;(function (angular) {
    'use strict';

    angular.module('superside.services')
        .factory('stateService', stateService);

    stateService.$inject = ['$http', '$q'];

    function stateService($http, $q) {

        var services = [];
        var deployments = [];
        var events = [];

        return {
			events:events,

            resolveInitialServices: function() {

                var servicesResolve = $q.defer();

                $http({
                    method: 'GET',
                    url: '/api/state/services',
                    dataType: 'json'
                }).then(function(response) {
                    services = response.data;
                    servicesResolve.resolve();
                }, function (error) {
                    console.log('ERROR: ' + error);
                    servicesResolve.reject();
                });

                return servicesResolve.promise;

            },

            resolveInitialDeployments: function() {

                var deploymentsResolve = $q.defer();

                $http({
                    method: 'GET',
                    url: '/api/state/deployments',
                    dataType: 'json'
                }).then(function(response) {
                    deployments = response.data;
                    deploymentsResolve.resolve();
                }, function (error) {
                    console.log('ERROR: ' + error);
                    deploymentsResolve.reject();
                });

                return deploymentsResolve.promise;

            },

            getCurrentServices: function() {
                return services;
            },

            getCurrentDeployments: function() {
                return deployments;
            }

        }

    }

})(angular);
