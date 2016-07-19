'use strict';

// Declare app level module which depends on views, and components
angular.module('superside', [
  'ngRoute',
  'superside.services'
]).
config(['$locationProvider', '$routeProvider', function($locationProvider, $routeProvider) {
    $locationProvider.hashPrefix('!');

    $routeProvider.when('/', {
        templateUrl: 'views/websockets-demo.html',
        controller: 'servicesCtrl'
    });

    $routeProvider.otherwise({redirectTo: '/'});
}]);
