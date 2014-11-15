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
            templateUrl: "/static/html/videos/detail_view.html"})
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
            if (response["ok"])
                $scope.rows = response["data"];
        }).error(display_error_message);
    }
]);


app.controller('videos_detail_view_controller', ['$scope', '$http', '$stateParams',
    function ($scope, $http, $stateParams) {
        $scope.model = {};
        $http.get("/videos/get/" + $stateParams.slug).success(function(response) {
            if (response["ok"])
                $scope.model = response["data"];
        }).error(display_error_message);

        $scope.media_link = function(slug) {
            if (slug)
                return "/media/videos/" + slug;
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
    $("#alert_box").children(".alert_content").text(data["err"] ? data["err"] : status);
    $("#alert_box").show();
}


$(function() {

    var slug = $("video").attr("data-slug");

    function resize_video() {
        $("video").width($(".video_container").width());
    }

    $(window).resize(resize_video);
    resize_video();

    var api_username = "";
    var api_password = "";
    var api_useragent = "OS Test User Agent";

    $("[role=search_subtitles]").click(function() {
        var data = {}

        get_hash(data);
        api_login(data);
        api_search_subtitles(data.hash, data.size, data.token);
        api_logout(data.token);
    });

    function api_login(d) {
        $.xmlrpc({

            url: 'http://api.opensubtitles.org/xml-rpc',
            methodName: 'LogIn',
            params: [api_username, api_password, "", api_useragent],
            async: false,

            success: function(response, status, jqXHR) {
                if (response.length > 0 && response[0].status == "200 OK") {
                    d.token = response[0].token;
                }
            }
        });
    }

    function api_logout(token) {
        $.xmlrpc({
            url: 'http://api.opensubtitles.org/xml-rpc',
            methodName: 'LogOut',
            params: [token]
        });
    }

    function api_search_subtitles(hash, size, token) {
        $.xmlrpc({

            url: 'http://api.opensubtitles.org/xml-rpc',
            methodName: 'SearchSubtitles',
            params: [token, [{"moviehash": hash, "moviebytesize": size}]],
            async: false,

            success: function(response, status, jqXHR) {
                console.log(response);
            }
        });
    }

    function get_hash(d) {
        $.ajax({

            url: "/videos/gethash/" + slug,
            async: false,

            success: function(data) {
                data = JSON.parse(data);
                if (data.ok) {
                    d.hash = data.hash;
                    d.size = data.size;
                }
            }
        });
    }
});
