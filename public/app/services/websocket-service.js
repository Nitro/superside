;(function (angular) {
    'use strict';

    angular.module('superside.services')
        .factory('websocketService', websocketService);

    websocketService.$inject = ['stateService', '$filter'];

    function websocketService(stateService, $filter) {

        var socket = null;

        return {

            setupSocket: function(options) {
                // Create a new WebSocket.
                var wsUrl = 'ws://' + window.location.hostname + ':7779/listen';
                socket = new WebSocket(wsUrl);

                // Handle any errors that occur.
                socket.onerror = function(event) {
                    console.log('WebSocket Error: ' + event);

                    if (typeof options.onError === 'function') {
                        options.onError(event);
                    }
                };

                // Show a connected message when the WebSocket is opened.
                socket.onopen = function(event) {

                    if (typeof options.onOpen === 'function') {
                        options.onOpen(event);
                    }

                };

                // Handle messages sent by the server.
                socket.onmessage = function(event) {
                    var message = event.data;
                    var evt = angular.fromJson(message);

                    var filteredEvent = $filter('uiEvent')(evt.Data);
                    stateService.events.push(filteredEvent);
                    stateService.addClusterName(filteredEvent);

                    if(evt.Type == 'Deployment') {
                        stateService.addDeployment(evt.Data)
                    }

                    if (typeof options.onMessage === 'function') {
                        options.onMessage(message);
                    }
                };


                // Show a disconnected message when the WebSocket is closed.
                socket.onclose = function(event) {
                    if (typeof options.onClose === 'function') {
                        options.onClose(event);
                    }

                };
            },

            sendMessage: function(message) {
                socket.send(message);
            },

            closeSocket: function() {
                socket.close();
            }

        }

    }

})(angular);
