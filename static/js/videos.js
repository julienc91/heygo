app.config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/grid");
    $stateProvider
        .state('grid_view', {
            url: "/grid",
            controller: "videos_grid_view_controller",
            templateUrl: "/static/html/videos/grid_view.html"})
        .state('thumbnail_view', {
            url: "/list",
            controller: "videos_thumbnail_view_controller",
            templateUrl: "/static/html/videos/thumbnail_view.html"})
        .state('detail_view', {
            url: "/watch/{slug:[a-z0-9_]+}",
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
                                            cellTemplate: '<a class="btn-link" ng-href="#/watch/{{row.entity.slug}}"><span class="glyphicon glyphicon-film"></span></a>'});

        $http.get("/videos/get").success(function(response) {
            if (response.ok)
                $scope.rows = response.data;
        }).error(display_error_message);
    }
]);

app.controller('videos_thumbnail_view_controller', ['$scope', '$http',
    function ($scope, $http) {
        $scope.rows = [];
        $scope.rows_to_display = [];
        $scope.current_filter = "";

        $http.get("/videos/get").success(function(response) {
            if (response.ok) {
                $scope.rows = response.data;
                $scope.rows_to_display = $scope.rows;
                $scope.load_thumbnails();
            }
        }).error(display_error_message);

        $scope.load_thumbnails = function() {
            angular.forEach($scope.rows_to_display, function(row) {
                if (row.imdb_id && !row.loaded) {
                    row.loaded = true;
                    $http.get("http://www.omdbapi.com/?i=" + row.imdb_id).success(function(response) {
                        if (response.Response != "False") {
                            if (response.Poster && response.Poster == "N/A")
                                row.thumbnail = "/static/images/no_thumbnail.png";
                            else
                                row.thumbnail = "/media/thumbnail/" + btoa(response.Poster);
                            if (response.Plot != "N/A")
                                row.resume = response.Plot;
                            row.year = response.Year;
                        }
                    }).error(display_error_message);
                } else if (!row.imdb_id && !row.loaded) {
                    row.loaded = true;
                    row.thumbnail = "/static/images/no_thumbnail.png";
                }
            });
        };

        $scope.filter_results = function() {
            var tmp_rows_to_display = [];
            var filter = $scope.current_filter.toLowerCase();

            angular.forEach($scope.rows, function(row) {
                if (row.title.toLowerCase().indexOf(filter) >= 0) {
                    tmp_rows_to_display.push(row);
                }
            });
            $scope.rows_to_display = tmp_rows_to_display;
            $scope.load_thumbnails();
        };
    }
]);


app.controller('videos_detail_view_controller', ['$scope', '$http', '$stateParams', '$sce',
    function ($scope, $http, $stateParams, $sce) {
        $scope.model = {};
        $scope.video_config = {theme: {url: "/static/css/videogular.css"}};
		$scope.API = null;
        $scope.tracks = [];
        $scope.current_track = -1;

		$scope.on_player_ready = function(API) {
			$scope.API = API;
		};

        $http.get("/videos/get/" + $stateParams.slug).success(function(response) {
            if (response.ok) {
                $scope.model = response.data;
                $scope.video_config.sources = [{src: "/media/videos/" + $scope.model.slug, type: "video/" + $scope.model.path.split('.').pop()}];
            }
        }).error(display_error_message);

        $scope.search_subtitles = function(lang) {
            $('#search_subtitles').hide();
            if ($scope.tracks.length > 0)
                return;
            var label = (lang == 'eng') ? "Anglais" : "Français";
            var srclang = (lang == 'eng') ? "en" : "fr";
            $http.get("/media/subtitles/list/" + $stateParams.slug + "/" + lang).success(function(response) {
                if (response.ok) {
                    for (var i in response.data) {
                        $scope.tracks.push({src: "/media/subtitles/get/" + response.data[i], kind: "subtitles", srclang: srclang, label: label + " " + (parseInt(i) + 1), trackid: i, default: true});
                    }
                    $scope.change_subtitles(0);
                }
            }).error(display_error_message);
        };

        $scope.change_subtitles = function(i) {
            if (i < 0 || i >= $scope.tracks.length) {
                $scope.video_config.tracks = [];
                $scope.current_track = -1;
            }
            else {
                $scope.video_config.tracks = [$scope.tracks[i]];
                $scope.current_track = i;
            }
        };
    }
]);

app.directive("vgSubtitles",
    function() {
        return {
            restrict: "E",
            templateUrl: "/static/html/videos/subtitles.html"
        };
    }
);


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
