app.config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");
    $stateProvider
        .state('grid_view', {
            url: "/",
            controller: "videos_grid_view_controller",
            templateUrl: "/static/html/videos/grid_view.html"})
        .state('detail_view', {
            url: "/{slug:[a-z0-9_]+}",
            controller: "videos_detail_view_controller",
            templateUrl: "/static/html/videos/detail_view.html"});
}]);

app.controller('videos_grid_view_controller', ['$scope', '$http',
    function ($scope, $http) {
        $scope.rows = [];
        $scope.grid_options = { enableSorting: true,
                               enableColumnMenus: false,
                               enableColumnResizing: false,
                               columnDefs: [{field: "id", width: "10%", displayName: "ID"}, {field: "title", width: "25%", displayName: "Titre"}],
                               data: 'rows' };
        $scope.grid_options.columnDefs.push({name: 'link',
                                            displayName: 'Lien',
                                            width: "15%",
                                            cellTemplate: '<a class="btn-link" ng-href="#/{{row.entity.slug}}"><span class="glyphicon glyphicon-film"></span></a>'});

        $http.get("/videos/get").success(function(response) {
            if (response.ok)
                $scope.rows = response.data;
        }).error(display_error_message);
    }
]);


app.controller('videos_detail_view_controller', ['$scope', '$http', '$stateParams', '$sce',
    function ($scope, $http, $stateParams, $sce) {
        $scope.model = {};
        $scope.video_config = {theme: {url: "/static/css/videogular.css"}};
		$scope.API = null;

		$scope.on_player_ready = function(API) {
			$scope.API = API;
		};

        $http.get("/videos/get/" + $stateParams.slug).success(function(response) {
            if (response.ok) {
                $scope.model = response.data;
                $scope.video_config.sources = [{src: "/media/videos/" + $scope.model.slug, type: "video/" + $scope.model.path.split('.').pop()}];
            }
        }).error(display_error_message);

        $scope.search_subtitles = function() {
            $scope.video_config.tracks = [{src: "/videos/getsubtitles/" + $scope.model.slug, kind: "subtitles", srclang: "fr", label: "French", default: "true"}];
        };
    }
]);


// Hide alerts instead of dismissing them
$("[data-hide]").on("click", function(){
    $(this).closest("." + $(this).attr("data-hide")).hide();
});

// Hide alert boxes
$("#alert_box").hide();

function display_error_message(data, status, headers, config) {
    $("#alert_box").children(".alert_content").text(data.err ? data.err : status);
    $("#alert_box").show();
}
