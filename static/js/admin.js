$(function() {
    
    config = {};
    config["display_fields"] = {"users": "login", "groups": "title", "videos": "title", "invitations": "value", "video_groups": "title"};
    config["ref_main_tab_filter"] = {"membership": "table[data-table=groups] td[data-field=title]", "classification": "table[data-table=video_groups] td[data-field=title]", "permissions": "table[data-table=groups] td[data-field=title]"};
    config["iter_main_tab_filter"] = {"membership": "table[data-table=users] td[data-field=login]", "classification": "table[data-table=videos] td[data-field=title]", "permissions": "table[data-table=video_groups] td[data-field=title]"};
    config["keyname"] = {"membership": "groups_id", "classification": "video_groups_id", "permissions": "groups_id" };
    config["valuename"] = {"membership": "users_id", "classification": "videos_id", "permissions": "video_groups_id" };


    // Generate random passwords
    function random_password() {
        return Math.random().toString(36).substring(2, 15);
    }
    
    // Initialize tooltips
    function init_tooltips() {
        $(".save_row").children().attr("data-toggle", "tooltip").attr("title", "[[ 'SAVE' | translate ]]");
        $(".reset_row").children().attr("data-toggle", "tooltip").attr("title", "[[ 'CANCEL' | translate ]]");
        $(".del_row").children().attr("data-toggle", "tooltip").attr("title", "[[ 'DELETE' | translate ]]");
        $(".add_row").children().attr("data-toggle", "tooltip").attr("title", "[[ 'ADD' | translate ]]");
        $("[data-toggle=tooltip]").tooltip({
            placement : 'top'
        });
    }
    
    // Hide alerts instead of dismissing them
    $("[data-hide]").on("click", function(){
        $(this).closest("." + $(this).attr("data-hide")).hide();
    });
    
    // Hide alert boxes and "reset password" inputs
    $("div[role=alert], input[data-table=users][data-field=password]").hide();
    
    // Disable "save" and "reset" buttons
    $(".save_row, .reset_row").attr("disabled", "disabled");
    
    // Click on "reset password": generate random password and activate the save button
    $(document.body).on("click", "[role=reset_password]", function() {
        $(this).parent().attr("role", "data").attr("contenteditable", "true");
        $(this).parent().trigger("input");
        $(this).parent().text(random_password());
    });
    
    // Save and Delete row buttons for new rows
    function new_row_buttons(table) {
        return "    <td>" +
            "        <button type=\"submit\" class=\"save_row btn-link\" data-table=\"" + table + "\" data-id=\"0\">" +
            "            <span class=\"glyphicon glyphicon-ok-sign\"></span>" +
            "        </button>" +
            "        <button type=\"submit\" class=\"reset_row btn-link\" data-table=\"" + table + "\" data-id=\"0\" disabled=\"disabled\">" +
            "            <span class=\"glyphicon glyphicon-ban-circle\"></span>" +
            "        </button>" +
            "        <button type=\"submit\" class=\"del_row btn-link\" data-table=\"" + table + "\" data-id=\"0\">" +
            "            <span class=\"glyphicon glyphicon-minus-sign\"></span>" +
            "        </button>" +
            "    </td>";
    }
    
    // Add a new row
    function new_row_cell(table, field, editable, role, content) {
        var res = "<td role=\"" + role + "\" data-id=\"0\" data-table=\"" + table + "\"";
        if (field != "")
            res += " data-field=\"" + field + "\"";
        if (editable)
            res += " data-original-value=\"\" contenteditable=\"true\"";
        res += ">" + content + "</td>";
        return res
    }
    
    // Add a new row for the users table
    $("button.add_row[data-table=users]").on("click", function() {
        $("table[data-table=users]").children("tbody").append(
            "<tr>" +
                new_row_cell("users", "", false, "id", "") +
                new_row_cell("users", "login", true, "data", "Login") +
                new_row_cell("users", "password", true, "data", random_password()) +
                new_row_buttons("users") +
            "</tr>");
        $(this).attr("disabled", "disabled");
        init_tooltips();
    });
    
    // Add a new row for the groups table
    $("button.add_row[data-table=groups],button.add_row[data-table=video_groups]").on("click", function() {
        var group = $(this).attr("data-table");
        $("table[data-table=" + group + "]").children("tbody").append(
            "<tr>" +
                new_row_cell(group, "", false, "id", "") +
                new_row_cell(group, "title", true, "data", "Nom") +
                new_row_buttons(group) +
            "</tr>");
        $(this).attr("disabled", "disabled");
        init_tooltips();
    });
    
    // Add a new row for the invitations table
    $("button.add_row[data-table=invitations]").on("click", function() {
        $("table[data-table=invitations]").children("tbody").append(
            "<tr>" +
                new_row_cell("invitations", "", false, "id", "") +
                new_row_cell("invitations", "value", true, "data", random_password()) +
                new_row_buttons("invitations") +
            "</tr>");
        $(this).attr("disabled", "disabled");
        init_tooltips();
    });
    
    // Add a new row for the videos table
    $("button.add_row[data-table=videos]").on("click", function() {
        $("table[data-table=videos]").children("tbody").append(
            "<tr>" +
                new_row_cell("videos", "", false, "id", "") +
                new_row_cell("videos", "title", true, "data", "Titre") +
                new_row_cell("videos", "path", true, "data", "Chemin") +
                new_row_cell("videos", "slug", true, "data", "Slug") +
                new_row_cell("videos", "url", false, "", "") +
                new_row_buttons("videos") +
            "</tr>");
        $(this).attr("disabled", "disabled");
        init_tooltips();
    });
    
    // Click on "delete row"
    $(document.body).on("click", "button.del_row", function() {
        
        var table = $(this).attr("data-table");
        var id = $(this).attr("data-id");
        var alert = $("#" + table + "_alert");
        
        // If the click is to delete an existing entry
        if (id != 0) {
            var elem = $(this);
            
            // Delete the entry from the database
            $.ajax({
                url: 'admin/delete',
                data: {"table": table, "id": id},
                    success: function(data) {
                        data = JSON.parse(data);
                        var success = data["ok"];
                        if(success) {
                            elem.closest("tr").remove();
                            elem.removeClass("info danger");
                            data["err"] = "Modifications enregistrées";
                            refresh_permissions();
                        } else {
                            elem.addClass("danger");
                        }
                        
                        print_alert(success, data["err"], alert);
                    }
            });
        // Else, delete the row from the table and show the "add row" button
        } else {
            $("button.add_row[data-table=" + table + "]").removeAttr("disabled");
            $(this).closest("tr").remove();
            refresh_permissions();
        }
    });
    
    // Disable the "enter" keypress in editable tags
    $(document.body).on("keypress", "[role=data]", function(e) {
        if (e.which == 13)
            return false;
    });
    
    // Enable the "save" button when editing an editable tags
    $(document.body).on("input", "[role=data]", function() {
        
        $(".save_row[data-table=" + $(this).attr("data-table") + "][data-id=" + $(this).attr("data-id") + "]").removeAttr("disabled");
        if ($(this).attr("data-id") != "0") {
            $(".reset_row[data-table=" + $(this).attr("data-table") + "][data-id=" + $(this).attr("data-id") + "]").removeAttr("disabled");
        }
    });
    
    function reset_password_button(id) {
        $("[role=data][data-table=users][data-id=" + id + "][data-field=password]").removeAttr("contenteditable role").html(
                        "<button class=\"btn btn-warning\" role=\"reset_password\">Réinitialiser</button>");
    }
    
    
    // Click on "reset row"
    $(document.body).on("click", ".reset_row", function() {
        $("[role=data][data-table=" + $(this).attr("data-table") + "][data-id=" + $(this).attr("data-id") + "][data-original-value]").each( function () {
            $(this).text($(this).attr("data-original-value"));
        });
        $("button.save_row[data-table=" + $(this).attr("data-table") + "][data-id=" + $(this).attr("data-id") + "]").attr("disabled", "disabled");
        $("button.reset_row[data-table=" + $(this).attr("data-table") + "][data-id=" + $(this).attr("data-id") + "]").attr("disabled", "disabled");
        
        if ($(this).attr("data-table") == "users" && $(this).attr("data-id") != "0") {
            reset_password_button($(this).attr("data-id"));
        }
    });
    
    // Click on "save row"
    $(document.body).on("click", ".save_row", function() {
            
        var id = $(this).attr("data-id");
        var table = $(this).attr("data-table");
        var alert = $("#" + table + "_alert");
        var elem = $(this);
        var is_insert = (id == 0);
        var dest = "";
        
        var parameters = {"table": table};
        $("[role=data][data-table=" + table + "][data-id=" + id + "]").each( function() {
            parameters[$(this).attr("data-field")] = $(this).text();
        });
        
        if (is_insert) {
            dest = "admin/insert";
        } else {
            parameters["id"] = id;
            dest = "admin/update";
        }
        
        $.ajax({
            url: dest,
            data: parameters,
            success: function(data) {
                
                data = JSON.parse(data);
                var success = data["ok"];
                
                if(success) {
                    data["err"] = "Modifications enregistrées";
                    if (is_insert) {
                        id = data["id"];
                        $("[data-table=" + table + "][data-id=0][role=id]").removeAttr("role data-id data-table").text(id);
                        $("[data-table=" + table + "][data-id=0]").attr("data-id", id);
                        $(".add_row[data-table=" + table + "]").removeAttr("disabled");
                        
                        $("[data-table=" + table + "][data-id=" + id + "][data-original-value]").each( function() {
                            $(this).attr("data-original-value", $(this).text());
                        });
                    }
                    
                    $(".save_row[data-table=" + table + "][data-id=" + id + "]").attr("disabled", "disabled");
                    $(".reset_row[data-table=" + table + "][data-id=" + id + "]").attr("disabled", "disabled");
                    
                    if (table == "users")
                        reset_password_button(id);
                    else if (table == "videos") {
                        var link = "/videos/watch/" + $("[data-table=videos][data-id=" + id + "][data-field=slug]").text();
                        $("[data-table=videos][data-id=" + id + "][data-field=url]").html("<a href=\"" + link + "\">" + link + "</a>");
                    }
                    
                    refresh_permissions();
                }
                
                print_alert(success, data["err"], alert);
            }
        });
    });
    
    // Refresh checkboxes
    function refresh_permissions() {
        $("#membership_ref, #classification_ref, #permissions_ref").each( function() {
            update_permissions_values($(this).attr("data-tab-id"));
        });
    }
    
    // Display a message in the given alert box
    function print_alert(success, text, alert) {
        
        if (!success) {
            alert.removeClass("alert-success").addClass("alert-danger");
        } else {
            alert.removeClass("alert-danger").addClass("alert-success");
        }
        alert.children(".alert_content").text(text);
        alert.show();
    }
    
    // Update permission views when selecting a new group
    function update_permissions(tab_id) {
        
        var data = json_values[tab_id];
        var ref_id = $("#" + tab_id + "_ref").val();
        $("#" + tab_id + "_iter").find("[type=checkbox]").prop("checked", false);
        if (ref_id in data) {
            for (var i in data[ref_id]) {
                var iter_id = data[ref_id][i];
                $("#" + tab_id + "_iter").find("[type=checkbox][value=" + iter_id + "]").prop("checked", true);
            }
        }
    }
    
    // Change values of select inputs and checkboxes
    function update_permissions_values(tab_id) {
        
        var ref_main_tab_filter = config["ref_main_tab_filter"][tab_id];
        var iter_main_tab_filter = config["iter_main_tab_filter"][tab_id];
        $("#" + tab_id + "_ref").html("");
        $("#" + tab_id + "_iter").html("");
        
        $(ref_main_tab_filter).each( function() {
            
            var ref_id = $(this).attr("data-id");
            if (ref_id == 0)
                return
            var text = $(this).text();
            
            $("#" + tab_id + "_ref").append("<option value=\"" + ref_id + "\">" + text + "</option>");
        });
        $(iter_main_tab_filter).each( function() {
            
            var iter_id = $(this).attr("data-id");
            if (iter_id == 0)
                return
            var text = $(this).text();
            
            $("#" + tab_id + "_iter").append("<div class=\"checkbox\"><label><input type=\"checkbox\" data-tab-id=\"" + tab_id + "\" value=\"" + iter_id + "\">" + text + "</label></div>");
        });
        
        update_permissions(tab_id);
    }
    
    $("[role=select_ref]").change( function() { update_permissions($(this).attr("data-tab-id")); });
    
    
    // Change checkboxes state
    $(document.body).on("change", "[type=checkbox][data-tab-id]", function() {
        
        var table = $(this).attr("data-tab-id");
        var alert = $("#permissions_alert");
        var dest = "";
        var is_insert = $(this).prop('checked');
        var ref = $("#" + table + "_ref").val();
        var iter = $(this).val();
        
        var parameters = {"table": table};
        parameters[config["keyname"][table]] = ref;
        parameters[config["valuename"][table]] = iter;
        
        if (is_insert) {
            dest = "admin/insert";
        } else {
            dest = "admin/delete_pivot";
        }
        
        $.ajax({
            url: dest,
            data: parameters,
            success: function(data) {
                
                data = JSON.parse(data);
                var success = data["ok"];
                
                if(success) {
                    data["err"] = "Modifications enregistrées";
                    if (is_insert)
                        json_values[table][ref].push(iter);
                    else
                        json_values[table][ref].splice(json_values[table][ref].indexOf(iter));
                }
                
                print_alert(success, data["err"], alert);
            }
        });
    });
    
    
    refresh_permissions();
    init_tooltips();
    
});
