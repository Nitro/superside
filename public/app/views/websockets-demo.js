;(function (angular) {
    'use strict';

    angular.module('superside.views')
        .controller('websocketCtrl', websocketCtrl);

    websocketCtrl.$inject = ['websocketService'];

    function websocketCtrl(websocketService) {

        // Get references to elements on the page.
        var form = document.getElementById('message-form');
        var messageField = document.getElementById('message');
        var messagesList = document.getElementById('messages');
        var socketStatus = document.getElementById('status');
        var closeBtn = document.getElementById('close');

        var callbacks = {
            onOpen: function(event) {
                socketStatus.innerHTML = 'Connected to: ' + event.currentTarget.URL;
                socketStatus.className = 'open';
            },

            onClose: function() {
                socketStatus.innerHTML = 'Disconnected from WebSocket.';
                socketStatus.className = 'closed';
            },

            onMessage: function(message) {

                messagesList.innerHTML += '<li class="received"><span>Received:</span>' + message + '</li>';

            }

        };

        websocketService.setupSocket(callbacks);

        // Send a message when the form is submitted.
        form.onsubmit = function(e) {
            e.preventDefault();

            var message = messageField.value;
            websocketService.sendMessage(message);

            messagesList.innerHTML += '<li class="sent"><span>Sent:</span>' + message +
                '</li>';

            messageField.value = '';

            return false;
        };


        // Close the WebSocket connection when the close button is clicked.
        closeBtn.onclick = function(e) {
            e.preventDefault();

            websocketService.closeSocket();

            return false;
        };

    }

})(angular);
