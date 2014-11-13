var config = {
    "titles": {
        "users": "Utilisateurs",
        "invitations": "Invitations",
        "groups": "Groupes",
        "videos": "Vidéos",
        "video_groups": "Groupes de vidéos"
    },
    "fields": {
        "users": ["id", "login", "password"],
        "invitations": ["id", "value"],
        "groups": ["id", "title"],
        "videos": ["id", "title", "path", "slug"],
        "video_groups": ["id", "title"]
    },
    "column_definitions": {
        "users": [{field: "id", width: "10%", displayName: "ID"}, {field: "login", width: "80%", displayName: "Login"}],
        "invitations": [{field: "id", width: "10%", displayName: "ID"}, {field: "value", width: "80%", displayName: "Valeur"}],
        "groups": [{field: "id", width: "10%", displayName: "ID"}, {field: "title", width: "80%", displayName: "Nom"}],
        "videos": [{field: "id", width: "10%", displayName: "ID"}, {field: "title", width: "25%", displayName: "Titre"}, {field: "path", width: "35%", displayName: "Chemin"}, {field: "slug", width: "20%", displayName: "Slug"}],
        "video_groups": [{field: "id", width: "10%", displayName: "ID"}, {field: "title", width: "80%", displayName: "Nom"}]
    },
    "main_fields": {
        "users": "login",
        "invitations": "value",
        "groups": "title",
        "videos": "title",
        "video_groups": "title"
    },
    "default_values": {
        "users": [0, "", "random_string"],
        "invitations": [0, "random_string"],
        "groups": [0, ""],
        "videos": [0, "", "", ""],
        "video_groups": [0, ""]
    }
};


app.config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/users");
    
    $stateProvider
        .state('generic_tables', {
            url: "/{table:users|invitations|groups|videos|video_groups}",
            controller: "generic_grid_view_controller",
            templateUrl: "/static/html/admin/generic_grid_view.html"})
        .state('users_new', {
            url: "/{table:users}/new/",
            controller: "generic_edit_view_controller",
            templateUrl: "/static/html/admin/users_add_view.html"})
        .state('users_edit', {
            url: "/{table:users}/{mode:edit}/{id:[0-9]+}",
            controller: "generic_edit_view_controller",
            templateUrl: "/static/html/admin/users_edit_view.html"})
        .state('invitations_edit', {
            url: "/{table:invitations}/{mode:new|edit}/{id:[0-9]*}",
            controller: "generic_edit_view_controller",
            templateUrl: "/static/html/admin/invitations_edit_view.html"})
        .state('generic_groups_edit', {
            url: "/{table:groups|video_groups}/{mode:new|edit}/{id:[0-9]*}",
            controller: "generic_edit_view_controller",
            templateUrl: "/static/html/admin/generic_groups_edit_view.html"})
        .state('videos_edit', {
            url: "/{table:videos}/{mode:new|edit}/{id:[0-9]*}",
            controller: "generic_edit_view_controller",
            templateUrl: "/static/html/admin/videos_edit_view.html"});
}]);


app.controller('generic_grid_view_controller', ['$scope', '$http', '$stateParams',
    function ($scope, $http, $stateParams) {
        $scope.table = $stateParams.table;
        $scope.title = config["titles"][$scope.table];
        $scope.rows = [];
        
        $scope.external_scope = {
            table: $scope.table,
            delete_row: function(row) {
                bootbox.confirm("Cette action supprimera définitivement l'entrée \"" + row[config["main_fields"][$scope.table]] + "\". Voulez-vous continuer ?",
                    function(result) {
                        if (result) {
                            $http.get("/admin/delete/" + $scope.table + "/" + row.id).success(function(response) {
                                if (response["ok"]) {
                                    var index_to_delete = -1;
                                    for (var i=0; i<$scope.rows.length; i++) {
                                        if ($scope.rows[i].id == row.id) {
                                            index_to_delete = i;
                                            break;
                                        }
                                    }
                                    if (index_to_delete >= 0)
                                        $scope.rows.splice(index_to_delete, 1);
                                }
                        }).error(display_error_message);
                    }
                });
            }
        };
        
        $scope.grid_options = { enableSorting: true,
                               enableColumnMenus: false,
                               enableColumnResizing: false,
                               columnDefs: config["column_definitions"][$scope.table],
                               data: 'rows' };
        $scope.grid_options.columnDefs.push({name: 'actions',
                                            displayName: 'Actions',
                                            width: "15%",
                                            cellTemplate: '<a class="edit btn-link" ng-href="#/{{getExternalScopes().table}}/edit/{{row.entity.id}}"><span class="glyphicon glyphicon-pencil"></span></a>\
                                                           <button class="delete btn-link" ng-click="getExternalScopes().delete_row(row.entity)"><span class="glyphicon glyphicon-remove"></span></button>'});
        $http.get("/admin/get/" + $scope.table).success(function(response) {
            if (response["ok"])
                $scope.rows = response["data"];
        }).error(display_error_message);
    }
]);


app.controller('generic_edit_view_controller', ['$scope', '$http', '$stateParams', '$location',
    function ($scope, $http, $stateParams, $location) {
        $scope.table = $stateParams.table;
        $scope.fields = config["fields"][$scope.table];
        $scope.valid_path = false;
        $scope.is_new = $stateParams.mode == "new";

        if (!$scope.is_new && (!$stateParams.id || $stateParams.id == 0))
            $location.path("/" + $scope.table + "/new/");
        
        $scope.random_string = function() {
            var random_string = "";
            for (; random_string.length < 14; random_string = Math.random().toString(36).slice(2)) {}
            return random_string;
        };
        
        $scope.model = {};
        if ($stateParams.id) {
            $http.get("/admin/get/" + $scope.table + "/" + $stateParams.id).success(function(response) {
                if (response["ok"])
                    $scope.model = response["data"];
            }).error(display_error_message);
        }
        
        $scope.init = function() {
            if (!$scope.is_new)
                return
            for (var i=0; i<$scope.fields.length; i++) {
                var default_value = config["default_values"][$scope.table][i];
                if ($scope[default_value])
                    default_value = $scope[default_value]();
                $scope.model[$scope.fields[i]] = default_value;
            }
        };

        $scope.check_file_on_server = function(model_attribute) {
            $http.get("/admin/media/check", {params: {"path": $scope.model[model_attribute]}}).success(function(response) {
                $scope.valid_path = response["ok"];
            }).error(display_error_message);
        };
        
        $scope.change_random_value = function(model_attribute) {
            $scope.model[model_attribute] = $scope.random_string();
        };

        $scope.generate_slug = function(model_attribute) {
            if (!$scope.model[model_attribute]) {
                $scope.model["slug"] = "";
                return;
            }
            var slug = $scope.model[model_attribute].toLowerCase();
            var from = "ãàáäâ@ẽèéëêìíïîõòóöôùúüûñç";
            var to   = "aaaaaaeeeeeiiiiooooouuuunc";
            for (var i=0, l=from.length ; i<l ; i++) {
                slug = slug.replace(new RegExp(from.charAt(i), 'g'), to.charAt(i));
            }
            $scope.model["slug"] = slug.replace(/[^\w]+/g, "_");
        };
        
        $scope.save_model = function() {
            if ($scope.edit.$valid) {
                if (!$scope.is_new && $scope.table == "users" && $scope.model.new_password)
                    $scope.model.password = $scope.model.new_password;
                var url = "/admin/" + ($scope.is_new ? "insert" : "update") + "/" + $scope.table + ($scope.model.is_new ? "" : "/" + $scope.model.id);
                $http.get(url, {params: $scope.model}).success(function(response) {
                    if (response["ok"]) {
                        $scope.model = response["data"];
                        $location.path("/" + $scope.table);
                    }
                }).error(display_error_message);
            }
        };

    }
]);



// Hide alerts instead of dismissing them
$("[data-hide]").on("click", function(){
    $(this).closest("." + $(this).attr("data-hide")).hide();
});

// Hide alert boxes
$("#alert_box").hide();

var set_active_tab = (function() {
                        var hash = $(location).attr('hash');
                        var tab_id = hash.split("/");
                        if (tab_id.length > 1 && tab_id[0] == "#" && $("#admin_tabs a[href=#" + tab_id[1] + "][data-toggle=pill]").length == 1)
                            $("#admin_tabs a[href=#" + tab_id[1] + "][data-toggle=pill]").parent().addClass("active");
                        })();


function display_error_message(data, status, headers, config) {
    $("#alert_box").children(".alert_content").text(data["err"] ? data["err"] : status);
    $("#alert_box").show();
}


(function() {

    $(".save_row").children().attr("data-toggle", "tooltip").attr("title", "Enregistrer");
    $(".reset_row").children().attr("data-toggle", "tooltip").attr("title", "Annuler");
    $(".del_row").children().attr("data-toggle", "tooltip").attr("title", "Supprimer");
    $(".add_row").children().attr("data-toggle", "tooltip").attr("title", "Ajouter");
    $("[data-toggle=tooltip]").tooltip({
        placement : 'top'
    });
    
    $("[data-hide]").on("click", function(){
        $(this).closest("." + $(this).attr("data-hide")).hide();
    });
})();
