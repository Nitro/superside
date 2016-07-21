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

		stateService.run();

		websocketService.setupSocket({
			onMessage: function() {
 				// WTF Angular... thanks for not updating the binding for us
				$scope.$apply();
			},
			onError: function(event) {
				console.log(event)
			}
		})

    }
})(angular);
