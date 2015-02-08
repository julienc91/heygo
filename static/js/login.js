app.controller('login_controller', ['$scope', '$http',
    function($scope, $http) {
        $scope.user = {};
        $scope.invitation = {};

        $scope.login = function() {
            $http.get("/login", {params: {"user": JSON.stringify($scope.user)}}).success(function(response) {
                if (response.ok) {
                    window.location.href = "/about";
                }
            });
        };

        $scope.signup = function() {
            $http.get("/signup", {params: {"user": JSON.stringify($scope.user)}}).success(function(response) {
                if (response.ok) {
                    window.location.href = "/about";
                }
            });
        };

    }
]);
