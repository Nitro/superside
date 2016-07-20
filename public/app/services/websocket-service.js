;(function (angular) {
    'use strict';

    angular.module('superside.services')
        .factory('websocketService', websocketService);

    websocketService.$inject = [];

    function websocketService() {

        var socket = null;

        return {

            setupSocket: function(options) {
                // Create a new WebSocket.
                socket = new WebSocket('ws://localhost:7778/listen');

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
