;(function (angular) {
    'use strict';

    // Declare app level module which depends on views, and components
    angular.module('superside', [
        'ngRoute',
		'nvd3',
        'superside.services',
        'superside.views',
        'superside.filters'
    ]);

    angular.module('superside.services', []);
    angular.module('superside.views', []);
    angular.module('superside.filters', []);

    angular.module('superside')
        .config(['$locationProvider', '$routeProvider', function ($locationProvider, $routeProvider) {
            $locationProvider.hashPrefix('!');

            $routeProvider.when('/', {
                templateUrl: 'views/events.html',
                controller: 'eventsController',
                controllerAs: 'eventsCtrl'
            });

            $routeProvider.when('/dashboard', {
                templateUrl: 'views/dashboard.html',
                controller: 'dashboardController',
                controllerAs: 'dashboardCtrl'
            });

            $routeProvider.when('/graph', {
                templateUrl: 'views/graph.html',
                controller: 'graphController',
                controllerAs: 'graphCtrl'
            });

            $routeProvider.when('/websockets-demo', {
                templateUrl: 'views/websockets-demo.html',
                controller: 'websocketCtrl'
            });

            $routeProvider.otherwise({redirectTo: '/'});
        }]);

})(angular);
