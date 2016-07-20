;(function (angular) {
    'use strict';

    angular.module('superside.views')
        .controller('eventsController', eventsController);

    eventsController.$inject = ['$filter', 'stateService', 'websocketService'];

    function eventsController($filter, stateService, websocketService) {

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
			onError: function(event) {
				console.log(event)
			},
		})

    }

})(angular);
