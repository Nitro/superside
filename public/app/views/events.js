;(function (angular) {
    'use strict';

    angular.module('superside.views')
        .controller('eventsController', eventsController);

    eventsController.$inject = ['$filter', 'stateService', 'websocketService'];

    function eventsController($filter, stateService, websocketService) {

        var self = this;
        self.events = [];

        stateService.resolveInitialServices()
            .then(function() {
                var initialEvents = stateService.getCurrentServices();

                _.forEach(initialEvents, function(event) {
                    self.events.push($filter('uiEvent')(event));
                });

            });

        stateService.resolveInitialDeployments()
            .then(function() {
                var initialDeployments = stateService.getCurrentDeployments();

                _.forEach(initialDeployments, function(services) {

                    _.forEach(services, function(service) {
                        self.events.push($filter('uiEvent')(service));
                    });

                });

            });

    }

})(angular);
