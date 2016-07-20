;(function (angular) {
    'use strict';

    angular.module('superside.views')
        .controller('eventsController', eventsController);

    eventsController.$inject = ['$scope', 'stateService', 'websocketService'];

    function eventsController($scope, stateService, websocketService) {

		var self = this;
		self.events = stateService.events;

        /*stateService.resolveInitialServices()
            .then(function() {
                var initialEvents = stateService.getCurrentServices();

                _.forEach(initialEvents, function(event) {
                    stateService.events.push($filter('uiEvent')(event));
                });

            });

        stateService.resolveInitialDeployments()
            .then(function() {
                var initialDeployments = stateService.getCurrentDeployments();

                _.forEach(initialDeployments, function(services) {

                    _.forEach(services, function(service) {
                        stateService.events.push($filter('uiEvent')(service));
                    });

                });

            });
			*/

		websocketService.setupSocket({
			onMessage: function() {
				$scope.$apply(); // WTF Angular... thanks for not updating the binding for us
			},
			onError: function(event) {
				console.log(event)
			},
		})

    }

})(angular);
